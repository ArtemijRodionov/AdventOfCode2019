package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Provide a file name as argument")
	}
	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	fuelSum := 0
	for _, s := range strings.Split(string(content), "\n") {
		v, err := strconv.Atoi(s)
		if err != nil {
			continue
		}

		fuelSum += int(v/3.0) - 2
	}
	fmt.Println("1:", fuelSum)

	fuelSum = 0
	for _, s := range strings.Split(string(content), "\n") {
		v, err := strconv.Atoi(s)
		if err != nil {
			continue
		}

		fuel := int(v/3.0) - 2
		for fuel > 0 {
			fuelSum += fuel
			fuel = int(fuel/3.0) - 2
		}
	}
	fmt.Println("2:", fuelSum)
}
