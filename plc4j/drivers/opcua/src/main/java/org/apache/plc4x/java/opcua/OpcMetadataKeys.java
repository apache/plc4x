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

package org.apache.plc4x.java.opcua;

import org.apache.plc4x.java.api.metadata.Metadata;
import org.apache.plc4x.java.api.metadata.Metadata.Key;
import org.apache.plc4x.java.opcua.tag.OpcuaQualityStatus;

/**
 * OPC UA level metadata keys.
 */
public interface OpcMetadataKeys {

    Key<OpcuaQualityStatus> QUALITY = Metadata.Key.of("opcua_quality", OpcuaQualityStatus.class);

    Key<Long> SERVER_TIMESTAMP = Metadata.Key.of("opcua_server_timestamp", Long.class);
    Key<Long> SOURCE_TIMESTAMP = Metadata.Key.of("opcua_source_timestamp", Long.class);

}
