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

import java.util.Properties;
import org.quartz.Scheduler;
import org.quartz.SchedulerException;
import org.quartz.core.QuartzSchedulerResources;
import org.quartz.impl.StdSchedulerFactory;


public class KarafStdSchedulerFactory extends StdSchedulerFactory {

    public KarafStdSchedulerFactory() {
        throw new IllegalStateException("Not supported. Use: org.apache.karaf.scheduler.core.KarafStdSchedulerFactory(java.util.Properties)");
    }

    public KarafStdSchedulerFactory(final Properties properties) throws SchedulerException {
        super(properties);
    }

    public KarafStdSchedulerFactory(final String fileName) throws SchedulerException {
        throw new IllegalStateException("Not supported. Use: org.apache.karaf.scheduler.core.KarafStdSchedulerFactory(java.util.Properties)");
    }

    public Scheduler instantiate(final QuartzSchedulerResources resources, final org.quartz.core.QuartzScheduler quartzScheduler) {
        final Scheduler scheduler = new KarafStdScheduler(quartzScheduler);
        return scheduler;
    }

}
