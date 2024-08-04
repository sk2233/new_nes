package main

import "github.com/hajimehoshi/ebiten/v2"

type Input struct {
	Buttons []bool
	Keys    []ebiten.Key
	Index   uint8
	Strobe  uint8
}

// keys顺序 A B Select Start Up Down Left Right
func NewInput(keys ...ebiten.Key) *Input {
	return &Input{Buttons: make([]bool, 8), Keys: keys}
}

func (c *Input) Read() uint8 {
	value := uint8(0)
	if c.Index < 8 && c.Buttons[c.Index] {
		value = 1
	}
	c.Index++
	if c.Strobe&1 == 1 {
		c.Index = 0
	}
	return value
}

func (c *Input) Write(val uint8) {
	c.Strobe = val
	if c.Strobe&1 == 1 {
		c.Index = 0
	}
}

func (c *Input) Step() {
	for i, key := range c.Keys {
		c.Buttons[i] = ebiten.IsKeyPressed(key)
	}
}
