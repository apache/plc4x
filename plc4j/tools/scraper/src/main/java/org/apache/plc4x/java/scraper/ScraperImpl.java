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

import org.apache.commons.collections4.MultiValuedMap;
import org.apache.commons.collections4.multimap.ArrayListValuedHashMap;
import org.apache.commons.lang3.Validate;
import org.apache.commons.lang3.concurrent.BasicThreadFactory;
import org.apache.commons.lang3.tuple.Triple;
import org.apache.commons.math3.stat.descriptive.DescriptiveStatistics;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.util.PercentageAboveThreshold;
import org.apache.plc4x.java.utils.cache.CachedPlcConnectionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.concurrent.*;

/**
 * Main class that orchestrates scraping.
 *
 * @deprecated Scraper is deprecated please use {@link org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperImpl} instead all functions are supplied as well
 */
@Deprecated
public class ScraperImpl implements Scraper {

    private static final Logger LOGGER = LoggerFactory.getLogger(ScraperImpl.class);

    private final ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(10,
        new BasicThreadFactory.Builder()
            .namingPattern("scheduler-thread-%d")
            .daemon(false)
            .build()
    );
    private final ExecutorService handlerPool = Executors.newFixedThreadPool(4,
        new BasicThreadFactory.Builder()
            .namingPattern("handler-thread-%d")
            .daemon(true)
            .build()
    );

    private final ResultHandler resultHandler;

    private final MultiValuedMap<ScrapeJob, ScraperTask> tasks = new ArrayListValuedHashMap<>();
    private final MultiValuedMap<ScraperTask, ScheduledFuture<?>> futures = new ArrayListValuedHashMap<>();
    private final PlcConnectionManager connectionManager;
    private final List<ScrapeJob> jobs;

    /**
     * default constructor
     * @param resultHandler handler for acquired data
     * @param connectionManager handler for Plc connection
     * @param jobs list of scrapings jobs to be executed
     */
    public ScraperImpl(ResultHandler resultHandler, PlcConnectionManager connectionManager, List<ScrapeJob> jobs) {
        this.resultHandler = resultHandler;
        Validate.notEmpty(jobs);
        this.connectionManager = connectionManager;
        this.jobs = jobs;
    }

    /**
     * Creates a Scraper instance from a configuration.
     * By default, a {@link CachedPlcConnectionManager} is used.
     * @param config Configuration to use.
     * @param resultHandler handler for acquired data
     * @throws ScraperException something went wrong ...
     */
    public ScraperImpl(ScraperConfiguration config, ResultHandler resultHandler) throws ScraperException {
        this(resultHandler, createCachedPlcConnectionManager(), config.getJobs());
    }

    /**
     * Min Idle per Key is set to 1 for situations where the network is broken.
     * Then, on reconnect we can fail all getConnection calls (in the ScraperTask) fast until
     * (in the background) the idle connection is created and the getConnection call returns fast.
     */
    private static CachedPlcConnectionManager createCachedPlcConnectionManager() {
        return CachedPlcConnectionManager.getBuilder().build();
    }

    @Override
    public void start() {
        // Schedule all jobs
        LOGGER.info("Starting jobs...");
        jobs.stream()
            .flatMap(job -> job.getSourceConnections().entrySet().stream()
                .map(entry -> Triple.of(job, entry.getKey(), entry.getValue()))
            )
            .forEach(
                tuple -> {
                    LOGGER.debug("Register task for job {} for conn {} ({}) at rate {} ms",
                        tuple.getLeft().getJobName(), tuple.getMiddle(), tuple.getRight(), tuple.getLeft().getScrapeRate());
                    ScraperTask task = new ScraperTaskImpl(connectionManager,
                        tuple.getLeft().getJobName(), tuple.getMiddle(), tuple.getRight(),
                        tuple.getLeft().getTags(),
                        1_000,
                        handlerPool, resultHandler);
                    // Add task to internal list
                    tasks.put(tuple.getLeft(), task);
                    ScheduledFuture<?> future = scheduler.scheduleAtFixedRate(task,
                        0, tuple.getLeft().getScrapeRate(), TimeUnit.MILLISECONDS);

                    // Store the handle for stopping, etc.
                    futures.put(task, future);
                }
            );

        // Add statistics tracker
        scheduler.scheduleAtFixedRate(() -> {
            for (Map.Entry<ScrapeJob, ScraperTask> entry : tasks.entries()) {
                DescriptiveStatistics statistics = entry.getValue().getLatencyStatistics();
                String msg = String.format(Locale.ENGLISH, "Job statistics (%s, %s) number of requests: %d (%d success, %.1f %% failed, %.1f %% too slow), min latency: %.2f ms, mean latency: %.2f ms, median: %.2f ms",
                    entry.getValue().getJobName(), entry.getValue().getConnectionAlias(),
                    entry.getValue().getRequestCounter(), entry.getValue().getSuccessfullRequestCounter(),
                    entry.getValue().getPercentageFailed(),
                    statistics.apply(new PercentageAboveThreshold(entry.getKey().getScrapeRate() * 1e6)),
                    statistics.getMin() * 1e-6, statistics.getMean() * 1e-6, statistics.getPercentile(50) * 1e-6);
                LOGGER.debug(msg);
            }
        }, 1_000, 1_000, TimeUnit.MILLISECONDS);
    }

    /**
     * For testing.
     */
    ScheduledExecutorService getScheduler() {
        return scheduler;
    }

    @Override
    public int getNumberOfActiveTasks() {
        return (int) futures.entries().stream().filter(entry -> !entry.getValue().isDone()).count();
    }

    @Override
    public void stop() {
        // Stop all futures
        LOGGER.info("Stopping scraper...");
        for (Map.Entry<ScraperTask, ScheduledFuture<?>> entry : futures.entries()) {
            LOGGER.debug("Stopping task {}...", entry.getKey());
            entry.getValue().cancel(true);
        }
        // Clear the map
        futures.clear();
    }

}
