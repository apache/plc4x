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

package model

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func ParseTiaTime(ctx context.Context, io utils.ReadBuffer) (uint32, error) {
	/*try {
	      int millisSinceMidnight = io.readInt(32);
	      return LocalTime.now().withHour(0).withMinute(0).withSecond(0).withNano(0).plus(
	      millisSinceMidnight, ChronoUnit.MILLIS);
	  } catch (ParseException e) {
	      return null;
	  }*/
	return 0, nil
}

func SerializeTiaTime(ctx context.Context, io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing TIME not implemented");
	return nil
}

func ParseS5Time(ctx context.Context, io utils.ReadBuffer) (uint32, error) {
	/*try {
	      int stuff = io.readInt(16);
	      // TODO: Implement this correctly.
	      throw new NotImplementedException("S5TIME not implemented");
	  } catch (ParseException e) {
	      return null;
	  }*/
	return 0, nil
}

func SerializeS5Time(ctx context.Context, io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing S5TIME not implemented");
	return nil
}

func ParseTiaLTime(ctx context.Context, io utils.ReadBuffer) (uint32, error) {
	//throw new NotImplementedException("LTIME not implemented");
	return 0, nil
}

func SerializeTiaLTime(ctx context.Context, io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing LTIME not implemented");
	return nil
}

func ParseTiaTimeOfDay(ctx context.Context, io utils.ReadBuffer) (time.Time, error) {
	/*try {
	      long millisSinceMidnight = io.readUnsignedLong(32);
	      return LocalTime.now().withHour(0).withMinute(0).withSecond(0).withNano(0).plus(
	          millisSinceMidnight, ChronoUnit.MILLIS);
	  } catch (ParseException e) {
	      return null;
	  }*/
	return time.Time{}, nil
}

func SerializeTiaTimeOfDay(ctx context.Context, io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing TIME_OF_DAY not implemented");
	return nil
}

func ParseTiaDate(ctx context.Context, io utils.ReadBuffer) (uint16, error) {
	/*try {
	      int daysSince1990 = io.readUnsignedInt(16);
	      return LocalDate.now().withYear(1990).withDayOfMonth(1).withMonth(1).plus(daysSince1990, ChronoUnit.DAYS);
	  } catch (ParseException e) {
	      return null;
	  }*/
	return 0, nil
}

func SerializeTiaDate(ctx context.Context, io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing DATE not implemented");
	return nil
}

func ParseTiaDateTime(ctx context.Context, io utils.ReadBuffer) (time.Time, error) {
	/*try {
	      int year = io.readUnsignedInt(16);
	      int month = io.readUnsignedInt(8);
	      int day = io.readUnsignedInt(8);
	      // Skip day-of-week
	      io.readByte(8);
	      int hour = io.readByte(8);
	      int minute = io.readByte(8);
	      int second = io.readByte(8);
	      int nanosecond = io.readUnsignedInt(24);

	      return LocalDateTime.of(year, month, day, hour, minute, second, nanosecond);
	  } catch (Exception e) {
	      return null;
	  }*/
	return time.Time{}, nil
}

func SerializeTiaDateTime(ctx context.Context, io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing DATE_AND_TIME not implemented");
	return nil
}

func parseTiaDate(ctx context.Context, io utils.ReadBuffer) (time.Time, error) {
	return time.Time{}, nil
}

func serializeTiaDate(ctx context.Context, io utils.WriteBuffer, value values.PlcValue) error {
	return nil
}

func ParseS7String(ctx context.Context, io utils.ReadBuffer, stringLength int32, encoding string) (string, error) {
	var multiplier int32
	switch encoding {
	case "UTF-8":
		multiplier = 8
	case "UTF-16":
		multiplier = 16
	}
	return io.ReadString("", uint32(stringLength*multiplier), utils.WithEncoding(encoding))
}

func SerializeS7String(ctx context.Context, io utils.WriteBuffer, value values.PlcValue, stringLength int32, encoding string) error {
	var multiplier int32
	switch encoding {
	case "UTF-8":
		multiplier = 8
	case "UTF-16":
		multiplier = 16
	}
	return io.WriteString("", uint32(stringLength*multiplier), value.GetString(), utils.WithEncoding(encoding))
}

func ParseS7Char(ctx context.Context, io utils.ReadBuffer, encoding string) (uint8, error) {
	return io.ReadUint8("", 8)
}

func SerializeS7Char(ctx context.Context, io utils.WriteBuffer, value values.PlcValue, encoding string) error {
	return io.WriteUint8("", 8, value.GetUint8())
}

func RightShift3(ctx context.Context, readBuffer utils.ReadBuffer, dataTransportSize DataTransportSize) (any, error) {
	return uint16(0), nil
}

func LeftShift3(ctx context.Context, writeBuffer utils.WriteBuffer, valueLength uint16) error {
	return nil
}

func EventItemLength(ctx context.Context, readBuffer utils.ReadBuffer, valueLength uint16) uint16 {
	return 0
}

func BcdToInt(ctx context.Context, readBuffer utils.ReadBuffer) (any, error) {
	return uint8(0), nil
}

func ByteToBcd(ctx context.Context, writeBuffer utils.WriteBuffer, value uint8) error {
	return nil
}

func S7msecToInt(ctx context.Context, readBuffer utils.ReadBuffer) (any, error) {
	return uint16(0), nil
}

func IntToS7msec(ctx context.Context, writeBuffer utils.WriteBuffer, value uint16) error {
	return nil
}

func ParseSiemensYear(_ context.Context, readBuffer utils.ReadBuffer) (uint16, error) {
	year, err := readBuffer.ReadUint16("year", 8, utils.WithEncoding("BCD"))
	if err != nil {
		return 0, errors.Wrap(err, "Error parsing year")
	}
	if year < 90 {
		return 2000 + year, nil
	} else {
		return 1900 + year, nil
	}
}

func SerializeSiemensYear(ctx context.Context, writeBuffer utils.WriteBuffer, dateTime values.PlcValue) error {
	year := dateTime.GetDateTime().Year()
	if year > 2000 {
		return writeBuffer.WriteUint16("year", 8, uint16(year-2000), utils.WithEncoding("BCD"))
	} else {
		return writeBuffer.WriteUint16("year", 8, uint16(year-1900), utils.WithEncoding("BCD"))
	}
}
