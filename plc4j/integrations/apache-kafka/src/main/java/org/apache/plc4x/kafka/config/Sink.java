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

import org.apache.plc4x.java.PlcDriverManager;

import org.apache.kafka.common.config.AbstractConfig;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.common.config.ConfigException;

import java.util.Map;
import java.util.ArrayList;
import java.util.List;

public class Sink extends AbstractConfig{

    private final String name;
    private final String connectionString;
    private final String topic;
    private final Integer retries;
    private final Integer timeout;
    private final List<Tag> tags;

    public Sink(String name, Map originals) {
        super(configDef(), originals);
        this.name = name;
        this.connectionString = getString(Constants.CONNECTION_STRING_CONFIG);
        this.topic = getString(Constants.TOPIC_CONFIG);
        this.retries = getInt(Constants.RETRIES_CONFIG);
        this.timeout = getInt(Constants.TIMEOUT_CONFIG);

        tags = new ArrayList<>(getList(Constants.TAGS_CONFIG).size());
        for (String tagName : getList(Constants.TAGS_CONFIG)) {
            tags.add(new Tag(tagName, (String) originals.get(Constants.TAGS_CONFIG + "." + tagName)));
        }
    }

    public void validate() throws ConfigException {
        if (this.connectionString == null) {
            throw new ConfigException(
                String.format("Connection string shouldn't be null for source %s ", this.name));
        }
        try {
            new PlcDriverManager().getDriverForUrl(connectionString);
        } catch (Exception e) {
            throw new ConfigException(
                String.format("Connection String format is incorrect %s ", Constants.SINKS_CONFIG + "." + this.name + "." + Constants.CONNECTION_STRING_CONFIG + "=" + connectionString));
        }
        for (Tag tag : tags) {
            tag.validate();
        }
    }

    public String getName() {
        return name;
    }

    public String getConnectionString() {
        return connectionString;
    }

    public String getTopic() {
        return topic;
    }

    public Integer getRetries() {
        return retries;
    }

    public Integer getTimeout() {
        return timeout;
    }

    public List<Tag> getTags() {
        return tags;
    }

    public static ConfigDef configDef() {
        return new ConfigDef()
            .define(Constants.CONNECTION_STRING_CONFIG,
                    ConfigDef.Type.STRING,
                    ConfigDef.Importance.HIGH,
                    Constants.CONNECTION_STRING_DOC)
            .define(Constants.TOPIC_CONFIG,
                    ConfigDef.Type.STRING,
                    ConfigDef.Importance.LOW,
                    Constants.TOPIC_DOC)
            .define(Constants.RETRIES_CONFIG,
                    ConfigDef.Type.INT,
                    Constants.RETRIES_DEFAULT,
                    ConfigDef.Importance.LOW,
                    Constants.RETRIES_CONFIG)
            .define(Constants.TIMEOUT_CONFIG,
                    ConfigDef.Type.INT,
                    Constants.TIMEOUT_DEFAULT,
                    ConfigDef.Importance.LOW,
                    Constants.JOBS_DOC)
            .define(Constants.TAGS_CONFIG,
                    ConfigDef.Type.LIST,
                    Constants.TAGS_DEFAULT,
                    ConfigDef.Importance.LOW,
                    Constants.TAGS_CONFIG);
    }

    @Override
    public String toString() {
        StringBuilder query = new StringBuilder();
        query.append(Constants.CONNECTION_STRING_CONFIG + "=" + connectionString + ",\n");
        query.append(Constants.TOPIC_CONFIG + "=" + topic + ",\n");
        query.append(Constants.RETRIES_CONFIG + "=" + retries + ",\n");
        query.append(Constants.TIMEOUT_CONFIG + "=" + timeout + ",\n");
        for (Tag tag : tags) {
            query.append(tag.toString());
        }
        return query.toString();
    }

}
