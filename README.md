# Image Optimizer

Ein Kommandozeilen-Tool zur Optimierung von Bildern für Web- und App-Nutzung.

## Features

- ✅ Verkleinert Bilder auf maximale Höhe/Breite (ohne Hochskalierung)
- ✅ Entfernt alle Metadaten (EXIF, Geo-Daten, Kamera-Informationen, etc.)
- ✅ Unterstützt einzelne Dateien oder ganze Ordner
- ✅ Speichert optimierte Bilder im `optimized/` Verzeichnis
- ✅ Unterstützt verschiedene Bildformate (JPEG, PNG, GIF, BMP, WebP)

## Installation

### Voraussetzungen

- Go 1.21 oder höher

### Build

```bash
go build -o image-optimizer main.go
```

Oder für verschiedene Plattformen:

```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o image-optimizer main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o image-optimizer main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o image-optimizer main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o image-optimizer.exe main.go
```

## Verwendung

### Einzelne Datei optimieren

```bash
./image-optimizer --path /pfad/zum/bild.jpg --type web --size 800
```

### Ordner mit Bildern optimieren

```bash
./image-optimizer --path /pfad/zum/ordner --type web --size 800
```

### Parameter

- `--path`: Pfad zur Datei oder zum Ordner (erforderlich)
- `--type`: Typ der Optimierung (`web` oder `app`, Standard: `web`)
  - `web`: JPEG-Qualität 85%
  - `app`: JPEG-Qualität 90%
- `--size`: Maximale Höhe/Breite in Pixeln (Standard: 800)

### Beispiele

```bash
# Bilder für Web auf max. 1200px optimieren
./image-optimizer --path ./images --type web --size 1200

# Einzelnes Bild für App auf max. 1600px optimieren
./image-optimizer --path photo.jpg --type app --size 1600

# Standard-Einstellungen (web, 800px)
./image-optimizer --path ./photos
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

## Technische Details

- Metadaten werden automatisch entfernt, da die Bilder neu encodiert werden
- Seitenverhältnis wird beim Verkleinern beibehalten
- Bilder werden nur verkleinert, nie vergrößert
- Verwendet Lanczos-Resampling für hohe Qualität

## Lizenz

MIT

