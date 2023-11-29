package main

import (
    "fmt"
    "log"
    "database/sql"
    "encoding/json"
    "os"
    "net/http"
    "time"
    _ "github.com/lib/pq"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
    "io/ioutil"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    
)

type GitHubIssueStruct []struct {
	Title                      string `json:"title"`
    IssueID                    int    `json:"id"`
	Creation                   string `json:"created_at"`
}

type StackOverflowPosts struct{
    Items [] struct{
        QuestionID int `json:"question_id"`
        Creation int `json:"creation_date"`
        Title string `json:"title"`
    } `json:"items"`
}

var apiCallsMade = promauto.NewCounter(
    prometheus.CounterOpts{
        Name: "api_calls_made",
        Help: "Number of API calls made by application",
    },
)

var dataCollectionCounter = promauto.NewCounter(
    prometheus.CounterOpts{
        Name: "data_collected",
        Help: "Amount of data collected per second (per json entry)",
    },
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

    // for {
    //     //GitHub Data
    //     collectPrometheusData(db)
    //     collectSeleniumData(db)
    //     collectOpenAIData(db)
    //     collectDockerData(db)
    //     collectMilvusData(db)
    //     collectGoData(db)

    //     //StackOverflowData
    //     collectPrometheusStack(db)
    //     collectSeleniumStack(db)
    //     collectOpenaiPosts(db)
    //     collectDockerPosts(db)
    //     collectMilvusPosts(db)
    //     collectGoPosts(db)
    //     time.Sleep(24 * time.Hour)
    // }
    startTime := time.Now()

    collectPrometheusData(db)
    collectSeleniumData(db)
    collectOpenAIData(db)
    collectDockerData(db)
    collectMilvusData(db)
    collectGoData(db)
    
    //StackOverflowData
    collectPrometheusStack(db)
    collectSeleniumStack(db)
    collectOpenaiPosts(db)
    collectDockerPosts(db)
    collectMilvusPosts(db)
    collectGoPosts(db)

    endTime := time.Now()

    executionDuration := endTime.Sub(startTime)

    executionTime := float64(executionDuration.Seconds())

    // fmt.Println("API Calls Per Second: %f", apiCallsPerSecond)
    // fmt.Println("Data Collected Per Second: %f", dataCollectedPerSecond)

    fmt.Println("Execution Time: ", executionTime)

    http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(":2112", nil)
    

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
                        "creation_date"  VARCHAR(255),
                        PRIMARY KEY("id")
                    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.github.com/repos/prometheus/prometheus/issues")
    apiCallsMade.Inc()
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var prometheuslist GitHubIssueStruct
	json.Unmarshal(body, &prometheuslist)

    for i:= 0; i < len(prometheuslist); i++ {
        dataCollectionCounter.Inc()
        var issue_name = prometheuslist[i].Title
        var creationDate = prometheuslist[i].Creation
        var issue_id = prometheuslist[i].IssueID

        sql := `INSERT INTO prometheus_issues ("title", "id", "creation_date") values($1, $2, $3)`
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
                        "creation_date"  VARCHAR(255),
                        PRIMARY KEY("id")
                    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.github.com/repos/seleniumhq/selenium/issues")
    apiCallsMade.Inc()
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var seleniumlist GitHubIssueStruct
	json.Unmarshal(body, &seleniumlist)

    for i:= 0; i < len(seleniumlist); i++ {
        dataCollectionCounter.Inc()
        var issue_name = seleniumlist[i].Title
        var creationDate = seleniumlist[i].Creation
        var issue_id = seleniumlist[i].IssueID

        sql := `INSERT INTO selenium_issues ("title", "id", "creation_date") values($1, $2, $3)`
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
                        "creation_date"  VARCHAR(255),
                        PRIMARY KEY("id")
                    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.github.com/repos/openai/openai-python/issues")
    apiCallsMade.Inc()
    if err != nil {
		panic(err)
	}
    
    body, _ := ioutil.ReadAll(req.Body)
	var openailist GitHubIssueStruct
	json.Unmarshal(body, &openailist)

    for i:= 0; i < len(openailist); i++ {
        dataCollectionCounter.Inc()
        var issue_name = openailist[i].Title
        var creationDate = openailist[i].Creation
        var issue_id = openailist[i].IssueID

        sql := `INSERT INTO openai_issues ("title", "id", "creation_date") values($1, $2, $3)`
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
                        "creation_date"  VARCHAR(255),
                        PRIMARY KEY("id")
                    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.github.com/repos/docker/cli/issues")
    apiCallsMade.Inc()
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var dockerlist GitHubIssueStruct
	json.Unmarshal(body, &dockerlist)

    for i:= 0; i < len(dockerlist); i++ {
        dataCollectionCounter.Inc()
        var issue_name = dockerlist[i].Title
        var creationDate = dockerlist[i].Creation
        var issue_id = dockerlist[i].IssueID

        sql := `INSERT INTO docker_issues ("title", "id", "creation_date") values($1, $2, $3)`
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
                        "creation_date"  VARCHAR(255),
                        PRIMARY KEY("id")
                    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.github.com/repos/milvus-io/milvus/issues")
    apiCallsMade.Inc()
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var milvuslist GitHubIssueStruct
	json.Unmarshal(body, &milvuslist)

    for i:= 0; i < len(milvuslist); i++ {
        dataCollectionCounter.Inc()
        var issue_name = milvuslist[i].Title
        var creationDate = milvuslist[i].Creation
        var issue_id = milvuslist[i].IssueID

        sql := `INSERT INTO milvus_issues ("title", "id", "creation_date") values($1, $2, $3)`
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
                        "creation_date"  VARCHAR(255),
                        PRIMARY KEY("id")
                    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.github.com/repos/golang/go/issues")
    apiCallsMade.Inc()
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var golist GitHubIssueStruct
	json.Unmarshal(body, &golist)

    for i:= 0; i < len(golist); i++ {
        dataCollectionCounter.Inc()
        var issue_name = golist[i].Title
        var creationDate = golist[i].Creation
        var issue_id = golist[i].IssueID

        sql := `INSERT INTO go_issues ("title", "id", "creation_date") values($1, $2, $3)`
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

func collectPrometheusStack (db *sql.DB) error{
    drop_table := `drop table if exists prometheus_posts`

    _, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}
    create_table := `CREATE TABLE IF NOT EXISTS "prometheus_posts" (
        "title"  VARCHAR(255),
        "id"        INT,
        "creation_date"  INT,
        PRIMARY KEY("id")
    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.stackexchange.com/2.3/search?order=desc&sort=activity&intitle=prometheus&site=stackoverflow&key=c5hL2NXUJ*1TIoPb27Qudg((")
    apiCallsMade.Inc()
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var promStackList StackOverflowPosts
	json.Unmarshal(body, &promStackList)

    for i := 0; i < len(promStackList.Items); i++{
        dataCollectionCounter.Inc()
        var question_name = promStackList.Items[i].Title
        var creation_date = promStackList.Items[i].Creation
        var question_id = promStackList.Items[i].QuestionID

        sql := `INSERT INTO prometheus_posts ("title", "id", "creation_date") values($1, $2, $3)`
        _, err = db.Exec(
			sql,
			question_name,
            question_id,
            creation_date)

		if err != nil {
			log.Fatalf("Error on prometheus stack insertion: %s", err.Error())
		}

    }

    return nil
}

func collectSeleniumStack (db *sql.DB) error{
    drop_table := `drop table if exists selenium_posts`

    _, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}
    create_table := `CREATE TABLE IF NOT EXISTS "selenium_posts" (
        "title"  VARCHAR(255),
        "id"        INT,
        "creation_date"  INT,
        PRIMARY KEY("id")
    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.stackexchange.com/2.3/search?order=desc&sort=activity&tagged=selenium&site=stackoverflow&key=c5hL2NXUJ*1TIoPb27Qudg((")
    apiCallsMade.Inc()
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var seleniumStackList StackOverflowPosts
	json.Unmarshal(body, &seleniumStackList)

    for i := 0; i < len(seleniumStackList.Items); i++{
        dataCollectionCounter.Inc()
        var question_name = seleniumStackList.Items[i].Title
        var creation_date = seleniumStackList.Items[i].Creation
        var question_id = seleniumStackList.Items[i].QuestionID

        sql := `INSERT INTO selenium_posts ("title", "id", "creation_date") values($1, $2, $3)`
        _, err = db.Exec(
			sql,
			question_name,
            question_id,
            creation_date)

		if err != nil {
			log.Fatalf("Error on selenium stack insertion: %s", err.Error())
		}

    }

    return nil
}

func collectOpenaiPosts (db *sql.DB) error{
    drop_table := `drop table if exists openai_posts`

    _, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}
    create_table := `CREATE TABLE IF NOT EXISTS "openai_posts" (
        "title"  VARCHAR(255),
        "id"        INT,
        "creation_date"  INT,
        PRIMARY KEY("id")
    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.stackexchange.com/2.3/search?order=desc&sort=activity&tagged=openai-api&site=stackoverflow&key=c5hL2NXUJ*1TIoPb27Qudg((")
    apiCallsMade.Inc()
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var openStackList StackOverflowPosts
	json.Unmarshal(body, &openStackList)

    for i := 0; i < len(openStackList.Items); i++{
        dataCollectionCounter.Inc()
        var question_name = openStackList.Items[i].Title
        var creation_date = openStackList.Items[i].Creation
        var question_id = openStackList.Items[i].QuestionID

        sql := `INSERT INTO openai_posts ("title", "id", "creation_date") values($1, $2, $3)`
        _, err = db.Exec(
			sql,
			question_name,
            question_id,
            creation_date)

		if err != nil {
			log.Fatalf("Error on openai stack insertion: %s", err.Error())
		}

    }

    return nil
}

func collectDockerPosts (db *sql.DB) error{
    drop_table := `drop table if exists docker_posts`

    _, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}
    create_table := `CREATE TABLE IF NOT EXISTS "docker_posts" (
        "title"  VARCHAR(255),
        "id"        INT,
        "creation_date"  INT,
        PRIMARY KEY("id")
    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.stackexchange.com/2.3/search?order=desc&sort=activity&tagged=docker&site=stackoverflow&key=c5hL2NXUJ*1TIoPb27Qudg((")
    apiCallsMade.Inc()
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var dockerStackList StackOverflowPosts
	json.Unmarshal(body, &dockerStackList)

    for i := 0; i < len(dockerStackList.Items); i++{
        dataCollectionCounter.Inc()
        var question_name = dockerStackList.Items[i].Title
        var creation_date = dockerStackList.Items[i].Creation
        var question_id = dockerStackList.Items[i].QuestionID

        sql := `INSERT INTO docker_posts ("title", "id", "creation_date") values($1, $2, $3)`
        _, err = db.Exec(
			sql,
			question_name,
            question_id,
            creation_date)

		if err != nil {
			log.Fatalf("Error on docker stack insertion: %s", err.Error())
		}

    }

    return nil
}

func collectMilvusPosts (db *sql.DB) error{
    drop_table := `drop table if exists milvus_posts`

    _, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}
    create_table := `CREATE TABLE IF NOT EXISTS "milvus_posts" (
        "title"  VARCHAR(255),
        "id"        INT,
        "creation_date"  INT,
        PRIMARY KEY("id")
    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.stackexchange.com/2.3/search?order=desc&sort=activity&tagged=milvus&site=stackoverflow&key=c5hL2NXUJ*1TIoPb27Qudg((")
    apiCallsMade.Inc()
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var milvusStackList StackOverflowPosts
	json.Unmarshal(body, &milvusStackList)

    for i := 0; i < len(milvusStackList.Items); i++{
        dataCollectionCounter.Inc()
        var question_name = milvusStackList.Items[i].Title
        var creation_date = milvusStackList.Items[i].Creation
        var question_id = milvusStackList.Items[i].QuestionID

        sql := `INSERT INTO milvus_posts ("title", "id", "creation_date") values($1, $2, $3)`
        _, err = db.Exec(
			sql,
			question_name,
            question_id,
            creation_date)

		if err != nil {
			log.Fatalf("Error on milvus stack insertion: %s", err.Error())
		}

    }

    return nil
}

func collectGoPosts (db *sql.DB) error{
    drop_table := `drop table if exists go_posts`

    _, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}
    create_table := `CREATE TABLE IF NOT EXISTS "go_posts" (
        "title"  VARCHAR(255),
        "id"        INT,
        "creation_date"  INT,
        PRIMARY KEY("id")
    );`

    _, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

    req, err := http.Get("https://api.stackexchange.com/2.3/search?order=desc&sort=activity&intitle=golang&site=stackoverflow&key=c5hL2NXUJ*1TIoPb27Qudg((")
    apiCallsMade.Inc()
    if err != nil {
		panic(err)
	}

    body, _ := ioutil.ReadAll(req.Body)
	var goStackList StackOverflowPosts
	json.Unmarshal(body, &goStackList)

    for i := 0; i < len(goStackList.Items); i++{
        dataCollectionCounter.Inc()
        var question_name = goStackList.Items[i].Title
        var creation_date = goStackList.Items[i].Creation
        var question_id = goStackList.Items[i].QuestionID

        sql := `INSERT INTO go_posts ("title", "id", "creation_date") values($1, $2, $3)`
        _, err = db.Exec(
			sql,
			question_name,
            question_id,
            creation_date)

		if err != nil {
			log.Fatalf("Error on go stack insertion: %s", err.Error())
		}

    }

    return nil
}