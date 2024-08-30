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

public class TelephonyDataDialOutFailure extends TelephonyData implements Message {

  // Accessors for discriminator values.

  // Properties.
  protected final DialOutFailureReason reason;

  public TelephonyDataDialOutFailure(
      TelephonyCommandTypeContainer commandTypeContainer,
      byte argument,
      DialOutFailureReason reason) {
    super(commandTypeContainer, argument);
    this.reason = reason;
  }

  public DialOutFailureReason getReason() {
    return reason;
  }

  @Override
  protected void serializeTelephonyDataChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("TelephonyDataDialOutFailure");

    // Simple Field (reason)
    writeSimpleEnumField(
        "reason",
        "DialOutFailureReason",
        reason,
        new DataWriterEnumDefault<>(
            DialOutFailureReason::getValue,
            DialOutFailureReason::name,
            writeUnsignedShort(writeBuffer, 8)));

    writeBuffer.popContext("TelephonyDataDialOutFailure");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    TelephonyDataDialOutFailure _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (reason)
    lengthInBits += 8;

    return lengthInBits;
  }

  public static TelephonyDataBuilder staticParseTelephonyDataBuilder(ReadBuffer readBuffer)
      throws ParseException {
    readBuffer.pullContext("TelephonyDataDialOutFailure");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    DialOutFailureReason reason =
        readEnumField(
            "reason",
            "DialOutFailureReason",
            readEnum(DialOutFailureReason::enumForValue, readUnsignedShort(readBuffer, 8)));

    readBuffer.closeContext("TelephonyDataDialOutFailure");
    // Create the instance
    return new TelephonyDataDialOutFailureBuilderImpl(reason);
  }

  public static class TelephonyDataDialOutFailureBuilderImpl
      implements TelephonyData.TelephonyDataBuilder {
    private final DialOutFailureReason reason;

    public TelephonyDataDialOutFailureBuilderImpl(DialOutFailureReason reason) {
      this.reason = reason;
    }

    public TelephonyDataDialOutFailure build(
        TelephonyCommandTypeContainer commandTypeContainer, byte argument) {
      TelephonyDataDialOutFailure telephonyDataDialOutFailure =
          new TelephonyDataDialOutFailure(commandTypeContainer, argument, reason);
      return telephonyDataDialOutFailure;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof TelephonyDataDialOutFailure)) {
      return false;
    }
    TelephonyDataDialOutFailure that = (TelephonyDataDialOutFailure) o;
    return (getReason() == that.getReason()) && super.equals(that) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getReason());
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
