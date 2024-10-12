package npvm

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func CleanUp() {
	mu.Lock()
	defer mu.Unlock()

	// Funktion zur Bereinigung des temporären Verzeichnisses
	cleanup := func() {
		err := os.RemoveAll(tempdirPath)
		if err != nil {
			log.Printf("Warnung: Fehler beim Löschen des temporären Verzeichnisses: %v", err)
		} else {
			fmt.Printf("Temporäres Verzeichnis %s gelöscht.\n", tempdirPath)
		}
	}

	// Signal-Handling für saubere Beendigung
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Printf("\nErhaltenes Signal: %s. Bereinige und beende...\n", sig)
		cleanup()
		os.Exit(0)
	}()
}
