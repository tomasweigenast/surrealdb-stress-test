package main

import (
	"encoding/json"
	"fmt"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/surrealdb/surrealdb.go"
	"github.com/valyala/fasthttp"
)

func main() {
	db, err := surrealdb.New("ws://0.0.0.0:8000/rpc")
	if err != nil {
		panic(err)
	}

	db.Signin(map[string]any{
		"user": "root",
		"pass": "root",
	})

	db.Use("test", "test")

	router := routing.New()

	router.Get("/", func(c *routing.Context) error {
		fmt.Println("received")
		_, err := db.Create("users:tomas", map[string]interface{}{
			"name":     "Tomas",
			"lastName": "Weigenast",
		})

		if err != nil {
			c.WriteString(err.Error())
		} else {
			c.WriteString("success")
		}

		return nil
	})

	router.Get("/read", func(c *routing.Context) error {
		fmt.Println("request received")
		res, err := db.Select("users:tomas")
		if err != nil {
			c.WriteString(err.Error())
		} else {
			j, _ := json.Marshal(res)
			c.WriteString(string(j))
		}

		return nil
	})

	panic(fasthttp.ListenAndServe(":8080", router.HandleRequest))
}
