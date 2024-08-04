package main

import (
	"fmt"
)

const (
	MirrorHorizontal = 0
	MirrorVertical   = 1
	MirrorSingle0    = 2
	MirrorSingle1    = 3
)

func NewMapper(bus *Bus) Mapper {
	cartridge := bus.Cartridge
	switch cartridge.Mapper {
	case 0, 2:
		return NewMapper2(cartridge)
	case 3:
		return NewMapper3(cartridge)
	case 7:
		return NewMapper7(cartridge)
	default:
		panic(fmt.Sprintf("unsupport mapper %d", cartridge.Mapper))
	}
}

type Mapper interface {
	Read(addr uint16, debug bool) uint8
	Write(addr uint16, val uint8)
}

//=====================Mapper2====================

type Mapper2 struct {
	*Cartridge
	PrgBanks int
	PrgBank1 int
	PrgBank2 int
}

func NewMapper2(cartridge *Cartridge) Mapper {
	prgBanks := len(cartridge.PRG) / 0x4000
	return &Mapper2{cartridge, prgBanks, 0, prgBanks - 1}
}

func (m *Mapper2) Read(addr uint16, _ bool) uint8 {
	switch {
	case addr < 0x2000:
		return m.CHR[addr]
	case addr >= 0xC000:
		index := m.PrgBank2*0x4000 + int(addr-0xC000)
		return m.PRG[index]
	case addr >= 0x8000:
		index := m.PrgBank1*0x4000 + int(addr-0x8000)
		return m.PRG[index]
	default:
		fmt.Printf("unsupport read addr %04X\n", addr)
	}
	return 0
}

func (m *Mapper2) Write(addr uint16, val uint8) {
	switch {
	case addr < 0x2000:
		m.CHR[addr] = val
	case addr >= 0x8000:
		m.PrgBank1 = int(val) % m.PrgBanks
	default:
		fmt.Printf("unsupport read addr %04X\n", addr)
	}
}

//======================Mapper3=====================

type Mapper3 struct {
	*Cartridge
	ChrBank  int
	PrgBank1 int
	PrgBank2 int
}

func NewMapper3(cartridge *Cartridge) Mapper {
	prgBanks := len(cartridge.PRG) / 0x4000
	return &Mapper3{cartridge, 0, 0, prgBanks - 1}
}

func (m *Mapper3) Read(addr uint16, _ bool) uint8 {
	switch {
	case addr < 0x2000:
		index := m.ChrBank*0x2000 + int(addr)
		return m.CHR[index]
	case addr >= 0xC000:
		index := m.PrgBank2*0x4000 + int(addr-0xC000)
		return m.PRG[index]
	case addr >= 0x8000:
		index := m.PrgBank1*0x4000 + int(addr-0x8000)
		return m.PRG[index]
	default:
		fmt.Printf("unsupport read addr %04X\n", addr)
	}
	return 0
}

func (m *Mapper3) Write(addr uint16, val uint8) {
	switch {
	case addr < 0x2000:
		index := m.ChrBank*0x2000 + int(addr)
		m.CHR[index] = val
	case addr >= 0x8000:
		m.ChrBank = int(val & 3)
	default:
		fmt.Printf("unsupport write addr %04X\n", addr)
	}
}

//========================Mapper7=======================

type Mapper7 struct {
	*Cartridge
	PrgBank int
}

func NewMapper7(cartridge *Cartridge) Mapper {
	return &Mapper7{cartridge, 0}
}

func (m *Mapper7) Read(addr uint16, _ bool) uint8 {
	switch {
	case addr < 0x2000:
		return m.CHR[addr]
	case addr >= 0x8000:
		index := m.PrgBank*0x8000 + int(addr-0x8000)
		return m.PRG[index]
	default:
		fmt.Printf("unsupport read addr %04X\n", addr)
	}
	return 0
}

func (m *Mapper7) Write(addr uint16, val uint8) {
	switch {
	case addr < 0x2000:
		m.CHR[addr] = val
	case addr >= 0x8000:
		m.PrgBank = int(val & 7)
		switch val & 0x10 {
		case 0x00:
			m.Cartridge.Mirror = MirrorSingle0
		case 0x10:
			m.Cartridge.Mirror = MirrorSingle1
		}
	default:
		fmt.Printf("unsupport write addr %04X\n", addr)
	}
}
