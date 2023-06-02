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

package testutils

import (
	"github.com/rs/zerolog/log"
	"os"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/ajankovic/xdiff"
	"github.com/ajankovic/xdiff/parser"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func CompareResults(t *testing.T, actualString []byte, referenceString []byte) error {
	localLog := ProduceTestingLogger(t)
	// Now parse the xml strings of the actual and the reference in xdiff's dom
	p := parser.New()
	actual, err := p.ParseBytes(actualString)
	if err != nil {
		return errors.Wrap(err, "Error parsing actual input")
	}
	reference, err := p.ParseBytes(referenceString)
	if err != nil {
		return errors.Wrap(err, "Error parsing reference input")
	}
	// Use XDiff to actually do the comparison
	diff, err := xdiff.Compare(actual, reference)
	if err != nil {
		return errors.Wrap(err, "Error comparing xml trees")
	}
	if diff == nil {
		// All good
		return nil
	}
	cleanDiff := make([]xdiff.Delta, 0)
	for _, delta := range diff {
		if delta.Operation == xdiff.Delete && delta.Subject.Value == nil || delta.Operation == xdiff.Insert && delta.Subject.Value == nil {
			localLog.Info().Msgf("We ignore empty elements which should be deleted %v", delta)
			continue
		}
		// Workaround for different precisions on float
		if delta.Operation == xdiff.Update &&
			string(delta.Subject.Parent.FirstChild.Name) == "dataType" &&
			string(delta.Subject.Parent.FirstChild.Value) == "float" &&
			string(delta.Object.Parent.FirstChild.Name) == "dataType" &&
			string(delta.Object.Parent.FirstChild.Value) == "float" {
			if strings.Contains(string(delta.Subject.Value), string(delta.Object.Value)) || strings.Contains(string(delta.Object.Value), string(delta.Subject.Value)) {
				localLog.Info().Msgf("We ignore precision diffs %v", delta)
				continue
			}
		}
		if delta.Operation == xdiff.Update &&
			string(delta.Subject.Parent.FirstChild.Name) == "dataType" &&
			string(delta.Subject.Parent.FirstChild.Value) == "string" &&
			string(delta.Object.Parent.FirstChild.Name) == "dataType" &&
			string(delta.Object.Parent.FirstChild.Value) == "string" {
			if diff, err := xdiff.Compare(delta.Subject, delta.Object); diff == nil && err == nil {
				localLog.Info().Msgf("We ignore newline diffs %v", delta)
				continue
			}
		}
		cleanDiff = append(cleanDiff, delta)
	}

	enc := xdiff.NewTextEncoder(os.Stdout)
	if err := enc.Encode(diff); err != nil {
		return errors.Wrap(err, "Error outputting results")
	}
	if len(cleanDiff) <= 0 {
		localLog.Warn().Msg("We only found non relevant changes")
		return nil
	}

	assert.Equal(t, string(referenceString), string(actualString))
	asciiBoxWriter := utils.NewAsciiBoxWriter()
	expectedBox := asciiBoxWriter.BoxString("expected", string(referenceString), 0)
	gotBox := asciiBoxWriter.BoxString("got", string(actualString), 0)
	boxSideBySide := asciiBoxWriter.BoxSideBySide(expectedBox, gotBox)
	_ = boxSideBySide // TODO: xml too distorted, we need a don't center option
	return errors.New("there were differences: Expected: \n" + string(referenceString) + "\nBut Got: \n" + string(actualString))
}

// ProduceTestingLogger produces a logger which redirects to testing.T
func ProduceTestingLogger(t *testing.T) zerolog.Logger {
	return zerolog.New(
		zerolog.NewConsoleWriter(
			zerolog.ConsoleTestWriter(t),
			func(w *zerolog.ConsoleWriter) {
				// TODO: this is really an issue with go-junit-report not sanitizing output before dumping into xml...
				onJenkins := os.Getenv("JENKINS_URL") != ""
				if onJenkins {
					w.NoColor = true
				}
			},
		),
	)
}

// SetToTestingLogger sets logger to  ProduceTestingLogger and resets it on cleanup
func SetToTestingLogger(t *testing.T, logger *zerolog.Logger) {
	oldLogger := *logger
	t.Cleanup(func() {
		*logger = oldLogger
	})
	newLogger := ProduceTestingLogger(t)
	*logger = newLogger
}

type _explodingGlobalLogger struct {
	hardExplode bool
}

func (e _explodingGlobalLogger) Write(_ []byte) (_ int, err error) {
	if e.hardExplode {
		debug.PrintStack()
		panic("found a global log usage")
	}
	return 0, errors.New("found a global log usage")
}

// ExplodingGlobalLogger Useful to find unredirected logs
func ExplodingGlobalLogger(hardExplode bool) {
	log.Logger = zerolog.New(_explodingGlobalLogger{hardExplode: hardExplode})
}
