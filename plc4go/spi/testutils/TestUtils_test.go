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
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func TestExplodingGlobalLogger(t *testing.T) {
	oldLog := log.Logger
	t.Cleanup(func() {
		log.Logger = oldLog
	})
	assert.Panics(t, func() {
		ExplodingGlobalLogger(true)
		log.Logger.Info().Msg("this should explode")
	})
}

type ASerializable struct {
	a string
	b int
	c float32
	D string
}

func (A *ASerializable) Serialize() ([]byte, error) {
	panic("not needed")
}

func (A *ASerializable) SerializeWithWriteBuffer(_ context.Context, writeBuffer utils.WriteBuffer) error {
	_ = writeBuffer.WriteString("a", 8, A.a)
	_ = writeBuffer.WriteInt64("b", 64, int64(A.b))
	_ = writeBuffer.WriteFloat32("c", 32, A.c)
	return nil
}

func (A *ASerializable) String() string {
	wbbb := utils.NewWriteBufferBoxBased()
	_ = A.SerializeWithWriteBuffer(nil, wbbb)
	return wbbb.GetBox().String()
}

func TestProduceTestingLogger_BetterStackRendering(t *testing.T) {
	got := ProduceTestingLogger(t)
	f1 := func() error {
		return errors.New("a error")
	}
	f2 := func() error {
		return errors.Wrap(f1(), "b error")
	}
	f3 := func() error {
		return errors.Wrap(f2(), "c error")
	}
	got.Error().Err(f3()).Msg("multiline error")
	got.Error().Err(nil).Msg("empty error")
}

func TestProduceTestingLogger_ASerializableLog(t *testing.T) {
	got := ProduceTestingLogger(t)
	aSerializable := &ASerializable{
		a: "a",
		b: 2,
		c: 3.1413,
	}
	got.Info().
		Interface("aSerializableInterface", aSerializable).
		Stringer("aSerializableStringer", aSerializable).
		Str("aString", "asdasdasd").
		Msg("something")
}

type TestProduceTestingLoggerImprovedOutputMultilineStringer struct {
}

func (TestProduceTestingLoggerImprovedOutputMultilineStringer) String() string {
	return "\nthis\nthis\nthis\nthis\nthis\nthis\n"
}

func TestProduceTestingLogger_Improved_Output(t *testing.T) {
	t.Run("multilinestring", func(t *testing.T) {
		logger := ProduceTestingLogger(t)
		logger.Info().Str("amultiline", "a\nb\nc").Msg("look at that")
	})
	t.Run("multilinestringer", func(t *testing.T) {
		logger := ProduceTestingLogger(t)
		logger.Info().Stringer("amultiline", TestProduceTestingLoggerImprovedOutputMultilineStringer{}).Msg("look at that")
	})
	t.Run("multilinestringers", func(t *testing.T) {
		logger := ProduceTestingLogger(t)
		logger.Info().Stringers("amultiline", []fmt.Stringer{
			TestProduceTestingLoggerImprovedOutputMultilineStringer{},
			TestProduceTestingLoggerImprovedOutputMultilineStringer{},
		}).Msg("look at that")
	})
}

func Test__explodingGlobalLogger_Write(t *testing.T) {
	type fields struct {
		hardExplode bool
	}
	type args struct {
		in0 []byte
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      int
		wantErr   bool
		wantPanic bool
	}{
		{
			name:    "write just errors",
			wantErr: true,
		},
		{
			name: "write panics",
			fields: fields{
				hardExplode: true,
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := _explodingGlobalLogger{
				hardExplode: tt.fields.hardExplode,
			}
			var got int
			var err error
			if tt.wantPanic {
				assert.Panics(t, func() {
					_, _ = e.Write(tt.args.in0)
				})
				return
			} else {
				got, err = e.Write(tt.args.in0)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOrLeaveBool(t *testing.T) {
	type args struct {
		key     string
		setting *bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getOrLeaveBool(tt.args.key, tt.args.setting)
		})
	}
}

// TestCompareResults_StringLeafDiffs pins down that dataType="string" leaf
// differences are no longer swallowed by CompareResults. Previously a
// "We ignore newline diffs" carve-out (added in 8828c6815e, meant to survive
// a CDATA-vs-entity-escaping rendering difference that no longer exists)
// unconditionally skipped every dataType="string" Update delta because its
// guard re-invoked xdiff.Compare directly on the two leaf value nodes, which
// always reports "no diff" for a leaf-vs-leaf root comparison regardless of
// their actual content. See GH-2611 EIP parity follow-up.
func TestCompareResults_StringLeafDiffs(t *testing.T) {
	t.Run("equal documents still pass", func(t *testing.T) {
		doc := `<root><code dataType="string" bitLength="16" encoding="UTF-8">OK</code></root>`
		fakeT := &testing.T{}
		err := CompareResults(fakeT, []byte(doc), []byte(doc))
		assert.NoError(t, err)
		assert.False(t, fakeT.Failed())
	})

	t.Run("differing string leaf now fails", func(t *testing.T) {
		expected := `<root><code dataType="string" bitLength="16" encoding="UTF-8">OK</code></root>`
		actual := `<root><code dataType="string" bitLength="16" encoding="UTF-8">NOTOK</code></root>`
		fakeT := &testing.T{}
		err := CompareResults(fakeT, []byte(actual), []byte(expected))
		assert.Error(t, err, "a string-leaf value difference must fail CompareResults")
	})

	t.Run("differing non-string leaf still fails (control)", func(t *testing.T) {
		expected := `<root><elementNb dataType="uint" bitLength="16">1</elementNb></root>`
		actual := `<root><elementNb dataType="uint" bitLength="16">2</elementNb></root>`
		fakeT := &testing.T{}
		err := CompareResults(fakeT, []byte(actual), []byte(expected))
		assert.Error(t, err)
	})
}

func Test_getOrLeaveDuration(t *testing.T) {
	type args struct {
		key     string
		setting *time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getOrLeaveDuration(tt.args.key, tt.args.setting)
		})
	}
}
