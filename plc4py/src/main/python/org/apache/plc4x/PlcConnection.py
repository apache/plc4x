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

from generated.org.apache.plc4x.interop.ttypes import PlcException, FieldResponse


class PlcConnection:

    def __init__(self, client, url):
        self.client = client
        # Try to get a handle
        try:
            self.handle = self.client.connect(url)
        except:
            raise PlcException("Unable to connect to the given url " + url)

    def execute(self, request):
        response = self.client.execute(self.handle, request)
        return PlcResponse(response)

    def close(self):
        self.client.close(self.handle)


class PlcResponse:
    fields: dict

    def __init__(self, response):
        self.fields = response.fields

    def get_fields(self):
        return list(self.fields.keys())

    def get_field(self, field_name):
        if field_name in self.fields:
            return PlcResponseItem(self.fields.get(field_name))
        else:
            raise Exception("A field with name '" + field_name + "' is not contained in the response");


class PlcResponseItem:

    fieldResponse: FieldResponse

    def __init__(self, fieldResponse):
        self.fieldResponse = fieldResponse

    def get_response_code(self):
        return self.fieldResponse.responseCode

    def get_int_value(self):
        return self.fieldResponse.longValue

    def get_bool_value(self):
        return self.fieldResponse.boolValue

    def get_double_value(self):
        return self.fieldResponse.doubleValue

    def get_string_value(self):
        return self.fieldResponse.stringValue