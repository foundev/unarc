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
	"flag"
	"fmt"
	"os"

	"github.com/foundev/unarc/pkg"
)

func main() {
	help := flag.Bool("help", false, "show command line help")
	verbose := flag.Bool("v", false, "show verbose extraction")
	delarc := flag.Bool("del", false, "delete the files that are successfully extracted")
	flag.Parse()
	if flag.NArg() == 0 || *help {
		flag.Usage()
		//I throw an error here because I do not like to have scripts silently fail if I pass the wrong args
		os.Exit(1)
	}
	dir := flag.Args()[0]
	err := pkg.WalkDir(dir, *delarc, *verbose)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir, err)
		os.Exit(2)
	}
}
