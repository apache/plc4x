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
package org.apache.plc4x.java.amsads.attic.protocol.util;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.ads.readwrite.types.AdsDataType;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.apache.plc4x.java.api.value.*;

import java.io.ByteArrayOutputStream;
import java.math.BigInteger;
import java.util.LinkedList;

// TODO: we might user ByteBuffer.wrap(buffer).order(ByteOrder.LITTLE_ENDIAN).putInt(port).asArray() etc
public class LittleEndianDecoder {

    private LittleEndianDecoder() {
        // Utility class
    }

    @SuppressWarnings("unchecked")
    public static PlcValue decodeData(AdsDataType adsDataType, byte[] adsData) {
        ByteBuf wrappedBuffer = Unpooled.wrappedBuffer(adsData);
        switch (adsDataType) {
            case BIT: {
                LinkedList<Boolean> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readUnsignedByte();
                    values.offer(aByte != 0);
                }
                if(values.size() == 1) {
                    return new PlcBOOL(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Boolean[values.size()]));
                }
            }
            case BIT8: {
                LinkedList<Boolean> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readUnsignedByte();
                    values.offer(aByte != 0);
                }
                if(values.size() == 1) {
                    return new PlcBOOL(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Boolean[values.size()]));
                }
            }
            case BITARR8: {
                LinkedList<Short> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readUnsignedByte();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    // TODO: Double-Check this ...
                    return new PlcBOOL(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Short[values.size()]));
                }
            }
            case BITARR16: {
                LinkedList<Integer> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    int aLong = wrappedBuffer.readUnsignedShortLE();
                    values.offer(aLong);
                }
                if(values.size() == 1) {
                    // TODO: Double-Check this ...
                    return new PlcBOOL(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Integer[values.size()]));
                }
            }
            case BITARR32: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aLong = wrappedBuffer.readUnsignedIntLE();
                    values.offer(aLong);
                }
                if(values.size() == 1) {
                    // TODO: Double-Check this ...
                    return new PlcBOOL(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Long[values.size()]));
                }
            }
            case INT8: {
                LinkedList<Byte> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    byte aLong = wrappedBuffer.readByte();
                    values.offer(aLong);
                }
                if(values.size() == 1) {
                    return new PlcSINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Byte[values.size()]));
                }
            }
            case INT16: {
                LinkedList<Short> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aLong = wrappedBuffer.readShortLE();
                    values.offer(aLong);
                }
                if(values.size() == 1) {
                    return new PlcINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Short[values.size()]));
                }
            }
            case INT32: {
                LinkedList<Integer> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    int intLE = wrappedBuffer.readIntLE();
                    values.offer(intLE);
                }
                if(values.size() == 1) {
                    return new PlcDINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Integer[values.size()]));
                }
            }
            case INT64: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long longLE = wrappedBuffer.readLongLE();
                    values.offer(longLE);
                }
                if(values.size() == 1) {
                    return new PlcLINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Long[values.size()]));
                }
            }
            case UINT8: {
                LinkedList<Short> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aLong = wrappedBuffer.readUnsignedByte();
                    values.offer(aLong);
                }
                if(values.size() == 1) {
                    return new PlcUSINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Short[values.size()]));
                }
            }
            case UINT16: {
                LinkedList<Integer> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    int aLong = wrappedBuffer.readUnsignedShortLE();
                    values.offer(aLong);
                }
                if(values.size() == 1) {
                    return new PlcUINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Integer[values.size()]));
                }
            }
            case UINT32: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aLong = wrappedBuffer.readUnsignedIntLE();
                    values.offer(aLong);
                }
                if(values.size() == 1) {
                    return new PlcUDINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Long[values.size()]));
                }
            }
            case ULINT:
            case UINT64: {
                LinkedList<BigInteger> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    byte[] bytes = new byte[8];
                    wrappedBuffer.readBytes(bytes);
                    ArrayUtils.reverse(bytes);
                    BigInteger bigInteger = new BigInteger(ArrayUtils.insert(0, bytes, (byte) 0x0));
                    values.offer(bigInteger);
                }
                if(values.size() == 1) {
                    return new PlcULINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new BigInteger[values.size()]));
                }
            }
            case FLOAT: {
                LinkedList<Float> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    // FIXME: Commented out till we are able to update to the latest netty version.
                    float aLong = wrappedBuffer.readFloat();
                    //float aLong = wrappedBuffer.readFloatLE();
                    values.offer(aLong);
                }
                if(values.size() == 1) {
                    return new PlcREAL(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Float[values.size()]));
                }
            }
            case DOUBLE: {
                LinkedList<Double> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    // FIXME: Commented out till we are able to update to the latest netty version.
                    double aLong = wrappedBuffer.readDouble();
                    //double aLong = wrappedBuffer.readDoubleLE();
                    values.offer(aLong);
                }
                if(values.size() == 1) {
                    return new PlcLREAL(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Double[values.size()]));
                }
            }
            case BOOL: {
                LinkedList<Boolean> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readUnsignedByte();
                    values.offer(aByte != 0);
                }
                if(values.size() == 1) {
                    return new PlcBOOL(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Boolean[values.size()]));
                }
            }
            case BYTE: {
                LinkedList<Short> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readUnsignedByte();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    // TODO: Double-Check this ...
                    return new PlcBOOL(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Short[values.size()]));
                }
            }
            case WORD: {
                LinkedList<Integer> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    int aByte = wrappedBuffer.readUnsignedShortLE();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    // TODO: Double-Check this ...
                    return new PlcBOOL(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Integer[values.size()]));
                }
            }
            case DWORD: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedIntLE();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    // TODO: Double-Check this ...
                    return new PlcBOOL(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Long[values.size()]));
                }
            }
            case SINT: {
                LinkedList<Byte> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    byte aByte = wrappedBuffer.readByte();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    return new PlcSINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Byte[values.size()]));
                }
            }
            case USINT: {
                LinkedList<Short> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readUnsignedByte();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    return new PlcUSINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Short[values.size()]));
                }
            }
            case INT: {
                LinkedList<Short> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    short aByte = wrappedBuffer.readShortLE();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    return new PlcINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Short[values.size()]));
                }
            }
            case UINT: {
                LinkedList<Integer> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    int aByte = wrappedBuffer.readUnsignedShortLE();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    return new PlcUINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Integer[values.size()]));
                }
            }
            case DINT: {
                LinkedList<Integer> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    int aByte = wrappedBuffer.readIntLE();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    return new PlcDINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Integer[values.size()]));
                }
            }
            case UDINT: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedIntLE();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    return new PlcUDINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Long[values.size()]));
                }
            }
            case LINT: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readLongLE();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    return new PlcLINT(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Long[values.size()]));
                }
            }
            case REAL: {
                LinkedList<Float> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    // FIXME: Commented out till we are able to update to the latest netty version.
                    float aByte = wrappedBuffer.readFloat();
                    //float aByte = wrappedBuffer.readFloatLE();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    return new PlcREAL(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Float[values.size()]));
                }
            }
            case LREAL: {
                LinkedList<Double> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    // FIXME: Commented out till we are able to update to the latest netty version.
                    double aByte = wrappedBuffer.readDouble();
                    //double aByte = wrappedBuffer.readDoubleLE();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    return new PlcLREAL(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new Double[values.size()]));
                }
            }
            case STRING: {
                LinkedList<String> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    ByteArrayOutputStream os = new ByteArrayOutputStream();
                    byte aByte;
                    while ((aByte = wrappedBuffer.readByte()) != 0x0) {
                        os.write(aByte);
                    }
                    values.offer(new String(os.toByteArray()));
                }
                if(values.size() == 1) {
                    return new PlcSTRING(values.get(0));
                } else {
                    return IEC61131ValueHandler.newPlcValue(values.toArray(new String[values.size()]));
                }
            }
/*            case TIME: {
                LinkedList<Long> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedIntLE();
                    values.offer(aByte);
                }
                if(values.size() == 1) {
                    return new PlcLong(values.get(0));
                } else {
                    return new PlcList(values);
                }
            }
            case TIME_OF_DAY: {
                LinkedList<LocalTime> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedIntLE();
                    values.offer(LocalTime.ofNanoOfDay(TimeUnit.MILLISECONDS.toNanos(aByte)));
                }
                if(values.size() == 1) {
                    return new PlcTime(values.get(0));
                } else {
                    return new PlcList(values);
                }
            }
            case DATE: {
                LinkedList<LocalDate> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedIntLE();
                    // TODO: where to get the zone offset from
                    LocalDateTime localDateTime = LocalDateTime.ofEpochSecond(aByte, 0, ZoneOffset.UTC);
                    values.offer(localDateTime.toLocalDate());
                }
                if(values.size() == 1) {
                    return new PlcDate(values.get(0));
                } else {
                    return new PlcList(values);
                }
            }
            case DATE_AND_TIME: {
                LinkedList<LocalDateTime> values = new LinkedList<>();
                while (wrappedBuffer.isReadable()) {
                    long aByte = wrappedBuffer.readUnsignedIntLE();
                    // TODO: where to get the zone offset from
                    LocalDateTime localDateTime = LocalDateTime.ofEpochSecond(aByte, 0, ZoneOffset.UTC);
                    values.offer(localDateTime);
                }
                if(values.size() == 1) {
                    return new PlcDateTime(values.get(0));
                } else {
                    return new PlcList(values);
                }
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
            case UNKNOWN:*/
            default:
                throw new PlcUnsupportedDataTypeException("Unsupported adsDataType " + adsDataType);
        }
    }

}
