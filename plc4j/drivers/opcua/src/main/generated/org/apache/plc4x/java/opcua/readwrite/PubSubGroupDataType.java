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

public class PubSubGroupDataType extends ExtensionObjectDefinition implements Message {

  // Accessors for discriminator values.
  public String getIdentifier() {
    return (String) "15611";
  }

  // Properties.
  protected final PascalString name;
  protected final boolean enabled;
  protected final MessageSecurityMode securityMode;
  protected final PascalString securityGroupId;
  protected final int noOfSecurityKeyServices;
  protected final List<ExtensionObjectDefinition> securityKeyServices;
  protected final long maxNetworkMessageSize;
  protected final int noOfGroupProperties;
  protected final List<ExtensionObjectDefinition> groupProperties;

  public PubSubGroupDataType(
      PascalString name,
      boolean enabled,
      MessageSecurityMode securityMode,
      PascalString securityGroupId,
      int noOfSecurityKeyServices,
      List<ExtensionObjectDefinition> securityKeyServices,
      long maxNetworkMessageSize,
      int noOfGroupProperties,
      List<ExtensionObjectDefinition> groupProperties) {
    super();
    this.name = name;
    this.enabled = enabled;
    this.securityMode = securityMode;
    this.securityGroupId = securityGroupId;
    this.noOfSecurityKeyServices = noOfSecurityKeyServices;
    this.securityKeyServices = securityKeyServices;
    this.maxNetworkMessageSize = maxNetworkMessageSize;
    this.noOfGroupProperties = noOfGroupProperties;
    this.groupProperties = groupProperties;
  }

  public PascalString getName() {
    return name;
  }

  public boolean getEnabled() {
    return enabled;
  }

  public MessageSecurityMode getSecurityMode() {
    return securityMode;
  }

  public PascalString getSecurityGroupId() {
    return securityGroupId;
  }

  public int getNoOfSecurityKeyServices() {
    return noOfSecurityKeyServices;
  }

  public List<ExtensionObjectDefinition> getSecurityKeyServices() {
    return securityKeyServices;
  }

  public long getMaxNetworkMessageSize() {
    return maxNetworkMessageSize;
  }

  public int getNoOfGroupProperties() {
    return noOfGroupProperties;
  }

  public List<ExtensionObjectDefinition> getGroupProperties() {
    return groupProperties;
  }

  @Override
  protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("PubSubGroupDataType");

    // Simple Field (name)
    writeSimpleField("name", name, writeComplex(writeBuffer));

    // Reserved Field (reserved)
    writeReservedField("reserved", (byte) 0x00, writeUnsignedByte(writeBuffer, 7));

    // Simple Field (enabled)
    writeSimpleField("enabled", enabled, writeBoolean(writeBuffer));

    // Simple Field (securityMode)
    writeSimpleEnumField(
        "securityMode",
        "MessageSecurityMode",
        securityMode,
        new DataWriterEnumDefault<>(
            MessageSecurityMode::getValue,
            MessageSecurityMode::name,
            writeUnsignedLong(writeBuffer, 32)));

    // Simple Field (securityGroupId)
    writeSimpleField("securityGroupId", securityGroupId, writeComplex(writeBuffer));

    // Simple Field (noOfSecurityKeyServices)
    writeSimpleField(
        "noOfSecurityKeyServices", noOfSecurityKeyServices, writeSignedInt(writeBuffer, 32));

    // Array Field (securityKeyServices)
    writeComplexTypeArrayField("securityKeyServices", securityKeyServices, writeBuffer);

    // Simple Field (maxNetworkMessageSize)
    writeSimpleField(
        "maxNetworkMessageSize", maxNetworkMessageSize, writeUnsignedLong(writeBuffer, 32));

    // Simple Field (noOfGroupProperties)
    writeSimpleField("noOfGroupProperties", noOfGroupProperties, writeSignedInt(writeBuffer, 32));

    // Array Field (groupProperties)
    writeComplexTypeArrayField("groupProperties", groupProperties, writeBuffer);

    writeBuffer.popContext("PubSubGroupDataType");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    PubSubGroupDataType _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (name)
    lengthInBits += name.getLengthInBits();

    // Reserved Field (reserved)
    lengthInBits += 7;

    // Simple field (enabled)
    lengthInBits += 1;

    // Simple field (securityMode)
    lengthInBits += 32;

    // Simple field (securityGroupId)
    lengthInBits += securityGroupId.getLengthInBits();

    // Simple field (noOfSecurityKeyServices)
    lengthInBits += 32;

    // Array field
    if (securityKeyServices != null) {
      int i = 0;
      for (ExtensionObjectDefinition element : securityKeyServices) {
        ThreadLocalHelper.lastItemThreadLocal.set(++i >= securityKeyServices.size());
        lengthInBits += element.getLengthInBits();
      }
    }

    // Simple field (maxNetworkMessageSize)
    lengthInBits += 32;

    // Simple field (noOfGroupProperties)
    lengthInBits += 32;

    // Array field
    if (groupProperties != null) {
      int i = 0;
      for (ExtensionObjectDefinition element : groupProperties) {
        ThreadLocalHelper.lastItemThreadLocal.set(++i >= groupProperties.size());
        lengthInBits += element.getLengthInBits();
      }
    }

    return lengthInBits;
  }

  public static ExtensionObjectDefinitionBuilder staticParseExtensionObjectDefinitionBuilder(
      ReadBuffer readBuffer, String identifier) throws ParseException {
    readBuffer.pullContext("PubSubGroupDataType");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    PascalString name =
        readSimpleField(
            "name", readComplex(() -> PascalString.staticParse(readBuffer), readBuffer));

    Byte reservedField0 =
        readReservedField("reserved", readUnsignedByte(readBuffer, 7), (byte) 0x00);

    boolean enabled = readSimpleField("enabled", readBoolean(readBuffer));

    MessageSecurityMode securityMode =
        readEnumField(
            "securityMode",
            "MessageSecurityMode",
            readEnum(MessageSecurityMode::enumForValue, readUnsignedLong(readBuffer, 32)));

    PascalString securityGroupId =
        readSimpleField(
            "securityGroupId", readComplex(() -> PascalString.staticParse(readBuffer), readBuffer));

    int noOfSecurityKeyServices =
        readSimpleField("noOfSecurityKeyServices", readSignedInt(readBuffer, 32));

    List<ExtensionObjectDefinition> securityKeyServices =
        readCountArrayField(
            "securityKeyServices",
            readComplex(
                () -> ExtensionObjectDefinition.staticParse(readBuffer, (String) ("314")),
                readBuffer),
            noOfSecurityKeyServices);

    long maxNetworkMessageSize =
        readSimpleField("maxNetworkMessageSize", readUnsignedLong(readBuffer, 32));

    int noOfGroupProperties = readSimpleField("noOfGroupProperties", readSignedInt(readBuffer, 32));

    List<ExtensionObjectDefinition> groupProperties =
        readCountArrayField(
            "groupProperties",
            readComplex(
                () -> ExtensionObjectDefinition.staticParse(readBuffer, (String) ("14535")),
                readBuffer),
            noOfGroupProperties);

    readBuffer.closeContext("PubSubGroupDataType");
    // Create the instance
    return new PubSubGroupDataTypeBuilderImpl(
        name,
        enabled,
        securityMode,
        securityGroupId,
        noOfSecurityKeyServices,
        securityKeyServices,
        maxNetworkMessageSize,
        noOfGroupProperties,
        groupProperties);
  }

  public static class PubSubGroupDataTypeBuilderImpl
      implements ExtensionObjectDefinition.ExtensionObjectDefinitionBuilder {
    private final PascalString name;
    private final boolean enabled;
    private final MessageSecurityMode securityMode;
    private final PascalString securityGroupId;
    private final int noOfSecurityKeyServices;
    private final List<ExtensionObjectDefinition> securityKeyServices;
    private final long maxNetworkMessageSize;
    private final int noOfGroupProperties;
    private final List<ExtensionObjectDefinition> groupProperties;

    public PubSubGroupDataTypeBuilderImpl(
        PascalString name,
        boolean enabled,
        MessageSecurityMode securityMode,
        PascalString securityGroupId,
        int noOfSecurityKeyServices,
        List<ExtensionObjectDefinition> securityKeyServices,
        long maxNetworkMessageSize,
        int noOfGroupProperties,
        List<ExtensionObjectDefinition> groupProperties) {
      this.name = name;
      this.enabled = enabled;
      this.securityMode = securityMode;
      this.securityGroupId = securityGroupId;
      this.noOfSecurityKeyServices = noOfSecurityKeyServices;
      this.securityKeyServices = securityKeyServices;
      this.maxNetworkMessageSize = maxNetworkMessageSize;
      this.noOfGroupProperties = noOfGroupProperties;
      this.groupProperties = groupProperties;
    }

    public PubSubGroupDataType build() {
      PubSubGroupDataType pubSubGroupDataType =
          new PubSubGroupDataType(
              name,
              enabled,
              securityMode,
              securityGroupId,
              noOfSecurityKeyServices,
              securityKeyServices,
              maxNetworkMessageSize,
              noOfGroupProperties,
              groupProperties);
      return pubSubGroupDataType;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof PubSubGroupDataType)) {
      return false;
    }
    PubSubGroupDataType that = (PubSubGroupDataType) o;
    return (getName() == that.getName())
        && (getEnabled() == that.getEnabled())
        && (getSecurityMode() == that.getSecurityMode())
        && (getSecurityGroupId() == that.getSecurityGroupId())
        && (getNoOfSecurityKeyServices() == that.getNoOfSecurityKeyServices())
        && (getSecurityKeyServices() == that.getSecurityKeyServices())
        && (getMaxNetworkMessageSize() == that.getMaxNetworkMessageSize())
        && (getNoOfGroupProperties() == that.getNoOfGroupProperties())
        && (getGroupProperties() == that.getGroupProperties())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getName(),
        getEnabled(),
        getSecurityMode(),
        getSecurityGroupId(),
        getNoOfSecurityKeyServices(),
        getSecurityKeyServices(),
        getMaxNetworkMessageSize(),
        getNoOfGroupProperties(),
        getGroupProperties());
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
