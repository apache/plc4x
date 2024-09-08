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

public class BACnetLogRecordLogDatumBooleanValue extends BACnetLogRecordLogDatum
    implements Message {

  // Accessors for discriminator values.

  // Properties.
  protected final BACnetContextTagBoolean booleanValue;

  // Arguments.
  protected final Short tagNumber;

  public BACnetLogRecordLogDatumBooleanValue(
      BACnetOpeningTag openingTag,
      BACnetTagHeader peekedTagHeader,
      BACnetClosingTag closingTag,
      BACnetContextTagBoolean booleanValue,
      Short tagNumber) {
    super(openingTag, peekedTagHeader, closingTag, tagNumber);
    this.booleanValue = booleanValue;
    this.tagNumber = tagNumber;
  }

  public BACnetContextTagBoolean getBooleanValue() {
    return booleanValue;
  }

  @Override
  protected void serializeBACnetLogRecordLogDatumChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("BACnetLogRecordLogDatumBooleanValue");

    // Simple Field (booleanValue)
    writeSimpleField("booleanValue", booleanValue, writeComplex(writeBuffer));

    writeBuffer.popContext("BACnetLogRecordLogDatumBooleanValue");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    BACnetLogRecordLogDatumBooleanValue _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (booleanValue)
    lengthInBits += booleanValue.getLengthInBits();

    return lengthInBits;
  }

  public static BACnetLogRecordLogDatumBuilder staticParseBACnetLogRecordLogDatumBuilder(
      ReadBuffer readBuffer, Short tagNumber) throws ParseException {
    readBuffer.pullContext("BACnetLogRecordLogDatumBooleanValue");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    BACnetContextTagBoolean booleanValue =
        readSimpleField(
            "booleanValue",
            readComplex(
                () ->
                    (BACnetContextTagBoolean)
                        BACnetContextTag.staticParse(
                            readBuffer, (short) (1), (BACnetDataType) (BACnetDataType.BOOLEAN)),
                readBuffer));

    readBuffer.closeContext("BACnetLogRecordLogDatumBooleanValue");
    // Create the instance
    return new BACnetLogRecordLogDatumBooleanValueBuilderImpl(booleanValue, tagNumber);
  }

  public static class BACnetLogRecordLogDatumBooleanValueBuilderImpl
      implements BACnetLogRecordLogDatum.BACnetLogRecordLogDatumBuilder {
    private final BACnetContextTagBoolean booleanValue;
    private final Short tagNumber;

    public BACnetLogRecordLogDatumBooleanValueBuilderImpl(
        BACnetContextTagBoolean booleanValue, Short tagNumber) {
      this.booleanValue = booleanValue;
      this.tagNumber = tagNumber;
    }

    public BACnetLogRecordLogDatumBooleanValue build(
        BACnetOpeningTag openingTag,
        BACnetTagHeader peekedTagHeader,
        BACnetClosingTag closingTag,
        Short tagNumber) {
      BACnetLogRecordLogDatumBooleanValue bACnetLogRecordLogDatumBooleanValue =
          new BACnetLogRecordLogDatumBooleanValue(
              openingTag, peekedTagHeader, closingTag, booleanValue, tagNumber);
      return bACnetLogRecordLogDatumBooleanValue;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BACnetLogRecordLogDatumBooleanValue)) {
      return false;
    }
    BACnetLogRecordLogDatumBooleanValue that = (BACnetLogRecordLogDatumBooleanValue) o;
    return (getBooleanValue() == that.getBooleanValue()) && super.equals(that) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getBooleanValue());
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
