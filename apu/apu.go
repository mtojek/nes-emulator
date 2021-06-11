package apu

import (
	"log"

	"github.com/mtojek/nes-emulator/bus"
)

const cpuFrequency = 1789773
const frameCounterRate = cpuFrequency / 240.0

var lengthTable = []byte{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}

var dutyTable = [][]byte{
	{0, 1, 0, 0, 0, 0, 0, 0},
	{0, 1, 1, 0, 0, 0, 0, 0},
	{0, 1, 1, 1, 1, 0, 0, 0},
	{1, 0, 0, 1, 1, 1, 1, 1},
}

var pulseTable [31]float32
var tndTable [203]float32

func init() {
	for i := 0; i < 31; i++ {
		pulseTable[i] = 95.52 / (8128.0/float32(i) + 100)
	}
	for i := 0; i < 203; i++ {
		tndTable[i] = 163.67 / (24329.0/float32(i) + 100)
	}
}

type APU2303 struct {
	TriggerIRQ bool
	dmcModer   dmcModer
	channel    chan float32

	pulse1   *Pulse
	pulse2   *Pulse
	noise    *Noise
	triangle *Triangle
	dmc      *DMC

	framePeriod byte
	frameValue  byte
	frameIRQ    bool

	filterChain FilterChain
	sampleRate float64
	cycle      uint64
}

var _ bus.ReadableWriteable = new(APU2303)

type dmcModer interface {
	DMCMode()
}

func Create(cpuBus bus.ReadableWriteable, dmcModer dmcModer) *APU2303 {
	return &APU2303{
		channel:  make(chan float32, 44100),
		dmcModer: dmcModer,
		sampleRate: float64(cpuFrequency)/48000,
		filterChain: FilterChain{
			HighPassFilter(48000, 90),
			HighPassFilter(48000, 440),
			LowPassFilter(48000, 14000),
		},

		pulse1:   &Pulse{
			channel: 1,
		},
		pulse2:   &Pulse{
			channel: 2,
		},
		noise:    &Noise{
			shiftRegister: 1,
		},
		triangle: new(Triangle),
		dmc: &DMC{
			cpuBus: cpuBus,
			dmcModer: dmcModer,
		},
	}
}

func (apu *APU2303) Read(addr uint16) uint8 {
	switch addr {
	case 0x4015:
		return apu.readStatus()
	}
	log.Printf("APU: read from unmapped address %04X\n", addr)
	return 0
}

func (apu *APU2303) Write(addr uint16, value uint8) {
	switch addr {
	case 0x4000:
		apu.pulse1.writeControl(value)
	case 0x4001:
		apu.pulse1.writeSweep(value)
	case 0x4002:
		apu.pulse1.writeTimerLow(value)
	case 0x4003:
		apu.pulse1.writeTimerHigh(value)
	case 0x4004:
		apu.pulse2.writeControl(value)
	case 0x4005:
		apu.pulse2.writeSweep(value)
	case 0x4006:
		apu.pulse2.writeTimerLow(value)
	case 0x4007:
		apu.pulse2.writeTimerHigh(value)
	case 0x4008:
		apu.triangle.writeControl(value)
	case 0x4009:
	case 0x4010:
		apu.dmc.writeControl(value)
	case 0x4011:
		apu.dmc.writeValue(value)
	case 0x4012:
		apu.dmc.writeAddress(value)
	case 0x4013:
		apu.dmc.writeLength(value)
	case 0x400A:
		apu.triangle.writeTimerLow(value)
	case 0x400B:
		apu.triangle.writeTimerHigh(value)
	case 0x400C:
		apu.noise.writeControl(value)
	case 0x400D:
	case 0x400E:
		apu.noise.writePeriod(value)
	case 0x400F:
		apu.noise.writeLength(value)
	case 0x4015:
		apu.writeControl(value)
	case 0x4017:
		apu.writeFrameCounter(value)
	}
}

func (apu *APU2303) writeControl(value uint8) {
	apu.pulse1.enabled = value&1 == 1
	apu.pulse2.enabled = value&2 == 2
	apu.triangle.enabled = value&4 == 4
	apu.noise.enabled = value&8 == 8
	apu.dmc.enabled = value&16 == 16
	if !apu.pulse1.enabled {
		apu.pulse1.lengthValue = 0
	}
	if !apu.pulse2.enabled {
		apu.pulse2.lengthValue = 0
	}
	if !apu.triangle.enabled {
		apu.triangle.lengthValue = 0
	}
	if !apu.noise.enabled {
		apu.noise.lengthValue = 0
	}
	if !apu.dmc.enabled {
		apu.dmc.currentLength = 0
	} else {
		if apu.dmc.currentLength == 0 {
			apu.dmc.restart()
		}
	}
}

func (apu *APU2303) writeFrameCounter(value byte) {
	apu.framePeriod = 4 + (value>>7)&1
	apu.frameIRQ = (value>>6)&1 == 0
	// apu.frameValue = 0
	if apu.framePeriod == 5 {
		apu.stepEnvelope()
		apu.stepSweep()
		apu.stepLength()
	}
}

func (apu *APU2303) readStatus() uint8 {
	var result byte
	if apu.pulse1.lengthValue > 0 {
		result |= 1
	}
	if apu.pulse2.lengthValue > 0 {
		result |= 2
	}
	if apu.triangle.lengthValue > 0 {
		result |= 4
	}
	if apu.noise.lengthValue > 0 {
		result |= 8
	}
	if apu.dmc.currentLength > 0 {
		result |= 16
	}
	return result
}

func (apu *APU2303) stepEnvelope() {
	apu.pulse1.stepEnvelope()
	apu.pulse2.stepEnvelope()
	apu.triangle.stepCounter()
	apu.noise.stepEnvelope()
}

func (apu *APU2303) stepSweep() {
	apu.pulse1.stepSweep()
	apu.pulse2.stepSweep()
}

func (apu *APU2303) stepLength() {
	apu.pulse1.stepLength()
	apu.pulse2.stepLength()
	apu.triangle.stepLength()
	apu.noise.stepLength()
}

func (apu *APU2303) Clock() {
	cycle1 := apu.cycle
	apu.cycle++
	cycle2 := apu.cycle
	apu.stepTimer()
	f1 := int(float64(cycle1) / frameCounterRate)
	f2 := int(float64(cycle2) / frameCounterRate)
	if f1 != f2 {
		apu.stepFrameCounter()
	}
	s1 := int(float64(cycle1) / apu.sampleRate)
	s2 := int(float64(cycle2) / apu.sampleRate)
	if s1 != s2 {
		apu.sendSample()
	}
}

func (apu *APU2303) stepTimer() {
	if apu.cycle%2 == 0 {
		apu.pulse1.stepTimer()
		apu.pulse2.stepTimer()
		apu.noise.stepTimer()
		apu.dmc.stepTimer()
	}
	apu.triangle.stepTimer()
}

func (apu *APU2303) stepFrameCounter() {
	switch apu.framePeriod {
	case 4:
		apu.frameValue = (apu.frameValue + 1) % 4
		switch apu.frameValue {
		case 0, 2:
			apu.stepEnvelope()
		case 1:
			apu.stepEnvelope()
			apu.stepSweep()
			apu.stepLength()
		case 3:
			apu.stepEnvelope()
			apu.stepSweep()
			apu.stepLength()
			apu.fireIRQ()
		}
	case 5:
		apu.frameValue = (apu.frameValue + 1) % 5
		switch apu.frameValue {
		case 0, 2:
			apu.stepEnvelope()
		case 1, 3:
			apu.stepEnvelope()
			apu.stepSweep()
			apu.stepLength()
		}
	}
}

func (apu *APU2303) fireIRQ() {
	if apu.frameIRQ {
		apu.TriggerIRQ = true
	}
}

func (apu *APU2303) sendSample() {
	output := apu.filterChain.Step(apu.output())
	select {
	case apu.channel <- output:
	default:
	}
}

func (apu *APU2303) output() float32 {
	p1 := apu.pulse1.output()
	p2 := apu.pulse2.output()
	t := apu.triangle.output()
	n := apu.noise.output()
	d := apu.dmc.output()
	pulseOut := pulseTable[p1+p2]
	tndOut := tndTable[3*t+2*n+d]
	return pulseOut + tndOut
}

func (apu *APU2303) AudioBuffer() chan float32 {
	return apu.channel
}
