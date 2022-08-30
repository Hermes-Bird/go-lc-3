package main

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

func read_image(path string) {
	file, err := os.Open(path)
	defer file.Close()

	bs := make([]byte, 2)
	_, err = file.Read(bs)

	if err != nil {
		log.Fatal("Failed to read image")
	}

	origin := binary.BigEndian.Uint16(bs)
	for i := int(origin); i < MEMORY_MAX; i++ {
		_, err = file.Read(bs)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		mem_store(uint16(i), binary.BigEndian.Uint16(bs))
	}
}
