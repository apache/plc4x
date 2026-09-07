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
package org.apache.plc4x.java.utils.testutils;

import org.apache.commons.lang3.SystemUtils;
import org.junit.jupiter.api.extension.ConditionEvaluationResult;
import org.junit.jupiter.api.extension.ExecutionCondition;
import org.junit.jupiter.api.extension.ExtensionContext;
import org.pcap4j.core.Pcaps;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class RequirePcapCondition implements ExecutionCondition {

    private static final Logger logger = LoggerFactory.getLogger(RequirePcapCondition.class);

    @Override
    public ConditionEvaluationResult evaluateExecutionCondition(ExtensionContext extensionContext) {
        // On macOS the bundled libpcap is too old, so point JNA at a Homebrew-installed one
        // before the native bindings are first loaded.
        configureMacOsLibraryPath();
        try {
            String libVersion = Pcaps.libVersion();
            Pattern pattern = Pattern.compile("^.*libpcap version (?<version>\\d+\\.\\d+(?:\\.\\d+)?)[^\\d]?.*$");
            Matcher matcher = pattern.matcher(libVersion);
            if (matcher.matches()) {
                String versionString = matcher.group("version");
                return ConditionEvaluationResult.enabled("Found libpcap version " + versionString);
            }
            logger.info("Could not parse a libpcap version from '{}'", libVersion);
        } catch (Throwable t) {
            // A missing native library surfaces as an Error (e.g. NoClassDefFoundError /
            // UnsatisfiedLinkError), not an Exception - so catch Throwable here, otherwise the
            // whole annotated test class errors out instead of being cleanly disabled
            // (e.g. on Windows CI without Npcap installed).
            logger.info("Error detecting libpcap version.", t);
        }
        if (SystemUtils.IS_OS_WINDOWS) {
            return ConditionEvaluationResult.disabled("Test disabled due to missing or invalid Npcap version. Please install from here: https://npcap.com/ as this version supports all needed features.");
        } else {
            return ConditionEvaluationResult.disabled("Test disabled due to missing or invalid libpcap version. Please install at least version 1.1.0 to support all features.");
        }
    }

    private static void configureMacOsLibraryPath() {
        if (!SystemUtils.IS_OS_MAC) {
            return;
        }
        // On an Intel Mac the libs are in: "/usr/local/Cellar/libpcap/<version>/lib"
        // On an M1 Mac the libs are in: "/opt/homebrew/opt/libpcap/lib"
        if (new File("/usr/local/Cellar/libpcap/1.10.1/lib").exists()) {
            System.getProperties().setProperty("jna.library.path", "/usr/local/Cellar/libpcap/1.10.1/lib");
        } else if (new File("/usr/local/Cellar/libpcap/1.10.5/lib").exists()) {
            System.getProperties().setProperty("jna.library.path", "/usr/local/Cellar/libpcap/1.10.5/lib");
        } else if (new File("/opt/homebrew/opt/libpcap/lib").exists()) {
            System.getProperties().setProperty("jna.library.path", "/opt/homebrew/opt/libpcap/lib");
        }
    }

}
