package Fractal

import (
	"fmt"
)

type DataBase map[string]Dollars

type Dollars float32

func (d Dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}