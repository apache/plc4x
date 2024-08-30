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

public class TriggerControlLabelOptions implements Message {

  // Properties.
  protected final TriggerControlLabelFlavour labelFlavour;
  protected final TriggerControlLabelType labelType;

  // Reserved Fields
  private Boolean reservedField0;
  private Boolean reservedField1;
  private Boolean reservedField2;
  private Boolean reservedField3;

  public TriggerControlLabelOptions(
      TriggerControlLabelFlavour labelFlavour, TriggerControlLabelType labelType) {
    super();
    this.labelFlavour = labelFlavour;
    this.labelType = labelType;
  }

  public TriggerControlLabelFlavour getLabelFlavour() {
    return labelFlavour;
  }

  public TriggerControlLabelType getLabelType() {
    return labelType;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("TriggerControlLabelOptions");

    // Reserved Field (reserved)
    writeReservedField(
        "reserved",
        reservedField0 != null ? reservedField0 : (boolean) false,
        writeBoolean(writeBuffer));

    // Simple Field (labelFlavour)
    writeSimpleEnumField(
        "labelFlavour",
        "TriggerControlLabelFlavour",
        labelFlavour,
        new DataWriterEnumDefault<>(
            TriggerControlLabelFlavour::getValue,
            TriggerControlLabelFlavour::name,
            writeUnsignedByte(writeBuffer, 2)));

    // Reserved Field (reserved)
    writeReservedField(
        "reserved",
        reservedField1 != null ? reservedField1 : (boolean) false,
        writeBoolean(writeBuffer));

    // Reserved Field (reserved)
    writeReservedField(
        "reserved",
        reservedField2 != null ? reservedField2 : (boolean) false,
        writeBoolean(writeBuffer));

    // Simple Field (labelType)
    writeSimpleEnumField(
        "labelType",
        "TriggerControlLabelType",
        labelType,
        new DataWriterEnumDefault<>(
            TriggerControlLabelType::getValue,
            TriggerControlLabelType::name,
            writeUnsignedByte(writeBuffer, 2)));

    // Reserved Field (reserved)
    writeReservedField(
        "reserved",
        reservedField3 != null ? reservedField3 : (boolean) false,
        writeBoolean(writeBuffer));

    writeBuffer.popContext("TriggerControlLabelOptions");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    TriggerControlLabelOptions _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Reserved Field (reserved)
    lengthInBits += 1;

    // Simple field (labelFlavour)
    lengthInBits += 2;

    // Reserved Field (reserved)
    lengthInBits += 1;

    // Reserved Field (reserved)
    lengthInBits += 1;

    // Simple field (labelType)
    lengthInBits += 2;

    // Reserved Field (reserved)
    lengthInBits += 1;

    return lengthInBits;
  }

  public static TriggerControlLabelOptions staticParse(ReadBuffer readBuffer)
      throws ParseException {
    readBuffer.pullContext("TriggerControlLabelOptions");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    Boolean reservedField0 =
        readReservedField("reserved", readBoolean(readBuffer), (boolean) false);

    TriggerControlLabelFlavour labelFlavour =
        readEnumField(
            "labelFlavour",
            "TriggerControlLabelFlavour",
            readEnum(TriggerControlLabelFlavour::enumForValue, readUnsignedByte(readBuffer, 2)));

    Boolean reservedField1 =
        readReservedField("reserved", readBoolean(readBuffer), (boolean) false);

    Boolean reservedField2 =
        readReservedField("reserved", readBoolean(readBuffer), (boolean) false);

    TriggerControlLabelType labelType =
        readEnumField(
            "labelType",
            "TriggerControlLabelType",
            readEnum(TriggerControlLabelType::enumForValue, readUnsignedByte(readBuffer, 2)));

    Boolean reservedField3 =
        readReservedField("reserved", readBoolean(readBuffer), (boolean) false);

    readBuffer.closeContext("TriggerControlLabelOptions");
    // Create the instance
    TriggerControlLabelOptions _triggerControlLabelOptions;
    _triggerControlLabelOptions = new TriggerControlLabelOptions(labelFlavour, labelType);
    _triggerControlLabelOptions.reservedField0 = reservedField0;
    _triggerControlLabelOptions.reservedField1 = reservedField1;
    _triggerControlLabelOptions.reservedField2 = reservedField2;
    _triggerControlLabelOptions.reservedField3 = reservedField3;
    return _triggerControlLabelOptions;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof TriggerControlLabelOptions)) {
      return false;
    }
    TriggerControlLabelOptions that = (TriggerControlLabelOptions) o;
    return (getLabelFlavour() == that.getLabelFlavour())
        && (getLabelType() == that.getLabelType())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getLabelFlavour(), getLabelType());
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
