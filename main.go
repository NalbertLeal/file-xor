package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func getPasswordFromArgs() string {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: Indique a senha")
		os.Exit(1)
	}

	return os.Args[1]
}

func getCurrentPath() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return pwd, err
	}
	return pwd, nil
}

func getAllFilesNameIntoDirectory(directoryPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return files, err
	}

	files = files[1:]

	return files, nil
}

func xorEncript(password []byte, binFile []byte) []byte {
	var result []byte

	counter := 0
	for _, b := range binFile {
		xored := b ^ password[counter]
		result = append(result, xored)

		counter++
		if counter == len(password) {
			counter = 0
		}
	}

	return result
}

func readFileBytes(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {

		return data, err
	}

	return data, nil
}

func writeFileBytes(filePath string, content []byte) error {
	err := ioutil.WriteFile(filePath, content, 0777)
	if err != nil {
		return err
	}

	return nil
}

func filesAreEncripted(files []string) (bool, string) {
	for _, file := range files {
		match, _ := regexp.MatchString("files_are_encripted.txt", file)
		if match {
			return true, file
		}
	}
	return false, ""
}

func createOrDeleteFilesEncriptedIndicator(current string, files []string) {
	filesEncripted, path := filesAreEncripted(files)
	if filesEncripted {
		err := os.Remove(path) // remove a single file
		if err != nil {
			os.Exit(1)
		}
	} else {
		writeFileBytes(current+"/files_are_encripted.txt", []byte{})
	}
}

func encriptFiles(files []string, password string) {
	binaryPassword := []byte(password)

	for _, file := range files {
		match, _ := regexp.MatchString("files_are_encripted.txt", file)
		if match {
			continue
		}

		binFile, err := readFileBytes(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		result := xorEncript(binaryPassword, binFile)

		writeFileBytes(file, result)
	}
}

func main() {
	password := getPasswordFromArgs()

	current, err := getCurrentPath()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	files, err := getAllFilesNameIntoDirectory(current)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	createOrDeleteFilesEncriptedIndicator(current, files)

	encriptFiles(files, password)
}
