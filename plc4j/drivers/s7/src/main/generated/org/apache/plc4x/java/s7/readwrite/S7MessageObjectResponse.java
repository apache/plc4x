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

public class S7MessageObjectResponse extends S7DataAlarmMessage implements Message {

  // Accessors for discriminator values.
  public Byte getCpuFunctionType() {
    return (byte) 0x08;
  }

  // Properties.
  protected final DataTransportErrorCode returnCode;
  protected final DataTransportSize transportSize;

  // Reserved Fields
  private Short reservedField0;

  public S7MessageObjectResponse(
      DataTransportErrorCode returnCode, DataTransportSize transportSize) {
    super();
    this.returnCode = returnCode;
    this.transportSize = transportSize;
  }

  public DataTransportErrorCode getReturnCode() {
    return returnCode;
  }

  public DataTransportSize getTransportSize() {
    return transportSize;
  }

  @Override
  protected void serializeS7DataAlarmMessageChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("S7MessageObjectResponse");

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

    // Reserved Field (reserved)
    writeReservedField(
        "reserved",
        reservedField0 != null ? reservedField0 : (short) 0x00,
        writeUnsignedShort(writeBuffer, 8));

    writeBuffer.popContext("S7MessageObjectResponse");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    S7MessageObjectResponse _value = this;

    // Simple field (returnCode)
    lengthInBits += 8;

    // Simple field (transportSize)
    lengthInBits += 8;

    // Reserved Field (reserved)
    lengthInBits += 8;

    return lengthInBits;
  }

  public static S7DataAlarmMessageBuilder staticParseS7DataAlarmMessageBuilder(
      ReadBuffer readBuffer, Byte cpuFunctionType) throws ParseException {
    readBuffer.pullContext("S7MessageObjectResponse");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    DataTransportErrorCode returnCode =
        readEnumField(
            "returnCode",
            "DataTransportErrorCode",
            new DataReaderEnumDefault<>(
                DataTransportErrorCode::enumForValue, readUnsignedShort(readBuffer, 8)));

    DataTransportSize transportSize =
        readEnumField(
            "transportSize",
            "DataTransportSize",
            new DataReaderEnumDefault<>(
                DataTransportSize::enumForValue, readUnsignedShort(readBuffer, 8)));

    Short reservedField0 =
        readReservedField("reserved", readUnsignedShort(readBuffer, 8), (short) 0x00);

    readBuffer.closeContext("S7MessageObjectResponse");
    // Create the instance
    return new S7MessageObjectResponseBuilderImpl(returnCode, transportSize, reservedField0);
  }

  public static class S7MessageObjectResponseBuilderImpl
      implements S7DataAlarmMessage.S7DataAlarmMessageBuilder {
    private final DataTransportErrorCode returnCode;
    private final DataTransportSize transportSize;
    private final Short reservedField0;

    public S7MessageObjectResponseBuilderImpl(
        DataTransportErrorCode returnCode, DataTransportSize transportSize, Short reservedField0) {
      this.returnCode = returnCode;
      this.transportSize = transportSize;
      this.reservedField0 = reservedField0;
    }

    public S7MessageObjectResponse build() {
      S7MessageObjectResponse s7MessageObjectResponse =
          new S7MessageObjectResponse(returnCode, transportSize);
      s7MessageObjectResponse.reservedField0 = reservedField0;
      return s7MessageObjectResponse;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof S7MessageObjectResponse)) {
      return false;
    }
    S7MessageObjectResponse that = (S7MessageObjectResponse) o;
    return (getReturnCode() == that.getReturnCode())
        && (getTransportSize() == that.getTransportSize())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getReturnCode(), getTransportSize());
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
