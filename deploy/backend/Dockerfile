# //1.- Use the Go Alpine image to keep the development container lightweight while retaining necessary tooling.
FROM golang:1.21-alpine

# //2.- Install OS packages and Go utilities required for live reloading with Air.
RUN apk add --no-cache build-base git && \
    go install github.com/cosmtrek/air@v1.49.0

# //3.- Prepare the workspace so source code can be mounted in by docker compose.
WORKDIR /workspace/backend

# //4.- Configure the container to watch the repository via the shared volume when it starts.
CMD ["air", "-c", "/workspace/deploy/backend/air.toml"]
