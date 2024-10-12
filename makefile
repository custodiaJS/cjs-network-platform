# Makefile

# Standard-Variablen
BINARY_NAME=cjsnpvm
BUILD_DIR=build
NPVM_DIR=npvm/qemu
NPVM_BIN_DIR=npvm/bin

# Gemeinsame Konfigurationsoptionen
CONFIGURE_OPTS=--prefix=/usr/local --disable-spice --disable-sdl --disable-gtk --disable-vnc --enable-virtfs

# Erkennung des Betriebssystems
OS := $(shell uname -s)
ifeq ($(OS),Darwin)
    OS := darwin
endif
ifeq ($(OS),Linux)
    OS := linux
endif
ifeq ($(OS),Windows_NT)
    OS := windows
endif

# Erkennung der Architektur
ARCH := $(shell uname -m)
ifeq ($(ARCH),x86_64)
    ARCH := amd64
endif
ifeq ($(ARCH),arm64)
    ARCH := arm64
endif

# Setze den Namen der ausführbaren Datei basierend auf OS
ifeq ($(OS),windows)
    EXECUTABLE=$(BINARY_NAME).exe
else
    EXECUTABLE=$(BINARY_NAME)
endif

# Setze den Namen der QEMU-Binärdatei basierend auf OS und ARCH
ifeq ($(OS),darwin)
    ifeq ($(ARCH),arm64)
        QEMU_BINARY=qemu-system-aarch64
    endif
    ifeq ($(ARCH),amd64)
        QEMU_BINARY=qemu-system-x86_64
    endif
endif

ifeq ($(OS),linux)
    ifeq ($(ARCH),arm64)
        QEMU_BINARY=qemu-system-aarch64
    endif
    ifeq ($(ARCH),amd64)
        QEMU_BINARY=qemu-system-x86_64
    endif
endif

ifeq ($(OS),windows)
    ifeq ($(ARCH),arm64)
        QEMU_BINARY=qemu-system-aarch64.exe
    endif
    ifeq ($(ARCH),amd64)
        QEMU_BINARY=qemu-system-x86_64.exe
    endif
endif

# Ziel: Standard build (generiert und baut)
.PHONY: all
all: build go

# Ziel: Baue das Projekt
.PHONY: build
build:
	@echo "Starting build for $(OS) $(ARCH)..."
	@mkdir -p $(BUILD_DIR) $(NPVM_BIN_DIR)

	# Navigiere zum NPVM-Verzeichnis und führe die Build-Schritte aus
	@cd $(NPVM_DIR) && \
	./configure $(CONFIGURE_OPTS) --target-list=aarch64-softmmu && \
	make -j$(shell sysctl -n hw.ncpu 2>/dev/null || echo 4)
	@cd ..

	# Verschiebe das erstellte Binary in das Ausgabe-Verzeichnis
	@echo "Moving QEMU binary..."
	@if [ "$(OS)" = "windows" ]; then \
		mv $(NPVM_DIR)/build/$(QEMU_BINARY) $(NPVM_BIN_DIR)/npvm-qemu.exe; \
	else \
		mv $(NPVM_DIR)/build/$(QEMU_BINARY) $(NPVM_BIN_DIR)/npvm-qemu; \
	fi

	# Führe go build aus, um das Go-Programm zu kompilieren
	@echo "Compiling Go program for $(OS) $(ARCH)..."
	@cd proc && \
	GOOS=$(OS) GOARCH=$(ARCH) go build -o $(EXECUTABLE) main.go parms.go

	# Verschiebe das ausführbare Programm in das Build-Verzeichnis
	@mv proc/$(EXECUTABLE) $(BUILD_DIR)/

	@echo "Build completed for $(OS) $(ARCH). Executables are in $(BUILD_DIR)/"
# Ziel: Bereinige generierte Dateien und Build-Ordner
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	go clean
	rm -rf $(BUILD_DIR) $(EXECUTABLE) $(NPVM_BIN_DIR)
	@echo "Clean completed."
# Erstellt nur die Go Datei neu
.PHONY: go
go:
	@echo "Compiling Go program for $(OS) $(ARCH)..."
	@cd proc && GOOS=$(OS) GOARCH=$(ARCH) go build -o $(EXECUTABLE) main.go parms.go
	@mv proc/$(EXECUTABLE) $(BUILD_DIR)/
	@chmod +x $(BUILD_DIR)/$(EXECUTABLE)
