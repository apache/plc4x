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

package org.apache.plc4x.language.java;

import org.apache.plc4x.language.definitions.DiscriminatedComplexTypeDefinition;
import org.apache.plc4x.language.definitions.TypeDefinition;
import org.apache.plc4x.language.references.ComplexTypeReference;
import org.apache.plc4x.language.references.SimpleTypeReference;
import org.apache.plc4x.language.references.TypeReference;

public class JavaLanguageTemplateHelper {

    public String getLanguageTypeNameForSpecType(TypeReference typeReference) {
        if(typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT: {
                    return "boolean";
                }
                case UINT: {
                    if (simpleTypeReference.getSize() <= 4) {
                        return "byte";
                    }
                    if (simpleTypeReference.getSize() <= 8) {
                        return "short";
                    }
                    if (simpleTypeReference.getSize() <= 16) {
                        return "int";
                    }
                    if (simpleTypeReference.getSize() <= 32) {
                        return "long";
                    }
                    return "BigInteger";
                }
                case INT: {
                    if (simpleTypeReference.getSize() <= 8) {
                        return "byte";
                    }
                    if (simpleTypeReference.getSize() <= 16) {
                        return "short";
                    }
                    if (simpleTypeReference.getSize() <= 32) {
                        return "int";
                    }
                    if (simpleTypeReference.getSize() <= 64) {
                        return "long";
                    }
                    return "BigInteger";
                }
                case FLOAT: {
                    if (simpleTypeReference.getSize() <= 32) {
                        return "float";
                    }
                    if (simpleTypeReference.getSize() <= 64) {
                        return "double";
                    }
                    return "BigDecimal";
                }
                case STRING: {
                    return "String";
                }
            }
            return "Hurz";
        } else {
            return ((ComplexTypeReference) typeReference).getName();
        }
    }

    public String getReadMethodName(SimpleTypeReference simpleTypeReference) {
        String languageTypeName = getLanguageTypeNameForSpecType(simpleTypeReference);
        languageTypeName = languageTypeName.substring(0, 1).toUpperCase() + languageTypeName.substring(1);
        if(simpleTypeReference.getBaseType().equals(SimpleTypeReference.SimpleBaseType.UINT)) {
            return "readUnsigned" + languageTypeName;
        } else {
            return "read" + languageTypeName;

        }
    }

    public boolean isDiscriminatedType(TypeDefinition typeDefinition) {
        return typeDefinition instanceof DiscriminatedComplexTypeDefinition;
    }

}
