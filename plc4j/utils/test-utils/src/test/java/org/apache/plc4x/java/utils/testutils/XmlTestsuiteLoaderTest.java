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
package org.apache.plc4x.java.utils.testutils;

import org.apache.plc4x.java.utils.testutils.utils.model.Location;
import org.apache.plc4x.java.utils.testutils.utils.model.LocationAware;
import org.junit.jupiter.api.Test;

import java.net.URI;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;

class XmlTestsuiteLoaderTest {

    static class TestLoader extends XmlTestsuiteLoader {
        public TestLoader() {
            super("/test-suite.xml");
        }

        public URI testGetSourceUri(Object model) {
            return getSourceUri(model);
        }
    }

    static class TestLocationAware implements LocationAware {
        private final Location location;

        public TestLocationAware(Location location) {
            this.location = location;
        }

        @Override
        public Optional<Location> getLocation() {
            return Optional.ofNullable(location);
        }
    }

    @Test
    void testLoaderCreation() {
        TestLoader loader = new TestLoader();
        assertNotNull(loader);
        assertNotNull(loader.testsuiteDocument);
        assertNotNull(loader.testsuiteDocumentXml);
        assertNotNull(loader.suiteUri);
    }

    @Test
    void testGetSourceUriWithLocationAware() {
        TestLoader loader = new TestLoader();
        Location location = new Location(10, 5);
        TestLocationAware model = new TestLocationAware(location);

        URI uri = loader.testGetSourceUri(model);

        assertNotNull(uri);
        assertTrue(uri.toString().contains("line=10"));
        assertTrue(uri.toString().contains("column=5"));
    }

    @Test
    void testGetSourceUriWithoutLocationAware() {
        TestLoader loader = new TestLoader();
        Object model = new Object();

        URI uri = loader.testGetSourceUri(model);

        assertNotNull(uri);
        assertFalse(uri.toString().contains("line="));
    }

    @Test
    void testGetSourceUriWithLocationAwareButNoLocation() {
        TestLoader loader = new TestLoader();
        TestLocationAware model = new TestLocationAware(null);

        URI uri = loader.testGetSourceUri(model);

        assertNotNull(uri);
        assertFalse(uri.toString().contains("line="));
    }

    @Test
    void testConstructorWithNonExistentFile() {
        assertThrows(IllegalArgumentException.class, () -> {
            new XmlTestsuiteLoader("/nonexistent-file.xml") {};
        });
    }
}
