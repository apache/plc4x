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

import org.apache.plc4x.java.ads.readwrite.AdsDataTypeArrayInfo;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.List;

/**
 * Result of resolving a symbolic ADS tag against the symbol- and data-type tables.
 *
 * <p>Carries everything needed to issue the AdsRead/AdsWrite and to decode/encode the
 * payload — no further table walking is needed at I/O time.
 *
 * @param indexGroup     ADS index group of the resolved location.
 * @param indexOffset    Absolute byte offset of the resolved location.
 * @param sizeInBytes    Total bytes to read/write at this location.
 * @param dataTypeName   Name of the leaf data type at the resolved location (e.g.
 *                       {@code "DINT"}, {@code "STRING(80)"}, {@code "TSimpleStruct"},
 *                       {@code "ARRAY [1..2,1..3] OF INT"}). For partial-array reads,
 *                       this is still the name of the array type — {@code remainingArrayInfo}
 *                       describes the remaining dimensions.
 * @param plcValueType   PLC4X value type the leaf decodes to. {@code List} when the
 *                       resolved address has remaining array dimensions; {@code Struct}
 *                       when the leaf is a struct; otherwise the scalar leaf type.
 * @param stringLength   Character count for STRING/WSTRING leaves; {@code 0} otherwise.
 * @param remainingArrayInfo Dimensions still un-indexed at the resolved location. Empty
 *                       for fully-indexed leaves; non-empty when the user asked for e.g.
 *                       {@code arr[1]} on a 2-D array, which yields a row.
 */
public record ResolvedAdsTag(
    long indexGroup,
    long indexOffset,
    long sizeInBytes,
    String dataTypeName,
    PlcValueType plcValueType,
    int stringLength,
    List<AdsDataTypeArrayInfo> remainingArrayInfo
) {
}
