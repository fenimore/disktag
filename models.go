package main

import (
	"database/sql"
	"os" // os?
	"time"

	"fmt"
	_ "github.com/bmizerany/pq" // HEROKU needs for parsing
)

const (
	DB_USER     = "trello"
	DB_PASSWORD = "trello"
	DB_NAME     = "trello"
	DB_SSL      = "disable" // "require" for HEROKU
)

// Database PostrgreSQL models

// like a 'Board'
type Document struct {
	Id      int
	Lists   []*List
	Members []*Member
	Cards   []*Card
}

type Member struct {
	Id   int
	Name string
}

type List struct {
	Id          int
	Title       string
	Cards       []*Card
	Subscribers []*Member

	// Unused
	Archived bool
}

type Card struct {
	Id      int
	Info    string
	File    *os.File
	Members []*Member
	Due     time.Time
	Stage   *List
	// if unused
	Archived bool
}

// NOTE:  Ignore archived in db tables
// TODO: Add blob file
var CardSchema = `
CREATE TABLE IF NOT EXISTS cards(
    card_id SERIAL PRIMARY KEY,
    info TEXT NOT NULL,
    due_date TIMESTAMP,
    list_id INT REFERENCES lists (list_id) ON UPDATE CASCADE
);`

var ListSchema = `
CREATE TABLE IF NOT EXISTS lists(
    list_id SERIAL PRIMARY KEY,
    title TEXT NOT NULL

);`

var MemberSchema = `
CREATE TABLE IF NOT EXISTS members(
    member_id SERIAL PRIMARY KEY,
    name VARCHAR(70)
);`

var SubscriptionSchema = `
CREATE TABLE IF NOT EXISTS subscriptions(
    member_id INT REFERENCES members (member_id) ON UPDATE CASCADE,
    list_id INT REFERENCES lists (list_id) ON UPDATE CASCADE,
    CONSTRAINT subs_key PRIMARY KEY (member_id, list_id)
);`

// NOTE: What cards belong to whom
var MembershipSchema = `
CREATE TABLE IF NOT EXISTS membership(
    member_id INT REFERENCES members (member_id) ON UPDATE CASCADE,
    card_id INT REFERENCES cards (card_id) ON UPDATE CASCADE,
    CONSTRAINT membership_key PRIMARY KEY (member_id, card_id)
);`

/* Functions */

func InitializeDB() (*sql.DB, error) {
	connection := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		DB_USER, DB_PASSWORD, DB_NAME, DB_SSL)

	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	var err error

	_, err = db.Exec(ListSchema)
	if err != nil {
		return err
	}

	_, err = db.Exec(CardSchema)
	if err != nil {
		return err
	}

	_, err = db.Exec(MemberSchema)
	if err != nil {
		return err
	}
	_, err = db.Exec(SubscriptionSchema)
	if err != nil {
		return err
	}
	_, err = db.Exec(MembershipSchema)
	if err != nil {
		return err
	}
	return nil
}

func InsertList(db *sql.DB, l *List) (int, error) {
	var lastInsertId int
	err := db.QueryRow("INSERT INTO lists(title) VALUES($1) returning list_id;",
		l.Title).Scan(&lastInsertId)
	if err != nil {
		return -1, err
	}
	l.Id = lastInsertId
	return lastInsertId, nil
}

// InsertCard returns an id of the inserted card.
func InsertCard(db *sql.DB, c *Card) (int, error) {
	var lastInsertId int
	err := db.QueryRow("INSERT INTO cards(info, due_date, list_id)"+
		" VALUES($1,$2, $3) returning card_id;",
		c.Info, c.Due, c.Stage).Scan(&lastInsertId)
	if err != nil {
		return -1, err
	}
	c.Id = lastInsertId
	return lastInsertId, nil
}

func InsertMember(db *sql.DB, m *Member) (int, error) {
	var lastInsertId int
	err := db.QueryRow("INSERT INTO members(name) VALUES($1) returning member_id;",
		m.Name).Scan(&lastInsertId)
	if err != nil {
		return -1, err
	}
	m.Id = lastInsertId
	return lastInsertId, nil
}