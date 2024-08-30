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

public class SetMonitoringModeRequest extends ExtensionObjectDefinition implements Message {

  // Accessors for discriminator values.
  public String getIdentifier() {
    return (String) "769";
  }

  // Properties.
  protected final ExtensionObjectDefinition requestHeader;
  protected final long subscriptionId;
  protected final MonitoringMode monitoringMode;
  protected final int noOfMonitoredItemIds;
  protected final List<Long> monitoredItemIds;

  public SetMonitoringModeRequest(
      ExtensionObjectDefinition requestHeader,
      long subscriptionId,
      MonitoringMode monitoringMode,
      int noOfMonitoredItemIds,
      List<Long> monitoredItemIds) {
    super();
    this.requestHeader = requestHeader;
    this.subscriptionId = subscriptionId;
    this.monitoringMode = monitoringMode;
    this.noOfMonitoredItemIds = noOfMonitoredItemIds;
    this.monitoredItemIds = monitoredItemIds;
  }

  public ExtensionObjectDefinition getRequestHeader() {
    return requestHeader;
  }

  public long getSubscriptionId() {
    return subscriptionId;
  }

  public MonitoringMode getMonitoringMode() {
    return monitoringMode;
  }

  public int getNoOfMonitoredItemIds() {
    return noOfMonitoredItemIds;
  }

  public List<Long> getMonitoredItemIds() {
    return monitoredItemIds;
  }

  @Override
  protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("SetMonitoringModeRequest");

    // Simple Field (requestHeader)
    writeSimpleField("requestHeader", requestHeader, writeComplex(writeBuffer));

    // Simple Field (subscriptionId)
    writeSimpleField("subscriptionId", subscriptionId, writeUnsignedLong(writeBuffer, 32));

    // Simple Field (monitoringMode)
    writeSimpleEnumField(
        "monitoringMode",
        "MonitoringMode",
        monitoringMode,
        new DataWriterEnumDefault<>(
            MonitoringMode::getValue, MonitoringMode::name, writeUnsignedLong(writeBuffer, 32)));

    // Simple Field (noOfMonitoredItemIds)
    writeSimpleField("noOfMonitoredItemIds", noOfMonitoredItemIds, writeSignedInt(writeBuffer, 32));

    // Array Field (monitoredItemIds)
    writeSimpleTypeArrayField(
        "monitoredItemIds", monitoredItemIds, writeUnsignedLong(writeBuffer, 32));

    writeBuffer.popContext("SetMonitoringModeRequest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    SetMonitoringModeRequest _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (requestHeader)
    lengthInBits += requestHeader.getLengthInBits();

    // Simple field (subscriptionId)
    lengthInBits += 32;

    // Simple field (monitoringMode)
    lengthInBits += 32;

    // Simple field (noOfMonitoredItemIds)
    lengthInBits += 32;

    // Array field
    if (monitoredItemIds != null) {
      lengthInBits += 32 * monitoredItemIds.size();
    }

    return lengthInBits;
  }

  public static ExtensionObjectDefinitionBuilder staticParseExtensionObjectDefinitionBuilder(
      ReadBuffer readBuffer, String identifier) throws ParseException {
    readBuffer.pullContext("SetMonitoringModeRequest");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    ExtensionObjectDefinition requestHeader =
        readSimpleField(
            "requestHeader",
            readComplex(
                () -> ExtensionObjectDefinition.staticParse(readBuffer, (String) ("391")),
                readBuffer));

    long subscriptionId = readSimpleField("subscriptionId", readUnsignedLong(readBuffer, 32));

    MonitoringMode monitoringMode =
        readEnumField(
            "monitoringMode",
            "MonitoringMode",
            readEnum(MonitoringMode::enumForValue, readUnsignedLong(readBuffer, 32)));

    int noOfMonitoredItemIds =
        readSimpleField("noOfMonitoredItemIds", readSignedInt(readBuffer, 32));

    List<Long> monitoredItemIds =
        readCountArrayField(
            "monitoredItemIds", readUnsignedLong(readBuffer, 32), noOfMonitoredItemIds);

    readBuffer.closeContext("SetMonitoringModeRequest");
    // Create the instance
    return new SetMonitoringModeRequestBuilderImpl(
        requestHeader, subscriptionId, monitoringMode, noOfMonitoredItemIds, monitoredItemIds);
  }

  public static class SetMonitoringModeRequestBuilderImpl
      implements ExtensionObjectDefinition.ExtensionObjectDefinitionBuilder {
    private final ExtensionObjectDefinition requestHeader;
    private final long subscriptionId;
    private final MonitoringMode monitoringMode;
    private final int noOfMonitoredItemIds;
    private final List<Long> monitoredItemIds;

    public SetMonitoringModeRequestBuilderImpl(
        ExtensionObjectDefinition requestHeader,
        long subscriptionId,
        MonitoringMode monitoringMode,
        int noOfMonitoredItemIds,
        List<Long> monitoredItemIds) {
      this.requestHeader = requestHeader;
      this.subscriptionId = subscriptionId;
      this.monitoringMode = monitoringMode;
      this.noOfMonitoredItemIds = noOfMonitoredItemIds;
      this.monitoredItemIds = monitoredItemIds;
    }

    public SetMonitoringModeRequest build() {
      SetMonitoringModeRequest setMonitoringModeRequest =
          new SetMonitoringModeRequest(
              requestHeader,
              subscriptionId,
              monitoringMode,
              noOfMonitoredItemIds,
              monitoredItemIds);
      return setMonitoringModeRequest;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof SetMonitoringModeRequest)) {
      return false;
    }
    SetMonitoringModeRequest that = (SetMonitoringModeRequest) o;
    return (getRequestHeader() == that.getRequestHeader())
        && (getSubscriptionId() == that.getSubscriptionId())
        && (getMonitoringMode() == that.getMonitoringMode())
        && (getNoOfMonitoredItemIds() == that.getNoOfMonitoredItemIds())
        && (getMonitoredItemIds() == that.getMonitoredItemIds())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getRequestHeader(),
        getSubscriptionId(),
        getMonitoringMode(),
        getNoOfMonitoredItemIds(),
        getMonitoredItemIds());
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
