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
package org.apache.plc4x.merlot.das.base.command;

import org.apache.karaf.shell.api.action.Action;
import org.apache.karaf.shell.api.action.Command;
import org.apache.karaf.shell.api.action.lifecycle.Reference;
import org.apache.karaf.shell.api.action.lifecycle.Service;
import org.osgi.framework.BundleContext;
import org.apache.plc4x.merlot.das.base.api.BaseDriver;


@Command(scope = "basic", name = "test", description = "Command for test.")
@Service
public class BaseDriverCommand implements Action {

    @Reference
    BundleContext bundleContext;
    
    @Reference
    BaseDriver theDriver;
    
	@Override
	public Object execute() throws Exception {
		// TODO Auto-generated method stub
		System.out.println("Comando de prueba");
		if (theDriver == null) {
			System.out.println("Nose encontro el driver.");
		} {
			//theDriver.InitializeDevice(null);
		};
		return null;
	}

}
