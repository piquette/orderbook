package orderbook

import "fmt"

// compare evals uints.
func compare(a, b uint) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// limitPriceTree is a limit level in the order book.
type limitPriceTree struct {
	root *limitPrice
	size int
}

// Put inserts node into the tree.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (t *limitPriceTree) addLimit(price uint, order *order) *limitPrice {
	lim := &limitPrice{price: price, orders: newOrderList(order)}
	t.put(lim, nil, &t.root)
	return lim
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (t *limitPriceTree) removeLimit(price uint) {
	t.remove(price, &t.root)
}

// getOrders searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
// func (t *limitPriceTree) getOrders(key uint) (value *orderList, found bool) {
// 	n := t.root
// 	for n != nil {
// 		cmp := compare(key, n.price)
// 		switch {
// 		case cmp == 0:
// 			return &n.orders, true
// 		case cmp < 0:
// 			n = n.children[0]
// 		case cmp > 0:
// 			n = n.children[1]
// 		}
// 	}
// 	return nil, false
// }

// Clear removes all nodes from the tree.
func (t *limitPriceTree) clear() {
	t.root = nil
	t.size = 0
}

func (t *limitPriceTree) put(lim *limitPrice, p *limitPrice, qp **limitPrice) bool {
	q := *qp
	if q == nil {
		t.size++
		lim.parent = p
		*qp = lim
		return true
	}

	c := compare(lim.price, q.price)
	if c == 0 {
		// q.price = key
		// q.orders.add(order)
		return false
	}

	if c < 0 {
		c = -1
	} else {
		c = 1
	}
	a := (c + 1) / 2
	var fix bool
	fix = t.put(lim, q, &q.children[a])
	if fix {
		return putFix(int8(c), qp)
	}
	return false
}

func (t *limitPriceTree) remove(key uint, qp **limitPrice) bool {
	q := *qp
	if q == nil {
		return false
	}

	c := compare(key, q.price)
	if c == 0 {
		t.size--
		if q.children[1] == nil {
			if q.children[0] != nil {
				q.children[0].parent = q.parent
			}
			*qp = q.children[0]
			return true
		}
		fix := removeMin(&q.children[1], &q.price, &q.orders)
		if fix {
			return removeFix(-1, qp)
		}
		return false
	}

	if c < 0 {
		c = -1
	} else {
		c = 1
	}
	a := (c + 1) / 2
	fix := t.remove(key, &q.children[a])
	if fix {
		return removeFix(int8(-c), qp)
	}
	return false
}

func removeMin(qp **limitPrice, minKey *uint, minVal *orderList) bool {
	q := *qp
	if q.children[0] == nil {
		*minKey = q.price
		minVal = &q.orders
		if q.children[1] != nil {
			q.children[1].parent = q.parent
		}
		*qp = q.children[1]
		return true
	}
	fix := removeMin(&q.children[0], minKey, minVal)
	if fix {
		return removeFix(1, qp)
	}
	return false
}

func putFix(c int8, t **limitPrice) bool {
	s := *t
	if s.b == 0 {
		s.b = c
		return true
	}

	if s.b == -c {
		s.b = 0
		return false
	}

	if s.children[(c+1)/2].b == c {
		s = singlerot(c, s)
	} else {
		s = doublerot(c, s)
	}
	*t = s
	return false
}

func removeFix(c int8, t **limitPrice) bool {
	s := *t
	if s.b == 0 {
		s.b = c
		return false
	}

	if s.b == -c {
		s.b = 0
		return true
	}

	a := (c + 1) / 2
	if s.children[a].b == 0 {
		s = rotate(c, s)
		s.b = -c
		*t = s
		return false
	}

	if s.children[a].b == c {
		s = singlerot(c, s)
	} else {
		s = doublerot(c, s)
	}
	*t = s
	return true
}

func singlerot(c int8, s *limitPrice) *limitPrice {
	s.b = 0
	s = rotate(c, s)
	s.b = 0
	return s
}

func doublerot(c int8, s *limitPrice) *limitPrice {
	a := (c + 1) / 2
	r := s.children[a]
	s.children[a] = rotate(-c, s.children[a])
	p := rotate(c, s)

	switch {
	default:
		s.b = 0
		r.b = 0
	case p.b == c:
		s.b = -c
		r.b = 0
	case p.b == -c:
		s.b = 0
		r.b = c
	}

	p.b = 0
	return p
}

func rotate(c int8, s *limitPrice) *limitPrice {
	a := (c + 1) / 2
	r := s.children[a]
	s.children[a] = r.children[a^1]
	if s.children[a] != nil {
		s.children[a].parent = s
	}
	r.children[a^1] = s
	r.parent = s.parent
	s.parent = r
	return r
}

// String returns a string representation of container
func (t *limitPriceTree) string() string {
	str := "\nTree:\n"
	if t.size > 0 {
		output(t.root, "", true, &str)
	}
	return str
}

func (lim *limitPrice) string() string {
	return fmt.Sprintf("%v (%v)", lim.price, lim.orders.size)
}

func output(lim *limitPrice, prefix string, isTail bool, str *string) {
	if lim.children[1] != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(lim.children[1], newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += lim.string() + "\n"
	if lim.children[0] != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(lim.children[0], newPrefix, true, str)
	}
}
