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
package org.apache.plc4x.java.opcua.readwrite;

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

public class HistoryData extends ExtensionObjectDefinition implements Message {

  // Accessors for discriminator values.
  public String getIdentifier() {
    return (String) "658";
  }

  // Properties.
  protected final int noOfDataValues;
  protected final List<DataValue> dataValues;

  public HistoryData(int noOfDataValues, List<DataValue> dataValues) {
    super();
    this.noOfDataValues = noOfDataValues;
    this.dataValues = dataValues;
  }

  public int getNoOfDataValues() {
    return noOfDataValues;
  }

  public List<DataValue> getDataValues() {
    return dataValues;
  }

  @Override
  protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("HistoryData");

    // Simple Field (noOfDataValues)
    writeSimpleField("noOfDataValues", noOfDataValues, writeSignedInt(writeBuffer, 32));

    // Array Field (dataValues)
    writeComplexTypeArrayField("dataValues", dataValues, writeBuffer);

    writeBuffer.popContext("HistoryData");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    HistoryData _value = this;

    // Simple field (noOfDataValues)
    lengthInBits += 32;

    // Array field
    if (dataValues != null) {
      int i = 0;
      for (DataValue element : dataValues) {
        boolean last = ++i >= dataValues.size();
        lengthInBits += element.getLengthInBits();
      }
    }

    return lengthInBits;
  }

  public static HistoryDataBuilder staticParseBuilder(ReadBuffer readBuffer, String identifier)
      throws ParseException {
    readBuffer.pullContext("HistoryData");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    int noOfDataValues = readSimpleField("noOfDataValues", readSignedInt(readBuffer, 32));

    List<DataValue> dataValues =
        readCountArrayField(
            "dataValues",
            new DataReaderComplexDefault<>(() -> DataValue.staticParse(readBuffer), readBuffer),
            noOfDataValues);

    readBuffer.closeContext("HistoryData");
    // Create the instance
    return new HistoryDataBuilder(noOfDataValues, dataValues);
  }

  public static class HistoryDataBuilder
      implements ExtensionObjectDefinition.ExtensionObjectDefinitionBuilder {
    private final int noOfDataValues;
    private final List<DataValue> dataValues;

    public HistoryDataBuilder(int noOfDataValues, List<DataValue> dataValues) {

      this.noOfDataValues = noOfDataValues;
      this.dataValues = dataValues;
    }

    public HistoryData build() {
      HistoryData historyData = new HistoryData(noOfDataValues, dataValues);
      return historyData;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof HistoryData)) {
      return false;
    }
    HistoryData that = (HistoryData) o;
    return (getNoOfDataValues() == that.getNoOfDataValues())
        && (getDataValues() == that.getDataValues())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getNoOfDataValues(), getDataValues());
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
