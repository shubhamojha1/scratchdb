package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	// SaveData("./sample.txt", [])
	fmt.Println("Type something...")
	var str string
	fmt.Scanln(&str)

	data := []byte(str)
	filepath := "./sample.txt"

	err := SaveData3(filepath, data)
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

	randomInt := rand.Intn(100-1+1) + 1
	tmp := fmt.Sprintf("%s.tmp.%d", path, randomInt)
	fp, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}

	_, err = fp.Write(data)
	if err != nil {
		os.Remove(tmp)
		return err
	}
	fp.Close() // defer will not work here, need to close file before renaming it.
	return os.Rename(tmp, path)
}

func SaveData3(path string, data []byte) error {
	rand.Seed(time.Now().UnixNano())

	randomInt := rand.Intn(100-1+1) + 1
	tmp := fmt.Sprintf("%s.tmp.%d", path, randomInt)
	fp, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}

	_, err = fp.Write(data)
	if err != nil {
		os.Remove(tmp)
		return err
	}

	err = fp.Sync() // linux fsync {flushes data to the disk before renaming it} https://youtu.be/JK2ZIx8jRu4?si=TqsYujUVBrFxKJVB
	if err != nil {
		os.Remove(tmp)
		return err
	}
	fp.Close()
	return os.Rename(tmp, path)

}

func LogCreate(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0664)
}

func LogAppend(fp *os.File, line string) error {
	buf := []byte(line)
	buf = append(buf, '\n')
	_, err := fp.Write(buf)
	if err != nil {
		return err
	}
	return fp.Sync()
}
