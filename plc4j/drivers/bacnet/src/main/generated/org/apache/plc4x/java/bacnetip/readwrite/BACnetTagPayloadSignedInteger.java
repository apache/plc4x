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

import java.math.BigInteger;
import java.time.*;
import java.util.*;
import org.apache.plc4x.java.api.exceptions.*;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.codegen.*;
import org.apache.plc4x.java.spi.codegen.fields.*;
import org.apache.plc4x.java.spi.codegen.io.*;
import org.apache.plc4x.java.spi.generation.*;

// Code generated by code-generation. DO NOT EDIT.

public class BACnetTagPayloadSignedInteger implements Message {

  // Properties.
  protected final Byte valueInt8;
  protected final Short valueInt16;
  protected final Integer valueInt24;
  protected final Integer valueInt32;
  protected final Long valueInt40;
  protected final Long valueInt48;
  protected final Long valueInt56;
  protected final Long valueInt64;

  // Arguments.
  protected final Long actualLength;

  public BACnetTagPayloadSignedInteger(
      Byte valueInt8,
      Short valueInt16,
      Integer valueInt24,
      Integer valueInt32,
      Long valueInt40,
      Long valueInt48,
      Long valueInt56,
      Long valueInt64,
      Long actualLength) {
    super();
    this.valueInt8 = valueInt8;
    this.valueInt16 = valueInt16;
    this.valueInt24 = valueInt24;
    this.valueInt32 = valueInt32;
    this.valueInt40 = valueInt40;
    this.valueInt48 = valueInt48;
    this.valueInt56 = valueInt56;
    this.valueInt64 = valueInt64;
    this.actualLength = actualLength;
  }

  public Byte getValueInt8() {
    return valueInt8;
  }

  public Short getValueInt16() {
    return valueInt16;
  }

  public Integer getValueInt24() {
    return valueInt24;
  }

  public Integer getValueInt32() {
    return valueInt32;
  }

  public Long getValueInt40() {
    return valueInt40;
  }

  public Long getValueInt48() {
    return valueInt48;
  }

  public Long getValueInt56() {
    return valueInt56;
  }

  public Long getValueInt64() {
    return valueInt64;
  }

  public boolean getIsInt8() {
    return (boolean) ((actualLength) == (1));
  }

  public boolean getIsInt16() {
    return (boolean) ((actualLength) == (2));
  }

  public boolean getIsInt24() {
    return (boolean) ((actualLength) == (3));
  }

  public boolean getIsInt32() {
    return (boolean) ((actualLength) == (4));
  }

  public boolean getIsInt40() {
    return (boolean) ((actualLength) == (5));
  }

  public boolean getIsInt48() {
    return (boolean) ((actualLength) == (6));
  }

  public boolean getIsInt56() {
    return (boolean) ((actualLength) == (7));
  }

  public boolean getIsInt64() {
    return (boolean) ((actualLength) == (8));
  }

  public BigInteger getActualValue() {
    Object o =
        ((getIsInt8())
            ? getValueInt8()
            : (((getIsInt16())
                ? getValueInt16()
                : (((getIsInt24())
                    ? getValueInt24()
                    : (((getIsInt32())
                        ? getValueInt32()
                        : (((getIsInt40())
                            ? getValueInt40()
                            : (((getIsInt48())
                                ? getValueInt48()
                                : (((getIsInt56()) ? getValueInt56() : getValueInt64())))))))))))));
    if (o instanceof BigInteger) return (BigInteger) o;
    return BigInteger.valueOf(((Number) o).longValue());
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("BACnetTagPayloadSignedInteger");

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    boolean isInt8 = getIsInt8();
    writeBuffer.writeVirtual("isInt8", isInt8);

    // Optional Field (valueInt8) (Can be skipped, if the value is null)
    writeOptionalField("valueInt8", valueInt8, writeSignedByte(writeBuffer, 8), getIsInt8());

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    boolean isInt16 = getIsInt16();
    writeBuffer.writeVirtual("isInt16", isInt16);

    // Optional Field (valueInt16) (Can be skipped, if the value is null)
    writeOptionalField("valueInt16", valueInt16, writeSignedShort(writeBuffer, 16), getIsInt16());

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    boolean isInt24 = getIsInt24();
    writeBuffer.writeVirtual("isInt24", isInt24);

    // Optional Field (valueInt24) (Can be skipped, if the value is null)
    writeOptionalField("valueInt24", valueInt24, writeSignedInt(writeBuffer, 24), getIsInt24());

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    boolean isInt32 = getIsInt32();
    writeBuffer.writeVirtual("isInt32", isInt32);

    // Optional Field (valueInt32) (Can be skipped, if the value is null)
    writeOptionalField("valueInt32", valueInt32, writeSignedInt(writeBuffer, 32), getIsInt32());

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    boolean isInt40 = getIsInt40();
    writeBuffer.writeVirtual("isInt40", isInt40);

    // Optional Field (valueInt40) (Can be skipped, if the value is null)
    writeOptionalField("valueInt40", valueInt40, writeSignedLong(writeBuffer, 40), getIsInt40());

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    boolean isInt48 = getIsInt48();
    writeBuffer.writeVirtual("isInt48", isInt48);

    // Optional Field (valueInt48) (Can be skipped, if the value is null)
    writeOptionalField("valueInt48", valueInt48, writeSignedLong(writeBuffer, 48), getIsInt48());

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    boolean isInt56 = getIsInt56();
    writeBuffer.writeVirtual("isInt56", isInt56);

    // Optional Field (valueInt56) (Can be skipped, if the value is null)
    writeOptionalField("valueInt56", valueInt56, writeSignedLong(writeBuffer, 56), getIsInt56());

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    boolean isInt64 = getIsInt64();
    writeBuffer.writeVirtual("isInt64", isInt64);

    // Optional Field (valueInt64) (Can be skipped, if the value is null)
    writeOptionalField("valueInt64", valueInt64, writeSignedLong(writeBuffer, 64), getIsInt64());

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    BigInteger actualValue = getActualValue();
    writeBuffer.writeVirtual("actualValue", actualValue);

    writeBuffer.popContext("BACnetTagPayloadSignedInteger");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    BACnetTagPayloadSignedInteger _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // A virtual field doesn't have any in- or output.

    // Optional Field (valueInt8)
    if (valueInt8 != null) {
      lengthInBits += 8;
    }

    // A virtual field doesn't have any in- or output.

    // Optional Field (valueInt16)
    if (valueInt16 != null) {
      lengthInBits += 16;
    }

    // A virtual field doesn't have any in- or output.

    // Optional Field (valueInt24)
    if (valueInt24 != null) {
      lengthInBits += 24;
    }

    // A virtual field doesn't have any in- or output.

    // Optional Field (valueInt32)
    if (valueInt32 != null) {
      lengthInBits += 32;
    }

    // A virtual field doesn't have any in- or output.

    // Optional Field (valueInt40)
    if (valueInt40 != null) {
      lengthInBits += 40;
    }

    // A virtual field doesn't have any in- or output.

    // Optional Field (valueInt48)
    if (valueInt48 != null) {
      lengthInBits += 48;
    }

    // A virtual field doesn't have any in- or output.

    // Optional Field (valueInt56)
    if (valueInt56 != null) {
      lengthInBits += 56;
    }

    // A virtual field doesn't have any in- or output.

    // Optional Field (valueInt64)
    if (valueInt64 != null) {
      lengthInBits += 64;
    }

    // A virtual field doesn't have any in- or output.

    return lengthInBits;
  }

  public static BACnetTagPayloadSignedInteger staticParse(ReadBuffer readBuffer, Long actualLength)
      throws ParseException {
    readBuffer.pullContext("BACnetTagPayloadSignedInteger");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    boolean isInt8 = readVirtualField("isInt8", boolean.class, (actualLength) == (1));

    Byte valueInt8 = readOptionalField("valueInt8", readSignedByte(readBuffer, 8), isInt8);
    boolean isInt16 = readVirtualField("isInt16", boolean.class, (actualLength) == (2));

    Short valueInt16 = readOptionalField("valueInt16", readSignedShort(readBuffer, 16), isInt16);
    boolean isInt24 = readVirtualField("isInt24", boolean.class, (actualLength) == (3));

    Integer valueInt24 = readOptionalField("valueInt24", readSignedInt(readBuffer, 24), isInt24);
    boolean isInt32 = readVirtualField("isInt32", boolean.class, (actualLength) == (4));

    Integer valueInt32 = readOptionalField("valueInt32", readSignedInt(readBuffer, 32), isInt32);
    boolean isInt40 = readVirtualField("isInt40", boolean.class, (actualLength) == (5));

    Long valueInt40 = readOptionalField("valueInt40", readSignedLong(readBuffer, 40), isInt40);
    boolean isInt48 = readVirtualField("isInt48", boolean.class, (actualLength) == (6));

    Long valueInt48 = readOptionalField("valueInt48", readSignedLong(readBuffer, 48), isInt48);
    boolean isInt56 = readVirtualField("isInt56", boolean.class, (actualLength) == (7));

    Long valueInt56 = readOptionalField("valueInt56", readSignedLong(readBuffer, 56), isInt56);
    boolean isInt64 = readVirtualField("isInt64", boolean.class, (actualLength) == (8));

    Long valueInt64 = readOptionalField("valueInt64", readSignedLong(readBuffer, 64), isInt64);
    // Validation
    if (!((((((((isInt8) || (isInt16)) || (isInt24)) || (isInt32)) || (isInt40)) || (isInt48))
            || (isInt56))
        || (isInt64))) {
      throw new ParseValidationException("unmapped integer length");
    }
    BigInteger actualValue =
        readVirtualField(
            "actualValue",
            BigInteger.class,
            ((isInt8)
                ? valueInt8
                : (((isInt16)
                    ? valueInt16
                    : (((isInt24)
                        ? valueInt24
                        : (((isInt32)
                            ? valueInt32
                            : (((isInt40)
                                ? valueInt40
                                : (((isInt48)
                                    ? valueInt48
                                    : (((isInt56) ? valueInt56 : valueInt64))))))))))))));

    readBuffer.closeContext("BACnetTagPayloadSignedInteger");
    // Create the instance
    BACnetTagPayloadSignedInteger _bACnetTagPayloadSignedInteger;
    _bACnetTagPayloadSignedInteger =
        new BACnetTagPayloadSignedInteger(
            valueInt8,
            valueInt16,
            valueInt24,
            valueInt32,
            valueInt40,
            valueInt48,
            valueInt56,
            valueInt64,
            actualLength);
    return _bACnetTagPayloadSignedInteger;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BACnetTagPayloadSignedInteger)) {
      return false;
    }
    BACnetTagPayloadSignedInteger that = (BACnetTagPayloadSignedInteger) o;
    return (getValueInt8() == that.getValueInt8())
        && (getValueInt16() == that.getValueInt16())
        && (getValueInt24() == that.getValueInt24())
        && (getValueInt32() == that.getValueInt32())
        && (getValueInt40() == that.getValueInt40())
        && (getValueInt48() == that.getValueInt48())
        && (getValueInt56() == that.getValueInt56())
        && (getValueInt64() == that.getValueInt64())
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        getValueInt8(),
        getValueInt16(),
        getValueInt24(),
        getValueInt32(),
        getValueInt40(),
        getValueInt48(),
        getValueInt56(),
        getValueInt64());
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
