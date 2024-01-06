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
package org.apache.plc4x.java.knxnetip.readwrite;

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

public class TunnelingRequest extends KnxNetIpMessage implements Message {

  // Accessors for discriminator values.
  public Integer getMsgType() {
    return (int) 0x0420;
  }

  // Properties.
  protected final TunnelingRequestDataBlock tunnelingRequestDataBlock;
  protected final CEMI cemi;

  public TunnelingRequest(TunnelingRequestDataBlock tunnelingRequestDataBlock, CEMI cemi) {
    super();
    this.tunnelingRequestDataBlock = tunnelingRequestDataBlock;
    this.cemi = cemi;
  }

  public TunnelingRequestDataBlock getTunnelingRequestDataBlock() {
    return tunnelingRequestDataBlock;
  }

  public CEMI getCemi() {
    return cemi;
  }

  @Override
  protected void serializeKnxNetIpMessageChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("TunnelingRequest");

    // Simple Field (tunnelingRequestDataBlock)
    writeSimpleField(
        "tunnelingRequestDataBlock",
        tunnelingRequestDataBlock,
        new DataWriterComplexDefault<>(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (cemi)
    writeSimpleField(
        "cemi",
        cemi,
        new DataWriterComplexDefault<>(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    writeBuffer.popContext("TunnelingRequest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    TunnelingRequest _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (tunnelingRequestDataBlock)
    lengthInBits += tunnelingRequestDataBlock.getLengthInBits();

    // Simple field (cemi)
    lengthInBits += cemi.getLengthInBits();

    return lengthInBits;
  }

  public static KnxNetIpMessageBuilder staticParseKnxNetIpMessageBuilder(
      ReadBuffer readBuffer, Integer totalLength) throws ParseException {
    readBuffer.pullContext("TunnelingRequest");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    TunnelingRequestDataBlock tunnelingRequestDataBlock =
        readSimpleField(
            "tunnelingRequestDataBlock",
            new DataReaderComplexDefault<>(
                () -> TunnelingRequestDataBlock.staticParse(readBuffer), readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    CEMI cemi =
        readSimpleField(
            "cemi",
            new DataReaderComplexDefault<>(
                () ->
                    CEMI.staticParse(
                        readBuffer,
                        (int)
                            ((totalLength)
                                - (((6) + (tunnelingRequestDataBlock.getLengthInBytes()))))),
                readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    readBuffer.closeContext("TunnelingRequest");
    // Create the instance
    return new TunnelingRequestBuilderImpl(tunnelingRequestDataBlock, cemi);
  }

  public static class TunnelingRequestBuilderImpl
      implements KnxNetIpMessage.KnxNetIpMessageBuilder {
    private final TunnelingRequestDataBlock tunnelingRequestDataBlock;
    private final CEMI cemi;

    public TunnelingRequestBuilderImpl(
        TunnelingRequestDataBlock tunnelingRequestDataBlock, CEMI cemi) {
      this.tunnelingRequestDataBlock = tunnelingRequestDataBlock;
      this.cemi = cemi;
    }

    public TunnelingRequest build() {
      TunnelingRequest tunnelingRequest = new TunnelingRequest(tunnelingRequestDataBlock, cemi);
      return tunnelingRequest;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof TunnelingRequest)) {
      return false;
    }
    TunnelingRequest that = (TunnelingRequest) o;
    return (getTunnelingRequestDataBlock() == that.getTunnelingRequestDataBlock())
        && (getCemi() == that.getCemi())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getTunnelingRequestDataBlock(), getCemi());
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
