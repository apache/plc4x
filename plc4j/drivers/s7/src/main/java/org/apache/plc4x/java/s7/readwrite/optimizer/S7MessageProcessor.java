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
package org.apache.plc4x.java.s7.readwrite.optimizer;

import io.vavr.control.Either;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.s7.readwrite.S7MessageRequest;
import org.apache.plc4x.java.s7.readwrite.S7MessageResponseData;

import java.util.Collection;
import java.util.Map;

/**
 * Some times the messages being sent have to be manipulated before
 * being able to send them. For example eventually a request has to
 * be split up into multiple ones to respect the maximum PDU size.
 */
public interface S7MessageProcessor {

    Collection<S7MessageRequest> processRequest(S7MessageRequest request, int pduSize) throws PlcException;

    S7MessageResponseData processResponse(S7MessageRequest originalRequest, Map<S7MessageRequest, Either<S7MessageResponseData, Throwable>> result) throws PlcException;

}
