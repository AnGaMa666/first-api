package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

const jsonFilePath = "albums.json"

type album struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Price  string `json:"price"`
}

func loadAlbums() []album {
	file, err := os.Open(jsonFilePath)
	if err != nil {
		log.Fatal("Fehler beim Öffnen der JSON-Datei:", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var albums []album
	if err := decoder.Decode(&albums); err != nil {
		log.Fatal("Fehler beim Parsen der JSON-Datei:", err)
	}
	return albums
}

func saveAlbums(albums []album) {
	file, err := os.Create(jsonFilePath)
	if err != nil {
		log.Fatal("Fehler beim Erstellen der JSON-Datei:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(albums); err != nil {
		log.Fatal("Fehler beim Schreiben der Albumdaten in die JSON-Datei:", err)
	}
}

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.GET("/albums", getAlbums)
	router.POST("/albums", createAlbum)
	router.DELETE("/albums", deleteAlbum)

	router.Run("localhost:8081")
}

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

func getAlbums(c *gin.Context) {
	albums := loadAlbums()
	c.IndentedJSON(http.StatusOK, albums)
}

func createAlbum(c *gin.Context) {
	var newAlbum album
	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	albums := loadAlbums()
	newAlbum.ID = generateUniqueID(albums)
	albums = append(albums, newAlbum)
	saveAlbums(albums)
	c.JSON(http.StatusCreated, newAlbum)
}

func deleteAlbum(c *gin.Context) {
	var deleteRequest struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&deleteRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	albums := loadAlbums()
	idIndex := make(map[string]bool)
	for _, id := range deleteRequest.IDs {
		idIndex[id] = true
	}

	for i := len(albums) - 1; i >= 0; i-- {
		if idIndex[albums[i].ID] {
			albums = append(albums[:i], albums[i+1:]...)
		}
	}
	saveAlbums(albums)

	c.JSON(http.StatusOK, gin.H{"message": "Selected items deleted"})
}

func generateUniqueID(albums []album) string {
	return strconv.Itoa(len(albums) + 1)
}
