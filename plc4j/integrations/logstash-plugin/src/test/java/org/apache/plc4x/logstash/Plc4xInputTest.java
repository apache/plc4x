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
package org.apache.plc4x.logstash;

import co.elastic.logstash.api.Configuration;
import org.assertj.core.util.Maps;
import org.junit.jupiter.api.Test;
import org.logstash.plugins.ConfigurationImpl;

import java.util.*;
import java.util.function.Consumer;

public class Plc4xInputTest {

    @Test
    public void testPlc4xInput() {
        Map<String, Object> configValues = new HashMap<>();
        Map<String, Object> jobValues = new HashMap<>();

        List<String> queries = Arrays.asList("testfield=RANDOM/foo:INTEGER");
        List<String> sources = Arrays.asList("TestConnection");

        jobValues.put("rate", 300);
        jobValues.put("queries", queries);
        jobValues.put("sources", sources);

        configValues.put(Plc4x.SOURCE_CONFIG.name(), Maps.newHashMap("TestConnection", "test:hurzpurzfurz"));
        configValues.put(Plc4x.JOB_CONFIG.name(),  Maps.newHashMap("job1", jobValues));


        Configuration config = new ConfigurationImpl(configValues);
        Plc4x input = new Plc4x("test-id", config, null);
        TestConsumer testConsumer = new TestConsumer();
        input.start(testConsumer);

        List<Map<String, Object>> events = testConsumer.getEvents();
        System.out.println(events.size());
    }

    private static class TestConsumer implements Consumer<Map<String, Object>> {

        private List<Map<String, Object>> events = new ArrayList<>();

        @Override
        public void accept(Map<String, Object> event) {
            synchronized (this) {
                events.add(event);
            }
        }

        public List<Map<String, Object>> getEvents() {
            return events;
        }
    }

}
