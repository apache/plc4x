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
package org.apache.plc4x;

import org.apache.calcite.schema.Table;
import org.apache.calcite.schema.impl.AbstractSchema;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.scraper.ResultHandler;
import org.apache.plc4x.java.scraper.Scraper;
import org.apache.plc4x.java.scraper.ScraperImpl;
import org.apache.plc4x.java.scraper.config.JobConfiguration;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;
import org.apache.plc4x.java.scraper.exception.ScraperException;

import java.time.Instant;
import java.util.Map;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.BlockingQueue;
import java.util.stream.Collectors;

public class Plc4xSchema extends AbstractSchema {

    protected final ScraperConfiguration configuration;
    protected final Scraper scraper;
    protected final QueueHandler handler;
    protected final Map<String, BlockingQueue<Record>> queues;
    protected final Map<String, Table> tableMap;

    public Plc4xSchema(ScraperConfiguration configuration, long tableCutoff) throws ScraperException {
        this.configuration = configuration;
        this.handler = new QueueHandler();
        this.scraper = new ScraperImpl(configuration, handler);
        this.queues = configuration.getJobConfigurations().stream()
            .collect(Collectors.toMap(
                JobConfiguration::getName,
                conf -> new ArrayBlockingQueue<Record>(1000)
            ));
        // Create the tables
        this.tableMap = configuration.getJobConfigurations().stream()
            .collect(Collectors.toMap(
                JobConfiguration::getName,
                conf -> defineTable(queues.get(conf.getName()), conf, tableCutoff)
            ));
        // Start the scraper
        this.scraper.start();
    }

    Table defineTable(BlockingQueue<Record> queue, JobConfiguration configuration, Long limit) {
        if (limit <= 0) {
            return new Plc4xStreamTable(queue, configuration);
        } else {
            return new Plc4xTable(queue, configuration, limit);
        }
    }

    @Override
    protected Map<String, Table> getTableMap() {
        // Return a map of all jobs
        return this.tableMap;
    }

    public static class Record {

        public final Instant timestamp;
        public final String source;
        public final Map<String, Object> values;

        public Record(Instant timestamp, String source, Map<String, Object> values) {
            this.timestamp = timestamp;
            this.source = source;
            this.values = values;
        }
    }

    class QueueHandler implements ResultHandler {

        @Override
        public void handle(String job, String alias, Map<String, Object> results) {
            try {
                Record record = new Record(Instant.now(), alias, results);
                queues.get(job).put(record);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                throw new PlcRuntimeException("Handling got interrupted", e);
            }
        }

    }
}
