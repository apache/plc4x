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

public class AdsDeviceNotificationRequest extends AmsPacket implements Message {

  // Accessors for discriminator values.
  public CommandId getCommandId() {
    return CommandId.ADS_DEVICE_NOTIFICATION;
  }

  public Boolean getResponse() {
    return (boolean) false;
  }

  // Properties.
  protected final long length;
  protected final long stamps;
  protected final List<AdsStampHeader> adsStampHeaders;

  public AdsDeviceNotificationRequest(
      AmsNetId targetAmsNetId,
      int targetAmsPort,
      AmsNetId sourceAmsNetId,
      int sourceAmsPort,
      long errorCode,
      long invokeId,
      long length,
      long stamps,
      List<AdsStampHeader> adsStampHeaders) {
    super(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, errorCode, invokeId);
    this.length = length;
    this.stamps = stamps;
    this.adsStampHeaders = adsStampHeaders;
  }

  public long getLength() {
    return length;
  }

  public long getStamps() {
    return stamps;
  }

  public List<AdsStampHeader> getAdsStampHeaders() {
    return adsStampHeaders;
  }

  @Override
  protected void serializeAmsPacketChild(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("AdsDeviceNotificationRequest");

    // Simple Field (length)
    writeSimpleField("length", length, writeUnsignedLong(writeBuffer, 32));

    // Simple Field (stamps)
    writeSimpleField("stamps", stamps, writeUnsignedLong(writeBuffer, 32));

    // Array Field (adsStampHeaders)
    writeComplexTypeArrayField("adsStampHeaders", adsStampHeaders, writeBuffer);

    writeBuffer.popContext("AdsDeviceNotificationRequest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    AdsDeviceNotificationRequest _value = this;

    // Simple field (length)
    lengthInBits += 32;

    // Simple field (stamps)
    lengthInBits += 32;

    // Array field
    if (adsStampHeaders != null) {
      int i = 0;
      for (AdsStampHeader element : adsStampHeaders) {
        boolean last = ++i >= adsStampHeaders.size();
        lengthInBits += element.getLengthInBits();
      }
    }

    return lengthInBits;
  }

  public static AmsPacketBuilder staticParseAmsPacketBuilder(ReadBuffer readBuffer)
      throws ParseException {
    readBuffer.pullContext("AdsDeviceNotificationRequest");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    long length = readSimpleField("length", readUnsignedLong(readBuffer, 32));

    long stamps = readSimpleField("stamps", readUnsignedLong(readBuffer, 32));

    List<AdsStampHeader> adsStampHeaders =
        readCountArrayField(
            "adsStampHeaders",
            new DataReaderComplexDefault<>(
                () -> AdsStampHeader.staticParse(readBuffer), readBuffer),
            stamps);

    readBuffer.closeContext("AdsDeviceNotificationRequest");
    // Create the instance
    return new AdsDeviceNotificationRequestBuilderImpl(length, stamps, adsStampHeaders);
  }

  public static class AdsDeviceNotificationRequestBuilderImpl
      implements AmsPacket.AmsPacketBuilder {
    private final long length;
    private final long stamps;
    private final List<AdsStampHeader> adsStampHeaders;

    public AdsDeviceNotificationRequestBuilderImpl(
        long length, long stamps, List<AdsStampHeader> adsStampHeaders) {
      this.length = length;
      this.stamps = stamps;
      this.adsStampHeaders = adsStampHeaders;
    }

    public AdsDeviceNotificationRequest build(
        AmsNetId targetAmsNetId,
        int targetAmsPort,
        AmsNetId sourceAmsNetId,
        int sourceAmsPort,
        long errorCode,
        long invokeId) {
      AdsDeviceNotificationRequest adsDeviceNotificationRequest =
          new AdsDeviceNotificationRequest(
              targetAmsNetId,
              targetAmsPort,
              sourceAmsNetId,
              sourceAmsPort,
              errorCode,
              invokeId,
              length,
              stamps,
              adsStampHeaders);
      return adsDeviceNotificationRequest;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof AdsDeviceNotificationRequest)) {
      return false;
    }
    AdsDeviceNotificationRequest that = (AdsDeviceNotificationRequest) o;
    return (getLength() == that.getLength())
        && (getStamps() == that.getStamps())
        && (getAdsStampHeaders() == that.getAdsStampHeaders())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getLength(), getStamps(), getAdsStampHeaders());
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
