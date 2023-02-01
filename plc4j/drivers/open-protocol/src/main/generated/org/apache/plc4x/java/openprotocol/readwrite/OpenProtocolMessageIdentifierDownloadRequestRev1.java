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
package org.apache.plc4x.java.openprotocol.readwrite;

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

public class OpenProtocolMessageIdentifierDownloadRequestRev1
    extends OpenProtocolMessageIdentifierDownloadRequest implements Message {

  // Accessors for discriminator values.
  public Long getRevision() {
    return (long) 1;
  }

  // Properties.
  protected final String identifierData;

  public OpenProtocolMessageIdentifierDownloadRequestRev1(
      Long midRevision,
      Short noAckFlag,
      Integer targetStationId,
      Integer targetSpindleId,
      Integer sequenceNumber,
      Short numberOfMessageParts,
      Short messagePartNumber,
      String identifierData) {
    super(
        midRevision,
        noAckFlag,
        targetStationId,
        targetSpindleId,
        sequenceNumber,
        numberOfMessageParts,
        messagePartNumber);
    this.identifierData = identifierData;
  }

  public String getIdentifierData() {
    return identifierData;
  }

  @Override
  protected void serializeOpenProtocolMessageIdentifierDownloadRequestChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("OpenProtocolMessageIdentifierDownloadRequestRev1");

    // Simple Field (identifierData)
    writeSimpleField(
        "identifierData",
        identifierData,
        writeString(writeBuffer, 800),
        WithOption.WithEncoding("ASCII"));

    writeBuffer.popContext("OpenProtocolMessageIdentifierDownloadRequestRev1");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    OpenProtocolMessageIdentifierDownloadRequestRev1 _value = this;

    // Simple field (identifierData)
    lengthInBits += 800;

    return lengthInBits;
  }

  public static OpenProtocolMessageIdentifierDownloadRequestBuilder
      staticParseOpenProtocolMessageIdentifierDownloadRequestBuilder(
          ReadBuffer readBuffer, Long revision) throws ParseException {
    readBuffer.pullContext("OpenProtocolMessageIdentifierDownloadRequestRev1");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    String identifierData =
        readSimpleField(
            "identifierData", readString(readBuffer, 800), WithOption.WithEncoding("ASCII"));

    readBuffer.closeContext("OpenProtocolMessageIdentifierDownloadRequestRev1");
    // Create the instance
    return new OpenProtocolMessageIdentifierDownloadRequestRev1BuilderImpl(identifierData);
  }

  public static class OpenProtocolMessageIdentifierDownloadRequestRev1BuilderImpl
      implements OpenProtocolMessageIdentifierDownloadRequest
          .OpenProtocolMessageIdentifierDownloadRequestBuilder {
    private final String identifierData;

    public OpenProtocolMessageIdentifierDownloadRequestRev1BuilderImpl(String identifierData) {
      this.identifierData = identifierData;
    }

    public OpenProtocolMessageIdentifierDownloadRequestRev1 build(
        Long midRevision,
        Short noAckFlag,
        Integer targetStationId,
        Integer targetSpindleId,
        Integer sequenceNumber,
        Short numberOfMessageParts,
        Short messagePartNumber) {
      OpenProtocolMessageIdentifierDownloadRequestRev1
          openProtocolMessageIdentifierDownloadRequestRev1 =
              new OpenProtocolMessageIdentifierDownloadRequestRev1(
                  midRevision,
                  noAckFlag,
                  targetStationId,
                  targetSpindleId,
                  sequenceNumber,
                  numberOfMessageParts,
                  messagePartNumber,
                  identifierData);
      return openProtocolMessageIdentifierDownloadRequestRev1;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof OpenProtocolMessageIdentifierDownloadRequestRev1)) {
      return false;
    }
    OpenProtocolMessageIdentifierDownloadRequestRev1 that =
        (OpenProtocolMessageIdentifierDownloadRequestRev1) o;
    return (getIdentifierData() == that.getIdentifierData()) && super.equals(that) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getIdentifierData());
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
