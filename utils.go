/*
Version 1.00
Date Created: 2022-11-15
Copyright (c) 2022, Akshay Singh Kanawat
Author: Akshay Singh Kanawat
*/
package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func readIndexFile(fileName string) (data map[string]string) {
	data = make(map[string]string)
	if getFileSize(fileName) == 0 {
		return
	}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error while reading index file: ", err)
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Error while reading index file: ", err)
	}
	return
}

func readDataFile(fileName string) (data map[string]string) {
	data = make(map[string]string)
	if getFileSize(fileName) == 0 {
		return
	}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error while reading data file: ", err)
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Error while reading data file: ", err)
	}
	return
}

func getFileSize(filename string) int64 {
	fInfo, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	fsize := fInfo.Size()
	//fmt.Printf("The file size is %d bytes\n", fsize)
	return fsize
}

func getLatestFile(configurableSourceDir string) (string, bool) {
	dir := fmt.Sprintf("%s/%s", configurableSourceDir, "dbFiles")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	var modTime time.Time
	var names []string
	for _, fi := range files {
		if fi.Mode().IsRegular() {
			if !fi.ModTime().Before(modTime) {
				if fi.ModTime().After(modTime) {
					modTime = fi.ModTime()
					names = names[:0]
				}
				names = append(names, fi.Name())
			}
		}
	}
	//If file exists return true
	if len(names) > 0 {
		fmt.Println(modTime, names)
		return names[len(names)-1], true
	}
	return "", false
}
