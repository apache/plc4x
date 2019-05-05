package org.apache.plc4x.codegen.util;

import org.apache.plc4x.codegen.ast.Expressions;
import org.apache.plc4x.codegen.ast.Method;
import org.apache.plc4x.codegen.ast.Primitive;
import org.apache.plc4x.codegen.ast.TypeDefinition;

import java.util.Collections;

/**
 * This class defines constants necessary for the code generation related to the
 * "Buffer API" which has to be implemented natively.
 */
public class BufferUtil {

    static final TypeDefinition BUFFER_TYPE = Expressions.typeOf("org.apache.plc4x.codegen.api.Buffer");

    // Read Methods
    static final Method READ_UINT8 = new Method(BUFFER_TYPE, "readUint8", Primitive.INTEGER, Collections.emptyList(), Collections.emptyList());
    static final Method READ_UINT16 = new Method(BUFFER_TYPE, "readUint16", Primitive.INTEGER, Collections.emptyList(), Collections.emptyList());
    static final Method READ_UINT32 = new Method(BUFFER_TYPE, "readUint32", Primitive.LONG, Collections.emptyList(), Collections.emptyList());

    private BufferUtil() {
    }


}
