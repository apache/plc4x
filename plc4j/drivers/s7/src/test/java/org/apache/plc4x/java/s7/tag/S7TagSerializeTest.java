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
package org.apache.plc4x.java.s7.tag;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Exercises {@code serialize()} on the three tag types so the WriteBuffer code paths
 * (writeUnsignedInt + writeString) are covered. Also verifies equals/hashCode/toString
 * branches that the parser tests don't hit.
 */
class S7TagSerializeTest {

    private static WriteBufferByteBased buffer() {
        return new WriteBufferByteBased(new byte[128],
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
    }

    @Test
    void s7Tag_serializeProducesNonEmptyOutput() throws Exception {
        S7Tag tag = S7Tag.of("%DB1.DBW20:INT");
        WriteBufferByteBased buf = buffer();
        tag.serialize(buf);
        assertTrue(buf.getPositionInBits() > 0);
    }

    @Test
    void s7Tag_serializeWithArrayCount() throws Exception {
        S7Tag tag = S7Tag.of("%DB1.DBW20:INT[5]");
        WriteBufferByteBased buf = buffer();
        tag.serialize(buf);
        assertEquals(5, tag.getNumberOfElements());
    }

    @Test
    void s7Tag_serializeFromInputArea() throws Exception {
        S7Tag tag = S7Tag.of("%I0.0:BOOL");
        WriteBufferByteBased buf = buffer();
        tag.serialize(buf);
        assertTrue(buf.getPositionInBits() > 0);
    }

    @Test
    void s7Tag_equalsAndHashCode() {
        S7Tag a = S7Tag.of("%DB1.DBW20:INT");
        S7Tag b = S7Tag.of("%DB1.DBW20:INT");
        S7Tag c = S7Tag.of("%DB1.DBW22:INT");
        assertEquals(a, a);
        assertEquals(a, b);
        assertNotEquals(a, c);
        assertNotEquals(a, "not a tag");
        assertEquals(a.hashCode(), b.hashCode());
    }

    @Test
    void s7Tag_toStringMentionsDataTypeAndArea() {
        S7Tag t = S7Tag.of("%DB10.DBD8:DINT[3]");
        String s = t.toString();
        assertTrue(s.contains("DINT"), () -> "expected dataType in toString: " + s);
        assertTrue(s.contains("DATA_BLOCKS") || s.contains("blockNumber"),
            () -> "expected area/block info in toString: " + s);
    }

    @Test
    void s7StringFixedLengthTag_serialize() throws Exception {
        S7StringFixedLengthTag tag = S7StringFixedLengthTag.of("%DB1.DB0:STRING(20)");
        WriteBufferByteBased buf = buffer();
        tag.serialize(buf);
        assertTrue(buf.getPositionInBits() > 0);
        assertEquals(20, tag.getStringLength());
    }

    @Test
    void s7StringFixedLengthTag_wstring() {
        S7StringFixedLengthTag tag = S7StringFixedLengthTag.of("%DB1.DB0:WSTRING(10)");
        assertEquals(10, tag.getStringLength());
    }

    @Test
    void s7StringFixedLengthTag_array() {
        S7StringFixedLengthTag tag = S7StringFixedLengthTag.of("%DB1.DB0:STRING(40)[3]");
        assertEquals(40, tag.getStringLength());
        assertEquals(3, tag.getNumberOfElements());
    }

    @Test
    void s7StringVarLengthTag_serialize() throws Exception {
        S7StringVarLengthTag tag = S7StringVarLengthTag.of("%DB1.DB0:STRING");
        WriteBufferByteBased buf = buffer();
        tag.serialize(buf);
        assertTrue(buf.getPositionInBits() > 0);
    }

    @Test
    void s7StringVarLengthTag_wstring() {
        S7StringVarLengthTag tag = S7StringVarLengthTag.of("%DB1.DB0:WSTRING");
        assertNotNull(tag);
    }

    @Test
    void s7StringVarLengthTag_equalsAndHashCode() {
        S7StringVarLengthTag a = S7StringVarLengthTag.of("%DB1.DB0:STRING");
        S7StringVarLengthTag b = S7StringVarLengthTag.of("%DB1.DB0:STRING");
        S7StringVarLengthTag c = S7StringVarLengthTag.of("%DB1.DB0:WSTRING");
        assertEquals(a, b);
        assertNotEquals(a, c);
        assertEquals(a.hashCode(), b.hashCode());
        assertNotNull(a.toString());
    }
}
