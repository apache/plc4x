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

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

import java.util.ArrayList;
import java.util.List;

@JsonTypeName("ClassDef")
@JsonIgnoreProperties(ignoreUnknown = true)
public class ClassDefNode extends LineEntryNode {

    @JsonProperty("bases")
    private List<Node> bases = new ArrayList<>();

    @JsonProperty("body")
    private List<Node> body = new ArrayList<>();

    @JsonProperty("name")
    private String name;

    @JsonCreator
    public ClassDefNode() {
    }

    public List<Node> getBases() {
        return bases;
    }

    public void setBases(List<Node> bases) {
        this.bases = bases;
    }

    public List<Node> getBody() {
        return body;
    }

    public void setBody(List<Node> body) {
        this.body = body;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    @Override
    public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
