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

import org.apache.plc4x.merlot.scheduler.api.ScheduleOptions;
import org.apache.plc4x.merlot.scheduler.api.Scheduler;
import org.apache.plc4x.merlot.scheduler.api.SchedulerMBean;
import org.apache.plc4x.merlot.scheduler.command.support.TriggerJob;
import java.util.Map;
import javax.management.MBeanException;
import javax.management.NotCompliantMBeanException;
import javax.management.StandardMBean;
import javax.management.openmbean.*;


public class SchedulerMBeanImpl extends StandardMBean implements SchedulerMBean {

    private Scheduler scheduler;

    public SchedulerMBeanImpl() throws NotCompliantMBeanException {
        super(SchedulerMBean.class);
    }

    @Override
    public TabularData getJobs() throws MBeanException {
        try {
            CompositeType jobType = new CompositeType("Job", "Scheduler job",
                    new String[]{ "Job", "Schedule" },
                    new String[]{ "Job Name", "Job Scheduling" },
                    new OpenType[]{ SimpleType.STRING, SimpleType.STRING });
            TabularType tableType = new TabularType("Jobs", "Tables of all jobs", jobType, new String[]{ "Job" });
            TabularData table = new TabularDataSupport(tableType);

            Map<String, ScheduleOptions> jobs = scheduler.getJobs();
            for (Map.Entry<String, ScheduleOptions> entry : jobs.entrySet()) {
                CompositeData data = new CompositeDataSupport(jobType,
                        new String[]{ "Job", "Schedule" },
                        new Object[]{ entry.getKey(), entry.getValue().schedule()});
                table.put(data);
            }
            return table;
        } catch (Exception e) {
            throw new MBeanException(null, e.toString());
        }
    }

    @Override
    public void trigger(String name, boolean background) throws MBeanException {
        try {
            if (background) {
                scheduler.schedule(new TriggerJob(scheduler, name), scheduler.NOW());
            } else {
                scheduler.trigger(name);
            }
        } catch (Exception e) {
            throw new MBeanException(null, e.toString());
        }
    }

    @Override
    public void unschedule(String name) throws MBeanException {
        try {
            scheduler.unschedule(name);
        } catch (Exception e) {
            throw new MBeanException(null, e.toString());
        }
    }

    public Scheduler getScheduler() {
        return scheduler;
    }

    public void setScheduler(Scheduler scheduler) {
        this.scheduler = scheduler;
    }

}
