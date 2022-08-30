package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func get_path() string {
	if len(os.Args) < 1 {
		log.Fatal("path to image not provided")
	}

	path := os.Args[1]
	return path
}

func sign_extend(x uint16, bCount int) uint16 {
	if ((x >> (bCount - 1)) & 1) == 1 {
		x |= 0xFFFF << bCount
	}

	return x
}

func update_flags(x uint16) {
	if x == 0 {
		reg[R_COND] = FL_ZRO
	} else if x>>15 == 1 {
		reg[R_COND] = FL_NEG
	} else {
		reg[R_COND] = FL_POS
	}
}

func halt() {
	fmt.Println("\nProgram halted...")
}

func unused_opcode() {
	fmt.Println("unused opcode")
}

func init_input() {
	exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-f", "/dev/tty", "-echo").Run()
}

func stop_input() {
	exec.Command("stty", "-f", "/dev/tty", "echo").Run()
}
