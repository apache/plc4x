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

public class COTPPacketDisconnectRequest extends COTPPacket implements Message {

  // Accessors for discriminator values.
  public Short getTpduCode() {
    return (short) 0x80;
  }

  // Properties.
  protected final int destinationReference;
  protected final int sourceReference;
  protected final COTPProtocolClass protocolClass;

  // Arguments.
  protected final Integer cotpLen;

  public COTPPacketDisconnectRequest(
      List<COTPParameter> parameters,
      S7Message payload,
      int destinationReference,
      int sourceReference,
      COTPProtocolClass protocolClass,
      Integer cotpLen) {
    super(parameters, payload, cotpLen);
    this.destinationReference = destinationReference;
    this.sourceReference = sourceReference;
    this.protocolClass = protocolClass;
    this.cotpLen = cotpLen;
  }

  public int getDestinationReference() {
    return destinationReference;
  }

  public int getSourceReference() {
    return sourceReference;
  }

  public COTPProtocolClass getProtocolClass() {
    return protocolClass;
  }

  @Override
  protected void serializeCOTPPacketChild(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("COTPPacketDisconnectRequest");

    // Simple Field (destinationReference)
    writeSimpleField(
        "destinationReference", destinationReference, writeUnsignedInt(writeBuffer, 16));

    // Simple Field (sourceReference)
    writeSimpleField("sourceReference", sourceReference, writeUnsignedInt(writeBuffer, 16));

    // Simple Field (protocolClass)
    writeSimpleEnumField(
        "protocolClass",
        "COTPProtocolClass",
        protocolClass,
        new DataWriterEnumDefault<>(
            COTPProtocolClass::getValue,
            COTPProtocolClass::name,
            writeUnsignedShort(writeBuffer, 8)));

    writeBuffer.popContext("COTPPacketDisconnectRequest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    COTPPacketDisconnectRequest _value = this;

    // Simple field (destinationReference)
    lengthInBits += 16;

    // Simple field (sourceReference)
    lengthInBits += 16;

    // Simple field (protocolClass)
    lengthInBits += 8;

    return lengthInBits;
  }

  public static COTPPacketBuilder staticParseCOTPPacketBuilder(
      ReadBuffer readBuffer, Integer cotpLen) throws ParseException {
    readBuffer.pullContext("COTPPacketDisconnectRequest");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    int destinationReference =
        readSimpleField("destinationReference", readUnsignedInt(readBuffer, 16));

    int sourceReference = readSimpleField("sourceReference", readUnsignedInt(readBuffer, 16));

    COTPProtocolClass protocolClass =
        readEnumField(
            "protocolClass",
            "COTPProtocolClass",
            new DataReaderEnumDefault<>(
                COTPProtocolClass::enumForValue, readUnsignedShort(readBuffer, 8)));

    readBuffer.closeContext("COTPPacketDisconnectRequest");
    // Create the instance
    return new COTPPacketDisconnectRequestBuilderImpl(
        destinationReference, sourceReference, protocolClass, cotpLen);
  }

  public static class COTPPacketDisconnectRequestBuilderImpl
      implements COTPPacket.COTPPacketBuilder {
    private final int destinationReference;
    private final int sourceReference;
    private final COTPProtocolClass protocolClass;
    private final Integer cotpLen;

    public COTPPacketDisconnectRequestBuilderImpl(
        int destinationReference,
        int sourceReference,
        COTPProtocolClass protocolClass,
        Integer cotpLen) {
      this.destinationReference = destinationReference;
      this.sourceReference = sourceReference;
      this.protocolClass = protocolClass;
      this.cotpLen = cotpLen;
    }

    public COTPPacketDisconnectRequest build(
        List<COTPParameter> parameters, S7Message payload, Integer cotpLen) {
      COTPPacketDisconnectRequest cOTPPacketDisconnectRequest =
          new COTPPacketDisconnectRequest(
              parameters, payload, destinationReference, sourceReference, protocolClass, cotpLen);
      return cOTPPacketDisconnectRequest;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof COTPPacketDisconnectRequest)) {
      return false;
    }
    COTPPacketDisconnectRequest that = (COTPPacketDisconnectRequest) o;
    return (getDestinationReference() == that.getDestinationReference())
        && (getSourceReference() == that.getSourceReference())
        && (getProtocolClass() == that.getProtocolClass())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(), getDestinationReference(), getSourceReference(), getProtocolClass());
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
