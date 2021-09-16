/*
   MIT License

   Copyright (c) 2021 Ryan SVIHLA

   Permission is hereby granted, free of charge, to any person obtaining a copy
   of this software and associated documentation files (the "Software"), to deal
   in the Software without restriction, including without limitation the rights
   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
   copies of the Software, and to permit persons to whom the Software is
   furnished to do so, subject to the following conditions:

   The above copyright notice and this permission notice shall be included in all
   copies or substantial portions of the Software.

   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
   SOFTWARE.
*/
package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var verbose bool

func main() {
	help := flag.Bool("help", false, "show command line help")
	flag.BoolVar(&verbose, "v", false, "show verbose extraction")
	delarc := flag.Bool("del", false, "delete the files that are successfully extracted")
	flag.Parse()
	if flag.NArg() == 0 || *help {
		flag.Usage()
		//I throw an error here because I do not like to have scripts silently fail if I pass the wrong args
		os.Exit(1)
	}
	cleanup := func(s string) {}
	if *delarc {
		cleanup = func(path string) {
			err := os.Remove(path)
			if err != nil {
				log.Printf("WARN unable to delete %v with error %v\n", path, err)
			} else if verbose {
				log.Printf("INFO file %v deleted\n", path)
			}
		}
	}
	dir := flag.Args()[0]
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("unable to access path %q: %v", path, err)
		}
		if info.Mode().IsRegular() {
			ext := filepath.Ext(path)
			switch ext {
			case ".zip":
				err = Unzip(path)
				if err != nil {
					log.Printf("ERROR unable to execute unzip on file '%v' with error '%v'", path, err)
				} else {
					cleanup(path)
				}
			case ".tar":
				err = UnTar(path)
				if err != nil {
					log.Printf("ERROR unable to execute untar on file '%v' with error '%v'", path, err)
				} else {
					cleanup(path)
				}
			case ".gz":
				newFile, err := UnGzip(path)
				if err != nil {
					log.Printf("ERROR unable to execute ungzip on file '%v' with error '%v'", path, err)
				} else {
					cleanup(path)
				}
				newExt := filepath.Ext(newFile)
				if newExt == ".tar" {
					err = UnTar(newFile)
					if err != nil {
						log.Printf("ERROR unable to execute untar on file '%v' with error '%v'", path, err)
					} else {
						cleanup(newFile)
					}
				}
			case ".tgz":
				newFile, err := UnGzip(path)
				if err != nil {
					log.Printf("ERROR unable to execute ungzip on file '%v' with error '%v'", path, err)
				} else {
					cleanup(path)
				}
				err = UnTar(newFile)
				if err != nil {
					log.Printf("ERROR unable to execute untar on file '%v' with error '%v'", path, err)
				} else {
					cleanup(newFile)
				}
			default:
				if verbose {
					log.Printf("INFO skipping %v as exentions is %v", path, ext)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir, err)
		os.Exit(2)
	}
}

func writeFile(path string, r io.Reader) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to write file %v with error %v", path, err)
	}
	defer f.Close()
	written, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	if verbose {
		log.Printf("INFO %v bytes written for file %v\n", written, path)
	}
	return nil
}

//UnTar a tarball and places all the files next to it
func UnTar(path string) error {
	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open path '%v' due to error '%v'", path, err)
	}
	defer r.Close()
	tarReader := tar.NewReader(r)
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("unable to create directory %v with error %v", dir, err)
	}
	// Iterate through the files in the archive
	// each file found is written out withthe full path
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			//we are done now
			log.Printf("file %v extracted\n", path)
			return nil
		}
		if err != nil {
			return fmt.Errorf("unable to read tar files for file %v with error %v", path, err)
		}
		newPath := filepath.Join(dir, header.Name)
		if header.Typeflag == tar.TypeReg {
			err = writeFile(newPath, tarReader)
			if err != nil {
				return fmt.Errorf("unable to write file %v in tar to %v with error %v", header.Name, newPath, err)
			}
			continue
		}

		if header.Typeflag == tar.TypeDir {
			err := os.MkdirAll(newPath, 0755)
			if err != nil {
				return fmt.Errorf("unable to create directory %v with error %v", dir, err)
			}
			fmt.Printf("created dir %v\n", newPath)
			continue
		}
		log.Printf("unknown file type %v for file %v", header, newPath)
	}
}

//Unzip unzips the file specified at the path and writes out all the files found next to it
func Unzip(path string) error {
	zipReader, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("unable to open path '%v' due to error '%v'", path, err)
	}
	defer zipReader.Close()
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("unable to create directory %v with error %v", dir, err)
	}

	// Iterate through the files in the archive
	// each file found is written out withthe full path
	for _, f := range zipReader.File {
		fileInZip, err := f.Open()
		if err != nil {
			return fmt.Errorf("unable to open file %v in zip with error %v", f.Name, err)
		}
		newFilePath := filepath.Join(dir, f.Name)
		err = writeFile(newFilePath, fileInZip)
		if err != nil {
			fileInZip.Close()
			return fmt.Errorf("unable to write file %v with error %v", newFilePath, err)
		} else {
			fileInZip.Close()
		}
	}
	return nil
}

//UnGzip unzips the gzip file specified at the path and writes out the file
func UnGzip(path string) (string, error) {
	r, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("unable to open %v with error %v", path, err)
	}

	defer r.Close()
	zr, err := gzip.NewReader(r)
	if err != nil {
		return "", fmt.Errorf("unable to open %v with error %v", path, err)
	}

	newFilePath := strings.TrimSuffix(path, filepath.Ext(path))
	err = writeFile(newFilePath, zr)

	if err != nil {
		return "", fmt.Errorf("unable to write file %v in gzip to %v with error %v", path, newFilePath, err)
	}
	defer func() {
		if err := zr.Close(); err != nil {
			log.Println(err)
		}
	}()
	return newFilePath, nil
}
