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
package org.apache.plc4x.java.df1;

import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.df1.configuration.Df1Configuration;
import org.apache.plc4x.java.df1.field.Df1FieldHandler;
import org.apache.plc4x.java.df1.protocol.Df1ProtocolLogic;
import org.apache.plc4x.java.df1.readwrite.DF1Command;
import org.apache.plc4x.java.df1.readwrite.io.DF1CommandIO;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;

public class DF1PlcDriver extends GeneratedDriverBase<DF1Command> {

    @Override
    public String getProtocolCode() {
        return "df1";
    }

    @Override
    public String getProtocolName() {
        return "Allen-Bradley DF1";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return Df1Configuration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "serial";
    }

    @Override
    protected Df1FieldHandler getFieldHandler() {
        return new Df1FieldHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new IEC61131ValueHandler();
    }

    @Override
    protected ProtocolStackConfigurer<DF1Command> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(DF1Command.class, DF1CommandIO.class)
            .withProtocol(Df1ProtocolLogic.class)
            .build();
    }

}
