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
package org.apache.plc4x.java.profinet.readwrite;

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

public class PascalString implements Message {

  // Properties.
  protected final String stringValue;

  public PascalString(String stringValue) {
    super();
    this.stringValue = stringValue;
  }

  public String getStringValue() {
    return stringValue;
  }

  public byte getStringLength() {
    return (byte) ((((getStringValue().length()) == (-(1))) ? 0 : getStringValue().length()));
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("PascalString");

    // Implicit Field (sLength) (Used for parsing, but its value is not stored as it's implicitly
    // given by the objects content)
    byte sLength =
        (byte) ((((getStringValue().length()) == (0)) ? -(1) : getStringValue().length()));
    writeImplicitField("sLength", sLength, writeSignedByte(writeBuffer, 8));

    // Simple Field (stringValue)
    writeSimpleField(
        "stringValue",
        stringValue,
        writeString(writeBuffer, (((sLength) == (-(1))) ? 0 : (sLength) * (8))));

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    byte stringLength = getStringLength();
    writeBuffer.writeVirtual("stringLength", stringLength);

    writeBuffer.popContext("PascalString");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    PascalString _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Implicit Field (sLength)
    lengthInBits += 8;

    // Simple field (stringValue)
    lengthInBits +=
        ((((((getStringValue().length()) == (0)) ? -(1) : getStringValue().length())) == (-(1)))
            ? 0
            : ((((getStringValue().length()) == (0)) ? -(1) : getStringValue().length())) * (8));

    // A virtual field doesn't have any in- or output.

    return lengthInBits;
  }

  public static PascalString staticParse(ReadBuffer readBuffer, Object... args)
      throws ParseException {
    PositionAware positionAware = readBuffer;
    return staticParse(readBuffer);
  }

  public static PascalString staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("PascalString");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    byte sLength = readImplicitField("sLength", readSignedByte(readBuffer, 8));

    String stringValue =
        readSimpleField(
            "stringValue", readString(readBuffer, (((sLength) == (-(1))) ? 0 : (sLength) * (8))));
    byte stringLength =
        readVirtualField(
            "stringLength",
            byte.class,
            (((stringValue.length()) == (-(1))) ? 0 : stringValue.length()));

    readBuffer.closeContext("PascalString");
    // Create the instance
    PascalString _pascalString;
    _pascalString = new PascalString(stringValue);
    return _pascalString;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof PascalString)) {
      return false;
    }
    PascalString that = (PascalString) o;
    return (getStringValue() == that.getStringValue()) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getStringValue());
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
