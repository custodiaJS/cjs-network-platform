package main

import (
	"fmt"

	"github.com/custodiaJs/cjs-network-platform/npvm"
)

func main() {
	vm, err := npvm.NewVmInstance()
	if err != nil {
		panic(err)
	}
	fmt.Println(vm)
}
