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
package org.apache.plc4x.java.canopen.readwrite.utils;

import org.apache.plc4x.java.canopen.readwrite.SDOInitiateExpeditedUploadResponse;
import org.apache.plc4x.java.canopen.readwrite.SDOInitiateUploadResponsePayload;

public class StaticHelper {

    /**
     * Helper invoked by generated SDO upload/download initiate code to compute the {@code size} field.
     *
     * <p>For an expedited+indicated transfer carrying ≤4 bytes of inline data, {@code size}
     * is the number of <em>unused</em> bytes (i.e. {@code 4 - data.length}). For all other
     * payload shapes the field is zero.</p>
     */
    public static int count(boolean expedited, boolean indicated, SDOInitiateUploadResponsePayload payload) {
        if (expedited && indicated && payload instanceof SDOInitiateExpeditedUploadResponse) {
            return 4 - ((SDOInitiateExpeditedUploadResponse) payload).getData().length;
        }
        return 0;
    }

}
