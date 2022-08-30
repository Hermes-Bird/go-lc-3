package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	TRAP_GETC  = 0x20 /* get character from keyboard, not echoed onto the terminal */
	TRAP_OUT   = 0x21 /* output a character */
	TRAP_PUTS  = 0x22 /* output a word string */
	TRAP_IN    = 0x23 /* get character from keyboard, echoed onto the terminal */
	TRAP_PUTSP = 0x24 /* output a byte string */
	TRAP_HALT  = 0x25 /* halt the program */
)

func trap_call(code uint16) bool {
	//m := debug_trap_map()
	//log.Println(m[code])
	switch code {
	case TRAP_GETC:
		r := bufio.NewReader(os.Stdin)
		c, _, _ := r.ReadRune()
		reg[R_R0] = uint16(c)
		update_flags(reg[R_R0])
	case TRAP_OUT:
		c := reg[R_R0] & 0xFF
		fmt.Printf("%c", c)
	case TRAP_PUTS:
		addr := reg[R_R0]
		for mem_read(addr) != 0 {
			fmt.Printf("%c", mem_read(addr))
			addr++
		}
	case TRAP_IN:
		print("enter char: ")
		r := bufio.NewReader(os.Stdin)
		c, _, _ := r.ReadRune()
		fmt.Printf("%c", c)
		reg[R_R0] = uint16(c)
		update_flags(reg[R_R0])
	case TRAP_PUTSP:
		addr := reg[R_R0]
		for mem_read(addr) != 0 {
			v := mem_read(addr)
			fmt.Printf("%c", v&0xFF)
			c2 := v >> 8
			if c2 != 0 {
				fmt.Printf("%c", v>>8)
			}
			addr++
		}
	case TRAP_HALT:
		halt()
		return false
	}

	return true
}
