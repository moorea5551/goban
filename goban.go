package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/job", getJobs)

	router.Run()
}

func getJobs(c *gin.Context) {
	c.String(http.StatusOK, "Hello world")
}
