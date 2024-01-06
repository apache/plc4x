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
package org.apache.plc4x.java.opcua.readwrite;

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

public class OpcuaMessageRequest extends MessagePDU implements Message {

  // Accessors for discriminator values.
  public String getMessageType() {
    return (String) "MSG";
  }

  public Boolean getResponse() {
    return (boolean) false;
  }

  // Properties.
  protected final String chunk;
  protected final int secureChannelId;
  protected final int secureTokenId;
  protected final int sequenceNumber;
  protected final int requestId;
  protected final byte[] message;

  public OpcuaMessageRequest(
      String chunk,
      int secureChannelId,
      int secureTokenId,
      int sequenceNumber,
      int requestId,
      byte[] message) {
    super();
    this.chunk = chunk;
    this.secureChannelId = secureChannelId;
    this.secureTokenId = secureTokenId;
    this.sequenceNumber = sequenceNumber;
    this.requestId = requestId;
    this.message = message;
  }

  public String getChunk() {
    return chunk;
  }

  public int getSecureChannelId() {
    return secureChannelId;
  }

  public int getSecureTokenId() {
    return secureTokenId;
  }

  public int getSequenceNumber() {
    return sequenceNumber;
  }

  public int getRequestId() {
    return requestId;
  }

  public byte[] getMessage() {
    return message;
  }

  @Override
  protected void serializeMessagePDUChild(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("OpcuaMessageRequest");

    // Simple Field (chunk)
    writeSimpleField("chunk", chunk, writeString(writeBuffer, 8));

    // Implicit Field (messageSize) (Used for parsing, but its value is not stored as it's
    // implicitly given by the objects content)
    int messageSize = (int) (getLengthInBytes());
    writeImplicitField("messageSize", messageSize, writeSignedInt(writeBuffer, 32));

    // Simple Field (secureChannelId)
    writeSimpleField("secureChannelId", secureChannelId, writeSignedInt(writeBuffer, 32));

    // Simple Field (secureTokenId)
    writeSimpleField("secureTokenId", secureTokenId, writeSignedInt(writeBuffer, 32));

    // Simple Field (sequenceNumber)
    writeSimpleField("sequenceNumber", sequenceNumber, writeSignedInt(writeBuffer, 32));

    // Simple Field (requestId)
    writeSimpleField("requestId", requestId, writeSignedInt(writeBuffer, 32));

    // Array Field (message)
    writeByteArrayField("message", message, writeByteArray(writeBuffer, 8));

    writeBuffer.popContext("OpcuaMessageRequest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    OpcuaMessageRequest _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (chunk)
    lengthInBits += 8;

    // Implicit Field (messageSize)
    lengthInBits += 32;

    // Simple field (secureChannelId)
    lengthInBits += 32;

    // Simple field (secureTokenId)
    lengthInBits += 32;

    // Simple field (sequenceNumber)
    lengthInBits += 32;

    // Simple field (requestId)
    lengthInBits += 32;

    // Array field
    if (message != null) {
      lengthInBits += 8 * message.length;
    }

    return lengthInBits;
  }

  public static MessagePDUBuilder staticParseMessagePDUBuilder(
      ReadBuffer readBuffer, Boolean response) throws ParseException {
    readBuffer.pullContext("OpcuaMessageRequest");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    String chunk = readSimpleField("chunk", readString(readBuffer, 8));

    int messageSize = readImplicitField("messageSize", readSignedInt(readBuffer, 32));

    int secureChannelId = readSimpleField("secureChannelId", readSignedInt(readBuffer, 32));

    int secureTokenId = readSimpleField("secureTokenId", readSignedInt(readBuffer, 32));

    int sequenceNumber = readSimpleField("sequenceNumber", readSignedInt(readBuffer, 32));

    int requestId = readSimpleField("requestId", readSignedInt(readBuffer, 32));

    byte[] message = readBuffer.readByteArray("message", Math.toIntExact((messageSize) - (24)));

    readBuffer.closeContext("OpcuaMessageRequest");
    // Create the instance
    return new OpcuaMessageRequestBuilderImpl(
        chunk, secureChannelId, secureTokenId, sequenceNumber, requestId, message);
  }

  public static class OpcuaMessageRequestBuilderImpl implements MessagePDU.MessagePDUBuilder {
    private final String chunk;
    private final int secureChannelId;
    private final int secureTokenId;
    private final int sequenceNumber;
    private final int requestId;
    private final byte[] message;

    public OpcuaMessageRequestBuilderImpl(
        String chunk,
        int secureChannelId,
        int secureTokenId,
        int sequenceNumber,
        int requestId,
        byte[] message) {
      this.chunk = chunk;
      this.secureChannelId = secureChannelId;
      this.secureTokenId = secureTokenId;
      this.sequenceNumber = sequenceNumber;
      this.requestId = requestId;
      this.message = message;
    }

    public OpcuaMessageRequest build() {
      OpcuaMessageRequest opcuaMessageRequest =
          new OpcuaMessageRequest(
              chunk, secureChannelId, secureTokenId, sequenceNumber, requestId, message);
      return opcuaMessageRequest;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof OpcuaMessageRequest)) {
      return false;
    }
    OpcuaMessageRequest that = (OpcuaMessageRequest) o;
    return (getChunk() == that.getChunk())
        && (getSecureChannelId() == that.getSecureChannelId())
        && (getSecureTokenId() == that.getSecureTokenId())
        && (getSequenceNumber() == that.getSequenceNumber())
        && (getRequestId() == that.getRequestId())
        && (getMessage() == that.getMessage())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getChunk(),
        getSecureChannelId(),
        getSecureTokenId(),
        getSequenceNumber(),
        getRequestId(),
        getMessage());
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