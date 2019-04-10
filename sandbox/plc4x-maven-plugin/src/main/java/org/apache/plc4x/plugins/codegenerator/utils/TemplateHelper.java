/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.plugins.codegenerator.utils;

import org.apache.plc4x.plugins.codegenerator.utils.model.ChoiceSegment;
import org.apache.plc4x.plugins.codegenerator.utils.model.FieldSegment;
import org.apache.plc4x.plugins.codegenerator.utils.model.RepeaterSegment;
import org.apache.plc4x.plugins.codegenerator.utils.model.Segment;

public class TemplateHelper {

    public boolean isFieldSegment(Segment segment) {
        return segment instanceof FieldSegment;
    }

    public boolean isChoiceSegment(Segment segment) {
        return segment instanceof ChoiceSegment;
    }

    public boolean isRepeaterSegment(Segment segment) {
        return segment instanceof RepeaterSegment;
    }

}
