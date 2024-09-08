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

public class BACnetFaultParameterFaultState extends BACnetFaultParameter implements Message {

  // Accessors for discriminator values.

  // Properties.
  protected final BACnetOpeningTag openingTag;
  protected final BACnetFaultParameterFaultStateListOfFaultValues listOfFaultValues;
  protected final BACnetClosingTag closingTag;

  public BACnetFaultParameterFaultState(
      BACnetTagHeader peekedTagHeader,
      BACnetOpeningTag openingTag,
      BACnetFaultParameterFaultStateListOfFaultValues listOfFaultValues,
      BACnetClosingTag closingTag) {
    super(peekedTagHeader);
    this.openingTag = openingTag;
    this.listOfFaultValues = listOfFaultValues;
    this.closingTag = closingTag;
  }

  public BACnetOpeningTag getOpeningTag() {
    return openingTag;
  }

  public BACnetFaultParameterFaultStateListOfFaultValues getListOfFaultValues() {
    return listOfFaultValues;
  }

  public BACnetClosingTag getClosingTag() {
    return closingTag;
  }

  @Override
  protected void serializeBACnetFaultParameterChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("BACnetFaultParameterFaultState");

    // Simple Field (openingTag)
    writeSimpleField("openingTag", openingTag, writeComplex(writeBuffer));

    // Simple Field (listOfFaultValues)
    writeSimpleField("listOfFaultValues", listOfFaultValues, writeComplex(writeBuffer));

    // Simple Field (closingTag)
    writeSimpleField("closingTag", closingTag, writeComplex(writeBuffer));

    writeBuffer.popContext("BACnetFaultParameterFaultState");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    BACnetFaultParameterFaultState _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (openingTag)
    lengthInBits += openingTag.getLengthInBits();

    // Simple field (listOfFaultValues)
    lengthInBits += listOfFaultValues.getLengthInBits();

    // Simple field (closingTag)
    lengthInBits += closingTag.getLengthInBits();

    return lengthInBits;
  }

  public static BACnetFaultParameterBuilder staticParseBACnetFaultParameterBuilder(
      ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("BACnetFaultParameterFaultState");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    BACnetOpeningTag openingTag =
        readSimpleField(
            "openingTag",
            readComplex(() -> BACnetOpeningTag.staticParse(readBuffer, (short) (4)), readBuffer));

    BACnetFaultParameterFaultStateListOfFaultValues listOfFaultValues =
        readSimpleField(
            "listOfFaultValues",
            readComplex(
                () ->
                    BACnetFaultParameterFaultStateListOfFaultValues.staticParse(
                        readBuffer, (short) (0)),
                readBuffer));

    BACnetClosingTag closingTag =
        readSimpleField(
            "closingTag",
            readComplex(() -> BACnetClosingTag.staticParse(readBuffer, (short) (4)), readBuffer));

    readBuffer.closeContext("BACnetFaultParameterFaultState");
    // Create the instance
    return new BACnetFaultParameterFaultStateBuilderImpl(openingTag, listOfFaultValues, closingTag);
  }

  public static class BACnetFaultParameterFaultStateBuilderImpl
      implements BACnetFaultParameter.BACnetFaultParameterBuilder {
    private final BACnetOpeningTag openingTag;
    private final BACnetFaultParameterFaultStateListOfFaultValues listOfFaultValues;
    private final BACnetClosingTag closingTag;

    public BACnetFaultParameterFaultStateBuilderImpl(
        BACnetOpeningTag openingTag,
        BACnetFaultParameterFaultStateListOfFaultValues listOfFaultValues,
        BACnetClosingTag closingTag) {
      this.openingTag = openingTag;
      this.listOfFaultValues = listOfFaultValues;
      this.closingTag = closingTag;
    }

    public BACnetFaultParameterFaultState build(BACnetTagHeader peekedTagHeader) {
      BACnetFaultParameterFaultState bACnetFaultParameterFaultState =
          new BACnetFaultParameterFaultState(
              peekedTagHeader, openingTag, listOfFaultValues, closingTag);
      return bACnetFaultParameterFaultState;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BACnetFaultParameterFaultState)) {
      return false;
    }
    BACnetFaultParameterFaultState that = (BACnetFaultParameterFaultState) o;
    return (getOpeningTag() == that.getOpeningTag())
        && (getListOfFaultValues() == that.getListOfFaultValues())
        && (getClosingTag() == that.getClosingTag())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getOpeningTag(), getListOfFaultValues(), getClosingTag());
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
