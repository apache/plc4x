/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"strings"
	"testing"
)

func TestBoxAnything(t *testing.T) {
	type args struct {
		name      string
		anything  interface{}
		charWidth int
	}
	tests := []struct {
		name string
		args args
		want AsciiBox
	}{
		{
			name: "test bool",
			args: args{
				name:      "exampleBool",
				anything:  true,
				charWidth: 0,
			},
			want: asciiBox(`
╔═exampleBool╗
║  b1 true   ║
╚════════════╝
`),
		},
		{
			name: "test int",
			args: args{
				name:      "exampleInt",
				anything:  1,
				charWidth: 0,
			},
			want: asciiBox(`
╔═exampleInt═════════╗
║0x0000000000000001 1║
╚════════════════════╝
`),
		},
		{
			name: "test int 123123123",
			args: args{
				name:      "exampleInt",
				anything:  123123123,
				charWidth: 0,
			},
			want: asciiBox(`
╔═exampleInt═════════════════╗
║0x000000000756b5b3 123123123║
╚════════════════════════════╝
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want = trimBox(tt.want)
			if got := BoxAnything(tt.args.name, tt.args.anything, tt.args.charWidth); got != tt.want {
				t.Errorf("BoxAnything() = '\n%v\n', want '\n%v\n'", got, tt.want)
			}
		})
	}
}

func TestBoxedDump(t *testing.T) {
	type args struct {
		name string
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Dump",
			args: args{
				name: "super nice data",
				data: []byte("1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aa1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aab"),
			},
			want: `
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
╚══════════════════════════════════════════════╝
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoxedDump(tt.args.name, tt.args.data); got != asciiBox(strings.Trim(tt.want, "\n")) {
				t.Errorf("Dump() = \n%v\n, want \n%v\n", got, tt.want)
			}
		})
	}
}

func TestDumpAnything(t *testing.T) {
	type args struct {
		anything interface{}
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
110|01 01 66 00 00                '..f..     '
`,
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
			if got := DumpAnything(tt.args.anything); got != strings.Trim(tt.want, "\n") {
				t.Errorf("Dump() = \n%v\n, want \n%v\n", got, tt.want)
			}
		})
	}
}

func TestBoxedDumpFixedWidth(t *testing.T) {
	type args struct {
		name      string
		data      []byte
		charWidth int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Dump",
			args: args{
				name:      "super nice data",
				data:      []byte("1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aa1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aab"),
				charWidth: 110,
			},
			want: `
╔═super nice data════════════════════════════════════════════════════════════════════════════════════════════╗
║ 000|31 32 33 34 35 36 37 38 39 30 61 62 63 64 65 66 67 68 69 6a 6b 6c 6d 6e 6f '1234567890abcdefghijklmno' ║
║ 025|70 71 72 73 74 75 76 77 78 79 7a d3 31 32 33 34 35 36 37 38 39 30 61 62 63 'pqrstuvwxyz.1234567890abc' ║
║ 050|64 65 66 67 68 69 6a 6b 6c 6d 6e 6f 70 71 72 73 74 75 76 77 78 79 7a d3 61 'defghijklmnopqrstuvwxyz.a' ║
║ 075|61 31 32 33 34 35 36 37 38 39 30 61 62 63 64 65 66 67 68 69 6a 6b 6c 6d 6e 'a1234567890abcdefghijklmn' ║
║ 100|6f 70 71 72 73 74 75 76 77 78 79 7a d3 31 32 33 34 35 36 37 38 39 30 61 62 'opqrstuvwxyz.1234567890ab' ║
║ 125|63 64 65 66 67 68 69 6a 6b 6c 6d 6e 6f 70 71 72 73 74 75 76 77 78 79 7a d3 'cdefghijklmnopqrstuvwxyz.' ║
║ 150|61 61 62                                                                   'aab                      ' ║
╚════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoxedDumpFixedWidth(tt.args.name, tt.args.data, tt.args.charWidth); got != asciiBox(strings.Trim(tt.want, "\n")) {
				t.Errorf("Dump() = \n%v\n, want \n%v\n", got, tt.want)
			}
		})
	}
}

func TestBoxedDumpAnythingFixedWidth(t *testing.T) {
	type args struct {
		name      string
		anything  interface{}
		charWidth int
	}
	tests := []struct {
		name string
		args args
		want AsciiBox
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoxedDumpAnythingFixedWidth(tt.args.name, tt.args.anything, tt.args.charWidth); got != tt.want {
				t.Errorf("BoxedDumpAnythingFixedWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDumpAnythingFixedWidth(t *testing.T) {
	type args struct {
		anything  interface{}
		charWidth int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DumpAnythingFixedWidth(tt.args.anything, tt.args.charWidth); got != tt.want {
				t.Errorf("DumpAnythingFixedWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}
