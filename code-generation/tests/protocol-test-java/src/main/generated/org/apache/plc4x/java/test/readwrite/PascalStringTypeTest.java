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
package org.apache.plc4x.java.test.readwrite;

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

public class PascalStringTypeTest implements Message {

  // Properties.
  protected final byte stringLength;
  protected final String stringField;

  public PascalStringTypeTest(byte stringLength, String stringField) {
    super();
    this.stringLength = stringLength;
    this.stringField = stringField;
  }

  public byte getStringLength() {
    return stringLength;
  }

  public String getStringField() {
    return stringField;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("PascalStringTypeTest");

    // Simple Field (stringLength)
    writeSimpleField("stringLength", stringLength, writeSignedByte(writeBuffer, 8));

    // Simple Field (stringField)
    writeSimpleField("stringField", stringField, writeString(writeBuffer, stringLength));

    writeBuffer.popContext("PascalStringTypeTest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    PascalStringTypeTest _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (stringLength)
    lengthInBits += 8;

    // Simple field (stringField)
    lengthInBits += getStringLength();

    return lengthInBits;
  }

  public static PascalStringTypeTest staticParse(ReadBuffer readBuffer, Object... args)
      throws ParseException {
    PositionAware positionAware = readBuffer;
    return staticParse(readBuffer);
  }

  public static PascalStringTypeTest staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("PascalStringTypeTest");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    byte stringLength = readSimpleField("stringLength", readSignedByte(readBuffer, 8));

    String stringField = readSimpleField("stringField", readString(readBuffer, stringLength));

    readBuffer.closeContext("PascalStringTypeTest");
    // Create the instance
    PascalStringTypeTest _pascalStringTypeTest;
    _pascalStringTypeTest = new PascalStringTypeTest(stringLength, stringField);
    return _pascalStringTypeTest;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof PascalStringTypeTest)) {
      return false;
    }
    PascalStringTypeTest that = (PascalStringTypeTest) o;
    return (getStringLength() == that.getStringLength())
        && (getStringField() == that.getStringField())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getStringLength(), getStringField());
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
