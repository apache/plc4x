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
package org.apache.plc4x.java.bacnetip.readwrite;

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

public class BACnetServiceAckGetEnrollmentSummary extends BACnetServiceAck implements Message {

  // Accessors for discriminator values.
  public BACnetConfirmedServiceChoice getServiceChoice() {
    return BACnetConfirmedServiceChoice.GET_ENROLLMENT_SUMMARY;
  }

  // Properties.
  protected final BACnetApplicationTagObjectIdentifier objectIdentifier;
  protected final BACnetEventTypeTagged eventType;
  protected final BACnetEventStateTagged eventState;
  protected final BACnetApplicationTagUnsignedInteger priority;
  protected final BACnetApplicationTagUnsignedInteger notificationClass;

  // Arguments.
  protected final Long serviceAckLength;

  public BACnetServiceAckGetEnrollmentSummary(
      BACnetApplicationTagObjectIdentifier objectIdentifier,
      BACnetEventTypeTagged eventType,
      BACnetEventStateTagged eventState,
      BACnetApplicationTagUnsignedInteger priority,
      BACnetApplicationTagUnsignedInteger notificationClass,
      Long serviceAckLength) {
    super(serviceAckLength);
    this.objectIdentifier = objectIdentifier;
    this.eventType = eventType;
    this.eventState = eventState;
    this.priority = priority;
    this.notificationClass = notificationClass;
    this.serviceAckLength = serviceAckLength;
  }

  public BACnetApplicationTagObjectIdentifier getObjectIdentifier() {
    return objectIdentifier;
  }

  public BACnetEventTypeTagged getEventType() {
    return eventType;
  }

  public BACnetEventStateTagged getEventState() {
    return eventState;
  }

  public BACnetApplicationTagUnsignedInteger getPriority() {
    return priority;
  }

  public BACnetApplicationTagUnsignedInteger getNotificationClass() {
    return notificationClass;
  }

  @Override
  protected void serializeBACnetServiceAckChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("BACnetServiceAckGetEnrollmentSummary");

    // Simple Field (objectIdentifier)
    writeSimpleField(
        "objectIdentifier", objectIdentifier, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (eventType)
    writeSimpleField("eventType", eventType, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (eventState)
    writeSimpleField("eventState", eventState, new DataWriterComplexDefault<>(writeBuffer));

    // Simple Field (priority)
    writeSimpleField("priority", priority, new DataWriterComplexDefault<>(writeBuffer));

    // Optional Field (notificationClass) (Can be skipped, if the value is null)
    writeOptionalField(
        "notificationClass", notificationClass, new DataWriterComplexDefault<>(writeBuffer));

    writeBuffer.popContext("BACnetServiceAckGetEnrollmentSummary");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    BACnetServiceAckGetEnrollmentSummary _value = this;

    // Simple field (objectIdentifier)
    lengthInBits += objectIdentifier.getLengthInBits();

    // Simple field (eventType)
    lengthInBits += eventType.getLengthInBits();

    // Simple field (eventState)
    lengthInBits += eventState.getLengthInBits();

    // Simple field (priority)
    lengthInBits += priority.getLengthInBits();

    // Optional Field (notificationClass)
    if (notificationClass != null) {
      lengthInBits += notificationClass.getLengthInBits();
    }

    return lengthInBits;
  }

  public static BACnetServiceAckGetEnrollmentSummaryBuilder staticParseBuilder(
      ReadBuffer readBuffer, Long serviceAckLength) throws ParseException {
    readBuffer.pullContext("BACnetServiceAckGetEnrollmentSummary");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    BACnetApplicationTagObjectIdentifier objectIdentifier =
        readSimpleField(
            "objectIdentifier",
            new DataReaderComplexDefault<>(
                () ->
                    (BACnetApplicationTagObjectIdentifier)
                        BACnetApplicationTag.staticParse(readBuffer),
                readBuffer));

    BACnetEventTypeTagged eventType =
        readSimpleField(
            "eventType",
            new DataReaderComplexDefault<>(
                () ->
                    BACnetEventTypeTagged.staticParse(
                        readBuffer, (short) (0), (TagClass) (TagClass.APPLICATION_TAGS)),
                readBuffer));

    BACnetEventStateTagged eventState =
        readSimpleField(
            "eventState",
            new DataReaderComplexDefault<>(
                () ->
                    BACnetEventStateTagged.staticParse(
                        readBuffer, (short) (0), (TagClass) (TagClass.APPLICATION_TAGS)),
                readBuffer));

    BACnetApplicationTagUnsignedInteger priority =
        readSimpleField(
            "priority",
            new DataReaderComplexDefault<>(
                () ->
                    (BACnetApplicationTagUnsignedInteger)
                        BACnetApplicationTag.staticParse(readBuffer),
                readBuffer));

    BACnetApplicationTagUnsignedInteger notificationClass =
        readOptionalField(
            "notificationClass",
            new DataReaderComplexDefault<>(
                () ->
                    (BACnetApplicationTagUnsignedInteger)
                        BACnetApplicationTag.staticParse(readBuffer),
                readBuffer));

    readBuffer.closeContext("BACnetServiceAckGetEnrollmentSummary");
    // Create the instance
    return new BACnetServiceAckGetEnrollmentSummaryBuilder(
        objectIdentifier, eventType, eventState, priority, notificationClass, serviceAckLength);
  }

  public static class BACnetServiceAckGetEnrollmentSummaryBuilder
      implements BACnetServiceAck.BACnetServiceAckBuilder {
    private final BACnetApplicationTagObjectIdentifier objectIdentifier;
    private final BACnetEventTypeTagged eventType;
    private final BACnetEventStateTagged eventState;
    private final BACnetApplicationTagUnsignedInteger priority;
    private final BACnetApplicationTagUnsignedInteger notificationClass;
    private final Long serviceAckLength;

    public BACnetServiceAckGetEnrollmentSummaryBuilder(
        BACnetApplicationTagObjectIdentifier objectIdentifier,
        BACnetEventTypeTagged eventType,
        BACnetEventStateTagged eventState,
        BACnetApplicationTagUnsignedInteger priority,
        BACnetApplicationTagUnsignedInteger notificationClass,
        Long serviceAckLength) {

      this.objectIdentifier = objectIdentifier;
      this.eventType = eventType;
      this.eventState = eventState;
      this.priority = priority;
      this.notificationClass = notificationClass;
      this.serviceAckLength = serviceAckLength;
    }

    public BACnetServiceAckGetEnrollmentSummary build(Long serviceAckLength) {

      BACnetServiceAckGetEnrollmentSummary bACnetServiceAckGetEnrollmentSummary =
          new BACnetServiceAckGetEnrollmentSummary(
              objectIdentifier,
              eventType,
              eventState,
              priority,
              notificationClass,
              serviceAckLength);
      return bACnetServiceAckGetEnrollmentSummary;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BACnetServiceAckGetEnrollmentSummary)) {
      return false;
    }
    BACnetServiceAckGetEnrollmentSummary that = (BACnetServiceAckGetEnrollmentSummary) o;
    return (getObjectIdentifier() == that.getObjectIdentifier())
        && (getEventType() == that.getEventType())
        && (getEventState() == that.getEventState())
        && (getPriority() == that.getPriority())
        && (getNotificationClass() == that.getNotificationClass())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getObjectIdentifier(),
        getEventType(),
        getEventState(),
        getPriority(),
        getNotificationClass());
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
