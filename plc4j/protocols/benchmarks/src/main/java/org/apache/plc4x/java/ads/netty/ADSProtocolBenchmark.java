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

package org.apache.plc4x.java.ads.netty;

import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.ads.api.commands.ADSReadRequest;
import org.apache.plc4x.java.ads.api.commands.ADSReadResponse;
import org.apache.plc4x.java.ads.api.commands.ADSWriteRequest;
import org.apache.plc4x.java.ads.api.commands.ADSWriteResponse;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.types.AMSNetId;
import org.apache.plc4x.java.ads.api.generic.types.AMSPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.openjdk.jmh.annotations.*;
import org.openjdk.jmh.infra.Blackhole;

import java.util.LinkedList;
import java.util.concurrent.TimeUnit;

public class ADSProtocolBenchmark {

    @State(Scope.Benchmark)
    public static class MyState {
        ADSProtocol SUT = new ADSProtocol();
        ADSWriteRequest adsWriteRequest = buildAdsWriteRequest();
        byte[] adsWriteResponse = buildADSWriteResponse();
        ADSReadRequest adsReadRequest = buildAdsReadRequest();
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
    public void measureEncodingAdsWriteRequest(Blackhole blackhole, MyState myState) throws Exception {
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
    public void measureDecodingAdsWriteResponse(Blackhole blackhole, MyState myState) throws Exception {
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
    public void measureEncodingAdsReadRequest(Blackhole blackhole, MyState myState) throws Exception {
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
    public void measureDecodingAdsReadResponse(Blackhole blackhole, MyState myState) throws Exception {
        LinkedList<Object> out = new LinkedList<>();
        myState.SUT.decode(null, Unpooled.wrappedBuffer(myState.adsReadResponse), out);
        blackhole.consume(out.remove());
    }

    private static ADSWriteRequest buildAdsWriteRequest() {
        AMSNetId targetAmsNetId = AMSNetId.of("1.2.3.4.5.6");
        AMSPort targetAmsPort = AMSPort.of(7);
        AMSNetId sourceAmsNetId = AMSNetId.of("8.9.10.11.12.13");
        AMSPort sourceAmsPort = AMSPort.of(14);
        Invoke invokeId = Invoke.of(15);
        Data data = Data.of("Hello World!".getBytes());
        return ADSWriteRequest.of(
            targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
            IndexGroup.of(1),
            IndexOffset.of(1),
            data
        );
    }

    private static byte[] buildADSWriteResponse() {
        AMSNetId targetAmsNetId = AMSNetId.of("1.2.3.4.5.6");
        AMSPort targetAmsPort = AMSPort.of(7);
        AMSNetId sourceAmsNetId = AMSNetId.of("8.9.10.11.12.13");
        AMSPort sourceAmsPort = AMSPort.of(14);
        Invoke invokeId = Invoke.of(15);
        Result result = Result.of(AdsReturnCode.ADS_CODE_0);
        return ADSWriteResponse.of(
            targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
            result
        ).getBytes();
    }

    private static ADSReadRequest buildAdsReadRequest() {
        AMSNetId targetAmsNetId = AMSNetId.of("1.2.3.4.5.6");
        AMSPort targetAmsPort = AMSPort.of(7);
        AMSNetId sourceAmsNetId = AMSNetId.of("8.9.10.11.12.13");
        AMSPort sourceAmsPort = AMSPort.of(14);
        Invoke invokeId = Invoke.of(15);
        return ADSReadRequest.of(
            targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
            IndexGroup.of(1),
            IndexOffset.of(1),
            Length.of(1)
        );
    }

    private static byte[] buildADSReadResponse() {
        AMSNetId targetAmsNetId = AMSNetId.of("1.2.3.4.5.6");
        AMSPort targetAmsPort = AMSPort.of(7);
        AMSNetId sourceAmsNetId = AMSNetId.of("8.9.10.11.12.13");
        AMSPort sourceAmsPort = AMSPort.of(14);
        Invoke invokeId = Invoke.of(15);
        Result result = Result.of(AdsReturnCode.ADS_CODE_0);
        Data data = Data.of("Hello World!".getBytes());
        return ADSReadResponse.of(
            targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, invokeId,
            result,
            data
        ).getBytes();
    }
}
