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
package org.apache.plc4x.java.knxnetip.readwrite;

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

public class GroupObjectDescriptorRealisationTypeB implements Message {

  // Properties.
  protected final boolean updateEnable;
  protected final boolean transmitEnable;
  protected final boolean segmentSelectorEnable;
  protected final boolean writeEnable;
  protected final boolean readEnable;
  protected final boolean communicationEnable;
  protected final CEMIPriority priority;
  protected final ComObjectValueType valueType;

  public GroupObjectDescriptorRealisationTypeB(
      boolean updateEnable,
      boolean transmitEnable,
      boolean segmentSelectorEnable,
      boolean writeEnable,
      boolean readEnable,
      boolean communicationEnable,
      CEMIPriority priority,
      ComObjectValueType valueType) {
    super();
    this.updateEnable = updateEnable;
    this.transmitEnable = transmitEnable;
    this.segmentSelectorEnable = segmentSelectorEnable;
    this.writeEnable = writeEnable;
    this.readEnable = readEnable;
    this.communicationEnable = communicationEnable;
    this.priority = priority;
    this.valueType = valueType;
  }

  public boolean getUpdateEnable() {
    return updateEnable;
  }

  public boolean getTransmitEnable() {
    return transmitEnable;
  }

  public boolean getSegmentSelectorEnable() {
    return segmentSelectorEnable;
  }

  public boolean getWriteEnable() {
    return writeEnable;
  }

  public boolean getReadEnable() {
    return readEnable;
  }

  public boolean getCommunicationEnable() {
    return communicationEnable;
  }

  public CEMIPriority getPriority() {
    return priority;
  }

  public ComObjectValueType getValueType() {
    return valueType;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("GroupObjectDescriptorRealisationTypeB");

    // Simple Field (updateEnable)
    writeSimpleField("updateEnable", updateEnable, writeBoolean(writeBuffer));

    // Simple Field (transmitEnable)
    writeSimpleField("transmitEnable", transmitEnable, writeBoolean(writeBuffer));

    // Simple Field (segmentSelectorEnable)
    writeSimpleField("segmentSelectorEnable", segmentSelectorEnable, writeBoolean(writeBuffer));

    // Simple Field (writeEnable)
    writeSimpleField("writeEnable", writeEnable, writeBoolean(writeBuffer));

    // Simple Field (readEnable)
    writeSimpleField("readEnable", readEnable, writeBoolean(writeBuffer));

    // Simple Field (communicationEnable)
    writeSimpleField("communicationEnable", communicationEnable, writeBoolean(writeBuffer));

    // Simple Field (priority)
    writeSimpleEnumField(
        "priority",
        "CEMIPriority",
        priority,
        writeEnum(CEMIPriority::getValue, CEMIPriority::name, writeUnsignedByte(writeBuffer, 2)));

    // Simple Field (valueType)
    writeSimpleEnumField(
        "valueType",
        "ComObjectValueType",
        valueType,
        writeEnum(
            ComObjectValueType::getValue,
            ComObjectValueType::name,
            writeUnsignedShort(writeBuffer, 8)));

    writeBuffer.popContext("GroupObjectDescriptorRealisationTypeB");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    GroupObjectDescriptorRealisationTypeB _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (updateEnable)
    lengthInBits += 1;

    // Simple field (transmitEnable)
    lengthInBits += 1;

    // Simple field (segmentSelectorEnable)
    lengthInBits += 1;

    // Simple field (writeEnable)
    lengthInBits += 1;

    // Simple field (readEnable)
    lengthInBits += 1;

    // Simple field (communicationEnable)
    lengthInBits += 1;

    // Simple field (priority)
    lengthInBits += 2;

    // Simple field (valueType)
    lengthInBits += 8;

    return lengthInBits;
  }

  public static GroupObjectDescriptorRealisationTypeB staticParse(ReadBuffer readBuffer)
      throws ParseException {
    readBuffer.pullContext("GroupObjectDescriptorRealisationTypeB");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    boolean updateEnable = readSimpleField("updateEnable", readBoolean(readBuffer));

    boolean transmitEnable = readSimpleField("transmitEnable", readBoolean(readBuffer));

    boolean segmentSelectorEnable =
        readSimpleField("segmentSelectorEnable", readBoolean(readBuffer));

    boolean writeEnable = readSimpleField("writeEnable", readBoolean(readBuffer));

    boolean readEnable = readSimpleField("readEnable", readBoolean(readBuffer));

    boolean communicationEnable = readSimpleField("communicationEnable", readBoolean(readBuffer));

    CEMIPriority priority =
        readEnumField(
            "priority",
            "CEMIPriority",
            readEnum(CEMIPriority::enumForValue, readUnsignedByte(readBuffer, 2)));

    ComObjectValueType valueType =
        readEnumField(
            "valueType",
            "ComObjectValueType",
            readEnum(ComObjectValueType::enumForValue, readUnsignedShort(readBuffer, 8)));

    readBuffer.closeContext("GroupObjectDescriptorRealisationTypeB");
    // Create the instance
    GroupObjectDescriptorRealisationTypeB _groupObjectDescriptorRealisationTypeB;
    _groupObjectDescriptorRealisationTypeB =
        new GroupObjectDescriptorRealisationTypeB(
            updateEnable,
            transmitEnable,
            segmentSelectorEnable,
            writeEnable,
            readEnable,
            communicationEnable,
            priority,
            valueType);
    return _groupObjectDescriptorRealisationTypeB;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof GroupObjectDescriptorRealisationTypeB)) {
      return false;
    }
    GroupObjectDescriptorRealisationTypeB that = (GroupObjectDescriptorRealisationTypeB) o;
    return (getUpdateEnable() == that.getUpdateEnable())
        && (getTransmitEnable() == that.getTransmitEnable())
        && (getSegmentSelectorEnable() == that.getSegmentSelectorEnable())
        && (getWriteEnable() == that.getWriteEnable())
        && (getReadEnable() == that.getReadEnable())
        && (getCommunicationEnable() == that.getCommunicationEnable())
        && (getPriority() == that.getPriority())
        && (getValueType() == that.getValueType())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        getUpdateEnable(),
        getTransmitEnable(),
        getSegmentSelectorEnable(),
        getWriteEnable(),
        getReadEnable(),
        getCommunicationEnable(),
        getPriority(),
        getValueType());
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
