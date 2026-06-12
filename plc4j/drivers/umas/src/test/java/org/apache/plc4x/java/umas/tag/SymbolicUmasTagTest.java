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
package org.apache.plc4x.java.umas.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.drivers.model.DefaultArrayInfo;
import org.junit.jupiter.api.Test;

import java.util.Collections;
import java.util.List;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class SymbolicUmasTagTest {

    @Test
    void parsesSimpleAddress() {
        SymbolicUmasTag tag = SymbolicUmasTag.of("MyVariable");
        assertThat(tag.getSymbolicAddress()).isEqualTo("MyVariable");
        assertThat(tag.getAddressString()).isEqualTo("MyVariable");
        assertThat(tag.getArrayInfo()).isEmpty();
    }

    @Test
    void parsesDottedNestedAddress() {
        SymbolicUmasTag tag = SymbolicUmasTag.of("Block.Field");
        assertThat(tag.getSymbolicAddress()).isEqualTo("Block.Field");
    }

    @Test
    void parsesArrayIndexAddress() {
        SymbolicUmasTag tag = SymbolicUmasTag.of("Arr[3]");
        assertThat(tag.getSymbolicAddress()).isEqualTo("Arr[3]");
    }

    @Test
    void matchesAgreesWithOf() {
        assertThat(SymbolicUmasTag.matches("MyVariable")).isTrue();
        assertThat(SymbolicUmasTag.matches("9bad-name")).isFalse();
    }

    @Test
    void invalidAddressThrows() {
        assertThatThrownBy(() -> SymbolicUmasTag.of("not a valid name"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void constructorDirectlyAcceptsTypeAndArrayInfo() {
        SymbolicUmasTag tag = new SymbolicUmasTag("foo", PlcValueType.DINT, Collections.emptyList());
        assertThat(tag.getPlcValueType()).isEqualTo(PlcValueType.DINT);
    }

    @Test
    void equalsAndHashCodeRespectAddressIdentity() {
        SymbolicUmasTag a = SymbolicUmasTag.of("Var");
        SymbolicUmasTag b = SymbolicUmasTag.of("Var");
        SymbolicUmasTag c = SymbolicUmasTag.of("Other");
        assertThat(a).isEqualTo(a).isEqualTo(b).isNotEqualTo(c).isNotEqualTo("string").isNotEqualTo(null);
        assertThat(a.hashCode()).isEqualTo(b.hashCode());
    }

    @Test
    void toStringRendersAddress() {
        assertThat(SymbolicUmasTag.of("Foo").toString()).contains("Foo");
    }

    @Test
    void getTotalNumberOfElementsReturnsOneForScalar() {
        // Default method on UmasTag — scalar (no array info) → 1.
        assertThat(SymbolicUmasTag.of("Var").getTotalNumberOfElements()).isEqualTo(1);
    }

    @Test
    void getTotalNumberOfElementsMultipliesArrayDimensions() {
        // DefaultArrayInfo.getSize() = upperBound - lowerBound + 1, so
        // (0..1)*(0..2) = 2 * 3 = 6.
        ArrayInfo dim2 = new DefaultArrayInfo(0, 1);
        ArrayInfo dim3 = new DefaultArrayInfo(0, 2);
        SymbolicUmasTag tag = new SymbolicUmasTag("Arr", PlcValueType.INT, List.of(dim2, dim3));
        assertThat(tag.getTotalNumberOfElements()).isEqualTo(6);
    }
}
