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
package org.apache.plc4x.test;

import org.junit.jupiter.api.extension.ConditionEvaluationResult;
import org.junit.jupiter.api.extension.ExecutionCondition;
import org.junit.jupiter.api.extension.ExtensionContext;

public class RequireAllTestsFlagCondition implements ExecutionCondition {

    @Override
    public ConditionEvaluationResult evaluateExecutionCondition(ExtensionContext extensionContext) {
        String propertyValue = System.getenv("ENABLE_ALL_TESTS");
        if("true".equalsIgnoreCase(propertyValue)) {
            return ConditionEvaluationResult.enabled("ENABLE_ALL_TESTS variable set to 'true'.");
        }
        System.out.println("DISABLED-RequireAllTestsFlagCondition");
        return ConditionEvaluationResult.disabled("ENABLE_ALL_TESTS variable not set to 'true', therefore disabling this test.");
    }

}
