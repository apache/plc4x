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

package primitivedata

import (
	"io"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type TagList struct {
	*DefaultRFormatter
	tagList []Tag
}

func NewTagList(arg Arg) *TagList {
	t := &TagList{
		DefaultRFormatter: NewDefaultRFormatter(),
	}
	switch arg := arg.(type) {
	case []any:
		args := arg
		for _, a := range args {
			t.tagList = append(t.tagList, a.(Tag))
		}
	case []Tag:
		args := arg
		for _, a := range args {
			t.tagList = append(t.tagList, a.(Tag))
		}
	case Args:
		args := arg
		for _, a := range args {
			t.tagList = append(t.tagList, a.(Tag))
		}
	}
	return t
}

func (b *TagList) Append(tag Tag) {
	b.tagList = append(b.tagList, tag)
}

func (b *TagList) Extend(tags ...Tag) {
	for _, tag := range tags {
		b.tagList = append(b.tagList, tag)
	}
}

func (b *TagList) Peek() Tag {
	if len(b.tagList) < 1 {
		return nil
	}
	return b.tagList[0]
}

func (b *TagList) Push(tag Tag) {
	b.tagList = append([]Tag{tag}, b.tagList...)
}

func (b *TagList) Pop() Tag {
	if len(b.tagList) < 1 {
		return nil
	}
	item := b.tagList[0]
	b.tagList = b.tagList[1:]
	return item
}

// GetContext Return a tag or a list of tags context encoded.
func (b *TagList) GetContext(context uint) (any, error) {
	// forward pass
	i := 0
	for i < len(b.tagList) {
		tag := b.tagList[i]

		switch tag.GetTagClass() {
		case model.TagClass_APPLICATION_TAGS: // skip application stuff
		case model.TagClass_CONTEXT_SPECIFIC_TAGS: // check for context encoded atomic value
			if tag.GetTagNumber() == context {
				return tag, nil
			}
		case TagOpeningTagClass:
			keeper := tag.GetTagNumber() == context
			var rslt []Tag
			i += 1
			lvl := 0
		innerSearch:
			for i < len(b.tagList) {
				tag := b.tagList[i]
				switch tag.GetTagClass() {
				case TagOpeningTagClass:
					lvl += 1
				case TagClosingTagClass:
					lvl -= 1
					if lvl < 0 {
						break innerSearch
					}
				}

				rslt = append(rslt, tag)
				i += 1
			}

			// make sure everything balances
			if lvl >= 1 {
				return nil, errors.New("mismatched open/close tag")
			}

			// get everything we need
			if keeper {
				return NewTagList(rslt), nil
			}
		}
		i += 1
	}
	return nil, nil
}

func (b *TagList) Encode(data PDUData) {
	for _, tag := range b.tagList {
		tag.Encode(data)
	}
}

func (b *TagList) Decode(data PDUData) error {
	for len(data.GetPduData()) != 0 {
		var tag Tag
		tag, err := NewTag(NA(data))
		if err != nil {
			return errors.Wrap(err, "error creating tag")
		}
		switch tag.GetTagClass() {
		case model.TagClass_APPLICATION_TAGS:
		case model.TagClass_CONTEXT_SPECIFIC_TAGS:
			tag, err = NewContextTag(NA(tag.GetTagNumber(), tag.GetTagData()))
			if err != nil {
				panic(err)
			}
		case model.TagClass(TagOpeningTagClass):
			tag, err = NewOpeningTag(tag.GetTagNumber())
			if err != nil {
				panic(err)
			}
		case model.TagClass(TagClosingTagClass):
			tag, err = NewClosingTag(tag.GetTagNumber())
			if err != nil {
				panic(err)
			}
		}
		b.tagList = append(b.tagList, tag)
	}
	return nil
}

func (b *TagList) GetTagList() []Tag {
	return b.tagList
}

func (b *TagList) PrintDebugContents(indent int, file io.Writer, _ids []uintptr) {
	for _, tag := range b.tagList {
		tag.PrintDebugContents(indent, file, _ids)
	}
}
