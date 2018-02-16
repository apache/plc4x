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

import org.apache.camel.Component;
import org.junit.Before;
import org.junit.Test;

import static org.hamcrest.core.Is.is;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.junit.Assert.assertThat;
import static org.mockito.Mockito.mock;

public class PLC4XEndpointTest {

    PLC4XEndpoint SUT;

    @Before
    public void setUp() throws Exception {
        SUT = new PLC4XEndpoint("plc4x:mock:10.10.10.1/1/1", mock(Component.class));
    }

    @Test
    public void createProducer() throws Exception {
        assertThat(SUT.createProducer(), notNullValue());
    }

    @Test(expected = UnsupportedOperationException.class)
    public void createConsumer() throws Exception {
        SUT.createConsumer(null);
    }

    @Test
    public void isSingleton() throws Exception {
        assertThat(SUT.isSingleton(), is(true));
    }

}