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
package org.apache.plc4x.java.spi.buffers.asciiboxbased;

import org.apache.plc4x.java.spi.buffers.api.AbstractBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.ascii.AsciiBox;
import org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.ascii.AsciiBoxWriter;
import org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.hex.Hex;
import org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.utils.StringUtils;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.*;

public class WriteBufferAsciiBoxBased extends AbstractBuffer implements WriteBuffer {

    private final Deque<Map.Entry<AsciiBox, Deque<AsciiBox>>> boxes = new LinkedList<>();
    private final AsciiBoxWriter asciiBoxWriter;
    private final AsciiBoxWriter asciiBoxWriterLight;
    private final int desiredWidth = 120;
    private final boolean mergeSingleBoxes;
    private final boolean omitEmptyBoxes;
    private int currentWidth = desiredWidth - 2;

    public WriteBufferAsciiBoxBased() {
        this(false, false);
    }

    public WriteBufferAsciiBoxBased(boolean mergeSingleBoxes, boolean omitEmptyBoxes) {
        this(AsciiBoxWriter.DEFAULT, AsciiBoxWriter.LIGHT, mergeSingleBoxes, omitEmptyBoxes);
    }

    private WriteBufferAsciiBoxBased(AsciiBoxWriter asciiBoxWriter, AsciiBoxWriter asciiBoxWriterLight, boolean mergeSingleBoxes, boolean omitEmptyBoxes) {
        this.asciiBoxWriter = asciiBoxWriter;
        this.asciiBoxWriterLight = asciiBoxWriterLight;
        this.mergeSingleBoxes = mergeSingleBoxes;
        this.omitEmptyBoxes = omitEmptyBoxes;
    }

    @Override
    public void pushContext(WithOption... writerArgs) {
        currentWidth -= Hex.boxLineOverheat;
        boxes.offerLast(new AbstractMap.SimpleEntry<>(null, new LinkedList<>()));
    }

    @Override
    public void popContext(WithOption... options) throws BufferException {
        String name = getName(WithOption.AddOptions(getContext(), options));
        currentWidth += Hex.boxLineOverheat;
        Deque<AsciiBox> finalBoxes = new LinkedList<>();
        findTheBox:
        for (Map.Entry<AsciiBox, Deque<AsciiBox>> back = boxes.pollLast(); back != null; back = boxes.pollLast()) {
            if (back.getKey() != null) {
                AsciiBox asciiBox = back.getKey();
                if (omitEmptyBoxes && asciiBox.isEmpty()) {
                    continue;
                }
                finalBoxes.offerFirst(asciiBox);
            } else {
                Deque<AsciiBox> asciiBoxes = back.getValue();
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
            onlyChild = onlyChild.changeBoxName(name + "/" + childName);
            if (omitEmptyBoxes && onlyChild.isEmpty()) {
                return;
            }
            boxes.offerLast(new AbstractMap.SimpleEntry<>(onlyChild, null));
            return;
        }
        AsciiBox asciiBox = asciiBoxWriter.boxBox(name, asciiBoxWriter.alignBoxes(finalBoxes, currentWidth), 0);
        if (omitEmptyBoxes && asciiBox.isEmpty()) {
            return;
        }
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBox, null));
    }

    @Override
    public void writeBit(boolean value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("b%d %b%s", value ? 1 : 0, value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeBits(int numBits, byte[] bytes, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        if (!StringUtils.isBlank(additionalStringRepresentation)) {
            additionalStringRepresentation += "\n";
        }
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("%s%s", Hex.dump(bytes), additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeUnsignedByte(int bitLength, byte value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeUnsignedShort(int bitLength, short value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeUnsignedInt(int bitLength, int value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeUnsignedLong(int bitLength, long value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeUnsignedBigInteger(int bitLength, BigInteger value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeSignedByte(int bitLength, byte value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeSignedShort(int bitLength, short value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeSignedInt(int bitLength, int value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeSignedLong(int bitLength, long value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeSignedBigInteger(int bitLength, BigInteger value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeFloat(int bitLength, float value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %f%s", Float.valueOf(value).longValue(), value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeDouble(int bitLength, double value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %f%s", Double.valueOf(value).longValue(), value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeBigDecimal(int bitLength, BigDecimal value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("0x%0" + Math.max(bitLength / 4, 1) + "x %d%s", value, value, additionalStringRepresentation), 0), null));
    }

    @Override
    public void writeString(int bitLength, String value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        boxes.offerLast(new AbstractMap.SimpleEntry<>(asciiBoxWriter.boxString(name, String.format("%s%s", value, additionalStringRepresentation), 0), null));
    }

    @Override
    public WriteBuffer createSubBuffer(int numBits, WithOption... options) throws BufferException {
        return null;
    }

    @Override
    public int getPositionInBits() {
        return 0;
    }

    @Override
    public int getRemainingBits() {
        return 0;
    }

    @Override
    public byte[] getBytes() {
        return new byte[0];
    }

    public void writeVirtual(Object value, WithOption... options) throws BufferException {
        String name = getName(options);
        String additionalStringRepresentation = WithOption.extractAdditionalStringRepresentation(options).map(s -> " " + s).orElse("");
        AsciiBox virtualBox = switch (value) {
            case String s ->
                asciiBoxWriterLight.boxString(name, String.format("%s%s", value, additionalStringRepresentation), 0);
            case Float number ->
                asciiBoxWriterLight.boxString(name, String.format("%f%s", number, additionalStringRepresentation), 0);
            case Double number ->
                asciiBoxWriterLight.boxString(name, String.format("%f%s", number, additionalStringRepresentation), 0);
            case Number number ->
                // TODO: adjust rendering
                asciiBoxWriterLight.boxString(name, String.format("0x%x %d%s", number.longValue(), number.longValue(), additionalStringRepresentation), 0);
            case Boolean b ->
                asciiBoxWriterLight.boxString(name, String.format("b%d %b%s", b ? 1 : 0, value, additionalStringRepresentation), 0);
            case Enum<?> enumValue ->
                asciiBoxWriterLight.boxString(name, String.format("%s%s", enumValue.name(), additionalStringRepresentation), 0);
        /*} else if (value instanceof Serializable) {
            Serializable serializable = (Serializable) value;
            try {
                WriteBufferBoxBased writeBuffer = new WriteBufferBoxBased(true, true);
                serializable.serialize(writeBuffer);
                virtualBox = asciiBoxWriterLight.boxBox(name, writeBuffer.getBox(), 0);
            } catch (SerializationException e) {
                virtualBox = asciiBoxWriterLight.boxString(name, e.getMessage(), 0);
            }*/
            case null, default -> asciiBoxWriterLight.boxString(name, "un-renderable", 0);
        };
        boxes.offerLast(new AbstractMap.SimpleEntry<>(virtualBox, null));
    }

    protected String getName(WithOption... options) throws BufferException {
        Optional<String> name = WithOption.extractName(options, getContext());
        if (name.isEmpty()) {
            throw new BufferException("Missing 'name' option.");
        }
        return name.get();
    }

    public AsciiBox getBox() {
        assert boxes.peek() != null;
        return boxes.peek().getKey();
    }

}
