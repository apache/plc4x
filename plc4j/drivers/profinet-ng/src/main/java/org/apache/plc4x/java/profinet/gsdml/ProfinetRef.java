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

public class ProfinetRef {

    private String dataType;

    private int byteOffset;

    private String defaultValue;

    private String allowedValues;

    private boolean changeable;

    private boolean visible;

    private String textId;

    public String getDataType() {
        return dataType;
    }

    public int getByteOffset() {
        return byteOffset;
    }

    public String getDefaultValue() {
        return defaultValue;
    }

    public String getAllowedValues() {
        return allowedValues;
    }

    public boolean isChangeable() {
        return changeable;
    }

    public boolean isVisible() {
        return visible;
    }

    public String getTextId() {
        return textId;
    }

    static ProfinetRef fromXml(Element element) throws IOException {
        ProfinetRef result = new ProfinetRef();
        for (GsdmlXml.Property property : GsdmlXml.properties(element, false)) {
            switch (property.name()) {
                case "DataType" -> result.dataType = property.string();
                case "ByteOffset" -> result.byteOffset = property.intValue();
                case "DefaultValue" -> result.defaultValue = property.string();
                case "AllowedValues" -> result.allowedValues = property.string();
                case "Changeable" -> result.changeable = property.booleanValue();
                case "Visible" -> result.visible = property.booleanValue();
                case "TextId" -> result.textId = property.string();
                default -> {
                    // Not used by the driver.
                }
            }
        }
        return result;
    }
}
