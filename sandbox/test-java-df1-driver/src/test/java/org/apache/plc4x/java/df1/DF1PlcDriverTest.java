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
package org.apache.plc4x.java.df1;

import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import java.util.regex.Matcher;

import static org.apache.plc4x.java.df1.DF1PlcDriver.DF1_URI_PATTERN;
import static org.apache.plc4x.java.df1.DF1PlcDriver.SERIAL_PATTERN;

public class DF1PlcDriverTest {

    @Test
    public void matchExpression() {
        Matcher matcher = SERIAL_PATTERN.matcher("serial:///COM4");

        Assertions.assertTrue(matcher.matches());
    }

    @Test
    public void matchExpression2() {
        Matcher matcher = DF1_URI_PATTERN.matcher("df1:serial:///COM4");

        Assertions.assertTrue(matcher.matches());
    }
}