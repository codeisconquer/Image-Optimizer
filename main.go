package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

const (
	TypeWeb       = "web"
	TypeApp       = "app"
	TypeBW        = "bw"
	TypeThumbnail = "thumbnail"
	TypeSepia     = "sepia"
	TypeInvert    = "invert"
)

func main() {
	var (
		inputPath = flag.String("path", "", "Pfad zur Datei oder zum Ordner mit Bildern")
		imgType   = flag.String("type", "web", "Typ der Optimierung: 'web', 'app', 'bw', 'thumbnail', 'sepia' oder 'invert'")
		maxSize   = flag.Int("size", 0, "Maximale Höhe/Breite in Pixeln (optional, 0 = keine Größenänderung)")
		outputDir = flag.String("output", "", "Ausgabeverzeichnis (optional, Standard: optimized/ im Quellverzeichnis)")
		overwrite = flag.Bool("overwrite", false, "Originaldateien überschreiben (Standard: false)")
	)
	flag.Parse()

	if *inputPath == "" {
		fmt.Fprintf(os.Stderr, "Fehler: --path ist erforderlich\n")
		flag.Usage()
		os.Exit(1)
	}

	validTypes := map[string]bool{
		TypeWeb: true, TypeApp: true, TypeBW: true,
		TypeThumbnail: true, TypeSepia: true, TypeInvert: true,
	}
	if !validTypes[*imgType] {
		fmt.Fprintf(os.Stderr, "Fehler: --type muss 'web', 'app', 'bw', 'thumbnail', 'sepia' oder 'invert' sein\n")
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

	// Ausgabeverzeichnis bestimmen
	var finalOutputDir string
	if *overwrite {
		// Bei Überschreiben: Ausgabe in dasselbe Verzeichnis wie Quelle
		if info.IsDir() {
			finalOutputDir = *inputPath
		} else {
			finalOutputDir = filepath.Dir(*inputPath)
		}
	} else if *outputDir != "" {
		// Benutzerdefiniertes Ausgabeverzeichnis
		finalOutputDir = *outputDir
	} else {
		// Standard: optimized/ im Quellverzeichnis
		if info.IsDir() {
			finalOutputDir = filepath.Join(*inputPath, "optimized")
		} else {
			finalOutputDir = filepath.Join(filepath.Dir(*inputPath), "optimized")
		}
	}

	// Ausgabeverzeichnis erstellen (nur wenn nicht überschreiben)
	if !*overwrite {
		err = os.MkdirAll(finalOutputDir, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fehler beim Erstellen des Ausgabeordners: %v\n", err)
			os.Exit(1)
		}
	}

	// Bilder verarbeiten
	successCount := 0
	for _, file := range files {
		err := processImage(file, finalOutputDir, *maxSize, *imgType, *overwrite)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fehler beim Verarbeiten von %s: %v\n", file, err)
			continue
		}
		successCount++
		if *overwrite {
			fmt.Printf("✓ Überschrieben: %s\n", file)
		} else {
			fmt.Printf("✓ Optimiert: %s\n", filepath.Base(file))
		}
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

func processImage(inputPath, outputDir string, maxSize int, imgType string, overwrite bool) error {
	// Bild öffnen
	img, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("Fehler beim Öffnen: %w", err)
	}

	// Farbkonvertierungen (wenn gewünscht)
	switch imgType {
	case TypeBW:
		img = imaging.Grayscale(img)
	case TypeSepia:
		img = applySepia(img)
	case TypeInvert:
		img = applyInvert(img)
	}

	// Größenänderung (wenn gewünscht)
	var resized image.Image = img
	finalMaxSize := maxSize
	
	// Für Thumbnails: Standardgröße 300px wenn nicht angegeben
	if imgType == TypeThumbnail && maxSize == 0 {
		finalMaxSize = 300
	}
	
	if finalMaxSize > 0 {
		bounds := img.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		// Nur verkleinern, nicht vergrößern
		if width > finalMaxSize || height > finalMaxSize {
			// Seitenverhältnis beibehalten
			// Für Apps: Höhere Qualität beim Resampling (CatmullRom statt Lanczos)
			if imgType == TypeApp {
				resized = imaging.Fit(img, finalMaxSize, finalMaxSize, imaging.CatmullRom)
			} else {
				resized = imaging.Fit(img, finalMaxSize, finalMaxSize, imaging.Lanczos)
			}
		}
	}

	// Ausgabepfad bestimmen
	var outputPath string
	if overwrite {
		// Originaldatei überschreiben
		outputPath = inputPath
	} else {
		// Neue Datei im Ausgabeverzeichnis
		outputPath = filepath.Join(outputDir, filepath.Base(inputPath))
	}
	
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
		// Für Web/BW/Thumbnail: Originalformat beibehalten
		switch ext {
		case ".jpg", ".jpeg":
			// JPEG mit Qualität basierend auf Typ
			quality := 85
			if imgType == TypeWeb {
				quality = 85 // Web-Optimierung
			} else if imgType == TypeThumbnail {
				quality = 75 // Niedrigere Qualität für Thumbnails (kleinere Dateigröße)
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

// applySepia wendet einen Sepia-Filter auf das Bild an
func applySepia(img image.Image) image.Image {
	bounds := img.Bounds()
	sepia := image.NewRGBA(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// Konvertiere von 16-bit zu 8-bit
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			
			// Sepia-Formel
			tr := float64(r8)*0.393 + float64(g8)*0.769 + float64(b8)*0.189
			tg := float64(r8)*0.349 + float64(g8)*0.686 + float64(b8)*0.168
			tb := float64(r8)*0.272 + float64(g8)*0.534 + float64(b8)*0.131
			
			// Begrenze auf 255
			if tr > 255 {
				tr = 255
			}
			if tg > 255 {
				tg = 255
			}
			if tb > 255 {
				tb = 255
			}
			
			sepia.Set(x, y, color.RGBA{
				R: uint8(tr),
				G: uint8(tg),
				B: uint8(tb),
				A: uint8(a >> 8),
			})
		}
	}
	return sepia
}

// applyInvert invertiert die Farben des Bildes
func applyInvert(img image.Image) image.Image {
	bounds := img.Bounds()
	inverted := image.NewRGBA(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// Konvertiere von 16-bit zu 8-bit und invertiere
			inverted.Set(x, y, color.RGBA{
				R: 255 - uint8(r>>8),
				G: 255 - uint8(g>>8),
				B: 255 - uint8(b>>8),
				A: uint8(a >> 8),
			})
		}
	}
	return inverted
}

