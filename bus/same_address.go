package bus

type SameAddress struct {
	readable Readable
	writeable Writeable
}

var _ ReadableWriteable = new(SameAddress)

func UseSameAdress(r Readable, w Writeable) *SameAddress {
	return &SameAddress{
		readable:  r,
		writeable: w,
	}
}

func (s *SameAddress) Read(addr uint16) uint8 {
	return s.readable.Read(addr)
}

func (s *SameAddress) Write(addr uint16, data uint8) {
	s.writeable.Write(addr, data)
}

