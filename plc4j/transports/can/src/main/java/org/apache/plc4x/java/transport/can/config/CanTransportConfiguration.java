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
package org.apache.plc4x.java.transport.can.config;

import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.transport.can.CanIdFilter;

/**
 * Base configuration class for CAN transports, providing common CAN ID filter parameters.
 * <p>
 * Concrete CAN transport configurations (SocketCAN, VirtualCAN) extend this class
 * to inherit filter support while adding transport-specific parameters.
 * <p>
 * Filter configuration allows transport instances to receive only CAN frames
 * with matching identifiers, reducing processing overhead.
 */
public class CanTransportConfiguration implements TransportConfiguration {

    /**
     * Comma-separated list of accepted CAN IDs (decimal or 0x hex notation).
     * Example: "0x100,0x200,0x300" or "256,512,768"
     * If empty, no ID-based filtering is applied.
     */
    @ConfigurationParameter("filter-ids")
    @Description("Comma-separated list of accepted CAN IDs (decimal or 0x hex). Empty means accept all.")
    public String filterIds;

    /**
     * Start of accepted CAN ID range (inclusive). -1 means no range filtering.
     */
    @ConfigurationParameter("filter-range-start")
    @Description("Start of accepted CAN ID range (inclusive). -1 means no range filtering.")
    @IntDefaultValue(-1)
    public int filterRangeStart = -1;

    /**
     * End of accepted CAN ID range (inclusive). -1 means no range filtering.
     */
    @ConfigurationParameter("filter-range-end")
    @Description("End of accepted CAN ID range (inclusive). -1 means no range filtering.")
    @IntDefaultValue(-1)
    public int filterRangeEnd = -1;

    /**
     * Builds a {@link org.apache.plc4x.java.transport.can.CanIdFilter} from the configured filter parameters.
     * <p>
     * Combines explicit IDs from {@link #filterIds} and the range from
     * {@link #filterRangeStart}/{@link #filterRangeEnd} into a single filter.
     * If no filter parameters are configured, returns an accept-all filter.
     *
     * @return the constructed filter
     * @throws IllegalArgumentException if any filter parameter value is invalid
     */
    public CanIdFilter buildFilter() {
        boolean hasIds = filterIds != null && !filterIds.trim().isEmpty();
        boolean hasRange = filterRangeStart >= 0 && filterRangeEnd >= 0;

        if (!hasIds && !hasRange) {
            return CanIdFilter.acceptAll();
        }

        CanIdFilter.Builder builder = CanIdFilter.builder();

        if (hasIds) {
            String[] parts = filterIds.split(",");
            for (String part : parts) {
                String trimmed = part.trim();
                if (!trimmed.isEmpty()) {
                    int id = parseCanId(trimmed);
                    builder.addId(id);
                }
            }
        }

        if (hasRange) {
            builder.addRange(filterRangeStart, filterRangeEnd);
        }

        return builder.build();
    }

    /**
     * Parses a CAN ID string that may be in decimal or hexadecimal (0x prefix) notation.
     *
     * @param idStr the string to parse
     * @return the parsed integer value
     * @throws IllegalArgumentException if the string cannot be parsed
     */
    private static int parseCanId(String idStr) {
        try {
            if (idStr.startsWith("0x") || idStr.startsWith("0X")) {
                return Integer.parseInt(idStr.substring(2), 16);
            }
            return Integer.parseInt(idStr);
        } catch (NumberFormatException e) {
            throw new IllegalArgumentException("Invalid CAN ID: " + idStr, e);
        }
    }
}
