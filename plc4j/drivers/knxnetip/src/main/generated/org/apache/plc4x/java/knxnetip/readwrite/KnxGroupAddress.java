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
package org.apache.plc4x.java.knxnetip.readwrite;

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

public abstract class KnxGroupAddress implements Message {

  // Abstract accessors for discriminator values.
  public abstract Byte getNumLevels();

  public KnxGroupAddress() {
    super();
  }

  protected abstract void serializeKnxGroupAddressChild(WriteBuffer writeBuffer)
      throws SerializationException;

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("KnxGroupAddress");

    // Switch field (Serialize the sub-type)
    serializeKnxGroupAddressChild(writeBuffer);

    writeBuffer.popContext("KnxGroupAddress");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    KnxGroupAddress _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Length of sub-type elements will be added by sub-type...

    return lengthInBits;
  }

  public static KnxGroupAddress staticParse(ReadBuffer readBuffer, Object... args)
      throws ParseException {
    PositionAware positionAware = readBuffer;
    if ((args == null) || (args.length != 1)) {
      throw new PlcRuntimeException(
          "Wrong number of arguments, expected 1, but got " + args.length);
    }
    Byte numLevels;
    if (args[0] instanceof Byte) {
      numLevels = (Byte) args[0];
    } else if (args[0] instanceof String) {
      numLevels = Byte.valueOf((String) args[0]);
    } else {
      throw new PlcRuntimeException(
          "Argument 0 expected to be of type Byte or a string which is parseable but was "
              + args[0].getClass().getName());
    }
    return staticParse(readBuffer, numLevels);
  }

  public static KnxGroupAddress staticParse(ReadBuffer readBuffer, Byte numLevels)
      throws ParseException {
    readBuffer.pullContext("KnxGroupAddress");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    KnxGroupAddressBuilder builder = null;
    if (EvaluationHelper.equals(numLevels, (byte) 1)) {
      builder = KnxGroupAddressFreeLevel.staticParseKnxGroupAddressBuilder(readBuffer, numLevels);
    } else if (EvaluationHelper.equals(numLevels, (byte) 2)) {
      builder = KnxGroupAddress2Level.staticParseKnxGroupAddressBuilder(readBuffer, numLevels);
    } else if (EvaluationHelper.equals(numLevels, (byte) 3)) {
      builder = KnxGroupAddress3Level.staticParseKnxGroupAddressBuilder(readBuffer, numLevels);
    }
    if (builder == null) {
      throw new ParseException(
          "Unsupported case for discriminated type"
              + " parameters ["
              + "numLevels="
              + numLevels
              + "]");
    }

    readBuffer.closeContext("KnxGroupAddress");
    // Create the instance
    KnxGroupAddress _knxGroupAddress = builder.build();
    return _knxGroupAddress;
  }

  public interface KnxGroupAddressBuilder {
    KnxGroupAddress build();
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof KnxGroupAddress)) {
      return false;
    }
    KnxGroupAddress that = (KnxGroupAddress) o;
    return true;
  }

  @Override
  public int hashCode() {
    return Objects.hash();
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