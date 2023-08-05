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
package org.apache.plc4x.merlot.scheduler.api;

import java.util.Date;
import java.util.Map;

/**
 * A scheduler to schedule time/cron based jobs.
 * A job is an object that is executed/fired by the scheduler. The object
 * should either implement the {@link Job} interface or the {@link Runnable}
 * interface.
 *
 * A job can be scheduled either by creating a {@link ScheduleOptions} instance
 * through one of the scheduler methods and then calling {@link #schedule(Object, ScheduleOptions)}
 * or
 * by using the whiteboard pattern and registering a Runnable service with either
 * the {@link #PROPERTY_SCHEDULER_EXPRESSION} or {@link #PROPERTY_SCHEDULER_PERIOD}
 * property. Services registered by the whiteboard pattern can by default run concurrently,
 * which usually is not wanted. Therefore it is advisable to also set the
 * {@link #PROPERTY_SCHEDULER_CONCURRENT} property with Boolean.FALSE.
 */
public interface Scheduler {

    /**
     * Name of the configuration property to define the period for a job.
     * The period is expressed in seconds.
     * This property needs to be of type Long.
     */
    String PROPERTY_SCHEDULER_PERIOD = "scheduler.period";

    /**
     * Name of the configuration property to defined the number of iterations for a job.
     * The times is expressed in iterations.
     * This property needs to be of numeric type.
     */
    String PROPERTY_SCHEDULER_TIMES = "scheduler.times";

    /**
     * Name of the configuration property to define if a periodically job should be scheduled immediate.
     * Default is to not startup immediate, the job is started the first time after the period has expired.
     * This property needs to be of type Boolean.
     */
    String PROPERTY_SCHEDULER_IMMEDIATE = "scheduler.immediate";

    /** Name of the configuration property to define the cron expression for a job. */
    String PROPERTY_SCHEDULER_EXPRESSION = "scheduler.expression";

    /** Name of the configuration property to define if the job can be run concurrently. */
    String PROPERTY_SCHEDULER_CONCURRENT = "scheduler.concurrent";

    /** Name of the configuration property to define the job name. */
    String PROPERTY_SCHEDULER_NAME = "scheduler.name";


    /**
     * Schedule a job based on the options.
     *
     * Note that if a job with the same name has already been added, the old job is cancelled and this new job replaces
     * the old job.
     *
     * The job object needs either to be a {@link Job} or a {@link Runnable}. The options have to be created
     * by one of the provided methods from this scheduler.
     *
     * @param job The job to execute (either {@link Job} or {@link Runnable}).
     * @param options Required options defining how to schedule the job.
     * @throws SchedulerError if the job can't be scheduled.
     * @throws IllegalArgumentException If the preconditions are not met.
     * @see #NOW()
     * @see #NOW(int, long)
     * @see #AT(Date)
     * @see #AT(Date, int, long)
     * @see #EXPR(String)
     */
    void schedule(Object job, ScheduleOptions options) throws IllegalArgumentException, SchedulerError;

    /**
     * Update the scheduling of an existing job.
     * @param jobName
     * @param options
     * @throws IllegalArgumentException
     * @throws SchedulerError
     */
    void reschedule(String jobName, ScheduleOptions options) throws IllegalArgumentException, SchedulerError;

    /**
     * Remove a scheduled job by name.
     *
     * @param jobName The name of the job.
     * @return <code>True</code> if the job existed and could be stopped, <code>false</code> otherwise.
     */
    boolean unschedule(String jobName);

    Map<String, ScheduleOptions> getJobs() throws SchedulerError;

    /**
     * Triggers a scheduled job.
     *
     * @param jobName The name of the job.
     * @return <code>true</code> if the job was triggered, otherwise <code>false</code>
     * @throws SchedulerError  if the job can't be triggered.
     */
    boolean trigger(String jobName) throws SchedulerError;

    /**
     * Create a schedule options to fire a job immediately and only once.
     *
     * @return The corresponding {@link ScheduleOptions}.
     */
    ScheduleOptions NOW();

    /**
     * Create a schedule options to fire a job immediately more than once.
     * @param times The number of times this job should be started (must be higher than 1 or -1 for endless).
     * @param period Every period seconds this job is started (must be at higher than 0).
     * @return The corresponding {@link ScheduleOptions}.
     */
    ScheduleOptions NOW(int times, long period);

    /**
     * Create a schedule options to fire a job once at a specific date.
     *
     * @param date The date this job should be run.
     * @return The corresponding {@link ScheduleOptions}.
     */
    ScheduleOptions AT(final Date date);

    /**
     * Create a schedule options to fire a job period starting at a specific date.
     *
     * @param date The date this job should be run.
     * @param times The number of times this job should be started (must be higher than 1 or -1 for endless).
     * @param period Every period seconds this job is started (must be at higher than 0).
     * @return The corresponding {@link ScheduleOptions}.
     */
    ScheduleOptions AT(final Date date, int times, long period);

    /**
     * Create a schedule options to schedule the job based on the expression.
     *
     * @param expression The cron exception.
     * @return The corresponding {@link ScheduleOptions}.
     */
    ScheduleOptions EXPR(final String expression);

}
