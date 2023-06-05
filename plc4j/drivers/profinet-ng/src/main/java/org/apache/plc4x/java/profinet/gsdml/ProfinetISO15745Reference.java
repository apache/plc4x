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

package org.apache.plc4x.java.profinet.gsdml;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonRootName;

@JsonIgnoreProperties(ignoreUnknown = true)
@JsonRootName("ISO15745Reference")
public class ProfinetISO15745Reference {

    @JsonProperty("ISO15745Part")
    private int iso15745Part;

    @JsonProperty("ISO15745Edition")
    private int iso15745Edition;

    @JsonProperty("ProfileTechnology")
    private String profileTechnology;

    public int getIso15745Part() {
        return iso15745Part;
    }

    public int getIso15745Edition() {
        return iso15745Edition;
    }

    public String getProfileTechnology() {
        return profileTechnology;
    }
}
