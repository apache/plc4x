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
package org.apache.plc4x.protocol.knxnetip.handlers;

import org.apache.commons.lang3.StringUtils;
import org.xml.sax.Attributes;
import org.xml.sax.SAXException;
import org.xml.sax.helpers.DefaultHandler;

import java.util.HashMap;
import java.util.Map;

public class ProductDescriptionHandler extends DefaultHandler {

    private Map<String, Integer> addresses = new HashMap<>();
    private String maskVersion = null;
    private String name = null;
    private String[] replacesVersions;
    private Integer comObjectTableAddress = null;

    @Override
    public void startElement(String uri, String localName, String qName, Attributes attributes) throws SAXException {
        if(qName.equalsIgnoreCase("ApplicationProgram")) {
            String maskVersionString = attributes.getValue("MaskVersion");
            maskVersion = maskVersionString.substring(maskVersionString.indexOf('-') + 1);
            name = attributes.getValue("Name");
            String replacesVersionsString = attributes.getValue("ReplacesVersions");
            if(StringUtils.isNotBlank(replacesVersionsString)) {
                replacesVersions = replacesVersionsString.split(" ");
            }
        } else if(qName.equalsIgnoreCase("AbsoluteSegment")) {
            String id = attributes.getValue("Id");
            Integer address = Integer.parseInt(attributes.getValue("Address"));
            addresses.put(id, address);
        } else if(qName.equalsIgnoreCase("ComObjectTable")) {
            String codeSegment = attributes.getValue("CodeSegment");
            comObjectTableAddress = addresses.get(codeSegment);
        }
    }

    public String getMaskVersion() {
        return maskVersion;
    }

    public String getName() {
        return name;
    }

    public String[] getReplacesVersions() {
        return replacesVersions;
    }

    public Integer getComObjectTableAddress() {
        return comObjectTableAddress;
    }

}
