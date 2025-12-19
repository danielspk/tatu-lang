package vm

// opcode ...
type opcode byte

// ...
const (
	OpHalt  opcode = 0x00 // stops the program
	OpConst opcode = 0x01 // pushes a const onto the stack
	OpAdd   opcode = 0x02 // addition instruction
	OpSub   opcode = 0x03 // substraction instruction
	OpMul   opcode = 0x04 // multiplication instruction
	OpDiv   opcode = 0x05 // division instruction
)
