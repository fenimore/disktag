package main

import (
	"database/sql"
	"os" // os?
	"time"

	"fmt"
	_ "github.com/bmizerany/pq" // HEROKU needs for parsing
)

const (
	DB_USER     = "disk"
	DB_PASSWORD = "disk"
	DB_NAME     = "disk"
	DB_SSL      = "disable" // "require" for HEROKU
)

// Ideas on list items?
const (
	Next = iota
	ToSend
	Waiting
	Confirm
	Done
)

// Database PostrgreSQL models

// like a 'Board'
// The Stages field will order the stages
type Document struct {
	Id      int       `json:"id"`
	Title   string    `json:"title"`
	Stages  []*Stage  `json:"stages"`
	Members []*Member `json:"members"`
	Cards   []*Card   `json:"cards"`
}

type Member struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Stage struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Cards       []*Card   `json:"cards"`
	Subscribers []*Member `json:"subscribers"`
}

type Card struct {
	Id          int           `json:"id"`
	Description string        `json:"info"`
	Due         time.Time     `json:"due_date"`
	Attachments []*Attachment `json:"attachments"`
	Labels      []*Label      `json:"labels"`
}

type Attachment struct {
	Attachment *os.File `json:"attachment"`
}

type Label struct {
	Label string `json:"label"`
}

// NOTE:  Ignore archived in db tables
// TODO: Add blob file
var CardSchema = `
CREATE TABLE IF NOT EXISTS cards(
    card_id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    due_date TIMESTAMP
);`

var StageSchema = `
CREATE TABLE IF NOT EXISTS stages(
    stage_id SERIAL PRIMARY KEY,
    title TEXT NOT NULL

);`

var MemberSchema = `
CREATE TABLE IF NOT EXISTS members(
    member_id SERIAL PRIMARY KEY,
    name VARCHAR(70)
);`

var StageCardsSchema = `
CREATE TABLE IF NOT EXISTS stage_cards(
    stage_id INT REFERENCES stages (stage_id),
    card_id INT REFERENCES cards (card_id),
    CONSTRAINT stage_card_key PRIMARY KEY (stage_id, card_id)
);`

// TODO: Document Schema and Relational tables

// NOTE: Members Subscribed to Stage
var SubscriptionSchema = `
CREATE TABLE IF NOT EXISTS subscriptions(
    member_id INT REFERENCES members (member_id),
    stage_id INT REFERENCES stages (stage_id),
    CONSTRAINT subs_key PRIMARY KEY (member_id, stage_id)
);`

// NOTE: What cards belong to whom
var MembershipSchema = `
CREATE TABLE IF NOT EXISTS membership(
    member_id INT REFERENCES members (member_id),
    card_id INT REFERENCES cards (card_id),
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
	_, err = db.Exec(StageSchema)
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
	_, err = db.Exec(StageCardsSchema)
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

// TODO: change to s
func InsertStage(db *sql.DB, s *Stage) (int, error) {
	var lastInsertId int
	err := db.QueryRow("INSERT INTO stages(title) VALUES($1) returning stage_id;",
		s.Title).Scan(&lastInsertId)
	if err != nil {
		return -1, err
	}
	s.Id = lastInsertId
	return lastInsertId, nil
}

// InsertCard returns an id of the inserted card.
func InsertCard(db *sql.DB, c *Card) (int, error) {
	var lastInsertId int
	err := db.QueryRow("INSERT INTO cards(description, due_date)"+
		" VALUES($1,$2, $3) returning card_id;",
		c.Description, c.Due).Scan(&lastInsertId)
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

// SelectStage returns a *Stage from db.
func SelectStage(db *sql.DB, id int) (*Stage, error) {
	s := new(Stage)
	stmt := "select * from stages where stage_id = $1"
	err := db.QueryRow(stmt, id).Scan(&s.Id, &s.Title)
	if err != nil {
		return s, err
	}
	// TODO: Find all cards that belong to said stage
	return s, nil
}

// SelectCard returns  a *Card from DB with id
func SelectCard(db *sql.DB, id int) (*Card, error) {
	c := new(Card)
	stmt := "select * from cards where card_id = $1"
	err := db.QueryRow(stmt, id).Scan(&c.Id, &c.Description, &c.Due)
	if err != nil {
		return c, err
	}
	// TODO: Get all members of a card
	return c, nil
}

// SelectStage returns a *Stage from db.
// TODO: select according to Documetn
func SelectAllStages(db *sql.DB, id int) ([]*Stage, error) {
	// NOTE: Id is document?
	stages := make([]*Stage, 0)
	stmt := "select * from stages"
	rows, err := db.Query(stmt)
	defer rows.Close()
	if err != nil {
		return stages, err
	}

	for rows.Next() {
		s := new(Stage)
		err = rows.Scan(&s.Id, &s.Title)
		if err != nil {
			return nil, err
		}

		stages = append(stages, s)
	}

	return stages, nil
}

// SelectCard returns  a *Card from DB with id
// TODO: Select according to Document
func SelectAllCards(db *sql.DB, id int) ([]*Card, error) {
	cards := make([]*Card, 0)
	stmt := "select * from cards"
	rows, err := db.Query(stmt)
	if err != nil {
		return cards, err
	}
	for rows.Next() {
		c := new(Card)
		err = rows.Scan(&c.Id, &c.Description, &c.Due)
		if err != nil {
			return nil, err
		}

		cards = append(cards, c)

	}
	// TODO: Get all members of a card// populate?
	return cards, nil
}

// SelectMember isn't used for now...
func SelectMember(db *sql.DB, id int) (*Member, error) {
	m := new(Member)
	stmt := "select * from members where member_id = $1"
	err := db.QueryRow(stmt, id).Scan(&m.Id, &m.Name)
	if err != nil {
		return m, err
	}

	return m, nil
}
