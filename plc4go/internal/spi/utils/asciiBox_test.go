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
	"fmt"
	"strings"
	"testing"
)

func init() {
	DebugAsciiBox = true
}

func TestAsciiBox_GetBoxName(t *testing.T) {
	type args struct {
		box AsciiBox
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple name",
			args: args{
				box: AsciiBoxWriterDefault.BoxString("someName", "some content", 0),
			},
			want: "someName",
		},
		{
			name: "no name",
			args: args{
				box: AsciiBoxWriterDefault.BoxString("", "some content", 0),
			},
			want: "",
		},
		{
			name: "long name",
			args: args{
				box: AsciiBoxWriterDefault.BoxString("veryLongName12_13", "some content", 0),
			},
			want: "veryLongName12_13",
		},
		{
			name: "name with spaces and slashes",
			args: args{
				box: AsciiBoxWriterDefault.BoxString("payload / Message / Concrete Message", "some content", 0),
			},
			want: "payload / Message / Concrete Message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.box.GetBoxName(); got != tt.want {
				t.Errorf("AsciiBox_GetBoxName() = '\n%v\n', want '\n%v\n'", got, tt.want)
			}
		})
	}
}

func TestAsciiBox_ChangeBoxName(t *testing.T) {
	type args struct {
		box     AsciiBox
		newName string
	}
	tests := []struct {
		name string
		args args
		want AsciiBox
	}{
		{
			name: "box with simple name",
			args: args{
				box:     AsciiBoxWriterDefault.BoxString("simpleName", "some content", 0),
				newName: "newSimpleName",
			},
			want: AsciiBoxWriterDefault.BoxString("newSimpleName", "some content", 0),
		},
		{
			name: "box with shorter name",
			args: args{
				box:     AsciiBoxWriterDefault.BoxString("veryLongName", "some content", 0),
				newName: "name",
			},
			want: AsciiBoxWriterDefault.BoxString("name", "some content", 0),
		},
		{
			name: "box getting dressed",
			args: args{
				box:     asciiBox("some content"),
				newName: "name",
			},
			want: AsciiBoxWriterDefault.BoxString("name", "some content", 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want = trimBox(tt.want)
			if got := tt.args.box.ChangeBoxName(tt.args.newName); got != tt.want {
				t.Errorf("BoxSideBySide() = '\n%v\n', want '\n%v\n'", got, tt.want)
			}
		})
	}
}

func TestAsciiBox_IsEmpty(t *testing.T) {
	type args struct {
		box AsciiBox
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty box",
			args: args{
				box: asciiBox(""),
			},
			want: true,
		},
		{
			name: "non empty box",
			args: args{
				box: asciiBox("a"),
			},
			want: false,
		},
		{
			name: "name empty box",
			args: args{
				box: AsciiBoxWriterDefault.BoxString("name", "", 0),
			},
			want: true,
		},
		{
			name: "name non empty box",
			args: args{
				box: AsciiBoxWriterDefault.BoxString("name", "a", 0),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.box.IsEmpty(); got != tt.want {
				t.Errorf("AsciiBox_IsEmpty() = '\n%v\n', want '\n%v\n'", got, tt.want)
			}
		})
	}
}

func TestBoxSideBySide(t *testing.T) {
	type args struct {
		box1 AsciiBox
		box2 AsciiBox
	}
	tests := []struct {
		name string
		args args
		want AsciiBox
	}{
		{
			name: "Test2Boxes",
			args: args{
				box1: asciiBox(`
000 0x: 31  32  33  34  35  36  37  38  '12345678'
008 0x: 39  30  61  62  63  64  65  66  '90abcdef'
016 0x: 67  68  69  6a  6b  6c  6d  6e  'ghijklmn'
024 0x: 6f  70  71  72  73  74  75  76  'opqrstuv'
032 0x: 77  78  79  7a                  'wxyz    '
`),
				box2: asciiBox(`
╔═super nice data══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║  000 0x: 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  '1234567890abcdefghijklmn'  ║
║  024 0x: 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  'opqrstuvwxyz.1234567890a'  ║
║  048 0x: 62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  'bcdefghijklmnopqrstuvwxy'  ║
║  072 0x: 7a  d3  61  61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  'z.aa1234567890abcdefghij'  ║
║  096 0x: 6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  'klmnopqrstuvwxyz.1234567'  ║
║  120 0x: 38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  '890abcdefghijklmnopqrstu'  ║
║  144 0x: 76  77  78  79  7a  d3  61  61  62                                                              'vwxyz.aab               '  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
`),
			},
			want: asciiBox(`
000 0x: 31  32  33  34  35  36  37  38  '12345678'╔═super nice data══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
008 0x: 39  30  61  62  63  64  65  66  '90abcdef'║  000 0x: 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  '1234567890abcdefghijklmn'  ║
016 0x: 67  68  69  6a  6b  6c  6d  6e  'ghijklmn'║  024 0x: 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  'opqrstuvwxyz.1234567890a'  ║
024 0x: 6f  70  71  72  73  74  75  76  'opqrstuv'║  048 0x: 62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  'bcdefghijklmnopqrstuvwxy'  ║
032 0x: 77  78  79  7a                  'wxyz    '║  072 0x: 7a  d3  61  61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  'z.aa1234567890abcdefghij'  ║
                                                  ║  096 0x: 6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  'klmnopqrstuvwxyz.1234567'  ║
                                                  ║  120 0x: 38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  '890abcdefghijklmnopqrstu'  ║
                                                  ║  144 0x: 76  77  78  79  7a  d3  61  61  62                                                              'vwxyz.aab               '  ║
                                                  ╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
`),
		},
		{
			name: "another 2 boxes",
			args: args{
				box1: asciiBox(`
╔═exampleInt╗
║     4     ║
╚═══════════╝
`),
				box2: asciiBox(`
╔═exampleInt╗
║     7     ║
╚═══════════╝
`),
			},
			want: asciiBox(`
╔═exampleInt╗╔═exampleInt╗
║     4     ║║     7     ║
╚═══════════╝╚═══════════╝
`),
		},
		{
			name: "size difference first box",
			args: args{
				box1: asciiBox(`
╔═exampleInt╗
║     4     ║
║     4     ║
╚═══════════╝
`),
				box2: asciiBox(`
╔═exampleInt╗
║     7     ║
╚═══════════╝
`),
			},
			want: asciiBox(`
╔═exampleInt╗╔═exampleInt╗
║     4     ║║     7     ║
║     4     ║╚═══════════╝
╚═══════════╝             
`),
		},
		{
			name: "size difference second box",
			args: args{
				box1: asciiBox(`
╔═exampleInt╗
║     4     ║
╚═══════════╝
`),
				box2: asciiBox(`
╔═exampleInt╗
║     7     ║
║     7     ║
╚═══════════╝
`),
			},
			want: asciiBox(`
╔═exampleInt╗╔═exampleInt╗
║     4     ║║     7     ║
╚═══════════╝║     7     ║
             ╚═══════════╝
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want = trimBox(tt.want)
			if got := AsciiBoxWriterDefault.BoxSideBySide(trimBox(tt.args.box1), trimBox(tt.args.box2)); got != tt.want {
				t.Errorf("BoxSideBySide() = '\n%v\n', want '\n%v\n'", got, tt.want)
			}
		})
	}
}

func BenchmarkBoxSideBySide(b *testing.B) {
	oldSetting := DebugAsciiBox
	DebugAsciiBox = false
	bigString := strings.Repeat(strings.Repeat("LoreIpsum", 100)+"\n", 100)
	box := AsciiBoxWriterDefault.BoxString("RandomBox", bigString, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AsciiBoxWriterDefault.BoxSideBySide(box, box)
	}
	DebugAsciiBox = oldSetting
}

func TestBoxBelowBox(t *testing.T) {
	type args struct {
		box1 AsciiBox
		box2 AsciiBox
	}
	tests := []struct {
		name string
		args args
		want AsciiBox
	}{
		{
			name: "Test2Boxes",
			args: args{
				box1: asciiBox(`
000 31  32  33  34  35  36  37  38  '12345678'
008 39  30  61  62  63  64  65  66  '90abcdef'
016 67  68  69  6a  6b  6c  6d  6e  'ghijklmn'
024 6f  70  71  72  73  74  75  76  'opqrstuv'
032 77  78  79  7a                  'wxyz    '
`),
				box2: asciiBox(`
╔═super nice data══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║  000 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  '1234567890abcdefghijklmn'  ║
║  024 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  'opqrstuvwxyz.1234567890a'  ║
║  048 62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  'bcdefghijklmnopqrstuvwxy'  ║
║  072 7a  d3  61  61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  'z.aa1234567890abcdefghij'  ║
║  096 6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  'klmnopqrstuvwxyz.1234567'  ║
║  120 38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  '890abcdefghijklmnopqrstu'  ║
║  144 76  77  78  79  7a  d3  61  61  62                                                              'vwxyz.aab               '  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
`),
			},
			want: asciiBox(`
000 31  32  33  34  35  36  37  38  '12345678'                                                                                      
008 39  30  61  62  63  64  65  66  '90abcdef'                                                                                      
016 67  68  69  6a  6b  6c  6d  6e  'ghijklmn'                                                                                      
024 6f  70  71  72  73  74  75  76  'opqrstuv'                                                                                      
032 77  78  79  7a                  'wxyz    '                                                                                      
╔═super nice data══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║  000 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  '1234567890abcdefghijklmn'  ║
║  024 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  'opqrstuvwxyz.1234567890a'  ║
║  048 62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  'bcdefghijklmnopqrstuvwxy'  ║
║  072 7a  d3  61  61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  'z.aa1234567890abcdefghij'  ║
║  096 6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  'klmnopqrstuvwxyz.1234567'  ║
║  120 38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  '890abcdefghijklmnopqrstu'  ║
║  144 76  77  78  79  7a  d3  61  61  62                                                              'vwxyz.aab               '  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
`),
		},
		{
			name: "different sized boxes",
			args: args{
				box1: asciiBox(`
╔═sampleField════════════╗
║123123123123123123123123║
╚════════════════════════╝
`),
				box2: asciiBox(`
╔═sampleField╗
║123123123123║
╚════════════╝
`),
			},
			want: asciiBox(`
╔═sampleField════════════╗
║123123123123123123123123║
╚════════════════════════╝
╔═sampleField╗            
║123123123123║            
╚════════════╝            
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want = trimBox(tt.want)
			if got := AsciiBoxWriterDefault.BoxBelowBox(trimBox(tt.args.box1), trimBox(tt.args.box2)); got != tt.want {
				t.Errorf("BoxSideBySide() = '\n%v\n', want '\n%v\n'", got, tt.want)
			}
		})
	}
}

func TestBoxString(t *testing.T) {
	type args struct {
		name      string
		data      string
		charWidth int
	}
	tests := []struct {
		name string
		args args
		want AsciiBox
	}{
		{
			name: "simplebox",
			args: args{
				name:      "sampleField",
				data:      "123123123123",
				charWidth: 1,
			},
			want: asciiBox(`
╔═sampleField╗
║123123123123║
╚════════════╝
`),
		},
		{
			name: "simplebox-unamed",
			args: args{
				name:      "",
				data:      "123123123123",
				charWidth: 1,
			},
			want: asciiBox(`
╔════════════╗
║123123123123║
╚════════════╝
`),
		},
		{
			name: "simplebox 2",
			args: args{
				name:      "sampleField",
				data:      "123123123123\n123123123123123123123123",
				charWidth: 1,
			},
			want: asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝
`),
		},
		{
			name: "simplebox with too long name",
			args: args{
				name:      "sampleFieldsampleFieldsampleFieldsampleField",
				data:      "123123123123\n123123123123123123123123",
				charWidth: 1,
			},
			want: asciiBox(`
╔═sampleFieldsampleFieldsampleFieldsampleField╗
║                123123123123                 ║
║          123123123123123123123123           ║
╚═════════════════════════════════════════════╝
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want = trimBox(tt.want)
			if got := AsciiBoxWriterDefault.BoxString(tt.args.name, tt.args.data, tt.args.charWidth); got != tt.want {
				t.Errorf("BoxString() = '\n%v\n', want '\n%v\n'", got, tt.want)
			}
		})
	}
}

func BenchmarkBoxString(b *testing.B) {
	oldSetting := DebugAsciiBox
	DebugAsciiBox = false
	bigString := strings.Repeat(strings.Repeat("LoreIpsum", 100)+"\n", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AsciiBoxWriterDefault.BoxString("randomName", bigString, 50)
	}
	DebugAsciiBox = oldSetting
}

func TestAlignBoxes(t *testing.T) {
	type args struct {
		boxes       []AsciiBox
		desiredWith int
	}
	tests := []struct {
		name string
		args args
		want AsciiBox
	}{
		{
			name: "enough space",
			args: args{
				boxes: []AsciiBox{
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝
`),
				},
				desiredWith: 1000,
			},
			want: asciiBox(`
╔═sampleField════════════╗╔═sampleField════════════╗
║      123123123123      ║║      123123123123      ║
║123123ABABABABABAB123123║║123123123123123123123123║
╚════════════════════════╝╚════════════════════════╝
`),
		},
		{
			name: "not enough space",
			args: args{
				boxes: []AsciiBox{
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝
`),
				},
				desiredWith: 0,
			},
			want: asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝
`),
		},
		{
			name: "not enough space should result in multiple rows",
			args: args{
				boxes: []AsciiBox{
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝
`),
				},
				desiredWith: 65,
			},
			want: asciiBox(`
╔═sampleField════════════╗╔═sampleField════════════╗
║      123123123123      ║║      123123123123      ║
║123123ABABABABABAB123123║║123123123123123123123123║
╚════════════════════════╝╚════════════════════════╝
╔═sampleField════════════╗╔═sampleField════════════╗
║      123123123123      ║║      123123123123      ║
║123123ABABABABABAB123123║║123123123123123123123123║
╚════════════════════════╝╚════════════════════════╝
╔═sampleField════════════╗╔═sampleField════════════╗
║      123123123123      ║║      123123123123      ║
║123123ABABABABABAB123123║║123123123123123123123123║
╚════════════════════════╝╚════════════════════════╝
╔═sampleField════════════╗╔═sampleField════════════╗
║      123123123123      ║║      123123123123      ║
║123123ABABABABABAB123123║║123123123123123123123123║
╚════════════════════════╝╚════════════════════════╝
`),
		},
		{
			name: "not enough space should result in multiple rows (3 columns)",
			args: args{
				boxes: []AsciiBox{
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝
`),
					asciiBox(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝
`),
				},
				desiredWith: 78,
			},
			want: asciiBox(`
╔═sampleField════════════╗╔═sampleField════════════╗╔═sampleField════════════╗
║      123123123123      ║║      123123123123      ║║      123123123123      ║
║123123ABABABABABAB123123║║123123123123123123123123║║123123ABABABABABAB123123║
╚════════════════════════╝╚════════════════════════╝╚════════════════════════╝
╔═sampleField════════════╗╔═sampleField════════════╗╔═sampleField════════════╗
║      123123123123      ║║      123123123123      ║║      123123123123      ║
║123123123123123123123123║║123123ABABABABABAB123123║║123123123123123123123123║
╚════════════════════════╝╚════════════════════════╝╚════════════════════════╝
╔═sampleField════════════╗╔═sampleField════════════╗                          
║      123123123123      ║║      123123123123      ║                          
║123123ABABABABABAB123123║║123123123123123123123123║                          
╚════════════════════════╝╚════════════════════════╝                          
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, box := range tt.args.boxes {
				tt.args.boxes[i] = trimBox(box)
			}
			tt.want = trimBox(tt.want)
			if got := AsciiBoxWriterDefault.AlignBoxes(tt.args.boxes, tt.args.desiredWith); got != tt.want {
				t.Errorf("AlignBoxes() = '\n%v\n', want '\n%v\n'", got, tt.want)
			}
		})
	}
}

func TestAsciiBox_width(t *testing.T) {
	tests := []struct {
		name string
		m    AsciiBox
		want int
	}{
		{
			name: "same width",
			m: asciiBox(`
123123123123123
123123123123123
123123123123123
`),
			want: 15,
		},
		{
			name: "different width",
			m: asciiBox(`
123123123123123
123123123123123123123123123123
123123123123123
`),
			want: 30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Width(); got != tt.want {
				t.Errorf("width() = '\n%v\n', want '\n%v\n'", got, tt.want)
			}
		})
	}
}

func Test_mergeHorizontal(t *testing.T) {
	type args struct {
		boxes []AsciiBox
	}
	tests := []struct {
		name string
		args args
		want AsciiBox
	}{
		{
			name: "3 same",
			args: args{
				boxes: []AsciiBox{
					asciiBox(`
123123123
123123123
123123123
`),
					asciiBox(`
abcabcabc
abcabcabc
abcabcabc
`),
					asciiBox(`
zxyzxyzxy
zxyzxyzxy
zxyzxyzxy
`),
				},
			},
			want: asciiBox(`
123123123abcabcabczxyzxyzxy
123123123abcabcabczxyzxyzxy
123123123abcabcabczxyzxyzxy
`),
		},
		{
			name: "3 different",
			args: args{
				boxes: []AsciiBox{
					asciiBox(`
123123123
123123123
123123123
`),
					asciiBox(`
abcabcabc
abcabcabcabcabcabcabcabcabc
abcabcabc
`),
					asciiBox(`
zxyzxyzxy
zxyzxyzxy
zxyzxyzxy
`),
				},
			},
			want: asciiBox(`
123123123abcabcabc                  zxyzxyzxy
123123123abcabcabcabcabcabcabcabcabczxyzxyzxy
123123123abcabcabc                  zxyzxyzxy
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, box := range tt.args.boxes {
				tt.args.boxes[i] = trimBox(box)
			}
			tt.want = trimBox(tt.want)
			if got := AsciiBoxWriterDefault.(*asciiBoxWriter).mergeHorizontal(tt.args.boxes); got != tt.want {
				t.Errorf("mergeHorizontal() = '\n%v\n', want '\n%v\n'", got, tt.want)
			}
		})
	}
}

func TestExpandBox(t *testing.T) {
	type args struct {
		box   AsciiBox
		width int
	}
	tests := []struct {
		name string
		args
		want AsciiBox
	}{
		{
			name: "Small expand",
			args: args{
				box: asciiBox(`
123123123
123123123
123123123
`),
				width: 100,
			},
			want: asciiBox(`
123123123                                                                                           
123123123                                                                                           
123123123                                                                                           
`),
		},
		{
			name: "Big expand",
			args: args{
				box: asciiBox(`
123123123
123123123
123123123
`),
				width: 10000,
			},
			want: asciiBox(fmt.Sprintf(`
123123123%[1]s
123123123%[1]s
123123123%[1]s
`, strings.Repeat(" ", 10000-9))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.box = trimBox(tt.args.box)
			tt.want = trimBox(tt.want)
			if got := AsciiBoxWriterDefault.(*asciiBoxWriter).expandBox(tt.args.box, tt.args.width); got != tt.want {
				t.Errorf("mergeHorizontal() = '\n%v\n', want '\n%v\n'", got, tt.want)
			}
		})
	}
}

func BenchmarkExpandBox(b *testing.B) {
	oldSetting := DebugAsciiBox
	DebugAsciiBox = false
	bigString := strings.Repeat(strings.Repeat("LoreIpsum", 100)+"\n", 100)
	box := AsciiBoxWriterDefault.BoxString("RandomBox", bigString, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AsciiBoxWriterDefault.(*asciiBoxWriter).expandBox(box, 10000)
	}
	DebugAsciiBox = oldSetting
}

func trimBox(box AsciiBox) AsciiBox {
	return AsciiBox{strings.Trim(box.String(), "\n"), AsciiBoxWriterDefault.(*asciiBoxWriter)}
}

func asciiBox(value string) AsciiBox {
	return AsciiBox{value, AsciiBoxWriterDefault.(*asciiBoxWriter)}
}
