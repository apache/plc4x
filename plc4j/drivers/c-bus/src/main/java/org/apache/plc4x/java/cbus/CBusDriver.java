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
package org.apache.plc4x.java.cbus;

import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.cbus.readwrite.CBusCommand;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;

public class CBusDriver extends GeneratedDriverBase<CBusCommand> {
    @Override
    public String getProtocolCode() {
        return null;
    }

    @Override
    public String getProtocolName() {
        return null;
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
        return null;
    }

    @Override
    protected ProtocolStackConfigurer<CBusCommand> getStackConfigurer() {
        return null;
    }
}
