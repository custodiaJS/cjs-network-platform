package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/custodiaJs/cjs-network-platform/npvm"
)

func main() {
	// Es wird eine neue Instanz erzeugt
	if err := npvm.InitVmInstance(); err != nil {
		panic(err)
	}

	// Alle Dateien der VM werden angezeigt
	fmt.Println(npvm.GetAllRunningFiles())

	// Die VM wird gestartet
	go func() {
		if err := npvm.ServeNpVm(); err != nil {
			panic(err)
		}
	}()

	// Es wird auf ein Enter gewartet
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Die Vm wird vollständig geschlossen
	npvm.ShutdownNpVm()

	// Es wird aufgeräumt
	npvm.CleanUp()
}
