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

package org.apache.plc4x.java.spi.generation;

import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.core.JsonFactoryBuilder;
import com.fasterxml.jackson.core.JsonGenerator;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.math.BigDecimal;
import java.math.BigInteger;

public class WriteBufferJsonBased implements WriteBuffer, BufferCommons {

    ByteArrayOutputStream byteArrayOutputStream;

    JsonGenerator generator;

    int depth = 0;

    boolean doRenderAttr = true;

    public WriteBufferJsonBased() {
        byteArrayOutputStream = new ByteArrayOutputStream();
        JsonFactory jsonFactory = new JsonFactoryBuilder()
            .build();
        try {
            generator = jsonFactory.createGenerator(byteArrayOutputStream).useDefaultPrettyPrinter();
        } catch (IOException e) {
            throw new PlcRuntimeException(e);
        }
    }

    @Override
    public int getPos() {
        return 0;
    }

    @Override
    public void pushContext(String logicalName) {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            if (depth == 0) {
                generator.writeStartObject();
            }
            depth++;
            // TODO: check if we need to write a array
            generator.writeObjectFieldStart(logicalName);
        } catch (IOException e) {
            throw new PlcRuntimeException(e);
        }
    }

    @Override
    public void writeBit(String logicalName, boolean value) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            writeAttr(logicalName, rwBitKey, 1);
            generator.writeBooleanField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeUnsignedByte(String logicalName, int bitLength, byte value) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            writeAttr(logicalName, rwUintKey, bitLength);
            generator.writeNumberField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeUnsignedShort(String logicalName, int bitLength, short value) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            writeAttr(logicalName, rwUintKey, bitLength);
            generator.writeNumberField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeUnsignedInt(String logicalName, int bitLength, int value) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            writeAttr(logicalName, rwUintKey, bitLength);
            generator.writeNumberField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeUnsignedLong(String logicalName, int bitLength, long value) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            writeAttr(logicalName, rwUintKey, bitLength);
            generator.writeNumberField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeUnsignedBigInteger(String logicalName, int bitLength, BigInteger value) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            writeAttr(logicalName, rwUintKey, bitLength);
            generator.writeNumberField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeByte(String logicalName, int bitLength, byte value) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            writeAttr(logicalName, rwIntKey, bitLength);
            generator.writeNumberField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeShort(String logicalName, int bitLength, short value) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            writeAttr(logicalName, rwIntKey, bitLength);
            generator.writeNumberField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeInt(String logicalName, int bitLength, int value) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            writeAttr(logicalName, rwIntKey, bitLength);
            generator.writeNumberField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeLong(String logicalName, int bitLength, long value) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            writeAttr(logicalName, rwIntKey, bitLength);
            generator.writeNumberField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeBigInteger(String logicalName, int bitLength, BigInteger value) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            writeAttr(logicalName, rwIntKey, bitLength);
            generator.writeNumberField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeFloat(String logicalName, float value, int bitsExponent, int bitsMantissa) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            int bitLength = (value < 0 ? 1 : 0) + bitsExponent + bitsMantissa;
            writeAttr(logicalName, rwFloatKey, bitLength);
            generator.writeNumberField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeDouble(String logicalName, double value, int bitsExponent, int bitsMantissa) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            int bitLength = (value < 0 ? 1 : 0) + bitsExponent + bitsMantissa;
            writeAttr(logicalName, rwFloatKey, bitLength);
            generator.writeNumberField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeBigDecimal(String logicalName, int bitLength, BigDecimal value) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            writeAttr(logicalName, rwFloatKey, bitLength);
            generator.writeNumberField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void writeString(String logicalName, int bitLength, String encoding, String value) throws ParseException {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            writeAttr(logicalName, rwStringKey, bitLength);
            generator.writeStringField(String.format("%s__plc4x_%s", logicalName, rwEncodingKey), encoding);
            generator.writeStringField(logicalName, value);
        } catch (IOException e) {
            throw new ParseException("error writing",e);
        }
    }

    @Override
    public void popContext(String logicalName) {
        logicalName = sanitizeLogicalName(logicalName);
        try {
            // TODO: check if we need to write a array
            generator.writeEndObject();
            depth--;
            if (depth == 0) {
                generator.writeEndObject();
            }
        } catch (IOException e) {
            throw new PlcRuntimeException(e);
        }
    }

    public String getJsonString() {
        try {
            generator.close();
            return byteArrayOutputStream.toString("UTF-8");
        } catch (IOException e) {
            throw new PlcRuntimeException(e);
        }
    }

    private void writeAttr(String logicalName , String dataType , int bitLength ) throws IOException {
        if( !doRenderAttr) {
            return;
        }
        generator.writeStringField(String.format("%s__plc4x_%s", logicalName, rwDataTypeKey), dataType);
        generator.writeNumberField(String.format("%s__plc4x_%s", logicalName, rwBitLengthKey), bitLength);
    }
}
