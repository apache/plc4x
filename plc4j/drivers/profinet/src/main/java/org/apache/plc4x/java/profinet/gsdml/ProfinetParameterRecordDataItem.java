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

package org.apache.plc4x.java.profinet.gsdml;

import org.w3c.dom.Element;

import java.io.IOException;

public class ProfinetParameterRecordDataItem {

    private int index;

    private int length;

    private ProfinetTextId name;

    private ProfinetRef ref;

    public int getIndex() {
        return index;
    }

    public int getLength() {
        return length;
    }

    public ProfinetTextId getName() {
        return name;
    }

    public ProfinetRef getRef() {
        return ref;
    }

    static ProfinetParameterRecordDataItem fromXml(Element element) throws IOException {
        ProfinetParameterRecordDataItem result = new ProfinetParameterRecordDataItem();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, false)) {
            switch (property.name()) {
                case "Index" -> result.index = property.intValue();
                case "Length" -> result.length = property.intValue();
                case "Name" -> result.name = property.object(ProfinetTextId::fromXml);
                case "Ref" -> result.ref = property.object(ProfinetRef::fromXml);
                default -> {
                    // Not used by the driver.
                }
            }
        }
        return result;
    }
}
