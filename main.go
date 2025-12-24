package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"projct/internal"
	"projct/model"
	"projct/postgres"
	"projct/stack"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

const NumWorkers = 2

func main() {
	l, err := os.OpenFile("errs.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.SetOutput(l)

	cnnStr := "user=postgres dbname=files_proektik password=LimitedEdition228 sslmode=disable"
	db, err := sql.Open("postgres", cnnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = postgres.PostgresTable(db)
	if err != nil {
		log.Fatal(err)
	}

	var mu sync.Mutex

	var wgParsers sync.WaitGroup
	var wgWorkers sync.WaitGroup
	var wgProcessor sync.WaitGroup
	var HTTP sync.WaitGroup
	var wgProcessorLoop sync.WaitGroup

	infos := make(chan model.FileInfo)
	dataForWorkers := make(chan model.Event)
	amm := make(chan model.Amm)
	ammsHTTP := make(chan model.Amm)

	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(2)

	ctxNeErrGroup, cancel := context.WithCancel(context.Background())
	defer cancel()

	filenames, err := os.ReadDir("./input")
	if err != nil {
		log.Fatal("Cant read dir:", err)
		return
	}

	// reader ammsHTTP
	go func() {
		for a := range ammsHTTP {
			am := a
			HTTP.Add(1)
			go func() {
				defer HTTP.Done()
				payloadMapa, err := stack.PayloadToMap(am.Payload, ctxNeErrGroup)
				if err != nil {
					log.Println(err)
					return
				}
				stack.Stack(payloadMapa, ctxNeErrGroup)

				cur := &am
				if cur.Payload.Currency == "" {
					cur.Payload.Currency = "unknown"
				}
				mu.Lock()
				model.Amms = append(model.Amms, *cur) // передается копия cur
				mu.Unlock()
			}()
		}
	}()

	// processor
	// wgProcessor.Add(len(amm)) // так нельзя, потму что в вейтгруппу добавляется единовременная длина канала, а это 0 либо 1 так, как канал небуферизованный, и записи в канал amm идут конкуретно так, что мы не знаем сколько всего будет записей
	wgProcessorLoop.Add(1)
	go func() {
		defer wgProcessorLoop.Done()
		for a := range amm {
			a := a
			wgProcessor.Add(1)

			go func() {
				defer wgProcessor.Done()

				ctxTime, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel2()
				err := internal.Process(a, ammsHTTP, db, ctxTime) // разобраться с указателем на ammsHTTP
				if err != nil {
					log.Println(err)
					return
				}
			}()
		}

		wgProcessor.Wait()
		close(ammsHTTP)
		// fmt.Println(model.Amms) // кароче так как model.Amms добавляется не в одной горутине читать его сразу же после цикла по каналу нельзя, так как нет гарантии, что больше не будут приходить данные, а если данные будут  приходить в закрытый канал - panic
	}()

	// workers
	go func() {
		for i := 1; i <= NumWorkers; i++ {
			wgWorkers.Add(1)

			go func(workerID int) {
				defer wgWorkers.Done()

				for d := range dataForWorkers {
					if err = internal.Worker(d, amm, ctxNeErrGroup); err != nil {
						log.Println(err)
					} else {
						continue
					}
				}
			}(i)

		}
	}()

	//parsers
	go func() {
		for info := range infos {
			wgParsers.Add(1)
			go func(info model.FileInfo) {
				defer wgParsers.Done()

				event, err := internal.Parse(info, ctxNeErrGroup)
				if err != nil {
					log.Println(err)
					return
				} else {
					dataForWorkers <- event
				}
			}(info)
		}
	}()

	// producer
	for _, filename := range filenames {
		g.Go(func() error {
			log.Println("Started openning", filename)
			return internal.Reader(infos, filename, ctx)
		})
	}

	// разделение

	// waiting infos
	if err = g.Wait(); err != nil {
		log.Fatal("Err in reading:", err)
		return
	} else {
		close(infos) // ждем пока все горутины producers доработают, потом закрываем канал, в который они пишут, если нет ошибки
	}

	go func() {
		wgParsers.Wait() // waiting parsers
		close(dataForWorkers)
	}()

	wgWorkers.Wait() // тоже самое с воркерами и со всеми остальными, до последнего
	close(amm)       // последний - должен блокировать main, чтобы основная горутина не закончилась

	wgProcessorLoop.Wait()

	HTTP.Wait()

	amms := model.Amms

	sort.Slice(amms, func(i, j int) bool {
		if amms[i].Timestamp.After(amms[j].Timestamp) {
			return true
		} else if amms[i].Timestamp.Before(amms[j].Timestamp) {
			return false
		}
		return amms[i].Type > amms[j].Type
	})

	for _, am := range model.Amms {
		fmt.Println(am)
		log.Println(am.EventID, "Ended")
	}

	r := mux.NewRouter()

	r.HandleFunc("/events", internal.GETevents).Methods("GET")

	err = http.ListenAndServe(":8081", r) // где-то пропадает 2 лога событий, должно быть 14
	if err != nil {
		log.Fatal("Server couldnt start:", err)
	}

}
