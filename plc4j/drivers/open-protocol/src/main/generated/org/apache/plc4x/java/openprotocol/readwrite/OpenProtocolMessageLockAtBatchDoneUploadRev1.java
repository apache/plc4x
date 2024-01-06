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

public class OpenProtocolMessageLockAtBatchDoneUploadRev1
    extends OpenProtocolMessageLockAtBatchDoneUpload implements Message {

  // Accessors for discriminator values.
  public Integer getRevision() {
    return (int) 1;
  }

  // Properties.
  protected final short relayStatus;

  public OpenProtocolMessageLockAtBatchDoneUploadRev1(
      Integer midRevision,
      Short noAckFlag,
      Integer targetStationId,
      Integer targetSpindleId,
      Integer sequenceNumber,
      Short numberOfMessageParts,
      Short messagePartNumber,
      short relayStatus) {
    super(
        midRevision,
        noAckFlag,
        targetStationId,
        targetSpindleId,
        sequenceNumber,
        numberOfMessageParts,
        messagePartNumber);
    this.relayStatus = relayStatus;
  }

  public short getRelayStatus() {
    return relayStatus;
  }

  @Override
  protected void serializeOpenProtocolMessageLockAtBatchDoneUploadChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("OpenProtocolMessageLockAtBatchDoneUploadRev1");

    // Simple Field (relayStatus)
    writeSimpleField(
        "relayStatus",
        relayStatus,
        writeUnsignedShort(writeBuffer, 8),
        WithOption.WithEncoding("ASCII"));

    writeBuffer.popContext("OpenProtocolMessageLockAtBatchDoneUploadRev1");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    OpenProtocolMessageLockAtBatchDoneUploadRev1 _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (relayStatus)
    lengthInBits += 8;

    return lengthInBits;
  }

  public static OpenProtocolMessageLockAtBatchDoneUploadBuilder
      staticParseOpenProtocolMessageLockAtBatchDoneUploadBuilder(
          ReadBuffer readBuffer, Integer revision) throws ParseException {
    readBuffer.pullContext("OpenProtocolMessageLockAtBatchDoneUploadRev1");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    short relayStatus =
        readSimpleField(
            "relayStatus", readUnsignedShort(readBuffer, 8), WithOption.WithEncoding("ASCII"));

    readBuffer.closeContext("OpenProtocolMessageLockAtBatchDoneUploadRev1");
    // Create the instance
    return new OpenProtocolMessageLockAtBatchDoneUploadRev1BuilderImpl(relayStatus);
  }

  public static class OpenProtocolMessageLockAtBatchDoneUploadRev1BuilderImpl
      implements OpenProtocolMessageLockAtBatchDoneUpload
          .OpenProtocolMessageLockAtBatchDoneUploadBuilder {
    private final short relayStatus;

    public OpenProtocolMessageLockAtBatchDoneUploadRev1BuilderImpl(short relayStatus) {
      this.relayStatus = relayStatus;
    }

    public OpenProtocolMessageLockAtBatchDoneUploadRev1 build(
        Integer midRevision,
        Short noAckFlag,
        Integer targetStationId,
        Integer targetSpindleId,
        Integer sequenceNumber,
        Short numberOfMessageParts,
        Short messagePartNumber) {
      OpenProtocolMessageLockAtBatchDoneUploadRev1 openProtocolMessageLockAtBatchDoneUploadRev1 =
          new OpenProtocolMessageLockAtBatchDoneUploadRev1(
              midRevision,
              noAckFlag,
              targetStationId,
              targetSpindleId,
              sequenceNumber,
              numberOfMessageParts,
              messagePartNumber,
              relayStatus);
      return openProtocolMessageLockAtBatchDoneUploadRev1;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof OpenProtocolMessageLockAtBatchDoneUploadRev1)) {
      return false;
    }
    OpenProtocolMessageLockAtBatchDoneUploadRev1 that =
        (OpenProtocolMessageLockAtBatchDoneUploadRev1) o;
    return (getRelayStatus() == that.getRelayStatus()) && super.equals(that) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getRelayStatus());
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
