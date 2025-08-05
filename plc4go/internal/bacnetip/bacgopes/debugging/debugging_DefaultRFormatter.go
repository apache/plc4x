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

package debugging

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

type DefaultRFormatter struct {
	header        string
	extraPrinters []DebugContentPrinter
}

func NewDefaultRFormatter(extraPrinters ...DebugContentPrinter) *DefaultRFormatter {
	pc, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("oh no")
	}
	dir := path.Dir(file)
	rootIndex := strings.Index(dir, projectName)
	dir = dir[rootIndex:]
	dirPrefix := path.Base(dir) + "_"
	base := path.Base(file)
	prefix := strings.TrimSuffix(base, ".go")
	prefix = strings.TrimPrefix(prefix, dirPrefix)
	qualifier := strings.ReplaceAll(dirPrefix, "_", ".")
	if strings.HasPrefix(qualifier, "test.") {
		qualifier = "tests.test_" + strings.TrimPrefix(qualifier, "test.")
	} else {
		qualifier = projectName + "." + qualifier
	}
	header := fmt.Sprintf("<%s object at 0x%x>", qualifier+prefix, pc)
	if customProjectName != "" {
		header = strings.ReplaceAll(header, projectName, customProjectName)
	}
	return &DefaultRFormatter{
		header:        header,
		extraPrinters: extraPrinters,
	}
}

func (d *DefaultRFormatter) StructHeader() []byte {
	return []byte(d.header)
}

func (d *DefaultRFormatter) Format(s fmt.State, v rune) {
	if d.header == "" && len(d.extraPrinters) == 0 {
		panic("misconfiguration")
	}
	switch v {
	case 'r':
		_, _ = s.Write(d.StructHeader())
		if len(d.extraPrinters) > 0 {
			_, _ = s.Write([]byte("\n"))
		}
		for _, printer := range d.extraPrinters {
			printer.PrintDebugContents(1, s, nil)
		}
	case 'v', 's':
		_, _ = s.Write(d.StructHeader())
	}
}

func (d *DefaultRFormatter) String() string {
	return d.header
}
