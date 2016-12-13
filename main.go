// trello clone

package main

import (
	"fmt"
	"os"
	"time"
)

type Member struct {
	Id   int
	Name string
}

type List struct {
	Title       string
	Cards       []*Card
	Subscribers []*Member

	// Unused
	Archived bool
}

type Card struct {
	Info    string
	File    *os.File
	Members []*Member
	Due     time.Time
	// if unused
	Archived bool
}

// NOTE: Fuck methods?
/* List Methods, Not currently Methods? */
func AddSubscriber(l *List, m *Member) {
	l.Subscribers = append(l.Subscribers, m)
}

/* Card Methods */
func AddMember(c *Card, m *Member) {
	c.Members = append(c.Members, m)
}

func main() {
	fmt.Println("Hello trello")
}
