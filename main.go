package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

const (
	TypeWeb = "web"
	TypeApp = "app"
)

func main() {
	var (
		inputPath = flag.String("path", "", "Pfad zur Datei oder zum Ordner mit Bildern")
		imgType   = flag.String("type", "web", "Typ der Optimierung: 'web' oder 'app'")
		maxSize   = flag.Int("size", 800, "Maximale Höhe/Breite in Pixeln (nicht hochskalieren)")
	)
	flag.Parse()

	if *inputPath == "" {
		fmt.Fprintf(os.Stderr, "Fehler: --path ist erforderlich\n")
		flag.Usage()
		os.Exit(1)
	}

	if *imgType != TypeWeb && *imgType != TypeApp {
		fmt.Fprintf(os.Stderr, "Fehler: --type muss 'web' oder 'app' sein\n")
		os.Exit(1)
	}

	if *maxSize <= 0 {
		fmt.Fprintf(os.Stderr, "Fehler: --size muss größer als 0 sein\n")
		os.Exit(1)
	}

	// Prüfe ob Pfad existiert
	info, err := os.Stat(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fehler: Pfad nicht gefunden: %v\n", err)
		os.Exit(1)
	}

	var files []string
	if info.IsDir() {
		// Ordner verarbeiten
		files, err = getImageFiles(*inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fehler beim Lesen des Ordners: %v\n", err)
			os.Exit(1)
		}
		if len(files) == 0 {
			fmt.Println("Keine Bilder im angegebenen Ordner gefunden.")
			os.Exit(0)
		}
	} else {
		// Einzelne Datei verarbeiten
		if !isImageFile(*inputPath) {
			fmt.Fprintf(os.Stderr, "Fehler: Datei ist kein unterstütztes Bildformat\n")
			os.Exit(1)
		}
		files = []string{*inputPath}
	}

	// Optimiertes Verzeichnis erstellen
	var outputDir string
	if info.IsDir() {
		// Bei Ordnern: optimized/ im angegebenen Ordner
		outputDir = filepath.Join(*inputPath, "optimized")
	} else {
		// Bei Einzeldateien: optimized/ im selben Ordner wie die Datei
		outputDir = filepath.Join(filepath.Dir(*inputPath), "optimized")
	}

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fehler beim Erstellen des Ausgabeordners: %v\n", err)
		os.Exit(1)
	}

	// Bilder verarbeiten
	successCount := 0
	for _, file := range files {
		err := processImage(file, outputDir, *maxSize, *imgType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fehler beim Verarbeiten von %s: %v\n", file, err)
			continue
		}
		successCount++
		fmt.Printf("✓ Optimiert: %s\n", filepath.Base(file))
	}

	fmt.Printf("\nFertig! %d von %d Bildern erfolgreich optimiert.\n", successCount, len(files))
}

func getImageFiles(dir string) ([]string, error) {
	var files []string
	extensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if extensions[ext] {
				files = append(files, path)
			}
		}
		return nil
	})

	return files, err
}

func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	extensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}
	return extensions[ext]
}

func processImage(inputPath, outputDir string, maxSize int, imgType string) error {
	// Bild öffnen
	img, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("Fehler beim Öffnen: %w", err)
	}

	// Originalgröße
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Nur verkleinern, nicht vergrößern
	var resized image.Image = img
	if width > maxSize || height > maxSize {
		// Seitenverhältnis beibehalten
		resized = imaging.Fit(img, maxSize, maxSize, imaging.Lanczos)
	}

	// Ausgabedatei erstellen
	outputPath := filepath.Join(outputDir, filepath.Base(inputPath))
	
	// Datei erstellen
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Fehler beim Erstellen der Ausgabedatei: %w", err)
	}
	defer outFile.Close()

	// Format basierend auf Dateiendung
	ext := strings.ToLower(filepath.Ext(inputPath))
	switch ext {
	case ".jpg", ".jpeg":
		// JPEG mit hoher Qualität für Web/App
		quality := 85
		if imgType == TypeApp {
			quality = 90 // Höhere Qualität für Apps
		}
		err = jpeg.Encode(outFile, resized, &jpeg.Options{Quality: quality})
	case ".png":
		// PNG ohne Kompression (kann später optimiert werden)
		err = png.Encode(outFile, resized)
	default:
		// Für andere Formate als PNG speichern
		err = png.Encode(outFile, resized)
		if err == nil {
			// Dateiname anpassen
			newPath := strings.TrimSuffix(outputPath, ext) + ".png"
			os.Rename(outputPath, newPath)
			outputPath = newPath
		}
	}

	if err != nil {
		return fmt.Errorf("Fehler beim Speichern: %w", err)
	}

	return nil
}

