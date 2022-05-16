package inventory

import (
	"github.com/pkg/errors"

	"deathbykeystroke.com/either/either"
)

var inventory map[int]*Inventory = map[int]*Inventory{
	10: &Inventory{onHand: 100},
	20: &Inventory{onHand: 50},
	30: &Inventory{onHand: 1},
	40: &Inventory{onHand: 0},
}

type OrderResult int

const (
	UnableToProcure OrderResult = iota
	Ordered
	Fulfilled
)

var OrderResultStrings = map[OrderResult]string{
	UnableToProcure: "Unable to procure",
	Ordered:         "Ordered",
	Fulfilled:       "Fulfilled",
}

type Inventory struct {
	onHand, onOrder int
}

type InventoryQuantity either.Either[*Inventory, error]
type InventoryOrderResult either.Either[OrderResult, error]
type InventoryProcurementResult either.Either[Inventory, error]

func (i InventoryQuantity) Which() either.WHICH { return either.Either[*Inventory, error](i).Which() }
func (i InventoryQuantity) Left() *Inventory    { return either.Either[*Inventory, error](i).Left() }
func (i InventoryQuantity) Right() error        { return either.Either[*Inventory, error](i).Right() }

func (i InventoryOrderResult) Which() either.WHICH {
	return either.Either[OrderResult, error](i).Which()
}
func (i InventoryOrderResult) Left() OrderResult { return either.Either[OrderResult, error](i).Left() }
func (i InventoryOrderResult) Right() error      { return either.Either[OrderResult, error](i).Right() }

func Quantity(id int) InventoryQuantity {
	if inv, ok := inventory[id]; ok {
		return InventoryQuantity(either.Left[*Inventory, error](inv))
	}
	return InventoryQuantity(either.Right[*Inventory, error](errors.Errorf("no record of id: %d", id)))
}

func (i InventoryQuantity) Procure(id, n int) InventoryQuantity {
	if i.Which() == either.RIGHT {
		return i
	}
	invN := i.Left()
	if id >= 30 && invN.onHand < n && invN.onOrder < n {
		return InventoryQuantity(either.Right[*Inventory, error](errors.Errorf("cannot order id: %d", id)))
	} else if n > invN.onHand {
		// do some procurement process
		(*invN).onOrder += n
		return InventoryQuantity(either.Left[*Inventory, error](invN))
	} else {
		return InventoryQuantity(either.Left[*Inventory, error](invN))
	}
	return InventoryQuantity(either.Right[*Inventory, error](errors.Errorf("no record of id: %d", id)))
}

func (i InventoryQuantity) Ship(id, n int) InventoryOrderResult {
	inv := i.Left()
	if i.Which() == either.RIGHT || (inv.onHand < n && inv.onOrder < n) {
		return InventoryOrderResult(either.Left[OrderResult, error](UnableToProcure))
	}
	if inv.onHand >= n {
		(*inv).onHand -= n
		return InventoryOrderResult(either.Left[OrderResult, error](Fulfilled))
	}

	(*inv).onOrder -= n
	return InventoryOrderResult(either.Left[OrderResult, error](Ordered))
}

func (i InventoryOrderResult) Invoice() InventoryOrderResult {
	return i
}

func (i InventoryOrderResult) String() string {
	if i.Which() == either.RIGHT {
		return i.Right().Error()
	}
	return OrderResultStrings[i.Left()]
}
