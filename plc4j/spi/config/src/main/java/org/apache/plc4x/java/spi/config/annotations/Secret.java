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
package org.apache.plc4x.java.spi.config.annotations;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * Marks a configuration field as carrying a secret - a password, a pre-shared key, a token, or any
 * other value that must never reach a log file, an audit-log entry, an exception message or a
 * {@code toString()} rendering.
 *
 * <p><strong>Annotate at the point of definition.</strong> This annotation is the single source of
 * truth for "is this value sensitive?". It replaces a pattern in {@code DriverBase} that guessed
 * from the parameter name ({@code password|passwd|secret|token|psk-key|passphrase}). That list
 * lived one module away from every declaration it protected, and could only ever be one parameter
 * behind - {@code psk-key} had to be added to it after the fact. A declaration that travels with
 * the parameter cannot drift the way a list in another module does.</p>
 *
 * <p><strong>What it drives:</strong></p>
 * <ol>
 *   <li>Connection-string redaction before logging, which derives the set of secret parameter
 *       names from the annotated fields of the driver's and the transports' configuration
 *       classes.</li>
 *   <li>{@code toString()} - an annotated field must render as {@code <redacted>}, never as its
 *       value. A reflective test plants a sentinel in every annotated field and fails the build if
 *       the sentinel appears in the output, so a hand-written {@code toString()} cannot silently
 *       leak a newly added secret.</li>
 * </ol>
 *
 * <p><strong>What is not a secret.</strong> An identifier that says <em>which</em> credential was
 * used is not one: {@code psk-identity} tells an operator which key the device refused, and hiding
 * it costs them the diagnosis while protecting nothing. Nor are store types, file paths, buffer
 * sizes or booleans. Marking everything credential-adjacent makes logs useless without making
 * anything safer.</p>
 *
 * <p><strong>Placement:</strong> on the field, beside {@link ConfigurationParameter}. A field
 * without a {@code @ConfigurationParameter} may still be annotated - it is then covered by the
 * {@code toString()} rule but contributes no parameter name, because there is nothing to match in
 * a connection string.</p>
 *
 * <p><strong>Limits - read before assuming coverage.</strong> This marks a <em>named parameter</em>.
 * A credential that reaches a driver by another route is not covered: the clear case is URI
 * userinfo ({@code s7://user:password@plc:102}), where the secret is part of the authority
 * component and has no parameter name at all. That is redacted structurally, independent of any
 * annotation.</p>
 *
 * @see ConfigurationParameter
 */
@Target(ElementType.FIELD)
@Retention(RetentionPolicy.RUNTIME)
public @interface Secret {
}
