package user

import (
	"time"
)

// Person .
type Person struct {
	doc      struct{} `doc:"some document for type person"`
	Name     *string
	Nickname string         `json:"nickname,omitempty"`
	Age      int            `doc:"some document for field age"`
	Estates  []PersonEstate `json:"estates"`
}

// PersonEstate .
type PersonEstate struct {
	Amount float64
	Time   XTime
}

// XTime .
type XTime time.Time
