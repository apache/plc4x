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
package org.apache.plc4x.java.cbus.readwrite;

import static org.apache.plc4x.java.spi.codegen.fields.FieldReaderFactory.*;
import static org.apache.plc4x.java.spi.codegen.fields.FieldWriterFactory.*;
import static org.apache.plc4x.java.spi.codegen.io.DataReaderFactory.*;
import static org.apache.plc4x.java.spi.codegen.io.DataWriterFactory.*;
import static org.apache.plc4x.java.spi.generation.StaticHelper.*;

import java.time.*;
import java.util.*;
import org.apache.plc4x.java.api.exceptions.*;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.codegen.*;
import org.apache.plc4x.java.spi.codegen.fields.*;
import org.apache.plc4x.java.spi.codegen.io.*;
import org.apache.plc4x.java.spi.generation.*;

// Code generated by code-generation. DO NOT EDIT.

public class MediaTransportControlDataEnumerateCategoriesSelectionTracks
    extends MediaTransportControlData implements Message {

  // Accessors for discriminator values.

  // Properties.
  protected final byte enumerateType;
  protected final short start;

  public MediaTransportControlDataEnumerateCategoriesSelectionTracks(
      MediaTransportControlCommandTypeContainer commandTypeContainer,
      byte mediaLinkGroup,
      byte enumerateType,
      short start) {
    super(commandTypeContainer, mediaLinkGroup);
    this.enumerateType = enumerateType;
    this.start = start;
  }

  public byte getEnumerateType() {
    return enumerateType;
  }

  public short getStart() {
    return start;
  }

  public boolean getIsListCategories() {
    return (boolean) ((getEnumerateType()) == (0x00));
  }

  public boolean getIsListSelections() {
    return (boolean) ((getEnumerateType()) == (0x01));
  }

  public boolean getIsListTracks() {
    return (boolean) ((getEnumerateType()) == (0x02));
  }

  public boolean getIsReserved() {
    return (boolean)
        (((!(getIsListCategories())) && (!(getIsListSelections()))) && (!(getIsListTracks())));
  }

  @Override
  protected void serializeMediaTransportControlDataChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("MediaTransportControlDataEnumerateCategoriesSelectionTracks");

    // Simple Field (enumerateType)
    writeSimpleField("enumerateType", enumerateType, writeByte(writeBuffer, 8));

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    boolean isListCategories = getIsListCategories();
    writeBuffer.writeVirtual("isListCategories", isListCategories);

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    boolean isListSelections = getIsListSelections();
    writeBuffer.writeVirtual("isListSelections", isListSelections);

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    boolean isListTracks = getIsListTracks();
    writeBuffer.writeVirtual("isListTracks", isListTracks);

    // Virtual field (doesn't actually serialize anything, just makes the value available)
    boolean isReserved = getIsReserved();
    writeBuffer.writeVirtual("isReserved", isReserved);

    // Simple Field (start)
    writeSimpleField("start", start, writeUnsignedShort(writeBuffer, 8));

    writeBuffer.popContext("MediaTransportControlDataEnumerateCategoriesSelectionTracks");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    MediaTransportControlDataEnumerateCategoriesSelectionTracks _value = this;

    // Simple field (enumerateType)
    lengthInBits += 8;

    // A virtual field doesn't have any in- or output.

    // A virtual field doesn't have any in- or output.

    // A virtual field doesn't have any in- or output.

    // A virtual field doesn't have any in- or output.

    // Simple field (start)
    lengthInBits += 8;

    return lengthInBits;
  }

  public static MediaTransportControlDataEnumerateCategoriesSelectionTracksBuilder
      staticParseBuilder(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("MediaTransportControlDataEnumerateCategoriesSelectionTracks");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    byte enumerateType = readSimpleField("enumerateType", readByte(readBuffer, 8));
    boolean isListCategories =
        readVirtualField("isListCategories", boolean.class, (enumerateType) == (0x00));
    boolean isListSelections =
        readVirtualField("isListSelections", boolean.class, (enumerateType) == (0x01));
    boolean isListTracks =
        readVirtualField("isListTracks", boolean.class, (enumerateType) == (0x02));
    boolean isReserved =
        readVirtualField(
            "isReserved",
            boolean.class,
            ((!(isListCategories)) && (!(isListSelections))) && (!(isListTracks)));

    short start = readSimpleField("start", readUnsignedShort(readBuffer, 8));

    readBuffer.closeContext("MediaTransportControlDataEnumerateCategoriesSelectionTracks");
    // Create the instance
    return new MediaTransportControlDataEnumerateCategoriesSelectionTracksBuilder(
        enumerateType, start);
  }

  public static class MediaTransportControlDataEnumerateCategoriesSelectionTracksBuilder
      implements MediaTransportControlData.MediaTransportControlDataBuilder {
    private final byte enumerateType;
    private final short start;

    public MediaTransportControlDataEnumerateCategoriesSelectionTracksBuilder(
        byte enumerateType, short start) {

      this.enumerateType = enumerateType;
      this.start = start;
    }

    public MediaTransportControlDataEnumerateCategoriesSelectionTracks build(
        MediaTransportControlCommandTypeContainer commandTypeContainer, byte mediaLinkGroup) {
      MediaTransportControlDataEnumerateCategoriesSelectionTracks
          mediaTransportControlDataEnumerateCategoriesSelectionTracks =
              new MediaTransportControlDataEnumerateCategoriesSelectionTracks(
                  commandTypeContainer, mediaLinkGroup, enumerateType, start);
      return mediaTransportControlDataEnumerateCategoriesSelectionTracks;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof MediaTransportControlDataEnumerateCategoriesSelectionTracks)) {
      return false;
    }
    MediaTransportControlDataEnumerateCategoriesSelectionTracks that =
        (MediaTransportControlDataEnumerateCategoriesSelectionTracks) o;
    return (getEnumerateType() == that.getEnumerateType())
        && (getStart() == that.getStart())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(super.hashCode(), getEnumerateType(), getStart());
  }

  @Override
  public String toString() {
    WriteBufferBoxBased writeBufferBoxBased = new WriteBufferBoxBased(true, true);
    try {
      writeBufferBoxBased.writeSerializable(this);
    } catch (SerializationException e) {
      throw new RuntimeException(e);
    }
    return "\n" + writeBufferBoxBased.getBox().toString() + "\n";
  }
}
