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

import org.apache.plc4x.java.opcua.readwrite.StatusCode;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class OpcuaQualityStatusTest {

    @Test
    void classifiesGood() {
        OpcuaQualityStatus q = new OpcuaQualityStatus(new StatusCode(0x00000000L));
        assertThat(q.isGood()).isTrue();
        assertThat(q.isBad()).isFalse();
        assertThat(q.isUncertain()).isFalse();
        assertThat(q.toString()).isEqualTo("good");
    }

    @Test
    void classifiesBad() {
        OpcuaQualityStatus q = new OpcuaQualityStatus(new StatusCode(0x80000000L));
        assertThat(q.isGood()).isFalse();
        assertThat(q.isBad()).isTrue();
        assertThat(q.isUncertain()).isFalse();
        assertThat(q.toString()).isEqualTo("bad");
    }

    @Test
    void classifiesUncertain() {
        OpcuaQualityStatus q = new OpcuaQualityStatus(new StatusCode(0x40000000L));
        assertThat(q.isGood()).isFalse();
        assertThat(q.isBad()).isFalse();
        assertThat(q.isUncertain()).isTrue();
        assertThat(q.toString()).isEqualTo("uncertain");
    }
}
