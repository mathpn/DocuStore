package main

import (
	"encoding/gob"
	"log"
	"os"
	"sync"
)

var lock sync.Mutex

func SaveText(path string, text string) error {
	lock.Lock()
	defer lock.Unlock()
	err := os.WriteFile(path, []byte(text), 0644)
	return err
}

func LoadText(path string) string {
	lock.Lock()
	defer lock.Unlock()
	content, err := os.ReadFile(path)
	check(err) // TODO improve
	return string(content)
}

// SaveStruct saves a representation of v to the file at path.
func SaveStruct(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(v)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	return err
}

// LoadStruct loads the file at path into v.
func LoadStruct(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	return dec.Decode(v)
}
