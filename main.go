package main

import (
    "fmt"
    "log"
    "database/sql"
    "encoding/json"
    "os"
    "net/http"
    //"context"
    _ "github.com/lib/pq"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
    //"github.com/google/go-github/v56/github"
    "io/ioutil"
    
)

type GitHubIssueStruct []struct {
	Title                      string `json:"title"`
    IssueID                    int    `json:"id"`
	Creation                   string `json:"created_at"`
}

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

    collectPrometheusData(db)
    collectSeleniumData(db)
    collectOpenAIData(db)
    collectDockerData(db)
    collectMilvusData(db)
    collectGoData(db)


}

func collectPrometheusData (db *sql.DB) error {

    drop_table := `drop table if exists prometheus_issues`

	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}
    create_table := `CREATE TABLE IF NOT EXISTS "prometheus_issues" (
                        "title"  VARCHAR(255),
                        "id"        INT,
                        "creationDate"  VARCHAR(255),
                        PRIMARY KEY("id")
                    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.github.com/repos/prometheus/prometheus/issues")
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var prometheuslist GitHubIssueStruct
	json.Unmarshal(body, &prometheuslist)

    for i:= 0; i < len(prometheuslist); i++ {
        var issue_name = prometheuslist[i].Title
        var creationDate = prometheuslist[i].Creation
        var issue_id = prometheuslist[i].IssueID

        sql := `INSERT INTO prometheus_issues ("title", "id", "creationDate") values($1, $2, $3)`
        _, err = db.Exec(
			sql,
			issue_name,
            issue_id,
            creationDate)

		if err != nil {
			log.Fatalf("Error on prometheus insertion: %s", err.Error())
		}

    }
    return nil
}

func collectSeleniumData (db *sql.DB) error {

    drop_table := `drop table if exists selenium_issues`

	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}
    create_table := `CREATE TABLE IF NOT EXISTS "selenium_issues" (
                        "title"  VARCHAR(255),
                        "id"        INT,
                        "creationDate"  VARCHAR(255),
                        PRIMARY KEY("id")
                    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.github.com/repos/seleniumhq/selenium/issues")
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var seleniumlist GitHubIssueStruct
	json.Unmarshal(body, &seleniumlist)

    for i:= 0; i < len(seleniumlist); i++ {
        var issue_name = seleniumlist[i].Title
        var creationDate = seleniumlist[i].Creation
        var issue_id = seleniumlist[i].IssueID

        sql := `INSERT INTO selenium_issues ("title", "id", "creationDate") values($1, $2, $3)`
        _, err = db.Exec(
			sql,
			issue_name,
            issue_id,
            creationDate)

		if err != nil {
			log.Fatalf("Error on selenium insertion: %s", err.Error())
		}

    }
    return nil
}

func collectOpenAIData (db *sql.DB) error {

    drop_table := `drop table if exists openai_issues`

	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}
    create_table := `CREATE TABLE IF NOT EXISTS "openai_issues" (
                        "title"  VARCHAR(255),
                        "id"        INT,
                        "creationDate"  VARCHAR(255),
                        PRIMARY KEY("id")
                    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.github.com/repos/openai/openai-python/issues")
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var openailist GitHubIssueStruct
	json.Unmarshal(body, &openailist)

    for i:= 0; i < len(openailist); i++ {
        var issue_name = openailist[i].Title
        var creationDate = openailist[i].Creation
        var issue_id = openailist[i].IssueID

        sql := `INSERT INTO openai_issues ("title", "id", "creationDate") values($1, $2, $3)`
        _, err = db.Exec(
			sql,
			issue_name,
            issue_id,
            creationDate)

		if err != nil {
			log.Fatalf("Error on openai insertion: %s", err.Error())
		}

    }
    return nil
}

func collectDockerData (db *sql.DB) error {

    drop_table := `drop table if exists docker_issues`

	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}
    create_table := `CREATE TABLE IF NOT EXISTS "docker_issues" (
                        "title"  VARCHAR(255),
                        "id"        INT,
                        "creationDate"  VARCHAR(255),
                        PRIMARY KEY("id")
                    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.github.com/repos/docker/cli/issues")
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var dockerlist GitHubIssueStruct
	json.Unmarshal(body, &dockerlist)

    for i:= 0; i < len(dockerlist); i++ {
        var issue_name = dockerlist[i].Title
        var creationDate = dockerlist[i].Creation
        var issue_id = dockerlist[i].IssueID

        sql := `INSERT INTO docker_issues ("title", "id", "creationDate") values($1, $2, $3)`
        _, err = db.Exec(
			sql,
			issue_name,
            issue_id,
            creationDate)

		if err != nil {
			log.Fatalf("Error on docker insertion: %s", err.Error())
		}

    }
    return nil
}

func collectMilvusData (db *sql.DB) error {

    drop_table := `drop table if exists milvus_issues`

	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}
    create_table := `CREATE TABLE IF NOT EXISTS "milvus_issues" (
                        "title"  VARCHAR(255),
                        "id"        INT,
                        "creationDate"  VARCHAR(255),
                        PRIMARY KEY("id")
                    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.github.com/repos/milvus-io/milvus/issues")
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var milvuslist GitHubIssueStruct
	json.Unmarshal(body, &milvuslist)

    for i:= 0; i < len(milvuslist); i++ {
        var issue_name = milvuslist[i].Title
        var creationDate = milvuslist[i].Creation
        var issue_id = milvuslist[i].IssueID

        sql := `INSERT INTO milvus_issues ("title", "id", "creationDate") values($1, $2, $3)`
        _, err = db.Exec(
			sql,
			issue_name,
            issue_id,
            creationDate)

		if err != nil {
			log.Fatalf("Error on milvus insertion: %s", err.Error())
		}

    }
    return nil
}

func collectGoData (db *sql.DB) error {

    drop_table := `drop table if exists go_issues`

	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}
    create_table := `CREATE TABLE IF NOT EXISTS "go_issues" (
                        "title"  VARCHAR(255),
                        "id"        INT,
                        "creationDate"  VARCHAR(255),
                        PRIMARY KEY("id")
                    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.github.com/repos/golang/go/issues")
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var golist GitHubIssueStruct
	json.Unmarshal(body, &golist)

    for i:= 0; i < len(golist); i++ {
        var issue_name = golist[i].Title
        var creationDate = golist[i].Creation
        var issue_id = golist[i].IssueID

        sql := `INSERT INTO go_issues ("title", "id", "creationDate") values($1, $2, $3)`
        _, err = db.Exec(
			sql,
			issue_name,
            issue_id,
            creationDate)

		if err != nil {
			log.Fatalf("Error on go insertion: %s", err.Error())
		}

    }
    return nil
}

// func collectStackOverflowData (db *sql.DB) error{

// }