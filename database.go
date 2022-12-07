package main

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/genjidb/genji"
	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/types"
)

type Db struct {
	db *genji.DB
}

type User struct {
	ID        string    `genji:"id"`
	CreatedAt time.Time `genji:"createdAt"`
	Name      string    `genji:"name"`
	LastName  string    `genji:"lastName"`
	Friends   []string  `genji:"friends"`
}

func NewDatabase() *Db {
	db, err := genji.Open("local.db")
	if err != nil {
		panic(err)
	}

	err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id              TEXT    PRIMARY KEY,
		createdAt       TEXT    NOT NULL,
		name            TEXT    NOT NULL,
		lastName 		TEXT 	NOT NULL,
		friends         ARRAY 	NOT NULL
	)`)

	if err != nil {
		panic(err)
	}

	return &Db{
		db: db,
	}
}

func (d *Db) InsertMany() {
	i := 0
	tx, _ := d.db.Begin(true)
	for i < 1_000_000 {
		u := User{
			ID:        fmt.Sprint(time.Now().UnixNano()),
			Name:      gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			CreatedAt: gofakeit.Date(),
			Friends: []string{
				gofakeit.Name(), gofakeit.Name(), gofakeit.Name(),
			},
		}
		err := tx.Exec(`INSERT INTO users VALUES ?`, &u)
		if err != nil {
			panic(err)
		}

		fmt.Println("inserted: ", i)

		i++
	}
	err := tx.Commit()
	if err != nil {
		panic(err)
	}
}

func (d *Db) InsertUser() {
	u := User{
		ID:        fmt.Sprint(time.Now().UnixNano()),
		CreatedAt: time.Now(),
		Name:      "Tomas",
		LastName:  "Wegenast",
		Friends:   []string{"Alejandro", "Magno"},
	}

	err := d.db.Exec(`INSERT INTO users VALUES ?`, &u)
	if err != nil {
		panic(err)
	}
}

func (d *Db) ReadUser(name string) []User {
	var users []User
	res, err := d.db.Query("SELECT * FROM users WHERE name LIKE ? LIMIT 100", name)
	if err != nil {
		panic(err)
	}

	defer res.Close()
	err = res.Iterate(func(d types.Document) error {
		var user User
		err := document.StructScan(d, &user)
		if err != nil {
			return err
		}

		users = append(users, user)
		return nil
	})

	if err != nil {
		panic(err)
	}

	return users
}
