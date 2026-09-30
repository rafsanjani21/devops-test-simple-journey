# DevOps Technical Test - Simple Journey

**Candidate:** Muhammad Rafsanjani  
**Go Version:** 1.22.2

---

## Part I: Build
- **Build Command:** `docker build --build-arg VERSION=1.0.0 -t devops-app:1.0.0 .`
- **Final Image Size:** ~15 MB.
- **Explanation:** I used a multi-stage build approach. The builder stage utilizes `golang:1.22-alpine` with `CGO_ENABLED=0` to compile a statically linked binary without dependencies on the host OS. The final stage uses `alpine:3.20` as the base image. Alpine was chosen over `scratch` because it maintains a minimal footprint while retaining basic utilities (like `sh` and user management), which are essential for executing `docker cp` commands and debugging effectively during hotfixes.

## Part II: Deploy & Zero-Rebuild Binary Swap
- **Run Command:** `docker run -d --name devops-service --restart unless-stopped -p 8080:8080 devops-app:1.0.0`
- **Hotfix Commands:** 
  1. `CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s -X main.version=2.0.0" -o server-hotfix .`
  2. `docker cp server-hotfix devops-service:/app/server`
  3. `docker restart devops-service`
- **Curl Output Before Swap:** `Hello, DevOps! version = 1.0.0`
- **Curl Output After Swap:** `Hello, DevOps! version = 2.0.0`
- **Explanation:** I chose the approach of compiling the static binary on the host/CI, transferring it via `docker cp`, and restarting the container. This method is highly suited for production hotfix scenarios because it completely bypasses the time-consuming process of rebuilding a full Docker image and pulling from a registry. The service downtime is restricted merely to the 1-2 seconds it takes for the container process to restart.

## Part III: CI/CD with Jenkins
- **Pipeline Snapshot:** See `pipeline-success.png` in the repository root.
- **Rollback Strategy:** If the deploy stage fails midway (e.g., the health check fails after `docker restart`), the Jenkins pipeline can catch this in a `post { failure }` block. An automated rollback can be implemented by executing a `docker cp` command to restore a pre-deploy backup of the binary (e.g., `server.bak`), or by spinning up the previous stable container image (`devops-app:${LAST_STABLE_HASH}`) to instantly restore the service.