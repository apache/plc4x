grammar imaginary;

file
 : statement*
 ;

statement
 : COMMENT
 | typeDefinition
 ;

typeDefinition
 : LBRACKET K_TYPE IDENTIFIER COMMENT? (fieldDefinition COMMENT?)+ RBRACKET
 | LBRACKET K_DISCRIMINATED_TYPE IDENTIFIER '[length]'? COMMENT? (fieldDefinition COMMENT?)+ RBRACKET
 ;


fieldDefinition
 : LBRACKET simpleField RBRACKET
 | typeSwitch
 ;

simpleField
 : K_CONST dataType IDENTIFIER expression
 | K_RESERVED dataType expression
 | K_IMPLICIT dataType IDENTIFIER expression
 | K_EMBEDDED IDENTIFIER LCBRACKET context RCBRACKET
 | K_DISCRIMINATOR dataType IDENTIFIER
 | 'context' (dataType | IDENTIFIER) expression
 | 'optionalField' (dataType | IDENTIFIER) expression
 | 'field' dataType IDENTIFIER
 | arrayField
 ;

typeSwitch
 : LBRACKET 'typeSwitch' IDENTIFIER COMMENT? caseStatement* RBRACKET
 ;

caseStatement
 : LBRACKET expression IDENTIFIER COMMENT? (fieldDefinition COMMENT?)+ RBRACKET
 ;

arrayField
 : 'arrayField' (IDENTIFIER | dataType) IDENTIFIER IDENTIFIER expression COMMENT*
 ;

dataType
 : 'bit'
 | 'uint7'
 | K_UINT8
 | K_UINT16
 ;

expression
 : TICK innerExpression TICK
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

fragment K_COMMENT : '<--';
K_CONST : 'const';
K_RESERVED : 'reserved';
K_IMPLICIT : 'implicit';
K_EMBEDDED : 'embedded';
K_DISCRIMINATOR : 'discriminator';
K_TYPE : 'type';
K_DISCRIMINATED_TYPE : 'discriminatedType';

K_UINT8 : 'uint8';
K_UINT16 : 'uint16';

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
 ;

ZERO : '0';
HEX_VALUE : [0-9A-F];

IDENTIFIER
 : [A-Za-z0-9_-]+
 ;

WS  :  [ \t\r\n\u000C]+ -> skip
;