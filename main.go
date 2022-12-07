package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/ravendb/ravendb-go-client"
	"github.com/valyala/fasthttp"
)

func main() {

	store := getDocumentStore("test")
	session, err := store.OpenSession("")
	if err != nil {
		log.Fatalf("store.OpenSession() failed with %s", err)
	}

	router := routing.New()

	router.Get("/", func(c *routing.Context) error {
		err = session.Store(&User{
			ID:        gonanoid.Must(),
			CreatedAt: time.Now(),
			Name:      "Tomas",
			LastName:  "Weigenast",
			Friends:   []string{"Jorge", "Diana", "Facundo"},
			Likes: map[string]bool{
				"cats": true,
				"dogs": false,
			},
		})

		if err != nil {
			c.WriteString(err.Error())
			return nil
		}

		fmt.Println("stored")

		err = session.SaveChanges()

		if err != nil {
			c.WriteString(err.Error())
		} else {
			fmt.Println("saved")
			c.WriteString("success")
		}

		return nil
	})

	router.Get("/read", func(c *routing.Context) error {
		var user *User
		err := session.Load(&user, "users/1-A")
		fmt.Println("user received")

		if err != nil {
			c.WriteString(err.Error())
		} else {
			j, _ := json.Marshal(user)
			c.WriteString(string(j))
		}

		return nil
	})

	fmt.Println("up and running")
	panic(fasthttp.ListenAndServe(":9000", router.HandleRequest))
}

func getDocumentStore(databaseName string) *ravendb.DocumentStore {
	serverNodes := []string{"http://localhost:8080"}
	store := ravendb.NewDocumentStore(serverNodes, databaseName)
	if err := store.Initialize(); err != nil {
		panic(err)
	}
	return store
}

type User struct {
	ID        string
	CreatedAt time.Time
	Name      string
	LastName  string
	Friends   []string
	Likes     map[string]bool
}
