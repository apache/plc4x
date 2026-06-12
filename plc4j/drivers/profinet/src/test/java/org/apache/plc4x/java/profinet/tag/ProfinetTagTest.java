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
import org.apache.plc4x.java.profinet.browse.ProfinetPlcQuery;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class ProfinetTagTest {

    @Test
    void parsesScalarTag() {
        ProfinetTag tag = ProfinetTag.of("foo.bar:DINT");
        assertThat(tag.getAddressString()).isEqualTo("foo.bar");
        assertThat(tag.getPlcValueType()).isEqualTo(PlcValueType.DINT);
        assertThat(tag.getArrayInfo()).isNotNull();
    }

    @Test
    void parsesQuantitySuffix() {
        ProfinetTag tag = ProfinetTag.of("foo:INT[4]");
        assertThat(tag.getPlcValueType()).isEqualTo(PlcValueType.INT);
    }

    @Test
    void rejectsZeroOrNegativeQuantityViaConstructor() {
        assertThatThrownBy(() -> ProfinetTagWithExposedCtor.of("foo:INT", 0))
            .isInstanceOf(IllegalArgumentException.class);
    }

    @Test
    void invalidAddressThrows() {
        assertThatThrownBy(() -> ProfinetTag.of("not-a-tag"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void tagHandlerRoutesParseQueryAndStubsParseTag() {
        ProfinetTagHandler handler = new ProfinetTagHandler();
        assertThat(handler.parseTag("anything")).isNull();
        assertThat(handler.parseQuery("anyquery")).isInstanceOf(ProfinetPlcQuery.class);
    }

    @Test
    void queryRoundtripsString() {
        assertThat(new ProfinetPlcQuery("q").getQueryString()).isEqualTo("q");
    }

    /**
     * Test-only subclass to access the protected constructor for the
     * zero-quantity validation path.
     */
    private static final class ProfinetTagWithExposedCtor extends ProfinetTag {
        private ProfinetTagWithExposedCtor(String address, Integer quantity, PlcValueType dataType) {
            super(address, quantity, dataType);
        }

        static ProfinetTag of(String address, int quantity) {
            return new ProfinetTagWithExposedCtor(address, quantity, PlcValueType.INT);
        }
    }
}
