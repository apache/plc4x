/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.ads.adslib;

import org.apache.plc4x.java.ads.api.generic.AmsPacket;
import org.apache.plc4x.java.base.messages.PlcProprietaryRequest;
import org.apache.plc4x.java.base.messages.DefaultPlcProprietaryRequest;

import java.util.concurrent.CompletableFuture;

/**
 * Ported from <a href="https://github.com/Beckhoff/ADS">github AdsLib</a>
 */
public class AmsRequest<REQUEST extends AmsPacket, RESPONSE extends AmsPacket> {

    private final PlcProprietaryRequest<REQUEST> request;
    private final CompletableFuture<RESPONSE> responseFuture;

    private AmsRequest(REQUEST amsPacket, CompletableFuture<RESPONSE> responseFuture) {
        this.request = new DefaultPlcProprietaryRequest<>(amsPacket);
        this.responseFuture = responseFuture;
    }

    public static <REQUEST extends AmsPacket, RESPONSE extends AmsPacket> AmsRequest<REQUEST, RESPONSE> of(REQUEST amsPacket) {
        return new AmsRequest<>(amsPacket, new CompletableFuture<>());
    }

    public PlcProprietaryRequest<REQUEST> getRequest() {
        return request;
    }

    public CompletableFuture<RESPONSE> getResponseFuture() {
        return responseFuture;
    }
}
