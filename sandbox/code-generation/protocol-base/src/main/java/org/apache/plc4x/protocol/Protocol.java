/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.protocol;

import java.io.InputStream;

public interface Protocol {

    /**
     * The name of the protocol what the plugin will use to select the correct protocol module.
     *
     * @return the name of the protocol.
     */
    String getName();

    /**
     * Returns an InputStream to the spec that defines the message format of packets of the current driver.
     *
     * @return the InputStream for reading the message format spec for the current driver.
     */
    InputStream getMessageFormatSchema();

}
