package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
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
	db, err := pgxpool.Connect(context.Background(), "postgres://localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var a Foo
	db.QueryRow(context.Background(), "SELECT 1").Scan(&a)
	fmt.Println(a) // {1}

	var b Foo
	db.QueryRow(context.Background(), "SELECT NULL").Scan(&b)
	fmt.Println(b) // {0}

	var c *Foo
	db.QueryRow(context.Background(), "SELECT 1").Scan(&c)
	fmt.Println(c) // &{0}

	var d *Foo
	db.QueryRow(context.Background(), "SELECT NULL").Scan(&d)
	fmt.Println(d) // <nil>

	var e FooBar
	db.QueryRow(context.Background(), "SELECT 2").Scan(&e.Id)
	x, _ := json.Marshal(e)
	fmt.Println(string(x)) // {"Id":{"Id":0}}

	f := &FooBar{}
	db.QueryRow(context.Background(), "SELECT 2").Scan(&f.Id)
	y, _ := json.Marshal(f)
	fmt.Println(string(y)) // {"Id":{"Id":0}}
}
