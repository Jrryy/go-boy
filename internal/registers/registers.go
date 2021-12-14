package registers

type Registers struct {
	A	 	byte
	F	 	byte
	B	 	byte
	C	 	byte
	D	 	byte
	E	 	byte
	H	 	byte
	L	 	byte
	PC 		int64 // It's not really a 64 bit integer but to read bytes from the game file we need it like this
	SP 		uint16
	stack 	[]byte
	flags 	byte
}

// This function initializes a new set of registers to their zero values (for the GB, ofc)
// after the checks that the GB is supposed to perform.
func InitializeRegisters() *Registers {
	r := new(Registers)
	r.A 	= 0x01
	r.F		= 0xB0
	r.B		= 0x00
	r.C		= 0x13
	r.D 	= 0x00
	r.E		= 0xD8
	r.H		= 0x01
	r.L		= 0x4D
	r.PC 	= 0x100
	r.SP 	= 0xE000
	r.stack = make([]byte, 0)
	r.flags = 0
	return r
}

func SetRegister(r* []byte, name rune, value byte) {

}