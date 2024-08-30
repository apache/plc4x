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

public class EndpointDescription extends ExtensionObjectDefinition implements Message {

  // Accessors for discriminator values.
  public String getIdentifier() {
    return (String) "314";
  }

  // Properties.
  protected final PascalString endpointUrl;
  protected final ExtensionObjectDefinition server;
  protected final PascalByteString serverCertificate;
  protected final MessageSecurityMode securityMode;
  protected final PascalString securityPolicyUri;
  protected final int noOfUserIdentityTokens;
  protected final List<ExtensionObjectDefinition> userIdentityTokens;
  protected final PascalString transportProfileUri;
  protected final short securityLevel;

  public EndpointDescription(
      PascalString endpointUrl,
      ExtensionObjectDefinition server,
      PascalByteString serverCertificate,
      MessageSecurityMode securityMode,
      PascalString securityPolicyUri,
      int noOfUserIdentityTokens,
      List<ExtensionObjectDefinition> userIdentityTokens,
      PascalString transportProfileUri,
      short securityLevel) {
    super();
    this.endpointUrl = endpointUrl;
    this.server = server;
    this.serverCertificate = serverCertificate;
    this.securityMode = securityMode;
    this.securityPolicyUri = securityPolicyUri;
    this.noOfUserIdentityTokens = noOfUserIdentityTokens;
    this.userIdentityTokens = userIdentityTokens;
    this.transportProfileUri = transportProfileUri;
    this.securityLevel = securityLevel;
  }

  public PascalString getEndpointUrl() {
    return endpointUrl;
  }

  public ExtensionObjectDefinition getServer() {
    return server;
  }

  public PascalByteString getServerCertificate() {
    return serverCertificate;
  }

  public MessageSecurityMode getSecurityMode() {
    return securityMode;
  }

  public PascalString getSecurityPolicyUri() {
    return securityPolicyUri;
  }

  public int getNoOfUserIdentityTokens() {
    return noOfUserIdentityTokens;
  }

  public List<ExtensionObjectDefinition> getUserIdentityTokens() {
    return userIdentityTokens;
  }

  public PascalString getTransportProfileUri() {
    return transportProfileUri;
  }

  public short getSecurityLevel() {
    return securityLevel;
  }

  @Override
  protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("EndpointDescription");

    // Simple Field (endpointUrl)
    writeSimpleField("endpointUrl", endpointUrl, writeComplex(writeBuffer));

    // Simple Field (server)
    writeSimpleField("server", server, writeComplex(writeBuffer));

    // Simple Field (serverCertificate)
    writeSimpleField("serverCertificate", serverCertificate, writeComplex(writeBuffer));

    // Simple Field (securityMode)
    writeSimpleEnumField(
        "securityMode",
        "MessageSecurityMode",
        securityMode,
        new DataWriterEnumDefault<>(
            MessageSecurityMode::getValue,
            MessageSecurityMode::name,
            writeUnsignedLong(writeBuffer, 32)));

    // Simple Field (securityPolicyUri)
    writeSimpleField("securityPolicyUri", securityPolicyUri, writeComplex(writeBuffer));

    // Simple Field (noOfUserIdentityTokens)
    writeSimpleField(
        "noOfUserIdentityTokens", noOfUserIdentityTokens, writeSignedInt(writeBuffer, 32));

    // Array Field (userIdentityTokens)
    writeComplexTypeArrayField("userIdentityTokens", userIdentityTokens, writeBuffer);

    // Simple Field (transportProfileUri)
    writeSimpleField("transportProfileUri", transportProfileUri, writeComplex(writeBuffer));

    // Simple Field (securityLevel)
    writeSimpleField("securityLevel", securityLevel, writeUnsignedShort(writeBuffer, 8));

    writeBuffer.popContext("EndpointDescription");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    EndpointDescription _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (endpointUrl)
    lengthInBits += endpointUrl.getLengthInBits();

    // Simple field (server)
    lengthInBits += server.getLengthInBits();

    // Simple field (serverCertificate)
    lengthInBits += serverCertificate.getLengthInBits();

    // Simple field (securityMode)
    lengthInBits += 32;

    // Simple field (securityPolicyUri)
    lengthInBits += securityPolicyUri.getLengthInBits();

    // Simple field (noOfUserIdentityTokens)
    lengthInBits += 32;

    // Array field
    if (userIdentityTokens != null) {
      int i = 0;
      for (ExtensionObjectDefinition element : userIdentityTokens) {
        ThreadLocalHelper.lastItemThreadLocal.set(++i >= userIdentityTokens.size());
        lengthInBits += element.getLengthInBits();
      }
    }

    // Simple field (transportProfileUri)
    lengthInBits += transportProfileUri.getLengthInBits();

    // Simple field (securityLevel)
    lengthInBits += 8;

    return lengthInBits;
  }

  public static ExtensionObjectDefinitionBuilder staticParseExtensionObjectDefinitionBuilder(
      ReadBuffer readBuffer, String identifier) throws ParseException {
    readBuffer.pullContext("EndpointDescription");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    PascalString endpointUrl =
        readSimpleField(
            "endpointUrl", readComplex(() -> PascalString.staticParse(readBuffer), readBuffer));

    ExtensionObjectDefinition server =
        readSimpleField(
            "server",
            readComplex(
                () -> ExtensionObjectDefinition.staticParse(readBuffer, (String) ("310")),
                readBuffer));

    PascalByteString serverCertificate =
        readSimpleField(
            "serverCertificate",
            readComplex(() -> PascalByteString.staticParse(readBuffer), readBuffer));

    MessageSecurityMode securityMode =
        readEnumField(
            "securityMode",
            "MessageSecurityMode",
            readEnum(MessageSecurityMode::enumForValue, readUnsignedLong(readBuffer, 32)));

    PascalString securityPolicyUri =
        readSimpleField(
            "securityPolicyUri",
            readComplex(() -> PascalString.staticParse(readBuffer), readBuffer));

    int noOfUserIdentityTokens =
        readSimpleField("noOfUserIdentityTokens", readSignedInt(readBuffer, 32));

    List<ExtensionObjectDefinition> userIdentityTokens =
        readCountArrayField(
            "userIdentityTokens",
            readComplex(
                () -> ExtensionObjectDefinition.staticParse(readBuffer, (String) ("306")),
                readBuffer),
            noOfUserIdentityTokens);

    PascalString transportProfileUri =
        readSimpleField(
            "transportProfileUri",
            readComplex(() -> PascalString.staticParse(readBuffer), readBuffer));

    short securityLevel = readSimpleField("securityLevel", readUnsignedShort(readBuffer, 8));

    readBuffer.closeContext("EndpointDescription");
    // Create the instance
    return new EndpointDescriptionBuilderImpl(
        endpointUrl,
        server,
        serverCertificate,
        securityMode,
        securityPolicyUri,
        noOfUserIdentityTokens,
        userIdentityTokens,
        transportProfileUri,
        securityLevel);
  }

  public static class EndpointDescriptionBuilderImpl
      implements ExtensionObjectDefinition.ExtensionObjectDefinitionBuilder {
    private final PascalString endpointUrl;
    private final ExtensionObjectDefinition server;
    private final PascalByteString serverCertificate;
    private final MessageSecurityMode securityMode;
    private final PascalString securityPolicyUri;
    private final int noOfUserIdentityTokens;
    private final List<ExtensionObjectDefinition> userIdentityTokens;
    private final PascalString transportProfileUri;
    private final short securityLevel;

    public EndpointDescriptionBuilderImpl(
        PascalString endpointUrl,
        ExtensionObjectDefinition server,
        PascalByteString serverCertificate,
        MessageSecurityMode securityMode,
        PascalString securityPolicyUri,
        int noOfUserIdentityTokens,
        List<ExtensionObjectDefinition> userIdentityTokens,
        PascalString transportProfileUri,
        short securityLevel) {
      this.endpointUrl = endpointUrl;
      this.server = server;
      this.serverCertificate = serverCertificate;
      this.securityMode = securityMode;
      this.securityPolicyUri = securityPolicyUri;
      this.noOfUserIdentityTokens = noOfUserIdentityTokens;
      this.userIdentityTokens = userIdentityTokens;
      this.transportProfileUri = transportProfileUri;
      this.securityLevel = securityLevel;
    }

    public EndpointDescription build() {
      EndpointDescription endpointDescription =
          new EndpointDescription(
              endpointUrl,
              server,
              serverCertificate,
              securityMode,
              securityPolicyUri,
              noOfUserIdentityTokens,
              userIdentityTokens,
              transportProfileUri,
              securityLevel);
      return endpointDescription;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof EndpointDescription)) {
      return false;
    }
    EndpointDescription that = (EndpointDescription) o;
    return (getEndpointUrl() == that.getEndpointUrl())
        && (getServer() == that.getServer())
        && (getServerCertificate() == that.getServerCertificate())
        && (getSecurityMode() == that.getSecurityMode())
        && (getSecurityPolicyUri() == that.getSecurityPolicyUri())
        && (getNoOfUserIdentityTokens() == that.getNoOfUserIdentityTokens())
        && (getUserIdentityTokens() == that.getUserIdentityTokens())
        && (getTransportProfileUri() == that.getTransportProfileUri())
        && (getSecurityLevel() == that.getSecurityLevel())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getEndpointUrl(),
        getServer(),
        getServerCertificate(),
        getSecurityMode(),
        getSecurityPolicyUri(),
        getNoOfUserIdentityTokens(),
        getUserIdentityTokens(),
        getTransportProfileUri(),
        getSecurityLevel());
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
