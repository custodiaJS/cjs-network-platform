package npvm

import (
	"embed"
	"sync"
)

//go:embed bin/***
var binaryFiles embed.FS

// Stellt den Programmweiten Mutex dar
var mu *sync.Mutex = new(sync.Mutex)

// Gibt an, wie oft die Init Funktion aufgerufen wurde
var totalCallInitFunctions uint8 = 0

// Speichert die Adresse dem Temp dirs ab
var tempdirPath string

// Speichert alle Running Files ab
var runningFiles []string
