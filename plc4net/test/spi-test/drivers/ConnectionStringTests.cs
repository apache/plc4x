//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

using System;
using System.Globalization;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.spi.drivers;
using Xunit;

namespace org.apache.plc4net.spi.test.drivers
{
    public class ConnectionStringTests
    {
        [Fact]
        public void Parses_the_simple_form()
        {
            var cs = ConnectionString.Parse("s7://192.168.0.1");

            Assert.Equal("s7", cs.ProtocolCode);
            Assert.Null(cs.TransportCode);
            Assert.Equal("192.168.0.1", cs.TransportConfig);
            Assert.Equal(string.Empty, cs.ParamString);
        }

        [Fact]
        public void Parses_the_form_with_an_explicit_transport()
        {
            var cs = ConnectionString.Parse("s7:cotp://10.0.0.5:102");

            Assert.Equal("s7", cs.ProtocolCode);
            Assert.Equal("cotp", cs.TransportCode);
            Assert.Equal("10.0.0.5:102", cs.TransportConfig);
        }

        [Fact]
        public void Parses_query_parameters()
        {
            var cs = ConnectionString.Parse("modbus-tcp://10.0.0.9:502?unit-identifier=1&retry=true");

            Assert.Equal("1", cs.GetParameter("unit-identifier"));
            Assert.Equal("true", cs.GetParameter("retry"));
        }

        [Fact]
        public void URL_decodes_parameter_values()
        {
            var cs = ConnectionString.Parse("s7://host?remote-rack=0&remote-slot=1");

            Assert.Equal("0", cs.GetParameter("remote-rack"));
            Assert.Equal("1", cs.GetParameter("remote-slot"));
        }

        [Fact]
        public void Percent_escapes_are_resolved()
        {
            var cs = ConnectionString.Parse("s7://host?path=%2Fsome%2Fpath&name=a%20b");

            Assert.Equal("/some/path", cs.GetParameter("path"));
            Assert.Equal("a b", cs.GetParameter("name"));
        }

        [Fact]
        public void A_plus_sign_decodes_to_a_space()
        {
            // Java hands parameter values to URLDecoder.decode(value, UTF_8), the
            // form-urlencoded decoder, so '+' means a space there. Uri.UnescapeDataString
            // would leave it as a literal '+' and the same connection string would
            // address a different device from each language.
            var cs = ConnectionString.Parse("s7://host?name=a+b");

            Assert.Equal("a b", cs.GetParameter("name"));
        }

        [Fact]
        public void Parameters_are_case_insensitive()
        {
            var cs = ConnectionString.Parse("x://h?Key=val");

            Assert.Equal("val", cs.GetParameter("KEY"));
            Assert.Equal("val", cs.GetParameter("key"));
        }

        [Fact]
        public void Bare_flags_are_read_as_true()
        {
            var cs = ConnectionString.Parse("s7://h?verbose");

            Assert.Equal("true", cs.GetParameter("verbose"));
        }

        [Fact]
        public void GetIntParameter_parses()
        {
            var cs = ConnectionString.Parse("s7://h?timeout=5000");

            Assert.Equal(5000, cs.GetIntParameter("timeout", 1000));
        }

        [Fact]
        public void GetIntParameter_does_not_depend_on_the_current_culture()
        {
            // A locale whose number formatting differs must not change how a connection
            // string is read. de-DE groups with '.', so a culture-sensitive parse would
            // read "5.000" as 5000 and could misread plain digits elsewhere.
            var previous = CultureInfo.CurrentCulture;
            try
            {
                CultureInfo.CurrentCulture = CultureInfo.GetCultureInfo("de-DE");
                var cs = ConnectionString.Parse("s7://h?timeout=5000");

                Assert.Equal(5000, cs.GetIntParameter("timeout", 1000));
                Assert.Equal(1000, ConnectionString.Parse("s7://h?timeout=5.000")
                    .GetIntParameter("timeout", 1000));
            }
            finally
            {
                CultureInfo.CurrentCulture = previous;
            }
        }

        [Fact]
        public void GetIntParameter_falls_back_when_absent()
        {
            var cs = ConnectionString.Parse("s7://h");

            Assert.Equal(1000, cs.GetIntParameter("timeout", 1000));
        }

        [Fact]
        public void GetBoolParameter_works()
        {
            var cs = ConnectionString.Parse("s7://h?enabled=true&flag=false");

            Assert.True(cs.GetBoolParameter("enabled", false));
            Assert.False(cs.GetBoolParameter("flag", true));
            Assert.True(cs.GetBoolParameter("missing", true));
        }

        [Fact]
        public void RedactSecrets_masks_password_tokens_and_secrets()
        {
            var input = "s7://192.168.0.1?remote-rack=0&password=hunter2&token=abc123&secret=shh";

            var redacted = ConnectionString.RedactSecrets(input);

            Assert.Contains("password=***", redacted);
            Assert.Contains("token=***", redacted);
            Assert.Contains("secret=***", redacted);
            // Non-secret parameters pass through.
            Assert.Contains("remote-rack=0", redacted);
        }

        [Theory]
        [InlineData("")]
        [InlineData("   ")]
        [InlineData(null)]
        public void Empty_connection_strings_throw(string? input)
        {
            Assert.Throws<PlcConnectionException>(() => ConnectionString.Parse(input));
        }

        [Fact]
        public void Missing_protocol_code_throws()
        {
            Assert.Throws<PlcConnectionException>(
                () => ConnectionString.Parse("://host"));
        }

        [Fact]
        public void ToString_redacts_secrets()
        {
            var cs = ConnectionString.Parse("s7://192.168.0.1?password=secret");

            var s = cs.ToString();

            Assert.DoesNotContain("secret", s);
            Assert.Contains("password=***", s);
        }
    }
}
