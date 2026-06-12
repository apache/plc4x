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
package org.apache.plc4x.java.spi.buffers.bytebased.encoding;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.ServiceLoader;

public class EncodingManager {

    /** Lazily-initialized default instance to avoid repeated ServiceLoader scans. */
    private static volatile EncodingManager defaultInstance;

    /**
     * Returns a shared default instance, creating it on first access.
     * This avoids expensive ServiceLoader classpath scans on every buffer creation.
     */
    public static EncodingManager getDefault() {
        if (defaultInstance == null) {
            synchronized (EncodingManager.class) {
                if (defaultInstance == null) {
                    defaultInstance = new EncodingManager();
                }
            }
        }
        return defaultInstance;
    }

    private final Map<String, Encoding> encodingMap;

    public EncodingManager() {
        this(Thread.currentThread().getContextClassLoader());
    }

    public EncodingManager(ClassLoader classLoader) {
        this.encodingMap = new HashMap<>();
        ServiceLoader<Encoding> encodings = ServiceLoader.load(Encoding.class, classLoader);
        for (Encoding encoding : encodings) {
            String name = encoding.getName();
            Encoding existing = encodingMap.get(name);
            if (existing != null && existing != encoding) {
                throw new IllegalStateException("Conflicting encoding registrations for name '"
                    + name + "': "
                    + existing.getClass().getName() + " vs " + encoding.getClass().getName());
            }
            encodingMap.put(name, encoding);
        }
    }

    public Optional<Encoding> getEncoding(String encodingName) {
        Encoding encoding = encodingMap.get(encodingName);
        return Optional.ofNullable(encoding);
    }

}
