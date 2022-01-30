/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.utils.connectionpool2;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;

import javax.management.ObjectName;
import java.lang.management.ManagementFactory;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Function;

/**
 * Pool that sits on top of the {@link CachedDriverManager}.
 * <p>
 * This class is thread safe!
 */
public class PooledDriverManager extends PlcDriverManager implements PooledDriverManagerMBean {

    private final Function<String, CachedDriverManager> factory;
    private final Map<String, CachedDriverManager> cachedManagers = new ConcurrentHashMap<>();

    public PooledDriverManager() {
        this(new PlcDriverManager());
    }

    public PooledDriverManager(PlcDriverManager driverManager) {
        this.factory = key -> {
            return new CachedDriverManager(key, () -> {
                try {
                    return driverManager.getConnection(key);
                } catch (PlcConnectionException e) {
                    throw new RuntimeException(e);
                }
            });
        };

        // Register as MBean
        try {
            ManagementFactory.getPlatformMBeanServer().registerMBean(this, new ObjectName("org.pragmaticindustries.cockpit.plc:name=connection-pool"));
        } catch (Exception e) {
            // Do nothing.
        }
    }

    PooledDriverManager(Function<String, CachedDriverManager> factory) {
        this.factory = factory;
    }

    @Override
    public PlcConnection getConnection(String url) throws PlcConnectionException {
        return cachedManagers.computeIfAbsent(url, this.factory).getConnection(url);
    }

    @Override
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new UnsupportedOperationException();
    }

    Map<String, CachedDriverManager> getCachedManagers() {
        return this.cachedManagers;
    }

    @Override
    public String[] getConnectedUrls() {
        return this.cachedManagers.keySet().toArray(new String[0]);
    }
}
