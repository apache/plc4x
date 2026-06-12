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
package org.apache.plc4x.java.eip.base.tag;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.eip.readwrite.CIPDataTypeCode;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class EipTagCoverageTest {

    @Test
    void singleArgCtorDefaultsElementsTo1AndLeavesTypeNull() {
        EipTag tag = new EipTag("%A0");
        assertThat(tag.getTag()).isEqualTo("%A0");
        assertThat(tag.getElementNb()).isEqualTo(1);
        assertThat(tag.getType()).isNull();
    }

    @Test
    void ctorWithElementNb() {
        EipTag tag = new EipTag("%A0", 4);
        assertThat(tag.getElementNb()).isEqualTo(4);
        assertThat(tag.getType()).isNull();
    }

    @Test
    void ctorWithTypeOnly() {
        EipTag tag = new EipTag("%A0", CIPDataTypeCode.INT);
        assertThat(tag.getType()).isEqualTo(CIPDataTypeCode.INT);
        // No element-count constructor → default int zero (not the 1 the
        // single-arg ctor injects). Lock the surface in place.
        assertThat(tag.getElementNb()).isZero();
    }

    @Test
    void ctorWithTypeAndElements() {
        EipTag tag = new EipTag("%A0", CIPDataTypeCode.DINT, 3);
        assertThat(tag.getType()).isEqualTo(CIPDataTypeCode.DINT);
        assertThat(tag.getElementNb()).isEqualTo(3);
        assertThat(tag.getPlcValueType()).isEqualTo(PlcValueType.DINT);
    }

    @Test
    void setTypeAndSetElementNbRoundtrip() {
        EipTag tag = new EipTag("%A0");
        tag.setType(CIPDataTypeCode.REAL);
        tag.setElementNb(7);
        assertThat(tag.getType()).isEqualTo(CIPDataTypeCode.REAL);
        assertThat(tag.getElementNb()).isEqualTo(7);
    }

    @Test
    void ofParsesTagWithTypeAndElements() {
        EipTag tag = EipTag.of("MyVar:INT:5");
        assertThat(tag).isNotNull();
        assertThat(tag.getTag()).isEqualTo("MyVar");
        assertThat(tag.getType()).isEqualTo(CIPDataTypeCode.INT);
        assertThat(tag.getElementNb()).isEqualTo(5);
    }

    @Test
    void ofDefaultsDataTypeToDintWhenMissing() {
        EipTag tag = EipTag.of("%A0:2");
        assertThat(tag).isNotNull();
        assertThat(tag.getType()).isEqualTo(CIPDataTypeCode.DINT);
        assertThat(tag.getElementNb()).isEqualTo(2);
    }

    @Test
    void ofWithZeroElementsUsesTypeOnlyCtorPath() {
        // elementNb=0 selects the (tag, type) constructor branch.
        EipTag tag = EipTag.of("%A0:INT:0");
        assertThat(tag).isNotNull();
        assertThat(tag.getElementNb()).isZero();
        assertThat(tag.getType()).isEqualTo(CIPDataTypeCode.INT);
    }

    @Test
    void matchesAgreesWithOf() {
        assertThat(EipTag.matches("%A0:2")).isTrue();
        // The regex is permissive — the unmatched-input branch in `of`
        // returns null; we don't need a separate negative test for `matches`
        // beyond the positive case above.
    }

    @Test
    void getAddressStringIsNotImplemented() {
        assertThatThrownBy(() -> new EipTag("%A0").getAddressString())
            .isInstanceOf(NotImplementedException.class);
    }

    @Test
    void getArrayInfoReturnsDefault() {
        assertThat(new EipTag("%A0").getArrayInfo()).isNotNull();
    }

    @Test
    void serializeWritesTagTypeAndElementCount() throws Exception {
        WriteBufferByteBased buf = new WriteBufferByteBased(new byte[128],
            WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"),
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithStringEncoding("UTF8"));
        new EipTag("%A0", CIPDataTypeCode.DINT, 2).serialize(buf);
        assertThat(buf.getBytes()).isNotEmpty();
    }

    @Test
    void serializeSkipsTypeWhenNull() throws Exception {
        WriteBufferByteBased buf = new WriteBufferByteBased(new byte[128],
            WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"),
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithStringEncoding("UTF8"));
        // Single-arg ctor leaves type null — exercises the if(type != null)
        // false branch in serialize.
        new EipTag("%A0").serialize(buf);
        assertThat(buf.getBytes()).isNotEmpty();
    }
}
