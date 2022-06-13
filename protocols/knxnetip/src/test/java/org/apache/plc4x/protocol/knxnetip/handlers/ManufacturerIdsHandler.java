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

import org.xml.sax.Attributes;
import org.xml.sax.SAXException;
import org.xml.sax.helpers.DefaultHandler;

import java.util.ArrayList;
import java.util.List;

public class ManufacturerIdsHandler extends DefaultHandler {

    private boolean inElement = false;

    private List<Integer> manufacturerIds = new ArrayList<>();

    @Override
    public void startElement(String uri, String localName, String qName, Attributes attributes) throws SAXException {
        inElement = qName.equalsIgnoreCase("unsignedShort");
    }

    @Override
    public void endElement(String uri, String localName, String qName) throws SAXException {
        inElement = ! (inElement && qName.equalsIgnoreCase("unsignedShort"));
    }

    public void characters(char[] ch, int start, int length) throws SAXException {
        if(inElement) {
            String content = new String(ch, start, length);
            int manufacturerId = Integer.parseInt(content);
            manufacturerIds.add(manufacturerId);
        }
    }

    public List<Integer> getManufacturerIds() {
        return manufacturerIds;
    }

}
