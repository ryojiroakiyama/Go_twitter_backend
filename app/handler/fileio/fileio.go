// Package fileio provides basic functions for file operations.
package fileio

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

//WriteToFile creates a filePath file and writes contents to the file.
//If successful, WriteToFile returns nil error.
//Else if faulse, WriteToFile returns any error encountered.
//Maybe same as os.WriteFile(filePath, contents, 066).
func WriteToFile(filePath string, contents []byte) (err error) {
	dstFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ToFile: %v", err)
	}
	defer func() {
		if cerr := dstFile.Close(); cerr != nil {
			err = fmt.Errorf("ToFile: %v", cerr)
		}
		if err != nil && dstFile != nil {
			os.Remove(dstFile.Name())
		}
	}()
	_, werr := dstFile.Write(contents)
	if werr != nil {
		return fmt.Errorf("ToFile: %v", werr)
	}
	err = dstFile.Sync()
	return
}

//WriteToTmpFile creates a tmpprary file and write contents to the file.
//If successful, writeTmpFile returns a file path and nil error, else returns error.
func WriteToTmpFile(src io.Reader) (filePath string, err error) {
	tmpfile, err := os.CreateTemp("", "")
	if err != nil {
		return "", fmt.Errorf("ToTmpFile: %v", err)
	}
	filePath = tmpfile.Name()
	defer func() {
		if cerr := tmpfile.Close(); cerr != nil {
			err = fmt.Errorf("ToTmpFile: %v", cerr)
		}
		if err != nil && filePath != "" {
			os.Remove(filePath)
		}
	}()
	_, err = io.Copy(tmpfile, src)
	if err != nil {
		return "", fmt.Errorf("ToTmpFile: %v", err)
	}
	err = tmpfile.Sync()
	return
}

//FileContents returns a byte array readed
//from the file named filePath or any error encountered.
func FileContents(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("FileBytes: %v", err)
	}
	defer file.Close()
	srcBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("FileBytes: %v", err)
	}
	return srcBytes, nil
}

//BindFiles binds contents of files stocked
//in srcNames to the file named dstName or any error encountered.
//The order to bind is the order of the slice index.
func BindFiles(srcPath []string, dstPath string) (err error) {
	dstfile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("BindFiles: %v", err)
	}
	defer func() {
		if cerr := dstfile.Close(); cerr != nil {
			err = fmt.Errorf("BindFiles: %v", cerr)
		}
		if err != nil {
			os.Remove(dstPath)
		}
	}()
	for _, srcName := range srcPath {
		srcfile, err := os.Open(srcName)
		if err != nil {
			return fmt.Errorf("BindFiles: %v", err)
		}
		defer srcfile.Close()
		_, err = io.Copy(dstfile, srcfile)
		if err != nil {
			return fmt.Errorf("BindFiles: %v", err)
		}
	}
	err = dstfile.Sync()
	return
}
