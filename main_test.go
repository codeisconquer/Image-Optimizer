package main

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/disintegration/imaging"
)

func TestIsImageFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"JPEG lowercase", "test.jpg", true},
		{"JPEG uppercase", "test.JPG", true},
		{"JPEG mixed case", "test.JpG", true},
		{"JPEG extension", "test.jpeg", true},
		{"PNG", "test.png", true},
		{"GIF", "test.gif", true},
		{"BMP", "test.bmp", true},
		{"WebP", "test.webp", true},
		{"Not an image", "test.txt", false},
		{"No extension", "test", false},
		{"PDF", "test.pdf", false},
		{"Empty string", "", false},
		{"Path with directory", "/path/to/image.jpg", true},
		{"Path with directory PNG", "/path/to/image.PNG", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isImageFile(tt.path)
			if result != tt.expected {
				t.Errorf("isImageFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetImageFiles(t *testing.T) {
	// Erstelle temporäres Verzeichnis
	tmpDir, err := os.MkdirTemp("", "image-optimizer-test-*")
	if err != nil {
		t.Fatalf("Fehler beim Erstellen des temporären Verzeichnisses: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Erstelle Test-Dateien
	testFiles := []struct {
		name     string
		isImage  bool
		content  string
	}{
		{"image1.jpg", true, "fake jpeg"},
		{"image2.png", true, "fake png"},
		{"document.txt", false, "text file"},
		{"image3.JPEG", true, "fake jpeg"},
		{"image4.gif", true, "fake gif"},
		{"readme.md", false, "markdown"},
		{"image5.bmp", true, "fake bmp"},
		{"image6.webp", true, "fake webp"},
	}

	var expectedImages []string
	for _, tf := range testFiles {
		filePath := filepath.Join(tmpDir, tf.name)
		err := os.WriteFile(filePath, []byte(tf.content), 0644)
		if err != nil {
			t.Fatalf("Fehler beim Erstellen der Test-Datei %s: %v", tf.name, err)
		}
		if tf.isImage {
			expectedImages = append(expectedImages, filePath)
		}
	}

	// Erstelle Unterordner mit Bildern
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Fehler beim Erstellen des Unterordners: %v", err)
	}

	subImagePath := filepath.Join(subDir, "subimage.jpg")
	err = os.WriteFile(subImagePath, []byte("fake jpeg"), 0644)
	if err != nil {
		t.Fatalf("Fehler beim Erstellen der Unterordner-Datei: %v", err)
	}
	expectedImages = append(expectedImages, subImagePath)

	// Test getImageFiles
	result, err := getImageFiles(tmpDir)
	if err != nil {
		t.Fatalf("getImageFiles() returned error: %v", err)
	}

	if len(result) != len(expectedImages) {
		t.Errorf("getImageFiles() found %d images, want %d", len(result), len(expectedImages))
	}

	// Prüfe ob alle erwarteten Bilder gefunden wurden
	found := make(map[string]bool)
	for _, img := range result {
		found[img] = true
	}

	for _, expected := range expectedImages {
		if !found[expected] {
			t.Errorf("getImageFiles() did not find expected image: %s", expected)
		}
	}
}

func TestGetImageFilesEmptyDir(t *testing.T) {
	// Erstelle leeres temporäres Verzeichnis
	tmpDir, err := os.MkdirTemp("", "image-optimizer-test-empty-*")
	if err != nil {
		t.Fatalf("Fehler beim Erstellen des temporären Verzeichnisses: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	result, err := getImageFiles(tmpDir)
	if err != nil {
		t.Fatalf("getImageFiles() returned error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("getImageFiles() found %d images in empty directory, want 0", len(result))
	}
}

func TestGetImageFilesNonExistentDir(t *testing.T) {
	nonExistentDir := "/non/existent/directory/that/does/not/exist"
	_, err := getImageFiles(nonExistentDir)
	if err == nil {
		t.Error("getImageFiles() should return error for non-existent directory")
	}
}

func TestProcessImage(t *testing.T) {
	// Erstelle temporäres Verzeichnis
	tmpDir, err := os.MkdirTemp("", "image-optimizer-test-*")
	if err != nil {
		t.Fatalf("Fehler beim Erstellen des temporären Verzeichnisses: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Erstelle ein einfaches Testbild (1x1 PNG)
	testImagePath := filepath.Join(tmpDir, "test.png")
	err = createTestImage(testImagePath)
	if err != nil {
		t.Fatalf("Fehler beim Erstellen des Testbildes: %v", err)
	}

	outputDir := filepath.Join(tmpDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		t.Fatalf("Fehler beim Erstellen des Ausgabeordners: %v", err)
	}

	// Test processImage mit verschiedenen Parametern
	tests := []struct {
		name      string
		maxSize   int
		imgType   string
		overwrite bool
	}{
		{"Web type, size 800", 800, TypeWeb, false},
		{"App type, size 800", 800, TypeApp, false},
		{"BW type, size 800", 800, TypeBW, false},
		{"Web type, size 100", 100, TypeWeb, false},
		{"App type, size 100", 100, TypeApp, false},
		{"BW type, size 100", 100, TypeBW, false},
		{"BW type, no resize", 0, TypeBW, false},
		{"Web type, no resize", 0, TypeWeb, false},
		{"App type, no resize", 0, TypeApp, false},
		{"Thumbnail type, default size", 0, TypeThumbnail, false},
		{"Thumbnail type, custom size", 200, TypeThumbnail, false},
		{"App type, overwrite", 800, TypeApp, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Für Overwrite-Test: Erstelle neue Testdatei
			testPath := testImagePath
			if tt.overwrite {
				testPath = filepath.Join(tmpDir, "test_overwrite.png")
				err = createTestImage(testPath)
				if err != nil {
					t.Fatalf("Fehler beim Erstellen des Testbildes: %v", err)
				}
			}

			err := processImage(testPath, outputDir, tt.maxSize, tt.imgType, tt.overwrite)
			if err != nil {
				t.Errorf("processImage() returned error: %v", err)
			}

			// Prüfe ob Ausgabedatei erstellt wurde
			var outputPath string
			if tt.overwrite {
				outputPath = testPath
			} else {
				outputPath = filepath.Join(outputDir, "test.png")
				if tt.name == "App type, overwrite" {
					outputPath = filepath.Join(outputDir, "test_overwrite.png")
				}
			}

			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				t.Errorf("processImage() did not create output file: %s", outputPath)
			}
		})
	}
}

func TestProcessImageNonExistentFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "image-optimizer-test-*")
	if err != nil {
		t.Fatalf("Fehler beim Erstellen des temporären Verzeichnisses: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	outputDir := filepath.Join(tmpDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		t.Fatalf("Fehler beim Erstellen des Ausgabeordners: %v", err)
	}

	nonExistentFile := filepath.Join(tmpDir, "nonexistent.jpg")
	err = processImage(nonExistentFile, outputDir, 800, TypeWeb, false)
	if err == nil {
		t.Error("processImage() should return error for non-existent file")
	}
}

// createTestImage erstellt ein einfaches 10x10 PNG-Bild für Tests
func createTestImage(path string) error {
	// Erstelle ein einfaches 10x10 rotes Bild
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	red := color.RGBA{255, 0, 0, 255}
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, red)
		}
	}
	return imaging.Save(img, path)
}
