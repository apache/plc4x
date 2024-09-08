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
package org.apache.plc4x.java.s7.readwrite;

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

public class AlarmMessageObjectQueryType implements Message {

  // Constant values.
  public static final Short VARIABLESPEC = 0x12;

  // Properties.
  protected final short lengthDataset;
  protected final State eventState;
  protected final State ackStateGoing;
  protected final State ackStateComing;
  protected final DateAndTime timeComing;
  protected final AssociatedValueType valueComing;
  protected final DateAndTime timeGoing;
  protected final AssociatedValueType valueGoing;

  public AlarmMessageObjectQueryType(
      short lengthDataset,
      State eventState,
      State ackStateGoing,
      State ackStateComing,
      DateAndTime timeComing,
      AssociatedValueType valueComing,
      DateAndTime timeGoing,
      AssociatedValueType valueGoing) {
    super();
    this.lengthDataset = lengthDataset;
    this.eventState = eventState;
    this.ackStateGoing = ackStateGoing;
    this.ackStateComing = ackStateComing;
    this.timeComing = timeComing;
    this.valueComing = valueComing;
    this.timeGoing = timeGoing;
    this.valueGoing = valueGoing;
  }

  public short getLengthDataset() {
    return lengthDataset;
  }

  public State getEventState() {
    return eventState;
  }

  public State getAckStateGoing() {
    return ackStateGoing;
  }

  public State getAckStateComing() {
    return ackStateComing;
  }

  public DateAndTime getTimeComing() {
    return timeComing;
  }

  public AssociatedValueType getValueComing() {
    return valueComing;
  }

  public DateAndTime getTimeGoing() {
    return timeGoing;
  }

  public AssociatedValueType getValueGoing() {
    return valueGoing;
  }

  public short getVariableSpec() {
    return VARIABLESPEC;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("AlarmMessageObjectQueryType");

    // Simple Field (lengthDataset)
    writeSimpleField("lengthDataset", lengthDataset, writeUnsignedShort(writeBuffer, 8));

    // Reserved Field (reserved)
    writeReservedField("reserved", (int) 0x0000, writeUnsignedInt(writeBuffer, 16));

    // Const Field (variableSpec)
    writeConstField("variableSpec", VARIABLESPEC, writeUnsignedShort(writeBuffer, 8));

    // Simple Field (eventState)
    writeSimpleField("eventState", eventState, writeComplex(writeBuffer));

    // Simple Field (ackStateGoing)
    writeSimpleField("ackStateGoing", ackStateGoing, writeComplex(writeBuffer));

    // Simple Field (ackStateComing)
    writeSimpleField("ackStateComing", ackStateComing, writeComplex(writeBuffer));

    // Simple Field (timeComing)
    writeSimpleField("timeComing", timeComing, writeComplex(writeBuffer));

    // Simple Field (valueComing)
    writeSimpleField("valueComing", valueComing, writeComplex(writeBuffer));

    // Simple Field (timeGoing)
    writeSimpleField("timeGoing", timeGoing, writeComplex(writeBuffer));

    // Simple Field (valueGoing)
    writeSimpleField("valueGoing", valueGoing, writeComplex(writeBuffer));

    writeBuffer.popContext("AlarmMessageObjectQueryType");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    AlarmMessageObjectQueryType _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (lengthDataset)
    lengthInBits += 8;

    // Reserved Field (reserved)
    lengthInBits += 16;

    // Const Field (variableSpec)
    lengthInBits += 8;

    // Simple field (eventState)
    lengthInBits += eventState.getLengthInBits();

    // Simple field (ackStateGoing)
    lengthInBits += ackStateGoing.getLengthInBits();

    // Simple field (ackStateComing)
    lengthInBits += ackStateComing.getLengthInBits();

    // Simple field (timeComing)
    lengthInBits += timeComing.getLengthInBits();

    // Simple field (valueComing)
    lengthInBits += valueComing.getLengthInBits();

    // Simple field (timeGoing)
    lengthInBits += timeGoing.getLengthInBits();

    // Simple field (valueGoing)
    lengthInBits += valueGoing.getLengthInBits();

    return lengthInBits;
  }

  public static AlarmMessageObjectQueryType staticParse(ReadBuffer readBuffer)
      throws ParseException {
    readBuffer.pullContext("AlarmMessageObjectQueryType");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    short lengthDataset = readSimpleField("lengthDataset", readUnsignedShort(readBuffer, 8));

    Integer reservedField0 =
        readReservedField("reserved", readUnsignedInt(readBuffer, 16), (int) 0x0000);

    short variableSpec =
        readConstField(
            "variableSpec",
            readUnsignedShort(readBuffer, 8),
            AlarmMessageObjectQueryType.VARIABLESPEC);

    State eventState =
        readSimpleField("eventState", readComplex(() -> State.staticParse(readBuffer), readBuffer));

    State ackStateGoing =
        readSimpleField(
            "ackStateGoing", readComplex(() -> State.staticParse(readBuffer), readBuffer));

    State ackStateComing =
        readSimpleField(
            "ackStateComing", readComplex(() -> State.staticParse(readBuffer), readBuffer));

    DateAndTime timeComing =
        readSimpleField(
            "timeComing", readComplex(() -> DateAndTime.staticParse(readBuffer), readBuffer));

    AssociatedValueType valueComing =
        readSimpleField(
            "valueComing",
            readComplex(() -> AssociatedValueType.staticParse(readBuffer), readBuffer));

    DateAndTime timeGoing =
        readSimpleField(
            "timeGoing", readComplex(() -> DateAndTime.staticParse(readBuffer), readBuffer));

    AssociatedValueType valueGoing =
        readSimpleField(
            "valueGoing",
            readComplex(() -> AssociatedValueType.staticParse(readBuffer), readBuffer));

    readBuffer.closeContext("AlarmMessageObjectQueryType");
    // Create the instance
    AlarmMessageObjectQueryType _alarmMessageObjectQueryType;
    _alarmMessageObjectQueryType =
        new AlarmMessageObjectQueryType(
            lengthDataset,
            eventState,
            ackStateGoing,
            ackStateComing,
            timeComing,
            valueComing,
            timeGoing,
            valueGoing);
    return _alarmMessageObjectQueryType;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof AlarmMessageObjectQueryType)) {
      return false;
    }
    AlarmMessageObjectQueryType that = (AlarmMessageObjectQueryType) o;
    return (getLengthDataset() == that.getLengthDataset())
        && (getEventState() == that.getEventState())
        && (getAckStateGoing() == that.getAckStateGoing())
        && (getAckStateComing() == that.getAckStateComing())
        && (getTimeComing() == that.getTimeComing())
        && (getValueComing() == that.getValueComing())
        && (getTimeGoing() == that.getTimeGoing())
        && (getValueGoing() == that.getValueGoing())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        getLengthDataset(),
        getEventState(),
        getAckStateGoing(),
        getAckStateComing(),
        getTimeComing(),
        getValueComing(),
        getTimeGoing(),
        getValueGoing());
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
