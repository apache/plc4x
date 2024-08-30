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
package org.apache.plc4x.java.s7.readwrite;

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

public class AssociatedQueryValueType implements Message {

  // Properties.
  protected final DataTransportErrorCode returnCode;
  protected final DataTransportSize transportSize;
  protected final int valueLength;
  protected final List<Short> data;

  public AssociatedQueryValueType(
      DataTransportErrorCode returnCode,
      DataTransportSize transportSize,
      int valueLength,
      List<Short> data) {
    super();
    this.returnCode = returnCode;
    this.transportSize = transportSize;
    this.valueLength = valueLength;
    this.data = data;
  }

  public DataTransportErrorCode getReturnCode() {
    return returnCode;
  }

  public DataTransportSize getTransportSize() {
    return transportSize;
  }

  public int getValueLength() {
    return valueLength;
  }

  public List<Short> getData() {
    return data;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("AssociatedQueryValueType");

    // Simple Field (returnCode)
    writeSimpleEnumField(
        "returnCode",
        "DataTransportErrorCode",
        returnCode,
        new DataWriterEnumDefault<>(
            DataTransportErrorCode::getValue,
            DataTransportErrorCode::name,
            writeUnsignedShort(writeBuffer, 8)));

    // Simple Field (transportSize)
    writeSimpleEnumField(
        "transportSize",
        "DataTransportSize",
        transportSize,
        new DataWriterEnumDefault<>(
            DataTransportSize::getValue,
            DataTransportSize::name,
            writeUnsignedShort(writeBuffer, 8)));

    // Simple Field (valueLength)
    writeSimpleField("valueLength", valueLength, writeUnsignedInt(writeBuffer, 16));

    // Array Field (data)
    writeSimpleTypeArrayField("data", data, writeUnsignedShort(writeBuffer, 8));

    writeBuffer.popContext("AssociatedQueryValueType");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    AssociatedQueryValueType _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (returnCode)
    lengthInBits += 8;

    // Simple field (transportSize)
    lengthInBits += 8;

    // Simple field (valueLength)
    lengthInBits += 16;

    // Array field
    if (data != null) {
      lengthInBits += 8 * data.size();
    }

    return lengthInBits;
  }

  public static AssociatedQueryValueType staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("AssociatedQueryValueType");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    DataTransportErrorCode returnCode =
        readEnumField(
            "returnCode",
            "DataTransportErrorCode",
            readEnum(DataTransportErrorCode::enumForValue, readUnsignedShort(readBuffer, 8)));

    DataTransportSize transportSize =
        readEnumField(
            "transportSize",
            "DataTransportSize",
            readEnum(DataTransportSize::enumForValue, readUnsignedShort(readBuffer, 8)));

    int valueLength = readSimpleField("valueLength", readUnsignedInt(readBuffer, 16));

    List<Short> data = readCountArrayField("data", readUnsignedShort(readBuffer, 8), valueLength);

    readBuffer.closeContext("AssociatedQueryValueType");
    // Create the instance
    AssociatedQueryValueType _associatedQueryValueType;
    _associatedQueryValueType =
        new AssociatedQueryValueType(returnCode, transportSize, valueLength, data);
    return _associatedQueryValueType;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof AssociatedQueryValueType)) {
      return false;
    }
    AssociatedQueryValueType that = (AssociatedQueryValueType) o;
    return (getReturnCode() == that.getReturnCode())
        && (getTransportSize() == that.getTransportSize())
        && (getValueLength() == that.getValueLength())
        && (getData() == that.getData())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getReturnCode(), getTransportSize(), getValueLength(), getData());
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
