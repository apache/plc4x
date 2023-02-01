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

public class ModifySubscriptionRequest extends ExtensionObjectDefinition implements Message {

  // Accessors for discriminator values.
  public String getIdentifier() {
    return (String) "793";
  }

  // Properties.
  protected final ExtensionObjectDefinition requestHeader;
  protected final long subscriptionId;
  protected final double requestedPublishingInterval;
  protected final long requestedLifetimeCount;
  protected final long requestedMaxKeepAliveCount;
  protected final long maxNotificationsPerPublish;
  protected final short priority;

  public ModifySubscriptionRequest(
      ExtensionObjectDefinition requestHeader,
      long subscriptionId,
      double requestedPublishingInterval,
      long requestedLifetimeCount,
      long requestedMaxKeepAliveCount,
      long maxNotificationsPerPublish,
      short priority) {
    super();
    this.requestHeader = requestHeader;
    this.subscriptionId = subscriptionId;
    this.requestedPublishingInterval = requestedPublishingInterval;
    this.requestedLifetimeCount = requestedLifetimeCount;
    this.requestedMaxKeepAliveCount = requestedMaxKeepAliveCount;
    this.maxNotificationsPerPublish = maxNotificationsPerPublish;
    this.priority = priority;
  }

  public ExtensionObjectDefinition getRequestHeader() {
    return requestHeader;
  }

  public long getSubscriptionId() {
    return subscriptionId;
  }

  public double getRequestedPublishingInterval() {
    return requestedPublishingInterval;
  }

  public long getRequestedLifetimeCount() {
    return requestedLifetimeCount;
  }

  public long getRequestedMaxKeepAliveCount() {
    return requestedMaxKeepAliveCount;
  }

  public long getMaxNotificationsPerPublish() {
    return maxNotificationsPerPublish;
  }

  public short getPriority() {
    return priority;
  }

  @Override
  protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("ModifySubscriptionRequest");

    // Simple Field (requestHeader)
    writeSimpleField("requestHeader", requestHeader, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (subscriptionId)
    writeSimpleField("subscriptionId", subscriptionId, writeUnsignedLong(writeBuffer, 32));

    // Simple Field (requestedPublishingInterval)
    writeSimpleField(
        "requestedPublishingInterval", requestedPublishingInterval, writeDouble(writeBuffer, 64));

    // Simple Field (requestedLifetimeCount)
    writeSimpleField(
        "requestedLifetimeCount", requestedLifetimeCount, writeUnsignedLong(writeBuffer, 32));

    // Simple Field (requestedMaxKeepAliveCount)
    writeSimpleField(
        "requestedMaxKeepAliveCount",
        requestedMaxKeepAliveCount,
        writeUnsignedLong(writeBuffer, 32));

    // Simple Field (maxNotificationsPerPublish)
    writeSimpleField(
        "maxNotificationsPerPublish",
        maxNotificationsPerPublish,
        writeUnsignedLong(writeBuffer, 32));

    // Simple Field (priority)
    writeSimpleField("priority", priority, writeUnsignedShort(writeBuffer, 8));

    writeBuffer.popContext("ModifySubscriptionRequest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    ModifySubscriptionRequest _value = this;

    // Simple field (requestHeader)
    lengthInBits += requestHeader.getLengthInBits();

    // Simple field (subscriptionId)
    lengthInBits += 32;

    // Simple field (requestedPublishingInterval)
    lengthInBits += 64;

    // Simple field (requestedLifetimeCount)
    lengthInBits += 32;

    // Simple field (requestedMaxKeepAliveCount)
    lengthInBits += 32;

    // Simple field (maxNotificationsPerPublish)
    lengthInBits += 32;

    // Simple field (priority)
    lengthInBits += 8;

    return lengthInBits;
  }

  public static ExtensionObjectDefinitionBuilder staticParseExtensionObjectDefinitionBuilder(
      ReadBuffer readBuffer, String identifier) throws ParseException {
    readBuffer.pullContext("ModifySubscriptionRequest");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    ExtensionObjectDefinition requestHeader =
        readSimpleField(
            "requestHeader",
            new DataReaderComplexDefault<>(
                () -> ExtensionObjectDefinition.staticParse(readBuffer, (String) ("391")),
                readBuffer));

    long subscriptionId = readSimpleField("subscriptionId", readUnsignedLong(readBuffer, 32));

    double requestedPublishingInterval =
        readSimpleField("requestedPublishingInterval", readDouble(readBuffer, 64));

    long requestedLifetimeCount =
        readSimpleField("requestedLifetimeCount", readUnsignedLong(readBuffer, 32));

    long requestedMaxKeepAliveCount =
        readSimpleField("requestedMaxKeepAliveCount", readUnsignedLong(readBuffer, 32));

    long maxNotificationsPerPublish =
        readSimpleField("maxNotificationsPerPublish", readUnsignedLong(readBuffer, 32));

    short priority = readSimpleField("priority", readUnsignedShort(readBuffer, 8));

    readBuffer.closeContext("ModifySubscriptionRequest");
    // Create the instance
    return new ModifySubscriptionRequestBuilderImpl(
        requestHeader,
        subscriptionId,
        requestedPublishingInterval,
        requestedLifetimeCount,
        requestedMaxKeepAliveCount,
        maxNotificationsPerPublish,
        priority);
  }

  public static class ModifySubscriptionRequestBuilderImpl
      implements ExtensionObjectDefinition.ExtensionObjectDefinitionBuilder {
    private final ExtensionObjectDefinition requestHeader;
    private final long subscriptionId;
    private final double requestedPublishingInterval;
    private final long requestedLifetimeCount;
    private final long requestedMaxKeepAliveCount;
    private final long maxNotificationsPerPublish;
    private final short priority;

    public ModifySubscriptionRequestBuilderImpl(
        ExtensionObjectDefinition requestHeader,
        long subscriptionId,
        double requestedPublishingInterval,
        long requestedLifetimeCount,
        long requestedMaxKeepAliveCount,
        long maxNotificationsPerPublish,
        short priority) {
      this.requestHeader = requestHeader;
      this.subscriptionId = subscriptionId;
      this.requestedPublishingInterval = requestedPublishingInterval;
      this.requestedLifetimeCount = requestedLifetimeCount;
      this.requestedMaxKeepAliveCount = requestedMaxKeepAliveCount;
      this.maxNotificationsPerPublish = maxNotificationsPerPublish;
      this.priority = priority;
    }

    public ModifySubscriptionRequest build() {
      ModifySubscriptionRequest modifySubscriptionRequest =
          new ModifySubscriptionRequest(
              requestHeader,
              subscriptionId,
              requestedPublishingInterval,
              requestedLifetimeCount,
              requestedMaxKeepAliveCount,
              maxNotificationsPerPublish,
              priority);
      return modifySubscriptionRequest;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof ModifySubscriptionRequest)) {
      return false;
    }
    ModifySubscriptionRequest that = (ModifySubscriptionRequest) o;
    return (getRequestHeader() == that.getRequestHeader())
        && (getSubscriptionId() == that.getSubscriptionId())
        && (getRequestedPublishingInterval() == that.getRequestedPublishingInterval())
        && (getRequestedLifetimeCount() == that.getRequestedLifetimeCount())
        && (getRequestedMaxKeepAliveCount() == that.getRequestedMaxKeepAliveCount())
        && (getMaxNotificationsPerPublish() == that.getMaxNotificationsPerPublish())
        && (getPriority() == that.getPriority())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getRequestHeader(),
        getSubscriptionId(),
        getRequestedPublishingInterval(),
        getRequestedLifetimeCount(),
        getRequestedMaxKeepAliveCount(),
        getMaxNotificationsPerPublish(),
        getPriority());
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
