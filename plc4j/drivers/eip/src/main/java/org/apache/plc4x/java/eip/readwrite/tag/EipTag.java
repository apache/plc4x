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
package org.apache.plc4x.java.eip.readwrite.tag;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.eip.readwrite.CIPDataTypeCode;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class EipTag implements PlcTag, Serializable {

    private static final Pattern ADDRESS_PATTERN =
        Pattern.compile("^%(?<tag>[a-zA-Z_.0-9]+\\[?[0-9]*\\]?):?(?<dataType>[A-Z]*):?(?<elementNb>[0-9]*)");

    private static final String TAG = "tag";
    private static final String ELEMENTS = "elementNb";
    private static final String TYPE = "dataType";


    private final String tag;
    private CIPDataTypeCode type;
    private int elementNb;

    public EipTag(String tag) {
        this.tag = tag;
    }

    public EipTag(String tag, int elementNb) {
        this.tag = tag;
        this.elementNb = elementNb;
    }

    public EipTag(String tag, CIPDataTypeCode type, int elementNb) {
        this.tag = tag;
        this.type = type;
        this.elementNb = elementNb;
    }

    public EipTag(String tag, CIPDataTypeCode type) {
        this.tag = tag;
        this.type = type;
    }

    @Override
    public String getAddressString() {
        throw new NotImplementedException("Need to implement this");
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.valueOf(type.name());
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return PlcTag.super.getArrayInfo();
    }

    public CIPDataTypeCode getType() {
        return type;
    }

    public void setType(CIPDataTypeCode type) {
        this.type = type;
    }

    public int getElementNb() {
        return elementNb;
    }

    public void setElementNb(int elementNb) {
        this.elementNb = elementNb;
    }

    public String getTag() {
        return tag;
    }

    public static boolean matches(String tagQuery) {
        return ADDRESS_PATTERN.matcher(tagQuery).matches();
    }

    public static EipTag of(String tagString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(tagString);
        if (matcher.matches()) {
            String tag = matcher.group(TAG);
            int nb = 0;
            CIPDataTypeCode type = null;
            if (!matcher.group(ELEMENTS).isEmpty()) {
                nb = Integer.parseInt(matcher.group(ELEMENTS));
            }
            if (!matcher.group(TYPE).isEmpty()) {
                type = CIPDataTypeCode.valueOf(matcher.group(TYPE));
            }
            if (nb != 0) {
                if (type != null) {
                    return new EipTag(tag, type, nb);
                }
                return new EipTag(tag, nb);
            } else {
                if (type != null) {
                    return new EipTag(tag, type);
                }
                return new EipTag(tag);
            }
        }
        return null;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());

        writeBuffer.writeString("node",
            tag.getBytes(StandardCharsets.UTF_8).length * 8,
            tag, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
        if (type != null) {
            writeBuffer.writeString("type",
                type.name().getBytes(StandardCharsets.UTF_8).length * 8,
                type.name(), WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
        }
        writeBuffer.writeUnsignedInt("elementNb", 16, elementNb);
        // TODO: remove this (not language agnostic)
        String defaultJavaType = (type == null ? Object.class : getPlcValueType().getDefaultJavaType()).getName();
        writeBuffer.writeString("defaultJavaType",
            defaultJavaType.getBytes(StandardCharsets.UTF_8).length * 8,
            defaultJavaType, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));

        writeBuffer.popContext(getClass().getSimpleName());
    }

}
