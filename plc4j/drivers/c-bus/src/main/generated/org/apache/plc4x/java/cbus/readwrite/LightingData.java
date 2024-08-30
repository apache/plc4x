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

public abstract class LightingData implements Message {

  // Abstract accessors for discriminator values.

  // Properties.
  protected final LightingCommandTypeContainer commandTypeContainer;

  public LightingData(LightingCommandTypeContainer commandTypeContainer) {
    super();
    this.commandTypeContainer = commandTypeContainer;
  }

  public LightingCommandTypeContainer getCommandTypeContainer() {
    return commandTypeContainer;
  }

  public LightingCommandType getCommandType() {
    return (LightingCommandType) (getCommandTypeContainer().getCommandType());
  }

  protected abstract void serializeLightingDataChild(WriteBuffer writeBuffer)
      throws SerializationException;

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("LightingData");

    // Simple Field (commandTypeContainer)
    writeSimpleEnumField(
        "commandTypeContainer",
        "LightingCommandTypeContainer",
        commandTypeContainer,
        new DataWriterEnumDefault<>(
            LightingCommandTypeContainer::getValue,
            LightingCommandTypeContainer::name,
            writeUnsignedShort(writeBuffer, 8)));

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    LightingCommandType commandType = getCommandType();
    writeBuffer.writeVirtual("commandType", commandType);

    // Switch field (Serialize the sub-type)
    serializeLightingDataChild(writeBuffer);

    writeBuffer.popContext("LightingData");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    LightingData _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (commandTypeContainer)
    lengthInBits += 8;

    // A virtual field doesn't have any in- or output.

    // Length of sub-type elements will be added by sub-type...

    return lengthInBits;
  }

  public static LightingData staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("LightingData");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    // Validation
    if (!(org.apache.plc4x.java.cbus.readwrite.utils.StaticHelper.knowsLightingCommandTypeContainer(
        readBuffer))) {
      throw new ParseAssertException("no command type could be found");
    }

    LightingCommandTypeContainer commandTypeContainer =
        readEnumField(
            "commandTypeContainer",
            "LightingCommandTypeContainer",
            readEnum(LightingCommandTypeContainer::enumForValue, readUnsignedShort(readBuffer, 8)));
    LightingCommandType commandType =
        readVirtualField(
            "commandType", LightingCommandType.class, commandTypeContainer.getCommandType());

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    LightingDataBuilder builder = null;
    if (EvaluationHelper.equals(commandType, LightingCommandType.OFF)) {
      builder = LightingDataOff.staticParseLightingDataBuilder(readBuffer);
    } else if (EvaluationHelper.equals(commandType, LightingCommandType.ON)) {
      builder = LightingDataOn.staticParseLightingDataBuilder(readBuffer);
    } else if (EvaluationHelper.equals(commandType, LightingCommandType.RAMP_TO_LEVEL)) {
      builder = LightingDataRampToLevel.staticParseLightingDataBuilder(readBuffer);
    } else if (EvaluationHelper.equals(commandType, LightingCommandType.TERMINATE_RAMP)) {
      builder = LightingDataTerminateRamp.staticParseLightingDataBuilder(readBuffer);
    } else if (EvaluationHelper.equals(commandType, LightingCommandType.LABEL)) {
      builder = LightingDataLabel.staticParseLightingDataBuilder(readBuffer, commandTypeContainer);
    }
    if (builder == null) {
      throw new ParseException(
          "Unsupported case for discriminated type"
              + " parameters ["
              + "commandType="
              + commandType
              + "]");
    }

    readBuffer.closeContext("LightingData");
    // Create the instance
    LightingData _lightingData = builder.build(commandTypeContainer);
    return _lightingData;
  }

  public interface LightingDataBuilder {
    LightingData build(LightingCommandTypeContainer commandTypeContainer);
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof LightingData)) {
      return false;
    }
    LightingData that = (LightingData) o;
    return (getCommandTypeContainer() == that.getCommandTypeContainer()) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getCommandTypeContainer());
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
