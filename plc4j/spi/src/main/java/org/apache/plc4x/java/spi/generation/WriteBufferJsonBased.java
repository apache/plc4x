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
package org.apache.plc4x.java.spi.generation;

import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.core.JsonFactoryBuilder;
import com.fasterxml.jackson.core.JsonGenerator;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.util.Optional;

public class WriteBufferJsonBased implements WriteBuffer, BufferCommons, AutoCloseable {

    public static final String PLC4X_ATTRIBUTE_FORMAT = "%s__plc4x_%s";

    private final ByteArrayOutputStream byteArrayOutputStream;

    private final JsonGenerator generator;

    private int pos = 1;
    private int depth = 0;

    private final boolean doRenderAttr;

    public WriteBufferJsonBased() {
        this(true);
    }

    public WriteBufferJsonBased(boolean doRenderAttr) {
        this.doRenderAttr = doRenderAttr;
        byteArrayOutputStream = new ByteArrayOutputStream();
        JsonFactory jsonFactory = new JsonFactoryBuilder()
            .build();
        try {
            generator = jsonFactory.createGenerator(byteArrayOutputStream);
            // Usually this is chained onto above creating of the generator but then sonar thinks this never gets closed
            this.generator.useDefaultPrettyPrinter();
        } catch (IOException e) {
            throw new PlcRuntimeException(e);
        }
    }

    @Override
    public ByteOrder getByteOrder() {
        // NO OP
        return ByteOrder.BIG_ENDIAN;
    }

    @Override
    public void setByteOrder(ByteOrder byteOrder) {
        // NO OP
    }

    @Override
    public int getPos() {
        return pos / 8;
    }

    @Override
    public void pushContext(String logicalName, WithWriterArgs... writerArgs) {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        try {
            if (depth == 0) {
                generator.writeStartObject();
            }
            depth++;
            if (isToBeRenderedAsList(writerArgs)) {
                generator.writeArrayFieldStart(sanitizedLogicalName);
            } else {
                if (generator.getOutputContext().inArray()) {
                    generator.writeStartObject();
                }
                generator.writeObjectFieldStart(sanitizedLogicalName);
            }
        } catch (IOException e) {
            throw new PlcRuntimeException(e);
        }
    }

    @Override
    public void writeBit(String logicalName, boolean value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwBitKey, 1, writerArgs);
            generator.writeBooleanField(sanitizedLogicalName, value);
        });
        move(1);
    }

    @Override
    public void writeByte(String logicalName, byte value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwByteKey, 8, writerArgs);
            generator.writeStringField(sanitizedLogicalName, String.format("0x%02x", value));
        });
        move(8);
    }

    @Override
    public void writeByteArray(String logicalName, byte[] bytes, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        StringBuilder hexString = new StringBuilder("0x");
        for (byte aByte : bytes) {
            hexString.append(String.format("%02x", aByte));
        }
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwByteKey, bytes.length * 8, writerArgs);
            generator.writeStringField(sanitizedLogicalName, hexString.toString());
        });
        move(8 * bytes.length);
    }

    @Override
    public void writeUnsignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwUintKey, bitLength, writerArgs);
            generator.writeNumberField(sanitizedLogicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void writeUnsignedShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwUintKey, bitLength, writerArgs);
            generator.writeNumberField(sanitizedLogicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void writeUnsignedInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwUintKey, bitLength, writerArgs);
            generator.writeNumberField(sanitizedLogicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void writeUnsignedLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwUintKey, bitLength, writerArgs);
            generator.writeNumberField(sanitizedLogicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void writeUnsignedBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwUintKey, bitLength, writerArgs);
            generator.writeNumberField(sanitizedLogicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void writeSignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwIntKey, bitLength, writerArgs);
            generator.writeNumberField(sanitizedLogicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void writeShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwIntKey, bitLength, writerArgs);
            generator.writeNumberField(sanitizedLogicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void writeInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwIntKey, bitLength, writerArgs);
            generator.writeNumberField(sanitizedLogicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void writeLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwIntKey, bitLength, writerArgs);
            generator.writeNumberField(sanitizedLogicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void writeBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwIntKey, bitLength, writerArgs);
            generator.writeNumberField(sanitizedLogicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void writeFloat(String logicalName, int bitLength, float value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwFloatKey, bitLength, writerArgs);
            generator.writeNumberField(logicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void writeDouble(String logicalName, int bitLength, double value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwFloatKey, bitLength, writerArgs);
            generator.writeNumberField(sanitizedLogicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void writeBigDecimal(String logicalName, int bitLength, BigDecimal value, WithWriterArgs... writerArgs) throws SerializationException {
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwFloatKey, bitLength, writerArgs);
            generator.writeNumberField(sanitizedLogicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void writeString(String logicalName, int bitLength, String value, WithWriterArgs... writerArgs) throws SerializationException {
        String encoding = extractEncoding(writerArgs).orElse("UTF-8");
        final String sanitizedLogicalName = sanitizeLogicalName(logicalName);
        wrapIfNecessary(() -> {
            writeAttr(sanitizedLogicalName, rwStringKey, bitLength, writerArgs);
            generator.writeStringField(String.format(PLC4X_ATTRIBUTE_FORMAT, sanitizedLogicalName, rwEncodingKey), encoding);
            generator.writeStringField(sanitizedLogicalName, value);
        });
        move(bitLength);
    }

    @Override
    public void popContext(String logicalName, WithWriterArgs... writerArgs) {
        try {
            if (isToBeRenderedAsList(writerArgs)) {
                generator.writeEndArray();
            } else {
                generator.writeEndObject();
                if (generator.getOutputContext().getParent().inArray()) {
                    generator.writeEndObject();
                }
            }
            depth--;
            if (depth == 0) {
                generator.writeEndObject();
                generator.flush();
            }
        } catch (IOException e) {
            throw new PlcRuntimeException(e);
        }
    }

    public void wrapIfNecessary(RunWrapped runnable) throws SerializationException {
        boolean inArray = generator.getOutputContext().inArray();
        if (inArray) {
            try {
                generator.writeStartObject();
            } catch (IOException e) {
                throw new SerializationException("Error opening wrap", e);
            }
        }
        try {
            runnable.run();
        } catch (IOException e) {
            throw new SerializationException("Error running wrap", e);
        }
        if (inArray) {
            try {
                generator.writeEndObject();
            } catch (IOException e) {
                throw new SerializationException("Error closing wrap", e);
            }
        }
    }

    @Override
    public void close() throws Exception {
        generator.close();
    }

    @FunctionalInterface
    private interface RunWrapped {
        void run() throws IOException;
    }

    public String getJsonString() {
        try {
            return byteArrayOutputStream.toString(StandardCharsets.UTF_8.name());
        } catch (IOException e) {
            throw new PlcRuntimeException(e);
        }
    }

    private void writeAttr(String logicalName, String dataType, int bitLength, WithWriterArgs... writerArgs) throws IOException {
        if (!doRenderAttr) {
            return;
        }
        generator.writeStringField(String.format(PLC4X_ATTRIBUTE_FORMAT, logicalName, rwDataTypeKey), dataType);
        generator.writeNumberField(String.format(PLC4X_ATTRIBUTE_FORMAT, logicalName, rwBitLengthKey), bitLength);
        Optional<String> stringRepresentation = extractAdditionalStringRepresentation(writerArgs);
        if (stringRepresentation.isPresent()) {
            generator.writeStringField(String.format(PLC4X_ATTRIBUTE_FORMAT, logicalName, rwStringRepresentationKey), stringRepresentation.get());
        }
    }

    private void move(int bits) {
        pos += bits;
    }
}
