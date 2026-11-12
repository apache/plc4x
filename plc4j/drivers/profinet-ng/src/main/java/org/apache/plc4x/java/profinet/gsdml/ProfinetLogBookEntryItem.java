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

public class ProfinetLogBookEntryItem {

    private String status;

    private ProfinetErrorCode2Value errorCode2Value;

    public String getStatus() {
        return status;
    }

    public ProfinetErrorCode2Value getErrorCode2Value() {
        return errorCode2Value;
    }

    static ProfinetLogBookEntryItem fromXml(Element element) throws IOException {
        ProfinetLogBookEntryItem result = new ProfinetLogBookEntryItem();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, false)) {
            switch (property.name()) {
                case "Status" -> result.status = property.string();
                case "ErrorCode2Value" -> result.errorCode2Value = property.object(ProfinetErrorCode2Value::fromXml);
                default -> {
                    // Not used by the driver.
                }
            }
        }
        return result;
    }
}


