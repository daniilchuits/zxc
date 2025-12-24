package internal

import (
	"context"
	"fmt"
	"projct/model"
	"strings"
)

func Worker(dataForWorkers model.Event, ammount chan<- model.Amm, ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	var BYN float64
	if strings.ToLower(dataForWorkers.Payload.Currency) == "usd" {
		BYN = float64(dataForWorkers.Payload.Amount) * 2.93
		ammount <- model.Amm{
			Ammount: BYN,
			Event:   dataForWorkers,
		}
		return nil
	} else if strings.ToLower(dataForWorkers.Payload.Currency) == "eur" {
		BYN = float64(dataForWorkers.Payload.Amount) * 3.4
		ammount <- model.Amm{
			Ammount: BYN,
			Event:   dataForWorkers,
		}
		return nil
	} else {
		return fmt.Errorf("Not availuble currency: %s", dataForWorkers.EventID)
	}
}
