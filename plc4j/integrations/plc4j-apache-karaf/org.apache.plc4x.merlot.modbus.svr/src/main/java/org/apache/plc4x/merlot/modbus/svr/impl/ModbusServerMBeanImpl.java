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
package org.apache.plc4x.merlot.modbus.svr.impl;

import org.apache.plc4x.merlot.modbus.svr.api.ModbusServer;
import org.apache.plc4x.merlot.modbus.svr.api.ModbusServerMBean;
import javax.management.NotCompliantMBeanException;
import javax.management.StandardMBean;


public class ModbusServerMBeanImpl extends StandardMBean implements ModbusServerMBean{

    private final ModbusServer mserver; 
    
    public ModbusServerMBeanImpl(ModbusServer mserver) throws NotCompliantMBeanException {
        super(ModbusServerMBean.class);
        this.mserver = mserver;
    }    
    
    @Override
    public void init() throws Exception {
        mserver.init();
    }

    @Override
    public void destroy() throws Exception {
        mserver.destroy();
    }

    @Override
    public void start() {
        mserver.start();
    }

    @Override
    public void stop() {
        mserver.stop();
    }

    @Override
    public void setPort(int port) {
        mserver.setPort(port);
    }

    @Override
    public int getPort() {
        return mserver.getPort();
    }

    @Override
    public void setHost(String host) {
        mserver.setHost(host);
    }

    @Override
    public String getHost() {
        return mserver.getHost();
    }

    @Override
    public boolean isStarted() {
        return mserver.isStarted();
    }

    @Override
    public long getCheckInterval() {
        return mserver.getCheckInterval();
    }

    @Override
    public void setCheckInterval(long newCheckInterval) {
        mserver.setCheckInterval(newCheckInterval);
    }

    @Override
    public long getLastReadThroughput() {
        return mserver.getLastReadThroughput();
    }

    @Override
    public long getLastWriteThroughput() {
        return mserver.getLastWriteThroughput();
    }

    @Override
    public long getCumulativeReadBytes() {
        return mserver.getCumulativeReadBytes();
    }

    @Override
    public long getCumulativeWrittenBytes() {
        return mserver.getCumulativeWrittenBytes();
    }

    @Override
    public long getCurrentReadBytes() {
        return mserver.getCurrentReadBytes();
    }

    @Override
    public long getCurrentWrittenBytes() {
        return mserver.getCurrentWrittenBytes();
    }
    
    
    
}
