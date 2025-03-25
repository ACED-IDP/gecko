package gecko

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/ACED-IDP/gecko/gecko/config"
	"github.com/jmoiron/sqlx"
)

type Document struct {
	ID      int             `db:"id"`
	Name    string          `db:"name"`
	Content json.RawMessage `db:"content"` // Store JSON as raw bytes
}

func configGET(db *sqlx.DB, name string) (map[string]any, error) {
	stmt := "SELECT name, content FROM documents WHERE name=$1"
	doc := &Document{}
	err := db.Get(doc, stmt, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var content []config.ConfigItem
	err = json.Unmarshal(doc.Content, &content)
	if err != nil {
		return nil, err
	}
	return map[string]any{"content": content, "id": doc.ID, "Name": doc.Name}, nil
}
func configDELETE(db *sqlx.DB, name string) (bool, error) {
	// First, let's check if the config even exists.
	stmt := "SELECT name FROM documents WHERE name=$1"
	doc := &Document{}
	err := db.Get(doc, stmt, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	deleteStmt := "DELETE FROM documents WHERE name=$1"
	_, err = db.Exec(deleteStmt, name)
	if err != nil {
		return false, err
	}
	return true, nil
}

func configPUT(db *sqlx.DB, name string, data []config.ConfigItem) error {
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
