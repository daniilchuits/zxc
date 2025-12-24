package stack

import (
	"context"
	"encoding/json"
	"fmt"
)

type StackItem struct {
	Value interface{}
	Path  string
}

func PayloadToMap(data any, ctx context.Context) (map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("Context done")
	default:
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var kefteme map[string]any
	err = json.Unmarshal(b, &kefteme)
	if err != nil {
		return nil, err
	}
	return kefteme, nil
}

func Stack(info map[string]any, ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	stack := []StackItem{
		{Value: info, Path: ""},
	}

	for len(stack) > 0 {
		n := len(stack) - 1 // индекс последнего элемента - длина стэка-1 (длина стека всегда будет либо 0, либо 1), выходит: 1 Value:detailsMap, path: details (потом этот элемент удаляется и разбор стека идет от корня к веткам, пока стек не будет не map[string]any, тогда к стеку не добавляется ничего и его len(stack) = 0)

		item := stack[n]  // последни элемент
		stack = stack[:n] // уменьшаем стэк с каждой иттерацией

		switch val := item.Value.(type) {

		case map[string]interface{}: // если мы не дошли до конца вложености, делаем это
			for k, v := range val {
				path := k
				if item.Path != "" {
					path = item.Path + "." + path // добавляем .path, к path'у до этого
				}
				stack = append(stack, StackItem{
					Value: v,    // значение мы выбрали из карты
					Path:  path, // путь добавляем каждый раз
				})
			}
		default: // когда путь к нашем info заканчивается, остается не мапа
			fmt.Printf("%s = %v\n", item.Path, val) // выбираем значение у последнего элемента вложености
		}
	}
	fmt.Println()
}
