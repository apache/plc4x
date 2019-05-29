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

package org.apache.plc4x.language.definitions;

import org.apache.plc4x.language.fields.ConstField;
import org.apache.plc4x.language.fields.Field;
import org.apache.plc4x.language.fields.SimpleField;

import java.util.List;

public interface ComplexTypeDefinition extends TypeDefinition {

    boolean isAbstract();

    /**
     * Get all fields no matter the type.
     *
     * @return all fields ;-)
     */
    List<Field> getFields();

    /**
     * Get only the fields which are of type SimpleField.
     *
     * @return all simple fields ;-)
     */
    List<SimpleField> getSimpleFields();

    List<ConstField> getConstFields();

    /**
     * In contrast to getSimpleFields, this also gets all simple fields of any parent type.
     *
     * @return all simple fields including any parents simple fields
     */
    List<SimpleField> getAllSimpleFields();

    /**
     * Returns all SimpleFields defined by any parent type.
     *
     * @return all parent types simple fields.
     */
    List<SimpleField> getParentSimpleFields();

}
