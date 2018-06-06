/*
 * Licensed under the Apache License, Version 2.0 (the &quot;License&quot;);
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an &quot;AS IS&quot; BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.apache.plc4x.utils.maven.site.asciidoctor.http;

import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.FileInputStream;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelFutureListener;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import io.netty.handler.codec.http.DefaultFullHttpResponse;
import io.netty.handler.codec.http.FullHttpRequest;
import io.netty.handler.codec.http.HttpHeaders;
import io.netty.handler.codec.http.HttpMethod;
import io.netty.handler.codec.http.HttpResponseStatus;
import io.netty.handler.codec.http.HttpVersion;
import io.netty.util.CharsetUtil;
import org.apache.commons.io.IOUtils;

public class AsciidoctorHandler extends SimpleChannelInboundHandler<FullHttpRequest> {
    private static final String HTML_MEDIA_TYPE = "text/html";
    public static final String HTML_EXTENSION = ".html";

    private final File directory;
    private final String defaultPage;

    public AsciidoctorHandler(final File workDir, final String defaultPage) {
        this.directory = workDir;

        if (defaultPage.contains(".")) {
            this.defaultPage = defaultPage;
        } else {
            this.defaultPage = addDefaultExtension(defaultPage);
        }
    }

    @Override
    public void channelRead0(final ChannelHandlerContext ctx, final FullHttpRequest msg) throws Exception {
        if (msg.getMethod() != HttpMethod.GET) {
            final DefaultFullHttpResponse response = new DefaultFullHttpResponse(HttpVersion.HTTP_1_1,
                    HttpResponseStatus.METHOD_NOT_ALLOWED,
                    Unpooled.copiedBuffer("<html><body>Only GET method allowed</body></html>", CharsetUtil.UTF_8));
            response.headers().set(HttpHeaders.Names.CONTENT_TYPE, HTML_MEDIA_TYPE);
            send(ctx, response);
            return;
        }

        final File file = deduceFile(msg.getUri());

        final HttpResponseStatus status;
        final ByteBuf body;
        final String mediaType;
        if (file.exists()) {
            final ByteArrayOutputStream baos = new ByteArrayOutputStream();
            final FileInputStream fileInputStream = new FileInputStream(file);
            IOUtils.copy(fileInputStream, baos);
            body = Unpooled.copiedBuffer(baos.toByteArray());
            IOUtils.closeQuietly(fileInputStream);

            mediaType = mediaType(file.getName());
            status = HttpResponseStatus.OK;
        } else {
            body = Unpooled.copiedBuffer("<body><html>File not found: " + file.getPath() + "<body></html>", CharsetUtil.UTF_8);
            status = HttpResponseStatus.NOT_FOUND;
            mediaType = HTML_MEDIA_TYPE;
        }

        final DefaultFullHttpResponse response = new DefaultFullHttpResponse(HttpVersion.HTTP_1_1, status, body);
        response.headers().set(HttpHeaders.Names.CONTENT_TYPE, mediaType);
        send(ctx, response);
    }

    private void send(final ChannelHandlerContext ctx, final DefaultFullHttpResponse response) {
        ctx.writeAndFlush(response).addListener(ChannelFutureListener.CLOSE);
    }

    private File deduceFile(final String path) {
        if (path.isEmpty() || "/".equals(path)) {
            return new File(directory, defaultPage);
        }

        if (!path.contains(".")) {
            return new File(directory, addDefaultExtension(path));
        }

        return new File(directory, path);
    }

    private static String addDefaultExtension(String path) {
        return path + HTML_EXTENSION;
    }

    private static String mediaType(final String name) {
        if (name.endsWith(".js")) {
            return "text/javascript";
        }
        if (name.endsWith(".css")) {
            return "text/css";
        }
        if (name.endsWith(".png")) {
            return "image/png";
        }
        if (name.endsWith(".gif")) {
            return "image/gif";
        }
        if (name.endsWith(".jpeg") || name.endsWith(".jpg")) {
            return "image/jpeg";
        }
        return HTML_MEDIA_TYPE;
    }
}
