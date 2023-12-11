package main

import (
	"context"
	"darren/moto"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

var ctx = context.Background()
var rdb *redis.Client

func main() {
	c := cron.New()
	c.AddFunc("@every 3m", func() {
		task()
	})
	c.Start()

	connectDB()

	http.HandleFunc("/ping", pong)

	fmt.Println(fmt.Sprintf("🏄‍♂️ : %d", 9390))
	http.ListenAndServe(":9390", nil)
}

func pong(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		word := req.PostFormValue("word")
		fmt.Fprintf(res, "pong=%s\n", word)
	default:
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func task() {
	pings := readPings()
	for _, link := range pings {
		err := ping(link)
		if err != nil {
			// fmt.Println("not good: " + link)
			moto.Send("hang", "lynk server down: "+link, "lynk server down: "+link, []string{}, "", "", "")
			continue
		}
		// fmt.Println("good: " + link)
	}
}

func ping(link string) error {
	word := strconv.Itoa(rand.Int())
	data := url.Values{
		"word": []string{word},
	}
	res, err := http.PostForm(link, data)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	body := string(bodyBytes)
	body = strings.ReplaceAll(body, "\n", "")
	if body == "pong="+word {
		return nil
	}
	return fmt.Errorf("lynk server down: %s", link)
}

func readPings() (pings []string) {
	iter := rdb.SScan(ctx, "pings", 0, "", 0).Iterator()
	for iter.Next(ctx) {
		pings = append(pings, iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	return pings
}

func connectDB() *redis.Client {
	if rdb != nil {
		return rdb
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil
	}
	fmt.Println("💾 rdb connected")
	return rdb
}
