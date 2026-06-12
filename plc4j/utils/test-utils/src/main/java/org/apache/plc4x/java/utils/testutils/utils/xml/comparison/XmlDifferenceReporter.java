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
package org.apache.plc4x.java.utils.testutils.utils.xml.comparison;

import org.slf4j.Logger;

import java.util.Arrays;
import java.util.List;

/**
 * Reports XML comparison differences in various formats.
 * Supports both string-based and streaming output to handle large XMLs.
 */
public class XmlDifferenceReporter {

    /**
     * Generates a detailed diff report as a string.
     * Use this for small to medium XMLs.
     *
     * @param result the comparison result
     * @return formatted diff report
     */
    public static String generateDiffReport(XmlComparisonResult result) {
        if (result.isIdentical()) {
            return "XMLs are identical";
        }

        StringBuilder sb = new StringBuilder();
        sb.append("\n");
        sb.append("=".repeat(185)).append("\n");
        sb.append("XML COMPARISON FAILED\n");
        sb.append("=".repeat(185)).append("\n\n");

        // List all differences
        sb.append("DIFFERENCES FOUND:\n");
        sb.append("-".repeat(50)).append("\n");

        int maxDiff = Math.min(result.getDifferenceCount(), 50);
        for (int i = 0; i < maxDiff; i++) {
            XmlDifference diff = result.getDifferences().get(i);
            sb.append(String.format("%d. %s%n", i + 1, diff.toString()));
        }

        if (result.getDifferenceCount() > maxDiff) {
            sb.append(String.format("... and %d more differences\n", result.getDifferenceCount() - maxDiff));
        }
        sb.append("\n");

        // Side-by-side comparison
        sb.append("SIDE-BY-SIDE COMPARISON:\n");
        sb.append("-".repeat(185)).append("\n");

        appendSideBySide(sb, result.getExpectedXml(), result.getActualXml());

        sb.append("\n");
        sb.append("=".repeat(185)).append("\n");

        return sb.toString();
    }

    /**
     * Logs the diff report line-by-line to avoid memory issues with huge XMLs.
     * Use this for very large XMLs.
     *
     * @param result the comparison result
     * @param logger the logger to use
     */
    public static void logDiffReport(XmlComparisonResult result, Logger logger) {
        if (result.isIdentical()) {
            logger.info("XMLs are identical");
            return;
        }

        logger.error("");
        logger.error("================================================================================");
        logger.error("XML COMPARISON FAILED");
        logger.error("================================================================================");
        logger.error("");

        logger.error("DIFFERENCES FOUND:");
        logger.error("--------------------------------------------------------------------------------");

        int maxDiff = Math.min(result.getDifferenceCount(), 50);
        for (int i = 0; i < maxDiff; i++) {
            XmlDifference diff = result.getDifferences().get(i);
            logger.error("{}. {}", i + 1, diff.toString());
        }

        if (result.getDifferenceCount() > maxDiff) {
            logger.error("... and {} more differences", result.getDifferenceCount() - maxDiff);
        }
        logger.error("");

        // Add side-by-side comparison
        logger.error("SIDE-BY-SIDE COMPARISON:");
        logger.error("--------------------------------------------------------------------------------");
        logSideBySide(result.getExpectedXml(), result.getActualXml(), logger);
        logger.error("");

        logger.error("Summary: {} differences found", result.getDifferenceCount());
    }

    /**
     * Logs a side-by-side comparison line-by-line without building huge strings.
     * This is memory-safe for huge XML documents.
     *
     * @param expected the expected XML string
     * @param actual the actual XML string
     * @param logger the logger to use
     */
    private static void logSideBySide(String expected, String actual, Logger logger) {
        // Split into lines but process them in chunks to avoid loading everything into memory
        String[] expectedLines = expected.split("\n");
        String[] actualLines = actual.split("\n");

        int maxLines = Math.max(expectedLines.length, actualLines.length);
        int colWidth = 90;

        // Log header
        logger.error(String.format("%-" + colWidth + "s | %s", "EXPECTED", "ACTUAL"));
        logger.error("-".repeat(colWidth) + " | " + "-".repeat(colWidth));

        // Limit output to prevent excessive logging for huge documents
        int maxLinesToShow = Math.min(maxLines, 500);

        for (int i = 0; i < maxLinesToShow; i++) {
            String expLineFull = i < expectedLines.length ? expectedLines[i].trim() : "";
            String actLineFull = i < actualLines.length ? actualLines[i].trim() : "";

            // Skip empty lines
            if (expLineFull.isEmpty() && actLineFull.isEmpty()) {
                continue;
            }

            String expLine = truncate(expLineFull, colWidth);
            String actLine = truncate(actLineFull, colWidth);

            // Mark differing lines with *
            String marker = expLineFull.equals(actLineFull) ? " " : "*";
            logger.error(String.format("%-" + colWidth + "s %s %s", expLine, marker, actLine));
        }

        if (maxLines > maxLinesToShow) {
            logger.error("");
            logger.error("... and {} more lines (truncated for readability)", maxLines - maxLinesToShow);
        }
    }

    /**
     * Appends a side-by-side comparison of two XML strings.
     */
    private static void appendSideBySide(StringBuilder sb, String expected, String actual) {
        List<String> expectedLines = Arrays.stream(expected.split("\n"))
            .map(String::trim)
            .filter(line -> !line.isEmpty())
            .toList();

        List<String> actualLines = Arrays.stream(actual.split("\n"))
            .map(String::trim)
            .filter(line -> !line.isEmpty())
            .toList();

        int maxLines = Math.max(expectedLines.size(), actualLines.size());
        int colWidth = 90;

        sb.append(String.format("%-" + colWidth + "s | %s%n", "EXPECTED", "ACTUAL"));
        sb.append("-".repeat(colWidth)).append(" | ").append("-".repeat(colWidth)).append("\n");

        for (int i = 0; i < maxLines; i++) {
            String expLineFull = i < expectedLines.size() ? expectedLines.get(i) : "";
            String actLineFull = i < actualLines.size() ? actualLines.get(i) : "";
            String expLine = truncate(expLineFull, colWidth);
            String actLine = truncate(actLineFull, colWidth);

            // Mark differing lines
            String marker = expLineFull.equals(actLineFull) ? " " : "*";
            sb.append(String.format("%-" + colWidth + "s %s %s%n", expLine, marker, actLine));
        }
    }

    /**
     * Truncates a string to a maximum length, adding "..." if truncated.
     */
    private static String truncate(String str, int maxLength) {
        if (str == null) {
            return "";
        }
        if (str.length() <= maxLength) {
            return str;
        }
        return str.substring(0, maxLength - 3) + "...";
    }
}
