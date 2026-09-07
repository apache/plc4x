/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.utils.testutils.driver.internal.handlers;

import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.utils.testutils.driver.internal.TestContext;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.dom4j.DocumentHelper;
import org.dom4j.Element;
import org.dom4j.QName;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.mockito.ArgumentCaptor;

import java.util.LinkedHashSet;
import java.util.List;
import java.util.concurrent.CompletableFuture;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

class ApiRequestHandlerTest {

    private PlcConnection mockConnection;

    @BeforeEach
    void setUp() {
        mockConnection = mock(PlcConnection.class);
    }

    @Test
    void testExecuteReadRequest() throws Exception {
        // Create request XML with child elements for name and address
        Element requestXml = DocumentHelper.createElement("PlcReadRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag1 = tagsElement.addElement("tag");
        tag1.addElement("name").setText("tag1");
        tag1.addElement("address").setText("address1");

        // Mock the read request builder chain
        PlcReadRequest.Builder mockBuilder = mock(PlcReadRequest.Builder.class);
        PlcReadRequest mockRequest = mock(PlcReadRequest.class);
        PlcReadResponse mockResponse = mock(PlcReadResponse.class);

        when(mockConnection.readRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addTagAddress(anyString(), anyString())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockRequest.getTagNames()).thenReturn(new LinkedHashSet<>(List.of("tag1")));
        when(mockRequest.getTagResponseCode("tag1")).thenReturn(PlcResponseCode.OK);
        when(mockRequest.getTag("tag1")).thenReturn(mock(org.apache.plc4x.java.api.model.PlcTag.class));

        // Execute
        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        assertDoesNotThrow(() -> handler.executeApiRequest(mockConnection, new TestContext()));

        // Verify
        verify(mockBuilder).addTagAddress("tag1", "address1");
        verify(mockRequest).execute();
    }

    @Test
    void testExecuteWriteRequest() throws Exception {
        // Create request XML
        Element requestXml = DocumentHelper.createElement("PlcWriteRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag1 = tagsElement.addElement("tag");
        tag1.addAttribute("name", "tag1");
        tag1.addAttribute("address", "address1");
        tag1.setText("42");

        // Mock the write request builder chain
        PlcWriteRequest.Builder mockBuilder = mock(PlcWriteRequest.Builder.class);
        PlcWriteRequest mockRequest = mock(PlcWriteRequest.class);
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);

        when(mockConnection.writeRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addTagAddress(anyString(), anyString(), any())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockResponse.getTagNames()).thenReturn(List.of("tag1"));
        when(mockResponse.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);

        // Execute
        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        assertDoesNotThrow(() -> handler.executeApiRequest(mockConnection, new TestContext()));

        // Verify
        verify(mockBuilder).addTagAddress(eq("tag1"), eq("address1"), eq(new PlcDINT(42)));
        verify(mockRequest).execute();
    }

    @Test
    @Disabled("needs fixing")
    void testExecuteBrowseRequest() throws Exception {
        // Create request XML
        Element requestXml = DocumentHelper.createElement("PlcBrowseRequest");
        Element queryElement = requestXml.addElement("query");
        queryElement.setText("*");

        // Mock the browse request builder chain
        PlcBrowseRequest.Builder mockBuilder = mock(PlcBrowseRequest.Builder.class);
        PlcBrowseRequest mockRequest = mock(PlcBrowseRequest.class);
        PlcBrowseResponse mockResponse = mock(PlcBrowseResponse.class);

        when(mockConnection.browseRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addQuery(anyString(), anyString())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockResponse.getQueryNames()).thenReturn(new LinkedHashSet<>(List.of("query")));
        when(mockResponse.getValues("query")).thenReturn(List.of());

        // Execute
        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        assertDoesNotThrow(() -> handler.executeApiRequest(mockConnection, new TestContext()));

        // Verify
        verify(mockBuilder).addQuery("query", "**");
        verify(mockRequest).execute();
    }

    @Test
    void testExecuteSubscriptionRequest() throws Exception {
        // Create request XML
        Element requestXml = DocumentHelper.createElement("PlcSubscriptionRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag = tagsElement.addElement("tag");
        tag.addElement(new QName("name")).addText("tag1");
        tag.addElement(new QName("address")).addText("address1");

        // Mock the subscription request builder chain
        PlcSubscriptionRequest.Builder mockBuilder = mock(PlcSubscriptionRequest.Builder.class);
        PlcSubscriptionRequest mockRequest = mock(PlcSubscriptionRequest.class);
        PlcSubscriptionResponse mockResponse = mock(PlcSubscriptionResponse.class);

        when(mockConnection.subscriptionRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addChangeOfStateTagAddress(anyString(), anyString())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));

        // Execute
        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        assertDoesNotThrow(() -> handler.executeApiRequest(mockConnection, new TestContext()));

        // Verify
        verify(mockBuilder).addChangeOfStateTagAddress("tag1", "address1");
        verify(mockRequest).execute();
    }

    @Test
    void testUnknownRequestType() {
        Element requestXml = DocumentHelper.createElement("UnknownRequest");

        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        assertThrows(DriverTestsuiteException.class, () -> {
            handler.executeApiRequest(mockConnection, new TestContext());
        });
    }

    @Test
    void testParseValueInteger() throws Exception {
        Element requestXml = DocumentHelper.createElement("PlcWriteRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag1 = tagsElement.addElement("tag");
        tag1.addAttribute("name", "tag1");
        tag1.addAttribute("address", "address1");
        tag1.setText("123");

        PlcWriteRequest.Builder mockBuilder = mock(PlcWriteRequest.Builder.class);
        PlcWriteRequest mockRequest = mock(PlcWriteRequest.class);
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);

        when(mockConnection.writeRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addTagAddress(anyString(), anyString(), any())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockResponse.getTagNames()).thenReturn(List.of("tag1"));
        when(mockResponse.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);

        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        handler.executeApiRequest(mockConnection, new TestContext());

        verify(mockBuilder).addTagAddress(eq("tag1"), eq("address1"), eq(new PlcDINT(123)));
    }

    @Test
    void testParseValueDouble() throws Exception {
        Element requestXml = DocumentHelper.createElement("PlcWriteRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag1 = tagsElement.addElement("tag");
        tag1.addAttribute("name", "tag1");
        tag1.addAttribute("address", "address1");
        tag1.setText("3.14");

        PlcWriteRequest.Builder mockBuilder = mock(PlcWriteRequest.Builder.class);
        PlcWriteRequest mockRequest = mock(PlcWriteRequest.class);
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);

        when(mockConnection.writeRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addTagAddress(anyString(), anyString(), any())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockResponse.getTagNames()).thenReturn(List.of("tag1"));
        when(mockResponse.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);

        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        handler.executeApiRequest(mockConnection, new TestContext());

        verify(mockBuilder).addTagAddress(eq("tag1"), eq("address1"), eq(new PlcLREAL(3.14)));
    }

    @Test
    void testParseValueBoolean() throws Exception {
        Element requestXml = DocumentHelper.createElement("PlcWriteRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag1 = tagsElement.addElement("tag");
        tag1.addAttribute("name", "tag1");
        tag1.addAttribute("address", "address1");
        tag1.setText("true");

        PlcWriteRequest.Builder mockBuilder = mock(PlcWriteRequest.Builder.class);
        PlcWriteRequest mockRequest = mock(PlcWriteRequest.class);
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);

        when(mockConnection.writeRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addTagAddress(anyString(), anyString(), any())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockResponse.getTagNames()).thenReturn(List.of("tag1"));
        when(mockResponse.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);

        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        handler.executeApiRequest(mockConnection, new TestContext());

        verify(mockBuilder).addTagAddress(eq("tag1"), eq("address1"), eq(new PlcBOOL(true)));
    }

    @Test
    void testParseValueString() throws Exception {
        Element requestXml = DocumentHelper.createElement("PlcWriteRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag1 = tagsElement.addElement("tag");
        tag1.addAttribute("name", "tag1");
        tag1.addAttribute("address", "address1");
        tag1.setText("hello");

        PlcWriteRequest.Builder mockBuilder = mock(PlcWriteRequest.Builder.class);
        PlcWriteRequest mockRequest = mock(PlcWriteRequest.class);
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);

        when(mockConnection.writeRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addTagAddress(anyString(), anyString(), any())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockResponse.getTagNames()).thenReturn(List.of("tag1"));
        when(mockResponse.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);

        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        handler.executeApiRequest(mockConnection, new TestContext());

        verify(mockBuilder).addTagAddress(eq("tag1"), eq("address1"), eq(new PlcSTRING("hello")));
    }

    @Test
    void testParseTypedValuePlcUINT() throws Exception {
        // Create request XML with nested PLC type element
        Element requestXml = DocumentHelper.createElement("PlcWriteRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag1 = tagsElement.addElement("tag");
        tag1.addElement("name").setText("tag1");
        tag1.addElement("address").setText("address1");
        Element valueElement = tag1.addElement("value");
        valueElement.addElement("PlcUINT").setText("1");

        PlcWriteRequest.Builder mockBuilder = mock(PlcWriteRequest.Builder.class);
        PlcWriteRequest mockRequest = mock(PlcWriteRequest.class);
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);

        when(mockConnection.writeRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addTagAddress(anyString(), anyString(), any())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockResponse.getTagNames()).thenReturn(List.of("tag1"));
        when(mockResponse.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);

        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        handler.executeApiRequest(mockConnection, new TestContext());

        // PlcUINT should be parsed as Integer
        verify(mockBuilder).addTagAddress(eq("tag1"), eq("address1"), eq(new PlcUINT(1)));
    }

    @Test
    void testParseTypedValuePlcINT() throws Exception {
        // Create request XML with nested PLC type element
        Element requestXml = DocumentHelper.createElement("PlcWriteRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag1 = tagsElement.addElement("tag");
        tag1.addElement("name").setText("tag1");
        tag1.addElement("address").setText("address1");
        Element valueElement = tag1.addElement("value");
        valueElement.addElement("PlcINT").setText("42");

        PlcWriteRequest.Builder mockBuilder = mock(PlcWriteRequest.Builder.class);
        PlcWriteRequest mockRequest = mock(PlcWriteRequest.class);
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);

        when(mockConnection.writeRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addTagAddress(anyString(), anyString(), any())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockResponse.getTagNames()).thenReturn(List.of("tag1"));
        when(mockResponse.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);

        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        handler.executeApiRequest(mockConnection, new TestContext());

        // PlcINT should be parsed as Short
        verify(mockBuilder).addTagAddress(eq("tag1"), eq("address1"), eq(new PlcINT((short) 42)));
    }

    @Test
    void testParseTypedValuePlcREAL() throws Exception {
        // Create request XML with nested PLC type element
        Element requestXml = DocumentHelper.createElement("PlcWriteRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag1 = tagsElement.addElement("tag");
        tag1.addElement("name").setText("tag1");
        tag1.addElement("address").setText("address1");
        Element valueElement = tag1.addElement("value");
        valueElement.addElement("PlcREAL").setText("3.14");

        PlcWriteRequest.Builder mockBuilder = mock(PlcWriteRequest.Builder.class);
        PlcWriteRequest mockRequest = mock(PlcWriteRequest.class);
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);

        when(mockConnection.writeRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addTagAddress(anyString(), anyString(), any())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockResponse.getTagNames()).thenReturn(List.of("tag1"));
        when(mockResponse.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);

        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        handler.executeApiRequest(mockConnection, new TestContext());

        // PlcREAL should be parsed as Float
        verify(mockBuilder).addTagAddress(eq("tag1"), eq("address1"), eq(new PlcREAL(3.14f)));
    }

    @Test
    void testParseTypedValuePlcBOOL() throws Exception {
        // Create request XML with nested PLC type element
        Element requestXml = DocumentHelper.createElement("PlcWriteRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag1 = tagsElement.addElement("tag");
        tag1.addElement("name").setText("tag1");
        tag1.addElement("address").setText("address1");
        Element valueElement = tag1.addElement("value");
        valueElement.addElement("PlcBOOL").setText("true");

        PlcWriteRequest.Builder mockBuilder = mock(PlcWriteRequest.Builder.class);
        PlcWriteRequest mockRequest = mock(PlcWriteRequest.class);
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);

        when(mockConnection.writeRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addTagAddress(anyString(), anyString(), any())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockResponse.getTagNames()).thenReturn(List.of("tag1"));
        when(mockResponse.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);

        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        handler.executeApiRequest(mockConnection, new TestContext());

        // PlcBOOL should be parsed as Boolean
        verify(mockBuilder).addTagAddress(eq("tag1"), eq("address1"), eq(new PlcBOOL(true)));
    }

    @Test
    void testParseTypedValuePlcList() throws Exception {
        // Create request XML with nested PLC type element for a list
        Element requestXml = DocumentHelper.createElement("PlcWriteRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag1 = tagsElement.addElement("tag");
        tag1.addElement("name").setText("tag1");
        tag1.addElement("address").setText("address1");
        Element valueElement = tag1.addElement("value");
        Element listElement = valueElement.addElement("PlcList");
        listElement.addElement("PlcUINT").setText("1");
        listElement.addElement("PlcUINT").setText("2");
        listElement.addElement("PlcUINT").setText("3");

        PlcWriteRequest.Builder mockBuilder = mock(PlcWriteRequest.Builder.class);
        PlcWriteRequest mockRequest = mock(PlcWriteRequest.class);
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);

        when(mockConnection.writeRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addTagAddress(anyString(), anyString(), any())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockResponse.getTagNames()).thenReturn(List.of("tag1"));
        when(mockResponse.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);

        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        handler.executeApiRequest(mockConnection, new TestContext());

        // PlcList should be parsed as Object array
        ArgumentCaptor<Object> valueCaptor = ArgumentCaptor.forClass(Object.class);
        verify(mockBuilder).addTagAddress(eq("tag1"), eq("address1"), valueCaptor.capture());

        Object capturedValue = valueCaptor.getValue();
        assertInstanceOf(PlcList.class, capturedValue);
        PlcList plcList = (PlcList) capturedValue;
        assertEquals(3, plcList.getLength());
        assertEquals(new PlcUINT(1), plcList.getIndex(0));
        assertEquals(new PlcUINT(2), plcList.getIndex(1));
        assertEquals(new PlcUINT(3), plcList.getIndex(2));
    }

    @Test
    void testParseTypedValuePlcDINT() throws Exception {
        // Create request XML with nested PLC type element
        Element requestXml = DocumentHelper.createElement("PlcWriteRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag1 = tagsElement.addElement("tag");
        tag1.addElement("name").setText("tag1");
        tag1.addElement("address").setText("address1");
        Element valueElement = tag1.addElement("value");
        valueElement.addElement("PlcDINT").setText("123456");

        PlcWriteRequest.Builder mockBuilder = mock(PlcWriteRequest.Builder.class);
        PlcWriteRequest mockRequest = mock(PlcWriteRequest.class);
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);

        when(mockConnection.writeRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addTagAddress(anyString(), anyString(), any())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockResponse.getTagNames()).thenReturn(List.of("tag1"));
        when(mockResponse.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);

        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        handler.executeApiRequest(mockConnection, new TestContext());

        // PlcDINT should be parsed as Integer
        verify(mockBuilder).addTagAddress(eq("tag1"), eq("address1"), eq(new PlcDINT(123456)));
    }

    @Test
    void testParseTypedValuePlcSTRING() throws Exception {
        // Create request XML with nested PLC type element
        Element requestXml = DocumentHelper.createElement("PlcWriteRequest");
        Element tagsElement = requestXml.addElement("tags");
        Element tag1 = tagsElement.addElement("tag");
        tag1.addElement("name").setText("tag1");
        tag1.addElement("address").setText("address1");
        Element valueElement = tag1.addElement("value");
        valueElement.addElement("PlcSTRING").setText("hello world");

        PlcWriteRequest.Builder mockBuilder = mock(PlcWriteRequest.Builder.class);
        PlcWriteRequest mockRequest = mock(PlcWriteRequest.class);
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);

        when(mockConnection.writeRequestBuilder()).thenReturn(mockBuilder);
        when(mockBuilder.addTagAddress(anyString(), anyString(), any())).thenReturn(mockBuilder);
        when(mockBuilder.build()).thenReturn(mockRequest);
        when(mockRequest.execute()).thenAnswer(invocation -> CompletableFuture.completedFuture(mockResponse));
        when(mockResponse.getTagNames()).thenReturn(List.of("tag1"));
        when(mockResponse.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);

        ApiRequestHandler handler = new ApiRequestHandler(requestXml);
        handler.executeApiRequest(mockConnection, new TestContext());

        // PlcSTRING should be parsed as String
        verify(mockBuilder).addTagAddress(eq("tag1"), eq("address1"), eq(new PlcSTRING("hello world")));
    }
}
