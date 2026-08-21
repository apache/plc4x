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

package model

import (
	"context"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

const parserSerializerTestsuite = "../../../../assets/testing/protocols/bacnet/ParserSerializerTestsuite.xml"

// A datagram arrives over UDP, so the peer decides where it ends. Every prefix of a well formed
// message is therefore something a driver can be handed, and none of them may take a parser down:
// an optional field whose read runs out of bytes is reported as absent, and anything deriving a
// value from it has to cope with that. Truncating the corpus of real messages at every length
// exercises exactly that, across every type the corpus reaches.
func TestBVLCTruncatedAtEveryLengthDoesNotPanic(t *testing.T) {
	messages := loadTestsuiteMessages(t)
	if len(messages) == 0 {
		t.Fatal("no messages loaded from the parser-serializer testsuite")
	}
	t.Logf("truncating %d messages from the parser-serializer testsuite", len(messages))

	parsed := 0
	for _, message := range messages {
		// The corpus is only useful as a truncation base if it parses whole.
		if _, err := parseBVLC(message.data); err != nil {
			t.Errorf("%s: the untruncated message did not parse: %v", message.name, err)
			continue
		}
		parsed++

		for length := 0; length < len(message.data); length++ {
			truncated := message.data[:length]
			if r := parseBVLCRecovering(truncated); r != nil {
				t.Errorf("%s truncated to %d of %d bytes (%s) panicked: %v",
					message.name, length, len(message.data), hex.EncodeToString(truncated), r)
			}
		}
	}
	t.Logf("%d messages parsed whole and were truncated at every length", parsed)
}

func parseBVLC(data []byte) (BVLC, error) {
	return BVLCParseWithBuffer[BVLC](context.Background(), utils.NewReadBufferByteBased(data))
}

// parseBVLCRecovering returns the recovered panic value, or nil if parsing did not panic. An error
// is a permitted outcome - the point is only that the failure stays a value we can handle.
func parseBVLCRecovering(data []byte) (recovered any) {
	defer func() { recovered = recover() }()
	_, _ = parseBVLC(data)
	return nil
}

type testsuiteMessage struct {
	name string
	data []byte
}

// loadTestsuiteMessages reads the raw hex of every BVLC rooted case in the parser-serializer
// testsuite. Going through encoding/xml rather than matching text keeps the commented out cases -
// which hold raw bytes the suite deliberately does not use - out of the corpus.
func loadTestsuiteMessages(t *testing.T) []testsuiteMessage {
	t.Helper()

	content, err := os.ReadFile(parserSerializerTestsuite)
	if err != nil {
		t.Fatalf("could not read the parser-serializer testsuite: %v", err)
	}

	var suite struct {
		Testcases []struct {
			Name     string `xml:"name"`
			Raw      string `xml:"raw"`
			RootType string `xml:"root-type"`
		} `xml:"testcase"`
	}
	if err := xml.Unmarshal(content, &suite); err != nil {
		t.Fatalf("could not parse the parser-serializer testsuite: %v", err)
	}

	var messages []testsuiteMessage
	for i, testcase := range suite.Testcases {
		if strings.TrimSpace(testcase.RootType) != "BVLC" {
			continue
		}
		data, err := hex.DecodeString(strings.TrimSpace(testcase.Raw))
		if err != nil {
			t.Errorf("testcase %d (%s) has raw bytes that are not hex: %v", i, testcase.Name, err)
			continue
		}
		messages = append(messages, testsuiteMessage{
			name: fmt.Sprintf("case %d %q", i, strings.TrimSpace(testcase.Name)),
			data: data,
		})
	}
	return messages
}
