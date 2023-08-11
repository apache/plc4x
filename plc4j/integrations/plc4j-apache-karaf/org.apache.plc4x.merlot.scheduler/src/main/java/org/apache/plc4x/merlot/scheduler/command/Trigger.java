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

import org.apache.plc4x.merlot.scheduler.api.Scheduler;
import org.apache.plc4x.merlot.scheduler.command.support.TriggerJob;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Argument;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.Option;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;


@Command(scope = "rt", name = "trigger", description = "Manually trigger a scheduled job")
@Service
public class Trigger implements Action {

    @Argument(description = "Name of the job to trigger", required = true)
    String name;

    @Option(name = "-b", aliases = "background", description = "schedule the trigger in the background", required = false)
    boolean background = false;

    @Reference
    Scheduler scheduler;

    @Override
    public Object execute() throws Exception {
        if (background) {
            System.out.println("Scheduling background trigger for job " + name);
            scheduler.schedule(new TriggerJob(scheduler, name), scheduler.NOW());
        } else {
            System.out.println("Triggering job " + name);
            if (!scheduler.trigger(name)) {
                System.out.println("Could not find a scheduled job with name " + name);
            }
        }
        return null;
    }


}
