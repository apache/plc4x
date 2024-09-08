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
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type TimeTuple struct {
	Hour      int
	Minute    int
	Second    int
	Hundredth int
}

type Time struct {
	value TimeTuple

	_appTag model.BACnetDataType
}

func NewTime(args Args) (*Time, error) {
	d := &Time{
		_appTag: model.BACnetDataType_TIME,
	}
	var arg any
	if len(args) > 0 {
		arg = args[0]
	}
	hour := 255
	if len(args) > 1 {
		hour = args[1].(int)
	}
	d.value.Hour = hour
	minute := 0xff
	if len(args) > 2 {
		minute = args[2].(int)
	}
	d.value.Minute = minute
	second := 0xff
	if len(args) > 3 {
		second = args[3].(int)
	}
	d.value.Second = second
	hundredth := 0xff
	if len(args) > 4 {
		hundredth = args[4].(int)
	}
	d.value.Hundredth = hundredth

	if arg == nil {
		return d, nil
	}
	switch arg := arg.(type) {
	case *tag:
		err := d.Decode(arg)
		if err != nil {
			return nil, errors.Wrap(err, "error decoding")
		}
		return d, nil
	case TimeTuple:
		d.value = arg
	case string:
		// lower case everything
		arg = strings.ToLower(arg)
		timeRegex := regexp.MustCompile(`^([*]|[0-9]+)[:]([*]|[0-9]+)(?:[:]([*]|[0-9]+)(?:[.]([*]|[0-9]+))?)?$`)

		if !timeRegex.MatchString(arg) {
			return nil, errors.New("invalid time pattern")
		}
		// make a list of the contents from matching patterns
		match := timeRegex.FindStringSubmatch(arg)[1:]
		if len(match) == 0 {
			return nil, errors.New("unmatched")
		}

		var tupList []int
		for _, s := range match {
			if s == "*" {
				tupList = append(tupList, 255)
			} else if s == "" {
				if slices.Contains(match, "*") {
					tupList = append(tupList, 255)
				} else {
					tupList = append(tupList, 0)
				}
			} else {
				i, _ := strconv.Atoi(s)
				tupList = append(tupList, i)
			}
		}
		if tupList[3] != 0xff {
			tupList[3] *= 10
		}
		d.value = TimeTuple{tupList[0], tupList[1], tupList[2], tupList[3]}
	case time.Duration:
		d.value = TimeTuple{int(arg.Hours()), int(arg.Minutes()), int(arg.Seconds()), int(arg.Milliseconds() * 10)}
	case *Time:
		d.value = arg.value
	case float32:
		d.now(arg)
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return d, nil
}

func (t *Time) GetAppTag() model.BACnetDataType {
	return t._appTag
}

func (t *Time) now(arg float32) {
	panic("implement me") // TODO
}

func (t *Time) Encode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	tag.setAppData(uint(t._appTag), []byte{byte(t.value.Hour), byte(t.value.Minute), byte(t.value.Second), byte(t.value.Hundredth)})
	return nil
}

func (t *Time) Decode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(t._appTag) {
		return errors.New("Time application tag required")
	}
	if len(tag.GetTagData()) != 4 {
		return errors.New("invalid tag length")
	}

	tagData := tag.GetTagData()

	t.value = TimeTuple{int(tagData[0]), int(tagData[1]), int(tagData[2]), int(tagData[3])}
	return nil
}

func (t *Time) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (t *Time) Compare(other any) int {
	switch other := other.(type) {
	case *Time:
		_ = other                // TODO: implement
		panic("not implemented") // TODO: implement me
		return -1
	default:
		return -1
	}
}

func (t *Time) LowerThan(other any) bool {
	switch other := other.(type) {
	case *Time:
		_ = other                // TODO: implement
		panic("not implemented") // TODO: implement me
		return false
	default:
		return false
	}
}

func (t *Time) Equals(other any) bool {
	return t.value == other
}

func (t *Time) GetValue() TimeTuple {
	return t.value
}

func (t *Time) Coerce(arg Time) TimeTuple {
	return arg.GetValue()
}

func (t *Time) String() string {
	// rip it apart
	hour, minute, second, hundredth := t.value.Hour, t.value.Minute, t.value.Second, t.value.Hundredth

	rslt := "Time("
	if hour == 255 {
		rslt += "*:"
	} else {
		rslt += fmt.Sprintf("%02d:", hour)
	}
	if minute == 255 {
		rslt += "*:"
	} else {
		rslt += fmt.Sprintf("%02d:", minute)
	}
	if second == 255 {
		rslt += "*."
	} else {
		rslt += fmt.Sprintf("%02d.", second)
	}
	if hundredth == 255 {
		rslt += "*)"
	} else {
		rslt += fmt.Sprintf("%02d)", hundredth)
	}
	return rslt
}
