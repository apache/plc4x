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
package org.apache.plc4x.java.ads.discovery.readwrite;

import static org.apache.plc4x.java.spi.codegen.fields.FieldReaderFactory.*;
import static org.apache.plc4x.java.spi.codegen.fields.FieldWriterFactory.*;
import static org.apache.plc4x.java.spi.codegen.io.DataReaderFactory.*;
import static org.apache.plc4x.java.spi.codegen.io.DataWriterFactory.*;
import static org.apache.plc4x.java.spi.generation.StaticHelper.*;

import java.time.*;
import java.util.*;
import org.apache.plc4x.java.api.exceptions.*;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.codegen.*;
import org.apache.plc4x.java.spi.codegen.fields.*;
import org.apache.plc4x.java.spi.codegen.io.*;
import org.apache.plc4x.java.spi.generation.*;

// Code generated by code-generation. DO NOT EDIT.

public class AmsString implements Message {

  // Properties.
  protected final String text;

  public AmsString(String text) {
    super();
    this.text = text;
  }

  public String getText() {
    return text;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("AmsString");

    // Implicit Field (strLen) (Used for parsing, but its value is not stored as it's implicitly
    // given by the objects content)
    int strLen = (int) ((STR_LEN(getText())) + (1));
    writeImplicitField("strLen", strLen, writeUnsignedInt(writeBuffer, 16));

    // Simple Field (text)
    writeSimpleField(
        "text",
        text,
        writeString(writeBuffer, (8) * (((strLen) - (1)))),
        WithOption.WithEncoding("UTF-8"));

    // Reserved Field (reserved)
    writeReservedField("reserved", (short) 0x00, writeUnsignedShort(writeBuffer, 8));

    writeBuffer.popContext("AmsString");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    AmsString _value = this;

    // Implicit Field (strLen)
    lengthInBits += 16;

    // Simple field (text)
    lengthInBits += (8) * ((((STR_LEN(getText())) + (1)) - (1)));

    // Reserved Field (reserved)
    lengthInBits += 8;

    return lengthInBits;
  }

  public static AmsString staticParse(ReadBuffer readBuffer, Object... args) throws ParseException {
    PositionAware positionAware = readBuffer;
    return staticParse(readBuffer);
  }

  public static AmsString staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("AmsString");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    int strLen = readImplicitField("strLen", readUnsignedInt(readBuffer, 16));

    String text =
        readSimpleField(
            "text",
            readString(readBuffer, (8) * (((strLen) - (1)))),
            WithOption.WithEncoding("UTF-8"));

    Short reservedField0 =
        readReservedField("reserved", readUnsignedShort(readBuffer, 8), (short) 0x00);

    readBuffer.closeContext("AmsString");
    // Create the instance
    AmsString _amsString;
    _amsString = new AmsString(text);
    return _amsString;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof AmsString)) {
      return false;
    }
    AmsString that = (AmsString) o;
    return (getText() == that.getText()) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getText());
  }

  @Override
  public String toString() {
    WriteBufferBoxBased writeBufferBoxBased = new WriteBufferBoxBased(true, true);
    try {
      writeBufferBoxBased.writeSerializable(this);
    } catch (SerializationException e) {
      throw new RuntimeException(e);
    }
    return "\n" + writeBufferBoxBased.getBox().toString() + "\n";
  }
}
