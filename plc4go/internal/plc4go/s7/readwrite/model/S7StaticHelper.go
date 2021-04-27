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

package model

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"time"
)

func StaticHelperParseTiaTime(io utils.ReadBuffer) (uint32, error) {
	/*try {
	      int millisSinceMidnight = io.readInt(32);
	      return LocalTime.now().withHour(0).withMinute(0).withSecond(0).withNano(0).plus(
	      millisSinceMidnight, ChronoUnit.MILLIS);
	  } catch (ParseException e) {
	      return null;
	  }*/
	return 0, nil
}

func StaticHelperSerializeTiaTime(io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing TIME not implemented");
	return nil
}

func StaticHelperParseS5Time(io utils.ReadBuffer) (uint32, error) {
	/*try {
	      int stuff = io.readInt(16);
	      // TODO: Implement this correctly.
	      throw new NotImplementedException("S5TIME not implemented");
	  } catch (ParseException e) {
	      return null;
	  }*/
	return 0, nil
}

func StaticHelperSerializeS5Time(io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing S5TIME not implemented");
	return nil
}

func StaticHelperParseTiaLTime(io utils.ReadBuffer) (uint32, error) {
	//throw new NotImplementedException("LTIME not implemented");
	return 0, nil
}

func StaticHelperSerializeTiaLTime(io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing LTIME not implemented");
	return nil
}

func StaticHelperParseTiaTimeOfDay(io utils.ReadBuffer) (time.Time, error) {
	/*try {
	      long millisSinceMidnight = io.readUnsignedLong(32);
	      return LocalTime.now().withHour(0).withMinute(0).withSecond(0).withNano(0).plus(
	          millisSinceMidnight, ChronoUnit.MILLIS);
	  } catch (ParseException e) {
	      return null;
	  }*/
	return time.Time{}, nil
}

func StaticHelperSerializeTiaTimeOfDay(io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing TIME_OF_DAY not implemented");
	return nil
}

func StaticHelperParseTiaDate(io utils.ReadBuffer) (time.Time, error) {
	/*try {
	      int daysSince1990 = io.readUnsignedInt(16);
	      return LocalDate.now().withYear(1990).withDayOfMonth(1).withMonth(1).plus(daysSince1990, ChronoUnit.DAYS);
	  } catch (ParseException e) {
	      return null;
	  }*/
	return time.Time{}, nil
}

func StaticHelperSerializeTiaDate(io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing DATE not implemented");
	return nil
}

func StaticHelperParseTiaDateTime(io utils.ReadBuffer) (time.Time, error) {
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

func StaticHelperSerializeTiaDateTime(io utils.WriteBuffer, value values.PlcValue) error {
	//throw new NotImplementedException("Serializing DATE_AND_TIME not implemented");
	return nil
}

func StaticHelperParseS7String(io utils.ReadBuffer, stringLength int32, encoding string) (string, error) {
	var multiplier int32
	switch encoding {
	case "UTF-8":
		multiplier = 0
	case "UTF-16":
		multiplier = 16
	}
	readString, err := io.ReadString("", uint32(stringLength*multiplier))
	if err != nil {
		return "", err
	}
	return readString, nil
}

func StaticHelperSerializeS7String(io utils.WriteBuffer, value values.PlcValue, stringLength int32, encoding string) error {
	var multiplier int32
	switch encoding {
	case "UTF-8":
		multiplier = 0
	case "UTF-16":
		multiplier = 16
	}
	return io.WriteString("", uint8(stringLength*multiplier), encoding, value.GetString())
}

func StaticHelperParseS7Char(io utils.ReadBuffer, encoding string) (uint8, error) {
	return io.ReadUint8("", 8)
}

func StaticHelperSerializeS7Char(io utils.WriteBuffer, value values.PlcValue, encoding string) error {
	return io.WriteUint8("", 8, value.GetUint8())
}
