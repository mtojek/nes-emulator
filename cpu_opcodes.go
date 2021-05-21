package main

var opcodes = map[uint8]instruction{}

type instruction struct {
	name   string
	cycles uint8

	addressingMode addressingModeFunc
	operate        operateFunc
}

type addressingModeFunc func() uint8

type operateFunc func() uint8