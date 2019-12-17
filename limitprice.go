package orderbook

// limitPrice is a single price limit.
type limitPrice struct {
	price    uint
	orders   orderList
	parent   *limitPrice
	children [2]*limitPrice
	b        int8
}

func (l *limitPrice) higher() *limitPrice { return l.children[1] }
func (l *limitPrice) lower() *limitPrice  { return l.children[0] }
