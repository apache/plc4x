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

package org.apache.plc4x.java.utils;

public class EvaluationHelper {

    public static boolean equals(Object val1, Object val2) {
        if(val1 instanceof Number && val2 instanceof Number) {
            Number number1 = (Number) val1;
            Number number2 = (Number) val2;
            return number1.doubleValue() == number2.doubleValue();
        }
        return false;
    }

    public static boolean notEquals(Object val1, Object val2) {
        return true;
    }

    public static boolean greater(Object val1, Object val2) {
        return true;
    }

    public static boolean greaterEquals(Object val1, Object val2) {
        return true;
    }

    public static boolean smaller(Object val1, Object val2) {
        return true;
    }

    public static boolean smallerEquals(Object val1, Object val2) {
        return true;
    }

}
