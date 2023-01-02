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
package org.apache.plc4x.java.eip.readwrite;

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

public class GetAttributeAllResponse extends CipService implements Message {

  // Accessors for discriminator values.
  public Short getService() {
    return (short) 0x01;
  }

  public Boolean getResponse() {
    return (boolean) true;
  }

  public Boolean getConnected() {
    return false;
  }

  // Properties.
  protected final short status;
  protected final short extStatus;
  protected final CIPAttributes attributes;

  // Arguments.
  protected final Integer serviceLen;
  protected final IntegerEncoding order;
  // Reserved Fields
  private Short reservedField0;

  public GetAttributeAllResponse(
      short status,
      short extStatus,
      CIPAttributes attributes,
      Integer serviceLen,
      IntegerEncoding order) {
    super(serviceLen, order);
    this.status = status;
    this.extStatus = extStatus;
    this.attributes = attributes;
    this.serviceLen = serviceLen;
    this.order = order;
  }

  public short getStatus() {
    return status;
  }

  public short getExtStatus() {
    return extStatus;
  }

  public CIPAttributes getAttributes() {
    return attributes;
  }

  @Override
  protected void serializeCipServiceChild(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("GetAttributeAllResponse");

    // Reserved Field (reserved)
    writeReservedField(
        "reserved",
        reservedField0 != null ? reservedField0 : (short) 0x00,
        writeUnsignedShort(writeBuffer, 8));

    // Simple Field (status)
    writeSimpleField(
        "status",
        status,
        writeUnsignedShort(writeBuffer, 8),
        WithOption.WithByteOrder(
            (((order) == (IntegerEncoding.BIG_ENDIAN))
                ? ByteOrder.BIG_ENDIAN
                : ByteOrder.LITTLE_ENDIAN)));

    // Simple Field (extStatus)
    writeSimpleField(
        "extStatus",
        extStatus,
        writeUnsignedShort(writeBuffer, 8),
        WithOption.WithByteOrder(
            (((order) == (IntegerEncoding.BIG_ENDIAN))
                ? ByteOrder.BIG_ENDIAN
                : ByteOrder.LITTLE_ENDIAN)));

    // Optional Field (attributes) (Can be skipped, if the value is null)
    writeOptionalField(
        "attributes",
        attributes,
        new DataWriterComplexDefault<>(writeBuffer),
        (((serviceLen) - (4))) > (0));

    writeBuffer.popContext("GetAttributeAllResponse");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    GetAttributeAllResponse _value = this;

    // Reserved Field (reserved)
    lengthInBits += 8;

    // Simple field (status)
    lengthInBits += 8;

    // Simple field (extStatus)
    lengthInBits += 8;

    // Optional Field (attributes)
    if (attributes != null) {
      lengthInBits += attributes.getLengthInBits();
    }

    return lengthInBits;
  }

  public static GetAttributeAllResponseBuilder staticParseBuilder(
      ReadBuffer readBuffer, Boolean connected, Integer serviceLen, IntegerEncoding order)
      throws ParseException {
    readBuffer.pullContext("GetAttributeAllResponse");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    Short reservedField0 =
        readReservedField(
            "reserved",
            readUnsignedShort(readBuffer, 8),
            (short) 0x00,
            WithOption.WithByteOrder(
                (((order) == (IntegerEncoding.BIG_ENDIAN))
                    ? ByteOrder.BIG_ENDIAN
                    : ByteOrder.LITTLE_ENDIAN)));

    short status =
        readSimpleField(
            "status",
            readUnsignedShort(readBuffer, 8),
            WithOption.WithByteOrder(
                (((order) == (IntegerEncoding.BIG_ENDIAN))
                    ? ByteOrder.BIG_ENDIAN
                    : ByteOrder.LITTLE_ENDIAN)));

    short extStatus =
        readSimpleField(
            "extStatus",
            readUnsignedShort(readBuffer, 8),
            WithOption.WithByteOrder(
                (((order) == (IntegerEncoding.BIG_ENDIAN))
                    ? ByteOrder.BIG_ENDIAN
                    : ByteOrder.LITTLE_ENDIAN)));

    CIPAttributes attributes =
        readOptionalField(
            "attributes",
            new DataReaderComplexDefault<>(
                () -> CIPAttributes.staticParse(readBuffer, (int) ((serviceLen) - (4))),
                readBuffer),
            (((serviceLen) - (4))) > (0),
            WithOption.WithByteOrder(
                (((order) == (IntegerEncoding.BIG_ENDIAN))
                    ? ByteOrder.BIG_ENDIAN
                    : ByteOrder.LITTLE_ENDIAN)));

    readBuffer.closeContext("GetAttributeAllResponse");
    // Create the instance
    return new GetAttributeAllResponseBuilder(
        status, extStatus, attributes, serviceLen, order, reservedField0);
  }

  public static class GetAttributeAllResponseBuilder implements CipService.CipServiceBuilder {
    private final short status;
    private final short extStatus;
    private final CIPAttributes attributes;
    private final Integer serviceLen;
    private final IntegerEncoding order;
    private final Short reservedField0;

    public GetAttributeAllResponseBuilder(
        short status,
        short extStatus,
        CIPAttributes attributes,
        Integer serviceLen,
        IntegerEncoding order,
        Short reservedField0) {
      this.status = status;
      this.extStatus = extStatus;
      this.attributes = attributes;
      this.serviceLen = serviceLen;
      this.order = order;
      this.reservedField0 = reservedField0;
    }

    public GetAttributeAllResponse build(Integer serviceLen, IntegerEncoding order) {

      GetAttributeAllResponse getAttributeAllResponse =
          new GetAttributeAllResponse(status, extStatus, attributes, serviceLen, order);
      getAttributeAllResponse.reservedField0 = reservedField0;
      return getAttributeAllResponse;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof GetAttributeAllResponse)) {
      return false;
    }
    GetAttributeAllResponse that = (GetAttributeAllResponse) o;
    return (getStatus() == that.getStatus())
        && (getExtStatus() == that.getExtStatus())
        && (getAttributes() == that.getAttributes())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getStatus(), getExtStatus(), getAttributes());
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
