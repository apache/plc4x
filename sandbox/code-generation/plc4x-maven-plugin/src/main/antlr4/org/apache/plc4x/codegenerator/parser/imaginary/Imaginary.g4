grammar Imaginary;

file
 : complexTypeDefinition* EOF
 ;

complexTypeDefinition
 : COMMENT
 | LBRACKET complexType RBRACKET
 ;

complexType
 : K_TYPE name=idExpression fieldDefinition+
 | K_DISCRIMINATED_TYPE name=idExpression (LBRACKET params=multipleExpressions RBRACKET)? fieldDefinition+
 ;


fieldDefinition
 : LBRACKET field (LCBRACKET context RCBRACKET)? RBRACKET
 ;

field
 : arrayField
 | constField
 | contextField
 | discriminatorField
 | embeddedField
 | simpleField
 | implicitField
 | optionalField
 | reservedField
 | typeSwitchField
 ;

arrayField
 : K_ARRAY type=typeReference name=idExpression lengthType=arrayType lengthExpression=expression
 ;

constField
 : K_CONST type=dataType name=idExpression expected=expression
 ;

contextField
 : K_CONTEXT type=dataType name=idExpression expression
 ;

discriminatorField
 : K_DISCRIMINATOR type=dataType name=idExpression
 ;

embeddedField
 : K_EMBEDDED name=idExpression LCBRACKET context RCBRACKET
 ;

simpleField
 : K_FIELD type=typeReference name=idExpression (LCBRACKET context RCBRACKET)?
 ;

implicitField
 : K_IMPLICIT type=dataType name=idExpression serializationExpression=expression
 ;

optionalField
 : K_OPTIONAL_FIELD type=typeReference name=idExpression condition=expression
 ;

reservedField
 : K_RESERVED type=dataType expected=expression
 ;

typeSwitchField
 : K_TYPE_SWITCH discriminators=multipleExpressions caseStatement*
 ;


typeReference
 : complexTypeReference=IDENTIFIER
 | simpleTypeReference=dataType
 ;

caseStatement
 : LBRACKET discriminatorValues=multipleExpressions name=IDENTIFIER fieldDefinition* RBRACKET
 ;

dataType
 : base=K_BIT
 | base=K_INT size=INTEGER_LITERAL
 | base=K_UINT size=INTEGER_LITERAL
 | base=K_FLOAT size=INTEGER_LITERAL
 | base=K_STRING
 ;

expression
 : TICK expr=innerExpression TICK
 ;

multipleExpressions
 : expression (',' expression)*
 ;

innerExpression
 : HEX_LITERAL
 | INTEGER_LITERAL
 | IDENTIFIER   // Variable reference
 | innerExpression '.' IDENTIFIER // Field Reference
 | innerExpression BinaryOperator innerExpression  // Addition
 | '(' innerExpression ')'
 ;

context
 : IDENTIFIER ':' expression (',' IDENTIFIER ':' expression)*
 ;

COMMENT
 : K_COMMENT [a-zA-Z0-9,.'":()/ \t\r\n\u000C-]*
 | '//' [a-zA-Z0-9,.'":()/ \t-]*
 ;

INTEGER_LITERAL
 : [0-9]+
 ;

HEX_LITERAL
 : HexNumeral
 ;

fragment HexNumeral
 : '0' [xX] HexDigit HexDigit?;

fragment HexDigit
 : [0-9a-fA-F]
;

modifier
 : K_CONST
 | K_RESERVED
 | K_IMPLICIT
 | K_EMBEDDED
 ;

arrayType
 : K_COUNT
 | K_LENGTH
 ;

idExpression
 : TICK id=IDENTIFIER TICK
 ;

fragment K_COMMENT : '<--';
K_ARRAY : 'arrayField';
K_CONST : 'const';
K_CONTEXT : 'context';
K_DISCRIMINATED_TYPE : 'discriminatedType';
K_DISCRIMINATOR : 'discriminator';
K_EMBEDDED : 'embedded';
K_FIELD : 'field';
K_IMPLICIT : 'implicit';
K_OPTIONAL_FIELD : 'optionalField';
K_RESERVED : 'reserved';
K_TYPE : 'type';
K_TYPE_SWITCH : 'typeSwitch';

K_COUNT : 'count';
K_LENGTH : 'length';

K_BIT : 'bit';
K_INT : 'int';
K_UINT : 'uint';
K_FLOAT : 'float';
K_UINT5 : 'uint5';
K_UINT8 : 'uint8';
K_UINT16 : 'uint16';
K_STRING: 'string';

TICK : '\'';
TIMES : 'x';
LBRACKET : '[';
RBRACKET : ']';
LCBRACKET : '{';
RCBRACKET : '}';

BinaryOperator
 : '+'
 | '-'
 | '=='
 | '!='
 | '>'
 | '<'
 | '>='
 | '<='
 ;

ZERO : '0';
HEX_VALUE : [0-9A-F];

IDENTIFIER
 : [A-Za-z0-9_-]+
 ;

WS  :  [ \t\r\n\u000C]+ -> skip
;