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
	TypeBW  = "bw"
)

func main() {
	var (
		inputPath = flag.String("path", "", "Pfad zur Datei oder zum Ordner mit Bildern")
		imgType   = flag.String("type", "web", "Typ der Optimierung: 'web', 'app' oder 'bw'")
		maxSize   = flag.Int("size", 0, "Maximale Höhe/Breite in Pixeln (optional, 0 = keine Größenänderung)")
	)
	flag.Parse()

	if *inputPath == "" {
		fmt.Fprintf(os.Stderr, "Fehler: --path ist erforderlich\n")
		flag.Usage()
		os.Exit(1)
	}

	if *imgType != TypeWeb && *imgType != TypeApp && *imgType != TypeBW {
		fmt.Fprintf(os.Stderr, "Fehler: --type muss 'web', 'app' oder 'bw' sein\n")
		os.Exit(1)
	}

	if *maxSize < 0 {
		fmt.Fprintf(os.Stderr, "Fehler: --size muss größer oder gleich 0 sein\n")
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

	// Schwarz-Weiß-Konvertierung (wenn gewünscht)
	if imgType == TypeBW {
		img = imaging.Grayscale(img)
	}

	// Größenänderung (wenn gewünscht)
	var resized image.Image = img
	if maxSize > 0 {
		bounds := img.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		// Nur verkleinern, nicht vergrößern
		if width > maxSize || height > maxSize {
			// Seitenverhältnis beibehalten
			// Für Apps: Höhere Qualität beim Resampling (CatmullRom statt Lanczos)
			if imgType == TypeApp {
				resized = imaging.Fit(img, maxSize, maxSize, imaging.CatmullRom)
			} else {
				resized = imaging.Fit(img, maxSize, maxSize, imaging.Lanczos)
			}
		}
	}

	// Ausgabedatei erstellen
	outputPath := filepath.Join(outputDir, filepath.Base(inputPath))
	
	// Datei erstellen
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Fehler beim Erstellen der Ausgabedatei: %w", err)
	}
	defer outFile.Close()

	// Format basierend auf Dateiendung und Typ
	ext := strings.ToLower(filepath.Ext(inputPath))
	
	// Für Apps: PNG bevorzugen (iOS/Android nutzen oft PNG)
	if imgType == TypeApp {
		// Konvertiere alle Formate zu PNG für Apps (bessere Qualität, Transparenz)
		outputPath = strings.TrimSuffix(outputPath, ext) + ".png"
		// Datei neu erstellen mit PNG-Endung
		outFile.Close()
		outFile, err = os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("Fehler beim Erstellen der Ausgabedatei: %w", err)
		}
		defer outFile.Close()
		err = png.Encode(outFile, resized)
	} else {
		// Für Web/BW: Originalformat beibehalten
		switch ext {
		case ".jpg", ".jpeg":
			// JPEG mit Qualität basierend auf Typ
			quality := 85
			if imgType == TypeWeb {
				quality = 85 // Web-Optimierung
			}
			err = jpeg.Encode(outFile, resized, &jpeg.Options{Quality: quality})
		case ".png":
			// PNG ohne Kompression
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
	}

	if err != nil {
		return fmt.Errorf("Fehler beim Speichern: %w", err)
	}

	return nil
}

