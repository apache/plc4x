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

public class ConnectionResponseDataBlockTunnelConnection extends ConnectionResponseDataBlock
    implements Message {

  // Accessors for discriminator values.
  public Short getConnectionType() {
    return (short) 0x04;
  }

  // Properties.
  protected final KnxAddress knxAddress;

  public ConnectionResponseDataBlockTunnelConnection(KnxAddress knxAddress) {
    super();
    this.knxAddress = knxAddress;
  }

  public KnxAddress getKnxAddress() {
    return knxAddress;
  }

  @Override
  protected void serializeConnectionResponseDataBlockChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("ConnectionResponseDataBlockTunnelConnection");

    // Simple Field (knxAddress)
    writeSimpleField("knxAddress", knxAddress, new DataWriterComplexDefault<>(writeBuffer));

    writeBuffer.popContext("ConnectionResponseDataBlockTunnelConnection");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    ConnectionResponseDataBlockTunnelConnection _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (knxAddress)
    lengthInBits += knxAddress.getLengthInBits();

    return lengthInBits;
  }

  public static ConnectionResponseDataBlockBuilder staticParseConnectionResponseDataBlockBuilder(
      ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("ConnectionResponseDataBlockTunnelConnection");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    KnxAddress knxAddress =
        readSimpleField(
            "knxAddress",
            new DataReaderComplexDefault<>(() -> KnxAddress.staticParse(readBuffer), readBuffer));

    readBuffer.closeContext("ConnectionResponseDataBlockTunnelConnection");
    // Create the instance
    return new ConnectionResponseDataBlockTunnelConnectionBuilderImpl(knxAddress);
  }

  public static class ConnectionResponseDataBlockTunnelConnectionBuilderImpl
      implements ConnectionResponseDataBlock.ConnectionResponseDataBlockBuilder {
    private final KnxAddress knxAddress;

    public ConnectionResponseDataBlockTunnelConnectionBuilderImpl(KnxAddress knxAddress) {
      this.knxAddress = knxAddress;
    }

    public ConnectionResponseDataBlockTunnelConnection build() {
      ConnectionResponseDataBlockTunnelConnection connectionResponseDataBlockTunnelConnection =
          new ConnectionResponseDataBlockTunnelConnection(knxAddress);
      return connectionResponseDataBlockTunnelConnection;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof ConnectionResponseDataBlockTunnelConnection)) {
      return false;
    }
    ConnectionResponseDataBlockTunnelConnection that =
        (ConnectionResponseDataBlockTunnelConnection) o;
    return (getKnxAddress() == that.getKnxAddress()) && super.equals(that) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getKnxAddress());
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
