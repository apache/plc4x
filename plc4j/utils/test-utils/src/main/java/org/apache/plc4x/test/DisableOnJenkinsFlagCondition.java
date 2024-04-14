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
package org.apache.plc4x.test;

import org.junit.jupiter.api.extension.ConditionEvaluationResult;
import org.junit.jupiter.api.extension.ExecutionCondition;
import org.junit.jupiter.api.extension.ExtensionContext;

public class DisableOnJenkinsFlagCondition implements ExecutionCondition {

    private static final boolean IS_JENKINS;
    static {
        // This environment variable is set in Jenkinsfile.
        String propertyValue = System.getenv("PLC4X_BUILD_ON_JENKINS");
        IS_JENKINS = "true".equalsIgnoreCase(propertyValue);
    }

    @Override
    public ConditionEvaluationResult evaluateExecutionCondition(ExtensionContext extensionContext) {
        if(IS_JENKINS) {
            return ConditionEvaluationResult.disabled("Jenkins detected");
        }
        return ConditionEvaluationResult.enabled("Jenkins not detected");
    }

}
