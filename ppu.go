package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	Palette      [64]color.RGBA
	MirrorLookup = [4][4]uint16{
		{0, 0, 1, 1},
		{0, 1, 0, 1},
		{0, 0, 0, 0},
		{1, 1, 1, 1},
	}
)

// 获取各种镜像模式下nametable的地址
func MirrorAddr(mode uint8, addr uint16) uint16 {
	addr = (addr - 0x2000) % 0x1000
	table := addr / 0x0400
	offset := addr % 0x0400
	return 0x2000 + MirrorLookup[mode][table]*0x0400 + offset
}

// 使用调色盘，颜色是固定的
func InitPalette() {
	colors := []uint32{
		0x666666, 0x002A88, 0x1412A7, 0x3B00A4, 0x5C007E, 0x6E0040, 0x6C0600, 0x561D00,
		0x333500, 0x0B4800, 0x005200, 0x004F08, 0x00404D, 0x000000, 0x000000, 0x000000,
		0xADADAD, 0x155FD9, 0x4240FF, 0x7527FE, 0xA01ACC, 0xB71E7B, 0xB53120, 0x994E00,
		0x6B6D00, 0x388700, 0x0C9300, 0x008F32, 0x007C8D, 0x000000, 0x000000, 0x000000,
		0xFFFEFF, 0x64B0FF, 0x9290FF, 0xC676FF, 0xF36AFF, 0xFE6ECC, 0xFE8170, 0xEA9E22,
		0xBCBE00, 0x88D800, 0x5CE430, 0x45E082, 0x48CDDE, 0x4F4F4F, 0x000000, 0x000000,
		0xFFFEFF, 0xC0DFFF, 0xD3D2FF, 0xE8C8FF, 0xFBC2FF, 0xFEC4EA, 0xFECCC5, 0xF7D8A5,
		0xE4E594, 0xCFEF96, 0xBDF4AB, 0xB3F3CC, 0xB5EBF2, 0xB8B8B8, 0x000000, 0x000000,
	}
	for i, c := range colors {
		r := uint8(c >> 16)
		g := uint8(c >> 8)
		b := uint8(c)
		Palette[i] = color.RGBA{R: r, G: g, B: b, A: 0xFF}
	}
}

type PPU struct {
	Bus      *Bus
	Cycle    int    // 0-340
	ScanLine int    // 0-261, 0-239 渲染可见内容, 240 后置处理, 241-260 vblank阶段, 261 前置处理
	Frame    uint64 // 绘制的多少帧了
	// 一些静态变量
	Palette   [32]uint8       // 调色盘
	NameTable [2 * 1024]uint8 // tile显示的样子
	OamData   [256]uint8      // 精灵属性数据
	FrontBuff *ebiten.Image   // 绘图双缓冲
	BackBuff  *ebiten.Image
	// PPU 寄存器
	CurrVRam  uint16 // 当前的 vram 地址 15 Bit   用来读取要显示内容信息的地址
	TempVRam  uint16 // 临时的 vram 地址 15 Bit   用来读取要显示内容信息的地址
	FineX     uint8  // fine_x 3 Bit
	WriteFlag uint8  // 地址等信息 16bit 需要写两次的标记 1 Bit
	FrameFlag uint8  // 奇数帧还是偶数帧标记 1 Bit
	Register  uint8  // PPU 寄存器
	// nmi 相关变量
	NmiOccur  bool
	NmiOutput bool
	NmiPre    bool
	NmiDelay  uint8
	// 背景tile相关变量
	NameTableData uint8
	AttrTableData uint8
	LowTileData   uint8
	HighTileData  uint8
	TileData      uint64
	// 前景 sprite相关变量
	SpriteCount      int
	SpritePatterns   [8]uint32
	SpritePos        [8]uint8
	SpritePriorities [8]uint8
	SpriteIdxes      [8]uint8
	// $2000 PPUCTRL    8 bit  还有 1bit控制 NmiOutput 不在这里
	FlagNameTable   uint8 // name起始地址映射 0: $2000; 1: $2400; 2: $2800; 3: $2C00    2 bit
	FlagIncrMode    uint8 // 写入后自动偏移的数量方便连续偏移 0: add 1; 1: add 32          1 bit
	FlagSpriteTable uint8 // 两张精灵tile用那个 0: $0000; 1: $1000;                     1 bit
	FlagBgTable     uint8 // 两张背景tile用那个 0: $0000; 1: $1000                      1 bit
	FlagSpriteSize  uint8 // 精灵大小是 8*8模式还是 8*16模式 0: 8x8; 1: 8x16              1 bit
	FlagMasterSlave uint8 // 读或写 暂时没有用                                           1 bit
	// $2001 PPUMASK 8bit
	FlagGray           uint8 // 是否为灰度模式 0: 不是; 1: 是           1 bit
	FlagShowLeftBg     uint8 // 是否显示LeftBg 0: 不显示; 1: 显示       1 bit
	FlagShowLeftSprite uint8 // 是否显示LeftSprite 0: 不显示; 1: 显示   1 bit
	FlagShowBg         uint8 // 是否显示Bg 0: 不显示; 1: 显示           1 bit
	FlagShowSprite     uint8 // 是否显示Sprite 0: 不显示; 1: 显示       1 bit
	FlagRedTint        uint8 // 是否进行红色染色 0: 不染色; 1: 染色      1 bit
	FlagGreenTint      uint8 // 是否进行绿色染色 0: 不染色; 1: 染色      1 bit
	FlagBlueTint       uint8 // 是否进行蓝色染色 0: 不染色; 1: 染色      1 bit
	// $2002 PPUSTATUS
	FlagSpriteHitZero  uint8 // 是否第一次前景与背景不透明部分重叠
	FlagSpriteOverflow uint8 // 同一行扫描线 sprite 是否超过了 8 个
	// $2003 OAMADDR
	OamAddr uint8 // Oam数据地址
	// $2007 PPUDATA
	BuffData uint8 // ppuData存在读延迟，需要记录上次读取的结果用于这次返回
}

func NewPPU(bus *Bus) *PPU {
	ppu := &PPU{Bus: bus}
	ppu.FrontBuff = ebiten.NewImage(256, 240)
	ppu.BackBuff = ebiten.NewImage(256, 240)
	InitPalette()
	ppu.Reset()
	return ppu
}

func (p *PPU) Read(addr uint16, debug bool) uint8 {
	addr = addr % 0x4000
	switch {
	case addr < 0x2000:
		return p.Bus.Mapper.Read(addr, debug)
	case addr < 0x3F00:
		mode := p.Bus.Cartridge.Mirror
		return p.Bus.PPU.NameTable[MirrorAddr(mode, addr)%2048]
	case addr < 0x4000:
		return p.Bus.PPU.ReadPalette(addr % 32)
	default:
		fmt.Printf("unsupport read addr %04X\n", addr)
	}
	return 0
}

func (p *PPU) Write(addr uint16, val uint8) {
	addr = addr % 0x4000
	switch {
	case addr < 0x2000:
		p.Bus.Mapper.Write(addr, val)
	case addr < 0x3F00:
		mode := p.Bus.Cartridge.Mirror
		p.Bus.PPU.NameTable[MirrorAddr(mode, addr)%2048] = val
	case addr < 0x4000:
		p.Bus.PPU.WritePalette(addr%32, val)
	default:
		fmt.Printf("unsupport write addr %04X\n", addr)
	}
}

func (p *PPU) Reset() {
	p.Cycle = 340
	p.ScanLine = 240
	p.Frame = 0
	p.WriteControl(0)
	p.WriteMask(0)
	p.WriteOAMAddr(0)
}

// 获取调色盘的索引值
func (p *PPU) ReadPalette(addr uint16) uint8 {
	if addr >= 16 && addr%4 == 0 {
		addr -= 16
	}
	return p.Palette[addr]
}

// 初始化调色盘
func (p *PPU) WritePalette(addr uint16, val uint8) {
	if addr >= 16 && addr%4 == 0 {
		addr -= 16
	}
	p.Palette[addr] = val
}

func (p *PPU) ReadR(addr uint16) uint8 {
	switch addr {
	case 0x2002:
		return p.ReadStatus()
	case 0x2004:
		return p.ReadOAMData() // 一般不会一个个读取数据，直接采用 DMA 数据拷贝
	case 0x2007:
		return p.ReadData()
	}
	return 0
}

func (p *PPU) WriteR(addr uint16, val uint8) {
	p.Register = val
	switch addr {
	case 0x2000:
		p.WriteControl(val)
	case 0x2001:
		p.WriteMask(val)
	case 0x2003:
		p.WriteOAMAddr(val)
	case 0x2004:
		p.WriteOAMData(val) // 一般不会一个个读取数据，直接采用 DMA 数据拷贝 与上面的配合使用
	case 0x2005:
		p.WriteScroll(val)
	case 0x2006:
		p.WriteAddr(val)
	case 0x2007:
		p.WriteData(val) // 一般 ppu 的数据读取，配合上面的方法使用
	case 0x4014:
		p.WriteDMA(val) // 启动 DMA
	}
}

// $2000: PPUCTRL
func (p *PPU) WriteControl(value uint8) {
	p.FlagNameTable = (value >> 0) & 3
	p.FlagIncrMode = (value >> 2) & 1
	p.FlagSpriteTable = (value >> 3) & 1
	p.FlagBgTable = (value >> 4) & 1
	p.FlagSpriteSize = (value >> 5) & 1
	p.FlagMasterSlave = (value >> 6) & 1
	p.NmiOutput = (value>>7)&1 == 1
	p.NmiChange() // 还要注意对 VRam 的影响
	// TempVRam: ....BA.. ........ = d: ......BA
	p.TempVRam = (p.TempVRam & 0xF3FF) | ((uint16(value) & 0x03) << 10)
}

// $2001: PPUMASK
func (p *PPU) WriteMask(value uint8) {
	p.FlagGray = (value >> 0) & 1
	p.FlagShowLeftBg = (value >> 1) & 1
	p.FlagShowLeftSprite = (value >> 2) & 1
	p.FlagShowBg = (value >> 3) & 1
	p.FlagShowSprite = (value >> 4) & 1
	p.FlagRedTint = (value >> 5) & 1
	p.FlagGreenTint = (value >> 6) & 1
	p.FlagBlueTint = (value >> 7) & 1
}

// $2002: PPUSTATUS
func (p *PPU) ReadStatus() uint8 {
	result := p.Register & 0x1F
	result |= p.FlagSpriteOverflow << 5
	result |= p.FlagSpriteHitZero << 6
	if p.NmiOccur {
		result |= 1 << 7
	}
	p.NmiOccur = false
	p.NmiChange()
	p.WriteFlag = 0
	return result
}

// $2003: OAMADDR
func (p *PPU) WriteOAMAddr(value uint8) {
	p.OamAddr = value
}

// $2004: OAMDATA (read)
func (p *PPU) ReadOAMData() uint8 {
	data := p.OamData[p.OamAddr]
	if (p.OamAddr & 0x03) == 0x02 {
		data = data & 0xE3
	}
	return data
}

// $2004: OAMDATA (write)
func (p *PPU) WriteOAMData(value uint8) {
	p.OamData[p.OamAddr] = value
	p.OamAddr++
}

// $2005: PPUSCROLL
func (p *PPU) WriteScroll(value uint8) {
	// 分两次写入滚动属性 使用 WriteFlag 区分
	if p.WriteFlag == 0 {
		// TempVRam: ........ ...HGFED = d: HGFED...
		// FineX:               CBA = d: .....CBA
		// WriteFlag:                   = 1
		p.TempVRam = (p.TempVRam & 0xFFE0) | (uint16(value) >> 3)
		p.FineX = value & 0x07
		p.WriteFlag = 1
	} else {
		// TempVRam: .CBA..HG FED..... = d: HGFEDCBA
		// WriteFlag:                   = 0
		p.TempVRam = (p.TempVRam & 0x8FFF) | ((uint16(value) & 0x07) << 12)
		p.TempVRam = (p.TempVRam & 0xFC1F) | ((uint16(value) & 0xF8) << 2)
		p.WriteFlag = 0
	}
}

// $2006: PPUADDR
func (p *PPU) WriteAddr(value uint8) {
	// 分两次写入 Addr 并切换 VRam
	if p.WriteFlag == 0 {
		// TempVRam: ..FEDCBA ........ = d: ..FEDCBA
		// TempVRam: .X...... ........ = 0
		// WriteFlag:                   = 1
		p.TempVRam = (p.TempVRam & 0x80FF) | ((uint16(value) & 0x3F) << 8)
		p.WriteFlag = 1
	} else {
		// TempVRam: ........ HGFEDCBA = d: HGFEDCBA
		// CurrVRam                    = TempVRam
		// WriteFlag:                   = 0
		p.TempVRam = (p.TempVRam & 0xFF00) | uint16(value)
		p.CurrVRam = p.TempVRam
		p.WriteFlag = 0
	}
}

// $2007: PPUDATA (read)
func (p *PPU) ReadData() uint8 {
	value := p.Read(p.CurrVRam, false)
	// 需要模拟延迟读
	if p.CurrVRam%0x4000 < 0x3F00 {
		temp := p.BuffData // 只有这部分需要延迟读
		p.BuffData = value
		value = temp
	} else {
		p.BuffData = p.Read(p.CurrVRam-0x1000, false)
	}
	// 读取完自动根据模式调整地址，方便连续读取
	if p.FlagIncrMode == 0 {
		p.CurrVRam += 1
	} else {
		p.CurrVRam += 32
	}
	return value
}

// $2007: PPUDATA (write)
func (p *PPU) WriteData(value uint8) {
	p.Write(p.CurrVRam, value)
	// 写入也要移动地址，方便连续写入
	if p.FlagIncrMode == 0 {
		p.CurrVRam += 1
	} else {
		p.CurrVRam += 32
	}
}

// $4014: OAMDMA
func (p *PPU) WriteDMA(value uint8) {
	cpu := p.Bus.CPU
	addr := uint16(value) << 8
	for i := 0; i < 256; i++ { // 使用DMA转移Oam数据，并模拟占据对应的 cpu周期
		p.OamData[p.OamAddr] = cpu.Read(addr, false)
		p.OamAddr++
		addr++
	}
	cpu.LastCycles += 513
	if cpu.Cycles%2 == 1 {
		cpu.LastCycles++
	}
}

// 扫描线 X 轴移动
func (p *PPU) IncrementX() {
	if p.CurrVRam&0x001F == 31 {
		// x 轴扫描完毕了，归零 x 部分bit 并切换 name_table
		p.CurrVRam &= 0xFFE0
		p.CurrVRam ^= 0x0400
	} else {
		// 简单累加 x x正好处于低位，可以直接累加整体
		p.CurrVRam++
	}
}

func (p *PPU) IncrementY() {
	// fineY部分 进行累加
	if p.CurrVRam&0x7000 != 0x7000 {
		p.CurrVRam += 0x1000 // 没有累加到7（一个tile）继续累加
	} else { // 累加到 7 了清零 fineY 部分
		p.CurrVRam &= 0x8FFF
		y := (p.CurrVRam & 0x03E0) >> 5 // 获取到 coarse_y
		if y == 29 {                    // y方向的 tile 处理完了归零 y 并切换 y 方向的name_table
			y = 0
			p.CurrVRam ^= 0x0800
		} else if y == 31 {
			y = 0 // 前面有 29 的判断，这里几乎不会走到
		} else { // 简单累加 y 方向的 tile
			y++
		}
		// 把 coarse_y 放回到 vram
		p.CurrVRam = (p.CurrVRam & 0xFC1F) | (y << 5)
	}
}

func (p *PPU) CopyX() {
	// 使用 TempVRam 复位 x 坐标
	// CurrVRam: .....F.. ...EDCBA = TempVRam: .....F.. ...EDCBA
	p.CurrVRam = (p.CurrVRam & 0xFBE0) | (p.TempVRam & 0x041F)
}

func (p *PPU) CopyY() {
	// 使用 TempVRam 复位 y 坐标
	// CurrVRam: .IHGF.ED CBA..... = TempVRam: .IHGF.ED CBA.....
	p.CurrVRam = (p.CurrVRam & 0x841F) | (p.TempVRam & 0x7BE0)
}

func (p *PPU) NmiChange() {
	nmi := p.NmiOutput && p.NmiOccur
	if nmi && !p.NmiPre { // 只有第一次需要模拟 nmi 的延迟
		p.NmiDelay = 15
	}
	p.NmiPre = nmi
}

func (p *PPU) SetVBlank() {
	p.FrontBuff, p.BackBuff = p.BackBuff, p.FrontBuff
	p.NmiOccur = true
	p.NmiChange()
}

func (p *PPU) ClearVBlank() {
	p.NmiOccur = false
	p.NmiChange()
}

func (p *PPU) FetchNameTableData() {
	// 更具地址获取要绘制的 name_table 数据
	v := p.CurrVRam
	address := 0x2000 | (v & 0x0FFF)
	p.NameTableData = p.Read(address, false)
}

func (p *PPU) FetchAttrTableData() {
	// 获取要绘制 tile 的属性信息，注意4个块共用一个属性
	v := p.CurrVRam
	address := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
	shift := ((v >> 4) & 4) | (v & 2)
	p.AttrTableData = ((p.Read(address, false) >> shift) & 3) << 2
}

func (p *PPU) FetchLowTileData() {
	// 获取 tile 低位索引
	fineY := (p.CurrVRam >> 12) & 7
	table := p.FlagBgTable
	tile := p.NameTableData
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	p.LowTileData = p.Read(address, false)
}

func (p *PPU) FetchHighTileData() {
	// 获取 tile 高位索引 必须与低位对应位结合形成 2bit 的色盘索引
	fineY := (p.CurrVRam >> 12) & 7
	table := p.FlagBgTable
	tile := p.NameTableData
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	p.HighTileData = p.Read(address+8, false)
}

func (p *PPU) StoreTileData() {
	var data uint32
	for i := 0; i < 8; i++ {
		// 2位色盘数据，2位色盘内颜色索引数据 4bit 8个 32bit
		a := p.AttrTableData
		p1 := (p.LowTileData & 0x80) >> 7
		p2 := (p.HighTileData & 0x80) >> 6
		p.LowTileData <<= 1
		p.HighTileData <<= 1
		data <<= 4
		data |= uint32(a | p1 | p2)
	} // 累加上去方便后面移位绘制
	p.TileData |= uint64(data)
}

func (p *PPU) FetchTileData() uint32 {
	return uint32(p.TileData >> 32) // 当前使用的是 高 32位
}

func (p *PPU) BgPixel() uint8 {
	if p.FlagShowBg == 0 {
		return 0
	} // 获取对应的色盘与色盘内的颜色索引 2位调色盘索引（只看背景的），2位调色盘内部颜色索引 正好组成整个颜色索引，所有调色盘顺序排放的
	data := p.FetchTileData() >> ((7 - p.FineX) * 4)
	return uint8(data & 0x0F)
}

func (p *PPU) SpritePixel() (uint8, uint8) {
	if p.FlagShowSprite == 0 {
		return 0, 0
	}
	for i := 0; i < p.SpriteCount; i++ {
		offset := (p.Cycle - 1) - int(p.SpritePos[i])
		if offset < 0 || offset > 7 {
			continue
		}
		offset = 7 - offset
		// sprite 的色盘与色盘内索引不是连续存放的，需要做一些移位操作来拼出色盘下标
		color0 := uint8((p.SpritePatterns[i] >> uint8(offset*4)) & 0x0F)
		if color0%4 == 0 { // 前景透明色忽略
			continue
		} // 返回是第几个精灵与其颜色
		return uint8(i), color0
	}
	return 0, 0
}

func (p *PPU) RenderPixel() {
	x := p.Cycle - 1
	y := p.ScanLine
	bg := p.BgPixel()            // 先获取对应像素的背景色
	i, sprite := p.SpritePixel() // 再获取前景色
	if x < 8 && p.FlagShowLeftBg == 0 {
		bg = 0
	}
	if x < 8 && p.FlagShowLeftSprite == 0 {
		sprite = 0
	}
	// 更具其是否为透明色进行最终颜色的仲裁
	b := bg%4 != 0
	s := sprite%4 != 0
	color0 := uint8(0)
	if !b && s { // 使用前景调色盘
		color0 = sprite | 0x10
	} else if b && !s { // 使用背景调色盘
		color0 = bg
	} else if b && s { // 都不是透明色
		if p.SpriteIdxes[i] == 0 && x < 255 { // 检查 HitZero
			p.FlagSpriteHitZero = 1
		}
		if p.SpritePriorities[i] == 0 { // 一般优先显示前景色
			color0 = sprite | 0x10
		} else {
			color0 = bg
		}
	}
	c := Palette[p.ReadPalette(uint16(color0))%64]
	p.BackBuff.Set(x, y, c)
}

func (p *PPU) FetchSpritePattern(i, row int) uint32 {
	tile := p.OamData[i*4+1]
	attr := p.OamData[i*4+2]
	addr := uint16(0)
	if p.FlagSpriteSize == 0 { // 8*8模式
		if attr&0x80 == 0x80 { // 判断y反转属性
			row = 7 - row
		}
		table := p.FlagSpriteTable
		// 各种数据占有不同位的偏移组合起来就是最终地址
		addr = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	} else { // 8*16模式
		if attr&0x80 == 0x80 { // 判断y反转属性
			row = 15 - row
		}
		table := tile & 1 // 8*16 模式是正好占据上下相连的 2 个 tile
		tile &= 0xFE
		if row > 7 {
			tile++
			row -= 8
		}
		addr = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	}
	// 只获取 attr 的色盘索引部分，且流出色盘内索引的 2bit
	a := (attr & 3) << 2
	lowTileData := p.Read(addr, false)
	highTileData := p.Read(addr+8, false)
	data := uint32(0)
	for j := 0; j < 8; j++ {
		// 读取 x方向的颜色信息一次存储 8 个
		var p1, p2 uint8
		if attr&0x40 == 0x40 { // 判断是否存在 x 反装属性
			p1 = (lowTileData & 1) << 0
			p2 = (highTileData & 1) << 1
			lowTileData >>= 1
			highTileData >>= 1
		} else {
			p1 = (lowTileData & 0x80) >> 7
			p2 = (highTileData & 0x80) >> 6
			lowTileData <<= 1
			highTileData <<= 1
		}
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	return data
}

func (p *PPU) EvaluateSprite() {
	var h int
	if p.FlagSpriteSize == 0 { // 获取对应 sprite 模式下的高度
		h = 8
	} else {
		h = 16
	}
	count := 0
	for i := 0; i < 64; i++ {
		// 最多 64 个精灵，遍历处理获取其对应的 x y 与 attr
		y := p.OamData[i*4]
		a := p.OamData[i*4+2]
		x := p.OamData[i*4+3]
		row := p.ScanLine - int(y)
		if row < 0 || row >= h { // y方向必须在范围内否则该精灵无需处理
			continue
		}
		if count < 8 { // 一条扫描线最多 8 个精灵
			p.SpritePatterns[count] = p.FetchSpritePattern(i, row)
			p.SpritePos[count] = x
			p.SpritePriorities[count] = (a >> 5) & 1
			p.SpriteIdxes[count] = uint8(i)
		}
		count++
	}
	if count > 8 { // 更新精灵溢出标记
		count = 8
		p.FlagSpriteOverflow = 1
	}
	p.SpriteCount = count
}

func (p *PPU) Tick() {
	if p.NmiDelay > 0 { // 处理中断，ppu通过中断与 cpu传输数据
		p.NmiDelay--
		if p.NmiDelay == 0 && p.NmiOutput && p.NmiOccur {
			p.Bus.CPU.TriggerNMI()
		}
	}
	// 判断是否渲染完了一帧
	if p.FlagShowBg != 0 || p.FlagShowSprite != 0 {
		if p.FrameFlag == 1 && p.ScanLine == 261 && p.Cycle == 339 {
			p.Cycle = 0
			p.ScanLine = 0
			p.Frame++
			p.FrameFlag ^= 1
			return
		}
	}
	// 逐像素推进渲染
	p.Cycle++
	if p.Cycle > 340 {
		p.Cycle = 0
		p.ScanLine++
		if p.ScanLine > 261 {
			p.ScanLine = 0
			p.Frame++
			p.FrameFlag ^= 1
		}
	}
}

func (p *PPU) Step() {
	p.Tick()
	// 大体是否需要渲染
	needRender := p.FlagShowBg != 0 || p.FlagShowSprite != 0
	preLine := p.ScanLine == 261 // 主要是为第 0 行渲染做准备
	visibleLine := p.ScanLine < 240
	// 改行是否需要渲染
	renderLine := preLine || visibleLine
	preFetchCycle := p.Cycle >= 321 && p.Cycle <= 336
	visibleCycle := p.Cycle >= 1 && p.Cycle <= 256
	fetchCycle := preFetchCycle || visibleCycle
	if needRender {
		if visibleLine && visibleCycle {
			p.RenderPixel() // 行列都是可见的，设置对应的像素信息
		}
		// 在需要准备渲染的 y 与需要拿数据的 x 更具周期预先拿对应的数据
		if renderLine && fetchCycle {
			p.TileData <<= 4
			switch p.Cycle % 8 {
			case 1:
				p.FetchNameTableData()
			case 3:
				p.FetchAttrTableData()
			case 5:
				p.FetchLowTileData()
			case 7:
				p.FetchHighTileData()
			case 0:
				p.StoreTileData()
			}
		} // preLine 要为下一帧渲染做准备了 先恢复 y
		if preLine && p.Cycle >= 280 && p.Cycle <= 304 {
			p.CopyY()
		}
		if renderLine {
			// 没渲染完 8 个增加 x 偏移
			if fetchCycle && p.Cycle%8 == 0 {
				p.IncrementX()
			} // 没渲染一行增加 y 偏移
			if p.Cycle == 256 {
				p.IncrementY()
			} // 渲染完一行恢复 x
			if p.Cycle == 257 {
				p.CopyX()
			}
		}
	}
	// 准备下一个扫描线的数据
	if needRender {
		if p.Cycle == 257 {
			if visibleLine {
				p.EvaluateSprite()
			} else {
				p.SpriteCount = 0
			}
		}
	}
	// 进入非渲染数据准备阶段
	if p.ScanLine == 241 && p.Cycle == 1 {
		p.SetVBlank()
	} // 退出非渲染数据准备阶段 进行下一帧的渲染
	if preLine && p.Cycle == 1 {
		p.ClearVBlank()
		p.FlagSpriteHitZero = 0
		p.FlagSpriteOverflow = 0
	}
}
