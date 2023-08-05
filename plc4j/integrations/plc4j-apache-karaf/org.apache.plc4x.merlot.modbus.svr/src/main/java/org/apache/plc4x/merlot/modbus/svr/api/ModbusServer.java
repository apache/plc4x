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
package org.apache.plc4x.merlot.modbus.svr.api;

import java.net.SocketAddress;
import org.osgi.framework.BundleContext;

public interface ModbusServer {

    public void init() throws Exception;

    public void destroy() throws Exception;

    public void start();

    public void stop();

    public void setPort(int port);

    public int getPort();

    public void setHost(String host);

    public String getHost();
    
    public boolean isStarted();
    
    public void setSocketAddress(SocketAddress[] saddress);
    
    public SocketAddress[] getSocketAddress();
    
    public String getDate();

    public String getElapseTime();
    
    public void setBundleContext(BundleContext bundleContext);   
    
    public long getCheckInterval();
    
    public void setCheckInterval(long newCheckInterval);
        
    public long getLastReadThroughput();
    
    public long getLastWriteThroughput();
    
    public long getCumulativeReadBytes();
    
    public long getCumulativeWrittenBytes();
    
    public long getCurrentReadBytes();
    
    public long getCurrentWrittenBytes();
    
}
