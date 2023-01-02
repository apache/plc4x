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

public class PnIoCm_IoCs implements Message {

  // Properties.
  protected final int slotNumber;
  protected final int subSlotNumber;
  protected final int ioFrameOffset;

  public PnIoCm_IoCs(int slotNumber, int subSlotNumber, int ioFrameOffset) {
    super();
    this.slotNumber = slotNumber;
    this.subSlotNumber = subSlotNumber;
    this.ioFrameOffset = ioFrameOffset;
  }

  public int getSlotNumber() {
    return slotNumber;
  }

  public int getSubSlotNumber() {
    return subSlotNumber;
  }

  public int getIoFrameOffset() {
    return ioFrameOffset;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("PnIoCm_IoCs");

    // Simple Field (slotNumber)
    writeSimpleField("slotNumber", slotNumber, writeUnsignedInt(writeBuffer, 16));

    // Simple Field (subSlotNumber)
    writeSimpleField("subSlotNumber", subSlotNumber, writeUnsignedInt(writeBuffer, 16));

    // Simple Field (ioFrameOffset)
    writeSimpleField("ioFrameOffset", ioFrameOffset, writeUnsignedInt(writeBuffer, 16));

    writeBuffer.popContext("PnIoCm_IoCs");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    PnIoCm_IoCs _value = this;

    // Simple field (slotNumber)
    lengthInBits += 16;

    // Simple field (subSlotNumber)
    lengthInBits += 16;

    // Simple field (ioFrameOffset)
    lengthInBits += 16;

    return lengthInBits;
  }

  public static PnIoCm_IoCs staticParse(ReadBuffer readBuffer, Object... args)
      throws ParseException {
    PositionAware positionAware = readBuffer;
    return staticParse(readBuffer);
  }

  public static PnIoCm_IoCs staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("PnIoCm_IoCs");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    int slotNumber = readSimpleField("slotNumber", readUnsignedInt(readBuffer, 16));

    int subSlotNumber = readSimpleField("subSlotNumber", readUnsignedInt(readBuffer, 16));

    int ioFrameOffset = readSimpleField("ioFrameOffset", readUnsignedInt(readBuffer, 16));

    readBuffer.closeContext("PnIoCm_IoCs");
    // Create the instance
    PnIoCm_IoCs _pnIoCm_IoCs;
    _pnIoCm_IoCs = new PnIoCm_IoCs(slotNumber, subSlotNumber, ioFrameOffset);
    return _pnIoCm_IoCs;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof PnIoCm_IoCs)) {
      return false;
    }
    PnIoCm_IoCs that = (PnIoCm_IoCs) o;
    return (getSlotNumber() == that.getSlotNumber())
        && (getSubSlotNumber() == that.getSubSlotNumber())
        && (getIoFrameOffset() == that.getIoFrameOffset())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getSlotNumber(), getSubSlotNumber(), getIoFrameOffset());
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
