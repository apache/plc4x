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
package org.apache.plc4x.java.iec608705104;

import org.apache.plc4x.java.iec608705104.tag.Iec608705104Tag;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertTrue;

class Iec60870514PlcDriverTest {

    private final Iec60870514PlcDriver driver = new Iec60870514PlcDriver();

    @Test
    void protocolIdentity() {
        assertEquals("iec-60870-5-104", driver.getProtocolCode());
        assertEquals("IEC 60870-5-104", driver.getProtocolName());
    }

    @Test
    void metadataAdvertisesTcp() {
        var metadata = driver.getMetadata();
        assertTrue(metadata.getDefaultTransportCode().isPresent());
        assertEquals("tcp", metadata.getDefaultTransportCode().get());
        assertTrue(metadata.getSupportedTransportCodes().contains("tcp"));
        assertTrue(metadata.getSupportedTransportCodes().contains("test"));
    }

    @Test
    void prepareTagReturnsIecTag() {
        // The grammar isn't wired yet — the handler returns a placeholder
        // tag for any input. This test pins the wiring through the driver.
        assertInstanceOf(Iec608705104Tag.class, driver.prepareTag("anything"));
    }

}
