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
package org.apache.plc4x.java.profinet.readwrite;

import static org.apache.plc4x.java.spi.codegen.fields.FieldReaderFactory.*;
import static org.apache.plc4x.java.spi.codegen.fields.FieldWriterFactory.*;
import static org.apache.plc4x.java.spi.codegen.io.DataReaderFactory.*;
import static org.apache.plc4x.java.spi.codegen.io.DataWriterFactory.*;
import static org.apache.plc4x.java.spi.generation.StaticHelper.*;

import java.time.*;
import java.util.*;
import org.apache.plc4x.java.api.exceptions.*;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.codegen.*;
import org.apache.plc4x.java.spi.codegen.fields.*;
import org.apache.plc4x.java.spi.codegen.io.*;
import org.apache.plc4x.java.spi.generation.*;

// Code generated by code-generation. DO NOT EDIT.

public class PnIoCm_ModuleDiffBlockApi implements Message {

  // Constant values.
  public static final Long API = 0x00000000L;

  // Properties.
  protected final List<PnIoCm_ModuleDiffBlockApi_Module> modules;

  public PnIoCm_ModuleDiffBlockApi(List<PnIoCm_ModuleDiffBlockApi_Module> modules) {
    super();
    this.modules = modules;
  }

  public List<PnIoCm_ModuleDiffBlockApi_Module> getModules() {
    return modules;
  }

  public long getApi() {
    return API;
  }

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    int startPos = positionAware.getPos();
    writeBuffer.pushContext("PnIoCm_ModuleDiffBlockApi");

    // Const Field (api)
    writeConstField("api", API, writeUnsignedLong(writeBuffer, 32));

    // Implicit Field (numModules) (Used for parsing, but its value is not stored as it's implicitly
    // given by the objects content)
    int numModules = (int) (COUNT(getModules()));
    writeImplicitField("numModules", numModules, writeUnsignedInt(writeBuffer, 16));

    // Array Field (modules)
    writeComplexTypeArrayField("modules", modules, writeBuffer);

    writeBuffer.popContext("PnIoCm_ModuleDiffBlockApi");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    PnIoCm_ModuleDiffBlockApi _value = this;

    // Const Field (api)
    lengthInBits += 32;

    // Implicit Field (numModules)
    lengthInBits += 16;

    // Array field
    if (modules != null) {
      int i = 0;
      for (PnIoCm_ModuleDiffBlockApi_Module element : modules) {
        boolean last = ++i >= modules.size();
        lengthInBits += element.getLengthInBits();
      }
    }

    return lengthInBits;
  }

  public static PnIoCm_ModuleDiffBlockApi staticParse(ReadBuffer readBuffer, Object... args)
      throws ParseException {
    PositionAware positionAware = readBuffer;
    return staticParse(readBuffer);
  }

  public static PnIoCm_ModuleDiffBlockApi staticParse(ReadBuffer readBuffer) throws ParseException {
    readBuffer.pullContext("PnIoCm_ModuleDiffBlockApi");
    PositionAware positionAware = readBuffer;
    int startPos = positionAware.getPos();
    int curPos;

    long api =
        readConstField("api", readUnsignedLong(readBuffer, 32), PnIoCm_ModuleDiffBlockApi.API);

    int numModules = readImplicitField("numModules", readUnsignedInt(readBuffer, 16));

    List<PnIoCm_ModuleDiffBlockApi_Module> modules =
        readCountArrayField(
            "modules",
            new DataReaderComplexDefault<>(
                () -> PnIoCm_ModuleDiffBlockApi_Module.staticParse(readBuffer), readBuffer),
            numModules);

    readBuffer.closeContext("PnIoCm_ModuleDiffBlockApi");
    // Create the instance
    PnIoCm_ModuleDiffBlockApi _pnIoCm_ModuleDiffBlockApi;
    _pnIoCm_ModuleDiffBlockApi = new PnIoCm_ModuleDiffBlockApi(modules);
    return _pnIoCm_ModuleDiffBlockApi;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof PnIoCm_ModuleDiffBlockApi)) {
      return false;
    }
    PnIoCm_ModuleDiffBlockApi that = (PnIoCm_ModuleDiffBlockApi) o;
    return (getModules() == that.getModules()) && true;
  }

  @Override
  public int hashCode() {
    return Objects.hash(getModules());
  }

  @Override
  public String toString() {
    WriteBufferBoxBased writeBufferBoxBased = new WriteBufferBoxBased(true, true);
    try {
      writeBufferBoxBased.writeSerializable(this);
    } catch (SerializationException e) {
      throw new RuntimeException(e);
    }
    return "\n" + writeBufferBoxBased.getBox().toString() + "\n";
  }
}
