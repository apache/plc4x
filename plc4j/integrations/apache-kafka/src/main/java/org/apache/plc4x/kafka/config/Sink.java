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
package org.apache.plc4x.kafka.config;

import java.util.Map;

public class Sink {

    private final String name;
    private final String connectionString;
    private final String topic;
    private final Integer retries;
    private final Integer timeout;
    private final Map<String, String> fields;

    public Sink(String name, String connectionString, String topic, Map<String, String> fields, Integer retries, Integer timeout) {
        this.name = name;
        this.connectionString = connectionString;
        this.topic = topic;
        this.fields = fields;
        this.retries = retries;
        this.timeout = timeout;
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

    public Map<String, String> getFields() {
        return fields;
    }

}
