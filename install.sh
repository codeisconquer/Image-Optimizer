#!/bin/bash

# Image Optimizer Install Script
# Installiert das Tool in ~/bin für systemweite Nutzung

set -e

echo "🚀 Image Optimizer Installation"
echo "================================"
echo ""

# Prüfe ob Go installiert ist
if ! command -v go &> /dev/null; then
    echo "❌ Fehler: Go ist nicht installiert!"
    echo "   Bitte installiere Go von https://golang.org/dl/"
    exit 1
fi

echo "✓ Go gefunden: $(go version)"
echo ""

# Prüfe ob wir im richtigen Verzeichnis sind
if [ ! -f "main.go" ]; then
    echo "❌ Fehler: main.go nicht gefunden!"
    echo "   Bitte führe das Script im Projektverzeichnis aus."
    exit 1
fi

# Dependencies installieren
echo "📦 Installiere Dependencies..."
go mod download
go mod tidy
echo "✓ Dependencies installiert"
echo ""

# Kompiliere das Projekt
echo "🔨 Kompiliere Image Optimizer..."
if go build -o image-optimizer main.go; then
    echo "✓ Kompilierung erfolgreich"
else
    echo "❌ Fehler beim Kompilieren!"
    exit 1
fi
echo ""

# Erstelle ~/bin Verzeichnis falls es nicht existiert
BIN_DIR="$HOME/bin"
if [ ! -d "$BIN_DIR" ]; then
    echo "📁 Erstelle ~/bin Verzeichnis..."
    mkdir -p "$BIN_DIR"
    echo "✓ Verzeichnis erstellt: $BIN_DIR"
    echo ""
fi

# Kopiere die ausführbare Datei
echo "📋 Kopiere image-optimizer nach $BIN_DIR..."
cp image-optimizer "$BIN_DIR/"
chmod +x "$BIN_DIR/image-optimizer"
echo "✓ Installation abgeschlossen!"
echo ""

# Prüfe ob ~/bin im PATH ist
if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
    echo "⚠️  Wichtig: ~/bin ist nicht in deinem PATH!"
    echo ""
    echo "Füge folgende Zeile zu deiner Shell-Konfiguration hinzu:"
    echo ""
    if [ -f "$HOME/.zshrc" ]; then
        echo "  echo 'export PATH=\"\$HOME/bin:\$PATH\"' >> ~/.zshrc"
        echo "  source ~/.zshrc"
    elif [ -f "$HOME/.bash_profile" ]; then
        echo "  echo 'export PATH=\"\$HOME/bin:\$PATH\"' >> ~/.bash_profile"
        echo "  source ~/.bash_profile"
    else
        echo "  echo 'export PATH=\"\$HOME/bin:\$PATH\"' >> ~/.bashrc"
        echo "  source ~/.bashrc"
    fi
    echo ""
else
    echo "✓ ~/bin ist bereits im PATH"
    echo ""
fi

echo "✅ Installation erfolgreich!"
echo ""
echo "Du kannst das Tool jetzt von überall ausführen:"
echo "  image-optimizer --path ./images --type web --size 800"
echo ""

