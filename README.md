# Image Optimizer

Ein Kommandozeilen-Tool zur Optimierung von Bildern für Web- und App-Nutzung.

## Features

- ✅ Verkleinert Bilder auf maximale Höhe/Breite (ohne Hochskalierung, optional)
- ✅ App-Optimierung: Konvertiert alle Bilder zu PNG (bessere Qualität, Transparenz)
- ✅ Konvertiert Bilder zu Schwarz-Weiß (optional)
- ✅ Entfernt alle Metadaten (EXIF, Geo-Daten, Kamera-Informationen, etc.)
- ✅ Unterstützt einzelne Dateien oder ganze Ordner
- ✅ Speichert optimierte Bilder im `optimized/` Verzeichnis
- ✅ Unterstützt verschiedene Bildformate (JPEG, PNG, GIF, BMP, WebP)

## Quick Start

```bash
# 1. Projekt klonen
git clone https://github.com/codeisconquer/Image-Optimizer.git
cd Image-Optimizer

# 2. Installieren (empfohlen - installiert in ~/bin)
./install.sh

# Oder mit npm Scripts:
npm run install:go    # Dependencies installieren
npm run build         # Kompilieren
npm run install:bin   # In ~/bin installieren

# Oder manuell:
# 2a. Dependencies installieren
go mod download

# 2b. Kompilieren
go build -o image-optimizer main.go

# 3. Ausführen
image-optimizer --path ./images --type web --size 800

# Oder Schwarz-Weiß ohne Größenänderung
image-optimizer --path ./images --type bw
```

## Installation

### Voraussetzungen

- **Go 1.21 oder höher** ([Download](https://golang.org/dl/))

Prüfe ob Go installiert ist:
```bash
go version
```

### Projekt klonen

```bash
git clone https://github.com/codeisconquer/Image-Optimizer.git
cd Image-Optimizer
```

### Dependencies installieren

```bash
go mod download
go mod tidy
```

### Kompilieren

#### Option 1: Einfacher Build (für aktuelle Plattform)

```bash
go build -o image-optimizer main.go
```

Dies erstellt eine ausführbare Datei `image-optimizer` (oder `image-optimizer.exe` auf Windows) im aktuellen Verzeichnis.

#### Option 2: Mit Makefile (empfohlen)

```bash
# Standard Build für aktuelle Plattform
make build

# Build für alle Plattformen (erstellt Binaries in dist/)
make build-all
```

#### Option 3: Cross-Compilation für verschiedene Plattformen

```bash
# macOS (Apple Silicon / M1/M2)
GOOS=darwin GOARCH=arm64 go build -o image-optimizer-darwin-arm64 main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o image-optimizer-darwin-amd64 main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o image-optimizer-linux-amd64 main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o image-optimizer-windows-amd64.exe main.go
```

### Installation systemweit (optional)

#### Option 1: Mit Install-Script (empfohlen für macOS/Linux)

Das `install.sh` Script kompiliert das Tool automatisch und installiert es in `~/bin`:

```bash
./install.sh
```

Das Script:
- Prüft ob Go installiert ist
- Installiert Dependencies
- Kompiliert das Tool
- Kopiert es nach `~/bin`
- Prüft ob `~/bin` im PATH ist (warnt falls nicht)

**Hinweis:** Falls `~/bin` nicht in deinem PATH ist, füge folgende Zeile zu deiner Shell-Konfiguration hinzu:
```bash
# Für zsh (macOS Standard)
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Für bash
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bash_profile
source ~/.bash_profile
```

#### Option 2: Mit Makefile

```bash
# Systemweite Installation (benötigt sudo)
make install
```

#### Option 3: Manuell

**macOS/Linux:**
```bash
# In ~/bin (empfohlen, kein sudo nötig)
mkdir -p ~/bin
cp image-optimizer ~/bin/

# Oder systemweit (benötigt sudo)
sudo cp image-optimizer /usr/local/bin/
```

**Windows:**
Füge das Verzeichnis mit der `.exe` Datei zu deinem PATH hinzu.

Nach der Installation kannst du das Tool von überall ausführen:
```bash
image-optimizer --path ./images --type web --size 800
```

## Verwendung

### Ausführung

**Wenn nicht systemweit installiert:**
```bash
./image-optimizer --path <pfad> [optionen]
```

**Wenn systemweit installiert:**
```bash
image-optimizer --path <pfad> [optionen]
```

### Einzelne Datei optimieren

```bash
./image-optimizer --path /pfad/zum/bild.jpg --type web --size 800
```

### Ordner mit Bildern optimieren

```bash
./image-optimizer --path /pfad/zum/ordner --type web --size 800
```

Das Tool durchsucht rekursiv alle Unterordner nach Bildern.

### Parameter

- `--path`: Pfad zur Datei oder zum Ordner (erforderlich)
- `--type`: Typ der Optimierung (`web`, `app` oder `bw`, Standard: `web`)
  - `web`: JPEG-Qualität 85%, behält Originalformat bei
  - `app`: Konvertiert alle Bilder zu PNG (bessere Qualität, Transparenz), höhere Resampling-Qualität
  - `bw`: Konvertiert Bilder zu Schwarz-Weiß
- `--size`: Maximale Höhe/Breite in Pixeln (optional, Standard: 0 = keine Größenänderung)
  - Wenn 0 oder nicht angegeben: Keine Größenänderung
  - Wenn > 0: Bilder werden auf max. Größe verkleinert (ohne Hochskalierung)

### Beispiele

```bash
# Bilder für Web auf max. 1200px optimieren
./image-optimizer --path ./images --type web --size 1200

# Einzelnes Bild für App auf max. 1600px optimieren (wird zu PNG konvertiert)
./image-optimizer --path photo.jpg --type app --size 1600

# Bilder für App optimieren (PNG, ohne Größenänderung)
./image-optimizer --path ./images --type app

# Bilder zu Schwarz-Weiß konvertieren (ohne Größenänderung)
./image-optimizer --path ./images --type bw

# Bilder zu Schwarz-Weiß konvertieren und auf 800px verkleinern
./image-optimizer --path ./images --type bw --size 800

# Nur Größe ändern (Web-Optimierung, 800px)
./image-optimizer --path ./photos --type web --size 800

# Nur Metadaten entfernen, keine Größenänderung
./image-optimizer --path ./photos --type web
```

## Ausgabe

Die optimierten Bilder werden im Verzeichnis `optimized/` gespeichert:
- Bei Einzeldateien: Im selben Verzeichnis wie die Quelldatei
- Bei Ordnern: Im angegebenen Ordner als Unterordner `optimized/`

## Unterstützte Formate

- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- BMP (.bmp)
- WebP (.webp)

## Entwicklung

### Projekt-Struktur

```
Image-Optimizer/
├── main.go          # Hauptprogramm
├── main_test.go     # Unit Tests
├── go.mod           # Go-Modul-Definition
├── go.sum           # Dependency-Checksums
├── package.json     # npm Scripts (optional)
├── Makefile         # Build-Hilfsmittel
├── install.sh       # Installations-Script (für ~/bin)
├── README.md        # Diese Datei
└── .gitignore       # Git-Ignore-Regeln
```

### Dependencies

- `github.com/disintegration/imaging` - Bildverarbeitung
- `golang.org/x/image` - Bildformate

### Build testen

**Mit npm Scripts:**
```bash
# Kompilieren
npm run build

# Hilfe anzeigen
npm run help
```

**Oder direkt:**
```bash
# Kompilieren
go build -o image-optimizer main.go

# Hilfe anzeigen
./image-optimizer --help

# Test mit Beispielbild
./image-optimizer --path test.jpg --type web --size 800
```

### npm Scripts

Das Projekt bietet npm Scripts als bequeme Wrapper für Go-Befehle:

```bash
# Tests
npm test              # Alle Tests ausführen
npm run test:verbose  # Tests mit detaillierter Ausgabe
npm run test:cover    # Tests mit Coverage-Report
npm run test:coverage # Coverage-Report als HTML generieren

# Build
npm run build         # Kompilieren
npm run build:all     # Build für alle Plattformen

# Installation
npm run install:go    # Go Dependencies installieren
npm run install:bin   # In ~/bin installieren

# Code-Qualität
npm run lint          # Code mit go vet prüfen
npm run fmt           # Code formatieren

# Sonstiges
npm run clean         # Build-Artefakte entfernen
npm run help          # Hilfe anzeigen
```

### Tests ausführen

Das Projekt enthält umfassende Unit Tests für alle wichtigen Funktionen:

**Mit Makefile (empfohlen):**
```bash
# Alle Tests ausführen
make test

# Tests mit detaillierter Ausgabe
make test-verbose

# Tests mit Coverage-Report
make test-cover

# Coverage-Report als HTML generieren
make test-coverage
```

**Oder direkt mit Go:**
```bash
# Alle Tests ausführen
go test

# Tests mit detaillierter Ausgabe
go test -v

# Tests mit Coverage-Report
go test -cover

# Coverage-Report als HTML generieren
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Getestete Funktionen:**
- `isImageFile()` - Erkennung von Bildformaten (JPEG, PNG, GIF, BMP, WebP)
- `getImageFiles()` - Rekursive Suche nach Bildern in Verzeichnissen
- `processImage()` - Bildverarbeitung und Optimierung

**Test-Coverage:**
Die Tests decken verschiedene Szenarien ab:
- Verschiedene Dateiendungen (Groß-/Kleinschreibung)
- Leere Verzeichnisse
- Nicht-existierende Pfade
- Verschiedene Bildtypen (web/app) und Größen

## Technische Details

- Metadaten werden automatisch entfernt, da die Bilder neu encodiert werden
- Seitenverhältnis wird beim Verkleinern beibehalten
- Bilder werden nur verkleinert, nie vergrößert
- Verwendet Lanczos-Resampling für hohe Qualität
- Single-Binary: Keine externen Dependencies nötig zur Laufzeit

## Lizenz

MIT

