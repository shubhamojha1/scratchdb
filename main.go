package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	// SaveData("./sample.txt", [])
	var str string
	fmt.Scanln(&str)

	data := []byte(str)
	filepath := "./sample.txt"

	err := SaveData2(filepath, data)
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

func SaveData2(path string, data []byte) error {
	rand.Seed(time.Now().UnixNano())

	// randomInt := rand.Intn(100-1+1) + 1
	tmp := fmt.Sprintf("%s.tmp.%d", path, 1)
	fp, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.Write(data)
	if err != nil {
		os.Remove(tmp)
		return err
	}
	// fp.Close()
	// return os.Rename(tmp, path)
	err = os.Rename(tmp, path)
	if err != nil {
		return err
	}
	return nil

}
