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
	"os"
	"path"
	"runtime"
	"strings"
)

type DebugPrinter = func(format string, a ...any)

const projectName = "bacgopes"

var customProjectName = os.Getenv("BACGOPES_DEBUG_CUSTOM_PROJECT_NAME")
var customReplaces map[string]string

func init() {
	customReplaces = map[string]string{}
	for _, replace := range strings.Split(os.Getenv("BACGOPES_DEBUG_CUSTOM_REPLACES"), ",") {
		if replace == "" || !strings.Contains(replace, "=") {
			continue
		}
		kv := strings.SplitN(replace, "=", 2)
		if len(kv) != 2 {
			panic("invalid replace " + replace)
		}
		customReplaces[kv[0]] = kv[1]
	}
}

func CreateDebugPrinter() DebugPrinter {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return nil
	}
	dir := path.Dir(file)
	rootIndex := strings.Index(dir, projectName)
	dir = dir[rootIndex:]
	qualifier := strings.ReplaceAll(dir, "/", ".")
	switch {
	case strings.HasPrefix(qualifier, projectName+".tests"):
		qualifier = "tests" + strings.TrimPrefix(qualifier, projectName+".tests")
	}
	dirPrefix := path.Base(dir) + "_"

	bacgopesDebug := os.Getenv("BACGOPES_DEBUG")
	if strings.Contains(bacgopesDebug, qualifier) {
		_isDebuggingActive = true
		return func(format string, a ...any) {
			_, file, _, ok := runtime.Caller(1)
			if !ok {
				return
			}
			base := path.Base(file)
			prefix := strings.TrimSuffix(base, ".go")
			prefix = strings.TrimPrefix(prefix, dirPrefix)
			formatString := "DEBUG:" + qualifier + "." + prefix + ":" + format + "\n"
			formatString = cleanupFormatString(formatString)
			if customProjectName != "" {
				formatString = strings.ReplaceAll(formatString, projectName, customProjectName)
			}
			formatString = fixVerbs(formatString, a...)
			output := fmt.Sprintf(formatString, a...)
			if strings.HasSuffix(output, "\n\n") { // TODO: another hacky workaround
				output = strings.TrimSuffix(output, "\n")
			}
			for k, v := range customReplaces {
				output = strings.ReplaceAll(output, k, v)
			}
			if customProjectName != "" {
				output = strings.ReplaceAll(output, projectName, customProjectName)
			}
			_, _ = os.Stdout.Write([]byte(output))
		}
	}
	return nil
}

func fixVerbs(formatString string, values ...any) string {
	length := len(formatString)
	verbNumber := -1
	for i, r := range formatString {
		switch r {
		case '%':
			nextIndex := i + 1
			if nextIndex >= length {
				continue
			}
			followRune := formatString[nextIndex]
			if followRune != '%' {
				verbNumber++
			}
			if followRune == 'r' && verbNumber < len(values) { // TODO: this completely breaks at indexed verbs... better fix assap
				runes := []rune(formatString)
				runes[nextIndex] = VerbForType(values[verbNumber], 'r')
				formatString = string(runes)
			}
		}
	}
	return formatString
}

func cleanupFormatString(s string) string {
	// TODO: investigate via comm.comm is happening
	s = strings.ReplaceAll(s, ".comm.comm:", ".comm:")
	s = strings.ReplaceAll(s, "pdu.comm_PCI:", "comm.PCI:")
	s = strings.ReplaceAll(s, "pdu.comm_PDUData:", "comm.PDUData:")
	s = strings.ReplaceAll(s, "DEBUG:"+projectName+".tests", "DEBUG:tests")
	s = strings.ReplaceAll(s, "appservice.app_DeviceInfoCache", "app.DeviceInfoCache")
	s = strings.ReplaceAll(s, "device_WhoIsIAmServices", "device.WhoIsIAmServices")
	return s
}
