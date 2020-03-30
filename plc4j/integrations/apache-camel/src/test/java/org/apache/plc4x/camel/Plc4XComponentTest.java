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

import org.apache.camel.Expression;
import org.apache.camel.builder.RouteBuilder;
import org.apache.camel.component.mock.MockEndpoint;
import org.apache.camel.test.junit5.CamelTestSupport;
import org.apache.plc4x.java.api.model.PlcField;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.TimeUnit;

public class Plc4XComponentTest extends CamelTestSupport {

    @Test
    public void testSimpleRouting() throws Exception {
        MockEndpoint mock = getMockEndpoint("mock:result");
        mock.expectedMinimumMessageCount(1);
        mock.expectedMessageCount(2);

        template.asyncSendBody("direct:plc4x", Collections.singletonList("irrelevant"));
        template.asyncSendBody("direct:plc4x2", Collections.singletonList("irrelevant"));

        assertMockEndpointsSatisfied(2, TimeUnit.SECONDS);
    }

    @Override
    protected RouteBuilder createRouteBuilder() {
        return new RouteBuilder() {
            public void configure() {
                List<TagData> tags = new ArrayList<>();
                tags.add(new TagData("testTagName","testTagAddress"));
                tags.add(new TagData("testTagName2","testTagAddress2"));
                Plc4XEndpoint producer = getContext().getEndpoint("plc4x:mock:10.10.10.1/1/1", Plc4XEndpoint.class);
                producer.setTags(tags);
                from("direct:plc4x")
                    .setBody(constant(Arrays.asList(new TagData("test","testAddress",false))))
                    .to("plc4x:mock:10.10.10.1/1/1")
                    .to("mock:result");
                from("direct:plc4x2")
                    .setBody(constant(Arrays.asList(new TagData("test2","testAddress2",0x05))))
                    .to("plc4x:mock:10.10.10.1/1/1")
                    .to("mock:result");
                from(producer)
                    .log("Got ${body}");
            }
        };
    }

}