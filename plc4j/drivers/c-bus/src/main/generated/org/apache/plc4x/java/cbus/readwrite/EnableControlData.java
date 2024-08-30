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

public class EnableControlData implements Message {

  // Properties.
  protected final EnableControlCommandTypeContainer commandTypeContainer;
  protected final byte enableNetworkVariable;
  protected final byte value;

  public EnableControlData(
      EnableControlCommandTypeContainer commandTypeContainer,
      byte enableNetworkVariable,
      byte value) {
    super();
    this.commandTypeContainer = commandTypeContainer;
    this.enableNetworkVariable = enableNetworkVariable;
    this.value = value;
  }

  public EnableControlCommandTypeContainer getCommandTypeContainer() {
    return commandTypeContainer;
  }

  public byte getEnableNetworkVariable() {
    return enableNetworkVariable;
  }

  public byte getValue() {
    return value;
  }

  public EnableControlCommandType getCommandType() {
    return (EnableControlCommandType) (getCommandTypeContainer().getCommandType());
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("EnableControlData");

    // Simple Field (commandTypeContainer)
    writeSimpleEnumField(
        "commandTypeContainer",
        "EnableControlCommandTypeContainer",
        commandTypeContainer,
        new DataWriterEnumDefault<>(
            EnableControlCommandTypeContainer::getValue,
            EnableControlCommandTypeContainer::name,
            writeUnsignedShort(writeBuffer, 8)));

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    EnableControlCommandType commandType = getCommandType();
    writeBuffer.writeVirtual("commandType", commandType);

    // Simple Field (enableNetworkVariable)
    writeSimpleField("enableNetworkVariable", enableNetworkVariable, writeByte(writeBuffer, 8));

    // Simple Field (value)
    writeSimpleField("value", value, writeByte(writeBuffer, 8));

    writeBuffer.popContext("EnableControlData");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    EnableControlData _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (commandTypeContainer)
    lengthInBits += 8;

    // A virtual field doesn't have any in- or output.

    // Simple field (enableNetworkVariable)
    lengthInBits += 8;

    // Simple field (value)
    lengthInBits += 8;

    return lengthInBits;
  }

  public static EnableControlData staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("EnableControlData");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    // Validation
    if (!(org.apache.plc4x.java.cbus.readwrite.utils.StaticHelper
        .knowsEnableControlCommandTypeContainer(readBuffer))) {
      throw new ParseAssertException("no command type could be found");
    }

    EnableControlCommandTypeContainer commandTypeContainer =
        readEnumField(
            "commandTypeContainer",
            "EnableControlCommandTypeContainer",
            readEnum(
                EnableControlCommandTypeContainer::enumForValue, readUnsignedShort(readBuffer, 8)));
    EnableControlCommandType commandType =
        readVirtualField(
            "commandType", EnableControlCommandType.class, commandTypeContainer.getCommandType());

    byte enableNetworkVariable = readSimpleField("enableNetworkVariable", readByte(readBuffer, 8));

    byte value = readSimpleField("value", readByte(readBuffer, 8));

    readBuffer.closeContext("EnableControlData");
    // Create the instance
    EnableControlData _enableControlData;
    _enableControlData = new EnableControlData(commandTypeContainer, enableNetworkVariable, value);
    return _enableControlData;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof EnableControlData)) {
      return false;
    }
    EnableControlData that = (EnableControlData) o;
    return (getCommandTypeContainer() == that.getCommandTypeContainer())
        && (getEnableNetworkVariable() == that.getEnableNetworkVariable())
        && (getValue() == that.getValue())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getCommandTypeContainer(), getEnableNetworkVariable(), getValue());
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
