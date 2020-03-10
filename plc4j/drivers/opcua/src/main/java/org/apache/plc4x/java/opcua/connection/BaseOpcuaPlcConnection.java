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
package org.apache.plc4x.java.opcua.connection;

import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.base.connection.AbstractPlcConnection;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.opcua.protocol.OpcuaPlcFieldHandler;
import org.eclipse.milo.opcua.stack.core.security.SecurityPolicy;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Created by Matthias Milan Strljic on 10.05.2019
 */
public abstract class BaseOpcuaPlcConnection extends AbstractPlcConnection implements PlcReader, PlcWriter, PlcSubscriber {

    private static final Logger logger = LoggerFactory.getLogger(BaseOpcuaPlcConnection.class);
    protected boolean skipDiscovery = false;
    protected String username = null;
    protected String password = null;
    protected String certFile = null;
    protected SecurityPolicy securityPolicy = null;
    protected String keyStoreFile = null;
//    protected String keyStorePassword;

    /**
     * @param params
     */
    BaseOpcuaPlcConnection(String params) {
        if (!StringUtils.isEmpty(params)) {
            for (String param : params.split("&")) {
                String[] paramElements = param.split("=");
                String paramName = paramElements[0];
                if (paramElements.length == 2) {
                    String paramValue = paramElements[1];
                    switch (paramName) {
                        case "discovery":
                            skipDiscovery = !Boolean.parseBoolean(paramValue);
                            logger.debug("Found Parameter 'skipDiscovery' with value {}", this.skipDiscovery);
                            break;
                        case "username":
                            username = paramValue;
                            logger.debug("Found Parameter 'username' with value {}", username);
                            break;
                        case "password":
                            password = paramValue;
                            logger.debug("Found Parameter 'password' with value {}", password);
                            break;
                        case "certFile":
                            certFile = paramValue;
                            logger.debug("Found Parameter 'certFile' with value {}", certFile);
                            break;
                        case "securityPolicy":
                            logger.debug("Got value for security policy: '{}', trying to parse", paramValue);
                            try {
                                securityPolicy = SecurityPolicy.valueOf(paramValue);
                                logger.debug("Using Security Policy {}", securityPolicy);
                            } catch (IllegalArgumentException e) {
                                logger.warn("Unable to parse policy {}", paramValue);
                            }
                            break;
                        case "keyStoreFile":
                            keyStoreFile = paramValue;
                            logger.debug("Found Parameter 'keyStoreFile' with value {}", keyStoreFile);
                            break;
//                        case "keyStorePassword":
//                            keyStorePassword = paramValue;
//                            logger.debug("Found Parameter 'keyStorePassword' with value {}", keyStorePassword);
//                            break;
                        default:
                            logger.debug("Unknown parameter {} with value {}", paramName, paramValue);
                    }
                } else {
                    logger.debug("Unknown no-value parameter {}", paramName);
                }
            }
        }
    }

    @Override
    public boolean canRead() {
        return true;
    }

    @Override
    public boolean canWrite() {
        return true;
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        return new DefaultPlcReadRequest.Builder(this, new OpcuaPlcFieldHandler());
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        return new DefaultPlcWriteRequest.Builder(this, new OpcuaPlcFieldHandler());
    }

    @Override
    public boolean canSubscribe() {
        return true;
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        return new DefaultPlcSubscriptionRequest.Builder(this, new OpcuaPlcFieldHandler());
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        return new DefaultPlcUnsubscriptionRequest.Builder(this);
    }

    public boolean isSkipDiscovery() {
        return skipDiscovery;
    }
}
