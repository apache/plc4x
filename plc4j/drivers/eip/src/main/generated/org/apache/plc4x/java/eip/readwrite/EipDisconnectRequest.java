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
package org.apache.plc4x.java.eip.readwrite;

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

public class EipDisconnectRequest extends EipPacket implements Message {

  // Accessors for discriminator values.
  public Integer getCommand() {
    return (int) 0x0066;
  }

  public Boolean getResponse() {
    return false;
  }

  public Integer getPacketLength() {
    return 0;
  }

  // Arguments.
  protected final IntegerEncoding order;

  public EipDisconnectRequest(
      long sessionHandle, long status, byte[] senderContext, long options, IntegerEncoding order) {
    super(sessionHandle, status, senderContext, options, order);
    this.order = order;
  }

  @Override
  protected void serializeEipPacketChild(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("EipDisconnectRequest");

    writeBuffer.popContext("EipDisconnectRequest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    EipDisconnectRequest _value = this;

    return lengthInBits;
  }

  public static EipDisconnectRequestBuilder staticParseBuilder(
      ReadBuffer readBuffer, IntegerEncoding order, Boolean response) throws ParseException {
    readBuffer.pullContext("EipDisconnectRequest");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    readBuffer.closeContext("EipDisconnectRequest");
    // Create the instance
    return new EipDisconnectRequestBuilder(order);
  }

  public static class EipDisconnectRequestBuilder implements EipPacket.EipPacketBuilder {
    private final IntegerEncoding order;

    public EipDisconnectRequestBuilder(IntegerEncoding order) {

      this.order = order;
    }

    public EipDisconnectRequest build(
        long sessionHandle,
        long status,
        byte[] senderContext,
        long options,
        IntegerEncoding order) {
      EipDisconnectRequest eipDisconnectRequest =
          new EipDisconnectRequest(sessionHandle, status, senderContext, options, order);
      return eipDisconnectRequest;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof EipDisconnectRequest)) {
      return false;
    }
    EipDisconnectRequest that = (EipDisconnectRequest) o;
    return super.equals(that) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode());
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
