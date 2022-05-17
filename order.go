package main

import (
	"fmt"
	"os"

	"deathbykeystroke.com/either/either"
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

var passFail = map[bool]string{false: "FAIL", true: "PASS"}

func main() {
	tests := []test{
		{id: 10, qty: 101},
		{id: 10, qty: 50},
		{id: 20, qty: 50},
		{id: 30, qty: 2},
		{id: 30, qty: 1},
		{id: 40, qty: 5},
	}
	expects := []inventory.OrderResult{inventory.Ordered, inventory.Fulfilled, inventory.Fulfilled, inventory.UnableToProcure, inventory.Fulfilled, inventory.UnableToProcure}
	pass := true
	for i, test := range tests {
		result := Order(test.id, test.qty)
		if result.Which() == either.RIGHT || result.Left() != expects[i] {
			pass = false
		}
		fmt.Printf("[%s] ordering %d, qty: %d => %+v\n", passFail[result.Which() == either.RIGHT || expects[i] == result.Left()], test.id, test.qty, result.String())
	}
	if !pass {
		fmt.Fprintf(os.Stderr, "\nFAIL\n")
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "\nPASS\n")
	os.Exit(0)
}
