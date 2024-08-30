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

public class StatusRequestBinaryStateDeprecated extends StatusRequest implements Message {

  // Accessors for discriminator values.

  // Properties.
  protected final ApplicationIdContainer application;

  // Reserved Fields
  private Byte reservedField0;
  private Byte reservedField1;

  public StatusRequestBinaryStateDeprecated(byte statusType, ApplicationIdContainer application) {
    super(statusType);
    this.application = application;
  }

  public ApplicationIdContainer getApplication() {
    return application;
  }

  @Override
  protected void serializeStatusRequestChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("StatusRequestBinaryStateDeprecated");

    // Reserved Field (reserved)
    writeReservedField(
        "reserved",
        reservedField0 != null ? reservedField0 : (byte) 0xFA,
        writeByte(writeBuffer, 8));

    // Simple Field (application)
    writeSimpleEnumField(
        "application",
        "ApplicationIdContainer",
        application,
        new DataWriterEnumDefault<>(
            ApplicationIdContainer::getValue,
            ApplicationIdContainer::name,
            writeUnsignedShort(writeBuffer, 8)));

    // Reserved Field (reserved)
    writeReservedField(
        "reserved",
        reservedField1 != null ? reservedField1 : (byte) 0x00,
        writeByte(writeBuffer, 8));

    writeBuffer.popContext("StatusRequestBinaryStateDeprecated");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    StatusRequestBinaryStateDeprecated _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Reserved Field (reserved)
    lengthInBits += 8;

    // Simple field (application)
    lengthInBits += 8;

    // Reserved Field (reserved)
    lengthInBits += 8;

    return lengthInBits;
  }

  public static StatusRequestBuilder staticParseStatusRequestBuilder(ReadBuffer readBuffer)
      throws ParseException {
    readBuffer.pullContext("StatusRequestBinaryStateDeprecated");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    Byte reservedField0 = readReservedField("reserved", readByte(readBuffer, 8), (byte) 0xFA);

    ApplicationIdContainer application =
        readEnumField(
            "application",
            "ApplicationIdContainer",
            readEnum(ApplicationIdContainer::enumForValue, readUnsignedShort(readBuffer, 8)));

    Byte reservedField1 = readReservedField("reserved", readByte(readBuffer, 8), (byte) 0x00);

    readBuffer.closeContext("StatusRequestBinaryStateDeprecated");
    // Create the instance
    return new StatusRequestBinaryStateDeprecatedBuilderImpl(
        application, reservedField0, reservedField1);
  }

  public static class StatusRequestBinaryStateDeprecatedBuilderImpl
      implements StatusRequest.StatusRequestBuilder {
    private final ApplicationIdContainer application;
    private final Byte reservedField0;
    private final Byte reservedField1;

    public StatusRequestBinaryStateDeprecatedBuilderImpl(
        ApplicationIdContainer application, Byte reservedField0, Byte reservedField1) {
      this.application = application;
      this.reservedField0 = reservedField0;
      this.reservedField1 = reservedField1;
    }

    public StatusRequestBinaryStateDeprecated build(byte statusType) {
      StatusRequestBinaryStateDeprecated statusRequestBinaryStateDeprecated =
          new StatusRequestBinaryStateDeprecated(statusType, application);
      statusRequestBinaryStateDeprecated.reservedField0 = reservedField0;
      statusRequestBinaryStateDeprecated.reservedField1 = reservedField1;
      return statusRequestBinaryStateDeprecated;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof StatusRequestBinaryStateDeprecated)) {
      return false;
    }
    StatusRequestBinaryStateDeprecated that = (StatusRequestBinaryStateDeprecated) o;
    return (getApplication() == that.getApplication()) && super.equals(that) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getApplication());
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
