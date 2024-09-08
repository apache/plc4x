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

public class BACnetLandingCallStatusCommandDestination extends BACnetLandingCallStatusCommand
    implements Message {

  // Accessors for discriminator values.

  // Properties.
  protected final BACnetContextTagUnsignedInteger destination;

  public BACnetLandingCallStatusCommandDestination(
      BACnetTagHeader peekedTagHeader, BACnetContextTagUnsignedInteger destination) {
    super(peekedTagHeader);
    this.destination = destination;
  }

  public BACnetContextTagUnsignedInteger getDestination() {
    return destination;
  }

  @Override
  protected void serializeBACnetLandingCallStatusCommandChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("BACnetLandingCallStatusCommandDestination");

    // Simple Field (destination)
    writeSimpleField("destination", destination, writeComplex(writeBuffer));

    writeBuffer.popContext("BACnetLandingCallStatusCommandDestination");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    BACnetLandingCallStatusCommandDestination _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Simple field (destination)
    lengthInBits += destination.getLengthInBits();

    return lengthInBits;
  }

  public static BACnetLandingCallStatusCommandBuilder
      staticParseBACnetLandingCallStatusCommandBuilder(ReadBuffer readBuffer)
          throws ParseException {
    readBuffer.pullContext("BACnetLandingCallStatusCommandDestination");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    BACnetContextTagUnsignedInteger destination =
        readSimpleField(
            "destination",
            readComplex(
                () ->
                    (BACnetContextTagUnsignedInteger)
                        BACnetContextTag.staticParse(
                            readBuffer,
                            (short) (2),
                            (BACnetDataType) (BACnetDataType.UNSIGNED_INTEGER)),
                readBuffer));

    readBuffer.closeContext("BACnetLandingCallStatusCommandDestination");
    // Create the instance
    return new BACnetLandingCallStatusCommandDestinationBuilderImpl(destination);
  }

  public static class BACnetLandingCallStatusCommandDestinationBuilderImpl
      implements BACnetLandingCallStatusCommand.BACnetLandingCallStatusCommandBuilder {
    private final BACnetContextTagUnsignedInteger destination;

    public BACnetLandingCallStatusCommandDestinationBuilderImpl(
        BACnetContextTagUnsignedInteger destination) {
      this.destination = destination;
    }

    public BACnetLandingCallStatusCommandDestination build(BACnetTagHeader peekedTagHeader) {
      BACnetLandingCallStatusCommandDestination bACnetLandingCallStatusCommandDestination =
          new BACnetLandingCallStatusCommandDestination(peekedTagHeader, destination);
      return bACnetLandingCallStatusCommandDestination;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BACnetLandingCallStatusCommandDestination)) {
      return false;
    }
    BACnetLandingCallStatusCommandDestination that = (BACnetLandingCallStatusCommandDestination) o;
    return (getDestination() == that.getDestination()) && super.equals(that) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getDestination());
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
