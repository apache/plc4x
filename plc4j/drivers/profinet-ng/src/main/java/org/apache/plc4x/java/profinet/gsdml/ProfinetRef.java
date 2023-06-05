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

package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonRootName;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("Ref")
public class ProfinetRef {

    @JacksonXmlProperty(isAttribute=true, localName="DataType")
    private String dataType;

    @JacksonXmlProperty(isAttribute=true, localName="ByteOffset")
    private int byteOffset;

    @JacksonXmlProperty(isAttribute=true, localName="DefaultValue")
    private String defaultValue;

    @JacksonXmlProperty(isAttribute=true, localName="AllowedValues")
    private String allowedValues;

    @JacksonXmlProperty(isAttribute=true, localName="Changeable")
    private boolean changeable;

    @JacksonXmlProperty(isAttribute=true, localName="Visible")
    private boolean visible;

    @JacksonXmlProperty(isAttribute=true, localName="TextId")
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
}
