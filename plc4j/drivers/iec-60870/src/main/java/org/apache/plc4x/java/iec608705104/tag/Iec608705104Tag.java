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

package org.apache.plc4x.java.iec608705104.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.List;
import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Address of one IEC 60870-5-104 information object.
 *
 * <p>The address is written as {@code <asdu>/<ioa>}:
 * <ul>
 *   <li>{@code <asdu>} is the 2 octet <em>common ASDU address</em>, given
 *       either as one decimal number ({@code 0..65535}) or as two
 *       octets {@code <low>/<high>} ({@code 0..255} each).</li>
 *   <li>{@code <ioa>} is the 3 octet <em>information object address</em>,
 *       given either as one decimal number ({@code 0..16777215}) or as
 *       three octets {@code <low>.<middle>.<high>} ({@code 0..255} each).</li>
 * </ul>
 *
 * <p><strong>Octet order:</strong> the octet-wise forms are written
 * least-significant octet first, matching the little-endian order the
 * octets appear in on the wire (see {@code asduAddressField} /
 * {@code InformationObject.address} in the mspec, both
 * {@code LITTLE_ENDIAN}). So {@code 1/2} as an ASDU address is
 * {@code 1 + 2 * 256 = 513}, and {@code 3.4.5} as an information object
 * address is {@code 3 + 4 * 256 + 5 * 65536 = 328707}. Note that this is the
 * <em>opposite</em> order of the identically shaped {@code a/b/c} address of
 * the KNX driver, where the most significant part comes first.
 *
 * <p>Surrounding whitespace is trimmed and leading zeros are accepted, but
 * only up to the number of digits the maximum of that component has: three
 * for an octet ({@code 255}), five for a whole ASDU address ({@code 65535})
 * and eight for a whole information object address ({@code 16777215}). So
 * {@code 007/002.03.004} parses while {@code 000007/2} does not. The Go driver
 * caps the digits the same way, so both languages accept the same addresses.
 *
 * <p>{@code *} is accepted in place of any number and makes that position a
 * wildcard, which is what lets a subscription cover a whole station or a
 * whole octet range. Examples:
 * <pre>
 *   1/2            ASDU 1, information object 2
 *   65535/16777215 the widest addresses both fields can carry
 *   1/2/3          ASDU 513 (low 1, high 2), information object 3
 *   1/2.3.4        ASDU 1, information object 262914 (low 2, middle 3, high 4)
 *   1/2/3.4.5      ASDU 513, information object 328707
 *   &#42;/&#42;            every information object of every station
 *   1/*            every information object of ASDU 1
 *   &#42;/3.4.5        information object 328707 of every station
 *   1/2/3.*.5      low and high octet pinned, middle octet wildcarded
 * </pre>
 */
public class Iec608705104Tag implements PlcTag {

    /** Value reported by {@link #getAdsuAddress()} / {@link #getObjectAddress()} for a wildcarded address. */
    public static final int WILDCARD_ADDRESS = -1;

    /** Widest value the 2 octet common ASDU address can carry. */
    public static final int MAX_ADSU_ADDRESS = 0xFFFF;

    /** Widest value the 3 octet information object address can carry. */
    public static final int MAX_OBJECT_ADDRESS = 0xFFFFFF;

    private static final int MAX_OCTET = 0xFF;
    private static final String WILDCARD = "*";
    private static final String NUMBER_OR_WILDCARD_OCTET = "\\d{1,3}|\\" + WILDCARD;

    private static final String ADSU_SINGLE =
        "(?<adsuAddress>\\d{1,5}|\\*)";
    private static final String ADSU_OCTETS =
        "(?<adsuAddressLowOctet>" + NUMBER_OR_WILDCARD_OCTET + ")/(?<adsuAddressHighOctet>" + NUMBER_OR_WILDCARD_OCTET + ")";
    private static final String OBJECT_SINGLE =
        "(?<objectAddress>\\d{1,8}|\\*)";
    private static final String OBJECT_OCTETS =
        "(?<objectAddressLowOctet>" + NUMBER_OR_WILDCARD_OCTET + ")" +
            "\\.(?<objectAddressMiddleOctet>" + NUMBER_OR_WILDCARD_OCTET + ")" +
            "\\.(?<objectAddressHighOctet>" + NUMBER_OR_WILDCARD_OCTET + ")";

    private static final Pattern ADDRESS_ADSU_OCTETS_OBJECT_OCTETS =
        Pattern.compile("^" + ADSU_OCTETS + "/" + OBJECT_OCTETS + "$");
    private static final Pattern ADDRESS_ADSU_OCTETS_OBJECT_SINGLE =
        Pattern.compile("^" + ADSU_OCTETS + "/" + OBJECT_SINGLE + "$");
    private static final Pattern ADDRESS_ADSU_SINGLE_OBJECT_OCTETS =
        Pattern.compile("^" + ADSU_SINGLE + "/" + OBJECT_OCTETS + "$");
    private static final Pattern ADDRESS_ADSU_SINGLE_OBJECT_SINGLE =
        Pattern.compile("^" + ADSU_SINGLE + "/" + OBJECT_SINGLE + "$");

    private final int adsuAddress;
    private final int adsuAddressMask;
    private final int objectAddress;
    private final int objectAddressMask;
    private final String addressString;

    /**
     * Fully specified (wildcard-free) address, as produced for every
     * information object arriving from the wire.
     */
    public Iec608705104Tag(int adsuAddress, int objectAddress) {
        if (adsuAddress < 0 || adsuAddress > MAX_ADSU_ADDRESS) {
            throw new PlcInvalidTagException("ASDU address " + adsuAddress + " out of range [0.." + MAX_ADSU_ADDRESS + "]");
        }
        if (objectAddress < 0 || objectAddress > MAX_OBJECT_ADDRESS) {
            throw new PlcInvalidTagException("Information object address " + objectAddress + " out of range [0.." + MAX_OBJECT_ADDRESS + "]");
        }
        this.adsuAddress = adsuAddress;
        this.adsuAddressMask = MAX_ADSU_ADDRESS;
        this.objectAddress = objectAddress;
        this.objectAddressMask = MAX_OBJECT_ADDRESS;
        this.addressString = adsuAddress + "/" + objectAddress;
    }

    private Iec608705104Tag(int adsuAddress, int adsuAddressMask,
                            int objectAddress, int objectAddressMask,
                            String addressString) {
        this.adsuAddress = adsuAddress;
        this.adsuAddressMask = adsuAddressMask;
        this.objectAddress = objectAddress;
        this.objectAddressMask = objectAddressMask;
        this.addressString = addressString;
    }

    /**
     * @return {@code true} if {@code tagString} is an address {@link #of(String)}
     * accepts, i.e. one that is both well-formed and in range. Meant as a guard
     * for {@code of(String)}: whenever this returns {@code true}, parsing the
     * same string succeeds.
     */
    public static boolean isValidAddress(String tagString) {
        try {
            of(tagString);
            return true;
        } catch (PlcInvalidTagException e) {
            return false;
        }
    }

    /**
     * @return {@code true} if {@code tagString} has the shape of an address,
     * ignoring whether its components are in range. Only useful for telling a
     * malformed address ({@code "1/abc"}) from an out-of-range one
     * ({@code "65536/0"}) — callers wanting to know whether an address is
     * usable want {@link #isValidAddress(String)}.
     */
    static boolean isSyntacticallyValid(String tagString) {
        if (tagString == null) {
            return false;
        }
        String address = tagString.trim();
        return ADDRESS_ADSU_OCTETS_OBJECT_OCTETS.matcher(address).matches()
            || ADDRESS_ADSU_OCTETS_OBJECT_SINGLE.matcher(address).matches()
            || ADDRESS_ADSU_SINGLE_OBJECT_OCTETS.matcher(address).matches()
            || ADDRESS_ADSU_SINGLE_OBJECT_SINGLE.matcher(address).matches();
    }

    /**
     * Parses an address in any of the accepted forms.
     *
     * @throws PlcInvalidTagException if the address doesn't match the syntax
     *                                or a component is out of range.
     */
    public static Iec608705104Tag of(String tagString) {
        if (tagString == null) {
            throw new PlcInvalidTagException("null");
        }
        String address = tagString.trim();

        // Try the octet-wise forms first (they carry more slashes / dots), so
        // that "1/2/3" is read as ASDU 1/2 rather than being rejected.
        Matcher matcher = ADDRESS_ADSU_OCTETS_OBJECT_OCTETS.matcher(address);
        if (matcher.matches()) {
            return of(address, adsuOctets(matcher), objectOctets(matcher));
        }
        matcher = ADDRESS_ADSU_OCTETS_OBJECT_SINGLE.matcher(address);
        if (matcher.matches()) {
            return of(address, adsuOctets(matcher), new String[]{matcher.group("objectAddress")});
        }
        matcher = ADDRESS_ADSU_SINGLE_OBJECT_OCTETS.matcher(address);
        if (matcher.matches()) {
            return of(address, new String[]{matcher.group("adsuAddress")}, objectOctets(matcher));
        }
        matcher = ADDRESS_ADSU_SINGLE_OBJECT_SINGLE.matcher(address);
        if (matcher.matches()) {
            return of(address, new String[]{matcher.group("adsuAddress")},
                new String[]{matcher.group("objectAddress")});
        }
        throw new PlcInvalidTagException(tagString);
    }

    private static String[] adsuOctets(Matcher matcher) {
        return new String[]{matcher.group("adsuAddressLowOctet"), matcher.group("adsuAddressHighOctet")};
    }

    private static String[] objectOctets(Matcher matcher) {
        return new String[]{
            matcher.group("objectAddressLowOctet"),
            matcher.group("objectAddressMiddleOctet"),
            matcher.group("objectAddressHighOctet")};
    }

    private static Iec608705104Tag of(String address, String[] adsuParts, String[] objectParts) {
        int adsuMask = mask(adsuParts, MAX_ADSU_ADDRESS);
        int objectMask = mask(objectParts, MAX_OBJECT_ADDRESS);
        return new Iec608705104Tag(
            value(address, adsuParts, MAX_ADSU_ADDRESS),
            adsuMask,
            value(address, objectParts, MAX_OBJECT_ADDRESS),
            objectMask,
            canonicalize(adsuParts, "/") + "/" + canonicalize(objectParts, "."));
    }

    /**
     * Folds the components of one address field into a single integer. A
     * single component is the whole value, multiple components are octets in
     * least-significant-first order. Wildcarded components contribute zero
     * bits (which is why the mask is needed to tell them from a literal 0).
     */
    private static int value(String address, String[] parts, int max) {
        if (parts.length == 1) {
            return WILDCARD.equals(parts[0]) ? 0 : number(address, parts[0], max);
        }
        int value = 0;
        for (int i = 0; i < parts.length; i++) {
            if (!WILDCARD.equals(parts[i])) {
                value |= number(address, parts[i], MAX_OCTET) << (8 * i);
            }
        }
        return value;
    }

    /**
     * Bits of an address field that are actually pinned down: everything but
     * the wildcarded components.
     */
    private static int mask(String[] parts, int max) {
        if (parts.length == 1) {
            return WILDCARD.equals(parts[0]) ? 0 : max;
        }
        int mask = 0;
        for (int i = 0; i < parts.length; i++) {
            if (!WILDCARD.equals(parts[i])) {
                mask |= MAX_OCTET << (8 * i);
            }
        }
        return mask;
    }

    private static int number(String address, String part, int max) {
        int value;
        try {
            value = Integer.parseInt(part);
        } catch (NumberFormatException e) {
            throw new PlcInvalidTagException(address, e);
        }
        if (value > max) {
            throw new PlcInvalidTagException(
                address + " (component " + part + " exceeds the maximum of " + max + ")");
        }
        return value;
    }

    /** Normalizes the components (dropping leading zeros) without changing the form. */
    private static String canonicalize(String[] parts, String separator) {
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < parts.length; i++) {
            if (i > 0) {
                sb.append(separator);
            }
            sb.append(WILDCARD.equals(parts[i]) ? WILDCARD : Integer.toString(Integer.parseInt(parts[i])));
        }
        return sb.toString();
    }

    /**
     * @return the common ASDU address, or {@link #WILDCARD_ADDRESS} if any of
     * its octets is wildcarded.
     */
    public int getAdsuAddress() {
        return isAdsuAddressWildcarded() ? WILDCARD_ADDRESS : adsuAddress;
    }

    /**
     * @return the information object address, or {@link #WILDCARD_ADDRESS} if
     * any of its octets is wildcarded.
     */
    public int getObjectAddress() {
        return isObjectAddressWildcarded() ? WILDCARD_ADDRESS : objectAddress;
    }

    /** @return the bits of the ASDU address this tag pins down. */
    public int getAdsuAddressMask() {
        return adsuAddressMask;
    }

    /** @return the bits of the information object address this tag pins down. */
    public int getObjectAddressMask() {
        return objectAddressMask;
    }

    public boolean isAdsuAddressWildcarded() {
        return adsuAddressMask != MAX_ADSU_ADDRESS;
    }

    public boolean isObjectAddressWildcarded() {
        return objectAddressMask != MAX_OBJECT_ADDRESS;
    }

    public boolean isWildcarded() {
        return isAdsuAddressWildcarded() || isObjectAddressWildcarded();
    }

    /**
     * Tests an information object arriving from the wire against this tag.
     * Wildcarded octets always match.
     */
    public boolean matches(int adsuAddress, int objectAddress) {
        return ((adsuAddress ^ this.adsuAddress) & adsuAddressMask) == 0
            && ((objectAddress ^ this.objectAddress) & objectAddressMask) == 0;
    }

    /**
     * @return the address in the form it was parsed from, with leading zeros
     * removed. {@code of(tag.getAddressString())} always yields a tag equal
     * to {@code tag}.
     */
    @Override
    public String getAddressString() {
        return addressString;
    }

    /**
     * IEC 60870-5-104 addresses carry no type information — the datatype of
     * an information object follows from the type identification of the ASDU
     * that delivers it, which is only known at runtime.
     */
    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.NULL;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return PlcTag.super.getArrayInfo();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof Iec608705104Tag other)) {
            return false;
        }
        // Two tags addressing the same set of information objects are equal,
        // no matter which of the equivalent notations they were written in.
        return adsuAddress == other.adsuAddress
            && adsuAddressMask == other.adsuAddressMask
            && objectAddress == other.objectAddress
            && objectAddressMask == other.objectAddressMask;
    }

    @Override
    public int hashCode() {
        return Objects.hash(adsuAddress, adsuAddressMask, objectAddress, objectAddressMask);
    }

    @Override
    public String toString() {
        return "Iec608705104Tag{address='" + addressString + "'}";
    }

}
