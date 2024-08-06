package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type Album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []Album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// get albums
func getAlbums( c *gin.Context){
	 c.IndentedJSON(http.StatusOK, albums)
}

func addAlbums(c *gin.Context){
    var newAlbum Album
     if err:= c.BindJSON(&newAlbum); err!=nil{
        return
     }

     albums = append(albums, newAlbum)
     c.IndentedJSON(http.StatusCreated, newAlbum)




}

func getAlbum(c *gin.Context){
    id:= c.Param("id")
    for _, value := range albums{
        if value.ID== id{
            c.IndentedJSON(http.StatusOK, value)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message" : "album not found"})
}

func main(){
    router:= gin.Default()
    router.GET("/albums", getAlbums)
    router.POST("/albums", addAlbums)
    router.GET("/albums/:id", getAlbum)
    router.Run("localhost:8000")

}