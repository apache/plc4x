/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.protocol.test.model;

import org.dom4j.Element;
import org.w3c.dom.Document;

public class Testcase {

    private final String name;
    private final String description;
    private final byte[] raw;
    private final String rootType;
    private final Element xml;

    public Testcase(String name, String description, byte[] raw, String rootType, Element xml) {
        this.name = name;
        this.description = description;
        this.raw = raw;
        this.rootType = rootType;
        this.xml = xml;
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

    public Element getXml() {
        return xml;
    }

}
