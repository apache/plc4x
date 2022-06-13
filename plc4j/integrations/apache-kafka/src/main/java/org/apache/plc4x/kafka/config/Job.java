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
package org.apache.plc4x.kafka.config;

import org.apache.kafka.common.config.AbstractConfig;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.common.config.ConfigException;

import java.time.Duration;
import java.util.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Job extends AbstractConfig{

    private static final Logger log = LoggerFactory.getLogger(Job.class);

    private final String name;
    private final int interval;
    private final List<Field> fields;

    private static final String INTERVAL_CONFIG = "interval";
    private static final String INTERVAL_DOC = "Polling Interval";

    private static final String FIELDS_CONFIG = "fields";
    private static final String FIELDS_DOC = "List of fields assigned to Job";

    public Job(String name, Map originals) {
        super(configDef(), originals);

        this.name = name;
        this.interval = getInt(INTERVAL_CONFIG);

        fields = new ArrayList<>(getList(FIELDS_CONFIG).size());
        for (String field : getList(FIELDS_CONFIG)) {
            fields.add(new Field(field, (String) originals.get(FIELDS_CONFIG + "." + field)));
        }
    }

    public void validate() throws ConfigException {
        for (Field field : fields) {
            field.validate();
        }
    }

    public String getName() {
        return name;
    }

    public int getInterval() {
        return interval;
    }

    public List<Field> getFields() {
        return fields;
    }

    protected static ConfigDef configDef() {
        return new ConfigDef()
            .define(INTERVAL_CONFIG,
                    ConfigDef.Type.INT,
                    ConfigDef.Importance.LOW,
                    INTERVAL_DOC)
            .define(FIELDS_CONFIG,
                    ConfigDef.Type.LIST,
                    ConfigDef.Importance.LOW,
                    FIELDS_DOC);
    }

    @Override
    public String toString() {
        StringBuilder query = new StringBuilder();
        query.append("\t\t" + name + "." + INTERVAL_CONFIG + "=" + interval + ",\n");
        for (Field field : fields) {
            query.append(field.toString());
        }
        return query.toString();
    }

}
