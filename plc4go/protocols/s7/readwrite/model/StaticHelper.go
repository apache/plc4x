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
	"strings"
	"time"
	"unicode/utf16"

	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
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

// ParseS5Time decodes the BCD encoded S5TIME format: the high nibble selects the time base
// (10ms/100ms/1s/10s), the lower three nibbles are a three digit BCD counter.
func ParseS5Time(ctx context.Context, io utils.ReadBuffer) (uint32, error) {
	s5time, err := io.ReadUint16("s5time", 16)
	if err != nil {
		return 0, errors.Wrap(err, "Error parsing s5time")
	}
	timeValue := uint32(s5time&0x000F) + uint32((s5time&0x00F0)>>4)*10 + uint32((s5time&0x0F00)>>8)*100
	timeBase := uint32(10)
	for i := uint16(0); i < (s5time&0xF000)>>12; i++ {
		timeBase *= 10
	}
	totalMs := timeValue * timeBase
	if totalMs > 9990000 {
		totalMs = 9990000
	}
	return totalMs, nil
}

func SerializeS5Time(ctx context.Context, io utils.WriteBuffer, value values.PlcValue) error {
	totalMs := value.GetDuration().Milliseconds()
	var timeBase, timeValue int64
	switch {
	case totalMs < 0 || totalMs > 9990000:
		// Out of the representable range, S5TIME stays 0.
	case totalMs <= 9990:
		timeBase, timeValue = 0, totalMs/10
	case totalMs <= 99900:
		timeBase, timeValue = 1, totalMs/100
	case totalMs <= 999000:
		timeBase, timeValue = 2, totalMs/1000
	default:
		timeBase, timeValue = 3, totalMs/10000
	}
	units := timeValue % 10
	tens := (timeValue / 10) % 10
	hundreds := (timeValue / 100) % 10
	s5time := uint16(timeBase<<12 | hundreds<<8 | tens<<4 | units)
	return io.WriteUint16("s5time", 16, s5time)
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

// daysBetweenUnixAndSiemensEpoch is the offset between 1970-01-01 and the Siemens epoch 1990-01-01.
const daysBetweenUnixAndSiemensEpoch = 7305

func ParseTiaDate(ctx context.Context, io utils.ReadBuffer) (uint16, error) {
	daysSinceSiemensEpoch, err := io.ReadUint16("daysSinceSiemensEpoch", 16)
	if err != nil {
		return 0, errors.Wrap(err, "Error parsing daysSinceSiemensEpoch")
	}
	if daysSinceSiemensEpoch > 0xFFFF-daysBetweenUnixAndSiemensEpoch {
		return 0xFFFF, nil
	}
	return daysSinceSiemensEpoch + daysBetweenUnixAndSiemensEpoch, nil
}

func SerializeTiaDate(ctx context.Context, io utils.WriteBuffer, value values.PlcValue) error {
	var daysSinceSiemensEpoch uint16
	if date, ok := value.(spiValues.PlcDATE); ok {
		daysSinceSiemensEpoch = date.GetDaysSinceSiemensEpoch()
	} else {
		siemensEpoch := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)
		daysSinceSiemensEpoch = uint16(value.GetDate().Sub(siemensEpoch).Hours() / 24)
	}
	return io.WriteUint16("daysSinceSiemensEpoch", 16, daysSinceSiemensEpoch)
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

// ParseS7String reads an S7 STRING/WSTRING: a max-length prefix, a current-length prefix and
// max-length characters of which only the first current-length carry data. Up to stringLength
// characters are consumed from the buffer so trailing padding is gobbled up as well.
func ParseS7String(ctx context.Context, io utils.ReadBuffer, stringLength int32, encoding string) (string, error) {
	switch {
	case strings.EqualFold(encoding, "UTF8"):
		if _, err := io.ReadUint8("maxLength", 8); err != nil {
			return "", errors.Wrap(err, "Error parsing max length")
		}
		totalStringLength, err := io.ReadUint8("totalStringLength", 8)
		if err != nil {
			return "", errors.Wrap(err, "Error parsing total string length")
		}
		data := make([]byte, 0, totalStringLength)
		for i := int32(0); i < stringLength && io.HasMore(8); i++ {
			curByte, err := io.ReadUint8("", 8)
			if err != nil {
				return "", errors.Wrap(err, "Error parsing character")
			}
			if i < int32(totalStringLength) {
				data = append(data, curByte)
			}
		}
		return string(data), nil
	case strings.EqualFold(encoding, "UTF16") || strings.EqualFold(encoding, "UTF16BE"):
		if _, err := io.ReadUint16("maxLength", 16); err != nil {
			return "", errors.Wrap(err, "Error parsing max length")
		}
		totalStringLength, err := io.ReadUint16("totalStringLength", 16)
		if err != nil {
			return "", errors.Wrap(err, "Error parsing total string length")
		}
		units := make([]uint16, 0, totalStringLength)
		for i := int32(0); i < stringLength && io.HasMore(16); i++ {
			curUnit, err := io.ReadUint16("", 16)
			if err != nil {
				return "", errors.Wrap(err, "Error parsing character")
			}
			if i < int32(totalStringLength) {
				units = append(units, curUnit)
			}
		}
		return string(utf16.Decode(units)), nil
	default:
		return "", errors.Errorf("unsupported string encoding %s", encoding)
	}
}

// SerializeS7String writes an S7 STRING/WSTRING: [maxLen][curLen][maxLen chars, zero padded].
func SerializeS7String(ctx context.Context, io utils.WriteBuffer, value values.PlcValue, stringLength int32, encoding string) error {
	switch {
	case strings.EqualFold(encoding, "UTF8"):
		maxStringLength := int(stringLength)
		if maxStringLength > 254 {
			maxStringLength = 254
		}
		data := []byte(value.GetString())
		if len(data) > maxStringLength {
			data = data[:maxStringLength]
		}
		if err := io.WriteUint8("maxLength", 8, uint8(maxStringLength)); err != nil {
			return errors.Wrap(err, "Error serializing max length")
		}
		if err := io.WriteUint8("totalStringLength", 8, uint8(len(data))); err != nil {
			return errors.Wrap(err, "Error serializing total string length")
		}
		padded := make([]byte, maxStringLength)
		copy(padded, data)
		return io.WriteByteArray("chars", padded)
	case strings.EqualFold(encoding, "UTF16") || strings.EqualFold(encoding, "UTF16BE"):
		maxStringLength := int(stringLength)
		if maxStringLength > 16382 {
			maxStringLength = 16382
		}
		units := utf16.Encode([]rune(value.GetString()))
		if len(units) > maxStringLength {
			units = units[:maxStringLength]
		}
		if err := io.WriteUint16("maxLength", 16, uint16(maxStringLength)); err != nil {
			return errors.Wrap(err, "Error serializing max length")
		}
		if err := io.WriteUint16("totalStringLength", 16, uint16(len(units))); err != nil {
			return errors.Wrap(err, "Error serializing total string length")
		}
		for i := 0; i < maxStringLength; i++ {
			var unit uint16
			if i < len(units) {
				unit = units[i]
			}
			if err := io.WriteUint16("", 16, unit); err != nil {
				return errors.Wrap(err, "Error serializing character")
			}
		}
		return nil
	default:
		return errors.Errorf("unsupported string encoding %s", encoding)
	}
}

func ParseS7Char(ctx context.Context, io utils.ReadBuffer, encoding string) (uint8, error) {
	return io.ReadUint8("", 8)
}

func SerializeS7Char(ctx context.Context, io utils.WriteBuffer, value values.PlcValue, encoding string) error {
	return io.WriteUint8("", 8, value.GetUint8())
}

// RightShift3 reads a 16 bit length field which is expressed in bits for the numeric transport
// sizes and in bytes for the octet-string like ones.
func RightShift3(ctx context.Context, readBuffer utils.ReadBuffer, dataTransportSize DataTransportSize) (any, error) {
	value, err := readBuffer.ReadUint16("valueLength", 16)
	if err != nil {
		return uint16(0), errors.Wrap(err, "Error parsing value length")
	}
	if dataTransportSize == DataTransportSize_OCTET_STRING ||
		dataTransportSize == DataTransportSize_REAL ||
		dataTransportSize == DataTransportSize_BIT {
		return value, nil
	}
	return value >> 3, nil
}

func LeftShift3(ctx context.Context, writeBuffer utils.WriteBuffer, valueLength uint16) error {
	return writeBuffer.WriteUint16("valueLength", 16, valueLength<<3)
}

// EventItemLength accounts for the pad byte after odd-length event payload items (as long as
// the buffer actually still contains it).
func EventItemLength(ctx context.Context, readBuffer utils.ReadBuffer, valueLength uint16) uint16 {
	if valueLength%2 == 0 {
		return valueLength
	}
	if rb, ok := readBuffer.(utils.ReadBufferByteBased); ok {
		remainingBytes := rb.GetTotalBytes() - uint64(rb.GetPos())
		if remainingBytes < uint64(valueLength)+1 {
			return valueLength
		}
	}
	return valueLength + 1
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
