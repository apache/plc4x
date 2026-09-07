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
package org.apache.plc4x.java.utils.testutils.driver.internal.api;

import org.apache.plc4x.java.spi.buffers.xmlbased.ReadBufferXmlBased;
import org.apache.plc4x.java.spi.buffers.xmlbased.WriteBufferXmlBased;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayInputStream;
import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.*;

class TestTagTest {

    @Test
    void testConstructorAndGetters() {
        TestTag tag = new TestTag("myTag", "holding-register:1:REAL");
        assertEquals("myTag", tag.getName());
        assertEquals("holding-register:1:REAL", tag.getAddress());
    }

    @Test
    void testGetLengthInBytes() {
        TestTag tag = new TestTag("test", "address");
        assertEquals(0, tag.getLengthInBytes());
    }

    @Test
    void testGetLengthInBits() {
        TestTag tag = new TestTag("test", "address");
        assertEquals(0, tag.getLengthInBits());
    }

    @Test
    void testSerialize() throws Exception {
        TestTag tag = new TestTag("tagName", "tagAddress");
        WriteBufferXmlBased writeBuffer = new WriteBufferXmlBased();

        assertDoesNotThrow(() -> tag.serialize(writeBuffer));
        String xml = writeBuffer.getXmlString();
        assertNotNull(xml);
        assertTrue(xml.contains("TestTag"));
    }

    @Test
    void testStaticParse() throws Exception {
        // The XML format requires specific dataType attributes for ReadBufferXmlBased
        String xml = "<TestTag><name dataType=\"string\" bitLength=\"512\">myTag</name><address dataType=\"string\" bitLength=\"512\">myAddress</address></TestTag>";
        ReadBufferXmlBased readBuffer = new ReadBufferXmlBased(
            new ByteArrayInputStream(xml.getBytes(StandardCharsets.UTF_8)));

        TestTag tag = TestTag.staticParse(readBuffer);
        assertNotNull(tag);
        assertEquals("myTag", tag.getName());
        assertEquals("myAddress", tag.getAddress());
    }
}
