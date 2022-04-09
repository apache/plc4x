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
package org.apache.plc4x.codegen.util;

public class MyPojo {

    private Double field1;

    private Double field2;

    private Double field3;

    public MyPojo() {
    }

    public Double getField1() {
        return this.field1;
    }

    public Double getField2() {
        return this.field2;
    }

    public Double getField3() {
        return this.field3;
    }

    public void setField1(Double field1) {
        this.field1 = field1;
    }

    public void setField2(Double field2) {
        this.field2 = field2;
    }

    public void setField3(Double field3) {
        this.field3 = field3;
    }

    public void encode(org.apache.plc4x.codegen.api.Buffer buffer) {
    }

    public static MyPojo decode(org.apache.plc4x.codegen.api.Buffer buffer) {
        MyPojo instance = new MyPojo();
        buffer.readUint8();
        return instance;
    }

}