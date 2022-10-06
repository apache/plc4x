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
package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;

public abstract class LineEntryNode extends Node {

    @JsonProperty("col_offset")
    private long colOffset;

    @JsonProperty("lineno")
    private long lineno;

    public long getColOffset() {
        return colOffset;
    }

    public void setColOffset(long colOffset) {
        this.colOffset = colOffset;
    }

    public long getLineno() {
        return lineno;
    }

    public void setLineno(long lineno) {
        this.lineno = lineno;
    }
}
