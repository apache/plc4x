//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package testutils

import (
	"errors"
	"github.com/ajankovic/xdiff"
	"github.com/ajankovic/xdiff/parser"
	"os"
)

func CompareResults(actualString []byte, referenceString []byte) error {
	// Now parse the xml strings of the actual and the reference in xdiff's dom
	p := parser.New()
	actual, err := p.ParseBytes(actualString)
	if err != nil {
		return errors.New("Error parsing actual input: " + err.Error())
	}
	reference, err := p.ParseBytes(referenceString)
	if err != nil {
		return errors.New("Error parsing reference input: " + err.Error())
	}
	// Use XDiff to actually do the comparison
	diff, err := xdiff.Compare(actual, reference)
	if err != nil {
		return errors.New("Error comparing xml trees: " + err.Error())
	}
	if diff != nil {
		enc := xdiff.NewTextEncoder(os.Stdout)
		if err := enc.Encode(diff); err != nil {
			return errors.New("Error outputting results: " + err.Error())
		}
		return errors.New("there were differences: Expected: \n" + string(referenceString) + "\nBut Got: \n" + string(actualString))
	}
	return nil
}
