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

package utils

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGenerateId(t *testing.T) {
	type args struct {
		numBytes int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "generate it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldRandomByteFiller := randomByteFiller
			t.Cleanup(func() {
				randomByteFiller = oldRandomByteFiller
			})
			randomByteFiller = func(_ []byte) (n int, err error) {
				return 0, nil
			}
			assert.Equalf(t, tt.want, GenerateId(produceTestLogger(t), tt.args.numBytes), "GenerateId(%v)", tt.args.numBytes)
		})
	}
}

// note: we can't use testutils here due to import cycle
func produceTestLogger(t *testing.T) zerolog.Logger {
	return zerolog.New(zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(t),
		func(w *zerolog.ConsoleWriter) {
			// TODO: this is really an issue with go-junit-report not sanitizing output before dumping into xml...
			onJenkins := os.Getenv("JENKINS_URL") != ""
			onGithubAction := os.Getenv("GITHUB_ACTIONS") != ""
			onCI := os.Getenv("CI") != ""
			if onJenkins || onGithubAction || onCI {
				w.NoColor = true
			}
		}))
}
