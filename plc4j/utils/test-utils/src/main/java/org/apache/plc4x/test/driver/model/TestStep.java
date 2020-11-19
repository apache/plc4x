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
package org.apache.plc4x.test.driver.model;

import org.dom4j.Element;

import java.util.List;

public class TestStep {

    private final StepType type;
    private final String name;
    private final List<String> parserArguments;
    private final Element payload;

    public TestStep(StepType type, String name, List<String> parserArguments, Element payload) {
        this.type = type;
        this.name = name;
        this.parserArguments = parserArguments;
        this.payload = payload;
    }

    public StepType getType() {
        return type;
    }

    public String getName() {
        return name;
    }

    public List<String> getParserArguments() {
        return parserArguments;
    }

    public Element getPayload() {
        return payload;
    }

}
