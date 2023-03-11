/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.apache.plc4x.hop.actions;

import org.apache.hop.metadata.api.HopMetadata;
import org.apache.hop.metadata.api.HopMetadataBase;
import org.apache.hop.metadata.api.IHopMetadata;

/**
 *
 * @author cgarcia
 */
@HopMetadata(
    key = "Xplc4xaction",
    name = "PLC4x Action",
    description = "A shared PLC4x connection to a PLC",
    image = "plc4x_toddy.svg",
    documentationUrl = "/metadata-types/neo4j/neo4j-connection.html")
public class ActionSampleMetaData extends HopMetadataBase implements IHopMetadata {
    
}
