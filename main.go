package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

func main() {
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r := gin.Default()

	r.GET("/getCount", func(c *gin.Context) {
		i := 0
		// if key "count" not exists set it value to 1
		// exists then increase by 1 and store in redis
		val, err := rdb.Get(ctx, "count").Result()
		if err == redis.Nil {
			err := rdb.Set(ctx, "count", 1, 0).Err()
			if err != nil {
				panic(err)
			}
			i = 1
		} else if err != nil {
			panic(err)
		} else {
			i, _ = strconv.Atoi(val)
			i++
			err := rdb.Set(ctx, "count", i, 0).Err()
			if err != nil {
				panic(err)
			}
		}

		c.String(http.StatusOK, strconv.Itoa(i))
	})

	r.GET("/go", func(c *gin.Context) {
		intChan := make(chan int, 1)
		for k := 0; k < 20; k++ {

			res, err := http.Get("https://localhost/getCount")
			if err != nil {
				panic(err)
			}

			body, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			i, _ := strconv.Atoi(string(body))
			intChan <- i
			fmt.Println(<-intChan)
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
