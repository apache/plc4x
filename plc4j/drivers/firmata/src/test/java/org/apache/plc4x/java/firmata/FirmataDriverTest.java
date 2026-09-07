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
package org.apache.plc4x.java.firmata;

import org.apache.plc4x.java.firmata.tag.FirmataTag;
import org.apache.plc4x.java.firmata.tag.FirmataTagAnalog;
import org.apache.plc4x.java.firmata.tag.FirmataTagDigital;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertTrue;

class FirmataDriverTest {

    private FirmataDriver driver;

    @BeforeEach
    void setUp() {
        driver = new FirmataDriver();
    }

    @Test
    void protocolIdentity() {
        assertEquals("firmata", driver.getProtocolCode());
        assertEquals("Firmata", driver.getProtocolName());
    }

    @Test
    void metadataAdvertisesSerialAndTcp() {
        var metadata = driver.getMetadata();
        assertTrue(metadata.getDefaultTransportCode().isPresent());
        assertEquals("serial", metadata.getDefaultTransportCode().get());
        assertTrue(metadata.getSupportedTransportCodes().contains("serial"));
        assertTrue(metadata.getSupportedTransportCodes().contains("tcp"));
        assertTrue(metadata.getSupportedTransportCodes().contains("test"));
    }

    @Test
    void discoveryIsNotSupported() {
        // Firmata is push-only, has no auto-discovery.
        assertFalse(driver.getMetadata().isDiscoverySupported());
    }

    @Test
    void prepareTagReturnsTheRightSubtype() {
        FirmataTag digital = driver.prepareTag("digital:7");
        assertInstanceOf(FirmataTagDigital.class, digital);

        FirmataTag analog = driver.prepareTag("analog:3");
        assertInstanceOf(FirmataTagAnalog.class, analog);
    }

}
