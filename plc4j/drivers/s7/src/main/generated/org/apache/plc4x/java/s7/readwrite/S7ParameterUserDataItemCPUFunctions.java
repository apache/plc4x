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

public class S7ParameterUserDataItemCPUFunctions extends S7ParameterUserDataItem
    implements Message {

  // Accessors for discriminator values.
  public Short getItemType() {
    return (short) 0x12;
  }

  // Properties.
  protected final short method;
  protected final byte cpuFunctionType;
  protected final byte cpuFunctionGroup;
  protected final short cpuSubfunction;
  protected final short sequenceNumber;
  protected final Short dataUnitReferenceNumber;
  protected final Short lastDataUnit;
  protected final Integer errorCode;

  public S7ParameterUserDataItemCPUFunctions(
      short method,
      byte cpuFunctionType,
      byte cpuFunctionGroup,
      short cpuSubfunction,
      short sequenceNumber,
      Short dataUnitReferenceNumber,
      Short lastDataUnit,
      Integer errorCode) {
    super();
    this.method = method;
    this.cpuFunctionType = cpuFunctionType;
    this.cpuFunctionGroup = cpuFunctionGroup;
    this.cpuSubfunction = cpuSubfunction;
    this.sequenceNumber = sequenceNumber;
    this.dataUnitReferenceNumber = dataUnitReferenceNumber;
    this.lastDataUnit = lastDataUnit;
    this.errorCode = errorCode;
  }

  public short getMethod() {
    return method;
  }

  public byte getCpuFunctionType() {
    return cpuFunctionType;
  }

  public byte getCpuFunctionGroup() {
    return cpuFunctionGroup;
  }

  public short getCpuSubfunction() {
    return cpuSubfunction;
  }

  public short getSequenceNumber() {
    return sequenceNumber;
  }

  public Short getDataUnitReferenceNumber() {
    return dataUnitReferenceNumber;
  }

  public Short getLastDataUnit() {
    return lastDataUnit;
  }

  public Integer getErrorCode() {
    return errorCode;
  }

  @Override
  protected void serializeS7ParameterUserDataItemChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("S7ParameterUserDataItemCPUFunctions");

    // Implicit Field (itemLength) (Used for parsing, but its value is not stored as it's implicitly
    // given by the objects content)
    short itemLength = (short) ((getLengthInBytes()) - (2));
    writeImplicitField("itemLength", itemLength, writeUnsignedShort(writeBuffer, 8));

    // Simple Field (method)
    writeSimpleField("method", method, writeUnsignedShort(writeBuffer, 8));

    // Simple Field (cpuFunctionType)
    writeSimpleField("cpuFunctionType", cpuFunctionType, writeUnsignedByte(writeBuffer, 4));

    // Simple Field (cpuFunctionGroup)
    writeSimpleField("cpuFunctionGroup", cpuFunctionGroup, writeUnsignedByte(writeBuffer, 4));

    // Simple Field (cpuSubfunction)
    writeSimpleField("cpuSubfunction", cpuSubfunction, writeUnsignedShort(writeBuffer, 8));

    // Simple Field (sequenceNumber)
    writeSimpleField("sequenceNumber", sequenceNumber, writeUnsignedShort(writeBuffer, 8));

    // Optional Field (dataUnitReferenceNumber) (Can be skipped, if the value is null)
    writeOptionalField(
        "dataUnitReferenceNumber", dataUnitReferenceNumber, writeUnsignedShort(writeBuffer, 8));

    // Optional Field (lastDataUnit) (Can be skipped, if the value is null)
    writeOptionalField("lastDataUnit", lastDataUnit, writeUnsignedShort(writeBuffer, 8));

    // Optional Field (errorCode) (Can be skipped, if the value is null)
    writeOptionalField("errorCode", errorCode, writeUnsignedInt(writeBuffer, 16));

    writeBuffer.popContext("S7ParameterUserDataItemCPUFunctions");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    S7ParameterUserDataItemCPUFunctions _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Implicit Field (itemLength)
    lengthInBits += 8;

    // Simple field (method)
    lengthInBits += 8;

    // Simple field (cpuFunctionType)
    lengthInBits += 4;

    // Simple field (cpuFunctionGroup)
    lengthInBits += 4;

    // Simple field (cpuSubfunction)
    lengthInBits += 8;

    // Simple field (sequenceNumber)
    lengthInBits += 8;

    // Optional Field (dataUnitReferenceNumber)
    if (dataUnitReferenceNumber != null) {
      lengthInBits += 8;
    }

    // Optional Field (lastDataUnit)
    if (lastDataUnit != null) {
      lengthInBits += 8;
    }

    // Optional Field (errorCode)
    if (errorCode != null) {
      lengthInBits += 16;
    }

    return lengthInBits;
  }

  public static S7ParameterUserDataItemBuilder staticParseS7ParameterUserDataItemBuilder(
      ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("S7ParameterUserDataItemCPUFunctions");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    short itemLength = readImplicitField("itemLength", readUnsignedShort(readBuffer, 8));

    short method = readSimpleField("method", readUnsignedShort(readBuffer, 8));

    byte cpuFunctionType = readSimpleField("cpuFunctionType", readUnsignedByte(readBuffer, 4));

    byte cpuFunctionGroup = readSimpleField("cpuFunctionGroup", readUnsignedByte(readBuffer, 4));

    short cpuSubfunction = readSimpleField("cpuSubfunction", readUnsignedShort(readBuffer, 8));

    short sequenceNumber = readSimpleField("sequenceNumber", readUnsignedShort(readBuffer, 8));

    Short dataUnitReferenceNumber =
        readOptionalField(
            "dataUnitReferenceNumber", readUnsignedShort(readBuffer, 8), (cpuFunctionType) == (8));

    Short lastDataUnit =
        readOptionalField(
            "lastDataUnit", readUnsignedShort(readBuffer, 8), (cpuFunctionType) == (8));

    Integer errorCode =
        readOptionalField("errorCode", readUnsignedInt(readBuffer, 16), (cpuFunctionType) == (8));

    readBuffer.closeContext("S7ParameterUserDataItemCPUFunctions");
    // Create the instance
    return new S7ParameterUserDataItemCPUFunctionsBuilderImpl(
        method,
        cpuFunctionType,
        cpuFunctionGroup,
        cpuSubfunction,
        sequenceNumber,
        dataUnitReferenceNumber,
        lastDataUnit,
        errorCode);
  }

  public static class S7ParameterUserDataItemCPUFunctionsBuilderImpl
      implements S7ParameterUserDataItem.S7ParameterUserDataItemBuilder {
    private final short method;
    private final byte cpuFunctionType;
    private final byte cpuFunctionGroup;
    private final short cpuSubfunction;
    private final short sequenceNumber;
    private final Short dataUnitReferenceNumber;
    private final Short lastDataUnit;
    private final Integer errorCode;

    public S7ParameterUserDataItemCPUFunctionsBuilderImpl(
        short method,
        byte cpuFunctionType,
        byte cpuFunctionGroup,
        short cpuSubfunction,
        short sequenceNumber,
        Short dataUnitReferenceNumber,
        Short lastDataUnit,
        Integer errorCode) {
      this.method = method;
      this.cpuFunctionType = cpuFunctionType;
      this.cpuFunctionGroup = cpuFunctionGroup;
      this.cpuSubfunction = cpuSubfunction;
      this.sequenceNumber = sequenceNumber;
      this.dataUnitReferenceNumber = dataUnitReferenceNumber;
      this.lastDataUnit = lastDataUnit;
      this.errorCode = errorCode;
    }

    public S7ParameterUserDataItemCPUFunctions build() {
      S7ParameterUserDataItemCPUFunctions s7ParameterUserDataItemCPUFunctions =
          new S7ParameterUserDataItemCPUFunctions(
              method,
              cpuFunctionType,
              cpuFunctionGroup,
              cpuSubfunction,
              sequenceNumber,
              dataUnitReferenceNumber,
              lastDataUnit,
              errorCode);
      return s7ParameterUserDataItemCPUFunctions;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof S7ParameterUserDataItemCPUFunctions)) {
      return false;
    }
    S7ParameterUserDataItemCPUFunctions that = (S7ParameterUserDataItemCPUFunctions) o;
    return (getMethod() == that.getMethod())
        && (getCpuFunctionType() == that.getCpuFunctionType())
        && (getCpuFunctionGroup() == that.getCpuFunctionGroup())
        && (getCpuSubfunction() == that.getCpuSubfunction())
        && (getSequenceNumber() == that.getSequenceNumber())
        && (getDataUnitReferenceNumber() == that.getDataUnitReferenceNumber())
        && (getLastDataUnit() == that.getLastDataUnit())
        && (getErrorCode() == that.getErrorCode())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getMethod(),
        getCpuFunctionType(),
        getCpuFunctionGroup(),
        getCpuSubfunction(),
        getSequenceNumber(),
        getDataUnitReferenceNumber(),
        getLastDataUnit(),
        getErrorCode());
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
