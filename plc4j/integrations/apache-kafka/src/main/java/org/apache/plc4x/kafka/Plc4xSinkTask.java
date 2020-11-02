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
package org.apache.plc4x.kafka;

import org.apache.kafka.common.config.AbstractConfig;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.connect.data.Date;
import org.apache.kafka.connect.data.*;
import org.apache.kafka.connect.errors.ConnectException;
import org.apache.kafka.connect.sink.SinkRecord;
import org.apache.kafka.connect.sink.SinkTask;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.scraper.config.triggeredscraper.JobConfigurationTriggeredImplBuilder;
import org.apache.plc4x.java.scraper.config.triggeredscraper.ScraperConfigurationTriggeredImpl;
import org.apache.plc4x.java.scraper.config.triggeredscraper.ScraperConfigurationTriggeredImplBuilder;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperImpl;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector.TriggerCollector;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector.TriggerCollectorImpl;
import org.apache.plc4x.java.utils.connectionpool.PooledPlcDriverManager;
import org.apache.plc4x.kafka.util.VersionUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigDecimal;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;
import java.util.concurrent.ArrayBlockingQueue;

/**
 * Source Connector Task polling the data source at a given rate.
 * A timer thread is scheduled which sets the fetch flag to true every rate milliseconds.
 * When poll() is invoked, the calling thread waits until the fetch flag is set for WAIT_LIMIT_MILLIS.
 * If the flag does not become true, the method returns null, otherwise a fetch is performed.
 */
public class Plc4xSinkTask extends SinkTask {

    private static final Logger log = LoggerFactory.getLogger(Plc4xSinkTask.class);

    /*
     * Config of the task.
     */
    static final String CONNECTION_NAME_CONFIG = "connection-name";
    private static final String CONNECTION_NAME_STRING_DOC = "Connection Name";

    static final String PLC4X_CONNECTION_STRING_CONFIG = "plc4x-connection-string";
    private static final String PLC4X_CONNECTION_STRING_DOC = "PLC4X Connection String";

    private static final ConfigDef CONFIG_DEF = new ConfigDef()
        .define(CONNECTION_NAME_CONFIG, ConfigDef.Type.STRING, ConfigDef.Importance.HIGH, CONNECTION_NAME_STRING_DOC)
        .define(PLC4X_CONNECTION_STRING_CONFIG, ConfigDef.Type.STRING, ConfigDef.Importance.HIGH, PLC4X_CONNECTION_STRING_DOC);

    /*
     * Configuration of the output.
     */
    private static final String SINK_NAME_FIELD = "sink-name";
    private static final String JOB_NAME_FIELD = "job-name";

    private static final Schema KEY_SCHEMA =
        new SchemaBuilder(Schema.Type.STRUCT)
            .field(SINK_NAME_FIELD, Schema.STRING_SCHEMA)
            .field(JOB_NAME_FIELD, Schema.STRING_SCHEMA)
            .build();

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

    private PlcDriverManager plcDriverManager;

    @Override
    public void start(Map<String, String> props) {
        AbstractConfig config = new AbstractConfig(CONFIG_DEF, props);
        String connectionName = config.getString(CONNECTION_NAME_CONFIG);
        String plc4xConnectionString = config.getString(PLC4X_CONNECTION_STRING_CONFIG);
        Map<String, String> topics = new HashMap<>();

        //try {
            log.info("Creating Pooled PLC4x driver manager");
            plcDriverManager = new PooledPlcDriverManager();
        //} catch (PlcConnectionException e) {
        //    log.error("Error starting the scraper", e.toString());
        //}
    }

    @Override
    public void stop() {
        synchronized (this) {
            notifyAll(); // wake up thread waiting in awaitFetch
        }
    }

    @Override
    public void put(Collection<SinkRecord> records) {
        if (records.isEmpty()) {
            return;
        }
        log.info(records.toString());
        return;
    }
}