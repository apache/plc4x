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

package org.apache.plc4x.java.opcuaserver;

import java.util.Optional;
import java.util.Set;
import java.util.function.Function;

import org.eclipse.milo.opcua.sdk.core.AccessLevel;
import org.eclipse.milo.opcua.sdk.server.Session;
import org.eclipse.milo.opcua.sdk.server.nodes.filters.AttributeFilter;
import org.eclipse.milo.opcua.sdk.server.nodes.filters.AttributeFilterContext.GetAttributeContext;
import org.eclipse.milo.opcua.stack.core.AttributeId;

public class RestrictedAccessFilter implements AttributeFilter {

    private static final Set<AccessLevel> INTERNAL_ACCESS = AccessLevel.READ_WRITE;

    private final Function<Object, Set<AccessLevel>> accessLevelsFn;

    public RestrictedAccessFilter(Function<Object, Set<AccessLevel>> accessLevelsFn) {
        this.accessLevelsFn = accessLevelsFn;
    }

    @Override
    public Object getAttribute(GetAttributeContext ctx, AttributeId attributeId) {
        if (attributeId == AttributeId.UserAccessLevel) {
            Optional<Object> identity = ctx.getSession().map(Session::getIdentityObject);

            Set<AccessLevel> accessLevels = identity.map(accessLevelsFn).orElse(INTERNAL_ACCESS);

            return AccessLevel.toValue(accessLevels);
        } else {
            return ctx.getAttribute(attributeId);
        }
    }

}
