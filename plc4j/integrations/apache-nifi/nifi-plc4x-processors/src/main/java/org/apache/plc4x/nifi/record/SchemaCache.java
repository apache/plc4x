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
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReferenceArray;

import org.apache.nifi.serialization.record.RecordSchema;
import org.apache.plc4x.java.api.model.PlcTag;

public class SchemaCache {
    private ConcurrentMap<String, SchemaContainer> schemaMap = new ConcurrentHashMap<>();
    private AtomicReferenceArray<String> schemaAppendOrder = new AtomicReferenceArray<>(0);
    private final AtomicInteger nextSchemaPosition = new AtomicInteger(0);
    private final AtomicInteger cacheSize = new AtomicInteger(0);

    /** Creates a schema cache with first-in-first-out replacement policy. Stores PlcTags and RecordSchema used for PlcResponse serialization
     * @param cacheSize initial cache size
     */
    public SchemaCache(int cacheSize) {
        this.cacheSize.set(cacheSize);
    }

    /** Empties and restart the cache with the given size
     * @param cacheSize size of schema cache
     */
    public void restartCache(int cacheSize) {
        this.cacheSize.set(cacheSize);
        this.schemaAppendOrder = new AtomicReferenceArray<>(cacheSize);
        this.schemaMap = new ConcurrentHashMap<>();
        this.nextSchemaPosition.set(0);
    }


    /** Adds the schema to the cache if not present. When the cache is full first-in-first-out replacement policy applies
     * @param schemaIdentifier tagName-address map used to store the schema 
     * @param tagsNames list of tag names
     * @param tagsList list of PlcTag's
     * @param schema record schema used for PlcResponse serialization. Can be null
     */
    public void addSchema(final Map<String,String> schemaIdentifier, final Set<String> tagsNames, final List<? extends PlcTag> tagsList,  final RecordSchema schema) {        
        if (!schemaMap.containsKey(schemaIdentifier.toString())){
            if (nextSchemaPosition.get() == cacheSize.get()){
                nextSchemaPosition.set(0);
            }
            removeSchema(schemaAppendOrder.get(nextSchemaPosition.get()));

            Map<String, PlcTag> tags = new HashMap<>();
            for (int i=0; i<tagsNames.size(); i++){
                tags.put(tagsNames.toArray(new String[]{})[i], tagsList.get(i));
            }
            schemaMap.put(schemaIdentifier.toString(), new SchemaContainer(tags, schema));
            schemaAppendOrder.set(nextSchemaPosition.get(), schemaIdentifier.toString());
            nextSchemaPosition.getAndAdd(1);
        }    
    }

    /** Removes the schema from the cache
     * @param schemaIdentifier tagName-address map used to store the schema 
     */
    public void removeSchema(final String schemaIdentifier) {
        if (schemaIdentifier == null)
            return;
        if (schemaMap.containsKey(schemaIdentifier)){
            schemaMap.remove(schemaIdentifier);
        }
    }


    /** Retrieves a schema from the cache if found
     * @param schemaIdentifier tagName-address map used to store the schema 
     * @return RecordSchema used for PlcResponse serialization. Null if not found
     */
    public RecordSchema retrieveSchema(final Map<String,String> schemaIdentifier) { 
        if (schemaMap.containsKey(schemaIdentifier.toString())){
            return schemaMap.get(schemaIdentifier.toString()).getSchema();
        }
        return null;
    }

    /** Retrieves tags from the cache if found
     * @param schemaIdentifier tagName-address map used to store the schema 
     * @return Map between tag names and the corresponding PlcTag. Null if not found
     */
    public Map<String, PlcTag> retrieveTags(final Map<String,String> schemaIdentifier) { 
        if (schemaMap.containsKey(schemaIdentifier.toString())){
            return schemaMap.get(schemaIdentifier.toString()).getTags();
        }
        return null;
    }

    protected int getNextSchemaPosition() {
        return nextSchemaPosition.get();
    }

    protected int getCacheSize() {
        return cacheSize.get();
    }

    static public class SchemaContainer {
        private final RecordSchema schema;
        private final Map<String, PlcTag> tags;

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
