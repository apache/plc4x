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
package org.apache.plc4x.java.plc4x;

import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.plc4x.readwrite.Plc4xRequest;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;

public class Plc4xDriver extends GeneratedDriverBase<Plc4xRequest> {

    @Override
    public String getProtocolCode() {
        return "plc4x";
    }

    @Override
    public String getProtocolName() {
        return "PLC4X (Proxy-Protocol)";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return null;
    }

    @Override
    protected PlcFieldHandler getFieldHandler() {
        return null;
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return null;
    }

    @Override
    protected String getDefaultTransport() {
        return "tls";
    }

    @Override
    protected ProtocolStackConfigurer<Plc4xRequest> getStackConfigurer() {
        return null;
    }

}
