# CHANGELOG

## [rel/0.12](https://github.com/apache/plc4x/releases/tag/rel/0.12) - 2024-02-16 11:48:19

## What's Changed
* build(deps): bump com.google.googlejavaformat:google-java-format from 1.17.0 to 1.18.0 by @dependabot in https://github.com/apache/plc4x/pull/1126
* build(deps): bump org.checkerframework:checker-qual from 3.38.0 to 3.39.0 by @dependabot in https://github.com/apache/plc4x/pull/1127
* build(deps): bump com.hivemq:hivemq-mqtt-client from 1.3.2 to 1.3.3 by @dependabot in https://github.com/apache/plc4x/pull/1129
* build(deps): bump com.google.protobuf:protobuf-java from 3.24.3 to 3.24.4 by @dependabot in https://github.com/apache/plc4x/pull/1131
* build(deps): bump github.com/spf13/viper from 1.16.0 to 1.17.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1135
* build(deps): bump mockito.version from 5.5.0 to 5.6.0 by @dependabot in https://github.com/apache/plc4x/pull/1137
* fix(opcua): Attempting to fix cyclic subscriptions by @hutcheb in https://github.com/apache/plc4x/pull/1124
* build(deps): bump golang.org/x/net from 0.15.0 to 0.17.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1140
* build(deps): bump net.bytebuddy:byte-buddy from 1.14.8 to 1.14.9 by @dependabot in https://github.com/apache/plc4x/pull/1143
* build(deps): bump org.eclipse.jetty:jetty-util from 12.0.1 to 12.0.2 by @dependabot in https://github.com/apache/plc4x/pull/1142
* build(deps): bump netty.version from 4.1.99.Final to 4.1.100.Final by @dependabot in https://github.com/apache/plc4x/pull/1141
* build(deps): bump io.jsonwebtoken:jjwt from 0.9.1 to 0.12.2 by @dependabot in https://github.com/apache/plc4x/pull/1132
* build(deps): bump com.google.googlejavaformat:google-java-format from 1.18.0 to 1.18.1 by @dependabot in https://github.com/apache/plc4x/pull/1133
* build(deps): bump golang.org/x/tools from 0.13.0 to 0.14.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1136
* plc4j-driver-opcua: Fix keepalive threads are never shut down by @takraj in https://github.com/apache/plc4x/pull/1139
* build(deps): bump jackson.version from 2.15.2 to 2.15.3 by @dependabot in https://github.com/apache/plc4x/pull/1145
* build(deps): bump com.google.guava:guava from 32.1.2-jre to 32.1.3-jre by @dependabot in https://github.com/apache/plc4x/pull/1146
* plc4j-driver-opcua: Move back to using session lifetime * 75% as keepalive period by @takraj in https://github.com/apache/plc4x/pull/1148
* Fix race condition causing `RejectedExecutionException` on `PlcConnection.close()` by @takraj in https://github.com/apache/plc4x/pull/1151
* build(deps): bump io.jsonwebtoken:jjwt-api from 0.12.2 to 0.12.3 by @dependabot in https://github.com/apache/plc4x/pull/1152
* build(deps-dev): bump org.json:json from 20230618 to 20231013 by @dependabot in https://github.com/apache/plc4x/pull/1153
* build(deps): bump io.swagger:swagger-annotations from 1.6.11 to 1.6.12 by @dependabot in https://github.com/apache/plc4x/pull/1154
* build(deps): bump org.jacoco:jacoco-maven-plugin from 0.8.10 to 0.8.11 by @dependabot in https://github.com/apache/plc4x/pull/1155
* plc4j-driver-opcua: Await `writeAndFlush(msg)` & send next msg async by @takraj in https://github.com/apache/plc4x/pull/1147
* OPC UA priority judgment using discovery parameter by @qtvbwfn in https://github.com/apache/plc4x/pull/1157
* build(deps): bump com.gradle:gradle-enterprise-maven-extension from 1.19.2 to 1.19.3 by @dependabot in https://github.com/apache/plc4x/pull/1158
* plc4j-driver-opcua: Re-enable a disabled test by @takraj in https://github.com/apache/plc4x/pull/1159
* plc4x-server: Fix logs not being written to screen by @takraj in https://github.com/apache/plc4x/pull/1161
* refactor(plc4x-server): Cleanup, add tests, and add option to specify port number by @takraj in https://github.com/apache/plc4x/pull/1162
* test(plc4j/opcua): Cleanup testcase manySubscriptionsOnSingleConnection by @takraj in https://github.com/apache/plc4x/pull/1160
* build(deps): bump com.google.errorprone:error_prone_annotations from 2.22.0 to 2.23.0 by @dependabot in https://github.com/apache/plc4x/pull/1164
* build(deps): bump org.jsoup:jsoup from 1.16.1 to 1.16.2 by @dependabot in https://github.com/apache/plc4x/pull/1169
* feat(plc4x-server): Build a standalone jar too by @takraj in https://github.com/apache/plc4x/pull/1167
* feat(plc4j/drivers): Create maven meta package to include all drivers by @takraj in https://github.com/apache/plc4x/pull/1166
* build(deps): bump org.codehaus.mojo:properties-maven-plugin from 1.2.0 to 1.2.1 by @dependabot in https://github.com/apache/plc4x/pull/1173
* build(deps): bump com.microsoft.azure.sdk.iot:iot-device-client from 2.3.0 to 2.3.1 by @dependabot in https://github.com/apache/plc4x/pull/1171
* build(deps): bump org.glassfish.jaxb:jaxb-runtime from 4.0.2 to 4.0.4 by @dependabot in https://github.com/apache/plc4x/pull/1170
* feat(plc4j/spi): Add option to synchronously await response from PLC by @takraj in https://github.com/apache/plc4x/pull/1163
* build(deps): bump github.com/google/uuid from 1.3.1 to 1.4.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1176
* feat: Moved the code-generation into a separate profile by @chrisdutz in https://github.com/apache/plc4x/pull/1172
* build(deps): bump org.cyclonedx:cyclonedx-maven-plugin from 2.7.9 to 2.7.10 by @dependabot in https://github.com/apache/plc4x/pull/1178
* Feature/transport configuration rework by @chrisdutz in https://github.com/apache/plc4x/pull/1179
* build(deps): bump org.eclipse.jetty:jetty-util from 12.0.2 to 12.0.3 by @dependabot in https://github.com/apache/plc4x/pull/1180
* build(deps): bump nl.jqno.equalsverifier:equalsverifier from 3.15.2 to 3.15.3 by @dependabot in https://github.com/apache/plc4x/pull/1184
* build(deps): bump com.google.protobuf:protobuf-java from 3.24.4 to 3.25.0 by @dependabot in https://github.com/apache/plc4x/pull/1183
* build(deps): bump org.checkerframework:checker-qual from 3.39.0 to 3.40.0 by @dependabot in https://github.com/apache/plc4x/pull/1182
* build(deps): bump mockito.version from 5.6.0 to 5.7.0 by @dependabot in https://github.com/apache/plc4x/pull/1185
* build(deps): bump github.com/fatih/color from 1.15.0 to 1.16.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1194
* build(deps): bump com.fazecast:jSerialComm from 2.10.3 to 2.10.4 by @dependabot in https://github.com/apache/plc4x/pull/1192
* build(deps): bump github.com/spf13/cobra from 1.7.0 to 1.8.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1189
* build(deps): bump com.googlecode.cmake-maven-project:cmake-maven-plugin from 3.26.3-b1 to 3.27.7-b1 by @dependabot in https://github.com/apache/plc4x/pull/1188
* build(deps): bump github.com/schollz/progressbar/v3 from 3.13.1 to 3.14.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1193
* build(deps): bump junit.jupiter.version from 5.10.0 to 5.10.1 by @dependabot in https://github.com/apache/plc4x/pull/1187
* build(deps): bump com.microsoft.azure.sdk.iot:iot-device-client from 2.3.1 to 2.4.0 by @dependabot in https://github.com/apache/plc4x/pull/1195
* build(deps): bump golang.org/x/net from 0.17.0 to 0.18.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1197
* build(deps): bump netty.version from 4.1.100.Final to 4.1.101.Final by @dependabot in https://github.com/apache/plc4x/pull/1201
* build(deps): bump github.com/schollz/progressbar/v3 from 3.14.0 to 3.14.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1200
* build(deps): bump golang.org/x/tools from 0.14.0 to 0.15.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1196
* upgrade Felix maven-bundle-plugin by @hboutemy in https://github.com/apache/plc4x/pull/1202
* build(deps): bump actions/github-script from 6.4.1 to 7.0.0 by @dependabot in https://github.com/apache/plc4x/pull/1205
* build(deps): bump bouncycastle.version from 1.76 to 1.77 by @dependabot in https://github.com/apache/plc4x/pull/1208
* build(deps): bump com.google.protobuf:protobuf-java from 3.25.0 to 3.25.1 by @dependabot in https://github.com/apache/plc4x/pull/1210
* build(deps): bump jackson.version from 2.15.3 to 2.16.0 by @dependabot in https://github.com/apache/plc4x/pull/1209
* build(deps): bump org.codehaus.mojo:exec-maven-plugin from 3.1.0 to 3.1.1 by @dependabot in https://github.com/apache/plc4x/pull/1211
* Issue/s7h by @glcj in https://github.com/apache/plc4x/pull/1214
* build(deps): bump actions/github-script from 7.0.0 to 7.0.1 by @dependabot in https://github.com/apache/plc4x/pull/1216
* build(deps): bump net.bytebuddy:byte-buddy from 1.14.9 to 1.14.10 by @dependabot in https://github.com/apache/plc4x/pull/1218
* fix(plc4j/spi) Make sure OPC UA discover event is fired prior connected event by @splatch in https://github.com/apache/plc4x/pull/1217
* Issue/s7h v2 by @glcj in https://github.com/apache/plc4x/pull/1219
* build(deps): bump org.codehaus.mojo:build-helper-maven-plugin from 3.4.0 to 3.5.0 by @dependabot in https://github.com/apache/plc4x/pull/1222
* build(deps): bump org.jsoup:jsoup from 1.16.2 to 1.17.1 by @dependabot in https://github.com/apache/plc4x/pull/1223
* build(deps): bump github.com/gopacket/gopacket from 1.1.1 to 1.2.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1224
* build(deps): bump logback.version from 1.4.11 to 1.4.12 by @dependabot in https://github.com/apache/plc4x/pull/1228
* build(deps): bump nifi.version from 1.23.2 to 1.24.0 by @dependabot in https://github.com/apache/plc4x/pull/1227
* build(deps): bump golang.org/x/net from 0.18.0 to 0.19.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1226
* build(deps): bump golang.org/x/tools from 0.15.0 to 0.16.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1225
* build(deps): bump logback.version from 1.4.12 to 1.4.13 by @dependabot in https://github.com/apache/plc4x/pull/1229
* build(deps): bump actions/setup-java from 3 to 4 by @dependabot in https://github.com/apache/plc4x/pull/1231
* build(deps): bump nl.jqno.equalsverifier:equalsverifier from 3.15.3 to 3.15.4 by @dependabot in https://github.com/apache/plc4x/pull/1232
* build(deps): bump net.sf.saxon:Saxon-HE from 12.3 to 12.4 by @dependabot in https://github.com/apache/plc4x/pull/1233
* build(deps): bump tel.schich:javacan-core from 3.2.4 to 3.3.0 by @dependabot in https://github.com/apache/plc4x/pull/1238
* build(deps): bump mockito.version from 5.7.0 to 5.8.0 by @dependabot in https://github.com/apache/plc4x/pull/1237
* build(deps): bump logback.version from 1.4.13 to 1.4.14 by @dependabot in https://github.com/apache/plc4x/pull/1236
* build(deps): bump groovy.version from 4.0.12 to 4.0.16 by @dependabot in https://github.com/apache/plc4x/pull/1239
* build(deps): bump commons-cli:commons-cli from 1.5.0 to 1.6.0 by @dependabot in https://github.com/apache/plc4x/pull/1241
* build(deps): bump org.checkerframework:checker-qual from 3.40.0 to 3.41.0 by @dependabot in https://github.com/apache/plc4x/pull/1242
* build(deps-dev): bump org.apache.commons:commons-compress from 1.23.0 to 1.25.0 by @dependabot in https://github.com/apache/plc4x/pull/1243
* build(deps): bump iotdb.version from 0.13.0 to 1.2.2 by @dependabot in https://github.com/apache/plc4x/pull/1240
* build(deps): bump com.gradle:common-custom-user-data-maven-extension from 1.12.4 to 1.12.5 by @dependabot in https://github.com/apache/plc4x/pull/1248
* build(deps): bump com.gradle:gradle-enterprise-maven-extension from 1.19.3 to 1.20 by @dependabot in https://github.com/apache/plc4x/pull/1246
* build(deps): bump org.apache.karaf.tooling:karaf-maven-plugin from 4.4.3 to 4.4.4 by @dependabot in https://github.com/apache/plc4x/pull/1247
* build(deps): bump org.apache.commons:commons-pool2 from 2.11.1 to 2.12.0 by @dependabot in https://github.com/apache/plc4x/pull/1245
* build(deps): bump org.apache.maven.plugins:maven-source-plugin from 3.2.1 to 3.3.0 by @dependabot in https://github.com/apache/plc4x/pull/1244
* build(deps): bump actions/setup-python from 4 to 5 by @dependabot in https://github.com/apache/plc4x/pull/1250
* build(deps): bump github.com/spf13/viper from 1.17.0 to 1.18.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1251
* build(deps): bump org.apache.commons:commons-text from 1.10.0 to 1.11.0 by @dependabot in https://github.com/apache/plc4x/pull/1252
* build(deps): bump calcite-core.version from 1.34.0 to 1.36.0 by @dependabot in https://github.com/apache/plc4x/pull/1253
* build(deps): bump org.apache.avro:avro from 1.11.1 to 1.11.3 by @dependabot in https://github.com/apache/plc4x/pull/1254
* build(deps): bump org.apache.maven.plugins:maven-surefire-plugin from 3.1.0 to 3.2.2 by @dependabot in https://github.com/apache/plc4x/pull/1256
* Cleanup of various handling of S7 types by @chrisdutz in https://github.com/apache/plc4x/pull/1259
* build(deps-dev): bump commons-io:commons-io from 2.11.0 to 2.15.1 by @dependabot in https://github.com/apache/plc4x/pull/1260
* build(deps): bump org.apache.maven.plugins:maven-release-plugin from 3.0.0 to 3.0.1 by @dependabot in https://github.com/apache/plc4x/pull/1261
* build(deps): bump org.apache.maven.plugins:maven-enforcer-plugin from 3.3.0 to 3.4.1 by @dependabot in https://github.com/apache/plc4x/pull/1262
* build(deps): bump commons-net:commons-net from 3.9.0 to 3.10.0 by @dependabot in https://github.com/apache/plc4x/pull/1263
* build(deps-dev): bump org.apache.commons:commons-lang3 from 3.12.0 to 3.14.0 by @dependabot in https://github.com/apache/plc4x/pull/1264
* build(deps): bump com.influxdb:influxdb-client-java from 6.10.0 to 6.11.0 by @dependabot in https://github.com/apache/plc4x/pull/1266
* build(deps): bump org.apache.maven.plugins:maven-assembly-plugin from 3.5.0 to 3.6.0 by @dependabot in https://github.com/apache/plc4x/pull/1267
* build(deps): bump org.eclipse.jetty:jetty-util from 12.0.3 to 12.0.4 by @dependabot in https://github.com/apache/plc4x/pull/1268
* build(deps): bump org.apache.maven.plugins:maven-failsafe-plugin from 3.1.0 to 3.2.2 by @dependabot in https://github.com/apache/plc4x/pull/1269
* build(deps): bump github.com/spf13/viper from 1.18.0 to 1.18.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1271
* build(deps): bump github.com/gdamore/tcell/v2 from 2.6.0 to 2.7.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1270
* build(deps): bump org.apache.maven.plugins:maven-javadoc-plugin from 3.5.0 to 3.6.3 by @dependabot in https://github.com/apache/plc4x/pull/1272
* build(deps): bump org.apache.maven:maven-core from 3.9.3 to 3.9.6 by @dependabot in https://github.com/apache/plc4x/pull/1274
* build(deps): bump net.java.dev.jna:jna from 5.13.0 to 5.14.0 by @dependabot in https://github.com/apache/plc4x/pull/1275
* build(deps): bump org.apache.maven.plugins:maven-invoker-plugin from 3.5.1 to 3.6.0 by @dependabot in https://github.com/apache/plc4x/pull/1276
* build(deps): bump github.com/google/uuid from 1.4.0 to 1.5.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1283
* build(deps): bump golang.org/x/tools from 0.16.0 to 0.16.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1282
* build(deps): bump commons-logging:commons-logging from 1.2 to 1.3.0 by @dependabot in https://github.com/apache/plc4x/pull/1280
* build(deps): bump org.apache.maven.plugins:maven-remote-resources-plugin from 3.0.0 to 3.1.0 by @dependabot in https://github.com/apache/plc4x/pull/1279
* build(deps): bump netty.version from 4.1.101.Final to 4.1.102.Final by @dependabot in https://github.com/apache/plc4x/pull/1277
* build(deps): bump org.apache.maven.plugins:maven-surefire-plugin from 3.2.2 to 3.2.3 by @dependabot in https://github.com/apache/plc4x/pull/1289
* build(deps): bump org.codehaus.plexus:plexus-compiler-api from 2.13.0 to 2.14.1 by @dependabot in https://github.com/apache/plc4x/pull/1287
* build(deps): bump netty.version from 4.1.102.Final to 4.1.103.Final by @dependabot in https://github.com/apache/plc4x/pull/1285
* build(deps): bump github/codeql-action from 2 to 3 by @dependabot in https://github.com/apache/plc4x/pull/1284
* build(deps): bump org.apache.maven.plugins:maven-failsafe-plugin from 3.2.2 to 3.2.3 by @dependabot in https://github.com/apache/plc4x/pull/1286
* build(deps): bump actions/upload-artifact from 3 to 4 by @dependabot in https://github.com/apache/plc4x/pull/1293
* build(deps): bump com.influxdb:influxdb-client-java from 6.11.0 to 6.12.0 by @dependabot in https://github.com/apache/plc4x/pull/1294
* build(deps): bump netty.version from 4.1.103.Final to 4.1.104.Final by @dependabot in https://github.com/apache/plc4x/pull/1295
* build(deps): bump org.checkerframework:checker-qual from 3.41.0 to 3.42.0 by @dependabot in https://github.com/apache/plc4x/pull/1296
* build(deps): bump github.com/spf13/viper from 1.18.1 to 1.18.2 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1300
* build(deps): bump com.google.googlejavaformat:google-java-format from 1.18.1 to 1.19.0 by @dependabot in https://github.com/apache/plc4x/pull/1299
* build(deps): bump com.google.guava:guava from 32.1.3-jre to 33.0.0-jre by @dependabot in https://github.com/apache/plc4x/pull/1298
* build(deps): bump org.apache.maven.plugins:maven-compiler-plugin from 3.11.0 to 3.12.0 by @dependabot in https://github.com/apache/plc4x/pull/1297
* build(deps): bump org.eclipse.jetty:jetty-util from 12.0.4 to 12.0.5 by @dependabot in https://github.com/apache/plc4x/pull/1301
* build(deps): bump com.google.googlejavaformat:google-java-format from 1.19.0 to 1.19.1 by @dependabot in https://github.com/apache/plc4x/pull/1302
* build(deps): bump org.codehaus.plexus:plexus-compiler-api from 2.14.1 to 2.14.2 by @dependabot in https://github.com/apache/plc4x/pull/1303
* build(deps): bump net.bytebuddy:byte-buddy from 1.14.10 to 1.14.11 by @dependabot in https://github.com/apache/plc4x/pull/1304
* build(deps): bump com.google.errorprone:error_prone_annotations from 2.23.0 to 2.24.0 by @dependabot in https://github.com/apache/plc4x/pull/1310
* build(deps): bump org.asciidoctor:asciidoctorj from 2.5.10 to 2.5.11 by @dependabot in https://github.com/apache/plc4x/pull/1309
* build(deps): bump org.jetbrains.kotlin:kotlin-stdlib-jdk8 from 1.9.21 to 1.9.22 by @dependabot in https://github.com/apache/plc4x/pull/1308
* build(deps): bump kotlin.version from 1.9.21 to 1.9.22 by @dependabot in https://github.com/apache/plc4x/pull/1307
* build(deps): bump groovy.version from 4.0.16 to 4.0.17 by @dependabot in https://github.com/apache/plc4x/pull/1306
* build(deps): bump org.apache.maven.plugins:maven-compiler-plugin from 3.12.0 to 3.12.1 by @dependabot in https://github.com/apache/plc4x/pull/1315
* build(deps): bump nl.jqno.equalsverifier:equalsverifier from 3.15.4 to 3.15.5 by @dependabot in https://github.com/apache/plc4x/pull/1314
* build(deps): bump jackson.version from 2.16.0 to 2.16.1 by @dependabot in https://github.com/apache/plc4x/pull/1313
* build(deps): bump com.googlecode.maven-download-plugin:download-maven-plugin from 1.7.1 to 1.8.0 by @dependabot in https://github.com/apache/plc4x/pull/1312
* build(deps): bump com.fasterxml.jackson.datatype:jackson-datatype-jsr310 from 2.16.0 to 2.16.1 by @dependabot in https://github.com/apache/plc4x/pull/1311
* build(deps): bump tel.schich:javacan-core from 3.3.0 to 3.3.2 by @dependabot in https://github.com/apache/plc4x/pull/1316
* build(deps): bump slf4j.version from 2.0.9 to 2.0.10 by @dependabot in https://github.com/apache/plc4x/pull/1317
* build(deps): bump org.jsoup:jsoup from 1.17.1 to 1.17.2 by @dependabot in https://github.com/apache/plc4x/pull/1318
* build(deps): bump org.asciidoctor:asciidoctorj-diagram from 2.2.13 to 2.2.14 by @dependabot in https://github.com/apache/plc4x/pull/1319
* build(deps): bump org.assertj:assertj-core from 3.24.2 to 3.25.0 by @dependabot in https://github.com/apache/plc4x/pull/1320
* build(deps): bump org.assertj:assertj-core from 3.25.0 to 3.25.1 by @dependabot in https://github.com/apache/plc4x/pull/1322
* build(deps): bump com.google.errorprone:error_prone_annotations from 2.24.0 to 2.24.1 by @dependabot in https://github.com/apache/plc4x/pull/1323
* feat(plc4py): Code Gen Update by @hutcheb in https://github.com/apache/plc4x/pull/1199
* build(deps): bump joda-time:joda-time from 2.12.5 to 2.12.6 by @dependabot in https://github.com/apache/plc4x/pull/1325
* build(deps): bump com.google.googlejavaformat:google-java-format from 1.19.1 to 1.19.2 by @dependabot in https://github.com/apache/plc4x/pull/1324
* Team/glcj by @glcj in https://github.com/apache/plc4x/pull/1326
* build(deps): bump org.apache.rat:apache-rat-plugin from 0.15 to 0.16 by @dependabot in https://github.com/apache/plc4x/pull/1329
* build(deps): bump golang.org/x/net from 0.19.0 to 0.20.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1327
* build(deps): bump slf4j.version from 2.0.10 to 2.0.11 by @dependabot in https://github.com/apache/plc4x/pull/1328
* build(deps): bump org.apache.karaf.tooling:karaf-maven-plugin from 4.4.4 to 4.4.5 by @dependabot in https://github.com/apache/plc4x/pull/1334
* build(deps): bump org.apache.maven.plugins:maven-surefire-plugin from 3.2.3 to 3.2.5 by @dependabot in https://github.com/apache/plc4x/pull/1333
* build(deps): bump org.apache.maven.plugins:maven-failsafe-plugin from 3.2.3 to 3.2.5 by @dependabot in https://github.com/apache/plc4x/pull/1332
* build(deps): bump nl.jqno.equalsverifier:equalsverifier from 3.15.5 to 3.15.6 by @dependabot in https://github.com/apache/plc4x/pull/1331
* build(deps): bump golang.org/x/tools from 0.16.1 to 0.17.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1337
* build(deps): bump com.microsoft.azure.sdk.iot:iot-device-client from 2.4.0 to 2.4.1 by @dependabot in https://github.com/apache/plc4x/pull/1336
* build(deps): bump com.google.protobuf:protobuf-java from 3.25.1 to 3.25.2 by @dependabot in https://github.com/apache/plc4x/pull/1335
* build(deps): bump mockito.version from 5.8.0 to 5.9.0 by @dependabot in https://github.com/apache/plc4x/pull/1341
* build(deps): bump org.asciidoctor:asciidoctor-maven-plugin from 2.2.4 to 2.2.5 by @dependabot in https://github.com/apache/plc4x/pull/1340
* build(deps): bump org.cyclonedx:cyclonedx-maven-plugin from 2.7.10 to 2.7.11 by @dependabot in https://github.com/apache/plc4x/pull/1345
* build(deps): bump netty.version from 4.1.104.Final to 4.1.105.Final by @dependabot in https://github.com/apache/plc4x/pull/1343
* build(deps): bump actions/cache from 3 to 4 by @dependabot in https://github.com/apache/plc4x/pull/1347
* build(deps): bump actions/dependency-review-action from 3 to 4 by @dependabot in https://github.com/apache/plc4x/pull/1348
* build(deps): bump groovy.version from 4.0.17 to 4.0.18 by @dependabot in https://github.com/apache/plc4x/pull/1349
* Feature/new UI tool by @chrisdutz in https://github.com/apache/plc4x/pull/1350
* build(deps): bump netty.version from 4.1.105.Final to 4.1.106.Final by @dependabot in https://github.com/apache/plc4x/pull/1353
* build(deps): bump com.github.eirslett:frontend-maven-plugin from 1.14.2 to 1.15.0 by @dependabot in https://github.com/apache/plc4x/pull/1352
* build(deps): bump org.aspectj:aspectjweaver from 1.9.20.1 to 1.9.21 by @dependabot in https://github.com/apache/plc4x/pull/1357
* build(deps): bump org.springframework.boot:spring-boot-maven-plugin from 3.1.2 to 3.2.2 by @dependabot in https://github.com/apache/plc4x/pull/1355
* build(deps): bump org.springframework.boot:spring-boot-dependencies from 3.2.1 to 3.2.2 by @dependabot in https://github.com/apache/plc4x/pull/1360
* build(deps): bump github.com/google/uuid from 1.5.0 to 1.6.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1358
* build(deps): bump mockito.version from 5.9.0 to 5.10.0 by @dependabot in https://github.com/apache/plc4x/pull/1365
* build(deps): bump org.apache.rat:apache-rat-plugin from 0.15 to 0.16 by @dependabot in https://github.com/apache/plc4x/pull/1354
* chore: update notice year to 2024 by @shoothzj in https://github.com/apache/plc4x/pull/1351
* build(deps): bump org.codehaus.mojo:license-maven-plugin from 2.3.0 to 2.4.0 by @dependabot in https://github.com/apache/plc4x/pull/1368
* build(deps): bump io.swagger:swagger-annotations from 1.6.12 to 1.6.13 by @dependabot in https://github.com/apache/plc4x/pull/1367
* build(deps): bump io.jsonwebtoken:jjwt-api from 0.12.3 to 0.12.4 by @dependabot in https://github.com/apache/plc4x/pull/1370
* build(deps): bump org.apache.rat:apache-rat-plugin from 0.15 to 0.16.1 by @dependabot in https://github.com/apache/plc4x/pull/1369
* remove unnecessary Thread.sleep by @schaebo in https://github.com/apache/plc4x/pull/1374
* build(deps): bump release-drafter/release-drafter from 5 to 6 by @dependabot in https://github.com/apache/plc4x/pull/1377
* build(deps): bump github.com/rs/zerolog from 1.31.0 to 1.32.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1378
* build(deps): bump com.influxdb:influxdb-client-java from 6.12.0 to 7.0.0 by @dependabot in https://github.com/apache/plc4x/pull/1372
* build(deps): bump nifi.version from 1.24.0 to 1.25.0 by @dependabot in https://github.com/apache/plc4x/pull/1371
* Updated DataIo Template (And working on the Siemens S7 L-types as well as  Temporal types) by @chrisdutz in https://github.com/apache/plc4x/pull/1376
* build(deps): bump org.assertj:assertj-core from 3.25.1 to 3.25.3 by @dependabot in https://github.com/apache/plc4x/pull/1381
* build(deps): bump org.asciidoctor:asciidoctorj-diagram from 2.2.14 to 2.2.17 by @dependabot in https://github.com/apache/plc4x/pull/1380
* build(deps): bump joda-time:joda-time from 2.12.6 to 2.12.7 by @dependabot in https://github.com/apache/plc4x/pull/1386
* build(deps): bump milo.version from 0.6.11 to 0.6.12 by @dependabot in https://github.com/apache/plc4x/pull/1385
* build(deps): bump io.jsonwebtoken:jjwt-api from 0.12.4 to 0.12.5 by @dependabot in https://github.com/apache/plc4x/pull/1384
* build(deps): bump junit.jupiter.version from 5.10.1 to 5.10.2 by @dependabot in https://github.com/apache/plc4x/pull/1383
* build(deps): bump golang.org/x/net from 0.20.0 to 0.21.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1387
* build(deps-dev): bump org.json:json from 20231013 to 20240205 by @dependabot in https://github.com/apache/plc4x/pull/1388
* build(deps): bump slf4j.version from 2.0.11 to 2.0.12 by @dependabot in https://github.com/apache/plc4x/pull/1389
* build(deps): bump org.eclipse.jetty:jetty-util from 11.0.19 to 11.0.20 by @dependabot in https://github.com/apache/plc4x/pull/1390
* build(deps): bump commons-codec:commons-codec from 1.16.0 to 1.16.1 by @dependabot in https://github.com/apache/plc4x/pull/1392
* build(deps): bump com.gradle:gradle-enterprise-maven-extension from 1.20 to 1.20.1 by @dependabot in https://github.com/apache/plc4x/pull/1393
* Fixed the last problems with writing complex types to ADS by @chrisdutz in https://github.com/apache/plc4x/pull/1394
* fix(plc4py/umas): Why not start implementing umas by @hutcheb in https://github.com/apache/plc4x/pull/1339
* Fix/enum discriminators by @chrisdutz in https://github.com/apache/plc4x/pull/1395
* build(deps): bump org.asciidoctor:asciidoctor-maven-plugin from 2.2.5 to 2.2.6 by @dependabot in https://github.com/apache/plc4x/pull/1398
* fix(plc4j): Stabilization of build after opcua security PR merges. by @splatch in https://github.com/apache/plc4x/pull/1401
* build(deps): bump org.asciidoctor:asciidoctorj-diagram from 2.2.17 to 2.3.0 by @dependabot in https://github.com/apache/plc4x/pull/1403
* build(deps): bump golang.org/x/tools from 0.17.0 to 0.18.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1402
* refactor: Refactored the project to allow inspecting the Drivers and … by @chrisdutz in https://github.com/apache/plc4x/pull/1397
* build(deps): bump org.aspectj:aspectjweaver from 1.9.21 to 1.9.21.1 by @dependabot in https://github.com/apache/plc4x/pull/1404
* Corrects the call to generate events that were commented. Tests again… by @glcj in https://github.com/apache/plc4x/pull/1400

## New Contributors
* @qtvbwfn made their first contribution in https://github.com/apache/plc4x/pull/1157
* @schaebo made their first contribution in https://github.com/apache/plc4x/pull/1374

**Full Changelog**: https://github.com/apache/plc4x/compare/rel/0.11...rel/0.12

### Feature

- general:
  - Added a pingAddress config option to the modbus driver, which allows providing an address string, that can be used for ping-operations (Defaults to reading holding-register:1) ([fa615dc](https://github.com/apache/plc4x/commit/fa615dc36617f02b9025194aeb6e44e3469d6e75))
  - Finished a first version of the release scripts. ([ec7a614](https://github.com/apache/plc4x/commit/ec7a6140490dd0632e6393fd869f32681666557b))
  - Added the code for signing all artifacts after the release-build was finished. ([2a3e5ea](https://github.com/apache/plc4x/commit/2a3e5eab9a7592eb21821bdf3faa9db2d441fb7b))
  - add support for nested configuration in the metadata reporting ([12f08be](https://github.com/apache/plc4x/commit/12f08bebbc3fb3937e036c03876cf08f07de5624))
  - add support for nested configuration in the metadata reporting ([aa8a3c9](https://github.com/apache/plc4x/commit/aa8a3c906aeb4d8c9640c4182e3bd36abda5723f))
  - add support for nested configuration in the metadata reporting ([387c44f](https://github.com/apache/plc4x/commit/387c44f01ceb94939405b381beb5c7e0c7fd9fbe))
  - OPC-UA mspec updates. ([6692f6f](https://github.com/apache/plc4x/commit/6692f6f8eaac07f896bea8fef5b4b287c116c12e))
  - Implementation of opc ua client security. ([c572e58](https://github.com/apache/plc4x/commit/c572e58db5df605ad7b1564f528f88f85c0fa43a))
  - Added support for all missing S7 64bit types (L-Types) as well as Duration/Time/Date types  ([a7310e3](https://github.com/apache/plc4x/commit/a7310e34ded164079d954753af18490677aaa975)) ([#1376](https://github.com/apache/plc4x/pull/1376))
  - Implemented the functionality to edit, add, delete devices. ([dc6dd21](https://github.com/apache/plc4x/commit/dc6dd21db3c0fa5aa977cb8cdcc91f2a6a3a42f8))
  - Continued working on getting the PN stuff working. ([43b5e7d](https://github.com/apache/plc4x/commit/43b5e7d2f490ec7aa7b3aa924e017d499f6e6d8c))
  - Added the ability to name expectations (helping debug missing requests or replies). Continued working on the PN stuff. ([6bfdda4](https://github.com/apache/plc4x/commit/6bfdda4a9bfbb8ad28ed4b23256c7d241a80bf53))
  - Added the ability to name expectations (helping debug missing requests or replies). Continued working on the PN stuff. ([c31148d](https://github.com/apache/plc4x/commit/c31148d1fba27d538db13c6f32a2a085522bea81))
  - Fixed some dependency problems. ([81d514e](https://github.com/apache/plc4x/commit/81d514e1c2bf96967edadf694b9610753c9d004e))
  - Updated the ProfinetDriver, to intercept mac-address connection strings and in this case to initially update the remote devices IP address using PN-DCP before actually initializing the PN connection. ([54fcaf6](https://github.com/apache/plc4x/commit/54fcaf63b64c1246a5b06cca596aea215cc9552e))
  - Updated the ProfinetDiscoverer to be able to handle PN devices without assigned IP addresses. ([9faa864](https://github.com/apache/plc4x/commit/9faa864c1e79d073ccbfd5208c499c04b3e534ad))
  - Continued working on implementing the PROFINET driver. ([67fb856](https://github.com/apache/plc4x/commit/67fb856ea09e211802aaa4cdd3cb2a9f121702f1))
  - Cleanup of S7 STRING/WSTRING/(various temporal types) ([d90c3dd](https://github.com/apache/plc4x/commit/d90c3dde1f394657299a3ffe67fe375661f80f73)) ([#1259](https://github.com/apache/plc4x/pull/1259))
  - Disabled a test that was made to run manually. ([0a3cec6](https://github.com/apache/plc4x/commit/0a3cec628dbebe3703c6bb5ea2c5780ebdd13f11))
  - Added a first working draft of a Boschrexroth CtlX driver. ([76cc647](https://github.com/apache/plc4x/commit/76cc647ebb56ac27fd66d2a8f24398e3dfcf6202))

- plc4j/opcua:
  - Chunking and encryption of request/response calls. ([d0893d7](https://github.com/apache/plc4x/commit/d0893d7ac786405783e0190676ba11e082511617))

- plc4py/umas:
  - UMAS implementation (#1339) ([780a104](https://github.com/apache/plc4x/commit/780a104cbafd5b0beb8b8649225c05b69de51205)) ([#1339](https://github.com/apache/plc4x/pull/1339))

- plc4go:
  - support for DATE_AND_LTIME ([51434bd](https://github.com/apache/plc4x/commit/51434bd8bde1b43785d83d5665566a4bdcb0561b))

- plc4go/connectioncache:
  - add GetConnectionWithContext to connection cache ([b6fbfbc](https://github.com/apache/plc4x/commit/b6fbfbcc1d2675b4db924f5cca6ad0932a7d35aa))

- plc4py:
  - Code Gen Update (#1199) - Incremental update to python ([9b7bf47](https://github.com/apache/plc4x/commit/9b7bf47b5b9ec53872c44db21c77be7494f9df0e)) ([#1199](https://github.com/apache/plc4x/pull/1199))

- plc4j:
  - Refactored the way transports are configured via the additional parameters of the connection string. ([554c083](https://github.com/apache/plc4x/commit/554c083b160b78abed246a1bf9c679c6ecd47fbe)) ([#1179](https://github.com/apache/plc4x/pull/1179))

- codegen:
  - Moved the code-generation into a separate profile (#1172) ([063524b](https://github.com/apache/plc4x/commit/063524b6e5931b8f3cfffbf7d087687db4f8fa81)) ([#1172](https://github.com/apache/plc4x/pull/1172))

- plc4j/spi:
  - Add option to synchronously await response from PLC (#1163) ([6073d4b](https://github.com/apache/plc4x/commit/6073d4bf4755c31600e858745bac45a12e937d6f)) ([#1163](https://github.com/apache/plc4x/pull/1163))

- plc4y:
  - Started to implement reading arrays ([e1ae587](https://github.com/apache/plc4x/commit/e1ae587ab3ffeec5a93c2c1e23cdd68b6b26b717))

- plc4j/drivers:
  - Create maven meta package to include all drivers (#1166) ([2d428ca](https://github.com/apache/plc4x/commit/2d428ca5acc1cfd64967f758c433cf02f96029e2)) ([#1166](https://github.com/apache/plc4x/pull/1166))

- plc4x-server:
  - Build a standalone jar too (#1167) ([18883b1](https://github.com/apache/plc4x/commit/18883b1e3eb5b31ad707a21a2ce9d9c3817393f2)) ([#1167](https://github.com/apache/plc4x/pull/1167))

### Bug Fixes

- general:
  - nested defaults should now be unwrapped ([ccd8bd2](https://github.com/apache/plc4x/commit/ccd8bd22662d5a39ab50c5e47f681895735bd7c2))
  - Fixed several issues with the new generated-documentation. ([2ab1795](https://github.com/apache/plc4x/commit/2ab1795b35597c10b7713d67ec2bb1cbacc3e4e4))
  - Fixed several issues with the new generated-documentation. ([f7a6a34](https://github.com/apache/plc4x/commit/f7a6a3406e62a01377f67e407a978132ea0ed421))
  - Enabled the "with-java" profile when building the plc4j-driver-all ([11de90a](https://github.com/apache/plc4x/commit/11de90a46ad9cb6d9aea92ae815d54bf881a600c))
  - Updated the connection string to the updated attribute names. ([dce42c4](https://github.com/apache/plc4x/commit/dce42c45a3d87c7a50db97ea5c94430b0751077c))
  - Fix remaining unit tests. ([7202ec6](https://github.com/apache/plc4x/commit/7202ec6530e09aa68f3106bbe26d45799d0a9c9a))
  - Made sure plc-values are serialized in LittleEndian format ([f132b0b](https://github.com/apache/plc4x/commit/f132b0b9439af8f29da70e97dade79bcca9bf6d7))
  - Fixed the configuration options for serial transports (the parity, was actually not the number of parity bits, but the type of parity being used)
feat: Now it's possible to use enums in configuration. ([77a9b12](https://github.com/apache/plc4x/commit/77a9b12f69fcbac488c7a21fc15fe2bc4630adf1))
  - fixed the build prior to java 19 ([464f510](https://github.com/apache/plc4x/commit/464f5106d98d0f4bd4ad3dedd203ccb3f9c09556))
  - Ensured everywhere a transaction manager is created, that it's also shutdown correctly. ([add6906](https://github.com/apache/plc4x/commit/add690627334dd3cc9616ce9f13e26d55710e059))
  - Fixed a build error. ([be54e22](https://github.com/apache/plc4x/commit/be54e22225f1b8514d3f946b3fbecafebc1c9d23))
  - Reduced the version of the rat plugin till the version 0.16.1 is released, which will fix the problem causing the build to pause for several minutes. ([1f9a949](https://github.com/apache/plc4x/commit/1f9a949a781955c6e38d355ba5f188b767062765))
  - Gave the initial connection request a bit more time ... also tried to find the reason for the reconnects. ([d54b88d](https://github.com/apache/plc4x/commit/d54b88de97eb76dc6e20aa37830dc0b6c8e03e69))
  - Gave the initial connection request a bit more time ... also tried to find the reason for the reconnects. ([54888ea](https://github.com/apache/plc4x/commit/54888eaa7efb87f6990669cb8de3374bf231fad2))
  - Found and fixed the reason for the second reconnect to fail. ([af6bdd5](https://github.com/apache/plc4x/commit/af6bdd53da5400f59a4769cf63da8eafbba0e279))
  - Replaced the while-loop in the ProfinetDevice with a timer triggered executor and this made the connection stable. ([e9ffa6b](https://github.com/apache/plc4x/commit/e9ffa6bb84b74a2b60006342c00db319ca7314a8))
  - Managed to acknowledge the ApplicationReady request ... seems that there's problems with the "expectRequest" functionality. ([c5fa7c3](https://github.com/apache/plc4x/commit/c5fa7c3e23b1e2007446064b8916ef2c4ebb3014))
  - Got the data flowing in the NG-Profinet driver. ([6df0fa5](https://github.com/apache/plc4x/commit/6df0fa59d20bd4e7994a65e6f590126106b9188d))
  - Fixed some issues in the old PN driver (Adjusted the Transport config and made sure it doesn't consume own outgoing messages) ([48f8e65](https://github.com/apache/plc4x/commit/48f8e65437d79ba04dc034db0e51f0164b7a9c27))
  - Implemented closing of the Profinet (and RawSocketChannel) ([850b9a1](https://github.com/apache/plc4x/commit/850b9a1297b2fe52048350d1659eb4947edd2899))
  - Updated the build to output failed go tests when being run in the maven build. ([8fd198c](https://github.com/apache/plc4x/commit/8fd198c25a0940ab36b8a961dbcf6ebc01739165))
  - Reduced the wait between the two executions as some times the first operation already is finished. ([edbce97](https://github.com/apache/plc4x/commit/edbce97442fb796a18bd672bd9b38fbdaaf8e605))
  - Updated the settings for nexus deployment as it seems the labels were changed. ([ed28a52](https://github.com/apache/plc4x/commit/ed28a52695cf0268f7713eefd8d24f51c6c413e6))
  - Ensured the XML transformation is executed before the execution of the resources-plugin ([8bf0d93](https://github.com/apache/plc4x/commit/8bf0d93c695a9aa5a287ee7c0b30c09e360091b4))
  - Resolved a compilation error with the examples after switching wot the "plc4j-drivers-all" module. ([400cbf5](https://github.com/apache/plc4x/commit/400cbf58a4559f70e5969ec2e8ed46ae0817f16a))
  - Make sure a leased-connection isn't double-closed ([f70f3e1](https://github.com/apache/plc4x/commit/f70f3e1bca07ec28efde0c3461da523171f02191))
  - Make sure a leased-connection isn't double-closed ([87a163c](https://github.com/apache/plc4x/commit/87a163cbba74129d6b28960ee1d6a5ef80edba09))
  - Made the project buildable on Java 21 (by excluding the Kafka-Connect integration) ([ce7a205](https://github.com/apache/plc4x/commit/ce7a205dfcb757300359821e1682eeebe37f2ad4))
  - Increased the memory for running unit tests, as the bacnet module was causing OOM errors. ([8f49123](https://github.com/apache/plc4x/commit/8f49123aaacdd6c4b2ca90c8f5e748f9ab3012cf))
  - Fixed the problems left over from the release ([4afa650](https://github.com/apache/plc4x/commit/4afa65037b47454d08bd56cc0ac6fed93f6b5125))
  - Fixed the problems left over from the release ([c48792e](https://github.com/apache/plc4x/commit/c48792ef7ad7e4edc033aa61ff6283e4fc9ec859))
  - Made sure the OpcuaSubscriptionHandleTest doesn't run on Docker. ([a5e78c8](https://github.com/apache/plc4x/commit/a5e78c8489b7b6811737e1767e6036f61b8fb0da))
  - Update the release-check-tools to the latest changes. ([042a4a1](https://github.com/apache/plc4x/commit/042a4a17fde3300a84db755821e88e1a7fb4810a))

- plc4j:
  - Stabilization of build after opcua security PR merges. ([271e1cb](https://github.com/apache/plc4x/commit/271e1cbfb0b65e02cd07735e22ce58764502eb17)) ([#1401](https://github.com/apache/plc4x/pull/1401))

- plc4go:
  - Fix golang compile errors. ([6b033d3](https://github.com/apache/plc4x/commit/6b033d3f7d29b8a20c716eea00d11acc2bbe70c1))
  - Update golang opcua structs. ([f8b0b3b](https://github.com/apache/plc4x/commit/f8b0b3bf7782f14915d05145cbf2a621185fcd76))
  - options should now correctly be applied ([9f658e6](https://github.com/apache/plc4x/commit/9f658e6ded38ab2fa4d621dbca56297e0f86b2d4))
  - wrong logger usage ([fb33c7b](https://github.com/apache/plc4x/commit/fb33c7b01189b0aa7d302717c9f383870095f96c))
  - port over s7 changes to golang ([faba6fe](https://github.com/apache/plc4x/commit/faba6fec57a1e4e4207e4e4450891432adb11691))
  - ensure discoverer respect context cancel ([fae748c](https://github.com/apache/plc4x/commit/fae748c3663352a82490b290993e1a9008e3fd47))

- knx:
  - Fixed how 16 bit floating-point numbers are parsed. ([1259eb1](https://github.com/apache/plc4x/commit/1259eb15d4607bd7812d7150bc7451b9411f10a9))
  - Fixed the typed for knx datatypes PDT_BITSET8 and PDT_BITSET16 ([fc89f78](https://github.com/apache/plc4x/commit/fc89f789a840c5b23b4b39b135bd31e15dcb9338))

- plc4go/s7:
  - port over some changes from plc4j ([fdda471](https://github.com/apache/plc4x/commit/fdda471da2e241729d8524846bfc2f00f3b84d5a))

- plc4j/eip:
  - Remove unnecessary Thread.sleep (#1374) ([7ca7602](https://github.com/apache/plc4x/commit/7ca76028c477f2cdd49a72876336ddde96da2555)) ([#1374](https://github.com/apache/plc4x/pull/1374))

- website:
  - Update hutcheb details ([7aa46a2](https://github.com/apache/plc4x/commit/7aa46a269603632baddebf1f7042b084b3c9b3d7))

- plc4py/modbus:
  - Tidy up the Modbus implementation ([c27ccbb](https://github.com/apache/plc4x/commit/c27ccbbcddda936db2f2e4790b6d3ac023473076))

- plc4py:
  - Add abstract decorators to read buffer ([f37f277](https://github.com/apache/plc4x/commit/f37f2775b8456a8c3b70597d0aacf34a6860f3f9))

- plc4j/profinet:
  - Fix for Array Count and Application Ready Ports ([c727b60](https://github.com/apache/plc4x/commit/c727b60f9c1e0a1f97a99b0c5ba0fed80fc3ff93))

- plc4j/spi:
  - fix(plc4j/spi) Make sure OPC UA discover event is fired prior connected event. ([d023f06](https://github.com/apache/plc4x/commit/d023f06ce30c2c7ca3c15620e4beb7cd59e6bab5)) ([#1217](https://github.com/apache/plc4x/pull/1217))

- plc4x-server:
  - Fix logs not being written to screen (#1161) ([0a62cb2](https://github.com/apache/plc4x/commit/0a62cb228d2548f7fd45905120f4fbd8854de944)) ([#1161](https://github.com/apache/plc4x/pull/1161))

- plc4j/opcua:
  - OPC UA priority judgment using discovery parameter (#1157) ([a17bcc8](https://github.com/apache/plc4x/commit/a17bcc865208a68f521d7688b26422e5c07bb7f2)) ([#1157](https://github.com/apache/plc4x/pull/1157))
  - Fix keepalive threads are never shut down (#1139) ([654929d](https://github.com/apache/plc4x/commit/654929dda0047b65c984379eca2adb43fb77252c)) ([#1139](https://github.com/apache/plc4x/pull/1139))

- opcua:
  - Await `writeAndFlush(msg)` & send next msg async (#1147) ([b0bc847](https://github.com/apache/plc4x/commit/b0bc847e923a2a59c15868b642d377738a082402)) ([#1147](https://github.com/apache/plc4x/pull/1147))
  - Add Null Variant (#1124) ([847fead](https://github.com/apache/plc4x/commit/847feadbfd3acc40c127e9d6b6802328ce42f56c)) ([#1124](https://github.com/apache/plc4x/pull/1124))

### Documentation

- general:
  - Updated the most recent changes. ([a1da03a](https://github.com/apache/plc4x/commit/a1da03aa013b00eb8802e69e1233016bfcef3098))
  - Fixed some invalid links on the webpage. ([526ca16](https://github.com/apache/plc4x/commit/526ca16b827ba652e4341351cedb82538851a7fe))
  - Updated the way the headlines of the transport options are highlighted to using h4 blocks. ([b5912bd](https://github.com/apache/plc4x/commit/b5912bd50c298228a962dffad30525492cecf24c))
  - Updated two more places, where we didn't output a concrete version in the docs. ([fd53a59](https://github.com/apache/plc4x/commit/fd53a59a68db6b94b5d24d0fb42f7c3678c404c9))
  - Made the current-last-released-version be output instead of the place-holder. ([e156be1](https://github.com/apache/plc4x/commit/e156be1c080e364e8a08a41ae05788db3055f423))
  - Updated the release-image-template.svg for the release of 0.11.0 ([db237f0](https://github.com/apache/plc4x/commit/db237f0aefb145a93311a0617efef09a0747d717))
  - Updated the getting started documentation to the new version. ([962374f](https://github.com/apache/plc4x/commit/962374fd2a37efdae9c2338056166d7f4d20e77c))
  - Added some more profiles to the documentation to really build everything. ([e4dec50](https://github.com/apache/plc4x/commit/e4dec500c1ca4d46b95b1f9294f0836142e37729))
  - Added the "enable-all-checks" to the release documentation. ([b869ce4](https://github.com/apache/plc4x/commit/b869ce4d63fb17e9e25d5089367b18c7bf10fd65))

### Refactor

- general:
  - Cleaned up in the API (moved the PlcConfiguration-related interfaces back to SPI), renamed the "canX()" methods to match the bean naming convention of "isXSupported()". Also removed one or two obsolete interfaces) ([0329142](https://github.com/apache/plc4x/commit/03291421260409bb2bfecdfa4a4acb92adedfa7a))
  - Refactored the API to provide metadata information directly. ([8c4bd28](https://github.com/apache/plc4x/commit/8c4bd28a43120dd651e9a1f475734cc04197216a))
  - Renamed the frontend directory inside the frontend module to "project" ([663bf95](https://github.com/apache/plc4x/commit/663bf957bc7af6a507ef8093e9e22a3bb6fc2611))
  - Updated the discovery example to not list every device for every local device that could reach it. ([08dcb72](https://github.com/apache/plc4x/commit/08dcb7221fa35fe7c48c542e4bcb2c2238b11eb6))

- plc4x-server:
  - Cleanup, add tests, and add option to specify port number (#1162) ([8091bfd](https://github.com/apache/plc4x/commit/8091bfd69218eb87665b3d291b12d46c0e5d0782)) ([#1162](https://github.com/apache/plc4x/pull/1162))

## [rel/0.11](https://github.com/apache/plc4x/releases/tag/rel/0.11) - 2023-10-02 07:23:01

## [rel/0.10](https://github.com/apache/plc4x/releases/tag/rel/0.10) - 2022-09-29 15:56:37

## [rel/0.9](https://github.com/apache/plc4x/releases/tag/rel/0.9) - 2021-09-17 09:25:12

## [rel/0.8](https://github.com/apache/plc4x/releases/tag/rel/0.8) - 2021-01-25 11:31:05

## [rel/0.7](https://github.com/apache/plc4x/releases/tag/rel/0.7) - 2020-05-15 10:51:06

## [rel/0.6](https://github.com/apache/plc4x/releases/tag/rel/0.6) - 2020-01-13 14:49:38

## [rel/0.5](https://github.com/apache/plc4x/releases/tag/rel/0.5) - 2019-10-21 15:00:13

## [rel/0.4](https://github.com/apache/plc4x/releases/tag/rel/0.4) - 2019-05-20 21:09:32

## [rel/0.3](https://github.com/apache/plc4x/releases/tag/rel/0.3) - 2019-01-16 15:32:56

## [rel/0.2](https://github.com/apache/plc4x/releases/tag/rel/0.2) - 2018-11-14 15:20:37

## [rel/0.1](https://github.com/apache/plc4x/releases/tag/rel/0.1) - 2018-09-17 16:31:32

\* *This CHANGELOG was automatically generated by [auto-generate-changelog](https://github.com/BobAnkh/auto-generate-changelog)*
