/*
@author: sk
@date: 2024/8/3
*/
package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestInst(t *testing.T) {
	str := "c.Brk, c.Ora, c.Nop, c.Nop, c.Nop, c.Ora, c.Asl, c.Nop, c.Php, c.Ora, c.Asl, c.Nop, c.Nop, c.Ora, c.Asl, c.Nop,\n\t\tc.Bpl, c.Ora, c.Nop, c.Nop, c.Nop, c.Ora, c.Asl, c.Nop, c.Clc, c.Ora, c.Nop, c.Nop, c.Nop, c.Ora, c.Asl, c.Nop,\n\t\tc.Jsr, c.And, c.Nop, c.Nop, c.Bit, c.And, c.rol, c.Nop, c.Plp, c.And, c.rol, c.Nop, c.Bit, c.And, c.rol, c.Nop,\n\t\tc.Bmi, c.And, c.Nop, c.Nop, c.Nop, c.And, c.rol, c.Nop, c.Sec, c.And, c.Nop, c.Nop, c.Nop, c.And, c.rol, c.Nop,\n\t\tc.Rti, c.Eor, c.Nop, c.Nop, c.Nop, c.Eor, c.Lsr, c.Nop, c.Pha, c.Eor, c.Lsr, c.Nop, c.Jmp, c.Eor, c.Lsr, c.Nop,\n\t\tc.Bvc, c.Eor, c.Nop, c.Nop, c.Nop, c.Eor, c.Lsr, c.Nop, c.Cli, c.Eor, c.Nop, c.Nop, c.Nop, c.Eor, c.Lsr, c.Nop,\n\t\tc.Rts, c.Adc, c.Nop, c.Nop, c.Nop, c.Adc, c.Ror, c.Nop, c.Pla, c.Adc, c.Ror, c.Nop, c.Jmp, c.Adc, c.Ror, c.Nop,\n\t\tc.Bvs, c.Adc, c.Nop, c.Nop, c.Nop, c.Adc, c.Ror, c.Nop, c.Sei, c.Adc, c.Nop, c.Nop, c.Nop, c.Adc, c.Ror, c.Nop,\n\t\tc.Nop, c.Sta, c.Nop, c.Nop, c.Sty, c.Sta, c.Stx, c.Nop, c.Dey, c.Nop, c.Txa, c.Nop, c.Sty, c.Sta, c.Stx, c.Nop,\n\t\tc.Bcc, c.Sta, c.Nop, c.Nop, c.Sty, c.Sta, c.Stx, c.Nop, c.Tya, c.Sta, c.Txs, c.Nop, c.Nop, c.Sta, c.Nop, c.Nop,\n\t\tc.Ldy, c.Lda, c.Ldx, c.Nop, c.Ldy, c.Lda, c.Ldx, c.Nop, c.Tay, c.Lda, c.Tax, c.Nop, c.Ldy, c.Lda, c.Ldx, c.Nop,\n\t\tc.Bcs, c.Lda, c.Nop, c.Nop, c.Ldy, c.Lda, c.Ldx, c.Nop, c.Clv, c.Lda, c.Tsx, c.Nop, c.Ldy, c.Lda, c.Ldx, c.Nop,\n\t\tc.Cpy, c.Cmp, c.Nop, c.Nop, c.Cpy, c.Cmp, c.Dec, c.Nop, c.Iny, c.Cmp, c.Dex, c.Nop, c.Cpy, c.Cmp, c.Dec, c.Nop,\n\t\tc.Bne, c.Cmp, c.Nop, c.Nop, c.Nop, c.Cmp, c.Dec, c.Nop, c.Cld, c.Cmp, c.Nop, c.Nop, c.Nop, c.Cmp, c.Dec, c.Nop,\n\t\tc.Cpx, c.Sbc, c.Nop, c.Nop, c.Cpx, c.Sbc, c.Inc, c.Nop, c.Inx, c.Sbc, c.Nop, c.Sbc, c.Cpx, c.Sbc, c.Inc, c.Nop,\n\t\tc.Beq, c.Sbc, c.Nop, c.Nop, c.Nop, c.Sbc, c.Inc, c.Nop, c.sed, c.Sbc, c.Nop, c.Nop, c.Nop, c.Sbc, c.Inc, c.Nop"
	buff := &strings.Builder{}
	items := strings.Split(str, "\n\t\t")
	for _, item := range items {
		insts := strings.Split(item, ",")
		for _, inst := range insts {
			inst = strings.TrimSpace(inst)
			if len(inst) > 0 {
				buff.WriteString(fmt.Sprintf("\"%s\",", strings.ToUpper(inst[2:])))
			}
		}
		buff.WriteString("\n")
	}
	fmt.Println(buff)
}
