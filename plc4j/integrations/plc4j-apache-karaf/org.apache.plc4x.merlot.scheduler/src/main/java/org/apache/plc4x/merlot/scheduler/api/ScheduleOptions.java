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

import java.io.Serializable;
import java.util.Map;

/**
 * Scheduler options provide an extensible way of defining how to schedule a job.
 * An option can be created via the scheduler.
 *
 * @since 2.3
 */
public interface ScheduleOptions extends Serializable {

    /**
     * Add optional configuration for the job.
     *
     * @param config An optional configuration object - this configuration is only passed to the job the job implements {@link Job}.
     * @return The {@code ScheduleOptions}.
     */
    ScheduleOptions config(final Map<String, Serializable> config);

    /**
     * Sets the name of the job.
     * A job only needs a name if it is scheduled and should be cancelled later on. The name can then be used to cancel the job.
     * If a second job with the same name is started, the second one replaces the first one.
     *
     * @param name The job name.
     * @return The {@code ScheduleOptions}.
     */
    ScheduleOptions name(final String name);

    /**
     * Flag indicating whether the job can be run concurrently.
     * This defaults to false.
     *
     * @param flag Whether this job can run even if previous scheduled runs are still running.
     * @return The {@code ScheduleOptions}.
     */
    ScheduleOptions canRunConcurrently(final boolean flag);

    String name();

    boolean canRunConcurrently();

    String schedule();

}