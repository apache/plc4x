/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.spi.fields.fields.reader;

import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.fields.TestFieldIoStubs.SimpleDataReader;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class FieldReaderChecksumTest {

    @Test
    void checksumMatchesReturnsValue() throws Exception {
        FieldReaderChecksum<Integer> reader = new FieldReaderChecksum<>();
        SimpleDataReader<Integer> dr = new SimpleDataReader<>(() -> 42);
        Integer val = reader.readChecksumField(dr, 42);
        assertEquals(42, val);
        assertEquals(1, dr.getReadCount());
    }

    @Test
    void checksumMismatchThrows() {
        FieldReaderChecksum<Integer> reader = new FieldReaderChecksum<>();
        SimpleDataReader<Integer> dr = new SimpleDataReader<>(() -> 41);
        assertThrows(BufferException.class, () -> reader.readChecksumField(dr, 42));
    }
}
