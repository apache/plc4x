package org.apache.plc4x.utils.maven.site.asciidoctor.extensions;

import org.apache.maven.plugin.MojoExecutionException;
import org.asciidoctor.Asciidoctor;
import org.asciidoctor.extension.BlockMacroProcessor;
import org.asciidoctor.extension.BlockProcessor;
import org.asciidoctor.extension.DocinfoProcessor;
import org.asciidoctor.extension.IncludeProcessor;
import org.asciidoctor.extension.InlineMacroProcessor;
import org.asciidoctor.extension.JavaExtensionRegistry;
import org.asciidoctor.extension.Postprocessor;
import org.asciidoctor.extension.Preprocessor;
import org.asciidoctor.extension.Processor;
import org.asciidoctor.extension.Treeprocessor;

import java.lang.reflect.Method;
import java.util.Arrays;

/**
 * Class responsible for registering extensions. This class is inspired by
 * {@link org.asciidoctor.extension.spi.ExtensionRegistry}
 * 
 * @author abelsromero
 * */
public class AsciidoctorJExtensionRegistry implements ExtensionRegistry {

    private JavaExtensionRegistry javaExtensionRegistry;

    public AsciidoctorJExtensionRegistry(Asciidoctor asciidoctorInstance) {
        javaExtensionRegistry = asciidoctorInstance.javaExtensionRegistry();
    }

    /*
     * (non-Javadoc)
     * 
     * @see
     * org.asciidoctor.maven.processors.ProcessorRegistry#register(java.lang.String, java.lang.String)
     */
    @Override
    @SuppressWarnings("unchecked")
    public void register(String extensionClassName, String blockName) throws MojoExecutionException {

        Class<? extends Processor> clazz;
        try {
            clazz = (Class<Processor>) Class.forName(extensionClassName);
        } catch (ClassCastException cce) {
            // Use RuntimeException to avoid catching, we only want the message in the Mojo
            throw new RuntimeException("'" + extensionClassName + "' is not a valid AsciidoctorJ processor class");
        } catch (ClassNotFoundException e) {
            throw new RuntimeException("'" + extensionClassName + "' not found in classpath");
        }

        // TODO: Replace with direct method calls again as soon as this project compiles against AsciidoctorJ 1.6.0
        if (DocinfoProcessor.class.isAssignableFrom(clazz)) {
            register(javaExtensionRegistry, "docinfoProcessor", clazz);
        } else if (Preprocessor.class.isAssignableFrom(clazz)) {
            register(javaExtensionRegistry, "preprocessor", clazz);
        } else if (Postprocessor.class.isAssignableFrom(clazz)) {
            register(javaExtensionRegistry, "postprocessor", clazz);
        } else if (Treeprocessor.class.isAssignableFrom(clazz)) {
            register(javaExtensionRegistry, "treeprocessor", clazz);
        } else if (BlockProcessor.class.isAssignableFrom(clazz)) {
            if (blockName == null) {
                register(javaExtensionRegistry, "block", clazz);
            } else {
                register(javaExtensionRegistry, "block", blockName, clazz);
            }
        } else if (IncludeProcessor.class.isAssignableFrom(clazz)) {
            register(javaExtensionRegistry, "includeProcessor", clazz);
        } else if (BlockMacroProcessor.class.isAssignableFrom(clazz)) {
            if (blockName == null) {
                register(javaExtensionRegistry, "blockMacro", clazz);
            } else {
                register(javaExtensionRegistry, "blockMacro", blockName, clazz);
            }
        } else if (InlineMacroProcessor.class.isAssignableFrom(clazz)) {
            if (blockName == null) {
                register(javaExtensionRegistry, "inlineMacro", clazz);
            } else {
                register(javaExtensionRegistry, "inlineMacro", blockName, clazz);
            }
        }
    }

    private void register(Object target, String methodName, Object... args) throws MojoExecutionException {
        for (Method method: javaExtensionRegistry.getClass().getMethods()) {

            if (isMethodMatching(method, methodName, args)) {
                try {
                    method.invoke(target, args);
                    return;
                } catch (Exception e) {
                    throw new MojoExecutionException("Unexpected exception while registering extensions", e);
                }
            }

        }
        throw new MojoExecutionException("Internal Error. Could not register " + methodName + " with arguments " + Arrays.asList(args));
    }

    private boolean isMethodMatching(Method method, String methodName, Object[] args) {
        if (!method.getName().equals(methodName)) {
            return false;
        }
        if (method.getParameterTypes().length != args.length) {
            return false;
        }
        // Don't care for primitives here, there's no method on JavaExtensionRegistry with primitives.
        for (int i = 0; i < method.getParameterTypes().length; i++) {
            if (args[i] != null && !method.getParameterTypes()[i].isAssignableFrom(args[i].getClass())) {
                return false;
            }
        }
        return true;
    }

}
