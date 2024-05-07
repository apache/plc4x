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
package org.apache.plc4x.language.python;

import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.math.NumberUtils;
import org.apache.commons.text.CaseUtils;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultArgument;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.*;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.DefaultStringLiteral;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.BaseFreemarkerLanguageTemplateHelper;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.FreemarkerException;
import org.apache.plc4x.plugins.codegenerator.protocol.freemarker.Tracer;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.references.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

@SuppressWarnings({"unused", "WeakerAccess"})
public class PythonLanguageTemplateHelper extends BaseFreemarkerLanguageTemplateHelper {

    private static final Logger LOGGER = LoggerFactory.getLogger(BaseFreemarkerLanguageTemplateHelper.class);

    // TODO: we could condense it to one import set as these can be emitted per template and are not hardcoded anymore

    public final SortedSet<String> requiredImports = new TreeSet<>();

    public final SortedSet<String> requiredImportsForDataIo = new TreeSet<>();

    public PythonLanguageTemplateHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types) {
        super(thisType, protocolName, flavorName, types);
    }

    public String fileName(String protocolName, String languageName, String languageFlavorName) {
        return String.join("", protocolName.split("\\-")) + "." +
            String.join("", languageFlavorName.split("\\-"));
    }

    public String getSanitizedPackageName() {
        String sanitizedName = getProtocolName().replaceAll("-", "");
        sanitizedName = sanitizedName.replaceAll("\\.", "/");
        sanitizedName = sanitizedName.toLowerCase();
        return sanitizedName;
    }

    // TODO: check if protocol name can be enforced to only contain valid chars
    public String getSanitizedProtocolName() {
        String sanitizedName = getProtocolName().replaceAll("-", "");
        sanitizedName = CaseUtils.toCamelCase(sanitizedName, false, '.');
        return sanitizedName;
    }

    public String packageName() {
        return packageName(protocolName, "python", flavorName);
    }

    public String packageName(String protocolName, String languageName, String languageFlavorName) {
        return String.join("", protocolName.split("-")) + "." +
                String.join("", languageFlavorName.split("-"));
    }

    @Override
    public String getLanguageTypeNameForField(Field field) {
        boolean optional = field instanceof OptionalField;
        // If the referenced type is a DataIo type, the value is of type PlcValue.
        if (field instanceof PropertyField) {
            PropertyField propertyField = (PropertyField) field;
            if (propertyField.getType() instanceof ComplexTypeReference) {
                ComplexTypeReference complexTypeReference = (ComplexTypeReference) propertyField.getType();
                final TypeDefinition typeDefinition = complexTypeReference.getTypeDefinition();
                if (typeDefinition instanceof DataIoTypeDefinition) {
                    return "PlcValue";
                }
            }
        }
        TypedField typedField = field.asTypedField().orElseThrow();
        String encoding = null;
        Optional<Term> encodingAttribute = field.getAttribute("encoding");
        if(encodingAttribute.isPresent()) {
            encoding = encodingAttribute.get().toString();
        }
        return getLanguageTypeNameForTypeReference(typedField.getType(), encoding);
    }

    public boolean isComplex(Field field) {
        return field instanceof PropertyField && ((PropertyField) field).getType() instanceof NonSimpleTypeReference;
    }

    @Override
    public String getLanguageTypeNameForTypeReference(TypeReference typeReference) {
        return getLanguageTypeNameForTypeReference(typeReference, null);
    }

    public String getNonPrimitiveLanguageTypeNameForField(TypedField field) {
        return getLanguageTypeNameForTypeReference(field.getType(), false);
    }

    public String getLanguageTypeNameForTypeReference(TypeReference typeReference, boolean allowPrimitives) {
        return getLanguageTypeNameForTypeReference(typeReference, null);
    }

    public String getLanguageTypeNameForTypeReference(TypeReference typeReference, String encoding) {
        if (typeReference == null) {
            // TODO: shouldn't this be an error case
            return "";
        }
        if (typeReference.isArrayTypeReference()) {
            final ArrayTypeReference arrayTypeReference = (ArrayTypeReference) typeReference;
            TypeReference elementTypeReference = arrayTypeReference.getElementTypeReference();
            emitRequiredImport("from typing import List");
            return "List[" + getLanguageTypeNameForTypeReference(elementTypeReference) + "]";
        }
        if (typeReference.isNonSimpleTypeReference()) {
            emitRequiredImport("from plc4py.protocols." + protocolName + "." + flavorName.replace("-", "") + "." + typeReference.asNonSimpleTypeReference().orElseThrow().getName() + " import " + typeReference.asNonSimpleTypeReference().orElseThrow().getName());
            return typeReference.asNonSimpleTypeReference().orElseThrow().getName();
        }
        if (typeReference instanceof ByteOrderTypeReference) {
            return "binary.byteOrder";
        }
        SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "bool";
            case BYTE:
                return "int";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = simpleTypeReference.asIntegerTypeReference().orElseThrow();
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "int";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "int";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "int";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 64) {
                    return "int";
                }
                return "int";
            case INT:
                IntegerTypeReference integerTypeReference = simpleTypeReference.asIntegerTypeReference().orElseThrow();
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "int";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "int";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "int";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "int";
                }
                return "int";
            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = simpleTypeReference.asFloatTypeReference().orElseThrow();
                int sizeInBits = floatTypeReference.getSizeInBits();
                if (sizeInBits <= 32) {
                    return "float";
                }
                if (sizeInBits <= 64) {
                    return "float";
                }
                return "float";
            case STRING:
            case VSTRING:
                return "str";
            case TIME:
            case DATE:
            case DATETIME:
                return "time.Time";
            default:
                throw new FreemarkerException("Unsupported simple type");
        }
    }

    public String getReservedValue(ReservedField reservedField) {
        final String languageTypeName = getLanguageTypeNameForTypeReference(reservedField.getType());
        return languageTypeName + "(" + reservedField.getReferenceValue() + ")";
    }

    public String getFieldOptions(TypedField field, List<Argument> parserArguments) {
        StringBuilder sb = new StringBuilder();
        final Optional<Term> encodingOptional = field.getEncoding();
        if (encodingOptional.isPresent()) {
            final String encoding = toParseExpression(field, field.getType(), encodingOptional.get(), parserArguments);
            sb.append(", encoding='").append(encoding).append("'");
        }
        final Optional<Term> byteOrderOptional = field.getByteOrder();
        if (byteOrderOptional.isPresent()) {
            final String byteOrder = toParseExpression(field, field.getType(), byteOrderOptional.get(), parserArguments);
            emitRequiredImport("from plc4py.utils.GenericTypes import ByteOrder");
            sb.append(", byte_order=ByteOrder.").append(byteOrder);
        }
        return sb.toString();
    }

    public String getDataReaderCall(TypeReference typeReference) {
        return getDataReaderCall(typeReference, "enumForValue", false);
    }

    public String getDataReaderCall(TypeReference typeReference, String resolverMethod, Boolean isArray) {
        if (typeReference.isEnumTypeReference()) {
            final String languageTypeName = getLanguageTypeNameForTypeReference(typeReference);
            final SimpleTypeReference enumBaseTypeReference = getEnumBaseTypeReference(typeReference);
            return "read_enum(read_function=" + languageTypeName + ",";
        } else if (typeReference.isArrayTypeReference()) {
            final ArrayTypeReference arrayTypeReference = typeReference.asArrayTypeReference().orElseThrow();
            return getDataReaderCall(arrayTypeReference.getElementTypeReference(), resolverMethod, true);
        } else if (typeReference.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow(IllegalStateException::new);
            return getDataReaderCall(simpleTypeReference);
        } else if (typeReference.isComplexTypeReference()) {
            StringBuilder paramsString = new StringBuilder();
            ComplexTypeReference complexTypeReference = typeReference.asComplexTypeReference().orElseThrow(IllegalStateException::new);
            ComplexTypeDefinition typeDefinition = complexTypeReference.getTypeDefinition();
            String parserCallString = getLanguageTypeNameForTypeReference(typeReference);
            // In case of DataIo we actually need to use the type name and not what above returns.
            // (In this case the mspec type name and the result type name differ)
            if (typeReference.isDataIoTypeReference()) {
                parserCallString = typeReference.asDataIoTypeReference().orElseThrow().getName();
            }
            if (typeDefinition.isDiscriminatedChildTypeDefinition()) {
                parserCallString = "(" + getLanguageTypeNameForTypeReference(typeReference) + ") " + typeDefinition.getParentType().orElseThrow().getName();
            }
            List<Term> paramTerms = complexTypeReference.getParams().orElse(Collections.emptyList());
            for (int i = 0; i < paramTerms.size(); i++) {
                Term paramTerm = paramTerms.get(i);
                final TypeReference argumentType = getArgumentType(complexTypeReference, i);
                paramsString
                    .append(", ")
                    .append(toParseExpression(null, argumentType, paramTerm, null))
                    .append("=")
                    .append(toParseExpression(null, argumentType, paramTerm, null));
            }
            if (isArray) {
                return parserCallString + ".static_parse, ";
            } else {
                return "read_complex(read_function=" + parserCallString + ".static_parse, ";
            }
        } else {
            throw new IllegalStateException("What is this type? " + typeReference);
        }
    }

    public String getDataReaderCall(SimpleTypeReference simpleTypeReference) {
        final int sizeInBits = simpleTypeReference.getSizeInBits();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "read_bit";
            case BYTE:
                return "read_byte";
            case UINT:
                if (sizeInBits <= 8) return "read_unsigned_byte";
                if (sizeInBits <= 16) return "read_unsigned_short";
                if (sizeInBits <= 32) return "read_unsigned_int";
                return "read_unsigned_long";
            case INT:
                if (sizeInBits <= 8) return "read_signed_byte";
                if (sizeInBits <= 16) return "read_signed_short";
                if (sizeInBits <= 32) return "read_signed_int";
                return "read_signed_long";
            case FLOAT:
                if (sizeInBits <= 32) return "read_float";
                return "read_double";
            case STRING:
                return "read_str";
            case VSTRING:
                VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                return "read_str";
            case TIME:
                return "read_time";
            case DATE:
                return "read_date";
            case DATETIME:
                return "read_date_time";
            default:
                throw new UnsupportedOperationException("Unsupported type " + simpleTypeReference.getBaseType());
        }
    }

    public String getDataWriterCall(TypeReference typeReference, String fieldName) {
        if (typeReference.isSimpleTypeReference()) {
            SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow(IllegalStateException::new);
            return getDataWriterCall(simpleTypeReference);
        } else if (typeReference.isArrayTypeReference()) {
            final ArrayTypeReference arrayTypeReference = typeReference.asArrayTypeReference().orElseThrow();
            return getDataWriterCall(arrayTypeReference.getElementTypeReference(), fieldName);
        } else if (typeReference.isComplexTypeReference()) {
            return "write_serializable";
        } else {
            throw new IllegalStateException("What is this type? " + typeReference);
        }
    }

    public String getEnumDataWriterCall(EnumTypeReference typeReference, String fieldName, String attributeName) {
        if (!typeReference.isEnumTypeReference()) {
            throw new IllegalArgumentException("this method should only be called for enum types");
        }
        final String languageTypeName = getLanguageTypeNameForTypeReference(typeReference);
        SimpleTypeReference outputTypeReference;
        if ("value".equals(attributeName)) {
            outputTypeReference = getEnumBaseTypeReference(typeReference);
        } else {
            outputTypeReference = getEnumFieldSimpleTypeReference(typeReference.asNonSimpleTypeReference().orElseThrow(), attributeName);
        }
        return getDataWriterCall(outputTypeReference, fieldName);
    }

    public String getDataWriterCall(SimpleTypeReference simpleTypeReference) {
        final int sizeInBits = simpleTypeReference.getSizeInBits();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "write_bit";
            case BYTE:
                return "write_byte";
            case UINT:
                if (sizeInBits <= 8) return "write_unsigned_byte";
                if (sizeInBits <= 16) return "write_unsigned_short";
                if (sizeInBits <= 32) return "write_unsigned_int";
                return "write_unsigned_long";
            case INT:
                if (sizeInBits <= 8) return "write_signed_byte";
                if (sizeInBits <= 16) return "write_signed_short";
                if (sizeInBits <= 32) return "write_signed_int";
                return "write_signed_long";
            case FLOAT:
                if (sizeInBits <= 32) return "write_float";
                return "write_double";
            case STRING:
                return "write_str";
            case VSTRING:
                VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                return "write_str";
            case TIME:
                return "write_time";
            case DATE:
                return "write_date";
            case DATETIME:
                return "write_date_time";
            default:
                return "";
        }
    }

    public String getPlcValueTypeForTypeReference(TypeReference typeReference) {
        if (typeReference == null) {
            // TODO: shouldn't this be an error case
            return "";
        }
        if (typeReference.isNonSimpleTypeReference()) {
            return ((NonSimpleTypeReference) typeReference).getName();
        }
        SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                emitRequiredImport("from plc4py.spi.values.PlcValues import PlcBOOL");
                return "PlcBOOL";
            case BYTE:
                emitRequiredImport("from plc4py.spi.values.PlcValues import PlcSINT");
                return "PlcSINT";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 4) {
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcUSINT");
                    return "PlcUSINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcUINT");
                    return "PlcUINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcUDINT");
                    return "PlcUDINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcULINT");
                    return "PlcULINT";
                }
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcSINT");
                    return "PlcSINT";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcINT");
                    return "PlcINT";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcDINT");
                    return "PlcDINT";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcLINT");
                    return "PlcLINT";
                }

            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                int sizeInBits = floatTypeReference.getSizeInBits();
                if (sizeInBits <= 32) {
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcREAL");
                    return "PlcREAL";
                }
                if (sizeInBits <= 64) {
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcLREAL");
                    return "PlcLREAL";
                }
            case STRING:
            case VSTRING:
                emitRequiredImport("from plc4py.spi.values.PlcValues import PlcSTRING");
                return "PlcSTRING";
            case TIME:
            case DATE:
            case DATETIME:
                emitRequiredImport("from plc4py.spi.values.PlcValues import PlcTIME");
                return "PlcTIME";
            default:
                return "";
        }
    }

    @Override
    public String getNullValueForTypeReference(TypeReference typeReference) {
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT:
                    return "False";
                case BYTE:
                case UINT:
                case INT:
                    return "0";
                case FLOAT:
                    return "0.0";
                case STRING:
                case VSTRING:
                    return "\"\"";
            }
        } else if (typeReference.isEnumTypeReference()) {
            return "0";
        }
        return "None";
    }

    public int getNumBits(SimpleTypeReference simpleTypeReference) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return 1;
            case BYTE:
                return 8;
            case UINT:
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                return integerTypeReference.getSizeInBits();
            case FLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                return floatTypeReference.getSizeInBits();
            case STRING:
            case VSTRING:
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return stringTypeReference.getSizeInBits();
            default:
                return 0;
        }
    }

    public boolean needsPointerAccess(PropertyField field) {
        boolean isAnTypeOfOptional = "optional".equals(field.getTypeName());
        return isAnTypeOfOptional && needPointerAccess(field.getType());
    }

    public boolean needPointerAccess(TypeReference typeReference) {
        boolean isNotAnComplexTypeReference = !typeReference.isComplexTypeReference();
        boolean arrayTypeIsNotAnComplexTypeReference = !(typeReference.isArrayTypeReference() && typeReference.asArrayTypeReference().orElseThrow().getElementTypeReference().isComplexTypeReference());
        return isNotAnComplexTypeReference && arrayTypeIsNotAnComplexTypeReference;
    }

    public String getSpecialReadBufferReadMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, TypedField field) {
        return getReadBufferReadMethodCall(logicalName, simpleTypeReference, null, field);
    }

    public String getReadBufferReadMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, TypedField field) {
        return getReadBufferReadMethodCall(logicalName, simpleTypeReference, null, field);
    }

    @Override
    public String getReadBufferReadMethodCall(SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        return getReadBufferReadMethodCall("", simpleTypeReference, valueString, field);
    }

    public String getReadBufferReadMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                String bitType = "bit";
                return "read_buffer.read_" + bitType + "(\"" + logicalName + "\")";
            case BYTE:
                String byteType = "byte";
                return "read_buffer.read_" + byteType + "(\"" + logicalName + "\")";
            case UINT:
                String unsignedIntegerType;
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 4) {
                    unsignedIntegerType = "unsigned_byte";
                } else if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    unsignedIntegerType = "unsigned_short";
                } else if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    unsignedIntegerType = "unsigned_int";
                } else if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    unsignedIntegerType = "unsigned_long";
                } else {
                    unsignedIntegerType = "unsigned_long";
                }
                return "read_buffer.read_" + unsignedIntegerType + "(" + simpleTypeReference.getSizeInBits() + ", logical_name=\"" + logicalName + "\")";
            case INT:
                String integerType;
                if (simpleTypeReference.getSizeInBits() <= 8) {
                    integerType = "signed_byte";
                } else if (simpleTypeReference.getSizeInBits() <= 16) {
                    integerType = "short";
                } else if (simpleTypeReference.getSizeInBits() <= 32) {
                    integerType = "int";
                } else if (simpleTypeReference.getSizeInBits() <= 64) {
                    integerType = "long";
                } else {
                    integerType = "long";
                }
                return "read_buffer.read_" + integerType + "(" + simpleTypeReference.getSizeInBits() + ", logical_name=\"" + logicalName + "\")";
            case FLOAT:
                String floatType = (simpleTypeReference.getSizeInBits() <= 32) ? "float" : "double";
                return "read_buffer.read_" + floatType + "(" + simpleTypeReference.getSizeInBits() + ", logical_name=\"" + logicalName + "\")";
            case STRING:
            case VSTRING:
                String stringType = "str";
                final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                if (!(encodingTerm instanceof StringLiteral)) {
                    throw new FreemarkerException("Encoding must be a quoted string value");
                }
                String encoding = ((StringLiteral) encodingTerm).getValue();
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                if (simpleTypeReference.getBaseType() == SimpleTypeReference.SimpleBaseType.VSTRING) {
                    VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                    length = toParseExpression(field, INT_TYPE_REFERENCE, vstringTypeReference.getLengthExpression(), null);
                }
                return "read_buffer.read_" + stringType + "(" + simpleTypeReference.getSizeInBits() + ", logical_name=\"" + logicalName + "\", encoding=" + "\"\")";

            default:
                return "";
        }
    }

    @Override
    public String getWriteBufferWriteMethodCall(SimpleTypeReference simpleTypeReference, String fieldName, TypedField field) {
        // Fallback if somewhere the method gets called without a name
        String logicalName = fieldName.replaceAll("[\"()*]", "").replaceFirst("_", "");
        return getWriteBufferWriteMethodCall(logicalName, simpleTypeReference, fieldName, field);
    }

    public String getWriteBufferWriteMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, Term valueTerm, TypedField field, String... writerArgs) {
        if (valueTerm instanceof BooleanLiteral) {
            return getWriteBufferWriteMethodCall(logicalName, simpleTypeReference, Boolean.toString(((BooleanLiteral) valueTerm).getValue()), field, writerArgs);
        }
        if (valueTerm instanceof NumericLiteral) {
            return getWriteBufferWriteMethodCall(logicalName, simpleTypeReference, ((NumericLiteral) valueTerm).getNumber().toString(), field, writerArgs);
        }
        if (valueTerm instanceof HexadecimalLiteral) {
            return getWriteBufferWriteMethodCall(logicalName, simpleTypeReference, ((HexadecimalLiteral) valueTerm).getHexString(), field, writerArgs);
        }
        if (valueTerm instanceof StringLiteral) {
            return getWriteBufferWriteMethodCall(logicalName, simpleTypeReference, "\"" + ((StringLiteral) valueTerm).getValue() + "\"", field, writerArgs);
        }
        throw new FreemarkerException("Outputting " + valueTerm.toString() + " not implemented yet. Please continue defining other types in the PythonLanguageHelper.getWriteBufferWriteMethodCall.");
    }

    public String getWriteBufferWriteMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, String fieldName, TypedField field, String... writerArgs) {
        String writerArgsString = "";
        if (writerArgs.length > 0) {
            writerArgsString += ", " + StringUtils.join(writerArgs, ", ");
        }
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                return "write_buffer.write_bit(" + fieldName + ", \"" + logicalName + "\"" + writerArgsString + ")";
            case BYTE:
                return "write_buffer.write_byte(" + fieldName + ", \"" + logicalName + "\"" + writerArgsString + ")";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    return "write_buffer.write_byte(" + fieldName + ", " + unsignedIntegerTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    return "write_buffer.write_unsigned_short(" + fieldName + ", " + unsignedIntegerTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    return "write_buffer.write_unsigned_int(" + fieldName + ", " + unsignedIntegerTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 64) {
                    return "write_buffer.write_unsigned_long(" + fieldName + ", " + unsignedIntegerTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
                }
                return "write_buffer.write_unsigned_long(" + fieldName + ", " + unsignedIntegerTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    return "write_buffer.write_signed_byte(" + fieldName + ", " + integerTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    return "write_buffer.write_short(" + fieldName + ", " + integerTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    return "write_buffer.write_int(" + fieldName + ", " + integerTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    return "write_buffer.write_long(" + fieldName + ", " + integerTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
                }
                return "write_buffer.write_long(" + fieldName + ", " + integerTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                if (floatTypeReference.getSizeInBits() <= 32) {
                    return "write_buffer.write_float(" + fieldName + ", " + floatTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
                }
                if (floatTypeReference.getSizeInBits() <= 64) {
                    return "write_buffer.write_double(" + fieldName + ", " + floatTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
                }
                return "write_buffer.write_double(" + fieldName + ", " + floatTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
            case STRING: {
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                String encoding = encodingTerm.asLiteral()
                    .orElseThrow(() -> new FreemarkerException("Encoding must be a literal"))
                    .asStringLiteral()
                    .orElseThrow(() -> new FreemarkerException("Encoding must be a quoted string value")).getValue();
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                return "write_buffer.write_str(" + fieldName + ", " + length + ", \"" +
                    encoding + "\", \"" + logicalName + "\"" + writerArgsString + ")";
            }
            case VSTRING: {
                VstringTypeReference vstringTypeReference = (VstringTypeReference) simpleTypeReference;
                final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                String encoding = encodingTerm.asLiteral()
                    .orElseThrow(() -> new FreemarkerException("Encoding must be a literal"))
                    .asStringLiteral()
                    .orElseThrow(() -> new FreemarkerException("Encoding must be a quoted string value")).getValue();
                String lengthExpression = toExpression(field, null, vstringTypeReference.getLengthExpression(), null, Collections.singletonList(new DefaultArgument("stringLength", new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.INT, 32))), true, false);
                String length = Integer.toString(simpleTypeReference.getSizeInBits());
                return "write_buffer.write_str(" + fieldName + ", " + lengthExpression + ", \"" +
                    encoding + "\", \"" + logicalName + "\"" + writerArgsString + ")";
            }
            case DATE:
            case TIME:
            case DATETIME:
                return "write_buffer.write_unsigned_int(uint32(" + fieldName + "), \"" + logicalName + "\")" + writerArgsString + ")";
            default:
                throw new FreemarkerException("Unsupported base type " + simpleTypeReference.getBaseType());
        }
    }


    public String toParseExpression(Field field, TypeReference resultType, Term term, List<Argument> parserArguments) {
        Tracer tracer = pythonTracerStart("toParseExpression");
        return tracer + toTypedParseExpression(field, resultType, term, parserArguments);
    }

    public String toParseExpression(Field field, TypeReference resultType, Term term, List<Argument> parserArguments, boolean suppressPointerAccess) {
        Tracer tracer = pythonTracerStart("toParseExpression");
        return tracer + toTypedParseExpression(field, resultType, term, parserArguments, suppressPointerAccess);
    }

    public String toSerializationExpression(Field field, TypeReference resultType, Term term, List<Argument> serializerArguments) {
        Tracer tracer = pythonTracerStart("toSerializationExpression");
        return tracer + toTypedSerializationExpression(field, resultType, term, serializerArguments);
    }

    public String toBooleanParseExpression(Field field, Term term, List<Argument> parserArguments) {
        Tracer tracer = pythonTracerStart("toBooleanParseExpression");
        return tracer + toTypedParseExpression(field, new DefaultBooleanTypeReference(), term, parserArguments);
    }

    public String toBooleanSerializationExpression(Field field, Term term, List<Argument> serializerArguments) {
        Tracer tracer = pythonTracerStart("toBooleanSerializationExpression");
        return tracer + toTypedSerializationExpression(field, new DefaultBooleanTypeReference(), term, serializerArguments);
    }

    public String toIntegerParseExpression(Field field, int sizeInBits, Term term, List<Argument> parserArguments) {
        Tracer tracer = pythonTracerStart("toIntegerParseExpression");
        return tracer + toTypedParseExpression(field, new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.UINT, sizeInBits), term, parserArguments);
    }

    public String toIntegerSerializationExpression(Field field, int sizeInBits, Term term, List<Argument> serializerArguments) {
        Tracer tracer = pythonTracerStart("toIntegerSerializationExpression");
        return tracer + toTypedSerializationExpression(field, new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.UINT, sizeInBits), term, serializerArguments);
    }

    public String toTypedParseExpression(Field field, TypeReference fieldType, Term term, List<Argument> parserArguments) {
        Tracer tracer = pythonTracerStart("toTypedParseExpression");
        return tracer + toExpression(field, fieldType, term, parserArguments, null, false, fieldType != null && fieldType.isComplexTypeReference());
    }

    public String toTypedParseExpression(Field field, TypeReference fieldType, Term term, List<Argument> parserArguments, boolean suppressPointerAccess) {
        Tracer tracer = pythonTracerStart("toTypedParseExpression");
        return tracer + toExpression(field, fieldType, term, parserArguments, null, false, suppressPointerAccess);
    }

    public String toTypedSerializationExpression(Field field, TypeReference fieldType, Term term, List<Argument> serializerArguments) {
        Tracer tracer = pythonTracerStart("toTypedSerializationExpression");
        return tracer + toExpression(field, fieldType, term, null, serializerArguments, true, false);
    }

    String getCastExpressionForTypeReference(TypeReference typeReference) {
        Tracer tracer = pythonTracerStart("castExpression");
        if (typeReference instanceof SimpleTypeReference) {
            return tracer.dive("simpleTypeRef") + getLanguageTypeNameForTypeReference(typeReference);
        } else if (typeReference instanceof ByteOrderTypeReference) {
            return tracer.dive( "byteOrderTypeRef") + "binary.ByteOrder";
        } else if (typeReference != null) {
            return tracer.dive("anyTypeRef") + "Cast" + getLanguageTypeNameForTypeReference(typeReference);
        } else {
            return tracer.dive("noTypeRef") + "";
        }
    }

    private String toExpression(Field field, TypeReference fieldType, Term term, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess) {
        Tracer tracer = pythonTracerStart("toExpression(suppressPointerAccess=" + suppressPointerAccess + ")");
        if (term == null) {
            return "";
        }
        if (term instanceof Literal) {
            return toLiteralTermExpression(field, fieldType, term, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        } else if (term instanceof UnaryTerm) {
            return toUnaryTermExpression(field, fieldType, (UnaryTerm) term, parserArguments, serializerArguments, serialize, tracer);
        } else if (term instanceof BinaryTerm) {
            return toBinaryTermExpression(field, fieldType, (BinaryTerm) term, parserArguments, serializerArguments, serialize, tracer);
        } else if (term instanceof TernaryTerm) {
            return toTernaryTermExpression(field, fieldType, (TernaryTerm) term, parserArguments, serializerArguments, serialize, tracer);
        } else {
            throw new FreemarkerException("Unsupported Term type " + term.getClass().getName());
        }
    }

    private String toTernaryTermExpression(Field field, TypeReference fieldType, TernaryTerm ternaryTerm, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("ternary term instanceOf");
        if ("if".equals(ternaryTerm.getOperation())) {
            Term a = ternaryTerm.getA();
            Term b = ternaryTerm.getB();
            Term c = ternaryTerm.getC();
            String castExpressionForTypeReference = getCastExpressionForTypeReference(fieldType);
            String inlineIf = toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, false) + " if " +
                toExpression(field, new DefaultBooleanTypeReference(), a, parserArguments, serializerArguments, serialize, false) + " else " +
                toExpression(field, fieldType, c, parserArguments, serializerArguments, serialize, false);

            return tracer + inlineIf;
        } else {
            throw new FreemarkerException("Unsupported ternary operation type " + ternaryTerm.getOperation());
        }
    }

    private String toBinaryTermExpression(Field field, TypeReference fieldType, BinaryTerm binaryTerm, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("binary term instanceOf");
        Term a = binaryTerm.getA();
        Term b = binaryTerm.getB();
        String operation = binaryTerm.getOperation();
        String castExpressionForTypeReference = getCastExpressionForTypeReference(fieldType);
        switch (operation) {
            case "^":
                tracer = tracer.dive("^");
                emitRequiredImport("from math import pow");
                return tracer + "pow(" +
                    castExpressionForTypeReference + "(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + "), " +
                    castExpressionForTypeReference + "(" + toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, false) + "))";
            // If we start casting for comparisons, equals or non equals, really messy things happen.
            case "==":
            case "!=":
            case ">":
            case "<":
            case ">=":
            case "<=":
                tracer = tracer.dive("compare");
                // For every access of optional elements we need pointer access ...
                // Except for doing a nil or not-nil check :-(
                // So in case of such a check, we need to suppress the pointer-access.
                boolean suppressPointerAccessOverride = (operation.equals("==") || operation.equals("!=")) && ((a instanceof NullLiteral) || (b instanceof NullLiteral));
                String aExpression = toExpression(field, null, a, parserArguments, serializerArguments, serialize, suppressPointerAccessOverride);
                String bExpression = toExpression(field, null, b, parserArguments, serializerArguments, serialize, suppressPointerAccessOverride);
                return tracer + "bool((" + aExpression + ") " + operation + " (" + bExpression + "))";
            case ">>":
            case "<<":
            case "|":
            case "&":
                tracer = tracer.dive("bitwise");
                // We don't want casts here
                return tracer +
                    toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) +
                    operation + " " +
                    toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, false);
            default:
                tracer = tracer.dive("default");
                if (fieldType instanceof StringTypeReference) {
                    tracer = tracer.dive("string type reference");
                    return tracer + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) +
                        operation + " " +
                        toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, false);
                }
                return tracer +
                    toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) +
                    operation + " " +
                    toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, false);
        }
    }

    private String toUnaryTermExpression(Field field, TypeReference fieldType, UnaryTerm unaryTerm, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("unary term instanceOf");
        Term a = unaryTerm.getA();
        switch (unaryTerm.getOperation()) {
            case "!":
                tracer = tracer.dive("case !");
                return tracer + "!(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + ")";
            case "-":
                tracer = tracer.dive("case -");
                return tracer + "-(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + ")";
            case "()":
                tracer = tracer.dive("case ()");
                return tracer + "(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + ")";
            default:
                throw new FreemarkerException("Unsupported unary operation type " + unaryTerm.getOperation());
        }
    }

    private String toLiteralTermExpression(Field field, TypeReference fieldType, Term term, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("literal term instanceOf");
        if (term instanceof NullLiteral) {
            tracer = tracer.dive("null literal instanceOf");
            return tracer + "None";
        } else if (term instanceof BooleanLiteral) {
            tracer = tracer.dive("boolean literal instanceOf");
            String bool = Boolean.toString(((BooleanLiteral) term).getValue());
            return tracer + bool.substring(0,1).toUpperCase() + bool.substring(1);
        } else if (term instanceof NumericLiteral) {
            tracer = tracer.dive("numeric literal instanceOf");
            if (getCastExpressionForTypeReference(fieldType).equals("string")) {
                tracer = tracer.dive("type reference string");
                return tracer + "(" + ((NumericLiteral) term).getNumber().toString() + ")";
            } else {
                return tracer + getCastExpressionForTypeReference(fieldType) + "(" + ((NumericLiteral) term).getNumber().toString() + ")";
            }
        } else if (term instanceof HexadecimalLiteral) {
            tracer = tracer.dive("hexadecimal literal instanceOf");
            return tracer + ((HexadecimalLiteral) term).getHexString();
        } else if (term instanceof StringLiteral) {
            tracer = tracer.dive("string literal instanceOf");
            return tracer + "\"" + ((StringLiteral) term).getValue() + "\"";
        } else if (term instanceof VariableLiteral) {
            tracer = tracer.dive("variable literal instanceOf");
            VariableLiteral variableLiteral = (VariableLiteral) term;
            if ("cur_pos".equals(((VariableLiteral) term).getName())) {
                return "(position_aware.get_pos() - startPos)";
            } else if ("BIG_ENDIAN".equals(((VariableLiteral) term).getName()) && (fieldType instanceof ByteOrderTypeReference)) {
                return "ByteOrder.BIG_ENDIAN";
            } else if ("LITTLE_ENDIAN".equals(((VariableLiteral) term).getName()) && (fieldType instanceof ByteOrderTypeReference)) {
                return "ByteOrder.LITTLE_ENDIAN";
            }
            return tracer + toVariableExpression(field, fieldType, (VariableLiteral) term, parserArguments, serializerArguments, serialize, suppressPointerAccess);
        } else {
            throw new FreemarkerException("Unsupported Literal type " + term.getClass().getName());
        }
    }

    private String toVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess) {
        return toVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, false);
    }

    private String toVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, boolean isChild) {
        Tracer tracer = pythonTracerStart("toVariableExpression(serialize=" + serialize + ")");
        String variableLiteralName = variableLiteral.getName();
        boolean isEnumTypeReference = typeReference != null && typeReference.isEnumTypeReference();
        if ("lengthInBytes".equals(variableLiteralName)) {
            return toLengthInBytesVariableExpression(typeReference, serialize, tracer);
        } else if ("lengthInBits".equals(variableLiteralName)) {
            return toLengthInBitsVariableExpression(typeReference, serialize, tracer);
        } else if ("_value".equals(variableLiteralName)) {
            return toValueVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        }
        if ("_lastItem".equals(variableLiteralName)) {
            return toLastItemVariableExpression(typeReference, serialize, tracer);
        }
        if ("length".equals(variableLiteral.getChild().map(VariableLiteral::getName).orElse(""))) {
            return toLengthVariableExpression(field, variableLiteral, serialize, tracer);
        }
        // If this literal references an Enum type, then we have to output it differently.
        else if (getTypeDefinitions().get(variableLiteralName) instanceof EnumTypeDefinition) {
            return toEnumVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, suppressPointerAccess, tracer);
        }
        // If we are accessing enum constants, these also need to be output differently.
        else if (thisType.isComplexTypeDefinition()
            && thisType.asComplexTypeDefinition()
            .orElseThrow(IllegalAccessError::new)
            .getPropertyFieldByName(variableLiteralName)
            .filter(EnumField.class::isInstance)
            .isPresent()
            && (variableLiteral.getChild().isPresent())
        ) {
            return toConstantVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, suppressPointerAccess, tracer);
        }
        // If we are accessing optional fields, (we might need to use pointer-access).
        else if (!serialize && thisType.isComplexTypeDefinition()
            && thisType.asComplexTypeDefinition()
            .orElseThrow(IllegalStateException::new)
            .getPropertyFieldByName(variableLiteralName)
            .filter(OptionalField.class::isInstance)
            .isPresent()
        ) {
            tracer = tracer.dive("non serialize optional fields");
            return toOptionalVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, suppressPointerAccess, tracer);
        }
        // If we are accessing optional fields, (we might need to use pointer-access).
        else if (thisType.isComplexTypeDefinition()
            && thisType.asComplexTypeDefinition()
            .orElseThrow(IllegalStateException::new)
            .getPropertyFieldByName(variableLiteralName)
            .filter(OptionalField.class::isInstance)
            .isPresent()
        ) {
            tracer = tracer.dive("optional fields");
            OptionalField optionalField = thisType.asComplexTypeDefinition().orElseThrow().getPropertyFieldByName(variableLiteralName).orElseThrow().asOptionalField().orElseThrow();
            return tracer + "(" + (suppressPointerAccess || optionalField.getType().isComplexTypeReference() ? "" : "*") + "self." + camelCaseToSnakeCase(variableLiteral.getName()) + ")" +
                variableLiteral.getChild().map(child -> "." + camelCaseToSnakeCase(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, serialize, suppressPointerAccess, true))).orElse("");
        }
        // If we are accessing implicit fields, we need to rely on local variable instead.
        //else if (isVariableLiteralImplicitField(vl)) {
        //    tracer = tracer.dive("implicit");
        //    return tracer + (serialize ? vl.getName() : vl.getName()) + ((vl.getChild() != null) ?
        //        "." + capitalize(toVariableExpression(typeReference, vl.getChild(), parserArguments, serializerArguments, false, suppressPointerAccess)) : "");
        //}
        // If we are accessing implicit fields, we need to rely on a local variable instead.

        // CAST expressions are special as we need to add a ".class" to the second parameter in Java.
        else if ("CAST".equals(variableLiteralName)) {
            return toCastVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        } else if ("STATIC_CALL".equals(variableLiteralName)) {
            return toStaticCallVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        } else if (!isEnumTypeReference && "COUNT".equals(variableLiteralName)) {
            return toCountVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        } else if (!isEnumTypeReference && "ARRAY_SIZE_IN_BYTES".equals(variableLiteralName)) {
            return toArraySizeInBytesVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, suppressPointerAccess, tracer);
        } else if ("CEIL".equals(variableLiteralName)) {
            return toCeilVariableExpression(field, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        } else if ("STR_LEN".equals(variableLiteralName)) {
            return toStrLenVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        }
        // All uppercase names are not fields, but utility methods.
        // TODO: It seems we also run into this, in case of using enum constants in type-switches.
        else if (variableLiteralName.equals(variableLiteralName.toUpperCase())) {
            tracer = tracer.dive("utility");
            return toUppercaseVariableExpression(field, typeReference, variableLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer);
        }
        // If the current property references a discriminator value, we have to serialize it differently.
        else if (thisType.isComplexTypeDefinition() && thisType.asComplexTypeDefinition()
            .orElseThrow(IllegalStateException::new)
            .getPropertyFieldFromThisOrParentByName(variableLiteralName)
            .filter(DiscriminatorField.class::isInstance)
            .isPresent()) {
            tracer = tracer.dive("discriminator value");
            // TODO: Should this return something?
        }
        // If the variable has a child element and we're able to find a type for this ... get the type.
        else if ((variableLiteral.getChild().isPresent()) && ((ComplexTypeDefinition) thisType).getTypeReferenceForProperty(variableLiteralName).isPresent()) {
            tracer = tracer.dive("child element");
            final Optional<NonSimpleTypeReference> typeReferenceForProperty = ((ComplexTypeDefinition) thisType).getTypeReferenceForProperty(variableLiteralName)
                .flatMap(TypeReferenceConversions::asNonSimpleTypeReference);
            if (typeReferenceForProperty.isPresent()) {
                tracer = tracer.dive("complex");
                final NonSimpleTypeReference nonSimpleTypeReference = typeReferenceForProperty.get();
                TypeDefinition typeDefinition = nonSimpleTypeReference.getTypeDefinition();
                if (typeDefinition instanceof ComplexTypeDefinition) {
                    tracer = tracer.dive("complex");
                    ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) typeDefinition;
                    String childProperty = variableLiteral.getChild()
                        .orElseThrow(() -> new FreemarkerException("complex needs a child"))
                        .getName();
                    final Optional<Field> matchingDiscriminatorField = complexTypeDefinition.getFields().stream()
                        .filter(curField -> (curField instanceof DiscriminatorField) && ((DiscriminatorField) curField).getName().equals(childProperty))
                        .findFirst();
                    if (matchingDiscriminatorField.isPresent()) {
                        return tracer + "Cast" + getLanguageTypeNameForTypeReference(nonSimpleTypeReference) + "(" + variableLiteralName + ").get_" + camelCaseToSnakeCase(childProperty) + "()";
                    }
                    // TODO: is this really meant to fall through?
                    tracer = tracer.dive("we fell through the complex complex");
                } else if (typeDefinition instanceof EnumTypeDefinition) {
                    tracer = tracer.dive("enum");
                    String variableAccess = variableLiteralName;
                    if (isChild) {
                        variableAccess = "" + camelCaseToSnakeCase(variableLiteralName);
                    }
                    return tracer + (serialize ? "self." + camelCaseToSnakeCase(variableLiteralName) + "" : variableAccess) +
                        "." + camelCaseToSnakeCase(variableLiteral.getChild().orElseThrow(() -> new FreemarkerException("enum needs a child")).getName()) + "()";
                }
            }
            // TODO: is this really meant to fall through?
            tracer = tracer.dive("we fell through the child complete");
        } else if (isVariableLiteralImplicitField(variableLiteral)) {
            tracer = tracer.dive("implicit");
            if (serialize) {
                tracer = tracer.dive("serialize");
                final ImplicitField referencedImplicitField = getReferencedImplicitField(variableLiteral);
                return tracer + toSerializationExpression(referencedImplicitField, referencedImplicitField.getType(), getReferencedImplicitField(variableLiteral).getSerializeExpression(), serializerArguments);
            } else {
                return tracer + camelCaseToSnakeCase(variableLiteralName);
                //return toParseExpression(getReferencedImplicitField(vl), getReferencedImplicitField(vl).getSerializeExpression(), serializerArguments);
            }
        }

        // This is a special case for DataIo string types, which need to access the stringLength
        if ((serializerArguments != null) && serializerArguments.stream()
            .anyMatch(argument -> argument.getName().equals(variableLiteralName)) && "stringLength".equals(variableLiteralName)) {
            tracer = tracer.dive("serialization argument");
            return tracer + camelCaseToSnakeCase(variableLiteralName) +
                variableLiteral.getChild()
                    .map(child -> "." + camelCaseToSnakeCase(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, serialize, suppressPointerAccess, true)))
                    .orElse("");
        } else if ((serializerArguments != null) && serializerArguments.stream()
            .anyMatch(argument -> argument.getName().equals(variableLiteralName))) {
            tracer = tracer.dive("serialization argument");
            return tracer + "self." + camelCaseToSnakeCase(variableLiteralName) +
                variableLiteral.getChild()
                    .map(child -> "." + camelCaseToSnakeCase(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, serialize, suppressPointerAccess, true)))
                    .orElse("");
        }
        if ((parserArguments != null) && parserArguments.stream()
            .anyMatch(argument -> argument.getName().equals(variableLiteralName))) {
            tracer = tracer.dive("parser argument");
            return tracer + camelCaseToSnakeCase(variableLiteralName) +
                variableLiteral.getChild()
                    .map(child -> "." + camelCaseToSnakeCase(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, serialize, suppressPointerAccess, true)))
                    .orElse("");
        }
        String indexCall = "";
        if (variableLiteral.getIndex().isPresent()) {
            tracer = tracer.dive("indexCall");
            // We have a index call
            indexCall = "[" + variableLiteral.getIndex().orElseThrow() + "]";
        }
        tracer = tracer.dive("else");
        Tracer tracer2 = tracer;
        String variableAccess = variableLiteralName;
        if (isChild) {
            variableAccess = "get_" + camelCaseToSnakeCase(variableAccess) + "()";
        }
        return tracer + (serialize ? "self." + camelCaseToSnakeCase(variableLiteralName) + "" : camelCaseToSnakeCase(variableAccess)) + indexCall +
            variableLiteral.getChild()
                .map(child -> "." + camelCaseToSnakeCase(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true)))
                .orElse("");
    }

    private String toUppercaseVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("toUppercaseVariableExpression");
        StringBuilder sb = new StringBuilder(variableLiteral.getName().toUpperCase());
        if (variableLiteral.getArgs().isPresent()) {
            sb.append("(");
            boolean firstArg = true;
            for (Term arg : variableLiteral.getArgs().get()) {
                if (!firstArg) {
                    sb.append(", ");
                }
                sb.append(toExpression(field, typeReference, arg, parserArguments, serializerArguments, serialize, suppressPointerAccess));
                firstArg = false;
            }
            sb.append(")");
        }
        if (variableLiteral.getIndex().isPresent()) {
            sb.append("[").append(variableLiteral.getIndex().orElseThrow()).append("]");
        }
        return tracer + sb.toString() + variableLiteral.getChild()
            .map(child -> "." + toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true))
            .orElse("");
    }

    private String toCeilVariableExpression(Field field, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("ceil");
        Term va = variableLiteral.getArgs().orElseThrow(() -> new FreemarkerException("CEIL needs at least one arg"))
            .stream().findFirst().orElseThrow(IllegalStateException::new);
        // The Ceil function expects 64 bit floating point values.
        TypeReference tr = new DefaultFloatTypeReference(SimpleTypeReference.SimpleBaseType.FLOAT, 64);
        emitRequiredImport("from math import ceil");
        return tracer + "ceil(" + toExpression(field, tr, va, parserArguments, serializerArguments, serialize, suppressPointerAccess) + ")";
    }

    private String toArraySizeInBytesVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("array size in bytes");
        VariableLiteral va = variableLiteral.getArgs()
            .orElseThrow(() -> new FreemarkerException("ARRAY_SIZE_IN_BYTES needs at least one arg"))
            .stream().findFirst().orElseThrow(IllegalStateException::new)
            .asLiteral()
            .orElseThrow(() -> new FreemarkerException("ARRAY_SIZE_IN_BYTES needs a literal"))
            .asVariableLiteral()
            .orElseThrow(() -> new FreemarkerException("ARRAY_SIZE_IN_BYTES needs a variable literal"));
        // "io" and "m" are always available in every parser.
        boolean isSerializerArg = "read_buffer".equals(va.getName()) || "write_buffer".equals(va.getName()) || "self".equals(va.getName()) || "element".equals(va.getName());
        if (!isSerializerArg && serializerArguments != null) {
            for (Argument serializerArgument : serializerArguments) {
                if (serializerArgument.getName().equals(va.getName())) {
                    isSerializerArg = true;
                    break;
                }
            }
        }
        StringBuilder sb = new StringBuilder();
        if (isSerializerArg) {
            sb.append(va.getName()).append(va.getChild().map(child -> "." + toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, true, suppressPointerAccess, true)).orElse(""));
        } else {
            sb.append(toVariableExpression(field, typeReference, va, parserArguments, serializerArguments, true, suppressPointerAccess));
        }
        emitRequiredImport("from sys import getsizeof");
        return tracer + getCastExpressionForTypeReference(typeReference) + "(getsizeof(" + sb + "))";
    }

    private String toCountVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("count");
        VariableLiteral countLiteral = variableLiteral.getArgs()
            .orElseThrow(() -> new FreemarkerException("Count needs at least one arg"))
            .get(0)
            .asLiteral()
            .orElseThrow(() -> new FreemarkerException("Count needs a literal"))
            .asVariableLiteral()
            .orElseThrow(() -> new FreemarkerException("Count needs a variable literal"));
        return tracer + (typeReference instanceof SimpleTypeReference ? getCastExpressionForTypeReference(typeReference) : "") + "(len(" +
            toVariableExpression(field, typeReference, countLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess) +
            "))";
    }

    private String toStrLenVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("str-len");
        VariableLiteral countLiteral = variableLiteral.getArgs()
            .orElseThrow(() -> new FreemarkerException("Str-len needs at least one arg"))
            .get(0)
            .asLiteral()
            .orElseThrow(() -> new FreemarkerException("Str-len needs a literal"))
            .asVariableLiteral()
            .orElseThrow(() -> new FreemarkerException("Str-len needs a variable literal"));
        return tracer + (typeReference instanceof SimpleTypeReference ? getCastExpressionForTypeReference(typeReference) : "") + "(len(" +
            toVariableExpression(field, typeReference, countLiteral, parserArguments, serializerArguments, serialize, suppressPointerAccess) +
            "))";
    }

    private String toStaticCallVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("STATIC_CALL");
        StringBuilder sb = new StringBuilder();
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new FreemarkerException("A STATIC_CALL expression needs arguments"));
        if (arguments.size() < 1) {
            throw new FreemarkerException("A STATIC_CALL expression expects at least one argument.");
        }
        // Get the class and method name
        String staticCall = arguments.get(0).asLiteral()
            .orElseThrow(() -> new FreemarkerException("First argument should be a literal"))
            .asStringLiteral()
            .orElseThrow(() -> new FreemarkerException("Expecting the first argument of a 'STATIC_CALL' to be a StringLiteral")).
            getValue();
        sb.append(camelCaseToSnakeCase(staticCall)).append("(");
        for (int i = 1; i < arguments.size(); i++) {
            Term arg = arguments.get(i);
            if (i > 1) {
                sb.append(", ");
            }
            if (arg instanceof UnaryTerm) {
                arg = ((UnaryTerm) arg).getA();
            }
            if (arg instanceof VariableLiteral) {
                tracer = tracer.dive("VariableLiteral nr." + i);
                VariableLiteral va = (VariableLiteral) arg;
                // "io" is the default name of the reader argument which is always available.
                boolean isParserArg = "readBuffer".equals(va.getName()) || "writeBuffer".equals(va.getName()) || ((thisType instanceof DataIoTypeDefinition) && "_value".equals(va.getName()));
                boolean isBufferArg = "readBuffer".equals(va.getName()) || "writeBuffer".equals(va.getName());
                boolean isTypeArg = "_type".equals(va.getName());
                if (!isParserArg && !isTypeArg && parserArguments != null) {
                    for (Argument parserArgument : parserArguments) {
                        if (parserArgument.getName().equals(va.getName())) {
                            isParserArg = true;
                            break;
                        }
                    }
                }
                if (isBufferArg) {
                    sb.append(toVariableExpression(field, typeReference, va, parserArguments, serializerArguments, false, suppressPointerAccess));
                } else if (isParserArg) {
                    tracer = tracer.dive("isParserArg");
                    if (va.getName().equals("_value")) {
                        tracer = tracer.dive("is _value");
                        sb.append(va.getName().substring(1)).append(va.getChild().map(child -> "." + toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, false)).orElse(""));
                    } else {
                        sb.append(camelCaseToSnakeCase(va.getName())).append((va.getChild().isPresent()) ?
                            ".get_" + camelCaseToSnakeCase(toVariableExpression(field, typeReference, va.getChild().orElseThrow(IllegalStateException::new), parserArguments, serializerArguments, false, suppressPointerAccess)) + "()" : "");
                    }
                }
                // We have to manually evaluate the type information at code-generation time.
                else if (isTypeArg) {
                    String part = va.getChild().map(VariableLiteral::getName).orElse("");
                    switch (part) {
                        case "name":
//                                sb.append("\"").append(field.getTypeName()).append("\"");
                            break;
                        case "length":
                            sb.append("\"").append(((SimpleTypeReference) typeReference).getSizeInBits()).append("\"");
                            break;
                        case "encoding":
                            final Term encodingTerm = field.getEncoding().orElse(new DefaultStringLiteral("UTF-8"));
                            if (!(encodingTerm instanceof StringLiteral)) {
                                throw new FreemarkerException("Encoding must be a quoted string value");
                            }
                            String encoding = ((StringLiteral) encodingTerm).getValue();
                            sb.append("\"").append(encoding).append("\"");
                            break;
                    }
                } else {
                    sb.append(toVariableExpression(field, typeReference, va, parserArguments, serializerArguments, serialize, suppressPointerAccess));
                }
            } else if (arg instanceof StringLiteral) {
                tracer = tracer.dive("StringLiteral");
                sb.append(((StringLiteral) arg).getValue());
            } else if (arg instanceof BooleanLiteral) {
                tracer = tracer.dive("BooleanLiteral");
                sb.append(((BooleanLiteral) arg).getValue());
            } else if (arg instanceof NumericLiteral) {
                tracer = tracer.dive("NumericLiteral");
                sb.append(((NumericLiteral) arg).getNumber());
            } else if (arg instanceof BinaryTerm) {
                tracer = tracer.dive("BinaryTerm");
                sb.append(toBinaryTermExpression(field, typeReference, (BinaryTerm) arg, parserArguments, serializerArguments, serialize, tracer));
            } else {
                throw new FreemarkerException(arg.getClass().getName());
            }
        }
        sb.append(")");
        return tracer + sb.toString();
    }

    private String toCastVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("CAST");
        List<Term> arguments = variableLiteral.getArgs().orElseThrow(() -> new FreemarkerException("A Cast expression needs arguments"));
        if (arguments.size() != 2) {
            throw new FreemarkerException("A CAST expression expects exactly two arguments.");
        }
        VariableLiteral firstArgument = arguments.get(0).asLiteral()
            .orElseThrow(() -> new FreemarkerException("First argument should be a literal"))
            .asVariableLiteral()
            .orElseThrow(() -> new FreemarkerException("First argument should be a Variable literal"));
        StringLiteral typeLiteral = arguments.get(1).asLiteral()
            .orElseThrow(() -> new FreemarkerException("Second argument should be a String literal"))
            .asStringLiteral()
            .orElseThrow(() -> new FreemarkerException("Second argument should be a String literal"));
        final TypeDefinition typeDefinition = getTypeDefinitions().get(typeLiteral.getValue());
        StringBuilder sb = new StringBuilder();
        if (typeDefinition.isComplexTypeDefinition()) {
            sb.append("Cast");
        }
        sb.append(typeLiteral.getValue());
        sb.append("(").append(toVariableExpression(field, typeReference, firstArgument, parserArguments, serializerArguments, serialize, suppressPointerAccess)).append(")");
        return tracer + sb.toString() + variableLiteral.getChild().map(child -> "." + camelCaseToSnakeCase(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true))).orElse("");
    }

    private String toOptionalVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("optional fields");
        return tracer + "(" + (suppressPointerAccess || (typeReference != null && typeReference.isComplexTypeReference()) ? "" : "*") + variableLiteral.getName() + ")" +
            variableLiteral.getChild().map(child -> "." + camelCaseToSnakeCase(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true))).orElse("");
    }

    private String toConstantVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("enum constant");
        VariableLiteral child = variableLiteral.getChild().orElseThrow(IllegalStateException::new);
        return tracer + variableLiteral.getName() + "." + camelCaseToSnakeCase(child.getName()) + "()" +
            child.getChild().map(childChild -> "." + toVariableExpression(field, typeReference, childChild, parserArguments, serializerArguments, false, suppressPointerAccess, true)).orElse("");
    }

    private String toEnumVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("enum");
        VariableLiteral child = variableLiteral.getChild().orElseThrow(() -> new FreemarkerException("Enum should have a child"));
        return tracer + variableLiteral.getName() + "_" + child.getName() +
            child.getChild().map(childChild -> "." + toVariableExpression(field, typeReference, childChild, parserArguments, serializerArguments, false, suppressPointerAccess, true)).orElse("");
    }

    private String toLastItemVariableExpression(TypeReference typeReference, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("lastItem");
        return tracer + "utils.GetLastItemFromContext(ctx)";
    }

    private String toLengthVariableExpression(Field field, VariableLiteral variableLiteral, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("length");
        return tracer + (serialize ? ("len(self." + camelCaseToSnakeCase(variableLiteral.getName()) + ")") : ("(" + variableLiteral.getName() + ")"));
    }

    private String toValueVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        final Tracer tracer2 = tracer.dive("_value");
        return variableLiteral.getChild()
            .map(child -> tracer2.dive("withChild") + "self." + toUppercaseVariableExpression(field, typeReference, child, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer2))
            .orElse(tracer2 + "value");
    }

    private String toLengthInBitsVariableExpression(TypeReference typeReference, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("lengthInBits");
        return tracer + (serialize ? getCastExpressionForTypeReference(typeReference) + "(" : "") + "length_in_bits" + (serialize ? "())" : "()");
    }

    private String toLengthInBytesVariableExpression(TypeReference typeReference, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("lengthInBytes");
        return tracer + (serialize ? getCastExpressionForTypeReference(typeReference) + "(" : "") + "length_in_bytes" + (serialize ? "())" : "()");
    }

    public String getSizeInBits(ComplexTypeDefinition complexTypeDefinition, List<Argument> parserArguments) {
        int sizeInBits = 0;
        StringBuilder sb = new StringBuilder();
        for (Field field : complexTypeDefinition.getFields()) {
            if (field instanceof ArrayField) {
                ArrayField arrayField = (ArrayField) field;
                final SimpleTypeReference type = (SimpleTypeReference) arrayField.getType();
                switch (arrayField.getLoopType()) {
                    case COUNT:
                        sb.append("(").append(toTypedSerializationExpression(field, type, arrayField.getLoopExpression(), parserArguments)).append(" * ").append(type.getSizeInBits()).append(") + ");
                        break;
                    case LENGTH:
                        sb.append("(").append(toTypedSerializationExpression(field, type, arrayField.getLoopExpression(), parserArguments)).append(" * 8) + ");
                        break;
                    case TERMINATED:
                        // No terminated.
                        break;
                }
            } else if (field instanceof TypedField) {
                TypedField typedField = (TypedField) field;
                final TypeReference type = typedField.getType();
                if (field instanceof ManualField) {
                    ManualField manualField = (ManualField) field;
                    sb.append("(").append(toSerializationExpression(manualField, getIntTypeReference(), manualField.getLengthExpression(), parserArguments)).append(") + ");
                } else if (type instanceof SimpleTypeReference) {
                    SimpleTypeReference simpleTypeReference = (SimpleTypeReference) type;
                    sizeInBits += simpleTypeReference.getSizeInBits();
                } else {
                    throw new IllegalStateException("No ComplexTypeReference supported");
                }
            }
        }
        return sb.toString() + sizeInBits;
    }

    public String escapeValue(TypeReference typeReference, String valueString) {
        if (valueString == null) {
            return null;
        }
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case UINT:
                case INT:
                    // If it's a one character string and is numeric, output it as char.
                    if (!NumberUtils.isParsable(valueString) && (valueString.length() == 1)) {
                        return "'" + valueString + "'";
                    }
                    break;
                case STRING:
                case VSTRING:
                    return "\"" + valueString + "\"";
            }
        }
        return valueString;
    }

    public String escapeEnumValue(TypeReference typeReference, String valueString) {
        // Currently the only case in which here complex type references are used are when referencing enum constants.
        if (typeReference != null && typeReference.isNonSimpleTypeReference()) {
            // C doesn't like NULL values for enums, so we have to return something else (we'll treat -1 as NULL)
            if ("null".equals(valueString) || valueString == null) {
                return "0";
            }
            if (valueString.contains(".")) {
                String typeName = valueString.substring(0, valueString.indexOf('.'));
                String constantName = valueString.substring(valueString.indexOf('.') + 1);
                return typeName + "_" + constantName;
            }
            return valueString;
        } else {
            return escapeValue(typeReference, valueString);
        }
    }

    public Collection<EnumValue> getUniqueEnumValues(List<EnumValue> enumValues) {
        Map<String, EnumValue> filteredEnumValues = new TreeMap<>();
        for (EnumValue enumValue : enumValues) {
            if (!filteredEnumValues.containsKey(enumValue.getValue())) {
                filteredEnumValues.put(enumValue.getValue(), enumValue);
            }
        }
        return filteredEnumValues.values();
    }

    public List<DiscriminatedComplexTypeDefinition> getUniqueSwitchCases(List<DiscriminatedComplexTypeDefinition> allSwitchCases) {
        Map<String, DiscriminatedComplexTypeDefinition> switchCases = new LinkedHashMap<>();
        for (DiscriminatedComplexTypeDefinition switchCase : allSwitchCases) {
            if (!switchCases.containsKey(switchCase.getName())) {
                switchCases.put(switchCase.getName(), switchCase);
            }
        }
        return new ArrayList<>(switchCases.values());
    }

    public void emitRequiredImport(String requiredImport) {
        LOGGER.debug("emitting import '\"{}\"'", requiredImport);
        requiredImports.add(requiredImport);
    }

    public void emitRequiredImport(String alias, String requiredImport) {
        LOGGER.debug("emitting import '{} \"{}'\"", alias, requiredImport);
        requiredImports.add(alias + ' ' + '"' + requiredImport + '"');
    }

    public Set<String> getRequiredImports() {
        return requiredImports;
    }

    public void emitDataIoRequiredImport(String requiredImport) {
        LOGGER.debug("emitting io import '\"{}\"'", requiredImport);
        requiredImportsForDataIo.add(requiredImport);
    }

    public void emitDataIoRequiredImport(String alias, String requiredImport) {
        LOGGER.debug("emitting data io import '{} \"{}'\"", alias, requiredImport);
        requiredImportsForDataIo.add(alias + ' ' + '"' + requiredImport + '"');
    }

    public Set<String> getRequiredImportsForDataIo() {
        return requiredImportsForDataIo;
    }

    public String getVariableName(Field field) {
        if (!(field instanceof NamedField)) {
            return "_";
        }
        NamedField namedField = (NamedField) field;

        String name = null;
        for (Field curField : ((ComplexTypeDefinition) thisType).getFields()) {
            if (curField == field) {
                name = namedField.getName();
            } else if (name != null) {
                if (curField instanceof ArrayField) {
                    ArrayField arrayField = (ArrayField) curField;
                    if (arrayField.getLoopExpression().contains(name)) {
                        return name;
                    }
                } else if (curField instanceof ChecksumField) {
                    ChecksumField checksumField = (ChecksumField) curField;
                    if (checksumField.getChecksumExpression().contains(name)) {
                        return name;
                    }
                } else if (curField instanceof ImplicitField) {
                    ImplicitField implicitField = (ImplicitField) curField;
                    if (implicitField.getSerializeExpression().contains(name)) {
                        return name;
                    }
                } else if (curField instanceof ManualArrayField) {
                    ManualArrayField manualArrayField = (ManualArrayField) curField;
                    if (manualArrayField.getLengthExpression().contains(name)) {
                        return name;
                    }
                    if (manualArrayField.getLoopExpression().contains(name)) {
                        return name;
                    }
                    if (manualArrayField.getParseExpression().contains(name)) {
                        return name;
                    }
                    if (manualArrayField.getSerializeExpression().contains(name)) {
                        return name;
                    }
                } else if (curField instanceof ManualField) {
                    ManualField manualField = (ManualField) curField;
                    if (manualField.getLengthExpression().contains(name)) {
                        return name;
                    }
                    if (manualField.getParseExpression().contains(name)) {
                        return name;
                    }
                    if (manualField.getSerializeExpression().contains(name)) {
                        return name;
                    }
                } else if (curField instanceof OptionalField) {
                    OptionalField optionalField = (OptionalField) curField;
                    if (optionalField.getConditionExpression().isPresent() && optionalField.getConditionExpression().orElseThrow(IllegalStateException::new).contains(name)) {
                        return name;
                    }
                } else if (curField instanceof SwitchField) {
                    SwitchField switchField = (SwitchField) curField;
                    for (Term discriminatorExpression : switchField.getDiscriminatorExpressions()) {
                        if (discriminatorExpression.contains(name)) {
                            return name;
                        }
                    }
                    for (DiscriminatedComplexTypeDefinition curCase : switchField.getCases()) {
                        for (Argument parserArgument : curCase.getParserArguments().orElse(Collections.emptyList())) {
                            if (parserArgument.getName().equals(name)) {
                                return name;
                            }
                        }
                    }
                } else if (curField instanceof VirtualField) {
                    VirtualField virtualField = (VirtualField) curField;
                    if (virtualField.getValueExpression().contains(name)) {
                        return name;
                    }
                }
                List<Term> params = field.asTypedField()
                    .map(typedField -> typedField.getType().asNonSimpleTypeReference()
                        .map(NonSimpleTypeReference::getParams)
                        .map(terms -> terms.orElse(Collections.emptyList()))
                        .orElse(Collections.emptyList())
                    )
                    .orElse(Collections.emptyList());
                for (Term param : params) {
                    if (param.contains(name)) {
                        return name;
                    }
                }
            }
        }

        return "_";
    }

    public boolean needsVariable(Field field, String variableName, boolean serialization) {
        if (!serialization) {
            if (field instanceof ArrayField) {
                ArrayField arrayField = (ArrayField) field;
                if (arrayField.getLoopExpression().contains(variableName)) {
                    return true;
                }
            }
        }
        if (field instanceof VirtualField) {
            VirtualField virtualField = (VirtualField) field;
            if (virtualField.getValueExpression().contains(variableName)) {
                return true;
            }
        }
        if (field instanceof PaddingField) {
            PaddingField paddingField = (PaddingField) field;
            if (paddingField.getPaddingCondition().contains(variableName)) {
                return true;
            }
            if (paddingField.getPaddingValue().contains(variableName)) {
                return true;
            }
        }
        return field.asTypedField()
            .map(typedField -> typedField.getType().asNonSimpleTypeReference()
                .map(nonSimpleTypeReference -> nonSimpleTypeReference.getParams()
                    .map(params -> params.stream()
                        .anyMatch(param -> param.contains(variableName))
                    )
                    .orElse(false)
                )
                .orElse(false)
            )
            .orElse(false);
    }

    /**
     * Right now only the ARRAY_SIZE_IN_BYTES requires helpers to be generated.
     * Also right now only the Modbus protocol requires this and here the referenced
     * properties are all also members of the current complex type,
     * so we'll simplify things here for now.
     *
     * @param functionName name of the
     * @return something
     */
    public Map<String, String> requiresHelperFunctions(String functionName) {
        Map<String, String> result = new HashMap<>();
        boolean usesFunction = false;
        // As the ARRAY_SIZE_IN_BYTES only applies to ArrayFields, search for these
        for (Field curField : ((ComplexTypeDefinition) thisType).getFields()) {
            if (curField instanceof ArrayField) {
                ArrayField arrayField = (ArrayField) curField;
                if (arrayField.getLoopExpression().contains(functionName)) {
                    usesFunction = true;
                }
                result.put(arrayField.getName(), getLanguageTypeNameForField(arrayField));
            } else if (curField instanceof ImplicitField) {
                ImplicitField implicitField = (ImplicitField) curField;
                if (implicitField.getSerializeExpression().contains(functionName)) {
                    usesFunction = true;
                }
            }
        }
        if (!usesFunction) {
            return Collections.emptyMap();
        }
        return result;
    }

    public boolean requiresCurPos() {
        if (thisType instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) this.thisType;
            for (Field curField : complexTypeDefinition.getFields()) {
                if (requiresVariable(curField, "curPos")) {
                    return true;
                }
            }
        }
        return false;
    }

    public boolean requiresStartPos() {
        if (thisType instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) this.thisType;
            for (Field curField : complexTypeDefinition.getFields()) {
                if (requiresVariable(curField, "startPos")) {
                    return true;
                }
            }
        }
        return false;
    }

    public boolean requiresVariable(Field curField, String variable) {
        if (curField.isArrayField()) {
            ArrayField arrayField = (ArrayField) curField;
            if (arrayField.getLoopExpression().contains(variable)) {
                return true;
            }
        } else if (curField.isOptionalField()) {
            OptionalField optionalField = (OptionalField) curField;
            if (optionalField.getConditionExpression().isPresent() && optionalField.getConditionExpression().orElseThrow(IllegalStateException::new).contains(variable)) {
                return true;
            }
        }
        return curField.asTypedField()
            .map(typedField -> typedField.getType().asNonSimpleTypeReference()
                .map(nonSimpleTypeReference -> nonSimpleTypeReference.getParams()
                    .map(params -> params.stream()
                        .anyMatch(param -> param.contains(variable))
                    )
                    .orElse(false)
                )
                .orElse(false)
            )
            .orElse(false);
    }

    public Term findTerm(Term baseTerm, String name) {
        if (baseTerm instanceof VariableLiteral) {
            VariableLiteral variableLiteral = (VariableLiteral) baseTerm;
            if (variableLiteral.getName().equals(name)) {
                return variableLiteral;
            }
            if (variableLiteral.getChild().isPresent()) {
                Term found = findTerm(variableLiteral.getChild().get(), name);
                if (found != null) {
                    return found;
                }
            }
            for (Term arg : variableLiteral.getArgs().orElse(Collections.emptyList())) {
                Term found = findTerm(arg, name);
                if (found != null) {
                    return found;
                }
            }
        } else if (baseTerm instanceof UnaryTerm) {
            UnaryTerm unaryTerm = (UnaryTerm) baseTerm;
            return findTerm(unaryTerm.getA(), name);
        } else if (baseTerm instanceof BinaryTerm) {
            BinaryTerm binaryTerm = (BinaryTerm) baseTerm;
            Term found = findTerm(binaryTerm.getA(), name);
            if (found != null) {
                return found;
            }
            found = findTerm(binaryTerm.getB(), name);
            return found;
        } else if (baseTerm instanceof TernaryTerm) {
            TernaryTerm ternaryTerm = (TernaryTerm) baseTerm;
            Term found = findTerm(ternaryTerm.getA(), name);
            if (found != null) {
                return found;
            }
            found = findTerm(ternaryTerm.getB(), name);
            if (found != null) {
                return found;
            }
            found = findTerm(ternaryTerm.getC(), name);
            return found;
        }
        return null;
    }

    public String getEnumExpression(String expression) {
        String enumName = expression.substring(0, expression.indexOf('.'));
        String enumConstant = expression.substring(expression.indexOf('.') + 1);
        return enumName + "_" + enumConstant;
    }

    public boolean needsReferenceForParserArgument(String propertyName, TypeReference argumentType) {
        return argumentType.asComplexTypeReference()
            .map(complexTypeReference -> {
                // Check if this is a local field.
                // FIXME: shouldn't this look onto the argumentType? this is nowhere used...
                return thisType.asComplexTypeDefinition()
                    .map(
                        complexTypeDefinition -> complexTypeDefinition.getPropertyFieldByName(propertyName)
                            .map(TypedField.class::cast)
                            .map(TypedField::getType)
                            .filter(NonSimpleTypeReference.class::isInstance)
                            .map(NonSimpleTypeReference.class::cast)
                            .map(NonSimpleTypeReference::getTypeDefinition)
                            .map(typeDefinition -> !(typeDefinition instanceof EnumTypeDefinition))
                            .orElse(false)
                    )
                    .orElse(false);
            })
            .orElse(false);
    }

    public String capitalize(String str) {
        Tracer dummyTracer = pythonTracerStart("");
        String extractedTrace = dummyTracer.extractTraces(str);
        String cleanedString = dummyTracer.removeTraces(str);
        return extractedTrace + StringUtils.capitalize(cleanedString);
    }


    /**
     * Converts a camel-case string to snake-case.
     *
     * @param camelCase camel-case string
     * @return snake-case string
     */
    public String camelCaseToSnakeCase(String camelCase) {
        StringBuilder snakeCase = new StringBuilder();
        final char[] chars = camelCase.toCharArray();
        for (int i = 0; i < chars.length; i++) {
            String lowerCaseChar = String.valueOf(chars[i]).toLowerCase();
            // If the previous letter is a lowercase letter and this one is uppercase, create a new snake-segment.
            if ((i > 0) && !Character.isUpperCase(chars[i - 1]) && Character.isUpperCase(chars[i])) {
                snakeCase.append('_').append(lowerCaseChar);
            } else if ((i < (chars.length - 2)) && Character.isUpperCase(chars[i]) && !Character.isUpperCase(chars[i + 1])) {
                snakeCase.append('_').append(lowerCaseChar);
            }
            // If this is uppercase and the previous one is too ... just make this letter lowercase.
            else if ((i > 0) && Character.isUpperCase(chars[i - 1]) && Character.isUpperCase(chars[i])) {
                snakeCase.append(lowerCaseChar);
            } else if (chars[i] == '-') {
                snakeCase.append("_");
            } else {
                snakeCase.append(lowerCaseChar);
            }
        }
        // If the first letter was a capital letter, the string will start with a "_".
        // In this case we cut that char off.
        if (snakeCase.indexOf("_") == 0) {
            return snakeCase.substring(1);
        }
        return snakeCase.toString();
    }

    public Tracer pythonTracerStart(String base) {
        return new Tracer(base) {
            protected String separator() {
                return "/";
            }

            protected String prefix() {
                return "\"\"\"";
            }

            protected String suffix() {
                return "\"\"\"";
            }
        };
    }

}