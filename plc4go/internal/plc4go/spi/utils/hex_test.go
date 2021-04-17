//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package utils

import (
	"math"
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
00 31  32  33  34  35  36  37  38  39  30  '1234567890'
10 61  62  63  64  65  66  67  68  69  6a  'abcdefghij'
20 6b  6c  6d  6e  6f  70  71  72  73  74  'klmnopqrst'
30 75  76  77  78  79  7a                  'uvwxyz    '
`,
		},
		{
			name: "Test Bigger Dump",
			args: args{
				data: []byte(strings.Repeat("Lorem ipsum", 90)),
			},
			want: `
000 4c  6f  72  65  6d  20  69  70  73  75  'Lorem ipsu'
010 6d  4c  6f  72  65  6d  20  69  70  73  'mLorem ips'
020 75  6d  4c  6f  72  65  6d  20  69  70  'umLorem ip'
030 73  75  6d  4c  6f  72  65  6d  20  69  'sumLorem i'
040 70  73  75  6d  4c  6f  72  65  6d  20  'psumLorem '
050 69  70  73  75  6d  4c  6f  72  65  6d  'ipsumLorem'
060 20  69  70  73  75  6d  4c  6f  72  65  ' ipsumLore'
070 6d  20  69  70  73  75  6d  4c  6f  72  'm ipsumLor'
080 65  6d  20  69  70  73  75  6d  4c  6f  'em ipsumLo'
090 72  65  6d  20  69  70  73  75  6d  4c  'rem ipsumL'
100 6f  72  65  6d  20  69  70  73  75  6d  'orem ipsum'
110 4c  6f  72  65  6d  20  69  70  73  75  'Lorem ipsu'
120 6d  4c  6f  72  65  6d  20  69  70  73  'mLorem ips'
130 75  6d  4c  6f  72  65  6d  20  69  70  'umLorem ip'
140 73  75  6d  4c  6f  72  65  6d  20  69  'sumLorem i'
150 70  73  75  6d  4c  6f  72  65  6d  20  'psumLorem '
160 69  70  73  75  6d  4c  6f  72  65  6d  'ipsumLorem'
170 20  69  70  73  75  6d  4c  6f  72  65  ' ipsumLore'
180 6d  20  69  70  73  75  6d  4c  6f  72  'm ipsumLor'
190 65  6d  20  69  70  73  75  6d  4c  6f  'em ipsumLo'
200 72  65  6d  20  69  70  73  75  6d  4c  'rem ipsumL'
210 6f  72  65  6d  20  69  70  73  75  6d  'orem ipsum'
220 4c  6f  72  65  6d  20  69  70  73  75  'Lorem ipsu'
230 6d  4c  6f  72  65  6d  20  69  70  73  'mLorem ips'
240 75  6d  4c  6f  72  65  6d  20  69  70  'umLorem ip'
250 73  75  6d  4c  6f  72  65  6d  20  69  'sumLorem i'
260 70  73  75  6d  4c  6f  72  65  6d  20  'psumLorem '
270 69  70  73  75  6d  4c  6f  72  65  6d  'ipsumLorem'
280 20  69  70  73  75  6d  4c  6f  72  65  ' ipsumLore'
290 6d  20  69  70  73  75  6d  4c  6f  72  'm ipsumLor'
300 65  6d  20  69  70  73  75  6d  4c  6f  'em ipsumLo'
310 72  65  6d  20  69  70  73  75  6d  4c  'rem ipsumL'
320 6f  72  65  6d  20  69  70  73  75  6d  'orem ipsum'
330 4c  6f  72  65  6d  20  69  70  73  75  'Lorem ipsu'
340 6d  4c  6f  72  65  6d  20  69  70  73  'mLorem ips'
350 75  6d  4c  6f  72  65  6d  20  69  70  'umLorem ip'
360 73  75  6d  4c  6f  72  65  6d  20  69  'sumLorem i'
370 70  73  75  6d  4c  6f  72  65  6d  20  'psumLorem '
380 69  70  73  75  6d  4c  6f  72  65  6d  'ipsumLorem'
390 20  69  70  73  75  6d  4c  6f  72  65  ' ipsumLore'
400 6d  20  69  70  73  75  6d  4c  6f  72  'm ipsumLor'
410 65  6d  20  69  70  73  75  6d  4c  6f  'em ipsumLo'
420 72  65  6d  20  69  70  73  75  6d  4c  'rem ipsumL'
430 6f  72  65  6d  20  69  70  73  75  6d  'orem ipsum'
440 4c  6f  72  65  6d  20  69  70  73  75  'Lorem ipsu'
450 6d  4c  6f  72  65  6d  20  69  70  73  'mLorem ips'
460 75  6d  4c  6f  72  65  6d  20  69  70  'umLorem ip'
470 73  75  6d  4c  6f  72  65  6d  20  69  'sumLorem i'
480 70  73  75  6d  4c  6f  72  65  6d  20  'psumLorem '
490 69  70  73  75  6d  4c  6f  72  65  6d  'ipsumLorem'
500 20  69  70  73  75  6d  4c  6f  72  65  ' ipsumLore'
510 6d  20  69  70  73  75  6d  4c  6f  72  'm ipsumLor'
520 65  6d  20  69  70  73  75  6d  4c  6f  'em ipsumLo'
530 72  65  6d  20  69  70  73  75  6d  4c  'rem ipsumL'
540 6f  72  65  6d  20  69  70  73  75  6d  'orem ipsum'
550 4c  6f  72  65  6d  20  69  70  73  75  'Lorem ipsu'
560 6d  4c  6f  72  65  6d  20  69  70  73  'mLorem ips'
570 75  6d  4c  6f  72  65  6d  20  69  70  'umLorem ip'
580 73  75  6d  4c  6f  72  65  6d  20  69  'sumLorem i'
590 70  73  75  6d  4c  6f  72  65  6d  20  'psumLorem '
600 69  70  73  75  6d  4c  6f  72  65  6d  'ipsumLorem'
610 20  69  70  73  75  6d  4c  6f  72  65  ' ipsumLore'
620 6d  20  69  70  73  75  6d  4c  6f  72  'm ipsumLor'
630 65  6d  20  69  70  73  75  6d  4c  6f  'em ipsumLo'
640 72  65  6d  20  69  70  73  75  6d  4c  'rem ipsumL'
650 6f  72  65  6d  20  69  70  73  75  6d  'orem ipsum'
660 4c  6f  72  65  6d  20  69  70  73  75  'Lorem ipsu'
670 6d  4c  6f  72  65  6d  20  69  70  73  'mLorem ips'
680 75  6d  4c  6f  72  65  6d  20  69  70  'umLorem ip'
690 73  75  6d  4c  6f  72  65  6d  20  69  'sumLorem i'
700 70  73  75  6d  4c  6f  72  65  6d  20  'psumLorem '
710 69  70  73  75  6d  4c  6f  72  65  6d  'ipsumLorem'
720 20  69  70  73  75  6d  4c  6f  72  65  ' ipsumLore'
730 6d  20  69  70  73  75  6d  4c  6f  72  'm ipsumLor'
740 65  6d  20  69  70  73  75  6d  4c  6f  'em ipsumLo'
750 72  65  6d  20  69  70  73  75  6d  4c  'rem ipsumL'
760 6f  72  65  6d  20  69  70  73  75  6d  'orem ipsum'
770 4c  6f  72  65  6d  20  69  70  73  75  'Lorem ipsu'
780 6d  4c  6f  72  65  6d  20  69  70  73  'mLorem ips'
790 75  6d  4c  6f  72  65  6d  20  69  70  'umLorem ip'
800 73  75  6d  4c  6f  72  65  6d  20  69  'sumLorem i'
810 70  73  75  6d  4c  6f  72  65  6d  20  'psumLorem '
820 69  70  73  75  6d  4c  6f  72  65  6d  'ipsumLorem'
830 20  69  70  73  75  6d  4c  6f  72  65  ' ipsumLore'
840 6d  20  69  70  73  75  6d  4c  6f  72  'm ipsumLor'
850 65  6d  20  69  70  73  75  6d  4c  6f  'em ipsumLo'
860 72  65  6d  20  69  70  73  75  6d  4c  'rem ipsumL'
870 6f  72  65  6d  20  69  70  73  75  6d  'orem ipsum'
880 4c  6f  72  65  6d  20  69  70  73  75  'Lorem ipsu'
890 6d  4c  6f  72  65  6d  20  69  70  73  'mLorem ips'
900 75  6d  4c  6f  72  65  6d  20  69  70  'umLorem ip'
910 73  75  6d  4c  6f  72  65  6d  20  69  'sumLorem i'
920 70  73  75  6d  4c  6f  72  65  6d  20  'psumLorem '
930 69  70  73  75  6d  4c  6f  72  65  6d  'ipsumLorem'
940 20  69  70  73  75  6d  4c  6f  72  65  ' ipsumLore'
950 6d  20  69  70  73  75  6d  4c  6f  72  'm ipsumLor'
960 65  6d  20  69  70  73  75  6d  4c  6f  'em ipsumLo'
970 72  65  6d  20  69  70  73  75  6d  4c  'rem ipsumL'
980 6f  72  65  6d  20  69  70  73  75  6d  'orem ipsum'
`,
		},
		{
			name: "minimum size",
			args: args{
				[]byte("a"),
			},
			want: "0 61                                      'a         '",
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
╔═super nice data════════════════════════════════════════╗
║000 31  32  33  34  35  36  37  38  39  30  '1234567890'║
║010 61  62  63  64  65  66  67  68  69  6a  'abcdefghij'║
║020 6b  6c  6d  6e  6f  70  71  72  73  74  'klmnopqrst'║
║030 75  76  77  78  79  7a  d3  31  32  33  'uvwxyz.123'║
║040 34  35  36  37  38  39  30  61  62  63  '4567890abc'║
║050 64  65  66  67  68  69  6a  6b  6c  6d  'defghijklm'║
║060 6e  6f  70  71  72  73  74  75  76  77  'nopqrstuvw'║
║070 78  79  7a  d3  61  61  31  32  33  34  'xyz.aa1234'║
║080 35  36  37  38  39  30  61  62  63  64  '567890abcd'║
║090 65  66  67  68  69  6a  6b  6c  6d  6e  'efghijklmn'║
║100 6f  70  71  72  73  74  75  76  77  78  'opqrstuvwx'║
║110 79  7a  d3  31  32  33  34  35  36  37  'yz.1234567'║
║120 38  39  30  61  62  63  64  65  66  67  '890abcdefg'║
║130 68  69  6a  6b  6c  6d  6e  6f  70  71  'hijklmnopq'║
║140 72  73  74  75  76  77  78  79  7a  d3  'rstuvwxyz.'║
║150 61  61  62                              'aab       '║
╚════════════════════════════════════════════════════════╝
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoxedDump(tt.args.name, tt.args.data); got != strings.Trim(tt.want, "\n") {
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
				charWidth: 136,
			},
			want: `
╔═super nice data══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║ 000 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  '1234567890abcdefghijklmno'  ║
║ 025 70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  62  63  'pqrstuvwxyz.1234567890abc'  ║
║ 050 64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  61  'defghijklmnopqrstuvwxyz.a'  ║
║ 075 61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  'a1234567890abcdefghijklmn'  ║
║ 100 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  62  'opqrstuvwxyz.1234567890ab'  ║
║ 125 63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  'cdefghijklmnopqrstuvwxyz.'  ║
║ 150 61  61  62                                                                                          'aab                      '  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoxedDumpFixedWidth(tt.args.name, tt.args.data, tt.args.charWidth); got != strings.Trim(tt.want, "\n") {
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
000 25  ff  81  03  01  02  ff  82  00  01  '%.........'
010 04  01  01  41  01  0c  00  01  01  42  '...A.....B'
020 01  0c  00  01  01  43  01  0c  00  01  '.....C....'
030 01  44  01  ff  84  00  00  00  37  ff  '.D......7.'
040 83  03  01  01  1d  73  74  72  75  63  '.....struc'
050 74  20  7b  20  45  20  73  74  72  69  't { E stri'
060 6e  67  3b  20  46  20  73  74  72  69  'ng; F stri'
070 6e  67  20  7d  01  ff  84  00  01  02  'ng }......'
080 01  01  45  01  0c  00  01  01  46  01  '..E.....F.'
090 0c  00  00  00  14  ff  82  01  01  61  '.........a'
100 01  01  62  01  01  63  01  01  01  65  '..b..c...e'
110 01  01  66  00  00                      '..f..     '
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
				charWidth: 136,
			},
			want: `
000 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  '1234567890abcdefghijklmnop'
026 71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  'qrstuvwxyz.1234567890abcde'
052 66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  61  61  31  32  'fghijklmnopqrstuvwxyz.aa12'
078 33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  '34567890abcdefghijklmnopqr'
104 73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  'stuvwxyz.1234567890abcdefg'
130 68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  61  61  62              'hijklmnopqrstuvwxyz.aab   '
`,
		},
		{
			name: "minimum size",
			args: args{
				[]byte("a"),
				1,
			},
			want: "0 61  'a'",
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
			want: "0 01  '.'",
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
			wantMaxBytesPerRow: 429496728,
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
			wantMaxBytesPerRow: 26,
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
