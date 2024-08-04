package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Bus struct {
	CPU       *CPU
	PPU       *PPU
	Cartridge *Cartridge
	Input1    *Input
	Input2    *Input
	Mapper    Mapper
	RAM       []byte
	CurrBuff  *ebiten.Image // 用来实现逐帧渲染的 使用 PPU.Frame 也可以
}

func NewBus(path string) *Bus {
	// 暂时 2p没有输入，只有 1p 有输入
	bus := &Bus{Cartridge: LoadCartridge(path), RAM: make([]byte, 2*1024), Input2: NewInput(),
		Input1: NewInput(ebiten.KeyK, ebiten.KeyJ, ebiten.KeyF, ebiten.KeyH, ebiten.KeyW, ebiten.KeyS, ebiten.KeyA, ebiten.KeyD)}
	bus.Mapper = NewMapper(bus)
	bus.CPU = NewCPU(bus)
	bus.PPU = NewPPU(bus)
	return bus
}

func (c *Bus) Reset() {
	c.CPU.Reset()
}

// 执行 cpu的一个指令
func (c *Bus) CpuStep() int {
	cpuCycles := c.CPU.Step()
	ppuCycles := cpuCycles * 3
	for i := 0; i < ppuCycles; i++ {
		c.PPU.Step()
	}
	return cpuCycles
}

// 绘制 ppu的一帧画面
func (c *Bus) PpuStep() {
	for c.CurrBuff == c.PPU.FrontBuff {
		c.CpuStep()
	}
	c.CurrBuff = c.PPU.FrontBuff
}

// 60帧每秒每帧的运行量
func (c *Bus) FrameStep() {
	cycles := CPUFreq / Fps
	for cycles > 0 {
		cycles -= c.CpuStep()
	}
}

func (c *Bus) Buffer() *ebiten.Image {
	return c.PPU.FrontBuff
}

func (c *Bus) UpdateInput() {
	c.Input1.Step()
	c.Input2.Step()
}
