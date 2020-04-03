package carbone

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func createTmpFile(newPath string, newContent string) (string, string) {
	pathFile := "tests"
	if newPath != "" {
		pathFile = newPath
	}
	filenamePath := ""
	time := getUnixTime()
	filename := "tmp." + time + ".test.html"
	filenamePath = filepath.Join(pathFile, filename)
	content := "<!DOCTYPE html><html><body> {d.name} date:" + time + "</body></html>"
	if newContent != "" {
		content = newContent
	}
	err := ioutil.WriteFile(filenamePath, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
	return filenamePath, filename
}

func deleteTmpFiles(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*tmp*test*"))
	if err != nil {
		return err
	}
	log.Printf("Cleaning '%v' directory: deleting %d tmp files.", dir, len(files))
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func getUnixTime() string {
	now := time.Now().Unix()
	return strconv.FormatInt(now, 10)
}
