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

public class PubSubConfigurationDataType extends ExtensionObjectDefinition implements Message {

  // Accessors for discriminator values.
  public Integer getExtensionId() {
    return (int) 15532;
  }

  // Properties.
  protected final List<PublishedDataSetDataType> publishedDataSets;
  protected final List<PubSubConnectionDataType> connections;
  protected final boolean enabled;

  public PubSubConfigurationDataType(
      List<PublishedDataSetDataType> publishedDataSets,
      List<PubSubConnectionDataType> connections,
      boolean enabled) {
    super();
    this.publishedDataSets = publishedDataSets;
    this.connections = connections;
    this.enabled = enabled;
  }

  public List<PublishedDataSetDataType> getPublishedDataSets() {
    return publishedDataSets;
  }

  public List<PubSubConnectionDataType> getConnections() {
    return connections;
  }

  public boolean getEnabled() {
    return enabled;
  }

  @Override
  protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("PubSubConfigurationDataType");

    // Implicit Field (noOfPublishedDataSets) (Used for parsing, but its value is not stored as it's
    // implicitly given by the objects content)
    int noOfPublishedDataSets =
        (int) ((((getPublishedDataSets()) == (null)) ? -(1) : COUNT(getPublishedDataSets())));
    writeImplicitField(
        "noOfPublishedDataSets", noOfPublishedDataSets, writeSignedInt(writeBuffer, 32));

    // Array Field (publishedDataSets)
    writeComplexTypeArrayField("publishedDataSets", publishedDataSets, writeBuffer);

    // Implicit Field (noOfConnections) (Used for parsing, but its value is not stored as it's
    // implicitly given by the objects content)
    int noOfConnections = (int) ((((getConnections()) == (null)) ? -(1) : COUNT(getConnections())));
    writeImplicitField("noOfConnections", noOfConnections, writeSignedInt(writeBuffer, 32));

    // Array Field (connections)
    writeComplexTypeArrayField("connections", connections, writeBuffer);

    // Reserved Field (reserved)
    writeReservedField("reserved", (byte) 0x00, writeUnsignedByte(writeBuffer, 7));

    // Simple Field (enabled)
    writeSimpleField("enabled", enabled, writeBoolean(writeBuffer));

    writeBuffer.popContext("PubSubConfigurationDataType");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    PubSubConfigurationDataType _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Implicit Field (noOfPublishedDataSets)
    lengthInBits += 32;

    // Array field
    if (publishedDataSets != null) {
      int i = 0;
      for (PublishedDataSetDataType element : publishedDataSets) {
        ThreadLocalHelper.lastItemThreadLocal.set(++i >= publishedDataSets.size());
        lengthInBits += element.getLengthInBits();
      }
    }

    // Implicit Field (noOfConnections)
    lengthInBits += 32;

    // Array field
    if (connections != null) {
      int i = 0;
      for (PubSubConnectionDataType element : connections) {
        ThreadLocalHelper.lastItemThreadLocal.set(++i >= connections.size());
        lengthInBits += element.getLengthInBits();
      }
    }

    // Reserved Field (reserved)
    lengthInBits += 7;

    // Simple field (enabled)
    lengthInBits += 1;

    return lengthInBits;
  }

  public static ExtensionObjectDefinitionBuilder staticParseExtensionObjectDefinitionBuilder(
      ReadBuffer readBuffer, Integer extensionId) throws ParseException {
    readBuffer.pullContext("PubSubConfigurationDataType");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    int noOfPublishedDataSets =
        readImplicitField("noOfPublishedDataSets", readSignedInt(readBuffer, 32));

    List<PublishedDataSetDataType> publishedDataSets =
        readCountArrayField(
            "publishedDataSets",
            readComplex(
                () ->
                    (PublishedDataSetDataType)
                        ExtensionObjectDefinition.staticParse(readBuffer, (int) (15580)),
                readBuffer),
            noOfPublishedDataSets);

    int noOfConnections = readImplicitField("noOfConnections", readSignedInt(readBuffer, 32));

    List<PubSubConnectionDataType> connections =
        readCountArrayField(
            "connections",
            readComplex(
                () ->
                    (PubSubConnectionDataType)
                        ExtensionObjectDefinition.staticParse(readBuffer, (int) (15619)),
                readBuffer),
            noOfConnections);

    Byte reservedField0 =
        readReservedField("reserved", readUnsignedByte(readBuffer, 7), (byte) 0x00);

    boolean enabled = readSimpleField("enabled", readBoolean(readBuffer));

    readBuffer.closeContext("PubSubConfigurationDataType");
    // Create the instance
    return new PubSubConfigurationDataTypeBuilderImpl(publishedDataSets, connections, enabled);
  }

  public static class PubSubConfigurationDataTypeBuilderImpl
      implements ExtensionObjectDefinition.ExtensionObjectDefinitionBuilder {
    private final List<PublishedDataSetDataType> publishedDataSets;
    private final List<PubSubConnectionDataType> connections;
    private final boolean enabled;

    public PubSubConfigurationDataTypeBuilderImpl(
        List<PublishedDataSetDataType> publishedDataSets,
        List<PubSubConnectionDataType> connections,
        boolean enabled) {
      this.publishedDataSets = publishedDataSets;
      this.connections = connections;
      this.enabled = enabled;
    }

    public PubSubConfigurationDataType build() {
      PubSubConfigurationDataType pubSubConfigurationDataType =
          new PubSubConfigurationDataType(publishedDataSets, connections, enabled);
      return pubSubConfigurationDataType;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof PubSubConfigurationDataType)) {
      return false;
    }
    PubSubConfigurationDataType that = (PubSubConfigurationDataType) o;
    return (getPublishedDataSets() == that.getPublishedDataSets())
        && (getConnections() == that.getConnections())
        && (getEnabled() == that.getEnabled())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getPublishedDataSets(), getConnections(), getEnabled());
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
