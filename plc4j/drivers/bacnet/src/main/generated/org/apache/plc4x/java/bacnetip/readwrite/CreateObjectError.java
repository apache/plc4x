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
package org.apache.plc4x.java.bacnetip.readwrite;

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

public class CreateObjectError extends BACnetError implements Message {

  // Accessors for discriminator values.
  public BACnetConfirmedServiceChoice getErrorChoice() {
    return BACnetConfirmedServiceChoice.CREATE_OBJECT;
  }

  // Properties.
  protected final ErrorEnclosed errorType;
  protected final BACnetContextTagUnsignedInteger firstFailedElementNumber;

  public CreateObjectError(
      ErrorEnclosed errorType, BACnetContextTagUnsignedInteger firstFailedElementNumber) {
    super();
    this.errorType = errorType;
    this.firstFailedElementNumber = firstFailedElementNumber;
  }

  public ErrorEnclosed getErrorType() {
    return errorType;
  }

  public BACnetContextTagUnsignedInteger getFirstFailedElementNumber() {
    return firstFailedElementNumber;
  }

  @Override
  protected void serializeBACnetErrorChild(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("CreateObjectError");

    // Simple Field (errorType)
    writeSimpleField("errorType", errorType, writeComplex(writeBuffer));

    // Simple Field (firstFailedElementNumber)
    writeSimpleField(
        "firstFailedElementNumber", firstFailedElementNumber, writeComplex(writeBuffer));

    writeBuffer.popContext("CreateObjectError");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    CreateObjectError _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (errorType)
    lengthInBits += errorType.getLengthInBits();

    // Simple field (firstFailedElementNumber)
    lengthInBits += firstFailedElementNumber.getLengthInBits();

    return lengthInBits;
  }

  public static BACnetErrorBuilder staticParseBACnetErrorBuilder(
      ReadBuffer readBuffer, BACnetConfirmedServiceChoice errorChoice) throws ParseException {
    readBuffer.pullContext("CreateObjectError");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    ErrorEnclosed errorType =
        readSimpleField(
            "errorType",
            readComplex(() -> ErrorEnclosed.staticParse(readBuffer, (short) (0)), readBuffer));

    BACnetContextTagUnsignedInteger firstFailedElementNumber =
        readSimpleField(
            "firstFailedElementNumber",
            readComplex(
                () ->
                    (BACnetContextTagUnsignedInteger)
                        BACnetContextTag.staticParse(
                            readBuffer,
                            (short) (1),
                            (BACnetDataType) (BACnetDataType.UNSIGNED_INTEGER)),
                readBuffer));

    readBuffer.closeContext("CreateObjectError");
    // Create the instance
    return new CreateObjectErrorBuilderImpl(errorType, firstFailedElementNumber);
  }

  public static class CreateObjectErrorBuilderImpl implements BACnetError.BACnetErrorBuilder {
    private final ErrorEnclosed errorType;
    private final BACnetContextTagUnsignedInteger firstFailedElementNumber;

    public CreateObjectErrorBuilderImpl(
        ErrorEnclosed errorType, BACnetContextTagUnsignedInteger firstFailedElementNumber) {
      this.errorType = errorType;
      this.firstFailedElementNumber = firstFailedElementNumber;
    }

    public CreateObjectError build() {
      CreateObjectError createObjectError =
          new CreateObjectError(errorType, firstFailedElementNumber);
      return createObjectError;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof CreateObjectError)) {
      return false;
    }
    CreateObjectError that = (CreateObjectError) o;
    return (getErrorType() == that.getErrorType())
        && (getFirstFailedElementNumber() == that.getFirstFailedElementNumber())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getErrorType(), getFirstFailedElementNumber());
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
