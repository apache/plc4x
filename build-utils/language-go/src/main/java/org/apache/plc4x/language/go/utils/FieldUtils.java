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
package org.apache.plc4x.language.go.utils;

import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

public class FieldUtils {

    public static boolean contains(Field field, String label) {
        switch(field.getTypeName()) {
            case "array": {
                ArrayField arrayField = (ArrayField) field;
                return arrayField.getLoopExpression().contains(label);
            }
            case "checksum": {
                ChecksumField checksumField = (ChecksumField) field;
                return checksumField.getChecksumExpression().contains(label);
            }
            case "const": {
                ConstField constField = (ConstField) field;
                return false;
            }
            case "discriminator": {
                DiscriminatorField discriminatorField = (DiscriminatorField) field;
                return false;
            }
            case "enum": {
                EnumField enumField = (EnumField) field;
                return false;
            }
            case "implicit": {
                ImplicitField implicitField = (ImplicitField) field;
                return implicitField.getSerializeExpression().contains(label);
            }
            case "manualArray": {
                ManualArrayField manualArrayField = (ManualArrayField) field;
                return
                    manualArrayField.getSerializeExpression().contains(label) ||
                        manualArrayField.getParseExpression().contains(label) ||
                        manualArrayField.getLoopExpression().contains(label) ||
                        manualArrayField.getLengthExpression().contains(label);
            }
            case "manual": {
                ManualField manualField = (ManualField) field;
                return
                    manualField.getSerializeExpression().contains(label) ||
                        manualField.getParseExpression().contains(label) ||
                        manualField.getLengthExpression().contains(label);
            }
            case "optional": {
                OptionalField optionalField = (OptionalField) field;
                return optionalField.getConditionExpression().contains(label);
            }
            case "padding": {
                PaddingField paddingField = (PaddingField) field;
                return
                    paddingField.getPaddingCondition().contains(label) ||
                        paddingField.getPaddingValue().contains(label);
            }
            case "reserved": {
                ReservedField reservedField = (ReservedField) field;
                return false;
            }
            case "simple": {
                SimpleField simpleField = (SimpleField) field;
                return false;
            }
            case "switch": {
                SwitchField switchField = (SwitchField) field;
                for (Term discriminatorExpression : switchField.getDiscriminatorExpressions()) {
                    if(discriminatorExpression.contains(label)) {
                        return true;
                    }
                }
                return false;
            }
            case "virtual": {
                VirtualField virtualField = (VirtualField) field;
                return virtualField.getValueExpression().contains(label);
            }
        }
        return false;
    }


}
