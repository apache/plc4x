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

public class ImplicitPascalStringTypeTest implements Message {

  // Properties.
  protected final String stringField;

  public ImplicitPascalStringTypeTest(String stringField) {
    super();
    this.stringField = stringField;
  }

  public String getStringField() {
    return stringField;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("ImplicitPascalStringTypeTest");

    // Implicit Field (stringLength) (Used for parsing, but its value is not stored as it's
    // implicitly given by the objects content)
    byte stringLength = (byte) (getStringField().length());
    writeImplicitField("stringLength", stringLength, writeSignedByte(writeBuffer, 8));

    // Simple Field (stringField)
    writeSimpleField("stringField", stringField, writeString(writeBuffer, stringLength));

    writeBuffer.popContext("ImplicitPascalStringTypeTest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    ImplicitPascalStringTypeTest _value = this;

    // Implicit Field (stringLength)
    lengthInBits += 8;

    // Simple field (stringField)
    lengthInBits += getStringField().length();

    return lengthInBits;
  }

  public static ImplicitPascalStringTypeTest staticParse(ReadBuffer readBuffer, Object... args)
      throws ParseException {
    PositionAware positionAware = readBuffer;
    return staticParse(readBuffer);
  }

  public static ImplicitPascalStringTypeTest staticParse(ReadBuffer readBuffer)
      throws ParseException {
    readBuffer.pullContext("ImplicitPascalStringTypeTest");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    byte stringLength = readImplicitField("stringLength", readSignedByte(readBuffer, 8));

    String stringField = readSimpleField("stringField", readString(readBuffer, stringLength));

    readBuffer.closeContext("ImplicitPascalStringTypeTest");
    // Create the instance
    ImplicitPascalStringTypeTest _implicitPascalStringTypeTest;
    _implicitPascalStringTypeTest = new ImplicitPascalStringTypeTest(stringField);
    return _implicitPascalStringTypeTest;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof ImplicitPascalStringTypeTest)) {
      return false;
    }
    ImplicitPascalStringTypeTest that = (ImplicitPascalStringTypeTest) o;
    return (getStringField() == that.getStringField()) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getStringField());
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
