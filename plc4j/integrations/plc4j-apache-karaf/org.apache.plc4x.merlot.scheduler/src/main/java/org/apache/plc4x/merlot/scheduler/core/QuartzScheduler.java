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
package org.apache.plc4x.merlot.scheduler.core;

import org.apache.plc4x.merlot.scheduler.api.Job;
import org.apache.plc4x.merlot.scheduler.api.ScheduleOptions;
import org.apache.plc4x.merlot.scheduler.api.Scheduler;
import org.apache.plc4x.merlot.scheduler.api.SchedulerError;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;
import java.util.Properties;
import java.util.UUID;
import org.quartz.*;
import org.quartz.impl.matchers.GroupMatcher;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


/**
 * The quartz based implementation of the scheduler.
 *
 */
public class QuartzScheduler implements Scheduler {

    /** Default logger. */
    private final Logger logger = LoggerFactory.getLogger(this.getClass());

    private static final String PREFIX = "Apache Karaf Quartz Scheduler ";

    /** Map key for the job object */
    static final String DATA_MAP_OBJECT = "QuartzJobScheduler.Object";

    /** Map key for the job name */
    static final String DATA_MAP_NAME = "QuartzJobScheduler.JobName";

    /** Map key for the scheduling options. */
    static final String DATA_MAP_OPTIONS = "QuartzJobScheduler.Options";

    /** Map key for non serializable context. */
    static final String DATA_MAP_CONTEXT = "QuartzJobScheduler.Context";

    /** Map key for the logger. */
    static final String DATA_MAP_LOGGER = "QuartzJobScheduler.Logger";

    /** The quartz scheduler. */
    private volatile org.quartz.Scheduler scheduler;

    public QuartzScheduler(Properties configuration) {
        // SLING-2261 Prevent Quartz from checking for updates
        System.setProperty("org.terracotta.quartz.skipUpdateCheck", Boolean.TRUE.toString());
        ClassLoader cl = Thread.currentThread().getContextClassLoader();
        try {
            Thread.currentThread().setContextClassLoader(QuartzScheduler.class.getClassLoader());
            KarafStdSchedulerFactory factory = new KarafStdSchedulerFactory(configuration);
            scheduler = factory.getScheduler();
            scheduler.start();
        } catch (Throwable t) {
            throw new RuntimeException("Unable to create quartz scheduler", t);
        } finally {
            Thread.currentThread().setContextClassLoader(cl);
        }
    }

    /**
     * Deactivate this component.
     * Stop the scheduler.
     */
    public void deactivate() {
        final org.quartz.Scheduler s = this.scheduler;
        this.scheduler = null;
        this.dispose(s);
    }

    /**
     * Dispose the quartz scheduler
     * @param s The scheduler.
     */
    private void dispose(final org.quartz.Scheduler s) {
        if ( s != null ) {
            try {
                s.shutdown();
            } catch (SchedulerException e) {
                this.logger.debug("Exception during shutdown of scheduler.", e);
            }
            if ( this.logger.isDebugEnabled() ) {
                this.logger.debug(PREFIX + "stopped.");
            }
        }
    }

    /**
     * Initialize the data map for the job executor.
     */
    private JobDataMap initDataMap(final String  jobName,
                                   final Object  job,
                                   final InternalScheduleOptions options) {
        final JobDataMap jobDataMap = new JobDataMap();
        final JobDataMap jobContextMap = new JobDataMap();

        // serializable data
        jobDataMap.put(DATA_MAP_NAME, jobName);
        jobDataMap.put(DATA_MAP_OPTIONS, options);

        // non serializable data
        jobContextMap.put(DATA_MAP_OBJECT, job);
        jobContextMap.put(DATA_MAP_LOGGER, this.logger);

        // temporary storage
        jobDataMap.put(DATA_MAP_CONTEXT, jobContextMap);

        return jobDataMap;
    }

    /**
     * Create the job detail.
     */
    private JobDetail createJobDetail(final String name,
                                      final JobDataMap jobDataMap,
                                      final boolean concurrent) {
        return JobBuilder.newJob((concurrent ? QuartzJobExecutor.class : NonParallelQuartzJobExecutor.class))
                .withIdentity(name)
                .usingJobData(jobDataMap)
                .build();
    }

    /**
     * Check the job object, only Merlot Job is allowed
     */
    private void checkJob(final Object job)
            throws IllegalArgumentException {
        if (!(job instanceof Job)) {
            System.out.println("Objeto: " + job.toString());
            System.out.println("Objeto: " + job.getClass().getName());
            throw new IllegalArgumentException("Job object is neither an instance of " + Runnable.class.getName() + " nor " + Job.class.getName());
        }
    }

    /** Used by the web console plugin. */
    org.quartz.Scheduler getScheduler() {
        return this.scheduler;
    }

    /**
     * @see org.apache.karaf.scheduler.Scheduler#NOW()
     */
    public ScheduleOptions NOW() {
        return AT(new Date());
    }

    /**
     * @see org.apache.karaf.scheduler.Scheduler#NOW(int, long)
     */
    public ScheduleOptions NOW(int times, long period) {
        return AT(new Date(), times, period);
    }

    /**
     * @see org.apache.karaf.scheduler.Scheduler#AT(java.util.Date)
     */
    public ScheduleOptions AT(Date date) {
        return new InternalScheduleOptions(date);
    }

    /**
     * @see org.apache.karaf.scheduler.Scheduler#AT(java.util.Date, int, long)
     */
    public ScheduleOptions AT(Date date, int times, long period) {
        return new InternalScheduleOptions(date, times, period);
    }

    /**
     * @see org.apache.karaf.scheduler.Scheduler#EXPR(java.lang.String)
     */
    public ScheduleOptions EXPR(String expression) {
        return new InternalScheduleOptions(expression);
    }

    /**
     * Schedule a job
     * @see org.apache.karaf.scheduler.Scheduler#schedule(java.lang.Object, org.apache.karaf.scheduler.ScheduleOptions)
     * @throws SchedulerError if the job can't be scheduled
     * @throws IllegalArgumentException If the preconditions are not met
     */
    public void schedule(final Object job, final ScheduleOptions options) throws IllegalArgumentException, SchedulerError {
        this.checkJob(job);

        if ( !(options instanceof InternalScheduleOptions)) {
            throw new IllegalArgumentException("Options has not been created via schedule or is null.");
        }
        final InternalScheduleOptions opts = (InternalScheduleOptions)options;

        // as this method might be called from unbind and during
        // unbind a deactivate could happen, we check the scheduler first
        final org.quartz.Scheduler s = this.scheduler;
        if ( s == null ) {
            throw new IllegalStateException("Scheduler is not available anymore.");
        }

        final String name;
        if ( opts.name != null ) {
            // if there is already a job with the name, remove it first
            try {
                final JobKey key = JobKey.jobKey(opts.name);
                final JobDetail jobdetail = s.getJobDetail(key);
                if (jobdetail != null) {
                    s.deleteJob(key);
                    this.logger.debug("Unscheduling job with name {}", opts.name);
                }
            } catch (final SchedulerException ignored) {
                // ignore
            }
            name = opts.name;
        } else {
            name = job.getClass().getName() + ':' + UUID.randomUUID();
            opts.name = name;
        }

        final Trigger trigger = opts.compile().withIdentity(name).build();

        // create the data map
        final JobDataMap jobDataMap = this.initDataMap(name, job, opts);

        final JobDetail detail = this.createJobDetail(name, jobDataMap, opts.canRunConcurrently);

        this.logger.debug("Scheduling job {} with name {} and trigger {}", job, name, trigger);
        try {
            s.scheduleJob(detail, trigger);
        } catch (SchedulerException ex) {
            throw new SchedulerError(ex);
        }
    }

    @Override
    public void reschedule(String jobName, ScheduleOptions options) throws SchedulerError {
        final org.quartz.Scheduler s = this.scheduler;
        if (jobName == null) {
            throw new IllegalArgumentException("Job name is mandatory");
        }
        JobKey key = JobKey.jobKey(jobName);
        if (key == null) {
            throw new IllegalStateException("No job found with name " + jobName);
        }
        try {
            JobDetail detail = s.getJobDetail(key);

            final String contextKey = key.toString();
            JobDataMap karafContext = ((KarafStdScheduler) s).getStorage().get(contextKey);
            Object job = karafContext.get(QuartzScheduler.DATA_MAP_OBJECT);

            s.deleteJob(key);

            final InternalScheduleOptions opts = (InternalScheduleOptions)options;
            Trigger trigger = opts.compile().withIdentity(jobName).build();
            JobDataMap jobDataMap = this.initDataMap(jobName, job, opts);
            detail = createJobDetail(jobName, jobDataMap, opts.canRunConcurrently);

            logger.debug("Update job scheduling {} with name {} and trigger {}", job, jobName, trigger);
            s.scheduleJob(detail, trigger);
        } catch (SchedulerException e) {
            throw new SchedulerError(e);
        }
    }

    /**
     * @see org.apache.karaf.scheduler.Scheduler#unschedule(java.lang.String)
     */
    public boolean unschedule(final String jobName) {
        final org.quartz.Scheduler s = this.scheduler;
        if (jobName != null && s != null) {
            try {
                final JobKey key = JobKey.jobKey(jobName);
                final JobDetail jobdetail = s.getJobDetail(key);
                if (jobdetail != null) {
                    s.deleteJob(key);
                    this.logger.debug("Unscheduling job with name {}", jobName);
                    return true;
                }
            } catch (final SchedulerException ignored) {
                // ignore
            }
        }
        return false;
    }

    @Override
    public Map<String, ScheduleOptions> getJobs() throws SchedulerError {
        try {
            Map<String, ScheduleOptions> jobs = new HashMap<>();
            org.quartz.Scheduler s = this.scheduler;
            if (s != null) {
                for (String group : s.getJobGroupNames()) {
                    for (JobKey key : s.getJobKeys(GroupMatcher.jobGroupEquals(group))) {
                        JobDetail detail = s.getJobDetail(key);
                        ScheduleOptions options = (ScheduleOptions) detail.getJobDataMap().get(DATA_MAP_OPTIONS);
                        jobs.put(key.getName(), options);
                    }
                }
            }
            return jobs;
        } catch (SchedulerException ex) {
            throw new SchedulerError(ex);
        }
    }

    @Override
    public boolean trigger(String jobName) throws SchedulerError {
        final org.quartz.Scheduler s = this.scheduler;
        if (jobName != null && s != null) {
            try {
                final JobKey key = JobKey.jobKey(jobName);
                final JobDetail jobdetail = s.getJobDetail(key);
                if (jobdetail != null) {
                    this.scheduler.triggerJob(key, jobdetail.getJobDataMap());
                    return true;
                }
            } catch (SchedulerException ex) {
                throw new SchedulerError(ex);
            }
        }
        return false;
    }

}
