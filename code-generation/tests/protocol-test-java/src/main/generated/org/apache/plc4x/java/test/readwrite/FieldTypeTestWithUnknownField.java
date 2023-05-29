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

public class FieldTypeTestWithUnknownField implements Message {

  // Constant values.
  public static final Short CONSTFIELD = 5;

  // Properties.
  protected final short simpleField;
  protected final List<Short> arrayField;
  protected final EnumTypeParameters enumField;
  protected final Short optionalField;

  public FieldTypeTestWithUnknownField(
      short simpleField,
      List<Short> arrayField,
      EnumTypeParameters enumField,
      Short optionalField) {
    super();
    this.simpleField = simpleField;
    this.arrayField = arrayField;
    this.enumField = enumField;
    this.optionalField = optionalField;
  }

  public short getSimpleField() {
    return simpleField;
  }

  public List<Short> getArrayField() {
    return arrayField;
  }

  public EnumTypeParameters getEnumField() {
    return enumField;
  }

  public Short getOptionalField() {
    return optionalField;
  }

  public short getConstField() {
    return CONSTFIELD;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    throw new SerializationException("Unknown field not serializable");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    FieldTypeTestWithUnknownField _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (simpleField)
    lengthInBits += 8;

    // Array field
    if (arrayField != null) {
      lengthInBits += 8 * arrayField.size();
    }

    // Checksum Field (checksum)
    lengthInBits += 8;

    // Const Field (constField)
    lengthInBits += 8;

    // Enum Field (enumField)
    lengthInBits += 8;

    // Implicit Field (implicitField)
    lengthInBits += 8;

    // Optional Field (optionalField)
    if (optionalField != null) {
      lengthInBits += 8;
    }

    // Padding Field (padding)
    int _timesPadding = (int) (simpleField);
    while (_timesPadding-- > 0) {
      lengthInBits += 8;
    }

    // Reserved Field (reserved)
    lengthInBits += 8;

    // Unknown field
    lengthInBits += 16;

    return lengthInBits;
  }

  public static FieldTypeTestWithUnknownField staticParse(ReadBuffer readBuffer, Object... args)
      throws ParseException {
    PositionAware positionAware = readBuffer;
    return staticParse(readBuffer);
  }

  public static FieldTypeTestWithUnknownField staticParse(ReadBuffer readBuffer)
      throws ParseException {
    readBuffer.pullContext("FieldTypeTestWithUnknownField");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    short simpleField = readSimpleField("simpleField", readUnsignedShort(readBuffer, 8));

    List<Short> arrayField = readCountArrayField("arrayField", readUnsignedShort(readBuffer, 8), 5);

    short checksumField =
        readChecksumField(
            "checksumField",
            readUnsignedShort(readBuffer, 8),
            (short) (org.apache.plc4x.java.test.readwrite.utils.StaticHelper.crcUint8(-(23))));

    short constField =
        readConstField(
            "constField",
            readUnsignedShort(readBuffer, 8),
            FieldTypeTestWithUnknownField.CONSTFIELD);

    EnumTypeParameters enumField =
        readEnumField(
            "enumField",
            "EnumTypeParameters",
            readEnum(EnumTypeParameters::firstEnumForFieldIntType, readSignedByte(readBuffer, 8)));

    short implicitField = readImplicitField("implicitField", readUnsignedShort(readBuffer, 8));

    Short optionalField =
        readOptionalField("optionalField", readUnsignedShort(readBuffer, 8), (simpleField) == (5));

    readPaddingField(readUnsignedShort(readBuffer, 8), (int) (simpleField));

    Short reservedField0 =
        readReservedField("reserved", readUnsignedShort(readBuffer, 8), (short) 0x00);

    readUnknownField("unknown", readUnsignedInt(readBuffer, 16));

    readBuffer.closeContext("FieldTypeTestWithUnknownField");
    // Create the instance
    FieldTypeTestWithUnknownField _fieldTypeTestWithUnknownField;
    _fieldTypeTestWithUnknownField =
        new FieldTypeTestWithUnknownField(simpleField, arrayField, enumField, optionalField);
    return _fieldTypeTestWithUnknownField;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof FieldTypeTestWithUnknownField)) {
      return false;
    }
    FieldTypeTestWithUnknownField that = (FieldTypeTestWithUnknownField) o;
    return (getSimpleField() == that.getSimpleField())
        && (getArrayField() == that.getArrayField())
        && (getEnumField() == that.getEnumField())
        && (getOptionalField() == that.getOptionalField())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getSimpleField(), getArrayField(), getEnumField(), getOptionalField());
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
