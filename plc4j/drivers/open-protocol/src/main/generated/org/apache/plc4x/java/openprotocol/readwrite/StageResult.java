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
package org.apache.plc4x.java.openprotocol.readwrite;

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

public class StageResult implements Message {

  // Properties.
  protected final long stageTorqueValue;
  protected final long stageTurningAngleValue;

  public StageResult(long stageTorqueValue, long stageTurningAngleValue) {
    super();
    this.stageTorqueValue = stageTorqueValue;
    this.stageTurningAngleValue = stageTurningAngleValue;
  }

  public long getStageTorqueValue() {
    return stageTorqueValue;
  }

  public long getStageTurningAngleValue() {
    return stageTurningAngleValue;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("StageResult");

    // Simple Field (stageTorqueValue)
    writeSimpleField(
        "stageTorqueValue",
        stageTorqueValue,
        writeUnsignedLong(writeBuffer, 48),
        WithOption.WithEncoding("ASCII"));

    // Simple Field (stageTurningAngleValue)
    writeSimpleField(
        "stageTurningAngleValue",
        stageTurningAngleValue,
        writeUnsignedLong(writeBuffer, 40),
        WithOption.WithEncoding("ASCII"));

    writeBuffer.popContext("StageResult");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    StageResult _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (stageTorqueValue)
    lengthInBits += 48;

    // Simple field (stageTurningAngleValue)
    lengthInBits += 40;

    return lengthInBits;
  }

  public static StageResult staticParse(ReadBuffer readBuffer, Object... args)
      throws ParseException {
    PositionAware positionAware = readBuffer;
    return staticParse(readBuffer);
  }

  public static StageResult staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("StageResult");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    long stageTorqueValue =
        readSimpleField(
            "stageTorqueValue", readUnsignedLong(readBuffer, 48), WithOption.WithEncoding("ASCII"));

    long stageTurningAngleValue =
        readSimpleField(
            "stageTurningAngleValue",
            readUnsignedLong(readBuffer, 40),
            WithOption.WithEncoding("ASCII"));

    readBuffer.closeContext("StageResult");
    // Create the instance
    StageResult _stageResult;
    _stageResult = new StageResult(stageTorqueValue, stageTurningAngleValue);
    return _stageResult;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof StageResult)) {
      return false;
    }
    StageResult that = (StageResult) o;
    return (getStageTorqueValue() == that.getStageTorqueValue())
        && (getStageTurningAngleValue() == that.getStageTurningAngleValue())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getStageTorqueValue(), getStageTurningAngleValue());
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
