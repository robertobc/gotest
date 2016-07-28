package testlib

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3" //sqlite3 driver
	_ "github.com/robertobc/gotest/sub"
)

//Bacon contains baconipsum
type Bacon struct {
	Val []byte
}

//BaconMe gets some Bacon
func BaconMe() (*Bacon, error) {
	resp, err := http.Get("https://baconipsum.com/api/?type=meat-and-filler")
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	b := &Bacon{Val: body}

	return b, nil
}

//Pixel gets an image of size width*height
func Pixel(width, height int) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://lorempixel.com/%d/%d/", width, height))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println(err)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(body), nil
}

//Fib calculates the fibonacci sequence using channels and go routines
func Fib(num int) int {
	c := make(chan int)
	quit := make(chan int)

	answer := 1

	go func() {
		for i := 0; i < num; i++ {
			answer += <-c
		}
		quit <- 0
	}()

	fibonacci(c, quit)
	return answer
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			return
		}
	}
}

type userResult struct {
	Res []*user `json:"results"`
}

type user struct {
	Name userName `json:"name"`
}

type userName struct {
	Title string `json:"title"`
	First string `json:"first"`
	Last  string `json:"last"`
}

//UsersN generates N random users and returns the json
func UsersN(num int) (string, error) {
	users := make([]*user, 0, num)
	lock := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(num)

	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			u := getUser()
			if u == nil {
				return
			}

			lock.Lock()
			users = append(users, u)
			lock.Unlock()
		}()
	}

	wg.Wait()
	jsonUsers, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(jsonUsers), nil

}

func getUser() *user {
	resp, err := http.Get("https://randomuser.me/api")
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println(err)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	r := &userResult{}
	err = json.Unmarshal(body, r)
	if err != nil {
		log.Println(err)
		return nil
	}

	return r.Res[0]
}

//TryDB creates a database and table (foo) at the specified path
func TryDB(DBPath string) error {
	os.Remove(DBPath)

	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Println("error creating database", err)
		return err
	}
	defer db.Close()

	sqlStmt := `CREATE TABLE users (id INTEGER NOT null PRIMARY KEY, name TEXT);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	_, err = db.Exec("INSERT INTO users(id, name) VALUES (1, 'Peter Parker'), (2, 'J Jonah Jameson'), (3, 'Cletus Cassidy')")
	if err != nil {
		log.Println("error inserting", err)
		return err
	}

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Println("error selecting", err)
		return err
	}
	defer rows.Close()

	type user struct {
		id   int
		name string
	}

	users := []user{}

	for rows.Next() {
		u := user{}
		if err := rows.Scan(&u.id, &u.name); err != nil {
			log.Println("error", err)
			break
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		log.Println("error iterating through results", err)
		return err
	}

	log.Printf("QUERY RESULTS: %+v\n", users)

	return nil
}
