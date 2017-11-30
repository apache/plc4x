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
package org.apache.plc4x.camel.util;

import co.unruly.matchers.StreamMatchers;
import org.junit.Test;

import java.util.Optional;
import java.util.stream.Stream;

import static org.junit.Assert.assertThat;

public class StreamUtilsTest {
    @Test
    public void streamOf_One_element_should_exaclty_contains_one_element() throws Exception {
        Stream<String> stringStream = StreamUtils.streamOf(Optional.of("String"));
        assertThat(stringStream, StreamMatchers.equalTo(Stream.of("String")));
    }

    @Test
    public void streamOf_should_be_empty() throws Exception {
        Stream<String> stringStream = StreamUtils.streamOf(Optional.empty());
        assertThat(stringStream, StreamMatchers.empty());
    }

}