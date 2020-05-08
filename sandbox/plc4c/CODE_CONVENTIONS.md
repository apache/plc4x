# Apache PLC4C Code Conventions

This is not a style guide, but rather a set of conventions to be used when creating variable, functions, and structs from an implementation point of view.
It is also concerned with naming conventions and other patterns to maintain a consistent, maintainable codebase.

## Structural Conventions

- Each module contains `src`, `include` and `test` directories
  - `src`: Sources (*.c files and resources) 
  - `include`: Header files (*.h files)
  - `test`: Test-Sources (and resources)
- All header files should be located in `plc4c` directories inside the `include` directories
- Types and structs that are cross-cutting (used in multiple type-domains) are defined in the `types.h` file
- Functions operating on a given type are declared in the corresponding header file: functions operating on `connections` are defined in `connection.h`

## Naming-Conventions

- All function- and type-names should have the prefix `plc4c_`
- The general pattern for functions should be `plc4c_{basetype}_{operation}`. Some examples:
  - `plc4c_connection_create_read_request`
- If the operations have an `effect` on a `property` of the base-type the `operation` separates the base-type and the property name (`plc4c_{basetype}_{operation}_{property-name}`). 
  - `plc4c_connection_get_connection_string`
  - `plc4c_connection_set_connection_string`
- Structures should all be created by using a function with a `_create` suffix
- Structures should all be freed/cleared/deleted by calling a function with a `_destroy` suffix
