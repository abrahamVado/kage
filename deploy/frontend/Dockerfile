# //1.- Start from a Node.js image that includes all tooling required for Next.js development.
FROM node:18-alpine

# //2.- Install dependencies commonly required for building native modules when the container installs packages.
RUN apk add --no-cache libc6-compat

# //3.- Create a working directory that will be shared with the host through bind mounts.
WORKDIR /workspace/frontend

# //4.- Expose Next.js' default development port for hot reloading.
EXPOSE 3000

# //5.- Run the Next.js development server when the container starts.
CMD ["npm", "run", "dev"]
