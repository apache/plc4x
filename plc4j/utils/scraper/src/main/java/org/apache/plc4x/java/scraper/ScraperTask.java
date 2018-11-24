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

package org.apache.plc4x.java.scraper;

import org.apache.commons.lang3.Validate;
import org.apache.commons.lang3.time.StopWatch;
import org.apache.commons.math3.stat.descriptive.DescriptiveStatistics;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.Function;
import java.util.stream.Collectors;

/**
 * Plc Scraper that scrapes one source.
 */
public class ScraperTask implements Runnable {

    private static final Logger LOGGER = LoggerFactory.getLogger(ScraperTask.class);

    private final PlcDriverManager driverManager;
    private final String jobName;
    private final String connectionAlias;
    private final String connectionString;
    private final Map<String, String> fields;
    private final long requestTimeoutMs;
    private final ExecutorService handlerService;

    private final AtomicLong requestCounter = new AtomicLong(0);
    private final AtomicLong successCounter = new AtomicLong(0);
    private final DescriptiveStatistics latencyStatistics = new DescriptiveStatistics(1000);
    private final DescriptiveStatistics failedStatistics = new DescriptiveStatistics(1000);

    public ScraperTask(PlcDriverManager driverManager, String jobName, String connectionAlias, String connectionString,
                       Map<String, String> fields, long requestTimeoutMs, ExecutorService handlerService) {
        Validate.notNull(driverManager);
        Validate.notBlank(jobName);
        Validate.notBlank(connectionAlias);
        Validate.notBlank(connectionString);
        Validate.notEmpty(fields);
        Validate.isTrue(requestTimeoutMs > 0);
        this.driverManager = driverManager;
        this.jobName = jobName;
        this.connectionAlias = connectionAlias;
        this.connectionString = connectionString;
        this.fields = fields;
        this.requestTimeoutMs = requestTimeoutMs;
        this.handlerService = handlerService;
    }

    @Override
    public void run() {
        // Does a single fetch
        LOGGER.trace("Start new scrape of task of job {} for connection {}", jobName, connectionAlias);
        requestCounter.incrementAndGet();
        StopWatch stopWatch = new StopWatch();
        stopWatch.start();
        try (PlcConnection connection = driverManager.getConnection(connectionString)) {
            LOGGER.trace("Connection to {} established: {}", connectionString, connection);
            PlcReadResponse response;
            try {
                PlcReadRequest.Builder builder = connection.readRequestBuilder();
                fields.forEach((alias,qry) -> {
                    LOGGER.trace("Requesting: {} -> {}", alias, qry);
                    builder.addItem(alias,qry);
                });
                response = builder
                    .build()
                    .execute()
                    .get(requestTimeoutMs, TimeUnit.MILLISECONDS);
            } catch (ExecutionException e) {
                // Handle execution exception
                handleException(e);
                return;
            }
            // Add statistics
            stopWatch.stop();
            latencyStatistics.addValue(stopWatch.getNanoTime());
            failedStatistics.addValue(0.0);
            successCounter.incrementAndGet();
            // Validate response
            validateResponse(response);
            // Handle response (Async)
            CompletableFuture.runAsync(() -> handle(transformResponseToMap(response)), handlerService);
        } catch (Exception e) {
            failedStatistics.addValue(1.0);
            LOGGER.debug("Exception during scrape", e);
            handleException(e);
        }
    }

    private void validateResponse(PlcReadResponse response) {
        Map<String, PlcResponseCode> failedFields = response.getFieldNames().stream()
            .filter(name -> !PlcResponseCode.OK.equals(response.getResponseCode(name)))
            .collect(Collectors.toMap(
                Function.identity(),
                response::getResponseCode
            ));
        if (failedFields.size() > 0) {
            handleErrorResponse(failedFields);
        }
    }

    private Map<String, Object> transformResponseToMap(PlcReadResponse response) {
        return response.getFieldNames().stream()
            .collect(Collectors.toMap(
                name -> name,
                response::getObject
            ));
    }

    public String getJobName() {
        return jobName;
    }

    public String getConnectionAlias() {
        return connectionAlias;
    }

    public long getRequestCounter() {
        return requestCounter.get();
    }

    public long getSuccessfullRequestCounter() {
        return successCounter.get();
    }

    public DescriptiveStatistics getLatencyStatistics() {
        return latencyStatistics;
    }

    public double getPercentageFailed() {
        return 100.0*failedStatistics.getMean();
    }

    public void handle(Map<String, Object> result) {
        LOGGER.debug("Handling result on gorgeous pool: {}", result);
    }

    public void handleException(Exception e) {
        failedStatistics.addValue(1.0);
    }

    public void handleErrorResponse(Map<String, PlcResponseCode> failed) {
        LOGGER.warn("Handling error responses: {}", failed);
    }

}
