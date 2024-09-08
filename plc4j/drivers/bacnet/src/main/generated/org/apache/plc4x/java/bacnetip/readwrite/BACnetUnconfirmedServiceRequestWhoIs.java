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

public class BACnetUnconfirmedServiceRequestWhoIs extends BACnetUnconfirmedServiceRequest
    implements Message {

  // Accessors for discriminator values.
  public BACnetUnconfirmedServiceChoice getServiceChoice() {
    return BACnetUnconfirmedServiceChoice.WHO_IS;
  }

  // Properties.
  protected final BACnetContextTagUnsignedInteger deviceInstanceRangeLowLimit;
  protected final BACnetContextTagUnsignedInteger deviceInstanceRangeHighLimit;

  // Arguments.
  protected final Integer serviceRequestLength;

  public BACnetUnconfirmedServiceRequestWhoIs(
      BACnetContextTagUnsignedInteger deviceInstanceRangeLowLimit,
      BACnetContextTagUnsignedInteger deviceInstanceRangeHighLimit,
      Integer serviceRequestLength) {
    super(serviceRequestLength);
    this.deviceInstanceRangeLowLimit = deviceInstanceRangeLowLimit;
    this.deviceInstanceRangeHighLimit = deviceInstanceRangeHighLimit;
    this.serviceRequestLength = serviceRequestLength;
  }

  public BACnetContextTagUnsignedInteger getDeviceInstanceRangeLowLimit() {
    return deviceInstanceRangeLowLimit;
  }

  public BACnetContextTagUnsignedInteger getDeviceInstanceRangeHighLimit() {
    return deviceInstanceRangeHighLimit;
  }

  @Override
  protected void serializeBACnetUnconfirmedServiceRequestChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("BACnetUnconfirmedServiceRequestWhoIs");

    // Optional Field (deviceInstanceRangeLowLimit) (Can be skipped, if the value is null)
    writeOptionalField(
        "deviceInstanceRangeLowLimit", deviceInstanceRangeLowLimit, writeComplex(writeBuffer));

    // Optional Field (deviceInstanceRangeHighLimit) (Can be skipped, if the value is null)
    writeOptionalField(
        "deviceInstanceRangeHighLimit",
        deviceInstanceRangeHighLimit,
        writeComplex(writeBuffer),
        (getDeviceInstanceRangeLowLimit()) != (null));

    writeBuffer.popContext("BACnetUnconfirmedServiceRequestWhoIs");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    BACnetUnconfirmedServiceRequestWhoIs _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Optional Field (deviceInstanceRangeLowLimit)
    if (deviceInstanceRangeLowLimit != null) {
      lengthInBits += deviceInstanceRangeLowLimit.getLengthInBits();
    }

    // Optional Field (deviceInstanceRangeHighLimit)
    if (deviceInstanceRangeHighLimit != null) {
      lengthInBits += deviceInstanceRangeHighLimit.getLengthInBits();
    }

    return lengthInBits;
  }

  public static BACnetUnconfirmedServiceRequestBuilder
      staticParseBACnetUnconfirmedServiceRequestBuilder(
          ReadBuffer readBuffer, Integer serviceRequestLength) throws ParseException {
    readBuffer.pullContext("BACnetUnconfirmedServiceRequestWhoIs");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    BACnetContextTagUnsignedInteger deviceInstanceRangeLowLimit =
        readOptionalField(
            "deviceInstanceRangeLowLimit",
            readComplex(
                () ->
                    (BACnetContextTagUnsignedInteger)
                        BACnetContextTag.staticParse(
                            readBuffer,
                            (short) (0),
                            (BACnetDataType) (BACnetDataType.UNSIGNED_INTEGER)),
                readBuffer));

    BACnetContextTagUnsignedInteger deviceInstanceRangeHighLimit =
        readOptionalField(
            "deviceInstanceRangeHighLimit",
            readComplex(
                () ->
                    (BACnetContextTagUnsignedInteger)
                        BACnetContextTag.staticParse(
                            readBuffer,
                            (short) (1),
                            (BACnetDataType) (BACnetDataType.UNSIGNED_INTEGER)),
                readBuffer),
            (deviceInstanceRangeLowLimit) != (null));

    readBuffer.closeContext("BACnetUnconfirmedServiceRequestWhoIs");
    // Create the instance
    return new BACnetUnconfirmedServiceRequestWhoIsBuilderImpl(
        deviceInstanceRangeLowLimit, deviceInstanceRangeHighLimit, serviceRequestLength);
  }

  public static class BACnetUnconfirmedServiceRequestWhoIsBuilderImpl
      implements BACnetUnconfirmedServiceRequest.BACnetUnconfirmedServiceRequestBuilder {
    private final BACnetContextTagUnsignedInteger deviceInstanceRangeLowLimit;
    private final BACnetContextTagUnsignedInteger deviceInstanceRangeHighLimit;
    private final Integer serviceRequestLength;

    public BACnetUnconfirmedServiceRequestWhoIsBuilderImpl(
        BACnetContextTagUnsignedInteger deviceInstanceRangeLowLimit,
        BACnetContextTagUnsignedInteger deviceInstanceRangeHighLimit,
        Integer serviceRequestLength) {
      this.deviceInstanceRangeLowLimit = deviceInstanceRangeLowLimit;
      this.deviceInstanceRangeHighLimit = deviceInstanceRangeHighLimit;
      this.serviceRequestLength = serviceRequestLength;
    }

    public BACnetUnconfirmedServiceRequestWhoIs build(Integer serviceRequestLength) {

      BACnetUnconfirmedServiceRequestWhoIs bACnetUnconfirmedServiceRequestWhoIs =
          new BACnetUnconfirmedServiceRequestWhoIs(
              deviceInstanceRangeLowLimit, deviceInstanceRangeHighLimit, serviceRequestLength);
      return bACnetUnconfirmedServiceRequestWhoIs;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BACnetUnconfirmedServiceRequestWhoIs)) {
      return false;
    }
    BACnetUnconfirmedServiceRequestWhoIs that = (BACnetUnconfirmedServiceRequestWhoIs) o;
    return (getDeviceInstanceRangeLowLimit() == that.getDeviceInstanceRangeLowLimit())
        && (getDeviceInstanceRangeHighLimit() == that.getDeviceInstanceRangeHighLimit())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(), getDeviceInstanceRangeLowLimit(), getDeviceInstanceRangeHighLimit());
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
