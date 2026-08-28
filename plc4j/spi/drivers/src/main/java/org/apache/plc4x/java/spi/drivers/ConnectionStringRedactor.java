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
package org.apache.plc4x.java.spi.drivers;

import org.apache.plc4x.java.spi.config.SecretParameters;
import org.apache.plc4x.java.spi.config.annotations.Secret;

import java.net.URLDecoder;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.Locale;
import java.util.LinkedHashSet;
import java.util.Set;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Removes secrets from a connection string before it is logged.
 *
 * <p>Which parameters are secret comes from the {@link Secret} markings on the configuration
 * classes involved, not from matching words against parameter names. The name-matching version
 * that preceded this could only ever be one parameter behind: it lived in {@code DriverBase},
 * modules away from every declaration it protected, and had to be extended by hand each time
 * somebody noticed a value in a log.</p>
 *
 * <p>Credentials in a URI's userinfo ({@code s7://user:secret@plc:102}) have no parameter name to
 * mark, so they are removed structurally - see {@link Secret}'s note on what the marking does not
 * cover.</p>
 *
 * <p><strong>A name-based backstop remains, demoted.</strong> The markings decide what a
 * configuration declares; they cannot decide anything about a parameter no configuration declares.
 * A user who misspells a secret parameter - {@code passwrod=hunter2}, or an {@code api-token} some
 * deployment passes through - has still typed a real credential into a string that is about to be
 * logged, and the marking has nothing to match it against. So a parameter whose name looks like a
 * secret is redacted too. This is no longer the source of truth, only a net under it, and it is
 * why the list may be conservative without being complete.</p>
 */
public final class ConnectionStringRedactor {

    /** What a redacted value is replaced with. Matches what the driver logs elsewhere. */
    private static final String REDACTED = "******";

    /**
     * The credentials of a URI authority: everything between the first {@code :} after the
     * scheme's {@code //} and the {@code @} that ends the userinfo component.
     *
     * <p>The user segment excludes {@code :} so that the <em>first</em> colon separates it from
     * the password. A greedy user segment would let a password containing a colon keep everything
     * up to its last one - {@code bob:pa:ss@} would be redacted to {@code bob:pa:******@}, which
     * publishes half the credential while looking like it hid it.</p>
     */
    private static final Pattern USERINFO = Pattern.compile("(//)([^/@\\s:]*:)([^/@\\s]*)(@)");

    private ConnectionStringRedactor() {
        // Utility class.
    }

    /**
     * The connection string with every secret value replaced.
     *
     * @param connectionString           the string to redact; {@code null} yields {@code null}
     * @param driverConfigurationClass   the driver's configuration class, or {@code null}
     * @param transportCode              the transport's code, which prefixes its parameters
     * @param transportConfigurationClass the transport's configuration class, or {@code null}
     */
    public static String redact(String connectionString, Class<?> driverConfigurationClass,
                                String transportCode, Class<?> transportConfigurationClass) {
        if (connectionString == null) {
            return null;
        }
        Set<String> secretNames = new LinkedHashSet<>();
        for (String name : secretParameterNames(driverConfigurationClass, transportCode, transportConfigurationClass)) {
            secretNames.add(name.toLowerCase(Locale.ROOT));
        }
        return redactParameters(redactUserinfo(connectionString), secretNames, true);
    }

    /**
     * Walks the parameters once, deciding each from its <em>decoded</em> name.
     *
     * <p>Deciding from the name as written would miss a percent-encoded one:
     * {@code ?%70assword=hunter2} is the {@code password} parameter by the time
     * {@code URI.getQuery()} has decoded it and the driver reads it, but no pattern over the raw
     * string sees the word. The value is a real credential either way.</p>
     */
    private static String redactParameters(String connectionString, Set<String> secretNames,
                                           boolean redactNestedValues) {
        Matcher matcher = PARAMETER.matcher(connectionString);
        StringBuilder redacted = new StringBuilder();
        while (matcher.find()) {
            String name = decode(matcher.group(2));
            String value = isSecret(name, secretNames) ? REDACTED
                : redactNestedValues ? redactNested(matcher.group(3)) : matcher.group(3);
            matcher.appendReplacement(redacted,
                Matcher.quoteReplacement(matcher.group(1) + matcher.group(2) + "=" + value));
        }
        matcher.appendTail(redacted);
        return redacted.toString();
    }

    /**
     * A parameter whose value is itself a connection string carries that string's credentials.
     *
     * <p>The PLC4X proxy driver takes a whole PLC URL as {@code remote-connection-string}, encoded
     * so the outer string can hold it. Nothing about the outer parameter's name says "secret", and
     * the encoded value hides the inner {@code password=} from every pattern here - so the inner
     * credential was logged in clear. Redacting the value as a connection string in its own right
     * keeps what an operator needs (which PLC the proxy talks to) and removes what they must not
     * see, which marking the whole parameter secret would not.</p>
     */
    private static String redactNested(String value) {
        String decoded;
        try {
            decoded = URLDecoder.decode(value, StandardCharsets.UTF_8);
        } catch (IllegalArgumentException e) {
            return value;
        }
        if (!decoded.contains("://")) {
            return value;
        }
        // One level: a connection string inside a connection string is the case that exists, and
        // a bound depth cannot be talked into recursing on crafted input.
        String redacted = redactParameters(redactUserinfo(decoded), Set.of(), false);
        if (decoded.equals(redacted)) {
            return value;
        }
        // Rendered back the way it was written, so the log line still looks like the string the
        // user supplied.
        return decoded.equals(value) ? redacted : URLEncoder.encode(redacted, StandardCharsets.UTF_8);
    }

    private static boolean isSecret(String name, Set<String> secretNames) {
        return secretNames.contains(name.toLowerCase(Locale.ROOT))
            || SECRET_LOOKING_NAME.matcher(name).find();
    }

    /** The name as the driver will read it, or as written when it is not valid encoding. */
    private static String decode(String name) {
        try {
            return URLDecoder.decode(name, StandardCharsets.UTF_8);
        } catch (IllegalArgumentException e) {
            return name;
        }
    }

    /**
     * Every parameter name that carries a secret, with the transport's parameters under its code -
     * a transport parameter is addressed as {@code tls.psk-key}, so that is the name to look for.
     */
    static Set<String> secretParameterNames(Class<?> driverConfigurationClass, String transportCode,
                                            Class<?> transportConfigurationClass) {
        Set<String> names = new LinkedHashSet<>(SecretParameters.namesFor(driverConfigurationClass));
        for (String name : SecretParameters.namesFor(transportConfigurationClass)) {
            names.add(((transportCode == null) ? "" : transportCode + ".") + name);
        }
        return names;
    }

    /**
     * Names that read like a credential. Only consulted for parameters the configurations do not
     * declare - see the class comment. {@code psk-identity} is deliberately absent: it says which
     * key was refused, which is the one thing an operator needs when a handshake fails.
     */
    private static final Pattern SECRET_LOOKING_NAME = Pattern.compile(
        "(?i)password|passwd|secret|token|psk-key|passphrase");

    /**
     * One parameter of a connection string: its separator, its name as written, and its value.
     * Names are compared without regard to case - a wrongly-cased parameter does not bind, but
     * the value the user typed is a real secret either way and must not reach the log.
     */
    private static final Pattern PARAMETER = Pattern.compile("([?&])([^=&]*)=([^&]*)");

    private static String redactUserinfo(String connectionString) {
        Matcher matcher = USERINFO.matcher(connectionString);
        return matcher.replaceAll("$1$2" + REDACTED + "$4");
    }
}
