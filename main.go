package main

import (
    "fmt"
    "log"
    "database/sql"
    //"encoding/json"
    "os"
    "net/http"
    _ "github.com/lib/pq"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
    _ "github.com/google/go-github/v56/github"
)

func main() {
    connectionName := "cs588assignment5-406403:us-central1:mypostgres"
	dbUser := "postgres"
	dbPass := "root"
	dbName := "cs588_assignment_5"

    dbURI := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		connectionName, dbName, dbUser, dbPass)

    log.Println("Initializing database connection")
	db, err := sql.Open("cloudsqlpostgres", dbURI)
	if err != nil {
		log.Fatalf("Error on initializing database connection: %s", err.Error())
	}
	defer db.Close()

	//Test the database connection
	log.Println("Testing database connection")
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error on database connection: %s", err.Error())
	}
	log.Println("Database connection established")

	log.Println("Database query done!")

	port := os.Getenv("PORT")
	if port == "" {
        port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, world!"))
    })
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	}()

    //collectGithubData(db)


}

// func collectGithubData (db *sql.DB) error {

//     client := github.NewClient(nil).WithAuthToken("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")

//     drop_table := `drop table if exists prometheus_issues`

// 	_, err := db.Exec(drop_table)
// 	if err != nil {
// 		panic(err)
// 	}

//     create_table := `CREATE TABLE IF NOT EXISTS "prometheus_issues" `

// }


// func collectStackOverflowData (db *sql.DB) error{

// }