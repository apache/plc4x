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

package org.apache.plc4x.language.expressions.terms;

public class VariableLiteral implements Literal {

    public static final int NO_INDEX = -1;

    private final String name;
    private final int index;
    private final VariableLiteral child;

    public VariableLiteral(String name, int index, VariableLiteral child) {
        this.name = name;
        this.index = index;
        this.child = child;
    }

    public String getName() {
        return name;
    }

    public int getIndex() {
        return index;
    }

    public VariableLiteral getChild() {
        return child;
    }

    public boolean isIndexed() {
        return index != NO_INDEX;
    }

}
