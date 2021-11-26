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

package org.apache.plc4x.java.spi.utils.ascii;

import org.apache.commons.lang3.StringUtils;

import java.util.Objects;
import java.util.regex.Matcher;

import static org.apache.plc4x.java.spi.utils.ascii.AsciiBoxUtils.*;

public class AsciiBox {
    private final String data;

    public AsciiBox(String data) {
        this.data = data;
    }

    /**
     * @return returns the width of the box without the newlines
     */
    public int width() {
        int maxWidth = 0;
        for (String line : lines()) {
            int currentLength = line.length();
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
        Matcher matcher = boxNameRegex.matcher(data);
        if (!matcher.find()) {
            return "";
        }
        return matcher.group("name");
    }

    public AsciiBox changeBoxName(String newName) {
        if (!hasBorders(this)) {
            return boxString(newName, this.toString(), 0);
        }
        int minimumWidthWithNewName = (upperLeftCorner + horizontalLine + newName + upperRightCorner).length();
        int nameLengthDifference = minimumWidthWithNewName - (unwrap(this).width() + borderWidth + borderWidth);
        return boxString(newName, unwrap(this).toString(), this.width() + nameLengthDifference);
    }

    public boolean isEmpty() {
        if (hasBorders(this)) {
            return StringUtils.isBlank(unwrap(this).toString());
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
