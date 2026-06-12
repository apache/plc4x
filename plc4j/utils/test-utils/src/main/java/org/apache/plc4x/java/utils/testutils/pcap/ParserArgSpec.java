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
package org.apache.plc4x.java.utils.testutils.pcap;

/**
 * Specification for a parser argument that is passed to the generated {@code staticParse()} method.
 *
 * <p>Parser arguments can be either fixed (same value for all packets) or direction-dependent
 * (computed from whether a packet is a request or response).
 */
public sealed interface ParserArgSpec permits ParserArgSpec.FixedArg, ParserArgSpec.DirectionDependentArg {

    /**
     * Returns the argument name (matches the parameter name in the mspec).
     */
    String name();

    /**
     * Resolves this argument to a string value for the given packet direction.
     *
     * @param isResponse true if the packet travels from server to client
     * @return the resolved argument value as a string
     */
    String resolve(boolean isResponse);

    /**
     * A parser argument with a constant value regardless of packet direction.
     *
     * @param name  the argument name
     * @param value the constant value
     */
    record FixedArg(String name, String value) implements ParserArgSpec {
        @Override
        public String resolve(boolean isResponse) {
            return value;
        }
    }

    /**
     * A parser argument whose boolean value depends on packet direction.
     * Resolves to {@code "true"} for response packets and {@code "false"} for request packets.
     *
     * @param name the argument name (typically "response")
     */
    record DirectionDependentArg(String name) implements ParserArgSpec {
        @Override
        public String resolve(boolean isResponse) {
            return Boolean.toString(isResponse);
        }
    }
}
