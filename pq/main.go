package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
)

type (
	Foo struct {
		Id int
	}

	FooBar struct {
		Id *Foo
	}
)

func (f *Foo) Scan(value interface{}) error {
	if i, ok := value.(int64); ok {
		f.Id = int(i)
	}
	return nil
}

func main() {
	db, err := sql.Open("postgres", "postgres://localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var a Foo
	db.QueryRow("SELECT 1").Scan(&a)
	fmt.Println(a) // {1}

	var b Foo
	db.QueryRow("SELECT NULL").Scan(&b)
	fmt.Println(b) // {0}

	var c *Foo
	db.QueryRow("SELECT 1").Scan(&c)
	fmt.Println(c) // &{1}

	var d *Foo
	db.QueryRow("SELECT NULL").Scan(&d)
	fmt.Println(d) // <nil>

	var e FooBar
	db.QueryRow("SELECT 2").Scan(&e.Id)
	x, _ := json.Marshal(e)
	fmt.Println(string(x)) // {"Id":{"Id":2}}

	f := &FooBar{}
	db.QueryRow("SELECT 2").Scan(&f.Id)
	y, _ := json.Marshal(f)
	fmt.Println(string(y)) // {"Id":{"Id":2}}
}
