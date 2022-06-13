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
package org.apache.plc4x.codegen.ast;

import java.util.Arrays;
import java.util.List;

public class IfStatement extends Statement {

    private List<Expression> condition;
    private List<Block> blocks;

    IfStatement(Expression condition, Block body, Block orElse) {
        this.condition = Arrays.asList(condition);
        if (orElse == null) {
            this.blocks = Arrays.asList(body);
        } else {
            this.blocks = Arrays.asList(body, orElse);
        }
    }

    IfStatement(List<Expression> condition, List<Block> blocks) {
        assert condition.size() == blocks.size() || condition.size() == (blocks.size() -1);
        this.condition = condition;
        this.blocks = blocks;
    }

    public List<Expression> getConditions() {
        return condition;
    }

    public List<Block> getBlocks() {
        return blocks;
    }

    @Override
    public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override
    public void write(Generator writer) {
        writer.generate(this);
    }

}
