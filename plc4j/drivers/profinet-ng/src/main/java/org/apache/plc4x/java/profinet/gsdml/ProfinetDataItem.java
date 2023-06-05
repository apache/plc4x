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
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlElementWrapper;
import com.fasterxml.jackson.dataformat.xml.annotation.JacksonXmlProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = false)
@JsonRootName("DataItem")
public class ProfinetDataItem {

    @JacksonXmlProperty(localName="BitDataItem")
    @JacksonXmlElementWrapper(useWrapping = false)
    private List<ProfinetBitDataItem> bitDataItem;

    @JacksonXmlProperty(isAttribute=true, localName="DataType")
    private String dataType;

    @JacksonXmlProperty(isAttribute=true, localName="TextId")
    private String textId;

    @JacksonXmlProperty(isAttribute=true, localName="Length")
    private int length;

    @JacksonXmlProperty(isAttribute=true, localName="UseAsBits")
    private boolean useAsBits;

    public String getDataType() {
        return dataType;
    }

    public String getTextId() {
        return textId;
    }

    public int getLength() {
        return length;
    }

    public boolean isUseAsBits() {
        return useAsBits;
    }

    public List<ProfinetBitDataItem> getBitDataItem() {
        return bitDataItem;
    }
}


