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
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

const _mm = `(?P<month>0?[1-9]|1[0-4]|odd|even|255|[*])`

const _dd = `(?P<day>[0-3]?\d|last|odd|even|255|[*])`

const _yy = `(?P<year>\d{2}|255|[*])`

const _yyyy = `(?P<year>\d{4}|255|[*])`

const _dow = `(?P<dow>[1-7]|mon|tue|wed|thu|fri|sat|sun|255|[*])`

var _special_mon = map[string]int{"*": 255, "odd": 13, "even": 14, "": 255}

var _special_mon_inv = map[int]string{255: "*", 13: "odd", 14: "even"}

var _special_day = map[string]int{"*": 255, "last": 32, "odd": 33, "even": 34, "": 255}

var _special_day_inv = map[int]string{255: "*", 32: "last", 33: "odd", 34: "even"}

var _special_dow = map[string]int{"*": 255, "mon": 1, "tue": 2, "wed": 3, "thu": 4, "fri": 5, "sat": 6, "sun": 7}

var _special_dow_inv = map[int]string{255: "*", 1: "mon", 2: "tue", 3: "wed", 4: "thu", 5: "fri", 6: "sat", 7: "sun"}

// Create a composite pattern and compile it.
func _merge(args ...string) *regexp.Regexp {
	return regexp.MustCompile(`^` + strings.Join(args, `[/-]`) + `(?:\s+` + _dow + `)?$`)
}

// make a list of compiled patterns
var _date_patterns = []*regexp.Regexp{
	_merge(_yyyy, _mm, _dd),
	_merge(_mm, _dd, _yyyy),
	_merge(_dd, _mm, _yyyy),
	_merge(_yy, _mm, _dd),
	_merge(_mm, _dd, _yy),
	_merge(_dd, _mm, _yy),
}

type DateTuple struct {
	Year      int
	Month     int
	Day       int
	DayOfWeek int
}

type Date struct {
	value DateTuple

	_appTag model.BACnetDataType
}

func NewDate(args Args) (*Date, error) {
	d := &Date{
		_appTag: model.BACnetDataType_DATE,
	}
	var arg any
	if len(args) > 0 {
		arg = args[0]
	}
	year := 255
	if len(args) > 1 {
		year = args[1].(int)
	}
	if year >= 1900 {
		year = year - 1900
	}
	d.value.Year = year
	month := 0xff
	if len(args) > 2 {
		month = args[2].(int)
	}
	d.value.Month = month
	day := 0xff
	if len(args) > 3 {
		day = args[3].(int)
	}
	d.value.Day = day
	dayOfWeek := 0xff
	if len(args) > 4 {
		dayOfWeek = args[4].(int)
	}
	d.value.DayOfWeek = dayOfWeek

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
	case DateTuple:
		d.value = arg
	case string:
		// lower case everything
		arg = strings.ToLower(arg)

		// make a list of the contents from matching patterns
		matches := [][]string{}
		for _, p := range _date_patterns {
			if p.MatchString(arg) {
				groups := CombinedPattern.FindStringSubmatch(arg)
				matches = append(matches, groups[1:])
			}
		}
		if len(matches) == 0 {
			return nil, errors.New("unmatched")
		}

		var match []string
		if len(matches) == 1 {
			match = matches[0]
		} else {
			// check to see if they really are the same
			panic("what to do here")
		}

		// extract the year and normalize
		matchedYear := match[0]
		if matchedYear == "*" || matchedYear == "" {
			year = 0xff
		} else {
			yearParse, err := strconv.ParseInt(matchedYear, 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "error parsing year")
			}
			year = int(yearParse)
			if year == 0xff {
				return d, nil
			}
			if year < 35 {
				year += 2000
			} else if year < 100 {
				year += 1900
			} else if year < 1900 {
				return nil, errors.New("invalid year")
			}
		}

		// extract the month and normalize
		matchedmonth := match[0]
		if specialMonth, ok := _special_mon[matchedmonth]; ok {
			month = specialMonth
		} else {
			monthParse, err := strconv.ParseInt(matchedmonth, 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "error parsing month")
			}
			month = int(monthParse)
			if month == 0xff {
				return d, nil
			}
			if month == 0 || month > 14 {
				return nil, errors.New("invalid month")
			}
		}

		// extract the day and normalize
		matchedday := match[0]
		if specialday, ok := _special_day[matchedday]; ok {
			day = specialday
		} else {
			dayParse, err := strconv.ParseInt(matchedday, 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "error parsing day")
			}
			day = int(dayParse)
			if day == 0xff {
				return d, nil
			}
			if day == 0 || day > 34 {
				return nil, errors.New("invalid day")
			}
		}

		// extract the dayOfWeek and normalize
		matcheddayOfWeek := match[0]
		if specialdayOfWeek, ok := _special_dow[matcheddayOfWeek]; ok {
			dayOfWeek = specialdayOfWeek
		} else if matcheddayOfWeek == "" {
			return d, nil
		} else {
			dayOfWeekParse, err := strconv.ParseInt(matcheddayOfWeek, 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "error parsing dayOfWeek")
			}
			dayOfWeek = int(dayOfWeekParse)
			if dayOfWeek == 0xff {
				return d, nil
			}
			if dayOfWeek > 7 {
				return nil, errors.New("invalid dayOfWeek")
			}
		}

		// year becomes the correct octet
		if year != 0xff {
			year -= 1900
		}

		// save the value
		d.value.Year = year
		d.value.Month = month
		d.value.Day = day
		d.value.DayOfWeek = dayOfWeek

		// calculate the day of the week
		if dayOfWeek == 0 {
			d.calcDayOfWeek()
		}
	case *Date:
		d.value = arg.value
	case float32:
		d.now(arg)
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", arg)
	}

	return d, nil
}

func (d *Date) GetAppTag() model.BACnetDataType {
	return d._appTag
}

func (d *Date) calcDayOfWeek() {
	year, month, day, dayOfWeek := d.value.Year, d.value.Month, d.value.Day, d.value.DayOfWeek

	// assume the worst
	dayOfWeek = 255

	// check for special values
	if year == 255 {
		return
	} else if _, ok := _special_mon_inv[month]; ok {
		return
	} else if _, ok := _special_day_inv[month]; ok {
		return
	} else {
		var today time.Time
		today = time.Date(year+1900, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		today.Add(24 * time.Hour)
		dayOfWeek = int(today.Weekday())
	}

	// put it back together
	d.value.Year = year
	d.value.Month = month
	d.value.Day = day
	d.value.DayOfWeek = dayOfWeek
}

func (d *Date) now(arg float32) {
	panic("implement me") // TODO
}

func (d *Date) Encode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	tag.setAppData(uint(d._appTag), []byte{byte(d.value.Year), byte(d.value.Month), byte(d.value.Day), byte(d.value.DayOfWeek)})

	return nil
}

func (d *Date) Decode(arg Arg) error {
	tag, ok := arg.(Tag)
	if !ok {
		return errors.Errorf("%T is not a Tag", arg)
	}
	if tag.GetTagClass() != model.TagClass_APPLICATION_TAGS || tag.GetTagNumber() != uint(d._appTag) {
		return errors.New("Date application tag required")
	}
	if len(tag.GetTagData()) != 4 {
		return errors.New("invalid tag length")
	}

	td := tag.GetTagData()
	year, month, day, dayOfWeek := td[0], td[1], td[2], td[3]
	d.value.Year, d.value.Month, d.value.Day, d.value.DayOfWeek = int(year), int(month), int(day), int(dayOfWeek)
	return nil
}

func (d *Date) IsValid(arg any) bool {
	_, ok := arg.(bool)
	return ok
}

func (d *Date) Compare(other any) int {
	switch other := other.(type) {
	case *Date:
		_ = other                //TODO: implement me
		panic("not implemented") // TODO: implement me
		return -1
	default:
		return -1
	}
}

func (d *Date) LowerThan(other any) bool {
	switch other := other.(type) {
	case *Date:
		// return d.getLong() < other.getLong()
		_ = other                // TODO: implement me
		panic("not implemented") // TODO: implement me
		return false
	default:
		return false
	}
}

func (d *Date) Equals(other any) bool {
	return d.value == other
}

func (d *Date) GetValue() DateTuple {
	return d.value
}

func (d *Date) Coerce(arg Date) DateTuple {
	return arg.GetValue()
}

func (d *Date) String() string {
	year, month, day, dayOfWeek := d.value.Year, d.value.Month, d.value.Day, d.value.DayOfWeek
	yearStr := "*"
	if year != 255 {
		yearStr = strconv.Itoa(year + 1900)
	}
	monthStr := strconv.Itoa(month)
	if ms, ok := _special_mon_inv[month]; ok {
		monthStr = ms
	}
	dayStr := strconv.Itoa(day)
	if ms, ok := _special_day_inv[day]; ok {
		dayStr = ms
	}
	dowStr := strconv.Itoa(dayOfWeek)
	if ms, ok := _special_dow_inv[dayOfWeek]; ok {
		dowStr = ms
	}

	return fmt.Sprintf("Date(%s-%s-%s %s)", yearStr, monthStr, dayStr, dowStr)
}
