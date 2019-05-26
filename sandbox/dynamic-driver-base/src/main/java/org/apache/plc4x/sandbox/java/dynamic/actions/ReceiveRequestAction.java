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

package org.apache.plc4x.sandbox.java.dynamic.actions;

import org.apache.commons.scxml2.ActionExecutionContext;
import org.jdom2.Document;

public class ReceiveRequestAction extends ReceiveAction {

    private String idExpression = null;

    public String getIdExpression() {
        return idExpression;
    }

    public void setIdExpression(String idExpression) {
        this.idExpression = idExpression;
    }

    @Override
    public void execute(ActionExecutionContext ctx) {
        if (idExpression == null) {
            fireFailureEvent(ctx, "'id' element not present");
            return;
        }
        super.execute(ctx);
    }

    @Override
    protected void processMessage(Document message, ActionExecutionContext ctx) {
        String requestId = getRuleText(message, idExpression);

    }

}
