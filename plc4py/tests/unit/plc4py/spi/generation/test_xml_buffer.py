#  Licensed to the Apache Software Foundation (ASF) under one
#  or more contributor license agreements.  See the NOTICE file
#  distributed with this work for additional information
#  regarding copyright ownership.  The ASF licenses this file
#  to you under the Apache License, Version 2.0 (the
#  "License"); you may not use this file except in compliance
#  with the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing,
#  software distributed under the License is distributed on an
#  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#  KIND, either express or implied.  See the License for the
#  specific language governing permissions and limitations
#  under the License.
from spi.generation.WriteBuffer import WriteBufferXmlBased


def test_xml_write_bit(mocker) -> None:

    wb: WriteBufferXmlBased = WriteBufferXmlBased()
    wb.write_bit(True, "error_flag")
    string: str = wb.to_xml_string()
    assert string == b'<test:testsuite xmlns:test="https://plc4x.apache.org/schemas/parser-serializer-testsuite.xsd" byteOrder="BIG_ENDIAN">\n\t<errorFlag dataType="bit" bitlength="1">true</errorFlag>\n</test:testsuite>'

