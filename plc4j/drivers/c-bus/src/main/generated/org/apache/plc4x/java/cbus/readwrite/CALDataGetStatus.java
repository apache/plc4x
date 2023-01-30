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
package org.apache.plc4x.java.cbus.readwrite;

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

public class CALDataGetStatus extends CALData implements Message {

  // Accessors for discriminator values.

  // Properties.
  protected final Parameter paramNo;
  protected final short count;

  // Arguments.
  protected final RequestContext requestContext;

  public CALDataGetStatus(
      CALCommandTypeContainer commandTypeContainer,
      CALData additionalData,
      Parameter paramNo,
      short count,
      RequestContext requestContext) {
    super(commandTypeContainer, additionalData, requestContext);
    this.paramNo = paramNo;
    this.count = count;
    this.requestContext = requestContext;
  }

  public Parameter getParamNo() {
    return paramNo;
  }

  public short getCount() {
    return count;
  }

  @Override
  protected void serializeCALDataChild(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("CALDataGetStatus");

    // Simple Field (paramNo)
    writeSimpleEnumField(
        "paramNo",
        "Parameter",
        paramNo,
        new DataWriterEnumDefault<>(
            Parameter::getValue, Parameter::name, writeUnsignedShort(writeBuffer, 8)));

    // Simple Field (count)
    writeSimpleField("count", count, writeUnsignedShort(writeBuffer, 8));

    writeBuffer.popContext("CALDataGetStatus");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    CALDataGetStatus _value = this;

    // Simple field (paramNo)
    lengthInBits += 8;

    // Simple field (count)
    lengthInBits += 8;

    return lengthInBits;
  }

  public static CALDataBuilder staticParseCALDataBuilder(
      ReadBuffer readBuffer, RequestContext requestContext) throws ParseException {
    readBuffer.pullContext("CALDataGetStatus");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    Parameter paramNo =
        readEnumField(
            "paramNo",
            "Parameter",
            new DataReaderEnumDefault<>(Parameter::enumForValue, readUnsignedShort(readBuffer, 8)));

    short count = readSimpleField("count", readUnsignedShort(readBuffer, 8));

    readBuffer.closeContext("CALDataGetStatus");
    // Create the instance
    return new CALDataGetStatusBuilderImpl(paramNo, count, requestContext);
  }

  public static class CALDataGetStatusBuilderImpl implements CALData.CALDataBuilder {
    private final Parameter paramNo;
    private final short count;
    private final RequestContext requestContext;

    public CALDataGetStatusBuilderImpl(
        Parameter paramNo, short count, RequestContext requestContext) {
      this.paramNo = paramNo;
      this.count = count;
      this.requestContext = requestContext;
    }

    public CALDataGetStatus build(
        CALCommandTypeContainer commandTypeContainer,
        CALData additionalData,
        RequestContext requestContext) {
      CALDataGetStatus cALDataGetStatus =
          new CALDataGetStatus(
              commandTypeContainer, additionalData, paramNo, count, requestContext);
      return cALDataGetStatus;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof CALDataGetStatus)) {
      return false;
    }
    CALDataGetStatus that = (CALDataGetStatus) o;
    return (getParamNo() == that.getParamNo())
        && (getCount() == that.getCount())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getParamNo(), getCount());
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
