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
package org.apache.plc4x.java.df1.readwrite;

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

public class DF1UnprotectedReadRequest extends DF1Command implements Message {

  // Accessors for discriminator values.
  public Short getCommandCode() {
    return (short) 0x01;
  }

  // Properties.
  protected final int address;
  protected final short size;

  public DF1UnprotectedReadRequest(short status, int transactionCounter, int address, short size) {
    super(status, transactionCounter);
    this.address = address;
    this.size = size;
  }

  public int getAddress() {
    return address;
  }

  public short getSize() {
    return size;
  }

  @Override
  protected void serializeDF1CommandChild(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("DF1UnprotectedReadRequest");

    // Simple Field (address)
    writeSimpleField("address", address, writeUnsignedInt(writeBuffer, 16));

    // Simple Field (size)
    writeSimpleField("size", size, writeUnsignedShort(writeBuffer, 8));

    writeBuffer.popContext("DF1UnprotectedReadRequest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    DF1UnprotectedReadRequest _value = this;

    // Simple field (address)
    lengthInBits += 16;

    // Simple field (size)
    lengthInBits += 8;

    return lengthInBits;
  }

  public static DF1CommandBuilder staticParseDF1CommandBuilder(ReadBuffer readBuffer)
      throws ParseException {
    readBuffer.pullContext("DF1UnprotectedReadRequest");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    int address = readSimpleField("address", readUnsignedInt(readBuffer, 16));

    short size = readSimpleField("size", readUnsignedShort(readBuffer, 8));

    readBuffer.closeContext("DF1UnprotectedReadRequest");
    // Create the instance
    return new DF1UnprotectedReadRequestBuilderImpl(address, size);
  }

  public static class DF1UnprotectedReadRequestBuilderImpl implements DF1Command.DF1CommandBuilder {
    private final int address;
    private final short size;

    public DF1UnprotectedReadRequestBuilderImpl(int address, short size) {
      this.address = address;
      this.size = size;
    }

    public DF1UnprotectedReadRequest build(short status, int transactionCounter) {
      DF1UnprotectedReadRequest dF1UnprotectedReadRequest =
          new DF1UnprotectedReadRequest(status, transactionCounter, address, size);
      return dF1UnprotectedReadRequest;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof DF1UnprotectedReadRequest)) {
      return false;
    }
    DF1UnprotectedReadRequest that = (DF1UnprotectedReadRequest) o;
    return (getAddress() == that.getAddress())
        && (getSize() == that.getSize())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getAddress(), getSize());
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
