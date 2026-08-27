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

import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.spi.drivers.model.AddressConstraints;
import org.apache.plc4x.java.spi.drivers.model.ArrayNotationParser;
import org.apache.plc4x.java.spi.drivers.model.DefaultArrayInfo;
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

    private static final Pattern ADDRESS_PATTERN = Pattern.compile(
        "^(?<tag>[%a-zA-Z_.0-9]+)"
            + ArrayNotationParser.ARRAY_GROUP
            + "(?::(?<dataType>[A-Z]+))?$");

    /** What a CIP request can encode of a selection: one dimension, starting no later than 255. */
    private static final AddressConstraints CONSTRAINTS = AddressConstraints.SINGLE_DIMENSION
        .withMaxIndex(255)
        .withOnlyTrailingDimensionMayBeRange(true);

    /**
     * Splits an address into the members it names. A qualifier of '[' introduces an array
     * member, anything else another symbol; characters outside a member name (a leading '%')
     * are skipped. Compiled once - a tag address is decomposed when the tag is built, not
     * again for every request that uses it.
     */
    private static final Pattern PATH_PATTERN = Pattern.compile("([.\\[\\]])*([A-Za-z_0-9]+)");

    private static final String GROUP_NAME_TAG = "tag";
    private static final String GROUP_NAME_ARRAY = "array";
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
    /**
     * What the address selects. Drives the CIP path and the element count. Not what
     * {@link #getArrayInfo()} reports - see there.
     */
    private final List<ArrayInfo> selection;
    private final boolean scalarSelection;
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

    /**
     * Builds a tag selecting {@code elementNb} elements from the start of {@code tag}, which is
     * the shape the older constructors describe. An address selecting a range is built through
     * {@link #of(String)}.
     */
    public EipTag(String tag, CIPDataTypeCode type, int elementNb) {
        this(tag, type, rangeOf(tag, Math.max(elementNb, 1)));
    }

    public EipTag(String tag, CIPDataTypeCode type, List<ArrayInfo> selection) {
        this.tag = tag;
        this.type = type;
        this.selection = selection == null ? Collections.emptyList() : List.copyOf(selection);
        this.scalarSelection =
            ArrayNotationParser.selectsSingleElement(ArrayNotationParser.expressionPart(tag));
        this.pathElements = decomposePath(tag, this.selection);
    }

    /**
     * The selection an element count describes on its own: the first {@code count} elements,
     * unless the address itself already names a starting index.
     */
    private static List<ArrayInfo> rangeOf(String tag, int count) {
        String expression = tag == null ? "" : ArrayNotationParser.expressionPart(tag);
        int start = 0;
        if (!expression.isEmpty()) {
            ArrayInfo first = ArrayNotationParser.parse(expression, tag, CONSTRAINTS).get(0);
            start = first.getLowerBound() - first.getBase();
        } else if (count <= 1) {
            return Collections.emptyList();
        }
        return List.of(new DefaultArrayInfo(start, start + count - 1));
    }

    /**
     * The CIP path of an address: one segment per member named in it, followed by a member
     * segment for the element the selection starts at. A structured address yields one segment
     * per member - {@code a.b} is two symbols, not one symbol named "a.b" - because that is how
     * a controller walks a path.
     *
     * <p>A selection given as a bare element count carries no member segment: the request starts
     * at the tag itself and asks for that many elements.
     */
    private static List<PathElement> decomposePath(String tag, List<ArrayInfo> arrayInfo) {
        if (tag == null) {
            return Collections.emptyList();
        }
        Matcher matcher = PATH_PATTERN.matcher(ArrayNotationParser.addressPart(tag));
        List<PathElement> elements = new ArrayList<>(2);
        while (matcher.find()) {
            if ("[".equals(matcher.group(1))) {
                // An index written inside the path rather than as a trailing selection, as in
                // "a[2].b". It addresses one element, so it is a member segment like any other.
                elements.add(new MemberElement(Short.parseShort(matcher.group(2))));
            } else {
                elements.add(new SymbolElement(matcher.group(2)));
            }
        }
        if (!ArrayNotationParser.expressionPart(tag).isEmpty() && !arrayInfo.isEmpty()) {
            ArrayInfo first = arrayInfo.get(0);
            short offset = (short) (first.getLowerBound() - first.getBase());
            // A member segment says where the selection starts. Starting at the first element is
            // what a request with no member segment already means, so emitting MemberID(0) would
            // add two bytes that say nothing - and every address of the form "tag:TYPE:n" used to
            // be sent without one.
            if (offset > 0) {
                elements.add(new MemberElement(offset));
            }
        }
        return Collections.unmodifiableList(elements);
    }

    @Override
    public String getAddressString() {
        // Mirrors the format ADDRESS_PATTERN accepts, so of(getAddressString()) round-trips:
        //   tag[range][:dataType]
        StringBuilder sb = new StringBuilder(ArrayNotationParser.addressPart(tag));
        sb.append(ArrayNotationParser.render(selection));
        if (type != null) {
            sb.append(':').append(type.name());
        }
        return sb.toString();
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.valueOf(type.name());
    }

    /**
     * The shape of the value the caller receives, so a consumer can tell a scalar from a list
     * without knowing the protocol: empty for a scalar, one entry per dimension for an array.
     * A bare index selects one element and so reports empty; a range reports its dimensions even
     * when it spans a single element.
     *
     * <p>Not the same as what the driver fetches - reading one element of an array still walks a
     * member path, which {@link #getPathElements()} describes.
     */
    @Override
    public List<ArrayInfo> getArrayInfo() {
        return scalarSelection ? Collections.emptyList() : selection;
    }

    public CIPDataTypeCode getType() {
        return type;
    }

    /**
     * How many elements the request asks the device for. Derived from the selection; a tag that
     * selects nothing explicitly reads a single element.
     */
    public int getElementNb() {
        int elements = 1;
        for (ArrayInfo dimension : selection) {
            elements *= dimension.getSize();
        }
        return Math.max(elements, 1);
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
        if (!matcher.matches()) {
            return null;
        }
        String tag = matcher.group(GROUP_NAME_TAG);
        String arrayExpression = matcher.group(GROUP_NAME_ARRAY);
        String typeString = matcher.group(GROUP_NAME_TYPE);

        CIPDataTypeCode type = typeString == null || typeString.isEmpty()
            ? CIPDataTypeCode.DINT
            : CIPDataTypeCode.valueOf(typeString);
        List<ArrayInfo> arrayInfo = arrayExpression == null
            ? Collections.emptyList()
            : ArrayNotationParser.parse(arrayExpression, tagString, CONSTRAINTS);

        return new EipTag(tag + (arrayExpression == null ? "" : arrayExpression), type, arrayInfo);
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
        writeBuffer.writeUnsignedInt(16, getElementNb(), WithOption.WithName("elementNb"));

        writeBuffer.popContext(WithOption.WithName(getClass().getSimpleName()));
    }

}
