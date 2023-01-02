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

public class BACnetPropertyStatesLiftFault extends BACnetPropertyStates implements Message {

  // Accessors for discriminator values.

  // Properties.
  protected final BACnetLiftFaultTagged liftFault;

  public BACnetPropertyStatesLiftFault(
      BACnetTagHeader peekedTagHeader, BACnetLiftFaultTagged liftFault) {
    super(peekedTagHeader);
    this.liftFault = liftFault;
  }

  public BACnetLiftFaultTagged getLiftFault() {
    return liftFault;
  }

  @Override
  protected void serializeBACnetPropertyStatesChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("BACnetPropertyStatesLiftFault");

    // Simple Field (liftFault)
    writeSimpleField("liftFault", liftFault, new DataWriterComplexDefault<>(writeBuffer));

    writeBuffer.popContext("BACnetPropertyStatesLiftFault");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    BACnetPropertyStatesLiftFault _value = this;

    // Simple field (liftFault)
    lengthInBits += liftFault.getLengthInBits();

    return lengthInBits;
  }

  public static BACnetPropertyStatesLiftFaultBuilder staticParseBuilder(
      ReadBuffer readBuffer, Short peekedTagNumber) throws ParseException {
    readBuffer.pullContext("BACnetPropertyStatesLiftFault");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    BACnetLiftFaultTagged liftFault =
        readSimpleField(
            "liftFault",
            new DataReaderComplexDefault<>(
                () ->
                    BACnetLiftFaultTagged.staticParse(
                        readBuffer,
                        (short) (peekedTagNumber),
                        (TagClass) (TagClass.CONTEXT_SPECIFIC_TAGS)),
                readBuffer));

    readBuffer.closeContext("BACnetPropertyStatesLiftFault");
    // Create the instance
    return new BACnetPropertyStatesLiftFaultBuilder(liftFault);
  }

  public static class BACnetPropertyStatesLiftFaultBuilder
      implements BACnetPropertyStates.BACnetPropertyStatesBuilder {
    private final BACnetLiftFaultTagged liftFault;

    public BACnetPropertyStatesLiftFaultBuilder(BACnetLiftFaultTagged liftFault) {

      this.liftFault = liftFault;
    }

    public BACnetPropertyStatesLiftFault build(BACnetTagHeader peekedTagHeader) {
      BACnetPropertyStatesLiftFault bACnetPropertyStatesLiftFault =
          new BACnetPropertyStatesLiftFault(peekedTagHeader, liftFault);
      return bACnetPropertyStatesLiftFault;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BACnetPropertyStatesLiftFault)) {
      return false;
    }
    BACnetPropertyStatesLiftFault that = (BACnetPropertyStatesLiftFault) o;
    return (getLiftFault() == that.getLiftFault()) && super.equals(that) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getLiftFault());
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
