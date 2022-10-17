package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateJsonFile(t *testing.T) {
	nf := "testFileToDelete.json"
	nd := "jsonTestDir"
	os.Mkdir(nd, fs.ModePerm)

	f, _, err := createJsonFile(nd, nf)
	if err != nil {
		log.Fatalf("Error : %v", err)
	}
	f.Close()
	os.Remove(filepath.Join(nd, nf))
	os.Remove(nd)

}

func TestCheckExistsDir(t *testing.T) {

	nf := "testFoldesToDelete"
	os.Remove(nf)
	checkExistsDir(nf)

	fs, err := os.Stat(nf)
	if err != nil {
		log.Fatalf("The file: %v could not be created.", fs.Name())
	}
	os.Remove(nf)
}

func TestGlobDir(t *testing.T) {

	path := "./csv_files"
	listFiles, _ := GlobDir(path, ".csv")

	if listFiles.Count() == 0 {
		t.Errorf("Expected != 0, got 0")
	}
}

func TestAdd(t *testing.T) {

	p := "path"
	n := "testName"
	var s int64 = 1

	fp := FilesProps{}

	fp = append(fp, struct {
		AbsPath string
		Name    string
		Size    int64
	}{p, n, s})

	if fp[0].AbsPath != "path" {
		t.Errorf("Expected path, got: %v", fp[0].AbsPath)
	}
	if fp[0].Name != "testName" {
		t.Errorf("Expected testName, got: %v", fp[0].Name)
	}
	if fp[0].Size != 1 {
		t.Errorf("Expected path, got: %v", fp[0].Size)
	}

}
