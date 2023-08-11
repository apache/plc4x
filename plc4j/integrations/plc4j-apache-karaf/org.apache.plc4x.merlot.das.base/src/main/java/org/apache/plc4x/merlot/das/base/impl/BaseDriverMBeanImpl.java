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
package org.apache.plc4x.merlot.das.base.impl;

import org.apache.plc4x.merlot.das.base.api.BaseDriver;
import org.apache.plc4x.merlot.das.base.api.BaseDriverMBean;
import java.time.LocalDateTime;
import javax.management.NotCompliantMBeanException;
import javax.management.StandardMBean;


public class BaseDriverMBeanImpl extends StandardMBean implements BaseDriverMBean {

    private final BaseDriver driver;
    
    public BaseDriverMBeanImpl(BaseDriver driver) throws NotCompliantMBeanException {
        super(BaseDriverMBean.class);
        this.driver = driver;
    };

    @Override
    public void init() throws Exception {
        driver.init();
    }

    @Override
    public void destroy() throws Exception {
        driver.destroy();
    }

    @Override
    public void start() {
        driver.start();
    }

    @Override
    public void stop() {
        driver.stop();
    }

    @Override
    public void setUrl(String url) {
        driver.setUrl(url);
    }

    @Override
    public String getUrl() {
        return driver.getUrl();
    }

    @Override
    public LocalDateTime getStartDateTime() {
        return driver.getStartDateTime();
    }

    @Override
    public LocalDateTime getLastFailDateTime() {
        return driver.getLastFailDateTime();
    }

    @Override
    public long getRunningTime() {
        return driver.getRunningTime();
    }

    @Override
    public long getStoppedTime() {
        return driver.getStoppedTime();
    }

    @Override
    public long getSendMessages() {
        return driver.getSendMessages();
    }

    @Override
    public long getFailedMessages() {
        return driver.getFailedMessages();
    }

    @Override
    public long getDelay() {
        return driver.getDelay();
    }

    @Override
    public long getMaxDelay() {
        return driver.getMaxDelay();
    }

    @Override
    public long getMinDelay() {
        return driver.getMinDelay();
    }  


    @Override
    public long getRequestQueueSize() {
        return driver.getRequestQueueSize();
    }

    @Override
    public long getRequestQueueItems() {
        return driver.getRequestQueueItems();
    }

    @Override
    public long getResponseQueueSize() {
        return driver.getResponseQueueSize();
    }

    @Override
    public long getResponseQueueItems() {
        return driver.getResponseQueueItems();
    }
    

    
    
}
