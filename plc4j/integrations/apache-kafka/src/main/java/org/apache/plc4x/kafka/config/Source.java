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

import org.apache.plc4x.java.PlcDriverManager;

import org.apache.kafka.common.config.AbstractConfig;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.common.config.Config;
import org.apache.kafka.common.config.ConfigException;
import org.apache.kafka.common.config.ConfigValue;
import java.util.Map;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

public class Source extends AbstractConfig{

    private final String name;
    private final String connectionString;
    private final Integer bufferSize;
    private final Integer pollReturnInterval;
    private final List<JobReference> jobReferences;

    private static final String JOB_REFERENCES_CONFIG = "jobReferences";
    private static final String JOB_REFERENCES_DOC = "List of job references";

    public static final String SOURCES_CONFIG = "sources";

    public static final String CONNECTION_STRING_CONFIG = "connectionString";
    public static final String CONNECTION_STRING_DOC = "PLC4X Connection String";

    public static final String KAFKA_POLL_RETURN_CONFIG = "pollReturnInterval";
    public static final String KAFKA_POLL_RETURN_DOC = "Default poll return interval to be used, if not otherwise configured.";
    public static final Integer KAFKA_POLL_RETURN_DEFAULT = 5000;

    public static final String BUFFER_SIZE_CONFIG = "bufferSize";
    public static final String BUFFER_SIZE_DOC = "Default buffer size to be used, if not otherwise configured.";
    public static final Integer BUFFER_SIZE_DEFAULT = 1000;

    public Source(String name, String defaultTopic, Map originals) {
        super(configDef(), originals);
        this.name = name;
        this.connectionString = getString(CONNECTION_STRING_CONFIG);
        this.bufferSize = getInt(BUFFER_SIZE_CONFIG);
        this.pollReturnInterval = getInt(KAFKA_POLL_RETURN_CONFIG);

        jobReferences = new ArrayList<>(getList(JOB_REFERENCES_CONFIG).size());
        for (String jobReference : getList(JOB_REFERENCES_CONFIG)) {
            jobReferences.add(new JobReference(jobReference, defaultTopic, originalsWithPrefix(JOB_REFERENCES_CONFIG + "." + jobReference + ".")));
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
                String.format("Connection String format is incorrect %s ", SOURCES_CONFIG + "." + this.name + "." + CONNECTION_STRING_CONFIG + "=" + connectionString));
        }
        for (JobReference jobReference : jobReferences) {
            jobReference.validate();
        }
    }

    public String getName() {
        return name;
    }

    public String getConnectionString() {
        return connectionString;
    }

    public Integer getBufferSize() {
        return bufferSize;
    }

    public Integer getPollReturnInterval() {
        return pollReturnInterval;
    }

    public List<JobReference> getJobReferences() {
        return jobReferences;
    }

    protected static ConfigDef configDef() {
        return new ConfigDef()
            .define(CONNECTION_STRING_CONFIG,
                    ConfigDef.Type.STRING,
                    ConfigDef.Importance.HIGH,
                    CONNECTION_STRING_DOC)
            .define(BUFFER_SIZE_CONFIG,
                    ConfigDef.Type.INT,
                    BUFFER_SIZE_DEFAULT,
                    ConfigDef.Importance.LOW,
                    BUFFER_SIZE_DOC)
            .define(KAFKA_POLL_RETURN_CONFIG,
                    ConfigDef.Type.INT,
                    KAFKA_POLL_RETURN_DEFAULT,
                    ConfigDef.Importance.LOW,
                    KAFKA_POLL_RETURN_DOC)
            .define(JOB_REFERENCES_CONFIG,
                    ConfigDef.Type.LIST,
                    ConfigDef.Importance.LOW,
                    JOB_REFERENCES_DOC);
    }

    @Override
    public String toString() {
        StringBuilder query = new StringBuilder();
        query.append("\t" + "Name" + "=" + name + ",\n");
        query.append("\t" + CONNECTION_STRING_CONFIG + "=" + connectionString + ",\n");
        query.append("\t" + BUFFER_SIZE_CONFIG + "=" + bufferSize + ",\n");
        query.append("\t" + KAFKA_POLL_RETURN_CONFIG + "=" + pollReturnInterval + ",\n");

        for (JobReference jobReference : jobReferences) {
            query.append(jobReference.toString());
        }
        return query.toString();
    }

}
