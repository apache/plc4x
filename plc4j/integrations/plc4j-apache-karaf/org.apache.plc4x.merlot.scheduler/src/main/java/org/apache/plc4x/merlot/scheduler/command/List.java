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
import java.util.Map;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.apache.karaf.shell.support.table.ShellTable;


@Command(scope = "rt", name = "list", description = "List scheduled jobs")
@Service
public class List implements Action {

    @Reference
    Scheduler scheduler;

    @Override
    public Object execute() throws Exception {
        ShellTable table = new ShellTable();
        table.column("Name");
        table.column("Schedule");
        Map<String, ScheduleOptions> jobs = scheduler.getJobs();
        for (Map.Entry<String, ScheduleOptions> entry : jobs.entrySet()) {
            table.addRow().addContent(entry.getKey(), entry.getValue().schedule());
        }
        table.print(System.out);
        return null;
    }
}
