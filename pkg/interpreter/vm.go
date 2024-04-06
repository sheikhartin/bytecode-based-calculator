package interpreter

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type VM struct {
	currOperands []string
	Stack        []interface{}
	Vars         map[string]interface{}
}

func (vm VM) String() string {
	return fmt.Sprintf("%v", vm.Stack[len(vm.Stack)-1])
}

func (vm *VM) insertNumber() error {
	value, err := strconv.ParseFloat(vm.currOperands[0], 64)
	if err != nil {
		return err
	}
	vm.Stack = append(vm.Stack, value)
	return nil
}

func (vm *VM) loadVariable() error {
	value, ok := vm.Vars[vm.currOperands[0]]
	if !ok {
		return fmt.Errorf("Undefined variable: %s", vm.currOperands[0])
	}
	vm.Stack = append(vm.Stack, value)
	return nil
}

func (vm *VM) callFunction() error {
	argCount, err := strconv.Atoi(vm.currOperands[1])
	if err != nil {
		return fmt.Errorf("Invalid argument count: %s", vm.currOperands[1])
	} else if len(vm.Stack) < argCount {
		return fmt.Errorf("Not enough arguments on stack for function `%s`!", vm.currOperands[0])
	}

	fn, found := vm.Vars[vm.currOperands[0]]
	if !found {
		return fmt.Errorf("Function `%s` not found!", vm.currOperands[0])
	}
	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		return fmt.Errorf("`%s` is not callable!", vm.currOperands[0])
	}

	var callArgs []reflect.Value
	for i := 0; i < argCount; i++ {
		callArgs = append(callArgs, reflect.ValueOf(vm.Stack[len(vm.Stack)-argCount+i]))
	}
	result := fnValue.Call(callArgs)
	if err, ok := result[1].Interface().(error); ok && err != nil {
		return err
	}

	vm.Stack = vm.Stack[:len(vm.Stack)-argCount]
	vm.Stack = append(vm.Stack, result[0].Interface())
	return nil
}

func (vm *VM) performUnaryOperation() error {
	if len(vm.Stack) < 1 {
		return fmt.Errorf("Stack underflow!")
	}
	operand := vm.Stack[len(vm.Stack)-1]

	var result interface{}
	var err error
	switch vm.currOperands[0] {
	case "FACT":
		result, err = Factorial(operand)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown unary operation: %s", vm.currOperands[0])
	}

	vm.Stack = vm.Stack[:len(vm.Stack)-1]
	vm.Stack = append(vm.Stack, result)
	return nil
}

func (vm *VM) performBinaryOperation() error {
	if len(vm.Stack) < 2 {
		return fmt.Errorf("Stack underflow!")
	}
	right := vm.Stack[len(vm.Stack)-1]
	left := vm.Stack[len(vm.Stack)-2]

	var result float64
	switch vm.currOperands[0] {
	case "ADD":
		result = left.(float64) + right.(float64)
	case "SUB":
		result = left.(float64) - right.(float64)
	case "MUL":
		result = left.(float64) * right.(float64)
	case "DIV":
		if right.(float64) == 0 {
			return fmt.Errorf("Division by zero!?")
		}
		result = left.(float64) / right.(float64)
	case "MOD":
		if right.(float64) == 0 {
			return fmt.Errorf("Division by zero!?")
		}
		result = math.Mod(left.(float64), right.(float64))
	case "POW":
		result = math.Pow(left.(float64), right.(float64))
	default:
		return fmt.Errorf("Unknown binary operation: %s", vm.currOperands[0])
	}

	vm.Stack = vm.Stack[:len(vm.Stack)-2]
	vm.Stack = append(vm.Stack, result)
	return nil
}

func (vm *VM) setVariable() error {
	if len(vm.Stack) < 1 {
		return fmt.Errorf("Stack underflow!")
	}
	vm.Vars[vm.currOperands[0]] = vm.Stack[len(vm.Stack)-1]
	return nil
}

func (vm *VM) Execute(instructions []string) error {
	for _, instr := range instructions {
		parts := strings.Split(instr, "\t")
		vm.currOperands = parts[1:]

		switch parts[0] {
		case "PUSH_NUM":
			if err := vm.insertNumber(); err != nil {
				return err
			}
		case "LOAD_VAR":
			if err := vm.loadVariable(); err != nil {
				return err
			}
		case "CALL_FUNC":
			if err := vm.callFunction(); err != nil {
				return err
			}
		case "UNARY_OP":
			if err := vm.performUnaryOperation(); err != nil {
				return err
			}
		case "BINARY_OP":
			if err := vm.performBinaryOperation(); err != nil {
				return err
			}
		case "STORE_VAR":
			if err := vm.setVariable(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unknown instruction: %s", parts[0])
		}
	}
	return nil
}

func NewVM() *VM {
	return &VM{
		Vars: map[string]interface{}{
			"PI":   PI,
			"E":    E,
			"rand": Random,
			"fact": Factorial,
			"min":  Min,
			"max":  Max,
		},
	}
}
