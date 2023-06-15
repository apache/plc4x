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

import org.apache.commons.io.IOUtils;
import org.apache.commons.lang3.SystemUtils;
import org.junit.jupiter.api.extension.ConditionEvaluationResult;
import org.junit.jupiter.api.extension.ExecutionCondition;
import org.junit.jupiter.api.extension.ExtensionContext;

import java.io.*;
import java.nio.charset.StandardCharsets;

public class DisableOnParallelsVmFlagCondition implements ExecutionCondition {

    private static final String PARALLELS_STRING_AMD = "Parallels Virtual Platform";
    private static final String PARALLELS_STRING_ARM = "Parallels ARM Virtual Machine";
    private static final boolean isParallels;
    static {
        boolean localIsParallels = false;
        if(SystemUtils.IS_OS_WINDOWS) {
            // TODO: If on Windows: Run "systeminfo /fo CSV /nh" command and check if the output contains "Parallels Virtual Platform"
            try {
                ProcessBuilder processBuilder = new ProcessBuilder();
                processBuilder.command("systeminfo", "/fo", "CSV", "/nh");
                Process process = processBuilder.start();
                try (BufferedReader reader = new BufferedReader(
                    new InputStreamReader(process.getInputStream()))) {
                    String line;
                    while ((line = reader.readLine()) != null) {
                        if(line.contains(PARALLELS_STRING_AMD) || line.contains(PARALLELS_STRING_ARM)) {
                            localIsParallels = true;
                            break;
                        }
                    }
                }
            } catch (Exception err) {
                // Ignore this...
            }
        } else if (SystemUtils.IS_OS_LINUX) {
            // If on Linux: check /sys/devices/virtual/dmi/id/product_name contains "Parallels Virtual Platform"
            File productNameFile = new File("/sys/devices/virtual/dmi/id/product_name");
            try(InputStream is = new FileInputStream(productNameFile)) {
                String content = IOUtils.toString(is, StandardCharsets.UTF_8);
                localIsParallels = content.contains(PARALLELS_STRING_AMD) || content.contains(PARALLELS_STRING_ARM);
            } catch (IOException e) {
                // Ignore this...
            }
        }
        isParallels = localIsParallels;
    }

    @Override
    public ConditionEvaluationResult evaluateExecutionCondition(ExtensionContext extensionContext) {
        if(isParallels) {
            return ConditionEvaluationResult.disabled("Parallels detected");
        }
        return ConditionEvaluationResult.enabled("Parallels not detected");
    }

}
