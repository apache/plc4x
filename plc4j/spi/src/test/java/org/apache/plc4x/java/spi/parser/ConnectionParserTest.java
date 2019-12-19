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

package org.apache.plc4x.java.spi.parser;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

class ConnectionParserTest {

    @Test
    void parse() {
        ConnectionParser parster = new ConnectionParser();
        PropertiesDescriptor properties = parster.parse("s7://192.168.167.1?rackId=1", PropertiesDescriptor.class);

        assertEquals(1, properties.rackId);
        assertEquals(1, properties.slotId);
    }

    static class PropertiesDescriptor {

        @ConfigurationParameter("rackId")
        @IntDefaultValue(1)
        private int rackId;

        @ConfigurationParameter("slotId")
        @IntDefaultValue(1)
        @Required
        private int slotId;

    }
}