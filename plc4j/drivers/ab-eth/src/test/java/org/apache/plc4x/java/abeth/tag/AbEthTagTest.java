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
package org.apache.plc4x.java.abeth.tag;

import org.apache.plc4x.java.abeth.types.FileType;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class AbEthTagTest {

    @Test
    void parsesWordTagDeriving2ByteSize() {
        AbEthTag tag = AbEthTag.of("N7:0:WORD");
        assertThat(tag.getFileNumber()).isEqualTo((short) 7);
        assertThat(tag.getElementNumber()).isEqualTo((short) 0);
        assertThat(tag.getBitNumber()).isEqualTo((short) 0);
        assertThat(tag.getFileType()).isEqualTo(FileType.WORD);
        // WORD/SINGLEBIT use a hard-coded 2-byte size.
        assertThat(tag.getByteSize()).isEqualTo((short) 2);
        assertThat(tag.getPlcValueType()).isEqualTo(PlcValueType.WORD);
    }

    @Test
    void parsesDwordTagWith4ByteSize() {
        AbEthTag tag = AbEthTag.of("N7:5:DWORD");
        assertThat(tag.getFileType()).isEqualTo(FileType.DWORD);
        assertThat(tag.getByteSize()).isEqualTo((short) 4);
    }

    @Test
    void parsesIntegerTagWithExplicitSize() {
        AbEthTag tag = AbEthTag.of("N7:3:INTEGER[8]");
        assertThat(tag.getFileType()).isEqualTo(FileType.INTEGER);
        assertThat(tag.getByteSize()).isEqualTo((short) 8);
    }

    @Test
    void parsesTagWithBitNumber() {
        AbEthTag tag = AbEthTag.of("N7:0/3:WORD");
        assertThat(tag.getBitNumber()).isEqualTo((short) 3);
    }

    @Test
    void matchesHelperAgreesWithOf() {
        assertThat(AbEthTag.matches("N7:0:WORD")).isTrue();
        assertThat(AbEthTag.matches("not-an-abeth-tag")).isFalse();
    }

    @Test
    void invalidTagThrows() {
        assertThatThrownBy(() -> AbEthTag.of("not-a-tag"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void getAddressStringRoundtripsForWord() {
        AbEthTag tag = AbEthTag.of("N7:0:WORD");
        // The serialized form contains the byteSize suffix when != 1.
        assertThat(tag.getAddressString()).isEqualTo("N7:0:WORD[2]");
    }

    @Test
    void getAddressStringIncludesBitNumberWhenSet() {
        AbEthTag tag = AbEthTag.of("N7:0/3:WORD");
        assertThat(tag.getAddressString()).isEqualTo("N7:0/3:WORD[2]");
    }

    @Test
    void getArrayInfoReturnsDefault() {
        assertThat(AbEthTag.of("N7:0:WORD").getArrayInfo()).isNotNull();
    }

    @Test
    void constructorDirectlySetsAllFields() {
        AbEthTag tag = new AbEthTag((short) 4, (short) 9, FileType.DWORD, (short) 2, (short) 1);
        assertThat(tag.getByteSize()).isEqualTo((short) 4);
        assertThat(tag.getFileNumber()).isEqualTo((short) 9);
        assertThat(tag.getFileType()).isEqualTo(FileType.DWORD);
        assertThat(tag.getElementNumber()).isEqualTo((short) 2);
        assertThat(tag.getBitNumber()).isEqualTo((short) 1);
    }
}
