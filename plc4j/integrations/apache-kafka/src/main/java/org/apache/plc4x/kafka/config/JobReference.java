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
import java.util.Map;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class JobReference extends AbstractConfig{

    private static final Logger log = LoggerFactory.getLogger(JobReference.class);

    private final String name;
    private final String topic;

    public JobReference(String name, String defaultTopic, Map originals) {
        super(configDef(), originals);
        this.name = name;
        this.topic = getString(Constants.TOPIC_CONFIG) == null ? defaultTopic : getString(Constants.TOPIC_CONFIG);
    }

    public void validate() {
        return;
    }

    public String getName() {
        return name;
    }

    public String getTopic() {
        return topic;
    }

    protected static ConfigDef configDef() {
        return new ConfigDef()
            .define(Constants.TOPIC_CONFIG,
                    ConfigDef.Type.STRING,
                    Constants.TOPIC_DEFAULT,
                    ConfigDef.Importance.LOW,
                    Constants.TOPIC_DOC);
    }

    @Override
    public String toString() {
        StringBuilder query = new StringBuilder();
        query.append("\t\t" + name + "." + Constants.TOPIC_CONFIG + "=" + topic + ",\n");
        return query.toString();
    }
}
