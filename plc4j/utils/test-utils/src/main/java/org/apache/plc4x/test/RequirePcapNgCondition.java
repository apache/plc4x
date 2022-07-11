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

import org.apache.commons.lang3.SystemUtils;
import org.apache.maven.artifact.versioning.DefaultArtifactVersion;
import org.junit.jupiter.api.extension.ConditionEvaluationResult;
import org.junit.jupiter.api.extension.ExecutionCondition;
import org.junit.jupiter.api.extension.ExtensionContext;
import org.pcap4j.core.Pcaps;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class RequirePcapNgCondition implements ExecutionCondition {

    private static final Logger logger = LoggerFactory.getLogger(RequirePcapNgCondition.class);

    @Override
    public ConditionEvaluationResult evaluateExecutionCondition(ExtensionContext extensionContext) {
        // Mac:     libpcap version 1.8.1 -- Apple version 79.200.4
        // Linux:
        // Windows: NPcap version 1.6.0
        try {
            // On Mac we need to force the usage of the updated libpcap
            if(System.getProperty( "os.name" ).startsWith( "Mac" )) {
                // TODO: Possibly find the highest version ...
                System.getProperties().setProperty("jna.library.path", "/usr/local/Cellar/libpcap/1.10.1/lib");
            }

            String libVersion = Pcaps.libVersion();
            Pattern pattern = Pattern.compile("^.*libpcap version (?<version>\\d+\\.\\d+(?:\\.\\d+)?)[^\\d]?.*$");
            Matcher matcher = pattern.matcher(libVersion);
            if (matcher.matches()) {
                String versionString = matcher.group("version");
                DefaultArtifactVersion curVersion = new DefaultArtifactVersion(versionString);
                // Macs ship with version 1.9.1, which causes exceptions.
                DefaultArtifactVersion minVersion = new DefaultArtifactVersion("1.10.1");
                if (curVersion.compareTo(minVersion) >= 0) {
                    return ConditionEvaluationResult.enabled("Found libpcap version " + versionString);
                } else if (SystemUtils.IS_OS_WINDOWS) {
                    return ConditionEvaluationResult.disabled("Test disabled due to too old Npcap version. Please install from here: https://npcap.com/ as this version supports all needed features.");
                } else {
                    return ConditionEvaluationResult.disabled("Test disabled due to too old libpcap version. Please install at least version 1.10.1 to support all features.");
                }
            }
        } catch(Throwable e) {
            logger.info("Error detecting libpcap version.", e);
        }
        if(SystemUtils.IS_OS_WINDOWS) {
            return ConditionEvaluationResult.disabled("Test disabled due to missing or invalid Npcap version. Please install from here: https://npcap.com/ as this version supports all needed features.");
        } else {
            return ConditionEvaluationResult.disabled("Test disabled due to missing or invalid libpcap version. Please install at least version 1.10.1 to support all features.");
        }
    }

}
