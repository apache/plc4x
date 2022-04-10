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

import static org.junit.jupiter.api.Assertions.assertEquals;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.Required;
import org.junit.jupiter.api.Test;

class ConnectionParserTest {

    @Test
    void parse() throws PlcConnectionException {
        PropertiesDescriptor properties = new ConfigurationFactory().createConfiguration(
            PropertiesDescriptor.class, "rackId=2");

        assertEquals(2, properties.getRackId());
        assertEquals(1, properties.getSlotId());
    }

    public static class PropertiesDescriptor implements Configuration {

        @ConfigurationParameter("rackId")
        @IntDefaultValue(1)
        private int rackId;

        @ConfigurationParameter("slotId")
        @IntDefaultValue(1)
        @Required
        private int slotId;

        public int getRackId() {
            return rackId;
        }

        public void setRackId(int rackId) {
            this.rackId = rackId;
        }

        public int getSlotId() {
            return slotId;
        }

        public void setSlotId(int slotId) {
            this.slotId = slotId;
        }

    }

}