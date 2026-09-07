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
package org.apache.plc4x.java.can.generic.tag;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.genericcan.readwrite.GenericCANDataType;
import org.junit.jupiter.api.Test;

import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class GenericCANTagTest {

    @Test
    void parsesScalarTag() {
        Optional<GenericCANTag> tag = GenericCANTag.matches("200:BYTE");
        assertThat(tag).isPresent();
        assertThat(tag.get().getNodeId()).isEqualTo(200);
        assertThat(tag.get().getDataType()).isEqualTo(GenericCANDataType.BYTE);
        assertThat(tag.get().getArraySize()).isZero();
        assertThat(tag.get().getPlcValueType()).isEqualTo(PlcValueType.BYTE);
    }

    @Test
    void parsesArrayTag() {
        Optional<GenericCANTag> tag = GenericCANTag.matches("200:BYTE[8]");
        assertThat(tag).isPresent();
        assertThat(tag.get().getArraySize()).isEqualTo(8);
    }

    @Test
    void unknownDataTypeThrows() {
        // The regex accepts arbitrary identifiers; the enum lookup fails late.
        assertThatThrownBy(() -> GenericCANTag.matches("200:NOT_A_TYPE"))
            .isInstanceOf(PlcRuntimeException.class);
    }

    @Test
    void nonMatchingTagReturnsEmpty() {
        assertThat(GenericCANTag.matches("not-a-can-tag")).isEmpty();
    }

    @Test
    void getAddressStringIncludesArraySuffixForArrays() {
        GenericCANTag arr = new GenericCANTag(7, GenericCANDataType.BYTE, 4);
        assertThat(arr.getAddressString()).isEqualTo("7:BYTE[4]");
    }

    @Test
    void getAddressStringDropsArraySuffixForArraySizeOne() {
        // The if (arraySize != 1) branch suppresses the [1] suffix.
        GenericCANTag scalar = new GenericCANTag(7, GenericCANDataType.BYTE, 1);
        assertThat(scalar.getAddressString()).isEqualTo("7:BYTE");
    }

    @Test
    void getArrayInfoNonEmptyForArrayTags() {
        GenericCANTag arr = new GenericCANTag(7, GenericCANDataType.BYTE, 4);
        assertThat(arr.getArrayInfo()).hasSize(1);
        // Scalar tag has no ArrayInfo entries.
        assertThat(new GenericCANTag(7, GenericCANDataType.BYTE, 1).getArrayInfo()).isEmpty();
    }

    @Test
    void toStringRendersAddress() {
        assertThat(new GenericCANTag(5, GenericCANDataType.BYTE, 3).toString())
            .contains("5").contains("BYTE").contains("3");
    }
}
