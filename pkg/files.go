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
package pkg

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

//WalkDir looks for files to unarchive, it current supports tar, zip , tgz, gz and tar.gz files
func WalkDir(dir string, delarc bool, verbose bool) error {
	cleanup := func(s string) {}
	if delarc {
		cleanup = func(path string) {
			err := os.Remove(path)
			if err != nil {
				log.Printf("WARN unable to delete %v with error %v\n", path, err)
			} else if verbose {
				log.Printf("INFO file %v deleted\n", path)
			}
		}
	}
	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("unable to access path %q: %v", path, err)
		}
		if info.Mode().IsRegular() {
			ext := filepath.Ext(path)
			switch ext {
			case ".zip":
				err = Unzip(path, verbose)
				if err != nil {
					log.Printf("ERROR unable to execute unzip on file '%v' with error '%v'", path, err)
				} else {
					cleanup(path)
				}
			case ".tar":
				err = UnTar(path, verbose)
				if err != nil {
					log.Printf("ERROR unable to execute untar on file '%v' with error '%v'", path, err)
				} else {
					cleanup(path)
				}
			case ".gz":
				newFile, err := UnGzip(path, verbose)
				if err != nil {
					log.Printf("ERROR unable to execute ungzip on file '%v' with error '%v'", path, err)
				} else {
					cleanup(path)
				}
				newExt := filepath.Ext(newFile)
				if newExt == ".tar" {
					err = UnTar(newFile, verbose)
					if err != nil {
						log.Printf("ERROR unable to execute untar on file '%v' with error '%v'", path, err)
					} else {
						cleanup(newFile)
					}
				}
			case ".tgz":
				newFile, err := UnGzip(path, verbose)
				if err != nil {
					log.Printf("ERROR unable to execute ungzip on file '%v' with error '%v'", path, err)
				} else {
					cleanup(path)
				}
				err = UnTar(newFile, verbose)
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
}
