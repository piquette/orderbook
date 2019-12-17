package orderbook

// Execution is an execution report.
type Execution struct {
	OrderID           OrderID
	FilledQuantity    uint
	RemainingQuantity uint
}
