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
import org.apache.commons.scxml2.model.NodeListValue;
import org.apache.commons.scxml2.model.NodeValue;
import org.apache.commons.scxml2.model.ParsedValue;
import org.jdom2.Document;
import org.w3c.dom.Element;
import org.w3c.dom.Node;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ReceiveExtractVerifyAction extends ReceiveAction {

    private final Map<String, String> valueVerificationRules;
    private final Map<String, String> contextVerificationRules;
    private final Map<String, String> extractionRules;

    public ReceiveExtractVerifyAction() {
        valueVerificationRules = new HashMap<>();
        contextVerificationRules = new HashMap<>();
        extractionRules = new HashMap<>();
    }

    @Override
    @SuppressWarnings("unchecked")
    public void setParsedValue(ParsedValue parsedValue) {
        super.setParsedValue(parsedValue);
        if(parsedValue != null) {
            if(parsedValue instanceof NodeListValue) {
                List<Node> ruleList = (List<Node>) parsedValue.getValue();
                for (Node node : ruleList) {
                    if(node instanceof Element) {
                        parseElement((Element) node);
                    }
                }
            } else if(parsedValue instanceof NodeValue) {
                parseElement((Element) parsedValue.getValue());
            }
        }
    }

    private void parseElement(Element ruleElement) {
        String expression = ruleElement.getAttribute("xpath-expression");
        if ("verification".equals(ruleElement.getTagName())) {
            if(ruleElement.hasAttribute("name")) {
                String name = ruleElement.getAttribute("name");
                contextVerificationRules.put(expression, name);
            } else if(ruleElement.hasAttribute("value")) {
                String value = ruleElement.getAttribute("value");
                valueVerificationRules.put(expression, value);
            }
        } else if ("extraction".equals(ruleElement.getTagName())) {
            String name = ruleElement.getAttribute("name");
            extractionRules.put(expression, name);
        } else {
            getLogger().error("unsupported rule type: " + ruleElement.getTagName());
        }
    }


    protected void processMessage(Document message, ActionExecutionContext ctx) {
        // First verify all verification conditions.
        for (Map.Entry<String, String> rule : contextVerificationRules.entrySet()) {
            Object reference = ctx.getGlobalContext().get(rule.getValue());
            verifyValue(message, rule.getKey(), reference);
        }
        for (Map.Entry<String, String> rule : valueVerificationRules.entrySet()) {
            Object reference = rule.getValue();
            verifyValue(message, rule.getKey(), reference);
        }

        // Then extract data from the document.
        for (Map.Entry<String, String> rule : extractionRules.entrySet()) {
            String current = getRuleText(message, rule.getKey());
            if(current == null) {
                throw new RuntimeException("Error extracting. Got null value");
            }
            ctx.getGlobalContext().set(rule.getValue(), current);
        }
    }

    private void verifyValue(Document message, String xpath, Object reference) {
        String current = getRuleText(message, xpath);
        if(current == null) {
            throw new RuntimeException("Error verifying. Expected: " + reference.toString() + " got null value");
        }
        if(!current.equals(reference)) {
            throw new RuntimeException("Error verifying. Expected: " + reference.toString() + " got: " + current);
        }
    }

}
