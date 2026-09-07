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
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class UmasTagHandlerTest {

    private final UmasTagHandler handler = new UmasTagHandler();

    @Test
    void parsesValidSymbolicTag() {
        assertThat(handler.parseTag("MyVar")).isInstanceOf(SymbolicUmasTag.class);
    }

    @Test
    void rejectsInvalidTag() {
        assertThatThrownBy(() -> handler.parseTag("9invalid"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void browseQueryReturnsNullForPassThrough() {
        // The driver routes browse queries straight through the connection
        // without parsing — confirmed by returning null here.
        assertThat(handler.parseQuery("anything")).isNull();
    }
}
