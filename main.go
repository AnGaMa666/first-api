package main

import (
	"net/http" // Importieren des "net/http"-Pakets für HTTP-Operationen

	"github.com/gin-gonic/gin" // Importieren des "github.com/gin-gonic/gin"-Pakets für das Gin-Framework
)

func main() {
	router := gin.Default()          // Erstelle einen neuen Gin-Router
	router.GET("/albums", getAlbums) // Behandle HTTP GET-Anfragen an "/albums" mit der getAlbums-Funktion

	router.Run("localhost:8080") // Starte den Server auf localhost:8080
}

// album repräsentiert Daten über ein Musikalbum.
type album struct {
	ID     string  `json:"id"`     // Eindeutige Kennung des Albums
	Title  string  `json:"title"`  // Titel des Albums
	Artist string  `json:"artist"` // Künstler des Albums
	Price  float64 `json:"price"`  // Preis des Albums
}

// albums ist eine Liste, um Beispieldaten von Musikalben zu speichern.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums gibt die Liste aller Alben als JSON zurück.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums) // Sendet die Liste der Alben als JSON-Antwort mit dem HTTP-Statuscode 200 (OK)
}
