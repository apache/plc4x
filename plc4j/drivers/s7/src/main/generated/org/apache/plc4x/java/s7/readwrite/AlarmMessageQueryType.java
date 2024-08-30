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

public class AlarmMessageQueryType implements Message {

  // Constant values.
  public static final Integer DATALENGTH = 0xFFFF;

  // Properties.
  protected final short functionId;
  protected final short numberOfObjects;
  protected final DataTransportErrorCode returnCode;
  protected final DataTransportSize transportSize;
  protected final List<AlarmMessageObjectQueryType> messageObjects;

  public AlarmMessageQueryType(
      short functionId,
      short numberOfObjects,
      DataTransportErrorCode returnCode,
      DataTransportSize transportSize,
      List<AlarmMessageObjectQueryType> messageObjects) {
    super();
    this.functionId = functionId;
    this.numberOfObjects = numberOfObjects;
    this.returnCode = returnCode;
    this.transportSize = transportSize;
    this.messageObjects = messageObjects;
  }

  public short getFunctionId() {
    return functionId;
  }

  public short getNumberOfObjects() {
    return numberOfObjects;
  }

  public DataTransportErrorCode getReturnCode() {
    return returnCode;
  }

  public DataTransportSize getTransportSize() {
    return transportSize;
  }

  public List<AlarmMessageObjectQueryType> getMessageObjects() {
    return messageObjects;
  }

  public int getDataLength() {
    return DATALENGTH;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("AlarmMessageQueryType");

    // Simple Field (functionId)
    writeSimpleField("functionId", functionId, writeUnsignedShort(writeBuffer, 8));

    // Simple Field (numberOfObjects)
    writeSimpleField("numberOfObjects", numberOfObjects, writeUnsignedShort(writeBuffer, 8));

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

    // Const Field (dataLength)
    writeConstField("dataLength", DATALENGTH, writeUnsignedInt(writeBuffer, 16));

    // Array Field (messageObjects)
    writeComplexTypeArrayField("messageObjects", messageObjects, writeBuffer);

    writeBuffer.popContext("AlarmMessageQueryType");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    AlarmMessageQueryType _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (functionId)
    lengthInBits += 8;

    // Simple field (numberOfObjects)
    lengthInBits += 8;

    // Simple field (returnCode)
    lengthInBits += 8;

    // Simple field (transportSize)
    lengthInBits += 8;

    // Const Field (dataLength)
    lengthInBits += 16;

    // Array field
    if (messageObjects != null) {
      int i = 0;
      for (AlarmMessageObjectQueryType element : messageObjects) {
        ThreadLocalHelper.lastItemThreadLocal.set(++i >= messageObjects.size());
        lengthInBits += element.getLengthInBits();
      }
    }

    return lengthInBits;
  }

  public static AlarmMessageQueryType staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("AlarmMessageQueryType");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    short functionId = readSimpleField("functionId", readUnsignedShort(readBuffer, 8));

    short numberOfObjects = readSimpleField("numberOfObjects", readUnsignedShort(readBuffer, 8));

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

    int dataLength =
        readConstField(
            "dataLength", readUnsignedInt(readBuffer, 16), AlarmMessageQueryType.DATALENGTH);

    List<AlarmMessageObjectQueryType> messageObjects =
        readCountArrayField(
            "messageObjects",
            readComplex(() -> AlarmMessageObjectQueryType.staticParse(readBuffer), readBuffer),
            numberOfObjects);

    readBuffer.closeContext("AlarmMessageQueryType");
    // Create the instance
    AlarmMessageQueryType _alarmMessageQueryType;
    _alarmMessageQueryType =
        new AlarmMessageQueryType(
            functionId, numberOfObjects, returnCode, transportSize, messageObjects);
    return _alarmMessageQueryType;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof AlarmMessageQueryType)) {
      return false;
    }
    AlarmMessageQueryType that = (AlarmMessageQueryType) o;
    return (getFunctionId() == that.getFunctionId())
        && (getNumberOfObjects() == that.getNumberOfObjects())
        && (getReturnCode() == that.getReturnCode())
        && (getTransportSize() == that.getTransportSize())
        && (getMessageObjects() == that.getMessageObjects())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        getFunctionId(),
        getNumberOfObjects(),
        getReturnCode(),
        getTransportSize(),
        getMessageObjects());
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
