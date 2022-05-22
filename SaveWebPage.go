package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

var urlStr *string
var path *string
var override *bool
var pageUrl *url.URL

func main() {
	parseParams()
	pageStr, err := getPage()
	checkErr(err)
	err = saveToFile(pageStr)
	checkErr(err)
	fmt.Printf("URL: %s\nPATH: %s\nOverride: %t\n", *urlStr, *path, *override)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func getPage() (string, error) {
	result, err := http.Get(pageUrl.String())
	if err != nil {
		return "", fmt.Errorf("Error getting page \"%s\": %w", pageUrl.String(), err)
	}
	defer result.Body.Close()
	content, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading from response body: %w", err)
	}
	return string(content), nil
}

func saveToFile(page string) error {
	_, err := os.Stat(*path)
	if err == nil && !*override {
		return fmt.Errorf("Cannot override existing file: %s", *path)
	}
	err = os.WriteFile(*path, []byte(page), 0444)
	//fmt.Println(page)
	return nil
}

func parseParams() {
	urlStr = flag.String("URL", "", "URL of web page to be saved")
	path = flag.String("path", "", "Full path to the file where to save page")
	override = flag.Bool("override", false, "Override file if it exists?")
	flag.Parse()
	var err error
	pageUrl, err = url.Parse(*urlStr)
	if err != nil || *urlStr == "" {
		log.Println(fmt.Errorf("Bad URL: \"%s\". Error parsing URL: %w", *urlStr, err))
		printUsage()
		return
	}
	dir := filepath.Dir(*path)
	_, err = os.Stat(dir)
	if *path == "" || errors.Is(err, os.ErrNotExist) {
		log.Println(fmt.Errorf("File undefined or directory does not exist: %w", err))
		printUsage()
		return
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("SaveWebPage <params>")
	fmt.Println("where <params> are:")
	flag.PrintDefaults()
}
