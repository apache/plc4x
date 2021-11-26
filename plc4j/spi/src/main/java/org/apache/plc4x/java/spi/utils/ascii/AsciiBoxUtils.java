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
import org.apache.plc4x.java.spi.utils.hex.Hex;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Collection;
import java.util.LinkedList;
import java.util.List;
import java.util.regex.Pattern;

public class AsciiBoxUtils {

    private static final Logger LOGGER = LoggerFactory.getLogger(AsciiBoxUtils.class);

    static final String upperLeftCorner = "╔";
    static final String upperRightCorner = "╗";
    static final String horizontalLine = "═";
    static final String verticalLine = "║";
    static final String lowerLeftCorner = "╚";
    static final String lowerRightCorner = "╝";
    static final String newLine = "\n";
    static final String emptyPadding = " ";
    // the name gets prefixed with a extra symbol for indent
    static final int extraNameCharIndent = 1;
    static final int borderWidth = 1;
    static final int newLineCharWidth = 1;
    static Pattern boxNameRegex = Pattern.compile("^" + upperLeftCorner + horizontalLine + "(?<name>[\\w /]+)" + horizontalLine + "*" + upperRightCorner);

    static boolean DebugAsciiBox;

    /**
     * BoxBox boxes a box
     *
     * @param name      the name for the new box
     * @param box       the box to box
     * @param charWidth the desired width
     * @return boxed data
     */
    public static AsciiBox boxBox(String name, AsciiBox box, int charWidth) {
        return boxString(name, box.toString(), charWidth);
    }

    /**
     * BoxString boxes a newline separated string into a beautiful box
     *
     * @param name      of the box (can be empty)
     * @param data      data to be boxed
     * @param charWidth desired width
     * @return boxed data
     */
    public static AsciiBox boxString(String name, String data, int charWidth) {
        AsciiBox rawBox = new AsciiBox(data);
        int longestLine = rawBox.width();
        if (charWidth < longestLine) {
            if (DebugAsciiBox) {
                LOGGER.debug("Overflow by {} chars", longestLine - charWidth);
            }
            charWidth = longestLine + borderWidth + borderWidth;
        }
        StringBuilder boxedString = new StringBuilder();
        int namePadding = (Math.max(charWidth - name.length() - borderWidth - extraNameCharIndent - borderWidth, 0));
        boxedString.append(upperLeftCorner + horizontalLine).append(name).append(StringUtils.repeat(horizontalLine, namePadding)).append(upperRightCorner);
        boxedString.append(newLine);
        // Name of the header stretches the box so we align to that
        charWidth = borderWidth + extraNameCharIndent + name.length() + namePadding + borderWidth;
        for (String line : rawBox.lines()) {
            int linePadding = charWidth - Hex.boxLineOverheat - line.length();
            if (linePadding < 0) {
                linePadding = 0;
            }
            int frontPadding = (int) Math.floor(linePadding / 2.0);
            int backPadding = (int) Math.ceil(linePadding / 2.0);
            boxedString.append(verticalLine).append(StringUtils.repeat(emptyPadding, frontPadding)).append(line).append(StringUtils.repeat(emptyPadding, backPadding)).append(verticalLine);
            boxedString.append(newLine);
        }
        int bottomPadding = namePadding + name.length() + extraNameCharIndent;
        boxedString.append(lowerLeftCorner).append(StringUtils.repeat(horizontalLine, bottomPadding)).append(lowerRightCorner);
        return new AsciiBox(boxedString.toString());
    }

    /**
     * AlignBoxes aligns all boxes to a desiredWidth and orders them from left to right and top to bottom (size will be at min the size of the biggest box)
     *
     * @param boxes        to be aligned.
     * @param desiredWidth width desired
     * @return the aligned box.
     */
    public static AsciiBox alignBoxes(Collection<AsciiBox> boxes, int desiredWidth) {
        if (boxes.size() == 0) {
            return new AsciiBox("");
        }
        int actualWidth = desiredWidth;
        for (AsciiBox box : boxes) {
            int boxWidth = box.width();
            if (boxWidth > actualWidth) {
                if (DebugAsciiBox) {
                    LOGGER.debug("Overflow by {} chars", boxWidth - desiredWidth);
                }
                actualWidth = boxWidth;
            }
        }
        if (DebugAsciiBox) {
            LOGGER.debug("Working with {} chars", actualWidth);
        }
        AsciiBox bigBox = new AsciiBox("");
        List<AsciiBox> currentBoxRow = new LinkedList<>();
        int currentRowLength = 0;
        for (AsciiBox box : boxes) {
            currentRowLength += box.width();
            if (currentRowLength > actualWidth) {
                AsciiBox mergedBoxes = mergeHorizontal(currentBoxRow);
                if (StringUtils.isBlank(bigBox.toString())) {
                    bigBox = mergedBoxes;
                } else {
                    bigBox = boxBelowBox(bigBox, mergedBoxes);
                }
                currentRowLength = box.width();
                currentBoxRow = new LinkedList<>();
            }
            currentBoxRow.add(box);
        }
        if (currentBoxRow.size() > 0) {
            // Special case where all boxes fit into one row
            AsciiBox mergedBoxes = mergeHorizontal(currentBoxRow);
            if (StringUtils.isBlank(bigBox.toString())) {
                bigBox = mergedBoxes;
            } else {
                bigBox = boxBelowBox(bigBox, mergedBoxes);
            }
        }
        return bigBox;
    }

    /**
     * BoxSideBySide renders two boxes side by side
     *
     * @param box1 left of box2
     * @param box2 right of box1
     * @return box1 left of box2
     */
    public static AsciiBox boxSideBySide(AsciiBox box1, AsciiBox box2) {
        StringBuilder aggregateBox = new StringBuilder();
        int box1Width = box1.width();
        String[] box1Lines = box1.lines();
        int box2Width = box2.width();
        String[] box2Lines = box2.lines();
        int maxRows = Math.max(box1Lines.length, box2Lines.length);
        for (int row = 0; row < maxRows; row++) {
            boolean ranOutOfLines = false;
            if (row >= box1Lines.length) {
                ranOutOfLines = true;
                aggregateBox.append(StringUtils.repeat(" ", box1Width));
            } else {
                String split1Row = box1Lines[row];
                int padding = box1Width - split1Row.length();
                aggregateBox.append(split1Row).append(StringUtils.repeat(" ", padding));
            }
            if (row >= box2Lines.length) {
                if (ranOutOfLines) {
                    break;
                }
                aggregateBox.append(StringUtils.repeat(" ", box2Width));
            } else {
                String split2Row = box2Lines[row];
                int padding = box2Width - split2Row.length();
                aggregateBox.append(split2Row).append(StringUtils.repeat(" ", padding));
            }
            if (row < maxRows - 1) {
                // Only write newline if we are not the last line
                aggregateBox.append('\n');
            }
        }
        return new AsciiBox(aggregateBox.toString());
    }

    /**
     * BoxBelowBox renders two boxes below
     *
     * @param box1 above box2
     * @param box2 below box1
     * @return box1 above box2
     */
    public static AsciiBox boxBelowBox(AsciiBox box1, AsciiBox box2) {
        int box1Width = box1.width();
        int box2Width = box2.width();
        if (box1Width < box2Width) {
            box1 = expandBox(box1, box2Width);
        } else if (box2Width < box1Width) {
            box2 = expandBox(box2, box1Width);
        }
        return new AsciiBox(box1.toString() + "\n" + box2.toString());
    }

    static AsciiBox mergeHorizontal(List<AsciiBox> boxes) {
        switch (boxes.size()) {
            case 0:
                return new AsciiBox("");
            case 1:
                return boxes.get(0);
            case 2:
                return boxSideBySide(boxes.get(0), boxes.get(1));
            default:
                return boxSideBySide(boxes.get(0), mergeHorizontal(new ArrayList<>(boxes).subList(1, boxes.size())));
        }
    }

    static AsciiBox expandBox(AsciiBox box, int desiredWidth) {
        if (box.width() >= desiredWidth) {
            return box;
        }
        String[] boxLines = box.lines();
        int numberOfLine = boxLines.length;
        int boxWidth = box.width();
        String padding = StringUtils.repeat(" ", desiredWidth - boxWidth);
        StringBuilder newBox = new StringBuilder();
        for (int i = 0; i < boxLines.length; i++) {
            String line = boxLines[i];
            newBox.append(line);
            newBox.append(padding);
            if (i < numberOfLine - 1) {
                newBox.append(newLine);
            }
        }
        return new AsciiBox(newBox.toString());
    }

    /**
     * Return true if this box has borders.
     *
     * @param box the box to be checked
     * @return true if it has borders
     */
    public static boolean hasBorders(AsciiBox box) {
        if (StringUtils.isBlank(box.toString())) {
            return false;
        }
        // Check if the first char is the upper left corner
        return upperLeftCorner.equals(box.toString().substring(0, 1));
    }

    public static AsciiBox unwrap(AsciiBox box) {
        if (!hasBorders(box)) {
            return box;
        }
        String[] originalLines = box.lines();
        String[] newLines = new String[originalLines.length - 2];
        for (int i = 0; i < originalLines.length; i++) {
            String line = originalLines[i];
            if (i == 0) {
                // we ignore the first line
                continue;
            }
            if (i == originalLines.length - 1) {
                // we ignore the last line
                break;
            }
            // Strip the vertical Lines and trim the padding
            String unwrappedLine = line.substring(1, line.length() - 1);

            if (!StringUtils.containsAny(unwrappedLine, verticalLine + horizontalLine)) {
                // only trim boxes witch don't contain other boxes
                unwrappedLine = StringUtils.trim(unwrappedLine);
            }
            newLines[i - 1] = unwrappedLine;
        }
        return new AsciiBox(StringUtils.join(newLines, newLine));
    }
}
