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
package org.apache.plc4x.java.ads.protocol;

import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.ads.api.commands.AdsReadRequest;
import org.apache.plc4x.java.ads.api.commands.AdsReadResponse;
import org.apache.plc4x.java.ads.api.commands.AdsWriteRequest;
import org.apache.plc4x.java.ads.api.commands.AdsWriteResponse;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.openjdk.jmh.annotations.*;
import org.openjdk.jmh.infra.Blackhole;

import java.util.LinkedList;
import java.util.concurrent.TimeUnit;

public class ADSProtocolBenchmark {

    @State(Scope.Benchmark)
    public static class MyState {
        Ads2PayloadProtocol SUT = new Ads2PayloadProtocol();
        AdsWriteRequest adsWriteRequest = buildAdsWriteRequest();
        byte[] adsWriteResponse = buildADSWriteResponse();
        AdsReadRequest adsReadRequest = buildAdsReadRequest();
        byte[] adsReadResponse = buildADSReadResponse();

        @Setup(Level.Trial)
        public void doSetup() {
            System.out.println("Do Setup");
        }

        @TearDown(Level.Iteration)
        public void doTearDown() {
            System.out.println("Do TearDown");
            SUT.reset();
        }

    }

    @Benchmark
    @Warmup(iterations = 5, time = 1)
    @Measurement(iterations = 5, time = 1)
    @Fork(3)
    @BenchmarkMode(Mode.AverageTime)
    @OutputTimeUnit(TimeUnit.NANOSECONDS)
    public void measureEncodingAdsWriteRequest(Blackhole blackhole, MyState myState) {
        LinkedList<Object> out = new LinkedList<>();
        myState.SUT.encode(null, myState.adsWriteRequest, out);
        blackhole.consume(out.remove());
    }

    @Benchmark
    @Warmup(iterations = 5, time = 1)
    @Measurement(iterations = 5, time = 1)
    @Fork(3)
    @BenchmarkMode(Mode.AverageTime)
    @OutputTimeUnit(TimeUnit.NANOSECONDS)
    public void measureDecodingAdsWriteResponse(Blackhole blackhole, MyState myState) {
        LinkedList<Object> out = new LinkedList<>();
        myState.SUT.decode(null, Unpooled.wrappedBuffer(myState.adsWriteResponse), out);
        blackhole.consume(out.remove());
    }

    @Benchmark
    @Warmup(iterations = 5, time = 1)
    @Measurement(iterations = 5, time = 1)
    @Fork(3)
    @BenchmarkMode(Mode.AverageTime)
    @OutputTimeUnit(TimeUnit.NANOSECONDS)
    public void measureEncodingAdsReadRequest(Blackhole blackhole, MyState myState) {
        LinkedList<Object> out = new LinkedList<>();
        myState.SUT.encode(null, myState.adsReadRequest, out);
        blackhole.consume(out.remove());
    }

    @Benchmark
    @Warmup(iterations = 5, time = 1)
    @Measurement(iterations = 5, time = 1)
    @Fork(3)
    @BenchmarkMode(Mode.AverageTime)
    @OutputTimeUnit(TimeUnit.NANOSECONDS)
    public void measureDecodingAdsReadResponse(Blackhole blackhole, MyState myState) {
        LinkedList<Object> out = new LinkedList<>();
        myState.SUT.decode(null, Unpooled.wrappedBuffer(myState.adsReadResponse), out);
        blackhole.consume(out.remove());
    }

    private static AdsWriteRequest buildAdsWriteRequest() {
        AmsNetId targetAmsNetId = AmsNetId.of("1.2.3.4.5.6");
        AmsPort targetAmsPort = AmsPort.of(7);
        AmsNetId sourceAmsNetId = AmsNetId.of("8.9.10.11.12.13");
        AmsPort sourceAmsPort = AmsPort.of(14);
        Invoke invokeId = Invoke.of(15);
        Data data = Data.of("Hello World!".getBytes());
        return AdsWriteRequest.of(
            targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
            IndexGroup.of(1),
            IndexOffset.of(1),
            data
        );
    }

    private static byte[] buildADSWriteResponse() {
        AmsNetId targetAmsNetId = AmsNetId.of("1.2.3.4.5.6");
        AmsPort targetAmsPort = AmsPort.of(7);
        AmsNetId sourceAmsNetId = AmsNetId.of("8.9.10.11.12.13");
        AmsPort sourceAmsPort = AmsPort.of(14);
        Invoke invokeId = Invoke.of(15);
        Result result = Result.of(AdsReturnCode.ADS_CODE_0);
        return AdsWriteResponse.of(
            targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
            result
        ).getBytes();
    }

    private static AdsReadRequest buildAdsReadRequest() {
        AmsNetId targetAmsNetId = AmsNetId.of("1.2.3.4.5.6");
        AmsPort targetAmsPort = AmsPort.of(7);
        AmsNetId sourceAmsNetId = AmsNetId.of("8.9.10.11.12.13");
        AmsPort sourceAmsPort = AmsPort.of(14);
        Invoke invokeId = Invoke.of(15);
        return AdsReadRequest.of(
            targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
            IndexGroup.of(1),
            IndexOffset.of(1),
            Length.of(1)
        );
    }

    private static byte[] buildADSReadResponse() {
        AmsNetId targetAmsNetId = AmsNetId.of("1.2.3.4.5.6");
        AmsPort targetAmsPort = AmsPort.of(7);
        AmsNetId sourceAmsNetId = AmsNetId.of("8.9.10.11.12.13");
        AmsPort sourceAmsPort = AmsPort.of(14);
        Invoke invokeId = Invoke.of(15);
        Result result = Result.of(AdsReturnCode.ADS_CODE_0);
        Data data = Data.of("Hello World!".getBytes());
        return AdsReadResponse.of(
            targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
            result,
            data
        ).getBytes();
    }
}
