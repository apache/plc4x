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
package org.apache.plc4x.java.plc4x.readwrite;

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

public abstract class Plc4xMessage implements Message {

  // Abstract accessors for discriminator values.
  public abstract Plc4xRequestType getRequestType();

  // Constant values.
  public static final Short VERSION = 0x01;

  // Properties.
  protected final int requestId;

  public Plc4xMessage(int requestId) {
    super();
    this.requestId = requestId;
  }

  public int getRequestId() {
    return requestId;
  }

  public short getVersion() {
    return VERSION;
  }

  protected abstract void serializePlc4xMessageChild(WriteBuffer writeBuffer)
      throws SerializationException;

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("Plc4xMessage");

    // Const Field (version)
    writeConstField(
        "version",
        VERSION,
        writeUnsignedShort(writeBuffer, 8),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Implicit Field (packetLength) (Used for parsing, but its value is not stored as it's
    // implicitly given by the objects content)
    int packetLength = (int) (getLengthInBytes());
    writeImplicitField(
        "packetLength",
        packetLength,
        writeUnsignedInt(writeBuffer, 16),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (requestId)
    writeSimpleField(
        "requestId",
        requestId,
        writeUnsignedInt(writeBuffer, 16),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Discriminator Field (requestType) (Used as input to a switch field)
    writeDiscriminatorEnumField(
        "requestType",
        "Plc4xRequestType",
        getRequestType(),
        new DataWriterEnumDefault<>(
            Plc4xRequestType::getValue, Plc4xRequestType::name, writeUnsignedShort(writeBuffer, 8)),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Switch field (Serialize the sub-type)
    serializePlc4xMessageChild(writeBuffer);

    writeBuffer.popContext("Plc4xMessage");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    Plc4xMessage _value = this;

    // Const Field (version)
    lengthInBits += 8;

    // Implicit Field (packetLength)
    lengthInBits += 16;

    // Simple field (requestId)
    lengthInBits += 16;

    // Discriminator Field (requestType)
    lengthInBits += 8;

    // Length of sub-type elements will be added by sub-type...

    return lengthInBits;
  }

  public static Plc4xMessage staticParse(ReadBuffer readBuffer, Object... args)
      throws ParseException {
    PositionAware positionAware = readBuffer;
    return staticParse(readBuffer);
  }

  public static Plc4xMessage staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("Plc4xMessage");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    short version =
        readConstField(
            "version",
            readUnsignedShort(readBuffer, 8),
            Plc4xMessage.VERSION,
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    int packetLength =
        readImplicitField(
            "packetLength",
            readUnsignedInt(readBuffer, 16),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    int requestId =
        readSimpleField(
            "requestId",
            readUnsignedInt(readBuffer, 16),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    Plc4xRequestType requestType =
        readDiscriminatorField(
            "requestType",
            new DataReaderEnumDefault<>(
                Plc4xRequestType::enumForValue, readUnsignedShort(readBuffer, 8)),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    Plc4xMessageBuilder builder = null;
    if (EvaluationHelper.equals(requestType, Plc4xRequestType.CONNECT_REQUEST)) {
      builder = Plc4xConnectRequest.staticParseBuilder(readBuffer);
    } else if (EvaluationHelper.equals(requestType, Plc4xRequestType.CONNECT_RESPONSE)) {
      builder = Plc4xConnectResponse.staticParseBuilder(readBuffer);
    } else if (EvaluationHelper.equals(requestType, Plc4xRequestType.READ_REQUEST)) {
      builder = Plc4xReadRequest.staticParseBuilder(readBuffer);
    } else if (EvaluationHelper.equals(requestType, Plc4xRequestType.READ_RESPONSE)) {
      builder = Plc4xReadResponse.staticParseBuilder(readBuffer);
    } else if (EvaluationHelper.equals(requestType, Plc4xRequestType.WRITE_REQUEST)) {
      builder = Plc4xWriteRequest.staticParseBuilder(readBuffer);
    } else if (EvaluationHelper.equals(requestType, Plc4xRequestType.WRITE_RESPONSE)) {
      builder = Plc4xWriteResponse.staticParseBuilder(readBuffer);
    }
    if (builder == null) {
      throw new ParseException(
          "Unsupported case for discriminated type"
              + " parameters ["
              + "requestType="
              + requestType
              + "]");
    }

    readBuffer.closeContext("Plc4xMessage");
    // Create the instance
    Plc4xMessage _plc4xMessage = builder.build(requestId);
    return _plc4xMessage;
  }

  public static interface Plc4xMessageBuilder {
    Plc4xMessage build(int requestId);
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof Plc4xMessage)) {
      return false;
    }
    Plc4xMessage that = (Plc4xMessage) o;
    return (getRequestId() == that.getRequestId()) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getRequestId());
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
