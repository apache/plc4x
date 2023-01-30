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
package org.apache.plc4x.java.ads.readwrite;

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

public class AdsReadDeviceInfoRequest extends AmsPacket implements Message {

  // Accessors for discriminator values.
  public CommandId getCommandId() {
    return CommandId.ADS_READ_DEVICE_INFO;
  }

  public Boolean getResponse() {
    return (boolean) false;
  }

  public AdsReadDeviceInfoRequest(
      AmsNetId targetAmsNetId,
      int targetAmsPort,
      AmsNetId sourceAmsNetId,
      int sourceAmsPort,
      long errorCode,
      long invokeId) {
    super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, errorCode, invokeId);
  }

  @Override
  protected void serializeAmsPacketChild(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("AdsReadDeviceInfoRequest");

    writeBuffer.popContext("AdsReadDeviceInfoRequest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    AdsReadDeviceInfoRequest _value = this;

    return lengthInBits;
  }

  public static AmsPacketBuilder staticParseAmsPacketBuilder(ReadBuffer readBuffer)
      throws ParseException {
    readBuffer.pullContext("AdsReadDeviceInfoRequest");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    readBuffer.closeContext("AdsReadDeviceInfoRequest");
    // Create the instance
    return new AdsReadDeviceInfoRequestBuilderImpl();
  }

  public static class AdsReadDeviceInfoRequestBuilderImpl implements AmsPacket.AmsPacketBuilder {

    public AdsReadDeviceInfoRequestBuilderImpl() {}

    public AdsReadDeviceInfoRequest build(
        AmsNetId targetAmsNetId,
        int targetAmsPort,
        AmsNetId sourceAmsNetId,
        int sourceAmsPort,
        long errorCode,
        long invokeId) {
      AdsReadDeviceInfoRequest adsReadDeviceInfoRequest =
          new AdsReadDeviceInfoRequest(
              targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, errorCode, invokeId);
      return adsReadDeviceInfoRequest;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof AdsReadDeviceInfoRequest)) {
      return false;
    }
    AdsReadDeviceInfoRequest that = (AdsReadDeviceInfoRequest) o;
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
