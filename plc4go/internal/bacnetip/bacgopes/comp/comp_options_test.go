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

package comp

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

func TestAddIfAbundant(t *testing.T) {
	type args[T any] struct {
		options []Option
		value   T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []Option
	}
	tests := []testCase[spi.Message]{
		{
			name: "nothing in nil out",
			want: []Option{genericOption[spi.Message]{value: nil}},
		},
		{
			name: "in nothing out",
			args: args[spi.Message]{
				options: []Option{
					genericOption[spi.Message]{value: nil},
				},
			},
			want: []Option{genericOption[spi.Message]{value: nil}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddIfAbundant(tt.args.options, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddIfAbundant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddIfAbundantMixed(t *testing.T) {
	var initialOptions []Option
	AddIfAbundant(initialOptions, genericOption[spi.Message]{})
	AddIfAbundant(initialOptions, genericOption[string]{})
	AddIfAbundant(initialOptions, genericOption[int]{})
	AddIfAbundant(initialOptions, genericOption[uint]{})
}

func TestAddIfAbundantCheckCollisions(t *testing.T) {
	t.Run("bigger to smaller", func(t *testing.T) {
		someMessage := genericOption[spi.Message]{}
		aNLM := readWriteModel.NLM(nil)
		initialOptions := []Option{someMessage}
		enrichedOptions := AddIfAbundant(initialOptions, aNLM)
		assert.Len(t, enrichedOptions, 2)
		assert.Contains(t, enrichedOptions, someMessage)
		assert.Contains(t, enrichedOptions, genericOption[readWriteModel.NLM]{value: aNLM})
		assert.Equal(t, enrichedOptions, []Option{someMessage, genericOption[readWriteModel.NLM]{value: aNLM}})
	})
	t.Run("smaller to bigger", func(t *testing.T) {
		someMessage := spi.Message(nil)
		aNLM := genericOption[readWriteModel.NLM]{}
		initialOptions := []Option{aNLM}
		enrichedOptions := AddIfAbundant(initialOptions, someMessage)
		assert.Len(t, enrichedOptions, 2)
		assert.Contains(t, enrichedOptions, genericOption[spi.Message]{value: someMessage})
		assert.Contains(t, enrichedOptions, aNLM)
		assert.Equal(t, enrichedOptions, []Option{aNLM, genericOption[spi.Message]{value: someMessage}})
	})
	t.Run("does not collide", func(t *testing.T) {
		someMessage := genericOption[spi.Message]{}
		aNLM := genericOption[readWriteModel.NLM]{}
		aAPDU := genericOption[readWriteModel.APDU]{}
		initialOptions := []Option{someMessage}
		enrichedOptions := AddIfAbundant(initialOptions, aNLM)
		enrichedOptions = AddIfAbundant(enrichedOptions, aAPDU)
		assert.Len(t, enrichedOptions, 3)
	})
	t.Run("empty ok", func(t *testing.T) {
		aNLM := genericOption[readWriteModel.NLM]{}
		aAPDU := genericOption[readWriteModel.APDU]{}
		var initialOptions []Option
		enrichedOptions := AddIfAbundant(initialOptions, aNLM)
		enrichedOptions = AddIfAbundant(enrichedOptions, aAPDU)
		assert.Len(t, enrichedOptions, 2)
	})
}

func TestApplyGenericOption(t *testing.T) {
	type args[T any] struct {
		options []Option
		applier func(T)
	}
	type testCase[T any] struct {
		name string
		args args[T]
	}
	tests := []testCase[string]{
		{
			name: "just apply",
			args: args[string]{
				options: []Option{genericOption[string]{}},
				applier: func(s string) {

				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyGenericOption(tt.args.options, tt.args.applier)
		})
	}
}

func TestExtractIfPresent(t *testing.T) {
	type A interface{}
	type B interface{ A }
	type C interface{ A }
	type D struct{}
	type E interface{}
	someA := genericOption[A]{value: "string"}
	someB := genericOption[B]{value: 2}
	someC := genericOption[C]{value: true}
	someD := genericOption[*D]{value: new(D)}
	options := []Option{someA, someB, someC, someD}
	assert.Len(t, options, 4)
	assert.Equal(t, options, []Option{someA, someB, someC, someD})
	extractedA, ok := ExtractIfPresent[A](options)
	assert.True(t, ok)
	assert.Equal(t, someA.value, extractedA)
	extractedB, ok := ExtractIfPresent[B](options)
	assert.True(t, ok)
	assert.Equal(t, someB.value, extractedB)
	extractedC, ok := ExtractIfPresent[C](options)
	assert.True(t, ok)
	assert.Equal(t, someC.value, extractedC)
	extractedD, ok := ExtractIfPresent[*D](options)
	assert.True(t, ok)
	assert.Same(t, someD.value, extractedD)
	extractedE, ok := ExtractIfPresent[*E](options)
	assert.False(t, ok)
	assert.Nil(t, extractedE)
}
