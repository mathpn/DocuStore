package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/errgroup"
)

func createTables(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS doc_summaries (doc_id TEXT PRIMARY KEY, summary BLOB)")
	return err
}

func InsertDocument(db *sql.DB, docSummary *DocSummary) (int64, error) {
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
	rows, err := insertDocTransaction(tx, docSummary)
	return rows, err
}

func insertDocTransaction(tx *sql.Tx, docSummary *DocSummary) (int64, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(docSummary)
	if err != nil {
		return 0, err
	}
	blob := buffer.Bytes()
	out, err := tx.Exec(
		"INSERT OR IGNORE INTO doc_summaries (doc_id, summary) VALUES (?, ?)",
		docSummary.DocID,
		blob,
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

func ListDocuments(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT doc_id FROM doc_summaries")
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

func LoadDocument(db *sql.DB, docID string) (*DocSummary, error) {
	row := db.QueryRow("SELECT summary FROM doc_summaries WHERE doc_id = ?", docID)
	// fmt.Printf("%+v\n", row)
	var blob []byte
	err := row.Scan(&blob)
	// fmt.Printf("%+v\n", blob)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(blob)
	decoder := gob.NewDecoder(buffer)
	docSummary := DocSummary{}
	decoder.Decode(&docSummary)
	return &docSummary, err
}

func LoadDocuments(ctx context.Context, db *sql.DB, docIDs ...string) ([]*DocSummary, error) {
	errs, ctx := errgroup.WithContext(ctx)
	out := make([]*DocSummary, len(docIDs))
	// var wg sync.WaitGroup
	for i := 0; i < len(docIDs); i++ {
		current := i
		errs.Go(
			func() error {
				// defer wg.Done()
				doc, err := LoadDocument(db, docIDs[current])
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
