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
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// RunParserTruncationTestsuite parses every prefix of every message in a parser-serializer
// testsuite and requires each one to fail, if it fails, as an error rather than as a panic.
//
// The remote end decides where a message ends: a datagram can arrive short, and a stream can be
// closed mid-message, so every prefix of a well formed message is something a driver can be handed.
// Truncation is worth exercising in particular because a field read that runs out of bytes is
// reported as an absent optional rather than as a failure, so anything deriving a value from that
// field sees it unset in a case the surrounding condition says cannot happen.
//
// The reference messages of a testsuite make a far better corpus for this than hand-written input,
// because they reach every type the protocol's own tests reach. Whether a prefix parses is left
// open - a truncated message may legitimately parse into a smaller one. Only the manner of failure
// is asserted.
func RunParserTruncationTestsuite(t *testing.T, testPath string, parser Parser, options ...options.WithOption) {
	skippedTestCases := map[string]bool{}
	for _, withOption := range options {
		if option, ok := withOption.(withSkippedTestCases); ok {
			for _, skippedTestCase := range option.skippedTestCases {
				t.Logf("Skipping %s", skippedTestCase)
				skippedTestCases[skippedTestCase] = true
			}
		}
	}

	rootNode := ParseParserSerializerTestSuiteXml(t, testPath)
	testsuite := ParseParserSerializerTestSuite(t, *rootNode, parser, nil)

	messages, prefixes := 0, 0
	for _, testcase := range testsuite.testcases {
		if skippedTestCases[testcase.name] {
			continue
		}
		rawInput, err := hex.DecodeString(testcase.rawInputText)
		if err != nil {
			t.Errorf("testcase %s has raw input that is not hex: %v", testcase.name, err)
			continue
		}
		messages++

		t.Run(testcase.name, func(t *testing.T) {
			// Up to and including the whole message, so the untruncated case is covered too.
			for length := 0; length <= len(rawInput); length++ {
				prefixes++
				if recovered := testsuite.parseRecovering(testcase, rawInput[:length]); recovered != nil {
					t.Errorf("truncated to %d of %d bytes (%s) panicked: %v",
						length, len(rawInput), hex.EncodeToString(rawInput[:length]), recovered)
				}
			}
		})
	}
	t.Logf("parsed every prefix of %d messages of testsuite %s (%d prefixes)", messages, testsuite.name, prefixes)
}

// parseRecovering returns the value recovered from a panic while parsing, or nil if parsing did not
// panic. A returned error is a permitted outcome; only a panic is not.
func (p *ParserSerializerTestsuite) parseRecovering(testcase ParserSerializerTestcase, rawInput []byte) (recovered any) {
	defer func() { recovered = recover() }()

	var readBuffer utils.ReadBuffer
	if p.byteOrder == binary.LittleEndian {
		readBuffer = utils.NewReadBufferByteBased(rawInput, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
	} else {
		readBuffer = utils.NewReadBufferByteBased(rawInput)
	}
	_, _ = p.parser.Parse(testcase.rootType, testcase.parserArguments, readBuffer)
	return nil
}
