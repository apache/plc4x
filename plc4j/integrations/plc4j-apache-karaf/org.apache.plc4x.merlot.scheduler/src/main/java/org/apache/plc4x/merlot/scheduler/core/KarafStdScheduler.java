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

import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import org.quartz.JobDataMap;
import org.quartz.JobDetail;
import org.quartz.JobKey;
import org.quartz.SchedulerException;
import org.quartz.Trigger;
import org.quartz.TriggerKey;
import org.quartz.impl.StdScheduler;


public class KarafStdScheduler extends StdScheduler {

    private final QuartzSchedulerStorage storage;

    public KarafStdScheduler(final org.quartz.core.QuartzScheduler scheduler) {
        super(scheduler);
        this.storage = new QuartzSchedulerStorage();
    }

    QuartzSchedulerStorage getStorage() {
        return this.storage;
    }

    @Override
    public Date scheduleJob(final JobDetail jobDetail, final Trigger trigger) throws SchedulerException {
        JobDataMap context = (JobDataMap) jobDetail.getJobDataMap().get(QuartzScheduler.DATA_MAP_CONTEXT);
        storage.put(jobDetail.getKey().toString(), context);

        jobDetail.getJobDataMap().remove(QuartzScheduler.DATA_MAP_CONTEXT);

        final Date date = super.scheduleJob(jobDetail, trigger);
        return date;
    }

    @Override
    public boolean deleteJob(JobKey jobKey) throws SchedulerException {
        final String contextKey = jobKey.toString();
        if (contextKey != null) {
            storage.release(contextKey);
        }
        return super.deleteJob(jobKey);
    }

    @Override
    public boolean deleteJobs(List<JobKey> jobKeys) throws SchedulerException {
        if (jobKeys != null) {
            final List<String> contextKeys = new ArrayList<>();
            for (JobKey jobKey : jobKeys) {
                contextKeys.add(jobKey.toString());
            }
            for (String contextKey : contextKeys) {
                storage.release(contextKey);
            }
        }
        return super.deleteJobs(jobKeys);
    }

    @Override
    public boolean unscheduleJob(TriggerKey triggerKey) throws SchedulerException {
        final Trigger trigger = getTrigger(triggerKey);
        final String contextKey = (trigger.getJobKey() != null) ? trigger.getJobKey().toString() : null;
        if (contextKey != null) {
            storage.release(contextKey);
        }
        return super.unscheduleJob(triggerKey);
    }

    @Override
    public boolean unscheduleJobs(List<TriggerKey> triggerKeys) throws SchedulerException {
        if (triggerKeys != null) {
            final List<String> contextKeys = new ArrayList<>();
            for (TriggerKey triggerKey : triggerKeys) {
                final Trigger trigger = getTrigger(triggerKey);
                final String contextKey = trigger.getJobKey().toString();
                if (contextKey != null) {
                    contextKeys.add(contextKey);
                }
            }
            for (String contextKey : contextKeys) {
                storage.release(contextKey);
            }
        }
        return super.unscheduleJobs(triggerKeys);
    }

}
