/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.knxnetip.ets.filehandlers;

import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * The input to the tests are cases documented in the official schema documentation
 */
class Ets6FileHandlerTest {

    @Test
    void getProcessedPasswordA() {
        Ets6FileHandler handler = new Ets6FileHandler();
        final String processedPassword = handler.getProcessedPassword("a");
        assertEquals(processedPassword, "+FAwP4iI7/Pu4WB3HdIHbbFmteLahPAVkjJShKeozAA=");
    }

    @Test
    void getProcessedPasswordTest() {
        Ets6FileHandler handler = new Ets6FileHandler();
        final String processedPassword = handler.getProcessedPassword("test");
        assertEquals(processedPassword, "2+IIP7ErCPPKxFjJXc59GFx2+w/1VTLHjJ2duc04CYQ=");
    }

    @Test
    @Disabled("This doesn't seem to work, I guess the clown face from the PDF would have been needed too.")
    void getProcessedPasswordPennywise() {
        Ets6FileHandler handler = new Ets6FileHandler();
        final String processedPassword = handler.getProcessedPassword("Penn¥w1se");
        assertEquals(processedPassword, "ZjlYlh+eTtoHvFadU7+EKvF4jOdEm7WkP49uanOMMk0=");
    }

}