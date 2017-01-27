package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"database/sql"
	"gopkg.in/gorp.v2"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"io"
	"os"
)

var dbMap = initDb()

var (
	Trace	*log.Logger
	Info	*log.Logger
	Warning	*log.Logger
	Error	*log.Logger
)

func main () {
	initLog(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	defer dbMap.Db.Close()
	router := gin.Default()

	router.GET("/job", getJobs)
	router.POST("/job", postJob)
	router.DELETE("/job", deleteJob)

	router.Run()
}

func getJobs (c *gin.Context) {
	var jobs []Job
	_, err := dbMap.Select(&jobs, "select * from jobs")
	checkErr(err, "Job Select Failed")

	c.JSON(200, jobs)
}

func postJob (c *gin.Context) {
	title := c.PostForm("title")
	reporter := c.PostForm("reporter")
	description := c.PostForm("description")
	assignee := c.PostForm("assignee")

	job := newJob(title, reporter, description, assignee)

	var existingJobs []Job
	_, selectErr := dbMap.Select(&existingJobs, "select * from jobs " +
		"						where title = ? " +
		"						and reporter = ? " +
		"						and description = ?",
								job.Title, job.Reporter, job.Description)
	checkErr(selectErr, "Jobs Select Failed")

	if (len(existingJobs) > 0) {
		_, err := dbMap.Update(&job)
		checkErr(err, "Job Update Failed")
	} else {
		err := dbMap.Insert(&job)
		checkErr(err, "Job Create Failed")
	}

	content := job
	c.JSON(200, content)
}

func saveJob (c *gin.Context) {
	//TODO implement saveJob
}

func deleteJob (c *gin.Context) {
	//TODO implement deleteJob
}

type Job struct {
	Id          int64
	Title       string
	Reporter    string
	Description string
	Assignee    string
}

func newJob (title, reporter, description, assignee string) Job {
	return Job {
		Title: title,
		Reporter: reporter,
		Description: description,
		Assignee: assignee,
	}
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

	// add a table, setting the table name to 'jobs' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(Job{}, "jobs").SetKeys(true, "Id")
	dbmap.AddTableWithName(Assignee{}, "assignees").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func initLog (traceHandle, infoHandle, warningHandle, errorHandle io.Writer) {
	Trace = log.New(traceHandle, "Trace: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "Warning: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func checkErr (err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}