package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type Ingredient struct {
	Ing_id           int
	Ing_name         string
	Quantity         float32
	QuantityType     string
	Price            float32
	PriceForQuantity float32
	Tag              string
	Remains          float32
}

type Receipt struct {
	Rec_id           int
	Name             string
	Quantity         float32
	QuantityType     string
	Price            float32
	PriceForQunatity float32
	Tag              string
}
