package npvm

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

//go:embed bin/*
var binaryFiles embed.FS

func _LoadVmImageResources() (*_NpVmImageResources, error) {
	// Erstellen eines temporären Verzeichnisses
	tempDir, err := os.MkdirTemp("", "embedded_files_*")
	if err != nil {
		log.Fatalf("Fehler beim Erstellen des temporären Verzeichnisses: %v", err)
	}

	// Funktion zur Bereinigung des temporären Verzeichnisses
	cleanup := func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			log.Printf("Warnung: Fehler beim Löschen des temporären Verzeichnisses: %v", err)
		} else {
			fmt.Printf("Temporäres Verzeichnis %s gelöscht.\n", tempDir)
		}
	}

	// Registrieren der Bereinigung bei Programmende
	defer cleanup()

	// Signal-Handling für saubere Beendigung
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Printf("\nErhaltenes Signal: %s. Bereinige und beende...\n", sig)
		cleanup()
		os.Exit(0)
	}()

	// Durchlaufen des eingebetteten Dateisystems
	err = fs.WalkDir(binaryFiles, "", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Bestimmen des relativen Pfads zur 'assets' Wurzel
		relPath, err := filepath.Rel("", path)
		if err != nil {
			return err
		}

		// Pfad im temporären Verzeichnis
		targetPath := filepath.Join(tempDir, relPath)

		if d.IsDir() {
			// Erstellen des Verzeichnisses im temporären Verzeichnis
			err := os.MkdirAll(targetPath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			// Lesen der Datei aus dem eingebetteten Dateisystem
			data, err := binaryFiles.ReadFile(path)
			if err != nil {
				return err
			}

			// Sicherstellen, dass das Zielverzeichnis existiert
			err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
			if err != nil {
				return err
			}

			// Schreiben der Datei in das temporäre Verzeichnis
			err = os.WriteFile(targetPath, data, 0644)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Fehler beim Extrahieren der Dateien: %v", err)
	}

	return nil, nil
}
