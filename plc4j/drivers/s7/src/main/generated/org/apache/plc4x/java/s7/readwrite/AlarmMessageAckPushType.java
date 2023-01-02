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

public class AlarmMessageAckPushType implements Message {

  // Properties.
  protected final DateAndTime TimeStamp;
  protected final short functionId;
  protected final short numberOfObjects;
  protected final List<AlarmMessageAckObjectPushType> messageObjects;

  public AlarmMessageAckPushType(
      DateAndTime TimeStamp,
      short functionId,
      short numberOfObjects,
      List<AlarmMessageAckObjectPushType> messageObjects) {
    super();
    this.TimeStamp = TimeStamp;
    this.functionId = functionId;
    this.numberOfObjects = numberOfObjects;
    this.messageObjects = messageObjects;
  }

  public DateAndTime getTimeStamp() {
    return TimeStamp;
  }

  public short getFunctionId() {
    return functionId;
  }

  public short getNumberOfObjects() {
    return numberOfObjects;
  }

  public List<AlarmMessageAckObjectPushType> getMessageObjects() {
    return messageObjects;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("AlarmMessageAckPushType");

    // Simple Field (TimeStamp)
    writeSimpleField("TimeStamp", TimeStamp, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (functionId)
    writeSimpleField("functionId", functionId, writeUnsignedShort(writeBuffer, 8));

    // Simple Field (numberOfObjects)
    writeSimpleField("numberOfObjects", numberOfObjects, writeUnsignedShort(writeBuffer, 8));

    // Array Field (messageObjects)
    writeComplexTypeArrayField("messageObjects", messageObjects, writeBuffer);

    writeBuffer.popContext("AlarmMessageAckPushType");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    AlarmMessageAckPushType _value = this;

    // Simple field (TimeStamp)
    lengthInBits += TimeStamp.getLengthInBits();

    // Simple field (functionId)
    lengthInBits += 8;

    // Simple field (numberOfObjects)
    lengthInBits += 8;

    // Array field
    if (messageObjects != null) {
      int i = 0;
      for (AlarmMessageAckObjectPushType element : messageObjects) {
        boolean last = ++i >= messageObjects.size();
        lengthInBits += element.getLengthInBits();
      }
    }

    return lengthInBits;
  }

  public static AlarmMessageAckPushType staticParse(ReadBuffer readBuffer, Object... args)
      throws ParseException {
    PositionAware positionAware = readBuffer;
    return staticParse(readBuffer);
  }

  public static AlarmMessageAckPushType staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("AlarmMessageAckPushType");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    DateAndTime TimeStamp =
        readSimpleField(
            "TimeStamp",
            new DataReaderComplexDefault<>(() -> DateAndTime.staticParse(readBuffer), readBuffer));

    short functionId = readSimpleField("functionId", readUnsignedShort(readBuffer, 8));

    short numberOfObjects = readSimpleField("numberOfObjects", readUnsignedShort(readBuffer, 8));

    List<AlarmMessageAckObjectPushType> messageObjects =
        readCountArrayField(
            "messageObjects",
            new DataReaderComplexDefault<>(
                () -> AlarmMessageAckObjectPushType.staticParse(readBuffer), readBuffer),
            numberOfObjects);

    readBuffer.closeContext("AlarmMessageAckPushType");
    // Create the instance
    AlarmMessageAckPushType _alarmMessageAckPushType;
    _alarmMessageAckPushType =
        new AlarmMessageAckPushType(TimeStamp, functionId, numberOfObjects, messageObjects);
    return _alarmMessageAckPushType;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof AlarmMessageAckPushType)) {
      return false;
    }
    AlarmMessageAckPushType that = (AlarmMessageAckPushType) o;
    return (getTimeStamp() == that.getTimeStamp())
        && (getFunctionId() == that.getFunctionId())
        && (getNumberOfObjects() == that.getNumberOfObjects())
        && (getMessageObjects() == that.getMessageObjects())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getTimeStamp(), getFunctionId(), getNumberOfObjects(), getMessageObjects());
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
