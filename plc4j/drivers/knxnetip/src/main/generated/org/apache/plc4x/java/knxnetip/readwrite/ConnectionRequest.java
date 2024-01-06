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
package org.apache.plc4x.java.knxnetip.readwrite;

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

public class ConnectionRequest extends KnxNetIpMessage implements Message {

  // Accessors for discriminator values.
  public Integer getMsgType() {
    return (int) 0x0205;
  }

  // Properties.
  protected final HPAIDiscoveryEndpoint hpaiDiscoveryEndpoint;
  protected final HPAIDataEndpoint hpaiDataEndpoint;
  protected final ConnectionRequestInformation connectionRequestInformation;

  public ConnectionRequest(
      HPAIDiscoveryEndpoint hpaiDiscoveryEndpoint,
      HPAIDataEndpoint hpaiDataEndpoint,
      ConnectionRequestInformation connectionRequestInformation) {
    super();
    this.hpaiDiscoveryEndpoint = hpaiDiscoveryEndpoint;
    this.hpaiDataEndpoint = hpaiDataEndpoint;
    this.connectionRequestInformation = connectionRequestInformation;
  }

  public HPAIDiscoveryEndpoint getHpaiDiscoveryEndpoint() {
    return hpaiDiscoveryEndpoint;
  }

  public HPAIDataEndpoint getHpaiDataEndpoint() {
    return hpaiDataEndpoint;
  }

  public ConnectionRequestInformation getConnectionRequestInformation() {
    return connectionRequestInformation;
  }

  @Override
  protected void serializeKnxNetIpMessageChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("ConnectionRequest");

    // Simple Field (hpaiDiscoveryEndpoint)
    writeSimpleField(
        "hpaiDiscoveryEndpoint",
        hpaiDiscoveryEndpoint,
        new DataWriterComplexDefault<>(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (hpaiDataEndpoint)
    writeSimpleField(
        "hpaiDataEndpoint",
        hpaiDataEndpoint,
        new DataWriterComplexDefault<>(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    // Simple Field (connectionRequestInformation)
    writeSimpleField(
        "connectionRequestInformation",
        connectionRequestInformation,
        new DataWriterComplexDefault<>(writeBuffer),
        WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    writeBuffer.popContext("ConnectionRequest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    ConnectionRequest _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (hpaiDiscoveryEndpoint)
    lengthInBits += hpaiDiscoveryEndpoint.getLengthInBits();

    // Simple field (hpaiDataEndpoint)
    lengthInBits += hpaiDataEndpoint.getLengthInBits();

    // Simple field (connectionRequestInformation)
    lengthInBits += connectionRequestInformation.getLengthInBits();

    return lengthInBits;
  }

  public static KnxNetIpMessageBuilder staticParseKnxNetIpMessageBuilder(ReadBuffer readBuffer)
      throws ParseException {
    readBuffer.pullContext("ConnectionRequest");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    HPAIDiscoveryEndpoint hpaiDiscoveryEndpoint =
        readSimpleField(
            "hpaiDiscoveryEndpoint",
            new DataReaderComplexDefault<>(
                () -> HPAIDiscoveryEndpoint.staticParse(readBuffer), readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    HPAIDataEndpoint hpaiDataEndpoint =
        readSimpleField(
            "hpaiDataEndpoint",
            new DataReaderComplexDefault<>(
                () -> HPAIDataEndpoint.staticParse(readBuffer), readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    ConnectionRequestInformation connectionRequestInformation =
        readSimpleField(
            "connectionRequestInformation",
            new DataReaderComplexDefault<>(
                () -> ConnectionRequestInformation.staticParse(readBuffer), readBuffer),
            WithOption.WithByteOrder(ByteOrder.BIG_ENDIAN));

    readBuffer.closeContext("ConnectionRequest");
    // Create the instance
    return new ConnectionRequestBuilderImpl(
        hpaiDiscoveryEndpoint, hpaiDataEndpoint, connectionRequestInformation);
  }

  public static class ConnectionRequestBuilderImpl
      implements KnxNetIpMessage.KnxNetIpMessageBuilder {
    private final HPAIDiscoveryEndpoint hpaiDiscoveryEndpoint;
    private final HPAIDataEndpoint hpaiDataEndpoint;
    private final ConnectionRequestInformation connectionRequestInformation;

    public ConnectionRequestBuilderImpl(
        HPAIDiscoveryEndpoint hpaiDiscoveryEndpoint,
        HPAIDataEndpoint hpaiDataEndpoint,
        ConnectionRequestInformation connectionRequestInformation) {
      this.hpaiDiscoveryEndpoint = hpaiDiscoveryEndpoint;
      this.hpaiDataEndpoint = hpaiDataEndpoint;
      this.connectionRequestInformation = connectionRequestInformation;
    }

    public ConnectionRequest build() {
      ConnectionRequest connectionRequest =
          new ConnectionRequest(
              hpaiDiscoveryEndpoint, hpaiDataEndpoint, connectionRequestInformation);
      return connectionRequest;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof ConnectionRequest)) {
      return false;
    }
    ConnectionRequest that = (ConnectionRequest) o;
    return (getHpaiDiscoveryEndpoint() == that.getHpaiDiscoveryEndpoint())
        && (getHpaiDataEndpoint() == that.getHpaiDataEndpoint())
        && (getConnectionRequestInformation() == that.getConnectionRequestInformation())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getHpaiDiscoveryEndpoint(),
        getHpaiDataEndpoint(),
        getConnectionRequestInformation());
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
