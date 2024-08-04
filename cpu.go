package main

import (
	"fmt"
	"strings"
)

// CPU 主频
const CPUFreq = 1789773

// 中断类型
const (
	IntNone = 1
	IntNMI  = 2
	IntIRQ  = 3
)

// 寻址模式
const (
	AddrAbsolute        = 1
	AddrAbsoluteX       = 2
	AddrAbsoluteY       = 3
	AddrAccumulator     = 4
	AddrImmediate       = 5
	AddrImplied         = 6
	AddrIndexedIndirect = 7
	AddrIndirect        = 8
	AddrIndirectIndexed = 9
	AddrRelative        = 10
	AddrZeroPage        = 11
	AddrZeroPageX       = 12
	AddrZeroPageY       = 13
)

// 指令寻址模式表
var InstAddrModes = [256]uint8{
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	1, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	6, 7, 6, 7, 11, 11, 11, 11, 6, 5, 4, 5, 8, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 13, 13, 6, 3, 6, 3, 2, 2, 3, 3,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 13, 13, 6, 3, 6, 3, 2, 2, 3, 3,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
	5, 7, 5, 7, 11, 11, 11, 11, 6, 5, 6, 5, 1, 1, 1, 1,
	10, 9, 6, 9, 12, 12, 12, 12, 6, 3, 6, 3, 2, 2, 2, 2,
}

// 每条指令的大小，PC 需要移动多少
var InstSizes = [256]uint8{
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	3, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	1, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 0, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 0, 3, 0, 0,
	2, 2, 2, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
}

// 每条指令执行消耗cpu 周期数目
var InstCycles = [256]uint8{
	7, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 6, 2, 6, 4, 4, 4, 4, 2, 5, 2, 5, 5, 5, 5, 5,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 5, 2, 5, 4, 4, 4, 4, 2, 4, 2, 4, 4, 4, 4, 4,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
}

// 指令发生跨页读时额外需要的周期
var InstPageCycles = [256]uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0,
}

// 每条指令的名称
var InstNames = [256]string{
	"BRK", "ORA", "NOP", "NOP", "NOP", "ORA", "ASL", "NOP", "PHP", "ORA", "ASL", "NOP", "NOP", "ORA", "ASL", "NOP",
	"BPL", "ORA", "NOP", "NOP", "NOP", "ORA", "ASL", "NOP", "CLC", "ORA", "NOP", "NOP", "NOP", "ORA", "ASL", "NOP",
	"JSR", "AND", "NOP", "NOP", "BIT", "AND", "ROL", "NOP", "PLP", "AND", "ROL", "NOP", "BIT", "AND", "ROL", "NOP",
	"BMI", "AND", "NOP", "NOP", "NOP", "AND", "ROL", "NOP", "SEC", "AND", "NOP", "NOP", "NOP", "AND", "ROL", "NOP",
	"RTI", "EOR", "NOP", "NOP", "NOP", "EOR", "LSR", "NOP", "PHA", "EOR", "LSR", "NOP", "JMP", "EOR", "LSR", "NOP",
	"BVC", "EOR", "NOP", "NOP", "NOP", "EOR", "LSR", "NOP", "CLI", "EOR", "NOP", "NOP", "NOP", "EOR", "LSR", "NOP",
	"RTS", "ADC", "NOP", "NOP", "NOP", "ADC", "ROR", "NOP", "PLA", "ADC", "ROR", "NOP", "JMP", "ADC", "ROR", "NOP",
	"BVS", "ADC", "NOP", "NOP", "NOP", "ADC", "ROR", "NOP", "SEI", "ADC", "NOP", "NOP", "NOP", "ADC", "ROR", "NOP",
	"NOP", "STA", "NOP", "NOP", "STY", "STA", "STX", "NOP", "DEY", "NOP", "TXA", "NOP", "STY", "STA", "STX", "NOP",
	"BCC", "STA", "NOP", "NOP", "STY", "STA", "STX", "NOP", "TYA", "STA", "TXS", "NOP", "NOP", "STA", "NOP", "NOP",
	"LDY", "LDA", "LDX", "NOP", "LDY", "LDA", "LDX", "NOP", "TAY", "LDA", "TAX", "NOP", "LDY", "LDA", "LDX", "NOP",
	"BCS", "LDA", "NOP", "NOP", "LDY", "LDA", "LDX", "NOP", "CLV", "LDA", "TSX", "NOP", "LDY", "LDA", "LDX", "NOP",
	"CPY", "CMP", "NOP", "NOP", "CPY", "CMP", "DEC", "NOP", "INY", "CMP", "DEX", "NOP", "CPY", "CMP", "DEC", "NOP",
	"BNE", "CMP", "NOP", "NOP", "NOP", "CMP", "DEC", "NOP", "CLD", "CMP", "NOP", "NOP", "NOP", "CMP", "DEC", "NOP",
	"CPX", "SBC", "NOP", "NOP", "CPX", "SBC", "INC", "NOP", "INX", "SBC", "NOP", "SBC", "CPX", "SBC", "INC", "NOP",
	"BEQ", "SBC", "NOP", "NOP", "NOP", "SBC", "INC", "NOP", "SED", "SBC", "NOP", "NOP", "NOP", "SBC", "INC", "NOP",
}

type CPU struct {
	Bus    *Bus
	Cycles uint64 // 周期数
	PC     uint16 // 程序计数器
	SP     uint8  // 栈指针
	A      uint8  // 累加寄存器
	X      uint8  // 通用寄存器 FineX
	Y      uint8  // 通用寄存器 y
	// 各种标记，使用 uint8存储更方便使用可以直接累加
	C          uint8                // 进位标记
	Z          uint8                // 结果为 0 标记
	I          uint8                // 屏蔽中断标记
	D          uint8                // 10进制标记 无用
	B          uint8                // Break 标记 无用
	U          uint8                // 空标记 没有含义
	V          uint8                // 计算溢出标记
	N          uint8                // 结果为负数标记
	IntType    uint8                // 中断类型
	LastCycles int                  // 剩余还要执行的周期
	InstTable  [256]func(*StepInfo) // 指令表
}

func NewCPU(bus *Bus) *CPU {
	cpu := CPU{Bus: bus}
	cpu.CreateInstTable()
	cpu.Reset()
	return &cpu
}

func (c *CPU) Read(addr uint16, debug bool) uint8 {
	if debug { // debug模式下只允许无损读
		switch {
		case addr < 0x2000:
			return c.Bus.RAM[addr%0x0800]
		case addr >= 0x6000:
			return c.Bus.Mapper.Read(addr, debug)
		default:
			return 0
		}
	}

	switch {
	case addr < 0x2000:
		return c.Bus.RAM[addr%0x0800]
	case addr < 0x4000:
		return c.Bus.PPU.ReadR(0x2000 + addr%8)
	case addr == 0x4014:
		return c.Bus.PPU.ReadR(addr)
	case addr == 0x4016:
		return c.Bus.Input1.Read()
	case addr == 0x4017:
		return c.Bus.Input2.Read()
	case addr >= 0x6000:
		return c.Bus.Mapper.Read(addr, debug)
	default:
		//fmt.Printf("unsupport read addr %04X\n", addr)
	}
	return 0
}

func (c *CPU) Write(addr uint16, val uint8) {
	switch {
	case addr < 0x2000:
		c.Bus.RAM[addr%0x0800] = val
	case addr < 0x4000:
		c.Bus.PPU.WriteR(0x2000+addr%8, val)
	case addr == 0x4014:
		c.Bus.PPU.WriteR(addr, val)
	case addr == 0x4016:
		c.Bus.Input1.Write(val)
		c.Bus.Input2.Write(val)
	case addr >= 0x6000:
		c.Bus.Mapper.Write(addr, val)
	default:
		//fmt.Printf("unsupport write addr %04X\n", addr)
	}
}

// 初始化指令表
func (c *CPU) CreateInstTable() {
	c.InstTable = [256]func(*StepInfo){
		c.Brk, c.Ora, c.Nop, c.Nop, c.Nop, c.Ora, c.Asl, c.Nop, c.Php, c.Ora, c.Asl, c.Nop, c.Nop, c.Ora, c.Asl, c.Nop,
		c.Bpl, c.Ora, c.Nop, c.Nop, c.Nop, c.Ora, c.Asl, c.Nop, c.Clc, c.Ora, c.Nop, c.Nop, c.Nop, c.Ora, c.Asl, c.Nop,
		c.Jsr, c.And, c.Nop, c.Nop, c.Bit, c.And, c.Rol, c.Nop, c.Plp, c.And, c.Rol, c.Nop, c.Bit, c.And, c.Rol, c.Nop,
		c.Bmi, c.And, c.Nop, c.Nop, c.Nop, c.And, c.Rol, c.Nop, c.Sec, c.And, c.Nop, c.Nop, c.Nop, c.And, c.Rol, c.Nop,
		c.Rti, c.Eor, c.Nop, c.Nop, c.Nop, c.Eor, c.Lsr, c.Nop, c.Pha, c.Eor, c.Lsr, c.Nop, c.Jmp, c.Eor, c.Lsr, c.Nop,
		c.Bvc, c.Eor, c.Nop, c.Nop, c.Nop, c.Eor, c.Lsr, c.Nop, c.Cli, c.Eor, c.Nop, c.Nop, c.Nop, c.Eor, c.Lsr, c.Nop,
		c.Rts, c.Adc, c.Nop, c.Nop, c.Nop, c.Adc, c.Ror, c.Nop, c.Pla, c.Adc, c.Ror, c.Nop, c.Jmp, c.Adc, c.Ror, c.Nop,
		c.Bvs, c.Adc, c.Nop, c.Nop, c.Nop, c.Adc, c.Ror, c.Nop, c.Sei, c.Adc, c.Nop, c.Nop, c.Nop, c.Adc, c.Ror, c.Nop,
		c.Nop, c.Sta, c.Nop, c.Nop, c.Sty, c.Sta, c.Stx, c.Nop, c.Dey, c.Nop, c.Txa, c.Nop, c.Sty, c.Sta, c.Stx, c.Nop,
		c.Bcc, c.Sta, c.Nop, c.Nop, c.Sty, c.Sta, c.Stx, c.Nop, c.Tya, c.Sta, c.Txs, c.Nop, c.Nop, c.Sta, c.Nop, c.Nop,
		c.Ldy, c.Lda, c.Ldx, c.Nop, c.Ldy, c.Lda, c.Ldx, c.Nop, c.Tay, c.Lda, c.Tax, c.Nop, c.Ldy, c.Lda, c.Ldx, c.Nop,
		c.Bcs, c.Lda, c.Nop, c.Nop, c.Ldy, c.Lda, c.Ldx, c.Nop, c.Clv, c.Lda, c.Tsx, c.Nop, c.Ldy, c.Lda, c.Ldx, c.Nop,
		c.Cpy, c.Cmp, c.Nop, c.Nop, c.Cpy, c.Cmp, c.Dec, c.Nop, c.Iny, c.Cmp, c.Dex, c.Nop, c.Cpy, c.Cmp, c.Dec, c.Nop,
		c.Bne, c.Cmp, c.Nop, c.Nop, c.Nop, c.Cmp, c.Dec, c.Nop, c.Cld, c.Cmp, c.Nop, c.Nop, c.Nop, c.Cmp, c.Dec, c.Nop,
		c.Cpx, c.Sbc, c.Nop, c.Nop, c.Cpx, c.Sbc, c.Inc, c.Nop, c.Inx, c.Sbc, c.Nop, c.Sbc, c.Cpx, c.Sbc, c.Inc, c.Nop,
		c.Beq, c.Sbc, c.Nop, c.Nop, c.Nop, c.Sbc, c.Inc, c.Nop, c.Sed, c.Sbc, c.Nop, c.Nop, c.Nop, c.Sbc, c.Inc, c.Nop,
	}
}

func (c *CPU) Reset() {
	c.PC = c.Read16(0xFFFC, false) // 读取入口地址
	c.SP = 0xFD                    // 初始化堆栈与标记
	c.SetFlags(0x24)
}

// 判断是否发生的切页
func IsPageDiff(a, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}

// 分支跳转本来就需要一个周期，若是发生跳转需要一个额外周期
func (c *CPU) AddBranchCycle(info *StepInfo) {
	c.Cycles++
	if IsPageDiff(info.PC, info.Addr) {
		c.Cycles++
	}
}

// 通用比较操作，修改 flag 寄存器
func (c *CPU) Compare(a, b uint8) {
	c.SetZN(a - b)
	if a >= b {
		c.C = 1
	} else {
		c.C = 0
	}
}

func (c *CPU) Read16(addr uint16, debug bool) uint16 {
	lo := uint16(c.Read(addr, debug))
	hi := uint16(c.Read(addr+1, debug))
	return hi<<8 | lo
}

// 模拟 6502 CPU 的 BUG 低位进位不进到高位去
func (c *CPU) Read16Bug(addr uint16, debug bool) uint16 {
	a := addr
	b := (a & 0xFF00) | uint16(uint8(a)+1)
	lo := uint16(c.Read(a, debug))
	hi := uint16(c.Read(b, debug))
	return hi<<8 | lo
}

func (c *CPU) Push(val uint8) {
	c.Write(0x100|uint16(c.SP), val)
	c.SP--
}

func (c *CPU) Pop() uint8 {
	c.SP++
	return c.Read(0x100|uint16(c.SP), false)
}

func (c *CPU) Push16(value uint16) {
	hi := uint8(value >> 8)
	lo := uint8(value & 0xFF)
	c.Push(hi)
	c.Push(lo)
}

func (c *CPU) Pop16() uint16 {
	lo := uint16(c.Pop())
	hi := uint16(c.Pop())
	return hi<<8 | lo
}

func (c *CPU) Flags() uint8 {
	flags := uint8(0)
	flags |= c.C << 0
	flags |= c.Z << 1
	flags |= c.I << 2
	flags |= c.D << 3
	flags |= c.B << 4
	flags |= c.U << 5
	flags |= c.V << 6
	flags |= c.N << 7
	return flags
}

func (c *CPU) SetFlags(flags uint8) {
	c.C = (flags >> 0) & 1
	c.Z = (flags >> 1) & 1
	c.I = (flags >> 2) & 1
	c.D = (flags >> 3) & 1
	c.B = (flags >> 4) & 1
	c.U = (flags >> 5) & 1
	c.V = (flags >> 6) & 1
	c.N = (flags >> 7) & 1
}

func (c *CPU) SetZ(val uint8) {
	if val == 0 {
		c.Z = 1
	} else {
		c.Z = 0
	}
}

func (c *CPU) SetN(val uint8) {
	if val&0x80 != 0 {
		c.N = 1
	} else {
		c.N = 0
	}
}

func (c *CPU) SetZN(val uint8) {
	c.SetZ(val)
	c.SetN(val)
}

func (c *CPU) TriggerNMI() {
	c.IntType = IntNMI
}

func (c *CPU) TriggerIRQ() {
	if c.I == 0 {
		c.IntType = IntIRQ
	}
}

type StepInfo struct {
	Addr uint16
	PC   uint16
	Mode uint8
}

func (c *CPU) Step() int {
	if c.LastCycles > 0 {
		c.LastCycles--
		return 1
	}

	oldCycles := c.Cycles

	switch c.IntType {
	case IntNMI:
		c.Nmi()
	case IntIRQ:
		c.Irq()
	}
	c.IntType = IntNone

	opCode := c.Read(c.PC, false)
	mode := InstAddrModes[opCode]
	addr := uint16(0)
	pageDiff := false
	switch mode {
	case AddrAbsolute:
		addr = c.Read16(c.PC+1, false)
	case AddrAbsoluteX:
		addr = c.Read16(c.PC+1, false) + uint16(c.X)
		pageDiff = IsPageDiff(addr-uint16(c.X), addr)
	case AddrAbsoluteY:
		addr = c.Read16(c.PC+1, false) + uint16(c.Y)
		pageDiff = IsPageDiff(addr-uint16(c.Y), addr)
	case AddrAccumulator:
		addr = 0 // 跟 AddrImplied 类似
	case AddrImmediate:
		addr = c.PC + 1
	case AddrImplied:
		addr = 0
	case AddrIndexedIndirect:
		addr = c.Read16Bug(uint16(c.Read(c.PC+1, false))+uint16(c.X), false)
	case AddrIndirect:
		addr = c.Read16Bug(c.Read16(c.PC+1, false), false)
	case AddrIndirectIndexed:
		addr = c.Read16Bug(uint16(c.Read(c.PC+1, false)), false) + uint16(c.Y)
		pageDiff = IsPageDiff(addr-uint16(c.Y), addr)
	case AddrRelative:
		offset := uint16(c.Read(c.PC+1, false))
		if offset < 0x80 {
			addr = c.PC + 2 + offset
		} else {
			addr = c.PC + 2 + offset - 0x100
		}
	case AddrZeroPage:
		addr = uint16(c.Read(c.PC+1, false))
	case AddrZeroPageX:
		addr = uint16(c.Read(c.PC+1, false)+c.X) & 0xff
	case AddrZeroPageY:
		addr = uint16(c.Read(c.PC+1, false)+c.Y) & 0xff
	}

	c.PC += uint16(InstSizes[opCode])
	c.Cycles += uint64(InstCycles[opCode])
	if pageDiff {
		c.Cycles += uint64(InstPageCycles[opCode])
	}
	info := &StepInfo{addr, c.PC, mode}
	c.InstTable[opCode](info)
	return int(c.Cycles - oldCycles)
}

func (c *CPU) Disassemble(start uint16, end uint16) map[uint16]string {
	res := make(map[uint16]string)
	pc := start
	for pc < end {
		oldPC := pc
		opCode := c.Read(pc, true)
		pc++
		msg := fmt.Sprintf("$%04X: %s ", oldPC, InstNames[opCode])
		switch InstAddrModes[opCode] {
		case AddrAbsolute:
			msg += fmt.Sprintf("$%04X {ABS}", c.Read16(pc, true))
			pc += 2
		case AddrAbsoluteX:
			msg += fmt.Sprintf("$%04X {ABX}", c.Read16(pc, true)+uint16(c.X))
			pc += 2
		case AddrAbsoluteY:
			msg += fmt.Sprintf("$%04X {ABY}", c.Read16(pc, true)+uint16(c.Y))
			pc += 2
		case AddrAccumulator:
			msg += "{ACC}"
		case AddrImmediate:
			msg += fmt.Sprintf("$%04X {IMM}", pc)
			pc++ // 这里会使用下一个 pc 需要自动后移
		case AddrImplied:
			msg += "{IMP}"
		case AddrIndexedIndirect:
			msg += fmt.Sprintf("$%04X {IZX}", c.Read16Bug(uint16(c.Read(pc, true))+uint16(c.X), true))
			pc++
		case AddrIndirect:
			msg += fmt.Sprintf("$%04X {IND}", c.Read16Bug(c.Read16(pc, true), true))
			pc += 2
		case AddrIndirectIndexed:
			msg += fmt.Sprintf("$%04X {IZY}", c.Read16Bug(uint16(c.Read(pc, true))+uint16(c.Y), true))
			pc++
		case AddrRelative:
			addr := uint16(c.Read(pc, true))
			pc++
			if addr < 0x80 {
				addr = uint16(pc) + 2 + addr
			} else {
				addr = uint16(pc) + 2 + addr - 0x100
			}
			msg += fmt.Sprintf("$%04X {REL}", addr)
		case AddrZeroPage:
			msg += fmt.Sprintf("$%04X {ZP0}", c.Read(pc, true))
			pc++
		case AddrZeroPageX:
			msg += fmt.Sprintf("$%04X {ZPX}", (c.Read(pc, true)+c.X)&0xff)
			pc++
		case AddrZeroPageY:
			msg += fmt.Sprintf("$%04X {ZPY}", (c.Read(pc, true)+c.Y)&0xff)
			pc++
		}
		res[oldPC] = msg
	}
	return res
}

func (c *CPU) DisassembleCode(line int) string {
	pc := c.PC
	buff := &strings.Builder{}
	for pc < 0xFFFF && line > 0 {
		opCode := c.Read(pc, true)
		buff.WriteString(fmt.Sprintf("$%04X: %s ", pc, InstNames[opCode]))
		pc++
		switch InstAddrModes[opCode] {
		case AddrAbsolute:
			buff.WriteString(fmt.Sprintf("$%04X {ABS}\n", c.Read16(pc, true)))
			pc += 2
		case AddrAbsoluteX:
			buff.WriteString(fmt.Sprintf("$%04X {ABX}\n", c.Read16(pc, true)+uint16(c.X)))
			pc += 2
		case AddrAbsoluteY:
			buff.WriteString(fmt.Sprintf("$%04X {ABY}\n", c.Read16(pc, true)+uint16(c.Y)))
			pc += 2
		case AddrAccumulator:
			buff.WriteString("{ACC}\n")
		case AddrImmediate:
			buff.WriteString(fmt.Sprintf("$%04X {IMM}\n", pc))
			pc++ // 这里会使用下一个 pc 需要自动后移
		case AddrImplied:
			buff.WriteString("{IMP}\n")
		case AddrIndexedIndirect:
			buff.WriteString(fmt.Sprintf("$%04X {IZX}\n", c.Read16Bug(uint16(c.Read(pc, true))+uint16(c.X), true)))
			pc++
		case AddrIndirect:
			buff.WriteString(fmt.Sprintf("$%04X {IND}\n", c.Read16Bug(c.Read16(pc, true), true)))
			pc += 2
		case AddrIndirectIndexed:
			buff.WriteString(fmt.Sprintf("$%04X {IZY}\n", c.Read16Bug(uint16(c.Read(pc, true))+uint16(c.Y), true)))
			pc++
		case AddrRelative:
			addr := uint16(c.Read(pc, true))
			pc++
			if addr < 0x80 {
				addr = uint16(pc) + 2 + addr
			} else {
				addr = uint16(pc) + 2 + addr - 0x100
			}
			buff.WriteString(fmt.Sprintf("$%04X {REL}\n", addr))
		case AddrZeroPage:
			buff.WriteString(fmt.Sprintf("$%04X {ZP0}\n", c.Read(pc, true)))
			pc++
		case AddrZeroPageX:
			buff.WriteString(fmt.Sprintf("$%04X {ZPX}\n", (c.Read(pc, true)+c.X)&0xff))
			pc++
		case AddrZeroPageY:
			buff.WriteString(fmt.Sprintf("$%04X {ZPY}\n", (c.Read(pc, true)+c.Y)&0xff))
			pc++
		}
		line--
	}
	return buff.String()
}

func (c *CPU) Nmi() {
	c.Push16(c.PC)
	c.Php(nil)
	c.PC = c.Read16(0xFFFA, false)
	c.I = 1
	c.Cycles += 7
}

func (c *CPU) Irq() {
	c.Push16(c.PC)
	c.Php(nil)
	c.PC = c.Read16(0xFFFE, false)
	c.I = 1
	c.Cycles += 7
}

func (c *CPU) Adc(info *StepInfo) {
	a := c.A
	b := c.Read(info.Addr, false)
	c0 := c.C
	c.A = a + b + c0
	c.SetZN(c.A)
	if int(a)+int(b)+int(c0) > 0xFF {
		c.C = 1
	} else {
		c.C = 0
	}
	if (a^b)&0x80 == 0 && (a^c.A)&0x80 != 0 {
		c.V = 1
	} else {
		c.V = 0
	}
}

func (c *CPU) And(info *StepInfo) {
	c.A = c.A & c.Read(info.Addr, false)
	c.SetZN(c.A)
}

func (c *CPU) Asl(info *StepInfo) {
	if info.Mode == AddrAccumulator {
		c.C = (c.A >> 7) & 1
		c.A <<= 1
		c.SetZN(c.A)
	} else {
		value := c.Read(info.Addr, false)
		c.C = (value >> 7) & 1
		value <<= 1
		c.Write(info.Addr, value)
		c.SetZN(value)
	}
}

func (c *CPU) Bcc(info *StepInfo) {
	if c.C == 0 {
		c.PC = info.Addr
		c.AddBranchCycle(info)
	}
}

func (c *CPU) Bcs(info *StepInfo) {
	if c.C != 0 {
		c.PC = info.Addr
		c.AddBranchCycle(info)
	}
}

func (c *CPU) Beq(info *StepInfo) {
	if c.Z != 0 {
		c.PC = info.Addr
		c.AddBranchCycle(info)
	}
}

func (c *CPU) Bit(info *StepInfo) {
	value := c.Read(info.Addr, false)
	c.V = (value >> 6) & 1
	c.SetZ(value & c.A)
	c.SetN(value)
}

func (c *CPU) Bmi(info *StepInfo) {
	if c.N != 0 {
		c.PC = info.Addr
		c.AddBranchCycle(info)
	}
}

func (c *CPU) Bne(info *StepInfo) {
	if c.Z == 0 {
		c.PC = info.Addr
		c.AddBranchCycle(info)
	}
}

func (c *CPU) Bpl(info *StepInfo) {
	if c.N == 0 {
		c.PC = info.Addr
		c.AddBranchCycle(info)
	}
}

func (c *CPU) Brk(info *StepInfo) {
	c.Push16(c.PC)
	c.Php(info)
	c.Sei(info)
	c.PC = c.Read16(0xFFFE, false)
}

func (c *CPU) Bvc(info *StepInfo) {
	if c.V == 0 {
		c.PC = info.Addr
		c.AddBranchCycle(info)
	}
}

func (c *CPU) Bvs(info *StepInfo) {
	if c.V != 0 {
		c.PC = info.Addr
		c.AddBranchCycle(info)
	}
}

func (c *CPU) Clc(_ *StepInfo) {
	c.C = 0
}

func (c *CPU) Cld(_ *StepInfo) {
	c.D = 0
}

func (c *CPU) Cli(_ *StepInfo) {
	c.I = 0
}

func (c *CPU) Clv(_ *StepInfo) {
	c.V = 0
}

func (c *CPU) Cmp(info *StepInfo) {
	value := c.Read(info.Addr, false)
	c.Compare(c.A, value)
}

func (c *CPU) Cpx(info *StepInfo) {
	value := c.Read(info.Addr, false)
	c.Compare(c.X, value)
}

func (c *CPU) Cpy(info *StepInfo) {
	value := c.Read(info.Addr, false)
	c.Compare(c.Y, value)
}

func (c *CPU) Dec(info *StepInfo) {
	value := c.Read(info.Addr, false) - 1
	c.Write(info.Addr, value)
	c.SetZN(value)
}

func (c *CPU) Dex(_ *StepInfo) {
	c.X--
	c.SetZN(c.X)
}

func (c *CPU) Dey(_ *StepInfo) {
	c.Y--
	c.SetZN(c.Y)
}

func (c *CPU) Eor(info *StepInfo) {
	c.A = c.A ^ c.Read(info.Addr, false)
	c.SetZN(c.A)
}

func (c *CPU) Inc(info *StepInfo) {
	value := c.Read(info.Addr, false) + 1
	c.Write(info.Addr, value)
	c.SetZN(value)
}

func (c *CPU) Inx(_ *StepInfo) {
	c.X++
	c.SetZN(c.X)
}

func (c *CPU) Iny(_ *StepInfo) {
	c.Y++
	c.SetZN(c.Y)
}

func (c *CPU) Jmp(info *StepInfo) {
	c.PC = info.Addr
}

func (c *CPU) Jsr(info *StepInfo) {
	c.Push16(c.PC - 1)
	c.PC = info.Addr
}

func (c *CPU) Lda(info *StepInfo) {
	c.A = c.Read(info.Addr, false)
	c.SetZN(c.A)
}

func (c *CPU) Ldx(info *StepInfo) {
	c.X = c.Read(info.Addr, false)
	c.SetZN(c.X)
}

func (c *CPU) Ldy(info *StepInfo) {
	c.Y = c.Read(info.Addr, false)
	c.SetZN(c.Y)
}

func (c *CPU) Lsr(info *StepInfo) {
	if info.Mode == AddrAccumulator {
		c.C = c.A & 1
		c.A >>= 1
		c.SetZN(c.A)
	} else {
		value := c.Read(info.Addr, false)
		c.C = value & 1
		value >>= 1
		c.Write(info.Addr, value)
		c.SetZN(value)
	}
}

func (c *CPU) Ora(info *StepInfo) {
	c.A = c.A | c.Read(info.Addr, false)
	c.SetZN(c.A)
}

func (c *CPU) Pha(_ *StepInfo) {
	c.Push(c.A)
}

func (c *CPU) Php(_ *StepInfo) {
	c.Push(c.Flags() | 0x10)
}

func (c *CPU) Pla(_ *StepInfo) {
	c.A = c.Pop()
	c.SetZN(c.A)
}

func (c *CPU) Plp(_ *StepInfo) {
	c.SetFlags(c.Pop()&0xEF | 0x20)
}

func (c *CPU) Rol(info *StepInfo) {
	if info.Mode == AddrAccumulator {
		c0 := c.C
		c.C = (c.A >> 7) & 1
		c.A = (c.A << 1) | c0
		c.SetZN(c.A)
	} else {
		c0 := c.C
		value := c.Read(info.Addr, false)
		c.C = (value >> 7) & 1
		value = (value << 1) | c0
		c.Write(info.Addr, value)
		c.SetZN(value)
	}
}

func (c *CPU) Ror(info *StepInfo) {
	if info.Mode == AddrAccumulator {
		c0 := c.C
		c.C = c.A & 1
		c.A = (c.A >> 1) | (c0 << 7)
		c.SetZN(c.A)
	} else {
		c0 := c.C
		value := c.Read(info.Addr, false)
		c.C = value & 1
		value = (value >> 1) | (c0 << 7)
		c.Write(info.Addr, value)
		c.SetZN(value)
	}
}

func (c *CPU) Rti(_ *StepInfo) {
	c.SetFlags(c.Pop()&0xEF | 0x20)
	c.PC = c.Pop16()
}

func (c *CPU) Rts(_ *StepInfo) {
	c.PC = c.Pop16() + 1
}

func (c *CPU) Sbc(info *StepInfo) {
	a := c.A
	b := c.Read(info.Addr, false)
	c0 := c.C
	c.A = a - b - (1 - c0)
	c.SetZN(c.A)
	if int(a)-int(b)-int(1-c0) >= 0 {
		c.C = 1
	} else {
		c.C = 0
	}
	if (a^b)&0x80 != 0 && (a^c.A)&0x80 != 0 {
		c.V = 1
	} else {
		c.V = 0
	}
}

func (c *CPU) Sec(_ *StepInfo) {
	c.C = 1
}

func (c *CPU) Sed(_ *StepInfo) {
	c.D = 1
}

func (c *CPU) Sei(_ *StepInfo) {
	c.I = 1
}

func (c *CPU) Sta(info *StepInfo) {
	c.Write(info.Addr, c.A)
}

func (c *CPU) Stx(info *StepInfo) {
	c.Write(info.Addr, c.X)
}

func (c *CPU) Sty(info *StepInfo) {
	c.Write(info.Addr, c.Y)
}

func (c *CPU) Tax(_ *StepInfo) {
	c.X = c.A
	c.SetZN(c.X)
}

func (c *CPU) Tay(_ *StepInfo) {
	c.Y = c.A
	c.SetZN(c.Y)
}

func (c *CPU) Tsx(_ *StepInfo) {
	c.X = c.SP
	c.SetZN(c.X)
}

func (c *CPU) Txa(_ *StepInfo) {
	c.A = c.X
	c.SetZN(c.A)
}

func (c *CPU) Txs(_ *StepInfo) {
	c.SP = c.X
}

func (c *CPU) Tya(_ *StepInfo) {
	c.A = c.Y
	c.SetZN(c.A)
}

func (c *CPU) Nop(_ *StepInfo) {
}
