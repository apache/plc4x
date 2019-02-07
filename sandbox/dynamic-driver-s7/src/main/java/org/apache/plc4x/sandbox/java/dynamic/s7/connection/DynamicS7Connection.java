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

package org.apache.plc4x.sandbox.java.dynamic.s7.connection;

import org.apache.commons.lang3.StringUtils;
import org.apache.commons.scxml2.model.CustomAction;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedOperationException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.sandbox.java.dynamic.connection.DynamicDriverConnectionBase;
import org.apache.plc4x.sandbox.java.dynamic.s7.actions.S7DecodeArticleNumber;
import org.apache.plc4x.sandbox.java.dynamic.s7.utils.S7TsapIdEncoder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

public class DynamicS7Connection extends DynamicDriverConnectionBase {

    private static final Logger logger = LoggerFactory.getLogger(DynamicS7Connection.class);

    private final InetAddress address;
    private final short calledTsapId;
    private final short callingTsapId;

    private short paramPduSize;
    private short paramMaxAmqCaller;
    private short paramMaxAmqCallee;
    private String paramControllerType;

    public DynamicS7Connection(InetAddress address, int rack, int slot, String params) {
        super("org/apache/plc4x/protocols/s7/protocol.scxml.xml",
            "org/apache/plc4x/protocols/s7/protocol.dfdl.xsd");

        this.address = address;
        this.calledTsapId = S7TsapIdEncoder.encodeS7TsapId((byte) 0x02, 0, 0);
        this.callingTsapId = S7TsapIdEncoder.encodeS7TsapId((byte) 0x01, rack, slot);

        short curParamPduSize = 1024;
        short curParamMaxAmqCaller = 8;
        short curParamMaxAmqCallee = 8;
        String curParamControllerType = null;

        if (!StringUtils.isEmpty(params)) {
            for (String param : params.split("&")) {
                String[] paramElements = param.split("=");
                String paramName = paramElements[0];
                if (paramElements.length == 2) {
                    String paramValue = paramElements[1];
                    switch (paramName) {
                        case "pdu-size":
                            curParamPduSize = Short.parseShort(paramValue);
                            break;
                        case "max-amq-caller":
                            curParamMaxAmqCaller = Short.parseShort(paramValue);
                            break;
                        case "max-amq-callee":
                            curParamMaxAmqCallee = Short.parseShort(paramValue);
                            break;
                        case "controller-type":
                            curParamControllerType = paramValue;
                            break;
                        default:
                            logger.debug("Unknown parameter {} with value {}", paramName, paramValue);
                    }
                } else {
                    logger.debug("Unknown no-value parameter {}", paramName);
                }
            }
        }

        // It seems that the LOGO devices are a little picky about the pdu-size.
        // Instead of handling this out, they just hang up without any error message.
        // So in case of a LOGO controller, set this to a known working value.
        if((curParamControllerType != null) && curParamControllerType.equalsIgnoreCase("logo") && curParamPduSize == 1024) {
            curParamPduSize = 480;
        }

        // IsoTP uses pre defined sizes. Find the smallest box,
        // that would be able to contain the requested pdu size.
        this.paramPduSize = curParamPduSize;
        this.paramMaxAmqCaller = curParamMaxAmqCaller;
        this.paramMaxAmqCallee = curParamMaxAmqCallee;
        this.paramControllerType = curParamControllerType;
    }

    @Override
    protected String getConnectedStateName() {
        return "connected";
    }

    @Override
    protected String getDisconnectTransitionName() {
        return "disconnect";
    }

    @Override
    protected Collection<CustomAction> getAdditionalCustomActions() {
        return Collections.singleton(
            new CustomAction("https://plc4x.apache.org/scxml-extension", "S7DecodeArticleNumber",
                S7DecodeArticleNumber.class));
    }

    @Override
    protected Map<String, Object> getAdditionalContextDataItems() {
        Map<String, Object> dataItems = new HashMap<>();

        dataItems.put("hostname", address.getHostAddress());
        dataItems.put("port", "102");
        dataItems.put("plcType", paramControllerType);

        dataItems.put("cotpLocalReference", "15");
        dataItems.put("cotpCalledTsap", Short.toString(calledTsapId));
        dataItems.put("cotpCallingTsap", Short.toString(callingTsapId));
        dataItems.put("cotpTpduSize", "10");

        dataItems.put("s7PduLength", Short.toString(paramPduSize));
        dataItems.put("s7MaxAmqCaller", Short.toString(paramMaxAmqCaller));
        dataItems.put("s7MaxAmqCallee", Short.toString(paramMaxAmqCallee));

        return dataItems;
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        throw new PlcUnsupportedOperationException("The connection does not support reading");
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        throw new PlcUnsupportedOperationException("The connection does not support writing");
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        throw new PlcUnsupportedOperationException("The connection does not support subscription");
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        throw new PlcUnsupportedOperationException("The connection does not support subscription");
    }

}
