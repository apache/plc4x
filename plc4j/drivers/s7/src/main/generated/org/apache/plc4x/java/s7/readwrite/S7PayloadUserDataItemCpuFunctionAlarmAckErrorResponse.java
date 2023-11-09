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

public class S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse extends S7PayloadUserDataItem
    implements Message {

  // Accessors for discriminator values.
  public Byte getCpuFunctionGroup() {
    return (byte) 0x04;
  }

  public Byte getCpuFunctionType() {
    return (byte) 0x08;
  }

  public Short getCpuSubfunction() {
    return (short) 0x0b;
  }

  public S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse(
      DataTransportErrorCode returnCode, DataTransportSize transportSize, int dataLength) {
    super(returnCode, transportSize, dataLength);
  }

  @Override
  protected void serializeS7PayloadUserDataItemChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse");

    writeBuffer.popContext("S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    return lengthInBits;
  }

  public static S7PayloadUserDataItemBuilder staticParseS7PayloadUserDataItemBuilder(
      ReadBuffer readBuffer, Byte cpuFunctionGroup, Byte cpuFunctionType, Short cpuSubfunction)
      throws ParseException {
    readBuffer.pullContext("S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    readBuffer.closeContext("S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse");
    // Create the instance
    return new S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponseBuilderImpl();
  }

  public static class S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponseBuilderImpl
      implements S7PayloadUserDataItem.S7PayloadUserDataItemBuilder {

    public S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponseBuilderImpl() {}

    public S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse build(
        DataTransportErrorCode returnCode, DataTransportSize transportSize, int dataLength) {
      S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse
          s7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse =
              new S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse(
                  returnCode, transportSize, dataLength);
      return s7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse)) {
      return false;
    }
    S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse that =
        (S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse) o;
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