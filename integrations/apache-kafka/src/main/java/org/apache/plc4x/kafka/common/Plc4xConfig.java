/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.kafka.common;

import org.apache.kafka.common.config.AbstractConfig;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.common.config.ConfigDef.Importance;
import org.apache.kafka.common.config.ConfigDef.Type;
import org.apache.kafka.common.config.ConfigException;

import java.util.Map;


public class Plc4xConfig extends AbstractConfig {

    public static final String PLC_CONNECTION_STRING_CONFIG = "my.setting";
    public static final String PLC_CONNECTION_STRING_DISPLAY = "PLC Connection String";
    public static final String PLC_CONNECTION_STRING_DOC = "Connection string used by PLC4X to connect to the PLC.";

    public static final String PLC_TOPIC = "topic";
    public static final String PLC_TOPIC_DOC = "Kafka topic to publish messages to.";

    public static final String PLC_DATATYPE_CONFIG = "type";
    public static final String PLC_DATATYPE_DOC = "Data type of values sent or received by PLC.";

    public static final String PLC_ADDRESS = "address";
    public static final String PLC_ADDRESS_DOC = "PLC address to sent to or receive data from.";

    public static ConfigDef baseConfigDef() {
        ConfigDef config = new ConfigDef();
        addPlcOptions(config);
        return config;
    }

    private static final void addPlcOptions(ConfigDef config) {
        config.define(
            PLC_CONNECTION_STRING_CONFIG,
            Type.STRING,
            Importance.HIGH,
            PLC_CONNECTION_STRING_DOC)
        .define(
            PLC_DATATYPE_CONFIG,
            Type.CLASS,
            Importance.HIGH,
            PLC_DATATYPE_DOC)
        .define(
            PLC_TOPIC,
            Type.STRING,
            Importance.HIGH,
            PLC_TOPIC_DOC)
        .define(
            PLC_ADDRESS,
            Type.STRING,
            Importance.HIGH,
            PLC_ADDRESS_DOC);
    }

    public static final ConfigDef CONFIG_DEF = baseConfigDef();

    public Plc4xConfig(ConfigDef config, Map<String, String> parsedConfig) {
        super(config, parsedConfig);
        String plcConnectionString = getString(PLC_CONNECTION_STRING_CONFIG);
        if (plcConnectionString == null) {
            throw new ConfigException("'PLC Connection String' must be specified");
        }
    }

    public Plc4xConfig(Map<String, String> parsedConfig) {
        this(CONFIG_DEF, parsedConfig);
    }

    public String getPlcConnectionString() {
        return this.getString(PLC_CONNECTION_STRING_CONFIG);
    }

}
