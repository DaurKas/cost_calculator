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
	Quantity         int
	QuantityType     string
	Price            int
	PriceForQunatity float32
	Tag              string
	Remains          int
}

type Receipt struct {
	Rec_id           int
	Name             string
	Quantity         int
	QuantityType     string
	Price            int
	PriceForQunatity float32
	Tag              string
}
