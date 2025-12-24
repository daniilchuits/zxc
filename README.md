короче я кратко описал китайскому интелекту по каким темам желательно было бы сделатть мне задание и он чота высрал, я доделал его вообще позавчера, но вчера за компутероф моожно сказать не сидел, а осталось только перекинуть на github, а уменя просто все настройки слетели и сосать)))))))))))))))

Короче скачиваешь компилятор golang'а, переходишь в директорию гомна и пишешь "go run ." все

Вот задание, я на самом деле вафельно проверял его в целом, я проверял по пунктам, когда делал:

# Проверочный проект: **Log Ingestion & Processing Service**

## 1. Общая идея проекта

Ты реализуешь сервис, который:

1. **Читает файлы логов конкурентно**
    
2. **Парсит данные (JSON / XML / YAML)**
    
3. **Обрабатывает их через worker pool**
    
4. **Отправляет часть данных во внешний HTTP-API**
    
5. **Сохраняет результат в PostgreSQL**
    
6. **Ведёт структурированные логи**
    
7. **Умеет корректно завершаться по context**
    
8. **Сортирует, агрегирует и валидирует данные**
    
9. **Использует стек для обработки вложенных структур**
    
10. **Разбит на пакеты как продакшен-код**
    

---

## 2. Входные данные

### 2.1 Файлы

В директории `./input/` лежат файлы:

`events_001.json events_002.xml events_003.yaml events_004.json ...`

### 2.2 Форматы

#### JSON

`{   "event_id": "e123",   "type": "payment",   "timestamp": "2025-12-01T10:15:30Z",   "payload": {     "amount": 1200,     "currency": "USD"   } }`

#### XML

`<event>   <event_id>e124</event_id>   <type>refund</type>   <timestamp>2025-12-01T11:00:00Z</timestamp>   <payload>     <amount>200</amount>     <currency>USD</currency>   </payload> </event>`

#### YAML

`event_id: e125 type: payment timestamp: 2025-12-01T12:00:00Z payload:   amount: 500   currency: EUR`

---

## 3. Архитектурные требования (ключевая часть)

### 3.1 Пакеты (обязательно)

```powershell
/cmd/app
 /internal/ 
   reader        // чтение файлов   
   parser        // json/xml/yaml   
   worker        // worker pool   
   processor     // бизнес-логика   
   storage       // postgres   
   client        // HTTP клиент   
   logger        // логирование   
   model         // структуры   
   util          // стек, сортировки и т.п.
```

---

## 4. Чтение файлов (io / bufio / concurrency)

### Требования

- Использовать `errgroup.WithContext`
    
- Файлы читать **конкурентно**
    
- Использовать `bufio.Scanner`
    
- Поддержать отмену через `context.Context`
    
- Ошибка в одном файле → отмена всей обработки
    

**Ограничение конкурентности** — не более `N` файлов одновременно.

---

## 5. Парсинг (json / xml / yaml)

### Требования

- Автоматически определять формат по расширению
    
- Отдельный парсер на каждый формат
    
- Общая структура `model.Event`
    
- Все ошибки оборачивать (`fmt.Errorf("parse xml: %w", err)`)
    

---

## 6. Worker Pool

### Требования

- После парсинга события отправляются в `worker pool`
    
- Количество воркеров конфигурируемое
    
- Использовать:
    
    - каналы
        
    - `WaitGroup`
        
    - закрытие каналов **строго корректно**
        
- Каждый воркер:
    
    - валидирует данные
        
    - обрабатывает payload
        
    - отправляет результат дальше
        

---

## 7. Бизнес-логика (processor)

### Логика

- Если `type == "payment"`:
    
    - отправить событие во внешний HTTP-API
        
- Если `type == "refund"`:
    
    - только сохранить в БД
        
- Невалидные события:
    
    - логировать
        
    - не падать
        

---

## 8. HTTP-клиент

### Требования

- `net/http`
    
- Таймауты
    
- Контекст
    
- Повторы (1–2 retries)
    
- Парсинг JSON-ответа
    

---

## 9. PostgreSQL

### Таблица

`CREATE TABLE events (   event_id TEXT PRIMARY KEY,   type TEXT,   timestamp TIMESTAMPTZ,   payload JSONB,   source_file TEXT );`

### Требования

- `database/sql`
    
- Prepared statements
    
- Контекст
    
- Транзакция на батч событий
    
- Обработка конфликтов (`ON CONFLICT DO NOTHING`)
    

---

## 10. Стек (обязательное, не формально)

Использовать **stack** для обработки вложенного payload:

Пример:

`payload: {   "a": {     "b": {       "c": 10     }   } }`

Задача:

- пройти вложенность
    
- собрать ключи вида `a.b.c`
    
- сохранить в лог или структуру
Со стэком было прилично дрочки, по этому я вынесу его отдельно
```go
package stack

import (
	"encoding/json"
	"fmt"
)

type StackItem struct { // структура для нашего стэка, он будет из массива этих                                                      структур
	Value interface{}
	Path  string
}

func PayloadToMap(data any) (map[string]interface{}, error) { // так как наш Stack принимает map[string]interface{}, а передавать мы будем model.Amm, нужно перевезти из model.Amm -> any
	b, err := json.Marshal(data) // получаем его в json формате
	if err != nil {
		return nil, err
	}

	var kefteme map[string]interface{}
	err = json.Unmarshal(b, &kefteme) // получаем из json формата, а оттуда данные всегда и приходят map[string]interface{}
	if err != nil {
		return nil, err
	}
	return kefteme, nil
}

func Stack(info map[string]interface{}) {
	stack := []StackItem{
		{Value: info, Path: ""},
	}

	for len(stack) > 0 {
		n := len(stack) - 1 // индекс последнего элемента - длина стэка-1

		item := stack[n]  // последни элемент
		stack = stack[:n] // уменьшаем стэк с каждой иттерацией

		switch val := item.Value.(type) {

		case map[string]interface{}: // если мы не дошли до конца вложености,                                                          делаем это:
			for k, v := range val {
				path := k
				if item.Path != "" {
					path = item.Path + "." + path // добавляем .path, к path'у до                                                                   этого
				}
				stack = append(stack, StackItem{
					Value: v,    // значение мы выбрали из мапы
					Path:  path, // путь добавляем каждый раз
				})
			}
		default: // когда путь к нашем info заканчивается (остается не мапа)
			fmt.Printf("%s = %v\n", item.Path, val) // выбираем значение у                                           последнего элемента вложености и скопленный путь
		}
	}
}
```

---

## 11. Сортировки

После обработки всех событий:

- Отсортировать:
    
    - по timestamp
        
    - по типу
        
- Использовать `sort.Slice`
    

---

## 12. Логирование

### Требования

- Единый логгер
    
- Уровни:
    
    - INFO
        
    - WARN
        
    - ERROR
        
- Логи:
    
    - старт/завершение
        
    - ошибки файлов
        
    - ошибки парсинга
        
    - отмена по context
        

---

## 13. Context (критично)

Контекст должен:

- передаваться **везде**
    
- отменять:
    
    - чтение файлов
        
    - worker pool
        
    - HTTP-клиент
        
    - работу с БД
