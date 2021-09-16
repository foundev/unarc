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
	"bytes"
	"os"
	"testing"
)

func TestUnzip(t *testing.T) {
	err := Unzip("testdata/withdir.zip", true)
	if err != nil {
		t.Errorf("unable to unzip with error %v", err)
	}
	defer os.RemoveAll("testdata/withdir")
	testWithDir(t)
}

func TestUnTar(t *testing.T) {
	err := UnTar("testdata/withdir-onlytar.tar", true)
	if err != nil {
		t.Errorf("unable to untar with error %v", err)
	}
	defer os.RemoveAll("testdata/withdir")
	testWithDir(t)
}

func TestGunzip(t *testing.T) {
	newFile, err := UnGzip("testdata/withdir.tar.gz", true)
	if err != nil {
		t.Errorf("unable to ungzip with error %v", err)
	}
	defer os.RemoveAll(newFile)
	f1, err := os.ReadFile(newFile)
	if err != nil {
		t.Errorf("unable to open new tar with error %v", err)
	}
	f2, err := os.ReadFile("testdata/withdir-onlytar.tar")
	if err != nil {
		t.Errorf("unable to open expected tar with error %v", err)
	}
	if !bytes.Equal(f1, f2) {
		t.Error("files to not match")
	}
}

func testWithDir(t *testing.T) {
	//test files
	f, err := os.ReadFile("testdata/withdir/ftar-nest.txt")
	if err != nil {
		t.Errorf("unable to read file ftar-nest.txt with error %v", err)
	}
	content := string(f)
	expected := "file for nested dir\n"
	if content != expected {
		t.Errorf("expected file ftar-nest.txt to have '%v' but had '%v'", expected, content)
	}
	f, err = os.ReadFile("testdata/withdir/nested/ftar-nest.txt")
	if err != nil {
		t.Errorf("unable to open file nested/ftar-nest.txt with error %v", err)
	}
	content = string(f)
	expected = "file for nested dir really nested\n"
	if content != expected {
		t.Errorf("expected file ftar-nest.txt to have '%v' but had '%v'", expected, content)
	}
}
