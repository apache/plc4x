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

import org.apache.commons.collections4.MultiValuedMap;
import org.apache.commons.collections4.multimap.ArrayListValuedHashMap;
import org.apache.commons.lang3.Validate;
import org.apache.commons.lang3.concurrent.BasicThreadFactory;
import org.apache.commons.lang3.tuple.Triple;
import org.apache.commons.math3.stat.descriptive.DescriptiveStatistics;
import org.apache.commons.pool2.impl.GenericKeyedObjectPool;
import org.apache.commons.pool2.impl.GenericKeyedObjectPoolConfig;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.scraper.*;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.config.triggeredscraper.TriggeredScraperConfiguration;
import org.apache.plc4x.java.scraper.util.PercentageAboveThreshold;
import org.apache.plc4x.java.utils.connectionpool.PooledPlcDriverManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.concurrent.*;

/**
 * replaces the old Scraper that only could do scheduled scraping jobs
 * Triggers have been introduced, so that in configuration "scrapeTime" has been exchanged by "triggerConfig"
 *
 * Some example:
 *   - 200ms scheduling is now performed by "triggerConfig: (SCHEDULED,200)" in scraper-configuration
 *   - a triggered S7 variable can be used as follows:
 *     "triggerConfig: (S7_TRIGGER_VAR,10,(%M0.3:BOOL)==(true))" meaning that Boolean in Marker-Block in Byte-Offset 0, Bit-Offset 3 is scanned every 10ms, when trigger has a rising-edge the acquirement of data-block is triggered
 *     the trigger variable must be a valid address as defined with PLC4X-S7-Driver
 *     right now boolean variables as well as numeric variables could be used as data-types
 *     available comparators are ==,!= for all data-types and &gt;,&gt;=,&lt;,&lt;= for numeric data-types
 */
public class TriggeredScraperImpl implements Scraper {

    private static final Logger LOGGER = LoggerFactory.getLogger(TriggeredScraperImpl.class);

    private final ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(10,
        new BasicThreadFactory.Builder()
            .namingPattern("triggeredscraper-scheduling-thread-%d")
            .daemon(false)
            .build()
    );
    private final ExecutorService executorService = Executors.newFixedThreadPool(4,
        new BasicThreadFactory.Builder()
            .namingPattern("triggeredscraper-executer-thread-%d")
            .daemon(true)
            .build()
    );

    private final ResultHandler resultHandler;

    private final MultiValuedMap<ScrapeJob, ScraperTask> tasks = new ArrayListValuedHashMap<>();
    private final MultiValuedMap<ScraperTask, ScheduledFuture<?>> futures = new ArrayListValuedHashMap<>();
    private final PlcDriverManager driverManager;
    private final List<ScrapeJob> jobs;

    /**
     * Creates a Scraper instance from a configuration.
     * By default a {@link PooledPlcDriverManager} is used.
     * @param config Configuration to use.
     * @param resultHandler
     */
    public TriggeredScraperImpl(TriggeredScraperConfiguration config, ResultHandler resultHandler) throws ScraperException {
        this(resultHandler, createPooledDriverManager(), config.getJobs());
    }

    /**
     * Min Idle per Key is set to 1 for situations where the network is broken.
     * Then, on reconnect we can fail all getConnection calls (in the ScraperTask) fast until
     * (in the background) the idle connection is created and the getConnection call returns fast.
     */
    private static PooledPlcDriverManager createPooledDriverManager() {
        return new PooledPlcDriverManager(pooledPlcConnectionFactory -> {
            GenericKeyedObjectPoolConfig<PlcConnection> poolConfig = new GenericKeyedObjectPoolConfig<>();
            poolConfig.setMinIdlePerKey(1);  // This should avoid problems with long running connect attempts??
            poolConfig.setTestOnBorrow(true);
            poolConfig.setTestOnReturn(true);
            return new GenericKeyedObjectPool<>(pooledPlcConnectionFactory, poolConfig);
        });
    }


    public TriggeredScraperImpl(ResultHandler resultHandler, PlcDriverManager driverManager, List<ScrapeJob> jobs) {
        this.resultHandler = resultHandler;
        Validate.notEmpty(jobs);
        this.driverManager = driverManager;
        this.jobs = jobs;
    }

    /**
     * Start the scraping.
     */
    //ToDo code-refactoring and improved testing --> PLC4X-90
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
                    TriggeredScraperTask task =
                        null;
                    try {
                        task = new TriggeredScraperTask(driverManager,
                            tuple.getLeft().getJobName(),
                            tuple.getMiddle(),
                            tuple.getRight(),
                            tuple.getLeft().getFields(),
                            1_000,
                            executorService,
                            resultHandler,
                            (TriggeredScrapeJobImpl) tuple.getLeft());
                        // Add task to internal list
                        tasks.put(tuple.getLeft(), task);
                        ScheduledFuture<?> future = scheduler.scheduleAtFixedRate(task,
                            0, tuple.getLeft().getScrapeRate(), TimeUnit.MILLISECONDS);

                        // Store the handle for stopping, etc.
                        futures.put(task, future);
                    } catch (ScraperException e) {
                        LOGGER.warn("Error executing the job {} for conn {} ({}) at rate {} ms",tuple.getLeft().getJobName(), tuple.getMiddle(), tuple.getRight(), tuple.getLeft().getScrapeRate(),e);
                    }

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

    @Override
    public int getNumberOfActiveTasks() {
        return 0;
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
