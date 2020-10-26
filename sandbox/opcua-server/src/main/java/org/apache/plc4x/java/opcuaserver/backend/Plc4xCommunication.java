/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.opcuaserver.backend;


import java.lang.reflect.Array;
import java.util.List;
import java.util.Arrays;
import java.util.Random;
import java.util.UUID;

import org.eclipse.milo.opcua.sdk.core.AccessLevel;
import org.eclipse.milo.opcua.sdk.core.Reference;
import org.eclipse.milo.opcua.sdk.core.ValueRank;
import org.eclipse.milo.opcua.sdk.core.ValueRanks;
import org.eclipse.milo.opcua.sdk.server.Lifecycle;
import org.eclipse.milo.opcua.sdk.server.OpcUaServer;
import org.eclipse.milo.opcua.sdk.server.api.DataItem;
import org.eclipse.milo.opcua.sdk.server.api.DataTypeDictionaryManager;
import org.eclipse.milo.opcua.sdk.server.api.ManagedNamespaceWithLifecycle;
import org.eclipse.milo.opcua.sdk.server.api.MonitoredItem;
import org.eclipse.milo.opcua.sdk.server.model.nodes.objects.BaseEventTypeNode;
import org.eclipse.milo.opcua.sdk.server.model.nodes.objects.ServerTypeNode;
import org.eclipse.milo.opcua.sdk.server.model.nodes.variables.AnalogItemTypeNode;
import org.eclipse.milo.opcua.sdk.server.nodes.UaFolderNode;
import org.eclipse.milo.opcua.sdk.server.nodes.UaMethodNode;
import org.eclipse.milo.opcua.sdk.server.nodes.UaNode;
import org.eclipse.milo.opcua.sdk.server.nodes.UaObjectNode;
import org.eclipse.milo.opcua.sdk.server.nodes.UaObjectTypeNode;
import org.eclipse.milo.opcua.sdk.server.nodes.UaVariableNode;
import org.eclipse.milo.opcua.sdk.server.nodes.factories.NodeFactory;
import org.eclipse.milo.opcua.sdk.server.nodes.filters.AttributeFilters;
import org.eclipse.milo.opcua.sdk.server.util.SubscriptionModel;
import org.eclipse.milo.opcua.stack.core.AttributeId;
import org.eclipse.milo.opcua.stack.core.BuiltinDataType;
import org.eclipse.milo.opcua.stack.core.Identifiers;
import org.eclipse.milo.opcua.stack.core.UaException;
import org.eclipse.milo.opcua.stack.core.types.builtin.ByteString;
import org.eclipse.milo.opcua.stack.core.types.builtin.DataValue;
import org.eclipse.milo.opcua.stack.core.types.builtin.DateTime;
import org.eclipse.milo.opcua.stack.core.types.builtin.ExtensionObject;
import org.eclipse.milo.opcua.stack.core.types.builtin.LocalizedText;
import org.eclipse.milo.opcua.stack.core.types.builtin.NodeId;
import org.eclipse.milo.opcua.stack.core.types.builtin.QualifiedName;
import org.eclipse.milo.opcua.stack.core.types.builtin.Variant;
import org.eclipse.milo.opcua.stack.core.types.builtin.XmlElement;
import org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.UInteger;
import org.eclipse.milo.opcua.stack.core.types.enumerated.StructureType;
import org.eclipse.milo.opcua.stack.core.types.structured.EnumDefinition;
import org.eclipse.milo.opcua.stack.core.types.structured.EnumDescription;
import org.eclipse.milo.opcua.stack.core.types.structured.EnumField;
import org.eclipse.milo.opcua.stack.core.types.structured.Range;
import org.eclipse.milo.opcua.stack.core.types.structured.StructureDefinition;
import org.eclipse.milo.opcua.stack.core.types.structured.StructureDescription;
import org.eclipse.milo.opcua.stack.core.types.structured.StructureField;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.utils.connectionpool.*;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;

import static org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.Unsigned.ubyte;
import static org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.Unsigned.uint;
import static org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.Unsigned.ulong;
import static org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.Unsigned.ushort;

public class Plc4xCommunication {

    private PlcDriverManager driverManager;
    private final Logger logger = LoggerFactory.getLogger(getClass());

    public Plc4xCommunication () {
        driverManager = new PooledPlcDriverManager();
    }

    public PlcField getField(String tag, String connectionString) throws PlcConnectionException {
        return driverManager.getDriver(connectionString).prepareField(tag);
    }

    public static NodeId getNodeId(String plcValue) {
        switch (plcValue) {
            case "BOOL":
            case "BIT":
                return Identifiers.Boolean;
            case "BYTE":
            case "BITARR8":
                return Identifiers.Byte;
            case "SINT":
            case "INT8":
                return Identifiers.SByte;
            case "USINT":
            case "UINT8":
            case "BIT8":
                return Identifiers.Byte;
            case "INT":
            case "INT16":
                return Identifiers.Int16;
            case "UINT":
            case "UINT16":
                return Identifiers.UInt16;
            case "WORD":
            case "BITARR16":
                return Identifiers.UInt16;
            case "DINT":
            case "INT32":
                return Identifiers.Int32;
            case "UDINT":
            case "UINT32":
                return Identifiers.UInt32;
            case "DWORD":
            case "BITARR32":
                return Identifiers.UInt32;
            case "LINT":
            case "INT64":
                return Identifiers.Int64;
            case "ULINT":
            case "UINT64":
                return Identifiers.UInt64;
            case "LWORD":
            case "BITARR64":
                return Identifiers.UInt64;
            case "REAL":
            case "FLOAT":
                return Identifiers.Float;
            case "LREAL":
            case "DOUBLE":
                return Identifiers.Double;
            case "CHAR":
                return Identifiers.String;
            case "WCHAR":
                return Identifiers.String;
            case "STRING":
                return Identifiers.String;
            case "WSTRING":
            case "STRING16":
                return Identifiers.String;
            default:
                return Identifiers.BaseDataType;
        }
    }

    public Variant getValue(String tag, String connectionString) {
        PlcConnection connection = null;
        try {
          connection = driverManager.getConnection(connectionString);
        } catch (PlcConnectionException e) {
          System.out.println("Failed" + e);
        }

        // Create a new read request:
        // - Give the single item requested an alias name
        PlcReadRequest.Builder builder = connection.readRequestBuilder();
        builder.addItem("value-1", tag);
        PlcReadRequest readRequest = builder.build();

        PlcReadResponse response = null;
        try {
          response = readRequest.execute().get();
        } catch (InterruptedException | ExecutionException e) {
          System.out.println("Failed" + e);
        }

        Variant resp = null;

        for (String fieldName : response.getFieldNames()) {
          if(response.getResponseCode(fieldName) == PlcResponseCode.OK) {
              int numValues = response.getNumberOfValues(fieldName);
              if(numValues == 1) {
                  resp = new Variant(response.getObject(fieldName));
              } else {
                Object array = Array.newInstance(response.getObject(fieldName, 0).getClass(), numValues);
                for (int i = 0; i < numValues; i++) {
                    Array.set(array, i, response.getObject(fieldName, i));
                }
                resp = new Variant(array);
              }
          }
        }
        try {
          connection.close();
        } catch (Exception e) {
          System.out.println("Close Failed" + e);
        }
        return resp;
    }

    public void setValue(String tag, String value, String connectionString) {
        PlcConnection connection = null;
        try {
          connection = driverManager.getConnection(connectionString);
        } catch (PlcConnectionException e) {
          System.out.println("Failed" + e);
        }

        // Create a new read request:
        // - Give the single item requested an alias name
        final PlcWriteRequest.Builder builder = connection.writeRequestBuilder();

        //If an array value is passed instead of a single value then convert to a String array
        if ((value.charAt(0) == '[') && (value.charAt(value.length() - 1) == ']')) {
            String[] values = value.substring(1,value.length() - 1).split(",");
            logger.info("Adding Tag " + Arrays.toString(values));
            builder.addItem(tag, tag, values);
        } else {
            builder.addItem(tag, tag, value);
        }

        PlcWriteRequest writeRequest = builder.build();

        try {
          writeRequest.execute().get();
        } catch (InterruptedException | ExecutionException e) {
          System.out.println("Failed" + e);
        }

        try {
          connection.close();
        } catch (Exception e) {
          System.out.println("Close Failed" + e);
        }
        return;
    }
}
