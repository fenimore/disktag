package main

import (
	"database/sql"
	"os" // os?
	"time"

	"fmt"
	_ "github.com/bmizerany/pq" // HEROKU needs for parsing
)

const (
	DB_USER     = ""
	DB_PASSWORD = ""
	DB_NAME     = ""
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
	// if unused
	Archived bool
}

// NOTE:  Ignore archived in db tables
var CardSchema = `
CREATE TABLE IF NOT EXISTS cards(
    card_id SERIAL PRIMARY KEY,
    info TEXT NOT NULL,
    file BLOB,
    due_date TIMESTAMP,
    list_id INT REFERENCES lists (list_id) ON UPDATE CASCADE
);`

var ListSchema = `
CREATE TABLE IF NOT EXISTS lists(
    list_id SERIAL PRIMARY KEY,
    title TEXT NOT NULL

);`

var MemberSchema = `
CREATE TABLE IF NOT EXISTS member(
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
    card_id INT REFERENCES lists (card_id) ON UPDATE CASCADE,
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
	_, err = db.Exec(CardSchema)
	if err != nil {
		return err
	}
	_, err = db.Exec(ListSchema)
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
