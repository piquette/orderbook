package orderbook

import "errors"

// Book is a limit-price orderbook for a particular instrument,
// that matches buys and sells in continuous time.
type Book struct {
	bidTree  *limitPriceTree
	askTree  *limitPriceTree
	bestBid  *limitPrice
	bestAsk  *limitPrice
	orderMap map[OrderID]*order
	bidMap   map[uint]*limitPrice
	askMap   map[uint]*limitPrice
}

// Init initializes a new order book.
func Init() *Book {
	return &Book{
		bidTree:  &limitPriceTree{},
		askTree:  &limitPriceTree{},
		orderMap: make(map[OrderID]*order),
		bidMap:   make(map[uint]*limitPrice),
		askMap:   make(map[uint]*limitPrice),
	}
}

// Submit submits a single order to the order book. The order must have a
// specified Side, either Buy or Sell, a nonzero price and nonzero size.
//
// The return values are:
// * The order id for this order created during the matching process.
// * A list a order executions that, if order matching was possible, will include the execution details
//   of the orders that are matched, including the originally submitted order.
// * An optional error.
func (b *Book) Submit(side Side, price uint, size uint) (OrderID, []Execution, error) {
	var (
		matchedQty uint
		matches    []Execution
	)
	if price == 0 || size == 0 {
		return 0, matches, errors.New("price/size cannot be zero")
	}

	// General methodology:
	// Check if we can match immediately at the best bid/offer,
	// taking liquidity up to the price limit specified.
	// handle partial fills,
	// report order executions
	// handle inserting into the book if we cant fill the entire order.

	newOrderID := genID()

	if side == Bid {
		// Looking to buy, match aginst existing asks.
		matchedQty, matches = b.matchBid(price, size)
	} else {
		// Looking to sell, match against existing bids.
		matchedQty, matches = b.matchAsk(price, size)
	}

	// Report the fact that the submitted order
	// has been at least partially matched.
	if matchedQty != 0 {
		matches = append(matches, Execution{OrderID: newOrderID,
			FilledQuantity:    matchedQty,
			RemainingQuantity: size - matchedQty})
	}

	// Add new order to the book if the new order wasn't completely filled.
	if matchedQty != size {
		o := &order{
			id:    newOrderID,
			price: price,
			size:  size - matchedQty,
		}
		b.orderMap[newOrderID] = o

		if side == Bid {
			// Check if the price limit already exists.
			lim, ok := b.bidMap[price]
			if !ok {
				// Doesn't exist, add new limit/order to tree.
				lim = b.bidTree.addLimit(price, o)
			} else {
				// Exists, just add the order to it.
				lim.orders.add(o)
			}
			// Adjust best bid if needed.
			if b.bestBid == nil || price > b.bestBid.price {
				b.bestBid = lim
			}
		} else {
			// Check if the price limit already exists.
			lim, ok := b.askMap[price]
			if !ok {
				// Doesn't exist, add new limit/order to tree.
				lim = b.askTree.addLimit(price, o)
			} else {
				// Exists, just add the order to it.
				lim.orders.add(o)
			}
			// Adjust best bid if needed.
			if b.bestAsk == nil || price > b.bestAsk.price {
				b.bestAsk = lim
			}
		}
	}

	return newOrderID, matches, nil
}

func (b *Book) matchBid(bidPrice, bidSize uint) (uint, []Execution) {

	// Matching methodology:
	//
	// find best ask price.
	// iterate through existing ask orders at that price.
	// if an ask order is filled, remove from the
	// order list and remove from the order map.
	// keep iterating until the bid is filled.
	// if the buy cannot be filled entirely at this price, advance to the next-best ask price.
	// if this happens, delete the old best ask limit from the ask tree
	// and delete the ask limit from the price map
	// repeat this process until the next ask limit is higher than the bidPrice,
	// there are no more ask limits, or the bid is filled.

	matches := []Execution{}

	remaining := bidSize
	for remaining != 0 {
		if b.bestAsk == nil || bidPrice < b.bestAsk.price {
			// Cant match, exit.
			break
		}

		potentialmatch := b.bestAsk.orders.first
		if potentialmatch == nil {
			// No orders left at this level. advance to the next limit price.
			next := b.bestAsk.higher()
			b.askTree.removeLimit(b.bestAsk.price)
			delete(b.askMap, b.bestAsk.price)
			b.bestAsk = next
			continue
		}

		if potentialmatch.size <= remaining {
			// Fill existing order and remove from order map.
			remaining -= potentialmatch.size
			delete(b.orderMap, potentialmatch.id)
			b.bestAsk.orders.remove(0)
			matches = append(matches, Execution{
				OrderID:           potentialmatch.id,
				FilledQuantity:    potentialmatch.size,
				RemainingQuantity: 0,
			})
		} else {
			// Partial fill the existing ask order and move on.
			matches = append(matches, Execution{
				OrderID:           potentialmatch.id,
				FilledQuantity:    remaining,
				RemainingQuantity: potentialmatch.size - remaining,
			})
			potentialmatch.size -= remaining
			remaining = 0
		}
	}

	return bidSize - remaining, matches
}

func (b *Book) matchAsk(askPrice, askSize uint) (uint, []Execution) {

	matches := []Execution{}
	remaining := askSize

	for remaining != 0 {
		if b.bestBid == nil || askPrice > b.bestBid.price {
			// Cant match, exit.
			break
		}

		potentialmatch := b.bestBid.orders.first
		if potentialmatch == nil {
			// No orders left at this price. advance to the next limit price.
			next := b.bestBid.lower()
			b.bidTree.removeLimit(b.bestBid.price)
			delete(b.bidMap, b.bestBid.price)
			b.bestBid = next
			continue
		}

		if potentialmatch.size <= remaining {
			// Fill existing bid order and remove from map.
			remaining -= potentialmatch.size
			delete(b.orderMap, potentialmatch.id)
			b.bestBid.orders.remove(0)
			matches = append(matches, Execution{
				OrderID:           potentialmatch.id,
				FilledQuantity:    potentialmatch.size,
				RemainingQuantity: 0,
			})
		} else {
			// Partial fill the existing bid order and move on.
			matches = append(matches, Execution{
				OrderID:           potentialmatch.id,
				FilledQuantity:    remaining,
				RemainingQuantity: potentialmatch.size - remaining,
			})
			potentialmatch.size -= remaining
			remaining = 0
		}
	}

	return askSize - remaining, matches
}

// Cancel order.
func (b *Book) Cancel(id OrderID) (bool, error) {

	// Check existence in map and return if not in.
	order, orderExists := b.orderMap[id]
	if !orderExists {
		return false, errors.New("order does not exist")
	}

	p := order.price
	s := order.side

	// Remove order from orders map and remove order from its list.
	var (
		lim         *limitPrice
		priceExists bool
	)
	if s == Bid {
		lim, priceExists = b.bidMap[p]
	} else {
		lim, priceExists = b.askMap[p]
	}
	if !priceExists {
		return false, errors.New("price does not exist, this is should not happen")
	}

	delete(b.orderMap, id)
	lim.orders.removeID(id)

	// If that order was the last in its price level, remove that price level
	// from the price level map and from the bid/ask tree.
	//
	// If a removed price level is the best bid/ask, the best bid/ask need to
	// be replaced with the next best.
	if lim.orders.Size() == 0 {
		if s == Bid {
			//
			if b.bestBid.price == p {
				b.bestBid = b.bestBid.lower()
			}
			delete(b.bidMap, p)
			b.bidTree.removeLimit(p)
		} else {
			//
			if b.bestAsk.price == p {
				b.bestAsk = b.bestAsk.higher()
			}
			delete(b.askMap, p)
			b.askTree.removeLimit(p)
		}
	}
	return true, nil
}

// Top of the book.
func (b *Book) Top() (bid, ask uint) {
	return b.bestBid.price, b.bestAsk.price
}
