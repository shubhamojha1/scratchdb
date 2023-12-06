package main

import (
	"fmt"
	"os"
)

func main() {
	// SaveData("./sample.txt", [])
	data := []byte("Hello world")
	filepath := "sample.txt"

	err := SaveData(filepath, data)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Data saved successfully")
	}
}

func SaveData(path string, data []byte) error {
	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)

	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.Write(data)
	return err
}
