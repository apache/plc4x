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
package org.apache.plc4x.java.cbus.readwrite;

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

public class AirConditioningDataSetZoneHvacMode extends AirConditioningData implements Message {

  // Accessors for discriminator values.

  // Properties.
  protected final byte zoneGroup;
  protected final HVACZoneList zoneList;
  protected final HVACModeAndFlags hvacModeAndFlags;
  protected final HVACType hvacType;
  protected final HVACTemperature level;
  protected final HVACRawLevels rawLevel;
  protected final HVACAuxiliaryLevel auxLevel;

  public AirConditioningDataSetZoneHvacMode(
      AirConditioningCommandTypeContainer commandTypeContainer,
      byte zoneGroup,
      HVACZoneList zoneList,
      HVACModeAndFlags hvacModeAndFlags,
      HVACType hvacType,
      HVACTemperature level,
      HVACRawLevels rawLevel,
      HVACAuxiliaryLevel auxLevel) {
    super(commandTypeContainer);
    this.zoneGroup = zoneGroup;
    this.zoneList = zoneList;
    this.hvacModeAndFlags = hvacModeAndFlags;
    this.hvacType = hvacType;
    this.level = level;
    this.rawLevel = rawLevel;
    this.auxLevel = auxLevel;
  }

  public byte getZoneGroup() {
    return zoneGroup;
  }

  public HVACZoneList getZoneList() {
    return zoneList;
  }

  public HVACModeAndFlags getHvacModeAndFlags() {
    return hvacModeAndFlags;
  }

  public HVACType getHvacType() {
    return hvacType;
  }

  public HVACTemperature getLevel() {
    return level;
  }

  public HVACRawLevels getRawLevel() {
    return rawLevel;
  }

  public HVACAuxiliaryLevel getAuxLevel() {
    return auxLevel;
  }

  @Override
  protected void serializeAirConditioningDataChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("AirConditioningDataSetZoneHvacMode");

    // Simple Field (zoneGroup)
    writeSimpleField("zoneGroup", zoneGroup, writeByte(writeBuffer, 8));

    // Simple Field (zoneList)
    writeSimpleField("zoneList", zoneList, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (hvacModeAndFlags)
    writeSimpleField(
        "hvacModeAndFlags", hvacModeAndFlags, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (hvacType)
    writeSimpleEnumField(
        "hvacType",
        "HVACType",
        hvacType,
        new DataWriterEnumDefault<>(
            HVACType::getValue, HVACType::name, writeUnsignedShort(writeBuffer, 8)));

    // Optional Field (level) (Can be skipped, if the value is null)
    writeOptionalField(
        "level",
        level,
        new DataWriterComplexDefault<>(writeBuffer),
        getHvacModeAndFlags().getIsLevelTemperature());

    // Optional Field (rawLevel) (Can be skipped, if the value is null)
    writeOptionalField(
        "rawLevel",
        rawLevel,
        new DataWriterComplexDefault<>(writeBuffer),
        getHvacModeAndFlags().getIsLevelRaw());

    // Optional Field (auxLevel) (Can be skipped, if the value is null)
    writeOptionalField(
        "auxLevel",
        auxLevel,
        new DataWriterComplexDefault<>(writeBuffer),
        getHvacModeAndFlags().getIsAuxLevelUsed());

    writeBuffer.popContext("AirConditioningDataSetZoneHvacMode");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    AirConditioningDataSetZoneHvacMode _value = this;

    // Simple field (zoneGroup)
    lengthInBits += 8;

    // Simple field (zoneList)
    lengthInBits += zoneList.getLengthInBits();

    // Simple field (hvacModeAndFlags)
    lengthInBits += hvacModeAndFlags.getLengthInBits();

    // Simple field (hvacType)
    lengthInBits += 8;

    // Optional Field (level)
    if (level != null) {
      lengthInBits += level.getLengthInBits();
    }

    // Optional Field (rawLevel)
    if (rawLevel != null) {
      lengthInBits += rawLevel.getLengthInBits();
    }

    // Optional Field (auxLevel)
    if (auxLevel != null) {
      lengthInBits += auxLevel.getLengthInBits();
    }

    return lengthInBits;
  }

  public static AirConditioningDataBuilder staticParseAirConditioningDataBuilder(
      ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("AirConditioningDataSetZoneHvacMode");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    byte zoneGroup = readSimpleField("zoneGroup", readByte(readBuffer, 8));

    HVACZoneList zoneList =
        readSimpleField(
            "zoneList",
            new DataReaderComplexDefault<>(() -> HVACZoneList.staticParse(readBuffer), readBuffer));

    HVACModeAndFlags hvacModeAndFlags =
        readSimpleField(
            "hvacModeAndFlags",
            new DataReaderComplexDefault<>(
                () -> HVACModeAndFlags.staticParse(readBuffer), readBuffer));

    HVACType hvacType =
        readEnumField(
            "hvacType",
            "HVACType",
            new DataReaderEnumDefault<>(HVACType::enumForValue, readUnsignedShort(readBuffer, 8)));

    HVACTemperature level =
        readOptionalField(
            "level",
            new DataReaderComplexDefault<>(
                () -> HVACTemperature.staticParse(readBuffer), readBuffer),
            hvacModeAndFlags.getIsLevelTemperature());

    HVACRawLevels rawLevel =
        readOptionalField(
            "rawLevel",
            new DataReaderComplexDefault<>(() -> HVACRawLevels.staticParse(readBuffer), readBuffer),
            hvacModeAndFlags.getIsLevelRaw());

    HVACAuxiliaryLevel auxLevel =
        readOptionalField(
            "auxLevel",
            new DataReaderComplexDefault<>(
                () -> HVACAuxiliaryLevel.staticParse(readBuffer), readBuffer),
            hvacModeAndFlags.getIsAuxLevelUsed());

    readBuffer.closeContext("AirConditioningDataSetZoneHvacMode");
    // Create the instance
    return new AirConditioningDataSetZoneHvacModeBuilderImpl(
        zoneGroup, zoneList, hvacModeAndFlags, hvacType, level, rawLevel, auxLevel);
  }

  public static class AirConditioningDataSetZoneHvacModeBuilderImpl
      implements AirConditioningData.AirConditioningDataBuilder {
    private final byte zoneGroup;
    private final HVACZoneList zoneList;
    private final HVACModeAndFlags hvacModeAndFlags;
    private final HVACType hvacType;
    private final HVACTemperature level;
    private final HVACRawLevels rawLevel;
    private final HVACAuxiliaryLevel auxLevel;

    public AirConditioningDataSetZoneHvacModeBuilderImpl(
        byte zoneGroup,
        HVACZoneList zoneList,
        HVACModeAndFlags hvacModeAndFlags,
        HVACType hvacType,
        HVACTemperature level,
        HVACRawLevels rawLevel,
        HVACAuxiliaryLevel auxLevel) {
      this.zoneGroup = zoneGroup;
      this.zoneList = zoneList;
      this.hvacModeAndFlags = hvacModeAndFlags;
      this.hvacType = hvacType;
      this.level = level;
      this.rawLevel = rawLevel;
      this.auxLevel = auxLevel;
    }

    public AirConditioningDataSetZoneHvacMode build(
        AirConditioningCommandTypeContainer commandTypeContainer) {
      AirConditioningDataSetZoneHvacMode airConditioningDataSetZoneHvacMode =
          new AirConditioningDataSetZoneHvacMode(
              commandTypeContainer,
              zoneGroup,
              zoneList,
              hvacModeAndFlags,
              hvacType,
              level,
              rawLevel,
              auxLevel);
      return airConditioningDataSetZoneHvacMode;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof AirConditioningDataSetZoneHvacMode)) {
      return false;
    }
    AirConditioningDataSetZoneHvacMode that = (AirConditioningDataSetZoneHvacMode) o;
    return (getZoneGroup() == that.getZoneGroup())
        && (getZoneList() == that.getZoneList())
        && (getHvacModeAndFlags() == that.getHvacModeAndFlags())
        && (getHvacType() == that.getHvacType())
        && (getLevel() == that.getLevel())
        && (getRawLevel() == that.getRawLevel())
        && (getAuxLevel() == that.getAuxLevel())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getZoneGroup(),
        getZoneList(),
        getHvacModeAndFlags(),
        getHvacType(),
        getLevel(),
        getRawLevel(),
        getAuxLevel());
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
