package main

import (
	"fmt"
	//"github.com/ziutek/mymysql/mysql"
	"encoding/json"
	"log"
	"net/http"

	//"time"
	//_ "github.com/ziutek/mymysql/native" // Native engine
	//_ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Video struct {
	Title      string `json:title`
	Url        string `json:url`
	Poster     string `json:poster`
	DatePosted string `json:datePosted`
	CatVidId   int64  `json:catVidId`
	UpMeows    uint   `json:upMeows`
	DownMeows  uint   `json:downMeows`
}

func main() {
	//The db.sql object is meant to be long lived. It does not create a connection to the source
	//sql.Open create a connection to the db. instead it only prepares database abstraction for later use
	db, err := sql.Open("mysql", "jflewis:jflewis2015!@tcp(mysqlcs.millersville.edu:3306)/jflewis")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("database connected to succesfully")
	}
	r := mux.NewRouter()

	r.HandleFunc("/", getRandVid(db))
	s := r.Host("www.api.localhost.com").Subrouter()
	s.HandleFunc("/randomVideo", getRandVid(db))

	//r.HandleFunc("/addVideo/", f)
	//addVideo(db)
	http.ListenAndServe(":8080", r)

}

func getRandVid(db *sql.DB) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var video Video
		err := db.QueryRow(`SELECT CatVid.title,CatVid.url,CatVid.video_poster, 
		CatVid.date_posted ,CatVid.catVidID, Vote.upmeows, Vote.downmeows FROM CatVid,
		 Vote where Vote.catVidID = CatVid.catVidID ORDER BY RAND() LIMIT 1`).Scan(&video.Title, &video.Url, &video.Poster, &video.DatePosted, &video.CatVidId, &video.UpMeows, &video.DownMeows)
		if err != nil {
			log.Fatal(err)
		}
		js, err := json.Marshal(video)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("Access-Control-Allow-Origin", "*")

		rw.Write(js)
	}
}

/*
func addVideo(db *sql.DB) http.Handler {
	return http.HandleFunc(func(rw http.ResponseWriter, req *http.Request) {
		//var tag string
		_, err := db.Exec("call UploadVideo(?,?,?,?)", "a new title", "some fake url", "this,thing,may,work,", "username")
		if err != nil {
			log.Fatal(err)
		}
	})
}
*/
