package main

import (
	"encoding/binary"
	"io"
	"os"
)

const NESMagic = 0x1A53454E

type Cartridge struct {
	PRG    []byte // PRG-ROM 程序代码
	CHR    []byte // CHR-ROM 图块数据
	Mapper uint8  // mapper 类型
	Mirror uint8  // mirroring 类型
}

type NESHeader struct {
	Magic    uint32 // NES魔数
	PRGNum   uint8  // 程序块数 每个  16k
	CHRNum   uint8  // tile块数 每个  8k
	Control1 uint8  // 控制位 1
	Control2 uint8  // 控制位 2
	Unused   [8]byte
}

func LoadCartridge(path string) *Cartridge {
	file, err := os.Open(path)
	HandleErr(err)
	defer file.Close()

	header := NESHeader{}
	err = binary.Read(file, binary.LittleEndian, &header)
	HandleErr(err)
	if header.Magic != NESMagic {
		panic("not nes file")
	}

	mapper1 := header.Control1 >> 4
	mapper2 := header.Control2 >> 4
	mapper := mapper1 | mapper2<<4
	mirror1 := header.Control1 & 1
	mirror2 := (header.Control1 >> 3) & 1
	mirror := mirror1 | mirror2<<1
	if header.Control1&4 == 4 { // 这部分信息不需要直接舍弃
		_, err = file.Seek(512, 1)
		HandleErr(err)
	}
	prg := make([]byte, int(header.PRGNum)*16*1024)
	_, err = io.ReadFull(file, prg)
	HandleErr(err)
	chr := make([]byte, 8*1024) // 可能 tile 没有直接存储是后面加载的 预留 8k
	if header.CHRNum > 0 {
		chr = make([]byte, int(header.CHRNum)*8*1024)
		_, err = io.ReadFull(file, chr)
		HandleErr(err)
	}
	return &Cartridge{PRG: prg, CHR: chr, Mapper: mapper, Mirror: mirror}
}
