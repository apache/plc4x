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

import java.math.BigInteger;
import java.time.*;
import java.util.*;
import org.apache.plc4x.java.api.exceptions.*;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.codegen.*;
import org.apache.plc4x.java.spi.codegen.fields.*;
import org.apache.plc4x.java.spi.codegen.io.*;
import org.apache.plc4x.java.spi.generation.*;

// Code generated by code-generation. DO NOT EDIT.

public class BACnetConstructedDataValueSourceArray extends BACnetConstructedData
    implements Message {

  // Accessors for discriminator values.
  public BACnetObjectType getObjectTypeArgument() {
    return null;
  }

  public BACnetPropertyIdentifier getPropertyIdentifierArgument() {
    return BACnetPropertyIdentifier.VALUE_SOURCE_ARRAY;
  }

  // Properties.
  protected final BACnetApplicationTagUnsignedInteger numberOfDataElements;
  protected final List<BACnetValueSource> vtClassesSupported;

  // Arguments.
  protected final Short tagNumber;
  protected final BACnetTagPayloadUnsignedInteger arrayIndexArgument;

  public BACnetConstructedDataValueSourceArray(
      BACnetOpeningTag openingTag,
      BACnetTagHeader peekedTagHeader,
      BACnetClosingTag closingTag,
      BACnetApplicationTagUnsignedInteger numberOfDataElements,
      List<BACnetValueSource> vtClassesSupported,
      Short tagNumber,
      BACnetTagPayloadUnsignedInteger arrayIndexArgument) {
    super(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument);
    this.numberOfDataElements = numberOfDataElements;
    this.vtClassesSupported = vtClassesSupported;
    this.tagNumber = tagNumber;
    this.arrayIndexArgument = arrayIndexArgument;
  }

  public BACnetApplicationTagUnsignedInteger getNumberOfDataElements() {
    return numberOfDataElements;
  }

  public List<BACnetValueSource> getVtClassesSupported() {
    return vtClassesSupported;
  }

  public BigInteger getZero() {
    Object o = 0L;
    if (o instanceof BigInteger) return (BigInteger) o;
    return BigInteger.valueOf(((Number) o).longValue());
  }

  @Override
  protected void serializeBACnetConstructedDataChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("BACnetConstructedDataValueSourceArray");

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    BigInteger zero = getZero();
    writeBuffer.writeVirtual("zero", zero);

    // Optional Field (numberOfDataElements) (Can be skipped, if the value is null)
    writeOptionalField(
        "numberOfDataElements",
        numberOfDataElements,
        writeComplex(writeBuffer),
        ((arrayIndexArgument) != (null)) && ((arrayIndexArgument.getActualValue()) == (getZero())));

    // Array Field (vtClassesSupported)
    writeComplexTypeArrayField("vtClassesSupported", vtClassesSupported, writeBuffer);

    writeBuffer.popContext("BACnetConstructedDataValueSourceArray");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    BACnetConstructedDataValueSourceArray _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // A virtual field doesn't have any in- or output.

    // Optional Field (numberOfDataElements)
    if (numberOfDataElements != null) {
      lengthInBits += numberOfDataElements.getLengthInBits();
    }

    // Array field
    if (vtClassesSupported != null) {
      for (Message element : vtClassesSupported) {
        lengthInBits += element.getLengthInBits();
      }
    }

    return lengthInBits;
  }

  public static BACnetConstructedDataBuilder staticParseBACnetConstructedDataBuilder(
      ReadBuffer readBuffer,
      Short tagNumber,
      BACnetObjectType objectTypeArgument,
      BACnetPropertyIdentifier propertyIdentifierArgument,
      BACnetTagPayloadUnsignedInteger arrayIndexArgument)
      throws ParseException {
    readBuffer.pullContext("BACnetConstructedDataValueSourceArray");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    BigInteger zero = readVirtualField("zero", BigInteger.class, 0L);

    BACnetApplicationTagUnsignedInteger numberOfDataElements =
        readOptionalField(
            "numberOfDataElements",
            readComplex(
                () ->
                    (BACnetApplicationTagUnsignedInteger)
                        BACnetApplicationTag.staticParse(readBuffer),
                readBuffer),
            ((arrayIndexArgument) != (null)) && ((arrayIndexArgument.getActualValue()) == (zero)));

    List<BACnetValueSource> vtClassesSupported =
        readTerminatedArrayField(
            "vtClassesSupported",
            readComplex(() -> BACnetValueSource.staticParse(readBuffer), readBuffer),
            () ->
                ((boolean)
                    (org.apache.plc4x.java.bacnetip.readwrite.utils.StaticHelper
                        .isBACnetConstructedDataClosingTag(readBuffer, false, tagNumber))));
    // Validation
    if (!(((arrayIndexArgument) != (null)) || ((COUNT(vtClassesSupported)) == (16)))) {
      throw new ParseValidationException("vtClassesSupported should have exactly 16 values");
    }

    readBuffer.closeContext("BACnetConstructedDataValueSourceArray");
    // Create the instance
    return new BACnetConstructedDataValueSourceArrayBuilderImpl(
        numberOfDataElements, vtClassesSupported, tagNumber, arrayIndexArgument);
  }

  public static class BACnetConstructedDataValueSourceArrayBuilderImpl
      implements BACnetConstructedData.BACnetConstructedDataBuilder {
    private final BACnetApplicationTagUnsignedInteger numberOfDataElements;
    private final List<BACnetValueSource> vtClassesSupported;
    private final Short tagNumber;
    private final BACnetTagPayloadUnsignedInteger arrayIndexArgument;

    public BACnetConstructedDataValueSourceArrayBuilderImpl(
        BACnetApplicationTagUnsignedInteger numberOfDataElements,
        List<BACnetValueSource> vtClassesSupported,
        Short tagNumber,
        BACnetTagPayloadUnsignedInteger arrayIndexArgument) {
      this.numberOfDataElements = numberOfDataElements;
      this.vtClassesSupported = vtClassesSupported;
      this.tagNumber = tagNumber;
      this.arrayIndexArgument = arrayIndexArgument;
    }

    public BACnetConstructedDataValueSourceArray build(
        BACnetOpeningTag openingTag,
        BACnetTagHeader peekedTagHeader,
        BACnetClosingTag closingTag,
        Short tagNumber,
        BACnetTagPayloadUnsignedInteger arrayIndexArgument) {
      BACnetConstructedDataValueSourceArray bACnetConstructedDataValueSourceArray =
          new BACnetConstructedDataValueSourceArray(
              openingTag,
              peekedTagHeader,
              closingTag,
              numberOfDataElements,
              vtClassesSupported,
              tagNumber,
              arrayIndexArgument);
      return bACnetConstructedDataValueSourceArray;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BACnetConstructedDataValueSourceArray)) {
      return false;
    }
    BACnetConstructedDataValueSourceArray that = (BACnetConstructedDataValueSourceArray) o;
    return (getNumberOfDataElements() == that.getNumberOfDataElements())
        && (getVtClassesSupported() == that.getVtClassesSupported())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getNumberOfDataElements(), getVtClassesSupported());
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
