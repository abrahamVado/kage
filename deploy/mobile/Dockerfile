# //1.- Use the Flutter stable image so mobile engineers have the expected SDK pre-installed.
FROM ghcr.io/cirruslabs/flutter:3.16.9

# //2.- Create a workspace that will be mounted from the host for live reload operations.
WORKDIR /workspace/mobile

# //3.- Expose the web-server device port so Flutter web hot reload works in containers.
EXPOSE 8082

# //4.- Run flutter in web-server mode to enable fast developer feedback loops.
CMD ["flutter", "run", "-d", "web-server", "--web-hostname", "0.0.0.0", "--web-port", "8082"]
