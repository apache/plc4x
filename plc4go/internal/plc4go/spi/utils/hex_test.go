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

import "testing"

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
			want: `000 0x: 31  32  33  34  35  36  37  38  '12345678'
008 0x: 39  30  61  62  63  64  65  66  '90abcdef'
016 0x: 67  68  69  6a  6b  6c  6d  6e  'ghijklmn'
024 0x: 6f  70  71  72  73  74  75  76  'opqrstuv'
032 0x: 77  78  79  7a                  'wxyz    '`,
		},
		{
			name: "minimum size",
			args: args{
				[]byte("a"),
			},
			want: "000 0x: 61                              'a       '",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Dump(tt.args.data); got != tt.want {
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
			want: `╔═super nice data══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║  000 0x: 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  '1234567890abcdefghijklmn'  ║
║  024 0x: 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  'opqrstuvwxyz.1234567890a'  ║
║  048 0x: 62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  'bcdefghijklmnopqrstuvwxy'  ║
║  072 0x: 7a  d3  61  61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  'z.aa1234567890abcdefghij'  ║
║  096 0x: 6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  'klmnopqrstuvwxyz.1234567'  ║
║  120 0x: 38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  '890abcdefghijklmnopqrstu'  ║
║  144 0x: 76  77  78  79  7a  d3  61  61  62                                                              'vwxyz.aab               '  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoxedDumpFixedWidth(tt.args.name, tt.args.data, tt.args.charWidth); got != tt.want {
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
			want: `000 0x: 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  '1234567890abcdefghijklmno'
025 0x: 70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  62  63  'pqrstuvwxyz.1234567890abc'
050 0x: 64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  61  'defghijklmnopqrstuvwxyz.a'
075 0x: 61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  'a1234567890abcdefghijklmn'
100 0x: 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  62  'opqrstuvwxyz.1234567890ab'
125 0x: 63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  'cdefghijklmnopqrstuvwxyz.'
150 0x: 61  61  62                                                                                          'aab                      '`,
		},
		{
			name: "minimum size",
			args: args{
				[]byte("a"),
				1,
			},
			want: "000 0x: 61  'a'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DumpFixedWidth(tt.args.data, tt.args.charWidth); got != tt.want {
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
