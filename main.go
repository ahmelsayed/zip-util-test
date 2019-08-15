package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var dirName string

func init() {
	flag.StringVar(&dirName, "dir", "", "dir to zip")
	flag.Parse()
}

func main() {
	output := "output.zip"
	// baseFolder := "/tmp/foo/"
	baseFolder := dirName
	fmt.Println(baseFolder)

	if err := zipWriter(baseFolder, output); err != nil {
		panic(err)
	}
	fmt.Println("Zipped File:", output)
}

func zipWriter(dir, outputFile string) error {
	// Get a Buffer to Write To
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	err = addFiles(w, dir, "")

	if err != nil {
		return err
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}

func addFiles(w *zip.Writer, basePath, baseInZip string) error {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := basePath + file.Name()
			fileToZip, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer fileToZip.Close()

			header, err := zip.FileInfoHeader(file)
			if err != nil {
				return err
			}

			// Using FileInfoHeader() above only uses the basename of the file. If we want
			// to preserve the folder structure we can overwrite this with the full path.
			header.Name = baseInZip + file.Name()

			// Change to deflate to gain better compression
			// see http://golang.org/pkg/archive/zip/#pkg-constants
			header.Method = zip.Deflate

			writer, err := w.CreateHeader(header)
			if err != nil {
				return err
			}
			_, err = io.Copy(writer, fileToZip)

			if err != nil {
				return err
			}
		} else if file.IsDir() {
			// Recurse
			newBase := basePath + file.Name() + "/"
			// fmt.Println("Recursing and Adding SubDir: " + file.Name())
			// fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
	return nil
}
