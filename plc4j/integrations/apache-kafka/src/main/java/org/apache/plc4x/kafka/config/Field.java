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
package org.apache.plc4x.kafka.config;

import org.apache.kafka.common.config.AbstractConfig;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.common.config.Config;
import org.apache.kafka.common.config.ConfigException;
import java.util.Map;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

public class Field {

    private final String name;
    private final String address;
    private List<Config> configs;

    private static final String FIELD_DOC = "Field Address";

    public Field(String name, String address) {
        this.name = name;
        this.address = address;
    }

    public void validate() throws ConfigException {
        if (this.address == null) {
            throw new ConfigException(
                String.format("Field Address for field '%s' is missing", this.name));
        }
    }

    public String getName() {
        return name;
    }

    public String getAddress() {
        return address;
    }

    @Override
    public String toString() {
        StringBuilder query = new StringBuilder();
        query.append("\t\t\t" + name + "=" + address + ",\n");
        return query.toString();
    }

}
