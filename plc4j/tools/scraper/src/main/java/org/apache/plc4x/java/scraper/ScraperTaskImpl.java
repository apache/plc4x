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
package org.apache.plc4x.java.scraper;

import org.apache.commons.lang3.Validate;
import org.apache.commons.lang3.time.StopWatch;
import org.apache.commons.math3.stat.descriptive.DescriptiveStatistics;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.scraper.config.JobConfigurationClassicImpl;
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
 * Plc Scraper Task that scrapes one source.
 * One {@link ScrapeJobImpl} gets split into multiple tasks.
 * One task for each source that is defined in the {@link JobConfigurationClassicImpl}.
 *
 * @deprecated Scraper is deprecated please use {@link org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScrapeJobImpl} instead all functions are supplied as well see java-doc of {@link org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperImpl}
 */
@Deprecated
public class ScraperTaskImpl implements ScraperTask {

    private static final Logger LOGGER = LoggerFactory.getLogger(ScraperTaskImpl.class);

    private final PlcConnectionManager connectionManager;
    private final String jobName;
    private final String connectionAlias;
    private final String connectionString;
    private final Map<String, String> tags;
    private final long requestTimeoutMs;
    private final ExecutorService handlerService;
    private final ResultHandler resultHandler;

    private final AtomicLong requestCounter = new AtomicLong(0);
    private final AtomicLong successCounter = new AtomicLong(0);
    private final DescriptiveStatistics latencyStatistics = new DescriptiveStatistics(1000);
    private final DescriptiveStatistics failedStatistics = new DescriptiveStatistics(1000);

    public ScraperTaskImpl(PlcConnectionManager connectionManager,
                           String jobName,
                           String connectionAlias,
                           String connectionString,
                           Map<String, String> tags,
                           long requestTimeoutMs,
                           ExecutorService handlerService,
                           ResultHandler resultHandler) {
        Validate.notNull(connectionManager);
        Validate.notBlank(jobName);
        Validate.notBlank(connectionAlias);
        Validate.notBlank(connectionString);
        Validate.notEmpty(tags);
        Validate.isTrue(requestTimeoutMs > 0);
        Validate.notNull(resultHandler);
        this.connectionManager = connectionManager;
        this.jobName = jobName;
        this.connectionAlias = connectionAlias;
        this.connectionString = connectionString;
        this.tags = tags;
        this.requestTimeoutMs = requestTimeoutMs;
        this.handlerService = handlerService;
        this.resultHandler = resultHandler;
    }


    @Override
    public void run() {
        // Does a single fetch
        if(LOGGER.isDebugEnabled()) {
            LOGGER.debug("Start new scrape of task of job {} for connection {}", jobName, connectionAlias);
        }
        requestCounter.incrementAndGet();
        StopWatch stopWatch = new StopWatch();
        stopWatch.start();
        PlcConnection connection = null;
        try {
            CompletableFuture<PlcConnection> future = CompletableFuture.supplyAsync(() -> {
                try {
                    return connectionManager.getConnection(connectionString);
                } catch (PlcConnectionException e) {
                    LOGGER.warn("Unable to instantiate connection to " + connectionString, e);
                    throw new PlcRuntimeException(e);
                }
            }, handlerService);
            connection = future.get(10*requestTimeoutMs, TimeUnit.MILLISECONDS);
            LOGGER.debug("Connection to {} established: {}", connectionString, connection);
            PlcReadResponse plcReadResponse;
            try {
                //build read request
                PlcReadRequest.Builder readRequestBuilder = connection.readRequestBuilder();
                //add tags to be acquired to builder
                tags.forEach((alias, qry) -> {
                    LOGGER.trace("Requesting: {} -> {}", alias, qry);
                    readRequestBuilder.addTagAddress(alias, qry);
                });
                plcReadResponse = readRequestBuilder
                    .build()
                    .execute()
                    .get(requestTimeoutMs, TimeUnit.MILLISECONDS);
            } catch (ExecutionException e) {
                // Handle execution exception
                handleException(e);
                return;
            }

            LOGGER.debug("Performing statistics");
            // Add some statistics
            stopWatch.stop();
            latencyStatistics.addValue(stopWatch.getNanoTime());
            failedStatistics.addValue(0.0);
            successCounter.incrementAndGet();
            // Validate response
            validateResponse(plcReadResponse);

            // Handle response (Async)
            CompletableFuture.runAsync(() -> resultHandler.handle(jobName, connectionAlias, transformResponseToMap(plcReadResponse)), handlerService);
        } catch (Exception e) {
            LOGGER.warn("Exception during scraping of Job {}, Connection-Alias {}: Error-message: {} - for stack-trace change logging to DEBUG", jobName,connectionAlias,e.getMessage());
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

    /**
     * validate read response due to failed tags
     * @param response acquired response
     */
    private void validateResponse(PlcReadResponse response) {
        Map<String, PlcResponseCode> failedTags = response.getTagNames().stream()
            .filter(name -> !PlcResponseCode.OK.equals(response.getResponseCode(name)))
            .collect(Collectors.toMap(
                Function.identity(),
                response::getResponseCode
            ));
        if (failedTags.size() > 0) {
            handleErrorResponse(failedTags);
        }
    }

    /**
     * transforms the read-response to a Map of String (Key) and Object(Value)
     * @param response response from PLC
     * @return transformed Map
     */
    private Map<String, Object> transformResponseToMap(PlcReadResponse response) {
        return response.getTagNames().stream()
            .collect(Collectors.toMap(
                name -> name,
                response::getObject
            ));
    }

    @Override
    public String getJobName() {
        return jobName;
    }

    @Override
    public String getConnectionAlias() {
        return connectionAlias;
    }

    @Override
    public long getRequestCounter() {
        return requestCounter.get();
    }

    @Override
    public long getSuccessfullRequestCounter() {
        return successCounter.get();
    }

    @Override
    public DescriptiveStatistics getLatencyStatistics() {
        return latencyStatistics;
    }

    @Override
    public double getPercentageFailed() {
        return 100.0*failedStatistics.getMean();
    }

    @Override
    public void handleException(Exception e) {
        if(LOGGER.isDebugEnabled()) {
            LOGGER.debug("Detailed exception occurred at scraping", e);
        }
        failedStatistics.addValue(1.0);
    }

    @Override
    public void handleErrorResponse(Map<String, PlcResponseCode> failed) {
        LOGGER.warn("Handling error responses: {}", failed);
    }

}
