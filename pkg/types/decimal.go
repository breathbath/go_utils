package types

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

var ZERO = NewDecimalFromInt(0)

const (
	defaultDecimalLength = 13
	percentNumber        = 100
)

type Decimal struct {
	dec decimal.Decimal
}

func NewDecimalFromString(input string) Decimal {
	dec, e := decimal.NewFromString(input)
	if e != nil {
		dec = decimal.New(0, 0)
	}

	return Decimal{dec}
}

func NewDecimalFromFloat(input float64) Decimal {
	dec := decimal.NewFromFloat(input)

	return Decimal{dec}
}

func NewDecimalFromInt(input int64) Decimal {
	dec := decimal.New(input, 0)
	return Decimal{dec}
}

func (d Decimal) String() string {
	output := d.dec.StringFixed(defaultDecimalLength)
	output = strings.TrimRight(output, "0")
	output = strings.TrimRight(output, ".")

	return output
}

func (d Decimal) Add(dec Decimal) Decimal {
	newDec := d.dec.Add(dec.dec)
	return Decimal{newDec}
}

func (d Decimal) Sub(dec Decimal) Decimal {
	newDec := d.dec.Sub(dec.dec)
	return Decimal{newDec}
}

func (d Decimal) Div(dec Decimal) Decimal {
	newDec := d.dec.Div(dec.dec)
	return Decimal{newDec}
}

func (d Decimal) Mul(dec Decimal) Decimal {
	newDec := d.dec.Mul(dec.dec)
	return Decimal{newDec}
}

func (d Decimal) CeilToValue(dec Decimal) Decimal {
	if dec.EqualZero() {
		return d
	}

	newDec := d.dec.Div(dec.dec).Ceil().Mul(dec.dec)

	return Decimal{dec: newDec}
}

func (d Decimal) FloorToValue(dec Decimal) Decimal {
	if dec.EqualZero() {
		return d
	}

	newDec := d.dec.Div(dec.dec).Floor().Mul(dec.dec)

	return Decimal{dec: newDec}
}

func (d Decimal) Floor() Decimal {
	return Decimal{dec: d.dec.Floor()}
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error {
	return d.dec.UnmarshalJSON(decimalBytes)
}

// Scan implements the Scanner interface.
func (d *Decimal) Scan(value interface{}) error {
	if value == nil {
		d.dec = decimal.New(0, 0)
		return nil
	}
	// we force to convert to string because values coming from db produce the following problems:
	/**
	decimal.NewFromFloat(0.357).Div(decimal.NewFromFloat(0.001)).Floor() == 356 or
	in other words floor(0.357 / 0.001) == 356 which should be 357

	but if we do it with strings
	p1,_ := decimal.NewFromString("0.357")
	one,_ := decimal.NewFromString("0.001")
	p1.Div(one).Floor() == 357
	*/
	switch value.(type) {
	case float32:
		value = fmt.Sprintf("%.13f", value)
	case float64:
		value = fmt.Sprintf("%.13f", value)
	}
	return d.dec.Scan(value)
}

// Value implements the driver Valuer interface.
func (d Decimal) Value() (driver.Value, error) {
	return d.dec.Value()
}

func (d Decimal) LessOrEqual(dec Decimal) bool {
	return d.dec.LessThanOrEqual(dec.dec)
}

func (d Decimal) Equal(dec Decimal) bool {
	return d.dec.Equal(dec.dec)
}

func (d Decimal) GreaterOrEqual(dec Decimal) bool {
	return d.dec.GreaterThanOrEqual(dec.dec)
}

func (d Decimal) Less(dec Decimal) bool {
	return d.dec.LessThan(dec.dec)
}

func (d Decimal) Greater(dec Decimal) bool {
	return d.dec.GreaterThan(dec.dec)
}

func (d Decimal) LessOrEqualZero() bool {
	return d.dec.LessThanOrEqual(ZERO.dec)
}

func (d Decimal) GreaterOrEqualZero() bool {
	return d.dec.GreaterThanOrEqual(ZERO.dec)
}

func (d Decimal) EqualZero() bool {
	return d.dec.Equal(ZERO.dec)
}

func (d Decimal) LessZero() bool {
	return d.dec.LessThan(ZERO.dec)
}

func (d Decimal) GreaterZero() bool {
	return d.dec.GreaterThan(ZERO.dec)
}

func (d Decimal) IncrementByPercent(dec Decimal) Decimal {
	realPercent := dec.Add(NewDecimalFromInt(1))
	return d.Mul(realPercent)
}

func (d Decimal) DecrementByPercent(dec Decimal) Decimal {
	realPercent := NewDecimalFromInt(1).Sub(dec)
	return d.Mul(realPercent)
}

func (d Decimal) EqualInt(input int64) bool {
	dec := NewDecimalFromInt(input)
	return d.Equal(dec)
}

func (d Decimal) GreaterInt(input int64) bool {
	dec := NewDecimalFromInt(input)
	return d.Greater(dec)
}

func (d Decimal) LowerOrEqualInt(input int64) bool {
	dec := NewDecimalFromInt(input)
	return d.LessOrEqual(dec)
}

func (d Decimal) GreaterOrEqualInt(input int64) bool {
	dec := NewDecimalFromInt(input)
	return d.GreaterOrEqual(dec)
}

func (d Decimal) Round(places int64) Decimal {
	if places < 0 {
		return d
	}

	roundResult := d.dec.Round(int32(places))
	return Decimal{dec: roundResult}
}

func (d Decimal) ToPercent(places int64) Decimal {
	percent := d.Mul(NewDecimalFromInt(percentNumber))
	return percent.Round(places)
}

func (d Decimal) ToFloat() float64 {
	fl, _ := d.dec.Float64()
	return fl
}
