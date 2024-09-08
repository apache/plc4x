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
package org.apache.plc4x.java.opcua.readwrite;

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

public abstract class UserIdentityTokenDefinition implements Message {

  // Abstract accessors for discriminator values.
  public abstract String getIdentifier();

  public UserIdentityTokenDefinition() {
    super();
  }

  protected abstract void serializeUserIdentityTokenDefinitionChild(WriteBuffer writeBuffer)
      throws SerializationException;

  public void serialize(WriteBuffer writeBuffer) throws SerializationException {
    PositionAware positionAware = writeBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();
    writeBuffer.pushContext("UserIdentityTokenDefinition");

    // Switch field (Serialize the sub-type)
    serializeUserIdentityTokenDefinitionChild(writeBuffer);

    writeBuffer.popContext("UserIdentityTokenDefinition");
  }

  @Override
  public int getLengthInBytes() {
    return (int) Math.ceil((float) getLengthInBits() / 8.0);
  }

  @Override
  public int getLengthInBits() {
    int lengthInBits = 0;
    UserIdentityTokenDefinition _value = this;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Length of sub-type elements will be added by sub-type...

    return lengthInBits;
  }

  public static UserIdentityTokenDefinition staticParse(ReadBuffer readBuffer, String identifier)
      throws ParseException {
    readBuffer.pullContext("UserIdentityTokenDefinition");
    PositionAware positionAware = readBuffer;
    boolean _lastItem = ThreadLocalHelper.lastItemThreadLocal.get();

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    UserIdentityTokenDefinitionBuilder builder = null;
    if (EvaluationHelper.equals(identifier, (String) "anonymous")) {
      builder =
          AnonymousIdentityToken.staticParseUserIdentityTokenDefinitionBuilder(
              readBuffer, identifier);
    } else if (EvaluationHelper.equals(identifier, (String) "username")) {
      builder =
          UserNameIdentityToken.staticParseUserIdentityTokenDefinitionBuilder(
              readBuffer, identifier);
    } else if (EvaluationHelper.equals(identifier, (String) "certificate")) {
      builder =
          X509IdentityToken.staticParseUserIdentityTokenDefinitionBuilder(readBuffer, identifier);
    } else if (EvaluationHelper.equals(identifier, (String) "identity")) {
      builder =
          IssuedIdentityToken.staticParseUserIdentityTokenDefinitionBuilder(readBuffer, identifier);
    }
    if (builder == null) {
      throw new ParseException(
          "Unsupported case for discriminated type"
              + " parameters ["
              + "identifier="
              + identifier
              + "]");
    }

    readBuffer.closeContext("UserIdentityTokenDefinition");
    // Create the instance
    UserIdentityTokenDefinition _userIdentityTokenDefinition = builder.build();
    return _userIdentityTokenDefinition;
  }

  public interface UserIdentityTokenDefinitionBuilder {
    UserIdentityTokenDefinition build();
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (!(o instanceof UserIdentityTokenDefinition)) {
      return false;
    }
    UserIdentityTokenDefinition that = (UserIdentityTokenDefinition) o;
    return true;
  }

  @Override
  public int hashCode() {
    return Objects.hash();
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
