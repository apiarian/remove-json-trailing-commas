package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s filename.json", os.Args[0])
	}

	filename := os.Args[1]

	data, err := readFile(filename)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	if quickCheckJSON(data) == nil {
		log.Printf("already valid json")
		return
	}

	processed, err := cleanTrailingCommas(data)
	if err != nil {
		log.Fatalf("failed to clean file: %v", err)
	}

	if err := quickCheckJSON(processed); err != nil {
		log.Fatalf("json still invalid: %v\n%s", err, processed)
	}

	bkpFilename, err := writeBackup(filename, data)
	if err != nil {
		log.Printf("failed to write backup: %v", err)
	}

	err = writeFile(filename, processed)
	if err != nil {
		log.Fatalf("failed to write processed data: %v", err)
	}

	os.Remove(bkpFilename)
}

func readFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

func writeFile(filename string, data []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		f.Close()
		return err
	}

	return f.Close()
}

func writeBackup(filename string, data []byte) (string, error) {
	bkpFilename := filename + ".bkp"

	f, err := os.Create(bkpFilename)
	if err != nil {
		return "", err
	}

	_, err = f.Write(data)
	if err != nil {
		f.Close()
		return "", err
	}

	return bkpFilename, f.Close()
}

func quickCheckJSON(data []byte) error {
	var stuff interface{}
	return json.Unmarshal(data, &stuff)
}

func cleanTrailingCommas(data []byte) ([]byte, error) {
	re, err := regexp.Compile(`,(\s*([}\]]|$))`)
	if err != nil {
		return nil, err
	}

	return re.ReplaceAll(data, []byte("$1")), nil
}
