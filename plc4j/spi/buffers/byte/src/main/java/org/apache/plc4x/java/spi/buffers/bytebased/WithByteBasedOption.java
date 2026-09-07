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
package org.apache.plc4x.java.spi.buffers.bytebased;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.byteorder.ByteOrderManager;
import org.apache.plc4x.java.spi.buffers.bytebased.encoding.EncodingManager;

import java.util.Optional;

public interface WithByteBasedOption extends WithOption {

    static WithByteBasedOption WithByteOrder(String byteOrderName) {
        return (withOptionByteOrder) () -> byteOrderName;
    }

    static WithByteBasedOption WithPaddingChar(String paddingChar) {
        if (paddingChar.length() != 1) {
            throw new IllegalArgumentException("Padding character length should be exactly 1");
        }
        return (withOptionPaddingChar) () -> paddingChar.charAt(0);
    }

    static Optional<ByteOrderManager> extractByteOrderManager(WithOption[] options) {
        if (options != null) {
            for (WithOption option : options) {
                if (option instanceof withOptionByteOrderManager) {
                    return Optional.of(((withOptionByteOrderManager) option).byteOrderManager());
                }
            }
        }
        return Optional.empty();
    }

    static Optional<EncodingManager> extractEncodingManager(WithOption[] options) {
        if (options != null) {
            for (WithOption option : options) {
                if (option instanceof withOptionEncodingManager) {
                    return Optional.of(((withOptionEncodingManager) option).encodingManager());
                }
            }
        }
        return Optional.empty();
    }

    static Optional<String> extractByteOrder(WithOption[] options) {
        if (options != null) {
            for (WithOption option : options) {
                if (option instanceof withOptionByteOrder) {
                    return Optional.of(((withOptionByteOrder) option).byteOrderName());
                }
            }
        }
        return Optional.empty();
    }

    static Optional<String> extractByteOrder(WithOption[] options, WithOption[] defaultOptions) {
        Optional<String> byteOrder = extractByteOrder(options);
        if (byteOrder.isPresent()) {
            return byteOrder;
        }
        return extractByteOrder(defaultOptions);
    }

    static Optional<Character> extractPaddingChar(WithOption[] options) {
        if (options != null) {
            for (WithOption option : options) {
                if (option instanceof withOptionPaddingChar) {
                    return Optional.of(((withOptionPaddingChar) option).paddingChar());
                }
            }
        }
        return Optional.empty();
    }

    static Optional<Character> extractPaddingChar(WithOption[] options, WithOption[] defaultOptions) {
        Optional<Character> byteOrder = extractPaddingChar(options);
        if (byteOrder.isPresent()) {
            return byteOrder;
        }
        return extractPaddingChar(defaultOptions);
    }
}

interface withOptionByteOrderManager extends WithByteBasedOption {
    ByteOrderManager byteOrderManager();

    default boolean isSticky() {
        return true;
    }
}

interface withOptionEncodingManager extends WithByteBasedOption {
    EncodingManager encodingManager();

    default boolean isSticky() {
        return true;
    }
}

interface withOptionByteOrder extends WithByteBasedOption {
    String byteOrderName();

    default boolean isSticky() {
        return true;
    }
}

interface withOptionPaddingChar extends WithByteBasedOption {
    Character paddingChar();

    default boolean isSticky() {
        return true;
    }
}
