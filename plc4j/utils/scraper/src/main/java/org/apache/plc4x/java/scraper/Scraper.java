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

import org.apache.commons.collections4.MultiValuedMap;
import org.apache.commons.collections4.multimap.ArrayListValuedHashMap;
import org.apache.commons.lang3.Validate;
import org.apache.commons.lang3.concurrent.BasicThreadFactory;
import org.apache.commons.lang3.tuple.Triple;
import org.apache.commons.math3.exception.MathIllegalArgumentException;
import org.apache.commons.math3.stat.descriptive.DescriptiveStatistics;
import org.apache.commons.math3.stat.descriptive.UnivariateStatistic;
import org.apache.plc4x.java.PlcDriverManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.concurrent.*;
import java.util.stream.IntStream;

/**
 * Main class that orchestrates scraping.
 */
public class Scraper {

    private static final Logger LOGGER = LoggerFactory.getLogger(Scraper.class);

    private final ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(10,
        new BasicThreadFactory.Builder()
            .namingPattern("scheduler-thread-%d")
            .daemon(true)
            .build()
    );
    private final ExecutorService handlerPool = Executors.newFixedThreadPool(4,
        new BasicThreadFactory.Builder()
            .namingPattern("handler-thread-%d")
            .daemon(true)
            .build()
    );
    private final MultiValuedMap<ScrapeJob, ScraperTask> tasks = new ArrayListValuedHashMap<>();
    private final MultiValuedMap<ScraperTask, ScheduledFuture<?>> futures = new ArrayListValuedHashMap<>();
    private final PlcDriverManager driverManager;
    private final List<ScrapeJob> jobs;

    public Scraper(PlcDriverManager driverManager, List<ScrapeJob> jobs) {
        Validate.notEmpty(jobs);
        this.driverManager = driverManager;
        this.jobs = jobs;
    }

    /**
     * Start the scraping.
     */
    public void start() {
        // Schedule all jobs
        LOGGER.info("Starting jobs...");
        jobs.stream()
            .flatMap(job -> job.connections.entrySet().stream()
                .map(entry -> Triple.of(job, entry.getKey(), entry.getValue()))
            )
            .forEach(
                tuple -> {
                    LOGGER.debug("Register task for job {} for conn {} ({}) at rate {} ms",
                        tuple.getLeft().name, tuple.getMiddle(), tuple.getRight(), tuple.getLeft().scrapeRate);
                    ScraperTask task = new ScraperTask(driverManager,
                        tuple.getLeft().name, tuple.getMiddle(), tuple.getRight(),
                        tuple.getLeft().fields,
                        1_000,
                        handlerPool);
                    // Add task to internal list
                    tasks.put(tuple.getLeft(), task);
                    ScheduledFuture<?> future = scheduler.scheduleAtFixedRate(task,
                        0, tuple.getLeft().scrapeRate, TimeUnit.MILLISECONDS);

                    // Store the handle for stopping, etc.
                    futures.put(task, future);
                }
            );

        // Add statistics tracker
        scheduler.scheduleAtFixedRate(() -> {
            for (Map.Entry<ScrapeJob, ScraperTask> entry : tasks.entries()) {
                DescriptiveStatistics statistics = entry.getValue().getLatencyStatistics();
                String msg = String.format(Locale.ENGLISH, "Job statistics (%s, %s) number of requests: %d (%d success, %.1f %% failed, %.1f %% too slow), mean latency: %.2f ms, median: %.2f ms",
                    entry.getValue().getJobName(), entry.getValue().getConnectionAlias(), entry.getValue().getRequestCounter(), entry.getValue().getSuccessfullRequestCounter(), entry.getValue().getPercentageFailed(), statistics.apply(new PercentageAboveThreshold(entry.getKey().scrapeRate * 1e6)), statistics.getMean() * 1e-6, statistics.getPercentile(50) * 1e-6);
                LOGGER.info(msg);
            }
        }, 1_000, 1_000, TimeUnit.MILLISECONDS);
    }

    /**
     * For testing.
     */
    ScheduledExecutorService getScheduler() {
        return scheduler;
    }

    public int getNumberOfActiveTasks() {
        return (int) futures.entries().stream().filter(entry -> !entry.getValue().isDone()).count();
    }

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

    public static class ScrapeJob {

        private final String name;
        private final long scrapeRate;
        /**
         * alias -> connection-string
         */
        private final Map<String, String> connections;
        /**
         * alias -> field-query
         */
        private final Map<String, String> fields;

        public ScrapeJob(String name, long scrapeRate, Map<String, String> connections, Map<String, String> fields) {
            this.name = name;
            this.scrapeRate = scrapeRate;
            this.connections = connections;
            this.fields = fields;
        }
    }

    private static class PercentageAboveThreshold implements UnivariateStatistic {

        private final double threshold;

        public PercentageAboveThreshold(double threshold) {
            this.threshold = threshold;
        }

        @Override
        public double evaluate(double[] values) throws MathIllegalArgumentException {
            long below = Arrays.stream(values)
                .filter(val -> val <= threshold)
                .count();
            return (double) below / values.length;
        }

        @Override
        public double evaluate(double[] values, int begin, int length) throws MathIllegalArgumentException {
            long below = IntStream.range(begin, length)
                .mapToDouble(i -> values[i])
                .filter(val -> val > threshold)
                .count();
            return 100.0 * below / length;
        }

        @Override
        public UnivariateStatistic copy() {
            return new PercentageAboveThreshold(threshold);
        }
    }
}
