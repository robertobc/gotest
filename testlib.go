package testlib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

//Bacon contains baconipsum
type Bacon struct {
	val []byte
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

	b := &Bacon{val: body}

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

	answer := 0

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

type userName struct {
	Title string `json:"title"`
	First string `json:"first"`
	Last  string `json:"last"`
}
type user struct {
	Name userName `json:"name"`
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
			if u != nil {
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

	u := &user{}
	err = json.Unmarshal(body, u)
	if err != nil {
		log.Println(err)
		return nil
	}

	return u
}
