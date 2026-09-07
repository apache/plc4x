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

package org.apache.plc4x.java.spi.fields.utils;

public class EvaluationHelper {

    public static boolean equals(Object val1, Object val2) {
        if ((val1 == null) && (val2 == null)) {
            return true;
        }
        if ((val1 == null) || (val2 == null)) {
            return false;
        }
        if (val1 instanceof Number number1 && val2 instanceof Number number2) {
            return number1.doubleValue() == number2.doubleValue();
        }
        if (val1 instanceof Boolean boolean1 && val2 instanceof Boolean boolean2) {
            return boolean1.equals(boolean2);
        }
        if (val1 instanceof String string1 && val2 instanceof String string2) {
            return string1.equals(string2);
        }
        if (val1.getClass().isEnum() && val2.getClass().isEnum()) {
            return val1.equals(val2);
        }
        return false;
    }

}
