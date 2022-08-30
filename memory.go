package main

import (
	"bufio"
	"os"
	"time"
)

var keyboard_ch = make(chan uint16, 1)

func listen_keyboard() {
	r := bufio.NewReader(os.Stdin)
	for {
		c, _, _ := r.ReadRune()
		keyboard_ch <- uint16(c)
	}
}

const MEMORY_MAX = 1 << 16

const (
	MR_KBSR uint16 = 0xFE00 /* keyboard status */
	MR_KBDR uint16 = 0xFE02 /* keyboard data */
)

var memory = [MEMORY_MAX]uint16{}

func mem_read(i uint16) uint16 {
	if i == MR_KBSR {
		if key := check_key(); key != 0 {
			memory[MR_KBSR] = 1 << 15
			memory[MR_KBDR] = key
		} else {
			memory[MR_KBSR] = 0
		}
	}
	return memory[i]
}

func mem_store(i uint16, v uint16) {
	memory[i] = v
}

func check_key() uint16 {
	t := time.NewTimer(10 * time.Millisecond)
	select {
	case <-t.C:
		return 0
	case key := <-keyboard_ch:
		return key
	}
}
