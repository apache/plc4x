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

package org.apache.plc4x.java.spi.utils.ascii;

import org.apache.commons.lang3.StringUtils;

import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class AsciiBox {

    // source: https://github.com/chalk/ansi-regex/blob/main/index.js#L3
    private final Pattern ANSI_PATTERN = Pattern.compile("[\u001b\u009b][\\[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]");

    private final AsciiBoxWriter asciiBoxWriter;

    private final String data;

    // TODO: should be final but for the moment we mutate it in change box name... Maybe we find another solution
    String compressedBoxSet;

    protected AsciiBox(String data) {
        this(AsciiBoxWriter.DEFAULT, data);
    }

    protected AsciiBox(AsciiBoxWriter asciiBoxWriter, String data) {
        Objects.requireNonNull(data);
        this.asciiBoxWriter = asciiBoxWriter;
        this.data = data;
        this.compressedBoxSet = asciiBoxWriter.boxSet.compressBoxSet();
    }

    /**
     * @return returns the width of the box without the newlines
     */
    public int width() {
        int maxWidth = 0;
        for (String line : lines()) {
            int currentLength = ANSI_PATTERN.matcher(line).replaceAll("").length();
            if (maxWidth < currentLength) {
                maxWidth = currentLength;
            }
        }
        return maxWidth;
    }

    /**
     * @return returns the height of the box without
     */
    public int height() {
        return lines().length;
    }

    /**
     * @return returns the lines of the box
     */
    public String[] lines() {
        return data.split("\n");
    }

    public String getBoxName() {
        Matcher matcher = asciiBoxWriter.boxNameRegex.matcher(data);
        if (!matcher.find()) {
            return "";
        }
        return matcher.group("name");
    }

    public AsciiBox changeBoxName(String newName) {
        if (!asciiBoxWriter.hasBorders(this)) {
            return asciiBoxWriter.boxString(newName, this.toString(), 0);
        }
        int minimumWidthWithNewName = (asciiBoxWriter.boxSet.upperLeftCorner + asciiBoxWriter.boxSet.horizontalLine + newName + asciiBoxWriter.boxSet.upperRightCorner).length();
        int nameLengthDifference = minimumWidthWithNewName - (asciiBoxWriter.unwrap(this).width() + asciiBoxWriter.borderWidth + asciiBoxWriter.borderWidth);
        AsciiBox asciiBox = asciiBoxWriter.boxString(newName, asciiBoxWriter.unwrap(this).toString(), this.width() + nameLengthDifference);
        asciiBox.compressedBoxSet = asciiBoxWriter.boxSet.contributeToCompressedBoxSet(this);
        return asciiBox;
    }

    public boolean isEmpty() {
        if (asciiBoxWriter.hasBorders(this)) {
            return StringUtils.isBlank(asciiBoxWriter.unwrap(this).toString());
        }
        return StringUtils.isBlank(this.toString());
    }

    @Override
    public String toString() {
        return data;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        AsciiBox asciiBox = (AsciiBox) o;
        return Objects.equals(data, asciiBox.data);
    }

    @Override
    public int hashCode() {
        return Objects.hash(data);
    }
}
