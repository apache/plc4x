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
package org.apache.plc4x.java.profinet.context;

import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.profinet.gsdml.ProfinetISO15745Profile;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.util.Objects;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Slot-range handling for vendor-supplied GSDML files.
 * <p>
 * A GSDML crosses a provenance boundary: it is downloaded from a vendor site or handed over by an
 * integrator, and dropped into the configured directory. The slot range it declares sizes an array
 * directly, so the file gets to choose an allocation unless the driver bounds it. PROFINET numbers
 * slots within a 16 bit space, so a range beyond that describes no device that can exist.
 */
class GsdmlSlotBoundsTest {

    /** The value the shipped test fixture declares, i.e. a range a real device would use. */
    private static final String PLAUSIBLE_RANGE = "0..32";

    private static ProfinetDeviceContext contextFor(String physicalSlots) throws Exception {
        String gsdml = fixture().replace(
            "PhysicalSlots=\"" + PLAUSIBLE_RANGE + "\"",
            "PhysicalSlots=\"" + physicalSlots + "\"");

        ProfinetDeviceContext context = new ProfinetDeviceContext();
        context.setDeviceAccess("PLC4X_1");
        context.setSubModules("");
        context.setGsdFile(new XmlMapper().readValue(gsdml, ProfinetISO15745Profile.class));
        return context;
    }

    private static String fixture() throws IOException {
        try (InputStream in = Objects.requireNonNull(
            GsdmlSlotBoundsTest.class.getClassLoader().getResourceAsStream("gsdml.xml"))) {
            return new String(in.readAllBytes(), StandardCharsets.ISO_8859_1);
        }
    }

    /**
     * Ranges that no PROFINET device can have. The upper end of the 16 bit slot space is 0x7FFF, so
     * everything here describes more slots than the protocol can address - and each one sizes an
     * array before anything else happens.
     */
    @ParameterizedTest
    @ValueSource(strings = {"0..100000", "0..1000000", "0..2147483647"})
    @DisplayName("a slot range beyond the protocol's space is rejected")
    void implausibleSlotRangeIsRejected(String physicalSlots) {
        PlcConnectionException thrown = assertThrows(PlcConnectionException.class,
            () -> contextFor(physicalSlots),
            "a GSDML declaring " + physicalSlots + " slots must be rejected before it sizes anything");

        assertTrue(thrown.getMessage().contains(physicalSlots.substring(3)),
            "the error should name the offending value, was: " + thrown.getMessage());
    }

    /**
     * A slot count too large for {@code int} at all. This reaches {@code Integer.parseInt} before
     * anything else looks at it, so without a guard it surfaces as a {@link NumberFormatException}
     * from inside the connect path rather than as a connection error naming the file.
     */
    @ParameterizedTest
    @ValueSource(strings = {"0..99999999999", "0..2147483648", "0..9999999999999999999"})
    @DisplayName("a slot count that does not fit an int is a connection error, not a parse crash")
    void oversizedSlotCountIsAConnectionError(String physicalSlots) {
        assertThrows(PlcConnectionException.class, () -> contextFor(physicalSlots),
            "an unparseable slot count must be reported as a malformed GSDML");
    }

    /**
     * The bound must not have been drawn so tightly that real files stop working. The fixture is a
     * working device description, and the top of the addressable slot space must remain usable.
     */
    @Test
    @DisplayName("plausible slot ranges still load")
    void plausibleSlotRangesStillLoad() {
        assertDoesNotThrow(() -> contextFor(PLAUSIBLE_RANGE));
        assertDoesNotThrow(() -> contextFor("0..32767"),
            "the top of the 16 bit slot space is legal and must keep working");
    }
}
