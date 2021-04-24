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
package org.apache.plc4x.test.mapper;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.SerializerProvider;
import com.fasterxml.jackson.databind.deser.std.StdDeserializer;
import com.fasterxml.jackson.databind.module.SimpleModule;
import com.fasterxml.jackson.databind.ser.std.StdSerializer;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;

import java.io.IOException;

/**
 * Custom module for Jackson to handle types in the way we wish - turning byte and byte[] into hex strings.
 * @deprecated replaced by custom serializer
 */
@Deprecated
public class TestSuiteMappingModule extends SimpleModule {

    public TestSuiteMappingModule() {
        addSerializer(new ByteArraySerializer());
        addSerializer(new ByteSerializer());
        addDeserializer(byte[].class, new ByteArrayDeserializer());
        addDeserializer(Byte.class, new ByteDeserializer());
    }

    // Specific serializers, these are quite compact so let them stay here for now

    static class ByteArraySerializer extends StdSerializer<byte[]> {

        protected ByteArraySerializer() {
            super(byte[].class);
        }

        @Override
        public void serialize(byte[] value, JsonGenerator gen, SerializerProvider provider) throws IOException {
            gen.writeString(Hex.encodeHexString(value, false));
        }
    }

    static class ByteSerializer extends StdSerializer<Byte> {
        protected ByteSerializer() {
            super(Byte.class);
        }

        @Override
        public void serialize(Byte value, JsonGenerator gen, SerializerProvider provider) throws IOException {
            gen.writeString(Hex.encodeHexString(new byte[]{value}, false));
        }
    }

    // Specific serializers, these are quite compact so let them stay here for now

    static class ByteArrayDeserializer extends StdDeserializer<byte[]> {

        protected ByteArrayDeserializer() {
            super(byte[].class);
        }

        @Override
        public byte[] deserialize(JsonParser jsonParser, DeserializationContext deserializationContext) throws IOException, JsonProcessingException {
            try {
                return Hex.decodeHex(jsonParser.getText());
            } catch (DecoderException e) {
                throw new IOException(e);
            }
        }

    }

    static class ByteDeserializer extends StdDeserializer<Byte> {
        protected ByteDeserializer() {
            super(Byte.class);
        }

        @Override
        public Byte deserialize(JsonParser jsonParser, DeserializationContext deserializationContext) throws IOException, JsonProcessingException {
            try {
                return Hex.decodeHex(jsonParser.getText())[0];
            } catch (DecoderException e) {
                throw new IOException(e);
            }
        }

    }
}
