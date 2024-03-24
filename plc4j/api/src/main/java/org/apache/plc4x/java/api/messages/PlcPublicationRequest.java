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
package org.apache.plc4x.java.api.messages;

import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.PlcValue;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;

public interface PlcPublicationRequest extends PlcPublicationTagRequest {

    @Override
    CompletableFuture<? extends PlcPublicationResponse> execute();

    interface Builder extends PlcRequestBuilder {

        @Override
        PlcPublicationRequest build();

        /**
         * Adds a new tag to the to be constructed request which should be published cyclically.
         * In this case will the driver regularly publish the given value, if it has changed or not.
         *
         * @param name                alias of the tag.
         * @param tagAddress          tag address string for accessing the tag.
         * @param publicationInterval interval, in which the tag should be published.
         * @param initialValue        initial value of the tag
         * @return builder.
         */
        PlcPublicationRequest.Builder addCyclicTagAddress(String name, String tagAddress, Duration publicationInterval, PlcValue initialValue);

        /**
         * Adds a new tag to the to be constructed request which should be published cyclically.
         * In this case will the driver regularly publish the given value, if it has changed or not.
         *
         * @param name                alias of the tag.
         * @param tag                 tag instance for accessing the tag.
         * @param publicationInterval interval, in which the tag should be published.
         * @param initialValue        initial value of the tag
         * @return builder.
         */
        PlcPublicationRequest.Builder addCyclicTag(String name, PlcTag tag, Duration publicationInterval, PlcValue initialValue);

        /**
         * Adds a new tag to the to be constructed request which should be published as soon as
         * a value changes locally.
         *
         * @param name         alias of the tag.
         * @param tagAddress   tag address string for accessing the tag.
         * @param initialValue initial value of the tag
         * @return builder.
         */
        PlcPublicationRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress, PlcValue initialValue);

        /**
         * Adds a new tag to the to be constructed request which should be published as soon as
         * a value changes locally.
         *
         * @param name         alias of the tag.
         * @param tag          tag instance for accessing the tag.
         * @param initialValue initial value of the tag
         * @return builder.
         */
        PlcPublicationRequest.Builder addChangeOfStateTag(String name, PlcTag tag, PlcValue initialValue);

    }

}
