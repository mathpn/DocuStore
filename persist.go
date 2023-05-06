package main

import (
	"os"
	"sync"
)

var lock sync.Mutex

func SaveText(path string, text string) {
	lock.Lock()
	defer lock.Unlock()
	err := os.WriteFile(path, []byte(text), 0644)
	check(err) // TODO improve
}

func LoadText(path string) string {
	lock.Lock()
	defer lock.Unlock()
	content, err := os.ReadFile(path)
	check(err) // TODO improve
	return string(content)
}

// SaveStruct saves a representation of v to the file at path.
func SaveStruct(path string, v interface{}, marshal func(interface{}) ([]byte, error)) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := marshal(v)
	if err != nil {
		return err
	}
	_, err = f.Write(r)
	return err
}

// LoadStruct loads the file at path into v.
func LoadStruct(path string, v interface{}, unmarshal func([]byte, interface{}) error) error {
	lock.Lock()
	defer lock.Unlock()
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return unmarshal(bytes, v)
}
