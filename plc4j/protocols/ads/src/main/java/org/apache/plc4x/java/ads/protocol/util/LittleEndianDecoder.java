/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.ads.protocol.util;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.ads.model.AdsDataType;
import org.apache.plc4x.java.base.messages.items.*;

import java.io.ByteArrayOutputStream;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.ZoneOffset;
import java.util.LinkedList;
import java.util.concurrent.TimeUnit;

// TODO: we might user ByteBuffer.wrap(buffer).order(ByteOrder.LITTLE_ENDIAN).putInt(port).asArray() etc
public class LittleEndianDecoder {

    private LittleEndianDecoder() {
        // Utility class
    }

    @SuppressWarnings("unchecked")
    public static FieldItem<?> decodeData(AdsDataType adsDataType, byte[] adsData) {
        ByteBuf wrappedBuffer = Unpooled.wrappedBuffer(adsData);
        switch (adsDataType) {
            case BIT: {
                LinkedList<Boolean> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readUnsignedByte();
                    values.offer(aByte != 0);
                }
                return new DefaultBooleanFieldItem(values.toArray(new Boolean[0]));
            }
            case BIT8: {
                LinkedList<Boolean> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readUnsignedByte();
                    values.offer(aByte != 0);
                }
                return new DefaultBooleanFieldItem(values.toArray(new Boolean[0]));
            }
            case BITARR8: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readUnsignedByte();
                    values.offer((long) aByte);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case BITARR16: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aLong = wrappedBuffer.readUnsignedShortLE();
                    values.offer(aLong);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case BITARR32: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aLong = wrappedBuffer.readUnsignedIntLE();
                    values.offer(aLong);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case INT8: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aLong = wrappedBuffer.readByte();
                    values.offer(aLong);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case INT16: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aLong = wrappedBuffer.readShortLE();
                    values.offer(aLong);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case INT32: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aLong = wrappedBuffer.readIntLE();
                    values.offer(aLong);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case INT64: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    byte[] bytes = new byte[8];
                    wrappedBuffer.readBytes(bytes);
                    BigInteger bigInteger = new BigInteger(bytes);
                    // TODO: potential dataloss here.
                    values.offer(bigInteger.longValue());
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case UINT8: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aLong = wrappedBuffer.readUnsignedByte();
                    values.offer(aLong);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case UINT16: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aLong = wrappedBuffer.readUnsignedShortLE();
                    values.offer(aLong);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case UINT32: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aLong = wrappedBuffer.readUnsignedIntLE();
                    values.offer(aLong);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case UINT64: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    byte[] bytes = new byte[64];
                    wrappedBuffer.readBytes(bytes);
                    BigInteger bigInteger = new BigInteger(ArrayUtils.add(bytes, (byte) 0x0));
                    // TODO: potential dataloss here.
                    values.offer(bigInteger.longValue());
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case FLOAT: {
                LinkedList<Double> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    double aLong = wrappedBuffer.readFloatLE();
                    values.offer(aLong);
                }
                return new DefaultFloatingPointFieldItem(values.toArray(new Double[0]));
            }
            case DOUBLE: {
                LinkedList<Double> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    double aLong = wrappedBuffer.readDoubleLE();
                    values.offer(aLong);
                }
                return new DefaultFloatingPointFieldItem(values.toArray(new Double[0]));
            }
            case BOOL: {
                LinkedList<Boolean> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readUnsignedByte();
                    values.offer(aByte != 0);
                }
                return new DefaultBooleanFieldItem(values.toArray(new Boolean[0]));
            }
            case BYTE: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedByte();
                    values.offer(aByte);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case WORD: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedShortLE();
                    values.offer(aByte);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case DWORD: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedIntLE();
                    values.offer(aByte);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case SINT: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readByte();
                    values.offer(aByte);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case USINT: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedByte();
                    values.offer(aByte);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case INT: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readShortLE();
                    values.offer(aByte);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case UINT: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedShortLE();
                    values.offer(aByte);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case DINT: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readIntLE();
                    values.offer(aByte);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case UDINT: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedIntLE();
                    values.offer(aByte);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case LINT: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readLongLE();
                    values.offer(aByte);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case ULINT: {
                LinkedList<BigInteger> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    byte[] bytes = new byte[64];
                    wrappedBuffer.readBytes(bytes);
                    BigInteger bigInteger = new BigInteger(ArrayUtils.add(bytes, (byte) 0x0));
                    values.offer(bigInteger);
                }
                return new DefaultBigIntegerFieldItem(values.toArray(new BigInteger[0]));
            }
            case REAL: {
                LinkedList<Double> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    double aByte = wrappedBuffer.readFloatLE();
                    values.offer(aByte);
                }
                return new DefaultFloatingPointFieldItem(values.toArray(new Double[0]));
            }
            case LREAL: {
                LinkedList<Double> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    double aByte = wrappedBuffer.readDoubleLE();
                    values.offer(aByte);
                }
                return new DefaultFloatingPointFieldItem(values.toArray(new Double[0]));
            }
            case STRING: {
                LinkedList<String> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    ByteArrayOutputStream os = new ByteArrayOutputStream();
                    Byte aByte;
                    while ((aByte = wrappedBuffer.readByte()) != 0x0) {
                        os.write(aByte);
                    }
                    values.offer(new String(os.toByteArray()));
                }
                return new DefaultStringFieldItem(values.toArray(new String[0]));
            }
            case TIME: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedIntLE();
                    values.offer(aByte);
                }
                return new DefaultIntegerFieldItem(values.toArray(new Long[0]));
            }
            case TIME_OF_DAY: {
                LinkedList<LocalTime> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedIntLE();
                    values.offer(LocalTime.ofNanoOfDay(TimeUnit.MILLISECONDS.toNanos(aByte)));
                }
                return new DefaultLocalTimeFieldItem(values.toArray(new LocalTime[0]));
            }
            case DATE: {
                LinkedList<LocalDate> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedIntLE();
                    values.offer(LocalDate.ofEpochDay(aByte));
                }
                return new DefaultLocalDateFieldItem(values.toArray(new LocalDate[0]));
            }
            case DATE_AND_TIME: {
                LinkedList<LocalDateTime> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedIntLE();
                    values.offer(LocalDateTime.ofEpochSecond(aByte, 0, ZoneOffset.UTC));
                }
                return new DefaultLocalDateTimeFieldItem(values.toArray(new LocalDateTime[0]));
            }
            case ARRAY: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case POINTER: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case ENUM: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case STRUCT: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case ALIAS: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case SUB_RANGE_DATA_TYPE: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            case UNKNOWN: {
                throw new NotImplementedException("not implemented yet " + adsDataType);
            }
            default:
                throw new IllegalArgumentException("Unsupported adsDataType " + adsDataType);
        }
    }

}
