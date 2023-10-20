// Dies ist die Hauptanweisung, die das Go-Programm startet.
package main

// Wir importieren Funktionen, um mit dem Internet und HTTP zu arbeiten.
import (
	"net/http"

	"github.com/gin-gonic/gin" // Wir importieren das Gin-Framework, um Webseiten zu erstellen und zu bedienen.
)

func main() {
	// Hier erstellen wir einen neuen Webserver.
	router := gin.Default()

	// Hier fügen wir eine spezielle Regel hinzu, die es dem Server erlaubt, Anfragen von verschiedenen Orten zu akzeptieren.
	// Dies ist wichtig, wenn Sie Anfragen von einer anderen Webseite erhalten.
	router.Use(CORSMiddleware())

	// Hier definieren wir, dass der Server auf Anfragen an die Adresse "/albums" reagiert.
	// Wenn jemand nach "/albums" fragt, zeigen wir ihm die Liste der Alben an.
	router.GET("/albums", getAlbums)

	// Hier sagen wir dem Server, dass er auf Anfragen an "/albums" reagieren soll, aber dieses Mal
	// um ein neues Album zu erstellen.
	router.POST("/albums", createAlbum)

	// Hier starten wir den Server auf dem Computer, damit er auf Anfragen von anderen Computern hören kann.
	router.Run("localhost:8081")
}

// Diese Funktion konfiguriert die CORS-Regeln, die den Zugriff von verschiedenen Orten erlauben.
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

// Hier definieren wir die Struktur eines Musikalbums.
type album struct {
	ID     string  `json:"id"`     // Jedes Album hat eine eindeutige Kennung (ID).
	Title  string  `json:"title"`  // Es hat auch einen Titel.
	Artist string  `json:"artist"` // Der Name des Künstlers oder der Band.
	Price  float64 `json:"price"`  // Und den Preis des Albums.
}

// Hier haben wir Beispieldaten für einige Alben.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// Hier sagen wir dem Server, wie er auf Anfragen nach der Liste der Alben reagieren soll.
func getAlbums(c *gin.Context) {
	// Wir senden die Liste der Alben zurück, wenn jemand danach fragt.
	c.IndentedJSON(http.StatusOK, albums)
}

// Hier zeigen wir dem Server, wie er auf Anfragen zum Erstellen eines neuen Albums reagieren soll.
func createAlbum(c *gin.Context) {
	// Wir nehmen die Daten, die uns jemand gesendet hat, um ein neues Album zu erstellen.
	// Wenn es Fehler gibt, teilen wir dem Sender das mit.
	var newAlbum album
	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Wenn alles in Ordnung ist, fügen wir das neue Album zur Liste der Alben hinzu.
	albums = append(albums, newAlbum)

	// Und wir sagen dem Sender, dass das Album erfolgreich erstellt wurde.
	c.JSON(http.StatusCreated, newAlbum)
}
