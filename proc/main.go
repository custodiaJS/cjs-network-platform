package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/custodiaJs/cjs-network-platform/host"
	"github.com/custodiaJs/cjs-network-platform/npvm"
)

func main() {
	// Es wird geprüft ob die benötigten Benutzerrechte vorhanden sind
	if !host.HasPrivilegesToRunNpVm() {
		fmt.Println("You are not authorized to run CJS-NpVm.")
		os.Exit(1)
	}

	// Es wird geprüft ob ausreichend Speicher Festplatten speichert verfügbar ist
	if !host.CheckHostHasEnoughDiskSpace(10000) {
		fmt.Println("Not enough free disk space to run NpVm.")
		os.Exit(1)
	}

	// Es wird geprüft ob Mindestens 256 Mb RAM verfüggbar sind
	if !host.CheckHostHasEnoughFreeRam(256) {
		fmt.Println("Not enough free memory to run NpVm.")
		os.Exit(1)
	}

	// Es wird eine neue Instanz erzeugt
	if err := npvm.InitVmInstance(); err != nil {
		panic(err)
	}

	// Die Parameter werden eingelesen
	if err := parseAndUseArgs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
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
