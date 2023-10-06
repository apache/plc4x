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
	"github.com/stretchr/testify/assert"
	"regexp"
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
			if got := tt.args.box.GetBoxName(); !assert.Equal(t, tt.want, got) {
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
				box:     asciiBoxForTest("some content"),
				newName: "name",
			},
			want: AsciiBoxWriterDefault.BoxString("name", "some content", 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.box.ChangeBoxName(tt.args.newName); !assert.Equal(t, tt.want, got) {
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
				box: asciiBoxForTest(""),
			},
			want: true,
		},
		{
			name: "non empty box",
			args: args{
				box: asciiBoxForTest("a"),
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
			if got := tt.args.box.IsEmpty(); !assert.Equal(t, tt.want, got) {
				t.Errorf("AsciiBox_IsEmpty() = '\n%v\n', want '\n%v\n'", got, tt.want)
			}
		})
	}
}

func TestAsciiBox_Height(t *testing.T) {
	type fields struct {
		data             string
		asciiBoxWriter   *asciiBoxWriter
		compressedBoxSet string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "test height",
			want: 1,
		},
		{
			name: "a bit higher",
			fields: fields{
				data: "\n\n\n\n",
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := AsciiBox{
				data:             tt.fields.data,
				asciiBoxWriter:   tt.fields.asciiBoxWriter,
				compressedBoxSet: tt.fields.compressedBoxSet,
			}
			assert.Equalf(t, tt.want, m.Height(), "Height()")
		})
	}
}

func TestAsciiBox_Lines(t *testing.T) {
	type fields struct {
		data             string
		asciiBoxWriter   *asciiBoxWriter
		compressedBoxSet string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "test height",
			want: []string{""},
		},
		{
			name: "a bit higher",
			fields: fields{
				data: "\n\n\n\n",
			},
			want: []string{"", "", "", "", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := AsciiBox{
				data:             tt.fields.data,
				asciiBoxWriter:   tt.fields.asciiBoxWriter,
				compressedBoxSet: tt.fields.compressedBoxSet,
			}
			assert.Equalf(t, tt.want, m.Lines(), "Lines()")
		})
	}
}

func TestAsciiBox_String(t *testing.T) {
	type fields struct {
		data             string
		asciiBoxWriter   *asciiBoxWriter
		compressedBoxSet string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string returns data",
			fields: fields{
				data: "data",
			},
			want: "data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := AsciiBox{
				data:             tt.fields.data,
				asciiBoxWriter:   tt.fields.asciiBoxWriter,
				compressedBoxSet: tt.fields.compressedBoxSet,
			}
			assert.Equalf(t, tt.want, m.String(), "String()")
		})
	}
}

func TestAsciiBox_Width(t *testing.T) {
	type fields struct {
		data             string
		asciiBoxWriter   *asciiBoxWriter
		compressedBoxSet string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "width simple",
			fields: fields{
				data: "     ",
			},
			want: 5,
		},
		{
			name: "width is longest",
			fields: fields{
				data: `     
                             `,
			},
			want: 29,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := AsciiBox{
				data:             tt.fields.data,
				asciiBoxWriter:   tt.fields.asciiBoxWriter,
				compressedBoxSet: tt.fields.compressedBoxSet,
			}
			assert.Equalf(t, tt.want, m.Width(), "Width()")
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
				box1: asciiBoxForTest(`
000 0x: 31  32  33  34  35  36  37  38  '12345678'
008 0x: 39  30  61  62  63  64  65  66  '90abcdef'
016 0x: 67  68  69  6a  6b  6c  6d  6e  'ghijklmn'
024 0x: 6f  70  71  72  73  74  75  76  'opqrstuv'
032 0x: 77  78  79  7a                  'wxyz    '`[1:]),
				box2: asciiBoxForTest(`
╔═super nice data══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║  000 0x: 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  '1234567890abcdefghijklmn'  ║
║  024 0x: 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  'opqrstuvwxyz.1234567890a'  ║
║  048 0x: 62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  'bcdefghijklmnopqrstuvwxy'  ║
║  072 0x: 7a  d3  61  61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  'z.aa1234567890abcdefghij'  ║
║  096 0x: 6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  'klmnopqrstuvwxyz.1234567'  ║
║  120 0x: 38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  '890abcdefghijklmnopqrstu'  ║
║  144 0x: 76  77  78  79  7a  d3  61  61  62                                                              'vwxyz.aab               '  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝`[1:]),
			},
			want: asciiBoxForTest(`
000 0x: 31  32  33  34  35  36  37  38  '12345678'╔═super nice data══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
008 0x: 39  30  61  62  63  64  65  66  '90abcdef'║  000 0x: 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  '1234567890abcdefghijklmn'  ║
016 0x: 67  68  69  6a  6b  6c  6d  6e  'ghijklmn'║  024 0x: 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  'opqrstuvwxyz.1234567890a'  ║
024 0x: 6f  70  71  72  73  74  75  76  'opqrstuv'║  048 0x: 62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  'bcdefghijklmnopqrstuvwxy'  ║
032 0x: 77  78  79  7a                  'wxyz    '║  072 0x: 7a  d3  61  61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  'z.aa1234567890abcdefghij'  ║
                                                  ║  096 0x: 6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  'klmnopqrstuvwxyz.1234567'  ║
                                                  ║  120 0x: 38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  '890abcdefghijklmnopqrstu'  ║
                                                  ║  144 0x: 76  77  78  79  7a  d3  61  61  62                                                              'vwxyz.aab               '  ║
                                                  ╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝`[1:]),
		},
		{
			name: "another 2 boxes",
			args: args{
				box1: asciiBoxForTest(`
╔═exampleInt╗
║     4     ║
╚═══════════╝`[1:]),
				box2: asciiBoxForTest(`
╔═exampleInt╗
║     7     ║
╚═══════════╝`[1:]),
			},
			want: asciiBoxForTest(`
╔═exampleInt╗╔═exampleInt╗
║     4     ║║     7     ║
╚═══════════╝╚═══════════╝`[1:]),
		},
		{
			name: "size difference first box",
			args: args{
				box1: asciiBoxForTest(`
╔═exampleInt╗
║     4     ║
║     4     ║
╚═══════════╝`[1:]),
				box2: asciiBoxForTest(`
╔═exampleInt╗
║     7     ║
╚═══════════╝`[1:]),
			},
			want: asciiBoxForTest(`
╔═exampleInt╗╔═exampleInt╗
║     4     ║║     7     ║
║     4     ║╚═══════════╝
╚═══════════╝             `[1:]),
		},
		{
			name: "size difference second box",
			args: args{
				box1: asciiBoxForTest(`
╔═exampleInt╗
║     4     ║
╚═══════════╝`[1:]),
				box2: asciiBoxForTest(`
╔═exampleInt╗
║     7     ║
║     7     ║
╚═══════════╝`[1:]),
			},
			want: asciiBoxForTest(`
╔═exampleInt╗╔═exampleInt╗
║     4     ║║     7     ║
╚═══════════╝║     7     ║
             ╚═══════════╝`[1:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsciiBoxWriterDefault.BoxSideBySide(tt.args.box1, tt.args.box2); !assert.Equal(t, tt.want, got) {
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
				box1: asciiBoxForTest(`
000 31  32  33  34  35  36  37  38  '12345678'
008 39  30  61  62  63  64  65  66  '90abcdef'
016 67  68  69  6a  6b  6c  6d  6e  'ghijklmn'
024 6f  70  71  72  73  74  75  76  'opqrstuv'
032 77  78  79  7a                  'wxyz    '`[1:]),
				box2: asciiBoxForTest(`
╔═super nice data══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║  000 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  '1234567890abcdefghijklmn'  ║
║  024 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  'opqrstuvwxyz.1234567890a'  ║
║  048 62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  'bcdefghijklmnopqrstuvwxy'  ║
║  072 7a  d3  61  61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  'z.aa1234567890abcdefghij'  ║
║  096 6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  'klmnopqrstuvwxyz.1234567'  ║
║  120 38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  '890abcdefghijklmnopqrstu'  ║
║  144 76  77  78  79  7a  d3  61  61  62                                                              'vwxyz.aab               '  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝`[1:]),
			},
			want: asciiBoxForTest(`
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
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝`[1:]),
		},
		{
			name: "different sized boxes",
			args: args{
				box1: asciiBoxForTest(`
╔═sampleField════════════╗
║123123123123123123123123║
╚════════════════════════╝`[1:]),
				box2: asciiBoxForTest(`
╔═sampleField╗
║123123123123║
╚════════════╝`[1:]),
			},
			want: asciiBoxForTest(`
╔═sampleField════════════╗
║123123123123123123123123║
╚════════════════════════╝
╔═sampleField╗            
║123123123123║            
╚════════════╝            `[1:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsciiBoxWriterDefault.BoxBelowBox(tt.args.box1, tt.args.box2); !assert.Equal(t, tt.want, got) {
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
			want: asciiBoxForTest(`
╔═sampleField╗
║123123123123║
╚════════════╝`[1:]),
		},
		{
			name: "simplebox-unamed",
			args: args{
				name:      "",
				data:      "123123123123",
				charWidth: 1,
			},
			want: asciiBoxForTest(`
╔════════════╗
║123123123123║
╚════════════╝`[1:]),
		},
		{
			name: "simplebox 2",
			args: args{
				name:      "sampleField",
				data:      "123123123123\n123123123123123123123123",
				charWidth: 1,
			},
			want: asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:]),
		},
		{
			name: "simplebox with too long name",
			args: args{
				name:      "sampleFieldsampleFieldsampleFieldsampleField",
				data:      "123123123123\n123123123123123123123123",
				charWidth: 1,
			},
			want: asciiBoxForTest(`
╔═sampleFieldsampleFieldsampleFieldsampleField╗
║                123123123123                 ║
║          123123123123123123123123           ║
╚═════════════════════════════════════════════╝`[1:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsciiBoxWriterDefault.BoxString(tt.args.name, tt.args.data, tt.args.charWidth); !assert.Equal(t, tt.want, got) {
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
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:]),
				},
				desiredWith: 1000,
			},
			want: asciiBoxForTest(`
╔═sampleField════════════╗╔═sampleField════════════╗
║      123123123123      ║║      123123123123      ║
║123123ABABABABABAB123123║║123123123123123123123123║
╚════════════════════════╝╚════════════════════════╝`[1:]),
		},
		{
			name: "not enough space",
			args: args{
				boxes: []AsciiBox{
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:]),
				},
				desiredWith: 0,
			},
			want: asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:]),
		},
		{
			name: "not enough space should result in multiple rows",
			args: args{
				boxes: []AsciiBox{
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:]),
				},
				desiredWith: 65,
			},
			want: asciiBoxForTest(`
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
╚════════════════════════╝╚════════════════════════╝`[1:]),
		},
		{
			name: "not enough space should result in multiple rows (3 columns)",
			args: args{
				boxes: []AsciiBox{
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝`[1:]),
					asciiBoxForTest(`
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:]),
				},
				desiredWith: 78,
			},
			want: asciiBoxForTest(`
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
╚════════════════════════╝╚════════════════════════╝                          `[1:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsciiBoxWriterDefault.AlignBoxes(tt.args.boxes, tt.args.desiredWith); !assert.Equal(t, tt.want, got) {
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
			m: asciiBoxForTest(`
123123123123123
123123123123123
123123123123123`[1:]),
			want: 15,
		},
		{
			name: "different width",
			m: asciiBoxForTest(`
123123123123123
123123123123123123123123123123
123123123123123`[1:]),
			want: 30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Width(); !assert.Equal(t, tt.want, got) {
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
					asciiBoxForTest(`
123123123
123123123
123123123`[1:]),
					asciiBoxForTest(`
abcabcabc
abcabcabc
abcabcabc`[1:]),
					asciiBoxForTest(`
zxyzxyzxy
zxyzxyzxy
zxyzxyzxy`[1:]),
				},
			},
			want: asciiBoxForTest(`
123123123abcabcabczxyzxyzxy
123123123abcabcabczxyzxyzxy
123123123abcabcabczxyzxyzxy`[1:]),
		},
		{
			name: "3 different",
			args: args{
				boxes: []AsciiBox{
					asciiBoxForTest(`
123123123
123123123
123123123`[1:]),
					asciiBoxForTest(`
abcabcabc
abcabcabcabcabcabcabcabcabc
abcabcabc`[1:]),
					asciiBoxForTest(`
zxyzxyzxy
zxyzxyzxy
zxyzxyzxy`[1:]),
				},
			},
			want: asciiBoxForTest(`
123123123abcabcabc                  zxyzxyzxy
123123123abcabcabcabcabcabcabcabcabczxyzxyzxy
123123123abcabcabc                  zxyzxyzxy`[1:]),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsciiBoxWriterDefault.(*asciiBoxWriter).mergeHorizontal(tt.args.boxes); !assert.Equal(t, tt.want, got) {
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
				box: asciiBoxForTest(`
123123123
123123123
123123123`[1:]),
				width: 100,
			},
			want: asciiBoxForTest(`
123123123                                                                                           
123123123                                                                                           
123123123                                                                                           `[1:]),
		},
		{
			name: "Big expand",
			args: args{
				box: asciiBoxForTest(`
123123123
123123123
123123123`[1:]),
				width: 10000,
			},
			want: asciiBoxForTest(fmt.Sprintf(`
123123123%[1]s
123123123%[1]s
123123123%[1]s`[1:], strings.Repeat(" ", 10000-9))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsciiBoxWriterDefault.(*asciiBoxWriter).expandBox(tt.args.box, tt.args.width); !assert.Equal(t, tt.want, got) {
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

func asciiBoxForTest(value string) AsciiBox {
	return AsciiBox{value, AsciiBoxWriterDefault.(*asciiBoxWriter), AsciiBoxWriterDefault.(*asciiBoxWriter).compressBoxSet()}
}

func TestNewAsciiBoxWriter(t *testing.T) {
	upperLeftCorner := "╔"
	upperRightCorner := "╗"
	horizontalLine := "═"
	verticalLine := "║"
	lowerLeftCorner := "╚"
	lowerRightCorner := "╝"
	tests := []struct {
		name string
		want AsciiBoxWriter
	}{
		{
			name: "create one",
			want: &asciiBoxWriter{
				boxSet: boxSet{
					upperLeftCorner:  upperLeftCorner,
					upperRightCorner: upperRightCorner,
					horizontalLine:   horizontalLine,
					verticalLine:     verticalLine,
					lowerLeftCorner:  lowerLeftCorner,
					lowerRightCorner: lowerRightCorner,
				},
				newLine:      '\n',
				emptyPadding: " ",
				// the name gets prefixed with a extra symbol for indent
				extraNameCharIndent: 1,
				borderWidth:         1,
				newLineCharWidth:    1,
				boxNameRegex:        regexp.MustCompile(`^` + upperLeftCorner + horizontalLine + `(?P<name>[\w /]+)` + horizontalLine + `*` + upperRightCorner),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewAsciiBoxWriter(), "NewAsciiBoxWriter()")
		})
	}
}

func TestNewAsciiBoxWriterWithCustomBorders(t *testing.T) {
	upperLeftCorner := "p"
	upperRightCorner := "l"
	horizontalLine := "c"
	verticalLine := "4"
	lowerLeftCorner := "x"
	lowerRightCorner := "!"
	type args struct {
		upperLeftCorner  string
		upperRightCorner string
		horizontalLine   string
		verticalLine     string
		lowerLeftCorner  string
		lowerRightCorner string
	}
	tests := []struct {
		name string
		args args
		want AsciiBoxWriter
	}{
		{
			name: "create one",
			args: args{
				upperLeftCorner:  upperLeftCorner,
				upperRightCorner: upperRightCorner,
				horizontalLine:   horizontalLine,
				verticalLine:     verticalLine,
				lowerLeftCorner:  lowerLeftCorner,
				lowerRightCorner: lowerRightCorner,
			},
			want: &asciiBoxWriter{
				boxSet: boxSet{
					upperLeftCorner:  upperLeftCorner,
					upperRightCorner: upperRightCorner,
					horizontalLine:   horizontalLine,
					verticalLine:     verticalLine,
					lowerLeftCorner:  lowerLeftCorner,
					lowerRightCorner: lowerRightCorner,
				},
				newLine:      '\n',
				emptyPadding: " ",
				// the name gets prefixed with a extra symbol for indent
				extraNameCharIndent: 1,
				borderWidth:         1,
				newLineCharWidth:    1,
				boxNameRegex:        regexp.MustCompile(`^` + upperLeftCorner + horizontalLine + `(?P<name>[\w /]+)` + horizontalLine + `*` + upperRightCorner),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewAsciiBoxWriterWithCustomBorders(tt.args.upperLeftCorner, tt.args.upperRightCorner, tt.args.horizontalLine, tt.args.verticalLine, tt.args.lowerLeftCorner, tt.args.lowerRightCorner), "NewAsciiBoxWriterWithCustomBorders(%v, %v, %v, %v, %v, %v)", tt.args.upperLeftCorner, tt.args.upperRightCorner, tt.args.horizontalLine, tt.args.verticalLine, tt.args.lowerLeftCorner, tt.args.lowerRightCorner)
		})
	}
}

func Test_asciiBoxWriter_AlignBoxes(t *testing.T) {
	type fields struct {
		boxSet              boxSet
		newLine             rune
		emptyPadding        string
		extraNameCharIndent int
		borderWidth         int
		newLineCharWidth    int
		boxNameRegex        *regexp.Regexp
	}
	type args struct {
		boxes        []AsciiBox
		desiredWidth int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   AsciiBox
	}{
		{
			name: "no boxes",
			want: AsciiBox{
				asciiBoxWriter: &asciiBoxWriter{},
			},
		},
		{
			name: "enough space",
			args: args{
				boxes: []AsciiBox{
					{
						data: `
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝`[1:],
						asciiBoxWriter: &asciiBoxWriter{},
					},
					{
						data: `
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:],
						asciiBoxWriter: &asciiBoxWriter{},
					},
				},
				desiredWidth: 1000,
			},
			want: AsciiBox{
				data: `
╔═sampleField════════════╗╔═sampleField════════════╗
║      123123123123      ║║      123123123123      ║
║123123ABABABABABAB123123║║123123123123123123123123║
╚════════════════════════╝╚════════════════════════╝`[1:],
				asciiBoxWriter: &asciiBoxWriter{},
			},
		},
		{
			name: "not enough space",
			fields: fields{
				boxSet: boxSet{
					upperLeftCorner:  "p",
					upperRightCorner: "l",
					horizontalLine:   "c",
					verticalLine:     "4",
					lowerLeftCorner:  "x",
					lowerRightCorner: "!",
				},
			},
			args: args{
				boxes: []AsciiBox{
					{
						data: `
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝`[1:],
						asciiBoxWriter: &asciiBoxWriter{
							boxSet: boxSet{
								upperLeftCorner:  "p",
								upperRightCorner: "l",
								horizontalLine:   "c",
								verticalLine:     "4",
								lowerLeftCorner:  "x",
								lowerRightCorner: "!",
							},
						},
					},
					{
						data: `
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:],
						asciiBoxWriter: &asciiBoxWriter{
							boxSet: boxSet{
								upperLeftCorner:  "p",
								upperRightCorner: "l",
								horizontalLine:   "c",
								verticalLine:     "4",
								lowerLeftCorner:  "x",
								lowerRightCorner: "!",
							},
						},
					},
				},
				desiredWidth: 0,
			},
			want: AsciiBox{
				data: `
╔═sampleField════════════╗
║      123123123123      ║
║123123ABABABABABAB123123║
╚════════════════════════╝
╔═sampleField════════════╗
║      123123123123      ║
║123123123123123123123123║
╚════════════════════════╝`[1:],
				asciiBoxWriter: &asciiBoxWriter{
					boxSet: boxSet{
						upperLeftCorner:  "p",
						upperRightCorner: "l",
						horizontalLine:   "c",
						verticalLine:     "4",
						lowerLeftCorner:  "x",
						lowerRightCorner: "!",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &asciiBoxWriter{
				boxSet:              tt.fields.boxSet,
				newLine:             tt.fields.newLine,
				emptyPadding:        tt.fields.emptyPadding,
				extraNameCharIndent: tt.fields.extraNameCharIndent,
				borderWidth:         tt.fields.borderWidth,
				newLineCharWidth:    tt.fields.newLineCharWidth,
				boxNameRegex:        tt.fields.boxNameRegex,
			}
			assert.Equalf(t, tt.want, a.AlignBoxes(tt.args.boxes, tt.args.desiredWidth), "AlignBoxes(%v, %v)", tt.args.boxes, tt.args.desiredWidth)
		})
	}
}

func Test_asciiBoxWriter_BoxBelowBox(t *testing.T) {
	type fields struct {
		boxSet              boxSet
		newLine             rune
		emptyPadding        string
		extraNameCharIndent int
		borderWidth         int
		newLineCharWidth    int
		boxNameRegex        *regexp.Regexp
	}
	type args struct {
		box1 AsciiBox
		box2 AsciiBox
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   AsciiBox
	}{
		{
			name: "put it below",
			want: AsciiBox{
				data: "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &asciiBoxWriter{
				boxSet:              tt.fields.boxSet,
				newLine:             tt.fields.newLine,
				emptyPadding:        tt.fields.emptyPadding,
				extraNameCharIndent: tt.fields.extraNameCharIndent,
				borderWidth:         tt.fields.borderWidth,
				newLineCharWidth:    tt.fields.newLineCharWidth,
				boxNameRegex:        tt.fields.boxNameRegex,
			}
			tt.want.asciiBoxWriter = a // hard to do that above otherwise
			assert.Equalf(t, tt.want, a.BoxBelowBox(tt.args.box1, tt.args.box2), "BoxBelowBox(%v, %v)", tt.args.box1, tt.args.box2)
		})
	}
}

func Test_asciiBoxWriter_BoxBox(t *testing.T) {
	type fields struct {
		boxSet              boxSet
		newLine             rune
		emptyPadding        string
		extraNameCharIndent int
		borderWidth         int
		newLineCharWidth    int
		boxNameRegex        *regexp.Regexp
	}
	type args struct {
		name      string
		box       AsciiBox
		charWidth int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   AsciiBox
	}{
		{
			name: "box a box",
			want: AsciiBox{
				data: "\x00\x00",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &asciiBoxWriter{
				boxSet:              tt.fields.boxSet,
				newLine:             tt.fields.newLine,
				emptyPadding:        tt.fields.emptyPadding,
				extraNameCharIndent: tt.fields.extraNameCharIndent,
				borderWidth:         tt.fields.borderWidth,
				newLineCharWidth:    tt.fields.newLineCharWidth,
				boxNameRegex:        tt.fields.boxNameRegex,
			}
			tt.want.asciiBoxWriter = a
			assert.Equalf(t, tt.want, a.BoxBox(tt.args.name, tt.args.box, tt.args.charWidth), "BoxBox(%v, %v, %v)", tt.args.name, tt.args.box, tt.args.charWidth)
		})
	}
}

func Test_asciiBoxWriter_BoxSideBySide(t *testing.T) {
	type fields struct {
		boxSet              boxSet
		newLine             rune
		emptyPadding        string
		extraNameCharIndent int
		borderWidth         int
		newLineCharWidth    int
		boxNameRegex        *regexp.Regexp
	}
	type args struct {
		box1 AsciiBox
		box2 AsciiBox
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   AsciiBox
	}{
		{
			name: "side by side",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &asciiBoxWriter{
				boxSet:              tt.fields.boxSet,
				newLine:             tt.fields.newLine,
				emptyPadding:        tt.fields.emptyPadding,
				extraNameCharIndent: tt.fields.extraNameCharIndent,
				borderWidth:         tt.fields.borderWidth,
				newLineCharWidth:    tt.fields.newLineCharWidth,
				boxNameRegex:        tt.fields.boxNameRegex,
			}
			tt.want.asciiBoxWriter = a
			assert.Equalf(t, tt.want, a.BoxSideBySide(tt.args.box1, tt.args.box2), "BoxSideBySide(%v, %v)", tt.args.box1, tt.args.box2)
		})
	}
}

func Test_asciiBoxWriter_BoxString(t *testing.T) {
	type fields struct {
		boxSet              boxSet
		newLine             rune
		emptyPadding        string
		extraNameCharIndent int
		borderWidth         int
		newLineCharWidth    int
		boxNameRegex        *regexp.Regexp
	}
	type args struct {
		name      string
		data      string
		charWidth int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   AsciiBox
	}{
		{
			name: "box a string",
			want: AsciiBox{
				data: "\x00\x00",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &asciiBoxWriter{
				boxSet:              tt.fields.boxSet,
				newLine:             tt.fields.newLine,
				emptyPadding:        tt.fields.emptyPadding,
				extraNameCharIndent: tt.fields.extraNameCharIndent,
				borderWidth:         tt.fields.borderWidth,
				newLineCharWidth:    tt.fields.newLineCharWidth,
				boxNameRegex:        tt.fields.boxNameRegex,
			}
			tt.want.asciiBoxWriter = a
			assert.Equalf(t, tt.want, a.BoxString(tt.args.name, tt.args.data, tt.args.charWidth), "BoxString(%v, %v, %v)", tt.args.name, tt.args.data, tt.args.charWidth)
		})
	}
}

func Test_asciiBoxWriter_boxString(t *testing.T) {
	type fields struct {
		boxSet              boxSet
		newLine             rune
		emptyPadding        string
		extraNameCharIndent int
		borderWidth         int
		newLineCharWidth    int
		boxNameRegex        *regexp.Regexp
	}
	type args struct {
		name      string
		data      string
		charWidth int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   AsciiBox
	}{
		{
			name: "box a string",
			want: AsciiBox{data: "\x00\x00"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &asciiBoxWriter{
				boxSet:              tt.fields.boxSet,
				newLine:             tt.fields.newLine,
				emptyPadding:        tt.fields.emptyPadding,
				extraNameCharIndent: tt.fields.extraNameCharIndent,
				borderWidth:         tt.fields.borderWidth,
				newLineCharWidth:    tt.fields.newLineCharWidth,
				boxNameRegex:        tt.fields.boxNameRegex,
			}
			tt.want.asciiBoxWriter = a
			assert.Equalf(t, tt.want, a.boxString(tt.args.name, tt.args.data, tt.args.charWidth), "boxString(%v, %v, %v)", tt.args.name, tt.args.data, tt.args.charWidth)
		})
	}
}

func Test_asciiBoxWriter_changeBoxName(t *testing.T) {
	type fields struct {
		boxSet              boxSet
		newLine             rune
		emptyPadding        string
		extraNameCharIndent int
		borderWidth         int
		newLineCharWidth    int
		boxNameRegex        *regexp.Regexp
	}
	type args struct {
		box     AsciiBox
		newName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   AsciiBox
	}{
		{
			name: "change a name",
			want: AsciiBox{data: "\x00\x00"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &asciiBoxWriter{
				boxSet:              tt.fields.boxSet,
				newLine:             tt.fields.newLine,
				emptyPadding:        tt.fields.emptyPadding,
				extraNameCharIndent: tt.fields.extraNameCharIndent,
				borderWidth:         tt.fields.borderWidth,
				newLineCharWidth:    tt.fields.newLineCharWidth,
				boxNameRegex:        tt.fields.boxNameRegex,
			}
			tt.want.asciiBoxWriter = a
			assert.Equalf(t, tt.want, a.changeBoxName(tt.args.box, tt.args.newName), "changeBoxName(%v, %v)", tt.args.box, tt.args.newName)
		})
	}
}

func Test_asciiBoxWriter_expandBox(t *testing.T) {
	type fields struct {
		boxSet              boxSet
		newLine             rune
		emptyPadding        string
		extraNameCharIndent int
		borderWidth         int
		newLineCharWidth    int
		boxNameRegex        *regexp.Regexp
	}
	type args struct {
		box          AsciiBox
		desiredWidth int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   AsciiBox
	}{
		{
			name: "expand",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &asciiBoxWriter{
				boxSet:              tt.fields.boxSet,
				newLine:             tt.fields.newLine,
				emptyPadding:        tt.fields.emptyPadding,
				extraNameCharIndent: tt.fields.extraNameCharIndent,
				borderWidth:         tt.fields.borderWidth,
				newLineCharWidth:    tt.fields.newLineCharWidth,
				boxNameRegex:        tt.fields.boxNameRegex,
			}
			assert.Equalf(t, tt.want, a.expandBox(tt.args.box, tt.args.desiredWidth), "expandBox(%v, %v)", tt.args.box, tt.args.desiredWidth)
		})
	}
}

func Test_asciiBoxWriter_getBoxName(t *testing.T) {
	type fields struct {
		boxSet              boxSet
		newLine             rune
		emptyPadding        string
		extraNameCharIndent int
		borderWidth         int
		newLineCharWidth    int
		boxNameRegex        *regexp.Regexp
	}
	type args struct {
		box AsciiBox
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "get a name",
			fields: fields{
				boxNameRegex: regexp.MustCompile("(?P<name>[\\w /]+)"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &asciiBoxWriter{
				boxSet:              tt.fields.boxSet,
				newLine:             tt.fields.newLine,
				emptyPadding:        tt.fields.emptyPadding,
				extraNameCharIndent: tt.fields.extraNameCharIndent,
				borderWidth:         tt.fields.borderWidth,
				newLineCharWidth:    tt.fields.newLineCharWidth,
				boxNameRegex:        tt.fields.boxNameRegex,
			}
			assert.Equalf(t, tt.want, a.getBoxName(tt.args.box), "getBoxName(%v)", tt.args.box)
		})
	}
}

func Test_asciiBoxWriter_getBoxNameWrongRegex(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			assert.Equal(t, "should never occur as we only have one named group", err)
		} else {
			t.Error("We should catch a panic")
		}
	}()
	(&asciiBoxWriter{boxNameRegex: regexp.MustCompile(".*")}).getBoxName(AsciiBox{})
	t.Error("the code above should panic")
}

func Test_asciiBoxWriter_hasBorders(t *testing.T) {
	type fields struct {
		boxSet              boxSet
		newLine             rune
		emptyPadding        string
		extraNameCharIndent int
		borderWidth         int
		newLineCharWidth    int
		boxNameRegex        *regexp.Regexp
	}
	type args struct {
		box AsciiBox
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "has no borders",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &asciiBoxWriter{
				boxSet:              tt.fields.boxSet,
				newLine:             tt.fields.newLine,
				emptyPadding:        tt.fields.emptyPadding,
				extraNameCharIndent: tt.fields.extraNameCharIndent,
				borderWidth:         tt.fields.borderWidth,
				newLineCharWidth:    tt.fields.newLineCharWidth,
				boxNameRegex:        tt.fields.boxNameRegex,
			}
			assert.Equalf(t, tt.want, a.hasBorders(tt.args.box), "hasBorders(%v)", tt.args.box)
		})
	}
}

func Test_asciiBoxWriter_mergeHorizontal(t *testing.T) {
	type fields struct {
		boxSet              boxSet
		newLine             rune
		emptyPadding        string
		extraNameCharIndent int
		borderWidth         int
		newLineCharWidth    int
		boxNameRegex        *regexp.Regexp
	}
	type args struct {
		boxes []AsciiBox
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   AsciiBox
	}{
		{
			name: "merge it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &asciiBoxWriter{
				boxSet:              tt.fields.boxSet,
				newLine:             tt.fields.newLine,
				emptyPadding:        tt.fields.emptyPadding,
				extraNameCharIndent: tt.fields.extraNameCharIndent,
				borderWidth:         tt.fields.borderWidth,
				newLineCharWidth:    tt.fields.newLineCharWidth,
				boxNameRegex:        tt.fields.boxNameRegex,
			}
			tt.want.asciiBoxWriter = a
			assert.Equalf(t, tt.want, a.mergeHorizontal(tt.args.boxes), "mergeHorizontal(%v)", tt.args.boxes)
		})
	}
}

func Test_asciiBoxWriter_unwrap(t *testing.T) {
	type fields struct {
		boxSet              boxSet
		newLine             rune
		emptyPadding        string
		extraNameCharIndent int
		borderWidth         int
		newLineCharWidth    int
		boxNameRegex        *regexp.Regexp
	}
	type args struct {
		box AsciiBox
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   AsciiBox
	}{
		{
			name: "unwrap",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &asciiBoxWriter{
				boxSet:              tt.fields.boxSet,
				newLine:             tt.fields.newLine,
				emptyPadding:        tt.fields.emptyPadding,
				extraNameCharIndent: tt.fields.extraNameCharIndent,
				borderWidth:         tt.fields.borderWidth,
				newLineCharWidth:    tt.fields.newLineCharWidth,
				boxNameRegex:        tt.fields.boxNameRegex,
			}
			assert.Equalf(t, tt.want, a.unwrap(tt.args.box), "unwrap(%v)", tt.args.box)
		})
	}
}

func Test_boxSet_compressBoxSet(t *testing.T) {
	type fields struct {
		upperLeftCorner  string
		upperRightCorner string
		horizontalLine   string
		verticalLine     string
		lowerLeftCorner  string
		lowerRightCorner string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "compress",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := boxSet{
				upperLeftCorner:  tt.fields.upperLeftCorner,
				upperRightCorner: tt.fields.upperRightCorner,
				horizontalLine:   tt.fields.horizontalLine,
				verticalLine:     tt.fields.verticalLine,
				lowerLeftCorner:  tt.fields.lowerLeftCorner,
				lowerRightCorner: tt.fields.lowerRightCorner,
			}
			assert.Equalf(t, tt.want, b.compressBoxSet(), "compressBoxSet()")
		})
	}
}

func Test_boxSet_contributeToCompressedBoxSet(t *testing.T) {
	type fields struct {
		upperLeftCorner  string
		upperRightCorner string
		horizontalLine   string
		verticalLine     string
		lowerLeftCorner  string
		lowerRightCorner string
	}
	type args struct {
		box AsciiBox
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "contribute nothing",
		},
		{
			name: "they don't contribute",
			fields: fields{
				upperLeftCorner:  "p",
				upperRightCorner: "l",
				horizontalLine:   "c",
				verticalLine:     "4",
				lowerLeftCorner:  "x",
			},
			want: "plc4x",
		},
		{
			name: "I don't contribute",
			args: args{
				AsciiBox{
					compressedBoxSet: "plc4x",
				},
			},
			want: "plc4x",
		},
		{
			name: "We are equal",
			fields: fields{
				upperLeftCorner:  "p",
				upperRightCorner: "l",
				horizontalLine:   "c",
				verticalLine:     "4",
				lowerLeftCorner:  "x",
			},
			args: args{
				AsciiBox{
					compressedBoxSet: "plc4x",
				},
			},
			want: "plc4x",
		},
		{
			name: "We add up",
			fields: fields{
				upperLeftCorner:  "p",
				upperRightCorner: "l",
				horizontalLine:   "c",
				verticalLine:     "4",
				lowerLeftCorner:  "x",
				lowerRightCorner: "!",
			},
			args: args{
				AsciiBox{
					compressedBoxSet: "plc4x",
				},
			},
			want: "plc4x,plc4x!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := boxSet{
				upperLeftCorner:  tt.fields.upperLeftCorner,
				upperRightCorner: tt.fields.upperRightCorner,
				horizontalLine:   tt.fields.horizontalLine,
				verticalLine:     tt.fields.verticalLine,
				lowerLeftCorner:  tt.fields.lowerLeftCorner,
				lowerRightCorner: tt.fields.lowerRightCorner,
			}
			assert.Equalf(t, tt.want, b.contributeToCompressedBoxSet(tt.args.box), "contributeToCompressedBoxSet(%v)", tt.args.box)
		})
	}
}

func Test_combineCompressedBoxSets(t *testing.T) {
	type args struct {
		box1 AsciiBox
		box2 AsciiBox
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "combine",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, combineCompressedBoxSets(tt.args.box1, tt.args.box2), "combineCompressedBoxSets(%v, %v)", tt.args.box1, tt.args.box2)
		})
	}
}

func Test_countChars(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "count",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, countChars(tt.args.s), "countChars(%v)", tt.args.s)
		})
	}
}
