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

public class SetMonitoringModeResponse extends ExtensionObjectDefinition implements Message {

  // Accessors for discriminator values.
  public Integer getExtensionId() {
    return (int) 772;
  }

  // Properties.
  protected final ResponseHeader responseHeader;
  protected final List<StatusCode> results;
  protected final List<DiagnosticInfo> diagnosticInfos;

  public SetMonitoringModeResponse(
      ResponseHeader responseHeader,
      List<StatusCode> results,
      List<DiagnosticInfo> diagnosticInfos) {
    super();
    this.responseHeader = responseHeader;
    this.results = results;
    this.diagnosticInfos = diagnosticInfos;
  }

  public ResponseHeader getResponseHeader() {
    return responseHeader;
  }

  public List<StatusCode> getResults() {
    return results;
  }

  public List<DiagnosticInfo> getDiagnosticInfos() {
    return diagnosticInfos;
  }

  @Override
  protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("SetMonitoringModeResponse");

    // Simple Field (responseHeader)
    writeSimpleField("responseHeader", responseHeader, writeComplex(writeBuffer));

    // Implicit Field (noOfResults) (Used for parsing, but its value is not stored as it's
    // implicitly given by the objects content)
    int noOfResults = (int) ((((getResults()) == (null)) ? -(1) : COUNT(getResults())));
    writeImplicitField("noOfResults", noOfResults, writeSignedInt(writeBuffer, 32));

    // Array Field (results)
    writeComplexTypeArrayField("results", results, writeBuffer);

    // Implicit Field (noOfDiagnosticInfos) (Used for parsing, but its value is not stored as it's
    // implicitly given by the objects content)
    int noOfDiagnosticInfos =
        (int) ((((getDiagnosticInfos()) == (null)) ? -(1) : COUNT(getDiagnosticInfos())));
    writeImplicitField("noOfDiagnosticInfos", noOfDiagnosticInfos, writeSignedInt(writeBuffer, 32));

    // Array Field (diagnosticInfos)
    writeComplexTypeArrayField("diagnosticInfos", diagnosticInfos, writeBuffer);

    writeBuffer.popContext("SetMonitoringModeResponse");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    SetMonitoringModeResponse _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (responseHeader)
    lengthInBits += responseHeader.getLengthInBits();

    // Implicit Field (noOfResults)
    lengthInBits += 32;

    // Array field
    if (results != null) {
      int i = 0;
      for (StatusCode element : results) {
        ThreadLocalHelper.lastItemThreadLocal.set(++i >= results.size());
        lengthInBits += element.getLengthInBits();
      }
    }

    // Implicit Field (noOfDiagnosticInfos)
    lengthInBits += 32;

    // Array field
    if (diagnosticInfos != null) {
      int i = 0;
      for (DiagnosticInfo element : diagnosticInfos) {
        ThreadLocalHelper.lastItemThreadLocal.set(++i >= diagnosticInfos.size());
        lengthInBits += element.getLengthInBits();
      }
    }

    return lengthInBits;
  }

  public static ExtensionObjectDefinitionBuilder staticParseExtensionObjectDefinitionBuilder(
      ReadBuffer readBuffer, Integer extensionId) throws ParseException {
    readBuffer.pullContext("SetMonitoringModeResponse");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    ResponseHeader responseHeader =
        readSimpleField(
            "responseHeader",
            readComplex(
                () ->
                    (ResponseHeader) ExtensionObjectDefinition.staticParse(readBuffer, (int) (394)),
                readBuffer));

    int noOfResults = readImplicitField("noOfResults", readSignedInt(readBuffer, 32));

    List<StatusCode> results =
        readCountArrayField(
            "results",
            readComplex(() -> StatusCode.staticParse(readBuffer), readBuffer),
            noOfResults);

    int noOfDiagnosticInfos =
        readImplicitField("noOfDiagnosticInfos", readSignedInt(readBuffer, 32));

    List<DiagnosticInfo> diagnosticInfos =
        readCountArrayField(
            "diagnosticInfos",
            readComplex(() -> DiagnosticInfo.staticParse(readBuffer), readBuffer),
            noOfDiagnosticInfos);

    readBuffer.closeContext("SetMonitoringModeResponse");
    // Create the instance
    return new SetMonitoringModeResponseBuilderImpl(responseHeader, results, diagnosticInfos);
  }

  public static class SetMonitoringModeResponseBuilderImpl
      implements ExtensionObjectDefinition.ExtensionObjectDefinitionBuilder {
    private final ResponseHeader responseHeader;
    private final List<StatusCode> results;
    private final List<DiagnosticInfo> diagnosticInfos;

    public SetMonitoringModeResponseBuilderImpl(
        ResponseHeader responseHeader,
        List<StatusCode> results,
        List<DiagnosticInfo> diagnosticInfos) {
      this.responseHeader = responseHeader;
      this.results = results;
      this.diagnosticInfos = diagnosticInfos;
    }

    public SetMonitoringModeResponse build() {
      SetMonitoringModeResponse setMonitoringModeResponse =
          new SetMonitoringModeResponse(responseHeader, results, diagnosticInfos);
      return setMonitoringModeResponse;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof SetMonitoringModeResponse)) {
      return false;
    }
    SetMonitoringModeResponse that = (SetMonitoringModeResponse) o;
    return (getResponseHeader() == that.getResponseHeader())
        && (getResults() == that.getResults())
        && (getDiagnosticInfos() == that.getDiagnosticInfos())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getResponseHeader(), getResults(), getDiagnosticInfos());
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
