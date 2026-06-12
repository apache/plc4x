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
package org.apache.plc4x.java.ctrlx.readwrite;

import com.hrakaroo.glob.GlobPattern;
import com.hrakaroo.glob.MatchingEngine;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.ctrlx.readwrite.configuration.CtrlXConfiguration;
import org.apache.plc4x.java.ctrlx.readwrite.tag.CtrlXQuery;
import org.apache.plc4x.java.ctrlx.readwrite.tag.CtrlXTag;
import org.apache.plc4x.java.ctrlx.readwrite.tag.CtrlXTagHandler;
import org.junit.jupiter.api.Test;

import java.util.Collections;

import static org.assertj.core.api.Assertions.assertThat;

/**
 * The handwritten ctrlX surface is mostly stubs today; OpenAPI-generated REST
 * client code under {@code rest/datalayer} is excluded from coverage via the
 * jacoco plugin. These tests cover the small handwritten layer that remains.
 */
class CtrlXUnitsTest {

    @Test
    void tagAccessorsRoundtrip() {
        CtrlXTag tag = new CtrlXTag("/sys/some/path", PlcValueType.DINT, Collections.emptyList());
        assertThat(tag.getPath()).isEqualTo("/sys/some/path");
        assertThat(tag.getAddressString()).isEqualTo("/sys/some/path");
        assertThat(tag.getPlcValueType()).isEqualTo(PlcValueType.DINT);
        assertThat(tag.getArrayInfo()).isEmpty();
    }

    @Test
    void queryAccessorsRoundtrip() throws Exception {
        MatchingEngine matcher = GlobPattern.compile("*", '/', '\\', 0);
        CtrlXQuery query = new CtrlXQuery("*", matcher);
        assertThat(query.getQueryString()).isEqualTo("*");
        assertThat(query.getMatcher()).isSameAs(matcher);
    }

    @Test
    void tagHandlerParseTagIsStubReturningNull() {
        // parseTag is intentionally a stub today; pinning the behaviour so a
        // future replacement also removes this assertion.
        assertThat(new CtrlXTagHandler().parseTag("anything")).isNull();
    }

    @Test
    void tagHandlerParseQueryCompilesGlob() {
        CtrlXQuery q = (CtrlXQuery) new CtrlXTagHandler().parseQuery("/sys/*");
        assertThat(q).isNotNull();
        assertThat(q.getQueryString()).isEqualTo("/sys/*");
        assertThat(q.getMatcher()).isNotNull();
    }

    @Test
    void configurationInstantiates() {
        // Configuration is currently empty — just exercise the no-arg ctor.
        assertThat(new CtrlXConfiguration()).isNotNull();
    }
}
