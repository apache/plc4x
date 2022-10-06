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
package org.apache.plc4x.test.migration;

import org.apache.commons.lang3.StringUtils;

import java.util.Arrays;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import java.util.stream.Collectors;

/**
 * Small util to help with patching of testcases
 */
public class TestCasePatcher {

    /**
     * Indents a xml fragment with the supplied indent (every line gets ident prefixed).
     *
     * @param xmlDocument the document to be patched.
     * @param indent      the indent to be applied
     * @return the indented document
     */
    public static String indent(String xmlDocument, String indent) {
        return Arrays.stream(xmlDocument.split("\n")).map(s -> indent + s).collect(Collectors.joining("\n"));
    }

    /**
     * Tries to find the base indent for supplied fragment
     *
     * @param xmlDocument where to look for the fragment
     * @param xmlFragment the fragment to look up
     * @return the found indent
     */
    public static String determineIndent(String xmlDocument, String xmlFragment) {
        assert xmlDocument != null;
        assert xmlFragment != null;

        Pattern pattern = getPatternForFragment(xmlFragment);
        Matcher matcher = pattern.matcher(xmlDocument);
        if (!matcher.find()) {
            throw new RuntimeException("Could not match");
        }
        return matcher.group(1);
    }

    /**
     * Returns a pattern for a xmlFragment which ignores leading indents.
     * The first group can be used to determine ident.
     *
     * @param xmlFragment the fragment where the pattern should be build for
     * @return the created pattern
     */
    public static Pattern getPatternForFragment(String xmlFragment) {
        assert xmlFragment != null;

        StringBuilder patternString = new StringBuilder();
        String[] lines = xmlFragment.split("\n");
        for (int i = 0; i < lines.length; i++) {
            String line = lines[i];
            line = StringUtils.replace(line, "\"", "\\\"");
            line = StringUtils.replace(line, ".", "\\.");
            if (i == 0) {
                patternString.append("([ ]*)").append(line).append("\\n");
                continue;
            }
            patternString.append("[ ]*").append(line).append("\\n");
        }
        return Pattern.compile(patternString.toString());
    }
}
