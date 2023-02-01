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

public class COTPParameterCalledTsap extends COTPParameter implements Message {

  // Accessors for discriminator values.
  public Short getParameterType() {
    return (short) 0xC2;
  }

  // Properties.
  protected final int tsapId;

  // Arguments.
  protected final Short rest;

  public COTPParameterCalledTsap(int tsapId, Short rest) {
    super(rest);
    this.tsapId = tsapId;
    this.rest = rest;
  }

  public int getTsapId() {
    return tsapId;
  }

  @Override
  protected void serializeCOTPParameterChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("COTPParameterCalledTsap");

    // Simple Field (tsapId)
    writeSimpleField("tsapId", tsapId, writeUnsignedInt(writeBuffer, 16));

    writeBuffer.popContext("COTPParameterCalledTsap");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    COTPParameterCalledTsap _value = this;

    // Simple field (tsapId)
    lengthInBits += 16;

    return lengthInBits;
  }

  public static COTPParameterBuilder staticParseCOTPParameterBuilder(
      ReadBuffer readBuffer, Short rest) throws ParseException {
    readBuffer.pullContext("COTPParameterCalledTsap");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    int tsapId = readSimpleField("tsapId", readUnsignedInt(readBuffer, 16));

    readBuffer.closeContext("COTPParameterCalledTsap");
    // Create the instance
    return new COTPParameterCalledTsapBuilderImpl(tsapId, rest);
  }

  public static class COTPParameterCalledTsapBuilderImpl
      implements COTPParameter.COTPParameterBuilder {
    private final int tsapId;
    private final Short rest;

    public COTPParameterCalledTsapBuilderImpl(int tsapId, Short rest) {
      this.tsapId = tsapId;
      this.rest = rest;
    }

    public COTPParameterCalledTsap build(Short rest) {

      COTPParameterCalledTsap cOTPParameterCalledTsap = new COTPParameterCalledTsap(tsapId, rest);
      return cOTPParameterCalledTsap;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof COTPParameterCalledTsap)) {
      return false;
    }
    COTPParameterCalledTsap that = (COTPParameterCalledTsap) o;
    return (getTsapId() == that.getTsapId()) && super.equals(that) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getTsapId());
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
