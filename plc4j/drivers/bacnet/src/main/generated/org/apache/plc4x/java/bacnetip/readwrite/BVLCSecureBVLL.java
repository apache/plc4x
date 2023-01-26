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

public class BVLCSecureBVLL extends BVLC implements Message {

  // Accessors for discriminator values.
  public Short getBvlcFunction() {
    return (short) 0x0C;
  }

  // Properties.
  protected final byte[] securityWrapper;

  // Arguments.
  protected final Integer bvlcPayloadLength;

  public BVLCSecureBVLL(byte[] securityWrapper, Integer bvlcPayloadLength) {
    super();
    this.securityWrapper = securityWrapper;
    this.bvlcPayloadLength = bvlcPayloadLength;
  }

  public byte[] getSecurityWrapper() {
    return securityWrapper;
  }

  @Override
  protected void serializeBVLCChild(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("BVLCSecureBVLL");

    // Array Field (securityWrapper)
    writeByteArrayField(
        "securityWrapper",
        securityWrapper,
        writeByteArray(writeBuffer, 8),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    writeBuffer.popContext("BVLCSecureBVLL");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    BVLCSecureBVLL _value = this;

    // Array field
    if (securityWrapper != null) {
      lengthInBits += 8 * securityWrapper.length;
    }

    return lengthInBits;
  }

  public static BVLCSecureBVLLBuilder staticParseBuilder(
      ReadBuffer readBuffer, Integer bvlcPayloadLength) throws ParseException {
    readBuffer.pullContext("BVLCSecureBVLL");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    byte[] securityWrapper =
        readBuffer.readByteArray(
            "securityWrapper",
            Math.toIntExact(bvlcPayloadLength),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    readBuffer.closeContext("BVLCSecureBVLL");
    // Create the instance
    return new BVLCSecureBVLLBuilder(securityWrapper, bvlcPayloadLength);
  }

  public static class BVLCSecureBVLLBuilder implements BVLC.BVLCBuilder {
    private final byte[] securityWrapper;
    private final Integer bvlcPayloadLength;

    public BVLCSecureBVLLBuilder(byte[] securityWrapper, Integer bvlcPayloadLength) {

      this.securityWrapper = securityWrapper;
      this.bvlcPayloadLength = bvlcPayloadLength;
    }

    public BVLCSecureBVLL build() {
      BVLCSecureBVLL bVLCSecureBVLL = new BVLCSecureBVLL(securityWrapper, bvlcPayloadLength);
      return bVLCSecureBVLL;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BVLCSecureBVLL)) {
      return false;
    }
    BVLCSecureBVLL that = (BVLCSecureBVLL) o;
    return (getSecurityWrapper() == that.getSecurityWrapper()) && super.equals(that) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getSecurityWrapper());
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
