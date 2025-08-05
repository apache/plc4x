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

package common

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

func GetLicenseFileContent(licenseFileName string, verbose bool) []byte {
	isAbs := path.IsAbs(licenseFileName)
	var licenceContent []byte
	rootReached := false
	currentDir, _ := os.Getwd()
	var licenseFileNameWithPath string
	for !isAbs && !rootReached {
		if verbose {
			fmt.Printf("Looking for %s in %s\n", licenseFileName, currentDir)
		}
		licenseFileNameWithPath = filepath.Join(currentDir, licenseFileName)
		if _, err := os.Stat(licenseFileNameWithPath); errors.Is(err, os.ErrNotExist) {
			if currentDir == "/" {
				return nil
			}
			// path/to/whatever does not exist
			currentDir = filepath.Join(currentDir, "..")
			continue
		}
		if verbose {
			fmt.Printf("Found %s in %s\n", licenseFileName, currentDir)
		}
		rootReached = true
		var err error
		licenceContent, err = os.ReadFile(licenseFileNameWithPath)
		if err != nil {
			log.Fatalf("error reading license file: %v", err)
		}
	}
	return licenceContent
}
