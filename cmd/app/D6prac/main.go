// эксперименты с встроенными структурами
package main

import (
	"fmt"
	"time"
)

type Entity struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (en *Entity) Init(id int) {
	en.ID = id
	en.CreatedAt = time.Now()
	en.UpdatedAt = time.Now()
}

func (en *Entity) Update() {
	//logic ...
	en.UpdatedAt = time.Now()
}

func (en *Entity) InitUpdate(id int) {
	en.Init(id)
	//logic ...
	en.UpdatedAt = time.Now()
}

type Person struct {
	Name string
	*Entity
}

type PersonAge struct {
	Age int
	Person
}

type A struct {
	Name string
}
type B struct {
	Name string
}

type C struct {
	A
	B
}

func (en *PersonAge) Init(id int) {
	en.ID = id
	en.CreatedAt = time.Now()
	en.UpdatedAt = time.Now()
	en.Age = 18
	en.Name = "Bob"
	fmt.Println("WOW")
}

func main() {
	//var pers PersonAge
	pers := PersonAge{Person: Person{Entity: &Entity{}}}
	pers.ID = 1
	fmt.Printf("%+v", pers)
}
