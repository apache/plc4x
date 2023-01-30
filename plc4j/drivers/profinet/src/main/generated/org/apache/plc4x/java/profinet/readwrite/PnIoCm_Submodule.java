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
package org.apache.plc4x.java.profinet.readwrite;

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

public abstract class PnIoCm_Submodule implements Message {

  // Abstract accessors for discriminator values.
  public abstract PnIoCm_SubmoduleType getSubmoduleType();

  // Properties.
  protected final int slotNumber;
  protected final long submoduleIdentNumber;
  protected final boolean discardIoxs;
  protected final boolean reduceOutputModuleDataLength;
  protected final boolean reduceInputModuleDataLength;
  protected final boolean sharedInput;

  public PnIoCm_Submodule(
      int slotNumber,
      long submoduleIdentNumber,
      boolean discardIoxs,
      boolean reduceOutputModuleDataLength,
      boolean reduceInputModuleDataLength,
      boolean sharedInput) {
    super();
    this.slotNumber = slotNumber;
    this.submoduleIdentNumber = submoduleIdentNumber;
    this.discardIoxs = discardIoxs;
    this.reduceOutputModuleDataLength = reduceOutputModuleDataLength;
    this.reduceInputModuleDataLength = reduceInputModuleDataLength;
    this.sharedInput = sharedInput;
  }

  public int getSlotNumber() {
    return slotNumber;
  }

  public long getSubmoduleIdentNumber() {
    return submoduleIdentNumber;
  }

  public boolean getDiscardIoxs() {
    return discardIoxs;
  }

  public boolean getReduceOutputModuleDataLength() {
    return reduceOutputModuleDataLength;
  }

  public boolean getReduceInputModuleDataLength() {
    return reduceInputModuleDataLength;
  }

  public boolean getSharedInput() {
    return sharedInput;
  }

  protected abstract void serializePnIoCm_SubmoduleChild(WriteBuffer writeBuffer)
      throws SerializationException;

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("PnIoCm_Submodule");

    // Simple Field (slotNumber)
    writeSimpleField("slotNumber", slotNumber, writeUnsignedInt(writeBuffer, 16));

    // Simple Field (submoduleIdentNumber)
    writeSimpleField(
        "submoduleIdentNumber", submoduleIdentNumber, writeUnsignedLong(writeBuffer, 32));

    // Reserved Field (reserved)
    writeReservedField("reserved", (int) 0x000, writeUnsignedInt(writeBuffer, 10));

    // Simple Field (discardIoxs)
    writeSimpleField("discardIoxs", discardIoxs, writeBoolean(writeBuffer));

    // Simple Field (reduceOutputModuleDataLength)
    writeSimpleField(
        "reduceOutputModuleDataLength", reduceOutputModuleDataLength, writeBoolean(writeBuffer));

    // Simple Field (reduceInputModuleDataLength)
    writeSimpleField(
        "reduceInputModuleDataLength", reduceInputModuleDataLength, writeBoolean(writeBuffer));

    // Simple Field (sharedInput)
    writeSimpleField("sharedInput", sharedInput, writeBoolean(writeBuffer));

    // Discriminator Field (submoduleType) (Used as input to a switch field)
    writeDiscriminatorEnumField(
        "submoduleType",
        "PnIoCm_SubmoduleType",
        getSubmoduleType(),
        new DataWriterEnumDefault<>(
            PnIoCm_SubmoduleType::getValue,
            PnIoCm_SubmoduleType::name,
            writeUnsignedByte(writeBuffer, 2)));

    // Switch field (Serialize the sub-type)
    serializePnIoCm_SubmoduleChild(writeBuffer);

    writeBuffer.popContext("PnIoCm_Submodule");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    PnIoCm_Submodule _value = this;

    // Simple field (slotNumber)
    lengthInBits += 16;

    // Simple field (submoduleIdentNumber)
    lengthInBits += 32;

    // Reserved Field (reserved)
    lengthInBits += 10;

    // Simple field (discardIoxs)
    lengthInBits += 1;

    // Simple field (reduceOutputModuleDataLength)
    lengthInBits += 1;

    // Simple field (reduceInputModuleDataLength)
    lengthInBits += 1;

    // Simple field (sharedInput)
    lengthInBits += 1;

    // Discriminator Field (submoduleType)
    lengthInBits += 2;

    // Length of sub-type elements will be added by sub-type...

    return lengthInBits;
  }

  public static PnIoCm_Submodule staticParse(ReadBuffer readBuffer, Object... args)
      throws ParseException {
    PositionAware positionAware = readBuffer;
    return staticParse(readBuffer);
  }

  public static PnIoCm_Submodule staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("PnIoCm_Submodule");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    int slotNumber = readSimpleField("slotNumber", readUnsignedInt(readBuffer, 16));

    long submoduleIdentNumber =
        readSimpleField("submoduleIdentNumber", readUnsignedLong(readBuffer, 32));

    Integer reservedField0 =
        readReservedField("reserved", readUnsignedInt(readBuffer, 10), (int) 0x000);

    boolean discardIoxs = readSimpleField("discardIoxs", readBoolean(readBuffer));

    boolean reduceOutputModuleDataLength =
        readSimpleField("reduceOutputModuleDataLength", readBoolean(readBuffer));

    boolean reduceInputModuleDataLength =
        readSimpleField("reduceInputModuleDataLength", readBoolean(readBuffer));

    boolean sharedInput = readSimpleField("sharedInput", readBoolean(readBuffer));

    PnIoCm_SubmoduleType submoduleType =
        readDiscriminatorField(
            "submoduleType",
            new DataReaderEnumDefault<>(
                PnIoCm_SubmoduleType::enumForValue, readUnsignedByte(readBuffer, 2)));

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    PnIoCm_SubmoduleBuilder builder = null;
    if (EvaluationHelper.equals(submoduleType, PnIoCm_SubmoduleType.NO_INPUT_NO_OUTPUT_DATA)) {
      builder = PnIoCm_Submodule_NoInputNoOutputData.staticParsePnIoCm_SubmoduleBuilder(readBuffer);
    } else if (EvaluationHelper.equals(submoduleType, PnIoCm_SubmoduleType.INPUT_AND_OUTPUT_DATA)) {
      builder = PnIoCm_Submodule_InputAndOutputData.staticParsePnIoCm_SubmoduleBuilder(readBuffer);
    }
    if (builder == null) {
      throw new ParseException(
          "Unsupported case for discriminated type"
              + " parameters ["
              + "submoduleType="
              + submoduleType
              + "]");
    }

    readBuffer.closeContext("PnIoCm_Submodule");
    // Create the instance
    PnIoCm_Submodule _pnIoCm_Submodule =
        builder.build(
            slotNumber,
            submoduleIdentNumber,
            discardIoxs,
            reduceOutputModuleDataLength,
            reduceInputModuleDataLength,
            sharedInput);
    return _pnIoCm_Submodule;
  }

  public interface PnIoCm_SubmoduleBuilder {
    PnIoCm_Submodule build(
        int slotNumber,
        long submoduleIdentNumber,
        boolean discardIoxs,
        boolean reduceOutputModuleDataLength,
        boolean reduceInputModuleDataLength,
        boolean sharedInput);
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof PnIoCm_Submodule)) {
      return false;
    }
    PnIoCm_Submodule that = (PnIoCm_Submodule) o;
    return (getSlotNumber() == that.getSlotNumber())
        && (getSubmoduleIdentNumber() == that.getSubmoduleIdentNumber())
        && (getDiscardIoxs() == that.getDiscardIoxs())
        && (getReduceOutputModuleDataLength() == that.getReduceOutputModuleDataLength())
        && (getReduceInputModuleDataLength() == that.getReduceInputModuleDataLength())
        && (getSharedInput() == that.getSharedInput())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        getSlotNumber(),
        getSubmoduleIdentNumber(),
        getDiscardIoxs(),
        getReduceOutputModuleDataLength(),
        getReduceInputModuleDataLength(),
        getSharedInput());
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
