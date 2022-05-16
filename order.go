package main

import (
	"fmt"

	"deathbykeystroke.com/either/inventory"
)

func Order(id int, n int) inventory.InventoryOrderResult {
	return inventory.
		Quantity(id).
		Procure(id, n).
		Ship(id, n).
		Invoice()
}

type test struct {
	id, qty int
}

func main() {
	tests := []test{
		{id: 10, qty: 101},
		{id: 10, qty: 50},
		{id: 20, qty: 50},
		{id: 30, qty: 2},
		{id: 30, qty: 1},
		{id: 40, qty: 5},
	}
	for _, test := range tests {
		fmt.Printf("ordering %d, qty: %d => %+v\n", test.id, test.qty, Order(test.id, test.qty).String())
	}
}
