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

    /**
     * Returns a file name for a generated Python file.
     *
     * The file name is constructed by taking the sanitized protocol name and
     * the language flavor name and concatenating them with a dot in between.
     *
     * @param protocolName the name of the protocol
     * @param languageName the name of the language
     * @param languageFlavorName the name of the language flavor
     * @return the file name for the generated Python file
     */
    public String fileName(String protocolName, String languageName, String languageFlavorName) {
        return String.join("", protocolName.split("\\-")) + "." +
            String.join("", languageFlavorName.split("\\-"));
    }

    /**
     * Returns a sanitized version of the protocol name.
     *
     * The protocol name is used to construct the package name and is thus
     * subject to a set of constraints:
     *
     * - It must only contain valid Java identifier characters.
     * - It must not contain any periods.
     *
     * This method will replace all periods with forward slashes, convert to
     * lower case and remove any dashes.
     *
     * @return a sanitized version of the protocol name
     */
    public String getSanitizedPackageName() {
        String sanitizedName = getProtocolName().replaceAll("-", "");
        // replace periods with forward slashes
        sanitizedName = sanitizedName.replaceAll("\\.", "/");
        // convert to lower case
        sanitizedName = sanitizedName.toLowerCase();
        return sanitizedName;
    }

    /**
     * Returns a sanitized version of the protocol name.
     *
     * The protocol name is used to construct the package name and is thus
     * subject to a set of constraints:
     *
     * - It must only contain valid Java identifier characters.
     * - It must not contain any periods.
     *
     * This method takes the protocol name and removes any hyphens, and then
     * converts the resulting string to camel case.
     *
     * @return the sanitized protocol name
     */
    // TODO: check if protocol name can be enforced to only contain valid chars
    public String getSanitizedProtocolName() {
        String sanitizedName = getProtocolName().replaceAll("-", "");
        sanitizedName = CaseUtils.toCamelCase(sanitizedName, false, '.');
        return sanitizedName;
    }

    /**
     * Returns the package name for the generated code.
     *
     * The package name is composed of the sanitized protocol name, the language name
     * ("python"), and the sanitized language flavor name.
     *
     * @return the package name.
     */
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

    /**
     * Checks whether the given field is a complex field.
     * <p>
     * A complex field is a field which has a type which is not a simple type (like int or boolean).
     * A complex type is a type which is not a primitive type and is not a {@link ComplexTypeReference}.
     * Examples of complex types are {@link java.util.List}, {@link java.util.Map}, etc.
     * @param field the field to check
     * @return true if the field is a complex field, false otherwise
     */
    public boolean isComplex(Field field) {
        return field instanceof PropertyField && ((PropertyField) field).getType() instanceof NonSimpleTypeReference;
    }

    /**
     * Returns the language type name for a given type reference.
     * <p>
     * This method is used in the generated code to generate type names for fields. The returned type name
     * is the name of the type as it should be used in the generated code. If the type is a primitive
     * type, the returned type name is usually more specific and will indicate that a
     * com.plc4x.common.value.Value of the corresponding type should be used.
     * @param typeReference the type reference to get the language type name for
     * @return the language type name for the given type reference
     */
    @Override
    public String getLanguageTypeNameForTypeReference(TypeReference typeReference) {
        return getLanguageTypeNameForTypeReference(typeReference, null);
    }

    /**
     * Get the language name for a given field, but return a non-primitive type name if the type is a
     * primitive type.
     *
     * This method is used in the generated code to generate type names for fields. If the type is a
     * primitive type, the returned type name is not the same as the one returned by
     * {@link #getLanguageTypeNameForField(Field) getLanguageTypeNameForField}. Instead, the returned
     * type name is usually more specific and will indicate that a complex type is expected.
     *
     * @param field the field to get the language name for
     * @return the language name for the field
     */
    public String getNonPrimitiveLanguageTypeNameForField(TypedField field) {
        return getLanguageTypeNameForTypeReference(field.getType(), false);
    }

    public String getLanguageTypeNameForTypeReference(TypeReference typeReference, boolean allowPrimitives) {
        return getLanguageTypeNameForTypeReference(typeReference, null);
    }

    /**
     * Get the language name for a given type reference.
     * This method is used to generate code for type references
     * in the generated code.
     *
     * @param typeReference the type reference for which the language name should be generated
     * @param encoding the encoding for the type reference, if applicable
     * @return the language name for the given type reference
     */
    public String getLanguageTypeNameForTypeReference(TypeReference typeReference, String encoding) {
        if (typeReference == null) {
            // TODO: shouldn't this be an error case
            return "";
        }
        if (typeReference.isArrayTypeReference()) {
            // We can use the array type reference to generate code for the type of the array.
            final ArrayTypeReference arrayTypeReference = (ArrayTypeReference) typeReference;
            TypeReference elementTypeReference = arrayTypeReference.getElementTypeReference();
            emitRequiredImport("from typing import List");
            // We use List to represent an array in python, as the type system doesn't support arrays.
            // We also use the name of the element type as the type parameter.
            return "List[" + getLanguageTypeNameForTypeReference(elementTypeReference) + "]";
        }
        if (typeReference.isNonSimpleTypeReference()) {
            // We have a complex type, such as a data io type or a custom type.
            emitRequiredImport("from plc4py.protocols." + protocolName + "." + flavorName.replace("-", "") + "." + typeReference.asNonSimpleTypeReference().orElseThrow().getName() + " import " + typeReference.asNonSimpleTypeReference().orElseThrow().getName());
            // We import the name of the complex type, and then use the name of the complex type as the type name.
            return typeReference.asNonSimpleTypeReference().orElseThrow().getName();
        }
        if (typeReference instanceof ByteOrderTypeReference) {
            // Byte order is represented by a class in the protocol, so we import that class and use its name as the type name.
            return "binary.ByteOrder";
        }
        SimpleTypeReference simpleTypeReference = typeReference.asSimpleTypeReference().orElseThrow();
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                // Bit is represented by a boolean in python.
                return "bool";
            case BYTE:
                // Byte is represented by an integer in python.
                return "int";
            case UINT:
                // Unsigned integer is represented by an integer in python, with the appropriate size.
                IntegerTypeReference unsignedIntegerTypeReference = simpleTypeReference.asIntegerTypeReference().orElseThrow();
                int sizeInBits = unsignedIntegerTypeReference.getSizeInBits();
                if (sizeInBits <= 8) {
                    // Size <= 8 bit is represented by an int in python.
                    return "int";
                }
                if (sizeInBits <= 16) {
                    // Size <= 16 bit is represented by an int in python.
                    return "int";
                }
                if (sizeInBits <= 32) {
                    // Size <= 32 bit is represented by an int in python.
                    return "int";
                }
                if (sizeInBits <= 64) {
                    // Size <= 64 bit is represented by an int in python.
                    return "int";
                }
                // Size > 64 bit is not supported in python.
                throw new UnsupportedOperationException("Size > 64 bit is not supported in python");
            case INT:
                // Integer is represented by an integer in python, with the appropriate size.
                IntegerTypeReference integerTypeReference = simpleTypeReference.asIntegerTypeReference().orElseThrow();
                int sizeInIntegerBits = integerTypeReference.getSizeInBits();
                if (sizeInIntegerBits <= 8) {
                    // Size <= 8 bit is represented by an int in python.
                    return "int";
                }
                if (sizeInIntegerBits <= 16) {
                    // Size <= 16 bit is represented by an int in python.
                    return "int";
                }
                if (sizeInIntegerBits <= 32) {
                    // Size <= 32 bit is represented by an int in python.
                    return "int";
                }
                if (sizeInIntegerBits <= 64) {
                    // Size <= 64 bit is represented by an int in python.
                    return "int";
                }
                // Size > 64 bit is not supported in python.
                throw new UnsupportedOperationException("Size > 64 bit is not supported in python");
            case FLOAT:
            case UFLOAT:
                // Float is represented by a float in python, with the appropriate size.
                FloatTypeReference floatTypeReference = simpleTypeReference.asFloatTypeReference().orElseThrow();
                int sizeInFloatBits = floatTypeReference.getSizeInBits();
                if (sizeInFloatBits <= 32) {
                    // Size <= 32 bit is represented by a float in python.
                    return "float";
                }
                if (sizeInFloatBits <= 64) {
                    // Size <= 64 bit is represented by a float in python.
                    return "float";
                }
                // Size > 64 bit is not supported in python.
                throw new UnsupportedOperationException("Size > 64 bit is not supported in python");
            case STRING:
            case VSTRING:
                // String is represented by a str in python.
                return "str";
            case TIME:
            case DATE:
            case DATETIME:
                // Time is represented by a datetime in python.
                return "datetime";
            default:
                // For all other types, an error should be thrown.
                throw new FreemarkerException("Unsupported simple type");
        }
    }

    public String getReservedValue(ReservedField reservedField) {
        final String languageTypeName = getLanguageTypeNameForTypeReference(reservedField.getType());
        return languageTypeName + "(" + reservedField.getReferenceValue() + ")";
    }

    /**
     * Generate code for options of a field, such as encoding and byte order.
     * 
     * @param field The field for which the options should be generated
     * @param parserArguments The arguments of the parser, which may contain
     * information about the encoding and byte order.
     * @return A string containing the options, or an empty string if no options
     * are needed.
     */
    public String getFieldOptions(TypedField field, List<Argument> parserArguments) {
        StringBuilder sb = new StringBuilder();
        // If the field has an encoding, add it as an option.
        final Optional<Term> encodingOptional = field.getEncoding();
        if (encodingOptional.isPresent()) {
            final String encoding = toParseExpression(field, field.getType(), encodingOptional.get(), parserArguments);
            sb.append(", encoding='").append(encoding).append("'");
        }
        // If the field has a byte order, add it as an option.
        final Optional<Term> byteOrderOptional = field.getByteOrder();
        if (byteOrderOptional.isPresent()) {
            final String byteOrder = toParseExpression(field, field.getType(), byteOrderOptional.get(), parserArguments);
            // We need to import the ByteOrder class, so add an import statement.
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

    /**
     * Returns the name of the Python data reader method that should be used to read a value of the given type.
     *
     * <p>
     * This method is used to generate code for converting a type reference to a PLC4Py type name.
     *
     * @param simpleTypeReference the type reference to get the PLC4Py type name for
     * @return the name of the Python data reader method
     */
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
    /**
     * Returns the name of the Python data writer method that should be used to write a value of the given type.
     *
     * <p>
     * This method is used to generate code for converting a type reference to a PLC4Py type name.
     *
     * @param typeReference the type reference to get the PLC4Py type name for
     * @param fieldName the name of the field to write
     * @return the name of the Python data writer method
     */

    /**
     * Returns the name of the Python data writer method that should be used to write a value of the given type.
     *
     * <p>
     * This method is used to generate code for converting a type reference to a PLC4Py type name.
     *
     * @param typeReference the type reference to get the PLC4Py type name for
     * @param fieldName the name of the field to write
     * @return the name of the Python data writer method
     */
    public String getDataWriterCall(TypeReference typeReference, String fieldName) {
        if (typeReference.isSimpleTypeReference()) {
            // The type is a simple type, such as a byte or int.
            // We can directly use the methods provided by the simple type.
            return getDataWriterCall(typeReference.asSimpleTypeReference().orElseThrow());
        } else if (typeReference.isArrayTypeReference()) {
            // The type is an array, such as an array of bytes or an array of strings.
            // We need to recursively call getDataWriterCall to determine the name of the PLC4Py
            // type name for the element type of the array.
            return getDataWriterCall(typeReference.asArrayTypeReference().orElseThrow().getElementTypeReference(), fieldName);
        } else if (typeReference.isComplexTypeReference()) {
            // The type is a complex type, such as a struct.
            // We can directly use the write_serializable method.
            return "write_serializable";
        } else {
            throw new IllegalStateException("What is this type? " + typeReference);
        }
    }

    /**
     * Returns the name of the Python data writer method that should be used to write a value of the given enum type.
     *
     * <p>
     * This method is used to generate code for converting a type reference to a PLC4Py type name.
     *
     * @param typeReference the enum type reference to get the PLC4Py type name for
     * @param fieldName the name of the field to write
     * @param attributeName the name of the enum attribute to write, defaults to "value"
     * @return the PLC4Py type name or an empty string if the type reference is not supported
     */
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

    /**
     * Returns the name of the Python data writer method that should be used to write a value of the given type.
     *
     * <p>
     * This method is used to generate code for converting a type reference to a PLC4Py type name.
     *
     * @param simpleTypeReference the type reference to get the PLC4Py type name for
     *
     * @return the PLC4Py type name or an empty string if the type reference is not supported
     */
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
                // Vstring is serialized as string
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

    /**
     * Returns the name of the PLC4Py value type that corresponds to the given type reference.
     * 
     * <p>
     * This method is used to generate code for converting a type reference to a PLC4Py type name.
     * 
     * @param typeReference the type reference to get the PLC4Py type name for
     * 
     * @return the PLC4Py type name or an empty string if the type reference is not supported
     */
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
                // Bit values are represented as boolean in PLC4Py
                emitRequiredImport("from plc4py.spi.values.PlcValues import PlcBOOL");
                return "PlcBOOL";
            case BYTE:
                // Byte values are represented as signed integers in PLC4Py
                emitRequiredImport("from plc4py.spi.values.PlcValues import PlcSINT");
                return "PlcBYTE";
            case UINT:
                IntegerTypeReference unsignedIntegerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (unsignedIntegerTypeReference.getSizeInBits() <= 4) {
                    // Ubyte values are represented as unsigned short integers in PLC4Py
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcUSINT");
                    return "PlcUSINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 8) {
                    // Uword values are represented as unsigned integers in PLC4Py
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcUINT");
                    return "PlcUINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 16) {
                    // Uword16 values are represented as unsigned long integers in PLC4Py
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcUDINT");
                    return "PlcUDINT";
                }
                if (unsignedIntegerTypeReference.getSizeInBits() <= 32) {
                    // Uword32 values are represented as unsigned long long integers in PLC4Py
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcULINT");
                    return "PlcULINT";
                }
            case INT:
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                if (integerTypeReference.getSizeInBits() <= 8) {
                    // Byte values are represented as signed integers in PLC4Py
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcSINT");
                    return "PlcSINT";
                }
                if (integerTypeReference.getSizeInBits() <= 16) {
                    // Word values are represented as signed long integers in PLC4Py
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcINT");
                    return "PlcINT";
                }
                if (integerTypeReference.getSizeInBits() <= 32) {
                    // Dword values are represented as signed long long integers in PLC4Py
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcDINT");
                    return "PlcDINT";
                }
                if (integerTypeReference.getSizeInBits() <= 64) {
                    // Lword values are represented as signed long long integers in PLC4Py
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcLINT");
                    return "PlcLINT";
                }
            case FLOAT:
            case UFLOAT:
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                int sizeInBits = floatTypeReference.getSizeInBits();
                if (sizeInBits <= 32) {
                    // Float values are represented as float in PLC4Py
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcREAL");
                    return "PlcREAL";
                }
                if (sizeInBits <= 64) {
                    // Double values are represented as double in PLC4Py
                    emitRequiredImport("from plc4py.spi.values.PlcValues import PlcLREAL");
                    return "PlcLREAL";
                }
            case STRING:
            case VSTRING:
                // String values are represented as a string in PLC4Py
                emitRequiredImport("from plc4py.spi.values.PlcValues import PlcSTRING");
                return "PlcSTRING";
            case TIME:
            case DATE:
            case DATETIME:
                // Time values are represented as a time in PLC4Py
                emitRequiredImport("from plc4py.spi.values.PlcValues import PlcTIME");
                return "PlcTIME";
            default:
                return "";
        }
    }

    /**
     * Returns the null value for the given type reference.
     * 
     * <p>
     * This method is used to generate code for accessing the null value of a type.
     * The null value is returned as a string that can be used directly in Python
     * code.
     * 
     * @param typeReference the type reference to get the null value for
     * @return the null value for the given type reference
     */
    @Override
    public String getNullValueForTypeReference(TypeReference typeReference) {
        if (typeReference instanceof SimpleTypeReference) {
            SimpleTypeReference simpleTypeReference = (SimpleTypeReference) typeReference;
            switch (simpleTypeReference.getBaseType()) {
                case BIT:
                    // Boolean type: False is the null value
                    return "False";
                case BYTE:
                case UINT:
                case INT:
                    // Integer types: 0 is the null value
                    return "0";
                case FLOAT:
                    // Float types: 0.0 is the null value
                    return "0.0";
                case STRING:
                case VSTRING:
                    // String types: an empty string is the null value
                    return "\"\"";
            }
        } else if (typeReference.isEnumTypeReference()) {
            // Enum types: 0 is the null value
            return "0";
        }
        // For all other types, return None as the null value
        return "None";
    }
    

    /**
     * Returns the number of bits for the given simple type reference.
     * 
     * @param simpleTypeReference the simple type reference to get the number of bits for
     * @return the number of bits for the given simple type reference
     */
    public int getNumBits(SimpleTypeReference simpleTypeReference) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                // bit is always 1 bit
                return 1;
            case BYTE:
                // byte is always 8 bits
                return 8;
            case UINT:
            case INT:
                // integer type references have a variable number of bits depending on the
                // size specified by the type reference
                IntegerTypeReference integerTypeReference = (IntegerTypeReference) simpleTypeReference;
                return integerTypeReference.getSizeInBits();
            case FLOAT:
                // float type references have a variable number of bits depending on the
                // precision specified by the type reference
                FloatTypeReference floatTypeReference = (FloatTypeReference) simpleTypeReference;
                return floatTypeReference.getSizeInBits();
            case STRING:
            case VSTRING:
                // string type references have a variable number of bits depending on the
                // length of the string
                StringTypeReference stringTypeReference = (StringTypeReference) simpleTypeReference;
                return stringTypeReference.getSizeInBits();
            default:
                // for all other types, return 0 as there is no meaningfull number of bits
                return 0;
        }
    }

    /**
     * Returns true if the given property field needs pointer access or not.
     * 
     * <p>
     * This method returns true if the property field has an optional attribute and its type reference
     * needs pointer access or not. This is the case if the type reference is not a complex type reference
     * and the type reference is not an array type reference or if the array type reference has an
     * element type reference that is not a complex type reference.
     * 
     * @param field the property field to check
     * @return true if the property field needs pointer access, false otherwise
     */
    public boolean needsPointerAccess(PropertyField field) {
        boolean isAnTypeOfOptional = "optional".equals(field.getTypeName());
        return isAnTypeOfOptional && needPointerAccess(field.getType());
    }

    /**
     * Returns true if the given type reference needs pointer access or not. This
     * method is used to determine whether the generated code needs to use a
     * pointer for accessing data of the given type.
     *
     * <p>
     * This method returns true if the type reference is not a complex type reference
     * and the type reference is not an array type reference or if the array type
     * reference has an element type reference that is not a complex type reference.
     *
     * @param typeReference the type reference to check
     * @return true if the type reference needs pointer access, false otherwise
     */
    public boolean needPointerAccess(TypeReference typeReference) {
        boolean isNotAnComplexTypeReference = !typeReference.isComplexTypeReference();
        boolean arrayTypeIsNotAnComplexTypeReference = !(typeReference.isArrayTypeReference() && typeReference.asArrayTypeReference().orElseThrow().getElementTypeReference().isComplexTypeReference());
        return isNotAnComplexTypeReference && arrayTypeIsNotAnComplexTypeReference;
    }

    /**
     * Gets a special read buffer read method call for a simple type reference and a typed field.
     * The method is a convenience method that is used to call one of the read buffer methods depending
     * on the type of the field. The method is used to generate code for reading data from a buffer.
     *
     * @param logicalName the logical name of the field
     * @param simpleTypeReference the simple type reference of the field
     * @param parserArguments the parser arguments used for the field
     * @param field the typed field
     * @return a string representing the read buffer read method call
     */
    public String getSpecialReadBufferReadMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, TypedField field) {
        return getReadBufferReadMethodCall(logicalName, simpleTypeReference, null, field);
    }

    public String getReadBufferReadMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, TypedField field) {
        return getReadBufferReadMethodCall(logicalName, simpleTypeReference, null, field);
    }

    /**
     * Gets a read buffer read method call for a simple type reference and a typed field.
     * The method will handle cases where the value string is present and in other cases it will default to the read buffer read method call without any arguments.
     *
     * @param logicalName the logical name of the field
     * @param simpleTypeReference the simple type reference of the field
     * @param valueString the value string to be passed to the read buffer read method call
     * @param field the typed field
     * @return the read buffer read method call
     */
    @Override
    public String getReadBufferReadMethodCall(SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        return getReadBufferReadMethodCall("", simpleTypeReference, valueString, field);
    }


    public String getReadBufferReadMethodCall(String logicalName, SimpleTypeReference simpleTypeReference, String valueString, TypedField field) {
        switch (simpleTypeReference.getBaseType()) {
            case BIT:
                String bitType = "bit";
                return "read_buffer.read_" + bitType + "(\"" + logicalName + "\"";
            case BYTE:
                String byteType = "byte";
                return "read_buffer.read_" + byteType + "(\"" + logicalName + "\"";
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
                return "read_buffer.read_" + unsignedIntegerType + "(" + simpleTypeReference.getSizeInBits() + ", logical_name=\"" + logicalName + "\"";
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
                return "read_buffer.read_" + integerType + "(" + simpleTypeReference.getSizeInBits() + ", logical_name=\"" + logicalName + "\"";
            case FLOAT:
                String floatType = (simpleTypeReference.getSizeInBits() <= 32) ? "float" : "double";
                return "read_buffer.read_" + floatType + "(" + simpleTypeReference.getSizeInBits() + ", logical_name=\"" + logicalName + "\"";
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
                return "read_buffer.read_" + stringType + "(" + simpleTypeReference.getSizeInBits() + ", logical_name=\"" + logicalName + "\"";

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
                ByteTypeReference byteTypeReference = (ByteTypeReference) simpleTypeReference;
                return "write_buffer.write_byte(" + fieldName + ", " + byteTypeReference.getSizeInBits() + ", \"" + logicalName + "\"" + writerArgsString + ")";
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
                    logicalName + "\", \"" + encoding + "\"" + writerArgsString + ")";
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
                    logicalName + "\", \"" + encoding + "\"" + writerArgsString + ")";
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

    /**
     * Returns a Python expression that represents the given term parsed according to the given type reference.
     *
     * @param field the field of the given type reference
     * @param resultType the type reference of the given data
     * @param term the term to parse
     * @param parserArguments the parser arguments
     * @param suppressPointerAccess whether to suppress pointer access
     * @return a Python expression representing the parsed term
     */
    public String toParseExpression(Field field, TypeReference resultType, Term term, List<Argument> parserArguments, boolean suppressPointerAccess) {
        Tracer tracer = pythonTracerStart("toParseExpression");
        return tracer + toTypedParseExpression(field, resultType, term, parserArguments, suppressPointerAccess);
    }

    /**
     * Returns a Python expression that represents the given term serialized according to the given type reference.
     *
     * @param field the field of the given type reference
     * @param resultType the type reference of the given data
     * @param term the term to serialize
     * @param serializerArguments the serializer arguments
     * @return a Python expression representing the serialized term
     */
    public String toSerializationExpression(Field field, TypeReference resultType, Term term, List<Argument> serializerArguments) {
        Tracer tracer = pythonTracerStart("toSerializationExpression");
        return tracer + toTypedSerializationExpression(field, resultType, term, serializerArguments);
    }

    /**
     * Returns a Python expression that represents the given term parsed according to the given type reference.
     *
     * @param field the field of the given type reference
     * @param term the term to parse
     * @param parserArguments the parser arguments
     * @return a Python expression representing the parsed term
     */
    public String toBooleanParseExpression(Field field, Term term, List<Argument> parserArguments) {
        Tracer tracer = pythonTracerStart("toBooleanParseExpression");
        return tracer + toTypedParseExpression(field, new DefaultBooleanTypeReference(), term, parserArguments);
    }

    /**
     * Returns a Python expression that represents the given term serialized according to the given type reference.
     *
     * @param field the field of the given type reference
     * @param term the term to serialize
     * @param serializerArguments the serializer arguments
     * @return a Python expression representing the serialized term
     */
    public String toBooleanSerializationExpression(Field field, Term term, List<Argument> serializerArguments) {
        Tracer tracer = pythonTracerStart("toBooleanSerializationExpression");
        return tracer + toTypedSerializationExpression(field, new DefaultBooleanTypeReference(), term, serializerArguments);
    }

    /**
     * Returns a Python expression that represents the given term parsed according to the given type reference.
     *
     * @param field the field of the given type reference
     * @param sizeInBits the size of the integer in bits
     * @param term the term to parse
     * @param parserArguments the parser arguments
     * @return a Python expression representing the parsed term
     */
    public String toIntegerParseExpression(Field field, int sizeInBits, Term term, List<Argument> parserArguments) {
        Tracer tracer = pythonTracerStart("toIntegerParseExpression");
        return tracer + toTypedParseExpression(field, new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.UINT, sizeInBits), term, parserArguments);
    }


    /**
     * Returns a Python expression that represents the given term serialized according to the given type reference.
     *
     * @param field the field of the given type reference
     * @param sizeInBits the size of the integer in bits
     * @param term the term to serialize
     * @param serializerArguments the serializer arguments
     * @return a Python expression representing the serialized term
     */
    public String toIntegerSerializationExpression(Field field, int sizeInBits, Term term, List<Argument> serializerArguments) {
        Tracer tracer = pythonTracerStart("toIntegerSerializationExpression");
        return tracer + toTypedSerializationExpression(field, new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.UINT, sizeInBits), term, serializerArguments);
    }

    public String toTypedParseExpression(Field field, TypeReference fieldType, Term term, List<Argument> parserArguments) {
        Tracer tracer = pythonTracerStart("toTypedParseExpression");
        return tracer + toExpression(field, fieldType, term, parserArguments, null, false, fieldType != null && fieldType.isComplexTypeReference());
    }

    /**
     * Returns a Python expression that represents the given term parsed according to the given type reference.
     *
     * @param field the field of the given type reference
     * @param fieldType the type reference of the given data
     * @param term the term to parse
     * @param parserArguments the parser arguments
     * @param suppressPointerAccess whether to suppress pointer access
     * @return a Python expression representing the parsed term
     */
    public String toTypedParseExpression(Field field, TypeReference fieldType, Term term, List<Argument> parserArguments, boolean suppressPointerAccess) {
        Tracer tracer = pythonTracerStart("toTypedParseExpression");
        return tracer + toExpression(field, fieldType, term, parserArguments, null, false, suppressPointerAccess);
    }

    /**
     * Returns a Python expression that represents the given term serialized according to
     * the given type reference.
     *
     * @param field the field of the given type reference
     * @param fieldType the type reference of the given data
     * @param term the term to serialize
     * @param serializerArguments the serializer arguments
     * @return a Python expression representing the serialized term
     */
    public String toTypedSerializationExpression(Field field, TypeReference fieldType, Term term, List<Argument> serializerArguments) {
        Tracer tracer = pythonTracerStart("toTypedSerializationExpression");
        return tracer + toExpression(field, fieldType, term, null, serializerArguments, true, false);
    }

    /**
     * Returns a Python expression that represents a cast operation for the given
     * type reference.
     *
     * @param typeReference the type reference to cast the data to
     * @return a Python expression representing the cast operation
     */
    private String getCastExpressionForTypeReference(TypeReference typeReference) {
        Tracer tracer = pythonTracerStart("castExpression");
        if (typeReference instanceof SimpleTypeReference) {
            // Simple type reference, just use the language name
            // e.g. int, char, etc.
            return tracer.dive("simpleTypeRef") + getLanguageTypeNameForTypeReference(typeReference);
        } else if (typeReference instanceof ByteOrderTypeReference) {
            // Byte order, use the binary.ByteOrder enum
            return tracer.dive( "byteOrderTypeRef") + "binary.ByteOrder";
        } else if (typeReference != null) {
            // Complex type reference, use the Cast class
            // e.g. Cast[int], Cast[Char], etc.
            return tracer.dive("anyTypeRef") + "Cast" + getLanguageTypeNameForTypeReference(typeReference);
        } else {
            // No type reference, return an empty string
            return tracer.dive("noTypeRef") + "";
        }
    }

    /**
     * Returns a Python expression that represents the given term.
     *
     * @param field the field of the given type reference
     * @param fieldType the type reference of the given data
     * @param term the term to represent
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param serialize whether the data is serialized
     * @param suppressPointerAccess whether to suppress pointer access
     * @return a Python expression
     */
    private String toExpression(
            Field field,
            TypeReference fieldType,
            Term term,
            List<Argument> parserArguments,
            List<Argument> serializerArguments,
            boolean serialize,
            boolean suppressPointerAccess) {
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

    /**
     * Returns a Python expression that represents the given ternary term.
     *
     * @param field the field of the given type reference
     * @param fieldType the type reference of the given data
     * @param ternaryTerm the ternary term to represent
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param serialize whether the data is serialized
     * @param tracer the tracer to use
     * @return a Python expression
     */
    private String toTernaryTermExpression(Field field, TypeReference fieldType, TernaryTerm ternaryTerm, List<Argument> parserArguments,
                                              List<Argument> serializerArguments, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("ternary term instanceOf");
        if ("if".equals(ternaryTerm.getOperation())) {
            // Ternary operation is "if", so we can build an inline if statement
            Term a = ternaryTerm.getA();
            Term b = ternaryTerm.getB();
            Term c = ternaryTerm.getC();
            String castExpressionForTypeReference = getCastExpressionForTypeReference(fieldType);
            String inlineIf = toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, false) + " if " +
                toExpression(field, new DefaultBooleanTypeReference(), a, parserArguments, serializerArguments, serialize, false) +
                " else " +
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
                String aExpression = toExpression(field, null, a, parserArguments, serializerArguments, serialize, true);
                String bExpression = toExpression(field, null, b, parserArguments, serializerArguments, serialize, true);
                return tracer + "bool((" + aExpression + ") " + operation + " (" + bExpression + "))";
            case ">>":
            case "<<":
            case "|":
            case "&":
                tracer = tracer.dive("bitwise");
                // We don't want casts here
                return tracer +
                    toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, true) +
                    operation + " " +
                    toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, true);
            default:
                tracer = tracer.dive("default");
                if (fieldType instanceof StringTypeReference) {
                    tracer = tracer.dive("string type reference");
                    return tracer + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, true) +
                        operation + " " +
                        toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, true);
                }
                return tracer +
                    toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, true) +
                    operation + " " +
                    toExpression(field, fieldType, b, parserArguments, serializerArguments, serialize, true);
        }
    }

    /**
     * Returns a Python expression that represents the given unary term.
     *
     * @param field the field of the given type reference
     * @param typeReference the type reference of the given data
     * @param unaryTerm the unary term to represent
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param serialize whether the data is serialized
     * @param tracer the tracer to use
     * @return a Python expression
     */
    private String toUnaryTermExpression(Field field, TypeReference fieldType, UnaryTerm unaryTerm, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("unary term instanceOf");
        Term a = unaryTerm.getA();
        switch (unaryTerm.getOperation()) {
            case "!":
                // Not
                // !x
                // => "!(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + ")"
                tracer = tracer.dive("case !");
                return tracer + "!(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + ")";
            case "-":
                // Negate
                // -x
                // => "-(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + ")"
                tracer =tracer.dive("case -");
                return tracer + "-(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + ")";
            case "()":
                // Identity
                // x()
                // => "(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + ")"
                tracer = tracer.dive("case ()");
                return tracer + "(" + toExpression(field, fieldType, a, parserArguments, serializerArguments, serialize, false) + ")";
            default:
                throw new FreemarkerException("Unsupported unary operation type " + unaryTerm.getOperation());
        }
    }

    /**
     * Returns a Python expression that represents the given literal term.
     *
     * @param field the field of the given type reference
     * @param typeReference the type reference of the given data
     * @param term the literal term to represent
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param serialize whether the data is serialized
     * @param suppressPointerAccess whether to suppress pointer 
     *  access
     * @param tracer the tracer to extend
     * @return a Python expression that represents the given literal term
     */
    private String toLiteralTermExpression(Field field, TypeReference fieldType, Term term, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("literal term instanceOf");
        if (term instanceof NullLiteral) {
            // Return None for null literals
            tracer = tracer.dive("null literal instanceOf");
            return tracer + "None";
        } else if (term instanceof BooleanLiteral) {
            // Return a string literal with the boolean value
            String bool = Boolean.toString(((BooleanLiteral) term).getValue());
            return tracer + bool.substring(0,1).toUpperCase() + bool.substring(1);
        } else if (term instanceof NumericLiteral) {
            // Return a numeric literal with the appropriate type cast,
            // unless the type reference is a string type reference
            if (getCastExpressionForTypeReference(fieldType).equals("string")) {
                // If the type reference is a string type reference, return a string literal
                tracer = tracer.dive("type reference string");
                return tracer + "(" + ((NumericLiteral) term).getNumber().toString() + ")";
            } else {
                // Otherwise, return a numeric literal with the appropriate type cast
                return tracer + getCastExpressionForTypeReference(fieldType) + "(" + ((NumericLiteral) term).getNumber().toString() + ")";
            }
        } else if (term instanceof HexadecimalLiteral) {
            // Return a string literal with the hexadecimal value
            return tracer + ((HexadecimalLiteral) term).getHexString();
        } else if (term instanceof StringLiteral) {
            // Return a string literal with the value of the string literal
            return tracer + "\"" + ((StringLiteral) term).getValue() + "\"";
        } else if (term instanceof VariableLiteral) {
            // Return a variable expression for the variable literal
            VariableLiteral variableLiteral = (VariableLiteral) term;
            if ("cur_pos".equals(((VariableLiteral) term).getName())) {
                // If the variable is named "cur_pos", return an expression
                // that computes the position from the position aware and the start position
                return "(position_aware.get_pos() - startPos)";
            } else if ("BIG_ENDIAN".equals(((VariableLiteral) term).getName()) && (fieldType instanceof ByteOrderTypeReference)) {
                // If the variable is named "BIG_ENDIAN" and the type reference is a byte order type reference,
                // return the constant "ByteOrder.BIG_ENDIAN"
                return "ByteOrder.BIG_ENDIAN";
            } else if ("LITTLE_ENDIAN".equals(((VariableLiteral) term).getName()) && (fieldType instanceof ByteOrderTypeReference)) {
                // If the variable is named "LITTLE_ENDIAN" and the type reference is a byte order type reference,
                // return the constant "ByteOrder.LITTLE_ENDIAN"
                return "ByteOrder.LITTLE_ENDIAN";
            }
            return tracer + toVariableExpression(field, fieldType, (VariableLiteral) term, parserArguments, serializerArguments, serialize, suppressPointerAccess);
        } else {
            // Throw an exception for unsupported Literal types
            throw new FreemarkerException("Unsupported Literal type " + term.getClass().getName());
        }
    }

    /**
     * Returns a Python expression that accesses the given variable literal.
     * @param field the field of the given type reference
     * @param typeReference the type reference of the given data
     * @param variableLiteral the variable literal to access
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param serialize whether the data is serialized
     * @param suppressPointerAccess whether to suppress pointer access
     * @return a Python expression that accesses the given variable literal.
     */
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
            return tracer + camelCaseToSnakeCase(variableLiteralName) +
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

    /**
     * Returns a Python expression that accesses the given variable literal
     * in an upper-case format.
     *
     * If the given variable literal has arguments, they are enclosed in
     * parentheses and separated by commas.
     *
     * If the given variable literal has an index, it is enclosed in square
     * brackets.
     *
     * @param field the field of the given type reference
     * @param typeReference the type reference of the given data
     * @param variableLiteral the variable literal to access
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param serialize whether the data is serialized
     * @param suppressPointerAccess whether to suppress pointer access
     * @param tracer the tracer for this expression
     * @return a Python expression that accesses the given variable literal
     *         in an upper-case format
     */
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

    /**
     * Returns a Python expression that computes the ceiling of the given
     * variable literal in Python.
     *
     * @param field the field of the given type reference
     * @param variableLiteral the variable literal to access
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param serialize whether to serialize
     * @param suppressPointerAccess whether to suppress pointer access
     * @param tracer the tracer for this expression
     * @return a Python expression that computes the ceiling of the given
     * variable literal in Python
     */
    private String toCeilVariableExpression(Field field, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("ceil");
        Term va = variableLiteral.getArgs().orElseThrow(() -> new FreemarkerException("CEIL needs at least one arg"))
            .stream().findFirst().orElseThrow(IllegalStateException::new);
        // The Ceil function expects 64 bit floating point values.
        TypeReference tr = new DefaultFloatTypeReference(SimpleTypeReference.SimpleBaseType.FLOAT, 64);
        emitRequiredImport("from math import ceil");
        return tracer + "ceil(" + toExpression(field, tr, va, parserArguments, serializerArguments, serialize, suppressPointerAccess) + ")";
    }


    /**
     * Returns an expression that computes the size in bytes of the given
     * variable literal in Python.
     *
     * @param field the field of the given type reference
     * @param typeReference the type reference of the given data
     * @param variableLiteral the variable literal to access
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param suppressPointerAccess whether to suppress pointer access
     * @param tracer the tracer for this expression
     * @return an expression that computes the size in bytes of the given
     * variable literal in Python.
     */
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
            // Access the given variable literal via a serializer argument. This
            // might be the case if the argument is called "read_buffer" or
            // "write_buffer".
            sb.append(va.getName()).append(va.getChild().map(child -> "." + toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, true, suppressPointerAccess, true)).orElse(""));
        } else {
            // Access the given variable literal directly.
            sb.append(toVariableExpression(field, typeReference, va, parserArguments, serializerArguments, true, suppressPointerAccess));
        }

        // We need the getsizeof function from the sys module.
        emitRequiredImport("from plc4py.spi.values.Common import get_size_of_array");
        // Cast the result to the correct type if necessary.
        return tracer + getCastExpressionForTypeReference(typeReference) + "(get_size_of_array(" + sb + "))";
    }

    /**
     * Returns an expression that computes the count of the given variable literal in Python.
     *
     * @param field the field of the given type reference
     * @param typeReference the type reference of the given data
     * @param variableLiteral the variable literal to access
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param serialize whether the data is serialized
     * @param suppressPointerAccess whether to suppress pointer access
     * @param tracer the tracer for this expression
     * @return an expression that computes the count of the given variable literal in Python
     */
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

    /**
     * Returns an expression that computes the length of the given string variable
     * literal in Python.
     *
     * @param field the field of the given type reference
     * @param typeReference the type reference of the given data
     * @param variableLiteral the variable literal to access
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param serialize whether the data is serialized
     * @param suppressPointerAccess whether to suppress pointer access
     * @param tracer the tracer for this expression
     * @return an expression that computes the length of the given string variable
     * literal in Python.
     */
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

    /**
     * Returns a Python expression that casts the given variable literal to the
     * given type.
     *
     * @param field the field of the given type reference
     * @param typeReference the type reference of the given data
     * @param variableLiteral the variable literal to cast
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param serialize whether to serialize
     * @param suppressPointerAccess whether to suppress pointer access
     * @param tracer the tracer for this expression
     * @return a Python expression that casts the given variable literal to the
     * given type
     */
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

    /**
     * Returns an expression that accesses the given optional field in Python.
     * @param field the field of the given type reference
     * @param typeReference the type reference of the given data
     * @param variableLiteral the variable literal to access
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param suppressPointerAccess whether to suppress pointer access
     * @param tracer the tracer for this expression
     * @return an expression that accesses the given optional field in Python
     */
    private String toOptionalVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("optional fields");
        return tracer + (suppressPointerAccess || (typeReference != null && typeReference.isComplexTypeReference()) ? "" : "*") + camelCaseToSnakeCase(variableLiteral.getName()) +
            variableLiteral.getChild().map(child -> "." + camelCaseToSnakeCase(toVariableExpression(field, typeReference, child, parserArguments, serializerArguments, false, suppressPointerAccess, true))).orElse("");
    }

    /**
     * Returns an expression that accesses the given enum constant in Python.
     * @param field the field of the given type reference
     * @param typeReference the type reference of the given data
     * @param variableLiteral the variable literal to access
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param suppressPointerAccess whether to suppress pointer access
     * @param tracer the tracer for this expression
     * @return an expression that accesses the given enum constant in Python
     */
    private String toConstantVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("enum constant");
        VariableLiteral child = variableLiteral.getChild().orElseThrow(() -> new IllegalStateException("A constant expression should have a child"));
        return tracer + variableLiteral.getName() + "." + camelCaseToSnakeCase(child.getName()) + "()" +
            child.getChild().map(childChild -> "." + toVariableExpression(field, typeReference, childChild, parserArguments, serializerArguments, false, suppressPointerAccess, true)).orElse("");
    }

    /**
     * Returns an expression that accesses the given enum constant in Python.
     * @param field the field of the given type reference
     * @param typeReference the type reference of the given data
     * @param variableLiteral the variable literal to access
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param suppressPointerAccess whether to suppress pointer access
     * @param tracer the tracer for this expression
     * @return an expression that accesses the given enum constant in Python
     */
    private String toEnumVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean suppressPointerAccess, Tracer tracer) {
        tracer = tracer.dive("enum");
        VariableLiteral child = variableLiteral.getChild().orElseThrow(() -> new FreemarkerException("Enum should have a child"));
        return tracer + variableLiteral.getName() + "_" + child.getName() +
            child.getChild().map(childChild -> "." + toVariableExpression(field, typeReference, childChild, parserArguments, serializerArguments, false, suppressPointerAccess, true)).orElse("");
    }

    /**
     * Returns an expression that accesses the last item of the context in Python.
     * @param typeReference the type reference of the given data
     * @param serialize whether the data is serialized
     * @param tracer the tracer for this expression
     * @return an expression that accesses the last item of the context in Python
     */
    private String toLastItemVariableExpression(TypeReference typeReference, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("lastItem");
        // The last item is accessed in the context via the utils.GetLastItemFromContext(ctx) method
        return tracer + "utils.GetLastItemFromContext(ctx)";
    }

    /**
     * Returns an expression that accesses the length of the given variable literal in the given type reference.
     * If the variable literal has a child, this expression will be a method call of the form
     * <code>self.length(&lt;child expression&gt;)</code>.
     * Otherwise, this expression will be simply <code>self.length</code>.
     * @param field the field of the given type reference
     * @param variableLiteral the variable literal to access
     * @param serialize whether the data is serialized
     * @param tracer the tracer for this expression
     * @return the expression to access the length of the given variable literal in the given type reference
     */
    private String toLengthVariableExpression(Field field, VariableLiteral variableLiteral, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("length");
        return tracer + (serialize ? ("len(self." + camelCaseToSnakeCase(variableLiteral.getName()) + ")") : ("(" + variableLiteral.getName() + ")"));
    }

    /**
     * Returns the expression to access the value of the given variable literal in the given type reference.
     * If the given variable literal has a child, this expression will be a method call of the form
     * <code>self._value(&lt;child expression&gt;)</code>.
     * Otherwise, this expression will be simply <code>self.value</code>.
     * @param field the field of the given type reference
     * @param typeReference the type reference of the given data
     * @param variableLiteral the variable literal to access
     * @param parserArguments the parser arguments
     * @param serializerArguments the serializer arguments
     * @param serialize whether the data is serialized
     * @param suppressPointerAccess whether to suppress pointer access
     * @param tracer the tracer for this expression
     * @return the expression to access the value of the given variable literal in the given type reference
     */
    private String toValueVariableExpression(Field field, TypeReference typeReference, VariableLiteral variableLiteral, List<Argument> parserArguments, List<Argument> serializerArguments, boolean serialize, boolean suppressPointerAccess, Tracer tracer) {
        final Tracer tracer2 = tracer.dive("_value");
        return variableLiteral.getChild()
            .map(child -> tracer2.dive("withChild") + "self." + toUppercaseVariableExpression(field, typeReference, child, parserArguments, serializerArguments, serialize, suppressPointerAccess, tracer2))
            .orElse(tracer2 + "value");
    }

    /**
     * Returns the expression to access the length in bits of the given data.
     * If the data is serialized, this will be a cast expression.
     * Otherwise it will be just 'length_in_bits'.
     * @param typeReference the type reference of the data
     * @param serialize whether the data is serialized
     * @param tracer the tracer for this expression
     * @return the expression to access the length in bits of the given data
     */
    private String toLengthInBitsVariableExpression(TypeReference typeReference, boolean serialize, Tracer tracer) {
        tracer = tracer.dive("lengthInBits");
        return tracer + (serialize ? getCastExpressionForTypeReference(typeReference) + "(" : "") + "length_in_bits" + (serialize ? "())" : "()");
    }


    /**
     * Returns the expression to access the length in bytes of the given data.
     * If the data is serialized, this will be a cast expression.
     * Otherwise it will be just 'length_in_bytes'.
     * @param typeReference the type reference of the data
     * @param serialize whether the data is serialized
     * @param tracer the tracer for this expression
     * @return the expression to access the length in bytes of the given data
     */
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

    /**
     * Escapes a string value according to the type reference given.
     * <p>
     * This method is used by the generated Python code to properly handle string values when
     * writing data to a PLC device.
     * 
     * @param typeReference the type reference of the field the value is being written to.
     * @param valueString the string value to be escaped.
     * @return an escaped string value.
     */
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
                    // This is necessary because in Python strings are Unicode by default
                    // and the "char" type is not a subset of the string type, but rather
                    // a distinct type. If a one character string is written as a string,
                    // it will be encoded as Unicode and not as a single byte.
                    if (!NumberUtils.isParsable(valueString) && (valueString.length() == 1)) {
                        return "'" + valueString + "'";
                    }
                    break;
                case STRING:
                case VSTRING:
                    // The Python string type is Unicode by default, so just add quotes
                    // to the value to make it a string literal.
                    return "\"" + valueString + "\"";
            }
        }
        return valueString;
    }

    /**
     * Escapes a string value according to the type reference given.
     * <p>
     * This method is used by the generated Python code to properly handle string values when
     * writing data to a PLC device.
     * 
     * @param typeReference the type reference of the field the value is being written to.
     * @param valueString the string value to be escaped.
     * @return an escaped string value.
     */
    public String escapeEnumValue(TypeReference typeReference, String valueString) {
        // Currently the only case in which here complex type references are used are when referencing enum constants.
        if (typeReference != null && typeReference.isNonSimpleTypeReference()) {
            // C doesn't like NULL values for enums, so we have to return something else (we'll treat -1 as NULL)
            if ("null".equals(valueString) || valueString == null) {
                return "0";
            }
            if (valueString.contains(".")) {
                // If the value string contains a dot, it's assumed to be a string
                // representation of an enum constant with a name in the form:
                // "typeName_constantName" (e.g. "MyEnum_MY_CONSTANT")
                String typeName = valueString.substring(0, valueString.indexOf('.'));
                String constantName = valueString.substring(valueString.indexOf('.') + 1);
                return typeName + "_" + constantName;
            }
            // Otherwise, return the value string as is.
            return valueString;
        } else {
            // For simple types, just use the {@link #escapeValue(TypeReference, String)}
            // method.
            return escapeValue(typeReference, valueString);
        }
    }

    /**
     * Returns a list of unique {@link EnumValue}s from the given list.
     * <p>
     * This is used by the generated Python code to filter out duplicates in a list of {@link EnumValue}s.
     * 
     * @param enumValues the list of {@link EnumValue}s to filter.
     * @return a list of unique {@link EnumValue}s.
     */
    public Collection<EnumValue> getUniqueEnumValues(List<EnumValue> enumValues) {
        Map<String, EnumValue> filteredEnumValues = new TreeMap<>();
        // Iterate over the list of enum values and add only the first occurrence of each value to the map
        for (EnumValue enumValue : enumValues) {
            if (!filteredEnumValues.containsKey(enumValue.getValue())) {
                filteredEnumValues.put(enumValue.getValue(), enumValue);
            }
        }
        // Return the values of the map as a list
        return filteredEnumValues.values();
    }

    /**
     * Returns a list of unique {@link DiscriminatedComplexTypeDefinition}s from the given list.
     * <p>
     * This is used by the generated Python code to filter out duplicates in a {@link DiscriminatedComplexTypeDefinition} list.
     * 
     * @param allSwitchCases the list of {@link DiscriminatedComplexTypeDefinition}s to filter.
     * @return a list of unique {@link DiscriminatedComplexTypeDefinition}s.
     */
    public List<DiscriminatedComplexTypeDefinition> getUniqueSwitchCases(List<DiscriminatedComplexTypeDefinition> allSwitchCases) {
        Map<String, DiscriminatedComplexTypeDefinition> switchCases = new LinkedHashMap<>();
        for (DiscriminatedComplexTypeDefinition switchCase : allSwitchCases) {
            // Check if the given switch case already exists in the map.
            // If not, add it to the map.
            if (!switchCases.containsKey(switchCase.getName())) {
                switchCases.put(switchCase.getName(), switchCase);
            }
        }
        // Convert the map values to a list and return it.
        return new ArrayList<>(switchCases.values());
    }

    /**
     * Emits a required import statement for a data io class.
     * <p>
     * This is used by the generated Python code to import classes needed to read and write PLC
     * data.
     *
     * @param requiredImport the name of the data io class to import
     */
    public void emitRequiredImport(String requiredImport) {
        /**
         * Logs the emitted import statement at debug level. This is useful for debugging
         * purposes, as it allows you to quickly identify which classes are being imported for
         * reading and writing PLC data.
         */
        LOGGER.debug("emitting import \"{}\"", requiredImport);
        /**
         * Keeps track of all the required imports that have been emitted. This collection is used
         * when generating the __all__ list for the generated module.
         */
        requiredImports.add(requiredImport);
    }

    /**
     * Emits a required import statement for a data io class.
     * <p>
     * This is used by the generated Python code to import classes needed to read and write PLC
     * data.
     *
     * @param alias the alias of the imported class
     * @param requiredImport the name of the data io class to import
     */
    public void emitRequiredImport(String alias, String requiredImport) {
        LOGGER.debug("emitting import '{} \"{}'\"", alias, requiredImport);
        requiredImports.add(alias + ' ' + '"' + requiredImport + '"');
    }

    public Set<String> getRequiredImports() {
        return requiredImports;
    }

    /**
     * Emits a required import statement for a data io class.
     * <p>
     * This is used by the generated Python code to import classes needed to read and write PLC
     * data.
     * <p>
     * The imported class must be a subclass of {@link PlcDataIo}.
     * @param requiredImport the name of the data io class to import
     */
    public void emitDataIoRequiredImport(String requiredImport) {
        LOGGER.debug("emitting io import '\"{}\"'", requiredImport);
        requiredImportsForDataIo.add(requiredImport);
    }


    /**
     * Emits a required import statement for a data io class.
     * <p>
     * This is used by the generated Python code to import classes needed to read and write PLC
     * data.
     * <p>
     * The alias parameter is used to specify an alias for the imported class. The alias is
     * used in the generated Python code to refer to the imported class.
     * @param alias an alias for the imported class
     * @param requiredImport the class name to import
     */
    public void emitDataIoRequiredImport(String alias, String requiredImport) {
        LOGGER.debug("emitting data io import '{} \"{}'\"", alias, requiredImport);
        requiredImportsForDataIo.add(alias + ' ' + '"' + requiredImport + '"');
    }

    /**
     * Returns a set of required import statements for data io classes.
     * <p>
     * This is used by the generated Python code to import classes needed to read and write PLC
     * data.
     * <p>
     * The set returned is not modified by the code generator and is intended to be used by the
     * generated code.
     * @return a set of required import statements for data io classes
     */
    public Set<String> getRequiredImportsForDataIo() {
        return requiredImportsForDataIo;
    }

    /**
     * Returns the variable name for a given field.
     * <p>
     * If the field is a named field, the name is used.
     * Otherwise, the default is "_".
     * <p>
     * If there is a conflict between fields, the first one found is used.
     * <p>
     * This is used to generate the names of variables in the generated Python code.
     * @param field the field to get the variable name for
     * @return the variable name
     */
    public String getVariableName(Field field) {
        if (!(field instanceof NamedField)) {
            return "_";
        }
        NamedField namedField = (NamedField) field;

        String name = null;
        for (Field curField : ((ComplexTypeDefinition) thisType).getFields()) {
            // If this is the current field, use its name
            if (curField == field) {
                name = namedField.getName();
            } else if (name != null) {
                // Otherwise, check if this field's expressions contain the given name
                // We check for:
                // - Array fields and their loop expressions
                // - Checksum fields and their checksum expressions
                // - Implicit fields and their serialize expressions
                // - Manual array fields and their length, loop, and serialize expressions
                // - Manual fields and their length, parse, and serialize expressions
                // - Optional fields and their condition expressions
                // - Switch fields and their discriminator expressions and parser arguments
                // - Virtual fields and their value expressions
                // - Typed fields and their type's params (if applicable)
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
                    // Check discriminator expressions first
                    for (Term discriminatorExpression : switchField.getDiscriminatorExpressions()) {
                        if (discriminatorExpression.contains(name)) {
                            return name;
                        }
                    }
                    // Then check parser arguments
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
                // Finally, if the field is a typed field, check its type's params (if applicable)
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

    /**
     * Checks whether a given field requires a variable.
     *
     * @param field the field to check
     * @param variable the variable to check for
     * @param serialization whether the field is used in serialization or not
     * @return true if the field requires the variable, false otherwise
     */
    public boolean needsVariable(Field field, String variableName, boolean serialization) {
        // If it's not a serialization, check if the field is an array and the loop
        // expression contains the variable name, or if it's a virtual field and
        // the value expression contains the variable name.
        if (!serialization) {
            if (field instanceof ArrayField) {
                ArrayField arrayField = (ArrayField) field;
                if (arrayField.getLoopExpression().contains(variableName)) {
                    return true;
                }
            }
        }
        // If it's not a serialization, check if the field is a virtual field and
        // the value expression contains the variable name.
        if (field instanceof VirtualField) {
            VirtualField virtualField = (VirtualField) field;
            if (virtualField.getValueExpression().contains(variableName)) {
                return true;
            }
        }
        // If it's not a serialization, check if the field is a padding field and
        // the padding condition or padding value expression contains the variable
        // name.
        if (field instanceof PaddingField) {
            PaddingField paddingField = (PaddingField) field;
            if (paddingField.getPaddingCondition().contains(variableName)) {
                return true;
            }
            if (paddingField.getPaddingValue().contains(variableName)) {
                return true;
            }
        }
        // Otherwise, check if the field's type reference has a parameter that
        // contains the variable name.
        return field.asTypedField()
            .map(typedField -> typedField.getType().asNonSimpleTypeReference()
                .map(nonSimpleTypeReference -> nonSimpleTypeReference.getParams()
                    .map(params -> params.stream()
                        // Stream over all parameters of the type reference and
                        // returns true if any of them contains the variable name.
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

    /**
     * Checks whether the complex type definition requires a "curPos" variable.
     *
     * @return true if a "curPos" variable is required, false otherwise
     */
    public boolean requiresCurPos() {
        // Iterate over all fields of the complex type definition.
        if (thisType instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) this.thisType;

            // If any of the fields require a "curPos" variable, return true.
            for (Field curField : complexTypeDefinition.getFields()) {
                if (requiresVariable(curField, "curPos")) {
                    return true;
                }
            }

            // If none of the fields require a "curPos" variable, return false.
            return false;
        }

        // If the complex type definition is null, return false.
        return false;
    }

    /**
     * Checks whether the complex type definition requires a "startPos" variable.
     *
     * @return true if a "startPos" variable is required, false otherwise
     */
    public boolean requiresStartPos() {
        if (thisType instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) this.thisType;
            // Iterate over all fields of the complex type definition.
            for (Field curField : complexTypeDefinition.getFields()) {
                // If a field requires a "startPos" variable, return true.
                if (requiresVariable(curField, "startPos")) {
                    return true;
                }
            }
        }
        // If none of the fields require a "startPos" variable, return false.
        return false;
    }

    /**
     * Checks whether a given field requires a variable.
     *
     * @param curField the field to check
     * @param variable the variable to check for
     * @return true if the field requires the variable, false otherwise
     */
    public boolean requiresVariable(Field curField, String variable) {
        if (curField.isArrayField()) {
            ArrayField arrayField = (ArrayField) curField;
            // Check if the loop expression of the array field contains the variable
            // name.
            if (arrayField.getLoopExpression().contains(variable)) {
                return true;
            }
        } else if (curField.isOptionalField()) {
            // Check if the optional field has a condition expression that contains
            // the variable name.
            OptionalField optionalField = (OptionalField) curField;
            if (optionalField.getConditionExpression().isPresent() &&
                optionalField.getConditionExpression().orElseThrow(IllegalStateException::new).contains(variable)) {
                return true;
            }
        }
        // Check if the type reference of the field's type has a parameter that
        // contains the variable name.
        return curField.asTypedField()
            .map(typedField -> typedField.getType().asNonSimpleTypeReference()
                .map(nonSimpleTypeReference -> nonSimpleTypeReference.getParams()
                    .map(params -> params.stream()
                        // Stream over all parameters of the type reference and
                        // returns true if any of them contains the variable name.
                        .anyMatch(param -> param.contains(variable))
                    )
                    .orElse(false)
                )
                .orElse(false)
            )
            .orElse(false);
    }

    /**
     * Finds a term by name in a given term.
     *
     * @param baseTerm the term to search
     * @param name the name to search for
     * @return the term with the given name, or null if not found.
     */
    public Term findTerm(Term baseTerm, String name) {
        if (baseTerm instanceof VariableLiteral) {
            VariableLiteral variableLiteral = (VariableLiteral) baseTerm;
            // Check if the variable literal itself has the desired name
            if (variableLiteral.getName().equals(name)) {
                return variableLiteral;
            }
            // Check if any of the arguments of the variable literal have the
            // desired name.
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
            // Recursively search the argument of the unary term.
            UnaryTerm unaryTerm = (UnaryTerm) baseTerm;
            return findTerm(unaryTerm.getA(), name);
        } else if (baseTerm instanceof BinaryTerm) {
            // Recursively search both arguments of the binary term.
            BinaryTerm binaryTerm = (BinaryTerm) baseTerm;
            Term found = findTerm(binaryTerm.getA(), name);
            if (found != null) {
                return found;
            }
            found = findTerm(binaryTerm.getB(), name);
            return found;
        } else if (baseTerm instanceof TernaryTerm) {
            // Recursively search all three arguments of the ternary term.
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

    /**
     * Converts a Python expression that refers to an enum constant to a constant
     * reference.
     *
     * @param expression the expression to convert
     * @return a string containing the enum name and constant name separated by
     *         underscore.
     */
    public String getEnumExpression(String expression) {
        // Split the expression into enum name and enum constant name
        String enumName = expression.substring(0, expression.indexOf('.'));
        String enumConstant = expression.substring(expression.indexOf('.') + 1);

        // Return enum name and constant name separated by underscore
        return enumName + "_" + enumConstant;
    }

    /**
     * Determines if a parser argument needs to be passed as a reference.

     * <p>This method is a bit confusing: it checks if the argument type is
     * complex and then checks if one of the properties of the complex type
     * definition has the same name as the property name given as argument. If
     * such a property exists and its type is not an enum, then the argument
     * needs to be passed as a reference.

     * <p>This method is implemented as it is because the information the
     * {@link ComplexTypeDefinition} provides is not sufficient to determine
     * whether an argument needs to be passed as a reference: the information
     * about the fields of the complex type does not contain information about
     * the type of the field. Instead, it contains information about the type
     * definition of the field, which does not necessarily reflect the type of
     * the field in the context of the given property.

     * <p>This method is a bit broken because it does not consider the type
     * of the argument. Instead, it only looks at the type of the property.
     * This is problematic because the type of the argument could be a
     * reference to a type that is not a complex type, but still needs to be
     * passed as a reference. This is not the case here, because we only
     * consider complex types.

     * @param propertyName the name of the property
     * @param argumentType the type of the parser argument
     * @return whether the parser argument needs to be passed as a reference
     */
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

    /**
     * Capitalizes a string by removing traces and capitalizing the first character
     * of the remaining string.
     *
     * @param str the string to be capitalized
     * @return the capitalized string
     *
     * <p>This method works by first extracting traces from the string using
     * {@link #pythonTracerStart(String)}, then removing all traces from the
     * string using {@link #removeTraces(String)}, and finally capitalizing the
     * first character of the remaining string using
     * {@link org.apache.commons.lang3.StringUtils#capitalize(String)}.
     */
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
     *
     * <p>Method that takes a camel-case string and converts it to a snake-case
     * string. This is done by iterating over the characters in the string and
     * adding an underscore in between uppercase and lowercase letters or when
     * a dash is encountered.
     */
    public String camelCaseToSnakeCase(String camelCase) {
        StringBuilder snakeCase = new StringBuilder();
        final char[] chars = camelCase.toCharArray();
        for (int i = 0; i < chars.length; i++) {
            String lowerCaseChar = String.valueOf(chars[i]).toLowerCase();
            // If the previous letter is a lowercase letter and this one is uppercase, create a new snake-segment.
            if ((i > 0) && !Character.isUpperCase(chars[i - 1]) && Character.isUpperCase(chars[i])) {
                snakeCase.append('_').append(lowerCaseChar);
            }
            // If the next letter is lowercase, or if this is the last letter, append this character directly.
            else if ((i < (chars.length - 2)) && (Character.isLowerCase(chars[i + 1]) || i == (chars.length - 1))) {
                snakeCase.append(lowerCaseChar);
            }
            // If this is uppercase and the previous one is too ..., append this letter in lowercase.
            else if ((i > 0) && Character.isUpperCase(chars[i - 1]) && Character.isUpperCase(chars[i])) {
                snakeCase.append(lowerCaseChar);
            }
            // If this is a dash, append it directly to the snake-case string.
            else if (chars[i] == '-') {
                snakeCase.append("_");
            }
            // Otherwise, append the character to the snake-case string.
            else {
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

    /**
     * Creates a Python tracing object with a given base name. The Tracer object
     * will use a '/' as separator, and prefix and suffix strings to include
     * tracing info.
     *
     * @param base base name of the tracing object
     * @return a Python Tracer object
     */
    public Tracer pythonTracerStart(String base) {
        return new Tracer(base) {
            /**
             * Returns the separator used in the tracing object.
             *
             * @return '/' as separator
             */
            protected String separator() {
                return "/";
            }

            /**
             * Returns the prefix string used in the tracing object.
             * @return a string used to prefix tracing information
             */
            protected String prefix() {
                return "\"\"\"";
            }

            /**
             * Returns the suffix string used in the tracing object.
             * @return a string used to suffix tracing information
             */
            protected String suffix() {
                return "\"\"\"";
            }
        };
    }

}