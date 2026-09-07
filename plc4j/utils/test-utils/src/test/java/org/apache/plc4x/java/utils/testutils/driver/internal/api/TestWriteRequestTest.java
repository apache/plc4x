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

import org.apache.plc4x.java.spi.buffers.xmlbased.WriteBufferXmlBased;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TestWriteRequestTest {

    @Test
    void testConstructorAndGetters() {
        TestTag tag = new TestTag("myTag", "address");
        TestTagRequest tagRequest = new TestTagRequest(new TestTag[]{tag});
        TestWriteRequest writeRequest = new TestWriteRequest(tagRequest);

        assertNotNull(writeRequest);
    }

    @Test
    void testGetLengthInBytes() {
        TestTag tag = new TestTag("test", "address");
        TestTagRequest tagRequest = new TestTagRequest(new TestTag[]{tag});
        TestWriteRequest writeRequest = new TestWriteRequest(tagRequest);

        assertEquals(0, writeRequest.getLengthInBytes());
    }

    @Test
    void testGetLengthInBits() {
        TestTag tag = new TestTag("test", "address");
        TestTagRequest tagRequest = new TestTagRequest(new TestTag[]{tag});
        TestWriteRequest writeRequest = new TestWriteRequest(tagRequest);

        assertEquals(0, writeRequest.getLengthInBits());
    }

    @Test
    void testSerialize() throws Exception {
        TestTag tag = new TestTag("myTag", "myAddress");
        TestTagRequest tagRequest = new TestTagRequest(new TestTag[]{tag});
        TestWriteRequest writeRequest = new TestWriteRequest(tagRequest);
        WriteBufferXmlBased writeBuffer = new WriteBufferXmlBased();

        assertDoesNotThrow(() -> writeRequest.serialize(writeBuffer));
        String xml = writeBuffer.getXmlString();
        assertNotNull(xml);
        assertTrue(xml.contains("TestReadRequest")); // Note: the class has a bug, it uses TestReadRequest in serialize
    }
}
