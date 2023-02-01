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
package org.apache.plc4x.java.bacnetip.readwrite;

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

public class NLMSecurityPayload extends NLM implements Message {

  // Accessors for discriminator values.
  public Short getMessageType() {
    return (short) 0x0B;
  }

  // Properties.
  protected final int payloadLength;
  protected final byte[] payload;

  // Arguments.
  protected final Integer apduLength;

  public NLMSecurityPayload(int payloadLength, byte[] payload, Integer apduLength) {
    super(apduLength);
    this.payloadLength = payloadLength;
    this.payload = payload;
    this.apduLength = apduLength;
  }

  public int getPayloadLength() {
    return payloadLength;
  }

  public byte[] getPayload() {
    return payload;
  }

  @Override
  protected void serializeNLMChild(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("NLMSecurityPayload");

    // Simple Field (payloadLength)
    writeSimpleField("payloadLength", payloadLength, writeUnsignedInt(writeBuffer, 16));

    // Array Field (payload)
    writeByteArrayField("payload", payload, writeByteArray(writeBuffer, 8));

    writeBuffer.popContext("NLMSecurityPayload");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    NLMSecurityPayload _value = this;

    // Simple field (payloadLength)
    lengthInBits += 16;

    // Array field
    if (payload != null) {
      lengthInBits += 8 * payload.length;
    }

    return lengthInBits;
  }

  public static NLMBuilder staticParseNLMBuilder(ReadBuffer readBuffer, Integer apduLength)
      throws ParseException {
    readBuffer.pullContext("NLMSecurityPayload");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    int payloadLength = readSimpleField("payloadLength", readUnsignedInt(readBuffer, 16));

    byte[] payload = readBuffer.readByteArray("payload", Math.toIntExact(payloadLength));

    readBuffer.closeContext("NLMSecurityPayload");
    // Create the instance
    return new NLMSecurityPayloadBuilderImpl(payloadLength, payload, apduLength);
  }

  public static class NLMSecurityPayloadBuilderImpl implements NLM.NLMBuilder {
    private final int payloadLength;
    private final byte[] payload;
    private final Integer apduLength;

    public NLMSecurityPayloadBuilderImpl(int payloadLength, byte[] payload, Integer apduLength) {
      this.payloadLength = payloadLength;
      this.payload = payload;
      this.apduLength = apduLength;
    }

    public NLMSecurityPayload build(Integer apduLength) {

      NLMSecurityPayload nLMSecurityPayload =
          new NLMSecurityPayload(payloadLength, payload, apduLength);
      return nLMSecurityPayload;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof NLMSecurityPayload)) {
      return false;
    }
    NLMSecurityPayload that = (NLMSecurityPayload) o;
    return (getPayloadLength() == that.getPayloadLength())
        && (getPayload() == that.getPayload())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getPayloadLength(), getPayload());
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
