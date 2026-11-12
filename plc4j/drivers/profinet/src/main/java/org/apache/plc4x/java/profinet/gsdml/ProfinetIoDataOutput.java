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
import java.util.List;

public class ProfinetIoDataOutput {

    private String consistency;

    private List<ProfinetDataItem> dataItemList;

    public String getConsistency() {
        return consistency;
    }

    public List<ProfinetDataItem> getDataItemList() {
        return dataItemList;
    }

    static ProfinetIoDataOutput fromXml(Element element) throws IOException {
        ProfinetIoDataOutput result = new ProfinetIoDataOutput();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, true)) {
            switch (property.name()) {
                case "Consistency" -> result.consistency = property.string();
                case "DataItem" -> result.dataItemList = property.unwrappedList(ProfinetDataItem::fromXml);
                // Strict, as under the former jackson mapping: anything unknown here rejects the whole file.
                default -> throw property.unknownIn(ProfinetIoDataOutput.class);
            }
        }
        return result;
    }
}
