package main

import (
	"testing"

	"github.com/shopspring/decimal"
)

type DT struct {
	A decimal.Decimal
	B decimal.Decimal
}

func TestDecimal(t *testing.T) {
	decimal.DivisionPrecision = 32
	price, err := decimal.NewFromString("136.02")
	if err != nil {
		panic(err)
	}

	quantity := decimal.NewFromInt(3)

	var big_2 decimal.Decimal
	big_2, _ = decimal.NewFromString("2.000000")
	var big_3 decimal.Decimal = decimal.NewFromInt(16)

	fee, _ := decimal.NewFromString(".035")
	taxRate, _ := decimal.NewFromString(".08875")

	subtotal := price.Mul(quantity)

	preTax := subtotal.Mul(fee.Add(decimal.NewFromFloat(1)))

	total := preTax.Mul(taxRate.Add(decimal.NewFromFloat(1)))

	bigBig := big_2.Div(big_3)

	dt := DT{
		A: big_2,
		B: big_3,
	}

	t.Logf("Subtotal:%s", subtotal)                      // Subtotal: 408.06
	t.Logf("Pre-tax:%s", preTax)                         // Pre-tax: 422.3421
	t.Logf("Taxes:%s", total.Sub(preTax))                // Taxes: 37.482861375
	t.Logf("Total:%s", total)                            // Total: 459.824961375
	t.Logf("Tax rate:%s", total.Sub(preTax).Div(preTax)) // Tax rate: 0.08875
	t.Logf("big_2:%s", big_2)                            // Tax rate: 0.08875
	t.Logf("bigBig:%s", bigBig)
	t.Logf("dt result:%s", dt.A.Mul(dt.B))
	t.Logf("%s", decimal.NewFromFloat(-1e13).String())
}
