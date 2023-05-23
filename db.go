package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/errgroup"
)

func NewDBConnection(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	err = createTables(db)
	return db, err
}

func createTables(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS documents (doc_id TEXT PRIMARY KEY, timestamp INTEGER, summary BLOB, content BLOB)")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS doc_timestamps ON documents (timestamp)")
	return err
}

func InsertDocument(db *sql.DB, docSummary *DocSummary, content string) (int64, error) {
	timestamp := time.Now().Unix()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	rows, err := insertDocTransaction(tx, docSummary, content, timestamp)
	return rows, err
}

func insertDocTransaction(tx *sql.Tx, docSummary *DocSummary, content string, timestamp int64) (int64, error) {
	var buffer bytes.Buffer
	var err error
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(docSummary)
	if err != nil {
		return 0, err
	}

	byteContent := []byte(content)
	blob := buffer.Bytes()
	out, err := tx.Exec(
		"INSERT OR IGNORE INTO documents (doc_id, summary, content, timestamp) VALUES (?, ?, ?, ?)",
		docSummary.DocID,
		blob,
		byteContent,
		timestamp,
	)
	if err != nil {
		return 0, err
	}
	rows, err := out.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, err
}

func GetLatestTimestamp(db *sql.DB) (int64, error) {
	row := db.QueryRow("SELECT coalesce(max(timestamp), 0) FROM documents")
	var timestamp int64
	err := row.Scan(&timestamp)
	return timestamp, err
}

func ListDocuments(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT doc_id FROM documents")
	if err != nil {
		return nil, err
	}
	var docIDs []string
	for rows.Next() {
		var docID string
		err := rows.Scan(&docID)
		if err != nil {
			return nil, err
		}
		docIDs = append(docIDs, docID)
	}
	return docIDs, nil
}

func LoadText(db *sql.DB, docID string) (string, error) {
	row := db.QueryRow("SELECT content FROM documents WHERE doc_id = ?", docID)
	var byteContent []byte
	err := row.Scan(&byteContent)
	if err != nil {
		return "", err
	}
	content := string(byteContent)
	return content, nil
}

func LoadDocSummary(db *sql.DB, docID string) (*DocSummary, error) {
	row := db.QueryRow("SELECT summary FROM documents WHERE doc_id = ?", docID)
	var blob []byte
	err := row.Scan(&blob)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(blob)
	decoder := gob.NewDecoder(buffer)
	docSummary := DocSummary{}
	decoder.Decode(&docSummary)
	return &docSummary, err
}

func LoadDocSummaries(ctx context.Context, db *sql.DB, docIDs ...string) ([]*DocSummary, error) {
	errs, ctx := errgroup.WithContext(ctx)
	out := make([]*DocSummary, len(docIDs))
	for i := 0; i < len(docIDs); i++ {
		current := i
		errs.Go(
			func() error {
				doc, err := LoadDocSummary(db, docIDs[current])
				if err != nil {
					return err
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					out[current] = doc
					return nil
				}
			},
		)
	}
	// go func() {
	// 	errs.Wait()
	// }()
	return out, errs.Wait()
}
