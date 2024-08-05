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
package org.apache.plc4x.java.opcua.readwrite;

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

public class AggregateConfiguration extends ExtensionObjectDefinition implements Message {

  // Accessors for discriminator values.
  public Integer getExtensionId() {
    return (int) 950;
  }

  // Properties.
  protected final boolean treatUncertainAsBad;
  protected final boolean useServerCapabilitiesDefaults;
  protected final short percentDataBad;
  protected final short percentDataGood;
  protected final boolean useSlopedExtrapolation;

  public AggregateConfiguration(
      boolean treatUncertainAsBad,
      boolean useServerCapabilitiesDefaults,
      short percentDataBad,
      short percentDataGood,
      boolean useSlopedExtrapolation) {
    super();
    this.treatUncertainAsBad = treatUncertainAsBad;
    this.useServerCapabilitiesDefaults = useServerCapabilitiesDefaults;
    this.percentDataBad = percentDataBad;
    this.percentDataGood = percentDataGood;
    this.useSlopedExtrapolation = useSlopedExtrapolation;
  }

  public boolean getTreatUncertainAsBad() {
    return treatUncertainAsBad;
  }

  public boolean getUseServerCapabilitiesDefaults() {
    return useServerCapabilitiesDefaults;
  }

  public short getPercentDataBad() {
    return percentDataBad;
  }

  public short getPercentDataGood() {
    return percentDataGood;
  }

  public boolean getUseSlopedExtrapolation() {
    return useSlopedExtrapolation;
  }

  @Override
  protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("AggregateConfiguration");

    // Reserved Field (reserved)
    writeReservedField("reserved", (byte) 0x00, writeUnsignedByte(writeBuffer, 6));

    // Simple Field (treatUncertainAsBad)
    writeSimpleField("treatUncertainAsBad", treatUncertainAsBad, writeBoolean(writeBuffer));

    // Simple Field (useServerCapabilitiesDefaults)
    writeSimpleField(
        "useServerCapabilitiesDefaults", useServerCapabilitiesDefaults, writeBoolean(writeBuffer));

    // Simple Field (percentDataBad)
    writeSimpleField("percentDataBad", percentDataBad, writeUnsignedShort(writeBuffer, 8));

    // Simple Field (percentDataGood)
    writeSimpleField("percentDataGood", percentDataGood, writeUnsignedShort(writeBuffer, 8));

    // Reserved Field (reserved)
    writeReservedField("reserved", (byte) 0x00, writeUnsignedByte(writeBuffer, 7));

    // Simple Field (useSlopedExtrapolation)
    writeSimpleField("useSlopedExtrapolation", useSlopedExtrapolation, writeBoolean(writeBuffer));

    writeBuffer.popContext("AggregateConfiguration");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    AggregateConfiguration _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Reserved Field (reserved)
    lengthInBits += 6;

    // Simple field (treatUncertainAsBad)
    lengthInBits += 1;

    // Simple field (useServerCapabilitiesDefaults)
    lengthInBits += 1;

    // Simple field (percentDataBad)
    lengthInBits += 8;

    // Simple field (percentDataGood)
    lengthInBits += 8;

    // Reserved Field (reserved)
    lengthInBits += 7;

    // Simple field (useSlopedExtrapolation)
    lengthInBits += 1;

    return lengthInBits;
  }

  public static ExtensionObjectDefinitionBuilder staticParseExtensionObjectDefinitionBuilder(
      ReadBuffer readBuffer, Integer extensionId) throws ParseException {
    readBuffer.pullContext("AggregateConfiguration");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    Byte reservedField0 =
        readReservedField("reserved", readUnsignedByte(readBuffer, 6), (byte) 0x00);

    boolean treatUncertainAsBad = readSimpleField("treatUncertainAsBad", readBoolean(readBuffer));

    boolean useServerCapabilitiesDefaults =
        readSimpleField("useServerCapabilitiesDefaults", readBoolean(readBuffer));

    short percentDataBad = readSimpleField("percentDataBad", readUnsignedShort(readBuffer, 8));

    short percentDataGood = readSimpleField("percentDataGood", readUnsignedShort(readBuffer, 8));

    Byte reservedField1 =
        readReservedField("reserved", readUnsignedByte(readBuffer, 7), (byte) 0x00);

    boolean useSlopedExtrapolation =
        readSimpleField("useSlopedExtrapolation", readBoolean(readBuffer));

    readBuffer.closeContext("AggregateConfiguration");
    // Create the instance
    return new AggregateConfigurationBuilderImpl(
        treatUncertainAsBad,
        useServerCapabilitiesDefaults,
        percentDataBad,
        percentDataGood,
        useSlopedExtrapolation);
  }

  public static class AggregateConfigurationBuilderImpl
      implements ExtensionObjectDefinition.ExtensionObjectDefinitionBuilder {
    private final boolean treatUncertainAsBad;
    private final boolean useServerCapabilitiesDefaults;
    private final short percentDataBad;
    private final short percentDataGood;
    private final boolean useSlopedExtrapolation;

    public AggregateConfigurationBuilderImpl(
        boolean treatUncertainAsBad,
        boolean useServerCapabilitiesDefaults,
        short percentDataBad,
        short percentDataGood,
        boolean useSlopedExtrapolation) {
      this.treatUncertainAsBad = treatUncertainAsBad;
      this.useServerCapabilitiesDefaults = useServerCapabilitiesDefaults;
      this.percentDataBad = percentDataBad;
      this.percentDataGood = percentDataGood;
      this.useSlopedExtrapolation = useSlopedExtrapolation;
    }

    public AggregateConfiguration build() {
      AggregateConfiguration aggregateConfiguration =
          new AggregateConfiguration(
              treatUncertainAsBad,
              useServerCapabilitiesDefaults,
              percentDataBad,
              percentDataGood,
              useSlopedExtrapolation);
      return aggregateConfiguration;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof AggregateConfiguration)) {
      return false;
    }
    AggregateConfiguration that = (AggregateConfiguration) o;
    return (getTreatUncertainAsBad() == that.getTreatUncertainAsBad())
        && (getUseServerCapabilitiesDefaults() == that.getUseServerCapabilitiesDefaults())
        && (getPercentDataBad() == that.getPercentDataBad())
        && (getPercentDataGood() == that.getPercentDataGood())
        && (getUseSlopedExtrapolation() == that.getUseSlopedExtrapolation())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getTreatUncertainAsBad(),
        getUseServerCapabilitiesDefaults(),
        getPercentDataBad(),
        getPercentDataGood(),
        getUseSlopedExtrapolation());
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
