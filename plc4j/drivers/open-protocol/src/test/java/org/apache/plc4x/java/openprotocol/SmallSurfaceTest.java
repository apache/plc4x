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
package org.apache.plc4x.java.openprotocol;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.openprotocol.config.OpenProtocolConfiguration;
import org.apache.plc4x.java.openprotocol.config.OpenProtocolTcpTransportConfiguration;
import org.apache.plc4x.java.openprotocol.readwrite.utils.StaticHelper;
import org.apache.plc4x.java.openprotocol.tag.OpenProtocolTag;
import org.apache.plc4x.java.openprotocol.tag.OpenProtocolTagHandler;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

/**
 * Open Protocol is still a stub driver — the handwritten surface is intentionally
 * thin. These tests exercise the small classes that exist today; the placeholders
 * (Tag.of returns null, etc.) are locked in deliberately.
 */
class SmallSurfaceTest {

    @Test
    void tagOfSaysThereIsNoTagAddressingYet() {
        // Marker behaviour — kept here so anyone replacing this stub knows to
        // also remove this test (or replace it with a real parsing assertion).
        // It used to return null, which the driver's prepareTag() handed to the caller as
        // though it were a tag; the failure then surfaced later as a NullPointerException,
        // somewhere with nothing left to say about the address that caused it.
        assertThatThrownBy(() -> OpenProtocolTag.of("any"))
            .isInstanceOf(PlcInvalidTagException.class)
            .hasMessageContaining("any");
    }

    @Test
    void tagDelegatesArrayInfoAndValueTypeToPlcTagDefaults() {
        OpenProtocolTag tag = new OpenProtocolTag();
        // Defaults from PlcTag — non-null collections and a default value type.
        assertThat(tag.getArrayInfo()).isNotNull();
        assertThat(tag.getPlcValueType()).isNotNull();
        // But there is no address to render, and saying so beats reporting a blank one.
        assertThatThrownBy(tag::getAddressString)
            .isInstanceOf(UnsupportedOperationException.class);
    }

    @Test
    void handlerRejectsAnyTagAddress() {
        assertThatThrownBy(() -> new OpenProtocolTagHandler().parseTag("anything"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void handlerDoesNotSupportBrowsing() {
        assertThatThrownBy(() -> new OpenProtocolTagHandler().parseQuery("anything"))
            .isInstanceOf(UnsupportedOperationException.class);
    }

    @Test
    void configRequestTimeoutRoundtrips() {
        OpenProtocolConfiguration cfg = new OpenProtocolConfiguration();
        cfg.setRequestTimeout(5_000);
        assertThat(cfg.getRequestTimeout()).isEqualTo(5_000);
        assertThat(cfg.toString()).isNotEmpty();
    }

    @Test
    void tcpTransportDefaultPortIsPositive() {
        assertThat(new OpenProtocolTcpTransportConfiguration().getDefaultPort()).isPositive();
    }

    @Test
    void isBitSetReadsIndividualBitsFromDecimalString() {
        // 0b1010 = 10 → bits 1 and 3 set, bits 0 and 2 clear.
        assertThat(StaticHelper.isBitSet("10", 0)).isFalse();
        assertThat(StaticHelper.isBitSet("10", 1)).isTrue();
        assertThat(StaticHelper.isBitSet("10", 2)).isFalse();
        assertThat(StaticHelper.isBitSet("10", 3)).isTrue();
    }
}
