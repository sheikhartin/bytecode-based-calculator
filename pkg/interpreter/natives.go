package interpreter

import (
	"fmt"
	"math/rand"
)

const (
	PI = 3.141592653589793238462643
	E  = 2.718281828459045235360287
)

func Min(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("At least two arguments are expected!")
	}
	min_ := args[0].(float64)
	for _, arg := range args[1:] {
		if arg.(float64) < min_ {
			min_ = arg.(float64)
		}
	}
	return min_, nil
}

func Max(args ...interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("At least two arguments are expected!")
	}
	max_ := args[0].(float64)
	for _, arg := range args[1:] {
		if arg.(float64) > max_ {
			max_ = arg.(float64)
		}
	}
	return max_, nil
}

func Random(args ...interface{}) (interface{}, error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("Random does not take any arguments!")
	}
	return rand.Float64(), nil
}

func Factorial(args ...interface{}) (interface{}, error) {
	if len(args) != 1 || args[0].(float64) < 0 {
		return nil, fmt.Errorf("Factorial requires a non-negative integer...")
	} else if args[0].(float64) < 2 {
		return 1.0, nil
	}
	var res int = 1
	for i := 2; i <= int(args[0].(float64)); i++ {
		res *= i
	}
	return float64(res), nil
}
