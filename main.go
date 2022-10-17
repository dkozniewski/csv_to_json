package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {

	// Flags
	csvDir := flag.String("csv", "./csv_files", "a string")
	jsonDir := flag.String("json", "./json_files", "a string")

	flag.Parse()

	// Make channel
	c := make(chan string)

	// Read .csv files
	dirWithCsv, err := GlobDir(toAbsPath(*csvDir), ".csv")
	if err != nil {
		log.Fatal(err)
	}

	// Check exist Dir for output result .json
	checkExistsDir(toAbsPath(*jsonDir))

	// Transform CSV to JSON files
	for _, path := range dirWithCsv {
		fmt.Printf("Read file: %v\n", path.Name)
		go etl(path.AbsPath, path.Name, *jsonDir, c)
	}

	// Channel result time
	for i := 0; i < dirWithCsv.Count(); i++ {
		fmt.Println(<-c)
	}

}

func etl(absPath string, nameFile string, jsonDir string, c chan string) {

	start := time.Now()

	var line int
	var recordHead []string

	// Open CSV file
	f, err := os.OpenFile(absPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("I can't open file: %v\n", err)
	}
	defer f.Close()

	// Generate JSON file like CSV file
	fJson, fileNameJson, err := createJsonFile(jsonDir, nameFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fJson.Close()

	// Scan CSV file
	scan := bufio.NewScanner(f)

	// Scanning CSV line by line
	for scan.Scan() {
		text := scan.Text()
		r := csv.NewReader(strings.NewReader(text))
		record, err := r.Read()
		if err != nil {
			fmt.Printf("I can't read this: %v\n", err)
			line++
			os.Remove(fileNameJson)
			break
		}
		if line == 0 {
			recordHead = record
		}

		if line > 0 {
			fJson.Write(createJson(recordHead, record))
			fJson.Write([]byte("\n"))
		}

		line++

	}

	// Create time resume task
	timeTask := strconv.Itoa(int(time.Since(start).Seconds()))

	// Send string -> channel
	c <- strings.Join([]string{"Time:", timeTask, "seconds, Copy", nameFile, strconv.Itoa(line), "rows"}, " ")

}

func createJsonFile(jsonDir string, nameFile string) (*os.File, string, error) {

	fileNameJson := filepath.Join(jsonDir, strings.Replace(nameFile, ".csv", ".json", 1))

	// Check exists path, if exists => delete file.
	checkExistsPath(fileNameJson)

	// Open JSON file
	fJson, err := os.OpenFile(fileNameJson, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, "", err
	}

	return fJson, fileNameJson, nil
}

func checkExistsPath(path string) error {
	// if file exists => delete file
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	return os.Remove(path)

}

func checkExistsDir(path string) error {
	// if directory exists => nothing
	if _, err := os.Stat(path); err != nil {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil

}

func createJson(h []string, r []string) []byte {

	// Create map with JSON fileds
	m := make(map[string]interface{})

	// Generate JSON data
	for i := 0; i < len(h); i++ {
		m[h[i]] = convertNumber(r[i])
	}

	// Encoding data to byte
	bs, _ := json.Marshal(m)
	return bs
}

func convertNumber(str string) interface{} {

	// Convert string to int/float number
	i, err := strconv.Atoi(str)
	if err == nil {
		return i
	}
	f, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return f
	}

	return str

}

func toAbsPath(path string) string {

	// return absolut path
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("I can't generate absolute path: %v", err)
	}

	return absPath
}

func Add(fp FilesProps, absPath string, info fs.FileInfo) FilesProps {
	fp = append(fp, struct {
		AbsPath string
		Name    string
		Size    int64
	}{absPath, info.Name(), info.Size()})
	return fp
}
