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
package org.apache.plc4x.merlot.scheduler.command;

import org.apache.plc4x.merlot.scheduler.api.ScheduleOptions;
import org.apache.plc4x.merlot.scheduler.api.Scheduler;
import org.apache.plc4x.merlot.scheduler.command.support.ScriptJob;
import java.util.Date;
import javax.xml.bind.DatatypeConverter;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Argument;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.karaf.shell.api.console.Function;
import org.apache.karaf.shell.api.console.Session;
import org.apache.karaf.shell.api.console.SessionFactory;


@Command(scope = "rt", name = "schedule-script", description = "Schedule a script execution")
@Service
public class ScheduleScript implements Action {

    @Option(name = "--name", description = "Name of this job")
    String name;

    @Option(name = "--concurrent", description = "Should jobs run concurrently or not (defaults to false)")
    boolean concurrent;

    @Option(name = "--cron", description = "The cron expression")
    String cron;

    @Option(name = "--at", description = "Absolute date in ISO format (ex: 2014-05-13T13:56:45)")
    String at;

    @Option(name = "--times", description = "Number of times this job should be executed")
    int times = -1;

    @Option(name = "--period", description = "Time during executions (in seconds)")
    long period;

    @Argument(name = "script", required = true, description = "The script to schedule")
    Function script;

    @Reference
    Scheduler scheduler;

    @Reference
    Session session;

    @Reference
    SessionFactory sessionFactory;

    @Override
    public Object execute() throws Exception {
        if (cron != null && (at != null || times != -1 || period != 0)) {
            throw new IllegalArgumentException("Both cron expression and explicit execution time can not be specified");
        }
        ScheduleOptions options;
        if (cron != null) {
            options = scheduler.EXPR(cron);
        } else {
            Date date;
            if (at != null) {
                date = DatatypeConverter.parseDateTime(at).getTime();
            } else {
                date = new Date();
            }
            if (period > 0) {
                options = scheduler.AT(date, times, period);
            } else {
                options = scheduler.AT(date);
            }
        }
        if (name != null) {
            options.name(name);
        }
        if (concurrent) {
            options.canRunConcurrently(concurrent);
        }
        scheduler.schedule(new ScriptJob(sessionFactory, session, script), options);
        return null;
    }
}
