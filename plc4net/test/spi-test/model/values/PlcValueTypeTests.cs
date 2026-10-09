//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

using System;
using System.Linq;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.spi.model.values;
using org.apache.plc4net.types;
using Xunit;

namespace org.apache.plc4net.spi.test.model.values
{
    /// <summary>
    /// Locks <see cref="PlcValueType.DefaultDotNetType"/> to the types the value
    /// classes actually use. The table was first copied from Java's generated
    /// defaultJavaType enum, which carries Java-specific choices: unsigned PLC types
    /// widen to the next larger signed type because Java has no unsigned numerics,
    /// and dates and times of day map to java.time types that .NET splits into
    /// DateOnly and TimeOnly. The .NET value classes use native unsigned storage,
    /// DateOnly and TimeOnly, so the defaults have to follow them.
    /// </summary>
    public class PlcValueTypeTests
    {
        [Theory]
        [InlineData("BOOL")]
        [InlineData("BYTE")]
        [InlineData("WORD")]
        [InlineData("DWORD")]
        [InlineData("LWORD")]
        [InlineData("USINT")]
        [InlineData("UINT")]
        [InlineData("UDINT")]
        [InlineData("SINT")]
        [InlineData("INT")]
        [InlineData("DINT")]
        [InlineData("LINT")]
        [InlineData("REAL")]
        [InlineData("LREAL")]
        [InlineData("CHAR")]
        [InlineData("WCHAR")]
        [InlineData("STRING")]
        [InlineData("WSTRING")]
        [InlineData("TIME")]
        [InlineData("DATE")]
        [InlineData("TIME_OF_DAY")]
        [InlineData("DATE_AND_TIME")]
        public void Default_type_is_the_type_the_value_class_is_built_from(string name)
        {
            var valueType = (PlcValueType)typeof(PlcValueType).GetField(name)!.GetValue(null)!;
            var valueClass = typeof(PlcBOOL).Assembly.GetType("org.apache.plc4net.spi.model.values.Plc" + name);
            Assert.NotNull(valueClass);

            var constructorArgumentTypes = valueClass!.GetConstructors()
                .Select(constructor => constructor.GetParameters())
                .Where(parameters => parameters.Length == 1)
                .Select(parameters => parameters[0].ParameterType);

            Assert.Contains(valueType.DefaultDotNetType!, constructorArgumentTypes);
        }

        [Fact]
        public void Nanosecond_and_long_forms_follow_the_IPlcValue_accessors()
        {
            // These value classes are built from a raw nanosecond count (or have no
            // dedicated class), so the default is the type the matching IPlcValue
            // accessor hands back.
            Assert.Equal(new PlcLTIME(1UL).GetDuration().GetType(), PlcValueType.LTIME.DefaultDotNetType);
            Assert.Equal(new PlcLTIME_OF_DAY(0UL).GetTime().GetType(), PlcValueType.LTIME_OF_DAY.DefaultDotNetType);
            Assert.Equal(
                new PlcDATE_AND_LTIME(new DateTime(2024, 1, 2), 0u).GetDateTime().GetType(),
                PlcValueType.DATE_AND_LTIME.DefaultDotNetType);

            // LDATE and LDATE_AND_TIME have no dedicated value class; they are the
            // long forms of DATE and DATE_AND_TIME.
            Assert.Equal(PlcValueType.DATE.DefaultDotNetType, PlcValueType.LDATE.DefaultDotNetType);
            Assert.Equal(PlcValueType.DATE_AND_TIME.DefaultDotNetType, PlcValueType.LDATE_AND_TIME.DefaultDotNetType);
        }

        [Fact]
        public void Incompatible_accessors_throw_instead_of_returning_null()
        {
            // Mirrors the Java SPI3 PlcValueAdapter, which throws
            // PlcIncompatibleDatatypeException from every accessor the concrete
            // value cannot serve. Check the IsXxx() guards first.
            // (GetString and GetBoolArray are deliberately absent: PlcBOOL
            // legitimately serves both, same as Java's PlcBOOLEAN.)
            var value = new PlcBOOL(true);

            Assert.Throws<PlcIncompatibleDatatypeException>(() => value.GetRaw());
            Assert.Throws<PlcIncompatibleDatatypeException>(() => value.GetIndex(0));
            Assert.Throws<PlcIncompatibleDatatypeException>(() => value.GetList());
            Assert.Throws<PlcIncompatibleDatatypeException>(() => value.GetKeys());
            Assert.Throws<PlcIncompatibleDatatypeException>(() => value.GetValue("any"));
            Assert.Throws<PlcIncompatibleDatatypeException>(() => value.GetStruct());
        }
    }
}