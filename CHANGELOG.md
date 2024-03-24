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

## What's Changed
* build(deps): bump actions/github-script from 6.3.0 to 6.3.1 by @dependabot in https://github.com/apache/plc4x/pull/530
* build(deps): bump logback.version from 1.4.1 to 1.4.3 by @dependabot in https://github.com/apache/plc4x/pull/529
* build(deps): bump kotlin.version from 1.7.10 to 1.7.20 by @dependabot in https://github.com/apache/plc4x/pull/528
* build(deps): bump protobuf-java from 3.21.6 to 3.21.7 by @dependabot in https://github.com/apache/plc4x/pull/527
* build(deps): bump influxdb-client-java from 6.5.0 to 6.6.0 by @dependabot in https://github.com/apache/plc4x/pull/526
* build(deps): bump spock-bom from 2.2-groovy-4.0 to 2.3-groovy-4.0 by @dependabot in https://github.com/apache/plc4x/pull/525
* build(deps): bump nifi.version from 1.17.0 to 1.18.0 by @dependabot in https://github.com/apache/plc4x/pull/535
* build(deps): bump checker-qual from 3.25.0 to 3.26.0 by @dependabot in https://github.com/apache/plc4x/pull/532
* build(deps): bump camel.version from 3.18.2 to 3.19.0 by @dependabot in https://github.com/apache/plc4x/pull/531
* fix(plc4go): pass in the correct context by @hongjinlin in https://github.com/apache/plc4x/pull/537
* build(deps): bump logback.version from 1.4.3 to 1.4.4 by @dependabot in https://github.com/apache/plc4x/pull/539
* build(deps): bump gmavenplus-plugin from 1.13.1 to 2.0.0 by @dependabot in https://github.com/apache/plc4x/pull/536
* build(deps): bump jackson.version from 2.14.0-rc1 to 2.14.0-rc2 by @dependabot in https://github.com/apache/plc4x/pull/538
* build(deps): bump actions/github-script from 6.3.1 to 6.3.2 by @dependabot in https://github.com/apache/plc4x/pull/543
* build(deps): bump netty-bom from 4.1.82.Final to 4.1.84.Final by @dependabot in https://github.com/apache/plc4x/pull/541
* build(deps): bump byte-buddy from 1.12.17 to 1.12.18 by @dependabot in https://github.com/apache/plc4x/pull/540
* build(deps): bump error_prone_annotations from 2.15.0 to 2.16 by @dependabot in https://github.com/apache/plc4x/pull/542
* build(deps): bump github.com/spf13/cobra from 1.5.0 to 1.6.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/544
* build(deps): bump joda-time from 2.11.2 to 2.12.0 by @dependabot in https://github.com/apache/plc4x/pull/547
* build(deps): bump iot-device-client from 2.1.1 to 2.1.2 by @dependabot in https://github.com/apache/plc4x/pull/546
* build(deps): bump actions/github-script from 6.3.2 to 6.3.3 by @dependabot in https://github.com/apache/plc4x/pull/548
* Setup the Channel Pipeline by @hutcheb in https://github.com/apache/plc4x/pull/362
* build(deps): bump swagger-annotations from 1.6.7 to 1.6.8 by @dependabot in https://github.com/apache/plc4x/pull/549
* build(deps): bump groovy.version from 4.0.5 to 4.0.6 by @dependabot in https://github.com/apache/plc4x/pull/550
* build(deps): bump protobuf-java from 3.21.7 to 3.21.8 by @dependabot in https://github.com/apache/plc4x/pull/551
* build(deps): bump golang.org/x/tools from 0.1.12 to 0.2.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/553
* build(deps): bump gmavenplus-plugin from 2.0.0 to 2.1.0 by @dependabot in https://github.com/apache/plc4x/pull/552
* build(deps): bump github.com/stretchr/testify from 1.8.0 to 1.8.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/557
* build(deps): bump karaf-maven-plugin from 4.4.1 to 4.4.2 by @dependabot in https://github.com/apache/plc4x/pull/554
* build(deps): bump mockito.version from 4.8.0 to 4.8.1 by @dependabot in https://github.com/apache/plc4x/pull/555
* build(deps): bump asciidoctorj from 2.5.6 to 2.5.7 by @dependabot in https://github.com/apache/plc4x/pull/556
* build(deps): bump woodstox-core from 6.3.1 to 6.4.0 by @dependabot in https://github.com/apache/plc4x/pull/559
* build(deps): bump gson from 2.9.1 to 2.10 by @dependabot in https://github.com/apache/plc4x/pull/560
* build(deps): bump github.com/spf13/cobra from 1.6.0 to 1.6.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/561
* build(deps): bump protobuf-java from 3.21.8 to 3.21.9 by @dependabot in https://github.com/apache/plc4x/pull/562
* build(deps): bump joda-time from 2.12.0 to 2.12.1 by @dependabot in https://github.com/apache/plc4x/pull/565
* build(deps): bump jackson.version from 2.14.0-rc2 to 2.14.0-rc3 by @dependabot in https://github.com/apache/plc4x/pull/564
* build(deps-dev): bump commons-compress from 1.21 to 1.22 by @dependabot in https://github.com/apache/plc4x/pull/566
* build(deps): bump influxdb-client-java from 6.6.0 to 6.7.0 by @dependabot in https://github.com/apache/plc4x/pull/572
* build(deps): bump github.com/schollz/progressbar/v3 from 3.11.0 to 3.12.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/571
* build(deps): bump equalsverifier from 3.10.1 to 3.11 by @dependabot in https://github.com/apache/plc4x/pull/570
* build(deps): bump maven-release-plugin from 3.0.0-M6 to 3.0.0-M7 by @dependabot in https://github.com/apache/plc4x/pull/569
* build(deps): bump checker-qual from 3.26.0 to 3.27.0 by @dependabot in https://github.com/apache/plc4x/pull/567
* build(deps): bump crc from 1.0.3 to 1.1.0 by @dependabot in https://github.com/apache/plc4x/pull/575
* build(deps): bump jackson.version from 2.14.0-rc3 to 2.14.0 by @dependabot in https://github.com/apache/plc4x/pull/573
* build(deps): bump github.com/spf13/viper from 1.13.0 to 1.14.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/578
* build(deps): bump github.com/schollz/progressbar/v3 from 3.12.0 to 3.12.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/577
* Feature/cdutz/go ads ng (Streamlining of PLC4X API in PLC4Go and PLC4J) by @chrisdutz in https://github.com/apache/plc4x/pull/576
* build(deps): bump golang.org/x/tools from 0.2.0 to 0.3.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/651
* build(deps): bump netty-bom from 4.1.84.Final to 4.1.85.Final by @dependabot in https://github.com/apache/plc4x/pull/650
* build(deps): bump kotlin.version from 1.7.20 to 1.7.21 by @dependabot in https://github.com/apache/plc4x/pull/579
* build(deps): bump actions/dependency-review-action from 2 to 3 by @dependabot in https://github.com/apache/plc4x/pull/652
* build(deps): bump mockito.version from 4.8.1 to 4.9.0 by @dependabot in https://github.com/apache/plc4x/pull/654
* build(deps): bump slf4j.version from 2.0.3 to 2.0.4 by @dependabot in https://github.com/apache/plc4x/pull/657
* build(deps): bump byte-buddy from 1.12.18 to 1.12.19 by @dependabot in https://github.com/apache/plc4x/pull/658
* build(deps): bump maven-dependency-tree from 3.2.0 to 3.2.1 by @dependabot in https://github.com/apache/plc4x/pull/661
* build(deps): bump logback.version from 1.4.4 to 1.4.5 by @dependabot in https://github.com/apache/plc4x/pull/662
* build(deps): bump equalsverifier from 3.11 to 3.11.1 by @dependabot in https://github.com/apache/plc4x/pull/660
* build(deps): bump jackson.version from 2.14.0 to 2.14.1 by @dependabot in https://github.com/apache/plc4x/pull/665
* build(deps): bump apache from 27 to 28 by @dependabot in https://github.com/apache/plc4x/pull/666
* fix(plc4j/opcua): not to get String.length but to calculate length in… by @modraedlau in https://github.com/apache/plc4x/pull/668
* build(deps): bump BobAnkh/auto-generate-changelog from 1.2.2 to 1.2.3 by @dependabot in https://github.com/apache/plc4x/pull/676
* build(deps): bump kotlin.version from 1.7.21 to 1.7.22 by @dependabot in https://github.com/apache/plc4x/pull/674
* build(deps): bump nifi.version from 1.18.0 to 1.19.0 by @dependabot in https://github.com/apache/plc4x/pull/673
* build(deps): bump slf4j.version from 2.0.4 to 2.0.5 by @dependabot in https://github.com/apache/plc4x/pull/672
* build(deps): bump httpcore from 4.4.15 to 4.4.16 by @dependabot in https://github.com/apache/plc4x/pull/675
* build(deps): bump equalsverifier from 3.11.1 to 3.12.1 by @dependabot in https://github.com/apache/plc4x/pull/678
* build(deps): bump jSerialComm from 2.9.2 to 2.9.3 by @dependabot in https://github.com/apache/plc4x/pull/679
* build(deps): bump checker-qual from 3.27.0 to 3.28.0 by @dependabot in https://github.com/apache/plc4x/pull/677
* build(deps): bump commons-net from 3.8.0 to 3.9.0 by @dependabot in https://github.com/apache/plc4x/pull/680
* build(deps-dev): bump httpclient from 4.5.13 to 4.5.14 by @dependabot in https://github.com/apache/plc4x/pull/684
* build(deps): bump joda-time from 2.12.1 to 2.12.2 by @dependabot in https://github.com/apache/plc4x/pull/682
* build(deps): bump protobuf-java from 3.21.9 to 3.21.10 by @dependabot in https://github.com/apache/plc4x/pull/681
* build(deps): bump github.com/schollz/progressbar/v3 from 3.12.1 to 3.12.2 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/687
* build(deps): bump golang.org/x/tools from 0.3.0 to 0.4.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/686
* build(deps): bump protobuf-java from 3.21.10 to 3.21.11 by @dependabot in https://github.com/apache/plc4x/pull/690
* build(deps): bump equalsverifier from 3.12.1 to 3.12.2 by @dependabot in https://github.com/apache/plc4x/pull/689
* build(deps): bump slf4j.version from 2.0.5 to 2.0.6 by @dependabot in https://github.com/apache/plc4x/pull/693
* build(deps): bump jetty-util from 11.0.12 to 11.0.13 by @dependabot in https://github.com/apache/plc4x/pull/692
* build(deps): bump nifi.version from 1.19.0 to 1.19.1 by @dependabot in https://github.com/apache/plc4x/pull/688
* build(deps): bump ecj from 3.31.0 to 3.32.0 by @dependabot in https://github.com/apache/plc4x/pull/685
* build(deps): bump actions/setup-python from 3 to 4 by @dependabot in https://github.com/apache/plc4x/pull/655
* build(deps): bump swagger-annotations from 1.6.8 to 1.6.9 by @dependabot in https://github.com/apache/plc4x/pull/656
* build(deps): bump iot-device-client from 2.1.2 to 2.1.3 by @dependabot in https://github.com/apache/plc4x/pull/700
* build(deps): bump netty-bom from 4.1.85.Final to 4.1.86.Final by @dependabot in https://github.com/apache/plc4x/pull/699
* build(deps): bump protobuf-java from 3.21.11 to 3.21.12 by @dependabot in https://github.com/apache/plc4x/pull/696
* build(deps): bump apache from 28 to 29 by @dependabot in https://github.com/apache/plc4x/pull/697
* build(deps): bump mockito.version from 4.9.0 to 4.10.0 by @dependabot in https://github.com/apache/plc4x/pull/698
* build(deps): bump equalsverifier from 3.12.2 to 3.12.3 by @dependabot in https://github.com/apache/plc4x/pull/707
* build(deps): bump byte-buddy from 1.12.19 to 1.12.20 by @dependabot in https://github.com/apache/plc4x/pull/706
* build(deps): bump maven-invoker-plugin from 3.3.0 to 3.4.0 by @dependabot in https://github.com/apache/plc4x/pull/705
* fix: Always install sources by @nielsbasjes in https://github.com/apache/plc4x/pull/717
* fix(protocols/modbus): fix write requests for coils always set to false (#710) by @SteinOv in https://github.com/apache/plc4x/pull/711
* Plc4py codegen by @hutcheb in https://github.com/apache/plc4x/pull/720
* Limit modbus quantity by @nielsbasjes in https://github.com/apache/plc4x/pull/721
* build(deps): bump mockito.version from 4.10.0 to 4.11.0 by @dependabot in https://github.com/apache/plc4x/pull/719
* build(deps): bump github.com/gdamore/tcell/v2 from 2.5.3 to 2.5.4 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/723
* build(deps): bump error_prone_annotations from 2.16 to 2.17.0 by @dependabot in https://github.com/apache/plc4x/pull/725
* build(deps): bump groovy.version from 4.0.6 to 4.0.7 by @dependabot in https://github.com/apache/plc4x/pull/712
* build(deps): bump golang.org/x/tools from 0.4.0 to 0.5.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/727
* build(deps): bump byte-buddy from 1.12.20 to 1.12.21 by @dependabot in https://github.com/apache/plc4x/pull/726
* build(deps): bump github.com/schollz/progressbar/v3 from 3.12.2 to 3.13.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/728
* build(deps): bump assertj-core from 3.23.1 to 3.24.0 by @dependabot in https://github.com/apache/plc4x/pull/730
* fix(plc4j/modbus): Cleanup of ModbusTag by @nielsbasjes in https://github.com/apache/plc4x/pull/732
* build(deps): bump assertj-core from 3.24.0 to 3.24.1 by @dependabot in https://github.com/apache/plc4x/pull/733
* build(deps): bump gson from 2.10 to 2.10.1 by @dependabot in https://github.com/apache/plc4x/pull/734
* build(deps): bump junit.jupiter.version from 5.9.1 to 5.9.2 by @dependabot in https://github.com/apache/plc4x/pull/735
* build(deps): bump error_prone_annotations from 2.17.0 to 2.18.0 by @dependabot in https://github.com/apache/plc4x/pull/736
* build(deps): bump checker-qual from 3.28.0 to 3.29.0 by @dependabot in https://github.com/apache/plc4x/pull/729
* Make sure s7 packet len is specified to avoid serialization errors by @splatch in https://github.com/apache/plc4x/pull/691
* build(deps): bump karaf-maven-plugin from 4.4.2 to 4.4.3 by @dependabot in https://github.com/apache/plc4x/pull/737
* build(deps): bump jakarta.activation-api from 2.1.0 to 2.1.1 by @dependabot in https://github.com/apache/plc4x/pull/738
* build(deps): bump maven-surefire-plugin from 3.0.0-M7 to 3.0.0-M8 by @dependabot in https://github.com/apache/plc4x/pull/739
* build(deps): bump maven-dependency-plugin from 3.1.2 to 3.5.0 by @dependabot in https://github.com/apache/plc4x/pull/742
* build(deps): bump byte-buddy from 1.12.21 to 1.12.22 by @dependabot in https://github.com/apache/plc4x/pull/741
* build(deps): bump xmlunit.version from 2.9.0 to 2.9.1 by @dependabot in https://github.com/apache/plc4x/pull/740
* build(deps): bump mockito.version from 4.11.0 to 5.0.0 by @dependabot in https://github.com/apache/plc4x/pull/746
* build(deps): bump freemarker from 2.3.31 to 2.3.32 by @dependabot in https://github.com/apache/plc4x/pull/744
* build(deps): bump netty-bom from 4.1.86.Final to 4.1.87.Final by @dependabot in https://github.com/apache/plc4x/pull/745
* build(deps): bump kotlin.version from 1.7.22 to 1.8.0 by @dependabot in https://github.com/apache/plc4x/pull/718
* build(deps): bump assertj-core from 3.24.1 to 3.24.2 by @dependabot in https://github.com/apache/plc4x/pull/751
* Feature/nifi integration type mapping by @QuanticPony in https://github.com/apache/plc4x/pull/752
* build(deps): bump github.com/spf13/viper from 1.14.0 to 1.15.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/756
* build(deps): bump woodstox-core from 6.4.0 to 6.5.0 by @dependabot in https://github.com/apache/plc4x/pull/754
* build(deps): bump javafx.version from 19 to 19.0.2 by @dependabot in https://github.com/apache/plc4x/pull/753
* build(deps): bump Saxon-HE from 11.4 to 12.0 by @dependabot in https://github.com/apache/plc4x/pull/750
* New Implementation of the Connection-Cache by @chrisdutz in https://github.com/apache/plc4x/pull/747
* docs(user): Marked the connection-pool as removed in the website by @chrisdutz in https://github.com/apache/plc4x/pull/758
* build(deps): bump dom4j from 2.1.3 to 2.1.4 by @dependabot in https://github.com/apache/plc4x/pull/760
* build(deps): bump jna from 5.12.1 to 5.13.0 by @dependabot in https://github.com/apache/plc4x/pull/759
* build(deps): bump maven-failsafe-plugin from 3.0.0-M7 to 3.0.0-M8 by @dependabot in https://github.com/apache/plc4x/pull/762
* build(deps): bump github.com/fatih/color from 1.13.0 to 1.14.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/764
* build(deps): bump javafx.version from 19.0.2 to 19.0.2.1 by @dependabot in https://github.com/apache/plc4x/pull/761
* build(deps): bump groovy.version from 4.0.7 to 4.0.8 by @dependabot in https://github.com/apache/plc4x/pull/765
* build(deps): bump github.com/rs/zerolog from 1.28.0 to 1.29.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/769
* build(deps): bump actions/github-script from 6.3.3 to 6.4.0 by @dependabot in https://github.com/apache/plc4x/pull/768
* build(deps): bump zip4j from 2.11.2 to 2.11.3 by @dependabot in https://github.com/apache/plc4x/pull/767
* build(deps): bump iot-device-client from 2.1.3 to 2.1.4 by @dependabot in https://github.com/apache/plc4x/pull/766
* Fix issue-602, completely kills the tasks associated with the S7 driver. by @glcj in https://github.com/apache/plc4x/pull/771
* Fix issue-701 for S7 driver. by @glcj in https://github.com/apache/plc4x/pull/770
* build(deps): bump equalsverifier from 3.12.3 to 3.12.4 by @dependabot in https://github.com/apache/plc4x/pull/774
* build(deps): bump jackson.version from 2.14.1 to 2.14.2 by @dependabot in https://github.com/apache/plc4x/pull/773
* build(deps): bump mockito.version from 5.0.0 to 5.1.0 by @dependabot in https://github.com/apache/plc4x/pull/772
* Plc4j/Profinet by @hutcheb in https://github.com/apache/plc4x/pull/534
* build(deps): bump mockito.version from 5.1.0 to 5.1.1 by @dependabot in https://github.com/apache/plc4x/pull/775
* fix: Improve java example by @nielsbasjes in https://github.com/apache/plc4x/pull/724
* build(deps): bump equalsverifier from 3.12.4 to 3.13 by @dependabot in https://github.com/apache/plc4x/pull/783
* build(deps): bump checker-qual from 3.29.0 to 3.30.0 by @dependabot in https://github.com/apache/plc4x/pull/781
* build(deps): bump BobAnkh/auto-generate-changelog from 1.2.3 to 1.2.4 by @dependabot in https://github.com/apache/plc4x/pull/780
* build(deps): bump commons-csv from 1.9.0 to 1.10.0 by @dependabot in https://github.com/apache/plc4x/pull/779
* build(deps): bump maven-enforcer-plugin from 3.1.0 to 3.2.1 by @dependabot in https://github.com/apache/plc4x/pull/777
* build(deps): bump jaxb-runtime from 4.0.1 to 4.0.2 by @dependabot in https://github.com/apache/plc4x/pull/789
* build(deps): bump byte-buddy from 1.12.22 to 1.12.23 by @dependabot in https://github.com/apache/plc4x/pull/788
* build(deps): bump calcite-core.version from 1.32.0 to 1.33.0 by @dependabot in https://github.com/apache/plc4x/pull/787
* build(deps): bump asciidoctorj-diagram from 2.2.3 to 2.2.4 by @dependabot in https://github.com/apache/plc4x/pull/786
* Feature/nifi integration address text by @QuanticPony in https://github.com/apache/plc4x/pull/755
* build(deps): bump groovy.version from 4.0.8 to 4.0.9 by @dependabot in https://github.com/apache/plc4x/pull/792
* build(deps): bump zip4j from 2.11.3 to 2.11.4 by @dependabot in https://github.com/apache/plc4x/pull/794
* build(deps): bump nifi.version from 1.19.1 to 1.20.0 by @dependabot in https://github.com/apache/plc4x/pull/793
* build(deps): bump nifi-nar-maven-plugin from 1.3.5 to 1.4.0 by @dependabot in https://github.com/apache/plc4x/pull/790
* build(deps): bump golang.org/x/tools from 0.5.0 to 0.6.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/800
* build(deps): bump equalsverifier from 3.13 to 3.13.1 by @dependabot in https://github.com/apache/plc4x/pull/799
* build(deps): bump netty-bom from 4.1.87.Final to 4.1.88.Final by @dependabot in https://github.com/apache/plc4x/pull/798
* fix(plc4j/connection-cache): fix issue with timing of thread by @hutcheb in https://github.com/apache/plc4x/pull/796
* build(deps): bump byte-buddy from 1.12.23 to 1.13.0 by @dependabot in https://github.com/apache/plc4x/pull/797
* Fix code gen concurrent modification by @hutcheb in https://github.com/apache/plc4x/pull/795
* build(deps): bump netty-bom from 4.1.88.Final to 4.1.89.Final by @dependabot in https://github.com/apache/plc4x/pull/804
* build(deps): bump maven-failsafe-plugin from 3.0.0-M8 to 3.0.0-M9 by @dependabot in https://github.com/apache/plc4x/pull/802
* build(deps): bump maven-surefire-plugin from 3.0.0-M8 to 3.0.0-M9 by @dependabot in https://github.com/apache/plc4x/pull/803
* build(deps): bump maven-invoker-plugin from 3.4.0 to 3.5.0 by @dependabot in https://github.com/apache/plc4x/pull/807
* build(deps): bump cyclonedx-maven-plugin from 2.7.4 to 2.7.5 by @dependabot in https://github.com/apache/plc4x/pull/806
* build(deps): bump maven-javadoc-plugin from 3.4.1 to 3.5.0 by @dependabot in https://github.com/apache/plc4x/pull/805
* build(deps): bump protobuf-java from 3.21.12 to 3.22.0 by @dependabot in https://github.com/apache/plc4x/pull/810
* build(deps): bump equalsverifier from 3.13.1 to 3.13.2 by @dependabot in https://github.com/apache/plc4x/pull/811
* build(deps): bump antlr.version from 4.11.1 to 4.12.0 by @dependabot in https://github.com/apache/plc4x/pull/813
* build(deps): bump github.com/gdamore/tcell/v2 from 2.5.4 to 2.6.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/817
* build(deps): bump byte-buddy from 1.13.0 to 1.14.0 by @dependabot in https://github.com/apache/plc4x/pull/816
* build(deps): bump checker-qual from 3.30.0 to 3.31.0 by @dependabot in https://github.com/apache/plc4x/pull/814
* build(deps): bump jsoup from 1.15.3 to 1.15.4 by @dependabot in https://github.com/apache/plc4x/pull/815
* build(deps): bump zip4j from 2.11.4 to 2.11.5 by @dependabot in https://github.com/apache/plc4x/pull/820
* build(deps): bump maven-assembly-plugin from 3.4.2 to 3.5.0 by @dependabot in https://github.com/apache/plc4x/pull/819
* Feature/nifi integration sink record processor and minor fixes by @QuanticPony in https://github.com/apache/plc4x/pull/809
* build(deps): bump google-java-format from 1.15.0 to 1.16.0 by @dependabot in https://github.com/apache/plc4x/pull/824
* build(deps): bump github.com/stretchr/testify from 1.8.1 to 1.8.2 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/825
* build(deps): bump jetty-util from 11.0.13 to 11.0.14 by @dependabot in https://github.com/apache/plc4x/pull/829
* build(deps-dev): bump json from 20220924 to 20230227 by @dependabot in https://github.com/apache/plc4x/pull/827
* build(deps): bump equalsverifier from 3.13.2 to 3.14 by @dependabot in https://github.com/apache/plc4x/pull/826
* build(deps): bump maven-compiler-plugin from 3.10.1 to 3.11.0 by @dependabot in https://github.com/apache/plc4x/pull/828
* build(deps): bump checker-qual from 3.31.0 to 3.32.0 by @dependabot in https://github.com/apache/plc4x/pull/830
* build(deps): bump protobuf-java from 3.22.0 to 3.22.1 by @dependabot in https://github.com/apache/plc4x/pull/833
* build(deps): bump byte-buddy from 1.14.0 to 1.14.1 by @dependabot in https://github.com/apache/plc4x/pull/834
* build(deps): bump golang.org/x/tools from 0.6.0 to 0.7.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/835
* build(deps): bump mockito.version from 5.1.1 to 5.2.0 by @dependabot in https://github.com/apache/plc4x/pull/837
* build(deps): bump github.com/fatih/color from 1.14.1 to 1.15.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/843
* build(deps): bump groovy.version from 4.0.9 to 4.0.10 by @dependabot in https://github.com/apache/plc4x/pull/840
* build(deps): bump logback.version from 1.4.5 to 1.4.6 by @dependabot in https://github.com/apache/plc4x/pull/846
* build(deps): bump byte-buddy from 1.14.1 to 1.14.2 by @dependabot in https://github.com/apache/plc4x/pull/848
* build(deps): bump calcite-core.version from 1.33.0 to 1.34.0 by @dependabot in https://github.com/apache/plc4x/pull/847
* build(deps): bump github.com/schollz/progressbar/v3 from 3.13.0 to 3.13.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/844
* build(deps): bump protobuf-java from 3.22.1 to 3.22.2 by @dependabot in https://github.com/apache/plc4x/pull/842
* build(deps): bump iot-device-client from 2.1.4 to 2.1.5 by @dependabot in https://github.com/apache/plc4x/pull/841
* build(deps): bump ecj from 3.32.0 to 3.33.0 by @dependabot in https://github.com/apache/plc4x/pull/853
* build(deps): bump nifi-nar-maven-plugin from 1.4.0 to 1.5.0 by @dependabot in https://github.com/apache/plc4x/pull/850
* build(deps): bump equalsverifier from 3.14 to 3.14.1 by @dependabot in https://github.com/apache/plc4x/pull/851
* build(deps): bump maven-failsafe-plugin from 3.0.0-M9 to 3.0.0 by @dependabot in https://github.com/apache/plc4x/pull/852
* Fix reading UTF-8 strings (from OPC UA nodes) by @Planet-X in https://github.com/apache/plc4x/pull/832
* build(deps): bump maven-release-plugin from 3.0.0-M7 to 3.0.0 by @dependabot in https://github.com/apache/plc4x/pull/857
* build(deps): bump asciidoctor-maven-plugin from 2.2.2 to 2.2.3 by @dependabot in https://github.com/apache/plc4x/pull/856
* build(deps): bump maven-surefire-plugin from 3.0.0-M9 to 3.0.0 by @dependabot in https://github.com/apache/plc4x/pull/855
* build(deps): bump slf4j.version from 2.0.6 to 2.0.7 by @dependabot in https://github.com/apache/plc4x/pull/854
* build(deps): bump netty-bom from 4.1.89.Final to 4.1.90.Final by @dependabot in https://github.com/apache/plc4x/pull/861
* build(deps): bump Saxon-HE from 12.0 to 12.1 by @dependabot in https://github.com/apache/plc4x/pull/859
* build(deps): bump swagger-annotations from 1.6.9 to 1.6.10 by @dependabot in https://github.com/apache/plc4x/pull/858
* build(deps): bump javafx.version from 19.0.2.1 to 20 by @dependabot in https://github.com/apache/plc4x/pull/860
* Revert "build(deps): bump ecj from 3.32.0 to 3.33.0" by @sruehl in https://github.com/apache/plc4x/pull/862
* build(deps-dev): bump commons-compress from 1.22 to 1.23.0 by @dependabot in https://github.com/apache/plc4x/pull/863
* build(deps): bump joda-time from 2.12.2 to 2.12.3 by @dependabot in https://github.com/apache/plc4x/pull/865
* build(deps): bump joda-time from 2.12.3 to 2.12.4 by @dependabot in https://github.com/apache/plc4x/pull/867
* build(deps): bump maven-resources-plugin from 3.3.0 to 3.3.1 by @dependabot in https://github.com/apache/plc4x/pull/868
* build(deps): bump github.com/gopacket/gopacket from 1.0.0 to 1.1.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/869
* build(deps): bump extra-enforcer-rules from 1.6.1 to 1.6.2 by @dependabot in https://github.com/apache/plc4x/pull/871
* build(deps): bump byte-buddy from 1.14.2 to 1.14.3 by @dependabot in https://github.com/apache/plc4x/pull/872
* build(deps): bump commons-configuration2 from 2.8.0 to 2.9.0 by @dependabot in https://github.com/apache/plc4x/pull/873
* build(deps): bump influxdb-client-java from 6.7.0 to 6.8.0 by @dependabot in https://github.com/apache/plc4x/pull/874
* build(deps): bump joda-time from 2.12.4 to 2.12.5 by @dependabot in https://github.com/apache/plc4x/pull/875
* build(deps): bump groovy.version from 4.0.10 to 4.0.11 by @dependabot in https://github.com/apache/plc4x/pull/877
* build(deps): bump maven-invoker-plugin from 3.5.0 to 3.5.1 by @dependabot in https://github.com/apache/plc4x/pull/878
* build(deps): bump cyclonedx-maven-plugin from 2.7.5 to 2.7.6 by @dependabot in https://github.com/apache/plc4x/pull/879
* build(deps): bump cmake-maven-plugin from 3.23.2-b1 to 3.25.2-b1 by @dependabot in https://github.com/apache/plc4x/pull/836
* build(deps): bump checker-qual from 3.32.0 to 3.33.0 by @dependabot in https://github.com/apache/plc4x/pull/881
* build(deps): bump jacoco-maven-plugin from 0.8.8 to 0.8.9 by @dependabot in https://github.com/apache/plc4x/pull/882
* build(deps): bump milo.version from 0.6.8 to 0.6.9 by @dependabot in https://github.com/apache/plc4x/pull/883
* build(deps): bump netty-bom from 4.1.90.Final to 4.1.91.Final by @dependabot in https://github.com/apache/plc4x/pull/884
* Remove SonarCloud cache and threads configuration and rely on default by @mpaladin in https://github.com/apache/plc4x/pull/885
* build(deps): bump github.com/spf13/cobra from 1.6.1 to 1.7.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/887
* build(deps): bump maven-enforcer-plugin from 3.2.1 to 3.3.0 by @dependabot in https://github.com/apache/plc4x/pull/886
* build(deps): bump byte-buddy from 1.14.3 to 1.14.4 by @dependabot in https://github.com/apache/plc4x/pull/889
* build(deps): bump golang.org/x/tools from 0.7.0 to 0.8.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/894
* build(deps): bump nifi.version from 1.20.0 to 1.21.0 by @dependabot in https://github.com/apache/plc4x/pull/893
* build(deps): bump actions/github-script from 6.4.0 to 6.4.1 by @dependabot in https://github.com/apache/plc4x/pull/891
* build(deps): bump asciidoctorj-diagram from 2.2.4 to 2.2.7 by @dependabot in https://github.com/apache/plc4x/pull/890
* build(deps): bump mockito.version from 5.2.0 to 5.3.0 by @dependabot in https://github.com/apache/plc4x/pull/895
* build(deps): bump protobuf-java from 3.22.2 to 3.22.3 by @dependabot in https://github.com/apache/plc4x/pull/897
* build(deps): bump hivemq-mqtt-client from 1.3.0 to 1.3.1 by @dependabot in https://github.com/apache/plc4x/pull/896
* build(deps): bump github.com/rs/zerolog from 1.29.0 to 1.29.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/899
* build(deps): bump jetty-util from 11.0.14 to 11.0.15 by @dependabot in https://github.com/apache/plc4x/pull/898
* build(deps): bump asciidoctorj from 2.5.7 to 2.5.8 by @dependabot in https://github.com/apache/plc4x/pull/901
* build(deps): bump cmake-maven-plugin from 3.25.2-b1 to 3.26.3-b1 by @dependabot in https://github.com/apache/plc4x/pull/902
* build(deps): bump cyclonedx-maven-plugin from 2.7.6 to 2.7.7 by @dependabot in https://github.com/apache/plc4x/pull/903
* build(deps): bump woodstox-core from 6.5.0 to 6.5.1 by @dependabot in https://github.com/apache/plc4x/pull/904
* Revert "build(deps): bump cmake-maven-plugin from 3.25.2-b1 to 3.26.3-b1" by @sruehl in https://github.com/apache/plc4x/pull/905
* build(deps): bump logback.version from 1.4.6 to 1.4.7 by @dependabot in https://github.com/apache/plc4x/pull/907
* build(deps): bump javafx.version from 20 to 20.0.1 by @dependabot in https://github.com/apache/plc4x/pull/906
* build(deps): bump cmake-maven-plugin from 3.25.2-b1 to 3.26.3-b1 by @dependabot in https://github.com/apache/plc4x/pull/908
* build(deps): bump mockito.version from 5.3.0 to 5.3.1 by @dependabot in https://github.com/apache/plc4x/pull/910
* build(deps): bump jackson.version from 2.14.2 to 2.15.0 by @dependabot in https://github.com/apache/plc4x/pull/909
* build(deps): bump github.com/libp2p/go-reuseport from 0.2.0 to 0.3.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/915
* build(deps): bump jacoco-maven-plugin from 0.8.9 to 0.8.10 by @dependabot in https://github.com/apache/plc4x/pull/914
* build(deps): bump cyclonedx-maven-plugin from 2.7.7 to 2.7.8 by @dependabot in https://github.com/apache/plc4x/pull/913
* build(deps): bump netty-bom from 4.1.91.Final to 4.1.92.Final by @dependabot in https://github.com/apache/plc4x/pull/912
* build(deps): bump junit.jupiter.version from 5.9.2 to 5.9.3 by @dependabot in https://github.com/apache/plc4x/pull/916
* Feat/s7ha by @glcj in https://github.com/apache/plc4x/pull/918
* Feat/s7ha by @glcj in https://github.com/apache/plc4x/pull/919
* Site: Small Scraper improvements, add mvn dependency, fix interval format, fix typo by @kubo44 in https://github.com/apache/plc4x/pull/917
* build(deps): bump google-java-format from 1.16.0 to 1.17.0 by @dependabot in https://github.com/apache/plc4x/pull/923
* build(deps): bump download-maven-plugin from 1.6.8 to 1.7.0 by @dependabot in https://github.com/apache/plc4x/pull/922
* build(deps): bump jsoup from 1.15.4 to 1.16.1 by @dependabot in https://github.com/apache/plc4x/pull/921
* build(deps): bump jakarta.activation-api from 2.1.1 to 2.1.2 by @dependabot in https://github.com/apache/plc4x/pull/920
* build(deps): bump checker-qual from 3.33.0 to 3.34.0 by @dependabot in https://github.com/apache/plc4x/pull/925
* build(deps): bump protobuf-java from 3.22.3 to 3.22.4 by @dependabot in https://github.com/apache/plc4x/pull/928
* Feat/profinet ip set by @hutcheb in https://github.com/apache/plc4x/pull/927
* build(deps): bump golang.org/x/net from 0.9.0 to 0.10.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/936
* build(deps): bump protobuf-java from 3.22.4 to 3.23.0 by @dependabot in https://github.com/apache/plc4x/pull/931
* build(deps): bump asciidoctorj-diagram from 2.2.7 to 2.2.8 by @dependabot in https://github.com/apache/plc4x/pull/934
* build(deps): bump error_prone_annotations from 2.18.0 to 2.19.0 by @dependabot in https://github.com/apache/plc4x/pull/933
* build(deps): bump golang.org/x/tools from 0.8.0 to 0.9.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/935
* build(deps): bump golang.org/x/tools from 0.9.0 to 0.9.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/937
* build(deps): bump maven-failsafe-plugin from 3.0.0 to 3.1.0 by @dependabot in https://github.com/apache/plc4x/pull/939
* build(deps): bump maven-surefire-plugin from 3.0.0 to 3.1.0 by @dependabot in https://github.com/apache/plc4x/pull/940
* build(deps): bump groovy.version from 4.0.11 to 4.0.12 by @dependabot in https://github.com/apache/plc4x/pull/942
* build(deps): bump error_prone_annotations from 2.19.0 to 2.19.1 by @dependabot in https://github.com/apache/plc4x/pull/941
* build(deps): bump buildnumber-maven-plugin from 3.0.0 to 3.1.0 by @dependabot in https://github.com/apache/plc4x/pull/944
* build(deps): bump build-helper-maven-plugin from 3.3.0 to 3.4.0 by @dependabot in https://github.com/apache/plc4x/pull/945
* build(deps): bump swagger-annotations from 1.6.10 to 1.6.11 by @dependabot in https://github.com/apache/plc4x/pull/947
* build(deps): bump cyclonedx-maven-plugin from 2.7.8 to 2.7.9 by @dependabot in https://github.com/apache/plc4x/pull/950
* build(deps): bump jackson.version from 2.15.0 to 2.15.1 by @dependabot in https://github.com/apache/plc4x/pull/951
* build(deps): bump protobuf-java from 3.23.0 to 3.23.1 by @dependabot in https://github.com/apache/plc4x/pull/953
* build(deps): bump github.com/stretchr/testify from 1.8.2 to 1.8.3 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/955
* build(deps): bump javacan-core from 3.2.2 to 3.2.3 by @dependabot in https://github.com/apache/plc4x/pull/956
* build(deps): bump antlr.version from 4.12.0 to 4.13.0 by @dependabot in https://github.com/apache/plc4x/pull/957
* build(deps): bump gmavenplus-plugin from 2.1.0 to 3.0.0 by @dependabot in https://github.com/apache/plc4x/pull/960
* build(deps): bump netty.version from 4.1.92.Final to 4.1.93.Final by @dependabot in https://github.com/apache/plc4x/pull/961
* build(deps): bump guava from 31.1-jre to 32.0.0-jre by @dependabot in https://github.com/apache/plc4x/pull/963
* build(deps): bump protobuf-java from 3.23.1 to 3.23.2 by @dependabot in https://github.com/apache/plc4x/pull/964
* build(deps): bump asciidoctor-maven-plugin from 2.2.3 to 2.2.4 by @dependabot in https://github.com/apache/plc4x/pull/965
* Feature/nifi integration record listener by @QuanticPony in https://github.com/apache/plc4x/pull/958
* build(deps): bump jackson.version from 2.15.1 to 2.15.2 by @dependabot in https://github.com/apache/plc4x/pull/966
* build(deps): bump github.com/stretchr/testify from 1.8.3 to 1.8.4 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/969
* build(deps): bump influxdb-client-java from 6.8.0 to 6.9.0 by @dependabot in https://github.com/apache/plc4x/pull/967
* build(deps): bump github.com/spf13/viper from 1.15.0 to 1.16.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/968
* build(deps): bump equalsverifier from 3.14.1 to 3.14.2 by @dependabot in https://github.com/apache/plc4x/pull/970
* build(deps): bump asciidoctorj from 2.5.8 to 2.5.9 by @dependabot in https://github.com/apache/plc4x/pull/974
* build(deps): bump extra-enforcer-rules from 1.6.2 to 1.7.0 by @dependabot in https://github.com/apache/plc4x/pull/975
* build(deps): bump golang.org/x/tools from 0.9.1 to 0.9.3 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/976
* build(deps): bump checker-qual from 3.34.0 to 3.35.0 by @dependabot in https://github.com/apache/plc4x/pull/977
* feat(plc4j) Better handling of timeouts in plc4j (#821). by @splatch in https://github.com/apache/plc4x/pull/822
* build(deps): bump byte-buddy from 1.14.4 to 1.14.5 by @dependabot in https://github.com/apache/plc4x/pull/978
* build(deps): bump asciidoctorj from 2.5.9 to 2.5.10 by @dependabot in https://github.com/apache/plc4x/pull/979
* Add support of PlcDINT for BigInteger in PlcValueHandler by @PatrykGala in https://github.com/apache/plc4x/pull/962
* build(deps): bump iot-device-client from 2.1.5 to 2.2.0 by @dependabot in https://github.com/apache/plc4x/pull/981
* build(deps): bump buildnumber-maven-plugin from 3.1.0 to 3.2.0 by @dependabot in https://github.com/apache/plc4x/pull/980
* Fix/www by @glcj in https://github.com/apache/plc4x/pull/984
* build(deps): bump guava from 32.0.0-jre to 32.0.1-jre by @dependabot in https://github.com/apache/plc4x/pull/983
* build(deps): bump nifi.version from 1.21.0 to 1.22.0 by @dependabot in https://github.com/apache/plc4x/pull/985
* Fix/www by @glcj in https://github.com/apache/plc4x/pull/987
* build(deps): bump ecj from 3.32.0 to 3.34.0 by @dependabot in https://github.com/apache/plc4x/pull/989
* build(deps): bump asciidoctorj-diagram from 2.2.8 to 2.2.9 by @dependabot in https://github.com/apache/plc4x/pull/986
* Revert "build(deps): bump ecj from 3.32.0 to 3.34.0" by @sruehl in https://github.com/apache/plc4x/pull/990
* build(deps): bump logback.version from 1.4.7 to 1.4.8 by @dependabot in https://github.com/apache/plc4x/pull/992
* build(deps): bump golang.org/x/net from 0.10.0 to 0.11.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/994
* build(deps): bump protobuf-java from 3.23.2 to 3.23.3 by @dependabot in https://github.com/apache/plc4x/pull/996
* build(deps): bump golang.org/x/tools from 0.9.3 to 0.10.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/998
* build(deps): bump mockito.version from 5.3.1 to 5.4.0 by @dependabot in https://github.com/apache/plc4x/pull/1000
* build(deps-dev): bump json from 20230227 to 20230618 by @dependabot in https://github.com/apache/plc4x/pull/1001
* build(deps): bump error_prone_annotations from 2.19.1 to 2.20.0 by @dependabot in https://github.com/apache/plc4x/pull/1002
* build(deps): bump nifi-nar-maven-plugin from 1.5.0 to 1.5.1 by @dependabot in https://github.com/apache/plc4x/pull/1003
* build(deps): bump milo.version from 0.6.9 to 0.6.10 by @dependabot in https://github.com/apache/plc4x/pull/1004
* build(deps): bump netty.version from 4.1.93.Final to 4.1.94.Final by @dependabot in https://github.com/apache/plc4x/pull/1005
* build(deps): bump xml-maven-plugin from 1.0.2 to 1.1.0 by @dependabot in https://github.com/apache/plc4x/pull/1006
* build(deps): bump equalsverifier from 3.14.2 to 3.14.3 by @dependabot in https://github.com/apache/plc4x/pull/1009
* build(deps): bump jSerialComm from 2.9.3 to 2.10.1 by @dependabot in https://github.com/apache/plc4x/pull/1010
* build(deps): bump guava from 32.0.1-jre to 32.1.0-jre by @dependabot in https://github.com/apache/plc4x/pull/1011
* build(deps): bump BobAnkh/auto-generate-changelog from 1.2.4 to 1.2.5 by @dependabot in https://github.com/apache/plc4x/pull/1014
* build(deps): bump guava from 32.1.0-jre to 32.1.1-jre by @dependabot in https://github.com/apache/plc4x/pull/1013
* build(deps): bump checker-qual from 3.35.0 to 3.36.0 by @dependabot in https://github.com/apache/plc4x/pull/1015
* build(deps): bump Saxon-HE from 12.1 to 12.3 by @dependabot in https://github.com/apache/plc4x/pull/1017
* build(deps): bump golang.org/x/net from 0.11.0 to 0.12.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1019
* build(deps): bump golang.org/x/tools from 0.10.0 to 0.11.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1020
* build(deps): bump protobuf-java from 3.23.3 to 3.23.4 by @dependabot in https://github.com/apache/plc4x/pull/1022
* build(deps): bump github.com/gopacket/gopacket from 1.1.0 to 1.1.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1023
* build(deps): bump jSerialComm from 2.10.1 to 2.10.2 by @dependabot in https://github.com/apache/plc4x/pull/1025
* build(deps): bump equalsverifier from 3.14.3 to 3.15 by @dependabot in https://github.com/apache/plc4x/pull/1026
* fix(plc4go/modbus): Delete elements in the loop, and the index is dec… by @hongjinlin in https://github.com/apache/plc4x/pull/1028
* build(deps): bump github.com/google/uuid from 1.1.2 to 1.3.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1029
* build(deps): bump javafx.version from 20.0.1 to 20.0.2 by @dependabot in https://github.com/apache/plc4x/pull/1034
* build(deps): bump netty.version from 4.1.94.Final to 4.1.95.Final by @dependabot in https://github.com/apache/plc4x/pull/1035
* build(deps): bump junit.jupiter.version from 5.9.3 to 5.10.0 by @dependabot in https://github.com/apache/plc4x/pull/1038
* Revert "build(deps): bump junit.jupiter.version from 5.9.3 to 5.10.0" by @sruehl in https://github.com/apache/plc4x/pull/1039
* build(deps): bump com.googlecode.maven-download-plugin:download-maven-plugin from 1.6.8 to 1.7.1 by @dependabot in https://github.com/apache/plc4x/pull/1040
* build(deps): bump com.influxdb:influxdb-client-java from 6.9.0 to 6.10.0 by @dependabot in https://github.com/apache/plc4x/pull/1043
* build(deps): bump netty.version from 4.1.95.Final to 4.1.96.Final by @dependabot in https://github.com/apache/plc4x/pull/1044
* build(deps): bump org.codehaus.mojo:properties-maven-plugin from 1.1.0 to 1.2.0 by @dependabot in https://github.com/apache/plc4x/pull/1045
* build(deps): bump junit.jupiter.version from 5.9.3 to 5.10.0 by @dependabot in https://github.com/apache/plc4x/pull/1041
* build(deps): bump github.com/rs/zerolog from 1.29.1 to 1.30.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1047
* build(deps): bump org.asciidoctor:asciidoctorj-diagram from 2.2.9 to 2.2.10 by @dependabot in https://github.com/apache/plc4x/pull/1048
* build(deps): bump com.fazecast:jSerialComm from 2.10.2 to 2.10.3 by @dependabot in https://github.com/apache/plc4x/pull/1049
* build(deps): bump tel.schich:javacan-core from 3.2.3 to 3.2.4 by @dependabot in https://github.com/apache/plc4x/pull/1050
* build(deps): bump golang.org/x/tools from 0.11.0 to 0.11.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1052
* build(deps): bump com.google.guava:guava from 32.1.1-jre to 32.1.2-jre by @dependabot in https://github.com/apache/plc4x/pull/1051
* build(deps): bump org.checkerframework:checker-qual from 3.36.0 to 3.37.0 by @dependabot in https://github.com/apache/plc4x/pull/1056
* build(deps): bump com.google.errorprone:error_prone_annotations from 2.20.0 to 2.21.0 by @dependabot in https://github.com/apache/plc4x/pull/1055
* build(deps): bump golang.org/x/net from 0.12.0 to 0.13.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1057
* Update StreamPipes integration on website by @bossenti in https://github.com/apache/plc4x/pull/1061
* build(deps): bump org.yaml:snakeyaml from 2.0 to 2.1 by @dependabot in https://github.com/apache/plc4x/pull/1062
* build(deps): bump com.google.errorprone:error_prone_annotations from 2.21.0 to 2.21.1 by @dependabot in https://github.com/apache/plc4x/pull/1063
* build(deps): bump logback.version from 1.4.8 to 1.4.9 by @dependabot in https://github.com/apache/plc4x/pull/1064
* build(deps): bump nl.jqno.equalsverifier:equalsverifier from 3.15 to 3.15.1 by @dependabot in https://github.com/apache/plc4x/pull/1058
* build(deps): bump golang.org/x/net from 0.13.0 to 0.14.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1065
* build(deps): bump golang.org/x/tools from 0.11.1 to 0.12.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1066
* build(deps): bump org.eclipse.jetty:jetty-util from 11.0.15 to 12.0.0 by @dependabot in https://github.com/apache/plc4x/pull/1067
* build(deps): bump com.google.protobuf:protobuf-java from 3.23.4 to 3.24.0 by @dependabot in https://github.com/apache/plc4x/pull/1069
* build(deps): bump logback.version from 1.4.9 to 1.4.11 by @dependabot in https://github.com/apache/plc4x/pull/1070
* build(deps): bump org.asciidoctor:asciidoctorj-diagram from 2.2.10 to 2.2.11 by @dependabot in https://github.com/apache/plc4x/pull/1071
* build(deps): bump net.bytebuddy:byte-buddy from 1.14.5 to 1.14.6 by @dependabot in https://github.com/apache/plc4x/pull/1072
* build(deps): bump github.com/libp2p/go-reuseport from 0.3.0 to 0.4.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1073
* build(deps): bump nifi.version from 1.23.0 to 1.23.1 by @dependabot in https://github.com/apache/plc4x/pull/1074
* build(deps): bump com.google.protobuf:protobuf-java from 3.24.0 to 3.24.1 by @dependabot in https://github.com/apache/plc4x/pull/1075
* build(deps): bump github.com/google/uuid from 1.3.0 to 1.3.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1076
* build(deps): bump mockito.version from 5.4.0 to 5.5.0 by @dependabot in https://github.com/apache/plc4x/pull/1077
* build(deps): bump nifi.version from 1.23.1 to 1.23.2 by @dependabot in https://github.com/apache/plc4x/pull/1078
* build(deps): bump netty.version from 4.1.96.Final to 4.1.97.Final by @dependabot in https://github.com/apache/plc4x/pull/1079
* build(deps): bump com.google.protobuf:protobuf-java from 3.24.1 to 3.24.2 by @dependabot in https://github.com/apache/plc4x/pull/1084
* build(deps): bump net.bytebuddy:byte-buddy from 1.14.6 to 1.14.7 by @dependabot in https://github.com/apache/plc4x/pull/1085
* build(deps): bump com.hivemq:hivemq-mqtt-client from 1.3.1 to 1.3.2 by @dependabot in https://github.com/apache/plc4x/pull/1087
* build(deps): bump org.yaml:snakeyaml from 2.1 to 2.2 by @dependabot in https://github.com/apache/plc4x/pull/1083
* build(deps): bump org.eclipse.jetty:jetty-util from 12.0.0 to 12.0.1 by @dependabot in https://github.com/apache/plc4x/pull/1090
* Feature/iec 60870 5 104 by @chrisdutz in https://github.com/apache/plc4x/pull/1088
* build(deps): bump golang.org/x/net from 0.14.0 to 0.15.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1096
* build(deps): bump antlr.version from 4.13.0 to 4.13.1 by @dependabot in https://github.com/apache/plc4x/pull/1094
* build(deps): bump actions/checkout from 3 to 4 by @dependabot in https://github.com/apache/plc4x/pull/1093
* build(deps): bump com.google.protobuf:protobuf-java from 3.24.2 to 3.24.3 by @dependabot in https://github.com/apache/plc4x/pull/1098
* build(deps): bump org.checkerframework:checker-qual from 3.37.0 to 3.38.0 by @dependabot in https://github.com/apache/plc4x/pull/1091
* build(deps): bump slf4j.version from 2.0.7 to 2.0.9 by @dependabot in https://github.com/apache/plc4x/pull/1092
* plc4j-driver-opcua: Fix incorrectly handled GUID tags by @takraj in https://github.com/apache/plc4x/pull/1099
* plc4j-driver-opcua: Add support for PlcUsernamePasswordAuthentication by @takraj in https://github.com/apache/plc4x/pull/1107
* plc4j-driver-opcua: Adapt error handling of reads, to writes and subscriptions by @takraj in https://github.com/apache/plc4x/pull/1108
* build(deps): bump netty.version from 4.1.97.Final to 4.1.98.Final by @dependabot in https://github.com/apache/plc4x/pull/1109
* build(deps): bump net.bytebuddy:byte-buddy from 1.14.7 to 1.14.8 by @dependabot in https://github.com/apache/plc4x/pull/1110
* build(deps): bump javafx.version from 20.0.2 to 21 by @dependabot in https://github.com/apache/plc4x/pull/1113
* build(deps): bump org.sonarsource.scanner.maven:sonar-maven-plugin from 3.9.1.2184 to 3.10.0.2594 by @dependabot in https://github.com/apache/plc4x/pull/1111
* build(deps): bump bouncycastle.version from 1.75 to 1.76 by @dependabot in https://github.com/apache/plc4x/pull/1112
* Closes #801 by @DmitriiMukhin in https://github.com/apache/plc4x/pull/888
* build(deps): bump nl.jqno.equalsverifier:equalsverifier from 3.15.1 to 3.15.2 by @dependabot in https://github.com/apache/plc4x/pull/1114
* build(deps): bump com.google.errorprone:error_prone_annotations from 2.21.1 to 2.22.0 by @dependabot in https://github.com/apache/plc4x/pull/1115
* build(deps): bump org.asciidoctor:asciidoctorj-diagram from 2.2.11 to 2.2.13 by @dependabot in https://github.com/apache/plc4x/pull/1116
* build(deps): bump jakarta.xml.bind:jakarta.xml.bind-api from 4.0.0 to 4.0.1 by @dependabot in https://github.com/apache/plc4x/pull/1118
* build(deps): bump github.com/rs/zerolog from 1.30.0 to 1.31.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/1119
* build(deps): bump netty.version from 4.1.98.Final to 4.1.99.Final by @dependabot in https://github.com/apache/plc4x/pull/1120
* build(deps): bump org.codehaus.gmavenplus:gmavenplus-plugin from 3.0.0 to 3.0.2 by @dependabot in https://github.com/apache/plc4x/pull/1121
* Nifi integration revision by @QuanticPony in https://github.com/apache/plc4x/pull/1122
* build(deps): bump com.microsoft.azure.sdk.iot:iot-device-client from 2.2.0 to 2.3.0 by @dependabot in https://github.com/apache/plc4x/pull/1125

## New Contributors
* @modraedlau made their first contribution in https://github.com/apache/plc4x/pull/668
* @nielsbasjes made their first contribution in https://github.com/apache/plc4x/pull/717
* @SteinOv made their first contribution in https://github.com/apache/plc4x/pull/711
* @QuanticPony made their first contribution in https://github.com/apache/plc4x/pull/752
* @Planet-X made their first contribution in https://github.com/apache/plc4x/pull/832
* @mpaladin made their first contribution in https://github.com/apache/plc4x/pull/885
* @kubo44 made their first contribution in https://github.com/apache/plc4x/pull/917
* @PatrykGala made their first contribution in https://github.com/apache/plc4x/pull/962
* @bossenti made their first contribution in https://github.com/apache/plc4x/pull/1061
* @takraj made their first contribution in https://github.com/apache/plc4x/pull/1099
* @DmitriiMukhin made their first contribution in https://github.com/apache/plc4x/pull/888

**Full Changelog**: https://github.com/apache/plc4x/compare/rel/0.10...rel/0.11

### Feature

- integration/nifi:
  - Various improvements for Nifi integration ([4156cc9](https://github.com/apache/plc4x/commit/4156cc9f1e65485ad7d21d3630835603727a5104)) ([#1122](https://github.com/apache/plc4x/pull/1122))

- general:
  - Implemented the PLC4J Ping API for a number of protocols:
- ADS
- EIP
- KNX
- Mock
- Modbus
- Simulated ([4abbf26](https://github.com/apache/plc4x/commit/4abbf260f071139811c4d57c891ceebb0b0b0a9f))
  - Added a modbuspal project file to the repo, which can be used by the Manual modbus test. ([14531ee](https://github.com/apache/plc4x/commit/14531ee3fe9f59a6f6a1c0b7812e1ee8a1ab1006))
  - Implemented the writing of arrays of coils in one request for Modbus ([48a29ca](https://github.com/apache/plc4x/commit/48a29ca1442a0074635a3223b0982e572322465a))
  - Added a modbuspal project file to the repo, which can be used by the Manual modbus test. ([4febbb6](https://github.com/apache/plc4x/commit/4febbb6be2bb52cb53a684ad8acfaf80a4509e3a))
  - Added some more test-cases to the manual modbus test ([d2adee9](https://github.com/apache/plc4x/commit/d2adee9acfc6e500265a84ffab760f5c234fb442))
  - Added methods to manually remove connections from the cache. ([c8e7fc9](https://github.com/apache/plc4x/commit/c8e7fc9913a21b5489a987e1a14b38d613529b17))
  - Added support in the ByteBased read- and write-buffers for WINDOWS-1252 encoding ([002da8a](https://github.com/apache/plc4x/commit/002da8a2d3cbef0bf450873b56226ac00a7fe114))
  - Added some commented out code that dumps the decrypted parts of ETS6+ files (helps debug ETS parsing) ([688c85e](https://github.com/apache/plc4x/commit/688c85e54c46f7c60d89bc81d18f6dfcfa9c7557))
  - Added support for ETS6.1 ([c268ae2](https://github.com/apache/plc4x/commit/c268ae260aa705f20dca17d45c82aa7fcf946dd7))
  - expose stringers in API ([cfad9d1](https://github.com/apache/plc4x/commit/cfad9d1370113dd92e30c02bd855f353a6ed53c4))
  - update vendors ([ea12846](https://github.com/apache/plc4x/commit/ea128461125e31f96502ea1b755d2c2a6b2b88ca))
  - update manufacturers ([13920f0](https://github.com/apache/plc4x/commit/13920f025be4add0e0a5aeaa4eb9b57fca922ef3))

- plc4j/opcua:
  - Adapt error handling of reads, to writes and subscriptions (#1108) ([bfee8d3](https://github.com/apache/plc4x/commit/bfee8d3e1f08cf6df99920da08489333b0e2c3f1)) ([#1108](https://github.com/apache/plc4x/pull/1108))
  - Add support for PlcUsernamePasswordAuthentication (#1107) ([41d82d8](https://github.com/apache/plc4x/commit/41d82d8cef71944c5b48fe7f453ee00899e8fb97)) ([#1107](https://github.com/apache/plc4x/pull/1107))
  - moved test to driver ([33aa723](https://github.com/apache/plc4x/commit/33aa723b5b66c1ea522a4ac82f9e58d4bbb72099))
  - Add support of PlcDINT for BigInteger in PlcValueHandler (#962) ([bd13295](https://github.com/apache/plc4x/commit/bd1329579af0b7df5f3cd9104d614bc3ccbffcf9)) ([#962](https://github.com/apache/plc4x/pull/962))
  - update node id services ([0a19df2](https://github.com/apache/plc4x/commit/0a19df2e749df6adb4c9d746394a3afc27493d0d))

- plc4j/iec-60870:
  - Fixed the problem decoding subsequent incoming messages. ([f4a6891](https://github.com/apache/plc4x/commit/f4a689189d3c9e93223e84b65a99522331464d8d))
  - Continued implementing the portocol. ([20246c7](https://github.com/apache/plc4x/commit/20246c7f8aa0260a44d68ba8b9006895511cf8a6))
  - Continued implementing the portocol and creating the ParserSerializer testsuite. ([a564eeb](https://github.com/apache/plc4x/commit/a564eebf5787b95b3830193765e460a2432cd802))

- knx:
  - update vendors ([39ff22f](https://github.com/apache/plc4x/commit/39ff22f208cc64d347e9fbdf2dc49cb81ef4ff52))
  - update vendor ([c29ea5b](https://github.com/apache/plc4x/commit/c29ea5b7e06dc05460e1d7893ef64f677038fc78))
  - update vendors ([90e30e3](https://github.com/apache/plc4x/commit/90e30e30ecb78dd649f9d91bebe1c3487cede3d4))
  - update vendors ([f02c4c5](https://github.com/apache/plc4x/commit/f02c4c53a64212945dec2730029aa17d654dc5f7))

- plc4j/spi:
  - add protocolCode/transportCode and transportConfig as injectable parameters ([943e4a1](https://github.com/apache/plc4x/commit/943e4a1d415321cfbbf7a6a11d9f92ee33c65626))

- plc4go:
  - reworked default codec sleep time to a min 10ms latency ([9d62036](https://github.com/apache/plc4x/commit/9d620361cf0e0e9c8df13843020b29d42fd4b5bc))
  - upgrade to golang 1.20.7 ([953d07a](https://github.com/apache/plc4x/commit/953d07a32466b756193894aa1ab56e766bab2ec5))
  - add option WithTraceTransportInstance to limit tracing of test.TransportInstance ([16ea3d3](https://github.com/apache/plc4x/commit/16ea3d369395f33dade70ef8184b39b9ec88e099))
  - use structured loggin when possible ([95c4983](https://github.com/apache/plc4x/commit/95c4983455ce553b789a6edc7b4bf34d692de4f8))
  - implement unsubscription requests ([dcc5edb](https://github.com/apache/plc4x/commit/dcc5edba2c96e361a5811a879fffc82ad76a5876))
  - expose PlcSubscriptionType and add type assertion for default implementations ([fbf49c1](https://github.com/apache/plc4x/commit/fbf49c159b8bfb001c5ca4ff01931f433ce7d264))
  - improved options handling ([d2a14ca](https://github.com/apache/plc4x/commit/d2a14cac943ba6fdfeedd9a2f8d7fc1abe04cec0))
  - define custom Logger and tracer for bacnet ([9bb5a88](https://github.com/apache/plc4x/commit/9bb5a88be7791d0e9b877cd9cc04abff5968f419))
  - add receive timeout option ([5eca784](https://github.com/apache/plc4x/commit/5eca78479cd542acbf0478f87606adf47798e6ea))
  - added new options to replace global constants ([adeac04](https://github.com/apache/plc4x/commit/adeac0498bc45aca675ca1e55c4d6ae58f42ad5f))
  - add io.Closer to PlcDriverManager and driver ([15d8e66](https://github.com/apache/plc4x/commit/15d8e665677fca0f845c37fc19ac66d2133c03a0))
  - hook in WithPassLoggerToModel option ([0a77277](https://github.com/apache/plc4x/commit/0a77277f7998e5e1664d06d7ab37fbe0731e7eb3))
  - new config.WithPassLoggerToModel option ([54dbdec](https://github.com/apache/plc4x/commit/54dbdec05f6317b77643fb7c22623d57ce2a64aa))
  - always pass context to static helpers ([16b47a4](https://github.com/apache/plc4x/commit/16b47a45fac2947ef143c79a15d098695e8b6765))
  - pass context into models whenever possible ([21ed91d](https://github.com/apache/plc4x/commit/21ed91df54cfdef017fef3dc5f26e93a5013dca0))
  - improve logging on browse ([6fa5bc1](https://github.com/apache/plc4x/commit/6fa5bc147f8ee0687ec994e36579d620329fbd28))
  - improve logging for subscription ([7615701](https://github.com/apache/plc4x/commit/7615701e13dddd4138ea04b4fca1d721afda5b73))
  - update to go 1.20.5 ([84be4cb](https://github.com/apache/plc4x/commit/84be4cbb1f67cc6f37b4dc73094849022e465ff9))
  - ensure right logger is used when creating a transport instance ([9dd16ed](https://github.com/apache/plc4x/commit/9dd16ed55a1f5c3bc971ff630180cd1aec56f651))
  - add new WithCustomLogger option ([3193782](https://github.com/apache/plc4x/commit/3193782fb966cabfd7877bee3840c5d604e518b7))
  - upgrade go version ([9e842d7](https://github.com/apache/plc4x/commit/9e842d72d5acf927d268e24e8e7e14d2fa8c6ab7))
  - update enum interfaces ([33f08a3](https://github.com/apache/plc4x/commit/33f08a3c9e8527d61fd60becc8267f80c14605b8))
  - streamline connect calls ([e8462ca](https://github.com/apache/plc4x/commit/e8462caecea21e97a6c56498d312ac9582707b47))
  - improve logging of discoverer ([53aa21b](https://github.com/apache/plc4x/commit/53aa21b845b71083d2d2aac8b7f04bc7876e6fda))
  - upgrade to golang 1.19 ([e3030a4](https://github.com/apache/plc4x/commit/e3030a49ac29db8b5339fca492bf5f704ac34846))
  - added generic min method ([4d8ead9](https://github.com/apache/plc4x/commit/4d8ead97043bca22f45eb313f322073273bcc5df))

- plc4go/opcua:
  - port fireDiscoverEvent to plc4goa ([3d66b11](https://github.com/apache/plc4x/commit/3d66b1189742205da773a1973b212585e39ac856))
  - add a bit more logging ([da3b11e](https://github.com/apache/plc4x/commit/da3b11e09f5f3c551e1f628780035a5bcf009058))
  - fix issues with test ([265fdd1](https://github.com/apache/plc4x/commit/265fdd1e146695fb3b943dfba0bc2e8af0363354))
  - add more stringers to structs ([7bb17de](https://github.com/apache/plc4x/commit/7bb17deee22e0932e962c261ca5cf4952f3641ae))
  - implement unsubscription ([5c014a0](https://github.com/apache/plc4x/commit/5c014a0536e7e3e37fbf3eecc1eb3c260a3b9883))
  - work on subscribing ([9664025](https://github.com/apache/plc4x/commit/966402586727a63138084e6fe44da4e1b11bd857))
  - fix stringer ([b7ad5c1](https://github.com/apache/plc4x/commit/b7ad5c1817049af911d0df05157548e6db0c6a08))
  - add writer ([b74e9f3](https://github.com/apache/plc4x/commit/b74e9f334622ff92efe038ab27b6d443ec463516))
  - work on encryption part ([7c21ea8](https://github.com/apache/plc4x/commit/7c21ea88ce2db4e36b4f9100c08d9eee521ed51e))
  - implement onDisconnect and onDiscover ([431472e](https://github.com/apache/plc4x/commit/431472ede50643f96b794e72d4f1e6def24fa59d))
  - implement keepalive and connect event ([fb1a6d6](https://github.com/apache/plc4x/commit/fb1a6d6bab895bb6083758a9456dea91b906a9af))
  - some progress on secure channel ([36673bd](https://github.com/apache/plc4x/commit/36673bdcbcd0cd3d8306eafc74c095f852b91612))
  - implemented first protocol stub ([461340a](https://github.com/apache/plc4x/commit/461340a1c6e0849af34c1a9629603be0c791ce32))
  - generate models for OPC-UA ([490dfc6](https://github.com/apache/plc4x/commit/490dfc67b9972adc331f7e59d258e911a63bd33d))

- plc4go/plc4xbrowser:
  - fix multiline key value output for logger ([76ae69a](https://github.com/apache/plc4x/commit/76ae69adfbfb02a0e57604bb469371e0ba76852f))
  - add support for opcua ([2fea480](https://github.com/apache/plc4x/commit/2fea4804601ddcbab52a6235784b36cadf5801cd))

- opcua:
  - add support for OpcuaMessageError ([0dff535](https://github.com/apache/plc4x/commit/0dff53557827fabb9189d61ad1d0193c43be1135))
  - add support for OpcuaMessageError ([21e76ba](https://github.com/apache/plc4x/commit/21e76ba626d67fff4aa749ee06e7fbe0577199c7))
  - add support for OpcuaMessageError ([b88d4eb](https://github.com/apache/plc4x/commit/b88d4eb3002f898b264bebc1cfed081980d64ce5))

- bacnet:
  - update vendor ([333b72b](https://github.com/apache/plc4x/commit/333b72bf47d002b09ecfc4b67928b29a68695f33))
  - update vendor ([b7774ae](https://github.com/apache/plc4x/commit/b7774aeb903551cb35c6db542b840b97f5ac8c84))
  - update vendors ([faf7dc1](https://github.com/apache/plc4x/commit/faf7dc160e258e077d1ff5a28c64f817dfc3c46e))
  - vendor update ([5ce69b5](https://github.com/apache/plc4x/commit/5ce69b511f4599d1ae15ec1d32af4153e5446c7e))
  - update vendors ([4f51eef](https://github.com/apache/plc4x/commit/4f51eef495b107a05766d1ef913cb7051d0b7c2e))
  - update vendors ([8eeec1e](https://github.com/apache/plc4x/commit/8eeec1e33b6f4aa997f753294738b546d899820f))
  - add more info to the enums to access values from code ([807419e](https://github.com/apache/plc4x/commit/807419e9e5c0822b40b372b99ad922dffdd96c92))

- plc4go/tools:
  - add nil check for stringers ([d343a1b](https://github.com/apache/plc4x/commit/d343a1b020e2c829f38f1e840f12d03eb7d68e27))
  - reworked alpha support to use defer to always unlock lock ([d7c7fe9](https://github.com/apache/plc4x/commit/d7c7fe9fd69c62272e842dfd64fa079478c9bb1f))
  - add locker support to gen.go ([f6d53d8](https://github.com/apache/plc4x/commit/f6d53d88f0bba2cde29411bfe798cc4ef39c100d))
  - add support for atomic.Pointer to plc4xgenerator ([915b6b5](https://github.com/apache/plc4x/commit/915b6b54089bcd157ac9103f460b85dd1e0a241d))
  - add plc4xlicenser to add ASF headers ([26aee4f](https://github.com/apache/plc4x/commit/26aee4f4bac66ab2c907cd67c05466dc6690300b))

- plc4go/gen:
  - add support for non string keyed maps ([07130ed](https://github.com/apache/plc4x/commit/07130ed03ac441565638f3c4c7c9b91721c8f7e8))
  - add support for []byte ([2c34096](https://github.com/apache/plc4x/commit/2c34096fe9b11201630b1b64451d5161c180c998))
  - added byte support ([d06433b](https://github.com/apache/plc4x/commit/d06433b7b5a8094f8278d6f6a8661e0f06bc0207))
  - added couple of missing features ([2931bd6](https://github.com/apache/plc4x/commit/2931bd6df94db3390b446c71a647f3184aabe082))
  - add support for chan and func fields ([f347cfc](https://github.com/apache/plc4x/commit/f347cfc3831162ee816152592522fbc95e6ef469))

- codegen/plc4go:
  - store virtual field to local field for serialization ([5bf72f2](https://github.com/apache/plc4x/commit/5bf72f2b0d2aee4acd7991c1daa56c77120eb05c))
  - add support for abstarct field ([a2b6d6b](https://github.com/apache/plc4x/commit/a2b6d6b20d89b5ba8556ebdf41b6be96d9aae1b4))

- plc4go/spi:
  - log when a expectation was created ([2430acf](https://github.com/apache/plc4x/commit/2430acf6ac93dfab06355c8f546edd050298e36c))
  - add a uuid to expectations to better understand what is going on ([3a5433d](https://github.com/apache/plc4x/commit/3a5433de2ea93eb4c1b855ff2789dbb95e6183c2))
  - improve output from DriverTestRunner ([6000069](https://github.com/apache/plc4x/commit/6000069a02ffbe83ce2438e1cd1a5f45730e2883))
  - added remaining options as env to test utils ([d993e8c](https://github.com/apache/plc4x/commit/d993e8c61f90b0346bbe28b933ffe4ce73a0e7fb))
  - added env test util option for traceExecutorWorkers ([9a21a7e](https://github.com/apache/plc4x/commit/9a21a7e03e8bf2a89e353ddd90c5a7875c688e28))
  - allow overriding of options by always use the last option found. ([89d210f](https://github.com/apache/plc4x/commit/89d210f4dda0ba557e8deeb19efe477a85a4eaed))
  - refined logging of Default Connection ([e6c897c](https://github.com/apache/plc4x/commit/e6c897c97a59060cc7074f58ee58198e73b7ab52))
  - time custom message handling and don't block when default message channel is not being drained ([8e0f219](https://github.com/apache/plc4x/commit/8e0f2193356967a11df6e3e0b160b458aa8152b9))
  - improve logging for request transaction ([65796b0](https://github.com/apache/plc4x/commit/65796b03e9dbba26ab533ecafb3703e727f15d04))
  - test transport should now properly output the url ([b8c4bf2](https://github.com/apache/plc4x/commit/b8c4bf24c86a7d7d54121e6a792e30c10d2c7767))
  - added more Stringer implementations ([caa9718](https://github.com/apache/plc4x/commit/caa9718473faa3ebb4fc71951fd3cf4ead000c17))
  - use more local loggers ([7c14c99](https://github.com/apache/plc4x/commit/7c14c99a104e0ee82db390f2c02b7befa1f8c7c1))
  - Introduce new WithCustomLogger option ([b330c7c](https://github.com/apache/plc4x/commit/b330c7ce33647b379089d23a91d5788bc7b98a27))
  - implement GetConnectionUrl for options ([ab8bfd8](https://github.com/apache/plc4x/commit/ab8bfd8a17acf890e77f779b3dec6ec230f63ac3))
  - introduce interfaces for request transaction manager ([ab70d1b](https://github.com/apache/plc4x/commit/ab70d1b11a264209515e4f1c9e7c6b03d68ccbd8))
  - improve string outputs of transports ([a618193](https://github.com/apache/plc4x/commit/a6181936a52af9c75b0ab4c21daa68f6bae91244))
  - initial dynamic executor for worker pool ([ea79a34](https://github.com/apache/plc4x/commit/ea79a34ba0b718dab1432b4ea672719ccfadf94c))
  - use atomic.Bool for state changes on WorkerPool ([48c398c](https://github.com/apache/plc4x/commit/48c398c916a4dd07f7ee9024fc11b19db3330664))
  - added DiffHex and hex highlighting ([34ac8a1](https://github.com/apache/plc4x/commit/34ac8a13420fa2147e69726e0591714b80a3feb4))
  - use pre-allocated byte arrays for default writing ([5c28603](https://github.com/apache/plc4x/commit/5c28603020fd0978a0b1cd6b7b3a7e92ac74ef74))

- plc4go/eip:
  - output message before to better understand test
failure ([35bf97a](https://github.com/apache/plc4x/commit/35bf97a71c3d9c7865ba7821c5b0f0952431ef08))

- plc4j/bacnet:
  - update vendor ([fa268e8](https://github.com/apache/plc4x/commit/fa268e8586075cbe79605af69f445023032049bc))

- plc4j/knx:
  - update vendor ([50184f6](https://github.com/apache/plc4x/commit/50184f67470f2e82d411a5e6b5a1e87dc89a70a8))
  - update manufacturers ([c623c80](https://github.com/apache/plc4x/commit/c623c80f63ddab6a79a8c04d4be75fd8b855ca09))

- plc4go/knx:
  - update vendor ([9fe80f9](https://github.com/apache/plc4x/commit/9fe80f9cfea910bfffdda124227c48ee9347d51a))
  - update manufacturers ([36e4f4c](https://github.com/apache/plc4x/commit/36e4f4c5394cfa751d60b4170624348b57dbb5ff))
  - update manufacturer ([a96b8d2](https://github.com/apache/plc4x/commit/a96b8d24a949a7f77b20029c0101e4cc5700f851))

- plc4go/bacnet:
  - update vendor ([405c8f8](https://github.com/apache/plc4x/commit/405c8f83edd2c6fd6906c42beceb8882a2450088))
  - implement BIPNetworkApplication ([1860ffe](https://github.com/apache/plc4x/commit/1860ffef4734aa8c47fb393e62450bdf3673a68c))
  - implement BIPForeignApplication ([f655587](https://github.com/apache/plc4x/commit/f655587bc2c2a11e421f37b739f597646d68e1d6))
  - basic mapping to reader/writer ([55f7913](https://github.com/apache/plc4x/commit/55f7913ed320147b91eb1d8f23d8081d6262ea83))
  - fix some open issues regarding task processing ([eb59f7c](https://github.com/apache/plc4x/commit/eb59f7c1113c2e068a547910b6cc93d604bac2fc))
  - basic comm working ([574dd3f](https://github.com/apache/plc4x/commit/574dd3f633a8c256f03bb35554eb5b517275540a))
  - first comm went through new stack ([2265a6a](https://github.com/apache/plc4x/commit/2265a6ad5b4b947c2c66d027cd7028796d8c197d))
  - port task manager ([b93d2e4](https://github.com/apache/plc4x/commit/b93d2e48ecbbeba7cccaddc4d0fa436b2795ca81))
  - update udp comms ([86b07fb](https://github.com/apache/plc4x/commit/86b07fb9e9e21dfd7d6801013b60e04750c146aa))
  - implemented BIPSimple ([198882b](https://github.com/apache/plc4x/commit/198882b0399ba4e1d57094d7ebabdf889b9f102e))
  - progress on network stack ([09796b0](https://github.com/apache/plc4x/commit/09796b07b20a52f4cf244e798a3a5579a48b9ded))
  - progress on IOQController ([b367da7](https://github.com/apache/plc4x/commit/b367da732b874da2747004dfeeb995eb2d397f54))
  - use upstream device info cache ([83606be](https://github.com/apache/plc4x/commit/83606befb06f9012aa68339dd328d068a83b9a54))
  - added building structure to BIPSimpleApplication ([e43a48f](https://github.com/apache/plc4x/commit/e43a48f12eb40b911522feffbc6547c9fb4a27c2))
  - partial port of application layer, application module ([94e73b0](https://github.com/apache/plc4x/commit/94e73b0b5f4f0190e09cf1f1f8c4821f3cc807b2))
  - Client, Server, ServiceAccessPoint and ApplicationServiceElement ([9b2e27f](https://github.com/apache/plc4x/commit/9b2e27f09fb8f6fa56604f35fd11e70d7b751e23))
  - ported StateMachineAccessPoint ([a1a7746](https://github.com/apache/plc4x/commit/a1a7746cb2d14e1bacf27dca92d678444ddf3e2f))
  - ported ClientSSM and ServerSSM ([21776cc](https://github.com/apache/plc4x/commit/21776cc95f97b30bd7f509fbca78e7690b4d6df5))
  - partial transaction state machine ([863eb4e](https://github.com/apache/plc4x/commit/863eb4e22a6e960c630ce144232f80786b773696))
  - initial skeleton of TransactionStateMachine ([4955117](https://github.com/apache/plc4x/commit/49551174ac54c88b4a93c70ae4697c3ca6ab8444))

- plc4go/cbus:
  - improved logging ([04662cf](https://github.com/apache/plc4x/commit/04662cf6d5dc35daa389460cb51e76b32eea6aa7))
  - more log details on the connection process ([dfe1459](https://github.com/apache/plc4x/commit/dfe1459f024ed101ef353149cb6033720eeb46a3))
  - indicate if we handle a subscription event ([cccdf75](https://github.com/apache/plc4x/commit/cccdf758268a029da1dc5b64b7a43c22c0542ad6))
  - added more Stringer implementations ([5422176](https://github.com/apache/plc4x/commit/54221760969bf088eb07df4eb8a46f933dff6ca0))
  - improve logging of discoverer ([14ea38c](https://github.com/apache/plc4x/commit/14ea38c66e49275938fe3fe14e5f8bdc33a70a6d))
  - improve logging of discoverer ([61045e6](https://github.com/apache/plc4x/commit/61045e6eb58e8f240617d35ee44b925b5796b21e))
  - add more tracing, increase receive timeout ([f66dd42](https://github.com/apache/plc4x/commit/f66dd42324682acf0eac7dae1d5670b2945d626d))
  - add more tracing, increase receive timeout ([39d5e59](https://github.com/apache/plc4x/commit/39d5e594cb2ea14fe0fbf999649aa2c240c7d88e))
  - expose bridge addressing in tags ([65414b2](https://github.com/apache/plc4x/commit/65414b2e741e2c3cfd4aef744ac328976d93d96c))
  - implemented bridge support in message mapper ([3ec51ec](https://github.com/apache/plc4x/commit/3ec51ecc82acbbd4ec974a446d00a362d40d6078))
  - prepare bridge support ([74976da](https://github.com/apache/plc4x/commit/74976dae913edf20cf9c3d0a14909f95fa9b7cc0))

- plc4j:
  - feat(plc4j) Better handling of timeouts in plc4j (#821). ([efbb79c](https://github.com/apache/plc4x/commit/efbb79c58b18cceab52d426db34411315a4713b9)) ([#822](https://github.com/apache/plc4x/pull/822))
  - update vendors ([af61985](https://github.com/apache/plc4x/commit/af6198531c6c063b40cf382e83ba09d832de68ef))
  - update vendors ([d8c332f](https://github.com/apache/plc4x/commit/d8c332f044bab4adb7dfd3b3289614b8a3ab160a))

- eip:
  - Started adding auto-discovery to the EIP driver. ([974138e](https://github.com/apache/plc4x/commit/974138ec4ac34f2698cc5a15d34636cf3275656b))

- plc4py:
  - Working through the read buffer template section ([9d8ff34](https://github.com/apache/plc4x/commit/9d8ff340650674252b20d2987cdfb885dde887de))
  - Finished first round of the read buffer. There are issues with it though. ([ce94334](https://github.com/apache/plc4x/commit/ce943349d81ff7f73f4412b60352fa4a5a5b0464))
  - Change data types that the crc checks in Modbus use to native types ([7529d84](https://github.com/apache/plc4x/commit/7529d84b6b66c01cca7dcbabcbc6bb07aeccd66f))
  - Write Buffer now accepts only python native data types ([70f8f3b](https://github.com/apache/plc4x/commit/70f8f3b5058abafc5091ad9a7b675edf8be8a163))
  - Implement Tracer for Python Template Helper ([a07ad26](https://github.com/apache/plc4x/commit/a07ad26acc09e2bff7add182b67f3ab559cf8318))
  - remove start_pos variable as it never gets used. ([df766ac](https://github.com/apache/plc4x/commit/df766ac02551d364283e75d8fe432548fbbacf6b))
  - remove start_pos variable as it never gets used. ([7c2d17c](https://github.com/apache/plc4x/commit/7c2d17c5eab1b6087c6bd0e479ca1f4efc52e977))
  - Clean up debug info ([67272f7](https://github.com/apache/plc4x/commit/67272f7a11d15c3543360da8cef8f6eb5dd4a846))
  - Fix issues with serializing data for Modbus ([4350fb2](https://github.com/apache/plc4x/commit/4350fb28ee3767bd8a4256ea2a0379ebf968e2dc))
  - Work on the templates ([ad7fbdb](https://github.com/apache/plc4x/commit/ad7fbdb97d5fac6c2c3a3160a9da5ca3a4074c83))
  - Missed a change ([c447eb5](https://github.com/apache/plc4x/commit/c447eb580eb4860138c83364f34a957b22414913))
  - migrated the python helper to be based on the go helper ([3cbd8f5](https://github.com/apache/plc4x/commit/3cbd8f52e09763f277940777b512370c63b06ef5))
  - WriteBuffer Add local byte_order check ([3f7001b](https://github.com/apache/plc4x/commit/3f7001b18472645e45b1082eb10de7466b829e3c))
  - Finish WriteBufferByteBased ([1c3fa84](https://github.com/apache/plc4x/commit/1c3fa84cf356b9a2364220313e1d10798cb1d8b3))
  - Add initial signed int ([e14a3d7](https://github.com/apache/plc4x/commit/e14a3d7dc8bae8ead824f019d5e87767c4460adc))
  - Finish unsigned integer writes for ByteBuffer ([5ae793e](https://github.com/apache/plc4x/commit/5ae793ea3c5c9aa6cf067d48816c0d36f88734c0))
  - Add more tests for the write buffer ([9c8431c](https://github.com/apache/plc4x/commit/9c8431c2fd46f4f35271252d1d5eed139c7a0dba))
  - Update license header and implement the start of write_unsigned_short ([1c892bd](https://github.com/apache/plc4x/commit/1c892bd9b095ae17e95f83310317001595223a11))
  - WriteBuffer interface and WriteBufferByteBased implementation ([99ed37c](https://github.com/apache/plc4x/commit/99ed37c86330a2ef0421d9ed81e918cce87ceb42))
  - Formatting and more dummy methods ([2db5613](https://github.com/apache/plc4x/commit/2db561306ee49f2ee17363dd7db5306e0e0c4595))
  - Fix issue with naming of optional enum code-gen ([6b5cdf7](https://github.com/apache/plc4x/commit/6b5cdf700c5a2e3073ef87c17548120d35749870))
  - Implement the WriteBuffer interface ([9af6b0e](https://github.com/apache/plc4x/commit/9af6b0e69bc5536d626ee509fc8e99537e8e4228))
  - Build dependencies as well ([df10651](https://github.com/apache/plc4x/commit/df106518fb45dc83e192a04fc3dcfb8bf9aa76d5))
  - Update plc4py package version ([5c4db61](https://github.com/apache/plc4x/commit/5c4db612cb5c8961c61639703b7fafb00a6b7b7c))
  - now using async generators for connection ([87d795f](https://github.com/apache/plc4x/commit/87d795fcb23ff98134eb428901bd0fe05f52863d)) ([#362](https://github.com/apache/plc4x/pull/362))

- plc4x/bacnet:
  - update vendors ([2456668](https://github.com/apache/plc4x/commit/2456668ae06fb3b8333286ffa9ba4ac08f9d6f4c))

- plc4go/test:
  - add write interceptor to test transport instance ([8a7edba](https://github.com/apache/plc4x/commit/8a7edbadbec284c4b30e5d4bc3efb38e745ce70c))

- plc4go/simulated:
  - add preregistered (shared) instances ([7aec29b](https://github.com/apache/plc4x/commit/7aec29bfe7675625cb3af88d1de0dbf2cf070f13))
  - new subscriber stub ([63c1641](https://github.com/apache/plc4x/commit/63c1641c9d27d197771ea9d1d97b47a1f23595eb))

- plc4go/codgen:
  - added fmt.Stringer to interface ([af77600](https://github.com/apache/plc4x/commit/af77600bc606193f0d2dfcec5daae8b108c76e4f))

- plc4c:
  - Added string-length support in data-io ([8385cb6](https://github.com/apache/plc4x/commit/8385cb694d6f073ad2619ef405ee7f95b8913fa3))

- code-gen:
  - typeSwitches in typeSwitches now allowed. ([5143917](https://github.com/apache/plc4x/commit/514391705cb4fa697df550dcd9fccd0447c97205))

- plc4j/driver/open-protocol:
  - Implemented a rough skeleton for an open-protocol driver (Not functional yet ... It just compiles the mspec generated types) ([e3cbcdf](https://github.com/apache/plc4x/commit/e3cbcdf4aeea7991e934cdcac664d85bacbd391c))

- protocol/open-protocol:
  - Started implementing the Open-Protocol mspec ([d8300da](https://github.com/apache/plc4x/commit/d8300da06ae9dd6a24aa8fbd1f3044381a5dd285))

- example/ads:
  - Continued working on the ADS EtherCAT Browser ([11d6cde](https://github.com/apache/plc4x/commit/11d6cdee490564f045b9f7b07d91e6d85bc935c5))

- plc4go/modbus:
  - Implementing the correct reading of BOOL type ([17d7f76](https://github.com/apache/plc4x/commit/17d7f765c670f86c3fd110f010a3faafe8ee1c5a))

- build:
  - Migrate CodeQL to build tasks. This should shorten the time taken as we aren't building twice ([9bf9451](https://github.com/apache/plc4x/commit/9bf94512121080e2551e9728db8fe7eb4f393ac2))
  - Update build files to only build a language if something changes ([2a508d6](https://github.com/apache/plc4x/commit/2a508d60ebd4dc356c8f6852e00e7ce1481c3878))

- plc4x:
  - Update Github Action badges ([f6a7074](https://github.com/apache/plc4x/commit/f6a7074bde026898ac8fd4c1bf374f9f4aa4092c))

- plc4j/connection-cache:
  - Added support for subscriptions to the plc connection cache. ([1f81060](https://github.com/apache/plc4x/commit/1f81060feccbe16df1333c933078322e8e1152cd))

- plc4j/s7:
  - add S7-200 (SMART) type support ([3f46e01](https://github.com/apache/plc4x/commit/3f46e01c49abce7d21c295ec4c14e3d1f10c2a43))

- plc4go/s7:
  - add S7-200 (SMART) type support ([9e88616](https://github.com/apache/plc4x/commit/9e886163bb15edefb32dfe8c20c3bedcdf08e3ac))

- plc4py/github:
  - Add Java ([0a7f3a1](https://github.com/apache/plc4x/commit/0a7f3a17998641d3630c6c0d50ec1ae762336cb8))
  - Add libpcap ([b14a87b](https://github.com/apache/plc4x/commit/b14a87b44f01f181bb8016c010aac9527e18a15b))
  - Oops ([6182503](https://github.com/apache/plc4x/commit/618250339ff68ba1c02052327e3f77545493239c))
  - Install flake8 ([aec6b0e](https://github.com/apache/plc4x/commit/aec6b0ef35f7fefb50750e6f37c2e1e1e8cb998c))
  - Add GitHub action for python project. Also added dependencies for the python maven project to include the protocols. ([50c3fb6](https://github.com/apache/plc4x/commit/50c3fb6d75e03b82606b0cef3b7b73e975cbab0e))

- plc4j/example:
  - Fixed class path in read example ([ce8fe7e](https://github.com/apache/plc4x/commit/ce8fe7e913bf09c38f698e109e41bb5e8391c249))

- ads:
  - ADS Auto-Discovery ([71b7977](https://github.com/apache/plc4x/commit/71b7977b2f0f7902281812f72fcae9e90ac45041))

- protocol/bacnet:
  - add missing NLM types ([955efbe](https://github.com/apache/plc4x/commit/955efbea0cf0653b88bbfd8d4765c1c144866eaa))
  - introduced new c-bus constant ([a63ebe1](https://github.com/apache/plc4x/commit/a63ebe1ba2536e85c26cf12cc5e0d1f59e2e3755))
  - allow big sized segmented messages ([8cfaa94](https://github.com/apache/plc4x/commit/8cfaa94b4693bcf59c10c54fbc4c4bb41b12fdbf))

- protocol/knx:
  - update manufacturer ([c325f34](https://github.com/apache/plc4x/commit/c325f341fa751ce02d323a457d11ab6311eed06c))
  - update manufacturer ([60ee4e3](https://github.com/apache/plc4x/commit/60ee4e3fedb2ac901ea95abb75376c2c1d8614e6))

- plc-simulator:
  - add I-AM response to whois ([6cd52d2](https://github.com/apache/plc4x/commit/6cd52d2087d02ac312d0cf7d7fc52e5241b83021))

- plc-simulator/bacnet:
  - bacnet simulator is now able to return a valid hard coded response ([c1f1700](https://github.com/apache/plc4x/commit/c1f1700932da9b2f0db523a8a2f10b1e85bde6e8))
  - initial bacnet plc-simulator ([37d8880](https://github.com/apache/plc4x/commit/37d8880cbe6f2418f669456060d7492d533ead1e))

- plc4xanalyzer:
  - handle panics in actions ([d136f7a](https://github.com/apache/plc4x/commit/d136f7a223e9f14e87bdd0278b47ba6c517c8b8d))
  - exposed cli options to the ui ([00d418f](https://github.com/apache/plc4x/commit/00d418f0e6af3146fb001c611ad1b105905b7930))
  - added option to abort a long running job ([f3fd13b](https://github.com/apache/plc4x/commit/f3fd13bdef2d19ec3acbcfbb7fa10827368f3f95))

- plc4xbrowser:
  - set default output lines for outputs to now crash application on long sessions ([23fcdce](https://github.com/apache/plc4x/commit/23fcdce5f739aab1ba3087af56dde5ca163036d5))

- plc4xpcapanalyzer:
  - set default output lines for outputs to now crash application on big pcaps ([decb3b9](https://github.com/apache/plc4x/commit/decb3b9ca523842ab109d41a3bbcfd3c7825ce37))

### Bug Fixes

- general:
  - Added missing file header. ([2a8311b](https://github.com/apache/plc4x/commit/2a8311b4d00f281d0a640e7127ad854d038d6839))
  - Hopefully fixed the thread leakage we were seeing in several issues. ([709fa7f](https://github.com/apache/plc4x/commit/709fa7f338c6288ba5e87b845123bf673ab66e53))
  - opm collection field set value error ([008dd3c](https://github.com/apache/plc4x/commit/008dd3c99c2f95515549632ab314e390ffe29fd7))
  - Changed the encoding of strings from utf8 to WINDOW-S1252. Also fixed a bug with reading strings after the ADS refactoring. ([46648c0](https://github.com/apache/plc4x/commit/46648c03f0eea1b1594c3d6faa0b15de74987001))
  - Changed the encoding of strings from utf8 to WINDOW-S1252. Also fixed a bug with reading strings after the ADS refactoring. ([a5bdd80](https://github.com/apache/plc4x/commit/a5bdd809984340e8f75ee573dc688ccde284abe4))
  - Fixed the logging dependencies ([45adf96](https://github.com/apache/plc4x/commit/45adf96fe0fbefbc06c3a7048efd41d395ddb5ee))
  - Fixed a package-visibility issue ([072d179](https://github.com/apache/plc4x/commit/072d17981d60dc0a4b7c80dd5bd13746246cfa99))
  - Fixed numerous issues when writing values to S7 ([cad4fd5](https://github.com/apache/plc4x/commit/cad4fd5f82443804167f9e52a509db747b760d56))
  - Removed all Edgent references and made the IoTDB example work again. ([7eb9b2d](https://github.com/apache/plc4x/commit/7eb9b2d78aab496c28d46e842d1de94dd83c0a96))
  - Made sure the ADS Driver sensibly uses the Futures ([7045669](https://github.com/apache/plc4x/commit/70456698be59b0684bc30622818f981369fe9c53))
  - Made the ParserSerializerTestsuiteGenerator a bit more resilient against empty tcp packets ([63294e0](https://github.com/apache/plc4x/commit/63294e03bf4d1c6412baa5e14f18cb67791dd495))
  - fixed a bug in the handling of little-endian unsigned integers ([5a8bf78](https://github.com/apache/plc4x/commit/5a8bf784643c95350c4fa8c9f5e607f65368f1ba))
  - Tried updating the jenkins sonarcloud config for plc4c ([b75a42a](https://github.com/apache/plc4x/commit/b75a42a4b89a44998814ce4c471e53ef79aaa5f1))
  - Disabled the sonar-hack in the plc4c build. ([6cec1a7](https://github.com/apache/plc4x/commit/6cec1a7d5217b214a7556f78336b9c5fc0078991))
  - Tried switching the build to Java 17 ([853fba8](https://github.com/apache/plc4x/commit/853fba8f79d8f2210b0640a513cceff4b8d3e519))
  - Made the read/write tests in OPCUA also not run on mac ([f87fb3c](https://github.com/apache/plc4x/commit/f87fb3c27f94640b6bf5f872b2b9dbd0a23c80b8))
  - Addressed a lot of code-smells sonarcloud reported. ([3c2978d](https://github.com/apache/plc4x/commit/3c2978d1dba4213ff42cd608af1d626349fc440c))
  - Always install sources (#717) ([7d368be](https://github.com/apache/plc4x/commit/7d368be87a936c967c030f49b3238e05d323e77c)) ([#717](https://github.com/apache/plc4x/pull/717))

- plc4j/opcua:
  - Fix incorrectly handled GUID tags (#1099) ([bd064a5](https://github.com/apache/plc4x/commit/bd064a594440160d98da61e6eae741fa2777a263)) ([#1099](https://github.com/apache/plc4x/pull/1099))
  - not to get String.length but to calculate length in UTF-8 for Pascal string length (support node id contains unicode characters) (#668) ([0bc5655](https://github.com/apache/plc4x/commit/0bc5655226a8d36ba4acb212aef6fc21e82e2899)) ([#668](https://github.com/apache/plc4x/pull/668))

- opcua:
  - re-enable enums ([b2f3acb](https://github.com/apache/plc4x/commit/b2f3acb97b3b209b0434c2441bdb09a9eee641b6))
  - enforce non gzip as it seems to fail from time to time ([f8a2740](https://github.com/apache/plc4x/commit/f8a2740295677e774c4dd4fa9afe647cdc0b42f4))
  - disable enum code generation for now ([7361a10](https://github.com/apache/plc4x/commit/7361a10438242a0e3115dbd5e4a8c9b3485ed150))
  - split up enums ([03a9174](https://github.com/apache/plc4x/commit/03a917484f8f79411c2307de4ea25b6b8c48913d))
  - convert OpcuaNodeIdServices back to enum ([9413b96](https://github.com/apache/plc4x/commit/9413b9658b074b32411a27b8b502ba45d3d58cce))
  - small fixes for enum generation ([2663b5d](https://github.com/apache/plc4x/commit/2663b5d08372039b33a58c15cebd5bc111f7732e))
  - use constants for OpcuaNodeIdServices as enum explodes in java for now ([673de2c](https://github.com/apache/plc4x/commit/673de2c62b5af77d13b72ffe9e4f3dc2ec7a34bd))
  - exclude PortableNodeId ([1b9239b](https://github.com/apache/plc4x/commit/1b9239b738973797e7b885914b6cdbd02cb581be))
  - exclude PubSubConfigurationValueDataType ([0dce106](https://github.com/apache/plc4x/commit/0dce10678e03c61b13085328c348ff6cc12ad5ea))
  - fix endpoint element mixup ([c1232bc](https://github.com/apache/plc4x/commit/c1232bcf4fe755f6beb95f8db1962f578789d771))
  - fix address pattern ([8464f58](https://github.com/apache/plc4x/commit/8464f58b25e8dc6bf28c8ef207adf46c5c452a11))
  - fix issues with PascalString ([4610930](https://github.com/apache/plc4x/commit/461093014c5c1c29dee4154d0af2283ec87d8a29))
  - fix issues with empty port ([40644f7](https://github.com/apache/plc4x/commit/40644f7e626fc90485d5efddf5cdea2021a162dd))

- plc4j/spi:
  - fireDiscoverEvent should default to false ([763b41c](https://github.com/apache/plc4x/commit/763b41c70559ffcb2ce2a40038086b62df7c203a))
  - fix error message ([4f4712f](https://github.com/apache/plc4x/commit/4f4712feba947788173c39e6424492176e4b93e4))
  - PLC4X-344 Handle concurrent access to transaction worklog. ([cd5df7e](https://github.com/apache/plc4x/commit/cd5df7e5b43c6725d4e427181c85bb4b4d68699a))

- plc4go/cbus:
  - rework the way conditional logging is made ([2c7e786](https://github.com/apache/plc4x/commit/2c7e786c8e11d3c5e44d81fc2b41965274b60e18))
  - ensure TransportInstances are properly synced ([73317f8](https://github.com/apache/plc4x/commit/73317f81ae3a410753775d713577378868b03e8c))
  - fix error when reader doesn't get a alpha capable response ([bf275e2](https://github.com/apache/plc4x/commit/bf275e23d1f043f70ff8328d320eb40d1ebd0180))
  - fix some concurrency issue when closing the codec ([1797abe](https://github.com/apache/plc4x/commit/1797abebc554422431e808f91a0a55efc0e8896f))
  - handle pre registered consumers ([c12edf3](https://github.com/apache/plc4x/commit/c12edf33d65c293dcfc8dfb5172905520cec6184))
  - remove some remaining global logs ([0025f9f](https://github.com/apache/plc4x/commit/0025f9fd730858585a6535bd35363e887e4abfc6))
  - decrease wait time if it is a confirmed request to server ([922c721](https://github.com/apache/plc4x/commit/922c7218ed49183aad83fbf2debdfa6090ae4121))
  - remove duplicated codition for media transport protocol ([5985eca](https://github.com/apache/plc4x/commit/5985eca3504c38ce2ca22c59a2a13a050d0c33b3))
  - concurrency might lead to a incomplete reporting ([fb7ff11](https://github.com/apache/plc4x/commit/fb7ff11c3fe2a771f86b8b58be49ad9f385e1473))
  - make read on unit installation non blocking ([2c488c7](https://github.com/apache/plc4x/commit/2c488c7b6c3364f806ed70979509cd812081d023))
  - fix broken interface selection ([9643831](https://github.com/apache/plc4x/commit/9643831b17e8bd08583c5afda5336a8ba62074bb))
  - fix wrong query string ([34cf92f](https://github.com/apache/plc4x/commit/34cf92fbb72881c5e024a0ff3ff739a11e8d58ee))
  - fix length calculation of cal commands ([2545dc7](https://github.com/apache/plc4x/commit/2545dc7453cd77e8773a22431557dd040601bd66))
  - fix broken wildcards ([eecde5a](https://github.com/apache/plc4x/commit/eecde5a63d11cf1ca37640e755d74cb5879e141f))
  - try read a bit more greedy on confirm ([57eacfb](https://github.com/apache/plc4x/commit/57eacfb43bbc20e0704d88a5ef972df028c46809))
  - fix bridge addressing ([c045c60](https://github.com/apache/plc4x/commit/c045c6036177501ccdd658d533d5b4b9f84d1a7f))
  - Added a missing file header ([4bc2c19](https://github.com/apache/plc4x/commit/4bc2c194115ce88b6bb1fe69a670b2942a9ef643))
  - discover should not block anymore ([c032651](https://github.com/apache/plc4x/commit/c0326517329ad43046f0ab2fd3166aa2bcb4b204))
  - disconnect message codec on connection error ([e4e92b0](https://github.com/apache/plc4x/commit/e4e92b07427b56ca8c9fc49d859d4eab02159a18))
  - use queues for discovery to not overwhelm small devices ([7701247](https://github.com/apache/plc4x/commit/7701247420650ea179cf305630efb6960f5018bc))
  - limit discoverer with semaphore ([cabc35d](https://github.com/apache/plc4x/commit/cabc35dd353132b8cecdf1438cd0c98f43d85957))

- plc4j:
  - add fireDiscoverEvent (replacing awaitDiscoverComplete) and moving awaitDiscoverComplete to it's true prupose ([8e20f3c](https://github.com/apache/plc4x/commit/8e20f3cf80c60ae6a791e3d94cca49cdacd8e0f3))
  - Remove startPos from template output when not being used ([431c442](https://github.com/apache/plc4x/commit/431c442708f81c41b463b1caef20741923f8f7f9))
  - fix logging issue ([71b6c20](https://github.com/apache/plc4x/commit/71b6c207584ca4a740c9314e21ebf7cfb7c27ece))
  - fix critical sonar issues ([a3ea768](https://github.com/apache/plc4x/commit/a3ea768d2f7b35211d553505c9509b2c28a000d2))
  - fix test for simulated driver. ([7a8d107](https://github.com/apache/plc4x/commit/7a8d1075d1b6a9fc26a97fab13521d4a9048a75e))
  - test out the new build files ([de34ec7](https://github.com/apache/plc4x/commit/de34ec73a1a98ffe8882458756e99f1b936bb4c5))

- plc4go/plc4xbrowser:
  - fix terminal output ([f03e07a](https://github.com/apache/plc4x/commit/f03e07ac73cfa44b4f75d478dc5cf2e9bf8732c6))
  - fix small woopsie ([d5f9c68](https://github.com/apache/plc4x/commit/d5f9c68f531b8e5c7c78688a62e4aa73b86983fa))
  - handle panics on actions ([dece48b](https://github.com/apache/plc4x/commit/dece48b3b588d0a9d23e91b26a0ecbac311cc2f0))

- plc4go/opcua:
  - removed code whoopsie ([bc9a0b1](https://github.com/apache/plc4x/commit/bc9a0b1bbd16210e6e831119f62a660ecf8a1fea))
  - fixed reading ([a197af1](https://github.com/apache/plc4x/commit/a197af1493766ec5e0297aa3d96ed738c3d3b2a6))
  - fixed several small issues in SecureChannel implementation ([5a5ed86](https://github.com/apache/plc4x/commit/5a5ed862a73433a5a09baa0e8a5cbcee73960e6d))
  - fixed issue regarding sending of messages ([37a4aed](https://github.com/apache/plc4x/commit/37a4aedb2016fbbd449350ec094aebbd666f6cf5))
  - fixed issue regarding host revolving ([a8c6d26](https://github.com/apache/plc4x/commit/a8c6d26e9b3f2b036456b12b5d3b718bd237e00e))
  - ensure right atomic initialization ([da34d61](https://github.com/apache/plc4x/commit/da34d61bcf8237667f57ca3f07672788b50588a6))
  - fix issues with increments ([71521eb](https://github.com/apache/plc4x/commit/71521ebca0336527bc1e8632641799ccafb597fd))
  - fix driver name ([f47d4d9](https://github.com/apache/plc4x/commit/f47d4d9aede03720af38d1618b4c3fc328335110))
  - fix uri pattern regex ([925b7fa](https://github.com/apache/plc4x/commit/925b7fafe5b832d2eb385861d7b15ade344a05e0))

- codgen/plc4go:
  - added count array guard ([f1f08fb](https://github.com/apache/plc4x/commit/f1f08fb58923d55256d23db1be4e3fe194646105))

- plc4go:
  - fix broken mock ([efbab28](https://github.com/apache/plc4x/commit/efbab28b73d5f78c287bf38eab8dac7a2d23ceff))
  - fix deadlock issue related to subscriptions ([959a207](https://github.com/apache/plc4x/commit/959a207bc2bad21fa3d40a9418553da421d87553))
  - ensure we are not blocking ourself on test ([0bafd59](https://github.com/apache/plc4x/commit/0bafd599b825e7c2f4816b63a272ba71e568ca0f))
  - fix issue with generator ([aff7788](https://github.com/apache/plc4x/commit/aff77887f2b76979a7c47115dff53836bf99f2dd))
  - sync remaining mocks ([b9545f5](https://github.com/apache/plc4x/commit/b9545f518b70e1d2a5a2ebca54a5c4e4bd30aa2d))
  - sync generated code ([c9e146e](https://github.com/apache/plc4x/commit/c9e146e9697e774507ce4b8db53b6060cf5aa9b0))
  - fixed some issues regarding subscription tags ([a86b285](https://github.com/apache/plc4x/commit/a86b2858727dbe77319c99ada2aa86ac455575f6))
  - ensure options are passed downstream ([840ca2a](https://github.com/apache/plc4x/commit/840ca2a2d22cb33544050e73373f890d66c6c998))
  - fall back to global logger if no logger is supplied ([c970c27](https://github.com/apache/plc4x/commit/c970c2739234be7142d532f218d50ca9881b2966))
  - return correct WithOption for WithExecutorOptionTracerWorkers ([9d1877f](https://github.com/apache/plc4x/commit/9d1877f38195ef759baaf2551b10ca51756f626c))
  - add more logging to hunt down race conditions ([333f9b9](https://github.com/apache/plc4x/commit/333f9b90d3fb62bc5c270212f44681b772268dc7))
  - fix concurrency issues ([e6f661e](https://github.com/apache/plc4x/commit/e6f661eab583e06e014c25ea92851fea8e3c48dc))
  - fixed several race conditions ([3297f63](https://github.com/apache/plc4x/commit/3297f63ffbb3fab8d59d98c7067f5c5e172a2ef7))
  - Made the Tracer synchronized ([9cc2445](https://github.com/apache/plc4x/commit/9cc2445a621ad4731b74b6cd7a7a31acda0c4174))
  - add missing context passing ([22dbc15](https://github.com/apache/plc4x/commit/22dbc15944730aef427cc716e0778db4b782c665))
  - always supply stack when recovering ([78ae3b6](https://github.com/apache/plc4x/commit/78ae3b6430d329b2db02b5226f251e6e68db7049))
  - fixed some quality issues ([dd568f9](https://github.com/apache/plc4x/commit/dd568f9ef4720f73cfeafa295aa5ab9e1dd8d56c))
  - transaction should now be properly handled ([25480b1](https://github.com/apache/plc4x/commit/25480b1d22a08f863ba15383d364e2b29605e35c))
  - pass loggers where missing ([2ff14f3](https://github.com/apache/plc4x/commit/2ff14f33c0a5fd2351d23128253f98804b38e3f1))
  - don't panic - catch panics and log them ([33f1d9a](https://github.com/apache/plc4x/commit/33f1d9a283dfcb91365bb03aed45238340c377ae))
  - Fixed one place where Sebastian's HexDiff wasn't used yet. ([386cade](https://github.com/apache/plc4x/commit/386cade4442a96b92cf8747ab379e447d872ee2d))
  - upgrade x/net ([6351dc6](https://github.com/apache/plc4x/commit/6351dc6e0f5e74323080fa7ba863a67551065377))
  - Fixed posrting the simple PlcValue types to the new context-aware serializers ([477e1af](https://github.com/apache/plc4x/commit/477e1afbbecb1a87b24561625d0506740c8149bd))
  - Stated fixing the driver testsuite to run properly ([3062da0](https://github.com/apache/plc4x/commit/3062da027ecaa8730ccc2ef3fbeb6062cb1dd9b3))
  - updated KNX Manufacturer generated file. Testing Go build. ([d266634](https://github.com/apache/plc4x/commit/d26663459577b086c4bc7f15044355eeae7c6792))
  - fix(plc4go) Fixed the reading of 32 bit floating point in Little Endian format. ([461946a](https://github.com/apache/plc4x/commit/461946acc8a4696fc50648e76cfef6a7fdd0d7a1))
  - pass in the correct context (#537) ([b6bf718](https://github.com/apache/plc4x/commit/b6bf718eb958fb5c3aa3cc59c8eb210fe4f1dc42)) ([#537](https://github.com/apache/plc4x/pull/537))

- codegen/plc4go:
  - fix issue with strings and virtual fields ([5704752](https://github.com/apache/plc4x/commit/57047521ecf614b8f0d0f894170acc1db802c9a1))
  - fix abstract fields ([1dccf57](https://github.com/apache/plc4x/commit/1dccf57f280e0ad7b45ce895d65fb7dfa9f02c10))
  - fix support vor vstring with tenary terms ([5c576e6](https://github.com/apache/plc4x/commit/5c576e6d371a003c442a95a6dcc764f4ef5e2ac7))
  - fix some string handling issues ([d4864ed](https://github.com/apache/plc4x/commit/d4864ed46a3ff2bb4a1a6ea186fd4c23e42513f7))
  - fix encoding retrieval with null fields ([20e5c91](https://github.com/apache/plc4x/commit/20e5c9124e1df494cae1c98e6f8dcbd50382679d))

- code-gen/java:
  - Updated the codegen for java to generally use smaller types for unsigned integers (a byte is able to keep an unsigned integer value to 7 bits and not up to 4 bits) ([3917dab](https://github.com/apache/plc4x/commit/3917dabd2351d59b1b3f7b95efa9fd144cdb7e09))
  - Updated the codegen for java to generally use smaller types for unsigned integers (a byte is able to keep an unsigned integer value to 7 bits and not up to 4 bits) ([8772f90](https://github.com/apache/plc4x/commit/8772f909d862003a558b378e5a8e1c882677f078))

- plc4go/modbus:
  - Delete elements in the loop, and the index is dec… (#1028) ([a359a2f](https://github.com/apache/plc4x/commit/a359a2f988d3bef7af5ffa3babdd3f80e3e9b355)) ([#1028](https://github.com/apache/plc4x/pull/1028))

- plc4go/spi:
  - ignore empty declaration of PLC4X_TEST_RECEIVE_TIMEOUT_MS ([2e65a27](https://github.com/apache/plc4x/commit/2e65a2759055899e289466d3cf0dbb49b200e188))
  - increase sleep times of DriverTestRunner ([77453bf](https://github.com/apache/plc4x/commit/77453bf1c465ae4c069e3152b9f605af5be07a0f))
  - avoid unnecessary read on DefaultCodec shutdown ([fe99681](https://github.com/apache/plc4x/commit/fe996814aebc9d585e35e23c93f9286ce1d9d215))
  - use LookupEnv instead of GetEnv in test utils ([31a8aa5](https://github.com/apache/plc4x/commit/31a8aa51c3691dfe4eed69fceab6ed12fbe11356))
  - sync tcp.TransportInstance state change ([14c59f0](https://github.com/apache/plc4x/commit/14c59f0f4ab5a636bd200b7d18fc655ba4f8132b))
  - fix concurrency issue when a executor is being started and stopped pretty fast ([5894b08](https://github.com/apache/plc4x/commit/5894b08ef7b7040de580e0011a5ef7e891597487))
  - fix race issues in request transaction ([0458529](https://github.com/apache/plc4x/commit/0458529ee1acfda55ae7ffaf75d6c95045a06fea))
  - fix worker logging on wrong logger ([dcf630a](https://github.com/apache/plc4x/commit/dcf630aa840b03e44549d88395535b5fbe613d6b))
  - fix race issues in worker pool ([5e51e66](https://github.com/apache/plc4x/commit/5e51e6606dc9fe58d8324a5e2e102e0c4537aeb4))
  - multierror only returns a error if it has an error ([ddda58b](https://github.com/apache/plc4x/commit/ddda58bb3e4b9d0d8a188a1b534d64a17c4a37d0))
  - DefaultCodec.go rendering fixed ([747e2ee](https://github.com/apache/plc4x/commit/747e2ee9fdfddcc95943cee3052d8a7889efc657))
  - potential fix with request transaction manager producing race conditions ([32c5531](https://github.com/apache/plc4x/commit/32c5531d2f65bfe86759f3a8a7fb6a62cc4bc26e))
  - avoid shutting down the shared executor ([0784d37](https://github.com/apache/plc4x/commit/0784d3786f066b61c4a637f2d0063faf017971b6))
  - re-order disconnection on Defaukt codec ([3456bca](https://github.com/apache/plc4x/commit/3456bcaa543f941c33f27a6e0b01d67a6c3be3e2))
  - test transport instance panics if worked with on disconnected state ([3c7acbd](https://github.com/apache/plc4x/commit/3c7acbdfa5f310cb1119f7fc029e5aa2ed5744b6))
  - fix timeout output using the wrong duration ([9db3034](https://github.com/apache/plc4x/commit/9db3034ea5d4d1e27819accac73aab4b48afcedf))
  - fix data race in executor ([860a15b](https://github.com/apache/plc4x/commit/860a15b0e2c7f9d676aa589363e55e2024c40dd9))
  - make shutdown of WorkerPool more reliable ([7d745da](https://github.com/apache/plc4x/commit/7d745dae3c663ab5c7b4a26a4a444c78fc8433ed))
  - harden request transaction manager implementation ([b9c89eb](https://github.com/apache/plc4x/commit/b9c89ebea3d57a1153919e05cb947d61c12f282e))
  - properly shutdown worker spawner and killer on shutdown ([d9584bc](https://github.com/apache/plc4x/commit/d9584bcde767b716a46f1554b2cf90f79335bab9))
  - fix WSTRING production ([c9db23c](https://github.com/apache/plc4x/commit/c9db23c018a5cfb48a019be0b52b799bfba9510e))
  - fix small refactoring woopsie ([da294ba](https://github.com/apache/plc4x/commit/da294bae2d7905bcf4723ac155df58ccf2968756))
  - gracefully handle tag names not found on WriteResponse ([0a14655](https://github.com/apache/plc4x/commit/0a14655948ed44bb2a62b4ce58a18bad652fb5f7))
  - gracefully handle tag names not found on DefaultTag ([fdce5b9](https://github.com/apache/plc4x/commit/fdce5b9aa7fafe83f85c9db794fe4b826f0d35e0))
  - gracefully handle tag names not found on SubscriptionEvent ([1f16e0f](https://github.com/apache/plc4x/commit/1f16e0f7beb8605e71284765bc22891bd915cd8d))
  - gracefully handle tag names not found on ReadResponse. ([7de8439](https://github.com/apache/plc4x/commit/7de8439f1d098a0e84024d560752698f83bf0416))
  - gracefully handle tag names not found. ([c240ade](https://github.com/apache/plc4x/commit/c240ade351764503a3373b073dd8893857642f36))
  - don't panic if sub handle is not found ([4c76916](https://github.com/apache/plc4x/commit/4c7691680d21c4caaab68a31ae284c02dc85fad2))
  - fix timing issue when closing cached connection ([c857f83](https://github.com/apache/plc4x/commit/c857f837fa4b00b6b2f1cbcd26f5ddf233521f85))
  - add test for WriteBufferPlcValueBased ([2c092fe](https://github.com/apache/plc4x/commit/2c092fe14be00fbc931b0f02ca38e67744bdd0f9))
  - don't panic when nil runnable is submitted to WorkerPool ([2b029b4](https://github.com/apache/plc4x/commit/2b029b45860f4c5d13a8865cd1e9ca0121e4f150))
  - fix issue with pcap close ([6bd2c57](https://github.com/apache/plc4x/commit/6bd2c57f353c7ea6adaa8ba1a9ff78f74d92bfa5))
  - fix issue with pcap close ([f4db5fb](https://github.com/apache/plc4x/commit/f4db5fbf4de5dcd47e5e1d63956dfb70dd5d491b))
  - print stack when panic is caught ([3011a3e](https://github.com/apache/plc4x/commit/3011a3e515ce63273a77da6fcb7a0b7ae8479151))
  - avoid test transport getting stuck on a endless loop when filling ([6b8da79](https://github.com/apache/plc4x/commit/6b8da79620d81a54dff7531c5ceddb14e6e219b3))
  - fix issue with UTF16 encodings cpu drain ([d28e704](https://github.com/apache/plc4x/commit/d28e704f766d10eed4b2086f22dfd0589b4e9918))
  - timeout expectation fails when expired and context errored ([ce079a0](https://github.com/apache/plc4x/commit/ce079a04cb6e04f7acbeff83998b7c2ea1e26e2f))
  - timeout expectation should not fail anymore ([bfea265](https://github.com/apache/plc4x/commit/bfea2658d366baf4b34ddfc2a12e8b90bd488d0e))
  - fix default driver delegation call ([07ae197](https://github.com/apache/plc4x/commit/07ae197fa60bd4aec9d90b113f79cf21491d20c4))
  - fix concurrency issue in DefaultCodec ([0cbf147](https://github.com/apache/plc4x/commit/0cbf14751240583d62f9a11100cc295c0613f222))
  - fix speedfactor output of pcap transport ([7efe9c9](https://github.com/apache/plc4x/commit/7efe9c9226280cbecda724279c45059d5c4e24f1))
  - fix RequestTransationManager using the wrong executor ([23f07b6](https://github.com/apache/plc4x/commit/23f07b6e61844259124e26e1315691fdad009637))
  - fixed npe while rendering fields ([7651ed1](https://github.com/apache/plc4x/commit/7651ed17485b0053b2d10cb5715d1b3cc33a9590))

- site:
  - Remove reference to the Confluent Hub in Readme.md ([12fc12c](https://github.com/apache/plc4x/commit/12fc12cb9f8bbfe2b61695c1f2caf91c372cc77d))
  - Remove reference to the Confluent Hub ([1e1e676](https://github.com/apache/plc4x/commit/1e1e676bf0b8bd7da9bc9fcbdfb09bc3526f5ccf))

- plc4go/tools:
  - licenser should now output the right file name ([16e91ad](https://github.com/apache/plc4x/commit/16e91ad87fe7c05eecd98039398d68a5a637ada4))
  - fix atomic.Pointer support ([7089ac3](https://github.com/apache/plc4x/commit/7089ac3b584ba7fd776905c8ed29d9eaa4a59eb4))

- plc4go/epi:
  - matches should work on exact basis ([442960c](https://github.com/apache/plc4x/commit/442960c2fe3e09a259c42c42371b7bd47b12cfee))

- plc4go/eip:
  - fix possible double match in switch ([096cb0d](https://github.com/apache/plc4x/commit/096cb0ddf85ab07bacab34149c610f654d1adf15))
  - fix possible contention by using 1 buffered chan ([85dcb43](https://github.com/apache/plc4x/commit/85dcb431983a6eeca7b3f05be4ea0a06c10f2ccb))
  - Tried to make the code more robust to no longer cause errors on Jenkins ([f283dc2](https://github.com/apache/plc4x/commit/f283dc2829041ab16b6ae0c10db1c7aa983e638f))

- eip:
  - Made the tagHandler actually parse the tag address instead of just passing it along as string, hereby setting the datatype, which was required for being able to write. ([f0fe0fa](https://github.com/apache/plc4x/commit/f0fe0faee69d7bef8dee266536e0b245c82f5325))
  - Adjusted the CIPAttributes type to allow the packet structure used by my AB CompactLogix controller. ([07fa836](https://github.com/apache/plc4x/commit/07fa8363791565a6b1328ae5602630a17cc4fab2))
  - Adjusted the CIPAttributes type to allow the packet structure used by my AB CompactLogix controller. ([693098b](https://github.com/apache/plc4x/commit/693098bbc00d465e72f0f172cd802ec12160ad4a))

- plc4j/profinet:
  - Sonar fix issue with changed name of parameter ([e57df53](https://github.com/apache/plc4x/commit/e57df53af3a0d97bd16af71c3b492ac791a02cda))
  - Sonar fix regex matching zeo chars ([30995f2](https://github.com/apache/plc4x/commit/30995f2c95ce619203815ea85ba778ce8c26856b))
  - Update for OctetString and F_MESSAGETRAILER4BYTE datatypes ([29e6fc1](https://github.com/apache/plc4x/commit/29e6fc1384949c4f2ade4274053471582d76cbb7))
  - Format of Allowed in slots string can be a single digit. ([09c5dc7](https://github.com/apache/plc4x/commit/09c5dc7792e0dab81c0fe154ed61ffd113dfc8bb))
  - Allow space chars in device, device access and submodule names ([261cb2e](https://github.com/apache/plc4x/commit/261cb2e531f06025b08e07b79268027f8661d134))
  - Updated documentation around the ip address usage ([24f7137](https://github.com/apache/plc4x/commit/24f7137b92fdbed5ae5c6993f1da7df4b0ca5d6d))
  - Implmented pop for the dcerpc response queue. Syncronized cyclic counter for pnio packets ([40e6a66](https://github.com/apache/plc4x/commit/40e6a66cd868e78d4890095dc2a0dd4874c3f18d))
  - split profinet mspec and checked endianess of PNIO packets ([635602d](https://github.com/apache/plc4x/commit/635602db1c74f2d5056d257a7c94274ea8c387c8))
  - Fixed issue with endianess - still need to specify the rest ([0900b34](https://github.com/apache/plc4x/commit/0900b34e153297eb9b4262f799d49049b1ff7220))
  - Error in IEE8023 Tlv for LLDP broadcast. ([33dc434](https://github.com/apache/plc4x/commit/33dc434a4a0316c7af5b8393b38682e11669b0f0))
  - Fixed issue with Write Parameter Size ([6203dc6](https://github.com/apache/plc4x/commit/6203dc6b07861169d2971855cfd92bde92070a64))
  - Fixed issue with subslot and ident numbers. ([ef74a15](https://github.com/apache/plc4x/commit/ef74a156bc34be4f216bfd71a79821d63eb2f505))
  - Implemented additional LLDP TLV's for Profibus and IEEE8023 ([97d9468](https://github.com/apache/plc4x/commit/97d9468a089aa118a0d054af0dccfc0e983a1f04))
  - Fixed some issues around parsing gsd files and vendor id formats ([0772198](https://github.com/apache/plc4x/commit/077219813d94068ea340cadf05f226546f56388d))
  - Removed dangling bracket from mspec ([0e6e0ac](https://github.com/apache/plc4x/commit/0e6e0ac1bd1a3605b2d9c37d33f16c9436b7df8d))
  - Updated website page ([e5b6cdd](https://github.com/apache/plc4x/commit/e5b6cdd66376d481cbbba84d45714ba39819c45d))
  - Updated magic numbers in discovery packets ([cb8402f](https://github.com/apache/plc4x/commit/cb8402fbf8cbf7bfefcde9fec0802538c267d7c5))

- plc4go/knx:
  - avoid panics ([9328974](https://github.com/apache/plc4x/commit/93289741b42d325d0d4f2f7a07ba163947911e8e))
  - use queues for discovery to not overwhelm small devices ([55066a7](https://github.com/apache/plc4x/commit/55066a78fe4f2dad5fce76fd54da0afa9457144c))

- plc4x/opcua:
  - Revert the download-maven-plugin to version 1.6.8 ([365edc9](https://github.com/apache/plc4x/commit/365edc91fce1858e6ccd4ddecd7a6d9dcc83a516))
  - remove annoying debug messages in xslt files ([99e234f](https://github.com/apache/plc4x/commit/99e234f99ff0fa299691b0a7861b1044728703a6))

- plc4j/s7:
  - fix build ([7647763](https://github.com/apache/plc4x/commit/7647763f50b571355f7714d6de5185755a39d513))
  - fix concurrency issue on tpduid rollover ([1f99989](https://github.com/apache/plc4x/commit/1f99989f6c3fcc8ea11962f32c54f94c7cb04c60))

- plc4go/gen:
  - star delegates should now work. ([ca5452c](https://github.com/apache/plc4x/commit/ca5452c4026ca1535a875c5121a50b9a96103ef7))

- plc4c:
  - repair plc4c ([b2be09f](https://github.com/apache/plc4x/commit/b2be09ff6b217d01aeb6078b9ec22f963ab1e56b))
  - test the c build files ([9279242](https://github.com/apache/plc4x/commit/92792422387c46b956ef65bb0bbb072ac1ee5227))

- plc4go/tcp:
  - fix tcp String() when local address is nil ([92039d0](https://github.com/apache/plc4x/commit/92039d09ea7fb85fb18b9e7f2dbc2b8e1eb4d5de))

- plc4j/connection-cache:
  - Fixed an issue with the timeout handler not being cancelled ([bf4f1e2](https://github.com/apache/plc4x/commit/bf4f1e2ac4e0fea617ef7f4792313133290df52c))
  - Fixed a problem, that after encountering an error with the PLC, connections weren't invalidated. ([9b06c2d](https://github.com/apache/plc4x/commit/9b06c2de0c77a7c1bbcb730bb5285c4435002c93))
  - fix issue with timing of thread during double connection test. (#796) ([b58ae5d](https://github.com/apache/plc4x/commit/b58ae5dac53ac99cf3827beab2f81c859e83df36)) ([#796](https://github.com/apache/plc4x/pull/796))
  - Increased the setup timeout to possibly get the test running on jenkins. ([ee18ebd](https://github.com/apache/plc4x/commit/ee18ebd734e11bb1c06af5f1ac865c23e7878640))

- plc4j/logix:
  - Include documentation for Logix Driver ([ee2a5f3](https://github.com/apache/plc4x/commit/ee2a5f3fc0e264c7b595c2b823ea05ce4cf1610a))

- plc4go/bacnet:
  - disable excessive logging on Task ([75ea6f8](https://github.com/apache/plc4x/commit/75ea6f8b6e840cf3a55234924a90b5dc1b978def))
  - fixes in BIP simple ([fa5e0ec](https://github.com/apache/plc4x/commit/fa5e0ecf2eac91a41a88f1dde7b9904986dd75b0))
  - fixes in application layer ([5eba21b](https://github.com/apache/plc4x/commit/5eba21b22ad73ed27de7a083a62be882d5bab669))
  - fixed a bunch of broken code parts ([609f6af](https://github.com/apache/plc4x/commit/609f6afbf21e45e713392ba0c475c98de3a7228d))
  - smaller fixes relating udp comm ([75872b4](https://github.com/apache/plc4x/commit/75872b46ef3d274d94d53265833e2e70d91ea76b))
  - several fixes ([b28801f](https://github.com/apache/plc4x/commit/b28801fc4a15a2c365103a106508a581fee45fc6))
  - fixed static helper creating wrong segmentation request ([ace9cc3](https://github.com/apache/plc4x/commit/ace9cc35e79bbc0bf27c416d145852212101a91f))

- plc4j/test-generator:
  - Fixed an issue causing core-dumps on Mac ([55d85be](https://github.com/apache/plc4x/commit/55d85be9182afa367454b1f3f2e7025e52068a8b))

- bacnet:
  - replace the dummy tag number so deep equals doesn't crap out ([488b457](https://github.com/apache/plc4x/commit/488b457cdc5eb3aa812688d242dae1d43c411739))
  - fix network layer message related to routing ([3311c89](https://github.com/apache/plc4x/commit/3311c89c84419925b2f531cee1923f0a46fc8d86))
  - fix segment ack ([e0a7639](https://github.com/apache/plc4x/commit/e0a76392a2d884d4be6ea972388871717c80b95d))
  - use BACnetConfirmedServiceChoice where appropriate ([ba86633](https://github.com/apache/plc4x/commit/ba86633aab7bb7f56473fe583c1d4954766cc83d))

- code-gen:
  - Handle the case where while dispatching types the consumer modify the consumer list. ([e369b2e](https://github.com/apache/plc4x/commit/e369b2ea12e3b0a863cc3b0b28ede33d45dae18d))

- plc4x/c-bus:
  - remove dangling parenthesis in virtual field ([92e5ac9](https://github.com/apache/plc4x/commit/92e5ac9365bb90a67be4f57ba75bb0d9a23b2eeb))

- plc4j/eip:
  - partial fixed length calculation ([2b0ae68](https://github.com/apache/plc4x/commit/2b0ae68e054ad3d4e99cae02e3e638b440996db1))

- plc4j/nifi:
  - removed explicit version reference in pom.xml for jackson. ([91b68ad](https://github.com/apache/plc4x/commit/91b68ad4b9864885351c6b286074d069d87c1efc))

- plc4j/examples/ads:
  - Fixed some wrongly named constants ([812adb8](https://github.com/apache/plc4x/commit/812adb8b29b6022ee80aeab1164cad68a309dc9f))
  - Fixed some wrongly named constants ([795fc77](https://github.com/apache/plc4x/commit/795fc77c93b9845d27a6ffa9d5e3d76366de28ec))

- plc4j/examples:
  - Improve java example (#724) ([63b77e1](https://github.com/apache/plc4x/commit/63b77e10e07bcc7bb525f0b3c56ed153411b7bde)) ([#724](https://github.com/apache/plc4x/pull/724))

- s7comm:
  - Commented out some recent changes in order to get the build running again. ([46ed026](https://github.com/apache/plc4x/commit/46ed02625c06fda709879473f8ac51e2114a0101))

- plc4py:
  - Fix for multi enums for python 3.7 ([179e8ee](https://github.com/apache/plc4x/commit/179e8eef794ed3e9c1bcd9da85d214e8ce63295c))
  - Fix compatibility with 3.7 when using multi value enums ([b8c4bf3](https://github.com/apache/plc4x/commit/b8c4bf36e35cb2e43ff200d086c1e5e3dddf7dc9))
  - Trying to fix pip being uninstalled on windows workflow ([7f5ea6f](https://github.com/apache/plc4x/commit/7f5ea6fdb4f8af7fc29ee79af309b9dd1d15e72b))
  - Trying to fix pip being uninstalled on windows workflow ([4b16655](https://github.com/apache/plc4x/commit/4b16655c394bd29faae2d1c17614779d61f82a8b))
  - Add wait for process output to python checks ([9c3b983](https://github.com/apache/plc4x/commit/9c3b983e1b9f2986a37042ea4f84ba96a86db68e))
  - Add wait for process output to python checks ([f18bc8b](https://github.com/apache/plc4x/commit/f18bc8bff01ad6585a6f8360e1628e9bd68605f2))
  - Add direct output of python --version during prereq check, troubleshooting ([d867414](https://github.com/apache/plc4x/commit/d867414f8d9a6de7a0ea130e500b3cd381d47cd6))
  - Disable CodeQL temporarily due to integration not being accessable ([55594ef](https://github.com/apache/plc4x/commit/55594efc594d229397bbc3cf1f8e6e2f91e74cf8))
  - Get rid of warning for the workflow about having branches in both pull and push sections ([b8b9641](https://github.com/apache/plc4x/commit/b8b96419dff88dfa9fe2bd95b0fb34a8af955a01))
  - Temporarily disable the lint check in the GH workflöow ([02ea818](https://github.com/apache/plc4x/commit/02ea81830e5c874d0d6c1e164e649d42a82f0833))
  - Add python/python3 switch for linux/windows executable ([a609edf](https://github.com/apache/plc4x/commit/a609edfefc1fcf3a2e72a4f622bbc2c65a092d9c))
  - Replace builtin types (list and dict) with class from the typing package. ([dd9d36f](https://github.com/apache/plc4x/commit/dd9d36ff17f9deab14c146d531779ce39dbfd8ef))
  - Replace builtin types (list and dict) with class from the typing package. ([135452d](https://github.com/apache/plc4x/commit/135452d103495202f363ebbe40a415d2881002f8))
  - Used Union for type hints, remove manual tests from build ([1e4e1d3](https://github.com/apache/plc4x/commit/1e4e1d395396cedcabeae7a23c1b5933f1cccbd5))
  - Found some inconsistancies in the PlcConnection classes ([add49ca](https://github.com/apache/plc4x/commit/add49ca88e9616e483e2924a3ad24d260c14843c)) ([#362](https://github.com/apache/plc4x/pull/362))

- plc4j/enums:
  - String comparison when evaluating enums ([3b84d5d](https://github.com/apache/plc4x/commit/3b84d5da8617f1b6f8f23374c7416200f2a2db9b))

- prereq:
  - Remove trailing -ea string in Java version string for version 20 ([d465153](https://github.com/apache/plc4x/commit/d465153d5e4db99a2fda42a1ef222e066ae54b73))

- plc-simulator/bacnet:
  - simulator should now respond to the right ip on bacnet ([7ff73a2](https://github.com/apache/plc4x/commit/7ff73a2cf63cd1d90fb7d33ee32a329554aa2bb0))

- readme:
  - Add PLC4X website link to the PLC4X image in the readme ([6038199](https://github.com/apache/plc4x/commit/60381995b68fd7190b2b0c7a86ab92f4aaada94a))

- plc4j/modbus:
  - Cleanup of ModbusTag (#732) ([8c49f4c](https://github.com/apache/plc4x/commit/8c49f4c5ada45e17f95cdfe892bd5e3a11db762f)) ([#732](https://github.com/apache/plc4x/pull/732))

- build:
  - fixed python build badge ([cb8e7e0](https://github.com/apache/plc4x/commit/cb8e7e026cf1bb74fc858283aa7c1c806d1b8167))
  - Also run the build when a pull request is craeted ([2c7ae0c](https://github.com/apache/plc4x/commit/2c7ae0c704b4ed14ba726a12d666f4739f724542))
  - Increased the timeout for getting the python version as the 500ms was way to fast ([97ac2a1](https://github.com/apache/plc4x/commit/97ac2a13f8eb31e7b816686307a4d01d6cec2da9))
  - Added a self-activating profile which adds the generated-sources for VSCode. ([f17780b](https://github.com/apache/plc4x/commit/f17780b640f2f9be9099e729790961673735d0dd))
  - Added a self-activating profile which adds the generated-sources for VSCode. ([44e1696](https://github.com/apache/plc4x/commit/44e16967d240ce9793103f652eab522917c61e70))
  - Added a self-activating profile which adds the generated-sources for VSCode. ([500732b](https://github.com/apache/plc4x/commit/500732b326422b50ce2236693f98a59e9a5e89c1))

- plc4j/camel:
  - Remove the camel integration as it has moved to the Apache Camel project ([b784ded](https://github.com/apache/plc4x/commit/b784dedaa787e52e7907085bf096e8ea2e093f3b))

- plc4x:
  - Update build status from Jenkins build ([1d165e8](https://github.com/apache/plc4x/commit/1d165e8f9fe68491a5aa040630ee6b9eb6037c12))
  - Fixed a naming issue with one of the source-files ([0bb2b4b](https://github.com/apache/plc4x/commit/0bb2b4ba9c34d62464c84dfb3e9509b10a6a4181))

- protocols/modbus:
  - fix write requests for coils always set to false (#710) (#711) ([c9a1938](https://github.com/apache/plc4x/commit/c9a1938eb19b7103e89173597e02de9f622d42ad)) ([#711](https://github.com/apache/plc4x/pull/711))

- plc4x/ads:
  - fix(plc4x/ads) Support for basic browse requests. ([25eb3e3](https://github.com/apache/plc4x/commit/25eb3e397aaa84140a4b03ae38512f702f148495))

- plc4go/ads:
  - fix(plc4go/ads) Got the subscriptions working for ADS. ([2c86e18](https://github.com/apache/plc4x/commit/2c86e18e70f7c9577c3e965fedc8230434ab41b3))

- plc4go/c-bus:
  - fix unit info query ([68d7a45](https://github.com/apache/plc4x/commit/68d7a4556b68e823a602c2adb513f58cf094cd8a))

- plc4j/bacnet:
  - fix static helper creating wrong requests ([3577f8b](https://github.com/apache/plc4x/commit/3577f8bfd2036eeda7319d5198212ad78f54d094))

- plc-simulator:
  - fix cBus using s7 port ([ddf17ee](https://github.com/apache/plc4x/commit/ddf17ee0db01d4b9d518bedf88e2dae537aa6838))

- plc4xanalyzer:
  - protocol filter should not override global filter ([f8b1de0](https://github.com/apache/plc4x/commit/f8b1de00e44e99489d0e0c9215e5b490017259e4))

- plc4xpcapanalyzer:
  - guard against unknown protocols ([3c11cbd](https://github.com/apache/plc4x/commit/3c11cbd7360d8c95dc640fc5d6483a83dfa17190))
  - fix file navigation ([e54ee39](https://github.com/apache/plc4x/commit/e54ee39ad4f7a9de31e6e16a49a3d843332a4918))
  - fix handling of cap files ([ac27254](https://github.com/apache/plc4x/commit/ac272549c245f6071afddf58e587585073aee6a6))

### Documentation

- general:
  - Updated the RELEASE_NOTES ([9369e7d](https://github.com/apache/plc4x/commit/9369e7d0e9b65b1c6e745ad8e3a9189d2b1fc6ac))
  - Updated the RELEASE_NOTES ([957a786](https://github.com/apache/plc4x/commit/957a786c13556dec4f54bafcaa4fdd25d7f40387))
  - Updated the general concepts and plc4j getting started guides. ([66649d4](https://github.com/apache/plc4x/commit/66649d4e99e99c99db164d73178720ec718846a2))
  - Added some links to important information on S7ComPlus [skip ci] ([06a68f1](https://github.com/apache/plc4x/commit/06a68f13d06a6a23f36d4c9fba66d1b4ecdafabf))
  - Added a comment to releases on Mac systems ([4b5f808](https://github.com/apache/plc4x/commit/4b5f808ee63cbe25d9d11c37278c283d0a190ce3))
  - Updated the docs for the knx protocol ([ff37b11](https://github.com/apache/plc4x/commit/ff37b1152773143bc9b6d1477d97f007ade525e4))
  - Added a paragraph on setting up IntelliJ to not get confused by the failsafeArgLine in the pom. ([3b695e5](https://github.com/apache/plc4x/commit/3b695e5e3399432c759e296c6a5deda2a32a8369))

- plc4go/bacnet:
  - Added some comments ([255ddc3](https://github.com/apache/plc4x/commit/255ddc3394cde0e889b425e78387583ba90d929b))

- test-generator:
  - Added a README for the test-generator ([680b618](https://github.com/apache/plc4x/commit/680b6184cabbc5cbcd748dfa6e6fae0fd8dd670a))
  - Added a README for the test-generator ([4ab57bd](https://github.com/apache/plc4x/commit/4ab57bda2ee9a6b2c1842286e6ee2e117a02f210))
  - Added a README for the test-generator ([e82f6e5](https://github.com/apache/plc4x/commit/e82f6e5900c54488ec3353eacea03a92dab3fa17))

- user:
  - Marked the connection-pool as removed stating with version 0.11.0 (#758) ([d74c0c6](https://github.com/apache/plc4x/commit/d74c0c6737e680ac96866eeffa2d8927bcb61daf)) ([#758](https://github.com/apache/plc4x/pull/758))

- api:
  - Started adding an alternate documentation for the SNAPSHOT version of the API ([2b687e9](https://github.com/apache/plc4x/commit/2b687e933d69fbc0cdcf06d6e431ab0c850f376c))

### Refactor

- plc4j/opcua:
  - reworked configuration/added driver context ([3ac65c0](https://github.com/apache/plc4x/commit/3ac65c04a0b08d1b6711f57b71fa51fb3ce620ad))
  - slight cleanup and qc fixing ([c9efef3](https://github.com/apache/plc4x/commit/c9efef33bc58128a388a176eae89be6fccbe1766))

- plc4go/opcua:
  - restructure secure channel ownership ([4cf782b](https://github.com/apache/plc4x/commit/4cf782b601d85a6661800c9224d9fe406ea03358))
  - use keyed logging ([a5a5e94](https://github.com/apache/plc4x/commit/a5a5e948343d79bb7c62f62e3b84524186ac83fc))

- opcua:
  - simplify pascal string ([93d7565](https://github.com/apache/plc4x/commit/93d7565a79f60c408562e4e2c65da927aced23b9))

- plc4go/spi:
  - add constructor to request transaction ([484af3b](https://github.com/apache/plc4x/commit/484af3bd711c1f68e35cd1e880194a830bfca6a8))
  - avoid issues when using executor and logging ([7d03458](https://github.com/apache/plc4x/commit/7d0345826641f97d4f49454635fc3a24946a9345))
  - slight cleanup of pool ([3ea774c](https://github.com/apache/plc4x/commit/3ea774cd4fa8a87bbd177c70401123157a1cc396))
  - move worker starting into a own method ([430655f](https://github.com/apache/plc4x/commit/430655fdf995364e800a7a5b5e938dd42cebcafb))
  - reorder methods for options ([1e2298b](https://github.com/apache/plc4x/commit/1e2298b3d8d82855eb3a4bca0fec2953c5dc6cab))
  - use getOrLeaveBool for high log precision too ([8535bef](https://github.com/apache/plc4x/commit/8535bef7c72829b9b4f019fec1f56cea5f461cbf))
  - abstract bufio.Reader through an interface ([62bc2ae](https://github.com/apache/plc4x/commit/62bc2ae77f8e416fe14caae9b0e588ce983b64e3))
  - move pool option to other options ([ac95770](https://github.com/apache/plc4x/commit/ac957705d69a694b5f8b3670cfa7e5e798df6cc9))
  - split up pool into multiple files ([ade5107](https://github.com/apache/plc4x/commit/ade510700aab3ef7200064501bd2f01152919aa9))
  - split up request transaction into separate file ([d7d5491](https://github.com/apache/plc4x/commit/d7d54912ee102858bd60a28aaea08713a6b32bcb))
  - improve code flow of read and write request ([ee99a37](https://github.com/apache/plc4x/commit/ee99a371eb168b897c6c1a0ee919d164207b2fcf))
  - move WorkerPool to own package ([c4bf3eb](https://github.com/apache/plc4x/commit/c4bf3ebd6206d470ab570c4d67710b95fe153c42))
  - move transaction manager to own package ([09147a8](https://github.com/apache/plc4x/commit/09147a8d953167c37a04eaaf49ecf4819f4b9be1))
  - simplify ReadBufferByteBased ([709866e](https://github.com/apache/plc4x/commit/709866edce4ef92c78139937dec5953ad020db82))
  - ignore callbacks in builder for code generation ([24eeddb](https://github.com/apache/plc4x/commit/24eeddbabe7182298de505f0334a90460b47ae79))
  - cleanup unused types ([794183a](https://github.com/apache/plc4x/commit/794183a15deaad0d6baa4ac8142da082a441edf8))
  - removed unsued type ([262cbc1](https://github.com/apache/plc4x/commit/262cbc1d9a4332048732ab7159dfb0fe7930801c))
  - introduce RequestTransactionRunnable ([fe482d9](https://github.com/apache/plc4x/commit/fe482d9305136b0a012f3dc86a005fdbab875ab1))
  - small cleanups ([0ddb758](https://github.com/apache/plc4x/commit/0ddb758cfd056cbd6507516383603781baa55807))
  - replace interface{} with any ([a24d2e1](https://github.com/apache/plc4x/commit/a24d2e18479630f272125c3494aa2e196d2eaa4b))
  - move ArrayContext to utils ([e36433a](https://github.com/apache/plc4x/commit/e36433a42513e1d83fd250635015342f8b91fc30))
  - cleanup browse query ([54095d0](https://github.com/apache/plc4x/commit/54095d0d629ec49960c294f7faddd3418027a0b7))
  - use Executor interface ([be227a7](https://github.com/apache/plc4x/commit/be227a7ad0989cad0e29e56fd197472477457272))
  - clean up interfaces of WorkerPool ([c9e0852](https://github.com/apache/plc4x/commit/c9e08520a7cddecf16f2e6cc76da1009bbdbf453))
  - clean up interfaces of WorkerPool ([12596b0](https://github.com/apache/plc4x/commit/12596b0828dfe18685eebf2602ac8925614dae7e))
  - generify WorkerPool ([7efcb36](https://github.com/apache/plc4x/commit/7efcb36f3ea0d61fb0e44f52902e23e6b843122c))
  - move worker related code into WorkerPool ([23f19db](https://github.com/apache/plc4x/commit/23f19dbd4c4b4625034e825320c0d24966239e64))
  - change the API of Parse ([c5fa0a3](https://github.com/apache/plc4x/commit/c5fa0a366a074bfb605dc94acf03a8733ba9f10a))
  - change the API of Serialize ([40ea49b](https://github.com/apache/plc4x/commit/40ea49b06d7d33b761fd43133205dba910ba3dc2))
  - change the API of Serialize ([f2510c2](https://github.com/apache/plc4x/commit/f2510c2c0c9a5b7bfaba1b50224690d378a4a6f2))
  - converted WriteBufferByteBased options to proper options ([e383adc](https://github.com/apache/plc4x/commit/e383adc494cea6f57d2b7e55c4cbe9da617c970f))

- plc4go/cbus:
  - split up browser code ([9036acb](https://github.com/apache/plc4x/commit/9036acbbd874341d43665a75aa4c7fec698cb936))
  - fix small code issues ([88fbf2c](https://github.com/apache/plc4x/commit/88fbf2cd0421363b42eacd5db989705f6c34355f))
  - split up reader into multiple methods ([fcd62ce](https://github.com/apache/plc4x/commit/fcd62ce597efa89e4064711b2ba30f1581995034))
  - align sal handler to mmi handler regarding the logic ([4f2d34f](https://github.com/apache/plc4x/commit/4f2d34f7921e75b9d668eda38c4154322629885e))
  - restructure Subscriber ([6c41c51](https://github.com/apache/plc4x/commit/6c41c51c557bac9571ae7a38c065e7c0469ab92b))
  - use *MessageCodec ([4eb09f8](https://github.com/apache/plc4x/commit/4eb09f8b70ad2367b5c23765ed8fd83bc00b70b3))
  - restructure Discoverer to make it testable ([553ec3d](https://github.com/apache/plc4x/commit/553ec3d85068079086c522a8eff5604414ff9518))
  - re-arrange message mapper and fix bug ([ee5ae76](https://github.com/apache/plc4x/commit/ee5ae76e4f248014812870a7d31a3b932f2dcbcd))
  - cleanup code ([a6ce77d](https://github.com/apache/plc4x/commit/a6ce77de59f045eadfc1643d5217e079b635cba8))
  - use spi.interface in browser ([2776597](https://github.com/apache/plc4x/commit/2776597304065d0e7d1e895e04d73100238bd64e))
  - move map encoded reply to message mapper ([874e2b9](https://github.com/apache/plc4x/commit/874e2b9565945d0a01bfbeeb3717df5d2ccfb43c))
  - struct fields should now start with a lower case ([f3e2130](https://github.com/apache/plc4x/commit/f3e2130c3f5fe5cd12974f38d6dc1978d143d684))
  - restructure reader ([7ffe3e8](https://github.com/apache/plc4x/commit/7ffe3e8f78a3da2be76ce0fcb11ec3e433d3e982))

- plc4j/profinet:
  - Added some comments and made the tests use the classloader to load the test-data instead of a fixed file-reference. ([3bf83a2](https://github.com/apache/plc4x/commit/3bf83a20a4e93211222ae224ada1bf5dd2d8fec9))

- plc4j/udp-transport:
  - Made it generally possible to open a UDP transport with a fixed local port ([23eee0a](https://github.com/apache/plc4x/commit/23eee0a89cb4039af3f6fddcb28a3fbf16bee61d))

- plc4go:
  - retire Connect on codec ([5b3633f](https://github.com/apache/plc4x/commit/5b3633f77dfe1b5af6ec907093434c19bcb62692))
  - use generated Stringers instead of hand written ones ([3aa6605](https://github.com/apache/plc4x/commit/3aa66052aa5bfd401ca8bfa103da828bd8aee380))
  - general QC pass ([41b82bd](https://github.com/apache/plc4x/commit/41b82bd16f83ed6e0b4f94e64fc6996b2e599523))
  - deduplicate package name out of struct names ([524c2c8](https://github.com/apache/plc4x/commit/524c2c836fa78e8843dc525471adf45da856fedd))
  - switch from global loggers to local loggers ([04d235d](https://github.com/apache/plc4x/commit/04d235d1a8e9472e00763eb9766b0ff0a39291ef))
  - switch from global loggers to local loggers ([95571e9](https://github.com/apache/plc4x/commit/95571e94942038cf712be7e07be1bebfaa8a8fc4))
  - streamline imports ([a69ada6](https://github.com/apache/plc4x/commit/a69ada6faa9ca58a22a5b481993a7bf48a24f33c))
  - streamline imports ([d89e40f](https://github.com/apache/plc4x/commit/d89e40f73f2c4000187cb75b8e7f8cdc9e6ecf80))
  - use constructors when possible for default types ([e9d689f](https://github.com/apache/plc4x/commit/e9d689fad485298b3f978f2c313dd15bbda8a95d))
  - use buffered channels when possible ([d915511](https://github.com/apache/plc4x/commit/d9155112910be36914c2190b4e2fe3d310708a8f))
  - convert some panics to errors ([1eaa3d5](https://github.com/apache/plc4x/commit/1eaa3d5897af3c788a6827617aafd4520b84b767))
  - code QC ([61a7ff0](https://github.com/apache/plc4x/commit/61a7ff02d10a6bf13caad64f572c4327b84db6c0))
  - avoid panics if possible ([7cc564f](https://github.com/apache/plc4x/commit/7cc564ff827c1bd1fcf0a0dc6788363290229000))
  - godoc and slight refactoring ([36d016b](https://github.com/apache/plc4x/commit/36d016ba20fe1e94f2ee4c876bb2ac82b39678e4))
  - fix hex formatting using the right fmt string ([0f3a330](https://github.com/apache/plc4x/commit/0f3a3304d0b4f94dc6c5e7c57e972e84d198c846))
  - switch from github.com/google/gopacket to github.com/gopacket/gopacket ([4318f4b](https://github.com/apache/plc4x/commit/4318f4b2b0b89a17c6f3276fda13e4c19af0cb64))

- plc4go/eip:
  - cleanup connection code ([39448aa](https://github.com/apache/plc4x/commit/39448aaf5cd1827255787d381c2b2580cee6f28a))

- s7:
  - Renamed some things and moved some classes to more reflect the structure of other drivers ([5a890b7](https://github.com/apache/plc4x/commit/5a890b77a448633333409d29653e3b97fe346c80))

- plc4j:
  - fixed some Q/A issues ([f605319](https://github.com/apache/plc4x/commit/f605319d210cec1a0165077a06ccebd38b87920f))
  - fixed some Q/A issues ([331f892](https://github.com/apache/plc4x/commit/331f892d34d108b8f478314c6b8a1a05f92afd01))
  - fixed some Q/A issues ([ab9026f](https://github.com/apache/plc4x/commit/ab9026f677d7ee4b36209afbaa1f93b214f1aa68))
  - fixed some Q/A issues ([0feee55](https://github.com/apache/plc4x/commit/0feee55a24e5bd9a09723823cbe064bddb453d43))
  - fixed some Q/A issues ([724c043](https://github.com/apache/plc4x/commit/724c043737cc324381b23e35c7f531a4794a6bfe))
  - get rid of jackson in test-utils ([5b72eaf](https://github.com/apache/plc4x/commit/5b72eaf8fe9d49cdeca7288d92e5ae2fe85079bd))
  - Continued streamlining the way Api types are parsed/serialized ([1997380](https://github.com/apache/plc4x/commit/199738041c7b493bb338b9e20962027f5ef9d225))
  - Commented in the stuff I commented out for testing. ([ec30925](https://github.com/apache/plc4x/commit/ec3092533fcc575807cbf141470c2de71177e3e5))
  - Continued porting the code to using contexts ... tried to fix the EIP Integration-Tests ([b583129](https://github.com/apache/plc4x/commit/b58312912f659565356c8fdacb1f0b503fa745c6))
  - Removed the Jackson Annotation stuff from all types ([edf5132](https://github.com/apache/plc4x/commit/edf5132bac0f0505a8eac257b33a43a5b59fbb76))

- plc4j/s7:
  - post merge cleanup ([3bb4097](https://github.com/apache/plc4x/commit/3bb40976a11ef574012bca8e04b9e69481284193))
  - post merge cleanup ([ee90580](https://github.com/apache/plc4x/commit/ee905808adf75aa30c9aae55544737214e3b5076))

- general:
  - streamline plc4go and plc4j spi serialization ([b303a33](https://github.com/apache/plc4x/commit/b303a331c060f5030800f62c91cbf7a5e5e596c0))

- plc4go/test:
  - simplify test TransportInstance signature ([1ad4fc7](https://github.com/apache/plc4x/commit/1ad4fc756f281f96d6bf42244344224be7ff7e01))

- plc4c:
  - Updated PLC4C to work with serializer and parse contexts in order to support variables like the "_latest" ([e139d4d](https://github.com/apache/plc4x/commit/e139d4de74159a797bc50aa13aa50ec644c8d601))

- code-gen:
  - Changed the code-gen to support a "_lastItem" variable inside array elements ([8dfef52](https://github.com/apache/plc4x/commit/8dfef5247a9b68c38a0847bfafaeac28edb64c38))
  - Changed the code-gen to support a "_lastItem" variable inside array elements ([138c619](https://github.com/apache/plc4x/commit/138c6195a4a1776ad34199bf5f3a4b58b2b86193))
  - Introduced two code-gen options: generate-properties-for-parser-arguments and generate-properties-for-reserved-fields ([d40f7d6](https://github.com/apache/plc4x/commit/d40f7d6097768845de1ed98e9ae646399d8a0b05))

- plc4j/driver/ads:
  - added the discovery types to the generated code ([9e57f6f](https://github.com/apache/plc4x/commit/9e57f6fbd0f3d19003587b675e3882bb51374ae8))

- protocol/open-protocol:
  - Added some more message-types ([41a110d](https://github.com/apache/plc4x/commit/41a110d77bcec6fbedcc9165803e4bd6fdb61237))

- plc4j/code-gen:
  - Refactored the way Encoding is handled and added "nullBytesHex" attribute support. ([102c041](https://github.com/apache/plc4x/commit/102c041d0662fbb8d9883e8d5f8bb9e52d2f9bbd))
  - Refactored the way Encoding is handled and added "nullBytesHex" attribute support. ([bcb14d3](https://github.com/apache/plc4x/commit/bcb14d398828bf54d68b9050c5a1861fa9ae000b))

- plc4go/bacnet:
  - cleanup package structure ([d1cbe41](https://github.com/apache/plc4x/commit/d1cbe41c7e14202e0e61a12ca6ca054d60484cf5))
  - smaller cleanups ([8e25dd2](https://github.com/apache/plc4x/commit/8e25dd229cba633173e6adb3a8fb27fb16e1ec35))
  - ported PDU object and encapsulate source and destination ([789403b](https://github.com/apache/plc4x/commit/789403bcb4ec45ac4a7f47c088440b11c3b838bd))
  - use generic messages for integration layer ([1e22184](https://github.com/apache/plc4x/commit/1e22184399e566a1d5663310f4c20d8267780951))
  - added more application code for protocol ([25c5e94](https://github.com/apache/plc4x/commit/25c5e941636f66d19ae70a0704a67362a24cd51f))
  - restructure code to hook in application layer ([c5639a6](https://github.com/apache/plc4x/commit/c5639a65fddbfa6f159649330ebe0cf31742e928))
  - remove useless code ([46cc1e1](https://github.com/apache/plc4x/commit/46cc1e13e3e73cdf296f9bbe414922109b106fe1))

- api:
  - refactor(api) PLC4Go API refactoring ([1424813](https://github.com/apache/plc4x/commit/142481333f54fd743d012d8520e6f4b088a608a9))

- plc4go/ads:
  - Refactoring of the go ADS drier ([554d756](https://github.com/apache/plc4x/commit/554d756cd215727869a5684d183fea5fd83caf84))
  - Refactoring of the go ADS drier ([4df0498](https://github.com/apache/plc4x/commit/4df04982253856a44ebc2ac90e33380821a34876))
  - Refactoring of the go ADS drier ([9c1296b](https://github.com/apache/plc4x/commit/9c1296b5d062f8b9aa528bb83c0f68af009d1af3))
  - Refactoring of the go ADS drier ([1612342](https://github.com/apache/plc4x/commit/1612342ae449d298520f8f433906d1970d20b9f7))
  - Refactoring of the go ADS drier ([2a9e84b](https://github.com/apache/plc4x/commit/2a9e84b439b46660f48a569aa9ecad82be5d36fa))
  - Refactoring of the go ADS drier ([8b0ba3f](https://github.com/apache/plc4x/commit/8b0ba3fd3da2254c70a8a7590ef7fd205ab4d9ed))
  - Refactoring of the go ADS drier ([b601c30](https://github.com/apache/plc4x/commit/b601c30c11076e537d66c4be9d823c0f56c7c920))

- times:
  - Refactoring of the was PlcValues are handled for TIME, DATE, TIME_OF_DAY and DATE_TIME types. ([76bcb48](https://github.com/apache/plc4x/commit/76bcb484b578dc0455e93127240675af02f89cfc))
  - Refactoring of the was PlcValues are handled for TIME, DATE, TIME_OF_DAY and DATE_TIME types. ([e73b569](https://github.com/apache/plc4x/commit/e73b569653dbb190d2e1136e53bc464286b89ad9))

- ads:
  - refactor(ads) PLC4Go Ads Driver update ([889511c](https://github.com/apache/plc4x/commit/889511c2928e3a2491192d42954a58dc188ff695))

- plc4xanalyzer:
  - moved override logic out of commands ([9de6d03](https://github.com/apache/plc4x/commit/9de6d0365f3cd6abb21486d1951bf358a2e7d45e))

## [rel/0.10](https://github.com/apache/plc4x/releases/tag/rel/0.10) - 2022-09-29 15:56:37

## What's Changed
* update readme to jdk 11 by @ottobackwards in https://github.com/apache/plc4x/pull/283
* Matching Pull request from type PR from the build tools. by @hutcheb in https://github.com/apache/plc4x/pull/285
* Fix/PLC4X 303 by @hutcheb in https://github.com/apache/plc4x/pull/284
* Fix documentation for S7 short address format by @markus-franke in https://github.com/apache/plc4x/pull/288
* Fix/plc4j-protocol-ads by @rmeister in https://github.com/apache/plc4x/pull/286
* build(deps): bump groovy.version from 3.0.8 to 3.0.9 by @dependabot in https://github.com/apache/plc4x/pull/278
* build(deps): bump jaxb-runtime from 2.3.3 to 3.0.2 by @dependabot in https://github.com/apache/plc4x/pull/282
* build(deps): bump nifi.version from 1.13.2 to 1.14.0 by @dependabot in https://github.com/apache/plc4x/pull/280
* build(deps): bump jna from 5.8.0 to 5.9.0 by @dependabot in https://github.com/apache/plc4x/pull/279
* build(deps): bump scala-library from 2.13.6 to 2.13.7 by @dependabot in https://github.com/apache/plc4x/pull/289
* Subscriptions for ADS in PLC4J by @rmeister in https://github.com/apache/plc4x/pull/265
* Fix for empty Passwords by @TorstenU in https://github.com/apache/plc4x/pull/294
* Add test for syncronized TransactionHandler by @hutcheb in https://github.com/apache/plc4x/pull/299
* Updated Notice file to 2022 by @Shoothzj in https://github.com/apache/plc4x/pull/306
* build(deps): bump karaf-maven-plugin from 4.3.2 to 4.3.5 by @dependabot in https://github.com/apache/plc4x/pull/302
* build(deps): bump asciidoctor-maven-plugin from 2.1.0 to 2.2.1 by @dependabot in https://github.com/apache/plc4x/pull/304
* URLDecoder.decode compatible with JDK 8 by @todoubaba in https://github.com/apache/plc4x/pull/307
* Possible fix for a race condition which lead a test to fail with some… by @JulianFeinauer in https://github.com/apache/plc4x/pull/309
* build(deps): bump maven-release-plugin from 3.0.0-M4 to 3.0.0-M5 by @dependabot in https://github.com/apache/plc4x/pull/310
* build(deps): bump calcite-core.version from 1.28.0 to 1.29.0 by @dependabot in https://github.com/apache/plc4x/pull/303
* build(deps): bump maven-enforcer-plugin from 3.0.0-M3 to 3.0.0 by @dependabot in https://github.com/apache/plc4x/pull/298
* build(deps): bump karaf-maven-plugin from 4.3.5 to 4.3.6 by @dependabot in https://github.com/apache/plc4x/pull/315
* build(deps): bump maven-compiler-plugin from 3.8.1 to 3.9.0 by @dependabot in https://github.com/apache/plc4x/pull/314
* build(deps): bump cmake-maven-plugin from 3.19.2-b1 to 3.22.1-b1 by @dependabot in https://github.com/apache/plc4x/pull/313
* build(deps): bump equalsverifier from 3.8.1 to 3.9 by @dependabot in https://github.com/apache/plc4x/pull/312
* build(deps): bump xmlunit.version from 2.8.4 to 2.9.0 by @dependabot in https://github.com/apache/plc4x/pull/311
* fix(docs): Updated the pcap docs for when setting up a developer env by @hutcheb in https://github.com/apache/plc4x/pull/318
* Dynamically determine ETS project number by @djexp in https://github.com/apache/plc4x/pull/321
* build(deps): bump asciidoctorj from 2.5.1 to 2.5.3 by @dependabot in https://github.com/apache/plc4x/pull/322
* feat(spi): First Draft of the encryption handler interface by @hutcheb in https://github.com/apache/plc4x/pull/319
* build(deps): bump sonar-maven-plugin from 3.9.0.2155 to 3.9.1.2184 by @dependabot in https://github.com/apache/plc4x/pull/324
* build(deps): bump animal-sniffer-maven-plugin from 1.20 to 1.21 by @dependabot in https://github.com/apache/plc4x/pull/323
* build(deps): bump httpcore from 4.4.14 to 4.4.15 by @dependabot in https://github.com/apache/plc4x/pull/320
* build(deps): bump annotations from 20.1.0 to 23.0.0 by @dependabot in https://github.com/apache/plc4x/pull/325
* build(deps): bump asciidoctorj-diagram from 2.1.2 to 2.2.1 by @dependabot in https://github.com/apache/plc4x/pull/316
* Fix typos by @cclauss in https://github.com/apache/plc4x/pull/326
* Delete .travis.yml by @cclauss in https://github.com/apache/plc4x/pull/327
* build(deps): bump guava from 31.0.1-jre to 31.1-jre by @dependabot in https://github.com/apache/plc4x/pull/329
* build(deps): bump maven-bundle-plugin from 5.1.2 to 5.1.4 by @dependabot in https://github.com/apache/plc4x/pull/330
* Fix typos by @cclauss in https://github.com/apache/plc4x/pull/328
* Upgrade Apache Camel to version 3.14.1 by @oscerd in https://github.com/apache/plc4x/pull/331
* build(deps): bump influxdb-client-java from 4.0.0 to 4.3.0 by @dependabot in https://github.com/apache/plc4x/pull/332
* build(deps): bump awaitility from 4.1.1 to 4.2.0 by @dependabot in https://github.com/apache/plc4x/pull/333
* build(deps): bump mockito.version from 4.2.0 to 4.4.0 by @dependabot in https://github.com/apache/plc4x/pull/334
* fix(plc4go/codegen): DefaultPlcWriteRequest interface conversion, cau… by @hongjinlin in https://github.com/apache/plc4x/pull/335
* fix(plc4go/codegen): Truly close connection by @hongjinlin in https://github.com/apache/plc4x/pull/338
* fix communication with LOGO 0AB7 and ISOonTCP tsap configuration by @sevenk in https://github.com/apache/plc4x/pull/308
* build(deps): bump slf4j.version from 1.7.32 to 1.7.36 by @dependabot in https://github.com/apache/plc4x/pull/342
* build(deps): bump spring-boot.version from 2.6.2 to 2.6.4 by @dependabot in https://github.com/apache/plc4x/pull/337
* Fix casting to byte for integer greater than 127 (byte is signed numb… by @alessandromnc94 in https://github.com/apache/plc4x/pull/339
* Feature/plc4py by @hutcheb in https://github.com/apache/plc4x/pull/343
* build(deps): bump zip4j from 2.9.1 to 2.10.0 by @dependabot in https://github.com/apache/plc4x/pull/345
* Feature/plc4py bh1 by @hutcheb in https://github.com/apache/plc4x/pull/341
* Bug Fix. Open too many files after complete by @Dustone-JavaWeb in https://github.com/apache/plc4x/pull/351
* Started with the PlcReadRequest by @hutcheb in https://github.com/apache/plc4x/pull/348
* build(deps): bump kotlin.version from 1.6.20 to 1.6.21 by @dependabot in https://github.com/apache/plc4x/pull/353
* Fix/plc4py venv by @hutcheb in https://github.com/apache/plc4x/pull/355
* Plc4py/plc field by @hutcheb in https://github.com/apache/plc4x/pull/358
* fix(plc4go): fix the error of value out of range when converting an i… by @hongjinlin in https://github.com/apache/plc4x/pull/363
* chore(deps): bump actions/cache from 2 to 3 by @dependabot in https://github.com/apache/plc4x/pull/368
* chore(deps): bump actions/checkout from 2 to 3 by @dependabot in https://github.com/apache/plc4x/pull/369
* chore(deps): bump github/codeql-action from 1 to 2 by @dependabot in https://github.com/apache/plc4x/pull/371
* chore(deps): bump actions/setup-java from 2 to 3 by @dependabot in https://github.com/apache/plc4x/pull/372
* chore(deps): bump BobAnkh/auto-generate-changelog from 1.1.0 to 1.1.1 by @dependabot in https://github.com/apache/plc4x/pull/370
* EIP adding read/write for STRINGS and LINT by @AndyGrebe in https://github.com/apache/plc4x/pull/367
* chore(deps): bump actions/github-script from 5.0.0 to 6.1.0 by @dependabot in https://github.com/apache/plc4x/pull/376
* chore(deps): bump github.com/rs/zerolog from 1.20.0 to 1.27.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/375
* chore(deps): bump github.com/icza/bitio from 1.0.0 to 1.1.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/373
* chore(deps): bump github.com/stretchr/testify from 1.7.1 to 1.7.2 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/374
* build(deps): bump jacoco-maven-plugin from 0.8.7 to 0.8.8 by @dependabot in https://github.com/apache/plc4x/pull/357
* chore(deps): bump maven-dependency-tree from 3.1.0 to 3.1.1 by @dependabot in https://github.com/apache/plc4x/pull/365
* chore(deps): bump maven-javadoc-plugin from 3.3.1 to 3.4.0 by @dependabot in https://github.com/apache/plc4x/pull/377
* chore(deps): bump download-maven-plugin from 1.6.7 to 1.6.8 by @dependabot in https://github.com/apache/plc4x/pull/379
* build(deps): bump maven-artifact from 3.6.3 to 3.8.5 by @dependabot in https://github.com/apache/plc4x/pull/359
* build(deps): bump spock-bom from 2.0-groovy-3.0 to 2.1-groovy-3.0 by @dependabot in https://github.com/apache/plc4x/pull/356
* build(deps): bump nifi.version from 1.16.1 to 1.16.2 by @dependabot in https://github.com/apache/plc4x/pull/380
* build(deps): bump actions/dependency-review-action from 1 to 2 by @dependabot in https://github.com/apache/plc4x/pull/385
* build(deps): bump nifi.version from 1.16.2 to 1.16.3 by @dependabot in https://github.com/apache/plc4x/pull/388
* build(deps): bump error_prone_annotations from 2.11.0 to 2.14.0 by @dependabot in https://github.com/apache/plc4x/pull/389
* build(deps): bump mvn-golang-wrapper from 2.3.9 to 2.3.10 by @dependabot in https://github.com/apache/plc4x/pull/390
* build(deps): bump dependency-check-maven from 7.1.0 to 7.1.1 by @dependabot in https://github.com/apache/plc4x/pull/391
* build(deps): bump maven-failsafe-plugin from 3.0.0-M5 to 3.0.0-M7 by @dependabot in https://github.com/apache/plc4x/pull/392
* build(deps): bump google-java-format from 1.11.0 to 1.15.0 by @dependabot in https://github.com/apache/plc4x/pull/395
* build(deps): bump maven-invoker-plugin from 3.2.2 to 3.3.0 by @dependabot in https://github.com/apache/plc4x/pull/394
* Update Dockerfile inline with requirments script changes by @vmpn in https://github.com/apache/plc4x/pull/387
* Updated java CIP write dataSize to match read dataSize by @AndyGrebe in https://github.com/apache/plc4x/pull/384
* fix plc4x_server build issues by @ottobackwards in https://github.com/apache/plc4x/pull/404
* Fix bug with select returning prematurely when device not sending any data by @vmpn in https://github.com/apache/plc4x/pull/386
* build(deps): bump maven-release-plugin from 3.0.0-M5 to 3.0.0-M6 by @dependabot in https://github.com/apache/plc4x/pull/400
* build(deps): bump mockito.version from 4.5.1 to 4.6.1 by @dependabot in https://github.com/apache/plc4x/pull/399
* build(deps): bump iotdb.version from 0.12.4 to 0.13.0 by @dependabot in https://github.com/apache/plc4x/pull/393
* build(deps): bump kotlin.version from 1.6.21 to 1.7.0 by @dependabot in https://github.com/apache/plc4x/pull/396
* build(deps): bump maven-artifact from 3.8.5 to 3.8.6 by @dependabot in https://github.com/apache/plc4x/pull/406
* build(deps): bump groovy.version from 3.0.10 to 3.0.11 by @dependabot in https://github.com/apache/plc4x/pull/408
* build(deps): bump maven-compiler-plugin from 3.10.0 to 3.10.1 by @dependabot in https://github.com/apache/plc4x/pull/409
* build(deps): bump jetty-util from 11.0.9 to 11.0.11 by @dependabot in https://github.com/apache/plc4x/pull/407
* build(deps): bump github.com/viney-shih/go-lock from 1.1.1 to 1.1.2 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/398
* build(deps): bump github.com/stretchr/testify from 1.7.2 to 1.7.5 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/403
* build(deps): bump maven-antrun-plugin from 3.0.0 to 3.1.0 by @dependabot in https://github.com/apache/plc4x/pull/411
* build(deps): bump buildnumber-maven-plugin from 1.4 to 3.0.0 by @dependabot in https://github.com/apache/plc4x/pull/410
* fixed copy-paste error in s7 config options by @tvormweg in https://github.com/apache/plc4x/pull/416
* build(deps): bump protobuf-java from 3.12.0 to 3.21.2 by @dependabot in https://github.com/apache/plc4x/pull/413
* build(deps): bump assertj-core from 3.22.0 to 3.23.1 by @dependabot in https://github.com/apache/plc4x/pull/415
* build(deps): bump byte-buddy from 1.12.10 to 1.12.12 by @dependabot in https://github.com/apache/plc4x/pull/414
* build(deps): bump maven-jar-plugin from 3.2.0 to 3.2.2 by @dependabot in https://github.com/apache/plc4x/pull/420
* build(deps): bump milo.version from 0.6.6 to 0.6.7 by @dependabot in https://github.com/apache/plc4x/pull/419
* build(deps): bump github.com/stretchr/testify from 1.7.5 to 1.8.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/417
* build(deps): bump properties-maven-plugin from 1.0.0 to 1.1.0 by @dependabot in https://github.com/apache/plc4x/pull/421
* build(deps): bump kafka.version from 3.2.0 to 7.1.2-ce by @dependabot in https://github.com/apache/plc4x/pull/422
* build(deps): bump BobAnkh/auto-generate-changelog from 1.1.1 to 1.2.1 by @dependabot in https://github.com/apache/plc4x/pull/426
* build(deps): bump github.com/schollz/progressbar/v3 from 3.8.6 to 3.8.7 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/427
* build(deps): bump maven-surefire-plugin from 3.0.0-M5 to 3.0.0-M7 by @dependabot in https://github.com/apache/plc4x/pull/424
* build(deps): bump jna from 5.10.0 to 5.12.1 by @dependabot in https://github.com/apache/plc4x/pull/425
* build(deps): bump jsoup from 1.15.1 to 1.15.2 by @dependabot in https://github.com/apache/plc4x/pull/423
* build(deps): bump BobAnkh/auto-generate-changelog from 1.2.1 to 1.2.2 by @dependabot in https://github.com/apache/plc4x/pull/431
* build(deps): bump camel.version from 3.17.0 to 3.18.0 by @dependabot in https://github.com/apache/plc4x/pull/430
* build(deps): bump maven-bundle-plugin from 5.1.6 to 5.1.7 by @dependabot in https://github.com/apache/plc4x/pull/429
* build(deps): bump woodstox-core from 6.2.8 to 6.3.0 by @dependabot in https://github.com/apache/plc4x/pull/428
* build(deps): bump influxdb-client-java from 4.3.0 to 6.4.0 by @dependabot in https://github.com/apache/plc4x/pull/436
* build(deps): bump crc from 1.0.2 to 1.0.3 by @dependabot in https://github.com/apache/plc4x/pull/435
* build(deps): bump github.com/gdamore/tcell/v2 from 2.5.1 to 2.5.2 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/438
* build(deps): bump github.com/schollz/progressbar/v3 from 3.8.7 to 3.9.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/437
* build(deps): bump kotlin.version from 1.7.0 to 1.7.10 by @dependabot in https://github.com/apache/plc4x/pull/432
* Apache Nifi plc4x record processor 2 by @inigoao in https://github.com/apache/plc4x/pull/439
* build(deps): bump protobuf-java from 3.21.2 to 3.21.5 by @dependabot in https://github.com/apache/plc4x/pull/443
* build(deps): bump actions/github-script from 6.1.0 to 6.1.1 by @dependabot in https://github.com/apache/plc4x/pull/449
* build(deps): bump maven-assembly-plugin from 3.3.0 to 3.4.2 by @dependabot in https://github.com/apache/plc4x/pull/448
* build(deps): bump maven-bundle-plugin from 5.1.7 to 5.1.8 by @dependabot in https://github.com/apache/plc4x/pull/447
* build(deps): bump camel.version from 3.18.0 to 3.18.1 by @dependabot in https://github.com/apache/plc4x/pull/446
* build(deps): bump mockito.version from 4.6.1 to 4.7.0 by @dependabot in https://github.com/apache/plc4x/pull/450
* build(deps): bump cmake-maven-plugin from 3.22.1-b1 to 3.23.2-b1 by @dependabot in https://github.com/apache/plc4x/pull/452
* build(deps): bump jakarta.activation-api from 1.2.2 to 2.1.0 by @dependabot in https://github.com/apache/plc4x/pull/451
* build(deps): bump Saxon-HE from 10.5 to 11.4 by @dependabot in https://github.com/apache/plc4x/pull/445
* build(deps): bump karaf-maven-plugin from 4.4.0 to 4.4.1 by @dependabot in https://github.com/apache/plc4x/pull/456
* build(deps): bump calcite-core.version from 1.30.0 to 1.31.0 by @dependabot in https://github.com/apache/plc4x/pull/455
* build(deps): bump byte-buddy from 1.12.12 to 1.12.13 by @dependabot in https://github.com/apache/plc4x/pull/454
* build(deps): bump woodstox-core from 6.3.0 to 6.3.1 by @dependabot in https://github.com/apache/plc4x/pull/453
* build(deps): bump maven-javadoc-plugin from 3.4.0 to 3.4.1 by @dependabot in https://github.com/apache/plc4x/pull/459
* build(deps): bump maven-site-plugin from 3.12.0 to 3.12.1 by @dependabot in https://github.com/apache/plc4x/pull/460
* build(deps): bump jsonassert from 1.5.0 to 1.5.1 by @dependabot in https://github.com/apache/plc4x/pull/457
* build(deps): bump error_prone_annotations from 2.14.0 to 2.15.0 by @dependabot in https://github.com/apache/plc4x/pull/458
* build(deps): bump nifi-nar-maven-plugin from 1.3.3 to 1.3.5 by @dependabot in https://github.com/apache/plc4x/pull/465
* build(deps): bump checker-qual from 3.21.4 to 3.24.0 by @dependabot in https://github.com/apache/plc4x/pull/464
* build(deps): bump milo.version from 0.6.6 to 0.6.7 by @dependabot in https://github.com/apache/plc4x/pull/463
* build(deps): bump gson from 2.9.0 to 2.9.1 by @dependabot in https://github.com/apache/plc4x/pull/462
* build(deps): bump joda-time from 2.10.14 to 2.11.0 by @dependabot in https://github.com/apache/plc4x/pull/461
* build(deps): bump maven-remote-resources-plugin from 1.7.0 to 3.0.0 by @dependabot in https://github.com/apache/plc4x/pull/466
* build(deps): bump github.com/gdamore/tcell/v2 from 2.5.2 to 2.5.3 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/470
* build(deps): bump actions/github-script from 6.1.1 to 6.2.0 by @dependabot in https://github.com/apache/plc4x/pull/473
* build(deps): bump github.com/rs/zerolog from 1.27.0 to 1.28.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/474
* build(deps): bump jsoup from 1.15.2 to 1.15.3 by @dependabot in https://github.com/apache/plc4x/pull/472
* build(deps): bump slf4j.version from 1.7.36 to 2.0.0 by @dependabot in https://github.com/apache/plc4x/pull/468
* build(deps): bump logback.version from 1.2.11 to 1.4.0 by @dependabot in https://github.com/apache/plc4x/pull/475
* build(deps): bump netty-bom from 4.1.79.Final to 4.1.80.Final by @dependabot in https://github.com/apache/plc4x/pull/476
* build(deps): bump exec-maven-plugin from 3.0.0 to 3.1.0 by @dependabot in https://github.com/apache/plc4x/pull/469
* build(deps): bump iot-device-client from 2.0.3 to 2.1.0 by @dependabot in https://github.com/apache/plc4x/pull/477
* build(deps): bump byte-buddy from 1.12.13 to 1.12.14 by @dependabot in https://github.com/apache/plc4x/pull/471
* build(deps): bump influxdb-client-java from 6.4.0 to 6.5.0 by @dependabot in https://github.com/apache/plc4x/pull/479
* build(deps): bump joda-time from 2.11.0 to 2.11.1 by @dependabot in https://github.com/apache/plc4x/pull/478
* build(deps): bump maven-dependency-tree from 3.1.1 to 3.2.0 by @dependabot in https://github.com/apache/plc4x/pull/480
* build(deps): bump milo.version from 0.6.6 to 0.6.7 by @dependabot in https://github.com/apache/plc4x/pull/481
* build(deps): bump milo.version from 0.6.7 to 0.6.8 by @dependabot in https://github.com/apache/plc4x/pull/483
* build(deps): bump github.com/schollz/progressbar/v3 from 3.9.0 to 3.10.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/484
* build(deps): bump github.com/schollz/progressbar/v3 from 3.10.0 to 3.10.1 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/488
* build(deps): bump antlr.version from 4.10.1 to 4.11.1 by @dependabot in https://github.com/apache/plc4x/pull/487
* build(deps): bump checker-qual from 3.24.0 to 3.25.0 by @dependabot in https://github.com/apache/plc4x/pull/486
* build(deps): bump jackson.version from 2.13.3 to 2.13.4 by @dependabot in https://github.com/apache/plc4x/pull/485
* build(deps): bump github.com/spf13/viper from 1.12.0 to 1.13.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/489
* build(deps): bump mockito.version from 4.7.0 to 4.8.0 by @dependabot in https://github.com/apache/plc4x/pull/490
* build(deps): bump netty-bom from 4.1.80.Final to 4.1.81.Final by @dependabot in https://github.com/apache/plc4x/pull/493
* build(deps): bump byte-buddy from 1.12.14 to 1.12.16 by @dependabot in https://github.com/apache/plc4x/pull/492
* build(deps): bump camel.version from 3.18.1 to 3.18.2 by @dependabot in https://github.com/apache/plc4x/pull/491
* build(deps): bump groovy.version from 4.0.4 to 4.0.5 by @dependabot in https://github.com/apache/plc4x/pull/495
* build(deps): bump calcite-core.version from 1.31.0 to 1.32.0 by @dependabot in https://github.com/apache/plc4x/pull/496
* build(deps): bump github.com/schollz/progressbar/v3 from 3.10.1 to 3.11.0 in /plc4go by @dependabot in https://github.com/apache/plc4x/pull/497
* build(deps): bump zip4j from 2.11.1 to 2.11.2 by @dependabot in https://github.com/apache/plc4x/pull/500
* build(deps): bump netty-bom from 4.1.81.Final to 4.1.82.Final by @dependabot in https://github.com/apache/plc4x/pull/501
* build(deps): bump dependency-check-maven from 7.1.2 to 7.2.0 by @dependabot in https://github.com/apache/plc4x/pull/502
* build(deps): bump slf4j.version from 2.0.0 to 2.0.1 by @dependabot in https://github.com/apache/plc4x/pull/503
* build(deps): bump logback.version from 1.4.0 to 1.4.1 by @dependabot in https://github.com/apache/plc4x/pull/504
* build(deps): bump protobuf-java from 3.21.5 to 3.21.6 by @dependabot in https://github.com/apache/plc4x/pull/507
* build(deps): bump iot-device-client from 2.1.0 to 2.1.1 by @dependabot in https://github.com/apache/plc4x/pull/506
* build(deps): bump jetty-util from 11.0.11 to 11.0.12 by @dependabot in https://github.com/apache/plc4x/pull/509
* build(deps): bump maven-jar-plugin from 3.2.2 to 3.3.0 by @dependabot in https://github.com/apache/plc4x/pull/510
* build(deps): bump apache-rat-plugin from 0.14 to 0.15 by @dependabot in https://github.com/apache/plc4x/pull/511
* build(deps): bump asciidoctorj from 2.5.5 to 2.5.6 by @dependabot in https://github.com/apache/plc4x/pull/512
* build(deps): bump dependency-check-maven from 7.2.0 to 7.2.1 by @dependabot in https://github.com/apache/plc4x/pull/513
* build(deps): bump byte-buddy from 1.12.16 to 1.12.17 by @dependabot in https://github.com/apache/plc4x/pull/514
* build(deps): bump junit.jupiter.version from 5.9.0 to 5.9.1 by @dependabot in https://github.com/apache/plc4x/pull/515
* build(deps): bump jaxb-runtime from 4.0.0 to 4.0.1 by @dependabot in https://github.com/apache/plc4x/pull/516
* build(deps): bump slf4j.version from 2.0.1 to 2.0.2 by @dependabot in https://github.com/apache/plc4x/pull/517
* build(deps): bump joda-time from 2.11.1 to 2.11.2 by @dependabot in https://github.com/apache/plc4x/pull/519
* build(deps-dev): bump json from 20220320 to 20220924 by @dependabot in https://github.com/apache/plc4x/pull/518
* build(deps): bump actions/github-script from 6.2.0 to 6.3.0 by @dependabot in https://github.com/apache/plc4x/pull/521
* build(deps): bump swagger-annotations from 1.6.6 to 1.6.7 by @dependabot in https://github.com/apache/plc4x/pull/522
* build(deps): bump commons-text from 1.9 to 1.10.0 by @dependabot in https://github.com/apache/plc4x/pull/523
* build(deps): bump slf4j.version from 2.0.2 to 2.0.3 by @dependabot in https://github.com/apache/plc4x/pull/524

## New Contributors
* @markus-franke made their first contribution in https://github.com/apache/plc4x/pull/288
* @rmeister made their first contribution in https://github.com/apache/plc4x/pull/286
* @TorstenU made their first contribution in https://github.com/apache/plc4x/pull/294
* @Shoothzj made their first contribution in https://github.com/apache/plc4x/pull/306
* @todoubaba made their first contribution in https://github.com/apache/plc4x/pull/307
* @djexp made their first contribution in https://github.com/apache/plc4x/pull/321
* @cclauss made their first contribution in https://github.com/apache/plc4x/pull/326
* @oscerd made their first contribution in https://github.com/apache/plc4x/pull/331
* @sevenk made their first contribution in https://github.com/apache/plc4x/pull/308
* @alessandromnc94 made their first contribution in https://github.com/apache/plc4x/pull/339
* @Dustone-JavaWeb made their first contribution in https://github.com/apache/plc4x/pull/351
* @AndyGrebe made their first contribution in https://github.com/apache/plc4x/pull/367
* @vmpn made their first contribution in https://github.com/apache/plc4x/pull/387
* @tvormweg made their first contribution in https://github.com/apache/plc4x/pull/416
* @inigoao made their first contribution in https://github.com/apache/plc4x/pull/439

**Full Changelog**: https://github.com/apache/plc4x/compare/rel/0.9...rel/0.10

### Feature

- plc4py:
  - Fixed python maven build. Next need to review and add comments. ([65684a0](https://github.com/apache/plc4x/commit/65684a026724e0d377ae15bf5c6d1a68421fcb8a)) ([#362](https://github.com/apache/plc4x/pull/362))
  - Added Generic for PlcConfiguration in PlcConnection ([906c075](https://github.com/apache/plc4x/commit/906c0754659af8368a19b7effae11eba7490cae0)) ([#362](https://github.com/apache/plc4x/pull/362))
  - Modbus Driver can now connect and disconnect ([ca564ce](https://github.com/apache/plc4x/commit/ca564ce31aefe7d0bb398830cd0a070b7c9f09e9)) ([#362](https://github.com/apache/plc4x/pull/362))
  - Fixed tests ([19baa51](https://github.com/apache/plc4x/commit/19baa51bb4650d53fa7acc26e30f39ff307808a6)) ([#362](https://github.com/apache/plc4x/pull/362))

- plc4go:
  - added remaining renderers for model ([1a54d16](https://github.com/apache/plc4x/commit/1a54d16a28625be50695fe6c88876f0629fc9a12))
  - implemente more serialize calls ([62c8c66](https://github.com/apache/plc4x/commit/62c8c66dd04bdf199090006d82f5492817844642))
  - added plc4xgenerator to generate the serialize functions ([7ae5070](https://github.com/apache/plc4x/commit/7ae50705a460e1f60f19eec947c8a09d57b6f3b8))
  - added GetSource to PlcSubscriptionEvent ([0e01d37](https://github.com/apache/plc4x/commit/0e01d379308048d845631fc5ba93ed300e2b3744))
  - added fmt.Stringer to plc_message ([7df4051](https://github.com/apache/plc4x/commit/7df4051865ec1b68fb99282debaf938267ac8264))
  - introduce Plc4xModelLog ([13809af](https://github.com/apache/plc4x/commit/13809afd28bf66eb6cc0ea0eed1558f4bd757f10))
  - added missing duration calls to plc_value ([b7fb8cf](https://github.com/apache/plc4x/commit/b7fb8cf995de0a99700c16a31d6b98891d3a885c))
  - added GetPLCValueType() PLCValueType to plc_value.go ([1f14f81](https://github.com/apache/plc4x/commit/1f14f8163d3e3c7bb70ca1a089f8d6ac6bcba26c))
  - Treat a single element PlcValueList as PlcValue ([88de075](https://github.com/apache/plc4x/commit/88de0758a4ba877094756d7e76fa4fec2fc460e3))
  - added DiscoverWithContext to driver manager ([5ff39f7](https://github.com/apache/plc4x/commit/5ff39f7a16ea0acee5f8e2817c02c471a1a125c2))
  - introduced ExecuteWithContext(ctx... calls ([a21cbc6](https://github.com/apache/plc4x/commit/a21cbc6c6d0715ad5f2b7b4505aef36af851b92d))
  - added net util for finding IPs ([1a041e8](https://github.com/apache/plc4x/commit/1a041e895908cede4e35ef73ae4bb68f402856bd))
  - exposed the SPI as it comes in handy when working with plc4x from time to time ([4c916f5](https://github.com/apache/plc4x/commit/4c916f58e65f82d9965146f8cf5c78f9375b037b))
  - plc4xbrowser should now remember last host and support more drivers ([54d0cf9](https://github.com/apache/plc4x/commit/54d0cf9d12c4be7be18e0aee6a35e0d5bf06419b))
  - first implementation of plc4xbrowser ([027ef30](https://github.com/apache/plc4x/commit/027ef30a90835d04f969dfc3161e440b271f42dd))
  - output progressbar on stderr for analyzer ([48a1a83](https://github.com/apache/plc4x/commit/48a1a83d2f4cd27f1fcbc4fdbf80627d7662c7f1))
  - add additional output for analyzer ([88f255b](https://github.com/apache/plc4x/commit/88f255b0777c152026c9f621a60ac14470c2b5fd))
  - added SO_REUSE support for udp ([f219a65](https://github.com/apache/plc4x/commit/f219a659dc07c2338bd128415bd84547570d55d8))
  - added protocol specific discovery options ([fdb52e5](https://github.com/apache/plc4x/commit/fdb52e5a2c58d2beb64f31989f6859fd36359c76))
  - expose Serializable ([1422a4e](https://github.com/apache/plc4x/commit/1422a4e9100045be7c6e821bc8583f11ae585b2e))
  - expose curstom write buffer ([7d0002a](https://github.com/apache/plc4x/commit/7d0002ac698c363bae408c03356fca92b00b5fcd))
  - expose read/write buffer via new interfaces ([badd0da](https://github.com/apache/plc4x/commit/badd0daa7663bdc838ea2d30a8de6d0b5d237ca9))
  - expose selected set of utils for public usage ([7e0ccdb](https://github.com/apache/plc4x/commit/7e0ccdb65726d53b341ee8680c03152556d64f5d))
  - fixes on golang ([c216cbe](https://github.com/apache/plc4x/commit/c216cbe3da782226476083b7e2c158130a6555b3))
  - add new factory method for write buffer ([a44a252](https://github.com/apache/plc4x/commit/a44a252ecb9181ae69c324f13830ab01e8e66b8c))

- plc-simulator/cbus:
  - unified cal data handling ([9878491](https://github.com/apache/plc4x/commit/9878491a9052fa1602239a13b46d50d041cd480e))
  - pick one of the available units for mmi ([cf1462e](https://github.com/apache/plc4x/commit/cf1462e87bd45089516e5bed97030ed24b44c700))
  - added support for installation mmi requests ([720c603](https://github.com/apache/plc4x/commit/720c603a8ac5b97bfec8bf200d1f52a5185f81a9))
  - cleanup code ([6ac271f](https://github.com/apache/plc4x/commit/6ac271f2f589495ebcf22eddeafd7fefbec1c425))
  - implemented direct command identify ([c9c9e77](https://github.com/apache/plc4x/commit/c9c9e772660ac204e70470f9529dbce29fd1f753))
  - added application filtering ([0a1817b](https://github.com/apache/plc4x/commit/0a1817b27b49333477d44d8b65a27e4068cd2c84))
  - implemented identify responses ([3c7065c](https://github.com/apache/plc4x/commit/3c7065cb47349ff041b749ddc6dfac73419f0fee))
  - implement proper option support ([8ddb301](https://github.com/apache/plc4x/commit/8ddb301752193e2c4de6bcdce260eacaa7a101fc))

- plc4go/connection-cache:
  - introduce connection cache options. ([5751ff3](https://github.com/apache/plc4x/commit/5751ff358c2b30425530e04edac6a789b7ee4327))
  - hook connection cache into a own logger ([b042f02](https://github.com/apache/plc4x/commit/b042f0253090dba12ee301ff0270a2c5834b10b0))

- cbus:
  - map sal data command type into address ([28708f9](https://github.com/apache/plc4x/commit/28708f9bf5fc6aec3f350661689d0588bc6fcf8b))
  - report levels as percentage (virtual) ([ef888b9](https://github.com/apache/plc4x/commit/ef888b9a5a71259e7841e17ce445e4828c0914f3))
  - introduced INFO_MESSAGES and NETWORK_CONTROL application ids ([5e61dc8](https://github.com/apache/plc4x/commit/5e61dc80dc997b69b0ab31b897cb67117a727224))
  - added number of arguments to sal command types ([86c4620](https://github.com/apache/plc4x/commit/86c462032279f9b581862ce69cf573c28664d305))
  - change some fields to c-bus now that they don't loose data ([0e48de1](https://github.com/apache/plc4x/commit/0e48de1c3b8914b64439889907a6d55e70005a11))
  - handle mmi status messages ([7d02a51](https://github.com/apache/plc4x/commit/7d02a51080c0e743236f95c2c9b3132d4ce82ff7))
  - properly parser parameter values ([f44d597](https://github.com/apache/plc4x/commit/f44d59737e38edd34e4950c478ed377830f5bc0f))
  - validate the checksum on read ([cfffc80](https://github.com/apache/plc4x/commit/cfffc808d294a3d4ca0766e970a967974c036be4))
  - implemented proper checksum support ([b0af4dc](https://github.com/apache/plc4x/commit/b0af4dc27f13c143018b34414807708992f78981))
  - implemented audio and video application ([409a709](https://github.com/apache/plc4x/commit/409a709cf8c5d92f63d5140d75d375935c93ae4d))
  - implemented heating application ([4d1b020](https://github.com/apache/plc4x/commit/4d1b020c27ffbf86bc8cec881d9b961806019d95))
  - implemented error hvac actuator application ([70a64ab](https://github.com/apache/plc4x/commit/70a64ab2f50128013b633b474bc263955b16f3da))
  - implemented error reporting application ([cc4f30d](https://github.com/apache/plc4x/commit/cc4f30d87e56475599ae8ba7e8dd8aec64cafcd6))
  - implemented measurement application ([238d80b](https://github.com/apache/plc4x/commit/238d80befe526880228ac2cbf26cda038358c8f3))
  - implemented pools spa ponds fountains control application ([61edca4](https://github.com/apache/plc4x/commit/61edca454de1ac051698100b862294a40112582a))
  - implemented measurement application ([65f7595](https://github.com/apache/plc4x/commit/65f7595859e44371c7231984994990477fe8f82b))
  - implemented irrigation control application ([16590cf](https://github.com/apache/plc4x/commit/16590cfaadb01dae1bf5be43f0aa618b5af12fa4))
  - implemented air conditioning application ([f9a1dfc](https://github.com/apache/plc4x/commit/f9a1dfcc7731d43ea89d82713480e61ec41a05be))
  - implemented air conditioning application ([f246821](https://github.com/apache/plc4x/commit/f24682115334e55339d33731cee8bd50201369e9))
  - implemented telephony application ([9edc926](https://github.com/apache/plc4x/commit/9edc9266343bb531270ede7dd33390cca70f2902))
  - implemented clock and timekeeping application ([218812e](https://github.com/apache/plc4x/commit/218812eb5dbf04d11f54d8d1f29b4419ba9721ba))
  - implemented media transport control application ([f4b45a4](https://github.com/apache/plc4x/commit/f4b45a446233ada6523fa0a67d7a2bd5de0ec1d6))
  - access control application ([819a6c4](https://github.com/apache/plc4x/commit/819a6c423fdb834e08039739b701291970380e19))
  - implemented ventilation application ([f8dbf89](https://github.com/apache/plc4x/commit/f8dbf89a9ad619a7efa175e88e1a5ef71dbcaf43))
  - implemented temperature broadcast application ([0119d0a](https://github.com/apache/plc4x/commit/0119d0a9258a9bc0200bc4bd1c40e4a0b3f57e89))
  - implemented enable control application ([79754ba](https://github.com/apache/plc4x/commit/79754bad42295480ab4388f3b778eebc34ba0c38))
  - implemented trigger control application ([9b9f882](https://github.com/apache/plc4x/commit/9b9f8820986461f9affefc6e007b78260b3f49ac))
  - implemented metering application ([ab42d2a](https://github.com/apache/plc4x/commit/ab42d2a951c64e7155871c4df0effd6cdb5a89b8))
  - implemented security application ([c18171e](https://github.com/apache/plc4x/commit/c18171e203e651198b072b58602aba9a7b35f47f))
  - added support for lighting label ([e2f61c4](https://github.com/apache/plc4x/commit/e2f61c4a583fd90889d0ef66efb7b28d8299d46c))
  - add support for cal data write ([a82a934](https://github.com/apache/plc4x/commit/a82a93435f6048f94040d751e74d5d4db929928c))
  - mapped more identify reply commands ([8a50d2a](https://github.com/apache/plc4x/commit/8a50d2a13689f26151bac695c60745087a539a0a))
  - added support for deprecated binary state status request ([b56ad32](https://github.com/apache/plc4x/commit/b56ad323c40c38923f53b5ce273f48b49e4c2dea))
  - added IdentifyReplyCommandOutputUnitSummary ([6d304e7](https://github.com/apache/plc4x/commit/6d304e7e02fd2254f8c36dbd5d99b63c4d5d3071))
  - added parameter relevant types ([de2ad06](https://github.com/apache/plc4x/commit/de2ad06146e7e44119630dd1217248b8007446fb))
  - added support for enum parameters ([7639143](https://github.com/apache/plc4x/commit/76391433e04c6437fc0c4009ded7501201b23136))

- plc4j/api:
  - Added a getArrayInfo() method to the PlcBrowseItem to provide information over array dimensions. ([cf2f2f5](https://github.com/apache/plc4x/commit/cf2f2f5aff13facf760eaba5378675e4bbfaa7a4))

- plc4go/spi:
  - improved virtual output from WriteBufferBoxBased ([0cfd466](https://github.com/apache/plc4x/commit/0cfd466c87a5414186598aa4a45b052f0e39bf41))
  - improved output of some types ([324ef77](https://github.com/apache/plc4x/commit/324ef77f3d93ad9e4a25502366ea3e9d13c0937d))
  - introduced WriteBufferPlcValueBased ([2a860ce](https://github.com/apache/plc4x/commit/2a860ced9333de0fdf4ea55270501bdefa9ae94f))
  - added new ConnectWithContext to message codec ([1f38cc9](https://github.com/apache/plc4x/commit/1f38cc9f8d0eccbd9f7dcebbcd19f1c300f1b641))
  - added new ConnectWithContext to transport instance ([de4fab5](https://github.com/apache/plc4x/commit/de4fab5f1fe2dfcf2b7a432081c47e80c5d190fe))
  - integrate ctx into DefaultCodec ([aa93c27](https://github.com/apache/plc4x/commit/aa93c27845fc415ec7a73af2c6950f41a673f88c))
  - propagate ctx down into implementations ([fbe964b](https://github.com/apache/plc4x/commit/fbe964bfaa0284d7a95b133e6cb2e84971f58bae))
  - implemented stringer for Default driver ([e8a3f67](https://github.com/apache/plc4x/commit/e8a3f67acec57b23c3fd945c8ac0f1a075383cc1))

- plc4go/cbus:
  - added source to events ([31438eb](https://github.com/apache/plc4x/commit/31438eba8c6a6c92b5a4ba42c337900747358be0))
  - added timeout to the read request on browse ([a0b7d0b](https://github.com/apache/plc4x/commit/a0b7d0b02806d1954a6f5b77adf3b7c52cdea477))
  - handle context in browse field ([a589161](https://github.com/apache/plc4x/commit/a58916131b43234e287da9ecab20f353970663ad))
  - map cal data using the WriteBufferPlcValueBased ([08ac93c](https://github.com/apache/plc4x/commit/08ac93ca0b1f4f8bad444fbc388dbd62f23cfed9))
  - SAL are now passed using the generic WriteBufferPlcValueBased ([4ca9925](https://github.com/apache/plc4x/commit/4ca9925e0a7235c96dfca949e6624971f38590b5))
  - added some more dummy implementations for value types ([6906355](https://github.com/apache/plc4x/commit/6906355d798538ddaaa35ddd309abb8dc275949a))
  - implemented value handler for writes ([96a452e](https://github.com/apache/plc4x/commit/96a452e4d539595a388818c4a0be3ddc9d5a44f7))
  - first implementation of writer support ([5fd053d](https://github.com/apache/plc4x/commit/5fd053d923cb483895b4ae9cb46fa9b54e64d8f1))
  - introduced sal field ([2651c29](https://github.com/apache/plc4x/commit/2651c29520b20f9047dcae873614ec504f069031))
  - implemented proper plc-value mapping for identify calls ([8a208fa](https://github.com/apache/plc4x/commit/8a208fa9952a8279543d8117df3041e8a7a10fc9))
  - properly implemented Discoverer ([87763de](https://github.com/apache/plc4x/commit/87763de69b0ca9272a3b9c0dd34bc712d657ac9c))
  - implemented Discovery ([840b324](https://github.com/apache/plc4x/commit/840b324701e63da8a0c38674c98ac6a47b518014))
  - fixed handling single confirmations ([0b222b0](https://github.com/apache/plc4x/commit/0b222b02b59d377e587f7aa98a60668f73a96b49))
  - added interceptor to browser ([78cd30b](https://github.com/apache/plc4x/commit/78cd30b6cef6ba0f266b0be411f16aacb59b727e))
  - added possibility to specify addresses in connection string ([fe18090](https://github.com/apache/plc4x/commit/fe18090840eb1ee9cbf527cd977307c41f44787e))
  - first draft of browse functionality ([c7e1d90](https://github.com/apache/plc4x/commit/c7e1d90c13d2b32e2e54686f2440e6fecd66bf93))
  - added new info field for browse ([456ba1d](https://github.com/apache/plc4x/commit/456ba1d50889e7b9c52621caff3a7eb824fa2a40))
  - first implementation of subscription for monitored SAL messages ([d516baa](https://github.com/apache/plc4x/commit/d516baae6790156d6aa865d77a56e154ca59d595))
  - properly map binary and level status responses to plcValue ([2bd067a](https://github.com/apache/plc4x/commit/2bd067a5eaf90753db364c2331b874a899aa6bbd))
  - progress on field handling ([38a8bdf](https://github.com/apache/plc4x/commit/38a8bdf4ec03b7e0705ec183c33108faaba6f0b0))
  - initial c-bus driver implementation ([6021471](https://github.com/apache/plc4x/commit/6021471df7036f2f285aa006404ae681fe231dd0))

- plc4j/ui:
  - Added a resource view for a selected PLC connection listing all resources it found using the "browse"-api in a tree-table. ([a4a7573](https://github.com/apache/plc4x/commit/a4a7573efecc148696805776609003262a3f9b58))

- general:
  - migrate to apache groovy ([36881ba](https://github.com/apache/plc4x/commit/36881ba512256ba95d4f62afc71948d540431221))
  - added new bacnet vendor ids ([a7c38fb](https://github.com/apache/plc4x/commit/a7c38fb6f604a542fa060af60cd92de654be9d2e))
  - Dynamically determine ETS project number (#321) ([7a6f03c](https://github.com/apache/plc4x/commit/7a6f03cf16e50e3476fe0af28085f3ef27ed209f)) ([#321](https://github.com/apache/plc4x/pull/321))
  - feature(plc4go): simulated driver, connection pool ([c76b28e](https://github.com/apache/plc4x/commit/c76b28ee125970b611ae608e5884121ef790facc))
  - mspec extension with parameterized type refs, assert, try, const ([22c9f00](https://github.com/apache/plc4x/commit/22c9f005abab73da1ab88e4b88d718657667fb3f))
  - add more tracer calls ([26bcca7](https://github.com/apache/plc4x/commit/26bcca72547d2a6f755567632a19664d9aeb9760))
  - inject tracer into freemarker templates ([212818e](https://github.com/apache/plc4x/commit/212818e71a428eced6de2fd35aac71708fa9d217))

- plc4j:
  - added addPreRegisteredConsumer as convenience method for pre register consumers ([8b8518c](https://github.com/apache/plc4x/commit/8b8518cd7b4d25e1fb3794d0af0005d0c850ea21))

- plc4go/bacnet:
  - update vendor ids with data from new url ([d4343a2](https://github.com/apache/plc4x/commit/d4343a2edd3b92783c5fbf266087b00a2b965092))
  - updated vendors ([dd11ad3](https://github.com/apache/plc4x/commit/dd11ad341172c72cd4c7162dea9c4127086bcc63))
  - update vendors ([8f8633a](https://github.com/apache/plc4x/commit/8f8633a71c70161b032316ef640c61e337f2293d))
  - first draft of reading single and multiple values ([6dc0a07](https://github.com/apache/plc4x/commit/6dc0a07fb217243a19e683b0ba58ce73f171d574))
  - added WhoHas and IHave to Discovery ([5335f6f](https://github.com/apache/plc4x/commit/5335f6ff8ac8b590990066b87464328521d24bd6))
  - Discovery should now display the remote ip at least for broadcast ([35d2185](https://github.com/apache/plc4x/commit/35d2185ab2316f98bca4388150b4d370316e8987))
  - Discovery should now parse the iam and produce a partial correct event ([2fbc348](https://github.com/apache/plc4x/commit/2fbc348dc7e00037e5fdb79688e8271e0bf3405f))
  - Discovery should now be able to handle broadcasts and unicasts ([f362cfa](https://github.com/apache/plc4x/commit/f362cfa5de1f97c14a939e37c644fca80ef3d296))
  - small changes on discovery using reuse of sockets ([1d069ad](https://github.com/apache/plc4x/commit/1d069addc03d380dab9af9fb459977796349200b))
  - initial implementation of discovery ([644d5c9](https://github.com/apache/plc4x/commit/644d5c9a37f3652d1d02f84ec53ab3123bb92b12))
  - fix broken static helper casts ([f883cc1](https://github.com/apache/plc4x/commit/f883cc1b5ed268c5d51d3236b635f2ec6278a9dd))
  - minor progress on bacnet ([7949d53](https://github.com/apache/plc4x/commit/7949d53b4c2627eff11494528e3d1ad393de02a1))

- plc-simulator:
  - graduate the plc-simulator into plc4j/utils ([9500fef](https://github.com/apache/plc4x/commit/9500fefa6b85f81aef2ea3720a8ae3a0eeb4eba2))
  - plc-simulator send out SAL messages every 5 seconds ([0e381aa](https://github.com/apache/plc4x/commit/0e381aa39e82f9ae62e8d05a80de88eea6a8112f))

- plc4go/knxnetip:
  - use context for discovery ([04fa66d](https://github.com/apache/plc4x/commit/04fa66d532cedf17c6a9dc518c5ef2b4bd1a38f9))

- plc4xbrowser:
  - added discover command ([e94ed04](https://github.com/apache/plc4x/commit/e94ed048d96daf25d2f54a7517db563b4936273c))
  - added discover command ([dc05bfa](https://github.com/apache/plc4x/commit/dc05bfa79a61b222bf7c94906da1b0a89a581d52))
  - increase display area ([ff7d444](https://github.com/apache/plc4x/commit/ff7d444e1fefbca98ef7db372d126197c88ff440))
  - set focus on ctrl+c ([1e2b13b](https://github.com/apache/plc4x/commit/1e2b13b23891ea9d4b816168daa93b7f46a51ea5))
  - added browse support ([a3d68a2](https://github.com/apache/plc4x/commit/a3d68a2d135d7f2d5ea4f190e316c1614507c438))
  - implemented write-direct ([8d2e713](https://github.com/apache/plc4x/commit/8d2e7134db61c303bdf7cdf4e48c28ba0b15bfb2))
  - made read produce a proper message ([5af52b7](https://github.com/apache/plc4x/commit/5af52b7652e823ddd2604848a40ecae3c126e2dd))
  - hook ctrl-c onto command input clear ([f2b5ec0](https://github.com/apache/plc4x/commit/f2b5ec0dc1d9104a6410bfecdc4cc26fc887d179))
  - added first iteration of read support ([5faa9e1](https://github.com/apache/plc4x/commit/5faa9e14acf366f9a58987329cb299950a0477db))
  - added subcommands for clear ([6aa9a57](https://github.com/apache/plc4x/commit/6aa9a57c5cd15d75962bd00d4591b4260e6c1aca))
  - auto register option for drivers ([34d70f7](https://github.com/apache/plc4x/commit/34d70f7414aa583afe31f0372bb3dbc0efa0690f))
  - add clear command ([a1e8b24](https://github.com/apache/plc4x/commit/a1e8b24f635521589425894a71d4aad24468f04f))
  - added several features ([0458e79](https://github.com/apache/plc4x/commit/0458e79860bdadaa5ff9e62097a19ed413a16ae9))
  - display last 10 commands and save them in config ([c53f4fa](https://github.com/apache/plc4x/commit/c53f4fafd87e71e993d1570b44b0958fc1270167))
  - added config options for logging and tracing ([55ded83](https://github.com/apache/plc4x/commit/55ded83c0491fa4b92327efb970e4c65742501ef))

- plc4xpcapanalyzer:
  - added frontend for interactive work with the cli ([decf45e](https://github.com/apache/plc4x/commit/decf45ef9ffaf6276db4f1abfbe42c841fdd2456))

- plc-simulator/s7:
  - support a alternate s7 port ([7adec29](https://github.com/apache/plc4x/commit/7adec29306c07da3679677751e3a7d7c7c710825))

- codegen/plc4j:
  - fixed issue with name collision ([ed60521](https://github.com/apache/plc4x/commit/ed60521dbab6be8f5aec727cde9371e25ffeb496))
  - save reserved in case they differ so they can be written out again using the same value ([8baeacf](https://github.com/apache/plc4x/commit/8baeacf1da74042dd41205cc9620b9523af37c0d))
  - add information of the type into the error message ([8636f3f](https://github.com/apache/plc4x/commit/8636f3fe8eb11ba7e6530da8435d0f40f25d47a5))

- codegen/plc4go:
  - save reserved in case they differ so they can be written out again using the same value ([a60ee58](https://github.com/apache/plc4x/commit/a60ee58e316211c4169311ff87f127d59b519a3b))
  - add information of the type into the error message ([5dfee90](https://github.com/apache/plc4x/commit/5dfee9081f763a5e262aa7ac1b78fb5b801089ee))
  - when a sub-type parse match fails output the actual values ([a7504a6](https://github.com/apache/plc4x/commit/a7504a6f4e4865e913c1ee92aecfa6352c2c3a53))
  - added an indirection through the write buffer for complex types. ([0801492](https://github.com/apache/plc4x/commit/08014922114185c25a68d924c52908c154d24db9))
  - added an indirection through the write buffer for complex types. ([77ff058](https://github.com/apache/plc4x/commit/77ff05801368faceb21bde5784f8e509241c1f1d))
  - wrap errors and added stacks where missing ([10764da](https://github.com/apache/plc4x/commit/10764daa1cb52f44e07a59a8b668afabd3b5705e))
  - log error to debug when optional is resetting for better traceability ([5b71b83](https://github.com/apache/plc4x/commit/5b71b83a8d2cd4274f1b6f4545ab88048e04b196))

- plc4xpcapanalyzer/cbus:
  - echo in merge requests should now be discarded so numbering won't get messed up ([c8893e0](https://github.com/apache/plc4x/commit/c8893e012624bc05192445d6d67e7df1e218e806))

- plc4j/cbus:
  - added simple cbus implementation to plc-simulator ([067ad4f](https://github.com/apache/plc4x/commit/067ad4f656f94b884db91dde46ece0180febf432))

- bacnet:
  - update vendors ([02f6191](https://github.com/apache/plc4x/commit/02f6191e4f9e12554c4347ae1176363b32e67efb))
  - use typeSwitch inherit feature for tags ([3694a28](https://github.com/apache/plc4x/commit/3694a28db8f37822e036abf8e230685f930a3854))
  - added missing accessors for constructed data types ([ddd8beb](https://github.com/apache/plc4x/commit/ddd8beb4ca9120e31de9f5daaeec5dcfbb0a1156))
  - added accessors for constructed data types ([2750970](https://github.com/apache/plc4x/commit/27509709ff19b2acd20f3b4ced249fa22fb2bd57))
  - renamed access credential all ([10de05b](https://github.com/apache/plc4x/commit/10de05bda6c82d941242b5d908554f9de84e116f))
  - added on more mapping for static helper ([b227400](https://github.com/apache/plc4x/commit/b22740023b5df5c1f120d1d05d1e5f3109bf97de))
  - added index 0 support for BACnetArray and validate lengths ([139832b](https://github.com/apache/plc4x/commit/139832bedae3ffac4c0d77667c75781afe0d8361))
  - pass index parameters to constructed data ([5e2003a](https://github.com/apache/plc4x/commit/5e2003a33296a29bafdf000cfd5afc5a70811e1b))
  - added safeguard against wrongly implemented constructed data ([d97c819](https://github.com/apache/plc4x/commit/d97c819ce83eadaf4d54a36c0e7c7108c98de185))
  - added shared property mapping ([b304355](https://github.com/apache/plc4x/commit/b304355310072ec288fd4f08230ab7a163282b0b))
  - added shared property mapping ([ef7538a](https://github.com/apache/plc4x/commit/ef7538ade321a5b00e9ec8aad5f3df26acf4a258))
  - added shared property mapping ([b7a5c45](https://github.com/apache/plc4x/commit/b7a5c45a5efdb82777cba6a0d8523544ada7fcb7))
  - added shared property mapping ([39f0757](https://github.com/apache/plc4x/commit/39f07579fac9355f916820589defd619fd8155d3))
  - added share property mapping ([4e32eb5](https://github.com/apache/plc4x/commit/4e32eb585427be7c9a43de07a7fdb6d8c3474568))
  - implemented unique properties for access zone type ([168884f](https://github.com/apache/plc4x/commit/168884ff52527868d8e7ae160b400f1a7ad80203))
  - implemented unique properties for accumulator type ([5fd804a](https://github.com/apache/plc4x/commit/5fd804a4478dbecabe249cf66b8ccd453cfab270))
  - implemented unique properties for averaging type ([c4f993b](https://github.com/apache/plc4x/commit/c4f993b8770646dcf7eb1d60dc6839f52980c3cb))
  - implemented unique properties for calendar type ([cc0fb1d](https://github.com/apache/plc4x/commit/cc0fb1d0da96818ea1893bd622c9738dc6a61f75))
  - implemented unique properties for channel type ([c365f85](https://github.com/apache/plc4x/commit/c365f851b8eb7be6fc494b2b51a663877f28580a))
  - implemented unique properties for command type ([fa695b0](https://github.com/apache/plc4x/commit/fa695b041e4a9b1d356d805daed945505c75ce7a))
  - implemented unique properties for credential data input type ([b9c9134](https://github.com/apache/plc4x/commit/b9c91345080d55a4128eba5ea97484650ea8b67f))
  - implemented unique properties for device type ([0be992d](https://github.com/apache/plc4x/commit/0be992d35b83cd054c527de93a5491414f30ae17))
  - implemented unique properties for elevator group type ([c94bf66](https://github.com/apache/plc4x/commit/c94bf669a76e64cc8df9e510e2059174621c0641))
  - implemented unique properties for escalator type ([ef8f17b](https://github.com/apache/plc4x/commit/ef8f17b60b62db0af2feea3debf8002bd5e9f692))
  - implemented unique properties for event enrollment type ([c2d39ad](https://github.com/apache/plc4x/commit/c2d39adcdb016d26be51ddc1827363a932d668ac))
  - implemented unique properties for file type ([c75cc7d](https://github.com/apache/plc4x/commit/c75cc7d3a632a3eae3f98450218c4e5a8ca78d23))
  - implemented unique properties for global group type ([7ed03bf](https://github.com/apache/plc4x/commit/7ed03bfc118e0efca155f7231fbfba25e99306af))
  - implemented unique properties for group type ([a8376a9](https://github.com/apache/plc4x/commit/a8376a9ae9b1eb22d6a4f2f45b0f4271477dd59c))
  - implemented unique properties for lift type ([9612626](https://github.com/apache/plc4x/commit/96126269331eefd441bdb0a7f0d94ce8e3267b10))
  - implemented unique properties for lightning output type ([7d74e34](https://github.com/apache/plc4x/commit/7d74e3419950de12066653b6f83d444760e766c1))
  - implemented unique properties for load control type ([b685af7](https://github.com/apache/plc4x/commit/b685af7a5054afcc73d8709ea7d467f8068627ab))
  - implemented unique properties for loop type ([89a1d3e](https://github.com/apache/plc4x/commit/89a1d3ea047ed21fa0ed35b3a1490023c21c9fbf))
  - implemented unique properties for network port type ([3ee1ecd](https://github.com/apache/plc4x/commit/3ee1ecdfaf0bb0d31976ec2977df05f387cc599e))
  - implemented unique properties for network security type ([d39d5d0](https://github.com/apache/plc4x/commit/d39d5d0c0858e0994a02921c91b5adcf35df03bd))
  - implemented unique properties for notification class type ([d679430](https://github.com/apache/plc4x/commit/d679430d565f7c92efd075d1a9756347e05671e5))
  - implemented unique properties for notification forwarder type ([334fc21](https://github.com/apache/plc4x/commit/334fc217dce73e77749d2e414a12c746ce845009))
  - implemented unique properties for program type ([c8276d5](https://github.com/apache/plc4x/commit/c8276d559fa8888ab24dbd059d0d495a6daa9890))
  - implemented unique properties for pulse converter type ([0bb2918](https://github.com/apache/plc4x/commit/0bb291893ab0e254a82b66863fd88d9fb8304928))
  - implemented unique properties for schedule object type ([176fb68](https://github.com/apache/plc4x/commit/176fb68759e511a72dfec90bcc130fe450f9e47e))
  - implemented unique properties for structured view object type ([8a683e2](https://github.com/apache/plc4x/commit/8a683e29d6867d089efeaf6487e5d2bdd2c8f184))
  - implemented unique properties for timer object type ([5994fd3](https://github.com/apache/plc4x/commit/5994fd32bc1e494c434f0e81b3dcdb5fad6d2712))
  - implemented unique properties for binary lightning output object type ([a70d1aa](https://github.com/apache/plc4x/commit/a70d1aaf0db305470d02bbd5a921ea15bcabf5d3))
  - implemented unique properties for bit string object type ([f1950f4](https://github.com/apache/plc4x/commit/f1950f481cb927f1e0f66b0b5b9714b0f1b02ba2))
  - implemented unique properties for access user object type ([737d078](https://github.com/apache/plc4x/commit/737d0787b201b74c2e472ac52ac27f0db8e116c6))
  - implemented unique properties for access point object type ([535da88](https://github.com/apache/plc4x/commit/535da8840f91c270b0aa0fbdb4bc2a6a721540b3))
  - implemented unique properties for access door object type ([003ddf1](https://github.com/apache/plc4x/commit/003ddf17d2fc3e94bbdefd48aa12f60597419d9e))
  - implemented unique properties for access credential object type ([8626057](https://github.com/apache/plc4x/commit/862605723ce7de09a55db72c1d9dc4657e1090db))
  - implemented several properties ([e8ba250](https://github.com/apache/plc4x/commit/e8ba2500bbbfa8c834d309e820541e65635526ce))
  - implemented some more property mappings and defaulted all properties to a validation error ([f4cd2a9](https://github.com/apache/plc4x/commit/f4cd2a922dbd8469105d08ec7ada1e62df23399d))
  - added first implementation of BACnetConstructedDataAnalogInputAll ([4d946c0](https://github.com/apache/plc4x/commit/4d946c0e26f849d7a98227812a2e16083952c624))
  - implemented all bit strings similar to enums ([65f5526](https://github.com/apache/plc4x/commit/65f55262db012d67a51d24bd7a54b5e3ce298966))
  - implemented BACnetPriorityArray ([53381f1](https://github.com/apache/plc4x/commit/53381f1c90d9b8a0f4625085590dea4b9fedd9a4))
  - mapped units to out of service ([5771964](https://github.com/apache/plc4x/commit/577196469f502e1ca1bf47f0b0ac072a64db8d50))
  - mapped controlled variable value ([fc22ff5](https://github.com/apache/plc4x/commit/fc22ff5277ef26336224d7cc2f82a837f020384d))
  - mapped units to BACnetEngineeringUnits ([a64de11](https://github.com/apache/plc4x/commit/a64de11e6f39b477ac4c919bbd386d0960c8f26c))
  - add unknown enums to ApduType enum ([ef1e4b3](https://github.com/apache/plc4x/commit/ef1e4b394f9a56b405937c46c028d9e2b1a94da1))
  - introduce ApduType enum ([5153cf7](https://github.com/apache/plc4x/commit/5153cf71df50fa81fc76ea1c5863331aa42f159e))
  - implement BACnetConstructedData*All for all objects as placeholder ([ee9fb24](https://github.com/apache/plc4x/commit/ee9fb247c3323f6117e45d2ad90c3754c837a70b))
  - implement BACnetConfirmedServiceRequestVTData ([e6b90c6](https://github.com/apache/plc4x/commit/e6b90c69f32aa88873bb908cad56c8cd66668c4e))
  - implement BACnetConfirmedServiceRequestVTOpen ([001eedb](https://github.com/apache/plc4x/commit/001eedb0b2bfc302228db2da6f9931dee73d7021))
  - implement BACnetConfirmedServiceRequestVTOpen ([ba8b2dd](https://github.com/apache/plc4x/commit/ba8b2dde6b6a2a12100254b6d54f660d9dbe56e3))
  - implement BACnetServiceAckVTData ([0003c82](https://github.com/apache/plc4x/commit/0003c82e2c1f258dbc5f28f1beaccf0eb216b497))
  - implement BACnetServiceAckVTOpen ([b0c9709](https://github.com/apache/plc4x/commit/b0c9709999386a4003330e664b8c288244479c1e))
  - implement BACnetServiceAckGetAlarmSummary ([1e5d6f5](https://github.com/apache/plc4x/commit/1e5d6f548c0285926a9b586c55ebcad9ef7af033))
  - implement BACnetUnconfirmedServiceRequestWriteGroup ([cbaedc9](https://github.com/apache/plc4x/commit/cbaedc9a39169862990680225c7c2be0620ee705))
  - implement BACnetUnconfirmedServiceRequestWriteGroup ([7d9c063](https://github.com/apache/plc4x/commit/7d9c063f5750fdbbe2da019a3cf40ddad9169eb5))
  - implemented BACnetConfirmedServiceRequestConfirmedTextMessage ([6e5cfee](https://github.com/apache/plc4x/commit/6e5cfeeadfc2bb9a1b7a022fb91536be6e301b93))
  - implemented BACnetConfirmedServiceRequestLifeSafetyOperation ([9709a9c](https://github.com/apache/plc4x/commit/9709a9cb6f195d6d3e05af485237113774eb8777))
  - implemented BACnetServiceAckGetEnrollmentSummary ([cd12935](https://github.com/apache/plc4x/commit/cd12935f3667af21304af7f772d5e252114743d7))
  - implemented BACnetConfirmedServiceRequestGetEnrollmentSummary ([1af0a40](https://github.com/apache/plc4x/commit/1af0a4028e23a2b165c3b4bd278f804adff8589b))
  - implemented BACnetNotificationParametersChangeOfTimer ([d8beb08](https://github.com/apache/plc4x/commit/d8beb08b32fa449b874df87eb0aceca8cfd4d5f5))
  - implemented BACnetNotificationParametersChangeOfDiscreteValue ([056a4b4](https://github.com/apache/plc4x/commit/056a4b42c24d93b3a1abf82df90a1acfcbaf072e))
  - implemented BACnetNotificationParametersChangeOfReliability ([4a727b1](https://github.com/apache/plc4x/commit/4a727b157c7d97327d389f0a8bb3c32d5ca2cf7f))
  - implemented BACnetNotificationParametersChangeOfStatusFlags ([839386e](https://github.com/apache/plc4x/commit/839386e61791f6043dd1f69ece974ba7eb2fe884))
  - implemented BACnetNotificationParametersChangeOfCharacterString ([c2d5017](https://github.com/apache/plc4x/commit/c2d50175496164bfcbc354f938668b3264ccb292))
  - implemented BACnetNotificationParametersUnsignedOutOfRange ([1505202](https://github.com/apache/plc4x/commit/15052026934a07ae80a36dd532a17c702bdaadb3))
  - implemented BACnetNotificationParametersSignedOutOfRange ([996bcbc](https://github.com/apache/plc4x/commit/996bcbc94c81dfb8cf06cc569701216179bc5c84))
  - implemented BACnetNotificationParametersDoubleOutOfRange ([01d0114](https://github.com/apache/plc4x/commit/01d01140dc7d26c9955101978c0072382fbbd2e0))
  - implemented BACnetNotificationParametersAccessEvent ([c964dc5](https://github.com/apache/plc4x/commit/c964dc5de61528ca6528124c2497a8c9582324b3))
  - implemented BACnetNotificationParametersChangeOfLifeSafety ([1e97df3](https://github.com/apache/plc4x/commit/1e97df332649abc907a26747c9ef5e32d2b5aba6))
  - put basic tags in own mspec ([852ed66](https://github.com/apache/plc4x/commit/852ed662f7250b9b3f570eaaa49fd63158b2f535))
  - implemented BACnetPropertyStates ([0a88392](https://github.com/apache/plc4x/commit/0a88392d30fdda1768c916afb29a8da4e973e0b3))
  - added helpers for vendor id ([3a188a5](https://github.com/apache/plc4x/commit/3a188a5b43d6cb4d9a15e5b5a1378d3aca3baed4))
  - use vendor id enum whenever possible ([337ed95](https://github.com/apache/plc4x/commit/337ed9511c4b1638a3ef47026a0682acd34ef1ca))
  - implement BACnetConstructedDataSubordinateList ([0f62ec0](https://github.com/apache/plc4x/commit/0f62ec0c46d850262fe718e326df5ccecf88722b))
  - small cleanups of static helpers ([6230a05](https://github.com/apache/plc4x/commit/6230a0505dea3d3d658ef0f4814b508a970026da))
  - added helper functions for segmentation ([09c1e79](https://github.com/apache/plc4x/commit/09c1e79bc747c07cff47de6dd2950ed2eeb287e9))
  - added mspec-code-generation for defined enums (DRY) ([115a4d7](https://github.com/apache/plc4x/commit/115a4d76e67ae870d3df65ca051f856f763ed618))
  - implemented remaining direct defined enums ([94ca4f0](https://github.com/apache/plc4x/commit/94ca4f0170d1c026bb22bd2da74d9988640b59ed))
  - split enums up into public and private enums ([11ae7e0](https://github.com/apache/plc4x/commit/11ae7e09814bc2fd253331e18594743c3eb7b2ea))
  - externalize enums in own mspec and add missing skeletons ([02d6616](https://github.com/apache/plc4x/commit/02d66166674b7b5885799e6ea6c859528d9e7e39))
  - pull vendor id list from the internet ([3b0a651](https://github.com/apache/plc4x/commit/3b0a651954b61e7ec6daeaf2cc4c99e4f7e1c4b9))
  - migrated BACnetAction to enum ([a392d38](https://github.com/apache/plc4x/commit/a392d382e6a6542caa357a74f3235170319d58c3))
  - implemented improved generic method to handle enums which can have extended values ([34e26f1](https://github.com/apache/plc4x/commit/34e26f171eb220769e41704ca502030e9b048088))
  - implemented support for BACnetConstructedDataReliability ([feb585b](https://github.com/apache/plc4x/commit/feb585b3f0bbf376ea6c1d5ebe1117b25dc8e043))
  - implemented support for BACnetConstructedDataLifeSafetyStateEntry ([675fb09](https://github.com/apache/plc4x/commit/675fb09ab2f3ba90554b59cbf087f0635b6ff8a2))
  - implemented support for ACCEPTED_MODES ([bc11e48](https://github.com/apache/plc4x/commit/bc11e48a89edefc27601b7cc7d4d371da56cee5f))
  - implemented support for LIFE_SAFETY_ALARM_VALUES ([7a75b2c](https://github.com/apache/plc4x/commit/7a75b2c57491640e0b74aad5cbd5983f051b0036))
  - added todos and validations for not yet implemented datatypes ([3b5dd93](https://github.com/apache/plc4x/commit/3b5dd93fae8d42a98617d266d0311eaad1c33527))
  - implemented BACnetConfirmedServiceRequestCreateObject and BACnetConfirmedServiceRequestDeleteObject ([4a78124](https://github.com/apache/plc4x/commit/4a78124dcbc0d4b1e8db88150e1c57e783a34d09))
  - implemented BACnetUnconfirmedServiceRequestUnconfirmedCOVNotificationMultiple ([a84ba0c](https://github.com/apache/plc4x/commit/a84ba0c406bea5033935277ddb994c2b90807221))
  - implemented BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple ([48521c3](https://github.com/apache/plc4x/commit/48521c33691c8ff42a72430c1871cb82d6b5206f))
  - implemented BACnetConfirmedServiceRequestSubscribeCOVProperty ([9a623d1](https://github.com/apache/plc4x/commit/9a623d16f9a4c2c832ed53803243dfcbde69f724))
  - map abort and reject reason to enums ([59437f7](https://github.com/apache/plc4x/commit/59437f7d5e7c0384b64b1722eaf71ca677eaedf9))
  - reworked error handling ([78b0326](https://github.com/apache/plc4x/commit/78b0326022608a7d2c5e4c8cfd4f4b8637cd285f))
  - implemented read range ([ddfad99](https://github.com/apache/plc4x/commit/ddfad995ae93cb33cdd61860c9b0b967ecf254e7))
  - implemented PrivateTransfer ([c75c2f0](https://github.com/apache/plc4x/commit/c75c2f05bf7e7288010d9f573211bf38027782bb))
  - major overhaul of type mappings ([b4d63f3](https://github.com/apache/plc4x/commit/b4d63f34b126f4c54fc843f8af396ccb0fc21170))
  - added property mapping placeholders ([e0d8f82](https://github.com/apache/plc4x/commit/e0d8f822a40430d4873267a7d38b517ebc4c645a))
  - implemented BACnetConstructedDataLifeSafetyZoneMembers and BACnetConstructedDataLifeSafetyZoneMembers ([986009e](https://github.com/apache/plc4x/commit/986009eccaedd211c4aad7836c8e926397a584a9))
  - implemented BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple ([1d39569](https://github.com/apache/plc4x/commit/1d395694d1bea0ae9836227c811442f78cefa48b))
  - implemented BACnetConfirmedServiceRequestAcknowledgeAlarm ([6a899a1](https://github.com/apache/plc4x/commit/6a899a134795eb8f117fac0d0d78be9c4ef322d6))
  - implemented BACnetUnconfirmedServiceRequestUnconfirmedEventNotification ([347196d](https://github.com/apache/plc4x/commit/347196de2dac87c62fe3b7b3a6d5d3c5a0a20ed5))
  - implemented BACnetUnconfirmedServiceRequestUTCTimeSynchronization ([e95d551](https://github.com/apache/plc4x/commit/e95d551936035d523298dca9866e4d87c9565125))
  - introduced BACnetUnconfirmedServiceChoice enum ([d5d7a1e](https://github.com/apache/plc4x/commit/d5d7a1eab5286884d3992739403195b819fb1e0c))
  - introduced BACnetConfirmedServiceChoice enum ([fec9e3e](https://github.com/apache/plc4x/commit/fec9e3eb0ef93460d83239cd34137660df06aafc))
  - added support for BACnetConfirmedServiceRequestAddListElement and BACnetConfirmedServiceRequestRemoveListElement ([d6e48ad](https://github.com/apache/plc4x/commit/d6e48adbf0ba7c4dfb629b629e405d0aea6cc56e))
  - minor cleanups ([69c255c](https://github.com/apache/plc4x/commit/69c255c9af6e47727c4dc2e9aafdc431c0ae18a5))
  - added support for more read replies ([0f46d4d](https://github.com/apache/plc4x/commit/0f46d4d6d3abf2b662ab7822d39acfb6df97cc2e))
  - refined BDT and FDT support ([5c895a3](https://github.com/apache/plc4x/commit/5c895a30046a96558a7dcefabab082ab303a7bf5))
  - initial BDT and FDT support ([261bfe7](https://github.com/apache/plc4x/commit/261bfe7e888ec228569c8446bbc2ee4998b4ca78))
  - added support for BACnetReadAccessResult ([ad729e2](https://github.com/apache/plc4x/commit/ad729e296175483914a8d1a292321cf25ecb04b6))
  - added support for BACnetConfirmedServiceRequestWritePropertyMultiple ([af95da3](https://github.com/apache/plc4x/commit/af95da34b963f556f1f6a708465593338cfded24))
  - introduce enums for MaxSegmentsAccepted and MaxApduLengthAccepted ([8024e07](https://github.com/apache/plc4x/commit/8024e07c3f72b6723bf4254e12cea3b4059e0518))
  - added helper for date and time ([201ded2](https://github.com/apache/plc4x/commit/201ded2fbe2a49da6c956d6a2e175cb82c415d1e))
  - added null context tag ([1e426c2](https://github.com/apache/plc4x/commit/1e426c270a8ab474dd749bb467f7a611ee67fbd7))
  - added helper for bit string data types ([3cb822b](https://github.com/apache/plc4x/commit/3cb822b8f79c6352bc0c0df427c96d03fbb582b6))
  - added helper for string data types ([b34ce11](https://github.com/apache/plc4x/commit/b34ce11a7de289024c7f5ccf2c2569f89d5f140b))
  - added helper for boolean point data types ([d2ec227](https://github.com/apache/plc4x/commit/d2ec2273c0d0ca95cf10683a522a48e62a09f9bd))
  - added helper for floating point data types ([1c229ce](https://github.com/apache/plc4x/commit/1c229ce66ae3afefe2e798e684a4e885fdc01868))
  - added helper for opening and closing tags. ([373d618](https://github.com/apache/plc4x/commit/373d618dc947768955f757f1a43dd8d51cf876c3))
  - fixed wrong calculations of length headers ([843cd77](https://github.com/apache/plc4x/commit/843cd77a46227c240fb2492536b89d42aa792469))
  - fixed wrong calculations of length headers ([bbee159](https://github.com/apache/plc4x/commit/bbee159e64d37232a1d8fdcc9db3ae01c1d1d0be))
  - added utility method for creating property identifiers ([e84103b](https://github.com/apache/plc4x/commit/e84103b0a7cbb4e9a7c56278120785409651c44d))
  - Add util methods in static helper for integer based values ([ba3154d](https://github.com/apache/plc4x/commit/ba3154d1e9f5de64501171892a054a5396db3f77))
  - Add util methods in static helper for object identifiers ([c22fedf](https://github.com/apache/plc4x/commit/c22fedf809c5a5171a9e8e5e77294d268b07a7b7))
  - initial support for segmented messages ([e89b6f4](https://github.com/apache/plc4x/commit/e89b6f4921d487818fa3c3b18b47c478d76a7d2e))
  - added BACnetConstructedDataEventTimestamps ([b45dbcd](https://github.com/apache/plc4x/commit/b45dbcd591e15297edbd1386a51e98a5d5c3106c))
  - added support for bacnet actions ([576d7ae](https://github.com/apache/plc4x/commit/576d7aeee8c61e03c949d7829020f86c5d144c89))
  - added BACnetNotificationParametersFloatingLimit ([fd417b8](https://github.com/apache/plc4x/commit/fd417b8059a68533e839270ec6282b569db6d539))
  - added more typing ([a57725a](https://github.com/apache/plc4x/commit/a57725a838eff5572a2f80a6c56bebb23e4073cf))
  - support for life support zones ([0806f51](https://github.com/apache/plc4x/commit/0806f517249153048ded8a2ee442e72909f42e00))
  - implemented BVLCWriteBroadcastDistributionTable ([18f7ac7](https://github.com/apache/plc4x/commit/18f7ac7c8ef0bfc5c7a75d76f0529bcbaae07194))
  - added missing ack ([4fd5f35](https://github.com/apache/plc4x/commit/4fd5f35665157ea6e4196974b2aaa81323c39e3a))
  - implementation improvements ([ea77517](https://github.com/apache/plc4x/commit/ea77517a2a716f7a8ae441d0ab33ccf38c65fc03))
  - added confirmed service request unknown ([28fa90f](https://github.com/apache/plc4x/commit/28fa90f69e5a5e326edd03a60cc1e14d34b534b2))
  - added ADPU unknown ([e208e62](https://github.com/apache/plc4x/commit/e208e62f22ca0000658f30be8324d156e1ad5b84))
  - add missing network routing messages ([43bcee8](https://github.com/apache/plc4x/commit/43bcee816a28456d7b44529a3e6bef5ca94b0dfe))
  - support for more object types ([56d85b3](https://github.com/apache/plc4x/commit/56d85b3c45de20f63c53344beae8e55f34219a23))
  - worked on notifications ([782518e](https://github.com/apache/plc4x/commit/782518eaba464a1d263a7cf1a76c2d6c9003d0e3))
  - constructed data generic parsing ([aa23ae7](https://github.com/apache/plc4x/commit/aa23ae7a0e4aa5ac93af36c4487b49e72da6a883))
  - initial support for constructed data ([80307c2](https://github.com/apache/plc4x/commit/80307c21b3c03ffe6681f9bada7aefe8a1bfd520))
  - fix merge messup ([3b0aa11](https://github.com/apache/plc4x/commit/3b0aa1112111d59e7be6467fc86ed5d0810503bc))
  - small optimizations ([d412f5c](https://github.com/apache/plc4x/commit/d412f5c7e2ebfa7398fd8f40f255124679304884))
  - worked on read property ([868bc10](https://github.com/apache/plc4x/commit/868bc1055b4ca6146e2bca8df2ae5cfd151ddb8c))
  - implemented BVLCResultCode ([85f56d5](https://github.com/apache/plc4x/commit/85f56d5f5c508353ae4917293f8c1323e5b9ffdb))
  - implemented BVLCRegisterForeignDevice ([52cd8ac](https://github.com/apache/plc4x/commit/52cd8ac8678eaeb9b47978a57616f1dda8a6f0d2))
  - add pcap based bacnet test ([ffb68a6](https://github.com/apache/plc4x/commit/ffb68a6f9441704e2f01c47041cf4ac8cad6e70d))
  - implement BACnetConfirmedServiceRequestReadProperty ([1972700](https://github.com/apache/plc4x/commit/1972700c229bbeefb861dd1f5c8fad21345248ba))
  - implement BACnetConfirmedServiceRequestReinitializeDevice ([5204018](https://github.com/apache/plc4x/commit/52040186852e1cfab64c6dd850a86512135230e1))
  - implement BACnetConfirmedServiceRequestAtomicWriteFile ([c7efc32](https://github.com/apache/plc4x/commit/c7efc32547d2fd4c292501879423bfd7e90a8ffe))
  - implement BACnetConfirmedServiceRequestAtomicWriteFile ([fa4c139](https://github.com/apache/plc4x/commit/fa4c139ea51b3ded524262bd1af7f17e433d2a22))
  - implement BACnetTagApplicationDate and BACnetTagApplicationTime ([f19d8f7](https://github.com/apache/plc4x/commit/f19d8f79f77ff18d356e858974396c785aa4042e))
  - implemented I-Have ([7006b32](https://github.com/apache/plc4x/commit/7006b32d6ac6524abf511d65f4ea5b195a847989))
  - work in BACnetUnconfirmedServiceRequestWhoIs ([7d41e7c](https://github.com/apache/plc4x/commit/7d41e7c26025d27981b43910f4c38889d3f580ca))

- plc4go/plc4xpcapanalyzer:
  - introduce package mapper which is used to "fake" a healthy communication pattern. ([93f523d](https://github.com/apache/plc4x/commit/93f523daae2b7bb001997c0da1ab79afeccd80f3))
  - analyzer should now be able to handle segmented messages (for c-bus) ([929558b](https://github.com/apache/plc4x/commit/929558b9c84f68dec444a9cd89c7c5b5d9d8e706))
  - added more options to cbus ([95c80a8](https://github.com/apache/plc4x/commit/95c80a8a40868397906d36f0da41c2c0f1ae7c07))

- c-bus:
  - added missing crc checks (parsing only) ([820d272](https://github.com/apache/plc4x/commit/820d27234e00e3148dfda22794446586343061ed))
  - proper identify support ([092276f](https://github.com/apache/plc4x/commit/092276f410d23a43b606e3f9f571ac30ebf9941e))
  - minor changes on c-bus ([6576810](https://github.com/apache/plc4x/commit/6576810a9f99b438c978ff64cf99d4901bc73253))
  - small improvements on spec ([913bd94](https://github.com/apache/plc4x/commit/913bd948a2d480ad9d7c0220dd41abf89bd00bf3))
  - refined c-bus mspec ([ace05fb](https://github.com/apache/plc4x/commit/ace05fb77d2211447c2c06db628853d3c89257c4))
  - add new root type ([fcae433](https://github.com/apache/plc4x/commit/fcae433c7a48b476522f4e31b4cab4be02eea9d5))

- spi:
  - changed the way a byte is rendered ([a243df0](https://github.com/apache/plc4x/commit/a243df0a92bc2753ddd8bd9d64a6ee7ee2703ba4))
  - First Draft of the encryption handler interface (#319) ([5ef3f31](https://github.com/apache/plc4x/commit/5ef3f317617764f665947a52ef2960e60e1d432a)) ([#319](https://github.com/apache/plc4x/pull/319))
  - Added highlight function to hex to highlight byte positions ([9450441](https://github.com/apache/plc4x/commit/945044127ca30db37ac6f2fdb9c8945800489063))

- bacnet/plc4go:
  - added new vendors ([f4645cb](https://github.com/apache/plc4x/commit/f4645cb47b14432d9788044dada3b6af71295347))

- codegen:
  - if a case name in a typeswitch start with * prefix the owner name ([c764f40](https://github.com/apache/plc4x/commit/c764f408de61628d8aedaa18814998b2ee28619e))
  - streamline protocol implementations by providing convenience methods ([fa0ce26](https://github.com/apache/plc4x/commit/fa0ce2624c64679a5e8fb6dafcb6224487a7a2d2))
  - added possibility to split up mspecs ([ae9c2e6](https://github.com/apache/plc4x/commit/ae9c2e6121377052f4fa5716ba129ab43f2e9698))
  - changed validation to fail parsing conditionally ([c21a184](https://github.com/apache/plc4x/commit/c21a1842f2dd9f0d4755bd634e1a38acc48b24bc))
  - more accurate error reportings ([ed71d06](https://github.com/apache/plc4x/commit/ed71d06cdd98460204869aaeddb500bb2614a379))
  - Fixed more issues in the C# code gent ([b7dcde9](https://github.com/apache/plc4x/commit/b7dcde96cb063838a0aa692977db04ea0ea2b278))
  - Fixed more issues in the C# code gent ([34d0887](https://github.com/apache/plc4x/commit/34d0887f2af52f268eb8c19cf48277d82852ce26))
  - Fixed more issues in the C# code gent ([71058d8](https://github.com/apache/plc4x/commit/71058d8936353e74def9e11f9a866ed038a0e7d8))
  - Fixed an issue in the C# code gent ([b080c87](https://github.com/apache/plc4x/commit/b080c8783ac61ab8b1712b3d22d8a79b9be77521))
  - Fixed an issue in the C# code gent ([83811e5](https://github.com/apache/plc4x/commit/83811e59a8f5e35f62cc40e184581916aaf42f92))
  - Got the C code-generation working again ([9777dc2](https://github.com/apache/plc4x/commit/9777dc26e842f17945008bb1d15459a2fc81cf97))
  - Got the C code-generation working again ([b2635e3](https://github.com/apache/plc4x/commit/b2635e3d5d75526eb854ffa1591aa9ed9463f1f8))
  - Tried getting stuff working in C ;-) ([efda339](https://github.com/apache/plc4x/commit/efda33948a4a3f020bee4dba976e97829c7434cf))
  - Removed some unreferenced code ([10389a2](https://github.com/apache/plc4x/commit/10389a22767b79378d987f7cc849ed14842c5a4d))
  - Fixed how byte-based manualArray fields are parsed ([b997544](https://github.com/apache/plc4x/commit/b9975440180547daf13824025cf6ec47225dff9f))
  - Fixed how data-io fields are referenced ([9ba8cf2](https://github.com/apache/plc4x/commit/9ba8cf28fc86c3341aa7158e6b6e895d784fc4bc))
  - Kept on working on the migration ([1426944](https://github.com/apache/plc4x/commit/14269444c218ffad5f8abbd756b852dcfd77624a))
  - Fixed go code generation after refactoring ([749e034](https://github.com/apache/plc4x/commit/749e034e0866c2f711bafc4fde53703fa55fee58))
  - Fixed array type handling and data-io code generation ([7c59bc3](https://github.com/apache/plc4x/commit/7c59bc310776f4ffc7e4e55e04134f9de0311463))
  - Started clening up the type references ([8afd22e](https://github.com/apache/plc4x/commit/8afd22ee362c369b77cbf3d14bc99f65108d6ae8))
  - Ensured the VariableLiterals are able to provide their type. ([7e21421](https://github.com/apache/plc4x/commit/7e2142156e75934b866d445bd5eee6adb883039e))
  - Ensured the VariableLiterals are able to provide their type. ([0a99ac6](https://github.com/apache/plc4x/commit/0a99ac62094c98abcd41071d5dd85b10c237afc7))
  - added info method to do logging during template development ([aa6175a](https://github.com/apache/plc4x/commit/aa6175a3d58f17f219132f5394067d403ce431d5))
  - introduce TypeContext to transport more information out of code generation ([e88ad99](https://github.com/apache/plc4x/commit/e88ad992ab566f1bd0789a24467f331556f82185))
  - added wildcard support for type switch ([30817ad](https://github.com/apache/plc4x/commit/30817ad49dde044b6bf614fc9e55970594f458c4))
  - add peek field ([3a7b261](https://github.com/apache/plc4x/commit/3a7b2611381fa667d2db083ddca18630148df8f4))
  - make the reader and writers a bit more verbose ([d8946d4](https://github.com/apache/plc4x/commit/d8946d4d6ce8b62c316f02122e410ae2b5d40009))
  - Added new validation field ([3a4d6de](https://github.com/apache/plc4x/commit/3a4d6de2b8104626f076bf77de9c7953afc0a450))

- plc4j/eip:
  - Updated java CIP write dataSize to match read dataSize (#384) [PLC4X-341] ([0489ec6](https://github.com/apache/plc4x/commit/0489ec6d12a05096675cfb96ebe208465cbef2d6)) ([#384](https://github.com/apache/plc4x/pull/384))

- protocols/eip:
  - EIP adding read/write for STRINGS and LINT (#367) ([5f75bba](https://github.com/apache/plc4x/commit/5f75bba748e027eb9146a46754e3307ab2a9de5d)) ([#367](https://github.com/apache/plc4x/pull/367))

- plc4go/tools:
  - initial version of plc4xpcapanalyzer which can be applied to pcaps to check if plc4go can handle it ([eda641e](https://github.com/apache/plc4x/commit/eda641e1d02d90dbae942a5a0a01a0d75854c743))

- modbus:
  - Added some initial tests for Modbus-RTU and Modbus-ASCII ([48839ac](https://github.com/apache/plc4x/commit/48839acfcf5be2bab2bb47e7661cde7a9e6dad83))
  - Updated the RELEASE_NOTES a bit ([4f01cd9](https://github.com/apache/plc4x/commit/4f01cd947ad190d9d935425015f80fd77c228e7f))
  - Added first working support for modbus-rtu and modbus-ascii (currently only tested with tcp transport) ([9798d9b](https://github.com/apache/plc4x/commit/9798d9b192c7ea9f85567e639f77a6e717fe4621))
  - Adjusted the modbus discovery to work with the renamed modbus-tcp driver ([9ae4a55](https://github.com/apache/plc4x/commit/9ae4a55faf8b2804987aa7f1a6d31ce7b539b4e7))
  - Split up "modbus" into "modbus-tcp", "modbus-rtu" and "modbus-ascii". Renamed "modbus" into "modbus-tcp". ([e3e46f3](https://github.com/apache/plc4x/commit/e3e46f3ec47b93187e3d21ad68b22e0870790c57))

- protocol/modbus:
  - Started working on ModbusRTU/ASCII ([6f38d2e](https://github.com/apache/plc4x/commit/6f38d2e8be91ec7b4f231ba41ee39e9d9b169bd3))

- protocol/c-bus:
  - Continued implementing the driver core ([e6f462f](https://github.com/apache/plc4x/commit/e6f462fcdfe42f57ca55aca68e026111b5383c71))
  - Implemented the ApplicationId part ([befcbd8](https://github.com/apache/plc4x/commit/befcbd8c27cae9a0315e06e2399b6409b6d4e7ea))
  - Got more CAL Data tests working ([42013eb](https://github.com/apache/plc4x/commit/42013ebbe10c2102e51dc432ca38e82dfb3884ea))
  - Got more CAL Data tests working ([fa465b3](https://github.com/apache/plc4x/commit/fa465b3bdc7ccd7663bd298f747b24c455db16bf))
  - Got the CAL Data tests working ([dcac0ea](https://github.com/apache/plc4x/commit/dcac0ea8d6648ce08a380e0d95ba703ca2778a12))
  - Fleshed out the CALData Header ([34f5719](https://github.com/apache/plc4x/commit/34f571955f848e869e3a5c1d7c139e2bd7773af7))
  - Fleshed out the SALData Header ([a55c297](https://github.com/apache/plc4x/commit/a55c2974946144a09c27941d5acedf08f3abcb81))
  - Got the build working again ([1292be0](https://github.com/apache/plc4x/commit/1292be0d01be3545d0cc1f0a932c0ab118598db4))
  - Got the mspec a bit more compilable. ([fb83aa6](https://github.com/apache/plc4x/commit/fb83aa63f2205eea6be108670d3f7bc2aeec14db))
  - Some finetuning (implemented srchk) ([16c3457](https://github.com/apache/plc4x/commit/16c34570f496f6a08e23d4489b122f9c8c7f8d74))
  - added more types, fixed some issues ([562f665](https://github.com/apache/plc4x/commit/562f6657c14c7b5c9e6a1e78119d9d877b15acf1))
  - added more types ([a29d692](https://github.com/apache/plc4x/commit/a29d6922de159bd92ca40672267b74e9204ad8e0))
  - added more base types ([44f06b3](https://github.com/apache/plc4x/commit/44f06b3b42474439704cd19f248b880a1fa01876))
  - Initial commit. ([53d0aef](https://github.com/apache/plc4x/commit/53d0aef9aa9a8b8de823b2ef2c9c817ac81f2291))

- plc4go/c-bus:
  - added IdentifyReplyCommand ([9dd8a07](https://github.com/apache/plc4x/commit/9dd8a07e968ccafeb4a5662bc2619b15f2fffcbe))
  - added c-bus protocol ([1fb2c8b](https://github.com/apache/plc4x/commit/1fb2c8b5a9fd787016e2637382f40a765fd46d70))

- protocols/c-bus:
  - added NetworkProtocolControlInformation ([abe3853](https://github.com/apache/plc4x/commit/abe38538929bd5919a4485647fe45d23935d7fe1))
  - added ExtendedFormatStatusReply ([c34636c](https://github.com/apache/plc4x/commit/c34636c86672532d7810b7b36e0e1d4e9708863b))
  - added StandardFormatStatusReply ([ed223d5](https://github.com/apache/plc4x/commit/ed223d56f4c3099896e69ebd6db206d67baa72e1))

- protocols/cbus:
  - added initial CALData ([7623887](https://github.com/apache/plc4x/commit/76238878429e2672e35ee19f7ad0b10d86dae547))
  - added initial SALData ([5232441](https://github.com/apache/plc4x/commit/5232441cd6ab6b5c1c0c53ed2ccfc55147217d65))

- protocol:
  - added sanity checks on protocols ([8e6551d](https://github.com/apache/plc4x/commit/8e6551d70801f5744d63b9ebe69785e6d607d865))

- plc4j/modbus:
  - Cleaned up and added some more code comments explaining why things are done the way they are. ([13d53e6](https://github.com/apache/plc4x/commit/13d53e6b00e3f87f424312a58cbed58c1a00bb38))
  - Updated the discovery to also find the unit-id for a modbus device. ([51657e8](https://github.com/apache/plc4x/commit/51657e8bcdb13856115dca8d12667bcc70df6631))
  - Made the draft threadsafe ([add57e7](https://github.com/apache/plc4x/commit/add57e74d68f099ceef683c0e1cdb3b7b88e2586))
  - Added a first experiment for modbus discovery ([395b975](https://github.com/apache/plc4x/commit/395b975254db369dcb3b3bde01b3d214a9866a35))

- bacnet/plc4go/codegen:
  - added factory methods to calculated headers ([cdb1e74](https://github.com/apache/plc4x/commit/cdb1e74712f578256e4b1cb08896b77d870088fb))

- plc4j/knx:
  - Enabled the EST parser to process encrypted projects (And fixed a bug resulting in it not correctly assigning datatypes to group-addresses) ([4d00f95](https://github.com/apache/plc4x/commit/4d00f95a8671ed97ede9a52d1dcafdf32fe596f3))
  - Enabled the EST parser to process encrypted projects (And fixed a bug resulting in it not correctly assigning datatypes to group-addresses) ([5ad93d3](https://github.com/apache/plc4x/commit/5ad93d372c33ca5662b702fbb5ba7e99e35ccae7))

- plc4j/codegen:
  - Store arguments as fields to make them usable by virtual fields ([6807c3a](https://github.com/apache/plc4x/commit/6807c3ab9b080d6c400eb7c7b563ffa47b45352c))
  - add function to test if enum is mappable ([907742a](https://github.com/apache/plc4x/commit/907742af8f750c82a0397b82b8944f5ced59cf74))
  - render virtual fields in boxes ([17086bd](https://github.com/apache/plc4x/commit/17086bd6b2df6ec283cffd7be5d4782d02196b84))
  - added support for date, time, datetime ([cf5a1ad](https://github.com/apache/plc4x/commit/cf5a1ad8a220c843778dff7673dae0f2591a816e))
  - use serialization based toString() ([072de49](https://github.com/apache/plc4x/commit/072de4903049fae370e61d0bce88cd5d33ff7b38))
  - re-add exception handling for expressionless optional field ([346d4ab](https://github.com/apache/plc4x/commit/346d4ab7c74b5f35047348995da7c2d1fb49ff84))
  - re-add exception handling for expressionless optional field ([07f9f7e](https://github.com/apache/plc4x/commit/07f9f7e93e8aeb8556c3a4370ac064b1e4c32292))
  - format java sources and optimize import after codegen ([4156659](https://github.com/apache/plc4x/commit/415665932c892ca8c044138ff1f884b1a73ac413))

- plc4j/plc4x:
  - Added a skeleton for a plc4x protocol driver, which will be used to relay PLC4X API requests from one system to another ([6652fa4](https://github.com/apache/plc4x/commit/6652fa40d4189a26a62cfa83dbf777ba56113dfa))

- code-generation:
  - Added the concept of a "protocol version" to the protocol modules and the code-generation plugin. ([d97d76b](https://github.com/apache/plc4x/commit/d97d76bacc26a256b06374ab8ded33fba0bc289f))

- plc4go/connection-pool:
  - Added a matching API component to make the pool usable outside the plc4go project ([019018f](https://github.com/apache/plc4x/commit/019018ffff965d28d0b449478abb476ca17e0dd1))
  - Added a PlcConnectionPool for go ([c9c0374](https://github.com/apache/plc4x/commit/c9c037491ed7abb9aa4fc290f7b057720da7903c))

- plc4j/pcap-transport:
  - added filter option ([e13ee55](https://github.com/apache/plc4x/commit/e13ee55de99d60177da675db92da2ed933784f4d))

- sast:
  - enable trivy scanner ([a77a399](https://github.com/apache/plc4x/commit/a77a399b3c1c2f9c6f5b6f65b3205a31968f32c7))

- plc4go/knx:
  - update manufactures ([fc50316](https://github.com/apache/plc4x/commit/fc5031630fb2e18e91f624c9498bc9d930dbc714))

- plc4j/ads:
  - Subscriptions for ADS in PLC4J (#265) ([5d4eb0a](https://github.com/apache/plc4x/commit/5d4eb0a7f32e1a959d6020214fb8f79db55e505b)) ([#265](https://github.com/apache/plc4x/pull/265))

- plc4go/codegen:
  - render virtual fields in boxes ([13d37b4](https://github.com/apache/plc4x/commit/13d37b41f29b2211c4e5479824afcf2914b07588))

- plc4go/pcap:
  - dump packages on debug for better inspection ([cfcd6e3](https://github.com/apache/plc4x/commit/cfcd6e363f2d47a87b7ce1c5ada70830f6191e60))
  - Pcap transport should work now ([cdc8419](https://github.com/apache/plc4x/commit/cdc841963b83b3c1f68526582f1450719597a740))

- plc4go/boxwriter:
  - write strings directly and omit hex values ([b7bf197](https://github.com/apache/plc4x/commit/b7bf197604677f62d84219ab043c0d914e61a063))

- plc4go/try-assert:
  - try/assert is now working in golang ([c076381](https://github.com/apache/plc4x/commit/c0763817575aeb27189149828d0d94d139f2cbd0))

- ci:
  - add changelog action for automated changelogs ([648409b](https://github.com/apache/plc4x/commit/648409bbff6fd41bf5633755894120ca95f53e31))
  - add workflow which ensures plc4x builds on every platform with every version ([4f7b609](https://github.com/apache/plc4x/commit/4f7b6095480953ff90c17a3932b2f76327fc839f))

- plc4x/codegen:
  - remove traces from single ticks on literals ([9864747](https://github.com/apache/plc4x/commit/986474753ee80adf89c79cfde67771b725689504))
  - reorder attributes to after params ([bfdb7d7](https://github.com/apache/plc4x/commit/bfdb7d7262bbd686ebe38115e28e6ac370d1aa50))
  - remove obsolete ticks ([24d3472](https://github.com/apache/plc4x/commit/24d3472ba991ff2d9aa4a9b9e68bd492c95d966f))
  - remove traces from single ticks on literals ([8580b51](https://github.com/apache/plc4x/commit/8580b51a01b9968d998d330e747abdb23eb8e304))
  - reorder attributes to after params ([ecb4bb0](https://github.com/apache/plc4x/commit/ecb4bb0c8e7e895124323ce9cf6dc5bd8e4703f6))
  - remove obsolete ticks ([2448eda](https://github.com/apache/plc4x/commit/2448edab8ae7fec526924ca8d9cda1f97a3e4b29))
  - virtual fields should now work much better ([8cc4e45](https://github.com/apache/plc4x/commit/8cc4e450ef186c7557bfcb991f24dc9f4019cf61))

- plc4j/codgen:
  - small code cleanup ([7acc53b](https://github.com/apache/plc4x/commit/7acc53b5a2c01932ae15f0440e90da92d26c1f16))
  - migrate simple type to new code generation ([1fae1f0](https://github.com/apache/plc4x/commit/1fae1f0dc3dec2ac631d3c84a044badf4d264564))
  - add additional optional factory method for expressionless optionals ([652af39](https://github.com/apache/plc4x/commit/652af396a4788aa2f747d130568b5ca7d520bc68))
  - small code cleanup ([75aa3cb](https://github.com/apache/plc4x/commit/75aa3cbca8c8ced7e2e81b96d0a738db08078ab8))
  - migrate simple type to new code generation ([8849b30](https://github.com/apache/plc4x/commit/8849b304bf9fb46c20cbb6d20c95f70e9dd55592))
  - add additional optional factory method for expressionless optionals ([b71830d](https://github.com/apache/plc4x/commit/b71830d8cf0233a9c21c7c24b5a7ab896ca49933))

- pl4go/codegen:
  - added more tracer-calls to codegen ([9d5ade8](https://github.com/apache/plc4x/commit/9d5ade8c617bf68917146afe18dc13d19c08aee1))

- pl4x:
  - Added actualTagNumber to bacnet for easier access to the tagnumber ([d88ea3b](https://github.com/apache/plc4x/commit/d88ea3bf4b106e7a9812296ecbf20c1ff03ecff6))
  - Reworked a bit of bacnet integration ([6146cd1](https://github.com/apache/plc4x/commit/6146cd1dff8f8334c77ef3f10575463951fce5d5))

- pl4go:
  - Added initial pcap transport ([580454f](https://github.com/apache/plc4x/commit/580454ffe9f5eb4ef9841cff236dabeb0d397e52))

### Bug Fixes

- plc4go/cbus:
  - fixed npe while rendering fields ([b318cbb](https://github.com/apache/plc4x/commit/b318cbb9b01f0963baf93ee6f6c1faa9e6ea27c8))
  - don't return connection instance on connection error ([b6d94c7](https://github.com/apache/plc4x/commit/b6d94c73c114362dc1b7c732a315f71d922bdf44))
  - fixed field addresses ([f409645](https://github.com/apache/plc4x/commit/f409645fefacd7e25f5f3a71e6e2cff369085c56))
  - fix error message reporting the unsuccessful cast as nil ([084073b](https://github.com/apache/plc4x/commit/084073bb4a9e703bd422315fa2bf0eb239016818))
  - try to catch installation mmi with offset 0 ([de97996](https://github.com/apache/plc4x/commit/de979966092ac2a82eb7e6b29ea3ad52b9d69c04))
  - fix wrongfully reporting of sal to mmi subscribers ([ff4c27c](https://github.com/apache/plc4x/commit/ff4c27c34749f31484b3f22a2dc54458fb39dc09))
  - remove debug statements ([2d9347d](https://github.com/apache/plc4x/commit/2d9347dd6fb028747706f7545e5de83c10d81be9))
  - fix build ([cbb656e](https://github.com/apache/plc4x/commit/cbb656e4e2fd3df89db69d9a2e997a9304842cfe))
  - change browser to not brute force all unit addresses rather use the installation mmi ([4359797](https://github.com/apache/plc4x/commit/4359797154f497aa0f63b18e84961ac73bcfa868))
  - fix filtering ([8334e20](https://github.com/apache/plc4x/commit/8334e20f47c1508abf5d177898aa881bf6a0f416))
  - fixed cal pattern ([b95a929](https://github.com/apache/plc4x/commit/b95a929de4a7849b455ed74d75479f169405e977))
  - fixed transaction not ending on reads ([2deddcd](https://github.com/apache/plc4x/commit/2deddcd57ba4489f862335f51d59c1d486c9bd5e))
  - fixed empty responses on read ([c38f164](https://github.com/apache/plc4x/commit/c38f164aca14e0bb96647b55a83d8d857156dcad))
  - don't misreport fields which are meant for sal or mmi ([23df43e](https://github.com/apache/plc4x/commit/23df43ea1e02503d1a9f2b2c68d7efc2de995f26))
  - avoid channel leak by adding wait groups ([ce2db6b](https://github.com/apache/plc4x/commit/ce2db6b314c397ffa092aa4046941a6813d3a798))
  - fixed go routine leak ([a641257](https://github.com/apache/plc4x/commit/a6412578904c05634bc6124aa7b3b1c1811e50a6))
  - removed connection closing on timeout ([438cd42](https://github.com/apache/plc4x/commit/438cd428b8f392300695f5036a81d0acaabddaff))
  - fixed issue when connection doesn't work on a power up notification ([e514fe1](https://github.com/apache/plc4x/commit/e514fe10e5fb3cda2a32587b90ceaf5221f02436))
  - removed wrong mapping ([4925d81](https://github.com/apache/plc4x/commit/4925d81a992fd684fd0f5abe645e7f0712ecd4db))
  - fix address string of cal identify field ([16d5348](https://github.com/apache/plc4x/commit/16d5348a4a4331232be12ca30de3c6797b0d46d6))
  - fix address string of unit info field ([2a5b518](https://github.com/apache/plc4x/commit/2a5b518bae07c1a96341316d28c74becc9e6b0c2))
  - fixed write ([191162d](https://github.com/apache/plc4x/commit/191162deb76d33f74cf26e5c27f69676e520a1f6))
  - fixed unprintable event ([46e4a9a](https://github.com/apache/plc4x/commit/46e4a9a004727eaaef077840bfee00a7be03ed62))
  - fixed reader using the wrong options ([1dcab3d](https://github.com/apache/plc4x/commit/1dcab3d7acc92bcb8c988c37f99bd309183b7770))
  - decreased log level of static helper ([7d98c6d](https://github.com/apache/plc4x/commit/7d98c6dc4b1e1a149bb71e3cae84046a831560a5))
  - fixed detection of server error ([5fc5d28](https://github.com/apache/plc4x/commit/5fc5d281beb08f359030737f398fde64bb4cc8ba))
  - fixed handling of server errors ([be71f0c](https://github.com/apache/plc4x/commit/be71f0cd218b64a87cc3d3ef390372e1d6f1a76d))
  - removed spamming log statement ([e812623](https://github.com/apache/plc4x/commit/e81262399c4359385ce239b968671c27938d6d1c))
  - change handling of error responses ([aae82e9](https://github.com/apache/plc4x/commit/aae82e92772b97fae0d7aa48555d1a1a0d8542fe))
  - fixed vet error ([758e078](https://github.com/apache/plc4x/commit/758e078e113e9fb27908848698757ffabc761521))
  - fix handling of "!" errors ([f55a931](https://github.com/apache/plc4x/commit/f55a931e5c0147bad24d13edb7391f9c472f3f92))
  - fixed FieldHandler ([55397a5](https://github.com/apache/plc4x/commit/55397a5ae282813e5eceedf1f7f36bcda7892a23))
  - added a workaround from a nasty bug in using PeekReadableBytes ([548e3b4](https://github.com/apache/plc4x/commit/548e3b43f2d8f06dcbbcffa3d874d0b005d1f3e7))
  - fixed c-bus reporting a request to pci falsely ([8eef885](https://github.com/apache/plc4x/commit/8eef8857fcb80aa78ab455d62fd6fbffea54ca4a))
  - fix application filters (if a1 is set to ff then a2 needs to be set too) ([abef0bf](https://github.com/apache/plc4x/commit/abef0bf2be51fed6c4fa627abf1216d9c3c8e1a4))
  - display MMI's with a default incoming message handler for now ([a5af69e](https://github.com/apache/plc4x/commit/a5af69e83bc8bee346a976af24d37d6b2c1b7394))
  - fixed reading of mmi ([9b7cee0](https://github.com/apache/plc4x/commit/9b7cee0b7f6aaba30b1135065683591002fc7572))
  - fixed broken code ([e0f86d9](https://github.com/apache/plc4x/commit/e0f86d9c697be07905bf8fa509e307072e898a7d))
  - reworked connection to use ack ([e5acb84](https://github.com/apache/plc4x/commit/e5acb84a1bb9d6a5be1a2590352ba83920b81a7b))
  - implement connection setup ([f20adc9](https://github.com/apache/plc4x/commit/f20adc91761cfc6c2acb1f1c678360b2daf37112))
  - T-0 implementation for STATUS, CAL_RECALL, CAL_IDENTIFY, CAL_GETSTATUS ([8cbe0ee](https://github.com/apache/plc4x/commit/8cbe0ee22ed610e66677f5c3f758300038e6d38f))

- plc4go/asciibox:
  - fix npe one empty writer ([8deea7b](https://github.com/apache/plc4x/commit/8deea7be59d54b639d71463ac83313609e7e3951))

- plc-simulator/cbus:
  - fix missing reset echo on smart or connect ([01a2ece](https://github.com/apache/plc4x/commit/01a2ecef479021ae40f90733e77313961a1fb119))
  - fix status level responses with exstat off ([91a9d2c](https://github.com/apache/plc4x/commit/91a9d2caf135afb656b073faad16f89eb537ea6e))
  - partially fix the output of status requests ([88b1b3f](https://github.com/apache/plc4x/commit/88b1b3fe5b31b15683613a4f265a817e411aff80))
  - when exstat is switch restart monitors ([d42ab0d](https://github.com/apache/plc4x/commit/d42ab0d497f4fc903ff9f3170835897fce200147))
  - fix mmi responses ([e1bc7e8](https://github.com/apache/plc4x/commit/e1bc7e872b46f2e11f15c334b462ec901b77b331))
  - fix plc simulator sending out wrong installation mmis ([5e4c77a](https://github.com/apache/plc4x/commit/5e4c77a8e1535d39749738e898f14301230d58c1))
  - avoid sending out the inner message ([b7db778](https://github.com/apache/plc4x/commit/b7db778f4506c34cdd3639d9f2c0d1251016fef6))
  - fixed broken outputs of text ([cfecae2](https://github.com/apache/plc4x/commit/cfecae2dcb40f20158ce436429f53a04ef594ed7))
  - fixed some wrong returns ([ff9fe8d](https://github.com/apache/plc4x/commit/ff9fe8d5fe964cfa17e75f21b9c7530f4a266a40))
  - fixed simulator returning wrong data ([fa96517](https://github.com/apache/plc4x/commit/fa9651754f31d112d005e029ba5859e72f1d2be3))
  - simulator should now use the right types for responses ([2aa5fb0](https://github.com/apache/plc4x/commit/2aa5fb0d1c5f6719c693a8c45cc50c0ab6fdd9ea))
  - fixed NPE while stopping monitor ([f8935f0](https://github.com/apache/plc4x/commit/f8935f069a7e72044bddebcba3b8577d71e30da9))
  - fixed simulator returning wrong status ([0056237](https://github.com/apache/plc4x/commit/0056237aee193bed70aefb6906bf3a103f7ff008))

- plc4j/modbus:
  - PLC4X-354 - ArrayIndexOutOfBoundsException when reading a single bit via Modbus ([c63919b](https://github.com/apache/plc4x/commit/c63919bf75cb93792c2e8040a1f3e66c677c46d9))

- plc4go/connection-cache:
  - guard against returning broken connections ([fa18004](https://github.com/apache/plc4x/commit/fa18004c81d22a08b725c6974c2310072a43b2f3))
  - fix panic when a initialized connection is returned ([8b13c4a](https://github.com/apache/plc4x/commit/8b13c4ab1af785e2e493c881133063bf0ebcf28b))
  - fix NPE on connection close ([cc322f0](https://github.com/apache/plc4x/commit/cc322f09f4295a7c5020e3fcc183c77986b568aa))
  - drain lease chan on timeout ([e6609cc](https://github.com/apache/plc4x/commit/e6609cc2424c1a1ff9fec8e5179e0ae81d3d86fb))

- plc-simulator:
  - added -public option to listen on all addresses ([3e82df2](https://github.com/apache/plc4x/commit/3e82df23c5c2591717903d8f4e9a98d6d9837ffa))
  - fixed issue where the mmi monitor would reset the srchk option ([43f4b87](https://github.com/apache/plc4x/commit/43f4b87f182f31e008755b47e778ff5fc520e3d7))
  - fixed some issues in the plc-simulator ([cbf8325](https://github.com/apache/plc4x/commit/cbf8325494890a0ee7c65652f68e3812357ee127))

- cbus:
  - correct wrong name of HVAC Actuator ([9d0ba2b](https://github.com/apache/plc4x/commit/9d0ba2b4a029093a6ffc9b6216410aaeeeae9466))
  - fixed typo in summary command ([ed5297d](https://github.com/apache/plc4x/commit/ed5297d1bbba8bd280c6ac06d80779b86e78c195))
  - fixed typo in network terminal levels ([385f6e6](https://github.com/apache/plc4x/commit/385f6e6c0223960f41156927e7b13cf12fd997e1))
  - fixed network voltage using the wrong bit size ([21b3d63](https://github.com/apache/plc4x/commit/21b3d63e6a5c19ef66afd547fd000a2631b2ae62))
  - fixed handling of error responses ([6b8b68e](https://github.com/apache/plc4x/commit/6b8b68ed0e6b6cdb348d1e86b59ea47ca40d7f09))
  - implementd InterfaceRequirementsTest and fixed bugs ([4e43951](https://github.com/apache/plc4x/commit/4e439513366a9a75139dd6d6593e176bce8038d5))
  - repaired some byte numberings ([274449d](https://github.com/apache/plc4x/commit/274449d66b64036e73aabe17da2026f7841e16fa))
  - fixed identify output unit sumary ([e4f60a3](https://github.com/apache/plc4x/commit/e4f60a3f1f81cab3fb52e2e69a87d8cfac3c6afa))
  - greatly simplified parsing by removing duplicate definition ([2a132d6](https://github.com/apache/plc4x/commit/2a132d6fce140bf0275673212dcf5488c0c660ab))
  - fixed ParameterValues respecting the additional data ([86d2ee3](https://github.com/apache/plc4x/commit/86d2ee35f75215d966ecb6e51e955101ae1cab4f))
  - fixed setting of c-bus options ([324868e](https://github.com/apache/plc4x/commit/324868edc250b5aa288782dde1a9e61773b339f6))
  - fixed power up notification ([42137e4](https://github.com/apache/plc4x/commit/42137e444d8cb69f5734c0bca154869feb7947e9))
  - fixed detection of monitored sal ([a82e700](https://github.com/apache/plc4x/commit/a82e7006ab3ecb3cce7e7edf02bbd9595d576067))
  - fixed level reporting for ExtendedFormatStatusReply ([0b713c6](https://github.com/apache/plc4x/commit/0b713c67bbb6c1fc86a01cd18f11a71a20536400))
  - fixed handling of extended status replies ([52f072c](https://github.com/apache/plc4x/commit/52f072cd93a4edf6e8275bc29bee19c7a29a8af2))
  - fixed several smaller issues ([6dfa26e](https://github.com/apache/plc4x/commit/6dfa26e7e1f4bc544bc399cc2e3176574d5f02fa))
  - fixed length calculations ([30aa269](https://github.com/apache/plc4x/commit/30aa2699c9167e5f4cf9202d7636c37b410838d6))
  - fixed naming error ([312705c](https://github.com/apache/plc4x/commit/312705c0d36951795d5f48ec1b1541f39c134ca4))
  - fixed several copy paste errors ([d2384a0](https://github.com/apache/plc4x/commit/d2384a05eabfd62983f6caac6f215194d84b76bc))
  - fixed auxiliary levels on temperature and humidity ([1041868](https://github.com/apache/plc4x/commit/1041868caf1357003b2a163519a06a3df5b69e79))
  - fix reset and write command ([de934ec](https://github.com/apache/plc4x/commit/de934ecc67f9cbb718109c01d7ac3dd700fb8c50))
  - fixed issue with GAVValuesCurrent, GAVValuesStored and GAVPhysicalAddress ([e90fdc6](https://github.com/apache/plc4x/commit/e90fdc65b4741423b84bd3cee10ef8e569761a62))
  - fixed issue when using a ReplyNetwork ([4b921f0](https://github.com/apache/plc4x/commit/4b921f0e1859c11622f69120a0352257618a0235))
  - only use request context once on cal data ([d317961](https://github.com/apache/plc4x/commit/d3179612d683b24106ea07660cd56fa13bbbb903))
  - fixed vstrings ([d607519](https://github.com/apache/plc4x/commit/d607519675397a97e1d68747ce74eeb2514946e9))
  - small fixes regarding loading of dynamic icons ([ac955c6](https://github.com/apache/plc4x/commit/ac955c617281ad68e44a76c832c0ae7216bbd21e))
  - relaxed the options parsing for labels ([84c732d](https://github.com/apache/plc4x/commit/84c732de6190c51a7234c995db3422d7d960306c))
  - fixed smaller issues related to extended status ([70ef662](https://github.com/apache/plc4x/commit/70ef662f19758a2508ab8e373da5be75ba1840b4))
  - switch SetParameter from const to reserved till it is clear why the delimiter has strange values sometimes ([14dcb25](https://github.com/apache/plc4x/commit/14dcb25861bbd096317bf095203d2ea9c3e5de9a))
  - fixed CALDataReply using the wrong datatype ([505eb67](https://github.com/apache/plc4x/commit/505eb67a735c93aa3ce49576665f33191ed978e8))
  - fixed response byte calculations ([376f950](https://github.com/apache/plc4x/commit/376f950e35eaf17d5fd8e96d0b161f1a7769a4c5))

- plc4go:
  - use upstream version of tview again ([454b053](https://github.com/apache/plc4x/commit/454b053e620f568fc301bc991d9e5a8a5b65d235))
  - fixed code smells ([0ff69c1](https://github.com/apache/plc4x/commit/0ff69c13c1e1446081e2135991a37d520ba4caf4))
  - fix plc4x_common not reacting to logger changes ([9961a04](https://github.com/apache/plc4x/commit/9961a0469c1252d4863cf5dd84cdbb8d5c060a11))
  - fixed PlcSubscriptionHandle registering the wrong handle ([e6a7255](https://github.com/apache/plc4x/commit/e6a7255b0a97c685e8271e74744797e5441b55d5))
  - fixed subscriber trying to add to a nil map part 2 ([4395053](https://github.com/apache/plc4x/commit/43950537227f6bc81196bb7608da6b471a8544a2))
  - fixed subscriber trying to add to a nil map ([1855082](https://github.com/apache/plc4x/commit/185508292c17666df6fe55e65dcda5889e5579f7))
  - fixed NPE when trying to print the SubscriptionResponse ([84fb3f2](https://github.com/apache/plc4x/commit/84fb3f27b7a68c7494aa5b3a51695dc2f682a1fe))
  - fixed some timer leaks ([4d22dd1](https://github.com/apache/plc4x/commit/4d22dd159756d97a75592b30adca8f63dff84083))
  - fixed issue with broken discovery ([514ae51](https://github.com/apache/plc4x/commit/514ae51feb0f8cd050ebfa2129f9172bcb66973d))
  - fixed close result of empty connection cache ([9f941a4](https://github.com/apache/plc4x/commit/9f941a4e55db66193ff91c34a56cdeaea802d3e8))
  - removed receiving trace calls as they are covered by the workerlog (DefaultCodec.go:273) ([b4c7fc5](https://github.com/apache/plc4x/commit/b4c7fc59cee2e76a40fa185cd644d92ec1cce620))
  - fixed browse query builder ([c142f78](https://github.com/apache/plc4x/commit/c142f78f8c18965c6582664c87e82ab8617fb539))
  - fix go install by changing to fork ([b67679f](https://github.com/apache/plc4x/commit/b67679f0b0f4e7b187131e590d70d79acd9de803))
  - fixed compile issue ([a57705d](https://github.com/apache/plc4x/commit/a57705d1513c07a296e486ed5b17b66f5e062094))
  - fixed wrong import ([5393a50](https://github.com/apache/plc4x/commit/5393a502afb71d1bfb96bca70d45eada5e8c5dbf))
  - fix linting errors ([8035c45](https://github.com/apache/plc4x/commit/8035c45ceccc39309e350a4bbc97f30b0bc1f649))
  - fixed some issues regarding custom message handlers ([7341b61](https://github.com/apache/plc4x/commit/7341b6122cf84612fd4a9d0a34c97cf3114eaea0))
  - fix so_reuse ([fd0eddd](https://github.com/apache/plc4x/commit/fd0eddda14eae0df3f75d061e37b6a1b7428a71a))
  - Ensured using "deploy" with the plc4go module doesn't break the build. ([3719894](https://github.com/apache/plc4x/commit/3719894065f5bcaac047947c870c57f6016723d3))
  - Got the build working with the latest go version by switching from "go get" to "go install" for installing tools. ([b4d01a7](https://github.com/apache/plc4x/commit/b4d01a79fe5d67716a66cab541929693b3f1b56e))
  - Got the go build working again after the plugin update. ([c238fd2](https://github.com/apache/plc4x/commit/c238fd266bb19503f653385eb0d92eb0a8f63b6e))
  - Got the go build working again after the plugin update. ([58210de](https://github.com/apache/plc4x/commit/58210ded43524d723e60c5a1eb6f83b29819739d))
  - fix import ([658162d](https://github.com/apache/plc4x/commit/658162d421c6b5ad799f403929ccb0b71598c748))
  - renamed plc4go/pkg/plc4go to plc4go/pkg/api ([d65201d](https://github.com/apache/plc4x/commit/d65201ddea035ff0a6d2c4e577632702c0afab82))
  - fix the error of value out of range when converting an in-range number string to IEC61131 unsigned integer or WORD-related type. (#363) ([773aba6](https://github.com/apache/plc4x/commit/773aba60bbc9a5a53400fce06cdf61ec38cd1158)) ([#363](https://github.com/apache/plc4x/pull/363))
  - adjusted xml writing ([5a05613](https://github.com/apache/plc4x/commit/5a05613d9c1a89ab9f63969a5dc7aacc7ac23644))
  - switch string reading ([e9d1dcf](https://github.com/apache/plc4x/commit/e9d1dcfb32a6e4f2e31abb1b8dc1865b379c95c0))
  - updated PlcValueAdapter to panic instead of returning 0, updated PlcBYTE to allow returning of uint 8 ([cfdf58a](https://github.com/apache/plc4x/commit/cfdf58a3087141dbfdcc079e824409ab17b1f558))
  - disable blocking test again ([d10d4f8](https://github.com/apache/plc4x/commit/d10d4f8927be5a6ba369919d4f0385fc40dc2c05))
  - fixed leaking spi interfaces in driverManager.go ([d1895d6](https://github.com/apache/plc4x/commit/d1895d660e69b28ba8b77378973eb55cd3e365a0))
  - fix segfaults on macos ([a97ecaf](https://github.com/apache/plc4x/commit/a97ecafbbefd0a4bd81e27836800b2497fcfabcc))
  - fixed issue where ordering was changed after latest refactor ([e0819a8](https://github.com/apache/plc4x/commit/e0819a89709a48943ec0a4324d156bcc3744a16e))
  - fixed driverManager.go omitting the transport path ([abcf1a5](https://github.com/apache/plc4x/commit/abcf1a5b975594c694feeda48bba0db6ff5fb75f))

- plc4j/ui:
  - Tried to get the build working on windows ([ae0dccf](https://github.com/apache/plc4x/commit/ae0dccf8436f4e9438fc849cb2dd0070a0c27e3d))
  - Tried to get the build working on windows ([56d0b88](https://github.com/apache/plc4x/commit/56d0b887703fc1404865f61254aefadadf4fea8a))

- plc4j:
  - remove e.printStackTrace() calls ([e898b7b](https://github.com/apache/plc4x/commit/e898b7b9542520f7f03908b327794c641657187d))
  - Changed the SPI to pass along PlcAuthentication information to the ProtocolLogic ([1045cf0](https://github.com/apache/plc4x/commit/1045cf01a525acddf118ecfb1916df3f84921853))
  - Added some code to avoid NPEs in the OPC-UA Testsuite ([5b70350](https://github.com/apache/plc4x/commit/5b70350fc682bd2fef57536d52354dae9db94755))
  - Made the ParserSerializerTestsuiteGeneratorSpec.groovy check if it uses the Apple pre 1.10.0 libpcap version and skip the test, if it does ([bb923af](https://github.com/apache/plc4x/commit/bb923aff2b1f5d95315a6f19c268b395c975cbe6))
  - Fixed broken code for detecting parllels ([0b2693e](https://github.com/apache/plc4x/commit/0b2693e35d9b9cd90d47e2b9a17a8e6d4b44f780))
  - reverted kafka connect version ([88ec8c4](https://github.com/apache/plc4x/commit/88ec8c489bc5b0d417d6676bddf085c169e18958))
  - Fix bug with select returning prematurely when device not sending any data (#386) ([d8fecd3](https://github.com/apache/plc4x/commit/d8fecd3ac5cc9e38764b2b4046472179d9a2259d)) ([#386](https://github.com/apache/plc4x/pull/386))
  - Set the class loader for for the driver to be the same as the PlcDriverManager ([10e942b](https://github.com/apache/plc4x/commit/10e942b74f4a4ebb8924c63c554400871918ad6b))
  - OPCUA more fixes for optional fields ([b3dcb67](https://github.com/apache/plc4x/commit/b3dcb67ce973490f794944d293906ac67b827121))
  - opcua Include Encoding Mask in NULL Extension Object ([4e14bc0](https://github.com/apache/plc4x/commit/4e14bc04a8b4fcbc3718ad81a2a6962f9bf8dc1e))
  - Kafka Connector fail to find transport fix ([63db99c](https://github.com/apache/plc4x/commit/63db99c6ce759b2933bfc813712096186bebb627))
  - Updated scraper pom ([e1629b9](https://github.com/apache/plc4x/commit/e1629b9352e683b978f99cf210af86f86c4fd0bd))
  - Replaced the PooledPlcDriverManger within the Scraper ([5531ca8](https://github.com/apache/plc4x/commit/5531ca85f829cf264e100c5ce95c6a0062d6034d))
  - Bumped minimum java compiler version to 11 ([a19ce13](https://github.com/apache/plc4x/commit/a19ce13f30fcb87f7621bd6c4e9d44f9c3983178))
  - Made the bacnet RandomPackagesTest.java execute on Windows and updated the documentation on setting up libpcap on Windows. ([0b7605e](https://github.com/apache/plc4x/commit/0b7605ef2333781e7d1c38ef94eb7b22834f3491))
  - Made the bacnet RandomPackagesTest.java execute on Windows and updated the documentation on setting up libpcap on Windows. ([4fc0433](https://github.com/apache/plc4x/commit/4fc04331abffce1d0c9f0fd6bb45c441aa9bd27a))
  - Hopefully fixed the code for detecting the libpcap version. ([e44a84a](https://github.com/apache/plc4x/commit/e44a84a38045404b9399716a7517bb4b7d8977ea))
  - Adjust OSGi framework versions PLC4X-336. ([ec0694a](https://github.com/apache/plc4x/commit/ec0694a700cb72c9a939834142c3488e711c3b0b))
  - Adjust endianess in ethernetip mspec. ([3453b60](https://github.com/apache/plc4x/commit/3453b60c4ac8b6dafa899e62ffff48c26f36ef2f))
  - Made the driver-testsuites run with the version 3.0.0-M5 of the failsafe plugin ([b0b8605](https://github.com/apache/plc4x/commit/b0b8605caea504fe236465f2e81c546861c3b7a5))
  - Disabled the ParserSerializerTestsuiteGeneratorSpec test as it wasn't runable on Windows ([0fd5d63](https://github.com/apache/plc4x/commit/0fd5d63d2dcc54855a1ff7927ff5eb563972c15d))
  - Made the integration tests work with failsafe 3.0.0-M5 ([1cd02bf](https://github.com/apache/plc4x/commit/1cd02bf95be1932b7ff2f98527b1379388be98dc))
  - removed files in the "attic" of the raw-sockets module ([37ebe6a](https://github.com/apache/plc4x/commit/37ebe6af6dd89237d50f4c21b77560b402902db4))

- plc4j/examples:
  - fix build part 2 ([332f506](https://github.com/apache/plc4x/commit/332f5065465ca8de8f3e92ca2daee13a860edf3d))
  - fix build ([326739e](https://github.com/apache/plc4x/commit/326739e35e9423347c154aa70fa74f80556888e8))
  - PLC4X-349 mqtt example config contains wrong plc protocol ([2e8c6eb](https://github.com/apache/plc4x/commit/2e8c6eb201adae23a1159eea198e4cadca9a7cb1))
  - fixed wrong path of example parent pom ([f6d146a](https://github.com/apache/plc4x/commit/f6d146afec20fead358f07d02075d6a94dbb78cd))

- plc4j/ads:
  - Refactored the ADS driver to support reading of complex types. ([9f3cc6d](https://github.com/apache/plc4x/commit/9f3cc6d0d4f5686e51db1780bf2df12b0ce351c5))
  - Updated the ADS connect logic to be a bit more asynchronous (At least on an API level) ([6fb1dea](https://github.com/apache/plc4x/commit/6fb1dea1ff8e6776c974794de02a1e640b676be3))
  - Implemented the automatic adding of AMS routes in Beckhoff PLCs ([6e6507b](https://github.com/apache/plc4x/commit/6e6507ba52614741539389619705c2fba545f3ba))

- plc4j/spi:
  - Added the missing "getObject" method to the PlcStruct ([386189f](https://github.com/apache/plc4x/commit/386189f4ddcedddb6c98aa105ec2927ff30289c5))
  - Open too many files after complete (#351) ([4b151e3](https://github.com/apache/plc4x/commit/4b151e30a6ed704399b68777af0a3b02247fb734)) ([#351](https://github.com/apache/plc4x/pull/351))
  - Made it simpler to detect configuration problems ([551868f](https://github.com/apache/plc4x/commit/551868fab9dbe6577414e2a21c74019ae88868cf))
  - fix issue when writing virtual float fields ([6e56e49](https://github.com/apache/plc4x/commit/6e56e49def592477f31e5b79abcac84564eb109c))
  - fixed asciibox breaking on CR/LF ([6fe7535](https://github.com/apache/plc4x/commit/6fe7535c1ece36fffe59d3c38ae6a42973e629cd))
  - fixed write buffer json based ([ada0abd](https://github.com/apache/plc4x/commit/ada0abdbc632bf8f83263a113db4b9e7b43d0771))
  - implement unknown field ([e882921](https://github.com/apache/plc4x/commit/e88292168b43017735d92f3fc0f324cda50395da))
  - implement reserved field ([1f55c4e](https://github.com/apache/plc4x/commit/1f55c4eac07532f509b51f3baf069724f5b6451c))
  - implement const field ([4c1a1eb](https://github.com/apache/plc4x/commit/4c1a1eb529fc6b3a7f3aa7699968b9158773b4c0))
  - implement assert field ([dd3cad5](https://github.com/apache/plc4x/commit/dd3cad5b0b723bf0e38a92c5d83ac2f40734ffaf))
  - change string type representation to optional ([c1912e0](https://github.com/apache/plc4x/commit/c1912e03ccfe1c6a57a0bb2f09db933c6c4fa4eb))
  - fixed write buffer json based ([7be4b36](https://github.com/apache/plc4x/commit/7be4b36916e16ec3e3fcd819f1e87e76f552c8b2))
  - implement unknown field ([d109d44](https://github.com/apache/plc4x/commit/d109d44fabece6e7634d4d7512468379bb9b57d8))
  - implement reserved field ([9d0a0e4](https://github.com/apache/plc4x/commit/9d0a0e4a6dce3053fa53577c7980b925967855cb))
  - implement const field ([0d697e3](https://github.com/apache/plc4x/commit/0d697e3cc63028a8bccde5a4dd2b7136e61b8445))
  - implement assert field ([e99fbc2](https://github.com/apache/plc4x/commit/e99fbc21c9a74d76ac3d0fe65d267fb71b3c0580))
  - change string type representation to optional ([58eeb0c](https://github.com/apache/plc4x/commit/58eeb0cfb029a92993da01c3dadf9047115e6412))

- plc4j/s7:
  - Fixed the PLC4X272 test and made it run automatically. ([865d5ef](https://github.com/apache/plc4x/commit/865d5effcda2ce3f9f6e505213e2150dfb6ad6af))
  - Removed some silly code I addded to provoke errors. ([4a771d3](https://github.com/apache/plc4x/commit/4a771d31332f780dda7888dc6e64f25075c23cac))
  - Fixed some issues in the freshly merged S7 PR ([6ab7c40](https://github.com/apache/plc4x/commit/6ab7c401d3d539fcad8b93754781c2e6a14e87a1))
  - Fixed some issues in the freshly merged S7 PR ([00a857a](https://github.com/apache/plc4x/commit/00a857af424cc4cabc1f6afeb18d5310dd200669))
  - communication with LOGO 0AB7 and ISOonTCP tsap configuration (#308) ([ece4af4](https://github.com/apache/plc4x/commit/ece4af4d260cbf4aa4d0884c40e87aa985ffa6fa)) ([#308](https://github.com/apache/plc4x/pull/308))

- protocols/bacnet:
  - Updated the URL for fetching the vendor ids ([16c7d37](https://github.com/apache/plc4x/commit/16c7d377a59e45ea933a48cd09a148c237f2c4f7))

- build:
  - Disabled the dependency-check for now completely as there was a cascade of reporting false positives ([5aaf204](https://github.com/apache/plc4x/commit/5aaf204e4499d6b7d5f857a51e714a58e87390a1))
  -  Got the build working on my Mac with M1 chip ([5c0ab1e](https://github.com/apache/plc4x/commit/5c0ab1e3c04ea7c13b86f7deb870c16403e330ee))
  - Tried making the Docker build successfully build plc4py ([639c03a](https://github.com/apache/plc4x/commit/639c03a922b4a4c3a9f091001450c17289fcdd64))
  - Bumped the version of the CVE checker plugin to hopefully fix issues during the build ([e2ca19d](https://github.com/apache/plc4x/commit/e2ca19d96c9d7a0faf4faae959522a5b68e7c336))
  - Enabled the "enable-all-checks" profile on Jenkins ([29f0dda](https://github.com/apache/plc4x/commit/29f0ddae0f0f33bcfd05cddc99ee7911425e5c69))
  - Extended the prerequisite check to check for "poetry" when building with python. ([c1ff8b9](https://github.com/apache/plc4x/commit/c1ff8b942b582184943e1146b13b51107f67a066))
  - Removed my double poetry.lock rat exclusions again ([5b837a1](https://github.com/apache/plc4x/commit/5b837a126bc67369fe66a3964af73781bd77436c))
  - Added poetry.lock to the rat exclusions ([92cfa9e](https://github.com/apache/plc4x/commit/92cfa9e12b554615b86ae9aaa8d833553457f166))
  - Addressed a CVE in jackson as well as bumped the version of Netty to the latest version ([49c9c48](https://github.com/apache/plc4x/commit/49c9c48662ed66f322835194f645af6a6b9f1bee))
  - Actually ensured the cve scan and the language tests are executed on jenkins ([3a11f39](https://github.com/apache/plc4x/commit/3a11f3907a9b9c329834da0b967383fa9cb3a7eb))
  - Moved the cve-scan to the jenkins profile ([a5e3c90](https://github.com/apache/plc4x/commit/a5e3c90344bd2580dd335de17a2b5826d3517506))
  - Streamlined the artifact names of the plc4j examples ([f6b466a](https://github.com/apache/plc4x/commit/f6b466acbd3d26e1e448ee8c325df16a2f9002f4))
  - Set the compiler code-level to 1.8 in the entire plc4j part. ([a263357](https://github.com/apache/plc4x/commit/a263357ffde2f2a46f8d669d0a338516909cb4f8))
  - Updated the maven-compiler-version to 3.10.0 (which was released today) as with the last update my build wasn't able to get 3.9.0 ([e908b7c](https://github.com/apache/plc4x/commit/e908b7c12a872146ea1360f88de25ecd21c984de))

- plc4j/opcua:
  - Made maven skip the dependency checker in the opc-ua module as it was ignoring all of my attempts to manage the CVEs ([0ddad63](https://github.com/apache/plc4x/commit/0ddad639236b9ed522d41ccc62b183c84e023110))
  - More rebust MonitoredItem handling within subscription handle. ([066c395](https://github.com/apache/plc4x/commit/066c3954afa0ae3b534432802a61c3f9d448d958))

- plc4go/spi:
  - fix AwaitCompletion of RequestTransactionManager ([10962f1](https://github.com/apache/plc4x/commit/10962f133609d4ed31b8efecd39f667d749612b1))
  - fixed transaction await never ending ([927baf3](https://github.com/apache/plc4x/commit/927baf3337f3fcf78a2c7acf657ab225a2d58357))
  - fixed minor issue with RequestTransaction ([33fdbda](https://github.com/apache/plc4x/commit/33fdbda7be4b1c8ec3e01d5ee2429bd64d7139e0))
  - disconnect message code on close ([7c814cc](https://github.com/apache/plc4x/commit/7c814cc27391658ea48bee3b663a0525316e7912))
  - fixed net command ip issues ([e28f2de](https://github.com/apache/plc4x/commit/e28f2ded79a8110fef00eb9c3a62c0e122deb929))
  - Fixed issues in the "Serialize" function of PLCValues and implemented the "GetRaw" for each of them. ([d191a8e](https://github.com/apache/plc4x/commit/d191a8e31351ad68840b41fcc342334e4096b2f6))

- spi:
  - fixed transaction await never ending ([de4d01e](https://github.com/apache/plc4x/commit/de4d01ea5aeb35a388e9acd5ef0784951c4ef575))
  - WriteBufferBoxBased should not switch out writer for sub ([35c2e26](https://github.com/apache/plc4x/commit/35c2e26940a126bdc94c7e5f343a7e601fe4443e))
  - remove beanutils ([deb42b7](https://github.com/apache/plc4x/commit/deb42b720c1e181d87bfdad05c451209e32187e5))

- plc4j/utils:
  - Added a check for libpcap and the version to the ArpUtils giving error messages if anything is missing. ([1216bbc](https://github.com/apache/plc4x/commit/1216bbcd6e9aa99a39264404d29aec40fc9bfd70))

- plc4j/opc-ua:
  - Added some exclusions to the depedency-check, as it was reporting two invalid CVEs (They were referring to an OPC-UA rust library) ([07378bd](https://github.com/apache/plc4x/commit/07378bd1f4b213782229e28abab3b57e2c9c17d7))

- asciibox:
  - ported the distortion bugfix from golang to java ([b944ea1](https://github.com/apache/plc4x/commit/b944ea173252897cecf3a8767f7db2f96c14183d))

- asciiBox:
  - fixed long outstanding bug which distorted boxes ([24edb1e](https://github.com/apache/plc4x/commit/24edb1eb4e6d08338058a78a650d7c5a53fa6452))

- plc4xtools:
  - set shutdown flag early to avoid corrupted files ([dde3b01](https://github.com/apache/plc4x/commit/dde3b0173a6491c4a053c81fcf10c043662f4bb6))

- plc4xpcapanalyzer:
  - remove ! in the middle of responses and return it as a "single" response ([9a1fb1c](https://github.com/apache/plc4x/commit/9a1fb1cef604c56ba48c89cab3a59d8ed8de029b))
  - fixed unnecessary replacing of payloads ([4d379b1](https://github.com/apache/plc4x/commit/4d379b194b8321700ea403ef02da4546ae1ab740))

- plc4xbrowser:
  - fixed command escaping when using brackets ([5ce93fb](https://github.com/apache/plc4x/commit/5ce93fbd9dca87b3c17e9eb676471e6d9c1a6aad))
  - fixed small commit accident ([e72a5ff](https://github.com/apache/plc4x/commit/e72a5ff06571b7f5cae82134537a168552443524))
  - replace tview with sruehl fork till https://github.com/rivo/tview/issues/751 is merged ([414e34d](https://github.com/apache/plc4x/commit/414e34db35dcee9871b8c98357958e5464e4bf0f))
  - don't wrap messages as this destroys layout ([36c2ce0](https://github.com/apache/plc4x/commit/36c2ce0fa08e7ea0be7a18dfc95d99c3fa51f66f))
  - fixed issue where output got overwhelmed by to many message ([1374543](https://github.com/apache/plc4x/commit/13745432b4f0f75ec62505a499f987e53618234c))
  - fixed log setup ([b75b229](https://github.com/apache/plc4x/commit/b75b22991e7508d0fd9511b0121b4f472ec59a0e))

- codegen/plc4go:
  - fix access to params ([4336ebe](https://github.com/apache/plc4x/commit/4336ebe8ef551373ca5dbba0676fa44788258099))
  - fixed issue with unnecessary casting when using bitwise operators ([514e03b](https://github.com/apache/plc4x/commit/514e03b25d75702bfe8a4a57f4bfa1550d741571))
  - guard values returned from manual fields against nil ([905649e](https://github.com/apache/plc4x/commit/905649e4e5d2773e11b239760a9623b9003d49fb))
  - guard values returned from manual fields against nil ([f956e15](https://github.com/apache/plc4x/commit/f956e158528b578373214c80ea2ae27d067b5a0f))
  - fixed issue when using parser argument childs ([c11481f](https://github.com/apache/plc4x/commit/c11481f0dc5746b69a0f60ad038e164e33369f12))
  - fixed issue when using vstring ([50ec883](https://github.com/apache/plc4x/commit/50ec883e5876569fd44534c8484cb607a521caa3))
  - fixed issue where failing enum read resulted in non error ([7b590d1](https://github.com/apache/plc4x/commit/7b590d18c747bb878ad1e768df3b7ce59083522f))
  - fixed build ([9e22a44](https://github.com/apache/plc4x/commit/9e22a4491f357cb1528efacf9a4e9fe9c62a995b))
  - fixed issue when multiple wildcards are used ([e282643](https://github.com/apache/plc4x/commit/e282643e26020e8df3275c2f23310d3a5a3cb05e))
  - enum by name should now tell if the name is unknown ([0f55269](https://github.com/apache/plc4x/commit/0f552695e6f75377e6c853d5b17800f0a628c5cb))
  - fixed issue were cast was done before error check on typeswitch ([cb93864](https://github.com/apache/plc4x/commit/cb93864b394346676296c66e60de75469c2a5140))
  - fixed issue with nil type cast ([150301c](https://github.com/apache/plc4x/commit/150301cfff02d85f63b292fb78fe2a24aa5c3ef4))
  - remove panic-ing serialize method from parent ([c99352a](https://github.com/apache/plc4x/commit/c99352a043e550fddb573be785ca1bcfff72e546))
  - fixed issue with complex type and length arrays ([6b7628f](https://github.com/apache/plc4x/commit/6b7628f91d3efe5146138dd377db03efcf7eaf15))
  - fix broken code comments ([aa0e848](https://github.com/apache/plc4x/commit/aa0e848c42e708dc2b7a8b9913e02208bf6d5847))
  - Fixed a problem using unary expressions in plc4go static calls ([afbdbc2](https://github.com/apache/plc4x/commit/afbdbc2d0371441faae9bd96d0f09740bc88a90c))
  - fix issue with broken optional fields ([7df1184](https://github.com/apache/plc4x/commit/7df1184c7073abbe7adc3ce30ac24d99e244ff5c))
  - fix issue with broken optional fields ([a3755c5](https://github.com/apache/plc4x/commit/a3755c527a0edfdeb5509604e81254ac9060748a))

- docker:
  - Fixed the Docker build ([fd1a0ef](https://github.com/apache/plc4x/commit/fd1a0efc9a63500a87c8b43539253e3b5d1057dc))

- general:
  - fix build ([55807bb](https://github.com/apache/plc4x/commit/55807bb1ffcd74d5768e0b1dd913a1483ed3df85))
  - fix issues ([43227e0](https://github.com/apache/plc4x/commit/43227e086d74a991b1cb6dd98db84fdccdb5d209))
  - fix build by switching to right test version ([93a1cff](https://github.com/apache/plc4x/commit/93a1cff1b3ea6367652a2eedc04ecf922434a78c))
  - fix build by switching to right test version ([3e502b4](https://github.com/apache/plc4x/commit/3e502b4c114b16fef608cfa3a44f51cd336fed14))
  - Fix casting to byte for integer greater than 127 (byte is signed number and for string greater than 127 characters it result to NegativeIndexException) (#339) ([6f5399e](https://github.com/apache/plc4x/commit/6f5399e4498fd016e9622bb90f69e0f0b906673b)) ([#339](https://github.com/apache/plc4x/pull/339))
  - Don't make errors in the analysis fail the build, as these seem to occur randomly. ([3ea3544](https://github.com/apache/plc4x/commit/3ea35442ac000b1a6967b2fb38884b6b4bffa513))
  - Enabled the owasp dependency analyzer and fixed all reported CVE problems, that are not minor. ([4a1e95c](https://github.com/apache/plc4x/commit/4a1e95c62f6091c58125b2f0559d0ce0af1e247e))
  - Bumped most dependencies and cleaned up dependency management throughout the project. ([f2c0578](https://github.com/apache/plc4x/commit/f2c0578a9244ae87c55c5eef441eff7d7bdd6f4e))
  - Updated to the latest version of Apache Calcite ([17ae381](https://github.com/apache/plc4x/commit/17ae38192a150022e9ffadec2bd141d42926f35b))
  - Updated to the latest version of Apache Kafka ([b45ec67](https://github.com/apache/plc4x/commit/b45ec67d1ab2c7faa9bec3e805e92b795b7e53b7))
  - Updated to the latest version of Apache NiFi ([9d0151c](https://github.com/apache/plc4x/commit/9d0151cf73346b6a2a8e55c426f0dd255ad75c57))
  - Fixed a potential exception if the serializer produces more bytes than the testsuite expects. ([13fc3d0](https://github.com/apache/plc4x/commit/13fc3d0d01057a2bd691fd7e6d08ec215650b823))
  - Updated mqtt client version ([08ce994](https://github.com/apache/plc4x/commit/08ce994a910c2a5ac364821d8cfa6ba071aa5fb1))
  - Updated list of knx manufacturers ([332163e](https://github.com/apache/plc4x/commit/332163e5b198e5dc1b0f9deb8b8e093918844ec8))
  - back ported the fix for a little error back to develop ([45913fc](https://github.com/apache/plc4x/commit/45913fc130acff6bc2a51675ecd9de2db329a8fb))
  - bumped the version of log4j ([7efe5f2](https://github.com/apache/plc4x/commit/7efe5f288b3ca703f4744fea6d2010f192f25f1f))
  - Changed the log level for Go Disconnect messages to trace ([9471ec7](https://github.com/apache/plc4x/commit/9471ec7e24de14d0f8153f7779b4cab247e75557))
  - made the worker terminate gracefully. ([6154562](https://github.com/apache/plc4x/commit/61545627ea2ac8253954cf3b0d92cf394c294c02))
  - fixed the worker log. ([e1c73a7](https://github.com/apache/plc4x/commit/e1c73a763049258d240456672854116d3fb4c879))
  - Reduced the log-level of the "keep running" message to trace. ([085d05a](https://github.com/apache/plc4x/commit/085d05aee863cdd0bf957c8911a908f87cfd866b))
  - Made the knx-driver actively close the transport in case of the connection not being successful ([49417d8](https://github.com/apache/plc4x/commit/49417d836e94462d07993c0bf66a8ee4e622b0f9))
  - Fixed a test to work with the new generated code ([0fc7ff2](https://github.com/apache/plc4x/commit/0fc7ff2ef7c23b4a5f8802d9dc698c118b05eda0))
  - bumped the version of log4j to the properly fixed one. ([f00c9a2](https://github.com/apache/plc4x/commit/f00c9a2c75826359365abdc88be0a04189c5f73c))
  - bumped the log4j version to a newer version ([ad68503](https://github.com/apache/plc4x/commit/ad685033de9d9cb3bc4736354c4d3a130807b623))
  - fix (knxnetip): Added the device-infos for replaced versions. ([a1c2cdc](https://github.com/apache/plc4x/commit/a1c2cdc77163168ee017b89f62c5d974297b3661))
  - Enabled the language tests again and fixed all the stuff that needed fixing after that. ([7c7b078](https://github.com/apache/plc4x/commit/7c7b078cde4ffb11e49359c3bf1a3d2f0ea2fc2d))
  - Updated profinet.mspec that correctly decodes the profinet connection response sent back from my device. ([366a30d](https://github.com/apache/plc4x/commit/366a30dafb504fa89dcc6f1c969608ec5143e0d7))
  - Made PLC4Go build again ([a156df6](https://github.com/apache/plc4x/commit/a156df6dd85e8bc40d79982eb66673cac4c64570))
  - Made PLC4C build again ([b17c6b4](https://github.com/apache/plc4x/commit/b17c6b4d460d44b263935e9e0be7801b4d876335))
  - Worked hard on getting the conditions for typeSwitches strongly typed. (WIP ... currently all Java modules work, C and Go still need some polishing) ([b99a88a](https://github.com/apache/plc4x/commit/b99a88aad30a619c48c5fefeebd05ecd28142ccf))
  - Changed typeSwitch to no longer use expressions as discriminators (Now variableLiterals are required) in const fields now no longer expressions can be used and instead only Literals can be used. ([7e289bd](https://github.com/apache/plc4x/commit/7e289bd950a613e3baee825e3092f6c84e3f1909))
  - Changed typeSwitch to no longer use expressions as discriminators (Now variableLiterals are required) in const fields now no longer expressions can be used and instead only Literals can be used. ([a89e1bd](https://github.com/apache/plc4x/commit/a89e1bd2b004c00b9a84a37f5b19e97340ff2c8a))
  - Changed typeSwitch to no longer use expressions as discriminators (Now variableLiterals are required) in const fields now no longer expressions can be used and instead only Literals can be used. ([a888473](https://github.com/apache/plc4x/commit/a888473c59a49abb7d36dc84ab7e62b9779bb946))
  - added support for typed expression handling in Java and added support for automatically adding "L" to uint expressions with 32 or more bits. ([ea0127c](https://github.com/apache/plc4x/commit/ea0127cd80fc6ef8bd38c0cefb7ff27b2023e668))
  - set maven version to the previous lts to give plugins a bit time to catch up ([2e7c543](https://github.com/apache/plc4x/commit/2e7c543ec142740be25b8e661effa976bad7fb01))
  - add float adjustment ([20ea701](https://github.com/apache/plc4x/commit/20ea701aefb55e0ab3a39145d7b6c8bf39909931))
  - Continued implementing the arguments in the parser
fix: Set a byteOrder argument on every root-type in our mspec grammars ([f2c17ff](https://github.com/apache/plc4x/commit/f2c17ff66b77605e8a7aab8e5ca417843b86e3af))
  - Some dusting off of the old simulator ([3e9d11f](https://github.com/apache/plc4x/commit/3e9d11ff3406c1db1d04f2562899b2b5e95b3aba))
  - fix build ([457b2a6](https://github.com/apache/plc4x/commit/457b2a6155103b48c36892a759376a8729624ac8))
  - Continued implementing the arguments in the parser
fix: Set a byteOrder argument on every root-type in our mspec grammars ([611bed9](https://github.com/apache/plc4x/commit/611bed990ce35c10f8e1d21c1bfd209304e368c4))
  - Some dusting off of the old simulator ([9e092a2](https://github.com/apache/plc4x/commit/9e092a217ce8ad0edac7a7c461a95f6bc70a37b3))
  - fix build: replace strange call in c-generator ([f902c74](https://github.com/apache/plc4x/commit/f902c74c8239c4159b75f509bd1444cc3db45a65))
  - fix stackoverflow on type referencing each other ([9470bbb](https://github.com/apache/plc4x/commit/9470bbb3d6b06da665e3db7fe0e25c8f04e516db))
  - set executable bits on sh script ([939a4cd](https://github.com/apache/plc4x/commit/939a4cd2a00d4ac573c4652b22855135e589e5ac))

- plc4go/s7:
  - fixed serialize of PlcStringField ([01e6c6d](https://github.com/apache/plc4x/commit/01e6c6d0fbc4d59bf66a43f7630afb0044f238c8))
  - fixed copy-paste error in s7 config options (#416) ([daa2914](https://github.com/apache/plc4x/commit/daa2914a2d8ecebc0a0cb9dc819b935c6f50710c)) ([#416](https://github.com/apache/plc4x/pull/416))

- plc4go/plc4xpcapanalyzer:
  - cbus-srchk only sets the checksum option towards the device ([b184705](https://github.com/apache/plc4x/commit/b1847051ca51ae14bd0070fedd681503ff9201b3))
  - fixed bug with payload replacement with equal length arrays ([a20c3b2](https://github.com/apache/plc4x/commit/a20c3b22aeb49b96d1d681677b38096678db073b))
  - first version of extractor which can be used to visualize communication (useful for ascii protocol like serial ones) ([67c97cd](https://github.com/apache/plc4x/commit/67c97cd9f9afa901f2a33a9d768a5928053d461c))
  - fixed log option and added parsing for missing types ([f2c5674](https://github.com/apache/plc4x/commit/f2c5674837eeac49529586e4f91d938fe1c90e06))

- c-bus:
  - fixed cal/sal data chaining ([e94abeb](https://github.com/apache/plc4x/commit/e94abebf9c5d44a662c44b9eba45db3818a5d9eb))
  - fixed more parsing issues ([0b8d3cb](https://github.com/apache/plc4x/commit/0b8d3cbad2ba634b7c8576c50a1ed6281a50cb1a))
  - add additional alpha for confirm ([6cba01a](https://github.com/apache/plc4x/commit/6cba01a8aa5daca4d1c5a57f89074bb8dc53e60b))
  - fixed several issues parsing messages ([5e4ac86](https://github.com/apache/plc4x/commit/5e4ac864397905f3486ba84449d2d918b7ddc1ac))
  - fixed that alpha char handling on commands ([3e2cd3c](https://github.com/apache/plc4x/commit/3e2cd3ccdab0aea6309ba67fe6140f5dc4e9bb98))
  - smaller refinements decoding the messages ([62240fb](https://github.com/apache/plc4x/commit/62240fbb2aab895e466e914ffee4d392342bba7b))
  - fixed hex de- and encoding ([512b6b0](https://github.com/apache/plc4x/commit/512b6b05185ee19c2a483ef0af8a3eb4945b31ab))

- codgen/plc4go:
  - fixed wrong variable reference ([d0df22b](https://github.com/apache/plc4x/commit/d0df22b48db9756958c731a4e133f8c6bb80b97f))

- plc4c:
  - plc4x_server build issues (#404) ([f8499cd](https://github.com/apache/plc4x/commit/f8499cde78ae74150c4659c7f37a3d8ded1a9e03)) ([#404](https://github.com/apache/plc4x/pull/404))
  - fix support for aarch64 ([981099d](https://github.com/apache/plc4x/commit/981099d376e9da50102c7e9d1b4a529a0cee01fc))
  - fix issues during codegen ([5f195de](https://github.com/apache/plc4x/commit/5f195de6ffa3b64cb092a33d6d60122545b44bbe))
  - fix ordering of const fields ([6c446a1](https://github.com/apache/plc4x/commit/6c446a1e1f922d045373d2d4b76cb509cf933aac))
  - fix tracer for enums ([c3e13b7](https://github.com/apache/plc4x/commit/c3e13b772fb337159b012d269e34d1a9a6a8f2d3))
  - include time.h so build won't fail on mac ([cec2547](https://github.com/apache/plc4x/commit/cec25479af8a26605386db5fad7343216eb7fcfc))

- protocol/eip:
  - Removed some single-quotes around the CIPStructTypeCode enum name ([9b2ebd0](https://github.com/apache/plc4x/commit/9b2ebd0afcd3c61c91248c04e77ab2723962d217))

- bacnet:
  - fixed cov notification ([9306d25](https://github.com/apache/plc4x/commit/9306d2581c2f199370f44684cf3502bd66145335))
  - ported fix regarding bit wise reading from plc4j ([66bfff0](https://github.com/apache/plc4x/commit/66bfff0bfefde00f2819269638da6f12070bf9db))
  - fixed status of BACnetAccumulatorRecord ([ea0f00d](https://github.com/apache/plc4x/commit/ea0f00dc2e5656bf6ef59ffdba89f5d56311424e))
  - added unmapped enums to static helper ([47c95bb](https://github.com/apache/plc4x/commit/47c95bbff067674b04c56acaf615b682a4ed81e7))
  - added unmapped enums to static helper ([6ffab0e](https://github.com/apache/plc4x/commit/6ffab0e4d6133356aebcc1c04b2a0222d6c3cf8c))
  - fixed several outstanding issues... ([85a9b48](https://github.com/apache/plc4x/commit/85a9b486f053cd7dfe84c1ea02ec9752c06dd842))
  - fixed memory overflow when receiving broken package ([9f524b8](https://github.com/apache/plc4x/commit/9f524b8f96d515fcb02de01a3cc5ec01c8dae7c4))
  - several small fixes regarding length calculations with proprietary values on serialize ([4cb82bd](https://github.com/apache/plc4x/commit/4cb82bd4dd4c72ef3ac4d019213b713e4c19c630))
  - fix apdu unknown not consuming enough bits ([1f7bed5](https://github.com/apache/plc4x/commit/1f7bed58814debcc996ba18ae816b581e236958d))
  - Removed the "()" around the optional \r as this introduced a new capture group and this broke accessing them by index. ([4214495](https://github.com/apache/plc4x/commit/421449519366878e28f294a933afb287a0f90ceb))
  - Adjusted the produceTaggedPrivateEnums.groovy script to also work on windows systems. ([3945014](https://github.com/apache/plc4x/commit/3945014147c16ffb7bcf9d850d5c81622be192b0))
  - Addressed the parser errors as a first step of tracking down the build errors ([1949a8e](https://github.com/apache/plc4x/commit/1949a8e5d4a976c51d2cfe495780773739c4cc77))
  - fix some open issues ([6bf9513](https://github.com/apache/plc4x/commit/6bf951376fd074e20e9d5fff4b366ce9ef7c7df9))
  - fixed production of vendor id ([8d94b4c](https://github.com/apache/plc4x/commit/8d94b4c0b79b85bf4e35deed703c4fc4f9c35047))
  - fixed several issues by adding validation to APDUConfirmedRequest and APDUComplexAck ([1d943f3](https://github.com/apache/plc4x/commit/1d943f36eadd25efc160c0deb1d64af7a5e63b16))
  - fixed several parsing errors ([09192fe](https://github.com/apache/plc4x/commit/09192feb5e83f3f96971e7bf0d39c146d80dea72))
  - fixed several bugs ([4bdd42a](https://github.com/apache/plc4x/commit/4bdd42af1ca7b4977faeb6ba4ab52e47f4c1971c))
  - fix support for segmented requests ([592c76d](https://github.com/apache/plc4x/commit/592c76d0442ba264eb9c93b3b008d34a25cb6071))
  - git commit didn't work properly again... ([4ed746a](https://github.com/apache/plc4x/commit/4ed746abfe6dd47aa685d694f460986fdbf40786))
  - fixed offset for service request ([a2b6332](https://github.com/apache/plc4x/commit/a2b633299abe38d42cff89cfe2f96ef238ae5620))
  - fixed BACnetReadAccessProperty ([97c28ee](https://github.com/apache/plc4x/commit/97c28ee21c2877f4be557cae73c5da12a0a3a5ac))
  - fixed length calculation for APDUComplexAck ([a761949](https://github.com/apache/plc4x/commit/a761949e0d2180956bc539f45a2725ee593f0103))
  - fixed string length calculation for helper ([00417e4](https://github.com/apache/plc4x/commit/00417e43f2def30389d135c8aafd42715881cbea))
  - fixed boolean tag ([368940d](https://github.com/apache/plc4x/commit/368940d20f73d974dfab0663f1c57745e0982c4a))
  - fixed opening/closing tag creation ([aa14db7](https://github.com/apache/plc4x/commit/aa14db76470986d944bd23052c65fbbfe05569b8))
  - fixed opening/closing tag creation ([67dcae9](https://github.com/apache/plc4x/commit/67dcae9b95ffe6b0dd1c7d3be646d61102da7108))
  - only supply proprietary value to objects if they are indeed proprietary ([1947a2a](https://github.com/apache/plc4x/commit/1947a2ab87ab3c079751d9a6257314f5be1187f3))
  - fixed COV Notifications ([4b02810](https://github.com/apache/plc4x/commit/4b028108be8ab9d9af391227f2305dbfcaf3d4e4))
  - use byte[] as base for enumerated for now ([cf8ebfe](https://github.com/apache/plc4x/commit/cf8ebfe7051c0624457580a70c90df85a11307e3))
  - several bugs fixed by enable testsuite in plc4go ([63c6301](https://github.com/apache/plc4x/commit/63c63017b2879442a2889954b2b195936de83cfb))
  - fixed multiread ([48b05ae](https://github.com/apache/plc4x/commit/48b05aed3faf3df97590b82ce86eff00dfa1198f))
  - fixed issue where context tags were too greedy ([9304e11](https://github.com/apache/plc4x/commit/9304e111435378e5162bc5a16b9ea6cd48a45f7d))
  - fixed atomic file reads ([6db9132](https://github.com/apache/plc4x/commit/6db9132c4a1e7c2c7e9e5d060e8746f2d73188f9))
  - fixed BACnetConfirmedServiceRequestDeviceCommunicationControl and error codes ([ce9aa60](https://github.com/apache/plc4x/commit/ce9aa604eebd01965d4a05d98220d320bdb6de44))
  - fix length calculations ([f7d0770](https://github.com/apache/plc4x/commit/f7d077013928e5ae28dac343459649f6fc09d3d6))
  - enable tag numbers > 14 ([dc9b50e](https://github.com/apache/plc4x/commit/dc9b50eb64e1d6dfb993fe02c90bfddbb2f84c4f))
  - cleaned up and worked on COV ([4d24d1c](https://github.com/apache/plc4x/commit/4d24d1c551c5c0cd8aef384dd08d6a7f06cc54d4))
  - cleaned up some strange constructs ([3cda25c](https://github.com/apache/plc4x/commit/3cda25cf6c11b94195edc40043fec7931ee359ab))
  - partially fix build ([4e8b274](https://github.com/apache/plc4x/commit/4e8b274b3e64a4bddb6a7d1ea6051be99516f003))
  - fixed typo in BACnetContextTag ([8b856ea](https://github.com/apache/plc4x/commit/8b856eaabf249202690d19c69943f2d338d450c3))
  - switch object type uint 10 to enum ([2a6dc22](https://github.com/apache/plc4x/commit/2a6dc222ab75e885b91ea159da2ecf396279cbd7))
  - fixed some issues in bacnet spec ([3e966c8](https://github.com/apache/plc4x/commit/3e966c8f84bdb141b1d7e405997081486cb70647))
  - implement BACnetErrorPasswordFailure ([457919b](https://github.com/apache/plc4x/commit/457919ba1324b33f936e25db8299c7b61ebce8e9))
  - BACnetComplexTagPropertyIdentifier should now properly render a enum ([15832cd](https://github.com/apache/plc4x/commit/15832cd8e9e7775156b11f6ebfe8b591dc02b091))
  - implement BACnetUnconfirmedServiceRequestIAm ([6563eb9](https://github.com/apache/plc4x/commit/6563eb9b5be2ab3abe289345a0746ebeec03a9b0))
  - substract encoding field from length ([aa31014](https://github.com/apache/plc4x/commit/aa31014b664430f5517074fdd126c0fa2e9f3f42))

- plc4go/codegen:
  - fixed argument handling of complex argument types ([d9a1c5d](https://github.com/apache/plc4x/commit/d9a1c5d5cba5d2eccf33f1180615de87d953fc33))
  - fixed issue where enum dicriminators omited filed names ([38b28dc](https://github.com/apache/plc4x/commit/38b28dce35398490779e7e85adb7498a30dc3a3d))
  - removed broken code ([72ee998](https://github.com/apache/plc4x/commit/72ee99887b8bca5f66106de0c4e52d12d940bd93))
  - Truly close connection (#338) ([5cda65f](https://github.com/apache/plc4x/commit/5cda65ff3a78791d521cc90b86d4e3b5ecd4f2ee)) ([#338](https://github.com/apache/plc4x/pull/338))
  - DefaultPlcWriteRequest interface conversion, cause it not implement (#335) ([cc013c6](https://github.com/apache/plc4x/commit/cc013c6b95bdee54dd5d9b22ed8d3b0279af7d4b)) ([#335](https://github.com/apache/plc4x/pull/335))
  - fixed issue were a discriminator field overshadows a virtual field ([526bedf](https://github.com/apache/plc4x/commit/526bedf5eccb8955b983f353d2c55efda0823b7f))
  - work on time support ([0c3f74e](https://github.com/apache/plc4x/commit/0c3f74efd613475693f6c849e41edd69e1829bb8))
  - work on time support ([e91f634](https://github.com/apache/plc4x/commit/e91f634277bde92364aa0fca1e6212d9df446844))
  - fixed checksum fields ([67cc89c](https://github.com/apache/plc4x/commit/67cc89c7ae929ded26e43e5b49d58ae834ad31ab))
  - treat EOF als optional reset ([908bfc7](https://github.com/apache/plc4x/commit/908bfc7f6eb115e56ed5a429b933090b9a789e83))
  - fixed broken tracer ([1bcd9d5](https://github.com/apache/plc4x/commit/1bcd9d563e538a7c8dd23da3347df74baed33c44))
  - fix build ([deb4a19](https://github.com/apache/plc4x/commit/deb4a19433441de8a92dd8a223106e88996bce6a))
  - avoid opening contexts for manual fields ([7d9f7b7](https://github.com/apache/plc4x/commit/7d9f7b772aae69f500f6026acfd5034fc8a49e79))
  - wrap arrays to avoid collisions ([20fd472](https://github.com/apache/plc4x/commit/20fd472917634a3d78323fc7a06161867e9deb15))

- plc4cs:
  - fix arrays access ([68a276e](https://github.com/apache/plc4x/commit/68a276ed6a294da6d56a45c077e04ee87e7a1693))

- plc4net:
  - Changed the KnxDatatypes to use BitStrings for BYTE, WORD, DWORD and LWORD ([06a9b34](https://github.com/apache/plc4x/commit/06a9b3414f4ecf8d2b02905f9264831351f2c508))

- plc4go/tools:
  - minor fixes on plc4xpcapanalyzer ([fd38251](https://github.com/apache/plc4x/commit/fd38251447fc3fc78974ff740a9cafad51e0fd3c))

- opcua:
  - Fixed a few issues with the OPCUA Encryption ([4951281](https://github.com/apache/plc4x/commit/49512816f666f58d7fcca4e08ce2d70e32d5d2b8))

- codegen/go:
  - Undid some of my changes and fixed the tests differently. ([542197b](https://github.com/apache/plc4x/commit/542197b9d37f8c7f10a3a9ee73b617b43be06b85))
  - Moved the StaticHelper.go to the new position and adjusted the pom to copy stuff from new locations to new locations ([6551982](https://github.com/apache/plc4x/commit/65519824dd4210f1460a4f2d3ffb73d98c16473d))
  - Made manual fields not get explicitly cast ([ea6cf34](https://github.com/apache/plc4x/commit/ea6cf344bc53699371058e4dfcf969234c2a9243))
  - Adjusted the packages of the ParserHelper and XmlParserHelper to the new structure ([dbaaccb](https://github.com/apache/plc4x/commit/dbaaccbb925117618556e79f3af3ca07baa36cb2))

- codegen/plcj4:
  - fixed missing break statement for unknown fields ([e9d6c9d](https://github.com/apache/plc4x/commit/e9d6c9da12e429474ed26892c6e41095f8729dca))

- plc4j/opcua-server:
  - Moved pinned jaxb-runtime from parent to nifi processor. It isn't needed and causes issue with the milo stack ([54c90db](https://github.com/apache/plc4x/commit/54c90dba81ed287308fb14c9dec06b8b23310a47))

- plc4go/bacnet:
  - ErrorCode and ErrorClass ([17d7dbe](https://github.com/apache/plc4x/commit/17d7dbe435c39800b51c93c421852eafafc7340f))
  - fixed enum write mapping ([83c31a4](https://github.com/apache/plc4x/commit/83c31a48a1c2fe9334959f9e5bcd879faa533175))
  - fixed enum mapping ([0fa1820](https://github.com/apache/plc4x/commit/0fa1820d036cc0b36f75da463a1054b64b88e2c4))
  - export tag functions ([1fd65f1](https://github.com/apache/plc4x/commit/1fd65f15e4501539631296ce395e4ba9964f3e04))
  - fix static helper ([78d068e](https://github.com/apache/plc4x/commit/78d068e4663a0dc56d1c679e155815a0d6af5915))
  - fix issue with wildcard handling ([0e58917](https://github.com/apache/plc4x/commit/0e58917b77fc21e338860a4fe87ff6b1ee900195))

- plc4j/scraper:
  - Handle broken connections when an exception occurs when calling getConnection() ([ec16e74](https://github.com/apache/plc4x/commit/ec16e7426c3f447ef4a8a038206f975eae92d4b9))

- plc4j/kafka:
  - revert change ([5418f8e](https://github.com/apache/plc4x/commit/5418f8e67a55be630925ed48aa576ac5d896a283))
  - Class Loader not finding transports correctly ([f68d980](https://github.com/apache/plc4x/commit/f68d98095abf40522d62bb09358a36ef60ab29bc))

- plc4j/bacnet:
  - only dump packages when flag is set to true ([9278ac7](https://github.com/apache/plc4x/commit/9278ac72a7f3307b305dc141a02e9e963b2f3d88))
  - Replaced the code for detecting libpcap with the annotation-based approach. ([d2bb8e4](https://github.com/apache/plc4x/commit/d2bb8e4938a636ffe7c0b10bd8052572149a6e12))
  - adjust artifact id ([c9a5bd5](https://github.com/apache/plc4x/commit/c9a5bd526085d94528cbc3baa7456d662b234a05))
  - re-enabled protocol logic ([9b6f7f2](https://github.com/apache/plc4x/commit/9b6f7f2bfcf299ca6cdfff061e94b53e114e4f4c))
  - disable RandomPackagesTest.java ([8edf938](https://github.com/apache/plc4x/commit/8edf9385ac40c7e58897ed7124370ced17c9ecfb))
  - only run on systems with installed pcap ([a3843b4](https://github.com/apache/plc4x/commit/a3843b480f623e60d685b85b4724055f7c0e8c4b))

- codegen:
  - optional fields respect condition now on write ([dd3b0cf](https://github.com/apache/plc4x/commit/dd3b0cf1f3d642c506cc56a729d16b4daff89ffb))
  - switch to lookup properties in parents too ([9c15f0d](https://github.com/apache/plc4x/commit/9c15f0d70c73275b72f4c109499f25baf181b41d))
  - Ensured the codegen works when referencing subtypes ([f64f02c](https://github.com/apache/plc4x/commit/f64f02c37f7cbb3041584fe3d974da748cea2a8b))
  - Enabled the virtual field tests and fixed a related bug in plc4go ([d573c0d](https://github.com/apache/plc4x/commit/d573c0de12b5f6ca23fad146422957d7e7cde5c8))
  - Enabled the virtual field tests and fixed a related bug in plc4go ([a5f8e6a](https://github.com/apache/plc4x/commit/a5f8e6a9c0bbd5a85f7fb221fcb0893b8930c182))
  - Added more testcases for the new field types and added support for them in Java, C and Go ([831db3e](https://github.com/apache/plc4x/commit/831db3ec70b45dc9655c1f93711109717c5bd211))
  - implement missing method for virtual fields ([7c963a0](https://github.com/apache/plc4x/commit/7c963a06c3c6c4fc6c9eb28c50f689e652c0a067))
  - Worked on implementing checksum fields ([845953f](https://github.com/apache/plc4x/commit/845953fdd17ac9f55ed08bccae8cef10be892050))
  - fix issue with double scheduling by removing unused code ([080a248](https://github.com/apache/plc4x/commit/080a24826528e7702b48554149484dce3ff4a68d))
  - fix lookup for variable literals by considering builtins ([278598f](https://github.com/apache/plc4x/commit/278598fc8c4cd65c91350ee8c2759e2505ef214a))
  - fix lookup for variable literals by considering all fields ([f18ae5c](https://github.com/apache/plc4x/commit/f18ae5ceec19c6bae05f9d13ad0a0b94fcead9b6))
  - include arguments by variable context lookup ([81fcbbb](https://github.com/apache/plc4x/commit/81fcbbbd844240f2a8e15852a0e579878e6fdf15))
  - small typo fix ([d93719a](https://github.com/apache/plc4x/commit/d93719a6e2ef93014ce35f759b1a37a31dd5c1db))
  - fix partial type references on variable literals ([91daf10](https://github.com/apache/plc4x/commit/91daf10fc1591f5294c544d948165c96e7bac0f1))
  - fixed remaining issues ([8b58275](https://github.com/apache/plc4x/commit/8b58275099f734eef335b8298182997b2235b1ba))
  - fixed problem with complex arguments ([e4be2ce](https://github.com/apache/plc4x/commit/e4be2ce29b243fbbe65972f4e6565e3998cd2ec1))
  - fixed typo in plc4go ([ddcd8a8](https://github.com/apache/plc4x/commit/ddcd8a8ea97f371f261b6a7a62ead716282ce465))
  - fixed typo in plc4go ([95c0d30](https://github.com/apache/plc4x/commit/95c0d30f530d0a496495882838a449ca1ce5e963))
  - fixed remaining helper calls ([b8c1881](https://github.com/apache/plc4x/commit/b8c1881209c746e0740ab5d909482d4d2fe0b703))
  - fixed some issues in plc4c ([5c89ba2](https://github.com/apache/plc4x/commit/5c89ba2570e7f57773188e2a0094b5019f072fcb))
  - fixed some issues in plc4c ([27feac1](https://github.com/apache/plc4x/commit/27feac1828d2165453034a070fb847179111963e))
  - fixed some complex references ([d5e3848](https://github.com/apache/plc4x/commit/d5e3848a7a5ea365d8a2ba5f61298ea214c876d9))
  - fixed small issues ([d1ce919](https://github.com/apache/plc4x/commit/d1ce91959dc4b8b0acd736125b26e95e45ccac87))
  - replace getComplexTypeReferences() with getNonSimpleTypeReferences() ([ffc9279](https://github.com/apache/plc4x/commit/ffc927989a4cf14397d7376a55ded9d897cbca5b))
  - fixed issue with dataio in plc4cs ([8de2719](https://github.com/apache/plc4x/commit/8de2719b715112802ec7ac807bfd6b25ef2a441d))
  - fixed issue with dataio in plc4j ([82af78d](https://github.com/apache/plc4x/commit/82af78dd058e1ca0532c08c6ee3605c522c1f23a))
  - fixed issue with dataio in plc4j ([4f275db](https://github.com/apache/plc4x/commit/4f275db14a839935f6467f3de7bf16530fcddd78))
  - fixed issue with array in plc4j ([0f2903b](https://github.com/apache/plc4x/commit/0f2903bcbd970516d0119dff4ee5d030698566fc))
  - fixed issue with manual arrays ([c84266a](https://github.com/apache/plc4x/commit/c84266ab34a793fd1d1df48c0624b5735b227b5c))
  - fixed small issues with c generation ([1ddd589](https://github.com/apache/plc4x/commit/1ddd5890c3a30b7f395d5331731d2a62b181b1b6))
  - fixed usage of deprecated utils ([e01cce8](https://github.com/apache/plc4x/commit/e01cce8ff2480392b09996aa6cbf4d5934abf73d))
  - fixed data-io c ([c874d78](https://github.com/apache/plc4x/commit/c874d78f43865480a50306258d1c1313ab38be13))
  - fixed issues with wrong type handling ([eaa18a6](https://github.com/apache/plc4x/commit/eaa18a6046220bad2673c50d91ab88cd2c6513b9))
  - fixed issues with wrong type handling ([7eb4c99](https://github.com/apache/plc4x/commit/7eb4c99360b6a39fa1ba6a91d4f4dd06e4f09e30))
  - fixed issues with wrong type handling ([aa61ce2](https://github.com/apache/plc4x/commit/aa61ce23a85061b20aa071b7b4badc923aa9ec5e))
  - fixed some issues in go codegen ([31db05d](https://github.com/apache/plc4x/commit/31db05df6e140c5f6d5db99ff9a9ec251e6ee806))
  - remove temporary fix ([b435e82](https://github.com/apache/plc4x/commit/b435e824fc9576ceb3c7a31cfa1a41175069ac6e))
  - temporary fix for duplicated parser args on dataio definition ([5b01ba0](https://github.com/apache/plc4x/commit/5b01ba05697b8a96c2ca71dfc14ca2c621ac7239))
  - fix order ([00570ea](https://github.com/apache/plc4x/commit/00570eaef74a1c56142fa83be5fc4d98d55bffd1))
  - put DataIo below ComplexType ([2443aed](https://github.com/apache/plc4x/commit/2443aeda78bd85236e15afab5d287acd4e1e5b7b))
  - fixed small issue in go tempate ([9afb345](https://github.com/apache/plc4x/commit/9afb345a1ef7e69c395b7562ebd10b474bb6612d))
  - fixed parent access issues ([b8c0c7b](https://github.com/apache/plc4x/commit/b8c0c7b70bb72f451fc21ea91ea907b61eba6200))
  - fixed issue in CLanguageTemplateHelper ([da3889f](https://github.com/apache/plc4x/commit/da3889fac8b18a9097868ea085637cf6ec9e44be))
  - fix issue where we forgot to set the type definition ([080fe95](https://github.com/apache/plc4x/commit/080fe9559aefe89f0fbcd2ae67f726602924aa4f))
  - small improvements ([d3a7113](https://github.com/apache/plc4x/commit/d3a71139494dd72e079f687671db7d89e0f555b3))
  - fixed some complex usages ([cb2db75](https://github.com/apache/plc4x/commit/cb2db75e8b73c0488c151aa077d99ff203d02eb5))
  - fixed some complex usages ([498413b](https://github.com/apache/plc4x/commit/498413bb616134ce1efa570d35bde3053436be9b))
  - fixed some optional usages ([dc49339](https://github.com/apache/plc4x/commit/dc493393ac0e3bc0f33ad69aba141e7184480dd8))
  - fixed some issues with parent types ([080640d](https://github.com/apache/plc4x/commit/080640dc4d6b93413bf85ceb7b503b96391eb62d))
  - fixed some places where Optional was ignored ([16b58b4](https://github.com/apache/plc4x/commit/16b58b4588f5af0f970fdad3a028ff8c6e61e1f7))
  - fixed problem with concurrent modification ([f352cbf](https://github.com/apache/plc4x/commit/f352cbf68dc9870a1c43ba3a68b1e83aa0024e5d))
  - use bitlengths for manual fields ([ae1d27a](https://github.com/apache/plc4x/commit/ae1d27a956acf40bc55f4ac7c068f15e367c9263))
  - fixed issues writing float and double ([42c494b](https://github.com/apache/plc4x/commit/42c494b8c9bd4e6bd42f3e76614403e68b0ba4d5))
  - assert field now should also map ArrayIndexOutOfBoundsException to a ParseAssertException ([30a2cb2](https://github.com/apache/plc4x/commit/30a2cb2f7a2537bb767c9d30c4b7009fe80c8b69))
  - fix several small issues ([d3d22f1](https://github.com/apache/plc4x/commit/d3d22f16aa86b5512fca06c2539b000e7fb84f2c))
  - fix several small issues ([92e29da](https://github.com/apache/plc4x/commit/92e29dad9832fa6b4fd3cea6d85be43fd3092708))

- kafka:
  - Update Kafka Connector to use the PooledDriverManager ([8af6661](https://github.com/apache/plc4x/commit/8af6661fd61824bde5e7922f85db29b765e4c686))

- plc4j/codegen:
  - make writers for complex and enum null safe ([6b6a8fd](https://github.com/apache/plc4x/commit/6b6a8fdcf65a87769e7914409d94ccaad2421f32))
  - fixed issue were virtual fields weren't properly generated ([75a57c0](https://github.com/apache/plc4x/commit/75a57c067fb970ae6af97a2c55948d3c37f5b70a))
  - fixed typo in write data time ([2f587a7](https://github.com/apache/plc4x/commit/2f587a78af30916161a126a7eb4a1d0d5a31a2dc))
  - fix issue with duplicated params ([3e6de53](https://github.com/apache/plc4x/commit/3e6de53335b4251ccc26792e3048bca2277adbd7))
  - fix overflow issues on byte[] ([43baaed](https://github.com/apache/plc4x/commit/43baaedd5b73d8ae887b85ee9e438e6e2065dc96))
  - fix build ([c9d4dd4](https://github.com/apache/plc4x/commit/c9d4dd450648c7b2b066354c6562deac6d86caf1))
  - fix extension ([9179cf1](https://github.com/apache/plc4x/commit/9179cf1a8ad659b1766d8e9408d1ff38108a1e91))
  - treat 0xff als valid value constant for byte ([5e6cc92](https://github.com/apache/plc4x/commit/5e6cc92778de9321cad08e7da78b1e5768bce708))
  - mark more methods as deprecated ([8d06fd7](https://github.com/apache/plc4x/commit/8d06fd7d82076464847b45f68e3c46513a439683))
  - fix issue when bitlength is < 3 ([85e65e2](https://github.com/apache/plc4x/commit/85e65e212aeac74bce8f50db6b3ef5ab90d7320f))
  - port ascii boxes from golang ([61c487c](https://github.com/apache/plc4x/commit/61c487c8283d9257574240653278698e7b0b3f32))
  - fix issue where unboxing didn't work properly ([5858893](https://github.com/apache/plc4x/commit/585889377626c7fae9a5cba323d00fb758ecedeb))
  - wrap manual fields with a context... ([28b5bd4](https://github.com/apache/plc4x/commit/28b5bd4879278367bc3b8dc06dc673d65ccd2224))
  - remove old codegen ([681517e](https://github.com/apache/plc4x/commit/681517edc7a32fab05fa9f702b5ad0936df1b151))
  - wrap manual fields with a context... ([14b7737](https://github.com/apache/plc4x/commit/14b773773a830123270f51e7ec1e6ee3d69bccf7))
  - remove old codegen ([7509183](https://github.com/apache/plc4x/commit/750918352806326920a9f1978331b5cae0d994f1))

- knx:
  - Fixed a second decoding bug in the 16bit floating point numbers of KNX ([0f2533c](https://github.com/apache/plc4x/commit/0f2533cbc91c8c808613d1bbddd728ca1f76834a))
  - Added the bug to the RELEASE_NOTES ([f23cd21](https://github.com/apache/plc4x/commit/f23cd21bb34baaa634e97a77b646f5f1ec0509f0))
  - Fixed a problem in the mspec-generation which caused problems in decoding mainly floating-point values ([375bff2](https://github.com/apache/plc4x/commit/375bff2c37c17b76ba44a5d4eb6d1637fae4c409))

- bacnet/plc4go:
  - use platform types to create integers ([82cbcfe](https://github.com/apache/plc4x/commit/82cbcfe144defcc94e2a4084b0840db9d6e8d782))
  - fix static_helper proprietary write ([1bf5884](https://github.com/apache/plc4x/commit/1bf5884d49ea0d934a3f6e0bead0261cf371b082))

- plc4py:
  - added PlcDriver and moved Mock Driver to plugin ([ce60af0](https://github.com/apache/plc4x/commit/ce60af02f545254eedd76c4390b5dbc710c56acf)) ([#358](https://github.com/apache/plc4x/pull/358))
  - change from using poetry to setuptools ([07a3903](https://github.com/apache/plc4x/commit/07a3903fda0178f7685413cdf33b6005d60c354b)) ([#355](https://github.com/apache/plc4x/pull/355))
  - Update the PlcRequest code ([1fb2d85](https://github.com/apache/plc4x/commit/1fb2d85c1e8537384e6b1c6b42698e3eb467bcac)) ([#348](https://github.com/apache/plc4x/pull/348))
  - Add a compile phase to install python dependencies ([ed5047e](https://github.com/apache/plc4x/commit/ed5047e4a13e7200ba3a61b961d232c04d20d318))

- plc4j/nifi:
  - Bumped the NiFi version to a version that is not vulnerable to CVE-2022-26850 ([8b3ef69](https://github.com/apache/plc4x/commit/8b3ef69a7867cbe46c2a40e2df02cc55a2e7d7da))

- codegen/plc4c:
  - Fixed a wrong include statement ([1cc51e3](https://github.com/apache/plc4x/commit/1cc51e333ca79324291187e901ff16dc85957298))
  - Fixed a wrong include statement ([ebe3bfc](https://github.com/apache/plc4x/commit/ebe3bfc9c9d0a9931ae58b20d30ff7cec00e93f9))
  - Fixed a wrong include statement ([29d8011](https://github.com/apache/plc4x/commit/29d8011486ccca717ac1ff4ace87ff3c8a5e7a0b))

- protocol/test:
  - Changed the checksum fields to actually call checksum functions ([18ea08e](https://github.com/apache/plc4x/commit/18ea08e06b3ca70d00dc210b6017d2a27af9ae57))

- plc4j/driver/modbus:
  - Renamed some internal variables (previous copy-paste error) ([5352166](https://github.com/apache/plc4x/commit/5352166f42402d59428d28528a43df8efd3127b6))

- plc4/integration/calcite:
  - Managed the kotlin version up to 1.6.10 to work around CVE-2022-24329 ([cb46fe2](https://github.com/apache/plc4x/commit/cb46fe271d7ee1dc04d06f9a8f410a04202d5fed))

- plc4go/c-bus:
  - used virtual field as discriminator ([e097816](https://github.com/apache/plc4x/commit/e097816a224cece0736ea1704a1b510dc3e9213b))

- protocols/c-bus:
  - fix power up fix ([9286215](https://github.com/apache/plc4x/commit/928621598ae6d525a73b3f2d518f236ffbeff2a6))
  - fix power up ([708c4e8](https://github.com/apache/plc4x/commit/708c4e8415b5ae4513fc78bd6829a51f4b9013b7))
  - associate type with Reply type ([edda7d2](https://github.com/apache/plc4x/commit/edda7d20a07c2c192cf8d862769514bf7f673304))

- codegen/plc4x:
  - fixed cases where byte and short didn't work for discriminators ([fcd6c8d](https://github.com/apache/plc4x/commit/fcd6c8d59532bd237f0e8fc41099d1b363ce02b4))

- cbus/plc4j:
  - fix compile issues ([115f0e1](https://github.com/apache/plc4x/commit/115f0e1ec04197f7a46e859da92b724008a2e91c))

- codgen:
  - lower log output by setting some errors on debug for now ([0cd6560](https://github.com/apache/plc4x/commit/0cd65608eca97f02cf9f67b0461daa77cee15d98))

- bacnet/plc4j:
  - enable uint64 support for virtual fields ([f4901fa](https://github.com/apache/plc4x/commit/f4901fa4c85daf0e36eccb6027560417f713fedb))

- codegen/plc4j:
  - fix cast exception relating BitInteger ([37eab1c](https://github.com/apache/plc4x/commit/37eab1c365ae2ad6f4b7393fefff03c3546e3ef9))
  - hotfix for missing BitInteger support ([54ec98f](https://github.com/apache/plc4x/commit/54ec98f2a834eea18d7c08ac9dbead016ef7f50a))
  - try to fix big integer conversion issues ([e492df7](https://github.com/apache/plc4x/commit/e492df7fea1052b4e890cedbbacbd7ecaf982017))

- plc4j/knx:
  - Fixed code where plc4j code was not Java 8 compatible ([7707680](https://github.com/apache/plc4x/commit/7707680d672bfbf9649f24b8cbd289fb6cd4e8ec))

- plc4j/test-utils:
  - Made the ManualTest a bit more robust against cardinality-errors. ([72a7da2](https://github.com/apache/plc4x/commit/72a7da20744a2726906a56ddb2f8f44ef16b1173))

- protocols/modbus:
  - PLC4X-335: Using single-byte modbus types resulted in unexpected behaviour ([a8ecff1](https://github.com/apache/plc4x/commit/a8ecff16abdb1356bd39ed2c2d54a040924904d9))
  - PLC4X-335: Using single-byte modbus types resulted in unexpected behaviour ([55044d6](https://github.com/apache/plc4x/commit/55044d64896e2e2eabac5c6a3dec8a85274ebf8e))

- docs:
  - Updated the pcap docs for when setting up a developer env (#318) ([d7c8599](https://github.com/apache/plc4x/commit/d7c85997178c0b9f3647b9e9a53b4db28358409f)) ([#318](https://github.com/apache/plc4x/pull/318))

- plc4j/example:
  - Changed the log level back to `info` so the application actually outputs something. ([434c2b4](https://github.com/apache/plc4x/commit/434c2b4cd0bd82bd525fa7eae5453dd680b9eadd))

- plc4j/testing:
  - Fixed a bug not correctly detecting the libpcap version on Windows. ([4ee0ca5](https://github.com/apache/plc4x/commit/4ee0ca567d32f0431af991aff125e8f0a22f1085))

- plc4c/codgen:
  - fixed test utils ([321da0b](https://github.com/apache/plc4x/commit/321da0b63c57edfbc92ad611e2a91ade47e59ae1))
  - fix const values go missing due to equals now implemented ([0485dbb](https://github.com/apache/plc4x/commit/0485dbb7c0b48f26ff80ef426766fab431b0a74a))
  - fix enum handling of bool on aarch64 ([f689a03](https://github.com/apache/plc4x/commit/f689a03a9cceac2b0d5128137bad571f1689da8b))

- plc4j/canopen:
  - Fixed a bug introduced by my recent refactoring ([5b81e70](https://github.com/apache/plc4x/commit/5b81e70ffeba51b057253421872815567994c8f5))

- plc4j/simulated:
  - Fixed the way the simulated driver outputs values on "STDOUT" ([39edc4d](https://github.com/apache/plc4x/commit/39edc4d7b1e46ff728939234baa7b4d54f8de69b))

- plc4j/connection-cache:
  - Implemented Writing for the connection cache ([3f92555](https://github.com/apache/plc4x/commit/3f92555c539705ebe4c4125fa312fc2f2f94c29c))

- plc4go/codgen:
  - added missing validation field ([5aa4e72](https://github.com/apache/plc4x/commit/5aa4e72001838ba1969154a72cbb9b4c21a571c9))
  - sync code ([15524ae](https://github.com/apache/plc4x/commit/15524ae3623021faf34cb7d6577332d290ce2e5c))
  - fixed issue with broken arrays ([3aef587](https://github.com/apache/plc4x/commit/3aef5876ad2ca6b09be937f4a729b349ca4545e7))
  - fix optional fields not emitting contexts ([4addf5e](https://github.com/apache/plc4x/commit/4addf5e5e2993699e14833adc29545aed2b51bc7))
  - update xml and json in s7 iotest ([740ecbc](https://github.com/apache/plc4x/commit/740ecbc229482480bd76c863ff2dca352fd10657))
  - migrate eip suite ([e8a4de5](https://github.com/apache/plc4x/commit/e8a4de506e0a5b26f3074aa71e9c8b49beffe05b))
  - migrate eip suite ([6b8c26d](https://github.com/apache/plc4x/commit/6b8c26d6e372a951e072f05bb79d5080b49b2996))
  - migrate eip suite ([d8720e6](https://github.com/apache/plc4x/commit/d8720e626a733a5908dc20e3feefc6ec89423f38))
  - fixed array and string problems ([342f3b2](https://github.com/apache/plc4x/commit/342f3b20edaf8fdfaa4259307fea43cf73439f01))
  - fix merge messup ([f7df8d1](https://github.com/apache/plc4x/commit/f7df8d1bc1f6f872841d08928f172e0edd528f0b))
  - progress on golang codegen ([1c98edf](https://github.com/apache/plc4x/commit/1c98edf0373cac2655bde7bc9d422d33d57950ab))
  - fixed some issues in GoLanguageTemplateHelper ([950b5c5](https://github.com/apache/plc4x/commit/950b5c526dd5b1fd4913fdc9dad0d19bb375b202))
  - fix optional fields not emitting contexts ([9b31501](https://github.com/apache/plc4x/commit/9b31501cea018fc743af8da47cdf0cf70e4a14f9))
  - update xml and json in s7 iotest ([0b8b738](https://github.com/apache/plc4x/commit/0b8b73897880bf7ab6286031e2c4274b00c21aee))
  - migrate eip suite ([b85fa34](https://github.com/apache/plc4x/commit/b85fa3427f5a52f137a318dd4881680244486cb9))
  - migrate eip suite ([8edfa90](https://github.com/apache/plc4x/commit/8edfa90bbd3fca34adfeb2544e0ff6ec0825dc53))
  - migrate eip suite ([30865e2](https://github.com/apache/plc4x/commit/30865e29904fcecf37f5bf6d9518f2326b8af64e))
  - fixed array and string problems ([f919215](https://github.com/apache/plc4x/commit/f919215cf926daaa693ea16e8c09d647a877f1ec))
  - fix merge messup ([743a7b2](https://github.com/apache/plc4x/commit/743a7b2952bbe68f556796d3e121dd766c25c698))
  - progress on golang codegen ([0c47784](https://github.com/apache/plc4x/commit/0c47784e88835c68a754cd2814248f4ce90bdd3e))
  - fixed some issues in GoLanguageTemplateHelper ([c7fe9db](https://github.com/apache/plc4x/commit/c7fe9db9c6cedb4b5ee260d4ffa238bfe3aee317))

- plc4j/codgen:
  - write encoding of strings fixed ([e7dcb11](https://github.com/apache/plc4x/commit/e7dcb11a73e96879c6e710d5966db25e9edef21b))
  - fixed json ([9dd393c](https://github.com/apache/plc4x/commit/9dd393c23cf8aca97b55bfa9b047177818ef7073))
  - fixed assert field ([19e898f](https://github.com/apache/plc4x/commit/19e898f1b2e16eb1ea73b464b5aaea8261d3fe5e))
  - fixed assert field ([51fadd5](https://github.com/apache/plc4x/commit/51fadd5bffec992c5f405e9b28c34f88f98b5bba))
  - avoid duplicating params ([9bda71f](https://github.com/apache/plc4x/commit/9bda71fdb453a77b3c4cb504005a67bb0232613c))
  - fix bacnet spec ([3e3d3da](https://github.com/apache/plc4x/commit/3e3d3da451ec88155e6e14d3affec24197f003e1))
  - fix build ([0e0fcb3](https://github.com/apache/plc4x/commit/0e0fcb363cc523581371bf81a72d70b37fad40fb))
  - fixed issue with missing type reference params ([d559bfe](https://github.com/apache/plc4x/commit/d559bfec64c686219f3f63abb87b7fdca4ff6900))
  - fixed issue with missing type reference params ([1329013](https://github.com/apache/plc4x/commit/1329013afc1afe2bb82866060e48c1cd7b71675e))
  - refined error message ([059ded7](https://github.com/apache/plc4x/commit/059ded788943469234aa78c64e8caf92494a4ec9))
  - added error when auto-migrate won't work ([d098ed2](https://github.com/apache/plc4x/commit/d098ed2a4333773960077ac12bc94df4c64af9aa))
  - fix edge case with broken patterns. ([012e319](https://github.com/apache/plc4x/commit/012e31902730da0df57279d21e1927181d1856f9))
  - fixed s7 io test ([736c0fe](https://github.com/apache/plc4x/commit/736c0fed5a41e47229ab46ae9fe37cc976f48c9e))
  - fixed data reader complex working with a empty logical name ([9bfbdfe](https://github.com/apache/plc4x/commit/9bfbdfe64f65fff4dde5aab948a6ef0f9fb3b150))
  - enum wip ([f02af0a](https://github.com/apache/plc4x/commit/f02af0a79a5c312f5824bec0bad2408eb03b716c))
  - implement array type reading with new readern ([1e819e6](https://github.com/apache/plc4x/commit/1e819e6421af887843c5f36b4b21952142549f67))
  - migrate the protocols ([5ade100](https://github.com/apache/plc4x/commit/5ade1001dd1484888a5695a692021bb2e918d072))
  - move to type ref args ([cad0de6](https://github.com/apache/plc4x/commit/cad0de64e7e292e595a1a6b33fc9d6f59f8400c2))
  - add byte order to read/writer, add enum example fixed ([d8b45c3](https://github.com/apache/plc4x/commit/d8b45c3eb28e09e9d99ee048b5eaea75afa0ae12))
  - add byte order to read/writer, add enum example ([c492327](https://github.com/apache/plc4x/commit/c492327b6c2053bf6680dbd5dc66b1942b4f8e2f))
  - draft a const call adjustments ([92c8551](https://github.com/apache/plc4x/commit/92c8551d231dc60827afa79be45b265396ed95d8))
  - draft a const call ([ee1a6d2](https://github.com/apache/plc4x/commit/ee1a6d252902e0b080caff15f8f3b0ee097ea7fd))
  - fixed json ([f124291](https://github.com/apache/plc4x/commit/f124291521206f00ae861efdd196d5b7baf17f89))
  - fixed assert field ([4a5f261](https://github.com/apache/plc4x/commit/4a5f261c38eb9e564d74d1a57c1e83b5b126ada6))
  - fixed assert field ([0be34a6](https://github.com/apache/plc4x/commit/0be34a6287bf7fbbadb3026a57af204124c90b39))
  - avoid duplicating params ([194cff3](https://github.com/apache/plc4x/commit/194cff3a56516bc6d2a3617e94d5447c6b97fbdf))
  - fix bacnet spec ([8a16e46](https://github.com/apache/plc4x/commit/8a16e461a7e9dddc7b949b872cf670cc1949aef1))
  - fix build ([b3ad562](https://github.com/apache/plc4x/commit/b3ad5622c4bf061b8ed775f03dcb27919fd2ea9f))
  - fixed issue with missing type reference params ([4a28fc0](https://github.com/apache/plc4x/commit/4a28fc0e33592b7f9a6a5234fc5a1a753310d6e1))
  - fixed issue with missing type reference params ([121fc78](https://github.com/apache/plc4x/commit/121fc789203125f873d586663c5b689b9d95c332))
  - refined error message ([8f06618](https://github.com/apache/plc4x/commit/8f06618c1dc30510e62efbd680ca667c1e38ccfa))
  - added error when auto-migrate won't work ([65bc019](https://github.com/apache/plc4x/commit/65bc0198534210f0f2feb1ffcd4e8b8a5761e4b6))
  - fix edge case with broken patterns. ([9d2956d](https://github.com/apache/plc4x/commit/9d2956d4618bd695bb459af2f72648bf5ace9c4d))
  - fixed s7 io test ([e786ae2](https://github.com/apache/plc4x/commit/e786ae2a8aeb7941405ab03adac1316f0487381c))
  - fixed data reader complex working with a empty logical name ([2ebc8a6](https://github.com/apache/plc4x/commit/2ebc8a609de7dc8e2d1d19cefea38ce305435298))
  - enum wip ([5753d78](https://github.com/apache/plc4x/commit/5753d784f091e03de2b60cb05b84037567c3286b))
  - implement array type reading with new readern ([90e40fd](https://github.com/apache/plc4x/commit/90e40fd5c02b0384ccfa205cfb8e109bf56f6e1b))
  - migrate the protocols ([a17dda1](https://github.com/apache/plc4x/commit/a17dda16699ca4cbd214df61bc4163e3b200ab01))
  - move to type ref args ([3e1e1a4](https://github.com/apache/plc4x/commit/3e1e1a4ef1c6ea1a90b8a44f378d5a85f0efc9a6))
  - add byte order to read/writer, add enum example fixed ([a8d835e](https://github.com/apache/plc4x/commit/a8d835e7b17ac09bdcb7827db4ad4bfb3ae16c65))
  - add byte order to read/writer, add enum example ([5f268d3](https://github.com/apache/plc4x/commit/5f268d3477f3c8e15dd0f5d6856f69457687d07b))
  - draft a const call adjustments ([72852a0](https://github.com/apache/plc4x/commit/72852a022e3034c50805cbcd6aa2cb303fa6450a))
  - draft a const call ([d9701e5](https://github.com/apache/plc4x/commit/d9701e5822e24d4fcb3467305749f766f2a6f244))

- plc4j/writeBufferByteBase:
  - only return data till pos ([12a5c19](https://github.com/apache/plc4x/commit/12a5c1957728d50206248546ad845ddd7bfba99f))

- plc4j/xmlWriter:
  - don't print invalid xml chars ([2c75a95](https://github.com/apache/plc4x/commit/2c75a954bfc011f3b4f236157e64eb5c37a90e3c))

- plc4go/connection-pool:
  - Tried hardening a flaky test ([3d5f184](https://github.com/apache/plc4x/commit/3d5f184a46f69cbbf5f2c85f12632dae283d7c05))
  - Tried hardening a flaky test ([01bb30a](https://github.com/apache/plc4x/commit/01bb30a9ad1c6be21b6065d6dcbc092de385ca41))

- plc4j/asciibox:
  - don't count control sequences ([bd66ca0](https://github.com/apache/plc4x/commit/bd66ca093f78e8b86bf15466d52ebdcc2062b428))

- ci:
  - fix prerequisite check for windows and linux ([45756a4](https://github.com/apache/plc4x/commit/45756a4168826143185f52da22786ff9d71502ae))
  - add prerequisite check for libpcap ([202fb9c](https://github.com/apache/plc4x/commit/202fb9c6de642f404b5fd3fd0eca122790e98b01))
  - add prerequisite check for libpcap ([1ebba46](https://github.com/apache/plc4x/commit/1ebba465d911222613286696222530fc83641423))
  - use libpcap from brew on mac as the builtin seems buggy ([769ef25](https://github.com/apache/plc4x/commit/769ef253de61cf7b9449b12ac4d75c1ef9c1784c))
  - use libpcap from brew on mac as the builtin seems buggy ([c8a4509](https://github.com/apache/plc4x/commit/c8a45096ba2ad52b195684f31ecf074f7aebe73a))
  - fix build for platform specifics ([685ca27](https://github.com/apache/plc4x/commit/685ca278776ec80793f6090afb6df599a6f49c0f))
  - cache maven .m2 for actions to speed up builds ([ff3176e](https://github.com/apache/plc4x/commit/ff3176edb8d2aad7d84c40a7dacbb592a0b8e527))
  - fix windows build ([c56e8bb](https://github.com/apache/plc4x/commit/c56e8bb1f4ff9cc9f70424c32b63c64d785f4583))
  - fix build (synchronize actions and Jenkinsfile) ([754f9f1](https://github.com/apache/plc4x/commit/754f9f17d592f5ff2d2ca913d75cd096c22365a8))
  - remove java 8 from platform ensurance ([98ce262](https://github.com/apache/plc4x/commit/98ce26262bbd630bae69023d6eab202796aa2595))
  - remove unused profiles again ([f66faf2](https://github.com/apache/plc4x/commit/f66faf2ef1ed20dcf79d9fc575f5ccad2f3742d0))
  - fix build for java 17 ([f4c2ccf](https://github.com/apache/plc4x/commit/f4c2ccfb23d62ee41b1aaf85aa67cd72f3d0b847))
  - fix build for java 17 (use version ranges) ([ace1c30](https://github.com/apache/plc4x/commit/ace1c30e8a6573f7ccd3ed04002b439302d85d63))
  - fix build for java 17 ([64a8ffa](https://github.com/apache/plc4x/commit/64a8ffae7529835ad821d352a74828bfa9f81cd9))
  - don't fail fast on platform builds. ([adaabd1](https://github.com/apache/plc4x/commit/adaabd189c7a47c9daba6a6518e361f717c362b0))
  - enable matrix build for windows ([f82697c](https://github.com/apache/plc4x/commit/f82697cf9c5fba752bbaf4b141f3609c2c4f60cd))
  - set executable bit on mvnw ([e005e6e](https://github.com/apache/plc4x/commit/e005e6eb52bdd4eb686863c7f8c2686432d0d95d))
  - use correct mvnw call for matrix build ([565fe63](https://github.com/apache/plc4x/commit/565fe63a27b47e5f93f0b5f3a192bdbb40990376))

- plc4j/netty:
  - fix NPE on failed connection ([6b26af6](https://github.com/apache/plc4x/commit/6b26af607bd6d89e8458cea236c64cc95448148e))
  - fix NPE introduced with ec3ff0cdfda097d3ae6e9f15f041e4c5e9c05815 ([42e0622](https://github.com/apache/plc4x/commit/42e062298bca0b7433e01e03454df426b2902d70))

- ci/codegen:
  - set streamLogsOnFailures to true ([5920d67](https://github.com/apache/plc4x/commit/5920d67765bd269dfd97eba8c3bfaa4b5c50a4d6))

- .mvn:
  - upgrading mvn wrapper to 3.6.3 (previous LTS version) ([d1c7c2f](https://github.com/apache/plc4x/commit/d1c7c2f31b59800ac783be60470238828b66a355))

- codegen/knx:
  - new manufacturer Can'nX breaks mspec codegen ([9e3693c](https://github.com/apache/plc4x/commit/9e3693cce3bfcb0fb183d0c658e5efb5bd8ee05a))

- plc4go/knx:
  - fieldHandler nil issue fixed. ([804520d](https://github.com/apache/plc4x/commit/804520d9bcdc485e92e3ba147cad2dc92a47423c))

- ci/ensure-platforms:
  - enable batch mode for mwnv ([c589ed4](https://github.com/apache/plc4x/commit/c589ed471d64770dd53a7fdd069d30d2a31990ab))

- plc4x/codegen:
  - fix failing build ([64108c0](https://github.com/apache/plc4x/commit/64108c0b8016b0bd4d12f487e5988611b5e954b5))
  - add tickles id expressions to mspec ([f9d737e](https://github.com/apache/plc4x/commit/f9d737e16ab86e34dc598d74f4ebf51dd2f2f8d0))
  - add missing separator call ([877f1f1](https://github.com/apache/plc4x/commit/877f1f1a45c5afebcc28b7fa29b9da195b9be151))
  - add possibility to extract traces ([ac6841f](https://github.com/apache/plc4x/commit/ac6841f47c0ef0a8c4f740e75297b94e95918469))
  - fix failing build ([3bed94d](https://github.com/apache/plc4x/commit/3bed94de670c9c5909b08c524c270b15755b8879))
  - add tickles id expressions to mspec ([1e1932c](https://github.com/apache/plc4x/commit/1e1932c8e25e3291e0178dc74aa10ca4b155004d))
  - add missing separator call ([233c1da](https://github.com/apache/plc4x/commit/233c1dab8f6e9c3d5adce7d51ebcee277393eb52))
  - add possibility to extract traces ([2e7a347](https://github.com/apache/plc4x/commit/2e7a3478420f8e3d585c193778021974c4740e3c))
  - added a Tracer to enrich generated code with traces ([57cfe3b](https://github.com/apache/plc4x/commit/57cfe3b273c68f0ef223189eba54aed2d1f3975e))

- plc4go/json:
  - fixed wrong string on empty byte arrays ([6ca28c9](https://github.com/apache/plc4x/commit/6ca28c905047191576a736577367e8c013438b31))
  - fixed wrong string on empty byte arrays ([d0817fc](https://github.com/apache/plc4x/commit/d0817fcb1754c002f0b72e6150d5ffe9c6ddfd21))

- plc4go/xml:
  - fixed wrong string on empty byte arrays ([d4d62fb](https://github.com/apache/plc4x/commit/d4d62fbe8443cd8d49d0d48523726ae7a59fa495))
  - fixed wrong string on empty byte arrays ([785ba4e](https://github.com/apache/plc4x/commit/785ba4eade9ccad80d357d702bcbdfed46870faf))

- plc4x/eip:
  - fixed eip tests ([9952386](https://github.com/apache/plc4x/commit/995238638d3a35b2580ee86d967408af1193eaed))
  - fixed eip tests ([20a6aee](https://github.com/apache/plc4x/commit/20a6aee6f04971103db509cd28a623944a2e4c8d))

- plc4go/test:
  - fixed s7io test ([73cc47f](https://github.com/apache/plc4x/commit/73cc47f61d03015999140b9d0f8e8d0df35b7e47))
  - fixed s7io test ([33b5a75](https://github.com/apache/plc4x/commit/33b5a75bf315ee4994d8079f21ad1c91d4cec581))

- java/codegen:
  - virtual field handling ([ef35531](https://github.com/apache/plc4x/commit/ef35531d5a872f29dccddb3a11a135b166958185))

- pl4go/codegen:
  - fixed broken inline if (adding neccessary function) (addon for e030d688) ([22ba15b](https://github.com/apache/plc4x/commit/22ba15bf8c1895177a28cfcff993c640859a76d2))
  - fixed broken inline if ([e030d68](https://github.com/apache/plc4x/commit/e030d688545e0182640f7c2ffefb2309946e9abf))

- plc4x:
  - post release cleanup ([f7941fc](https://github.com/apache/plc4x/commit/f7941fcf91038c32382b65b8af74ec63c6db09b3))

### Documentation

- plc4j/ads:
  - Updated the statemachine for ADS ([c0e73c9](https://github.com/apache/plc4x/commit/c0e73c9c20a896c9be0fcaa8d0b715fe482f5fa6))
  - Added a comment about stuff I learned recently. ([4663ff9](https://github.com/apache/plc4x/commit/4663ff9716bc54656a1e5b5678b400009e1a5b5f))

### Refactor

- plc4j/ads:
  - Continued testing the new ADS drivers Write abilities. ([f03243a](https://github.com/apache/plc4x/commit/f03243a6231735c6a4c6a9da06139b1c81bbd0b9))

- plc-values:
  - Refactored all usages of BitString to use BYTE, WORD, DWORD or LWORD instead. ([4ee794e](https://github.com/apache/plc4x/commit/4ee794ecf2f8afafaf9f8aba8d51238b9aee3184))

- codegen:
  - Refactored all usages of BitString to use BYTE, WORD, DWORD or LWORD instead. ([7e509a6](https://github.com/apache/plc4x/commit/7e509a6b16de8ed596e473c4138d45456e2940c6))
  - Updated how CHAR, WCHAR, STRING and WSTRING fields are handled. ([fe1a2b1](https://github.com/apache/plc4x/commit/fe1a2b1ce16106aa8e2123d709897e5ee7162639))
  - Updated how CHAR, WCHAR, STRING and WSTRING fields are handled. ([3d79f3d](https://github.com/apache/plc4x/commit/3d79f3d352c7f819935980ea2f99fd0eacf2bb58))
  - move WildcardTerm to build-tools ([df6aef1](https://github.com/apache/plc4x/commit/df6aef1016796fe832e94f4168ccbc5ff5d895d6))
  - move WildcardTerm to build-tools ([bd6f3d7](https://github.com/apache/plc4x/commit/bd6f3d770239bb4690599f0f526feb16faac7065))
  - change index on variable literal to make use of optional ([99c274a](https://github.com/apache/plc4x/commit/99c274ad8b84eb890d94de01b23e45680ecb8271))
  - removed inverted throwable logic ([09415e0](https://github.com/apache/plc4x/commit/09415e0801690e33ee94ecd82c5aaad432de336d))
  - delay resolving of type reference ([485dcf8](https://github.com/apache/plc4x/commit/485dcf82ef598a26065bec1b8068505c56d943df))
  - added TODO: about wrongly dispatched type reference ([e9aa564](https://github.com/apache/plc4x/commit/e9aa5643ef4204a20fd5ece8c815d21f16a75769))
  - improve error message ([d99e1e3](https://github.com/apache/plc4x/commit/d99e1e3a1e2948d42eeb5eaa1286bc42fa27e4b6))
  - cleanup type definitions ([9ac5c0c](https://github.com/apache/plc4x/commit/9ac5c0cef6b685d607846c1088ec2c4267714c24))
  - moved default types to mspec base ([6ec6837](https://github.com/apache/plc4x/commit/6ec683787b9b817499ec7d9cf9a924a21ac774f2))
  - cleanup helper and move more functions to their owners ([51f216f](https://github.com/apache/plc4x/commit/51f216f72c316887adf477b2cc5dfd662f5d12d1))
  - fix a bunch of deprecation warnings ([69d82df](https://github.com/apache/plc4x/commit/69d82dfb7e8401a1be600b434350c6acaf4b1edc))
  - add more tracing calls to CLanguageTemplateHelper ([cc0d864](https://github.com/apache/plc4x/commit/cc0d864fc6dc10e2a423a840edd13b303671ca0b))
  - cleanup *TemplateHelperCode ([effb1c1](https://github.com/apache/plc4x/commit/effb1c18ba4840733fce16ad3976679ef2d2dc51))

- plc4go/cbus:
  - added cleanups and logging ([04ff774](https://github.com/apache/plc4x/commit/04ff7740f6101f468d6d362e76a5832a394791eb))
  - use fireConnectionError on message codec fail ([ed401b8](https://github.com/apache/plc4x/commit/ed401b8d15e54d76ca4b432fa28e226daa02fbf1))
  - improved debug output ([80ba5d0](https://github.com/apache/plc4x/commit/80ba5d0cc156b98b92c5855245ab750d2edfa87e))
  - move MMI handling to subscription handling ([6bd288c](https://github.com/apache/plc4x/commit/6bd288c47fc5dda23c503ba18510885dadb8a2fd))
  - put connect parts into methods added todo about power up ([633794d](https://github.com/apache/plc4x/commit/633794d1ea110479012655c3e19e6f01754a93b5))
  - put connect parts into methods added todo about power up ([2f58834](https://github.com/apache/plc4x/commit/2f58834fe5c8d4db35fff93e6f43e86795829601))

- plc4go/connection-cache:
  - small cleanups ([cee041b](https://github.com/apache/plc4x/commit/cee041b96ce05ecda4c0e78ee3eca73740ec1ad5))
  - cleanup ([24e8bf9](https://github.com/apache/plc4x/commit/24e8bf91ea4d39ad187553436979a6f7e5b45e64))
  - Added log statements to the connection-cache ([8e1ab62](https://github.com/apache/plc4x/commit/8e1ab624e110df277a687ddca87a6bd1dec57a12))
  - Renamed the connection-pool to connection-cache and added a Close function to the cache for gracefully closing all connections it manages. ([4760b5e](https://github.com/apache/plc4x/commit/4760b5eae0785243a427762919bf69ed1a098b0e))

- plc4go:
  - remove unneeded methods ([7ac5ab6](https://github.com/apache/plc4x/commit/7ac5ab6030f50b6a7ce6a3ec98efc7a211f757e9))
  - align subscription implementation to plc4j ([19f6600](https://github.com/apache/plc4x/commit/19f660025b0da1545a1ce2fc7a54c8dcc7c9ae44))
  - renamed mashed up NewBoxedWriteBufferWithOptions ([0f32971](https://github.com/apache/plc4x/commit/0f3297114c2736ffafb8e63fdaa0b9c1a16164c9))
  - streamlined plc_browse api ([e0cc944](https://github.com/apache/plc4x/commit/e0cc944d4e20ae9e510fdd9468c0202856298d37))
  - added a IsAPlcMessage() bool method to plc_message ([65adaaa](https://github.com/apache/plc4x/commit/65adaaaa38e300e9af83d32c8f0ec92253c0bcb7))
  - introduced DefaultBufferedTransportInstance to consolidate common code ([4479cb0](https://github.com/apache/plc4x/commit/4479cb05643e9a81ff0987e1e74c515aefab5cc7))
  - reworked the GetNumReadableBytes ([a5fb77a](https://github.com/apache/plc4x/commit/a5fb77aa59b4b0ca4e03ad9306899beb956d022f))
  - restructured package (+moved protocols to a public importable place) ([adc23d2](https://github.com/apache/plc4x/commit/adc23d2f9d0caa8c901e49e1621107dfd44cc45c))
  - fix QA issues ([49a43db](https://github.com/apache/plc4x/commit/49a43db477aee8b0f112a2181387e8b3b4b71a23))
  - add accessors for property and virtual fields ([33e07bd](https://github.com/apache/plc4x/commit/33e07bd55a3a1b26f31dc5fcf81633e43b1bab43))

- spi:
  - optimize RequestTransactionManager ([49809c5](https://github.com/apache/plc4x/commit/49809c54def1c0542fe643c48b4f3c7b2af38249))

- cbus:
  - small adjustments on reset and optional alpha on direct command ([cbce5d2](https://github.com/apache/plc4x/commit/cbce5d212b9ee7e0a335ef372c4982f74c20f55f))
  - followup on improved message output ([d3f7132](https://github.com/apache/plc4x/commit/d3f71328325b5500262f87c1ffe1245363107fbc))
  - add virtual fields for hex encoded fields ([3626805](https://github.com/apache/plc4x/commit/362680547e039b16301e0995439b7112065cd12f))
  - removed unneeded types ([6e5abeb](https://github.com/apache/plc4x/commit/6e5abebda652af26103410f15d06cba196b9e4ac))
  - streamline network routing ([b93e7bf](https://github.com/apache/plc4x/commit/b93e7bffc520d527b4e59c6efee0ff563e40bc39))
  - calculate the remaining bytes instead of supplying them via argument ([37a6f63](https://github.com/apache/plc4x/commit/37a6f63d67ab0374c6f18c1b6bc3cfe8d4b12e40))
  - split up applications in sub mspecs ([b7fb2d6](https://github.com/apache/plc4x/commit/b7fb2d633b907eae62b1ac9f4fb97d6c6efeb164))
  - small cleanups and added a test for an unknown command ([3d6fbd0](https://github.com/apache/plc4x/commit/3d6fbd0f60f8eae636c96a11f3114014acb0634b))
  - get rid of CALDataOrSetParameter ([0304928](https://github.com/apache/plc4x/commit/03049281171e44fbaf103cf0247ca989b7e479da))
  - SALData is more than Lightning so the other types received a dummy implementation ([9ae9f29](https://github.com/apache/plc4x/commit/9ae9f293234fefd38c9319d76ca79d1b18738521))

- plc4xbrowser:
  - put ui parts into a ui package ([e8ab0a7](https://github.com/apache/plc4x/commit/e8ab0a7ed589db8cab7bc85f5a00fd3f0d6a453d))
  - split up main into ui and actions ([8b9da2a](https://github.com/apache/plc4x/commit/8b9da2abdc0d1766c2da9e1416f738ccd7d88044))
  - reworked the REPL command system ([6858376](https://github.com/apache/plc4x/commit/6858376df9571a5c8c131bab7096d81e74146f8e))

- codegen/plc4go:
  - remove unnecessary cast on - operator ([7291746](https://github.com/apache/plc4x/commit/7291746b1751026ad2ca3589977e81edca44cefa))
  - remove unnecessary cast on () operator ([2239dba](https://github.com/apache/plc4x/commit/2239dbace89ba322a99a7cb025695d458548ddb2))
  - switched enum back to value with a log statement as an error break more for now... ([c684d6c](https://github.com/apache/plc4x/commit/c684d6c4217eecb587bd32b091c4c1dab2429c5b))
  - if an enum is unknown return a ParseAssertError instead of a generic error ([2ba1332](https://github.com/apache/plc4x/commit/2ba1332921c2f9f8c3daf31adaeacdc4c4f1e8a3))
  - reworked enum parsing to include a flag if it is a unknown enum ([2f17300](https://github.com/apache/plc4x/commit/2f1730023fce7e12a81b7f87293018c5937a7da4))
  - initialize parent with parameter ([b84e75e](https://github.com/apache/plc4x/commit/b84e75ed8749ae02789fc940b7a93bb080182581))
  - avoid duplicating arguments as fields when present in parent ([67a1f6f](https://github.com/apache/plc4x/commit/67a1f6f2356003d95437387dbe93b129e2d6f2ae))
  - set parser argument field in child which shadows parent ([3fc44f9](https://github.com/apache/plc4x/commit/3fc44f950802fa7f65773960a435b3f677469886))
  - set arrays to nil if the result in an empty one ([3f0ae43](https://github.com/apache/plc4x/commit/3f0ae4368fde7e3087963878bbad21a9b15a6427))
  - ${type.name}Exactly interface extends now the type interfaces for ease of use ([5b372b0](https://github.com/apache/plc4x/commit/5b372b0ea19e7046733d38a77c7ad0ee9262ae81))
  - introduced a ${type.name}Exactly interface for exact type matches ([58f08fc](https://github.com/apache/plc4x/commit/58f08fcc7087cd8760e25904548886de008a9a77))
  - moved getLength methods in own interfaces and inherit Serialize from Serializable ([73baa76](https://github.com/apache/plc4x/commit/73baa76ec94585ce635f12bac3972151751d777b))
  - major refactoring how types are handled in plc4go ([e7fb792](https://github.com/apache/plc4x/commit/e7fb7923209a3ab80852dac6c01aca5b85a38618))

- plc4j/spi:
  - small optimization using the right datatype ([3a131f2](https://github.com/apache/plc4x/commit/3a131f2daf8ca42787f550fa87f1b36406bb7f8e))
  - minor cleanup ([5aca709](https://github.com/apache/plc4x/commit/5aca70984adfbcb5faedb1dc3e6b3af2d1b8748f))
  - use woodstox as stax ([4404a90](https://github.com/apache/plc4x/commit/4404a902d50912665fd0fa21497e0a48095232f2))

- plc4go/s7:
  - cleanup s7 code ([955055a](https://github.com/apache/plc4x/commit/955055aea56e037eb7564c5eee2698ad43ddeaac))

- plc4j/codegen:
  - moved logging from generated enums to DataReaderEnumDefault ([9b0eb6e](https://github.com/apache/plc4x/commit/9b0eb6eaaf51a0dd682766d70e38fbbfcc83d0c4))
  - remove unused generic qualifier ([5d082a5](https://github.com/apache/plc4x/commit/5d082a59acb05ee56bf7e3c13fa9fe6d15a30659))
  - remove IO classes ([be88a5e](https://github.com/apache/plc4x/commit/be88a5e543debab16aafe8ca2c9addcd51bd0f38))
  - move parser to model classes ([b08eb50](https://github.com/apache/plc4x/commit/b08eb50f458304448191bd7ad4a5a0c0392e191e))
  - cleanup field reader/writer ([93af95a](https://github.com/apache/plc4x/commit/93af95a6ff9224b51f6050f4511d1bd2edb85a76))

- bacnet:
  - DRY removed owner names from type switched ([31b5c38](https://github.com/apache/plc4x/commit/31b5c388960f90f725e02668d74d30e208c150f2))
  - fix who-has ([a746de1](https://github.com/apache/plc4x/commit/a746de1d140a240a1553cc0cd25eff5ad2d52c8f))
  - rename BACnetConstructedDataSecurityPduTimeout to BACnetConstructedDataSecurityPDUTimeout ([e611fd1](https://github.com/apache/plc4x/commit/e611fd13d96cd281a4458f8f0d80c44eea0e320a))
  - relax useage of BACnetPriorityArray ([0d1ace4](https://github.com/apache/plc4x/commit/0d1ace47f99964ced2a1c64f9acb5887b93e59ef))
  - switch BACnetTagPayloadOctetString to byte[] ([12a3a23](https://github.com/apache/plc4x/commit/12a3a233b3a2b24e584f512f6af75906e1f4c038))
  - externalized opening and closing tag ([080e2b9](https://github.com/apache/plc4x/commit/080e2b9e1a1ef4674e49524b99eace1404aed514))
  - major rework the way enums are handled ([d0255f4](https://github.com/apache/plc4x/commit/d0255f42889ad152af4c304caf2284a97a1c729c))
  - moved duplicate implementations into payload ([ffd996b](https://github.com/apache/plc4x/commit/ffd996b9a80dcdde6c3284bae34238d27975ae40))
  - moved catch all data block into BACnetConstructedDataUnspecified ([bddf87d](https://github.com/apache/plc4x/commit/bddf87d2d5915a090435b6c83c4364502dcf6c07))
  - removed unused virtual field ([a488906](https://github.com/apache/plc4x/commit/a488906aad57ea2e0dbddc1a299acb03fc71d9d0))
  - move bac net header parsing to a common BACnetTagHeader type ([f6b5502](https://github.com/apache/plc4x/commit/f6b550207f2130cb7509a68ec953b63a723cf7a6))
  - small optimizations ([bd71859](https://github.com/apache/plc4x/commit/bd718596b826ef8591c5795bf2fa3ba20dcd489f))
  - small improvements ([deb58b3](https://github.com/apache/plc4x/commit/deb58b3ad7d592a468b61e2f794715823fcc4ef1))
  - re-arrange some value ([e5aef05](https://github.com/apache/plc4x/commit/e5aef0517b6e5d7985109a7a82e235d2d7f036b7))

- plc4go/bacnet:
  - refined logging for Discoverer ([e016585](https://github.com/apache/plc4x/commit/e01658516953d3d826a93ee3b69856956f754507))

- plc4go/spi:
  - use spi.Message instead of interface{} to clean up intention ([bed9aa6](https://github.com/apache/plc4x/commit/bed9aa6bea588ae074094a184d9a92a84cd678d7))
  - Added a GetString implementation for RawPlcValue types ([835136a](https://github.com/apache/plc4x/commit/835136a3a168fd7485db6cbe026bc7bcf292637b))

- ab-eth:
  - don't use len as field name ([d45cbf3](https://github.com/apache/plc4x/commit/d45cbf3990d43834bfe617271ef61703453bc42a))

- eip:
  - don't use len as field name ([b2b8feb](https://github.com/apache/plc4x/commit/b2b8feb27e5dbc005c4b53efc75ce7be6855211f))

- protocols/knx:
  - Changed the KnxDatatypes to use BitStrings for BYTE, WORD, DWORD and LWORD ([57d5c1e](https://github.com/apache/plc4x/commit/57d5c1e80592d36485e1d43a7ee1e0a3bf84c587))

- plc4go/codegen:
  - flipped new functions of child to return the actual child ([2b61308](https://github.com/apache/plc4x/commit/2b6130872f7914d56ad17d88c1e328efe9e4328d))
  - flipped parse functions of child to return the actual child ([ca65d6c](https://github.com/apache/plc4x/commit/ca65d6cc66e4d67b5bebca035a0151de2be548aa))
  - add GetParent() retriever method ([1a234ed](https://github.com/apache/plc4x/commit/1a234ed7941e4b5d330dbf898569bcd9223d0ab0))
  - optimized code sections ([9d38eaf](https://github.com/apache/plc4x/commit/9d38eafcb54cab4bb7b4b922b6aac1aab5461a1c))
  - removed superfluous methods ([e413f26](https://github.com/apache/plc4x/commit/e413f2682915457b3300c0d3cba9883040b97555))
  - added string rendering error handling ([36942cf](https://github.com/apache/plc4x/commit/36942cf5b85542cccefb53f372de3e326a910d29))
  - simplified cast function ([fd8c81c](https://github.com/apache/plc4x/commit/fd8c81c888062f76a58415ad686bcdf3545aacc2))
  - add parent interface to child interface ([42b529f](https://github.com/apache/plc4x/commit/42b529f5bee9126b9035f48e149f30c3fb80dd56))
  - fixed virtual field implementation ([e02e64d](https://github.com/apache/plc4x/commit/e02e64d518b164e214fa7e6c753f1436a73f2262))

- general:
  - slightly adjust kotlin example (fix) ([394d759](https://github.com/apache/plc4x/commit/394d759bf84cb5c701038de7a100ade5314a6caa))
  - slightly adjust kotlin example ([2d3fd7c](https://github.com/apache/plc4x/commit/2d3fd7c1e4e4c0d276dc9510a3b79df75bed2ee1))
  - Changed the artifact-id of the test-generator ([77408e6](https://github.com/apache/plc4x/commit/77408e64fb7d0470652f82dd5489c86e79f0c9c6))
  - refactoring(plc4j): split up the raw socket transport into a "raw" and "raw-passive" transport, where the raw-passive is equivalent to the previous raw transport. ([27442e6](https://github.com/apache/plc4x/commit/27442e65e1619aaa358768fc985219a7020c70e7))
  - refactoring: Moved the plc4net module outside of the sandbox ([e8f89cb](https://github.com/apache/plc4x/commit/e8f89cb8513e683b6445336590bed66cc12507a2))
  - refactoring: Updated the plc4net branch to the latest changes on develop.
new feature: Now the DataIo generation is almost finished. ([2775ec9](https://github.com/apache/plc4x/commit/2775ec998eaba81764a87295e95201b80240ac26))
  - refactoring(documentation): Started updating the mspec documentation ([1a0e3e9](https://github.com/apache/plc4x/commit/1a0e3e9f44e21b1e583d73c2324a1ba97d2f9f0e))
  - refactor (documentation): Started updating the code-generation documentation. ([ec3fc42](https://github.com/apache/plc4x/commit/ec3fc42e7fdd26ecf1743270f36ae197f8ae7074))

- test-utils:
  - Made the ParserSerializerTestsuiteRunner play a bit more nicely with not fully implemented tests. ([f170c59](https://github.com/apache/plc4x/commit/f170c59fae984a92de345481633710f4a04d9152))

- code-generation:
  - Changed the implementation of getSizeInBytes to be a bit more correct. ([2ce32ef](https://github.com/apache/plc4x/commit/2ce32ef2d18d91e175617cab2c1b0ca43cc51c94))
  - Finished migrating the code to using the static getSizeInBits method of DataIo types ([447ef68](https://github.com/apache/plc4x/commit/447ef680f274d779b0777eb7c2538b7dc5159300))
  - Made the code generation work with DataIo types referenced from normal model types. ([68a9627](https://github.com/apache/plc4x/commit/68a96272d7d6082b801f8e3d6d20e86cdb4e06a8))
  - Made the code generation work with DataIo types referenced from normal model types. ([9ef24e1](https://github.com/apache/plc4x/commit/9ef24e1bba3e478c8c2a51ab59814d0bc0167e92))

- plc4j/profinet:
  - Refactored the profinet discovery to the latest changes in the protocol spec and added some more features. ([ae9bb4a](https://github.com/apache/plc4x/commit/ae9bb4abbceceba49db81022ea2c5e01948e2169))
  - Disabled the parser serializer test for now (as it's not yet supposed to work) ([083dce1](https://github.com/apache/plc4x/commit/083dce1d9db42570c499869e3421b0acc5c3a659))

- profinet:
  - Continued refactoring the PNIO Messages ([c3b1c0c](https://github.com/apache/plc4x/commit/c3b1c0c05477a1ab58b61f37528280dbcb396087))

- plc4x-protocol:
  - Added a length field to the plc4x protocol messages ([d871063](https://github.com/apache/plc4x/commit/d8710630231e32f09f583ab7f8b449bbaf6fe8b8))

- plc4c:
  - Changed the size calculation of manual fields from bytes to bits ([2057eb8](https://github.com/apache/plc4x/commit/2057eb81727cb628799961478091d78fa569e2cf))

- plc4j:
  - Changed the code-generation to generate the dataIo types into the same directory as the other types and removed the "IO" suffix to match the rest. ([af79344](https://github.com/apache/plc4x/commit/af793448bd3e74f643a91ee27bacddccf9d74a71))

- plc4go/drivers:
  - Made all driver's Connections implement the spi.PlcConnection interface (so we can use them in the connection pool). ([20be7b2](https://github.com/apache/plc4x/commit/20be7b27f3762d2a20d045c4948930dda73de41c))

- plc4j/writeBufferByteBase:
  - deprecate getData() method ([b753e27](https://github.com/apache/plc4x/commit/b753e270908cf9ac96355e52226c618f7f1d6fac))

- plc4j/codgen:
  - cleanup list handling/revert list entries to anonymous logical name ([a2b36a4](https://github.com/apache/plc4x/commit/a2b36a44f2fa163e9e88b9d5a5255fec85a69ac2))
  - small cleanups ([e1266ae](https://github.com/apache/plc4x/commit/e1266ae03bb8cce1d664564d076dbcd91e50c2e4))
  - cleanup list handling/revert list entries to anonymous logical name ([02e102c](https://github.com/apache/plc4x/commit/02e102c638d23e0492ab1a479eb00d5b26f41b67))
  - small cleanups ([ad6fe78](https://github.com/apache/plc4x/commit/ad6fe789fdf270c18b446d8806f289f10a7ad484))

- code-gen:
  - cleanup interfaces ([8499434](https://github.com/apache/plc4x/commit/84994340c026bb865a980004ba0ffe0d28c1add4))

## [rel/0.9](https://github.com/apache/plc4x/releases/tag/rel/0.9) - 2021-09-17 09:25:12

## What's Changed
* Added Pool2 (now called Connection Cache) to develop Branch by @JulianFeinauer in https://github.com/apache/plc4x/pull/217
* Some doc about the Julian plc4x-pool2 by @foxpluto in https://github.com/apache/plc4x/pull/209
* PLC4X-207 When a Handler Timeout occurs cancel the read future to not… by @JulianFeinauer in https://github.com/apache/plc4x/pull/170
* PLC4X-265 Support for hex mapping of byte/byte[] values in XML test by @splatch in https://github.com/apache/plc4x/pull/213
* added strtok_s define for windows by @thomas169 in https://github.com/apache/plc4x/pull/228
* Fix for opcua subscription by @hutcheb in https://github.com/apache/plc4x/pull/227
* Feature/string enum mspec by @hutcheb in https://github.com/apache/plc4x/pull/230
* feat: Add read/write support for string types within the Java dataio classes by @hutcheb in https://github.com/apache/plc4x/pull/235
* Update S7Field.java by @Meng5 in https://github.com/apache/plc4x/pull/239
* regress: removed string support for modbus by @hutcheb in https://github.com/apache/plc4x/pull/236
* Abstract field's getter doesn't get get declared in base type by @hutcheb in https://github.com/apache/plc4x/pull/240
* Remove merge button when merging PR's, this forces the use of squash by @hutcheb in https://github.com/apache/plc4x/pull/238
* plc4go: fixed passing parameters incorrectly, resulting in a null pointer by @hongjinlin in https://github.com/apache/plc4x/pull/243
* Address string with tokenisation, refactoring, loopback s7 c example added. s7c write variable works by @thomas169 in https://github.com/apache/plc4x/pull/233
* plc4c: memory plumbing by @thomas169 in https://github.com/apache/plc4x/pull/244
* [PLC4X-299] Fix for array handling in the Kafka source connector by @hutcheb in https://github.com/apache/plc4x/pull/255
* Fix link : plc4j getting start, graphviz by @shblue21 in https://github.com/apache/plc4x/pull/256
* Feature/native opua client by @hutcheb in https://github.com/apache/plc4x/pull/253
* PLC4X-307 Add support for custom generator options / customized packages by @splatch in https://github.com/apache/plc4x/pull/263
* fixing timer leaks by @shaunco in https://github.com/apache/plc4x/pull/267
* S7event by @glcj in https://github.com/apache/plc4x/pull/264
* changed delimiter of opc ua from : to ; by @nalim2 in https://github.com/apache/plc4x/pull/258
* plc4x: major dependency updates by @sruehl in https://github.com/apache/plc4x/pull/252
* build(deps): bump karaf-maven-plugin from 4.3.0 to 4.3.2 by @dependabot in https://github.com/apache/plc4x/pull/245
* Bump zookeeper from 3.4.13 to 3.4.14 by @dependabot in https://github.com/apache/plc4x/pull/270
* Bump milo.version from 0.6.1 to 0.6.3 by @dependabot in https://github.com/apache/plc4x/pull/269
* build(deps): bump gmavenplus-plugin from 1.6.2 to 1.12.1 by @dependabot in https://github.com/apache/plc4x/pull/246
* build(deps): bump junit.jupiter.version from 5.7.0 to 5.7.2 by @dependabot in https://github.com/apache/plc4x/pull/251
* build(deps): bump logback.version from 1.2.3 to 1.2.5 by @dependabot in https://github.com/apache/plc4x/pull/271
* Change PLC4GO docs read syntax by @NiklasMerz in https://github.com/apache/plc4x/pull/276
* PLC4X-294 / PLC4X-296 Reorganization of stack configurer / events support in public api by @splatch in https://github.com/apache/plc4x/pull/241

## New Contributors
* @thomas169 made their first contribution in https://github.com/apache/plc4x/pull/228
* @Meng5 made their first contribution in https://github.com/apache/plc4x/pull/239
* @shblue21 made their first contribution in https://github.com/apache/plc4x/pull/256
* @shaunco made their first contribution in https://github.com/apache/plc4x/pull/267
* @glcj made their first contribution in https://github.com/apache/plc4x/pull/264

**Full Changelog**: https://github.com/apache/plc4x/compare/rel/0.8...rel/0.9

### Bug Fixes

- general:
  - Updated junit.platform.version to match the updated jupiter version ([d767dc9](https://github.com/apache/plc4x/commit/d767dc9b9160ecdfc500be382415d05957af3f48))
  - PLC4X-312 - CAN NOT READ STRING FROM S7 PLC ([0e57493](https://github.com/apache/plc4x/commit/0e5749347def06d1f961a179695a5751653b2815)) ([#264](https://github.com/apache/plc4x/pull/264))
  - fixing timer leaks (#267) ([d059709](https://github.com/apache/plc4x/commit/d05970961a4dd92cd52e547a15ca5fdb3030a2ff)) ([#267](https://github.com/apache/plc4x/pull/267))
  - fix build: ([24e739e](https://github.com/apache/plc4x/commit/24e739e73c5ceb55785c8ff59444fee69813a649))
  - fix build: ([2ab5472](https://github.com/apache/plc4x/commit/2ab54726d05fde08b2fa23119943d4b6778756a0))
  - fixed memory leaks on read and write.
added a callback to clear io request (like  we had for responces)
renamed request destoy function to be in style with destory execution and responce
added lots of small delelte functions into s7_packets
big refactor of sm_read and sm_write
no leaks on IO now I think ([20ac37a](https://github.com/apache/plc4x/commit/20ac37a3267cfadb143c93f0b9d6285024fb919a)) ([#244](https://github.com/apache/plc4x/pull/244))
  - fix build ([557122b](https://github.com/apache/plc4x/commit/557122b93465065d76d353ddde0aed4b50701641))
  - Update the itemCount variable form uint 8 to uint 16 ([ef83443](https://github.com/apache/plc4x/commit/ef83443a4a1ab675c1cfea89b6daa6a6c98ac7e1))
  - Commented out the call to "install" which I incorrectly commented out
fix: Cleaned up in the pom structures ([868ef40](https://github.com/apache/plc4x/commit/868ef403e0747bfd974e523de61f75537d9f9925))
  - reverted ComplexTypeDefinition change on tmeplate ([568a8fc](https://github.com/apache/plc4x/commit/568a8fcf62ac8593e376bd0ee92b0337fdc55185)) ([#240](https://github.com/apache/plc4x/pull/240))
  - Made the NiFi modules use the PooledDriverManager more correctly. ([036a157](https://github.com/apache/plc4x/commit/036a157857dad552b51c56b60f7f622a964d6a6c))
  - fix for 32bit int writes, now works too ([c02a2d3](https://github.com/apache/plc4x/commit/c02a2d32d311ba8c233804b267a2d2fc503456a4)) ([#233](https://github.com/apache/plc4x/pull/233))
  - fix to offsetting of var paload and some exra fetures on loopback tests ([6d9c705](https://github.com/apache/plc4x/commit/6d9c705e6277f26c4ba42e2b8a57bc1a5cb16e5e)) ([#233](https://github.com/apache/plc4x/pull/233))
  - fix download link 0.7 ([fd4a8e4](https://github.com/apache/plc4x/commit/fd4a8e42ecbfc08f47ed9c550e54320e37a5429e))
  - fix logback dependency ([7cbd42b](https://github.com/apache/plc4x/commit/7cbd42b494ca83a120788b848d7ee017b05519c8)) ([#217](https://github.com/apache/plc4x/pull/217))

## [rel/0.8](https://github.com/apache/plc4x/releases/tag/rel/0.8) - 2021-01-25 11:31:05

## What's Changed
* Feature/plc simulator by @JulianFeinauer in https://github.com/apache/plc4x/pull/157
* [SITE] Add Logo Column. Add pragmatic minds and pragmatic industries … by @JulianFeinauer in https://github.com/apache/plc4x/pull/159
* have rat ignore .java-version files from jenv since we need java kung… by @ottobackwards in https://github.com/apache/plc4x/pull/160
* Updated Eclipse Milo to 0.3.7 by @patrickse in https://github.com/apache/plc4x/pull/158
* Updated the codecentric adopters page entry by @chrisdutz in https://github.com/apache/plc4x/pull/162
* C-Code-Generation (Please don't merge ... this is used for allowing code-reviews) by @chrisdutz in https://github.com/apache/plc4x/pull/161
* IndustryFusion as PLC4X adopter by @KonstantinKe in https://github.com/apache/plc4x/pull/165
* Added some test cases for OPC UA by @nalim2 in https://github.com/apache/plc4x/pull/164
* Adopter ISW addition by @nalim2 in https://github.com/apache/plc4x/pull/167
* Feature/c code generation tagged unions by @chrisdutz in https://github.com/apache/plc4x/pull/168
* - Started working on the parser and serializer code ... WIP by @chrisdutz in https://github.com/apache/plc4x/pull/163
* Team addition strljic by @nalim2 in https://github.com/apache/plc4x/pull/166
* [PLC4X-216]update IoTDB JDBC example and session API example; add the related doc on website by @jixuan1989 in https://github.com/apache/plc4x/pull/171
* Feature/modbus add additional address formats and change lowest register to 1. by @hutcheb in https://github.com/apache/plc4x/pull/172
* Fix documentation with updates to README by @ottlinger in https://github.com/apache/plc4x/pull/173
* See issue reported by Stefano Bossi. by @JulianFeinauer in https://github.com/apache/plc4x/pull/175
* Feature/extended register read by @hutcheb in https://github.com/apache/plc4x/pull/174
* PLC4X-223 Fix exception on Pool usage for all drivers due to syntax c… by @JulianFeinauer in https://github.com/apache/plc4x/pull/176
* Use Gradle with compiled library by @foxpluto in https://github.com/apache/plc4x/pull/183
* PLC4X-244 Support for variable length padding fields. by @splatch in https://github.com/apache/plc4x/pull/182
* PLC4X-248 Permit more locations for MSpec comments. by @splatch in https://github.com/apache/plc4x/pull/185
* Modbus/update documentation by @hutcheb in https://github.com/apache/plc4x/pull/184
* Feature/plc4c by @chrisdutz in https://github.com/apache/plc4x/pull/181
* Feature/write extended registers by @hutcheb in https://github.com/apache/plc4x/pull/190
* Modbus Data Type Documentation by @hutcheb in https://github.com/apache/plc4x/pull/188
* Refactor Field Handler Classes by @hutcheb in https://github.com/apache/plc4x/pull/192
* Add an OPC UA Server by @hutcheb in https://github.com/apache/plc4x/pull/194
* Beckhoff Device Discovery structures by @splatch in https://github.com/apache/plc4x/pull/195
* PLC4X-252 Support passing of PlcField in read/write builders. by @splatch in https://github.com/apache/plc4x/pull/198
* Opcaua server updates by @hutcheb in https://github.com/apache/plc4x/pull/196
* Record Evolution text adjustment by @markope in https://github.com/apache/plc4x/pull/199
* PLC4X-255 Fix for Kafka Source CPU Usage by @hutcheb in https://github.com/apache/plc4x/pull/200
* Bug/simulated device data types by @hutcheb in https://github.com/apache/plc4x/pull/204
* Update/GitHub metadata by @hutcheb in https://github.com/apache/plc4x/pull/205
* Test ASF file by @hutcheb in https://github.com/apache/plc4x/pull/206
* Feature/kafkasink - Add a kafka sink  by @hutcheb in https://github.com/apache/plc4x/pull/202
* Fix for connection pool, unable to reconnect after failed connection. by @hutcheb in https://github.com/apache/plc4x/pull/207
* Update/opcua server by @hutcheb in https://github.com/apache/plc4x/pull/208
* Fix ToAnsi method : bufferoverflow by @JohannaMillet in https://github.com/apache/plc4x/pull/178
* poll loop fetches data in a continous loop by @gg587998 in https://github.com/apache/plc4x/pull/67
* CANopen over socketcan transport by @splatch in https://github.com/apache/plc4x/pull/211
* Add ConnectorIO as Apache PLC4X adopter by @splatch in https://github.com/apache/plc4x/pull/212
* Bug/driver opcua read arrays [PLC4X-202] by @hutcheb in https://github.com/apache/plc4x/pull/210
* Add enum's for opc ua driver in mspec by @hutcheb in https://github.com/apache/plc4x/pull/214
* Bug/close connection by @hutcheb in https://github.com/apache/plc4x/pull/221
* Add confluent package to deployment so it appears in Nexus. by @hutcheb in https://github.com/apache/plc4x/pull/222
* Updated Notice and License files. by @hutcheb in https://github.com/apache/plc4x/pull/223
* Add encryption handler for OPC UA - Minor fix for Kafka Connector by @hutcheb in https://github.com/apache/plc4x/pull/225

## New Contributors
* @patrickse made their first contribution in https://github.com/apache/plc4x/pull/158
* @KonstantinKe made their first contribution in https://github.com/apache/plc4x/pull/165
* @ottlinger made their first contribution in https://github.com/apache/plc4x/pull/173
* @markope made their first contribution in https://github.com/apache/plc4x/pull/199
* @JohannaMillet made their first contribution in https://github.com/apache/plc4x/pull/178

**Full Changelog**: https://github.com/apache/plc4x/compare/rel/0.7...rel/0.8

## [rel/0.7](https://github.com/apache/plc4x/releases/tag/rel/0.7) - 2020-05-15 10:51:06

## What's Changed
* PLC4X-146 - Problem when using examples/hello-integration-edgent with… by @chrisdutz in https://github.com/apache/plc4x/pull/93
* Fixed a problem causing the karaf feature module not to sign the feat… by @chrisdutz in https://github.com/apache/plc4x/pull/95
* Fixed a problem with running the elasticsearch example as updating el… by @chrisdutz in https://github.com/apache/plc4x/pull/94
* - Moved the plc4cpp, plc4net and plc4py into the sandbox by @chrisdutz in https://github.com/apache/plc4x/pull/98
* add iotdb integration example. for just saving one field for a PLC. by @jixuan1989 in https://github.com/apache/plc4x/pull/99
* Feature/reproducible builds by @chrisdutz in https://github.com/apache/plc4x/pull/100
* Bump org.eclipse.paho.client.mqttv3 from 1.2.0 to 1.2.1 in /plc4j/examples/hello-cloud-google by @dependabot in https://github.com/apache/plc4x/pull/96
* Bump jackson.version from 2.9.9 to 2.10.0 by @dependabot in https://github.com/apache/plc4x/pull/97
* Fix streampipes adapters by @tenthe in https://github.com/apache/plc4x/pull/101
* Ab eth station address parameter by @vemmert in https://github.com/apache/plc4x/pull/102
* PLC4X-160: Fix css to make links to Apache events and ASF homepage work again by @dominikriemer in https://github.com/apache/plc4x/pull/104
* Feature/ams ads mpsec by @sruehl in https://github.com/apache/plc4x/pull/107
* Plc4 x 157 opc ua disablediscovery by @nalim2 in https://github.com/apache/plc4x/pull/105
* PLC4X-45 Add float support to Modbus protocol by @acs in https://github.com/apache/plc4x/pull/109
* Fix future chain for InternalPlcWriteRequest by @amrod- in https://github.com/apache/plc4x/pull/110
* improve comment matching by @sruehl in https://github.com/apache/plc4x/pull/111
* PLC4X-164: Fix wrong NOT FOUND exception in OPC-UA driver by @acs in https://github.com/apache/plc4x/pull/112
* Fixes in ADS protocol by @amrod- in https://github.com/apache/plc4x/pull/103
* Improve carousel for small and medium screen sizes by @dominikriemer in https://github.com/apache/plc4x/pull/114
* Introduced a new optional "slaveId" parameter to the modbus driver by @acs in https://github.com/apache/plc4x/pull/118
* [fix] ported to new base version by @nalim2 in https://github.com/apache/plc4x/pull/120
* PLC4X-176 resolve issues round running on linux platforms by @ottobackwards in https://github.com/apache/plc4x/pull/121
* Fix typo he -> we by @turbaszek in https://github.com/apache/plc4x/pull/126
* Update link for PLC4J in README by @turbaszek in https://github.com/apache/plc4x/pull/125
* Use existing base image in Dockerfile by @turbaszek in https://github.com/apache/plc4x/pull/123
* Change (P)PMC to PMC on team page by @NiklasMerz in https://github.com/apache/plc4x/pull/122
* [WIP] Refactor PlcDriverManager as context manager by @turbaszek in https://github.com/apache/plc4x/pull/124
* [fix] Transfer bug fixes from rel/0.6 to 0.7+ which includes correct … by @nalim2 in https://github.com/apache/plc4x/pull/131
* Change repo URLs from incubator-plc4x to plc4x by @NiklasMerz in https://github.com/apache/plc4x/pull/139
* remove Bnd-LastModified header that is not reproducible by @hboutemy in https://github.com/apache/plc4x/pull/137
* move structure defs into private header by @ottobackwards in https://github.com/apache/plc4x/pull/142
* c and clion git ignore by @ottobackwards in https://github.com/apache/plc4x/pull/144
* flatten includes, remove addr sample, use <> includes by @ottobackwards in https://github.com/apache/plc4x/pull/145
* PLC4X-192 Support for connection string parameter conversion. by @splatch in https://github.com/apache/plc4x/pull/140
* C api m4 feedback by @ottobackwards in https://github.com/apache/plc4x/pull/146
* add plc4c_data type by @ottobackwards in https://github.com/apache/plc4x/pull/147
* integrate writes with plc_data by @ottobackwards in https://github.com/apache/plc4x/pull/149
* complex object members should not leak to external entities, add apis… by @ottobackwards in https://github.com/apache/plc4x/pull/150
* add CODE_CONVENTIONS by @ottobackwards in https://github.com/apache/plc4x/pull/151
* More work on access to _t internals by @ottobackwards in https://github.com/apache/plc4x/pull/152
* Feature/c api subscriptions by @chrisdutz in https://github.com/apache/plc4x/pull/153
* Feature/osgi by @etiennerobinet in https://github.com/apache/plc4x/pull/154
* PLC4X-197 do not stop the global timer, remove and cancel the Timeouts instead by @ottobackwards in https://github.com/apache/plc4x/pull/155

## New Contributors
* @tenthe made their first contribution in https://github.com/apache/plc4x/pull/101
* @amrod- made their first contribution in https://github.com/apache/plc4x/pull/110
* @turbaszek made their first contribution in https://github.com/apache/plc4x/pull/126
* @hboutemy made their first contribution in https://github.com/apache/plc4x/pull/137
* @etiennerobinet made their first contribution in https://github.com/apache/plc4x/pull/154

**Full Changelog**: https://github.com/apache/plc4x/compare/rel/0.5...rel/0.7

## [rel/0.6](https://github.com/apache/plc4x/releases/tag/rel/0.6) - 2020-01-13 14:49:38

## What's Changed
* PLC4X-146 - Problem when using examples/hello-integration-edgent with… by @chrisdutz in https://github.com/apache/plc4x/pull/93
* Fixed a problem causing the karaf feature module not to sign the feat… by @chrisdutz in https://github.com/apache/plc4x/pull/95
* Fixed a problem with running the elasticsearch example as updating el… by @chrisdutz in https://github.com/apache/plc4x/pull/94
* - Moved the plc4cpp, plc4net and plc4py into the sandbox by @chrisdutz in https://github.com/apache/plc4x/pull/98
* add iotdb integration example. for just saving one field for a PLC. by @jixuan1989 in https://github.com/apache/plc4x/pull/99
* Feature/reproducible builds by @chrisdutz in https://github.com/apache/plc4x/pull/100
* Bump org.eclipse.paho.client.mqttv3 from 1.2.0 to 1.2.1 in /plc4j/examples/hello-cloud-google by @dependabot in https://github.com/apache/plc4x/pull/96
* Bump jackson.version from 2.9.9 to 2.10.0 by @dependabot in https://github.com/apache/plc4x/pull/97
* Fix streampipes adapters by @tenthe in https://github.com/apache/plc4x/pull/101
* Ab eth station address parameter by @vemmert in https://github.com/apache/plc4x/pull/102
* PLC4X-160: Fix css to make links to Apache events and ASF homepage work again by @dominikriemer in https://github.com/apache/plc4x/pull/104
* Feature/ams ads mpsec by @sruehl in https://github.com/apache/plc4x/pull/107
* Plc4 x 157 opc ua disablediscovery by @nalim2 in https://github.com/apache/plc4x/pull/105
* PLC4X-45 Add float support to Modbus protocol by @acs in https://github.com/apache/plc4x/pull/109
* Fix future chain for InternalPlcWriteRequest by @amrod- in https://github.com/apache/plc4x/pull/110
* improve comment matching by @sruehl in https://github.com/apache/plc4x/pull/111
* PLC4X-164: Fix wrong NOT FOUND exception in OPC-UA driver by @acs in https://github.com/apache/plc4x/pull/112
* Fixes in ADS protocol by @amrod- in https://github.com/apache/plc4x/pull/103
* PLC4X-163 Fixed Netty Buffer Leaks for S7 Driver. by @JulianFeinauer in https://github.com/apache/plc4x/pull/116
* [PLC4X-168] Shorter S7 Field Syntax by @JulianFeinauer in https://github.com/apache/plc4x/pull/113
* PLC4X-158 Add Warning if Scraper is not used with Pooled Driver Manager. by @JulianFeinauer in https://github.com/apache/plc4x/pull/115
* - Introduced a new optional "slaveId" parameter to the modbus driver by @acs in https://github.com/apache/plc4x/pull/127
* [fix] wrong Parameter handling by @nalim2 in https://github.com/apache/plc4x/pull/128
* PLC4x-142 OPC-UA read variables by @nalim2 in https://github.com/apache/plc4x/pull/130
* Feature/plc4 x 185 cert support opc ua by @JulianFeinauer in https://github.com/apache/plc4x/pull/132
* Feature/port simotion by @vemmert in https://github.com/apache/plc4x/pull/197
* Update license and notice file for derived works by @hutcheb in https://github.com/apache/plc4x/pull/224

## New Contributors
* @tenthe made their first contribution in https://github.com/apache/plc4x/pull/101
* @amrod- made their first contribution in https://github.com/apache/plc4x/pull/110

**Full Changelog**: https://github.com/apache/plc4x/compare/rel/0.5...rel/0.6

### Bug Fixes

- general:
  - fixed Release Notes ([dd74b34](https://github.com/apache/plc4x/commit/dd74b34b38988c56fe721b70afb631f14f571bed))
  - fixed retrieval of lists ([dab9308](https://github.com/apache/plc4x/commit/dab93088cadf02bbcb32cdd7639cd14b1a40ea1c))
  - fixed wrong equal check for discriminator in ads spec ([a6e4880](https://github.com/apache/plc4x/commit/a6e48805c941d892fde02c54753d91a1a5740547))
  - fixed wrong expression in ams mspec ([2921627](https://github.com/apache/plc4x/commit/29216274c7869d709c60654094b77a4049964de7))
  - fixed missing sandbox ams code ([408054c](https://github.com/apache/plc4x/commit/408054cb857a65a85c081410323771629c11c8f2))
  - fix build ([a7ca2bb](https://github.com/apache/plc4x/commit/a7ca2bbca2fe9efbe2a4a6ebb5694edac68113fd))
  - fixed generation for big integer ([0ee09fa](https://github.com/apache/plc4x/commit/0ee09faaec8c93121e41bc097ca9de394f82113a)) ([#107](https://github.com/apache/plc4x/pull/107))

### Refactor

- general:
  - refactored ads mspec to use proper names ([329b639](https://github.com/apache/plc4x/commit/329b63996b0e2484b7c6081534f8905fd0a6bb37))

## [rel/0.5](https://github.com/apache/plc4x/releases/tag/rel/0.5) - 2019-10-21 15:00:13

## What's Changed
* PLC4X-129 Improve version detection macro. by @splatch in https://github.com/apache/plc4x/pull/71
* fixed a typo on the plc4py/README.md by @martinLim45 in https://github.com/apache/plc4x/pull/74
* Fixed a typo on the index page by @jbarop in https://github.com/apache/plc4x/pull/73
* fix bug in opc ua driver response handling when fields are not found by @stefah in https://github.com/apache/plc4x/pull/75
* PLC4X-139 close the worker thread on connection abortion to avoid thr… by @JulianFeinauer in https://github.com/apache/plc4x/pull/76
* Feature/improve scraper tim by @JulianFeinauer in https://github.com/apache/plc4x/pull/77
* Example MSpec for the DF1 protocol added to the website by @vemmert in https://github.com/apache/plc4x/pull/79
* PLC4X-86 - Fix and re-enable tests that were disabled for Java 11 support by @thomasdarimont in https://github.com/apache/plc4x/pull/78
* A few fixes required to build Elastic example in a containerized environment by @rvs in https://github.com/apache/plc4x/pull/80
* Feature/implement logstash integration by @till1993 in https://github.com/apache/plc4x/pull/82
* Simplifying inheritance of shade plugin configurations by @rvs in https://github.com/apache/plc4x/pull/81
* Optimizing JRE for size and adding dynamic option setting for example jar by @rvs in https://github.com/apache/plc4x/pull/83
* [fixed] OPC UA Subscription item name is null by @nalim2 in https://github.com/apache/plc4x/pull/84
* Feature/logstash plugin improved fields by @stefah in https://github.com/apache/plc4x/pull/85
* Adding OSGi processing to the pom for java modules by @cschneider in https://github.com/apache/plc4x/pull/87
* Ab eth data types and doc by @vemmert in https://github.com/apache/plc4x/pull/86
* Migrate README to Github flavored markdown page by @dominikriemer in https://github.com/apache/plc4x/pull/88
* Update bit-io to newest version with OSGi support by @cschneider in https://github.com/apache/plc4x/pull/90
* Feature/resolve split package osgi by @JulianFeinauer in https://github.com/apache/plc4x/pull/89
* PLC4X-144 - When requesting invalid addresses, the DefaultS7MessagePr… by @chrisdutz in https://github.com/apache/plc4x/pull/91
* Update the Thrift version to the just released 0.13.0 to make it compatible with Boost again (Also bumped the Boost version to the latest 1.71.0) by @chrisdutz in https://github.com/apache/plc4x/pull/92

## New Contributors
* @martinLim45 made their first contribution in https://github.com/apache/plc4x/pull/74
* @jbarop made their first contribution in https://github.com/apache/plc4x/pull/73
* @stefah made their first contribution in https://github.com/apache/plc4x/pull/75
* @thomasdarimont made their first contribution in https://github.com/apache/plc4x/pull/78
* @rvs made their first contribution in https://github.com/apache/plc4x/pull/80
* @till1993 made their first contribution in https://github.com/apache/plc4x/pull/82
* @cschneider made their first contribution in https://github.com/apache/plc4x/pull/87

**Full Changelog**: https://github.com/apache/plc4x/compare/rel/0.4...rel/0.5

### Bug Fixes

- general:
  - fixed plugin build ([bf4ca89](https://github.com/apache/plc4x/commit/bf4ca89051bc578457f79660f4eeb7a98bd200e4))
  - fix bug in opc ua driver response handling when fields are not found ([da0b74e](https://github.com/apache/plc4x/commit/da0b74e662062aea63d77357efaedfc9e323c0a1)) ([#75](https://github.com/apache/plc4x/pull/75))
  - fixed a typo on the plc4py/README.md ([efec1a9](https://github.com/apache/plc4x/commit/efec1a91dd0e7b8599e9a5154bd8ee2a48e2a251)) ([#74](https://github.com/apache/plc4x/pull/74))

## [rel/0.4](https://github.com/apache/plc4x/releases/tag/rel/0.4) - 2019-05-20 21:09:32

## What's Changed
* Minor changes on getting-involved.adoc by @cptblaubaer in https://github.com/apache/plc4x/pull/45
* PLC4X-88: Add Triggering to PLC Scraper by @timbo2k in https://github.com/apache/plc4x/pull/46
* Extend field with type information by @JulianFeinauer in https://github.com/apache/plc4x/pull/47
* Correcting two typos in contributing.adoc by @cptblaubaer in https://github.com/apache/plc4x/pull/48
* INFRA-18056 by @clambertus in https://github.com/apache/plc4x/pull/49
* Preparing a little more complex PR for figuring out the Sync issues. by @chrisdutz in https://github.com/apache/plc4x/pull/50
* Added the content for Sebastian Rühl to the team page. by @chrisdutz in https://github.com/apache/plc4x/pull/51
* Yet another PR merge test by @chrisdutz in https://github.com/apache/plc4x/pull/52
* Moded the web-client unpacking to the maven prepare-package phase by @chrisdutz in https://github.com/apache/plc4x/pull/53
* Another  PR by @chrisdutz in https://github.com/apache/plc4x/pull/54
* Hopefully the last to get PRs working by @chrisdutz in https://github.com/apache/plc4x/pull/55
* Feature/plc4 x 108 ping method by @JulianFeinauer in https://github.com/apache/plc4x/pull/57
* Feature/alternate plc4 x 108 by @chrisdutz in https://github.com/apache/plc4x/pull/59
* Feature/plc4 x 111 interop server by @JulianFeinauer in https://github.com/apache/plc4x/pull/60
* Feature/plc4net by @bjoernhoeper in https://github.com/apache/plc4x/pull/61
* Fix link to banner by @NiklasMerz in https://github.com/apache/plc4x/pull/62
* Add FOSSA hackathon to event list by @NiklasMerz in https://github.com/apache/plc4x/pull/63
* [Fixed] DefaultPlcWriteRequest wrong return item for DefaultPlcWriteRequest.Builder addItem on Type Byte[] by @nalim2 in https://github.com/apache/plc4x/pull/65
* fixed typo in S7PlcConnection by @gg587998 in https://github.com/apache/plc4x/pull/69
* Feature/opcua by @nalim2 in https://github.com/apache/plc4x/pull/66

## New Contributors
* @cptblaubaer made their first contribution in https://github.com/apache/plc4x/pull/45
* @clambertus made their first contribution in https://github.com/apache/plc4x/pull/49
* @bjoernhoeper made their first contribution in https://github.com/apache/plc4x/pull/61

**Full Changelog**: https://github.com/apache/plc4x/compare/rel/0.3...rel/0.4

### Bug Fixes

- general:
  - fixed typo in S7PlcConnection ([7f94d6c](https://github.com/apache/plc4x/commit/7f94d6cc721eadc99d6115ed8cd6c9f040ac45ed)) ([#69](https://github.com/apache/plc4x/pull/69))
  - fixed a TriggeredScraper Bug ([67c4372](https://github.com/apache/plc4x/commit/67c4372581408248977fce5b24caaf137166b61f))
  - fixed one more sonar-bug ([c95f302](https://github.com/apache/plc4x/commit/c95f302a385a007986d63443d309639c13d5831b))
  - fixed sonar bugs ([de8877f](https://github.com/apache/plc4x/commit/de8877fa37378e5a7b406a69859efc76506fbb91))

## [rel/0.3](https://github.com/apache/plc4x/releases/tag/rel/0.3) - 2019-01-16 15:32:56

## What's Changed
* [PLC4X-75] - Fixing dependency to the wrap url-handler by @ANierbeck in https://github.com/apache/plc4x/pull/42
* S7 fix of array and String acquirement by @timbo2k in https://github.com/apache/plc4x/pull/43
* S7: changed byteLength and blockNumber from short to int by @timbo2k in https://github.com/apache/plc4x/pull/44

## New Contributors
* @ANierbeck made their first contribution in https://github.com/apache/plc4x/pull/42

**Full Changelog**: https://github.com/apache/plc4x/compare/rel/0.2...rel/0.3

## [rel/0.2](https://github.com/apache/plc4x/releases/tag/rel/0.2) - 2018-11-14 15:20:37

## What's Changed
* Added note for Service Loader. by @JulianFeinauer in https://github.com/apache/plc4x/pull/23
* Feature/top level item spliting by @sruehl in https://github.com/apache/plc4x/pull/24
* PLC4X-57 Bugfix by @JulianFeinauer in https://github.com/apache/plc4x/pull/25
* Type conversions for default byte array fiel item by @timbo2k in https://github.com/apache/plc4x/pull/28
* API Refactoring: add execute operation to requests, extract SPI package in https://github.com/apache/plc4x/pull/27
* Opm plcentitymanager by @JulianFeinauer in https://github.com/apache/plc4x/pull/29
* Provide connection metadata in https://github.com/apache/plc4x/pull/32
* Add simple mock driver by @JulianFeinauer in https://github.com/apache/plc4x/pull/38


**Full Changelog**: https://github.com/apache/plc4x/compare/rel/0.1...rel/0.2

### Bug Fixes

- general:
  - fixed hello plc4x example ([2bb806d](https://github.com/apache/plc4x/commit/2bb806d6644c65df3c63c2f9e3396094d10b55fd))
  - fix build ([2906972](https://github.com/apache/plc4x/commit/29069721ba50025834503cf682c0da9e471ada49))
  - fixed documentation ([08437d2](https://github.com/apache/plc4x/commit/08437d2d26e3eb691c4befcc2bfe2a5750377d5c))
  - fixed examples to use the new api ([7375d04](https://github.com/apache/plc4x/commit/7375d040917b2de0391611d3aef8c4316b625950))
  - fixed integrations to use the new api ([361432c](https://github.com/apache/plc4x/commit/361432c63b4d94701ea52ddd4a1cf8632a28c4db))

## [rel/0.1](https://github.com/apache/plc4x/releases/tag/rel/0.1) - 2018-09-17 16:31:32

## What's Changed
* Warning cleanup by @dlaboss in https://github.com/apache/plc4x/pull/1
* Change from AssertJ to Hamcrest by @niclash in https://github.com/apache/plc4x/pull/4
* Cleanup warnings by @dlaboss in https://github.com/apache/plc4x/pull/5
* cleanup sonar edgent complaints by @dlaboss in https://github.com/apache/plc4x/pull/6
* cleanup some warnings (in Eclipse) by @dlaboss in https://github.com/apache/plc4x/pull/7
* Added some javadoc to S7 communication path and several todos that ma… by @JulianFeinauer in https://github.com/apache/plc4x/pull/9
* fix a few typos by @bodewig in https://github.com/apache/plc4x/pull/12
* Implement basic example of connecting an S7 device to Azure IoT Hub by @DerDackel in https://github.com/apache/plc4x/pull/11
* Basic example to connect S7 device to Google Cloud IoT Core by @pisquaredover6 in https://github.com/apache/plc4x/pull/13
* Quick and dirty implementation of Apache Kafka Source Connector in https://github.com/apache/plc4x/pull/15
* Implemented throttling in Kafka Source Connector in https://github.com/apache/plc4x/pull/16
* Add support for multiple queries in kafka source connector in https://github.com/apache/plc4x/pull/17
* Add support for multiple tasks in kafka sink connector in https://github.com/apache/plc4x/pull/18
* Add URL Field to Key Schema in Kafka Source Connector in https://github.com/apache/plc4x/pull/19
* Remove old files in https://github.com/apache/plc4x/pull/22

## New Contributors
* @dlaboss made their first contribution in https://github.com/apache/plc4x/pull/1
* @niclash made their first contribution in https://github.com/apache/plc4x/pull/4
* @bodewig made their first contribution in https://github.com/apache/plc4x/pull/12
* @DerDackel made their first contribution in https://github.com/apache/plc4x/pull/11
* @pisquaredover6 made their first contribution in https://github.com/apache/plc4x/pull/13

**Full Changelog**: https://github.com/apache/plc4x/commits/rel/0.1

### Bug Fixes

- general:
  - fixed compare to messup in AdsPlcFieldHandler ([4bd85c7](https://github.com/apache/plc4x/commit/4bd85c72fe1dfc4ce6fd87e50feeaa4f4dbaa648))
  - fixed issues with site generation ([2587167](https://github.com/apache/plc4x/commit/258716723c1a3ae023349d64c8cee6aeebb1c2df))
  - fixed ethernetip field handler ([cb27228](https://github.com/apache/plc4x/commit/cb27228be4e923cdaae2335722e580583cc7f98f))
  - fixed missing methods for BigInteger Support ([5c3a127](https://github.com/apache/plc4x/commit/5c3a127763b919b11ccd0c7ec15564de1f60e65a))
  - fixed remaining tests in ADS. Still needs a lot of refactoring and todo tests ([16fc039](https://github.com/apache/plc4x/commit/16fc0398b44c19bfec9888eef2f622d7c53abefd))
  - fix some issues with generics and build ([ba34efd](https://github.com/apache/plc4x/commit/ba34efd7e50ba9118ea49cf06ff440c55b3d8dfd))
  - fix a few typos ([ec261d5](https://github.com/apache/plc4x/commit/ec261d570deaef60a222db22e6cde72484dc25e0)) ([#12](https://github.com/apache/plc4x/pull/12))
  - fixed some sonar bugs ([7a38a36](https://github.com/apache/plc4x/commit/7a38a361f43ab7504f48fd65f54c8b842b74df5e))
  - fixed some sonar bugs ([4f38a3f](https://github.com/apache/plc4x/commit/4f38a3fae4cd4bc481ee6b1514fb95625fd20587))
  - fixed double implementation in S7 ([810d649](https://github.com/apache/plc4x/commit/810d649f4912c5f720f44f6aed41440a7b99fb27))
  - fixed upper bounds for registers. ([dbc6812](https://github.com/apache/plc4x/commit/dbc6812517b065e5a000eee324b8548eb3107249))
  - fixed documentation of test class javadoc ([9efd47c](https://github.com/apache/plc4x/commit/9efd47c32a6d5d8125636b020dadd80f3097f982))
  - fixed calender tests and simplified assertions while at it ([6adebdf](https://github.com/apache/plc4x/commit/6adebdffbc44d7111e935b4cb550ccd19357999d))
  - fixed build problems ([036e394](https://github.com/apache/plc4x/commit/036e3942f599a5ab87210feffc2b9f8a46616d5d))
  - fixed produce coil value method. Added a basic test ([cdddd19](https://github.com/apache/plc4x/commit/cdddd19e5bef17f4ccdddf9f1cb27b5fc1b69fa3))
  - fixed small typo ([4fe91d0](https://github.com/apache/plc4x/commit/4fe91d0516c2194d4e4d704c38981cb4d85e926f))
  - fix build ([e2e735e](https://github.com/apache/plc4x/commit/e2e735ea68e3c909fa967d88b0b7fe22ada0bdf8))
  - fix build ([948bbba](https://github.com/apache/plc4x/commit/948bbba1c81ca0b649dfe64f4faaa065c64fdc86))
  - fix dependency problems ([0db64c8](https://github.com/apache/plc4x/commit/0db64c8c003f4735294529858cd184c5c2793faf))
  - fixed small warning ([ffc1680](https://github.com/apache/plc4x/commit/ffc16807010e234b6b8e107f30b54467dab87b64))
  - fixed serial driver regex to match the serial port the right way ([151a20f](https://github.com/apache/plc4x/commit/151a20ffbcf00082b2336f97cbbb1881425cdbc7))
  - fix remaining sonar bug ([de97867](https://github.com/apache/plc4x/commit/de97867f1aa5c85d92671a8ffcf301d9c59b25f4))
  - fixed some bugs while working on coverage ([8b8a938](https://github.com/apache/plc4x/commit/8b8a9388d7e9e0caefc6c42c57932c871ac27ad6))
  - fixed sonar warning about not serializable fields ([96eaca6](https://github.com/apache/plc4x/commit/96eaca6cff19802bc92c39ae6cfa1fd74baec436))
  - fixed some sonar bugs (equals, hashcode, serializable) ([864bd96](https://github.com/apache/plc4x/commit/864bd9626177e6cef1507c01895fb270ac9114a8))
  - fixed build ([89c034c](https://github.com/apache/plc4x/commit/89c034c9523cb2ab4de18b7f9d6c536d37887648))
  - fixed tests ([b0fad96](https://github.com/apache/plc4x/commit/b0fad962e23dc46143fb915f3a22ee42c9dba9c5))
  - fixed device name handling ([20090bb](https://github.com/apache/plc4x/commit/20090bb78993256448cb79c2b17580c65c859eb7))
  - fixed wrong padding of strings. ([4b7fc7e](https://github.com/apache/plc4x/commit/4b7fc7e9a7331eb5c5fd9cb4b07e1e858be890e5))
  - fixed failing test on travis ([976464d](https://github.com/apache/plc4x/commit/976464d3d80d332a5cd78fa90941364983e4ce7a))
  - fixed wrong log category ([db58135](https://github.com/apache/plc4x/commit/db5813586acc07b823e3b80768776ba5f5c22ec8))
  - fixed wrong constant for calendar ([238a2b0](https://github.com/apache/plc4x/commit/238a2b0f14cfeebd6788b16bb6a5d7f00ddf6e38))
  - fix build ([bca6a0e](https://github.com/apache/plc4x/commit/bca6a0e2970bf139d2fe58a6e08da41b05899468))
  - fixed bug where we remove a mapping too early ([4f1b37c](https://github.com/apache/plc4x/commit/4f1b37ca47d11d07ee06cce2c67384e4f202b358))
  - fixed sonar bugs ([4aa6e87](https://github.com/apache/plc4x/commit/4aa6e8768672bb7c919439669a6d0b16e5c148b3))
  - fix build ([da40d74](https://github.com/apache/plc4x/commit/da40d74d8dc1e935e95dd6d095162df0d08b3f6b))
  - fixed handling of variable response lengths. ([010c5cb](https://github.com/apache/plc4x/commit/010c5cbf029e07d0c73c9e631f926070e999afc5))
  - fix failing tests due to fix ([b0a0abd](https://github.com/apache/plc4x/commit/b0a0abd69e02e78764919aa63ff20b3657b75847))
  - fix failing build by adding netty-common dependency ([b5795a7](https://github.com/apache/plc4x/commit/b5795a7eb37e8b2ec5ea4a7ff0c0fce0b6c36c7a))
  - fixed inverted ads pipeline ([d4868e2](https://github.com/apache/plc4x/commit/d4868e2be0a9bb3eff71eb973b39f7f938e8f096))
  - fixed remaining Tests after refactoring ([3c0c4c9](https://github.com/apache/plc4x/commit/3c0c4c9267e3cfa4cb53f082ae37551b96269933))
  - fixed DigestUtil under the assumption the crc16 is LittleEndian ([df13820](https://github.com/apache/plc4x/commit/df1382063fbb78f6e4b02a2a998a17601c294eca))
  - fix DigestUtil by using unsigned int conversion ([1ba0e8b](https://github.com/apache/plc4x/commit/1ba0e8bc09c338e5d7173e987b4a97cf6261a61e))
  - fixed crc left-overs ([9a2c27e](https://github.com/apache/plc4x/commit/9a2c27e396d407358904cd38113b4e170cf14b15))
  - fixed Benchmark ([f79137e](https://github.com/apache/plc4x/commit/f79137ec0b528a5b969c3844ab2d41ee78ece4dd))
  - fixed digest calculation due to short overflow ([933adbe](https://github.com/apache/plc4x/commit/933adbe4a81dc996427a06c588a88b8aaf287cad))
  - fixed broken benchmark code ([a994a4e](https://github.com/apache/plc4x/commit/a994a4e6f85bcc23993b20d8392bd6911ae3a2c7))
  - fix build ([39fb481](https://github.com/apache/plc4x/commit/39fb481a559a6091605a184b37a7de9662653f12))
  - fix build ([13932e9](https://github.com/apache/plc4x/commit/13932e9bed5cba1479e0ce30266f4fc261f876be))
  - fixed renaming messup ([ab1a277](https://github.com/apache/plc4x/commit/ab1a277322ef157257239a49dcd15b60aa02795e))
  - fix leftover from refactoring ([1a3cfe6](https://github.com/apache/plc4x/commit/1a3cfe661b87f27a22e7c76146c99a09ad473341))
  - fix speling ([bd61054](https://github.com/apache/plc4x/commit/bd6105489f6d46a57c5b6962d76427a354e88982))
  - fixed remaining float todo. ([bd1c69f](https://github.com/apache/plc4x/commit/bd1c69f87f31cf266bd2995c3473ea61d338b50d))
  - fix sonarqube warnings ([e88e8d3](https://github.com/apache/plc4x/commit/e88e8d36f04c26c70b23df072955e0c4e2154ad2))
  - fixed problem with missing variable ([8360582](https://github.com/apache/plc4x/commit/8360582787b39af452333fc1c1e89730cfc494ea))
  - fixed missing MAX_VALUE support ([183b51d](https://github.com/apache/plc4x/commit/183b51d36ef3073c04388ab6eb58ca472df14185))
  - fixed build as the toString() from Device.java would include non printable
chars into <testcase name= in surefire plugin. This could be a bug with
surefire because it should escape non-printable chars here too (might be
the problem of junit4-parameterized too). ([0a4ed92](https://github.com/apache/plc4x/commit/0a4ed9204e473277ffe9895a75613bc8170690e6))
  - fixed a bunch of sonar warnings. ([976d15f](https://github.com/apache/plc4x/commit/976d15f371d673d873ea85436fc7cd65a368f7de))
  - fixed typo paket -> packet ([4e904ba](https://github.com/apache/plc4x/commit/4e904bab0d7d134b852228dd705c9946dc545419))
  - fixed test by letting the TcpHexDumper perform a gentle shutdown ([e01a7b5](https://github.com/apache/plc4x/commit/e01a7b5197071720e29a2cdeae60d0fe554d8df8))
  - fixed state mask and simplified generated requests. ([2cd0c2c](https://github.com/apache/plc4x/commit/2cd0c2c0ffb410c352c1bf85f19f9a6c8ce69404))
  - fix build ([b4f5570](https://github.com/apache/plc4x/commit/b4f5570f15535db9dd3210a0afd42028f154e66d))
  - fixed usage of wrong Type ([f609240](https://github.com/apache/plc4x/commit/f6092400a91c6b43bbd10a74a2fe3dc9855fe13a))
  - fixed remaining implementations and added tests for it ([8e0e8e6](https://github.com/apache/plc4x/commit/8e0e8e633c3c60ce2c99c062e0cf136da93eddd9))
  - fixed failing build ([7805d0c](https://github.com/apache/plc4x/commit/7805d0c1c89d12598250ea4a1677dec26ef6ffdd))
  - fixed ReadLength implementation ([d0e3445](https://github.com/apache/plc4x/commit/d0e34452450a8b76bf033aa85e8cd30746d4aad1))
  - fixed wrong implementation of Length ([46fa650](https://github.com/apache/plc4x/commit/46fa6504226d71010fe834ab9cfa9c8919e97860))
  - fixed wrong implementation of AMSPort ([5d3ecb8](https://github.com/apache/plc4x/commit/5d3ecb85fa0afe392e335cdd10ffd393dfd58619))
  - fixed wrong masking of integer as we masked to much ([013421c](https://github.com/apache/plc4x/commit/013421c050cb145e7082d611a1c1168a9f14c1fa))
  - fix logging so no side effects ([1c831ad](https://github.com/apache/plc4x/commit/1c831ad213de8699ab3766f70e615e936ed480a3))
  - fix logging so no side effects ([716ac61](https://github.com/apache/plc4x/commit/716ac61b06f268502e99c5dba036172a41277e3c))
  - fix logging so no side effects ([c6c50b1](https://github.com/apache/plc4x/commit/c6c50b1e79f02663534eff8518d72f2682764ab5))
  - fix missing coverage by switching to junit5 and make test compatible ([38e344e](https://github.com/apache/plc4x/commit/38e344e8158084d2062279dc0fc33174f4ea695d))
  - fixed wrong default reporting path as this variable already contains target ([e596961](https://github.com/apache/plc4x/commit/e5969616dcece9d696568798cada7a5c083f010a))
  - fix wrong scope. Meant to be optional true ([b8c0b06](https://github.com/apache/plc4x/commit/b8c0b063b4cfb1ba0be9aa262a37a87bf8afad7a))
  - fixed timing issues on get responseItem by signaling after item has
been set not before. ([5886538](https://github.com/apache/plc4x/commit/58865385e5028e4b113831685c1acfedd32c5cae))
  - fixed edgent test ([3ce27ea](https://github.com/apache/plc4x/commit/3ce27ea3b73614c86cddcdcd44880fe0c46586d3))
  - fixed test by cleaning up the generics. ([675a108](https://github.com/apache/plc4x/commit/675a10869faa495384e5539a68dd7f3a44e27d6c))
  - fixed one more builder bug ([b86dc42](https://github.com/apache/plc4x/commit/b86dc4211a2a46e0a40d34f3a150380241e2244b))
  - fixed firstType builder bug; added getValue to bulk; added getValue tests ([b3869eb](https://github.com/apache/plc4x/commit/b3869ebee63feb825bfbee2f6f21a7bace8c53ca))
  - fix tests ([80fa0ae](https://github.com/apache/plc4x/commit/80fa0ae3d67ff5b6c4abd6c4ebd47c3473e21a90))
  - fix up switch and endless loops ([ff882df](https://github.com/apache/plc4x/commit/ff882df35e2487502d7865d5f546ed235e225547))
  - fix build by adding missing apache headers ([ce7eaee](https://github.com/apache/plc4x/commit/ce7eaee3f5aca6e80abb48ede5665ef17b703ad3))
  - fix version and add tests and coverage ([6b64873](https://github.com/apache/plc4x/commit/6b64873bcc9e36653f07d811434c607d403c5430))
  - fix failing test ([ad6bf6f](https://github.com/apache/plc4x/commit/ad6bf6fcbfb7ae9ac35ab09dbaae14e275cc505d))
  - fixed site generation ([d0254d0](https://github.com/apache/plc4x/commit/d0254d07d9b277d59cb1f9fa43d1b050192cae7a))
  - fixed failing build due to missing documentation. ([a1b365b](https://github.com/apache/plc4x/commit/a1b365bd5212cf4c998d106a36dbe0725394b6c6))
  - fix pom for plx4x -> plc4x rename ([b1a2a92](https://github.com/apache/plc4x/commit/b1a2a9244b57f1ecfd05050e60be429d141cbf71))

### Refactor

- general:
  - refactoring ([28a70be](https://github.com/apache/plc4x/commit/28a70be675efff73215de8f72d3dce764ea67a6a))
  - refactored common used test code into driver-bases-test ([234c121](https://github.com/apache/plc4x/commit/234c1211e54ebdbf3dda83847b0146ead42da209))
  - refactored protocol tests ([bf6bac3](https://github.com/apache/plc4x/commit/bf6bac343afbd0157676d4a28792c27a6f7856d3))
  - refactor to make simplier ([a0a4ef6](https://github.com/apache/plc4x/commit/a0a4ef6197707c80cb13b21705939163e5626c80))
  - refactored common code ([444e16a](https://github.com/apache/plc4x/commit/444e16a8be7e32e2d88542e58ccc5b19b8dc55c8))
  - refactor to reduce complexity ([d88234a](https://github.com/apache/plc4x/commit/d88234a0af556cc3b19181dd3503ad5ced2ec6b5))
  - refactored much common code to a abstract class ([9f6fc95](https://github.com/apache/plc4x/commit/9f6fc95b43b2ee6e6a5e51ae3e0010ab44c05703))
  - refactored common code of request and responses into abstract classes ([af2ca39](https://github.com/apache/plc4x/commit/af2ca39c33a80180f0121a9ed105696fcdf0fae2))
  - refactor to reduce complexity and fix endless loops ([04245aa](https://github.com/apache/plc4x/commit/04245aa62de2545be239c3945fede4abf60d16d5))
  - refactor and clean up switches ([641c99f](https://github.com/apache/plc4x/commit/641c99f40b80bc7012602ec75a621a8871020037))

\* *This CHANGELOG was automatically generated by [auto-generate-changelog](https://github.com/BobAnkh/auto-generate-changelog)*
