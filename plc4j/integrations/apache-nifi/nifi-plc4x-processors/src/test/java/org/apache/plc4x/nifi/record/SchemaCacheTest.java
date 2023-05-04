/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.plc4x.nifi.record;

import java.util.ArrayList;
import java.util.Collection;
import java.util.HashSet;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;

import org.apache.avro.Schema;
import org.apache.avro.SchemaBuilder;
import org.apache.avro.SchemaBuilder.FieldAssembler;
import org.apache.nifi.avro.AvroTypeUtil;
import org.apache.nifi.serialization.record.RecordSchema;
import org.apache.plc4x.java.api.model.PlcTag;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

public class SchemaCacheTest {

    private static final SchemaCache schemaCache = new SchemaCache(0);
    private static final List<RecordSchema> schemas = new ArrayList<>();
    private static final List<List<PlcTag>> tags = new ArrayList<>();
    private static final List<Map<String, String>> addresses = new ArrayList<>();
    private static final List<LinkedHashSet<String>> tagNames = new ArrayList<>();

    @BeforeAll
    static public void innit() {
        // Create 5 unique Schemas

        for (int i = 0; i < 5; i++) {
            FieldAssembler<Schema> fieldAssembler = SchemaBuilder.builder().record("foo").fields();
            Collection<String> collection = new HashSet<>();
            List<PlcTag> innerTags = new ArrayList<>();
            String address = null;

            switch (i) {
                case 0:
                    address = "RANDOM/fooBOOL:BOOL";
                    addresses.add(Map.of("BOOL", address));
                    innerTags.add(new TestPlcTag(address));
                    tags.add(innerTags);
                    fieldAssembler.nullableBoolean("BOOL", false);
                    collection.add("fooBOOL");
                    tagNames.add(new LinkedHashSet<>(collection));
                    break;
                case 1:
                    address = "RANDOM/fooINT:INT";
                    addresses.add(Map.of("INT", address));
                    innerTags.add(new TestPlcTag(address));
                    tags.add(innerTags);
                    fieldAssembler.nullableInt("INT", 1);
                    collection.add("fooINT");
                    tagNames.add(new LinkedHashSet<>(collection));
                    break;
                case 2:
                    address = "RANDOM/fooLONG:LONG";
                    addresses.add(Map.of("LONG", address));
                    innerTags.add(new TestPlcTag(address));
                    tags.add(innerTags);
                    fieldAssembler.nullableLong("LONG", i * 100);
                    collection.add("fooLONG");
                    tagNames.add(new LinkedHashSet<>(collection));
                    break;
                case 3:
                    address = "RANDOM/fooFLOAT:FLOAT";
                    addresses.add(Map.of("FLOAT", address));
                    innerTags.add(new TestPlcTag(address));
                    tags.add(innerTags);
                    fieldAssembler.nullableFloat("FLOAT", i * 0.1F);
                    collection.add("fooFLOAT");
                    tagNames.add(new LinkedHashSet<>(collection));
                    break;
                case 4:
                    address = "RANDOM/fooDOUBLE:DOUBLE";
                    addresses.add(Map.of("DOUBLE", address));
                    innerTags.add(new TestPlcTag(address));
                    tags.add(innerTags);
                    fieldAssembler.nullableDouble("DOUBLE", i * 0.01);
                    collection.add("fooDOUBLE");
                    tagNames.add(new LinkedHashSet<>(collection));
                    break;
            }
            Schema avroSchema = fieldAssembler.endRecord();
            schemas.add(AvroTypeUtil.createSchema(avroSchema));
        }
    }

    // Cache size set to 4 < number of schemas: to check schema override.
    @BeforeEach
    public void testCacheSize() {
        schemaCache.restartCache(4);
        assert schemaCache.getCacheSize() == 4;
        assert schemaCache.getNextSchemaPosition() == 0;
    }

    // In this test we add 4 schemas and try to add schema 0 again. It should not be added.
    @Test
    public void testAddSchema() {
        for (int i = 0; i < 4; i++) {
            schemaCache.addSchema(addresses.get(i), tagNames.get(i), tags.get(i), schemas.get(i));
            assert schemaCache.getNextSchemaPosition() == i + 1;
        }
        int prev = schemaCache.getNextSchemaPosition();
        schemaCache.addSchema(addresses.get(0), tagNames.get(0), tags.get(0), schemas.get(0));
        assert prev == schemaCache.getNextSchemaPosition();
    }

    // In this test check schema overriding
    @Test
    public void testSchemaOverride() {
        for (int i = 0; i < 4; i++) {
            schemaCache.addSchema(addresses.get(i), tagNames.get(i), tags.get(i), schemas.get(i));
            assert schemaCache.getNextSchemaPosition() == i + 1;
        }
        // Override first schema
        schemaCache.addSchema(addresses.get(4), tagNames.get(4), tags.get(4), schemas.get(4));
        assert schemaCache.getNextSchemaPosition() == 1;

        // First schema should not be present in the cache
        assert schemaCache.retrieveSchema(addresses.get(0)) == null;

        // Check remaining schemas
        for (int i=1; i<5; i++){
            assert schemaCache.retrieveSchema(addresses.get(i)) == schemas.get(i);
        }
    }


    public static void main(String[] args) {
        SchemaCacheTest.innit();

        SchemaCacheTest instance = new SchemaCacheTest();

        instance.testCacheSize();
        instance.testAddSchema();

        instance.testCacheSize();
        instance.testSchemaOverride();
    }


    private static class TestPlcTag implements PlcTag {

        final String address;

        TestPlcTag(String address) {
            this.address = address;
        }

        @Override
        public String getAddressString() {
            return address;
        }
    }
}
