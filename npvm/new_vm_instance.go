package npvm

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func InitVmInstance() error {
	mu.Lock()
	defer mu.Unlock()

	// Es wird geprüft wie oft ein Init durchgeführt wurde
	if totalCallInitFunctions+1 != 1 {
		return fmt.Errorf("can only run one vm instance")
	}

	// Erstellen eines temporären Verzeichnisses
	tempDir, err := os.MkdirTemp("", "embedded_files_*")
	if err != nil {
		return fmt.Errorf("fehler beim Erstellen des temporären Verzeichnisses: %w", err)
	}

	// Erstellen eines Unter-Dateisystems für 'bin'
	binFS, err := fs.Sub(binaryFiles, "bin")
	if err != nil {
		CleanUp()
		return fmt.Errorf("fehler beim Erstellen des Unter-Dateisystems 'bin': %w", err)
	}

	// Durchlaufen des eingebetteten 'bin' Dateisystems
	err = fs.WalkDir(binFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Pfad im temporären Verzeichnis
		targetPath := filepath.Join(tempDir, path)

		// Lesen der Datei aus dem eingebetteten Dateisystem
		data, err := fs.ReadFile(binFS, path) // Korrigierte Zeile
		if err != nil {
			return fmt.Errorf("fehler beim Lesen der Datei %s: %w", path, err)
		}

		// Sicherstellen, dass das Zielverzeichnis existiert
		err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("fehler beim Erstellen des Zielverzeichnisses für %s: %w", targetPath, err)
		}

		// Schreiben der Datei in das temporäre Verzeichnis
		err = os.WriteFile(targetPath, data, 0644)
		if err != nil {
			return fmt.Errorf("fehler beim Schreiben der Datei %s: %w", targetPath, err)
		}
		runningFiles = append(runningFiles, targetPath)
		return nil
	})
	if err != nil {
		CleanUp()
		return fmt.Errorf("fehler beim Extrahieren der Dateien: %w", err)
	}

	return nil
}
