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
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func ParseTiaTime(io utils.ReadBuffer) (uint32, error) {
	/*try {
	      int millisSinceMidnight = io.readInt(32);
	      return LocalTime.now().withHour(0).withMinute(0).withSecond(0).withNano(0).plus(
	      millisSinceMidnight, ChronoUnit.MILLIS);
	  } catch (ParseException e) {
	      return null;
	  }*/
	return 0, nil
}

func SerializeTiaTime(io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing TIME not implemented");
	return nil
}

func ParseS5Time(io utils.ReadBuffer) (uint32, error) {
	/*try {
	      int stuff = io.readInt(16);
	      // TODO: Implement this correctly.
	      throw new NotImplementedException("S5TIME not implemented");
	  } catch (ParseException e) {
	      return null;
	  }*/
	return 0, nil
}

func SerializeS5Time(io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing S5TIME not implemented");
	return nil
}

func ParseTiaLTime(io utils.ReadBuffer) (uint32, error) {
	//throw new NotImplementedException("LTIME not implemented");
	return 0, nil
}

func SerializeTiaLTime(io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing LTIME not implemented");
	return nil
}

func ParseTiaTimeOfDay(io utils.ReadBuffer) (time.Time, error) {
	/*try {
	      long millisSinceMidnight = io.readUnsignedLong(32);
	      return LocalTime.now().withHour(0).withMinute(0).withSecond(0).withNano(0).plus(
	          millisSinceMidnight, ChronoUnit.MILLIS);
	  } catch (ParseException e) {
	      return null;
	  }*/
	return time.Time{}, nil
}

func SerializeTiaTimeOfDay(io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing TIME_OF_DAY not implemented");
	return nil
}

func ParseTiaDate(io utils.ReadBuffer) (time.Time, error) {
	/*try {
	      int daysSince1990 = io.readUnsignedInt(16);
	      return LocalDate.now().withYear(1990).withDayOfMonth(1).withMonth(1).plus(daysSince1990, ChronoUnit.DAYS);
	  } catch (ParseException e) {
	      return null;
	  }*/
	return time.Time{}, nil
}

func SerializeTiaDate(io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing DATE not implemented");
	return nil
}

func ParseTiaDateTime(io utils.ReadBuffer) (time.Time, error) {
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

func SerializeTiaDateTime(io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing DATE_AND_TIME not implemented");
	return nil
}

func ParseS7String(io utils.ReadBuffer, stringLength int32, defaultEncoding string, encoding string) (string, error) {
	usedEncoding := defaultEncoding
	if len(encoding) > 0 {
		usedEncoding = encoding
	}
	var multiplier int32
	switch usedEncoding {
	case "UTF-8":
		multiplier = 8
	case "UTF-16":
		multiplier = 16
	}
	return io.ReadString("", uint32(stringLength*multiplier), usedEncoding)
}

func SerializeS7String(io utils.WriteBuffer, value values.PlcValue, stringLength int32, defaultEncoding string, encoding string) error {
	usedEncoding := defaultEncoding
	if len(encoding) > 0 {
		usedEncoding = encoding
	}
	var multiplier int32
	switch usedEncoding {
	case "UTF-8":
		multiplier = 8
	case "UTF-16":
		multiplier = 16
	}
	return io.WriteString("", uint32(stringLength*multiplier), usedEncoding, value.GetString())
}

func ParseS7Char(io utils.ReadBuffer, defaultEncoding string, encoding string) (string, error) {
	usedEncoding := defaultEncoding
	if len(encoding) > 0 {
		usedEncoding = encoding
	}
	return io.ReadString("value", uint32(8), usedEncoding)
}

func SerializeS7Char(io utils.WriteBuffer, value values.PlcValue, defaultEncoding string, encoding string) error {
	// TODO: This sort of looks wrong.
	return io.WriteUint8("", 8, value.GetUint8())
}

func RightShift3(readBuffer utils.ReadBuffer) (interface{}, error) {
	return uint16(0), nil
}

func LeftShift3(writeBuffer utils.WriteBuffer, valueLength uint16) error {
	return nil
}

func EventItemLength(readBuffer utils.ReadBuffer, valueLength uint16) uint16 {
	return 0
}

func BcdToInt(readBuffer utils.ReadBuffer) (interface{}, error) {
	return uint8(0), nil
}

func ByteToBcd(writeBuffer utils.WriteBuffer, value uint8) error {
	return nil
}

func S7msecToInt(readBuffer utils.ReadBuffer) (interface{}, error) {
	return uint16(0), nil
}

func IntToS7msec(writeBuffer utils.WriteBuffer, value uint16) error {
	return nil
}
