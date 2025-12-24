package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"projct/model"
)

func GETevents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(model.Amms); err != nil {
		http.Error(w, "Error getting events", http.StatusInternalServerError)
	}
}

func HTTPresponse(info model.Amm, ctx context.Context) error {
	// time.Sleep(3 * time.Second) // проверил, работает ли отмена по контексту
	for i := 1; i <= 2; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("Ctx done. Event %s is not ready", info.EventID)
		default:
		}
		_, err := json.Marshal(info)
		if err != nil {
			continue
		}
	}
	return nil
}

func PostgreWriter(info model.Amm, db *sql.DB) error {
	var (
		method string
		last4  string
		issuer string
	)

	if info.Payload.Details != nil {
		method = info.Payload.Details.Method
		if info.Payload.Details.Card != nil {
			last4 = info.Payload.Details.Card.Last4
			issuer = info.Payload.Details.Card.Issuer
		}
	}

	query := `
	INSERT INTO files (
	event_id, type, event_time, byn, currency,
	originalevent, methods, last4, issuer
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := db.Exec(
		query,
		info.EventID,
		info.Type,
		info.Timestamp,
		info.Ammount,
		info.Payload.Currency,
		info.Payload.OriginalEvent,
		method,
		last4,
		issuer,
	)
	if err != nil {
		return fmt.Errorf("Cant insert into 'files': %s. Err: %v", info.EventID, err)
	}

	return err
}

func Process(info model.Amm, events chan<- model.Amm, db *sql.DB, ctx context.Context) error {
	switch info.Type {
	case "payment":
		err := HTTPresponse(info, ctx)
		if err != nil {
			return err
		}
		events <- info
		return nil
	case "refund":
		err := PostgreWriter(info, db)
		if err != nil {
			return err
		}
		events <- info
		return nil
	default:
		return fmt.Errorf("Unknown event type %s", info.Type)
	}
}
