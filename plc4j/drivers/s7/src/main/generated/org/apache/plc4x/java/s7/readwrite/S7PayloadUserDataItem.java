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

public abstract class S7PayloadUserDataItem implements Message {

  // Abstract accessors for discriminator values.
  public abstract Byte getCpuFunctionGroup();

  public abstract Byte getCpuFunctionType();

  public abstract Short getCpuSubfunction();

  // Properties.
  protected final DataTransportErrorCode returnCode;
  protected final DataTransportSize transportSize;
  protected final int dataLength;

  public S7PayloadUserDataItem(
      DataTransportErrorCode returnCode, DataTransportSize transportSize, int dataLength) {
    super();
    this.returnCode = returnCode;
    this.transportSize = transportSize;
    this.dataLength = dataLength;
  }

  public DataTransportErrorCode getReturnCode() {
    return returnCode;
  }

  public DataTransportSize getTransportSize() {
    return transportSize;
  }

  public int getDataLength() {
    return dataLength;
  }

  protected abstract void serializeS7PayloadUserDataItemChild(WriteBuffer writeBuffer)
      throws SerializationException;

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("S7PayloadUserDataItem");

    // Simple Field (returnCode)
    writeSimpleEnumField(
        "returnCode",
        "DataTransportErrorCode",
        returnCode,
        new DataWriterEnumDefault<>(
            DataTransportErrorCode::getValue,
            DataTransportErrorCode::name,
            writeUnsignedShort(writeBuffer, 8)));

    // Simple Field (transportSize)
    writeSimpleEnumField(
        "transportSize",
        "DataTransportSize",
        transportSize,
        new DataWriterEnumDefault<>(
            DataTransportSize::getValue,
            DataTransportSize::name,
            writeUnsignedShort(writeBuffer, 8)));

    // Simple Field (dataLength)
    writeSimpleField("dataLength", dataLength, writeUnsignedInt(writeBuffer, 16));

    // Switch field (Serialize the sub-type)
    serializeS7PayloadUserDataItemChild(writeBuffer);

    writeBuffer.popContext("S7PayloadUserDataItem");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    S7PayloadUserDataItem _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (returnCode)
    lengthInBits += 8;

    // Simple field (transportSize)
    lengthInBits += 8;

    // Simple field (dataLength)
    lengthInBits += 16;

    // Length of sub-type elements will be added by sub-type...

    return lengthInBits;
  }

  public static S7PayloadUserDataItem staticParse(
      ReadBuffer readBuffer, Byte cpuFunctionGroup, Byte cpuFunctionType, Short cpuSubfunction)
      throws ParseException {
    readBuffer.pullContext("S7PayloadUserDataItem");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    DataTransportErrorCode returnCode =
        readEnumField(
            "returnCode",
            "DataTransportErrorCode",
            readEnum(DataTransportErrorCode::enumForValue, readUnsignedShort(readBuffer, 8)));

    DataTransportSize transportSize =
        readEnumField(
            "transportSize",
            "DataTransportSize",
            readEnum(DataTransportSize::enumForValue, readUnsignedShort(readBuffer, 8)));

    int dataLength = readSimpleField("dataLength", readUnsignedInt(readBuffer, 16));

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    S7PayloadUserDataItemBuilder builder = null;
    if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x02)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x00)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x01)) {
      builder =
          S7PayloadUserDataItemCyclicServicesPush.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x02)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x00)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x05)) {
      builder =
          S7PayloadUserDataItemCyclicServicesChangeDrivenPush
              .staticParseS7PayloadUserDataItemBuilder(
                  readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x02)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x04)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x01)) {
      builder =
          S7PayloadUserDataItemCyclicServicesSubscribeRequest
              .staticParseS7PayloadUserDataItemBuilder(
                  readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x02)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x04)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x04)) {
      builder =
          S7PayloadUserDataItemCyclicServicesUnsubscribeRequest
              .staticParseS7PayloadUserDataItemBuilder(
                  readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x02)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x01)) {
      builder =
          S7PayloadUserDataItemCyclicServicesSubscribeResponse
              .staticParseS7PayloadUserDataItemBuilder(
                  readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x02)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x04)) {
      builder =
          S7PayloadUserDataItemCyclicServicesUnsubscribeResponse
              .staticParseS7PayloadUserDataItemBuilder(
                  readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x02)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x05)
        && EvaluationHelper.equals(dataLength, (int) 0x00)) {
      builder =
          S7PayloadUserDataItemCyclicServicesErrorResponse.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x02)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x05)) {
      builder =
          S7PayloadUserDataItemCyclicServicesChangeDrivenSubscribeResponse
              .staticParseS7PayloadUserDataItemBuilder(
                  readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x00)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x03)) {
      builder =
          S7PayloadDiagnosticMessage.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x00)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x05)) {
      builder =
          S7PayloadAlarm8.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x00)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x06)) {
      builder =
          S7PayloadNotify.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x00)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x0c)) {
      builder =
          S7PayloadAlarmAckInd.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x00)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x11)) {
      builder =
          S7PayloadAlarmSQ.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x00)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x12)) {
      builder =
          S7PayloadAlarmS.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x00)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x13)) {
      builder =
          S7PayloadAlarmSC.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x00)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x16)) {
      builder =
          S7PayloadNotify8.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x04)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x01)
        && EvaluationHelper.equals(dataLength, (int) 0x00)) {
      builder =
          S7PayloadUserDataItemCpuFunctionReadSzlNoDataRequest
              .staticParseS7PayloadUserDataItemBuilder(
                  readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x04)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x01)) {
      builder =
          S7PayloadUserDataItemCpuFunctionReadSzlRequest.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x01)) {
      builder =
          S7PayloadUserDataItemCpuFunctionReadSzlResponse.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, dataLength, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x04)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x02)) {
      builder =
          S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest
              .staticParseS7PayloadUserDataItemBuilder(
                  readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x02)
        && EvaluationHelper.equals(dataLength, (int) 0x00)) {
      builder =
          S7PayloadUserDataItemCpuFunctionMsgSubscriptionResponse
              .staticParseS7PayloadUserDataItemBuilder(
                  readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x02)
        && EvaluationHelper.equals(dataLength, (int) 0x02)) {
      builder =
          S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse
              .staticParseS7PayloadUserDataItemBuilder(
                  readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x02)
        && EvaluationHelper.equals(dataLength, (int) 0x05)) {
      builder =
          S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse
              .staticParseS7PayloadUserDataItemBuilder(
                  readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x04)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x0b)) {
      builder =
          S7PayloadUserDataItemCpuFunctionAlarmAckRequest.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x0b)
        && EvaluationHelper.equals(dataLength, (int) 0x00)) {
      builder =
          S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse
              .staticParseS7PayloadUserDataItemBuilder(
                  readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x0b)) {
      builder =
          S7PayloadUserDataItemCpuFunctionAlarmAckResponse.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x04)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x13)) {
      builder =
          S7PayloadUserDataItemCpuFunctionAlarmQueryRequest.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x04)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x13)) {
      builder =
          S7PayloadUserDataItemCpuFunctionAlarmQueryResponse
              .staticParseS7PayloadUserDataItemBuilder(
                  readBuffer, dataLength, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x07)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x04)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x01)) {
      builder =
          S7PayloadUserDataItemClkRequest.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x07)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x01)) {
      builder =
          S7PayloadUserDataItemClkResponse.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, dataLength, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x07)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x04)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x03)) {
      builder =
          S7PayloadUserDataItemClkFRequest.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x07)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x03)) {
      builder =
          S7PayloadUserDataItemClkFResponse.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, dataLength, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x07)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x04)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x04)) {
      builder =
          S7PayloadUserDataItemClkSetRequest.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    } else if (EvaluationHelper.equals(cpuFunctionGroup, (byte) 0x07)
        && EvaluationHelper.equals(cpuFunctionType, (byte) 0x08)
        && EvaluationHelper.equals(cpuSubfunction, (short) 0x04)) {
      builder =
          S7PayloadUserDataItemClkSetResponse.staticParseS7PayloadUserDataItemBuilder(
              readBuffer, cpuFunctionGroup, cpuFunctionType, cpuSubfunction);
    }
    if (builder == null) {
      throw new ParseException(
          "Unsupported case for discriminated type"
              + " parameters ["
              + "cpuFunctionGroup="
              + cpuFunctionGroup
              + " "
              + "cpuFunctionType="
              + cpuFunctionType
              + " "
              + "cpuSubfunction="
              + cpuSubfunction
              + " "
              + "dataLength="
              + dataLength
              + "]");
    }

    readBuffer.closeContext("S7PayloadUserDataItem");
    // Create the instance
    S7PayloadUserDataItem _s7PayloadUserDataItem =
        builder.build(returnCode, transportSize, dataLength);
    return _s7PayloadUserDataItem;
  }

  public interface S7PayloadUserDataItemBuilder {
    S7PayloadUserDataItem build(
        DataTransportErrorCode returnCode, DataTransportSize transportSize, int dataLength);
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof S7PayloadUserDataItem)) {
      return false;
    }
    S7PayloadUserDataItem that = (S7PayloadUserDataItem) o;
    return (getReturnCode() == that.getReturnCode())
        && (getTransportSize() == that.getTransportSize())
        && (getDataLength() == that.getDataLength())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getReturnCode(), getTransportSize(), getDataLength());
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
