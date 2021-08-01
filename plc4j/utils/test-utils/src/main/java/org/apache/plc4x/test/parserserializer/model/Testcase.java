/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.test.parserserializer.model;

import java.util.Map;
import org.apache.plc4x.test.model.Location;
import org.apache.plc4x.test.model.LocationAware;
import org.dom4j.Element;

import java.util.List;
import java.util.Optional;

public class Testcase implements LocationAware {

    private final String testSuiteName;
    private final String protocolName;
    private final String outputFlavor;
    private final String name;
    private final String description;
    private final byte[] raw;
    private final String rootType;
    private final List<String> parserArguments;
    private final Element xml;

    private Location location;

    public Testcase(String testSuiteName, String protocolName, String outputFlavor, String name, String description, byte[] raw, String rootType, List<String> parserArguments, Element xml) {
        this.testSuiteName = testSuiteName;
        this.protocolName = protocolName;
        this.outputFlavor = outputFlavor;
        this.name = name;
        this.description = description;
        this.raw = raw;
        this.rootType = rootType;
        this.parserArguments = parserArguments;
        this.xml = xml;
    }

    public String getTestSuiteName() {
        return testSuiteName;
    }

    public String getProtocolName() {
        return protocolName;
    }

    public String getOutputFlavor() {
        return outputFlavor;
    }

    public String getName() {
        return name;
    }

    public String getDescription() {
        return description;
    }

    public byte[] getRaw() {
        return raw;
    }

    public String getRootType() {
        return rootType;
    }

    public List<String> getParserArguments() {
        return parserArguments;
    }

    public Element getXml() {
        return xml;
    }

    public Optional<Location> getLocation() {
        return Optional.ofNullable(location);
    }

    public void setLocation(Location location) {
        this.location = location;
    }
}
