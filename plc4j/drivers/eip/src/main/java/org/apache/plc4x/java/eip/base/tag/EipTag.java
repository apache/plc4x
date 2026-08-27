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

package org.apache.plc4x.java.eip.base.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.eip.readwrite.CIPDataTypeCode;
import org.apache.plc4x.java.spi.buffers.api.Serializable;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class EipTag implements PlcTag, Serializable {

    private static final Pattern ADDRESS_PATTERN =
        Pattern.compile("^(?<tag>[%a-zA-Z_.0-9]+\\[?[0-9]*]?):?(?<dataType>[A-Z]*):?(?<elementNb>[0-9]*)");

    /**
     * Splits an address into the members it names. A qualifier of '[' introduces an array
     * member, anything else another symbol; characters outside a member name (a leading '%')
     * are skipped. Compiled once - a tag address is decomposed when the tag is built, not
     * again for every request that uses it.
     */
    private static final Pattern PATH_PATTERN = Pattern.compile("([.\\[\\]])*([A-Za-z_0-9]+)");

    private static final String GROUP_NAME_TAG = "tag";
    private static final String GROUP_NAME_GROUP_NAME_ELEMENTS = "elementNb";
    private static final String GROUP_NAME_TYPE = "dataType";

    /** One step of the CIP path an address describes. */
    public sealed interface PathElement permits SymbolElement, MemberElement {
    }

    /** A named symbol - a tag, or one member of a structured tag. */
    public record SymbolElement(String name) implements PathElement {
    }

    /** An array member, addressed by index. */
    public record MemberElement(short index) implements PathElement {
    }

    private final String tag;
    private final CIPDataTypeCode type;
    private final int elementNb;
    private final List<PathElement> pathElements;

    public EipTag(String tag) {
        this(tag, null, 1);
    }

    public EipTag(String tag, int elementNb) {
        this(tag, null, elementNb);
    }

    public EipTag(String tag, CIPDataTypeCode type) {
        this(tag, type, 1);
    }

    public EipTag(String tag, CIPDataTypeCode type, int elementNb) {
        this.tag = tag;
        this.type = type;
        // A request always asks for at least one element; normalising here keeps every use site
        // from having to guard against a count of zero.
        this.elementNb = Math.max(elementNb, 1);
        this.pathElements = decomposePath(tag);
    }

    private static List<PathElement> decomposePath(String tag) {
        if (tag == null) {
            return Collections.emptyList();
        }
        Matcher matcher = PATH_PATTERN.matcher(tag);
        List<PathElement> elements = new ArrayList<>(2);
        while (matcher.find()) {
            String identifier = matcher.group(2);
            if ("[".equals(matcher.group(1))) {
                // MemberID.instance is a uint 8 in the mspec, so it holds 0 to 255. Reject a
                // larger index here rather than letting it fail later during serialization.
                int index = Integer.parseInt(identifier);
                if (index > 255) {
                    throw new PlcInvalidTagException(String.format(
                        "Error parsing address %s, index %d is out of range 0 to 255", tag, index));
                }
                elements.add(new MemberElement((short) index));
            } else {
                elements.add(new SymbolElement(identifier));
            }
        }
        return Collections.unmodifiableList(elements);
    }

    @Override
    public String getAddressString() {
        // Mirrors the format ADDRESS_PATTERN accepts, so of(getAddressString()) round-trips:
        //   tag[:dataType[:elementNb]]
        StringBuilder sb = new StringBuilder(tag);
        if (type != null) {
            sb.append(':').append(type.name());
            if (elementNb > 1) {
                sb.append(':').append(elementNb);
            }
        }
        return sb.toString();
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

    public int getElementNb() {
        return elementNb;
    }

    public String getTag() {
        return tag;
    }

    /**
     * The CIP path this address describes, one element per member, decomposed when the tag was
     * built. A structured address yields one element per member - {@code a.b} is two symbols,
     * not one symbol named "a.b" - because that is how a controller walks the path.
     */
    public List<PathElement> getPathElements() {
        return pathElements;
    }

    public static boolean matches(String tagQuery) {
        return ADDRESS_PATTERN.matcher(tagQuery).matches();
    }

    public static EipTag of(String tagString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(tagString);
        if (matcher.matches()) {
            String tag = matcher.group(GROUP_NAME_TAG);
            int nb = 1;
            CIPDataTypeCode type;
            if (!matcher.group(GROUP_NAME_GROUP_NAME_ELEMENTS).isEmpty()) {
                nb = Integer.parseInt(matcher.group(GROUP_NAME_GROUP_NAME_ELEMENTS));
            }
            if (!matcher.group(GROUP_NAME_TYPE).isEmpty()) {
                type = CIPDataTypeCode.valueOf(matcher.group(GROUP_NAME_TYPE));
            } else {
                type = CIPDataTypeCode.DINT;
            }
            return new EipTag(tag, type, nb);
        }
        return null;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        writeBuffer.pushContext(WithOption.WithName(getClass().getSimpleName()));

        writeBuffer.writeString(tag.getBytes(StandardCharsets.UTF_8).length * 8, tag,
            WithOption.WithName("node"), WithOption.WithEncoding("UTF8"));
        if (type != null) {
            writeBuffer.writeString(type.name().getBytes(StandardCharsets.UTF_8).length * 8, type.name(),
                WithOption.WithName("type"), WithOption.WithEncoding("UTF8"));
        }
        writeBuffer.writeUnsignedInt(16, elementNb, WithOption.WithName("elementNb"));

        writeBuffer.popContext(WithOption.WithName(getClass().getSimpleName()));
    }

}
