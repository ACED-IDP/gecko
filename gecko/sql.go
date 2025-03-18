package gecko

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type Document struct {
	ID      int             `db:"id"`
	Name    string          `db:"name"`
	Content json.RawMessage `db:"content"` // Store JSON as raw bytes
}

func configGET(db *sqlx.DB, name string) (any, error) {
	stmt := "SELECT name, content FROM documents WHERE name=$1"
	doc := &Document{}
	err := db.Get(doc, stmt, name)
	if err != nil {
		return nil, err
	}

	var content map[string]any
	err = json.Unmarshal(doc.Content, &content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func configPUT(db *sqlx.DB, name string, data map[string]any) error {
	stmt := `
                INSERT INTO documents (name, content)
                VALUES ($1, $2)
                ON CONFLICT (name)
                DO UPDATE SET content = $2;
        `
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
