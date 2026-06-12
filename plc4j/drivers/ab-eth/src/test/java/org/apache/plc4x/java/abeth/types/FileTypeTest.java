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
package org.apache.plc4x.java.abeth.types;

import org.apache.plc4x.java.api.types.PlcValueType;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class FileTypeTest {

    @Test
    void exposesAbTypeCodesAndPlcValueTypes() {
        assertThat(FileType.STATUS.getTypeCode()).isEqualTo((short) 0x84);
        assertThat(FileType.STATUS.getPlcValueType()).isEqualTo(PlcValueType.RAW_BYTE_ARRAY);

        assertThat(FileType.BIT.getTypeCode()).isEqualTo((short) 0x85);
        assertThat(FileType.BIT.getPlcValueType()).isEqualTo(PlcValueType.BOOL);

        assertThat(FileType.INTEGER.getTypeCode()).isEqualTo((short) 0x89);
        assertThat(FileType.INTEGER.getPlcValueType()).isEqualTo(PlcValueType.INT);

        assertThat(FileType.FLOAT.getTypeCode()).isEqualTo((short) 0x8A);
        assertThat(FileType.FLOAT.getPlcValueType()).isEqualTo(PlcValueType.REAL);

        assertThat(FileType.STRING.getTypeCode()).isEqualTo((short) 0x8D);
        assertThat(FileType.STRING.getPlcValueType()).isEqualTo(PlcValueType.STRING);
    }

    @Test
    void valueOfShortResolvesBackToFileType() {
        // WORD/DWORD/SINGLEBIT all share 0x89 with INTEGER — last-wins in the
        // static map, but the lookup must round-trip a known type code.
        assertThat(FileType.valueOf((short) 0x84)).isEqualTo(FileType.STATUS);
        assertThat(FileType.valueOf((short) 0x85)).isEqualTo(FileType.BIT);
        assertThat(FileType.valueOf((short) 0x8A)).isEqualTo(FileType.FLOAT);
        assertThat(FileType.valueOf((short) 0x8D)).isEqualTo(FileType.STRING);
        assertThat(FileType.valueOf((short) 0x99)).isNull();
    }

    @Test
    void valuesCoversAllEntries() {
        // Ensure no enum entry is missing from the catalogue.
        assertThat(FileType.values()).hasSize(15);
    }
}
