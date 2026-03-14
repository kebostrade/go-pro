# Quiz: Containerization with Docker

## Question 1

What is the main advantage of multi-stage Docker builds?

A) Faster image push/pull
B) Smaller final image size by excluding build dependencies
C) Better security
D) Easier debugging

## Question 2

Which Docker flag is used to access NVIDIA GPUs?

A) --gpu
B) --gpus all
C) --nvidia
D) --device nvidia

## Question 3

What is the purpose of the `.dockerignore` file?

A) To exclude files from the final image
B) To exclude files from the build context
C) To disable Docker build caching
D) To specify which files to copy

## Question 4

In Docker, what is a "layer"?

A) A network interface
B) A filesystem snapshot that forms part of the image
C) A running container instance
D) A volume mount point

## Question 5

Which base image provides NVIDIA GPU support?

A) ubuntu:22.04
B) nvidia/cuda:12.1.0-runtime-ubuntu22.04
C) python:3.11-slim
D) alpine:3.19

## Question 6

What does `COPY --from=builder` do?

A) Copies files from another Docker image
B) Copies files from the build stage of the same Dockerfile
C) Copies files from a remote server
D) Copies files with root privileges

## Question 7

What is the recommended way to pass secrets to containers in production?

A) Environment variables in Dockerfile
B) Hardcoded in application code
C) Docker secrets or external secret management
D) Command line arguments

## Question 8

Which directive in Dockerfile creates a new layer?

A) FROM
B) RUN
C) LABEL
D) COMMENT

## Question 9

What is the purpose of EXPOSE in Dockerfile?

A) To publish a port to the host
B) To document which ports the container listens on
C) To allocate network bandwidth
D) To enable port forwarding

## Question 10

Which is the most secure way to run a container?

A) Running as root with all capabilities
B) Running as non-root user with read-only filesystem
C) Running with privileged mode
D) Running without any resource limits

## Question 11

In Docker Compose, what does the `depends_on` directive do?

A) Ensures containers start in the correct order
B) Links containers together
C) Shares volumes between containers
D) Enables inter-container communication

## Question 12

What is the purpose of healthcheck in docker-compose?

A) To restart containers automatically
B) To monitor container health and readiness
C) To load balance between containers
D) To cache container images
