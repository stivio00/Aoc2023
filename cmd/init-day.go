// Quick and dirty script to build a aoc basis structure
// Stephen Krol

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	year     = "2022"
	basisUrl = "https://adventofcode.com/%s"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: init-day [dayNumber]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func copyFile(src string, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0744)
}

func buildDirectoryStructure(dayNumber string) error {
	// Check if folder exist, if not create dir src/day-<dayNumber>
	dirName := "src/day-" + dayNumber
	err := os.MkdirAll(dirName, 0750)
	if err != nil {
		return err
	}

	// Check if .go file exists, if not add main.go from template dir
	if _, err := os.Stat(dirName + "/main.go"); err == nil {
		// path/to/whatever exists
		fmt.Println("main.go already exists: skipping.")
	} else if errors.Is(err, os.ErrNotExist) {
		//Copy file from template
		fmt.Println("Copy main.go from template.")
		copyFile("src/day-template/main.go", dirName+"/main.go")
	} else {
		return err
	}

	// check if problem.txt exists if not download problem.txt
	// TODO

	// check if input.txt exists if not download input.txt
	if _, err := os.Stat(dirName + "/input.txt"); err == nil {
		// path/to/whatever exists
		fmt.Println("input.txt already exists: skipping.")
	} else if errors.Is(err, os.ErrNotExist) {
		downloadFile("/day/"+dayNumber+"/input", dirName+"/input.txt")
	} else {
		return err
	}
	// check if input2.txt exists if not download input2.txt
	if _, err := os.Stat(dirName + "/input2.txt"); err == nil {
		// path/to/whatever exists
		fmt.Println("input2.txt already exists: skipping.")
	} else if errors.Is(err, os.ErrNotExist) {
		downloadFile("/day/"+dayNumber+"/input2", dirName+"/input2.txt")
	} else {
		return err
	}

	return nil
}

func downloadFile(path, file string) error {
	dat, err := os.ReadFile("token.txt")
	if err != nil {
		log.Fatal("Missing token.txt file.")
	}

	session := string(dat)

	baseUrl := fmt.Sprintf(basisUrl, year)
	url := baseUrl + path
	client := &http.Client{}
	fmt.Println("Hitting " + url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Fail to build request")
	}

	req.AddCookie(&http.Cookie{Name: "session", Value: session})
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Fail to request data")
	}

	fmt.Println("Response: " + res.Status)

	if res.StatusCode != 200 {
		fmt.Println("skipping")
		return nil
	}

	f, err := os.Create(file)
	if err != nil {
		log.Fatal("Unable to open the file")
	}
	defer f.Close()

	fmt.Print("Writing input file ")
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}
	fmt.Println(" OK")

	return nil
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("dayNumber missing")
		os.Exit(1)
	}

	dayNumber := args[0]
	fmt.Printf("Creating day #%s\n", dayNumber)

	err := buildDirectoryStructure(dayNumber)
	if err != nil {
		log.Fatal(err)
	}
}
