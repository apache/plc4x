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

import java.io.File;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ThreadFactory;
import java.util.concurrent.atomic.AtomicInteger;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelFutureListener;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelOption;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.SocketChannel;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.codec.http.HttpObjectAggregator;
import io.netty.handler.codec.http.HttpRequestDecoder;
import io.netty.handler.codec.http.HttpResponseEncoder;
import io.netty.handler.stream.ChunkedWriteHandler;
import io.netty.util.concurrent.Future;

import org.apache.maven.plugin.logging.Log;

public class AsciidoctorHttpServer {
    private static final String HOST = "localhost";
    private static final int THREAD_NUMBER = 3;
    private static final String THREAD_PREFIX = "asciidoctor-thread-";

    private final Log logger;
    private final int port;
    private final File workDir;
    private final String defaultPage;

    private ServerBootstrap bootstrap;
    private NioEventLoopGroup workerGroup;

    public AsciidoctorHttpServer(final Log logger, final int port, final File outputDirectory, final String defaultPage) {
        this.logger = logger;
        this.port = port;
        this.workDir = outputDirectory;
        this.defaultPage = defaultPage;
    }

    public void start() {
        final AtomicInteger threadId = new AtomicInteger(1);
        workerGroup = new NioEventLoopGroup(THREAD_NUMBER, new ThreadFactory() {
            @Override
            public Thread newThread(final Runnable r) {
                final Thread t = new Thread(r, THREAD_PREFIX + threadId.getAndIncrement());
                if (t.getPriority() != Thread.NORM_PRIORITY) {
                    t.setPriority(Thread.NORM_PRIORITY);
                }
                if (t.isDaemon()) {
                    t.setDaemon(false);
                }
                return t;
            }
        });

        try {
            bootstrap = new ServerBootstrap();
            bootstrap
                .option(ChannelOption.SO_REUSEADDR, true)
                .option(ChannelOption.SO_SNDBUF, 1024)
                .option(ChannelOption.TCP_NODELAY, true)
                .group(workerGroup)
                .channel(NioServerSocketChannel.class)
                .childHandler(new ChannelInitializer<SocketChannel>() {
                    @Override
                    protected void initChannel(final SocketChannel ch) {
                        ch.pipeline()
                            .addLast("decoder", new HttpRequestDecoder())
                            .addLast("aggregator", new HttpObjectAggregator(Integer.MAX_VALUE))
                            .addLast("encoder", new HttpResponseEncoder())
                            .addLast("chunked-writer", new ChunkedWriteHandler())
                            .addLast("asciidoctor", new AsciidoctorHandler(workDir, defaultPage));
                    }
                })
                .bind(HOST, port).addListener(new ChannelFutureListener() {
                @Override
                public void operationComplete(final ChannelFuture future) {
                    if (!future.isSuccess()) {
                        logger.error("Can't start HTTP server");
                    } else {
                        logger.info(String.format("Server started on http://%s:%s", HOST, port));
                    }
                }
            }).sync();
        } catch (final InterruptedException e) {
            logger.error(e.getMessage(), e);
        }
    }

   public void stop() {
      Future<?> shutdownGracefully = workerGroup.shutdownGracefully();
      logger.info("Server stopping...");
      try {
         shutdownGracefully.get();
         logger.info("Server stopped");
      } catch (InterruptedException e) {
         logger.error(e);
      } catch (ExecutionException e) {
         logger.error(e);
      }
   }

}
