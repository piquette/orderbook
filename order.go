package orderbook

import (
	"math/rand"
)

// OrderID is an order id.
type OrderID int

// Side is a side.
type Side bool

// Bid means you are buying.
const Bid Side = false

// Ask means you are selling.
const Ask Side = true

// order is a single order in the book.
type order struct {
	id    OrderID
	side  Side
	price uint
	size  uint
	next  *order
	prev  *order
}

// genID returns a psuedo-random 8-digit order id.
func genID() OrderID {
	return OrderID(10000000 + rand.Intn(99999999-10000000))
}
