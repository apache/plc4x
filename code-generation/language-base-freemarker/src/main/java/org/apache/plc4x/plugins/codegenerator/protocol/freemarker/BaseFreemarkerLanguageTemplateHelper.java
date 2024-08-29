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
package org.apache.plc4x.plugins.codegenerator.protocol.freemarker;

import net.objecthunter.exp4j.Expression;
import net.objecthunter.exp4j.ExpressionBuilder;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions.DefaultDataIoTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.DefaultBooleanTypeReference;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.DefaultIntegerTypeReference;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.references.DefaultUndefinedTypeReference;
import org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms.WildcardTerm;
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.references.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;
import org.apache.plc4x.plugins.codegenerator.types.terms.VariableLiteral;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.function.Function;
import java.util.function.Supplier;
import java.util.stream.Collectors;

public abstract class BaseFreemarkerLanguageTemplateHelper implements FreemarkerLanguageTemplateHelper {

    private static final Logger LOGGER = LoggerFactory.getLogger(BaseFreemarkerLanguageTemplateHelper.class);

    protected final TypeDefinition thisType;
    protected final String protocolName;
    protected final String flavorName;
    protected final Map<String, TypeDefinition> types;

    public static final TypeReference INT_TYPE_REFERENCE = new DefaultIntegerTypeReference(SimpleTypeReference.SimpleBaseType.INT, 32);

    public TypeReference getIntTypeReference() {
        return INT_TYPE_REFERENCE;
    }

    public static final TypeReference BOOL_TYPE_REFERENCE = new DefaultBooleanTypeReference();

    public TypeReference getBoolTypeReference() {
        return BOOL_TYPE_REFERENCE;
    }

    public static final TypeReference ANY_TYPE_REFERENCE = new DefaultUndefinedTypeReference();

    public TypeReference getAnyTypeReference() {
        return ANY_TYPE_REFERENCE;
    }

    protected BaseFreemarkerLanguageTemplateHelper(TypeDefinition thisType, String protocolName, String flavorName, Map<String, TypeDefinition> types) {
        this.thisType = thisType;
        this.protocolName = protocolName;
        this.flavorName = flavorName;
        this.types = types;
    }

    public String getProtocolName() {
        return protocolName;
    }

    public String getFlavorName() {
        return flavorName;
    }

    public Map<String, TypeDefinition> getTypeDefinitions() {
        return types;
    }

    public List<TypeDefinition> getComplexTypeRootDefinitions() {
        return types.values().stream()
            .filter(ComplexTypeDefinition.class::isInstance)
            .filter(typeDefinition -> !(typeDefinition instanceof DiscriminatedComplexTypeDefinition))
            .collect(Collectors.toList());
    }

    /* *********************************************************************************
     * Methods related to type-references.
     **********************************************************************************/

    protected EnumTypeDefinition getEnumTypeDefinition(TypeReference typeReference) {
        NonSimpleTypeReference nonSimpleTypeReference = typeReference.asNonSimpleTypeReference().orElseThrow(
            () -> new FreemarkerException("type reference for enum types must be of type non simple type"));
        String typeName = nonSimpleTypeReference.getName();
        final TypeDefinition typeDefinition = nonSimpleTypeReference.getTypeDefinition();
        if (typeDefinition == null) {
            throw new FreemarkerException("Couldn't find given enum type definition with name " + typeName);
        }
        // TODO: same here. It is named complex type reference but it references a enum...
        if (!typeDefinition.isEnumTypeDefinition()) {
            throw new FreemarkerException("Referenced type with name " + typeName + " is not an enum type");
        }
        return (EnumTypeDefinition) typeDefinition;
    }

    /**
     * Enums are always based on a main type. This helper accesses this information in a safe manner.
     *
     * @param typeReference type reference
     * @return simple type reference for the enum type referenced by the given type reference
     */
    public SimpleTypeReference getEnumBaseTypeReference(TypeReference typeReference) {
        // Enum types always have simple type references.
        return getEnumTypeDefinition(typeReference).getType().orElseThrow();
    }

    public SimpleTypeReference getEnumFieldTypeReference(TypeReference typeReference, String constantName) {
        return (SimpleTypeReference) getEnumTypeDefinition(typeReference).getConstantType(constantName);
    }

    /* *********************************************************************************
     * Methods related to fields.
     **********************************************************************************/

    public boolean hasFieldOfType(String fieldTypeName) {
        Objects.requireNonNull(fieldTypeName);
        if (thisType instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) this.thisType;
            return complexTypeDefinition.getFields().stream()
                .anyMatch(field -> fieldTypeName.equals(field.getTypeName()));
        }
        return false;
    }

    public boolean hasFieldsWithNames(List<Field> fields, String... names) {
        for (String name : names) {
            boolean foundName = false;
            for (Field field : fields) {
                if (field instanceof NamedField && name.equals(((NamedField) field).getName())) {
                    foundName = true;
                    break;
                }
            }
            if (!foundName) {
                return false;
            }
        }
        // TODO: document why true is returned here.
        return true;
    }

    // TODO: check or describe why a instanceOf EnumField is not sufficient here
    public boolean isEnumField(Field field) {
        if (!(field instanceof TypedField)) {
            return false;
        }
        TypedField typedField = (TypedField) field;
        TypeReference typeReference = typedField.getType();
        if (!typeReference.isNonSimpleTypeReference()) {
            return false;
        }
        TypeDefinition typeDefinition = typeReference.asNonSimpleTypeReference().orElseThrow()
            .getTypeDefinition();
        return typeDefinition instanceof EnumTypeDefinition;
    }

    /* *********************************************************************************
     * Methods related to terms and expressions.
     **********************************************************************************/
    protected int evaluateFixedValueExpression(Term term) {
        Objects.requireNonNull(term);
        final Expression expression = new ExpressionBuilder(term.stringRepresentation()).build();
        return (int) expression.evaluate();
    }

    /* *********************************************************************************
     * Methods related to discriminators.
     **********************************************************************************/

    /**
     * Get a list of the types for every discriminator name.
     *
     * @return Map mapping discriminator names to types.
     */
    public Map<String, TypeReference> getDiscriminatorTypes() {
        return getDiscriminatorTypes(thisType);
    }
    public Map<String, TypeReference> getDiscriminatorTypes(TypeDefinition type) {
        // Get the parent type (Which contains the typeSwitch field)
        SwitchField switchField = null;
        Function<String, TypeReference> typeRefRetriever = null;
        if (type.isDiscriminatedComplexTypeDefinition()) {
            DiscriminatedComplexTypeDefinition discriminatedComplexTypeDefinition = type.asDiscriminatedComplexTypeDefinition().orElseThrow();
            switchField = discriminatedComplexTypeDefinition.getSwitchField().orElse(null);
            typeRefRetriever = propertyName -> discriminatedComplexTypeDefinition.getTypeReferenceForProperty(propertyName).orElse(null);
            // Please forgive us, we didn't know what we were doing.
            if(switchField == null) {
                ComplexTypeDefinition parentType = type.asDiscriminatedComplexTypeDefinition().orElseThrow().getParentType().orElseThrow();
                switchField = parentType.getSwitchField().orElse(null);
                typeRefRetriever = propertyName -> parentType.getTypeReferenceForProperty(propertyName).orElse(null);
            }
        } else if (type.isDataIoTypeDefinition()) {
            final DefaultDataIoTypeDefinition dataIoTypeDefinition = (DefaultDataIoTypeDefinition) type;
            switchField = dataIoTypeDefinition.getSwitchField().orElseThrow();
            typeRefRetriever = propertyName -> type.getParserArguments()
                .orElse(Collections.emptyList())
                .stream()
                .filter(argument -> argument.getName().equals(propertyName))
                .findFirst()
                .map(Argument::getType)
                .orElse(null);
        } else if (type.isComplexTypeDefinition()) {
            switchField = ((ComplexTypeDefinition) type).getSwitchField().orElse(null);
            typeRefRetriever = propertyName -> ((ComplexTypeDefinition) type).getTypeReferenceForProperty(propertyName).orElse(null);
        }
        // Get the typeSwitch field from that.
        if (switchField == null) {
            return Collections.emptyMap();
        }
        Map<String, TypeReference> discriminatorTypes = new TreeMap<>();
        for (VariableLiteral variableLiteral : switchField.getDiscriminatorExpressions()) {
            // Get some symbolic name we can use.
            String discriminatorName = variableLiteral.getDiscriminatorName();
            final TypeReference typeReference = typeRefRetriever.apply(variableLiteral.getName());
            Optional<TypeReference> discriminatorType = typeReference.getDiscriminatorType(variableLiteral);
            if (discriminatorType.isEmpty()) {
                throw new RuntimeException("no type for " + discriminatorName);
            }
            discriminatorTypes.put(discriminatorName, discriminatorType.orElse(null));
        }
        return discriminatorTypes;
    }

    public TypeReference getArgumentType(TypeReference typeReference, int index) {
        Objects.requireNonNull(typeReference, "type reference must not be null");
        NonSimpleTypeReference complexTypeReference = typeReference.asNonSimpleTypeReference().orElseThrow(() -> new FreemarkerException("Only non simple type references supported here."));
        return complexTypeReference.getArgumentType(index);
    }

    public boolean discriminatorValueNeedsStringEqualityCheck(Term term) {
        return discriminatorValueNeedsStringEqualityCheck(term, thisType);
    }

    public boolean discriminatorValueNeedsStringEqualityCheck(Term term, TypeDefinition typeDefinition) {
        if (!(term instanceof VariableLiteral)) {
            return false;
        }

        VariableLiteral variableLiteral = (VariableLiteral) term;
        // If this literal references an Enum type, then we have to output it differently.
        if (getTypeDefinitions().get(variableLiteral.getName()) instanceof EnumTypeDefinition) {
            return false;
        }

        if (typeDefinition instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) typeDefinition;
            boolean found = false;
            for (Field field : complexTypeDefinition.getFields()) {
                if (field instanceof NamedField) {
                    if (((NamedField) field).getName().equals(variableLiteral.getName())) {
                        if (field instanceof TypedField) {
                            TypedField typedField = (TypedField) field;
                            TypeReference type = typedField.getType();
                            found = (type instanceof StringTypeReference) || (type instanceof VstringTypeReference);
                            break;
                        }
                    }
                }
            }
            if (found) {
                return true;
            }
        }
        for (Argument argument : typeDefinition.getParserArguments()
            .orElse(Collections.emptyList())) {
            if (argument.getName().equals(variableLiteral.getName())) {
                TypeReference type = argument.getType();
                return (type instanceof StringTypeReference) || (type instanceof VstringTypeReference);
            }
        }
        return false;
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

    public Collection<EnumValue> getEnumValuesForUniqueConstantValues(List<EnumValue> enumValues, String constantName) {
        Map<String, EnumValue> filteredEnumValues = new TreeMap<>();
        for (EnumValue enumValue : enumValues) {
            String key = enumValue.getConstant(constantName).orElseThrow(() -> new FreemarkerException("No constant name " + constantName + " found in enum value" + enumValue));
            if (!filteredEnumValues.containsKey(key)) {
                filteredEnumValues.put(key, enumValue);
            }
        }
        return filteredEnumValues.values();
    }

    public SimpleTypeReference getEnumFieldSimpleTypeReference(NonSimpleTypeReference type, String fieldName) {
        if (!(type.getTypeDefinition() instanceof EnumTypeDefinition)
            || !(((EnumTypeDefinition) type.getTypeDefinition()).getConstantType(fieldName) instanceof SimpleTypeReference)) {
            throw new IllegalArgumentException("not an enum type or enum constant is not a simple type");
        }
        return (SimpleTypeReference) ((EnumTypeDefinition) type.getTypeDefinition()).getConstantType(fieldName);
    }

    /**
     * Confirms if a variable is an implicit variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @return boolean returns true if the variable's name is an implicit field
     */
    protected boolean isVariableLiteralImplicitField(VariableLiteral variableLiteral) {
        return thisType.asComplexTypeDefinition()
            .map(complexTypeDefinition -> complexTypeDefinition.isVariableLiteralImplicitField(variableLiteral))
            .orElse(false);
    }

    /**
     * Confirms if a variable is an virtual variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @return boolean returns true if the variable's name is an virtual field
     */
    protected boolean isVariableLiteralVirtualField(VariableLiteral variableLiteral) {
        return thisType.asComplexTypeDefinition()
            .map(complexTypeDefinition -> complexTypeDefinition.isVariableLiteralVirtualField(variableLiteral))
            .orElse(false);
    }

    /**
     * Returns the implicit field that has the same name as the variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @return ImplicitField returns the implicit field that corresponds to the variable's name.
     */
    protected ImplicitField getReferencedImplicitField(VariableLiteral variableLiteral) {
        return thisType.asComplexTypeDefinition()
            .map(complexTypeDefinition -> complexTypeDefinition.getReferencedImplicitField(variableLiteral))
            .orElse(null);
    }


    // TODO: replace that with term.isWildcard() (once the referenced wildcard term from build utils is used)
    public boolean isWildcard(Term term) {
        return term instanceof WildcardTerm;
    }

    /**
     * can be used to throw an exception from the template
     *
     * @param message the message
     * @return the exception
     */
    public Supplier<FreemarkerException> fail(String message) {
        return () -> new FreemarkerException(message);
    }

    public void info(String message, Object... objects) {
        LOGGER.info(message, objects);
    }
}
