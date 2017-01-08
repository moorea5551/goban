package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"database/sql"
	"gopkg.in/gorp.v2"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main () {
	dbMap  := initDb()
	defer dbMap.Db.Close()
	router := gin.Default()

	router.GET("/job", getJobs)
	router.PUT("/job", createJob)
	router.POST("/job", saveJob)
	router.DELETE("/job", deleteJob)

	router.Run()
}

func getJobs (c *gin.Context) {
	c.String(http.StatusOK, "Hello world")
}

func createJob (c *gin.Context) {
	//TODO implement createJob
}

func saveJob (c *gin.Context) {
	//TODO implement saveJob
}

func deleteJob (c *gin.Context) {
	//TODO implement deleteJob
}

type Job struct {
	Id          int64
	WorkOrder   int64
	Title       string
	Reporter    string
	Location    string
	Description string
	Assignee    int64
}

type Assignee struct {
	Id        int64
	FirstName string
	LastName  string
}

func initDb () *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("sqlite3", "./jobsData.bin")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(Job{}, "jobs").SetKeys(true, "Id")
	dbmap.AddTableWithName(Assignee{}, "assignees").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr (err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}