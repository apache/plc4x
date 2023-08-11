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

import org.apache.plc4x.merlot.scheduler.api.Scheduler;
import java.util.Enumeration;
import java.util.Properties;
import org.apache.karaf.util.tracker.BaseActivator;
import org.apache.karaf.util.tracker.annotation.Managed;
import org.apache.karaf.util.tracker.annotation.ProvideService;
import org.apache.karaf.util.tracker.annotation.Services;
import org.osgi.service.cm.ManagedService;


@Services(provides = @ProvideService(Scheduler.class))
@Managed("org.apache.karaf.scheduler.quartz")
public class Activator extends BaseActivator implements ManagedService {

    private QuartzScheduler scheduler;
    private WhiteboardHandler whiteboardHandler;

    @Override
    protected void doStart() throws Exception {
        Properties properties = new Properties();
        if (getConfiguration() == null) {
            return;
        }
        Enumeration<String> keys = getConfiguration().keys();
        while (keys.hasMoreElements()) {
            String key = keys.nextElement();
            if (key.startsWith("org.quartz")) {
                Object value = getConfiguration().get(key);
                properties.put(key, value);
            }
        }
        scheduler = new QuartzScheduler(properties);
        register(Scheduler.class, scheduler);
        whiteboardHandler = new WhiteboardHandler(bundleContext, scheduler);

        SchedulerMBeanImpl mBean = new SchedulerMBeanImpl();
        mBean.setScheduler(scheduler);
        registerMBean(mBean, "type=rt");
    }

    @Override
    protected void doStop() {
        super.doStop();

        if (whiteboardHandler != null) {
            whiteboardHandler.deactivate();
            whiteboardHandler = null;
        }
        if (scheduler != null) {
            scheduler.deactivate();
            scheduler = null;
        }
        super.doStop();
    }

}
