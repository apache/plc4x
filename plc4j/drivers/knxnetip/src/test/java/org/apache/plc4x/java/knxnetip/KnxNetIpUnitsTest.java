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
package org.apache.plc4x.java.knxnetip;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.knxnetip.configuration.KnxNetIpConfiguration;
import org.apache.plc4x.java.knxnetip.configuration.KnxNetIpUdpTransportConfiguration;
import org.apache.plc4x.java.knxnetip.ets.model.GroupAddress;
import org.apache.plc4x.java.knxnetip.tag.KnxNetIpQuery;
import org.apache.plc4x.java.knxnetip.tag.KnxNetIpTag;
import org.apache.plc4x.java.knxnetip.tag.KnxNetIpTagHandler;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.io.File;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class KnxNetIpUnitsTest {

    @Test
    void tagParses3LevelAddress() {
        KnxNetIpTag tag = KnxNetIpTag.of("1/2/3");
        assertThat(tag.getLevels()).isEqualTo(3);
        assertThat(tag.getMainGroup()).isEqualTo("1");
        assertThat(tag.getMiddleGroup()).isEqualTo("2");
        assertThat(tag.getSubGroup()).isEqualTo("3");
        assertThat(tag.getDptId()).isNull();
        assertThat(tag.getAddressString()).isEqualTo("1/2/3");
        assertThat(tag.getPlcValueType()).isEqualTo(PlcValueType.NULL);
        assertThat(tag.getArrayInfo()).isNotNull();
    }

    @Test
    void tagParses2And1LevelAddresses() {
        KnxNetIpTag two = KnxNetIpTag.of("12/3000");
        assertThat(two.getLevels()).isEqualTo(2);
        assertThat(two.getMiddleGroup()).isNull();
        assertThat(two.getSubGroup()).isEqualTo("3000");

        KnxNetIpTag one = KnxNetIpTag.of("12345");
        assertThat(one.getLevels()).isEqualTo(1);
        assertThat(one.getMainGroup()).isEqualTo("12345");
    }

    @Test
    void tagParsesDptSuffix() {
        KnxNetIpTag tag = KnxNetIpTag.of("1/2/3:DPT9.001");
        assertThat(tag.getDptId()).isEqualTo("DPT9.001");
    }

    @Test
    void matchesAgreesWithOf() {
        assertThat(KnxNetIpTag.matches("1/2/3")).isTrue();
        assertThat(KnxNetIpTag.matches("not-a-tag")).isFalse();
    }

    @Test
    void invalidAddressThrows() {
        assertThatThrownBy(() -> KnxNetIpTag.of("not-a-tag"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void wildcardMatchingAgainstGroupAddress() {
        // Wildcard accepts any same-level address.
        KnxNetIpTag wildcard = KnxNetIpTag.of("*/*/*");
        KnxNetIpTag exact = KnxNetIpTag.of("1/2/3");
        KnxNetIpTag different = KnxNetIpTag.of("9/9/9");

        GroupAddress addr = Mockito.mock(GroupAddress.class);
        Mockito.when(addr.getGroupAddress()).thenReturn("1/2/3");

        assertThat(wildcard.matchesGroupAddress(addr)).isTrue();
        assertThat(exact.matchesGroupAddress(addr)).isTrue();
        assertThat(different.matchesGroupAddress(addr)).isFalse();
    }

    @Test
    void tagHandlerRoutesAddressesAndBrowseQueries() {
        KnxNetIpTagHandler handler = new KnxNetIpTagHandler();
        assertThat(handler.parseTag("1/2/3")).isInstanceOf(KnxNetIpTag.class);
        assertThat(handler.parseQuery("any")).isInstanceOf(KnxNetIpQuery.class);
    }

    @Test
    void tagHandlerRejectsInvalidAddress() {
        assertThatThrownBy(() -> new KnxNetIpTagHandler().parseTag("nope"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void queryRoundtripsString() {
        KnxNetIpQuery q = new KnxNetIpQuery("hello");
        assertThat(q.getQueryString()).isEqualTo("hello");
    }

    @Test
    void configAccessorsRoundtrip() {
        KnxNetIpConfiguration cfg = new KnxNetIpConfiguration();
        cfg.setKnxprojFile(new File("/tmp/x.knxproj"));
        cfg.setKnxprojPassword("pw");
        cfg.setGroupAddressNumLevels(2);
        cfg.setConnectionType("BUSMONITOR");
        cfg.setRequestTimeout(7_000);

        assertThat(cfg.getKnxprojFile()).hasName("x.knxproj");
        assertThat(cfg.getKnxprojPassword()).isEqualTo("pw");
        assertThat(cfg.getGroupAddressNumLevels()).isEqualTo(2);
        assertThat(cfg.getConnectionType()).isEqualTo("BUSMONITOR");
        assertThat(cfg.getRequestTimeout()).isEqualTo(7_000);
        assertThat(cfg.toString()).isNotEmpty();
    }

    @Test
    void udpTransportDefaultPortIsPositive() {
        assertThat(new KnxNetIpUdpTransportConfiguration().getDefaultPort()).isPositive();
    }
}
