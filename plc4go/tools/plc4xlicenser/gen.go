/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	typeNames   = flag.String("type", "", "comma-separated list of type names; must be set")
	output      = flag.String("output", "", "output file name; default srcdir/<type>_string.go")
	buildTags   = flag.String("tags", "", "comma-separated list of build tags to apply")
	licenseFile = flag.String("licenseFile", ".mockery.asl.header", "file containing the license (will be searched upwards)")
	verbose     = flag.Bool("verbose", false, "verbosity")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of plc4xlicenser:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tplc4xlicenser [flags] -type T [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tplc4xlicenser [flags] -type T files... # Must be a single package\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("stringer: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	types := strings.Split(*typeNames, ",")
	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	var dir string
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		dir = filepath.Dir(args[0])
	}

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_string.go", types[0])
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	inputFile, err2 := os.ReadFile(outputName)

	licenseFileName := *licenseFile
	isAbs := path.IsAbs(licenseFileName)
	var licenceContent []byte
	rootReached := false
	currentDir, _ := os.Getwd()
	var licenseFileNameWithPath string
	for !isAbs && !rootReached {
		if *verbose {
			fmt.Printf("Looking for %s in %s\n", licenseFileName, currentDir)
		}
		licenseFileNameWithPath = filepath.Join(currentDir, licenseFileName)
		if _, err := os.Stat(licenseFileNameWithPath); errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does not exist
			currentDir = filepath.Join(currentDir, "..")
			continue
		}
		if *verbose {
			fmt.Printf("Found %s in %s\n", licenseFileName, currentDir)
		}
		rootReached = true
		var err error
		licenceContent, err = os.ReadFile(licenseFileNameWithPath)
		if err != nil {
			panic(err2)
		}
	}

	if err := os.WriteFile(outputName, append(licenceContent, inputFile...), 0644); err != nil {
		log.Fatalf("writing output: %s", err)
	}
	fmt.Printf("Fixed plc4x license of %s\n", licenseFileNameWithPath)
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
