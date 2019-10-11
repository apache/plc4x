grammar Expression;

// Borrowed part of this Grammar:
// https://github.com/bkiers/tiny-language-antlr4/blob/master/src/main/antlr4/tl/antlr4/TL.g4
// Which is under UNLICENCE:
// https://unlicense.org/
// Which is regarded beint Category A:
// https://issues.apache.org/jira/browse/LEGAL-463

expressionString
 : expression EOF
 ;

expression
 : '-' expression                                       #unaryMinusExpression
 | '!' expression                                       #notExpression
 | <assoc=right> expression '^' expression              #powerExpression
 | expression op=( '*' | '/' | '%' ) expression         #multExpression
 | expression op=( '+' | '-' ) expression               #addExpression
 | expression op=( '>=' | '<=' | '>' | '<' ) expression #compExpression
 | expression op=( '==' | '!=' ) expression             #eqExpression
 | expression '&&' expression                           #andExpression
 | expression '||' expression                           #orExpression
 | expression '?' expression ':' expression             #ifExpression
 | Number                                               #numberExpression
 | Bool                                                 #boolExpression
 | Null                                                 #nullExpression
 | identifierSegment                                    #identifierExpression
 | String indexes?                                      #stringExpression
 | '(' expression ')' indexes?                          #expressionExpression
 ;

identifierSegment
 : name=Identifier args=identifierSegmentArguments? index=identifierSegmentIndexes? ('.' rest=identifierSegmentRest)?
 ;

identifierSegmentArguments
 : arguments
 ;

identifierSegmentIndexes
 : indexes
 ;

identifierSegmentRest
 : identifierSegment
 ;

arguments
 : '(' (expression (',' expression)*)? ')'
 ;

indexes
 : ( '[' expression ']' )+
 ;

Null     : 'null';

Bool
 : 'true'
 | 'false'
 ;

Number
 : Int ( '.' Digit* )?
 ;

Identifier
 : [a-zA-Z_] [a-zA-Z_0-9]*
 ;

String
 : ["] ( ~["\r\n\\] | '\\' ~[\r\n] )* ["]
 | ['] ( ~['\r\n\\] | '\\' ~[\r\n] )* [']
 ;

Space
 : [ \t\r\n\u000C] -> skip
 ;

fragment Int
 : [1-9] Digit*
 | '0'
 ;

fragment Digit
 : [0-9]
 ;