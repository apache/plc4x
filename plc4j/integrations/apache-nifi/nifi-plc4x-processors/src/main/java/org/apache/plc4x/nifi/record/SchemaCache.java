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
package org.apache.plc4x.nifi.record;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReferenceArray;

import org.apache.nifi.serialization.record.RecordSchema;

public class SchemaCache {
    private ConcurrentMap<String, RecordSchema> schemaMap = new ConcurrentHashMap<>();
    private AtomicReferenceArray<String> schemaAppendOrder = new AtomicReferenceArray<>(0);
    private final AtomicInteger lastSchemaPosition = new AtomicInteger(0);
    private final AtomicInteger cacheSize = new AtomicInteger(0);

    public SchemaCache(int cacheSize) {
        this.cacheSize.set(cacheSize);
    }

    public void setCacheSize(int cacheSize) {
        this.cacheSize.set(cacheSize);
        schemaAppendOrder = new AtomicReferenceArray<>(cacheSize);
        schemaMap = new ConcurrentHashMap<>();
    }

    public void addSchema(final Map<String,String> schemaIdentifier, final RecordSchema schema) {
        if (!schemaMap.containsKey(schemaIdentifier.toString())){
            if (lastSchemaPosition.get() == cacheSize.get()){
                lastSchemaPosition.set(0);
            }
            removeSchema(schemaAppendOrder.get(lastSchemaPosition.get()));

            schemaMap.put(schemaIdentifier.toString(), schema);
            schemaAppendOrder.set(lastSchemaPosition.get(), schemaIdentifier.toString());
            lastSchemaPosition.getAndAdd(1);
        }    
    }

    public void removeSchema(final String schemaIdentifier) {
        if (schemaIdentifier == null)
            return;
        if (schemaMap.containsKey(schemaIdentifier)){
            schemaMap.remove(schemaIdentifier);
        }
    }

    public RecordSchema retrieveSchema(final Map<String,String> schemaIdentifier) { 
        if (schemaMap.containsKey(schemaIdentifier.toString())){
            return schemaMap.get(schemaIdentifier.toString());
        }
        return null;
    }
}
