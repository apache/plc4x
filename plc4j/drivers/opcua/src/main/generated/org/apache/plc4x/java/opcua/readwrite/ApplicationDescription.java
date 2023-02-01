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

public class ApplicationDescription extends ExtensionObjectDefinition implements Message {

  // Accessors for discriminator values.
  public String getIdentifier() {
    return (String) "310";
  }

  // Properties.
  protected final PascalString applicationUri;
  protected final PascalString productUri;
  protected final LocalizedText applicationName;
  protected final ApplicationType applicationType;
  protected final PascalString gatewayServerUri;
  protected final PascalString discoveryProfileUri;
  protected final int noOfDiscoveryUrls;
  protected final List<PascalString> discoveryUrls;

  public ApplicationDescription(
      PascalString applicationUri,
      PascalString productUri,
      LocalizedText applicationName,
      ApplicationType applicationType,
      PascalString gatewayServerUri,
      PascalString discoveryProfileUri,
      int noOfDiscoveryUrls,
      List<PascalString> discoveryUrls) {
    super();
    this.applicationUri = applicationUri;
    this.productUri = productUri;
    this.applicationName = applicationName;
    this.applicationType = applicationType;
    this.gatewayServerUri = gatewayServerUri;
    this.discoveryProfileUri = discoveryProfileUri;
    this.noOfDiscoveryUrls = noOfDiscoveryUrls;
    this.discoveryUrls = discoveryUrls;
  }

  public PascalString getApplicationUri() {
    return applicationUri;
  }

  public PascalString getProductUri() {
    return productUri;
  }

  public LocalizedText getApplicationName() {
    return applicationName;
  }

  public ApplicationType getApplicationType() {
    return applicationType;
  }

  public PascalString getGatewayServerUri() {
    return gatewayServerUri;
  }

  public PascalString getDiscoveryProfileUri() {
    return discoveryProfileUri;
  }

  public int getNoOfDiscoveryUrls() {
    return noOfDiscoveryUrls;
  }

  public List<PascalString> getDiscoveryUrls() {
    return discoveryUrls;
  }

  @Override
  protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("ApplicationDescription");

    // Simple Field (applicationUri)
    writeSimpleField("applicationUri", applicationUri, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (productUri)
    writeSimpleField("productUri", productUri, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (applicationName)
    writeSimpleField(
        "applicationName", applicationName, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (applicationType)
    writeSimpleEnumField(
        "applicationType",
        "ApplicationType",
        applicationType,
        new DataWriterEnumDefault<>(
            ApplicationType::getValue, ApplicationType::name, writeUnsignedLong(writeBuffer, 32)));

    // Simple Field (gatewayServerUri)
    writeSimpleField(
        "gatewayServerUri", gatewayServerUri, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (discoveryProfileUri)
    writeSimpleField(
        "discoveryProfileUri", discoveryProfileUri, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (noOfDiscoveryUrls)
    writeSimpleField("noOfDiscoveryUrls", noOfDiscoveryUrls, writeSignedInt(writeBuffer, 32));

    // Array Field (discoveryUrls)
    writeComplexTypeArrayField("discoveryUrls", discoveryUrls, writeBuffer);

    writeBuffer.popContext("ApplicationDescription");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    ApplicationDescription _value = this;

    // Simple field (applicationUri)
    lengthInBits += applicationUri.getLengthInBits();

    // Simple field (productUri)
    lengthInBits += productUri.getLengthInBits();

    // Simple field (applicationName)
    lengthInBits += applicationName.getLengthInBits();

    // Simple field (applicationType)
    lengthInBits += 32;

    // Simple field (gatewayServerUri)
    lengthInBits += gatewayServerUri.getLengthInBits();

    // Simple field (discoveryProfileUri)
    lengthInBits += discoveryProfileUri.getLengthInBits();

    // Simple field (noOfDiscoveryUrls)
    lengthInBits += 32;

    // Array field
    if (discoveryUrls != null) {
      int i = 0;
      for (PascalString element : discoveryUrls) {
        boolean last = ++i >= discoveryUrls.size();
        lengthInBits += element.getLengthInBits();
      }
    }

    return lengthInBits;
  }

  public static ExtensionObjectDefinitionBuilder staticParseExtensionObjectDefinitionBuilder(
      ReadBuffer readBuffer, String identifier) throws ParseException {
    readBuffer.pullContext("ApplicationDescription");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    PascalString applicationUri =
        readSimpleField(
            "applicationUri",
            new DataReaderComplexDefault<>(() -> PascalString.staticParse(readBuffer), readBuffer));

    PascalString productUri =
        readSimpleField(
            "productUri",
            new DataReaderComplexDefault<>(() -> PascalString.staticParse(readBuffer), readBuffer));

    LocalizedText applicationName =
        readSimpleField(
            "applicationName",
            new DataReaderComplexDefault<>(
                () -> LocalizedText.staticParse(readBuffer), readBuffer));

    ApplicationType applicationType =
        readEnumField(
            "applicationType",
            "ApplicationType",
            new DataReaderEnumDefault<>(
                ApplicationType::enumForValue, readUnsignedLong(readBuffer, 32)));

    PascalString gatewayServerUri =
        readSimpleField(
            "gatewayServerUri",
            new DataReaderComplexDefault<>(() -> PascalString.staticParse(readBuffer), readBuffer));

    PascalString discoveryProfileUri =
        readSimpleField(
            "discoveryProfileUri",
            new DataReaderComplexDefault<>(() -> PascalString.staticParse(readBuffer), readBuffer));

    int noOfDiscoveryUrls = readSimpleField("noOfDiscoveryUrls", readSignedInt(readBuffer, 32));

    List<PascalString> discoveryUrls =
        readCountArrayField(
            "discoveryUrls",
            new DataReaderComplexDefault<>(() -> PascalString.staticParse(readBuffer), readBuffer),
            noOfDiscoveryUrls);

    readBuffer.closeContext("ApplicationDescription");
    // Create the instance
    return new ApplicationDescriptionBuilderImpl(
        applicationUri,
        productUri,
        applicationName,
        applicationType,
        gatewayServerUri,
        discoveryProfileUri,
        noOfDiscoveryUrls,
        discoveryUrls);
  }

  public static class ApplicationDescriptionBuilderImpl
      implements ExtensionObjectDefinition.ExtensionObjectDefinitionBuilder {
    private final PascalString applicationUri;
    private final PascalString productUri;
    private final LocalizedText applicationName;
    private final ApplicationType applicationType;
    private final PascalString gatewayServerUri;
    private final PascalString discoveryProfileUri;
    private final int noOfDiscoveryUrls;
    private final List<PascalString> discoveryUrls;

    public ApplicationDescriptionBuilderImpl(
        PascalString applicationUri,
        PascalString productUri,
        LocalizedText applicationName,
        ApplicationType applicationType,
        PascalString gatewayServerUri,
        PascalString discoveryProfileUri,
        int noOfDiscoveryUrls,
        List<PascalString> discoveryUrls) {
      this.applicationUri = applicationUri;
      this.productUri = productUri;
      this.applicationName = applicationName;
      this.applicationType = applicationType;
      this.gatewayServerUri = gatewayServerUri;
      this.discoveryProfileUri = discoveryProfileUri;
      this.noOfDiscoveryUrls = noOfDiscoveryUrls;
      this.discoveryUrls = discoveryUrls;
    }

    public ApplicationDescription build() {
      ApplicationDescription applicationDescription =
          new ApplicationDescription(
              applicationUri,
              productUri,
              applicationName,
              applicationType,
              gatewayServerUri,
              discoveryProfileUri,
              noOfDiscoveryUrls,
              discoveryUrls);
      return applicationDescription;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof ApplicationDescription)) {
      return false;
    }
    ApplicationDescription that = (ApplicationDescription) o;
    return (getApplicationUri() == that.getApplicationUri())
        && (getProductUri() == that.getProductUri())
        && (getApplicationName() == that.getApplicationName())
        && (getApplicationType() == that.getApplicationType())
        && (getGatewayServerUri() == that.getGatewayServerUri())
        && (getDiscoveryProfileUri() == that.getDiscoveryProfileUri())
        && (getNoOfDiscoveryUrls() == that.getNoOfDiscoveryUrls())
        && (getDiscoveryUrls() == that.getDiscoveryUrls())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getApplicationUri(),
        getProductUri(),
        getApplicationName(),
        getApplicationType(),
        getGatewayServerUri(),
        getDiscoveryProfileUri(),
        getNoOfDiscoveryUrls(),
        getDiscoveryUrls());
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
