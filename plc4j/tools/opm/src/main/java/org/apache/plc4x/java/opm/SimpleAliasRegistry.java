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
package org.apache.plc4x.java.opm;

import java.util.HashMap;
import java.util.Map;
import java.util.NoSuchElementException;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Simple Map Based Implementation of {@link AliasRegistry}.
 * It is not connection specific and forwards connection aware methods to the "simple" methods.
 */
public class SimpleAliasRegistry implements AliasRegistry {

    /**
     * Map from alias -> plc field address
     */
    private final Map<String, String> aliasMap;

    public SimpleAliasRegistry() {
        this(new ConcurrentHashMap<>());
    }

    public SimpleAliasRegistry(Map<String, String> aliasMap) {
        this.aliasMap = aliasMap;
    }

    /**
     * Register an Alias in the Registry.
     */
    public void register(String alias, String address) {
        this.aliasMap.put(alias, address);
    }

    @Override
    public boolean canResolve(String connection, String alias) {
        return canResolve(alias);
    }

    @Override
    public String resolve(String connection, String alias) {
        return resolve(alias);
    }

    @Override
    public boolean canResolve(String alias) {
        return aliasMap.containsKey(alias);
    }

    @Override
    public String resolve(String alias) {
        if (!canResolve(alias)) {
            throw new NoSuchElementException("Unable to resolve '" + alias + "'");
        }
        return aliasMap.get(alias);
    }
}
