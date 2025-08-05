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
package org.apache.plc4x.java.api.messages;

import org.apache.plc4x.java.api.metadata.Metadata;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.util.Collection;

/**
 * Base type for all response messages sent as response for a prior request
 * from a plc to the plc4x system.
 */
public interface PlcTagResponse extends PlcResponse {

    @Override
    PlcTagRequest getRequest();

    Collection<String> getTagNames();

    PlcTag getTag(String name);

    PlcResponseCode getResponseCode(String name);

    /**
     * Returns tag level metadata information.
     */
    Metadata getTagMetadata(String name);

}
