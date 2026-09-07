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
package org.apache.plc4x.java.spi.buffers.bytebased.byteorder;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.ServiceLoader;

public class ByteOrderManager {

    /** Lazily-initialized default instance to avoid repeated ServiceLoader scans. */
    private static volatile ByteOrderManager defaultInstance;

    /**
     * Returns a shared default instance, creating it on first access.
     * This avoids expensive ServiceLoader classpath scans on every buffer creation.
     */
    public static ByteOrderManager getDefault() {
        if (defaultInstance == null) {
            synchronized (ByteOrderManager.class) {
                if (defaultInstance == null) {
                    defaultInstance = new ByteOrderManager();
                }
            }
        }
        return defaultInstance;
    }

    private final Map<String, ByteOrder> byteOrderMap;

    public ByteOrderManager() {
        this(Thread.currentThread().getContextClassLoader());
    }

    public ByteOrderManager(ClassLoader classLoader) {
        this.byteOrderMap = new HashMap<>();
        ServiceLoader<ByteOrder> byteOrders = ServiceLoader.load(ByteOrder.class, classLoader);
        for (ByteOrder byteOrder : byteOrders) {
            if (byteOrderMap.containsKey(byteOrder.getName())) {
                throw new IllegalStateException(
                    "Multiple byte order implementations available for byte order name '" +
                        byteOrder.getName() + "'");
            }
            byteOrderMap.put(byteOrder.getName(), byteOrder);
        }
    }

    public Optional<ByteOrder> getByteOrder(String byteOrderName) {
        ByteOrder byteOrder = byteOrderMap.get(byteOrderName);
        return Optional.ofNullable(byteOrder);
    }

}
