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

package org.eclipse.milo.examples.server;

import static org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.Unsigned.ubyte;
import static org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.Unsigned.uint;
import static org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.Unsigned.ulong;
import static org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.Unsigned.ushort;

import java.util.List;
import org.eclipse.milo.opcua.sdk.core.AccessLevel;
import org.eclipse.milo.opcua.sdk.core.Reference;
import org.eclipse.milo.opcua.sdk.core.ValueRank;
import org.eclipse.milo.opcua.sdk.server.OpcUaServer;
import org.eclipse.milo.opcua.sdk.server.api.DataItem;
import org.eclipse.milo.opcua.sdk.server.api.ManagedNamespaceWithLifecycle;
import org.eclipse.milo.opcua.sdk.server.api.MonitoredItem;
import org.eclipse.milo.opcua.sdk.server.util.SubscriptionModel;
import org.eclipse.milo.opcua.sdk.server.nodes.UaFolderNode;
import org.eclipse.milo.opcua.sdk.server.nodes.UaVariableNode;
import org.eclipse.milo.opcua.stack.core.Identifiers;
import org.eclipse.milo.opcua.stack.core.types.builtin.DataValue;
import org.eclipse.milo.opcua.stack.core.types.builtin.LocalizedText;
import org.eclipse.milo.opcua.stack.core.types.builtin.NodeId;
import org.eclipse.milo.opcua.stack.core.types.builtin.Variant;
import org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.UInteger;
import org.eclipse.milo.opcua.stack.core.types.builtin.unsigned.UShort;

/**
 * A second server namespace (index 3) dedicated to PLC4X integration testing. Where Milo's
 * {@link ExampleNamespace} (ns=2, {@code HelloWorld/...}) exposes one node per built-in type, this
 * namespace mirrors the structure of the ADS manual test
 * ({@code ManualFactoryAdsDriverTestTC3}) so the OPC UA driver can be exercised against the same
 * matrix of addressing variants: scalars, 1-D arrays and multi-dimensional matrices/cubes (and,
 * in a later increment, custom structs / nested structs / arrays of structs).
 *
 * <p>All values are fixed and distinctive so a client test can assert exact round-trips. Node
 * addresses are string identifiers of the form {@code ns=3;s=Test/Scalar/Bool}.
 */
public class Plc4xTestNamespace extends ManagedNamespaceWithLifecycle {

    public static final String NAMESPACE_URI = "urn:apache:plc4x:test";

    private final SubscriptionModel subscriptionModel;

    public Plc4xTestNamespace(OpcUaServer server) {
        super(server, NAMESPACE_URI);

        subscriptionModel = new SubscriptionModel(server, this);
        getLifecycleManager().addLifecycle(subscriptionModel);
        getLifecycleManager().addStartupTask(this::createAndAddNodes);
    }

    private void createAndAddNodes() {
        UaFolderNode root = new UaFolderNode(
            getNodeContext(),
            newNodeId("Test"),
            newQualifiedName("Test"),
            LocalizedText.english("Test"));
        getNodeManager().addNode(root);
        // Surface the folder under the standard Objects folder so it is browsable from the root.
        root.addReference(new Reference(
            root.getNodeId(),
            Identifiers.Organizes,
            Identifiers.ObjectsFolder.expanded(),
            Reference.Direction.INVERSE));

        addScalarNodes(root);
        addArrayNodes(root);
        addMatrixNodes(root);
    }

    // =====================================================================================
    // Scalars — one node per built-in type the driver maps to an IEC 61131-3 type.
    // =====================================================================================
    private void addScalarNodes(UaFolderNode root) {
        UaFolderNode folder = childFolder(root, "Scalar");
        addVariable(folder, "Scalar/Bool", Identifiers.Boolean, new Variant(true));
        addVariable(folder, "Scalar/SInt", Identifiers.SByte, new Variant((byte) -12));
        addVariable(folder, "Scalar/USInt", Identifiers.Byte, new Variant(ubyte(250)));
        addVariable(folder, "Scalar/Int", Identifiers.Int16, new Variant((short) -1234));
        addVariable(folder, "Scalar/UInt", Identifiers.UInt16, new Variant(ushort(54321)));
        addVariable(folder, "Scalar/DInt", Identifiers.Int32, new Variant(-12345678));
        addVariable(folder, "Scalar/UDInt", Identifiers.UInt32, new Variant(uint(305419896L)));
        addVariable(folder, "Scalar/LInt", Identifiers.Int64, new Variant(-9223372036854770000L));
        addVariable(folder, "Scalar/ULInt", Identifiers.UInt64, new Variant(ulong(new java.math.BigInteger("18446744073709551000"))));
        addVariable(folder, "Scalar/Real", Identifiers.Float, new Variant(3.14159f));
        addVariable(folder, "Scalar/LReal", Identifiers.Double, new Variant(2.718281828459045d));
        addVariable(folder, "Scalar/String", Identifiers.String, new Variant("Hello PLC4X"));
    }

    // =====================================================================================
    // 1-D arrays — whole-array reads (PlcList) and, via IndexRange, element/slice reads.
    // =====================================================================================
    private void addArrayNodes(UaFolderNode root) {
        UaFolderNode folder = childFolder(root, "Array");
        addArray(folder, "Array/Bool", Identifiers.Boolean, ValueRank.OneDimension,
            new Boolean[]{true, false, true, true, false, false, true, false});
        addArray(folder, "Array/Int", Identifiers.Int16, ValueRank.OneDimension,
            new Short[]{-3, -1, 0, 1, 3});
        addArray(folder, "Array/UInt", Identifiers.UInt16, ValueRank.OneDimension,
            new UShort[]{ushort(1), ushort(10), ushort(100), ushort(1000), ushort(10000)});
        addArray(folder, "Array/DInt", Identifiers.Int32, ValueRank.OneDimension,
            new Integer[]{-1000, -500, 0, 1000000, 2000000});
        addArray(folder, "Array/LReal", Identifiers.Double, ValueRank.OneDimension,
            new Double[]{1.5d, -2.0d, 0.125d});
        addArray(folder, "Array/String", Identifiers.String, ValueRank.OneDimension,
            new String[]{"alpha", "beta", "gamma"});
    }

    // =====================================================================================
    // Multi-dimensional arrays — matrix (2-D) and cube (3-D), whole-array + IndexRange slices.
    // =====================================================================================
    private void addMatrixNodes(UaFolderNode root) {
        UaFolderNode folder = childFolder(root, "Matrix");

        // INT[2][3]
        short[][] matI16 = {{10, 11, 12}, {-10, -11, -12}};
        addMatrix(folder, "Matrix/Int_2x3", Identifiers.Int16, matI16, new UInteger[]{uint(2), uint(3)});

        // REAL[3][2]
        float[][] matR32 = {{1.0f, 1.5f}, {2.0f, 2.5f}, {3.0f, 3.5f}};
        addMatrix(folder, "Matrix/Real_3x2", Identifiers.Float, matR32, new UInteger[]{uint(3), uint(2)});

        // UINT[2][2][2]
        UShort[][][] cubeU16 = {
            {{ushort(1), ushort(2)}, {ushort(3), ushort(4)}},
            {{ushort(5), ushort(6)}, {ushort(7), ushort(8)}}
        };
        addMatrix(folder, "Matrix/UInt_2x2x2", Identifiers.UInt16, cubeU16,
            new UInteger[]{uint(2), uint(2), uint(2)});
    }

    // ---------------------------------------------------------------------------------------
    // Builders
    // ---------------------------------------------------------------------------------------
    private UaFolderNode childFolder(UaFolderNode parent, String name) {
        UaFolderNode folder = new UaFolderNode(
            getNodeContext(),
            newNodeId("Test/" + name),
            newQualifiedName(name),
            LocalizedText.english(name));
        getNodeManager().addNode(folder);
        parent.addOrganizes(folder);
        return folder;
    }

    private void addVariable(UaFolderNode folder, String id, NodeId dataType, Variant value) {
        UaVariableNode node = new UaVariableNode.UaVariableNodeBuilder(getNodeContext())
            .setNodeId(newNodeId("Test/" + id))
            .setAccessLevel(AccessLevel.READ_WRITE)
            .setUserAccessLevel(AccessLevel.READ_WRITE)
            .setBrowseName(newQualifiedName(leaf(id)))
            .setDisplayName(LocalizedText.english(leaf(id)))
            .setDataType(dataType)
            .setTypeDefinition(Identifiers.BaseDataVariableType)
            .build();
        node.setValue(new DataValue(value));
        getNodeManager().addNode(node);
        folder.addOrganizes(node);
    }

    private void addArray(UaFolderNode folder, String id, NodeId dataType, ValueRank valueRank, Object array) {
        UaVariableNode node = new UaVariableNode.UaVariableNodeBuilder(getNodeContext())
            .setNodeId(newNodeId("Test/" + id))
            .setAccessLevel(AccessLevel.READ_WRITE)
            .setUserAccessLevel(AccessLevel.READ_WRITE)
            .setBrowseName(newQualifiedName(leaf(id)))
            .setDisplayName(LocalizedText.english(leaf(id)))
            .setDataType(dataType)
            .setTypeDefinition(Identifiers.BaseDataVariableType)
            .setValueRank(valueRank.getValue())
            .setArrayDimensions(new UInteger[]{uint(0)})
            .build();
        node.setValue(new DataValue(new Variant(array)));
        getNodeManager().addNode(node);
        folder.addOrganizes(node);
    }

    private void addMatrix(UaFolderNode folder, String id, NodeId dataType, Object array, UInteger[] dimensions) {
        UaVariableNode node = new UaVariableNode.UaVariableNodeBuilder(getNodeContext())
            .setNodeId(newNodeId("Test/" + id))
            .setAccessLevel(AccessLevel.READ_WRITE)
            .setUserAccessLevel(AccessLevel.READ_WRITE)
            .setBrowseName(newQualifiedName(leaf(id)))
            .setDisplayName(LocalizedText.english(leaf(id)))
            .setDataType(dataType)
            .setTypeDefinition(Identifiers.BaseDataVariableType)
            .setValueRank(dimensions.length)
            .setArrayDimensions(dimensions)
            .build();
        node.setValue(new DataValue(new Variant(array)));
        getNodeManager().addNode(node);
        folder.addOrganizes(node);
    }

    private static String leaf(String id) {
        int slash = id.lastIndexOf('/');
        return slash < 0 ? id : id.substring(slash + 1);
    }

    // ---------------------------------------------------------------------------------------
    // Subscription plumbing — delegate to the SubscriptionModel (required by the base class).
    // ---------------------------------------------------------------------------------------
    @Override
    public void onDataItemsCreated(List<DataItem> dataItems) {
        subscriptionModel.onDataItemsCreated(dataItems);
    }

    @Override
    public void onDataItemsModified(List<DataItem> dataItems) {
        subscriptionModel.onDataItemsModified(dataItems);
    }

    @Override
    public void onDataItemsDeleted(List<DataItem> dataItems) {
        subscriptionModel.onDataItemsDeleted(dataItems);
    }

    @Override
    public void onMonitoringModeChanged(List<MonitoredItem> monitoredItems) {
        subscriptionModel.onMonitoringModeChanged(monitoredItems);
    }
}
