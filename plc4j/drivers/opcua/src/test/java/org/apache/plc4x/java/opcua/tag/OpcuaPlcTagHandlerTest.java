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
package org.apache.plc4x.java.opcua.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class OpcuaPlcTagHandlerTest {

    private final OpcuaPlcTagHandler handler = new OpcuaPlcTagHandler();

    @Test
    void parsesValidTag() {
        OpcuaTag tag = handler.parseTag("ns=2;i=10846");
        assertThat(tag).isNotNull();
        assertThat(tag.getNamespace()).isEqualTo(2);
    }

    @Test
    void rejectsInvalidTag() {
        assertThatThrownBy(() -> handler.parseTag("nonsense"))
            .isInstanceOf(PlcInvalidTagException.class);
    }

    @Test
    void parsesBrowseQuery() {
        // Browsing is supported: parseQuery turns the browse expression into an OpcuaQuery
        // carrying the start-address (or wildcard) string verbatim.
        PlcQuery query = handler.parseQuery("ns=2;s=HelloWorld");
        assertThat(query).isInstanceOf(OpcuaQuery.class);
        assertThat(query.getQueryString()).isEqualTo("ns=2;s=HelloWorld");
    }

}
