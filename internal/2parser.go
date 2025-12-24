package internal

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"path/filepath"
	"projct/model"

	"gopkg.in/yaml.v3"
)

func ParseJSON(infoData model.FileInfo) (model.Event, error) {
	var j model.Event
	err := json.Unmarshal(infoData.Data, &j)
	if err != nil {
		return model.Event{}, fmt.Errorf("parse json - %s error: %w", infoData.Path, err)
	}
	if j.Payload.Currency == "" {
		return model.Event{}, fmt.Errorf("Empty currency in %s", infoData.Path)
	}
	if j.Payload.Amount <= 0 {
		return model.Event{}, fmt.Errorf("Ammount above zero in %s", infoData.Path)
	}
	return j, nil
}

func ParseXML(infoData model.FileInfo) (model.Event, error) {
	var x model.Event
	err := xml.Unmarshal(infoData.Data, &x)
	if err != nil {
		return model.Event{}, fmt.Errorf("parse xml - %s error: %w", infoData.Path, err)
	}
	if x.Payload.Currency == "" {
		return model.Event{}, fmt.Errorf("Empty currency in %s", infoData.Path)
	}
	if x.Payload.Amount <= 0 {
		return model.Event{}, fmt.Errorf("Ammount above zero in %s", infoData.Path)
	}
	return x, nil
}

func ParseYAML(infoData model.FileInfo) (model.Event, error) {
	var y model.Event
	err := yaml.Unmarshal(infoData.Data, &y)
	if err != nil {
		return model.Event{}, fmt.Errorf("parse yaml - %s error: %w", infoData.Path, err)
	}
	if y.Payload.Currency == "" {
		return model.Event{}, fmt.Errorf("Empty currency in %s", infoData.Path)
	}
	if y.Payload.Amount <= 0 {
		return model.Event{}, fmt.Errorf("Ammount above zero in %s", infoData.Path)
	}
	return y, nil
}

func Parse(info model.FileInfo, ctx context.Context) (model.Event, error) {
	path := filepath.Ext(info.Path)

	select {
	case <-ctx.Done():
		return model.Event{}, nil
	default:
	}

	switch path {
	case ".json":
		return ParseJSON(info)
	case ".xml":
		return ParseXML(info)
	case ".yaml":
		return ParseYAML(info)
	default:
		return model.Event{}, fmt.Errorf("No available parser for file: %s", info.Path)
	}
}
