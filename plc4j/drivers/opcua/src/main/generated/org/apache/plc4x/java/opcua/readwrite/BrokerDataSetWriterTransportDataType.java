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

public class BrokerDataSetWriterTransportDataType extends ExtensionObjectDefinition
    implements Message {

  // Accessors for discriminator values.
  public Integer getExtensionId() {
    return (int) 15671;
  }

  // Properties.
  protected final PascalString queueName;
  protected final PascalString resourceUri;
  protected final PascalString authenticationProfileUri;
  protected final BrokerTransportQualityOfService requestedDeliveryGuarantee;
  protected final PascalString metaDataQueueName;
  protected final double metaDataUpdateTime;

  public BrokerDataSetWriterTransportDataType(
      PascalString queueName,
      PascalString resourceUri,
      PascalString authenticationProfileUri,
      BrokerTransportQualityOfService requestedDeliveryGuarantee,
      PascalString metaDataQueueName,
      double metaDataUpdateTime) {
    super();
    this.queueName = queueName;
    this.resourceUri = resourceUri;
    this.authenticationProfileUri = authenticationProfileUri;
    this.requestedDeliveryGuarantee = requestedDeliveryGuarantee;
    this.metaDataQueueName = metaDataQueueName;
    this.metaDataUpdateTime = metaDataUpdateTime;
  }

  public PascalString getQueueName() {
    return queueName;
  }

  public PascalString getResourceUri() {
    return resourceUri;
  }

  public PascalString getAuthenticationProfileUri() {
    return authenticationProfileUri;
  }

  public BrokerTransportQualityOfService getRequestedDeliveryGuarantee() {
    return requestedDeliveryGuarantee;
  }

  public PascalString getMetaDataQueueName() {
    return metaDataQueueName;
  }

  public double getMetaDataUpdateTime() {
    return metaDataUpdateTime;
  }

  @Override
  protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("BrokerDataSetWriterTransportDataType");

    // Simple Field (queueName)
    writeSimpleField("queueName", queueName, writeComplex(writeBuffer));

    // Simple Field (resourceUri)
    writeSimpleField("resourceUri", resourceUri, writeComplex(writeBuffer));

    // Simple Field (authenticationProfileUri)
    writeSimpleField(
        "authenticationProfileUri", authenticationProfileUri, writeComplex(writeBuffer));

    // Simple Field (requestedDeliveryGuarantee)
    writeSimpleEnumField(
        "requestedDeliveryGuarantee",
        "BrokerTransportQualityOfService",
        requestedDeliveryGuarantee,
        writeEnum(
            BrokerTransportQualityOfService::getValue,
            BrokerTransportQualityOfService::name,
            writeUnsignedLong(writeBuffer, 32)));

    // Simple Field (metaDataQueueName)
    writeSimpleField("metaDataQueueName", metaDataQueueName, writeComplex(writeBuffer));

    // Simple Field (metaDataUpdateTime)
    writeSimpleField("metaDataUpdateTime", metaDataUpdateTime, writeDouble(writeBuffer, 64));

    writeBuffer.popContext("BrokerDataSetWriterTransportDataType");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    BrokerDataSetWriterTransportDataType _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (queueName)
    lengthInBits += queueName.getLengthInBits();

    // Simple field (resourceUri)
    lengthInBits += resourceUri.getLengthInBits();

    // Simple field (authenticationProfileUri)
    lengthInBits += authenticationProfileUri.getLengthInBits();

    // Simple field (requestedDeliveryGuarantee)
    lengthInBits += 32;

    // Simple field (metaDataQueueName)
    lengthInBits += metaDataQueueName.getLengthInBits();

    // Simple field (metaDataUpdateTime)
    lengthInBits += 64;

    return lengthInBits;
  }

  public static ExtensionObjectDefinitionBuilder staticParseExtensionObjectDefinitionBuilder(
      ReadBuffer readBuffer, Integer extensionId) throws ParseException {
    readBuffer.pullContext("BrokerDataSetWriterTransportDataType");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    PascalString queueName =
        readSimpleField(
            "queueName", readComplex(() -> PascalString.staticParse(readBuffer), readBuffer));

    PascalString resourceUri =
        readSimpleField(
            "resourceUri", readComplex(() -> PascalString.staticParse(readBuffer), readBuffer));

    PascalString authenticationProfileUri =
        readSimpleField(
            "authenticationProfileUri",
            readComplex(() -> PascalString.staticParse(readBuffer), readBuffer));

    BrokerTransportQualityOfService requestedDeliveryGuarantee =
        readEnumField(
            "requestedDeliveryGuarantee",
            "BrokerTransportQualityOfService",
            readEnum(
                BrokerTransportQualityOfService::enumForValue, readUnsignedLong(readBuffer, 32)));

    PascalString metaDataQueueName =
        readSimpleField(
            "metaDataQueueName",
            readComplex(() -> PascalString.staticParse(readBuffer), readBuffer));

    double metaDataUpdateTime = readSimpleField("metaDataUpdateTime", readDouble(readBuffer, 64));

    readBuffer.closeContext("BrokerDataSetWriterTransportDataType");
    // Create the instance
    return new BrokerDataSetWriterTransportDataTypeBuilderImpl(
        queueName,
        resourceUri,
        authenticationProfileUri,
        requestedDeliveryGuarantee,
        metaDataQueueName,
        metaDataUpdateTime);
  }

  public static class BrokerDataSetWriterTransportDataTypeBuilderImpl
      implements ExtensionObjectDefinition.ExtensionObjectDefinitionBuilder {
    private final PascalString queueName;
    private final PascalString resourceUri;
    private final PascalString authenticationProfileUri;
    private final BrokerTransportQualityOfService requestedDeliveryGuarantee;
    private final PascalString metaDataQueueName;
    private final double metaDataUpdateTime;

    public BrokerDataSetWriterTransportDataTypeBuilderImpl(
        PascalString queueName,
        PascalString resourceUri,
        PascalString authenticationProfileUri,
        BrokerTransportQualityOfService requestedDeliveryGuarantee,
        PascalString metaDataQueueName,
        double metaDataUpdateTime) {
      this.queueName = queueName;
      this.resourceUri = resourceUri;
      this.authenticationProfileUri = authenticationProfileUri;
      this.requestedDeliveryGuarantee = requestedDeliveryGuarantee;
      this.metaDataQueueName = metaDataQueueName;
      this.metaDataUpdateTime = metaDataUpdateTime;
    }

    public BrokerDataSetWriterTransportDataType build() {
      BrokerDataSetWriterTransportDataType brokerDataSetWriterTransportDataType =
          new BrokerDataSetWriterTransportDataType(
              queueName,
              resourceUri,
              authenticationProfileUri,
              requestedDeliveryGuarantee,
              metaDataQueueName,
              metaDataUpdateTime);
      return brokerDataSetWriterTransportDataType;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BrokerDataSetWriterTransportDataType)) {
      return false;
    }
    BrokerDataSetWriterTransportDataType that = (BrokerDataSetWriterTransportDataType) o;
    return (getQueueName() == that.getQueueName())
        && (getResourceUri() == that.getResourceUri())
        && (getAuthenticationProfileUri() == that.getAuthenticationProfileUri())
        && (getRequestedDeliveryGuarantee() == that.getRequestedDeliveryGuarantee())
        && (getMetaDataQueueName() == that.getMetaDataQueueName())
        && (getMetaDataUpdateTime() == that.getMetaDataUpdateTime())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getQueueName(),
        getResourceUri(),
        getAuthenticationProfileUri(),
        getRequestedDeliveryGuarantee(),
        getMetaDataQueueName(),
        getMetaDataUpdateTime());
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
