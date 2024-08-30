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

public abstract class MeteringData implements Message {

  // Abstract accessors for discriminator values.

  // Properties.
  protected final MeteringCommandTypeContainer commandTypeContainer;
  protected final byte argument;

  public MeteringData(MeteringCommandTypeContainer commandTypeContainer, byte argument) {
    super();
    this.commandTypeContainer = commandTypeContainer;
    this.argument = argument;
  }

  public MeteringCommandTypeContainer getCommandTypeContainer() {
    return commandTypeContainer;
  }

  public byte getArgument() {
    return argument;
  }

  public MeteringCommandType getCommandType() {
    return (MeteringCommandType) (getCommandTypeContainer().getCommandType());
  }

  protected abstract void serializeMeteringDataChild(WriteBuffer writeBuffer)
      throws SerializationException;

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("MeteringData");

    // Simple Field (commandTypeContainer)
    writeSimpleEnumField(
        "commandTypeContainer",
        "MeteringCommandTypeContainer",
        commandTypeContainer,
        new DataWriterEnumDefault<>(
            MeteringCommandTypeContainer::getValue,
            MeteringCommandTypeContainer::name,
            writeUnsignedShort(writeBuffer, 8)));

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    MeteringCommandType commandType = getCommandType();
    writeBuffer.writeVirtual("commandType", commandType);

    // Simple Field (argument)
    writeSimpleField("argument", argument, writeByte(writeBuffer, 8));

    // Switch field (Serialize the sub-type)
    serializeMeteringDataChild(writeBuffer);

    writeBuffer.popContext("MeteringData");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    MeteringData _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (commandTypeContainer)
    lengthInBits += 8;

    // A virtual field doesn't have any in- or output.

    // Simple field (argument)
    lengthInBits += 8;

    // Length of sub-type elements will be added by sub-type...

    return lengthInBits;
  }

  public static MeteringData staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("MeteringData");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    // Validation
    if (!(org.apache.plc4x.java.cbus.readwrite.utils.StaticHelper.knowsMeteringCommandTypeContainer(
        readBuffer))) {
      throw new ParseAssertException("no command type could be found");
    }

    MeteringCommandTypeContainer commandTypeContainer =
        readEnumField(
            "commandTypeContainer",
            "MeteringCommandTypeContainer",
            readEnum(MeteringCommandTypeContainer::enumForValue, readUnsignedShort(readBuffer, 8)));
    MeteringCommandType commandType =
        readVirtualField(
            "commandType", MeteringCommandType.class, commandTypeContainer.getCommandType());

    byte argument = readSimpleField("argument", readByte(readBuffer, 8));

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    MeteringDataBuilder builder = null;
    if (EvaluationHelper.equals(commandType, MeteringCommandType.EVENT)
        && EvaluationHelper.equals(argument, (byte) 0x01)) {
      builder = MeteringDataMeasureElectricity.staticParseMeteringDataBuilder(readBuffer);
    } else if (EvaluationHelper.equals(commandType, MeteringCommandType.EVENT)
        && EvaluationHelper.equals(argument, (byte) 0x02)) {
      builder = MeteringDataMeasureGas.staticParseMeteringDataBuilder(readBuffer);
    } else if (EvaluationHelper.equals(commandType, MeteringCommandType.EVENT)
        && EvaluationHelper.equals(argument, (byte) 0x03)) {
      builder = MeteringDataMeasureDrinkingWater.staticParseMeteringDataBuilder(readBuffer);
    } else if (EvaluationHelper.equals(commandType, MeteringCommandType.EVENT)
        && EvaluationHelper.equals(argument, (byte) 0x04)) {
      builder = MeteringDataMeasureOtherWater.staticParseMeteringDataBuilder(readBuffer);
    } else if (EvaluationHelper.equals(commandType, MeteringCommandType.EVENT)
        && EvaluationHelper.equals(argument, (byte) 0x05)) {
      builder = MeteringDataMeasureOil.staticParseMeteringDataBuilder(readBuffer);
    } else if (EvaluationHelper.equals(commandType, MeteringCommandType.EVENT)
        && EvaluationHelper.equals(argument, (byte) 0x81)) {
      builder = MeteringDataElectricityConsumption.staticParseMeteringDataBuilder(readBuffer);
    } else if (EvaluationHelper.equals(commandType, MeteringCommandType.EVENT)
        && EvaluationHelper.equals(argument, (byte) 0x82)) {
      builder = MeteringDataGasConsumption.staticParseMeteringDataBuilder(readBuffer);
    } else if (EvaluationHelper.equals(commandType, MeteringCommandType.EVENT)
        && EvaluationHelper.equals(argument, (byte) 0x83)) {
      builder = MeteringDataDrinkingWaterConsumption.staticParseMeteringDataBuilder(readBuffer);
    } else if (EvaluationHelper.equals(commandType, MeteringCommandType.EVENT)
        && EvaluationHelper.equals(argument, (byte) 0x84)) {
      builder = MeteringDataOtherWaterConsumption.staticParseMeteringDataBuilder(readBuffer);
    } else if (EvaluationHelper.equals(commandType, MeteringCommandType.EVENT)
        && EvaluationHelper.equals(argument, (byte) 0x85)) {
      builder = MeteringDataOilConsumption.staticParseMeteringDataBuilder(readBuffer);
    }
    if (builder == null) {
      throw new ParseException(
          "Unsupported case for discriminated type"
              + " parameters ["
              + "commandType="
              + commandType
              + " "
              + "argument="
              + argument
              + "]");
    }

    readBuffer.closeContext("MeteringData");
    // Create the instance
    MeteringData _meteringData = builder.build(commandTypeContainer, argument);
    return _meteringData;
  }

  public interface MeteringDataBuilder {
    MeteringData build(MeteringCommandTypeContainer commandTypeContainer, byte argument);
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof MeteringData)) {
      return false;
    }
    MeteringData that = (MeteringData) o;
    return (getCommandTypeContainer() == that.getCommandTypeContainer())
        && (getArgument() == that.getArgument())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getCommandTypeContainer(), getArgument());
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
