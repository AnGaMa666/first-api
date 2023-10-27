package main

import (
	"encoding/json" // Importiert das "encoding/json"-Paket, um JSON-Verarbeitung zu ermöglichen.
	"log"          // Importiert das "log"-Paket für das Protokollieren von Fehlern.
	"net/http"      // Importiert das "net/http"-Paket für HTTP-Operationen.
	"os"           // Importiert das "os"-Paket für Dateioperationen.
	"strconv"      // Importiert das "strconv"-Paket zur Konvertierung von Zeichenfolgen in Zahlen.
	"github.com/gin-gonic/gin" // Importiert das "gin"-Paket, um einen HTTP-Server zu erstellen.
)

const jsonFilePath = "albums.json" // Die Konstante jsonFilePath definiert den Dateipfad für die JSON-Datei.

type album struct {
	ID     string `json:"id"`     // Die Struktur "album" repräsentiert ein Musikalbum mit ID, Titel, Künstler und Preis.
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Price  string `json:"price"`
}

// loadAlbums lädt die Albumdaten aus der JSON-Datei.
func loadAlbums() []album {
	file, err := os.Open(jsonFilePath) // Öffnet die JSON-Datei zum Lesen.
	if err != nil {
		log.Fatal("Fehler beim Öffnen der JSON-Datei:", err) // Protokolliert einen Fehler, falls beim Öffnen ein Fehler auftritt.
	}
	defer file.Close() // Stellt sicher, dass die Datei nach Verlassen der Funktion geschlossen wird.

	decoder := json.NewDecoder(file) // Erstellt einen JSON-Decoder, um die Datei zu parsen.
	var albums []album // Erstellt ein Slice, um die geladenen Alben zu speichern.
	if err := decoder.Decode(&albums); err != nil {
		log.Fatal("Fehler beim Parsen der JSON-Datei:", err) // Protokolliert einen Fehler, falls beim Parsen ein Fehler auftritt.
	}
	return albums // Gibt das Slice der geladenen Alben zurück.
}

// saveAlbums speichert die Albumdaten in der JSON-Datei.
func saveAlbums(albums []album) {
	file, err := os.Create(jsonFilePath) // Erstellt die JSON-Datei zum Schreiben.
	if err != nil {
		log.Fatal("Fehler beim Erstellen der JSON-Datei:", err) // Protokolliert einen Fehler, falls beim Erstellen ein Fehler auftritt.
	}
	defer file.Close() // Stellt sicher, dass die Datei nach Verlassen der Funktion geschlossen wird.

	encoder := json.NewEncoder(file) // Erstellt einen JSON-Encoder, um die Daten in die Datei zu schreiben.
	if err := encoder.Encode(albums); err != nil {
		log.Fatal("Fehler beim Schreiben der Albumdaten in die JSON-Datei:", err) // Protokolliert einen Fehler, falls beim Schreiben ein Fehler auftritt.
	}
}

func main() {
	router := gin.Default() // Erstellt einen neuen HTTP-Router mit den Standardeinstellungen.
	router.Use(CORSMiddleware()) // Fügt dem Router ein CORS-Middleware hinzu, um Cross-Origin-Anfragen zu verarbeiten.

	router.GET("/albums", getAlbums) // Definiert einen Endpunkt für das Abrufen von Albumdaten.
	router.POST("/albums", createAlbum) // Definiert einen Endpunkt zum Erstellen eines neuen Albums.
	router.DELETE("/albums", deleteAlbum) // Definiert einen Endpunkt zum Löschen von ausgewählten Alben.

	router.Run("localhost:8081") // Startet den HTTP-Server auf Port 8081 auf dem lokalen Host.
}

// CORSMiddleware fügt CORS-Header zu den Anfragen hinzu, um Cross-Origin-Anfragen zu unterstützen.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Max-Age", "600")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// getAlbums gibt die Liste der Alben als JSON zurück.
func getAlbums(c *gin.Context) {
	albums := loadAlbums() // Lädt die Albumdaten aus der JSON-Datei.
	c.IndentedJSON(http.StatusOK, albums) // Sendet die Albumdaten als JSON-Antwort zurück.
}

// createAlbum erstellt ein neues Album und fügt es zur Liste hinzu.
func createAlbum(c *gin.Context) {
	var newAlbum album
	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Sendet eine Fehlerantwort, wenn die Eingabe ungültig ist.
		return
	}

	albums := loadAlbums() // Lädt die vorhandenen Alben.
	newAlbum.ID = generateUniqueID(albums) // Generiert eine eindeutige ID für das neue Album.
	albums = append(albums, newAlbum) // Fügt das neue Album zur Liste hinzu.
	saveAlbums(albums) // Speichert die aktualisierte Liste der Alben.
	c.JSON(http.StatusCreated, newAlbum) // Sendet eine Bestätigungsantwort zurück.
}

// deleteAlbum löscht ausgewählte Alben aus der Liste.
func deleteAlbum(c *gin.Context) {
	var deleteRequest struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&deleteRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Sendet eine Fehlerantwort, wenn die Eingabe ungültig ist.
		return
	}

	albums := loadAlbums() // Lädt die vorhandenen Alben.
	idIndex := make(map[string]bool)
	for _, id := range deleteRequest.IDs {
		idIndex[id] = true
	}

	for i := len(albums) - 1; i >= 0; i-- {
		if idIndex[albums[i].ID] {
			albums = append(albums[:i], albums[i+1:]...) // Löscht die ausgewählten Alben aus der Liste.
		}
	}
	saveAlbums(albums) // Speichert die aktualisierte Liste der Alben.
	c.JSON(http.StatusOK, gin.H{"message": "Selected items deleted"}) // Sendet eine Bestätigungsantwort zurück.
}

// generateUniqueID generiert eine eindeutige ID für ein neues Album.
func generateUniqueID(albums []album) string {
	return strconv.Itoa(len(albums) + 1)
}
