package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//OpenLog open log file for writing
func OpenLog(path string, filename string) (*os.File, error) {
	if path == "" {
		//use current directory
		curdir, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
		dir, err := filepath.Abs(filepath.Dir(strings.Join([]string{curdir, filename}, "/")))
		if err != nil {
			log.Fatal(err)
		}
		return os.OpenFile(strings.Join([]string{dir, strings.Join([]string{filename, ".log"}, "")}, "/"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	}
	return &os.File{}, errors.New("Unable to determine path for writing")
}
