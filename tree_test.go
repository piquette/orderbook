package orderbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Add(t *testing.T) {

	tree := limitPriceTree{}
	tree.addLimit(1000, &order{})
	tt := tree.addLimit(1010, &order{})
	tt.orders.add(&order{})

	//tree.addLimit(1010, &order{})
	tree.addLimit(9010, &order{})
	tree.addLimit(900, &order{})
	tree.addLimit(1009, &order{})
	tree.addLimit(1008, &order{})

	s := tree.string()
	t.Log(s)

	cool := true

	assert.True(t, cool)
}
