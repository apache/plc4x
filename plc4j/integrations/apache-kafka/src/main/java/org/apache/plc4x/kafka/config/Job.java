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

import java.util.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Job extends AbstractConfig{

    private static final Logger log = LoggerFactory.getLogger(Job.class);

    private final String name;
    private final int interval;
    private final List<Tag> tags;

    private static final String INTERVAL_CONFIG = "interval";
    private static final String INTERVAL_DOC = "Polling Interval";

    private static final String TAGS_CONFIG = "tags";
    private static final String TAGS_DOC = "List of tags assigned to Job";

    public Job(String name, Map originals) {
        super(configDef(), originals);

        this.name = name;
        this.interval = getInt(INTERVAL_CONFIG);

        tags = new ArrayList<>(getList(TAGS_CONFIG).size());
        for (String tagName : getList(TAGS_CONFIG)) {
            tags.add(new Tag(tagName, (String) originals.get(TAGS_CONFIG + "." + tagName)));
        }
    }

    public void validate() throws ConfigException {
        for (Tag tag : tags) {
            tag.validate();
        }
    }

    public String getName() {
        return name;
    }

    public int getInterval() {
        return interval;
    }

    public List<Tag> getTags() {
        return tags;
    }

    protected static ConfigDef configDef() {
        return new ConfigDef()
            .define(INTERVAL_CONFIG,
                    ConfigDef.Type.INT,
                    ConfigDef.Importance.LOW,
                    INTERVAL_DOC)
            .define(TAGS_CONFIG,
                    ConfigDef.Type.LIST,
                    ConfigDef.Importance.LOW,
                TAGS_DOC);
    }

    @Override
    public String toString() {
        StringBuilder query = new StringBuilder();
        query.append("\t\t" + name + "." + INTERVAL_CONFIG + "=" + interval + ",\n");
        for (Tag tag : tags) {
            query.append(tag.toString());
        }
        return query.toString();
    }

}
