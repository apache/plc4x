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

import java.util.Arrays;
import java.util.HashSet;
import java.util.Objects;
import java.util.Set;

class BoxSet {
    final String upperLeftCorner;
    final String upperRightCorner;
    final String horizontalLine;
    final String verticalLine;
    final String lowerLeftCorner;
    final String lowerRightCorner;

    public BoxSet(String upperLeftCorner, String upperRightCorner, String horizontalLine, String verticalLine, String lowerLeftCorner, String lowerRightCorner) {
        this.upperLeftCorner = upperLeftCorner;
        this.upperRightCorner = upperRightCorner;
        this.horizontalLine = horizontalLine;
        this.verticalLine = verticalLine;
        this.lowerLeftCorner = lowerLeftCorner;
        this.lowerRightCorner = lowerRightCorner;
    }

    public String compressBoxSet() {
        return upperLeftCorner + upperRightCorner + horizontalLine + verticalLine + lowerLeftCorner + lowerRightCorner;
    }

    String contributeToCompressedBoxSet(AsciiBox box) {
        String actualSet = compressBoxSet();
        if (box.compressedBoxSet.contains(actualSet)) {
            // we have nothing to add
            return box.compressedBoxSet;
        }
        return box.compressedBoxSet + "," + actualSet;
    }

    static String combineCompressedBoxSets(AsciiBox box1, AsciiBox box2) {
        Set<String> allSets = new HashSet<>();
        allSets.addAll(Arrays.asList(box1.compressedBoxSet.split(",")));
        allSets.addAll(Arrays.asList(box2.compressedBoxSet.split(",")));
        return StringUtils.join(allSets, ",");
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        BoxSet boxSet = (BoxSet) o;
        return upperLeftCorner.equals(boxSet.upperLeftCorner) && upperRightCorner.equals(boxSet.upperRightCorner) && horizontalLine.equals(boxSet.horizontalLine) && verticalLine.equals(boxSet.verticalLine) && lowerLeftCorner.equals(boxSet.lowerLeftCorner) && lowerRightCorner.equals(boxSet.lowerRightCorner);
    }

    @Override
    public int hashCode() {
        return Objects.hash(upperLeftCorner, upperRightCorner, horizontalLine, verticalLine, lowerLeftCorner, lowerRightCorner);
    }
}
