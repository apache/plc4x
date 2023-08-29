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
package org.apache.plc4x.merlot.api;

import org.apache.plc4x.merlot.api.core.Merlot;
import org.apache.plc4x.merlot.api.core.Merlot.STATE;
import java.time.LocalDateTime;


public interface Driver extends org.osgi.service.device.Driver {
	
    public void init() throws Exception;

    public void destroy() throws Exception;

    
    public void reset();  
    public void start();   
    public void stop();
    public void abort();
    public void clear();    
    public void hold();  
    public void unhold();
    public void suspend();
    public void unsuspend();
    
    public void onStopped();
    public void onResetting();
    public void onIdle();
    public void onStarting();
    public void onExecute();
    public void onCompleting();
    public void onComplete();
    public void onHolding();
    public void onHeld();
    public void onUnholding();
    public void onSuspending();
    public void onSuspended();
    public void onUnsuspending();
    public void onAborting();
    public void onAborted();
    public void onClearing();
    public void onStopping();
    
    
    
    
    public STATE getState();
        
    public void setUrl(String url);
    
    public String getUrl();   
    
    public LocalDateTime getStartDateTime();
    
    public LocalDateTime getLastFailDateTime();
    
    public long getRunningTime();
    
    public long getStoppedTime();    
    
    public long getSendMessages();
    
    public long getFailedMessages();          
        
    public long getDelay();
    
    public long getMaxDelay();
    
    public long getMinDelay();
        
    public void putEvent(DriverEvent event);   

    public DriverEvent getEvent();

    public DriverEvent EventFactory(Merlot m);
    
    public long getRequestQueueSize();
    
    public long getRequestQueueItems();   

    public long getResponseQueueSize();
    
    public long getResponseQueueItems();      
        
      

}
