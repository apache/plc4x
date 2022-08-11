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
	"math"
	"reflect"
	"strings"
	"testing"
)

func init() {
	DebugHex = true
}

func TestDump(t *testing.T) {
	type args struct {
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
				data: []byte("1234567890abcdefghijklmnopqrstuvwxyz"),
			},
			want: `
00|31 32 33 34 35 36 37 38 39 30 '1234567890'
10|61 62 63 64 65 66 67 68 69 6a 'abcdefghij'
20|6b 6c 6d 6e 6f 70 71 72 73 74 'klmnopqrst'
30|75 76 77 78 79 7a             'uvwxyz    '
`,
		},
		{
			name: "Test Bigger Dump",
			args: args{
				data: []byte(strings.Repeat("Lorem ipsum", 90)),
			},
			want: `
000|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'
010|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'
020|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'
030|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'
040|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '
050|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'
060|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'
070|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'
080|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'
090|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'
100|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'
110|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'
120|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'
130|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'
140|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'
150|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '
160|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'
170|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'
180|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'
190|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'
200|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'
210|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'
220|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'
230|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'
240|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'
250|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'
260|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '
270|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'
280|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'
290|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'
300|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'
310|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'
320|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'
330|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'
340|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'
350|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'
360|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'
370|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '
380|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'
390|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'
400|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'
410|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'
420|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'
430|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'
440|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'
450|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'
460|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'
470|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'
480|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '
490|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'
500|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'
510|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'
520|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'
530|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'
540|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'
550|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'
560|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'
570|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'
580|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'
590|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '
600|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'
610|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'
620|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'
630|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'
640|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'
650|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'
660|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'
670|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'
680|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'
690|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'
700|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '
710|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'
720|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'
730|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'
740|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'
750|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'
760|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'
770|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'
780|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'
790|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'
800|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'
810|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '
820|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'
830|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'
840|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'
850|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'
860|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'
870|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'
880|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'
890|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'
900|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'
910|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'
920|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '
930|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'
940|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'
950|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'
960|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'
970|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'
980|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'
`,
		},
		{
			name: "no size",
			args: args{
				[]byte("a"),
			},
			want: "0|61                            'a         '",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Dump(tt.args.data); got != strings.Trim(tt.want, "\n") {
				t.Errorf("Dump() = \n%v\n, want \n%v\n", got, tt.want)
			}
		})
	}
}

func BenchmarkTestDump(b *testing.B) {
	DebugHex = false
	type args struct {
		data []byte
	}
	benchmarks := []struct {
		name string
		args args
	}{
		{
			"small",
			args{
				data: []byte(strings.Repeat("Lorem ipsum", 1)),
			},
		},
		{
			"medium",
			args{
				data: []byte(strings.Repeat("Lorem ipsum", 100)),
			},
		},
		{
			"big",
			args{
				data: []byte(strings.Repeat("Lorem ipsum", 10000)),
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Dump(bm.args.data)
			}
		})
	}
	DebugHex = true
}

func TestDumpFixedWidth(t *testing.T) {
	type args struct {
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
				data:      []byte("1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aa1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aab"),
				charWidth: 110,
			},
			want: `
000|31 32 33 34 35 36 37 38 39 30 61 62 63 64 65 66 67 68 69 6a 6b 6c 6d 6e 6f 70 '1234567890abcdefghijklmnop'
026|71 72 73 74 75 76 77 78 79 7a d3 31 32 33 34 35 36 37 38 39 30 61 62 63 64 65 'qrstuvwxyz.1234567890abcde'
052|66 67 68 69 6a 6b 6c 6d 6e 6f 70 71 72 73 74 75 76 77 78 79 7a d3 61 61 31 32 'fghijklmnopqrstuvwxyz.aa12'
078|33 34 35 36 37 38 39 30 61 62 63 64 65 66 67 68 69 6a 6b 6c 6d 6e 6f 70 71 72 '34567890abcdefghijklmnopqr'
104|73 74 75 76 77 78 79 7a d3 31 32 33 34 35 36 37 38 39 30 61 62 63 64 65 66 67 'stuvwxyz.1234567890abcdefg'
130|68 69 6a 6b 6c 6d 6e 6f 70 71 72 73 74 75 76 77 78 79 7a d3 61 61 62          'hijklmnopqrstuvwxyz.aab   '
`,
		},
		{
			name: "minimum size",
			args: args{
				[]byte("a"),
				1,
			},
			want: "0|61 'a'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DumpFixedWidth(tt.args.data, tt.args.charWidth); got != strings.Trim(tt.want, "\n") {
				t.Errorf("Dump() = \n%v\n, want \n%v\n", got, tt.want)
			}
		})
	}
}

func Test_maskString(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "printable chars",
			args: args{[]byte("1234567890abcdefghijklmnopqrstuvwxyz")},
			want: "1234567890abcdefghijklmnopqrstuvwxyz",
		},
		{
			name: "unprintable chars",
			args: args{[]byte("\3231234567890abcdefghijklmnopqrstuvwxyz.\323")},
			want: ".1234567890abcdefghijklmnopqrstuvwxyz..",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maskString(tt.args.data); got != tt.want {
				t.Errorf("maskString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_MinMax(t *testing.T) {
	type args struct {
		data      []byte
		charWidth int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil data",
			args: args{
				data:      nil,
				charWidth: math.MinInt32,
			},
			want: "",
		}, {
			name: "empty data",
			args: args{
				data:      []byte{},
				charWidth: math.MinInt32,
			},
			want: "",
		},
		{
			name: "-1 one byte",
			args: args{
				data:      []byte{0x1},
				charWidth: -1,
			},
			want: "0|01 '.'",
		},
		{
			name: "12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DumpFixedWidth(tt.args.data, tt.args.charWidth); got != tt.want {
				t.Errorf("maskString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark(b *testing.B) {
	DebugHex = false
	type args struct {
		numberOfBytes, desiredStringWidth int
	}
	benchmarks := []struct {
		name string
		args args
	}{
		{
			"small",
			args{
				numberOfBytes:      1,
				desiredStringWidth: 120,
			},
		},
		{
			"medium",
			args{
				numberOfBytes:      1000,
				desiredStringWidth: 120,
			},
		},
		{
			"big",
			args{
				numberOfBytes:      1000000,
				desiredStringWidth: 120,
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				calculateBytesPerRowAndIndexWidth(bm.args.numberOfBytes, bm.args.desiredStringWidth)
			}
		})
	}
	DebugHex = true
}

func Test_calculateBytesPerRowAndIndexWidth(t *testing.T) {
	type args struct {
		numberOfBytes    int
		desiredCharWidth int
	}
	tests := []struct {
		name               string
		args               args
		wantMaxBytesPerRow int
		wantIndexWidth     int
	}{
		{
			name: "1 byte MinInt32 width",
			args: args{
				numberOfBytes:    1,
				desiredCharWidth: math.MinInt32,
			},
			wantMaxBytesPerRow: 1,
			wantIndexWidth:     1,
		},
		{
			name: "10 byte MinInt32 width",
			args: args{
				numberOfBytes:    10,
				desiredCharWidth: math.MinInt32,
			},
			wantMaxBytesPerRow: 1,
			wantIndexWidth:     2,
		},
		{
			name: "100 byte MinInt32 width",
			args: args{
				numberOfBytes:    100,
				desiredCharWidth: math.MinInt32,
			},
			wantMaxBytesPerRow: 1,
			wantIndexWidth:     3,
		},
		{
			name: "100 byte MaxInt32 width",
			args: args{
				numberOfBytes:    100,
				desiredCharWidth: math.MaxInt32,
			},
			wantMaxBytesPerRow: 536870910,
			wantIndexWidth:     3,
		},
		{
			name: "100 byte 12 width",
			args: args{
				numberOfBytes:    100,
				desiredCharWidth: 12,
			},
			wantMaxBytesPerRow: 1,
			wantIndexWidth:     3,
		},
		{
			name: "153 byte 136 width",
			args: args{
				numberOfBytes:    153,
				desiredCharWidth: 136,
			},
			wantMaxBytesPerRow: 32,
			wantIndexWidth:     3,
		},
		{
			name: "153 byte calculated width",
			args: args{
				numberOfBytes: 153,
				desiredCharWidth: func() int {
					const quoteRune = 1
					const numberOfBytes = 153
					const indexWidth = 3
					const charRepresentation = 1

					// 000 AF FE AF FE ..... '....*'
					return indexWidth + blankWidth + (numberOfBytes * byteWidth) + quoteRune + (numberOfBytes * charRepresentation) + quoteRune
				}(),
			},
			wantMaxBytesPerRow: 153,
			wantIndexWidth:     3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := calculateBytesPerRowAndIndexWidth(tt.args.numberOfBytes, tt.args.desiredCharWidth)
			if got != tt.wantMaxBytesPerRow {
				t.Errorf("calculateBytesPerRowAndIndexWidth() got max bytes per row = %v, want %v", got, tt.wantMaxBytesPerRow)
			}
			if got1 != tt.wantIndexWidth {
				t.Errorf("calculateBytesPerRowAndIndexWidth() got index width = %v, want %v", got1, tt.wantIndexWidth)
			}
		})
	}
}

func Test_Immutability(t *testing.T) {
	inputBytes := []byte{0, 1, 2, 46, 56, 0, 200}
	_ = Dump(inputBytes)
	if !reflect.DeepEqual(inputBytes, []byte{0, 1, 2, 46, 56, 0, 200}) {
		t.Errorf("Dump has mutated the array got:=%x, want:=%x", inputBytes, []byte{0, 1, 2, 46, 56, 0, 200})
	}
}
