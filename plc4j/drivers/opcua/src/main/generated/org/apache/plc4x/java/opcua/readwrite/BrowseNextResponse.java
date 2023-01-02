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

public class BrowseNextResponse extends ExtensionObjectDefinition implements Message {

  // Accessors for discriminator values.
  public String getIdentifier() {
    return (String) "536";
  }

  // Properties.
  protected final ExtensionObjectDefinition responseHeader;
  protected final int noOfResults;
  protected final List<ExtensionObjectDefinition> results;
  protected final int noOfDiagnosticInfos;
  protected final List<DiagnosticInfo> diagnosticInfos;

  public BrowseNextResponse(
      ExtensionObjectDefinition responseHeader,
      int noOfResults,
      List<ExtensionObjectDefinition> results,
      int noOfDiagnosticInfos,
      List<DiagnosticInfo> diagnosticInfos) {
    super();
    this.responseHeader = responseHeader;
    this.noOfResults = noOfResults;
    this.results = results;
    this.noOfDiagnosticInfos = noOfDiagnosticInfos;
    this.diagnosticInfos = diagnosticInfos;
  }

  public ExtensionObjectDefinition getResponseHeader() {
    return responseHeader;
  }

  public int getNoOfResults() {
    return noOfResults;
  }

  public List<ExtensionObjectDefinition> getResults() {
    return results;
  }

  public int getNoOfDiagnosticInfos() {
    return noOfDiagnosticInfos;
  }

  public List<DiagnosticInfo> getDiagnosticInfos() {
    return diagnosticInfos;
  }

  @Override
  protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("BrowseNextResponse");

    // Simple Field (responseHeader)
    writeSimpleField("responseHeader", responseHeader, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (noOfResults)
    writeSimpleField("noOfResults", noOfResults, writeSignedInt(writeBuffer, 32));

    // Array Field (results)
    writeComplexTypeArrayField("results", results, writeBuffer);

    // Simple Field (noOfDiagnosticInfos)
    writeSimpleField("noOfDiagnosticInfos", noOfDiagnosticInfos, writeSignedInt(writeBuffer, 32));

    // Array Field (diagnosticInfos)
    writeComplexTypeArrayField("diagnosticInfos", diagnosticInfos, writeBuffer);

    writeBuffer.popContext("BrowseNextResponse");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    BrowseNextResponse _value = this;

    // Simple field (responseHeader)
    lengthInBits += responseHeader.getLengthInBits();

    // Simple field (noOfResults)
    lengthInBits += 32;

    // Array field
    if (results != null) {
      int i = 0;
      for (ExtensionObjectDefinition element : results) {
        boolean last = ++i >= results.size();
        lengthInBits += element.getLengthInBits();
      }
    }

    // Simple field (noOfDiagnosticInfos)
    lengthInBits += 32;

    // Array field
    if (diagnosticInfos != null) {
      int i = 0;
      for (DiagnosticInfo element : diagnosticInfos) {
        boolean last = ++i >= diagnosticInfos.size();
        lengthInBits += element.getLengthInBits();
      }
    }

    return lengthInBits;
  }

  public static BrowseNextResponseBuilder staticParseBuilder(
      ReadBuffer readBuffer, String identifier) throws ParseException {
    readBuffer.pullContext("BrowseNextResponse");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    ExtensionObjectDefinition responseHeader =
        readSimpleField(
            "responseHeader",
            new DataReaderComplexDefault<>(
                () -> ExtensionObjectDefinition.staticParse(readBuffer, (String) ("394")),
                readBuffer));

    int noOfResults = readSimpleField("noOfResults", readSignedInt(readBuffer, 32));

    List<ExtensionObjectDefinition> results =
        readCountArrayField(
            "results",
            new DataReaderComplexDefault<>(
                () -> ExtensionObjectDefinition.staticParse(readBuffer, (String) ("524")),
                readBuffer),
            noOfResults);

    int noOfDiagnosticInfos = readSimpleField("noOfDiagnosticInfos", readSignedInt(readBuffer, 32));

    List<DiagnosticInfo> diagnosticInfos =
        readCountArrayField(
            "diagnosticInfos",
            new DataReaderComplexDefault<>(
                () -> DiagnosticInfo.staticParse(readBuffer), readBuffer),
            noOfDiagnosticInfos);

    readBuffer.closeContext("BrowseNextResponse");
    // Create the instance
    return new BrowseNextResponseBuilder(
        responseHeader, noOfResults, results, noOfDiagnosticInfos, diagnosticInfos);
  }

  public static class BrowseNextResponseBuilder
      implements ExtensionObjectDefinition.ExtensionObjectDefinitionBuilder {
    private final ExtensionObjectDefinition responseHeader;
    private final int noOfResults;
    private final List<ExtensionObjectDefinition> results;
    private final int noOfDiagnosticInfos;
    private final List<DiagnosticInfo> diagnosticInfos;

    public BrowseNextResponseBuilder(
        ExtensionObjectDefinition responseHeader,
        int noOfResults,
        List<ExtensionObjectDefinition> results,
        int noOfDiagnosticInfos,
        List<DiagnosticInfo> diagnosticInfos) {

      this.responseHeader = responseHeader;
      this.noOfResults = noOfResults;
      this.results = results;
      this.noOfDiagnosticInfos = noOfDiagnosticInfos;
      this.diagnosticInfos = diagnosticInfos;
    }

    public BrowseNextResponse build() {
      BrowseNextResponse browseNextResponse =
          new BrowseNextResponse(
              responseHeader, noOfResults, results, noOfDiagnosticInfos, diagnosticInfos);
      return browseNextResponse;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BrowseNextResponse)) {
      return false;
    }
    BrowseNextResponse that = (BrowseNextResponse) o;
    return (getResponseHeader() == that.getResponseHeader())
        && (getNoOfResults() == that.getNoOfResults())
        && (getResults() == that.getResults())
        && (getNoOfDiagnosticInfos() == that.getNoOfDiagnosticInfos())
        && (getDiagnosticInfos() == that.getDiagnosticInfos())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getResponseHeader(),
        getNoOfResults(),
        getResults(),
        getNoOfDiagnosticInfos(),
        getDiagnosticInfos());
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
