// Package vm provides a stack-based virtual machine.
package vm

import (
	"fmt"

	"github.com/danielspk/tatu-lang/pkg/runtime"
)

const (
	StackLimit = 512
)

// VirtualMachine ...
type VirtualMachine struct {
	ip        uint
	sp        uint
	code      []byte
	stack     [StackLimit]runtime.Value
	constants []runtime.Value
}

// NewVirtualMachine ...
func NewVirtualMachine() *VirtualMachine {
	return &VirtualMachine{
		ip:        0,
		sp:        0,
		code:      make([]byte, 0),
		constants: make([]runtime.Value, 0),
	}
}

// Execute ...
func (vm *VirtualMachine) Execute(code *Code) (runtime.Value, error) {
	// TODO test >>>>
	//vm.constants = append(vm.constants, NewNumber(2))
	//vm.constants = append(vm.constants, NewNumber(3))
	//vm.constants = append(vm.constants, NewString("hola "))
	//vm.constants = append(vm.constants, NewString("mundo"))

	//program = []byte{
	//	byte(OpConst), 0, byte(OpConst), 1, byte(OpAdd), 0x00,
	//}
	// TODO test <<<<

	// TODO check this
	vm.code = code.Code
	vm.constants = code.Constants

	// parsing the program

	// compile the program

	// init instruction pointer (or program counter)
	//vm.ip = int(vm.code[0])

	return vm.eval()
}

// eval ...
func (vm *VirtualMachine) eval() (runtime.Value, error) {
	for {
		op := opcode(vm.readByte())

		switch op {
		case OpHalt:
			return vm.stackPop(), nil

		case OpConst:
			constIdx := vm.readByte()
			value := vm.constants[constIdx]
			vm.stackPush(value)

		case OpAdd:
			op2 := vm.stackPop()
			op1 := vm.stackPop()

			if op1.Type() == runtime.NumberType && op2.Type() == runtime.NumberType {
				vm.stackPush(runtime.NewNumber(op1.(runtime.Number).Value + op2.(runtime.Number).Value))
			} else {
				vm.stackPush(runtime.NewString(op1.String() + op2.String()))
			}

		case OpSub:
			num1, num2 := vm.binaryOperation()
			vm.stackPush(runtime.NewNumber(num1 - num2))

		case OpMul:
			num1, num2 := vm.binaryOperation()
			vm.stackPush(runtime.NewNumber(num1 * num2))

		case OpDiv:
			num1, num2 := vm.binaryOperation()
			if num2 == 0 {
				return nil, fmt.Errorf("division by zero")
			}

			vm.stackPush(runtime.NewNumber(num1 / num2))

		default:
			return nil, fmt.Errorf("unknown opcode 0x%X", op)
		}
	}
}

// readByte ...
func (vm *VirtualMachine) readByte() byte {
	b := vm.code[vm.ip]
	vm.ip++

	return b
}

// stackPush ...
func (vm *VirtualMachine) stackPush(value runtime.Value) {
	if vm.sp == StackLimit {
		// TODO error stack overflow
	}

	vm.stack[vm.sp] = value
	vm.sp++
}

// stackPop ...
func (vm *VirtualMachine) stackPop() runtime.Value {
	if vm.sp == 0 {
		// TODO error empty stack
	}

	vm.sp--

	return vm.stack[vm.sp]
}

// binaryOperation ...
func (vm *VirtualMachine) binaryOperation() (float64, float64) {
	op2 := vm.stackPop()
	op1 := vm.stackPop()

	return op1.(runtime.Number).Value, op2.(runtime.Number).Value
}
