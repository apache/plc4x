/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package org.apache.plc4x.merlot.das.api;

import java.util.List;
import java.util.function.Consumer;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;

public interface Device extends org.osgi.service.device.Device {
	
	public void init() throws Exception;
	
	public void destroy() throws Exception;
	
	public void start();
	
	public void stop();
	
        public void attach(Driver driver);
        
	public DriverEvent getEvent();
	
	public void putEvent(DriverEvent event);
        
        public void doEvent(DriverEvent event, String tag);
        
        public void ReadRequest(String index, String id, DriverCallback cb);
        
        public void WriteRequest(String scalar, String id, List<String> values, DriverCallback cb);
        
        public void SubscriptionRequest(String... events);
        
        public void UnsubscriptionRequest(String... events);
        
        public void ConsumerRegister(String event, Consumer<PlcSubscriptionEvent> consumer);
        
        public void ConsumerUnRegister(String event);
	
}
