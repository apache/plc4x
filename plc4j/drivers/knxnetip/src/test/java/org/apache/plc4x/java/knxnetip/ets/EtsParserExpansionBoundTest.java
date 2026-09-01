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
package org.apache.plc4x.java.knxnetip.ets;

import org.apache.plc4x.java.knxnetip.ets.EtsParser.BoundedInputStream;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A zip entry declares how large it will be once expanded, and that declaration is written by
 * whoever made the file. So the bytes are counted as they arrive as well.
 */
class EtsParserExpansionBoundTest {

    private static InputStream of(int bytes) {
        return new ByteArrayInputStream(new byte[bytes]);
    }

    @Test
    void anEntryWithinItsBudgetIsReadWhole() throws IOException {
        try (BoundedInputStream in = new BoundedInputStream(of(1000), 4096, "knx_master.xml")) {
            assertEquals(1000, in.readAllBytes().length);
        }
    }

    @Test
    void anEntryThatExpandsPastTheBudgetStops() {
        IOException thrown = assertThrows(IOException.class, () -> {
            try (BoundedInputStream in = new BoundedInputStream(of(8192), 4096, "0.xml")) {
                in.readAllBytes();
            }
        });
        assertTrue(thrown.getMessage().contains("0.xml"),
            "the failure should name the entry, but was: " + thrown.getMessage());
        assertTrue(thrown.getMessage().contains("whatever its header claimed"),
            "and should say why it is not trusting the header");
    }

    @Test
    void theBudgetIsEnforcedByteByByteToo() {
        assertThrows(IOException.class, () -> {
            try (BoundedInputStream in = new BoundedInputStream(of(64), 8, "single-byte-reads")) {
                for (int i = 0; i < 64; i++) {
                    in.read();
                }
            }
        });
    }

    @Test
    void readingExactlyTheBudgetIsAllowed() {
        assertDoesNotThrow(() -> {
            try (BoundedInputStream in = new BoundedInputStream(of(4096), 4096, "exact")) {
                assertEquals(4096, in.readAllBytes().length);
            }
        });
    }
}
