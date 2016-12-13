// trello clone

package main

import (
	"fmt"
)

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
	db, err := InitializeDB()
	defer db.Close()
	if err != nil {
		fmt.Println("Error DB", err)
	}

	err = CreateTables(db)
	if err != nil {
		fmt.Println("Creation error", err)
	}

	Serve() // TODO: pass in connection and port flag
}
