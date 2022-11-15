/*
Version 1.00
Date Created: 2022-11-15
Copyright (c) 2022, Akshay Singh Kanawat
Author: Akshay Singh Kanawat
*/
package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type dbConnection struct {
	StorageLoaction string
}

type DbClient interface {
	Insert(key, value string) error
	Update(key, value string) error
	Get(key string) (string, error)
}

var IndexFile map[string]string
var DataFile map[string]string

func (db *dbConnection) Insert(key, value string) error {
	//Get the latest file from the disk folder(configurable)
	sourceDir := fmt.Sprintf("%s/%s", db.StorageLoaction, "dbFiles")
	t := time.Now().Unix()
	location := fmt.Sprintf("%s/%d.json", sourceDir, t)
	latestFile, exists := getLatestFile(db.StorageLoaction)
	if exists {
		fileSize := getFileSize(fmt.Sprintf("%s/%s", sourceDir, latestFile))
		if fileSize >= 1000000 {
			_, err := os.Create(location)
			if err != nil {
				return err
			}
			latestFile = location
			//create new file and insert
		} else {
			location = fmt.Sprintf("%s/%s", sourceDir, latestFile)
		}
	} else {
		err := ioutil.WriteFile(location, nil, 0644)
		if err != nil {
			panic(err)
		}
	}
	err := db.insertData(key, value, location)
	if err != nil {
		return err
	}
	db.insertIndex(key, location)
	return nil
}

func (db *dbConnection) Get(key string) (val string, err error) {
	var dataFile string
	fileName := fmt.Sprintf("%s/%s", db.StorageLoaction, "index.json")
	indexFile := readIndexFile(fileName)
	dataFile, indexAlreadyExists := indexFile[key]
	if !indexAlreadyExists {
		return "", errors.New("Key Not Found")
	} else {
		data := readDataFile(dataFile)
		//TODO: Further optimisation by Reading file chunks concurrently by using multiple go routines
		val, _ = data[key]
		return val, nil
	}
}

func (db *dbConnection) Update(key, value string) error {
	fileName := fmt.Sprintf("%s/%s", db.StorageLoaction, "index.json")
	indexFile := readIndexFile(fileName)
	dataFile, exists := indexFile[key]
	if !exists {
		return errors.New("Key does not exist")
	}
	err := db.updateData(key, value, dataFile)
	if err != nil {
		return err
	}
	return nil
}

func (db *dbConnection) insertIndex(key, location string) {
	//configurable index file
	fileName := fmt.Sprintf("%s/%s", db.StorageLoaction, "index.json")
	indexFile := readIndexFile(fileName)
	indexFile[key] = location
	dataToWrite, _ := json.Marshal(indexFile)
	err := ioutil.WriteFile(fileName, dataToWrite, 0644)
	if err != nil {
		fmt.Println("Error while writing index to file: ", err)
	}
}

func (db *dbConnection) insertData(key, value, fileName string) error {
	//configurable index file
	dataFile := readDataFile(fileName)
	_, keyExists := dataFile[key]
	if keyExists {
		return errors.New("Data already exists for input key")
	}
	dataFile[key] = value
	dataToWrite, _ := json.Marshal(dataFile)
	err := ioutil.WriteFile(fileName, dataToWrite, 0644)
	if err != nil {
		fmt.Sprintf("Data already exists for input key: %s")
		return err
	}
	return nil
}

func (db *dbConnection) updateData(key, value, fileName string) error {
	//configurable index file
	dataFile := readDataFile(fileName)
	dataFile[key] = value
	dataToWrite, _ := json.Marshal(dataFile)
	err := ioutil.WriteFile(fileName, dataToWrite, 0644)
	if err != nil {
		fmt.Println("Error while writing data to file: ", err)
		return err
	}
	return nil
}

func SetupdbConnection(storageLoaction string) DbClient {
	var client DbClient = &dbConnection{storageLoaction}
	//For encapsulation only, user of this package will be only able to access 3 exposed methods.
	//Further we can create as many receiver functions as we can, which will be not accessible to the outside world
	indexLocation := fmt.Sprintf("%s/index.json", storageLoaction)
	fmt.Println("file ", indexLocation)
	if _, err := os.Stat(indexLocation); errors.Is(err, os.ErrNotExist) {
		err = ioutil.WriteFile(indexLocation, nil, 0644)
		if err != nil {
			panic(err)
		}
	}
	return client
}
