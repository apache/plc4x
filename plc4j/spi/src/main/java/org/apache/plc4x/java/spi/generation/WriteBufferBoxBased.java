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

package org.apache.plc4x.java.spi.generation;

import io.vavr.control.Either;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.spi.utils.ascii.AsciiBox;
import org.apache.plc4x.java.spi.utils.hex.Hex;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Collections;
import java.util.LinkedList;
import java.util.Deque;

import static org.apache.plc4x.java.spi.utils.ascii.AsciiBoxUtils.*;

public class WriteBufferBoxBased implements WriteBuffer, BufferCommons {

    private Deque<Either<AsciiBox, Deque<AsciiBox>>> boxes = new LinkedList<>();
    private int desiredWidth = 120;
    private int currentWidth = 118;
    private boolean mergeSingleBoxes;
    private boolean omitEmptyBoxes;

    public WriteBufferBoxBased() {

    }

    public WriteBufferBoxBased(boolean mergeSingleBoxes, boolean omitEmptyBoxes) {
        this.mergeSingleBoxes = mergeSingleBoxes;
        this.omitEmptyBoxes = omitEmptyBoxes;
    }

    @Override
    public ByteOrder getByteOrder() {
        // NO OP
        return ByteOrder.BIG_ENDIAN;
    }

    @Override
    public void setByteOrder(ByteOrder byteOrder) {
        // NO OP
    }

    @Override
    public int getPos() {
        return 0;
    }

    @Override
    public void pushContext(String logicalName, WithWriterArgs... writerArgs) {
        currentWidth -= Hex.boxLineOverheat;
        boxes.offerLast(Either.right(new LinkedList<>()));
    }

    @Override
    public void writeBit(String logicalName, boolean value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("b%d %b%s", value ? 1 : 0, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeByte(String logicalName, byte value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("%02x%s", value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeByteArray(String logicalName, byte[] bytes, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        if (StringUtils.isNotBlank(additionalStringRepresentation)) {
            additionalStringRepresentation += "\n";
        }
        boxes.offerLast(Either.left(boxString(logicalName, String.format("%s%s", Hex.dump(bytes), additionalStringRepresentation), 0)));
    }

    @Override
    public void writeUnsignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeUnsignedShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeUnsignedInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeUnsignedLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeUnsignedBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeSignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeFloat(String logicalName, int bitLength, float value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeDouble(String logicalName, int bitLength, double value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeBigDecimal(String logicalName, int bitLength, BigDecimal value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0)));
    }

    @Override
    public void writeString(String logicalName, int bitLength, String encoding, String value, WithWriterArgs... writerArgs) throws SerializationException {
        String additionalStringRepresentation = extractAdditionalStringRepresentation(writerArgs).map(s -> " " + s).orElse("");
        boxes.offerLast(Either.left(boxString(logicalName, String.format("%s%s", value, additionalStringRepresentation), 0)));
    }

    @Override
    public void popContext(String logicalName, WithWriterArgs... writerArgs) {
        currentWidth += Hex.boxLineOverheat;
        Deque<AsciiBox> finalBoxes = new LinkedList<>();
        findTheBox:
        for (Either<AsciiBox, Deque<AsciiBox>> back = boxes.pollLast(); back != null; back = boxes.pollLast()) {
            if (back.isLeft()) {
                AsciiBox asciiBox = back.getLeft();
                if (omitEmptyBoxes && asciiBox.isEmpty()) {
                    continue;
                }
                finalBoxes.offerFirst(asciiBox);
            } else {
                Deque<AsciiBox> asciiBoxes = back.get();
                LinkedList<AsciiBox> reversedList = new LinkedList<>(asciiBoxes);
                Collections.reverse(reversedList);
                for (AsciiBox box : asciiBoxes) {
                    finalBoxes.offerFirst(box);
                }
                break findTheBox;
            }
        }
        if (mergeSingleBoxes && finalBoxes.size() == 1) {
            AsciiBox onlyChild = finalBoxes.remove();
            String childName = onlyChild.getBoxName();
            onlyChild = onlyChild.changeBoxName(logicalName + "/" + childName);
            if (omitEmptyBoxes && onlyChild.isEmpty()) {
                return;
            }
            boxes.offerLast(Either.left(onlyChild));
            return;
        }
        AsciiBox asciiBox = boxBox(logicalName, alignBoxes(finalBoxes, currentWidth), 0);
        if (omitEmptyBoxes && asciiBox.isEmpty()) {
            return;
        }
        boxes.offerLast(Either.left(asciiBox));
    }

    public AsciiBox getBox() {
        return boxes.peek().getLeft();
    }
}
