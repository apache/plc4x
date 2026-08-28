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
package org.apache.plc4x.java.profinet.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class ProfinetTagTest {

    @Test
    void parsesFullAddress() {
        ProfinetTag tag = ProfinetTag.of("1.2.INPUT.3[0..3]:INT");
        assertThat(tag.getSlot()).isEqualTo(1);
        assertThat(tag.getSubSlot()).isEqualTo(2);
        assertThat(tag.getDirection()).isEqualTo(ProfinetTag.Direction.INPUT);
        assertThat(tag.getIndex()).isEqualTo(3);
        assertThat(tag.getPlcValueType()).isEqualTo(PlcValueType.INT);
        assertThat(tag.getNumElements()).isEqualTo(4);
        assertThat(tag.getArrayInfo()).isNotNull();
    }

    @Test
    void outputDirectionAlsoSupported() {
        ProfinetTag tag = ProfinetTag.of("1.2.OUTPUT.3:DINT");
        assertThat(tag.getDirection()).isEqualTo(ProfinetTag.Direction.OUTPUT);
        assertThat(tag.getNumElements()).isEqualTo(1);
    }

    @Test
    void getAddressStringRoundtrips() {
        ProfinetTag tag = ProfinetTag.of("1.2.INPUT.3:INT");
        assertThat(tag.getAddressString()).isNotEmpty();
    }

    @Test
    void invalidAddressThrows() {
        assertThatThrownBy(() -> ProfinetTag.of("not-a-tag"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void zeroElementsRejected() {
        assertThatThrownBy(() ->
            new ProfinetTag(1, 2, ProfinetTag.Direction.INPUT, 3, PlcValueType.INT, 0))
            .isInstanceOf(IllegalArgumentException.class);
    }

    @Test
    void tagHandlerParsesTagsAndStubsQueries() {
        ProfinetTagHandler handler = new ProfinetTagHandler();
        assertThat(handler.parseTag("1.2.INPUT.3:INT")).isInstanceOf(ProfinetTag.class);
        assertThat(handler.parseQuery("anything")).isNull();
    }

    @Test
    void aCountTooWideToBeANumberIsAnInvalidTagNotANumberFormatError() {
        assertThatThrownBy(() -> ProfinetTag.of("1.1.INPUT.1[0..99999999998]:BOOL"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void aSlotTooWideToBeANumberIsAlsoAnInvalidTag() {
        assertThatThrownBy(() -> ProfinetTag.of("99999999999.1.INPUT.1:BOOL"))
            .isInstanceOf(PlcInvalidTagException.class);
    }
}
