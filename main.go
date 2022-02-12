package main

import (
	"os"

	iprange "github.com/ryuichi1208/go-iprange-cloud/iprange"
)

func main() {
	iprange.Run(os.Args)
}
