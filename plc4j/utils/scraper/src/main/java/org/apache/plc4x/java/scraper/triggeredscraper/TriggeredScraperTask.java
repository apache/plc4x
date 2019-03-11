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

package org.apache.plc4x.java.scraper.triggeredscraper;

import org.apache.commons.lang3.time.StopWatch;
import org.apache.commons.math3.stat.descriptive.DescriptiveStatistics;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.scraper.ResultHandler;
import org.apache.plc4x.java.scraper.ScraperTask;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.TriggerHandler;
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
 * performs the triggered task from a job for one device based on the TriggerHandler as defined in Configuration
 * ToDo Implement the monitoring as well: PLC4X-90
 */
public class TriggeredScraperTask implements ScraperTask {
    private static final Logger LOGGER = LoggerFactory.getLogger(TriggeredScraperTask.class);

    private final PlcDriverManager driverManager;
    private final String jobName;
    private final String connectionAlias;
    private final String connectionString;
    private final Map<String, String> fields;
    private final long requestTimeoutMs;
    private final ExecutorService executorService;
    private final ResultHandler resultHandler;
    private final TriggerHandler triggerHandler;

    private final AtomicLong requestCounter = new AtomicLong(0);
    private final AtomicLong successCounter = new AtomicLong(0);
    private final DescriptiveStatistics latencyStatistics = new DescriptiveStatistics(1000);
    private final DescriptiveStatistics failedStatistics = new DescriptiveStatistics(1000);

    public TriggeredScraperTask(PlcDriverManager driverManager,
                                String jobName,
                                String connectionAlias,
                                String connectionString,
                                Map<String, String> fields,
                                long requestTimeoutMs,
                                ExecutorService executorService,
                                ResultHandler resultHandler,
                                TriggeredScrapeJobImpl triggeredScrapeJob) throws ScraperException {
        this.driverManager = driverManager;
        this.jobName = jobName;
        this.connectionAlias = connectionAlias;
        this.connectionString = connectionString;
        this.fields = fields;
        this.requestTimeoutMs = requestTimeoutMs;
        this.executorService = executorService;
        this.resultHandler = resultHandler;
        this.triggerHandler = new TriggerHandler(triggeredScrapeJob.getTriggerConfig(),triggeredScrapeJob,this);
    }

    @Override
    //ToDo code-refactoring and improved testing --> PLC4X-90
    public void run() {
        if(this.triggerHandler.checkTrigger()) {
            // Does a single fetch only when trigger is valid
            LOGGER.trace("Start new scrape of task of job {} for connection {}", jobName, connectionAlias);
            requestCounter.incrementAndGet();
            StopWatch stopWatch = new StopWatch();
            stopWatch.start();
            PlcConnection connection = null;
            try {
                CompletableFuture<PlcConnection> future = CompletableFuture.supplyAsync(() -> {
                    try {
                        return driverManager.getConnection(connectionString);
                    } catch (PlcConnectionException e) {
                        LOGGER.warn("Unable to instantiate connection to " + connectionString, e);
                        throw new PlcRuntimeException(e);
                    }
                }, executorService);
                connection = future.get(10 * requestTimeoutMs, TimeUnit.MILLISECONDS);
                LOGGER.trace("Connection to {} established: {}", connectionString, connection);
                PlcReadResponse response;
                try {
                    PlcReadRequest.Builder builder = connection.readRequestBuilder();
                    fields.forEach((alias, qry) -> {
                        LOGGER.trace("Requesting: {} -> {}", alias, qry);
                        builder.addItem(alias, qry);
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
                CompletableFuture.runAsync(() -> resultHandler.handle(jobName, connectionAlias, transformResponseToMap(response)), executorService);
            } catch (Exception e) {
                LOGGER.debug("Exception during scrape", e);
                handleException(e);
            } finally {
                if (connection != null) {
                    try {
                        connection.close();
                    } catch (Exception e) {
                        LOGGER.warn("Error on closing connection",e);
                    }
                }
            }
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

    @Override
    public String getJobName() {
        return null;
    }

    @Override
    public String getConnectionAlias() {
        return null;
    }

    @Override
    public long getRequestCounter() {
        return 0;
    }

    @Override
    public long getSuccessfullRequestCounter() {
        return 0;
    }

    @Override
    public DescriptiveStatistics getLatencyStatistics() {
        return null;
    }

    @Override
    public double getPercentageFailed() {
        return 0;
    }

    @Override
    public void handleException(Exception e) {

    }

    @Override
    public void handleErrorResponse(Map<String, PlcResponseCode> failed) {

    }

    public PlcDriverManager getDriverManager() {
        return driverManager;
    }

    public String getConnectionString() {
        return connectionString;
    }

    public ExecutorService getExecutorService() {
        return executorService;
    }

    public long getRequestTimeoutMs() {
        return requestTimeoutMs;
    }
}
