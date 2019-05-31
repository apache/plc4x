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
import org.apache.plc4x.language.fields.ArrayField;
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

    public String getNullValueForType(TypeReference typeReference) {
        if(typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT: {
                    return "false";
                }
                case UINT: {
                    if (simpleTypeReference.getSize() <= 16) {
                        return "0";
                    }
                    if (simpleTypeReference.getSize() <= 32) {
                        return "0l";
                    }
                    return "null";
                }
                case INT: {
                    if (simpleTypeReference.getSize() <= 32) {
                        return "0";
                    }
                    if (simpleTypeReference.getSize() <= 64) {
                        return "0l";
                    }
                    return "null";
                }
                case FLOAT: {
                    if (simpleTypeReference.getSize() <= 32) {
                        return "0.0f";
                    }
                    if (simpleTypeReference.getSize() <= 64) {
                        return "0.0";
                    }
                    return "null";
                }
                case STRING: {
                    return "null";
                }
            }
            return "Hurz";
        } else {
            return "null";
        }
    }

    public String getIoBufferReadMethodCall(SimpleTypeReference simpleTypeReference) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT: {
                return "readBit()";
            }
            case UINT: {
                if (simpleTypeReference.getSize() <= 4) {
                    return "readUnsignedByte(" + simpleTypeReference.getSize() + ")";
                }
                if (simpleTypeReference.getSize() <= 8) {
                    return "readUnsignedShort(" + simpleTypeReference.getSize() + ")";
                }
                if (simpleTypeReference.getSize() <= 16) {
                    return "readUnsignedInt(" + simpleTypeReference.getSize() + ")";
                }
                if (simpleTypeReference.getSize() <= 32) {
                    return "readUnsignedLong(" + simpleTypeReference.getSize() + ")";
                }
                return "readUnsignedBigInteger" + simpleTypeReference.getSize() + ")";
            }
            case INT: {
                if (simpleTypeReference.getSize() <= 8) {
                    return "readByte" + simpleTypeReference.getSize() + ")";
                }
                if (simpleTypeReference.getSize() <= 16) {
                    return "readShort" + simpleTypeReference.getSize() + ")";
                }
                if (simpleTypeReference.getSize() <= 32) {
                    return "readInt" + simpleTypeReference.getSize() + ")";
                }
                if (simpleTypeReference.getSize() <= 64) {
                    return "readLong" + simpleTypeReference.getSize() + ")";
                }
                return "readBigInteger" + simpleTypeReference.getSize() + ")";
            }
            case FLOAT: {
                if (simpleTypeReference.getSize() <= 32) {
                    return "readFloat" + simpleTypeReference.getSize() + ")";
                }
                if (simpleTypeReference.getSize() <= 64) {
                    return "readDouble" + simpleTypeReference.getSize() + ")";
                }
                return "readBigDecimal" + simpleTypeReference.getSize() + ")";
            }
            case STRING: {
                return "readString" + simpleTypeReference.getSize() + ")";
            }
        }
        return "Hurz";
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

    public boolean isSimpleType(TypeReference typeReference) {
        return typeReference instanceof SimpleTypeReference;
    }

    public boolean isDiscriminatedType(TypeDefinition typeDefinition) {
        return typeDefinition instanceof DiscriminatedComplexTypeDefinition;
    }

    public boolean isCountArray(ArrayField arrayField) {
        return arrayField.getLengthType() == ArrayField.LengthType.COUNT;
    }

}
