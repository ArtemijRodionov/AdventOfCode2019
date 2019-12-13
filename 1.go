package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
    "os"
)

func main() {
	content, err := ioutil.ReadAll(os.Stdin)
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
	fmt.Println(fuelSum)
}
