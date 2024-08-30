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

public class TemperatureBroadcastData implements Message {

  // Properties.
  protected final TemperatureBroadcastCommandTypeContainer commandTypeContainer;
  protected final byte temperatureGroup;
  protected final byte temperatureByte;

  public TemperatureBroadcastData(
      TemperatureBroadcastCommandTypeContainer commandTypeContainer,
      byte temperatureGroup,
      byte temperatureByte) {
    super();
    this.commandTypeContainer = commandTypeContainer;
    this.temperatureGroup = temperatureGroup;
    this.temperatureByte = temperatureByte;
  }

  public TemperatureBroadcastCommandTypeContainer getCommandTypeContainer() {
    return commandTypeContainer;
  }

  public byte getTemperatureGroup() {
    return temperatureGroup;
  }

  public byte getTemperatureByte() {
    return temperatureByte;
  }

  public TemperatureBroadcastCommandType getCommandType() {
    return (TemperatureBroadcastCommandType) (getCommandTypeContainer().getCommandType());
  }

  public float getTemperatureInCelsius() {
    return (float) ((getTemperatureByte()) / (4F));
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("TemperatureBroadcastData");

    // Simple Field (commandTypeContainer)
    writeSimpleEnumField(
        "commandTypeContainer",
        "TemperatureBroadcastCommandTypeContainer",
        commandTypeContainer,
        new DataWriterEnumDefault<>(
            TemperatureBroadcastCommandTypeContainer::getValue,
            TemperatureBroadcastCommandTypeContainer::name,
            writeUnsignedShort(writeBuffer, 8)));

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    TemperatureBroadcastCommandType commandType = getCommandType();
    writeBuffer.writeVirtual("commandType", commandType);

    // Simple Field (temperatureGroup)
    writeSimpleField("temperatureGroup", temperatureGroup, writeByte(writeBuffer, 8));

    // Simple Field (temperatureByte)
    writeSimpleField("temperatureByte", temperatureByte, writeByte(writeBuffer, 8));

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    float temperatureInCelsius = getTemperatureInCelsius();
    writeBuffer.writeVirtual("temperatureInCelsius", temperatureInCelsius);

    writeBuffer.popContext("TemperatureBroadcastData");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    TemperatureBroadcastData _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (commandTypeContainer)
    lengthInBits += 8;

    // A virtual field doesn't have any in- or output.

    // Simple field (temperatureGroup)
    lengthInBits += 8;

    // Simple field (temperatureByte)
    lengthInBits += 8;

    // A virtual field doesn't have any in- or output.

    return lengthInBits;
  }

  public static TemperatureBroadcastData staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("TemperatureBroadcastData");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    // Validation
    if (!(org.apache.plc4x.java.cbus.readwrite.utils.StaticHelper
        .knowsTemperatureBroadcastCommandTypeContainer(readBuffer))) {
      throw new ParseAssertException("no command type could be found");
    }

    TemperatureBroadcastCommandTypeContainer commandTypeContainer =
        readEnumField(
            "commandTypeContainer",
            "TemperatureBroadcastCommandTypeContainer",
            readEnum(
                TemperatureBroadcastCommandTypeContainer::enumForValue,
                readUnsignedShort(readBuffer, 8)));
    TemperatureBroadcastCommandType commandType =
        readVirtualField(
            "commandType",
            TemperatureBroadcastCommandType.class,
            commandTypeContainer.getCommandType());

    byte temperatureGroup = readSimpleField("temperatureGroup", readByte(readBuffer, 8));

    byte temperatureByte = readSimpleField("temperatureByte", readByte(readBuffer, 8));
    float temperatureInCelsius =
        readVirtualField("temperatureInCelsius", float.class, (temperatureByte) / (4F));

    readBuffer.closeContext("TemperatureBroadcastData");
    // Create the instance
    TemperatureBroadcastData _temperatureBroadcastData;
    _temperatureBroadcastData =
        new TemperatureBroadcastData(commandTypeContainer, temperatureGroup, temperatureByte);
    return _temperatureBroadcastData;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof TemperatureBroadcastData)) {
      return false;
    }
    TemperatureBroadcastData that = (TemperatureBroadcastData) o;
    return (getCommandTypeContainer() == that.getCommandTypeContainer())
        && (getTemperatureGroup() == that.getTemperatureGroup())
        && (getTemperatureByte() == that.getTemperatureByte())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getCommandTypeContainer(), getTemperatureGroup(), getTemperatureByte());
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
