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

public class NodeIdString extends NodeIdTypeDefinition implements Message {

  // Accessors for discriminator values.
  public NodeIdType getNodeType() {
    return NodeIdType.nodeIdTypeString;
  }

  // Properties.
  protected final int namespaceIndex;
  protected final PascalString id;

  public NodeIdString(int namespaceIndex, PascalString id) {
    super();
    this.namespaceIndex = namespaceIndex;
    this.id = id;
  }

  public int getNamespaceIndex() {
    return namespaceIndex;
  }

  public PascalString getId() {
    return id;
  }

  public String getIdentifier() {
    return String.valueOf(getId().getStringValue());
  }

  public short getNamespace() {
    return (short) (getNamespaceIndex());
  }

  @Override
  protected void serializeNodeIdTypeDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("NodeIdString");

    // Simple Field (namespaceIndex)
    writeSimpleField("namespaceIndex", namespaceIndex, writeUnsignedInt(writeBuffer, 16));

    // Simple Field (id)
    writeSimpleField("id", id, writeComplex(writeBuffer));

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    String identifier = getIdentifier();
    writeBuffer.writeVirtual("identifier", identifier);

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    short namespace = getNamespace();
    writeBuffer.writeVirtual("namespace", namespace);

    writeBuffer.popContext("NodeIdString");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    NodeIdString _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (namespaceIndex)
    lengthInBits += 16;

    // Simple field (id)
    lengthInBits += id.getLengthInBits();

    // A virtual field doesn't have any in- or output.

    // A virtual field doesn't have any in- or output.

    return lengthInBits;
  }

  public static NodeIdTypeDefinitionBuilder staticParseNodeIdTypeDefinitionBuilder(
      ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("NodeIdString");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    int namespaceIndex = readSimpleField("namespaceIndex", readUnsignedInt(readBuffer, 16));

    PascalString id =
        readSimpleField("id", readComplex(() -> PascalString.staticParse(readBuffer), readBuffer));
    String identifier = readVirtualField("identifier", String.class, id.getStringValue());
    short namespace = readVirtualField("namespace", short.class, namespaceIndex);

    readBuffer.closeContext("NodeIdString");
    // Create the instance
    return new NodeIdStringBuilderImpl(namespaceIndex, id);
  }

  public static class NodeIdStringBuilderImpl
      implements NodeIdTypeDefinition.NodeIdTypeDefinitionBuilder {
    private final int namespaceIndex;
    private final PascalString id;

    public NodeIdStringBuilderImpl(int namespaceIndex, PascalString id) {
      this.namespaceIndex = namespaceIndex;
      this.id = id;
    }

    public NodeIdString build() {
      NodeIdString nodeIdString = new NodeIdString(namespaceIndex, id);
      return nodeIdString;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof NodeIdString)) {
      return false;
    }
    NodeIdString that = (NodeIdString) o;
    return (getNamespaceIndex() == that.getNamespaceIndex())
        && (getId() == that.getId())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getNamespaceIndex(), getId());
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
