package main

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

const (
	Width  = 256
	Height = 240
	Fps    = 60
)

const (
	ModeNormal = 0
	ModeFrame  = 1
	ModeInst   = 2
)

var (
	ModeNames = []string{"NORMAL", "FRAME", "INST"}
)

type Game struct {
	Bus        *Bus
	Option     *ebiten.DrawImageOptions
	PaletteIdx uint8
	TileMaps   []*ebiten.Image
	Mode       uint8
	//CodeLines   []uint16
	//CodeMap     map[uint16]string
	//CodeLineIdx map[uint16]int
}

func NewGame(bus *Bus) *Game {
	tileMaps := []*ebiten.Image{ebiten.NewImage(128, 128), ebiten.NewImage(128, 128)}
	// 必须有无影响读接口
	//codeMap := bus.CPU.Disassemble(0x8000, 0xFFFF)
	//codeLines := make([]uint16, 0)
	//for line := range codeMap {
	//	codeLines = append(codeLines, line)
	//}
	//sort.Slice(codeLines, func(i, j int) bool {
	//	return codeLines[i] < codeLines[j]
	//})
	//codeLineIdx := make(map[uint16]int)
	//for idx, line := range codeLines {
	//	codeLineIdx[line] = idx
	//}
	return &Game{Bus: bus, Option: &ebiten.DrawImageOptions{}, PaletteIdx: 0, TileMaps: tileMaps, Mode: ModeNormal}
}

func (g *Game) Update() error {
	g.Bus.UpdateInput()
	if inpututil.IsKeyJustPressed(ebiten.KeyR) { // 重启
		g.Bus.Reset()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) { // 调整调色盘
		g.PaletteIdx = (g.PaletteIdx + 1) % 8
		g.UpdateTileMap()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.Mode = (g.Mode + 1) % 3 // MODE 切换
	}
	// 按不同的模式执行
	switch g.Mode {
	case ModeNormal: // 正常执行
		g.Bus.FrameStep()
	case ModeFrame: // debug 逐帧允许
		if inpututil.IsKeyJustPressed(ebiten.KeyN) {
			g.Bus.PpuStep()
		}
	case ModeInst: // debug 逐指令允许
		if inpututil.IsKeyJustPressed(ebiten.KeyN) {
			g.Bus.CpuStep()
		}
	}
	return nil
}

func (g *Game) UpdateTileMap() {
	for table := 0; table < 2; table++ {
		for tileY := uint16(0); tileY < 16; tileY++ {
			for tileX := uint16(0); tileX < 16; tileX++ {
				offset := (tileY*16 + tileX) * 8 * 2 // 一共过了tileY*16 + tileX 个Tile，每个Tile 8*8 需要 8*2 byte
				for row := uint16(0); row < 8; row++ {
					// 每行8个像素由2byte组成 i 是第几个tile表(一共2个) 一共8行，所以另外一个byte需要偏移8
					tileHi := g.Bus.PPU.Read(0x1000*uint16(table)+offset+row, false)
					tileLo := g.Bus.PPU.Read(0x1000*uint16(table)+offset+row+8, false)
					for col := uint16(0); col < 8; col++ {
						// 拼接高位与地位获取索引 获取颜色
						i := (tileHi&0x01)<<1 | (tileLo & 0x01)
						tileHi >>= 1
						tileLo >>= 1 // 之所以 7-col是因为 这里是从低位开始计算的，每次位移抹除的也是低位
						g.TileMaps[table].Set(int(tileX*8+(7-col)), int(tileY*8+row), Palette[g.Bus.PPU.Palette[g.PaletteIdx*4+i]])
					}
				}
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Blue)
	// 绘制游戏画面
	g.Option.GeoM.Reset()
	g.Option.GeoM.Scale(3, 3)
	screen.DrawImage(g.Bus.Buffer(), g.Option)
	// 绘制 CPU 状态  字母宽 6 高 16
	buff := &strings.Builder{}
	buff.WriteString("STATUS:")
	WriteStatus(buff, g.Bus.CPU.N, "N")
	WriteStatus(buff, g.Bus.CPU.V, "V")
	WriteStatus(buff, g.Bus.CPU.U, "U")
	WriteStatus(buff, g.Bus.CPU.B, "B")
	WriteStatus(buff, g.Bus.CPU.D, "D")
	WriteStatus(buff, g.Bus.CPU.I, "I")
	WriteStatus(buff, g.Bus.CPU.Z, "Z")
	WriteStatus(buff, g.Bus.CPU.C, "C")
	buff.WriteString(fmt.Sprintf("\nPC: $%04X\nA: $%02X\nX: $%02X\nY: $%02X\nSP: $%04X\nMODE: %s",
		g.Bus.CPU.PC, g.Bus.CPU.A, g.Bus.CPU.X, g.Bus.CPU.Y, g.Bus.CPU.PC, ModeNames[g.Mode]))
	ebitenutil.DebugPrintAt(screen, buff.String(), Width*3, 0)
	// 绘制汇编部分
	//buff.Reset()
	//_, ok := g.CodeLineIdx[g.Bus.CPU.PC]
	//if !ok {
	//	panic(fmt.Sprintf("not find pc of %d", g.Bus.CPU.PC))
	//}
	//for i := 0; i < 29; i++ {
	//	buff.WriteString("$XXXX\n")
	//}
	code := g.Bus.CPU.DisassembleCode(29)
	ebitenutil.DebugPrintAt(screen, code, Width*3, 115)
	// 绘制调色盘
	for i := 0; i < 8; i++ {
		for j := 0; j < 4; j++ {
			x := float32(Width*3 + i*32 + 2 + j*7)
			clr := Palette[g.Bus.PPU.Palette[i*4+j]]
			vector.DrawFilledRect(screen, x, 583, 7, 7, clr, false)
		}
	}
	vector.StrokeRect(screen, Width*3+float32(g.PaletteIdx*32)+1, 582, 30, 9, 2, colornames.White, false)
	// 绘制 tileMap
	g.Option.GeoM.Reset()
	g.Option.GeoM.Translate(Width*3, 592)
	screen.DrawImage(g.TileMaps[0], g.Option) // spriteTile
	g.Option.GeoM.Reset()
	g.Option.GeoM.Translate(Width*3+128, 592)
	screen.DrawImage(g.TileMaps[1], g.Option) // bgTile
}

func WriteStatus(buff *strings.Builder, flag uint8, name string) {
	if flag == 0 {
		buff.WriteString(fmt.Sprintf(" (%s)", name))
	} else {
		buff.WriteString(fmt.Sprintf(" [%s]", name))
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return Width * 4, Height * 3
}

// WSAD FH JK 手柄控制
// R 重启 P 调整调色盘(调色盘不会实时更新需要使用P触发更新) M 运行模式切换 N debug运行时推动执行
// https://www.bilibili.com/video/BV1Uv4y1v7T9
// https://www.nesdev.org/wiki/Nesdev_Wiki

func main() {
	path := "roms/魂斗罗美版.nes"
	ebiten.SetWindowSize(Width*4, Height*3)
	ebiten.SetTPS(Fps)
	index := strings.LastIndex(path, "/") + 1
	if index < 0 {
		index = 0
	}
	ebiten.SetWindowTitle(path[index:])
	bus := NewBus(path)
	err := ebiten.RunGame(NewGame(bus))
	HandleErr(err)
}
