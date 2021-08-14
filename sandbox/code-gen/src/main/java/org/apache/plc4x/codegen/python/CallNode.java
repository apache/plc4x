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
package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

import java.util.ArrayList;
import java.util.List;

@JsonTypeName("Call")
public class CallNode extends LineEntryNode {

    @JsonProperty("args")
    private List<Node> args = new ArrayList<>();

    @JsonProperty("func")
    private Node func;

    @JsonProperty("keywords")
    private List<Node> keywords = new ArrayList<>();

    public List<Node> getArgs() {
        return args;
    }

    public void setArgs(List<Node> args) {
        this.args = args;
    }

    public Node getFunc() {
        return func;
    }

    public void setFunc(Node func) {
        this.func = func;
    }

    public List<Node> getKeywords() {
        return keywords;
    }

    public void setKeywords(List<Node> keywords) {
        this.keywords = keywords;
    }

    @Override
    public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
