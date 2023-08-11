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
package org.apache.plc4x.merlot.scheduler.core;

import org.apache.plc4x.merlot.scheduler.api.SchedulerStorage;
import java.io.Serializable;
import java.util.HashMap;
import java.util.Map;


public class QuartzSchedulerStorage implements SchedulerStorage {

    private final Map<Serializable, Object> store = new HashMap<>();

    @Override
    public <T> T get(Serializable key) {
        return (T) this.store.get(key);
    }

    @Override
    public void put(Serializable key, Object value) {
        this.store.put(key, value);
    }

    @Override
    public boolean contains(Serializable key) {
        return this.store.containsKey(key);
    }

    @Override
    public void release(Serializable key) {
        this.store.remove(key);
    }

}
