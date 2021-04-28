package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
)

func URLToFile(url, path string) error {
	body, err := BodyFromURL(url)
	if err != nil {
		return err
	}

	return WriteToFile(path, body)
}

func WriteToFile(path, text string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
		return err
	}

	if _, err := f.Write([]byte(text)); err != nil {
		log.Fatal(err)
		return err
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func BodyFromURL(url string) (string, error) {
	bodyBuffer, err := BufferFromURL(url)
	if err != nil {
		return "", err
	}

	return string(bodyBuffer), nil
}

func BufferFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return make([]byte, 1), err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// some random []byte
		return make([]byte, 1), err
	}

	return body, nil
}

func TrimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func Handle(err error, message string) {
	if err != nil {
		fmt.Println(message)
		os.Exit(1)
	}
}

func Log(message string) {
	log.Println(message)
}

func CreateEmptyFile(name string) {
	d := []byte("")
	check(ioutil.WriteFile(name, d, 0644))
}

func Mkdir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}

	return nil
}

func CloneRepoToDir(url, path string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: url,
	})

	return err
}

func DirsInPath(path string) []string {
	files, err := ioutil.ReadDir(path)
	Handle(err, "Error reading catalogs")

	var packages []string

	for _, f := range files {
		if f.IsDir() && !(strings.HasPrefix(f.Name(), ".")) {
			packages = append(packages, f.Name())
		}
	}

	return packages
}

// Internal

func check(e error) {
	if e != nil {
		panic(e)
	}
}
