package main

import (
	"log"
	"os"
	"os/signal"
	"time"
)

type ConditionFlag = uint16

const (
	FL_POS ConditionFlag = 1 << 0 /* P */
	FL_ZRO               = 1 << 1 /* Z */
	FL_NEG               = 1 << 2 /* N */
)

const STEP = 0

func main() {
	init_input()

	c := make(chan os.Signal, 3)
	signal.Notify(c, os.Interrupt)

	read_image(get_path())

	go loop()

	<-c
	stop_input()
	log.Println("Process interrupted")
}

func loop() {
	running := true

	reg[R_COND] = FL_ZRO

	const PC_START = 0x3000
	reg[R_PC] = PC_START

	//m := debug_map()
	go listen_keyboard()
	for running {
		instr := mem_read(reg[R_PC])
		time.Sleep(STEP * time.Millisecond)
		reg[R_PC] += 1
		op := Opcode(instr >> 12)
		//fmt.Printf("%X) %v\n", reg[R_PC]-1, m[op])
		switch op {
		case OP_ADD:
			dr := (instr >> 9) & 0x7
			r1 := (instr >> 6) & 0x7
			immFlag := (instr >> 5) & 0x1
			if immFlag == 1 {
				v := sign_extend(instr&0x1F, 5)
				reg[dr] = reg[r1] + v
			} else {
				r2 := instr & 0x7
				reg[dr] = reg[r1] + reg[r2]
			}

			update_flags(reg[dr])
		case OP_AND:
			dr := (instr >> 9) & 0x7
			r1 := (instr >> 6) & 0x7

			immFlag := (instr >> 5) & 0x1

			if immFlag == 1 {
				v := sign_extend(instr&0x1F, 5)
				reg[dr] = reg[r1] & v
			} else {
				r2 := instr & 0x7
				reg[dr] = reg[r1] & reg[r2]
			}

			update_flags(reg[dr])
		case OP_NOT:
			dr := (instr >> 9) & 0x7
			sr := (instr >> 6) & 0x7

			reg[dr] = ^reg[sr]

			update_flags(reg[dr])
		case OP_BR:
			cf := (instr >> 9) & 0x7
			if (cf & reg[R_COND]) != 0 {
				reg[R_PC] += sign_extend(instr&0x1FF, 9)
			}
		case OP_JMP:
			br := (instr >> 6) & 0x7
			reg[R_PC] = reg[br]
		case OP_JSR:
			reg[R_R7] = reg[R_PC]
			flag := (instr >> 11) & 1

			if flag == 1 {
				reg[R_PC] += sign_extend(instr&0x7FF, 11)
			} else {
				br := (instr >> 6) & 0x7
				reg[R_PC] = reg[br]
			}
		case OP_LD:
			dr := (instr >> 9) & 0x7
			reg[dr] = mem_read(reg[R_PC] + sign_extend(instr&0x1FF, 9))

			update_flags(reg[dr])
		case OP_LDI:
			dr := (instr >> 9) & 0x7
			pcOffset9 := sign_extend(instr&0x1FF, 9)

			reg[dr] = mem_read(mem_read(reg[R_PC] + pcOffset9))

			update_flags(reg[dr])
		case OP_LDR:
			dr := (instr >> 9) & 0x7
			br := (instr >> 6) & 0x7
			reg[dr] = mem_read(reg[br] + sign_extend(instr&0x3F, 6))

			update_flags(reg[dr])
		case OP_LEA:
			dr := (instr >> 9) & 0x7
			pcOffset9 := sign_extend(instr&0x1FF, 9)

			reg[dr] = reg[R_PC] + pcOffset9

			update_flags(reg[dr])
		case OP_ST:
			sr := (instr >> 9) & 0x7
			pcOffset9 := sign_extend(instr&0x1FF, 9)

			mem_store(reg[R_PC]+pcOffset9, reg[sr])

		case OP_STI:
			sr := (instr >> 9) & 0x7
			pcOffset9 := sign_extend(instr&0x1FF, 9)

			mem_store(mem_read(reg[R_PC]+pcOffset9), reg[sr])
		case OP_STR:
			sr := (instr >> 9) & 0x7
			br := (instr >> 6) & 0x7
			pcOffset6 := sign_extend(instr&0x3F, 6)

			mem_store(reg[br]+pcOffset6, reg[sr])
		case OP_TRAP:
			reg[R_R7] = reg[R_PC]
			trapv := instr & 0xFF
			running = trap_call(trapv)
		case OP_RES:
		case OP_RTI:
		default:
			log.Println("Unused opcode")
			running = false
		}
		//print_reg_state()
	}

	os.Exit(0)
}
