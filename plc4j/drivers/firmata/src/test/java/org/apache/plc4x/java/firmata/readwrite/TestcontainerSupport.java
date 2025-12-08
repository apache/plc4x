package org.apache.plc4x.java.firmata.readwrite;

import com.github.pfichtner.testcontainers.virtualavr.VirtualAvrContainer;
import org.testcontainers.images.ImagePullPolicy;
import org.testcontainers.utility.DockerImageName;

import java.io.File;
import java.util.UUID;

public final class TestcontainerSupport {

	private static final String VIRTUALAVR_DOCKER_TAG_PROPERTY_NAME = "virtualavr.docker.tag";
	private static final String VIRTUALAVR_DOCKER_PULL_PROPERTY_NAME = "virtualavr.docker.pull";

	private TestcontainerSupport() {
		super();
	}

	@SuppressWarnings("resource")
	public static VirtualAvrContainer<?> virtualAvrContainer(File inoFile) {
		VirtualAvrContainer<?> container = new VirtualAvrContainer<>(imageName()) //
				.withImagePullPolicy(onlyPullIfEnabled()) //
				.withSketchFile(inoFile) //
				.withDeviceName("virtualavr" + UUID.randomUUID()) //
				.withDeviceGroup("root") //
				.withDeviceMode(666);

		// Use TCP serial mode on non-Linux systems (macOS/Windows with Docker Desktop)
		if (!isLinux()) {
			container.withTcpSerialMode();
		}

		return container;
	}

	public static ImagePullPolicy onlyPullIfEnabled() {
		return imageName -> System.getProperty(VIRTUALAVR_DOCKER_PULL_PROPERTY_NAME) != null;
	}

	public static DockerImageName imageName() {
		String dockerTagName = System.getProperty(VIRTUALAVR_DOCKER_TAG_PROPERTY_NAME, "dev");
		DockerImageName defaultImageName = VirtualAvrContainer.DEFAULT_IMAGE_NAME;
		return dockerTagName == null ? defaultImageName : defaultImageName.withTag(dockerTagName);
	}

	/**
	 * Checks if the current OS is Linux.
	 * On Linux, the standard PTY mode works because /dev can be bind-mounted.
	 * On macOS/Windows with Docker Desktop, TCP serial mode is needed.
	 *
	 * @return true if running on Linux
	 */
	public static boolean isLinux() {
		String os = System.getProperty("os.name", "").toLowerCase();
		return os.contains("linux");
	}

}
