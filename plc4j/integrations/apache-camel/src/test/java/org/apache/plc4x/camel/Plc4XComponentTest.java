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
package org.apache.plc4x.camel;

import org.apache.camel.builder.RouteBuilder;
import org.apache.camel.component.mock.MockEndpoint;
import org.apache.camel.test.junit4.CamelTestSupport;
import org.apache.plc4x.java.api.model.PlcField;
import org.junit.Test;

import java.util.Arrays;
import java.util.concurrent.TimeUnit;

public class Plc4XComponentTest extends CamelTestSupport {

    @Test
    public void testSimpleRouting() throws Exception {
        MockEndpoint mock = getMockEndpoint("mock:result");
        mock.expectedMinimumMessageCount(1);
        mock.expectedMessageCount(2);

        template.asyncSendBody("direct:plc4x", "irrelevant");
        template.asyncSendBody("direct:plc4x2", "irrelevant");

        assertMockEndpointsSatisfied(2, TimeUnit.SECONDS);
    }

    @Override
    protected RouteBuilder createRouteBuilder() {
        return new RouteBuilder() {
            public void configure() {
                from("direct:plc4x")
                    .setHeader(Constants.FIELD_QUERY_HEADER, constant(new PlcField() {
                    }))
                    .setBody(constant((byte) 0x0))
                    .to("plc4x:mock:10.10.10.1/1/1")
                    .to("mock:result");
                from("direct:plc4x2")
                    .setHeader(Constants.FIELD_QUERY_HEADER, constant(new PlcField() {
                    }))
                    .setBody(constant(Arrays.asList((byte) 0x0, (byte) 0x1, (byte) 0x2, (byte) 0x3)))
                    .to("plc4x:mock:10.10.10.1/1/1")
                    .to("mock:result");
                from("plc4x:mock:10.10.10.1/1/1?address=Main.by0&dataType=java.lang.String")
                    .log("Got ${body}");
            }
        };
    }

}