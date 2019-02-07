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

package org.apache.plc4x.sandbox.java.dynamic.s7.actions;

import org.apache.commons.scxml2.ActionExecutionContext;
import org.apache.commons.scxml2.model.ActionExecutionError;
import org.apache.plc4x.sandbox.java.dynamic.actions.BasePlc4xAction;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class S7DecodeArticleNumber extends BasePlc4xAction {

    private String articleNumberParameterName;
    private String plcTypeParameterName;

    public String getArticleNumberParameterName() {
        return articleNumberParameterName;
    }

    public void setArticleNumberParameterName(String articleNumberParameterName) {
        this.articleNumberParameterName = articleNumberParameterName;
    }

    public String getPlcTypeParameterName() {
        return plcTypeParameterName;
    }

    public void setPlcTypeParameterName(String plcTypeParameterName) {
        this.plcTypeParameterName = plcTypeParameterName;
    }

    @Override
    protected Logger getLogger() {
        return LoggerFactory.getLogger(S7DecodeArticleNumber.class);
    }

    @Override
    public void execute(ActionExecutionContext ctx) throws ActionExecutionError {
        String articleNumber = ctx.getGlobalContext().get(articleNumberParameterName).toString();
        if(articleNumber == null) {
            fireFailureEvent(ctx, "Couldn't find article number.");
            return;
        }

        String plcType = lookupControllerType(articleNumber);
        if(plcType == null) {
            fireFailureEvent(ctx, "Unknown PLC type for article number: " + articleNumber);
        }

        ctx.getGlobalContext().set(plcTypeParameterName, plcType);
        fireSuccessEvent(ctx);
    }

    private String lookupControllerType(String articleNumber) {
        if(!articleNumber.startsWith("6ES7 ")) {
            return null;
        }

        String model = articleNumber.substring(articleNumber.indexOf(' ') + 1, articleNumber.indexOf(' ') + 2);
        switch (model) {
            case "2":
                return "S7-1200";
            case "5":
                return "S7-1500";
            case "3":
                return "S7-300";
            case "4":
                return "S7-400";
            default:
                return null;
        }
    }

}
