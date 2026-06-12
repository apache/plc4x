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
package org.apache.plc4x.java.ads.resolution;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/**
 * Parses TwinCAT-style symbolic addresses into a recursive tree of segments.
 *
 * <p>An address like {@code MAIN.g_plant.channels[1].setpoints[4]} parses into:
 * <pre>
 *   AddressPart(baseSegment="MAIN.g_plant", indices=[],
 *     child=AddressPart("channels", [1],
 *       child=AddressPart("setpoints", [4], child=null)))
 * </pre>
 * The first {@code .} separates the namespace prefix (typically {@code MAIN}) from
 * the variable name; both join to form the symbol-table lookup key.
 */
public final class AddressParser {

    private AddressParser() {
    }

    /**
     * Address segment: a field name, zero or more {@code [i]} indices, and an optional child.
     *
     * @param baseSegment   Field-name part of the segment (or the symbol name for the root).
     * @param arrayIndices  Indices in source order — {@code [1][2][3]} → {@code [1, 2, 3]}.
     * @param child         Child segment, or {@code null} when this is the leaf.
     */
    public record AddressPart(String baseSegment, List<Integer> arrayIndices, AddressPart child) {
    }

    /**
     * Parse a full symbolic address. Throws {@link PlcInvalidTagException} on syntactic problems.
     */
    public static AddressPart parse(String symbolicAddress) {
        if (symbolicAddress == null || symbolicAddress.isEmpty()) {
            throw new PlcInvalidTagException("Empty symbolic address");
        }
        int dot = symbolicAddress.indexOf('.');
        if (dot < 0) {
            // No namespace prefix — uncommon, but accept it.
            return new AddressPart(extractField(symbolicAddress), extractIndices(symbolicAddress), null);
        }
        String rootName = symbolicAddress.substring(0, dot);
        AddressPart rest = parseRest(symbolicAddress.substring(dot + 1));
        // Glue the namespace onto the first segment so that the symbol-table lookup key is
        // formed (e.g. "MAIN" + "." + "g_plant" → "MAIN.g_plant").
        return new AddressPart(rootName + "." + rest.baseSegment, rest.arrayIndices, rest.child);
    }

    private static AddressPart parseRest(String s) {
        int dot = s.indexOf('.');
        if (dot < 0) {
            return new AddressPart(extractField(s), extractIndices(s), null);
        }
        String head = s.substring(0, dot);
        AddressPart child = parseRest(s.substring(dot + 1));
        return new AddressPart(extractField(head), extractIndices(head), child);
    }

    private static String extractField(String segment) {
        int bracket = segment.indexOf('[');
        return bracket >= 0 ? segment.substring(0, bracket) : segment;
    }

    private static List<Integer> extractIndices(String segment) {
        int bracket = segment.indexOf('[');
        if (bracket < 0) {
            return Collections.emptyList();
        }
        List<Integer> indices = new ArrayList<>();
        int start = bracket;
        while (true) {
            int open = segment.indexOf('[', start);
            if (open < 0) {
                return indices;
            }
            int close = segment.indexOf(']', open);
            if (close < 0) {
                throw new PlcInvalidTagException("Unmatched '[' in address segment: " + segment);
            }
            String num = segment.substring(open + 1, close);
            try {
                indices.add(Integer.parseInt(num.trim()));
            } catch (NumberFormatException e) {
                throw new PlcInvalidTagException("Non-integer array index in segment: " + segment);
            }
            start = close + 1;
        }
    }
}
