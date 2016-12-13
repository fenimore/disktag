// trello clone

package main

import (
	"database/sql"
	"fmt"
)

var (
	db *sql.DB
)

func main() {
	fmt.Println("Hello Backend")
	conn, err := InitializeDB()
	defer conn.Close()
	db = conn
	if err != nil {
		fmt.Println("Error DB", err)
	}

	err = CreateTables(conn)
	if err != nil {
		fmt.Println("Creation error", err)
	}

	Serve() // TODO: pass in connection and port flag
}
