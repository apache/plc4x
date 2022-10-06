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
package org.apache.plc4x.test;

import org.apache.plc4x.test.model.LocationAware;

import java.io.IOException;
import java.io.InputStream;
import java.net.MalformedURLException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;
import java.util.Objects;

/**
 * Parent class for test suites utilizing XML as representation mechanism.
 * <p>
 * It adds handling of resource loading and helps in navigating test framework to suite source files.
 */
public abstract class XmlTestsuiteLoader {

    protected final String testsuiteDocument;
    protected final InputStream testsuiteDocumentXml;
    protected final URI suiteUri;

    protected XmlTestsuiteLoader(String testsuiteDocument) {
        this.testsuiteDocument = testsuiteDocument;
        this.testsuiteDocumentXml = getClass().getResourceAsStream(testsuiteDocument);

        if (testsuiteDocumentXml == null) {
            throw new IllegalArgumentException("Suite " + testsuiteDocument + " not found");
        }

        try {
            URL resource = getClass().getResource(testsuiteDocument);
            Objects.requireNonNull(resource, testsuiteDocument + " not found in classpath");
            if ("file".equals(resource.getProtocol())) { // we run in IDE so lets swap "target" directory to source!
                String sourceLocation = resource.toString()
                    .replace("/target/classes/", "/src/main/resources/")
                    .replace("/target/test-classes/", "/src/test/resources/");
                resource = new URL(sourceLocation);
            }
            this.suiteUri = resource.toURI();
        } catch (URISyntaxException | MalformedURLException e) {
            throw new IllegalArgumentException("Failed to load suite resource " + testsuiteDocument, e);
        }
    }

    protected final URI getSourceUri(Object model) {
        if (model instanceof LocationAware) {
            return ((LocationAware) model).getLocation()
                .map(location -> URI.create(suiteUri + "?line=" + location.getLine() + "&column=" + location.getColumn()))
                .orElse(suiteUri);
        }
        return suiteUri;
    }
}
