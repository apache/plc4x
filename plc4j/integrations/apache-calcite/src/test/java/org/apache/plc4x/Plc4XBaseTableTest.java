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
package org.apache.plc4x;

import org.apache.calcite.linq4j.Enumerator;
import org.apache.plc4x.java.scraper.config.JobConfigurationImpl;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;

import java.time.Instant;
import java.util.Collections;
import java.util.Map;
import java.util.concurrent.ArrayBlockingQueue;

class Plc4XBaseTableTest implements WithAssertions {

    @Test
    void testOnBlockingQueue() {
        ArrayBlockingQueue<Plc4xSchema.Record> queue = new ArrayBlockingQueue<>(100);
        Plc4xStreamTable table = new Plc4xStreamTable(queue, new JobConfigurationImpl(
            "job1",
            null,
            100,
            Collections.emptyList(),
            Collections.singletonMap("key", "address")));

        Map<String, Object> objects = Collections.singletonMap("key", "value");
        queue.add(new Plc4xSchema.Record(Instant.now(), "", objects));

        Enumerator<Object[]> enumerator = table.scan(null).enumerator();

        assertThat(enumerator.moveNext()).isTrue();
        assertThat(enumerator.current()).contains("value");
    }

}