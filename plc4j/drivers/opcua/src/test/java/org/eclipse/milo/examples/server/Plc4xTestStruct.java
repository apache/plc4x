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

import org.eclipse.milo.opcua.stack.core.UaSerializationException;
import org.eclipse.milo.opcua.stack.core.serialization.SerializationContext;
import org.eclipse.milo.opcua.stack.core.serialization.UaDecoder;
import org.eclipse.milo.opcua.stack.core.serialization.UaEncoder;
import org.eclipse.milo.opcua.stack.core.serialization.UaStructure;
import org.eclipse.milo.opcua.stack.core.serialization.codecs.GenericDataTypeCodec;
import org.eclipse.milo.opcua.stack.core.types.builtin.ExpandedNodeId;

import java.util.Objects;

/**
 * A simple custom (user-defined) OPC UA structure used to exercise the PLC4X driver's struct
 * (PlcStruct) read/write support. Its four fields cover the common scalar cases (string, signed
 * int, boolean, floating point). Registered in {@link Plc4xTestNamespace} together with a
 * StructureDefinition so the server exposes both the DataTypeDefinition attribute and the legacy
 * type dictionary.
 */
public class Plc4xTestStruct implements UaStructure {

    public static final ExpandedNodeId TYPE_ID = ExpandedNodeId.parse(String.format(
        "nsu=%s;s=%s", Plc4xTestNamespace.NAMESPACE_URI, "DataType.Plc4xTestStruct"));

    public static final ExpandedNodeId BINARY_ENCODING_ID = ExpandedNodeId.parse(String.format(
        "nsu=%s;s=%s", Plc4xTestNamespace.NAMESPACE_URI, "DataType.Plc4xTestStruct.BinaryEncoding"));

    private final String foo;
    private final int bar;
    private final boolean baz;
    private final double qux;

    public Plc4xTestStruct() {
        this("", 0, false, 0.0d);
    }

    public Plc4xTestStruct(String foo, int bar, boolean baz, double qux) {
        this.foo = foo;
        this.bar = bar;
        this.baz = baz;
        this.qux = qux;
    }

    public String getFoo() {
        return foo;
    }

    public int getBar() {
        return bar;
    }

    public boolean isBaz() {
        return baz;
    }

    public double getQux() {
        return qux;
    }

    @Override
    public ExpandedNodeId getTypeId() {
        return TYPE_ID;
    }

    @Override
    public ExpandedNodeId getBinaryEncodingId() {
        return BINARY_ENCODING_ID;
    }

    @Override
    public ExpandedNodeId getXmlEncodingId() {
        return ExpandedNodeId.NULL_VALUE;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        Plc4xTestStruct that = (Plc4xTestStruct) o;
        return bar == that.bar && baz == that.baz
            && Double.compare(that.qux, qux) == 0 && Objects.equals(foo, that.foo);
    }

    @Override
    public int hashCode() {
        return Objects.hash(foo, bar, baz, qux);
    }

    @Override
    public String toString() {
        return "Plc4xTestStruct{foo=" + foo + ", bar=" + bar + ", baz=" + baz + ", qux=" + qux + '}';
    }

    public static class Codec extends GenericDataTypeCodec<Plc4xTestStruct> {
        @Override
        public Class<Plc4xTestStruct> getType() {
            return Plc4xTestStruct.class;
        }

        @Override
        public Plc4xTestStruct decode(SerializationContext context, UaDecoder decoder)
                throws UaSerializationException {
            String foo = decoder.readString("Foo");
            int bar = decoder.readInt32("Bar");
            boolean baz = decoder.readBoolean("Baz");
            double qux = decoder.readDouble("Qux");
            return new Plc4xTestStruct(foo, bar, baz, qux);
        }

        @Override
        public void encode(SerializationContext context, UaEncoder encoder, Plc4xTestStruct value)
                throws UaSerializationException {
            encoder.writeString("Foo", value.foo);
            encoder.writeInt32("Bar", value.bar);
            encoder.writeBoolean("Baz", value.baz);
            encoder.writeDouble("Qux", value.qux);
        }
    }
}
