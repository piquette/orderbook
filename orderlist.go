package orderbook

// orderList is a doubly-linked list of orders.
type orderList struct {
	first *order
	last  *order
	size  int
}

// New instantiates a new list and adds the passed values, if any, to the list
func newOrderList(orders ...*order) orderList {
	list := orderList{}
	if len(orders) > 0 {
		list.add(orders...)
	}
	return list
}

// Add appends a value (one or more) at the end of the list.
func (list *orderList) add(orders ...*order) {
	for _, o := range orders {
		//newElement := &element{value: value, prev: list.last}
		if list.size == 0 {
			list.first = o
			list.last = o
		} else {
			list.last.next = o
			list.last = o
		}
		list.size++
	}
}

// // Prepend prepends a values (or more)
// func (list *orderList) Prepend(values ...interface{}) {
// 	// in reverse to keep passed order i.e. ["c","d"] -> Prepend(["a","b"]) -> ["a","b","c",d"]
// 	for v := len(values) - 1; v >= 0; v-- {
// 		newElement := &element{value: values[v], next: list.first}
// 		if list.size == 0 {
// 			list.first = newElement
// 			list.last = newElement
// 		} else {
// 			list.first.prev = newElement
// 			list.first = newElement
// 		}
// 		list.size++
// 	}
// }

// get returns the element at index.
// Second return parameter is true if index is within bounds of the array and array is not empty, otherwise false.
func (list *orderList) get(index int) (*order, bool) {

	if !list.withinRange(index) {
		return nil, false
	}

	// determine traveral direction, last to first or first to last
	if list.size-index < index {
		element := list.last
		for e := list.size - 1; e != index; e, element = e-1, element.prev {
		}
		return element, true
	}
	element := list.first
	for e := 0; e != index; e, element = e+1, element.next {
	}
	return element, true
}

func (list *orderList) removeID(id OrderID) {

	var element *order

	element = list.first
	for {
		if element.id == id {
			break
		}
		element = element.next
	}

	if element == list.first {
		list.first = element.next
	}
	if element == list.last {
		list.last = element.prev
	}
	if element.prev != nil {
		element.prev.next = element.next
	}
	if element.next != nil {
		element.next.prev = element.prev
	}

	element = nil

	list.size--
}

// Remove removes the element at the given index from the list.
func (list *orderList) remove(index int) {

	if !list.withinRange(index) {
		return
	}

	if list.size == 1 {
		list.Clear()
		return
	}

	var element *order
	// determine traversal direction, last to first or first to last
	if list.size-index < index {
		element = list.last
		for e := list.size - 1; e != index; e, element = e-1, element.prev {
		}
	} else {
		element = list.first
		for e := 0; e != index; e, element = e+1, element.next {
		}
	}

	if element == list.first {
		list.first = element.next
	}
	if element == list.last {
		list.last = element.prev
	}
	if element.prev != nil {
		element.prev.next = element.next
	}
	if element.next != nil {
		element.next.prev = element.prev
	}

	element = nil

	list.size--
}

// Contains check if values (one or more) are present in the set.
// All values have to be present in the set for the method to return true.
// Performance time complexity of n^2.
// Returns true if no arguments are passed at all, i.e. set is always super-set of empty set.
// func (list *orderList) Contains(values ...interface{}) bool {

// 	if len(values) == 0 {
// 		return true
// 	}
// 	if list.size == 0 {
// 		return false
// 	}
// 	for _, value := range values {
// 		found := false
// 		for element := list.first; element != nil; element = element.next {
// 			if element.value == value {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			return false
// 		}
// 	}
// 	return true
// }

// Values returns all elements in the list.
func (list *orderList) Values() []*order {
	values := make([]*order, list.size, list.size)
	for e, element := 0, list.first; element != nil; e, element = e+1, element.next {
		values[e] = element
	}
	return values
}

//IndexOf returns index of provided element
func (list *orderList) IndexOf(value interface{}) int {
	if list.size == 0 {
		return -1
	}
	for index, element := range list.Values() {
		if element == value {
			return index
		}
	}
	return -1
}

// Empty returns true if list does not contain any elements.
func (list *orderList) Empty() bool {
	return list.size == 0
}

// Size returns number of elements within the list.
func (list *orderList) Size() int {
	return list.size
}

// Clear removes all elements from the list.
func (list *orderList) Clear() {
	list.size = 0
	list.first = nil
	list.last = nil
}

// // Sort sorts values (in-place) using.
// func (list *orderList) Sort(comparator utils.Comparator) {

// 	if list.size < 2 {
// 		return
// 	}

// 	values := list.Values()
// 	utils.Sort(values, comparator)

// 	list.Clear()

// 	list.Add(values...)

// }

// Swap swaps values of two elements at the given indices.
// func (list *orderList) Swap(i, j int) {
// 	if list.withinRange(i) && list.withinRange(j) && i != j {
// 		var element1, element2 *element
// 		for e, currentElement := 0, list.first; element1 == nil || element2 == nil; e, currentElement = e+1, currentElement.next {
// 			switch e {
// 			case i:
// 				element1 = currentElement
// 			case j:
// 				element2 = currentElement
// 			}
// 		}
// 		element1.value, element2.value = element2.value, element1.value
// 	}
// }

// Insert inserts values at specified index position shifting the value at that position (if any) and any subsequent elements to the right.
// Does not do anything if position is negative or bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
// func (list *orderList) Insert(index int, values ...interface{}) {

// 	if !list.withinRange(index) {
// 		// Append
// 		if index == list.size {
// 			list.add(values...)
// 		}
// 		return
// 	}

// 	list.size += len(values)

// 	var beforeElement *element
// 	var foundElement *element
// 	// determine traversal direction, last to first or first to last
// 	if list.size-index < index {
// 		foundElement = list.last
// 		for e := list.size - 1; e != index; e, foundElement = e-1, foundElement.prev {
// 			beforeElement = foundElement.prev
// 		}
// 	} else {
// 		foundElement = list.first
// 		for e := 0; e != index; e, foundElement = e+1, foundElement.next {
// 			beforeElement = foundElement
// 		}
// 	}

// 	if foundElement == list.first {
// 		oldNextElement := list.first
// 		for i, value := range values {
// 			newElement := &element{value: value}
// 			if i == 0 {
// 				list.first = newElement
// 			} else {
// 				newElement.prev = beforeElement
// 				beforeElement.next = newElement
// 			}
// 			beforeElement = newElement
// 		}
// 		oldNextElement.prev = beforeElement
// 		beforeElement.next = oldNextElement
// 	} else {
// 		oldNextElement := beforeElement.next
// 		for _, value := range values {
// 			newElement := &element{value: value}
// 			newElement.prev = beforeElement
// 			beforeElement.next = newElement
// 			beforeElement = newElement
// 		}
// 		oldNextElement.prev = beforeElement
// 		beforeElement.next = oldNextElement
// 	}
// }

// Set value at specified index position
// Does not do anything if position is negative or bigger than list's size
// Note: position equal to list's size is valid, i.e. append.
// func (list *orderList) Set(index int, value interface{}) {

// 	if !list.withinRange(index) {
// 		// Append
// 		if index == list.size {
// 			list.Add(value)
// 		}
// 		return
// 	}

// 	var foundElement *element
// 	// determine traversal direction, last to first or first to last
// 	if list.size-index < index {
// 		foundElement = list.last
// 		for e := list.size - 1; e != index; {
// 			fmt.Println("Set last", index, value, foundElement, foundElement.prev)
// 			e, foundElement = e-1, foundElement.prev
// 		}
// 	} else {
// 		foundElement = list.first
// 		for e := 0; e != index; {
// 			e, foundElement = e+1, foundElement.next
// 		}
// 	}

// 	foundElement.value = value
// }

// // String returns a string representation of container
// func (list *orderList) String() string {
// 	str := "DoublyLinkedList\n"
// 	values := []string{}
// 	for element := list.first; element != nil; element = element.next {
// 		values = append(values, fmt.Sprintf("%v", element.value))
// 	}
// 	str += strings.Join(values, ", ")
// 	return str
// }

// Check that the index is within bounds of the list
func (list *orderList) withinRange(index int) bool {
	return index >= 0 && index < list.size
}
