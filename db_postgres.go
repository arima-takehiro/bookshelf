package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

type Pdb struct{
	client *sql.DB
}

func NewPdb() *Pdb {
	db, err := sql.Open("postgres", fmt.Sprintf("user=root password=password host=%s port=5432 dbname=bookshelf sslmode=disable", os.Getenv("HOST")))
	if err != nil {
		log.Fatalf("postgres.NewClient: %v", err)
	}
	return &Pdb{db}
}

// get book data by specifying bookID
func (p *Pdb) GetBook(ctx context.Context, id string) (*Book, error) {
	ctx = context.TODO()
	cmd := `SELECT * FROM book WHERE id = $1`
	rows, err := p.client.Query(cmd, id)
	if err != nil{
		return nil, err
	}
	var b Book
	for rows.Next() {
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ImageURL, &b.PublishedDate, &b.Description); err != nil {
			return nil, err
		}
	}
	return &b, nil
}

func (p *Pdb) AddBook(ctx context.Context, b *Book) (id string, err error) {
	ctx = context.TODO()
	b.ID = createUID()
	cmd := "INSERT INTO book (id, title, author, published_date, image_url, description) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = p.client.Exec(cmd, b.ID, b.Title, b.Author, b.PublishedDate, b.ImageURL, b.Description)
	fmt.Println(b)
	if err != nil{
		return "", fmt.Errorf("PostgreSQL: Create: %v", err)
	}
	return b.ID, nil
}

func createUID()string{
	uid, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err.Error())
	}
	return uid.String()
}

func (p *Pdb) DeleteBook(ctx context.Context, id string) error {
	cmd := `DELETE FROM book WHERE id = $1`
	if _, err := p.client.Exec(cmd, id); err != nil{
		return fmt.Errorf("PostgreSQL: Delete: %v", err)
	}
	return nil
}

// UpdateBook updates the entry for a given book by specifying its ID
func (p *Pdb) UpdateBook(ctx context.Context, b *Book) error {
	cmd := `UPDATE book SET (id, title, author, publish_date, image_url, description) VALUES ($1, $2, $3, $4, $5, $6) WHERE id = ?`
	if _, err := p.client.Exec(cmd, b.ID, b.Title, b.Author, b.PublishedDate, b.ImageURL, b.Description, b.ID); err != nil {
		return fmt.Errorf("PostgreSQL: Set: %v", err)
	}
	return nil
}

// ListBooks return a list of books, ordered by title.
func(p *Pdb) ListBooks(ctx context.Context)([]*Book, error){
	ctx = context.TODO()
	cmd := `SELECT * FROM book ORDER BY title`
	rows, err := p.client.Query(cmd)
	if err != nil{
		return nil, err
	}

	var books []*Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ImageURL, &b.PublishedDate, &b.Description); err != nil {
			return nil, err
		}
		books = append(books, &b)
		log.Printf("Book %q ID: %q", b.Title, b.ID)
	}
	return books, nil
}
