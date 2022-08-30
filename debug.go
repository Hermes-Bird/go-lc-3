package main

import "fmt"

func debug_map() map[Opcode]string {
	m := map[Opcode]string{}

	m[OP_ADD] = "ADD"
	m[OP_AND] = "AND"
	m[OP_JMP] = "JMP"
	m[OP_NOT] = "NOT"
	m[OP_BR] = "BR"
	m[OP_JSR] = "JSR"
	m[OP_LD] = "LD"
	m[OP_LDI] = "LDI"
	m[OP_LDR] = "LDR"
	m[OP_LEA] = "LEA"
	m[OP_ST] = "ST"
	m[OP_STI] = "STI"
	m[OP_STR] = "STR"
	m[OP_TRAP] = "TRAP"

	return m
}

func debug_trap_map() map[uint16]string {
	m := map[uint16]string{}
	m[TRAP_GETC] = "t_GETC"   /* get character from keyboard, not echoed onto the terminal */
	m[TRAP_OUT] = "t_OUT"     /* output a character */
	m[TRAP_PUTS] = "t_PUTS"   /* output a word string */
	m[TRAP_IN] = "t_IN"       /* get character from keyboard, echoed onto the terminal */
	m[TRAP_PUTSP] = "t_PUTSP" /* output a byte string */
	m[TRAP_HALT] = "t_HALT"
	return m
}

func debug_reg_map() map[Register]string {
	m := map[Register]string{}
	m[R_R0] = "R_R0"
	m[R_R1] = "R_R1"
	m[R_R2] = "R_R2"
	m[R_R3] = "R_R3"
	m[R_R4] = "R_R4"
	m[R_R5] = "R_R5"
	m[R_R6] = "R_R6"
	m[R_R7] = "R_R7"
	m[R_PC] = "R_PC"
	m[R_COND] = "R_COND"
	m[R_COUNT] = "R_COUNT"
	return m
}

var reg_map = debug_reg_map()

func print_reg_state() {
	fmt.Println("{")
	for i, v := range reg {
		fmt.Printf("\t[%s] %x\n", reg_map[Register(i)], v)
	}
	fmt.Println("{")
}

func print_mem_map() {
	fmt.Println("{")
	for a, v := range memory {
		if v != 0 {
			fmt.Printf("\t%x) %c\n", a, v&0xFF)
		}
	}
	fmt.Println("}")
}
