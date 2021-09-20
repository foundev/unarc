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
	"os"
	"path/filepath"
	"testing"
)

func TestWalkDir(t *testing.T) {
	dir := "testdata/walk"
	err := WalkDir(dir, false, false)
	if err != nil {
		t.Errorf("unable to walk %v with error %v", dir, err)
	}

	fileName := filepath.Join(dir, "f.txt")
	defer os.Remove(fileName)
	b, err := os.ReadFile(fileName)
	if err != nil {
		t.Errorf("did not extract file %v with error %v", fileName, err)
	}
	content := string(b)
	expected := "this is a file\n"
	if expected != content {
		t.Errorf("expected '%q' but was '%q' for file %v", expected, content, fileName)
	}

	fileName = filepath.Join(dir, "fgzip.txt")
	defer os.Remove(fileName)
	b, err = os.ReadFile(fileName)
	if err != nil {
		t.Errorf("did not extract file %v with error %v", fileName, err)
	}
	content = string(b)
	expected = "gunziped this file\n"
	if expected != content {
		t.Errorf("expected '%q' but was '%q' for file %v", expected, content, fileName)
	}

	fileName = filepath.Join(dir, "ftar.txt")
	defer os.Remove(fileName)
	b, err = os.ReadFile(fileName)
	if err != nil {
		t.Errorf("did not extract file %v with error %v", fileName, err)
	}
	content = string(b)
	expected = "this is a file for tar\n"
	if expected != content {
		t.Errorf("expected '%q' but was '%q' for file %v", expected, content, fileName)
	}

    //final cleanup of tar file from tgz
	fileName = filepath.Join(dir, "ftar")
	defer os.Remove(fileName)
}
