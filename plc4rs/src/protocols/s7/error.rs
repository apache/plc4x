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

use thiserror::Error;

#[derive(Error, Debug)]
pub enum S7Error {
    #[error("Invalid message type: {0:#04x}")]
    InvalidMessageType(u8),

    #[error("Invalid function code: {0:#04x}")]
    InvalidFunctionCode(u8),

    #[error("Invalid parameter type: {0:#04x}")]
    InvalidParameterType(u8),

    #[error("Invalid protocol ID: expected 0x32, got {0:#04x}")]
    InvalidProtocolId(u8),

    #[error("Invalid TPKT version: expected 0x03, got {0:#04x}")]
    InvalidTpktVersion(u8),

    #[error("Invalid length: {0}")]
    InvalidLength(String),

    #[error("Parse error: {0}")]
    ParseError(String),

    #[error("Invalid variable specification: {0:#04x}")]
    InvalidVarSpec(u8),

    #[error("Invalid syntax ID: {0:#04x}")]
    InvalidSyntaxId(u8),

    #[error("Invalid area: {0:#04x}")]
    InvalidArea(u8),

    #[error("Invalid return code: {0:#04x}")]
    InvalidReturnCode(u8),

    #[error("Invalid transport size: {0:#04x}")]
    InvalidTransportSize(u8),
}
