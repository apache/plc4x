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

public class BACnetConstructedDataEventMessageTextsConfig extends BACnetConstructedData
    implements Message {

  // Accessors for discriminator values.
  public BACnetObjectType getObjectTypeArgument() {
    return null;
  }

  public BACnetPropertyIdentifier getPropertyIdentifierArgument() {
    return BACnetPropertyIdentifier.EVENT_MESSAGE_TEXTS_CONFIG;
  }

  // Properties.
  protected final BACnetApplicationTagUnsignedInteger numberOfDataElements;
  protected final List<BACnetOptionalCharacterString> eventMessageTextsConfig;

  // Arguments.
  protected final Short tagNumber;
  protected final BACnetTagPayloadUnsignedInteger arrayIndexArgument;

  public BACnetConstructedDataEventMessageTextsConfig(
      BACnetOpeningTag openingTag,
      BACnetTagHeader peekedTagHeader,
      BACnetClosingTag closingTag,
      BACnetApplicationTagUnsignedInteger numberOfDataElements,
      List<BACnetOptionalCharacterString> eventMessageTextsConfig,
      Short tagNumber,
      BACnetTagPayloadUnsignedInteger arrayIndexArgument) {
    super(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument);
    this.numberOfDataElements = numberOfDataElements;
    this.eventMessageTextsConfig = eventMessageTextsConfig;
    this.tagNumber = tagNumber;
    this.arrayIndexArgument = arrayIndexArgument;
  }

  public BACnetApplicationTagUnsignedInteger getNumberOfDataElements() {
    return numberOfDataElements;
  }

  public List<BACnetOptionalCharacterString> getEventMessageTextsConfig() {
    return eventMessageTextsConfig;
  }

  public BigInteger getZero() {
    Object o = 0L;
    if (o instanceof BigInteger) return (BigInteger) o;
    return BigInteger.valueOf(((Number) o).longValue());
  }

  public BACnetOptionalCharacterString getToOffnormalTextConfig() {
    return (BACnetOptionalCharacterString)
        ((((COUNT(getEventMessageTextsConfig())) == (3))
            ? getEventMessageTextsConfig().get(0)
            : null));
  }

  public BACnetOptionalCharacterString getToFaultTextConfig() {
    return (BACnetOptionalCharacterString)
        ((((COUNT(getEventMessageTextsConfig())) == (3))
            ? getEventMessageTextsConfig().get(1)
            : null));
  }

  public BACnetOptionalCharacterString getToNormalTextConfig() {
    return (BACnetOptionalCharacterString)
        ((((COUNT(getEventMessageTextsConfig())) == (3))
            ? getEventMessageTextsConfig().get(2)
            : null));
  }

  @Override
  protected void serializeBACnetConstructedDataChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("BACnetConstructedDataEventMessageTextsConfig");

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    BigInteger zero = getZero();
    writeBuffer.writeVirtual("zero", zero);

    // Optional Field (numberOfDataElements) (Can be skipped, if the value is null)
    writeOptionalField(
        "numberOfDataElements",
        numberOfDataElements,
        new DataWriterComplexDefault<>(writeBuffer),
        ((arrayIndexArgument) != (null)) && ((arrayIndexArgument.getActualValue()) == (getZero())));

    // Array Field (eventMessageTextsConfig)
    writeComplexTypeArrayField("eventMessageTextsConfig", eventMessageTextsConfig, writeBuffer);

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    BACnetOptionalCharacterString toOffnormalTextConfig = getToOffnormalTextConfig();
    writeBuffer.writeVirtual("toOffnormalTextConfig", toOffnormalTextConfig);

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    BACnetOptionalCharacterString toFaultTextConfig = getToFaultTextConfig();
    writeBuffer.writeVirtual("toFaultTextConfig", toFaultTextConfig);

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    BACnetOptionalCharacterString toNormalTextConfig = getToNormalTextConfig();
    writeBuffer.writeVirtual("toNormalTextConfig", toNormalTextConfig);

    writeBuffer.popContext("BACnetConstructedDataEventMessageTextsConfig");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    BACnetConstructedDataEventMessageTextsConfig _value = this;

    // A virtual field doesn't have any in- or output.

    // Optional Field (numberOfDataElements)
    if (numberOfDataElements != null) {
      lengthInBits += numberOfDataElements.getLengthInBits();
    }

    // Array field
    if (eventMessageTextsConfig != null) {
      for (Message element : eventMessageTextsConfig) {
        lengthInBits += element.getLengthInBits();
      }
    }

    // A virtual field doesn't have any in- or output.

    // A virtual field doesn't have any in- or output.

    // A virtual field doesn't have any in- or output.

    return lengthInBits;
  }

  public static BACnetConstructedDataEventMessageTextsConfigBuilder staticParseBuilder(
      ReadBuffer readBuffer,
      Short tagNumber,
      BACnetObjectType objectTypeArgument,
      BACnetPropertyIdentifier propertyIdentifierArgument,
      BACnetTagPayloadUnsignedInteger arrayIndexArgument)
      throws ParseException {
    readBuffer.pullContext("BACnetConstructedDataEventMessageTextsConfig");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;
    BigInteger zero = readVirtualField("zero", BigInteger.class, 0L);

    BACnetApplicationTagUnsignedInteger numberOfDataElements =
        readOptionalField(
            "numberOfDataElements",
            new DataReaderComplexDefault<>(
                () ->
                    (BACnetApplicationTagUnsignedInteger)
                        BACnetApplicationTag.staticParse(readBuffer),
                readBuffer),
            ((arrayIndexArgument) != (null)) && ((arrayIndexArgument.getActualValue()) == (zero)));

    List<BACnetOptionalCharacterString> eventMessageTextsConfig =
        readTerminatedArrayField(
            "eventMessageTextsConfig",
            new DataReaderComplexDefault<>(
                () -> BACnetOptionalCharacterString.staticParse(readBuffer), readBuffer),
            () ->
                ((boolean)
                    (org.apache.plc4x.java.bacnetip.readwrite.utils.StaticHelper
                        .isBACnetConstructedDataClosingTag(readBuffer, false, tagNumber))));
    BACnetOptionalCharacterString toOffnormalTextConfig =
        readVirtualField(
            "toOffnormalTextConfig",
            BACnetOptionalCharacterString.class,
            (((COUNT(eventMessageTextsConfig)) == (3)) ? eventMessageTextsConfig.get(0) : null));
    BACnetOptionalCharacterString toFaultTextConfig =
        readVirtualField(
            "toFaultTextConfig",
            BACnetOptionalCharacterString.class,
            (((COUNT(eventMessageTextsConfig)) == (3)) ? eventMessageTextsConfig.get(1) : null));
    BACnetOptionalCharacterString toNormalTextConfig =
        readVirtualField(
            "toNormalTextConfig",
            BACnetOptionalCharacterString.class,
            (((COUNT(eventMessageTextsConfig)) == (3)) ? eventMessageTextsConfig.get(2) : null));
    // Validation
    if (!(((arrayIndexArgument) != (null)) || ((COUNT(eventMessageTextsConfig)) == (3)))) {
      throw new ParseValidationException("eventMessageTextsConfig should have exactly 3 values");
    }

    readBuffer.closeContext("BACnetConstructedDataEventMessageTextsConfig");
    // Create the instance
    return new BACnetConstructedDataEventMessageTextsConfigBuilder(
        numberOfDataElements, eventMessageTextsConfig, tagNumber, arrayIndexArgument);
  }

  public static class BACnetConstructedDataEventMessageTextsConfigBuilder
      implements BACnetConstructedData.BACnetConstructedDataBuilder {
    private final BACnetApplicationTagUnsignedInteger numberOfDataElements;
    private final List<BACnetOptionalCharacterString> eventMessageTextsConfig;
    private final Short tagNumber;
    private final BACnetTagPayloadUnsignedInteger arrayIndexArgument;

    public BACnetConstructedDataEventMessageTextsConfigBuilder(
        BACnetApplicationTagUnsignedInteger numberOfDataElements,
        List<BACnetOptionalCharacterString> eventMessageTextsConfig,
        Short tagNumber,
        BACnetTagPayloadUnsignedInteger arrayIndexArgument) {

      this.numberOfDataElements = numberOfDataElements;
      this.eventMessageTextsConfig = eventMessageTextsConfig;
      this.tagNumber = tagNumber;
      this.arrayIndexArgument = arrayIndexArgument;
    }

    public BACnetConstructedDataEventMessageTextsConfig build(
        BACnetOpeningTag openingTag,
        BACnetTagHeader peekedTagHeader,
        BACnetClosingTag closingTag,
        Short tagNumber,
        BACnetTagPayloadUnsignedInteger arrayIndexArgument) {
      BACnetConstructedDataEventMessageTextsConfig bACnetConstructedDataEventMessageTextsConfig =
          new BACnetConstructedDataEventMessageTextsConfig(
              openingTag,
              peekedTagHeader,
              closingTag,
              numberOfDataElements,
              eventMessageTextsConfig,
              tagNumber,
              arrayIndexArgument);
      return bACnetConstructedDataEventMessageTextsConfig;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BACnetConstructedDataEventMessageTextsConfig)) {
      return false;
    }
    BACnetConstructedDataEventMessageTextsConfig that =
        (BACnetConstructedDataEventMessageTextsConfig) o;
    return (getNumberOfDataElements() == that.getNumberOfDataElements())
        && (getEventMessageTextsConfig() == that.getEventMessageTextsConfig())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getNumberOfDataElements(), getEventMessageTextsConfig());
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
