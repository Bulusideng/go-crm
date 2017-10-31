package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

func ExcelReadWrite() {
	//example type
	type structTest struct {
		IntVal     int     `xlsx:"0"`
		StringVal  string  `xlsx:"1"`
		FloatVal   float64 `xlsx:"2"`
		IgnoredVal int     `xlsx:"-"`
		BoolVal    bool    `xlsx:"4"`
	}
	structVal := structTest{
		IntVal:     16,
		StringVal:  "heyheyhey :)!",
		FloatVal:   3.14159216,
		IgnoredVal: 7,
		BoolVal:    true,
	}
	//create a new xlsx file and write a struct
	//in a new row
	f := xlsx.NewFile()
	sheet, _ := f.AddSheet("TestRead")
	row := sheet.AddRow()
	row.WriteStruct(&structVal, -1)

	//read the struct from the same row
	readStruct := &structTest{}
	err := row.ReadStruct(readStruct)
	if err != nil {
		fmt.Println(readStruct)
	} else {
		fmt.Printf("Success\n")
	}
	f.Save("./data/t1.xlsx")
}

func zipdb() {
	ExcelReadWrite()
	fn := strings.Replace(time.Now().Format(time.RFC3339), ":", "-", -1)
	// Files to Zip
	files := []string{"./data/data.db", "./data/data.xlsx"}
	output := "./data/DB_" + fn + ".zip"

	err := ZipFiles(output, files)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Zipped File: " + output)
}

// ZipFiles compresses one or many files into a single zip archive file
func ZipFiles(filename string, files []string) error {

	newfile, err := os.Create(filename)

	if err != nil {
		return err
	}
	defer newfile.Close()

	zipWriter := zip.NewWriter(newfile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// Get the file information
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, zipfile)
		if err != nil {
			return err
		}
	}
	return nil
}
