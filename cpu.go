package main

type cpu6502 struct {
	bus readableWriteable
}

func createCPU(b readableWriteable) *cpu6502 {
	return &cpu6502{
		bus: b,
	}
}
