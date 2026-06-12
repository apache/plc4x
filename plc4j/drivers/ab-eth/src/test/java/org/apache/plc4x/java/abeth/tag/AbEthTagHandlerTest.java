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

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class AbEthTagHandlerTest {

    private final AbEthTagHandler handler = new AbEthTagHandler();

    @Test
    void parsesValidTag() {
        assertThat(handler.parseTag("N7:0:WORD")).isInstanceOf(AbEthTag.class);
    }

    @Test
    void rejectsInvalidTag() {
        assertThatThrownBy(() -> handler.parseTag("nope"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void doesNotSupportBrowsing() {
        assertThatThrownBy(() -> handler.parseQuery("anything"))
            .isInstanceOf(UnsupportedOperationException.class);
    }
}
