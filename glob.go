package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FilesProps []struct {
	AbsPath string
	Name    string
	Size    int64
}

func GlobDir(path_dir string, ext string) (FilesProps, error) {

	fp := FilesProps{}
	e := strings.ToLower(ext)

	files, err := os.ReadDir(path_dir)
	if err != nil {
		return fp, err
	}

	abs_path, err := filepath.Abs(path_dir)
	if err != nil {
		return fp, err
	}

	for _, j := range files {

		absPath := strings.ToLower(filepath.Join(abs_path, j.Name()))
		info, err := os.Stat(absPath)
		if ext == "*" {
			fp = Add(fp, absPath, info)
		} else {
			if strings.Contains(absPath, e) {

				if err != nil {
					fmt.Println(err)
				}
				fp = Add(fp, absPath, info)

			}

		}

	}

	return fp, nil
}

func (g FilesProps) Count() int {
	return len(g)
}
