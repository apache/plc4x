/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x;

import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;

import java.util.Collections;

class Plc4xSchemaFactoryTest implements WithAssertions {

    @Test
    void create() {
        Plc4xSchemaFactory factory = new Plc4xSchemaFactory();
        assertThatThrownBy(() -> factory.create(null, "", Collections.emptyMap()))
            .isInstanceOf(NullPointerException.class)
            .hasMessageContaining("Please specify operand 'config'...");
    }
}