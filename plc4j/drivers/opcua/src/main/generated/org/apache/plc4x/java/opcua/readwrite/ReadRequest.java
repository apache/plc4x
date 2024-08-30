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

public class ReadRequest extends ExtensionObjectDefinition implements Message {

  // Accessors for discriminator values.
  public String getIdentifier() {
    return (String) "631";
  }

  // Properties.
  protected final ExtensionObjectDefinition requestHeader;
  protected final double maxAge;
  protected final TimestampsToReturn timestampsToReturn;
  protected final int noOfNodesToRead;
  protected final List<ExtensionObjectDefinition> nodesToRead;

  public ReadRequest(
      ExtensionObjectDefinition requestHeader,
      double maxAge,
      TimestampsToReturn timestampsToReturn,
      int noOfNodesToRead,
      List<ExtensionObjectDefinition> nodesToRead) {
    super();
    this.requestHeader = requestHeader;
    this.maxAge = maxAge;
    this.timestampsToReturn = timestampsToReturn;
    this.noOfNodesToRead = noOfNodesToRead;
    this.nodesToRead = nodesToRead;
  }

  public ExtensionObjectDefinition getRequestHeader() {
    return requestHeader;
  }

  public double getMaxAge() {
    return maxAge;
  }

  public TimestampsToReturn getTimestampsToReturn() {
    return timestampsToReturn;
  }

  public int getNoOfNodesToRead() {
    return noOfNodesToRead;
  }

  public List<ExtensionObjectDefinition> getNodesToRead() {
    return nodesToRead;
  }

  @Override
  protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("ReadRequest");

    // Simple Field (requestHeader)
    writeSimpleField("requestHeader", requestHeader, writeComplex(writeBuffer));

    // Simple Field (maxAge)
    writeSimpleField("maxAge", maxAge, writeDouble(writeBuffer, 64));

    // Simple Field (timestampsToReturn)
    writeSimpleEnumField(
        "timestampsToReturn",
        "TimestampsToReturn",
        timestampsToReturn,
        new DataWriterEnumDefault<>(
            TimestampsToReturn::getValue,
            TimestampsToReturn::name,
            writeUnsignedLong(writeBuffer, 32)));

    // Simple Field (noOfNodesToRead)
    writeSimpleField("noOfNodesToRead", noOfNodesToRead, writeSignedInt(writeBuffer, 32));

    // Array Field (nodesToRead)
    writeComplexTypeArrayField("nodesToRead", nodesToRead, writeBuffer);

    writeBuffer.popContext("ReadRequest");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    ReadRequest _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (requestHeader)
    lengthInBits += requestHeader.getLengthInBits();

    // Simple field (maxAge)
    lengthInBits += 64;

    // Simple field (timestampsToReturn)
    lengthInBits += 32;

    // Simple field (noOfNodesToRead)
    lengthInBits += 32;

    // Array field
    if (nodesToRead != null) {
      int i = 0;
      for (ExtensionObjectDefinition element : nodesToRead) {
        ThreadLocalHelper.lastItemThreadLocal.set(++i >= nodesToRead.size());
        lengthInBits += element.getLengthInBits();
      }
    }

    return lengthInBits;
  }

  public static ExtensionObjectDefinitionBuilder staticParseExtensionObjectDefinitionBuilder(
      ReadBuffer readBuffer, String identifier) throws ParseException {
    readBuffer.pullContext("ReadRequest");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    ExtensionObjectDefinition requestHeader =
        readSimpleField(
            "requestHeader",
            readComplex(
                () -> ExtensionObjectDefinition.staticParse(readBuffer, (String) ("391")),
                readBuffer));

    double maxAge = readSimpleField("maxAge", readDouble(readBuffer, 64));

    TimestampsToReturn timestampsToReturn =
        readEnumField(
            "timestampsToReturn",
            "TimestampsToReturn",
            readEnum(TimestampsToReturn::enumForValue, readUnsignedLong(readBuffer, 32)));

    int noOfNodesToRead = readSimpleField("noOfNodesToRead", readSignedInt(readBuffer, 32));

    List<ExtensionObjectDefinition> nodesToRead =
        readCountArrayField(
            "nodesToRead",
            readComplex(
                () -> ExtensionObjectDefinition.staticParse(readBuffer, (String) ("628")),
                readBuffer),
            noOfNodesToRead);

    readBuffer.closeContext("ReadRequest");
    // Create the instance
    return new ReadRequestBuilderImpl(
        requestHeader, maxAge, timestampsToReturn, noOfNodesToRead, nodesToRead);
  }

  public static class ReadRequestBuilderImpl
      implements ExtensionObjectDefinition.ExtensionObjectDefinitionBuilder {
    private final ExtensionObjectDefinition requestHeader;
    private final double maxAge;
    private final TimestampsToReturn timestampsToReturn;
    private final int noOfNodesToRead;
    private final List<ExtensionObjectDefinition> nodesToRead;

    public ReadRequestBuilderImpl(
        ExtensionObjectDefinition requestHeader,
        double maxAge,
        TimestampsToReturn timestampsToReturn,
        int noOfNodesToRead,
        List<ExtensionObjectDefinition> nodesToRead) {
      this.requestHeader = requestHeader;
      this.maxAge = maxAge;
      this.timestampsToReturn = timestampsToReturn;
      this.noOfNodesToRead = noOfNodesToRead;
      this.nodesToRead = nodesToRead;
    }

    public ReadRequest build() {
      ReadRequest readRequest =
          new ReadRequest(requestHeader, maxAge, timestampsToReturn, noOfNodesToRead, nodesToRead);
      return readRequest;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof ReadRequest)) {
      return false;
    }
    ReadRequest that = (ReadRequest) o;
    return (getRequestHeader() == that.getRequestHeader())
        && (getMaxAge() == that.getMaxAge())
        && (getTimestampsToReturn() == that.getTimestampsToReturn())
        && (getNoOfNodesToRead() == that.getNoOfNodesToRead())
        && (getNodesToRead() == that.getNodesToRead())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getRequestHeader(),
        getMaxAge(),
        getTimestampsToReturn(),
        getNoOfNodesToRead(),
        getNodesToRead());
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
