# Concepts (open for discussion)

## Requirements

- Maven >= 3.3.1 
- Java >= 8

## General

- All modules should use only one primary build system
- It should be super-easy for new contributors to get started (A new contributor should be able to checkout and build with a core testsuite with a simple: 'mvn package' run)
- New code should only be accepted, if there are tests (Currently the java part of the build is configured to fail if the code coverage is below 90%)

## Java Specific

- Development should be done in Java 8
- Providing Java 7 compatible versions should be possible by using the retrolambda plugin
  - Usage of default implementations does cause more problems than it solves in this case. 
  
 