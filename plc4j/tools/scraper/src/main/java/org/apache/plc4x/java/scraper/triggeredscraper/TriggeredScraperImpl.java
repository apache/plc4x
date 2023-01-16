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
package org.apache.plc4x.java.scraper.triggeredscraper;

import org.apache.commons.collections4.MultiValuedMap;
import org.apache.commons.collections4.multimap.ArrayListValuedHashMap;
import org.apache.commons.lang3.Validate;
import org.apache.commons.lang3.concurrent.BasicThreadFactory;
import org.apache.commons.math3.stat.descriptive.DescriptiveStatistics;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.scraper.*;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.config.triggeredscraper.ScraperConfigurationTriggeredImpl;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector.TriggerCollector;
import org.apache.plc4x.java.scraper.util.PercentageAboveThreshold;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.utils.cache.CachedPlcConnectionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.management.*;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.concurrent.*;
import java.util.stream.Collectors;

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
public class TriggeredScraperImpl implements Scraper, TriggeredScraperMBean {

    private static final Logger LOGGER = LoggerFactory.getLogger(TriggeredScraperImpl.class);
    private static final String MX_DOMAIN = "org.apache.plc4x.java";

    private static final int DEFAULT_FUTURE_TIME_OUT = 2000;

    private final ScheduledExecutorService scheduler;
    private final ExecutorService executorService;
    private ScheduledFuture<?> statisticsLogger;

    private final ResultHandler resultHandler;

    private final MultiValuedMap<ScrapeJob, ScraperTask> tasks = new ArrayListValuedHashMap<>();
    private final MultiValuedMap<ScraperTask, ScheduledFuture<?>> scraperTaskMap = new ArrayListValuedHashMap<>();
    private final PlcConnectionManager plcConnectionManager;
    private final List<ScrapeJob> jobs;
    private MBeanServer mBeanServer;

    private long futureTimeOut;

    private final TriggerCollector triggerCollector;

    /**
     * Creates a Scraper instance from a configuration.
     * By default a {@link CachedPlcConnectionManager} is used.
     * @param config Configuration to use.
     * @param resultHandler handler the defines the processing of acquired data
     * @param triggerCollector the trigger collector
     * @throws ScraperException something went wrong
     */
    public TriggeredScraperImpl(ScraperConfiguration config, ResultHandler resultHandler, TriggerCollector triggerCollector) throws ScraperException {
        this(resultHandler, createCachedPlcConnectionManager(), config.getJobs(),triggerCollector,DEFAULT_FUTURE_TIME_OUT);
    }

    /**
     * Creates a Scraper instance from a configuration.
     * @param config Configuration to use.
     * @param plcConnectionManager external DriverManager
     * @param resultHandler handler the defines the processing of acquired data
     * @param triggerCollector the trigger collector
     * @throws ScraperException something went wrong
     */
    public TriggeredScraperImpl(ScraperConfiguration config, PlcConnectionManager plcConnectionManager, ResultHandler resultHandler, TriggerCollector triggerCollector) throws ScraperException {
        this(resultHandler, plcConnectionManager, config.getJobs(),triggerCollector,DEFAULT_FUTURE_TIME_OUT);
    }

    /**
     * Creates a Scraper instance from a configuration.
     * @param config Configuration to use.
     * @param plcConnectionManager external DriverManager
     * @param resultHandler handler the defines the processing of acquired data
     * @param triggerCollector the trigger collector
     * @param poolSizeExecutor the pool size of the executor
     * @param poolSizeScheduler the pool size of the scheduler
     * @throws ScraperException something went wrong
     */
    public TriggeredScraperImpl(ScraperConfigurationTriggeredImpl config, PlcConnectionManager plcConnectionManager, ResultHandler resultHandler, TriggerCollector triggerCollector, int poolSizeScheduler, int poolSizeExecutor) throws ScraperException {
        this(resultHandler, plcConnectionManager, config.getJobs(),triggerCollector,DEFAULT_FUTURE_TIME_OUT,poolSizeScheduler,poolSizeExecutor);
    }

    /**
     * Creates a Scraper instance from a configuration.
     * @param plcConnectionManager external DriverManager
     * @param resultHandler handler the defines the processing of acquired data
     * @param jobs list of jobs that scraper shall handle
     * @param triggerCollector a collection that centralizes the trigger requests and joins them to grouped plc requests
     * @param futureTimeOut max duration of future to return a result
     */
    public TriggeredScraperImpl(ResultHandler resultHandler, PlcConnectionManager plcConnectionManager, List<ScrapeJob> jobs,TriggerCollector triggerCollector, long futureTimeOut) {
        this(resultHandler,plcConnectionManager,jobs,triggerCollector,futureTimeOut,20,5);
    }

    public TriggeredScraperImpl(ResultHandler resultHandler, PlcConnectionManager plcConnectionManager, List<ScrapeJob> jobs,TriggerCollector triggerCollector, long futureTimeOut, int poolSizeScheduler, int poolSizeExecutor) {
        this.resultHandler = resultHandler;
        Validate.notEmpty(jobs);
        if (!(plcConnectionManager instanceof CachedPlcConnectionManager)) {
            LOGGER.warn("The Triggered Scraper is intended to be used with a cached PlcConnectionManager. In other situations leaks could occur!");
        }
        this.plcConnectionManager = plcConnectionManager;
        this.jobs = jobs;
        this.triggerCollector = triggerCollector;
        this.futureTimeOut = futureTimeOut;

        this.scheduler = Executors.newScheduledThreadPool(poolSizeScheduler,
            new BasicThreadFactory.Builder()
                .namingPattern("triggeredscraper-scheduling-thread-%d")
                .daemon(false)
                .build()
        );

        this.executorService = Executors.newFixedThreadPool(poolSizeExecutor,
            new BasicThreadFactory.Builder()
                .namingPattern("triggeredscraper-executor-thread-%d")
                .daemon(true)
                .build()
        );


        // Register MBean
        /*mBeanServer = ManagementFactory.getPlatformMBeanServer();
        try {
            mBeanServer.registerMBean(this, new ObjectName(MX_DOMAIN, "scraper", "scraper"));
        } catch (InstanceAlreadyExistsException | MBeanRegistrationException | NotCompliantMBeanException | MalformedObjectNameException e) {
            LOGGER.debug("Unable to register Scraper as MBean", e);
        }*/
    }


    /**
     * Min Idle per Key is set to 1 for situations where the network is broken.
     * Then, on reconnect we can fail all getConnection calls (in the ScraperTask) fast until
     * (in the background) the idle connection is created and the getConnection call returns fast.
     */
    private static CachedPlcConnectionManager createCachedPlcConnectionManager() {
        return CachedPlcConnectionManager.getBuilder().build();
    }

    /**
     * Start the scraping.
     */
    //ToDo code-refactoring and improved testing --> PLC4X-90
    @Override
    public void start() {
        // Schedule all jobs
        LOGGER.info("Starting jobs...");
        //start iterating over all available jobs
        for(ScrapeJob job:jobs){
            //iterate over all source the jobs shall performed on
            for(Map.Entry<String,String> sourceEntry:job.getSourceConnections().entrySet()){
                if(LOGGER.isDebugEnabled()) {
                    LOGGER.debug("Register task for job {} for conn {} ({}) at rate {} ms",
                        job.getJobName(),
                        sourceEntry.getKey(),
                        sourceEntry.getValue(),
                        job.getScrapeRate());
                }

                //create the regarding triggered scraper task
                TriggeredScraperTask triggeredScraperTask;
                try {
                    triggeredScraperTask = new TriggeredScraperTask(
                        plcConnectionManager,
                        job.getJobName(),
                        sourceEntry.getKey(),
                        sourceEntry.getValue(),
                        job.getTags(),
                        futureTimeOut,
                        executorService,
                        resultHandler,
                        (TriggeredScrapeJobImpl) job,
                        triggerCollector);

                    // Add task to internal list
                    if(LOGGER.isInfoEnabled()) {
                        LOGGER.info("Task {} added to scheduling", triggeredScraperTask);
                    }
                    registerTaskMBean(triggeredScraperTask);
                    tasks.put(job, triggeredScraperTask);
                    ScheduledFuture<?> future = scheduler.scheduleAtFixedRate(triggeredScraperTask, 0, job.getScrapeRate(), TimeUnit.MILLISECONDS);

                    // Store the handle for stopping, etc.
                    scraperTaskMap.put(triggeredScraperTask, future);
                } catch (ScraperException e) {
                    LOGGER.warn("Error executing the job {} for conn {} ({}) at rate {} ms",job.getJobName(), sourceEntry.getKey(), sourceEntry.getValue(), job.getScrapeRate(),e);
                }
            }

        }

        // Add statistics tracker
        statisticsLogger = scheduler.scheduleAtFixedRate(() -> {
            for (Map.Entry<ScrapeJob, ScraperTask> entry : tasks.entries()) {
                DescriptiveStatistics statistics = entry.getValue().getLatencyStatistics();
                String msg = String.format(Locale.ENGLISH, "Job statistics (%s, %s) number of requests: %d (%d success, %.1f %% failed, %.1f %% too slow), min latency: %.2f ms, mean latency: %.2f ms, median: %.2f ms",
                    entry.getValue().getJobName(), entry.getValue().getConnectionAlias(),
                    entry.getValue().getRequestCounter(), entry.getValue().getSuccessfullRequestCounter(),
                    entry.getValue().getPercentageFailed(),
                    statistics.apply(new PercentageAboveThreshold(entry.getKey().getScrapeRate() * 1e6)),
                    statistics.getMin() * 1e-6, statistics.getMean() * 1e-6, statistics.getPercentile(50) * 1e-6);
                if (LOGGER.isDebugEnabled()) {
                    LOGGER.debug(msg);
                }
            }
        }, 1_000, 1_000, TimeUnit.MILLISECONDS);
    }

    /**
     * Register a task as MBean
     * @param task task to register
     */
    private void registerTaskMBean(ScraperTask task) {
        /*try {
            mBeanServer.registerMBean(task, new ObjectName(MX_DOMAIN + ":type=ScrapeTask,name=" + task.getJobName() + "-" + task.getConnectionAlias()));
        } catch (InstanceAlreadyExistsException | MBeanRegistrationException | NotCompliantMBeanException | MalformedObjectNameException e) {
            LOGGER.debug("Unable to register Task as MBean", e);
        }*/
    }

    @Override
    public void stop() {
        // Stop all futures
        LOGGER.info("Stopping scraper...");
        for (Map.Entry<ScraperTask, ScheduledFuture<?>> entry : scraperTaskMap.entries()) {
            LOGGER.debug("Stopping task {}...", entry.getKey());
            entry.getValue().cancel(true);
        }
        // Clear the map
        scraperTaskMap.clear();

        // Stop the statistics logger, if it is currently running.
        if((statisticsLogger != null) && (!statisticsLogger.isCancelled())) {
            statisticsLogger.cancel(false);
        }
    }

    /**
     * acquires a plc connection from connection pool
     * @param plcConnectionManager  Connection manager handling connection and pools
     * @param connectionString      Connection string as defined in the regarding implementation of {@link PlcDriver}
     * @param executorService       ExecutorService holding a pool as threads handling requests and stuff
     * @param requestTimeoutMs      maximum wait time for the future to return a result
     * @param info                  additional info for trace reasons
     * @return the {@link PlcConnection} used for acquiring data from PLC endpoint
     * @throws InterruptedException something went wrong
     * @throws ExecutionException something went wrong
     * @throws TimeoutException something went wrong
     */
    public static PlcConnection getPlcConnection(PlcConnectionManager plcConnectionManager,
                                                 String connectionString,
                                                 ExecutorService executorService,
                                                 long requestTimeoutMs,
                                                 String info) throws InterruptedException, ExecutionException, TimeoutException {
        if(!info.isEmpty() && LOGGER.isTraceEnabled()){
            LOGGER.trace("Additional Info from caller {}", info);
        }
        CompletableFuture<PlcConnection> future = CompletableFuture.supplyAsync(() -> {
            try {
                return plcConnectionManager.getConnection(connectionString);
            } catch (PlcConnectionException e) {
                LOGGER.warn("Unable to instantiate connection to " + connectionString, e);
                throw new PlcRuntimeException(e);
            }
            catch (Exception e){
                LOGGER.warn("Unable to instantiate connection to " + connectionString, e);
                throw new PlcRuntimeException(e);
            }
        }, executorService);
        if(LOGGER.isTraceEnabled()){
            LOGGER.trace("try to get a connection to {}", connectionString);
        }
        PlcConnection plcConnection=null;
        try {
            plcConnection = future.get(requestTimeoutMs, TimeUnit.MILLISECONDS);
        }
        catch (Exception e){
            LOGGER.trace("Additional Info from caller {}", info,e);
            throw e;
        }
        return plcConnection;
    }

    /**
     * acquires a plc connection from connection pool
     * @param plcDriverManager  Driver manager handling connection and pools
     * @param connectionString  Connection string as defined in the regarding implementation of {@link PlcDriver}
     * @param executorService   ExecuterService holding a pool as threads handling requests and stuff
     * @param requestTimeoutMs  maximum awaiting for the the future to return a result
     * @return the {@link PlcConnection} used for acquiring data from PLC endpoint
     * @throws InterruptedException something went wrong
     * @throws ExecutionException something went wrong
     * @throws TimeoutException something went wrong
     */
    public static PlcConnection getPlcConnection(PlcDriverManager plcDriverManager,
                                                 String connectionString,
                                                 ExecutorService executorService,
                                                 long requestTimeoutMs) throws InterruptedException, ExecutionException, TimeoutException {
        return getPlcConnection(plcDriverManager,connectionString,executorService,requestTimeoutMs,"");
    }

    /**
     * transforms the results from a {@link PlcReadResponse} into a map
     * @param plcReadResponse response that shall be converted to map for further processing
     * @return the converted map
     */
    public static Map<String, Object> convertPlcResponseToMap(PlcReadResponse plcReadResponse) {
        return plcReadResponse.getTagNames().stream()
            .collect(Collectors.toMap(
                name -> name,
                plcReadResponse::getObject
            ));
    }


    // MBean methods
    @Override
    public boolean isRunning() {
        // TODO is this okay so?
        return !scraperTaskMap.isEmpty();
    }

    @Override
    public int getNumberOfActiveTasks() {
        return (int) scraperTaskMap.entries().stream().filter(entry -> !entry.getValue().isDone()).count();
    }
}
