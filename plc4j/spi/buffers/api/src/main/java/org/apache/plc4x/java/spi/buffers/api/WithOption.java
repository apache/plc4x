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
package org.apache.plc4x.java.spi.buffers.api;

import java.util.Optional;

public interface WithOption {
    /**
     * @return true if the option should also be valid in sub-contexts.
     */
    boolean isSticky();

    static WithOption[] AddOptions(WithOption[] options, WithOption... newOptions) {
        WithOption[] newOptionsArray = new WithOption[newOptions.length + options.length];
        System.arraycopy(newOptions, 0, newOptionsArray, 0, newOptions.length);
        System.arraycopy(options, 0, newOptionsArray, newOptions.length, options.length);
        return newOptionsArray;
    }

    static WithOption[] UpdateOptions(WithOption[] options, WithOption... updateOptions) {
        WithOption[] newOptionsArray = new WithOption[options.length];
        // First, simply copy over all options.
        System.arraycopy(options, 0, newOptionsArray, 0, options.length);
        // Then update existing options.
        for (WithOption updateOption : updateOptions) {
            for (int i = 0; i < newOptionsArray.length; i++) {
                WithOption option = newOptionsArray[i];
                if (option.getClass() == updateOption.getClass()) {
                    newOptionsArray[i] = updateOption;
                    break;
                }
            }
        }
        return newOptionsArray;
    }

    static WithOption WithName(String name) {
        return (withOptionName) () -> name;
    }

    /**
     * This sets the encoding for all types of values. Use the With*Encoding methods to set encoding for specific types.
     *
     * @param encodingName name of the encoding to use for all operations.
     * @return option
     */
    static WithOption WithEncoding(String encodingName) {
        return (withOptionEncoding) () -> encodingName;
    }

    /**
     * This sets the encoding for reading unsigned integers.
     *
     * @param encodingName name of the encoding to use for unsigned integer operations.
     * @return option
     */
    static WithOption WithUnsignedIntegerEncoding(String encodingName) {
        return (withOptionUnsignedIntegerEncoding) () -> encodingName;
    }

    /**
     * This sets the encoding for reading signed integers.
     *
     * @param encodingName name of the encoding to use for signed integer operations.
     * @return option
     */
    static WithOption WithSignedIntegerEncoding(String encodingName) {
        return (withOptionSignedIntegerEncoding) () -> encodingName;
    }

    /**
     * This sets the encoding for reading floating point values.
     *
     * @param encodingName name of the encoding to use for floating point value operations.
     * @return option
     */
    static WithOption WithFloatEncoding(String encodingName) {
        return (withOptionFloatEncoding) () -> encodingName;
    }

    /**
     * This sets the encoding for reading strings.
     *
     * @param encodingName name of the encoding to use for string operations.
     * @return option
     */
    static WithOption WithStringEncoding(String encodingName) {
        return (withOptionStringEncoding) () -> encodingName;
    }

    static WithOption WithAdditionalStringRepresentation(String stringRepresentation) {
        return (withOptionAdditionalStringRepresentation) () -> stringRepresentation;
    }

    static WithOption WithRenderAsList(boolean renderAsList) {
        return (withOptionRenderAsList) () -> renderAsList;
    }

    static Optional<String> extractName(WithOption[] options) {
        if (options != null) {
            for (WithOption option : options) {
                if (option instanceof withOptionName) {
                    return Optional.of(((withOptionName) option).name());
                }
            }
        }
        return Optional.empty();
    }

    static Optional<String> extractName(WithOption[] options, WithOption[] defaultOptions) {
        Optional<String> byteOrder = extractName(options);
        if (byteOrder.isPresent()) {
            return byteOrder;
        }
        return extractName(defaultOptions);
    }

    static Optional<String> extractEncoding(WithOption[] options) {
        if (options != null) {
            for (WithOption option : options) {
                if (option instanceof withOptionEncoding) {
                    return Optional.of(((withOptionEncoding) option).encodingName());
                }
            }
        }
        return Optional.empty();
    }

    static Optional<String> extractEncoding(WithOption[] options, WithOption[] defaultOptions) {
        Optional<String> encoding = extractEncoding(options);
        if (encoding.isPresent()) {
            return encoding;
        }
        return extractEncoding(defaultOptions);
    }

    static Optional<String> extractUnsignedIntegerEncoding(WithOption[] options) {
        if (options != null) {
            for (WithOption option : options) {
                if (option instanceof withOptionUnsignedIntegerEncoding) {
                    return Optional.of(((withOptionUnsignedIntegerEncoding) option).encodingName());
                }
            }
        }
        return Optional.empty();
    }

    static Optional<String> extractUnsignedIntegerEncoding(WithOption[] options, WithOption[] defaultOptions) {
        Optional<String> encoding = extractUnsignedIntegerEncoding(options);
        if (encoding.isPresent()) {
            return encoding;
        }
        return extractUnsignedIntegerEncoding(defaultOptions);
    }

    static Optional<String> extractSignedIntegerEncoding(WithOption[] options) {
        if (options != null) {
            for (WithOption option : options) {
                if (option instanceof withOptionSignedIntegerEncoding) {
                    return Optional.of(((withOptionSignedIntegerEncoding) option).encodingName());
                }
            }
        }
        return Optional.empty();
    }

    static Optional<String> extractSignedIntegerEncoding(WithOption[] options, WithOption[] defaultOptions) {
        Optional<String> encoding = extractSignedIntegerEncoding(options);
        if (encoding.isPresent()) {
            return encoding;
        }
        return extractSignedIntegerEncoding(defaultOptions);
    }

    static Optional<String> extractFloatEncoding(WithOption[] options) {
        if (options != null) {
            for (WithOption option : options) {
                if (option instanceof withOptionFloatEncoding) {
                    return Optional.of(((withOptionFloatEncoding) option).encodingName());
                }
            }
        }
        return Optional.empty();
    }

    static Optional<String> extractFloatEncoding(WithOption[] options, WithOption[] defaultOptions) {
        Optional<String> encoding = extractFloatEncoding(options);
        if (encoding.isPresent()) {
            return encoding;
        }
        return extractFloatEncoding(defaultOptions);
    }

    static Optional<String> extractStringEncoding(WithOption[] options) {
        if (options != null) {
            for (WithOption option : options) {
                if (option instanceof withOptionStringEncoding) {
                    return Optional.of(((withOptionStringEncoding) option).encodingName());
                }
            }
        }
        return Optional.empty();
    }

    static Optional<String> extractStringEncoding(WithOption[] options, WithOption[] defaultOptions) {
        Optional<String> encoding = extractStringEncoding(options);
        if (encoding.isPresent()) {
            return encoding;
        }
        return extractStringEncoding(defaultOptions);
    }

    static Optional<String> extractAdditionalStringRepresentation(WithOption[] options) {
        if (options != null) {
            for (WithOption option : options) {
                if (option instanceof withOptionAdditionalStringRepresentation) {
                    return Optional.of(((withOptionAdditionalStringRepresentation) option).stringRepresentation());
                }
            }
        }
        return Optional.empty();
    }

    static Optional<Boolean> extractRenderAsList(WithOption[] options) {
        if (options != null) {
            for (WithOption option : options) {
                if (option instanceof withOptionRenderAsList) {
                    return Optional.of(((withOptionRenderAsList) option).renderAsList());
                }
            }
        }
        return Optional.empty();
    }

}

interface withOptionName extends WithOption {
    String name();

    default boolean isSticky() {
        return false;
    }
}

interface withOptionEncoding extends WithOption {
    String encodingName();

    default boolean isSticky() {
        return true;
    }
}

interface withOptionUnsignedIntegerEncoding extends WithOption {
    String encodingName();

    default boolean isSticky() {
        return true;
    }
}

interface withOptionSignedIntegerEncoding extends WithOption {
    String encodingName();

    default boolean isSticky() {
        return true;
    }
}

interface withOptionFloatEncoding extends WithOption {
    String encodingName();

    default boolean isSticky() {
        return true;
    }
}

interface withOptionStringEncoding extends WithOption {
    String encodingName();

    default boolean isSticky() {
        return true;
    }
}

interface withOptionAdditionalStringRepresentation extends WithOption {
    String stringRepresentation();

    default boolean isSticky() {
        return false;
    }
}

interface withOptionRenderAsList extends WithOption {
    boolean renderAsList();

    default boolean isSticky() {
        return false;
    }
}


