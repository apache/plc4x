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

public class BACnetConfirmedServiceRequestGetEnrollmentSummary extends BACnetConfirmedServiceRequest
    implements Message {

  // Accessors for discriminator values.
  public BACnetConfirmedServiceChoice getServiceChoice() {
    return BACnetConfirmedServiceChoice.GET_ENROLLMENT_SUMMARY;
  }

  // Properties.
  protected final BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterTagged
      acknowledgmentFilter;
  protected final BACnetRecipientProcessEnclosed enrollmentFilter;
  protected final BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilterTagged
      eventStateFilter;
  protected final BACnetEventTypeTagged eventTypeFilter;
  protected final BACnetConfirmedServiceRequestGetEnrollmentSummaryPriorityFilter priorityFilter;
  protected final BACnetContextTagUnsignedInteger notificationClassFilter;

  // Arguments.
  protected final Long serviceRequestLength;

  public BACnetConfirmedServiceRequestGetEnrollmentSummary(
      BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterTagged
          acknowledgmentFilter,
      BACnetRecipientProcessEnclosed enrollmentFilter,
      BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilterTagged eventStateFilter,
      BACnetEventTypeTagged eventTypeFilter,
      BACnetConfirmedServiceRequestGetEnrollmentSummaryPriorityFilter priorityFilter,
      BACnetContextTagUnsignedInteger notificationClassFilter,
      Long serviceRequestLength) {
    super(serviceRequestLength);
    this.acknowledgmentFilter = acknowledgmentFilter;
    this.enrollmentFilter = enrollmentFilter;
    this.eventStateFilter = eventStateFilter;
    this.eventTypeFilter = eventTypeFilter;
    this.priorityFilter = priorityFilter;
    this.notificationClassFilter = notificationClassFilter;
    this.serviceRequestLength = serviceRequestLength;
  }

  public BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterTagged
      getAcknowledgmentFilter() {
    return acknowledgmentFilter;
  }

  public BACnetRecipientProcessEnclosed getEnrollmentFilter() {
    return enrollmentFilter;
  }

  public BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilterTagged
      getEventStateFilter() {
    return eventStateFilter;
  }

  public BACnetEventTypeTagged getEventTypeFilter() {
    return eventTypeFilter;
  }

  public BACnetConfirmedServiceRequestGetEnrollmentSummaryPriorityFilter getPriorityFilter() {
    return priorityFilter;
  }

  public BACnetContextTagUnsignedInteger getNotificationClassFilter() {
    return notificationClassFilter;
  }

  @Override
  protected void serializeBACnetConfirmedServiceRequestChild(WriteBuffer writeBuffer)
      throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("BACnetConfirmedServiceRequestGetEnrollmentSummary");

    // Simple Field (acknowledgmentFilter)
    writeSimpleField(
        "acknowledgmentFilter", acknowledgmentFilter, new DataWriterComplexDefault<>(writeBuffer));

    // Optional Field (enrollmentFilter) (Can be skipped, if the value is null)
    writeOptionalField(
        "enrollmentFilter", enrollmentFilter, new DataWriterComplexDefault<>(writeBuffer));

    // Optional Field (eventStateFilter) (Can be skipped, if the value is null)
    writeOptionalField(
        "eventStateFilter", eventStateFilter, new DataWriterComplexDefault<>(writeBuffer));

    // Optional Field (eventTypeFilter) (Can be skipped, if the value is null)
    writeOptionalField(
        "eventTypeFilter", eventTypeFilter, new DataWriterComplexDefault<>(writeBuffer));

    // Optional Field (priorityFilter) (Can be skipped, if the value is null)
    writeOptionalField(
        "priorityFilter", priorityFilter, new DataWriterComplexDefault<>(writeBuffer));

    // Optional Field (notificationClassFilter) (Can be skipped, if the value is null)
    writeOptionalField(
        "notificationClassFilter",
        notificationClassFilter,
        new DataWriterComplexDefault<>(writeBuffer));

    writeBuffer.popContext("BACnetConfirmedServiceRequestGetEnrollmentSummary");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = super.getLengthInBits();
    BACnetConfirmedServiceRequestGetEnrollmentSummary _value = this;

    // Simple field (acknowledgmentFilter)
    lengthInBits += acknowledgmentFilter.getLengthInBits();

    // Optional Field (enrollmentFilter)
    if (enrollmentFilter != null) {
      lengthInBits += enrollmentFilter.getLengthInBits();
    }

    // Optional Field (eventStateFilter)
    if (eventStateFilter != null) {
      lengthInBits += eventStateFilter.getLengthInBits();
    }

    // Optional Field (eventTypeFilter)
    if (eventTypeFilter != null) {
      lengthInBits += eventTypeFilter.getLengthInBits();
    }

    // Optional Field (priorityFilter)
    if (priorityFilter != null) {
      lengthInBits += priorityFilter.getLengthInBits();
    }

    // Optional Field (notificationClassFilter)
    if (notificationClassFilter != null) {
      lengthInBits += notificationClassFilter.getLengthInBits();
    }

    return lengthInBits;
  }

  public static BACnetConfirmedServiceRequestBuilder
      staticParseBACnetConfirmedServiceRequestBuilder(
          ReadBuffer readBuffer, Long serviceRequestLength) throws ParseException {
    readBuffer.pullContext("BACnetConfirmedServiceRequestGetEnrollmentSummary");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterTagged
        acknowledgmentFilter =
            readSimpleField(
                "acknowledgmentFilter",
                new DataReaderComplexDefault<>(
                    () ->
                        BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterTagged
                            .staticParse(
                                readBuffer,
                                (short) (0),
                                (TagClass) (TagClass.CONTEXT_SPECIFIC_TAGS)),
                    readBuffer));

    BACnetRecipientProcessEnclosed enrollmentFilter =
        readOptionalField(
            "enrollmentFilter",
            new DataReaderComplexDefault<>(
                () -> BACnetRecipientProcessEnclosed.staticParse(readBuffer, (short) (1)),
                readBuffer));

    BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilterTagged eventStateFilter =
        readOptionalField(
            "eventStateFilter",
            new DataReaderComplexDefault<>(
                () ->
                    BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilterTagged
                        .staticParse(
                            readBuffer, (short) (2), (TagClass) (TagClass.CONTEXT_SPECIFIC_TAGS)),
                readBuffer));

    BACnetEventTypeTagged eventTypeFilter =
        readOptionalField(
            "eventTypeFilter",
            new DataReaderComplexDefault<>(
                () ->
                    BACnetEventTypeTagged.staticParse(
                        readBuffer, (short) (3), (TagClass) (TagClass.CONTEXT_SPECIFIC_TAGS)),
                readBuffer));

    BACnetConfirmedServiceRequestGetEnrollmentSummaryPriorityFilter priorityFilter =
        readOptionalField(
            "priorityFilter",
            new DataReaderComplexDefault<>(
                () ->
                    BACnetConfirmedServiceRequestGetEnrollmentSummaryPriorityFilter.staticParse(
                        readBuffer, (short) (4)),
                readBuffer));

    BACnetContextTagUnsignedInteger notificationClassFilter =
        readOptionalField(
            "notificationClassFilter",
            new DataReaderComplexDefault<>(
                () ->
                    (BACnetContextTagUnsignedInteger)
                        BACnetContextTag.staticParse(
                            readBuffer,
                            (short) (5),
                            (BACnetDataType) (BACnetDataType.UNSIGNED_INTEGER)),
                readBuffer));

    readBuffer.closeContext("BACnetConfirmedServiceRequestGetEnrollmentSummary");
    // Create the instance
    return new BACnetConfirmedServiceRequestGetEnrollmentSummaryBuilderImpl(
        acknowledgmentFilter,
        enrollmentFilter,
        eventStateFilter,
        eventTypeFilter,
        priorityFilter,
        notificationClassFilter,
        serviceRequestLength);
  }

  public static class BACnetConfirmedServiceRequestGetEnrollmentSummaryBuilderImpl
      implements BACnetConfirmedServiceRequest.BACnetConfirmedServiceRequestBuilder {
    private final BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterTagged
        acknowledgmentFilter;
    private final BACnetRecipientProcessEnclosed enrollmentFilter;
    private final BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilterTagged
        eventStateFilter;
    private final BACnetEventTypeTagged eventTypeFilter;
    private final BACnetConfirmedServiceRequestGetEnrollmentSummaryPriorityFilter priorityFilter;
    private final BACnetContextTagUnsignedInteger notificationClassFilter;
    private final Long serviceRequestLength;

    public BACnetConfirmedServiceRequestGetEnrollmentSummaryBuilderImpl(
        BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterTagged
            acknowledgmentFilter,
        BACnetRecipientProcessEnclosed enrollmentFilter,
        BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilterTagged eventStateFilter,
        BACnetEventTypeTagged eventTypeFilter,
        BACnetConfirmedServiceRequestGetEnrollmentSummaryPriorityFilter priorityFilter,
        BACnetContextTagUnsignedInteger notificationClassFilter,
        Long serviceRequestLength) {
      this.acknowledgmentFilter = acknowledgmentFilter;
      this.enrollmentFilter = enrollmentFilter;
      this.eventStateFilter = eventStateFilter;
      this.eventTypeFilter = eventTypeFilter;
      this.priorityFilter = priorityFilter;
      this.notificationClassFilter = notificationClassFilter;
      this.serviceRequestLength = serviceRequestLength;
    }

    public BACnetConfirmedServiceRequestGetEnrollmentSummary build(Long serviceRequestLength) {

      BACnetConfirmedServiceRequestGetEnrollmentSummary
          bACnetConfirmedServiceRequestGetEnrollmentSummary =
              new BACnetConfirmedServiceRequestGetEnrollmentSummary(
                  acknowledgmentFilter,
                  enrollmentFilter,
                  eventStateFilter,
                  eventTypeFilter,
                  priorityFilter,
                  notificationClassFilter,
                  serviceRequestLength);
      return bACnetConfirmedServiceRequestGetEnrollmentSummary;
    }
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof BACnetConfirmedServiceRequestGetEnrollmentSummary)) {
      return false;
    }
    BACnetConfirmedServiceRequestGetEnrollmentSummary that =
        (BACnetConfirmedServiceRequestGetEnrollmentSummary) o;
    return (getAcknowledgmentFilter() == that.getAcknowledgmentFilter())
        && (getEnrollmentFilter() == that.getEnrollmentFilter())
        && (getEventStateFilter() == that.getEventStateFilter())
        && (getEventTypeFilter() == that.getEventTypeFilter())
        && (getPriorityFilter() == that.getPriorityFilter())
        && (getNotificationClassFilter() == that.getNotificationClassFilter())
        && super.equals(that)
        && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(
        super.hashCode(),
        getAcknowledgmentFilter(),
        getEnrollmentFilter(),
        getEventStateFilter(),
        getEventTypeFilter(),
        getPriorityFilter(),
        getNotificationClassFilter());
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
