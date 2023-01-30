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

public class RepublishRequest extends ExtensionObjectDefinition implements Message {

  // Accessors for discriminator values.
  public String getIdentifier() {
    return (String) "832";
  }

  // Properties.
  protected final ExtensionObjectDefinition requestHeader;
  protected final long subscriptionId;
  protected final long retransmitSequenceNumber;

  public RepublishRequest(
      ExtensionObjectDefinition requestHeader, long subscriptionId, long retransmitSequenceNumber) {
    super();
    this.requestHeader = requestHeader;
    this.subscriptionId = subscriptionId;
    this.retransmitSequenceNumber = retransmitSequenceNumber;
  }

  public ExtensionObjectDefinition getRequestHeader() {
    return requestHeader;
  }

  public long getSubscriptionId() {
    return subscriptionId;
  }

  public long getRetransmitSequenceNumber() {
    return retransmitSequenceNumber;
  }

  @Override
  protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("RepublishRequest");

    // Simple Field (requestHeader)
    writeSimpleField("requestHeader", requestHeader, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (subscriptionId)
    writeSimpleField("subscriptionId", subscriptionId, writeUnsignedLong(writeBuffer, 32));

    // Simple Field (retransmitSequenceNumber)
    writeSimpleField(
        "retransmitSequenceNumber", retransmitSequenceNumber, writeUnsignedLong(writeBuffer, 32));

    writeBuffer.popContext("RepublishRequest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    RepublishRequest _value = this;

    // Simple field (requestHeader)
    lengthInBits += requestHeader.getLengthInBits();

    // Simple field (subscriptionId)
    lengthInBits += 32;

    // Simple field (retransmitSequenceNumber)
    lengthInBits += 32;

    return lengthInBits;
  }

  public static ExtensionObjectDefinitionBuilder staticParseExtensionObjectDefinitionBuilder(
      ReadBuffer readBuffer, String identifier) throws ParseException {
    readBuffer.pullContext("RepublishRequest");
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

    long retransmitSequenceNumber =
        readSimpleField("retransmitSequenceNumber", readUnsignedLong(readBuffer, 32));

    readBuffer.closeContext("RepublishRequest");
    // Create the instance
    return new RepublishRequestBuilderImpl(requestHeader, subscriptionId, retransmitSequenceNumber);
  }

  public static class RepublishRequestBuilderImpl
      implements ExtensionObjectDefinition.ExtensionObjectDefinitionBuilder {
    private final ExtensionObjectDefinition requestHeader;
    private final long subscriptionId;
    private final long retransmitSequenceNumber;

    public RepublishRequestBuilderImpl(
        ExtensionObjectDefinition requestHeader,
        long subscriptionId,
        long retransmitSequenceNumber) {
      this.requestHeader = requestHeader;
      this.subscriptionId = subscriptionId;
      this.retransmitSequenceNumber = retransmitSequenceNumber;
    }

    public RepublishRequest build() {
      RepublishRequest republishRequest =
          new RepublishRequest(requestHeader, subscriptionId, retransmitSequenceNumber);
      return republishRequest;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof RepublishRequest)) {
      return false;
    }
    RepublishRequest that = (RepublishRequest) o;
    return (getRequestHeader() == that.getRequestHeader())
        && (getSubscriptionId() == that.getSubscriptionId())
        && (getRetransmitSequenceNumber() == that.getRetransmitSequenceNumber())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(), getRequestHeader(), getSubscriptionId(), getRetransmitSequenceNumber());
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
