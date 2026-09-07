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
package org.apache.plc4x.java.spi.drivers.model;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * The array notation shared by every PLC4X tag address.
 *
 * <pre>
 * array-expression = dimension , { dimension } ;
 * dimension        = "[" , bounds , [ ";" , base ] , "]" ;
 * bounds           = index , [ ".." , index ] ;
 * </pre>
 *
 * <p>A range is inclusive of both bounds, so {@code [0..7]} is eight elements. A bare index is
 * one element, so {@code [4]} is the fifth. {@code ;base} states the array's declared lower
 * bound - {@code [4..7;1]} selects elements 4 to 7 of an array declared from 1, which sit at
 * offsets 3 to 6 - and defaults to 0. Each bracket group is one dimension, in written order.
 *
 * <p>This class is the single definition of that grammar. Drivers differ in what they can encode,
 * not in what they accept, so a driver states its limits in {@link AddressConstraints} and a
 * selection exceeding them is reported here rather than truncated when the request is serialized.
 *
 * <p>Patterns are compiled once: address parsing runs per tag per request in some drivers.
 */
public final class ArrayNotationParser {

    /** One dimension: an index or an inclusive range, with an optional declared lower bound. */
    private static final String DIMENSION = "\\d+(?:\\.\\.\\d+)?(?:;\\d+)?";

    /**
     * The array expression, as a regex fragment with no capturing groups. A driver embeds this in
     * its own address pattern so that the grammar has one definition rather than a copy per
     * driver - see {@link #ARRAY_GROUP}.
     */
    public static final String EXPRESSION_REGEX = "(?:\\[" + DIMENSION + "(?:," + DIMENSION + ")*])+";

    /**
     * The array expression as an optional named group, ready to splice into a driver's address
     * pattern between the address and the type. The group is named {@code array}.
     */
    public static final String ARRAY_GROUP = "(?<array>" + EXPRESSION_REGEX + ")?";

    /** A trailing run of strictly numeric bracket groups. */
    private static final Pattern EXPRESSION_PATTERN = Pattern.compile(EXPRESSION_REGEX + "$");

    /** One bracket group, whose content is one or more comma-separated dimensions. */
    private static final Pattern GROUP_PATTERN = Pattern.compile("\\[([^\\]]*)]");

    private static final Pattern DIMENSION_PATTERN =
        Pattern.compile("(\\d+)(?:\\.\\.(\\d+))?(?:;(\\d+))?");

    private ArrayNotationParser() {
    }

    /**
     * An address in the pre-migration shape {@code address:TYPE[n]}, where the selection came
     * after the type and meant a count.
     */
    private static final Pattern LEGACY_AFTER_TYPE = Pattern.compile(
        "^(?<head>.+?):(?<type>[A-Za-z_][A-Za-z_0-9]*(?:\\(\\d+\\))?)\\[(?<count>\\d+)]$");

    /** An address in the pre-migration shape {@code address:TYPE:n}, where a count trailed it. */
    private static final Pattern LEGACY_COUNT_SUFFIX = Pattern.compile(
        "^(?<head>.+?):(?<type>[A-Za-z_][A-Za-z_0-9]*)(?::(?<count>\\d+))$");

    /**
     * How to rewrite an address that was written before the array notation was unified, or empty
     * when it does not look like one.
     *
     * <p>The brackets moved from after the type to before it, and a count became a range. An
     * upgrading user who sees only "does not match pattern" has to work that out from a regex;
     * this hands them the address they meant.
     *
     * @param address the address that failed to parse
     * @return the equivalent address in the current notation, if one can be worked out
     */
    public static Optional<String> currentFormOf(String address) {
        if (address == null) {
            return Optional.empty();
        }
        Matcher afterType = LEGACY_AFTER_TYPE.matcher(address);
        if (afterType.matches()) {
            return Optional.of(afterType.group("head")
                + rangeFor(afterType.group("count")) + ":" + afterType.group("type"));
        }
        Matcher countSuffix = LEGACY_COUNT_SUFFIX.matcher(address);
        if (countSuffix.matches()) {
            return Optional.of(countSuffix.group("head")
                + rangeFor(countSuffix.group("count")) + ":" + countSuffix.group("type"));
        }
        return Optional.empty();
    }

    private static String rangeFor(String count) {
        int elements = Integer.parseInt(count);
        return elements <= 1 ? "[0]" : "[0.." + (elements - 1) + "]";
    }

    /**
     * Reports an address the driver could not parse, naming the form it expected and - when the
     * address looks like one written before the notation was unified - the address to write
     * instead.
     */
    public static PlcInvalidTagException invalidAddress(String address, String expectedForm) {
        String message = "Invalid address '" + address + "': expected " + expectedForm;
        Optional<String> current = currentFormOf(address);
        if (current.isPresent()) {
            message += ". The array notation moved before the type and a count became a range,"
                + " so this address is now written '" + current.get() + "'";
        }
        return new PlcInvalidTagException(message);
    }

    /**
     * The part of an address before any trailing array expression. An address with no such
     * expression is returned unchanged - including one whose brackets are not numeric, such as an
     * OPC UA string identifier that happens to contain them.
     */
    public static String addressPart(String address) {
        Matcher matcher = EXPRESSION_PATTERN.matcher(address);
        return matcher.find() ? address.substring(0, matcher.start()) : address;
    }

    /** The trailing array expression of an address, or the empty string if it has none. */
    public static String expressionPart(String address) {
        Matcher matcher = EXPRESSION_PATTERN.matcher(address);
        return matcher.find() ? matcher.group() : "";
    }

    /**
     * Whether an expression selects a single element rather than an array, which decides what a
     * tag reports from {@code getArrayInfo()}: a bare index yields a scalar, a range yields an
     * array even when it spans one element. {@code [1]} is a scalar and {@code [1..1]} is an
     * array of one, so equal bounds alone cannot tell them apart - only the written form can.
     *
     * <p>An expression is a single element when every one of its dimensions is a bare index.
     */
    public static boolean selectsSingleElement(String expression) {
        return expression != null && !expression.isEmpty() && !expression.contains("..");
    }

    /** Parses an array expression with no constraints beyond the grammar. */
    public static List<ArrayInfo> parse(String expression, String address) {
        return parse(expression, address, AddressConstraints.UNCONSTRAINED);
    }

    /**
     * Parses an array expression into one {@link ArrayInfo} per dimension, in written order.
     *
     * @param expression the bracket run, or empty for an address that selects no range
     * @param address the whole address, quoted back in any error so the user can find it
     * @param constraints what the calling driver's protocol can encode
     * @throws PlcInvalidTagException if the expression is malformed or exceeds the constraints
     */
    public static List<ArrayInfo> parse(String expression, String address, AddressConstraints constraints) {
        if (expression == null || expression.isEmpty()) {
            return Collections.emptyList();
        }
        if (!EXPRESSION_PATTERN.matcher(expression).matches()) {
            throw new PlcInvalidTagException(String.format(
                "Invalid array expression '%s' in tag '%s': expected [index], [lo..hi] or either "
                    + "with a ';base', repeated once per dimension", expression, address));
        }

        List<ArrayInfo> dimensions = new ArrayList<>(2);
        Matcher group = GROUP_PATTERN.matcher(expression);
        while (group.find()) {
            // A group may hold several dimensions separated by commas - the spelling Allen-Bradley
            // and others use. "[1..2,3..4]" and "[1..2][3..4]" are the same selection.
            for (String part : group.group(1).split(",")) {
                Matcher dimension = DIMENSION_PATTERN.matcher(part);
                if (!dimension.matches()) {
                    throw new PlcInvalidTagException(String.format(
                        "Invalid array dimension '%s' in tag '%s'", part, address));
                }
                dimensions.add(toDimension(dimension, address, constraints));
            }
        }

        if (dimensions.size() > constraints.maxDimensions()) {
            throw new PlcInvalidTagException(String.format(
                "Array expression '%s' in tag '%s' has %d dimensions, but this protocol carries "
                    + "at most %d", expression, address, dimensions.size(), constraints.maxDimensions()));
        }
        if (constraints.onlyTrailingDimensionMayBeRange()) {
            for (int i = 0; i < dimensions.size() - 1; i++) {
                // What matters is how the dimension was written, not how wide it turned out to
                // be: "[1..1]" is a range that happens to span one element, and letting it pass
                // here would hand the driver a leading range it has no element count for.
                if (dimensions.get(i).isRange()) {
                    throw new PlcInvalidTagException(String.format(
                        "Array expression '%s' in tag '%s' writes dimension %d as a range, but "
                            + "this protocol carries one element count for the whole address, so "
                            + "only the last dimension may be a range",
                        expression, address, i + 1));
                }
            }
        }
        return Collections.unmodifiableList(dimensions);
    }

    private static ArrayInfo toDimension(Matcher bracket, String address, AddressConstraints constraints) {
        int lowerBound = parseIndex(bracket.group(1), bracket.group(), address);
        boolean range = bracket.group(2) != null;
        int upperBound = range ? parseIndex(bracket.group(2), bracket.group(), address) : lowerBound;
        int base = bracket.group(3) != null
            ? parseIndex(bracket.group(3), bracket.group(), address)
            : 0;

        if (upperBound < lowerBound) {
            throw new PlcInvalidTagException(String.format(
                "Invalid array range '%s' in tag '%s': the upper bound %d is below the lower bound %d",
                bracket.group(), address, upperBound, lowerBound));
        }
        // The inclusive size is computed in an int, so a range spanning more than Integer.MAX_VALUE
        // elements would wrap to a negative count and be reported as a syntactically valid
        // selection of minus two billion elements. No protocol here can carry such a range anyway.
        if (((long) upperBound - lowerBound + 1) > Integer.MAX_VALUE) {
            throw new PlcInvalidTagException(String.format(
                "Invalid array range '%s' in tag '%s': it spans %d elements, more than can be counted",
                bracket.group(), address, (long) upperBound - lowerBound + 1));
        }
        if (lowerBound - base < 0) {
            throw new PlcInvalidTagException(String.format(
                "Invalid array range '%s' in tag '%s': index %d lies below the declared lower bound %d",
                bracket.group(), address, lowerBound, base));
        }
        // The bound applies to the offset the protocol actually encodes - the start of the
        // selection - not to the last element of it. A CIP request carries a start index and an
        // element count, so [0..300] is encodable where [300] is not.
        if (lowerBound - base > constraints.maxIndex()) {
            throw new PlcInvalidTagException(String.format(
                "Invalid array range '%s' in tag '%s': index %d is out of range 0 to %d for this protocol",
                bracket.group(), address, lowerBound - base, constraints.maxIndex()));
        }
        return new DefaultArrayInfo(lowerBound, upperBound, base, range);
    }

    private static int parseIndex(String value, String bracket, String address) {
        try {
            return Integer.parseInt(value);
        } catch (NumberFormatException e) {
            throw new PlcInvalidTagException(String.format(
                "Invalid array range '%s' in tag '%s': '%s' is not a number this protocol can address",
                bracket, address, value));
        }
    }

    /**
     * Renders dimensions back to their canonical form: one bracket per dimension, omitting what
     * is defaulted - a base of 0 is dropped, and equal bounds render as a bare index. The
     * comma-separated spelling is accepted on input but never produced, so
     * {@code [1..2,3..4]} renders as {@code [1..2][3..4]}. Parsing the result yields equal
     * dimensions, so an address round-trips by meaning rather than by spelling.
     */
    public static String render(List<ArrayInfo> dimensions) {
        if (dimensions == null || dimensions.isEmpty()) {
            return "";
        }
        StringBuilder sb = new StringBuilder();
        for (ArrayInfo dimension : dimensions) {
            sb.append('[').append(dimension.getLowerBound());
            if (dimension.isRange()) {
                // A one-element range still renders as a range: [8..8] is an array of one and
                // [8] is a scalar, so collapsing it would change what the address means.
                sb.append("..").append(dimension.getUpperBound());
            }
            if (dimension.getBase() != 0) {
                sb.append(';').append(dimension.getBase());
            }
            sb.append(']');
        }
        return sb.toString();
    }
}
