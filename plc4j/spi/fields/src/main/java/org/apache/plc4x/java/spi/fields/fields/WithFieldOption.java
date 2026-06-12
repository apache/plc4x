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

package org.apache.plc4x.java.spi.fields.fields;

import org.apache.plc4x.java.spi.buffers.api.WithOption;

import java.util.Optional;

public interface WithFieldOption extends WithOption {

    static WithFieldOption WithNullBytesHex(String nullBytesHex) {
        return (withNullBytesHex) () -> nullBytesHex;
    }

    static Optional<String> extractNullBytesHex(WithOption... options) {
        for (WithOption option : options) {
            if (option instanceof withNullBytesHex) {
                return Optional.of(((withNullBytesHex) option).nullBytesHex());
            }
        }
        return Optional.empty();
    }

}

interface withNullBytesHex extends WithFieldOption {
    String nullBytesHex();

    default boolean isSticky() {
        return false;
    }
}
