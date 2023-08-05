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
package org.apache.plc4x.merlot.uns.command;

import org.apache.plc4x.merlot.uns.api.Model;
import java.util.UUID;
import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Argument;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.osgi.framework.BundleContext;


@Command(scope = "model", name = "del", description = "Command for test tree.")
@Service
public class ModelDelCommand implements Action {

    @Reference
    BundleContext bundleContext;
    
    @Reference
    Model model;
    
    @Argument(index = 0, name = "uuid", description = "Node uuid in the model.", required = true, multiValued = false)
    String struuid;     
    
    @Override
    public Object execute() throws Exception {
        UUID uuid = UUID.fromString(struuid);
        model.delUnit(uuid);
        return null;
    }
    
}
