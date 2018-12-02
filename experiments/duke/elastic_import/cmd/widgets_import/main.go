package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "vivo_data"
	password = "experiment4"
	dbname   = "vivo_data"
)

type Resource struct {
	Uri   string         `db:"uri"`
	Type  string         `db:"type"`
	Hash  string         `db:"hash"`
	Data  types.JSONText `db:"data"`
	DataB types.JSONText `db:"data_b"`
}

var client *http.Client

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 50
)

var psqlInfo string

func init() {
	client = createHTTPClient()
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}
	return client
}

type Position struct {
	Uri        string `json:"uri"`
	Label      string `json:"label"`
	Attributes struct {
		OrganizationLabel string `json:"organizationLabel"`
	} `json:"attributes"`
}

type WidgetsPerson struct {
	Uri        string `json:"uri"`
	Attributes struct {
		FirstName  string  `json:"firstName"`
		LastName   string  `json:"lastName"`
		MiddleName *string `json:"middleName"`
	} `json:"attributes"`
	Positions []Position `json:"positions"`
}

func widgetsParse(duid string) WidgetsPerson {
	url := "https://scholars.duke.edu/widgets/api/v0.9/people/complete/all.json?uri=" + duid
	//fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		// NOTE: returning a 'blank' person
		fmt.Println("widgets", err)
		return WidgetsPerson{}
	}

	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		// NOTE: returning 'blank' person
		fmt.Println("widgets", err)
		return WidgetsPerson{}
	}

	defer res.Body.Close()

	var person WidgetsPerson
	json.Unmarshal([]byte(body), &person)
	return person
}

type ResourcePerson struct {
	Uri        string
	FirstName  string
	LastName   string
	MiddleName *string
}

//https://stackoverflow.com/questions/2377881/how-to-get-a-md5-hash-from-a-string-in-golang
//https://stackoverflow.com/questions/2377881/how-to-get-a-md5-hash-from-a-string-in-golang
func makeHash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

var db *sqlx.DB

func saveAsResource(person WidgetsPerson) {
	fmt.Printf("saving %v\n", person.Uri)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// NOTE: if person.Uri is null - should probably exit
	obj := ResourcePerson{person.Uri,
		person.Attributes.FirstName,
		person.Attributes.LastName,
		person.Attributes.MiddleName}

	str, err := json.Marshal(obj)
	if err != nil {
		log.Fatalln(err)
	}
	hash := makeHash(string(str))

	found := Resource{}
	res := &Resource{person.Uri, "Person", hash, str, str}
	fmt.Println(res)
	
	findSql := `SELECT uri, type, hash, data, data_b  FROM resources 
	  WHERE (uri = $1 AND type = $2)`

	err = db.Get(&found, findSql, person.Uri, "Person")

	tx := db.MustBegin()
	if err != nil {
		// NOTE: assuming the error means it doesn't exist
		log.Println("GET:%v", err)
		// must be an add?
		sql := `INSERT INTO resources (uri, type, hash, data, data_b) 
	      VALUES (:uri, :type, :hash, :data, :data_b)`
		_, err := tx.NamedExec(sql, res)
		if err != nil {
			log.Fatalln("INSERT:%v", err)
		}
	} else {
		fmt.Println("found!!" + found.Uri)
		sql := `UPDATE resources 
	    set uri = :uri, 
		type = :type, 
		hash = :hash, 
		data = :data, 
		data_b = :data_b
		WHERE uri = :uri and type = :type`
		_, err := tx.NamedExec(sql, res)

		if err != nil {
			log.Fatalln("UPDATE:%v", err)
		}
	}
	tx.Commit()
}

func processDuids(cin <-chan string) <-chan WidgetsPerson {
	out := make(chan WidgetsPerson)
	defer wg.Done()
	go func() {
		for line := range cin {
			out <- widgetsParse(line)
		}
		close(out)
	}()
	return out
}

func persistWidgets(cin <-chan WidgetsPerson) {
	go func() {
		for person := range cin {
			//fmt.Printf("saving %v\n", person.Uri)
			saveAsResource(person)
			// TODO: need to save positions
			//positions := person.Positions
			//for _, pos := range positions {
			//	fmt.Println(pos.Label)
			//}
		}
		// 'sink' so need to close waitgroup
		wg.Done()
	}()
}

func produceDuids(filename string) <-chan string {
	c := make(chan string)
	defer wg.Done()

	go func() {
		file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println("could not open file ", filename)
			close(c)
			return
		}
		defer file.Close()

		sc := bufio.NewScanner(file)
		for sc.Scan() {
			c <- sc.Text()
		}
		if err := sc.Err(); err != nil {
			close(c)
			return
		}
		// close
		close(c)
	}()
	return c
}

var wg sync.WaitGroup

func main() {
	start := time.Now()

	if len(os.Args) == 1 {
		fmt.Println("need filename arg")
		os.Exit(1)
	}

	filename := os.Args[1]

	wg.Add(3)
	duids := produceDuids(filename)
	widgets := processDuids(duids)
	persistWidgets(widgets)

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("%s\n", elapsed)
}
