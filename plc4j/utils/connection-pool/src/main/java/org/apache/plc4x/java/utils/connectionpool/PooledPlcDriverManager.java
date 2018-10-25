/*
 Licensed to the Apache Software Foundation (ASF) under one or more
 contributor license agreements.  See the NOTICE file distributed with
 this work for additional information regarding copyright ownership.
 The ASF licenses this file to You under the Apache License, Version 2.0
 (the "License"); you may not use this file except in compliance with
 the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
 */

package org.apache.plc4x.java.utils.connectionpool;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.commons.pool2.ObjectPool;
import org.apache.commons.pool2.impl.GenericObjectPool;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Proxy;
import java.util.HashMap;
import java.util.LinkedHashSet;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReadWriteLock;
import java.util.concurrent.locks.ReentrantReadWriteLock;

public class PooledPlcDriverManager extends PlcDriverManager {

    private final static Logger LOGGER = LoggerFactory.getLogger(PooledPlcDriverManager.class);

    private PoolCreator poolCreator;

    private ReadWriteLock readWriteLock = new ReentrantReadWriteLock();

    private ConcurrentMap<Pair<String, PlcAuthentication>, ObjectPool<PlcConnection>> poolMap = new ConcurrentHashMap<>();

    // Marker class do detected a non null value
    private static final NoPlcAuthentication noPlcAuthentication = new NoPlcAuthentication();

    public PooledPlcDriverManager() {
        this(GenericObjectPool::new);
    }

    public PooledPlcDriverManager(ClassLoader classLoader) {
        super(classLoader);
        this.poolCreator = GenericObjectPool::new;
    }

    public PooledPlcDriverManager(PoolCreator poolCreator) {
        this.poolCreator = poolCreator;
    }

    public PooledPlcDriverManager(ClassLoader classLoader, PoolCreator poolCreator) {
        super(classLoader);
        this.poolCreator = poolCreator;
    }

    @Override
    public PlcConnection getConnection(String url) throws PlcConnectionException {
        return getConnection(url, noPlcAuthentication);
    }

    @Override
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        Pair<String, PlcAuthentication> argPair = Pair.of(url, authentication);
        ObjectPool<PlcConnection> pool = retrieveFromPool(argPair);
        try {
            LOGGER.debug("Try to borrow an object for url {} and authentication {}", url, authentication);
            PlcConnection plcConnection = pool.borrowObject();
            // TODO 25-10-2018 jf: Return a real wrapper object. This implementation leaks the connection.
            // The connection can be reused by the pool but is still referenced and can thus still be used.
            return (PlcConnection) Proxy.newProxyInstance(classLoader, new Class[]{PlcConnection.class}, (o, method, objects) -> {
                if (method.getName().equals("close")) {
                    LOGGER.debug("close called on {}. Returning to {}", plcConnection, pool);
                    pool.returnObject(plcConnection);
                    return null;
                } else {
                    return method.invoke(plcConnection, objects);
                }
            });
        } catch (Exception e) {
            throw new PlcConnectionException(e);
        }
    }

    private ObjectPool<PlcConnection> retrieveFromPool(Pair<String, PlcAuthentication> argPair) {
        String url = argPair.getLeft();
        PlcAuthentication plcAuthentication = argPair.getRight();
        ObjectPool<PlcConnection> pool = poolMap.get(argPair);
        if (pool == null) {
            Lock lock = readWriteLock.writeLock();
            lock.lock();
            try {
                poolMap.computeIfAbsent(argPair, pair -> poolCreator.createPool(new PooledPlcConnectionFactory() {
                    @Override
                    public PlcConnection create() throws PlcConnectionException {
                        if (plcAuthentication == noPlcAuthentication) {
                            LOGGER.debug("getting actual connection for {}", url);
                            return PooledPlcDriverManager.super.getConnection(url);
                        } else {
                            LOGGER.debug("getting actual connection for {} and plcAuthentication {}", url, plcAuthentication);
                            return PooledPlcDriverManager.super.getConnection(url, plcAuthentication);
                        }
                    }
                }));
                pool = poolMap.get(argPair);
            } finally {
                lock.unlock();
            }
        }
        return pool;
    }

    @FunctionalInterface
    interface PoolCreator {
        ObjectPool<PlcConnection> createPool(PooledPlcConnectionFactory pooledPlcConnectionFactory);
    }

    // TODO: maybe add a Thread which calls this cyclic
    public void removedUnusedPools() {
        Lock lock = readWriteLock.writeLock();
        lock.lock();
        try {
            Set<Pair<String, PlcAuthentication>> itemsToBeremoved = new LinkedHashSet<>();
            poolMap.forEach((key, value) -> {
                // TODO: check if this pool has been used in the last time and if not remove it.
                // TODO: evicting empty pools for now
                if (value.getNumActive() == 0 && value.getNumIdle() == 0) {
                    LOGGER.info("Removing unused pool {}", value);
                    itemsToBeremoved.add(key);
                }
            });
            itemsToBeremoved.forEach(poolMap::remove);
        } finally {
            lock.unlock();
        }
    }

    // TODO: maybe export to jmx
    public Map<String, Number> getStatistics() {
        Lock lock = readWriteLock.readLock();
        try {
            lock.lock();
            HashMap<String, Number> statistics = new HashMap<>();
            for (Map.Entry<Pair<String, PlcAuthentication>, ObjectPool<PlcConnection>> poolEntry : poolMap.entrySet()) {
                Pair<String, PlcAuthentication> pair = poolEntry.getKey();
                ObjectPool<PlcConnection> objectPool = poolEntry.getValue();
                String url = pair.getLeft();
                PlcAuthentication plcAuthentication = pair.getRight();

                String authSuffix = plcAuthentication != noPlcAuthentication ? "/" + plcAuthentication : "";
                statistics.put(url + authSuffix + ".numActive", objectPool.getNumActive());
                statistics.put(url + authSuffix + ".numIdle", objectPool.getNumIdle());
            }

            return statistics;
        } finally {
            lock.unlock();
        }
    }

    private static final class NoPlcAuthentication implements PlcAuthentication {

    }
}
