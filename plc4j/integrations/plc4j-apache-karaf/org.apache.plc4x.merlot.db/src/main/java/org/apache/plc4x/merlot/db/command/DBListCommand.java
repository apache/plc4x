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
package org.apache.plc4x.merlot.db.command;

import java.util.Arrays;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.epics.pvdatabase.PVDatabase;
import org.osgi.framework.BundleContext;


@Command(scope = "db", name = "list", description = "List command for simple database.")
@Service
public class DBListCommand implements Action {
    @Reference
    BundleContext bundleContext;
    
    @Reference
    PVDatabase master;
    
    @Override
    public Object execute() throws Exception {
        String[] names = master.getRecordNames();
        System.out.println(Arrays.toString(names));      
        return null;
    };
    
}
