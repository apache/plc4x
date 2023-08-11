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
package org.apache.plc4x.merlot.scheduler.command.support;

import org.apache.plc4x.merlot.scheduler.api.Scheduler;
import org.apache.plc4x.merlot.scheduler.api.SchedulerError;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class TriggerJob implements Runnable {

    private static final Logger LOGGER = LoggerFactory.getLogger(TriggerJob.class);

    private final Scheduler scheduler;
    private final String name;

    public TriggerJob(Scheduler scheduler, String name) {
        this.scheduler = scheduler;
        this.name = name;
    }

    @Override
    public void run() {
        try {
            if (!scheduler.trigger(name)) {
                LOGGER.warn("Could not find a scheduled job with name " + name);
            }
        } catch (SchedulerError ex) {
            LOGGER.error("Failed to trigger job {}", name, ex);
        }
    }

}
