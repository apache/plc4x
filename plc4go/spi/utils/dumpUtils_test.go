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
	"testing"

	"github.com/stretchr/testify/assert"
)

type _TestBoxAnythingTestBoxer struct {
}

func (_TestBoxAnythingTestBoxer) Box(options ...func(*BoxOptions)) AsciiBox {
	return AsciiBox{data: "i did it"}
}

type _TestBoxAnythingTestBoxerStringer struct {
}

func (_TestBoxAnythingTestBoxerStringer) String() string {
	return "i did it"
}

type _TestBoxAnythingTestBoxerAnything struct {
	someField string
	someInt   int
}

func TestBoxAnything(t *testing.T) {
	type args struct {
		anything any
		options  []func(*BoxOptions)
	}
	tests := []struct {
		name string
		args args
		want AsciiBox
	}{
		{
			name: "test nil",
			want: AsciiBox{asciiBoxWriter: AsciiBoxWriterDefault.(*asciiBoxWriter)},
		},
		{
			name: "asciiBoxer",
			args: args{
				anything: _TestBoxAnythingTestBoxer{},
			},
			want: AsciiBox{data: "i did it"},
		},
		{
			name: "asciiBoxer as interface",
			args: args{
				anything: func() AsciiBoxer {
					var a AsciiBoxer
					return a
				}(),
			},
			want: AsciiBox{asciiBoxWriter: AsciiBoxWriterDefault.(*asciiBoxWriter)},
		},
		{
			name: "test bool",
			args: args{
				anything: true,
				options: []func(*BoxOptions){
					WithAsciiBoxName("exampleBool"),
				},
			},
			want: asciiBoxForTest(`
╔═exampleBool╗
║  b1 true   ║
╚════════════╝`[1:]),
		},
		{
			name: "test int",
			args: args{
				anything: 1,
				options: []func(*BoxOptions){
					WithAsciiBoxName("exampleInt"),
				},
			},
			want: asciiBoxForTest(`
╔═exampleInt═════════╗
║0x0000000000000001 1║
╚════════════════════╝`[1:]),
		},
		{
			name: "test int 123123123",
			args: args{
				anything: 123123123,
				options: []func(*BoxOptions){
					WithAsciiBoxName("exampleInt"),
				},
			},
			want: asciiBoxForTest(`
╔═exampleInt═════════════════╗
║0x000000000756b5b3 123123123║
╚════════════════════════════╝`[1:]),
		},
		{
			name: "test []byte",
			args: args{
				anything: []byte{1, 2, 3},
				options: []func(*BoxOptions){
					WithAsciiBoxName("example byte[]"),
				},
			},
			want: asciiBoxForTest(`
╔═example byte[]╗
║   0|01 '.'    ║
║   1|02 '.'    ║
║   2|03 '.'    ║
╚═══════════════╝`[1:]),
		},
		{
			name: "box string",
			args: args{
				anything: "asdasdasdasdasd",
				options: []func(*BoxOptions){
					WithAsciiBoxName("example string"),
				},
			},
			want: asciiBoxForTest(`
╔═example string╗
║asdasdasdasdasd║
╚═══════════════╝`[1:]),
		},
		{
			name: "box stringer",
			args: args{
				anything: _TestBoxAnythingTestBoxerStringer{},
				options: []func(*BoxOptions){
					WithAsciiBoxName("example stringer"),
				},
			},
			want: asciiBoxForTest(`
╔═example stringer╗
║    i did it     ║
╚═════════════════╝`[1:]),
		},
		{
			name: "anything struct",
			args: args{
				anything: _TestBoxAnythingTestBoxerAnything{
					someField: "a field value",
					someInt:   234243,
				},
				options: []func(*BoxOptions){
					WithAsciiBoxName("any struct"),
				},
			},
			want: asciiBoxForTest(`
╔═any struct═══════════╗
║{a field value 234243}║
╚══════════════════════╝`[1:]),
		},
		{
			name: "anything struct pointer",
			args: args{
				anything: &_TestBoxAnythingTestBoxerAnything{
					someField: "a field value",
					someInt:   234243,
				},
				options: []func(*BoxOptions){
					WithAsciiBoxName("any struct"),
				},
			},
			want: asciiBoxForTest(`
╔═any struct═══════════╗
║{a field value 234243}║
╚══════════════════════╝`[1:]),
		},
		{
			name: "int pointer",
			args: args{
				anything: func() *int {
					a := 13
					return &a
				}(),
				options: []func(*BoxOptions){
					WithAsciiBoxName("int pointer"),
				},
			},
			want: asciiBoxForTest(`
╔═int pointer═════════╗
║0x000000000000000d 13║
╚═════════════════════╝`[1:]),
		},
		{
			name: "array",
			args: args{
				anything: []any{1, "1,2,3", 2, _TestBoxAnythingTestBoxerAnything{
					someField: "a field value",
					someInt:   234243,
				}, &_TestBoxAnythingTestBoxerAnything{
					someField: "a field value",
					someInt:   234243,
				}},
				options: []func(*BoxOptions){
					WithAsciiBoxName("array"),
				},
			},
			want: asciiBoxForTest(`
╔═array══════════════════╗
║╔════════════════════╗  ║
║║0x0000000000000001 1║  ║
║╚════════════════════╝  ║
║╔═════╗                 ║
║║1,2,3║                 ║
║╚═════╝                 ║
║╔════════════════════╗  ║
║║0x0000000000000002 2║  ║
║╚════════════════════╝  ║
║╔══════════════════════╗║
║║{a field value 234243}║║
║╚══════════════════════╝║
║╔══════════════════════╗║
║║{a field value 234243}║║
║╚══════════════════════╝║
╚════════════════════════╝`[1:]),
		},
		{
			name: "slice",
			args: args{
				anything: []any{1, "1,2,3", 2, _TestBoxAnythingTestBoxerAnything{
					someField: "a field value",
					someInt:   234243,
				}, &_TestBoxAnythingTestBoxerAnything{
					someField: "a field value",
					someInt:   234243,
				}}[:],
				options: []func(*BoxOptions){
					WithAsciiBoxName("array"),
				},
			},
			want: asciiBoxForTest(`
╔═array══════════════════╗
║╔════════════════════╗  ║
║║0x0000000000000001 1║  ║
║╚════════════════════╝  ║
║╔═════╗                 ║
║║1,2,3║                 ║
║╚═════╝                 ║
║╔════════════════════╗  ║
║║0x0000000000000002 2║  ║
║╚════════════════════╝  ║
║╔══════════════════════╗║
║║{a field value 234243}║║
║╚══════════════════════╝║
║╔══════════════════════╗║
║║{a field value 234243}║║
║╚══════════════════════╝║
╚════════════════════════╝`[1:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoxAnything(tt.args.anything, tt.args.options...); !assert.Equal(t, tt.want, got) {
				t.Errorf("BoxAnything() = '\n%v'\n, want '\n%v'", got, tt.want)
			}
		})
	}
}

func TestBoxedDump(t *testing.T) {
	type args struct {
		data    []byte
		options []func(*BoxOptions)
	}
	tests := []struct {
		name string
		args args
		want AsciiBox
	}{
		{
			name: "Test Dump 1",
			args: args{
				data: []byte("1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aa1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aab"),
				options: []func(*BoxOptions){
					WithAsciiBoxName("super nice data"),
				},
			},
			want: asciiBoxForTest(`
╔═super nice data══════════════════════════════╗
║000|31 32 33 34 35 36 37 38 39 30 '1234567890'║
║010|61 62 63 64 65 66 67 68 69 6a 'abcdefghij'║
║020|6b 6c 6d 6e 6f 70 71 72 73 74 'klmnopqrst'║
║030|75 76 77 78 79 7a d3 31 32 33 'uvwxyz.123'║
║040|34 35 36 37 38 39 30 61 62 63 '4567890abc'║
║050|64 65 66 67 68 69 6a 6b 6c 6d 'defghijklm'║
║060|6e 6f 70 71 72 73 74 75 76 77 'nopqrstuvw'║
║070|78 79 7a d3 61 61 31 32 33 34 'xyz.aa1234'║
║080|35 36 37 38 39 30 61 62 63 64 '567890abcd'║
║090|65 66 67 68 69 6a 6b 6c 6d 6e 'efghijklmn'║
║100|6f 70 71 72 73 74 75 76 77 78 'opqrstuvwx'║
║110|79 7a d3 31 32 33 34 35 36 37 'yz.1234567'║
║120|38 39 30 61 62 63 64 65 66 67 '890abcdefg'║
║130|68 69 6a 6b 6c 6d 6e 6f 70 71 'hijklmnopq'║
║140|72 73 74 75 76 77 78 79 7a d3 'rstuvwxyz.'║
║150|61 61 62                      'aab       '║
╚══════════════════════════════════════════════╝`[1:]),
		},
		{
			name: "Test Dump 2",
			args: args{
				data: []byte("1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aa1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aab"),
				options: []func(*BoxOptions){
					WithAsciiBoxName("super nice data"),
					WithAsciiBoxCharWidth(110),
				},
			},
			want: asciiBoxForTest(`
╔═super nice data════════════════════════════════════════════════════════════════════════════════════════════╗
║ 000|31 32 33 34 35 36 37 38 39 30 61 62 63 64 65 66 67 68 69 6a 6b 6c 6d 6e 6f '1234567890abcdefghijklmno' ║
║ 025|70 71 72 73 74 75 76 77 78 79 7a d3 31 32 33 34 35 36 37 38 39 30 61 62 63 'pqrstuvwxyz.1234567890abc' ║
║ 050|64 65 66 67 68 69 6a 6b 6c 6d 6e 6f 70 71 72 73 74 75 76 77 78 79 7a d3 61 'defghijklmnopqrstuvwxyz.a' ║
║ 075|61 31 32 33 34 35 36 37 38 39 30 61 62 63 64 65 66 67 68 69 6a 6b 6c 6d 6e 'a1234567890abcdefghijklmn' ║
║ 100|6f 70 71 72 73 74 75 76 77 78 79 7a d3 31 32 33 34 35 36 37 38 39 30 61 62 'opqrstuvwxyz.1234567890ab' ║
║ 125|63 64 65 66 67 68 69 6a 6b 6c 6d 6e 6f 70 71 72 73 74 75 76 77 78 79 7a d3 'cdefghijklmnopqrstuvwxyz.' ║
║ 150|61 61 62                                                                   'aab                      ' ║
╚════════════════════════════════════════════════════════════════════════════════════════════════════════════╝`[1:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoxedDump(tt.args.data, tt.args.options...); !assert.Equal(t, tt.want, got) {
				t.Errorf("Dump() = \n%v\n, want \n%v\n", got, tt.want)
			}
		})
	}
}

func TestDumpAnything(t *testing.T) {
	t.Skip("That seems to be platform dependent, so not sure how to fix that")
	type args struct {
		anything any
		options  []func(*BoxOptions)
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Random struct",
			args: args{
				anything: struct {
					A string
					B string
					C string
					D struct {
						E string
						F string
					}
				}{A: "a", B: "b", C: "c", D: struct {
					E string
					F string
				}{
					E: "e",
					F: "f",
				}},
			},
			want: `
000|25 ff 81 03 01 02 ff 82 00 01 '%.........'
010|04 01 01 41 01 0c 00 01 01 42 '...A.....B'
020|01 0c 00 01 01 43 01 0c 00 01 '.....C....'
030|01 44 01 ff 84 00 00 00 37 ff '.D......7.'
040|83 03 01 01 1d 73 74 72 75 63 '.....struc'
050|74 20 7b 20 45 20 73 74 72 69 't { E stri'
060|6e 67 3b 20 46 20 73 74 72 69 'ng; F stri'
070|6e 67 20 7d 01 ff 84 00 01 02 'ng }......'
080|01 01 45 01 0c 00 01 01 46 01 '..E.....F.'
090|0c 00 00 00 14 ff 82 01 01 61 '.........a'
100|01 01 62 01 01 63 01 01 01 65 '..b..c...e'
110|01 01 66 00 00                '..f..     '`[1:],
		},
		{
			name: "unexported struct gob error",
			args: args{
				anything: struct {
					a string
				}{a: "a"},
			},
			want: "<undumpable>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DumpAnything(tt.args.anything, tt.args.options...); !assert.Equal(t, tt.want, got) {
				t.Errorf("Dump() = \n%v\n, want \n%v\n", got, tt.want)
			}
		})
	}
}

func TestBoxedDumpAnything(t *testing.T) {
	type args struct {
		anything any
		options  []func(*BoxOptions)
	}
	tests := []struct {
		name string
		args args
		want AsciiBox
	}{
		{
			name: "dump it",
			want: asciiBoxForTest(`
╔════════════╗
║<undumpable>║
╚════════════╝`[1:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoxedDumpAnything(tt.args.anything, tt.args.options...); !assert.Equal(t, tt.want, got) {
				t.Errorf("BoxedDumpAnything() = \n%v\n, want \n%v", got, tt.want)
			}
		})
	}
}
