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

import java.util.HashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReferenceArray;

import org.apache.nifi.serialization.record.RecordSchema;
import org.apache.plc4x.java.api.model.PlcTag;

public class SchemaCache {
    private ConcurrentMap<String, SchemaContainer> schemaMap = new ConcurrentHashMap<>();
    private AtomicReferenceArray<String> schemaAppendOrder = new AtomicReferenceArray<>(0);
    private final AtomicInteger lastSchemaPosition = new AtomicInteger(0);
    private final AtomicInteger cacheSize = new AtomicInteger(0);

    public SchemaCache(int cacheSize) {
        this.cacheSize.set(cacheSize);
    }

    protected int getLastSchemaPosition() {
        return lastSchemaPosition.get();
    }

    protected int getCacheSize() {
        return cacheSize.get();
    }

    public void setCacheSize(int cacheSize) {
        this.cacheSize.set(cacheSize);
        this.schemaAppendOrder = new AtomicReferenceArray<>(cacheSize);
        this.schemaMap = new ConcurrentHashMap<>();
        this.lastSchemaPosition.set(0);
    }

    public void addSchema(final Map<String,String> schemaIdentifier, final LinkedHashSet<String> tagsNames, final List<PlcTag> tagsList,  final RecordSchema schema) {
        if (!schemaMap.containsKey(schemaIdentifier.toString())){
            if (lastSchemaPosition.get() == cacheSize.get()){
                lastSchemaPosition.set(0);
            }
            removeSchema(schemaAppendOrder.get(lastSchemaPosition.get()));

            Map<String, PlcTag> tags = new HashMap<>();
            for (int i=0; i<tagsNames.size(); i++){
                tags.put(tagsNames.toArray(new String[]{})[i], tagsList.get(i));
            }
            schemaMap.put(schemaIdentifier.toString(), new SchemaContainer(tags, schema));
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
            return schemaMap.get(schemaIdentifier.toString()).getSchema();
        }
        return null;
    }

    public Map<String, PlcTag> retrieveTags(final Map<String,String> schemaIdentifier) { 
        if (schemaMap.containsKey(schemaIdentifier.toString())){
            return schemaMap.get(schemaIdentifier.toString()).getTags();
        }
        return null;
    }

    public class SchemaContainer {
        private RecordSchema schema;
        private Map<String, PlcTag> tags;

        public Map<String, PlcTag> getTags() {
            return tags;
        }

        public RecordSchema getSchema() {
            return schema;
        }

        SchemaContainer(Map<String, PlcTag> tags, RecordSchema schema){
            this.tags = tags;
            this.schema = schema;
        }
    }
}
