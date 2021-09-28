/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
import org.apache.plc4x.plugins.codegenerator.types.definitions.*;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.fields.*;
import org.apache.plc4x.plugins.codegenerator.types.references.*;
import org.apache.plc4x.plugins.codegenerator.types.terms.*;

import java.util.*;
import java.util.function.Supplier;
import java.util.stream.Collectors;

public abstract class BaseFreemarkerLanguageTemplateHelper implements FreemarkerLanguageTemplateHelper {

    protected final TypeDefinition thisType;
    protected final String protocolName;
    protected final String flavorName;
    protected final Map<String, TypeDefinition> types;

    // In mspec we are using some virtual virtual fields that are useful for code generation.
    // As they should be shared over all language template implementations,
    // these are defined here manually.
    private static final Map<String, SimpleTypeReference> builtInFields;

    public static final String CUR_POS = "curPos";
    public static final String START_POS = "startPos";
    public static final String LAST_ITEM = "lastItem";
    public static final String IMPLICIT = "implicit";
    public static final String VIRTUAL = "virtual";
    public static final String DISCRIMINATOR = "discriminator";

    static {
        builtInFields = new HashMap<>();
        builtInFields.put(CUR_POS, new SimpleTypeReference() {
            @Override
            public SimpleBaseType getBaseType() {
                return SimpleBaseType.UINT;
            }

            @Override
            public int getSizeInBits() {
                return 16;
            }
        });
        builtInFields.put(START_POS, new SimpleTypeReference() {
            @Override
            public SimpleBaseType getBaseType() {
                return SimpleBaseType.UINT;
            }

            @Override
            public int getSizeInBits() {
                return 16;
            }
        });
        builtInFields.put(LAST_ITEM, new SimpleTypeReference() {
            @Override
            public SimpleBaseType getBaseType() {
                return SimpleBaseType.BIT;
            }

            @Override
            public int getSizeInBits() {
                return 1;
            }
        });
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

    protected static Map<String, SimpleTypeReference> getBuiltInFieldTypes() {
        return builtInFields;
    }

    /* *********************************************************************************
     * Methods related to type-references.
     **********************************************************************************/

    /**
     * @param typeReference type reference
     * @return true if the given type reference is a simple type reference.
     * @deprecated use method of {@link TypeReference}
     */
    @Deprecated
    public boolean isSimpleTypeReference(TypeReference typeReference) {
        if (typeReference == null) {
            return false;
        }
        return typeReference.isSimpleTypeReference();
    }

    /**
     * @param typeReference type reference
     * @return true if the given type reference is a byte based type reference.
     * @deprecated use method of {@link TypeReference}
     */
    @Deprecated
    public boolean isByteBased(TypeReference typeReference) {
        if (typeReference == null) {
            return false;
        }
        return typeReference.isByteBased();
    }

    /**
     * @param typeReference type reference
     * @return true if the given type reference is a complex type reference.
     * @deprecated use method of {@link TypeReference}
     */
    @Deprecated
    public boolean isComplexTypeReference(TypeReference typeReference) {
        if (typeReference == null) {
            return false;
        }
        return typeReference.isComplexTypeReference();
    }

    public boolean isEnumTypeReference(TypeReference typeReference) {
        if (!isComplexTypeReference(typeReference)) {
            return false;
        }
        return getTypeDefinitionForTypeReference(typeReference) instanceof EnumTypeDefinition;
    }

    /**
     * @deprecated use method of {@link TypeReference}
     */
    @Deprecated
    public boolean isStringTypeReference(TypeReference typeReference) {
        if (typeReference == null) {
            return false;
        }
        return typeReference.isStringTypeReference();
    }

    /**
     * Helper for collecting referenced complex types as these usually ned to be
     * imported in some way.
     *
     * @return Collection of all complex type references used in fields or enum constants.
     */
    public Collection<String> getComplexTypeReferences() {
        return getComplexTypeReferences(thisType);
    }

    /**
     * Helper for collecting referenced complex types as these usually need to be
     * imported in some way.
     *
     * @param baseType the base type we want to get the type references from
     * @return collection of complex type references used in the type.
     */
    public Collection<String> getComplexTypeReferences(TypeDefinition baseType) {
        return getComplexTypeReferences(baseType, new HashSet<>());
    }

    public Collection<String> getComplexTypeReferences(TypeDefinition baseType, Set<String> complexTypeReferences) {
        // We add ourselves to avoid a stackoverflow
        complexTypeReferences.add(baseType.getName());
        // If this is a subtype of a discriminated type, we have to add a reference to the parent type.
        if (baseType instanceof DiscriminatedComplexTypeDefinition) {
            DiscriminatedComplexTypeDefinition discriminatedComplexTypeDefinition = (DiscriminatedComplexTypeDefinition) baseType;
            if (!discriminatedComplexTypeDefinition.isAbstract()) {
                String typeReferenceName = ((ComplexTypeReference) discriminatedComplexTypeDefinition.getParentType().getTypeReference()).getName();
                complexTypeReferences.add(typeReferenceName);
            }
        }
        // If it's a complex type definition, add all the types referenced by any property fields
        // (Includes any types referenced by sub-types in case this is a discriminated type parent)
        if (baseType instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) baseType;
            for (Field field : complexTypeDefinition.getFields()) {
                if (field instanceof PropertyField) {
                    PropertyField propertyField = (PropertyField) field;
                    if (propertyField.getType() instanceof ComplexTypeReference) {
                        ComplexTypeReference complexTypeReference = (ComplexTypeReference) propertyField.getType();
                        complexTypeReferences.add(complexTypeReference.getName());
                    }
                } else if (field instanceof SwitchField) {
                    SwitchField switchField = (SwitchField) field;
                    for (DiscriminatedComplexTypeDefinition switchCase : switchField.getCases()) {
                        if (complexTypeReferences.contains(switchCase.getName())) {
                            continue;
                        }
                        complexTypeReferences.addAll(getComplexTypeReferences(switchCase, complexTypeReferences));
                    }
                }
            }
        } else if (baseType instanceof EnumTypeDefinition) {// In case this is an enum type, we have to check all the constant types.
            EnumTypeDefinition enumTypeDefinition = (EnumTypeDefinition) baseType;
            for (String constantName : enumTypeDefinition.getConstantNames()) {
                final TypeReference constantType = enumTypeDefinition.getConstantType(constantName);
                if (constantType instanceof ComplexTypeReference) {
                    ComplexTypeReference complexTypeReference = (ComplexTypeReference) constantType;
                    complexTypeReferences.add(complexTypeReference.getName());
                }
            }
        }
        // If the type has any parser arguments, these have to be checked too.
        baseType.getParserArguments().ifPresent(arguments -> arguments.stream()
            .map(Argument::getType)
            .map(TypeReferenceConversions::asComplexTypeReference)
            .filter(Optional::isPresent)
            .map(Optional::get)
            .map(ComplexTypeReference::getName)
            .forEach(complexTypeReferences::add)
        );

        // We remove ourselves to avoid a stackoverflow
        complexTypeReferences.remove(baseType.getName());
        return complexTypeReferences;
    }

    /**
     * Little helper to return the {@link TypeReference} of a given property.
     *
     * @param baseType     base type definition that contains the given property.
     * @param propertyName name of the property
     * @return the type reference of the given property
     */
    public Optional<TypeReference> getTypeReferenceForProperty(ComplexTypeDefinition baseType, String propertyName) {
        // If this is a built-in type, use that.
        if (builtInFields.containsKey(propertyName)) {
            return Optional.of(builtInFields.get(propertyName));
        }
        // Check if the expression root is referencing a field
        final Optional<PropertyField> propertyFieldOptional = baseType.getPropertyFields().stream()
            .filter(propertyField -> propertyField.getName().equals(propertyName))
            .findFirst();
        if (propertyFieldOptional.isPresent()) {
            return propertyFieldOptional.map(PropertyField::getType);
        }
        // Check if the expression is a ImplicitField
        final Optional<ImplicitField> implicitFieldOptional = baseType.getFields().stream()
            .filter(ImplicitField.class::isInstance)
            .map(ImplicitField.class::cast)
            .filter(implicitField -> implicitField.getName().equals(propertyName))
            .findFirst();
        if (implicitFieldOptional.isPresent()) {
            return implicitFieldOptional.map(ImplicitField::getType);
        }
        // Check if the expression is a VirtualField
        final Optional<VirtualField> virtualFieldOptional = baseType.getFields().stream()
            .filter(VirtualField.class::isInstance)
            .map(VirtualField.class::cast)
            .filter(virtualField -> virtualField.getName().equals(propertyName))
            .findFirst();
        if (virtualFieldOptional.isPresent()) {
            return virtualFieldOptional.map(VirtualField::getType);
        }
        // Check if the expression root is referencing an argument
        final Optional<Argument> argumentOptional = baseType.getParserArguments()
            .orElse(Collections.emptyList())
            .stream()
            .filter(argument -> argument.getName().equals(propertyName))
            .findFirst();
        if (argumentOptional.isPresent()) {
            return argumentOptional.map(Argument::getType);
        }
        // Check if the expression is a DiscriminatorField
        // This is a more theoretical case where the expression is referencing a discriminator value of the current type
        final Optional<DiscriminatorField> discriminatorFieldOptional = baseType.getFields().stream()
            .filter(DiscriminatorField.class::isInstance)
            .map(DiscriminatorField.class::cast)
            .filter(discriminatorField -> discriminatorField.getName().equals(propertyName))
            .findFirst();
        if (discriminatorFieldOptional.isPresent()) {
            return discriminatorFieldOptional.map(DiscriminatorField::getType);
        }
        return Optional.empty();
    }

    /**
     * Enums are always based on a main type. This helper accesses this information in a safe manner.
     *
     * @param typeReference type reference
     * @return simple type reference for the enum type referenced by the given type reference
     */
    public SimpleTypeReference getEnumBaseTypeReference(TypeReference typeReference) {
        if (!(typeReference instanceof ComplexTypeReference)) {
            throw new FreemarkerException("type reference for enum types must be of type complex type");
        }
        String typeName = ((ComplexTypeReference) typeReference).getName();
        final TypeDefinition typeDefinition = types.get(typeName);
        if (typeDefinition == null) {
            throw new FreemarkerException("Couldn't find given enum type definition with name " + typeName);
        }
        if (!(typeDefinition instanceof EnumTypeDefinition)) {
            throw new FreemarkerException("Referenced type with name " + typeName + " is not an enum type");
        }
        EnumTypeDefinition enumTypeDefinition = (EnumTypeDefinition) typeDefinition;
        // Enum types always have simple type references.
        return (SimpleTypeReference) enumTypeDefinition.getType();
    }

    /* *********************************************************************************
     * Methods related to fields.
     **********************************************************************************/

    public boolean hasFieldOfType(String fieldTypeName) {
        if (thisType instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) this.thisType;
            return complexTypeDefinition.getFields().stream()
                .anyMatch(field -> field.getTypeName().equals(fieldTypeName));
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

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public Field getFieldForNameFromCurrentOrParent(String fieldName) {
        if (!(thisType instanceof ComplexTypeDefinition)) {
            return null;
        }
        ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) this.thisType;
        return complexTypeDefinition.getPropertyFieldFromThisOrParentByName(fieldName).orElse(null);
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public Field getFieldForNameFromCurrent(String fieldName) {
        if (!(thisType instanceof ComplexTypeDefinition)) {
            return null;
        }
        ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) this.thisType;
        return complexTypeDefinition.getPropertyFieldByName(fieldName).orElse(null);
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isAbstractField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isAbstractField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isArrayField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isArrayField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isChecksumField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isChecksumField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isConstField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isConstField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isDiscriminatorField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isDiscriminatorField();
    }

    // TODO: check or describe why a instanceOf EnumField is not sufficient here
    public boolean isEnumField(Field field) {
        if (!(field instanceof TypedField)) {
            return false;
        }
        TypedField typedField = (TypedField) field;
        TypeReference typeReference = typedField.getType();
        if (isSimpleTypeReference(typeReference)) {
            return false;
        }
        TypeDefinition typeDefinition = getTypeDefinitionForTypeReference(typedField.getType());
        return typeDefinition instanceof EnumTypeDefinition;
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isImplicitField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isImplicitField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isManualArrayField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isManualArrayField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isNamedField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isNamedField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isOptionalField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isOptionalField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isPaddingField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isPaddingField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isPropertyField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isPropertyField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isReservedField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isReservedField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isSimpleField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isSimpleField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isSwitchField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isSwitchField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isTypedField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isTypedField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isVirtualField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isVirtualField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isCountArrayField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isCountArrayField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isLengthArrayField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isLengthArrayField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isTerminatedArrayField(Field field) {
        if (field == null) {
            return false;
        }
        return field.isTerminatedArrayField();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public SwitchField getSwitchField() {
        if (!(thisType instanceof ComplexTypeDefinition)) {
            return null;
        }
        ComplexTypeDefinition complexTypeDefinition = (ComplexTypeDefinition) this.thisType;
        return complexTypeDefinition.getSwitchField().orElse(null);
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public Collection<Field> getPropertyAndSwitchFields() {
        if (!(thisType instanceof ComplexTypeDefinition)) {
            return Collections.emptyList();
        }
        return ((ComplexTypeDefinition) thisType).getPropertyAndSwitchFields();
    }

    /* *********************************************************************************
     * Methods related to type-definitions.
     **********************************************************************************/

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isDiscriminatedParentTypeDefinition() {
        if (thisType == null) {
            return false;
        }
        return thisType.isDiscriminatedParentTypeDefinition();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isDiscriminatedParentTypeDefinition(TypeDefinition typeDefinition) {
        if (typeDefinition == null) {
            return false;
        }
        return typeDefinition.isDiscriminatedParentTypeDefinition();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isDiscriminatedChildTypeDefinition() {
        if (thisType == null) {
            return false;
        }
        return thisType.isDiscriminatedChildTypeDefinition();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isDiscriminatedChildTypeDefinition(TypeDefinition typeDefinition) {
        if (typeDefinition == null) {
            return false;
        }
        return typeDefinition.isDiscriminatedChildTypeDefinition();
    }

    public TypeDefinition getTypeDefinitionForTypeReference(TypeReference typeReference) {
        if (!isComplexTypeReference(typeReference)) {
            throw new FreemarkerException("Type reference must be a complex type reference");
        }
        ComplexTypeReference complexTypeReference = (ComplexTypeReference) typeReference;
        return getTypeDefinitions().get(complexTypeReference.getName());
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public List<DiscriminatedComplexTypeDefinition> getSubTypeDefinitions() {
        if (!(thisType instanceof ComplexTypeDefinition)) {
            return Collections.emptyList();
        }
        return ((ComplexTypeDefinition) thisType).getSubTypeDefinitions();
    }

    /* *********************************************************************************
     * Methods related to terms and expressions.
     **********************************************************************************/

    /**
     * @deprecated use field method.
     */
    @Deprecated
    protected boolean isFixedValueExpression(Term term) {
        if (term == null) {
            return false;
        }
        return term.isFixedValueExpression();
    }

    protected int evaluateFixedValueExpression(Term term) {
        final Expression expression = new ExpressionBuilder(toString(term)).build();
        return (int) expression.evaluate();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    protected String toString(Term term) {
        if (term == null) {
            return "";
        }
        return term.stringRepresentation();
    }

    /* *********************************************************************************
     * Methods related to discriminators.
     **********************************************************************************/

    private Optional<TypeReference> getDiscriminatorType(ComplexTypeDefinition parentType, Term disciminatorExpression) {
        if (!(disciminatorExpression instanceof VariableLiteral)) {
            throw new FreemarkerException("Currently no arithmetic expressions are supported as discriminator expressions.");
        }
        VariableLiteral variableLiteral = (VariableLiteral) disciminatorExpression;
        Optional<TypeReference> type = getTypeReferenceForProperty(parentType, variableLiteral.getName());
        // If we found something but there's a "rest" left, we got to use the type we
        // found in this level, get that type's definition and continue from there.
        if (type.isPresent() && (variableLiteral.getChild().isPresent())) {
            TypeReference typeReference = type.get();
            if (typeReference instanceof ComplexTypeReference) {
                ComplexTypeReference complexTypeReference = (ComplexTypeReference) typeReference;
                final TypeDefinition typeDefinition = this.types.get(complexTypeReference.getName());
                if (typeDefinition instanceof ComplexTypeDefinition) {
                    return getDiscriminatorType((ComplexTypeDefinition) typeDefinition, variableLiteral.getChild().get());
                }
            }
        }
        return type;
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public List<String> getDiscriminatorNames() {
        if (thisType == null) {
            return Collections.emptyList();
        }
        return thisType.getDiscriminatorNames();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isNonDiscriminatorField(String discriminatorName) {
        if (!(thisType instanceof ComplexTypeDefinition)) {
            throw new ClassCastException(thisType + " not a" + ComplexTypeDefinition.class.getName());
        }
        return ((ComplexTypeDefinition) thisType).isNonDiscriminatorField(discriminatorName);
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public boolean isDiscriminatorField(String discriminatorName) {
        if (thisType == null) {
            return false;
        }
        return thisType.isDiscriminatorField(discriminatorName);
    }

    /**
     * Get a list of the types for every discriminator name.
     *
     * @return Map mapping discriminator names to types.
     */
    public Map<String, TypeReference> getDiscriminatorTypes() {
        // Get the parent type (Which contains the typeSwitch field)
        ComplexTypeDefinition parentType;
        if (thisType instanceof DiscriminatedComplexTypeDefinition) {
            parentType = (ComplexTypeDefinition) thisType.getParentType();
        } else {
            parentType = (ComplexTypeDefinition) thisType;
        }
        // Get the typeSwitch field from that.
        // TODO: map
        final SwitchField switchField = parentType.getSwitchField().orElse(null);
        if (switchField == null) {
            return Collections.emptyMap();
        }
        Map<String, TypeReference> discriminatorTypes = new TreeMap<>();
        for (Term discriminatorExpression : switchField.getDiscriminatorExpressions()) {
            // Get some symbolic name we can use.
            String discriminatorName = discriminatorExpression.getDiscriminatorName();
            Optional<TypeReference> discriminatorType = getDiscriminatorType(parentType, discriminatorExpression);
            discriminatorTypes.put(discriminatorName, discriminatorType.orElse(null));
        }
        return discriminatorTypes;
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public Map<String, String> getDiscriminatorValues(TypeDefinition type) {
        if (!(type instanceof DiscriminatedComplexTypeDefinition)) {
            return Collections.emptyMap();
        }
        return ((DiscriminatedComplexTypeDefinition) type).getDiscriminatorMap();
    }

    /**
     * @deprecated use field method.
     */
    @Deprecated
    public Map<String, Map<String, String>> getDiscriminatorValues() {
        if (thisType == null) {
            return Collections.emptyMap();
        }
        return thisType.getDiscriminatorCaseToKeyValueMap();
    }

    public TypeReference getArgumentType(TypeReference typeReference, int index) {
        Objects.requireNonNull(typeReference, "type reference must not be null");
        ComplexTypeReference complexTypeReference = typeReference.asComplexTypeReference().orElseThrow(() -> new FreemarkerException("Only complex type references supported here."));
        if (!getTypeDefinitions().containsKey(complexTypeReference.getName())) {
            throw new FreemarkerException("Could not find definition of complex type " + complexTypeReference.getName());
        }
        TypeDefinition complexTypeDefinition = getTypeDefinitions().get(complexTypeReference.getName());
        List<Argument> parserArguments = complexTypeDefinition.getParserArguments().orElseThrow(() -> new FreemarkerException("No parser arguments present"));
        if (parserArguments.size() <= index) {
            throw new FreemarkerException("Type " + complexTypeReference.getName() + " specifies too few parser arguments. Available:" + parserArguments.size() + " index:" + index);
        }
        return parserArguments.get(index).getType();
    }

    /**
     * Filters out the arguments that are user for serializiation.
     *
     * @param arguments list of all arguments.
     * @return list of arguments that are used during serialization.
     */
    public List<Argument> getSerializerArguments(List<Argument> arguments) {
        if (arguments == null) {
            return Collections.emptyList();
        }
        List<Argument> serializerArguments = new LinkedList<>();
        for (Argument argument : arguments) {
            if (LAST_ITEM.equals(argument.getName())) {
                serializerArguments.add(argument);
            }
        }
        return serializerArguments;
    }

    public List<Term> getSerializerTerms(List<Term> terms) {
        if (terms == null) {
            return Collections.emptyList();
        }
        List<Term> serializerTerms = new LinkedList<>();
        for (Term term : terms) {
            if (term.contains(LAST_ITEM)) {
                serializerTerms.add(term);
            }
        }
        return serializerTerms;
    }

    public boolean hasLastItemTerm(List<Term> terms) {
        if (terms == null) {
            return false;
        }
        for (Term term : terms) {
            if (term.contains(LAST_ITEM)) {
                return true;
            }
        }
        return false;
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
            boolean needsStringEquals = complexTypeDefinition.getFields().stream()
                .filter(NamedField.class::isInstance)
                .filter(field -> ((NamedField) field).getName().equals(variableLiteral.getName()))
                .filter(TypedField.class::isInstance)
                .map(TypedField.class::cast)
                .map(TypedField::getType)
                .map(StringTypeReference.class::isInstance)
                .findFirst()
                .orElse(false);
            if (needsStringEquals) {
                return true;
            }
        }
        return typeDefinition.getParserArguments()
            .orElse(Collections.emptyList())
            .stream()
            .filter(argument -> argument.getName().equals(variableLiteral.getName()))
            .map(Argument::getType)
            .map(StringTypeReference.class::isInstance)
            .findFirst()
            .orElse(false);
    }

    public boolean isEnumExpression(String expression) {
        if (!expression.contains(".")) {
            return false;
        }
        String enumName = expression.substring(0, expression.indexOf('.'));
        TypeDefinition typeDefinition = this.getTypeDefinitions().get(enumName);
        return (typeDefinition instanceof EnumTypeDefinition);
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

    public Collection<EnumValue> getEnumValuesForConstantValue(List<EnumValue> enumValues, String constantName, String constantValue) {
        List<EnumValue> filteredEnumValues = new ArrayList<>();
        for (EnumValue enumValue : enumValues) {
            if (enumValue.getConstant(constantName).orElseThrow(() -> new FreemarkerException("No constant name " + constantName + " found in enum value" + enumValue)).equals(constantValue)) {
                filteredEnumValues.add(enumValue);
            }
        }
        return filteredEnumValues;
    }

    public SimpleTypeReference getEnumFieldSimpleTypeReference(TypeReference type, String fieldName) {
        TypeDefinition typeDefinition = getTypeDefinitionForTypeReference(type);

        if (typeDefinition instanceof EnumTypeDefinition
            && ((EnumTypeDefinition) typeDefinition).getConstantType(fieldName) instanceof SimpleTypeReference) {
            return (SimpleTypeReference) ((EnumTypeDefinition) typeDefinition).getConstantType(fieldName);
        }
        return null;
    }

    /**
     * Confirms if a variable is an implicit variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @return boolean returns true if the variable's name is an implicit field
     */
    protected boolean isVariableLiteralImplicitField(VariableLiteral variableLiteral) {
        return isVariableLiteralImplicitField(variableLiteral, thisType);
    }

    /**
     * Confirms if a variable is an implicit variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @param typeDefinition  Type definition to check
     * @return boolean returns true if the variable's name is an implicit field
     */
    protected boolean isVariableLiteralImplicitField(VariableLiteral variableLiteral, TypeDefinition typeDefinition) {
        List<Field> fields = null;
        if (typeDefinition instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexType = (ComplexTypeDefinition) typeDefinition;
            fields = complexType.getFields();
        }
        if (fields == null) {
            return false;
        }
        for (Field field : fields) {
            if (field.getTypeName().equals(IMPLICIT)) {
                ImplicitField implicitField = (ImplicitField) field;
                if (variableLiteral.getName().equals(implicitField.getName())) {
                    return true;
                }
            }
        }
        return false;
    }

    /**
     * Confirms if a variable is an virtual variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @return boolean returns true if the variable's name is an virtual field
     */
    protected boolean isVariableLiteralVirtualField(VariableLiteral variableLiteral) {
        return isVariableLiteralVirtualField(variableLiteral, thisType);
    }

    /**
     * Confirms if a variable is an virtual variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @param typeDefinition  Type definition to check
     * @return boolean returns true if the variable's name is an virtual field
     */
    protected boolean isVariableLiteralVirtualField(VariableLiteral variableLiteral, TypeDefinition typeDefinition) {
        List<Field> fields = new ArrayList<>();
        if (typeDefinition instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexType = (ComplexTypeDefinition) typeDefinition;
            fields.addAll(complexType.getFields());
            if (complexType.getParentType() != null) {
                fields.addAll(((ComplexTypeDefinition) complexType.getParentType()).getFields());
            }
        }
        for (Field field : fields) {
            if (field.getTypeName().equals(VIRTUAL)) {
                VirtualField virtualField = (VirtualField) field;
                if (variableLiteral.getName().equals(virtualField.getName())) {
                    return true;
                }
            }
        }
        return false;
    }

    /**
     * Confirms if a variable is a discriminator variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @return boolean returns true if the variable's name is an discriminator field
     */
    protected boolean isVariableLiteralDiscriminatorField(VariableLiteral variableLiteral) {
        return isVariableLiteralDiscriminatorField(variableLiteral, thisType);
    }

    /**
     * Confirms if a variable is an discriminator variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @param typeDefinition  Type definition to check
     * @return boolean returns true if the variable's name is an discriminator field
     */
    protected boolean isVariableLiteralDiscriminatorField(VariableLiteral variableLiteral, TypeDefinition typeDefinition) {
        List<Field> fields = null;
        if (typeDefinition instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexType = (ComplexTypeDefinition) typeDefinition;
            fields = complexType.getFields();
        }
        if (fields == null) {
            return false;
        }
        for (Field field : fields) {
            if (field.getTypeName().equals(DISCRIMINATOR)) {
                DiscriminatorField discriminatorField = (DiscriminatorField) field;
                if (variableLiteral.getName().equals(discriminatorField.getName())) {
                    return true;
                }
            }
        }
        return false;
    }

    /**
     * Returns the implicit field that has the same name as the variable. These need to be handled differently when serializing and parsing.
     *
     * @param variableLiteral The variable to search for.
     * @return ImplicitField returns the implicit field that corresponds to the variable's name.
     */
    protected ImplicitField getReferencedImplicitField(VariableLiteral variableLiteral) {
        return getReferencedImplicitField(variableLiteral, thisType);
    }

    /**
     * Returns the implicit field that has the same name as the variable. These need to be handled differently when serializing and parsing.
     *
     * @param vl             The variable to search for.
     * @param typeDefinition Type definition to check
     * @return ImplicitField returns the implicit field that corresponds to the variable's name.
     */
    protected ImplicitField getReferencedImplicitField(VariableLiteral vl, TypeDefinition typeDefinition) {
        List<Field> fields = null;
        if (typeDefinition instanceof ComplexTypeDefinition) {
            ComplexTypeDefinition complexType = (ComplexTypeDefinition) typeDefinition;
            fields = complexType.getFields();
        }
        if (fields == null) {
            return null;
        }
        for (Field field : fields) {
            if (field.getTypeName().equals(IMPLICIT)) {
                ImplicitField implicitField = (ImplicitField) field;
                if (vl.getName().equals(implicitField.getName())) {
                    return implicitField;
                }
            }
        }
        return null;
    }

}
