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
package org.apache.plc4x.test.dom4j;

import org.dom4j.DocumentFactory;
import org.dom4j.ElementHandler;
import org.dom4j.io.SAXContentHandler;
import org.xml.sax.Locator;

public class LocationAwareSAXContentHandler extends SAXContentHandler {

    private final DocumentFactory documentFactory;
    private Locator locator;

    public LocationAwareSAXContentHandler(DocumentFactory documentFactory, ElementHandler elementHandler) {
        super(documentFactory, elementHandler);
        this.documentFactory = documentFactory;
    }

    @Override
    public void setDocumentLocator(Locator documentLocator) {
        super.setDocumentLocator(documentLocator);
        this.locator = documentLocator;
        if (documentFactory instanceof LocationAwareDocumentFactory) {
            ((LocationAwareDocumentFactory)documentFactory).setLocator(documentLocator);
        }
    }

    public Locator getLocator() {
        return locator;
    }

}
