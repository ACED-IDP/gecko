package chameleon

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Document struct {
	ID      int             `db:"id"`
	Name    string          `db:"name"`
	Content json.RawMessage `db:"content"` // Store JSON as raw bytes
}

func configGET(db *sqlx.DB, name string) (any, error) {
	stmt := "SELECT id, name, content FROM documents WHERE name=$1"
	doc := &Document{}
	err := db.Select(&doc, stmt, name)
	if err != nil {
		return nil, err
	}

	// optionally print the value for the user
	var content map[string]any
	err = json.Unmarshal(doc.Content, &content)
	if err != nil {
		return nil, err
	}
	fmt.Println("VALUE OF DOC: ", content)

	return content, nil
}

/*
init db
CREATE TABLE documents (

	id SERIAL PRIMARY KEY,
	name TEXT,
	content JSONB

);
*/

func configPUT(db *sqlx.DB, name string, data map[string]any) error {
	stmt := "INSERT INTO documents (name, content) VALUES ($1, $2)"
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = db.Exec(stmt, name, jsonData)
	if err != nil {
		return err
	}
	return nil
}
