package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Provide an input file name as the first argument")
	}
	input, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("Can't parse file: ", err)
	}

	const (
		H = 6
		W = 25
	)

	L := len(input) / (H * W)
	counts := make(map[int]map[byte]int, L)

	for i := 0; i < L; i++ {
		counts[i] = map[byte]int{'0': 0, '1': 0, '2': 0}
		for j := 0; j < H; j++ {
			for z := 0; z < W; z++ {
				counts[i][input[z+j*W+i*H*W]]++
			}
		}
	}

	var min map[byte]int
	for _, v := range counts {
		if min == nil || v['0'] < min['0'] {
			min = v
		}
	}

	fmt.Println("Part 1:", min['1']*min['2'])

	image := make([][H][W]byte, L)
	for i := range image {
		for j := range image[i] {
			for z := range image[i][j] {
				image[i][j][z] = input[z+j*W+i*H*W]
			}
		}
	}

	fmt.Println("Part 2:")
	writer := bufio.NewWriter(os.Stdout)
	for j := 0; j < H; j++ {
		for z := 0; z < W; z++ {
			var pixel byte
			for i := 0; i < L; i++ {
				if layerPixel := image[i][j][z]; layerPixel != '2' {
					pixel = layerPixel
					break
				}
			}
			writer.WriteByte(pixel)
		}
		writer.WriteByte('\n')
	}
	writer.Flush()
}
