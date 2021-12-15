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

using System;

namespace org.apache.plc4net.spi.generation
{
    public class EvaluationHelper
    {
        /*public static bool equals(Object val1, Object val2) {
            if((val1 == null) && (val2 == null)) {
                return true;
            }
            if((val1 == null) || (val2 == null)) {
                return false;
            }
            if(typeof(val1).IsInstanceOfType(Number) && val2 instanceof Number) {
                Number number1 = (Number) val1;
                Number number2 = (Number) val2;
                return number1.doubleValue() == number2.doubleValue();
            }
            if(val1 instanceof Boolean && val2 instanceof Boolean) {
                Boolean boolean1 = (Boolean) val1;
                Boolean boolean2 = (Boolean) val2;
                return boolean1.equals(boolean2);
            }
            if(val1 instanceof String && val2 instanceof String) {
                String string1 = (String) val1;
                String string2 = (String) val2;
                return string1.equals(string2);
            }
            if(val1.getClass().isEnum() && val2.getClass().isEnum()) {
                return val1.equals(val2);
            }
            return false;
        }*/

        public static bool notEquals(Object val1, Object val2) {
            return true;
        }

        public static bool greater(Object val1, Object val2) {
            return true;
        }

        public static bool greaterEquals(Object val1, Object val2) {
            return true;
        }

        public static bool smaller(Object val1, Object val2) {
            return true;
        }

        public static bool smallerEquals(Object val1, Object val2) {
            return true;
        }
    }
}