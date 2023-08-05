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

/**
 * A job is executed by the {@link Scheduler} service.
 * If the implementation of the job requires certain environment information
 * it can implement this interface to get additional information
 * through the provided {@link JobContext}.
 * If no additional information is required, implementing {@link Runnable} is
 * sufficient.
 */
public interface Job {

    /**
     * Execute this job.
     * @param context The context of the job.
     */
    void execute(JobContext context);
}
