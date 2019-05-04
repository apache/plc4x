package org.apache.plc4x.codegen.version2;

public class TypeUtil {

    public static TypeNode infer(Object o) {
        if (o instanceof Double) {
            return Primitive.DOUBLE;
        }
        throw new UnsupportedOperationException("No type available for object " + o);
    }
}
