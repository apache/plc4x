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

class TestTagRequestTest {

    @Test
    void testConstructorAndGetters() {
        TestTag tag1 = new TestTag("tag1", "address1");
        TestTag tag2 = new TestTag("tag2", "address2");
        TestTagRequest tagRequest = new TestTagRequest(new TestTag[]{tag1, tag2});

        assertNotNull(tagRequest);
        assertEquals(2, tagRequest.getTags().length);
        assertEquals("tag1", tagRequest.getTags()[0].getName());
        assertEquals("tag2", tagRequest.getTags()[1].getName());
    }

    @Test
    void testGetLengthInBytes() {
        TestTag tag = new TestTag("test", "address");
        TestTagRequest tagRequest = new TestTagRequest(new TestTag[]{tag});

        assertEquals(0, tagRequest.getLengthInBytes());
    }

    @Test
    void testGetLengthInBits() {
        TestTag tag = new TestTag("test", "address");
        TestTagRequest tagRequest = new TestTagRequest(new TestTag[]{tag});

        assertEquals(0, tagRequest.getLengthInBits());
    }

    @Test
    void testSerialize() throws Exception {
        TestTag tag = new TestTag("myTag", "myAddress");
        TestTagRequest tagRequest = new TestTagRequest(new TestTag[]{tag});
        WriteBufferXmlBased writeBuffer = new WriteBufferXmlBased();

        assertDoesNotThrow(() -> tagRequest.serialize(writeBuffer));
        String xml = writeBuffer.getXmlString();
        assertNotNull(xml);
        assertTrue(xml.contains("TestTagRequest"));
    }

    @Test
    void testEmptyTags() {
        TestTagRequest tagRequest = new TestTagRequest(new TestTag[]{});

        assertEquals(0, tagRequest.getTags().length);
    }
}
