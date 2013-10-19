package main
import (
	"fmt"
)

/* packet type 1, containing client info */
type PktFileChunkMap struct {
	pktid int
	path string
	file_map []byte
}

func (p* PktFileChunkMap) GetPkid() int {
	return p.pktid;
}

func (p* PktFileChunkMap) Print() {
	fmt.Println("I'm a PktFileChunkMap struct.")
}