/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.apache.plc4x.merlot.das.base.api;

import java.time.LocalDateTime;


public interface BaseDriverMBean {
    
    public void init() throws Exception;

    public void destroy() throws Exception;

    public void start();

    public void stop();	
        
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
       
    public long getRequestQueueSize();
    
    public long getRequestQueueItems();   

    public long getResponseQueueSize();
    
    public long getResponseQueueItems();      
    
}
