package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// Configuration : structure for config.json
type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
	Posts        string
}

var config Configuration

func init() {
	loadConfig()
	loadPosts()
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
}

func loadPosts() {
	files := filepaths(config.Posts)
	for i, filepath := range files {
		file, err := ioutil.ReadFile(filepath)
		if err != nil {
			panic(err)
		}
		temp := Post{
			Markdown: string(file),
			Title:    fileName(filepath),
			URL:      strconv.Itoa(i),
		}
		parseMarkdown(&temp)
		posts = append(posts, temp)
	}
}

func fileName(filepath string) string {
	info, err := os.Stat(filepath)
	if err != nil {
		panic(err)
	}
	return info.Name()
}

func filepaths(filePath string) []string {
	var files []string
	re, _ := regexp.Compile(`([a-zA-Z0-9\s_\\.\-\(\):])+(.md)$`)
	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		b := re.MatchString(path)
		if b {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
