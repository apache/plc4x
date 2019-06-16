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

public class BinaryTerm implements Term {

    private final Term a;
    private final Term b;
    private final String operation;

    public BinaryTerm(Term a, Term b, String operation) {
        this.a = a;
        this.b = b;
        this.operation = operation;
    }

    public Term getA() {
        return a;
    }

    public Term getB() {
        return b;
    }

    public String getOperation() {
        return operation;
    }

    @Override
    public boolean contains(String str) {
        return ((a != null) && a.contains(str)) || ((b != null) && b.contains(str));
    }

}
