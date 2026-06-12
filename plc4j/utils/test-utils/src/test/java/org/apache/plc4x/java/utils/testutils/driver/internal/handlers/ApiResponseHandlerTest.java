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

import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.utils.testutils.driver.internal.TestContext;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.dom4j.DocumentHelper;
import org.dom4j.Element;
import org.junit.jupiter.api.Test;

import java.util.List;
import java.util.concurrent.CompletableFuture;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class ApiResponseHandlerTest {

    @Test
    void testExecuteApiResponse_noPendingResponse() {
        Element responseXml = DocumentHelper.createElement("PlcReadResponse");

        ApiResponseHandler handler = new ApiResponseHandler(responseXml);
        TestContext context = new TestContext();

        // Should throw because there's no pending response
        DriverTestsuiteException exception = assertThrows(DriverTestsuiteException.class,
            () -> handler.executeApiResponse(context));
        assertEquals("No pending response found in context", exception.getMessage());
    }

    @Test
    void testConstructor() {
        Element responseXml = DocumentHelper.createElement("PlcReadResponse");
        ApiResponseHandler handler = new ApiResponseHandler(responseXml);
        assertNotNull(handler);
    }

    @Test
    void testExecuteApiResponse_readResponse_withValidation() {
        // Create PlcReadResponse XML with expected values
        // Note: The actual XML generated includes type attributes, so we match that format
        Element responseXml = DocumentHelper.createElement("PlcReadResponse");
        Element valuesElement = responseXml.addElement("values").addAttribute("isList", "true");
        Element tagElement = valuesElement.addElement("testTag");
        Element responseItem = tagElement.addElement("PlcResponseItem");
        Element codeElement = responseItem.addElement("code");
        codeElement.addAttribute("dataType", "string");
        codeElement.addAttribute("bitLength", "16");
        codeElement.addAttribute("encoding", "UTF-8");
        codeElement.setText("OK");
        Element valueElement = responseItem.addElement("value");
        Element plcIntElement = valueElement.addElement("PlcINT");
        plcIntElement.addAttribute("dataType", "int");
        plcIntElement.addAttribute("bitLength", "16");
        plcIntElement.setText("42");

        // Mock the read response
        PlcReadResponse mockResponse = mock(PlcReadResponse.class);
        PlcValue mockValue = mock(PlcValue.class);

        when(mockResponse.getTagNames()).thenReturn(List.of("testTag"));
        when(mockResponse.getResponseCode("testTag")).thenReturn(PlcResponseCode.OK);
        when(mockResponse.getPlcValue("testTag")).thenReturn(mockValue);
        when(mockValue.isList()).thenReturn(false);
        when(mockValue.isStruct()).thenReturn(false);
        when(mockValue.getPlcValueType()).thenReturn(PlcValueType.INT);
        when(mockValue.toString()).thenReturn("42");

        CompletableFuture<PlcReadResponse> future = CompletableFuture.completedFuture(mockResponse);

        ApiResponseHandler handler = new ApiResponseHandler(responseXml);
        TestContext context = new TestContext();
        context.setPendingResponse(future);

        // Should not throw - values match
        assertDoesNotThrow(() -> handler.executeApiResponse(context));
    }

    @Test
    void testExecuteApiResponse_readResponse_withMismatchedValue() {
        // Create PlcReadResponse XML with expected values
        // Note: The actual XML generated includes type attributes, so we match that format
        Element responseXml = DocumentHelper.createElement("PlcReadResponse");
        Element valuesElement = responseXml.addElement("values").addAttribute("isList", "true");
        Element tagElement = valuesElement.addElement("testTag");
        Element responseItem = tagElement.addElement("PlcResponseItem");
        Element codeElement = responseItem.addElement("code");
        codeElement.addAttribute("dataType", "string");
        codeElement.addAttribute("bitLength", "16");
        codeElement.addAttribute("encoding", "UTF-8");
        codeElement.setText("OK");
        Element valueElement = responseItem.addElement("value");
        Element plcIntElement = valueElement.addElement("PlcINT");
        plcIntElement.addAttribute("dataType", "int");
        plcIntElement.addAttribute("bitLength", "16");
        plcIntElement.setText("100"); // Expected 100

        // Mock the read response returning 42 (mismatch)
        PlcReadResponse mockResponse = mock(PlcReadResponse.class);
        PlcValue mockValue = mock(PlcValue.class);

        when(mockResponse.getTagNames()).thenReturn(List.of("testTag"));
        when(mockResponse.getResponseCode("testTag")).thenReturn(PlcResponseCode.OK);
        when(mockResponse.getPlcValue("testTag")).thenReturn(mockValue);
        when(mockValue.isList()).thenReturn(false);
        when(mockValue.isStruct()).thenReturn(false);
        when(mockValue.getPlcValueType()).thenReturn(PlcValueType.INT);
        when(mockValue.toString()).thenReturn("42"); // Actual is 42

        CompletableFuture<PlcReadResponse> future = CompletableFuture.completedFuture(mockResponse);

        ApiResponseHandler handler = new ApiResponseHandler(responseXml);
        TestContext context = new TestContext();
        context.setPendingResponse(future);

        // Should throw - values don't match
        assertThrows(AssertionError.class, () -> handler.executeApiResponse(context));
    }

    @Test
    void testExecuteApiResponse_writeResponse_withValidation() {
        // Create PlcWriteResponse XML with expected response codes
        Element responseXml = DocumentHelper.createElement("PlcWriteResponse");
        Element responseCodesElement = responseXml.addElement("responseCodes");
        Element tagElement = responseCodesElement.addElement("testTag");
        tagElement.addElement("ResponseCode").addAttribute("stringRepresentation", "OK");

        // Mock the write response
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);
        when(mockResponse.getTagNames()).thenReturn(List.of("testTag"));
        when(mockResponse.getResponseCode("testTag")).thenReturn(PlcResponseCode.OK);

        CompletableFuture<PlcWriteResponse> future = CompletableFuture.completedFuture(mockResponse);

        ApiResponseHandler handler = new ApiResponseHandler(responseXml);
        TestContext context = new TestContext();
        context.setPendingResponse(future);

        // Should not throw - response codes match
        assertDoesNotThrow(() -> handler.executeApiResponse(context));
    }

    @Test
    void testExecuteApiResponse_writeResponse_withMismatchedCode() {
        // Create PlcWriteResponse XML expecting OK
        Element responseXml = DocumentHelper.createElement("PlcWriteResponse");
        Element responseCodesElement = responseXml.addElement("responseCodes");
        Element tagElement = responseCodesElement.addElement("testTag");
        tagElement.addElement("ResponseCode").addAttribute("stringRepresentation", "OK");

        // Mock the write response returning NOT_FOUND
        PlcWriteResponse mockResponse = mock(PlcWriteResponse.class);
        when(mockResponse.getTagNames()).thenReturn(List.of("testTag"));
        when(mockResponse.getResponseCode("testTag")).thenReturn(PlcResponseCode.NOT_FOUND);

        CompletableFuture<PlcWriteResponse> future = CompletableFuture.completedFuture(mockResponse);

        ApiResponseHandler handler = new ApiResponseHandler(responseXml);
        TestContext context = new TestContext();
        context.setPendingResponse(future);

        // Should throw - response codes don't match
        assertThrows(AssertionError.class, () -> handler.executeApiResponse(context));
    }
}
