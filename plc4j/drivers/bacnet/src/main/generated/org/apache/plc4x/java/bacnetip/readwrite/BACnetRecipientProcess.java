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

public class BACnetRecipientProcess implements Message {

  // Properties.
  protected final BACnetRecipientEnclosed recipient;
  protected final BACnetContextTagUnsignedInteger processIdentifier;

  public BACnetRecipientProcess(
      BACnetRecipientEnclosed recipient, BACnetContextTagUnsignedInteger processIdentifier) {
    super();
    this.recipient = recipient;
    this.processIdentifier = processIdentifier;
  }

  public BACnetRecipientEnclosed getRecipient() {
    return recipient;
  }

  public BACnetContextTagUnsignedInteger getProcessIdentifier() {
    return processIdentifier;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("BACnetRecipientProcess");

    // Simple Field (recipient)
    writeSimpleField("recipient", recipient, writeComplex(writeBuffer));

    // Optional Field (processIdentifier) (Can be skipped, if the value is null)
    writeOptionalField("processIdentifier", processIdentifier, writeComplex(writeBuffer));

    writeBuffer.popContext("BACnetRecipientProcess");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    BACnetRecipientProcess _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (recipient)
    lengthInBits += recipient.getLengthInBits();

    // Optional Field (processIdentifier)
    if (processIdentifier != null) {
      lengthInBits += processIdentifier.getLengthInBits();
    }

    return lengthInBits;
  }

  public static BACnetRecipientProcess staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("BACnetRecipientProcess");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    BACnetRecipientEnclosed recipient =
        readSimpleField(
            "recipient",
            readComplex(
                () -> BACnetRecipientEnclosed.staticParse(readBuffer, (short) (0)), readBuffer));

    BACnetContextTagUnsignedInteger processIdentifier =
        readOptionalField(
            "processIdentifier",
            readComplex(
                () ->
                    (BACnetContextTagUnsignedInteger)
                        BACnetContextTag.staticParse(
                            readBuffer,
                            (short) (1),
                            (BACnetDataType) (BACnetDataType.UNSIGNED_INTEGER)),
                readBuffer));

    readBuffer.closeContext("BACnetRecipientProcess");
    // Create the instance
    BACnetRecipientProcess _bACnetRecipientProcess;
    _bACnetRecipientProcess = new BACnetRecipientProcess(recipient, processIdentifier);
    return _bACnetRecipientProcess;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BACnetRecipientProcess)) {
      return false;
    }
    BACnetRecipientProcess that = (BACnetRecipientProcess) o;
    return (getRecipient() == that.getRecipient())
        && (getProcessIdentifier() == that.getProcessIdentifier())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getRecipient(), getProcessIdentifier());
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
