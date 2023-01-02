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
package org.apache.plc4x.java.eip.readwrite;

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

public class TransportType implements Message {

  // Properties.
  protected final boolean direction;
  protected final byte trigger;
  protected final byte classTransport;

  // Arguments.
  protected final IntegerEncoding order;

  public TransportType(
      boolean direction, byte trigger, byte classTransport, IntegerEncoding order) {
    super();
    this.direction = direction;
    this.trigger = trigger;
    this.classTransport = classTransport;
    this.order = order;
  }

  public boolean getDirection() {
    return direction;
  }

  public byte getTrigger() {
    return trigger;
  }

  public byte getClassTransport() {
    return classTransport;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("TransportType");

    // Simple Field (direction)
    writeSimpleField(
        "direction",
        direction,
        writeBoolean(writeBuffer),
        WithOption.WithByteOrder(
            (((order) == (IntegerEncoding.BIG_ENDIAN))
                ? ByteOrder.BIG_ENDIAN
                : ByteOrder.LITTLE_ENDIAN)));

    // Simple Field (trigger)
    writeSimpleField(
        "trigger",
        trigger,
        writeUnsignedByte(writeBuffer, 3),
        WithOption.WithByteOrder(
            (((order) == (IntegerEncoding.BIG_ENDIAN))
                ? ByteOrder.BIG_ENDIAN
                : ByteOrder.LITTLE_ENDIAN)));

    // Simple Field (classTransport)
    writeSimpleField(
        "classTransport",
        classTransport,
        writeUnsignedByte(writeBuffer, 4),
        WithOption.WithByteOrder(
            (((order) == (IntegerEncoding.BIG_ENDIAN))
                ? ByteOrder.BIG_ENDIAN
                : ByteOrder.LITTLE_ENDIAN)));

    writeBuffer.popContext("TransportType");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    TransportType _value = this;

    // Simple field (direction)
    lengthInBits += 1;

    // Simple field (trigger)
    lengthInBits += 3;

    // Simple field (classTransport)
    lengthInBits += 4;

    return lengthInBits;
  }

  public static TransportType staticParse(ReadBuffer readBuffer, Object... args)
      throws ParseException {
    PositionAware positionAware = readBuffer;
    if ((args == null) || (args.length != 1)) {
      throw new PlcRuntimeException(
          "Wrong number of arguments, expected 1, but got " + args.length);
    }
    IntegerEncoding order;
    if (args[0] instanceof IntegerEncoding) {
      order = (IntegerEncoding) args[0];
    } else if (args[0] instanceof String) {
      order = IntegerEncoding.valueOf((String) args[0]);
    } else {
      throw new PlcRuntimeException(
          "Argument 0 expected to be of type IntegerEncoding or a string which is parseable but was"
              + " "
              + args[0].getClass().getName());
    }
    return staticParse(readBuffer, order);
  }

  public static TransportType staticParse(ReadBuffer readBuffer, IntegerEncoding order)
      throws ParseException {
    readBuffer.pullContext("TransportType");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    boolean direction =
        readSimpleField(
            "direction",
            readBoolean(readBuffer),
            WithOption.WithByteOrder(
                (((order) == (IntegerEncoding.BIG_ENDIAN))
                    ? ByteOrder.BIG_ENDIAN
                    : ByteOrder.LITTLE_ENDIAN)));

    byte trigger =
        readSimpleField(
            "trigger",
            readUnsignedByte(readBuffer, 3),
            WithOption.WithByteOrder(
                (((order) == (IntegerEncoding.BIG_ENDIAN))
                    ? ByteOrder.BIG_ENDIAN
                    : ByteOrder.LITTLE_ENDIAN)));

    byte classTransport =
        readSimpleField(
            "classTransport",
            readUnsignedByte(readBuffer, 4),
            WithOption.WithByteOrder(
                (((order) == (IntegerEncoding.BIG_ENDIAN))
                    ? ByteOrder.BIG_ENDIAN
                    : ByteOrder.LITTLE_ENDIAN)));

    readBuffer.closeContext("TransportType");
    // Create the instance
    TransportType _transportType;
    _transportType = new TransportType(direction, trigger, classTransport, order);
    return _transportType;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof TransportType)) {
      return false;
    }
    TransportType that = (TransportType) o;
    return (getDirection() == that.getDirection())
        && (getTrigger() == that.getTrigger())
        && (getClassTransport() == that.getClassTransport())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getDirection(), getTrigger(), getClassTransport());
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
