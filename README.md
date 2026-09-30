# DevOps Technical Test - Simple Journey

## Part I: Build
- **Build Command:** `docker build --build-arg VERSION=1.0.0 -t devops-app:1.0.0 .`
- **Final Image Size:** ~15 MB.
- **Explanation:** I used a multi-stage build. The builder stage utilizes `golang:1.22-alpine` with `CGO_ENABLED=0` to create a statically linked binary. The final stage uses `alpine:3.20` as the base image. Alpine was chosen over `scratch` because it is minimal yet retains basic utilities (like `sh` and user management) which are useful for debugging and executing `docker cp` commands effectively.

## Part II: Deploy & Zero-Rebuild Binary Swap
- **Run Command:** `docker run -d --name devops-service --restart unless-stopped -p 8080:8080 devops-app:1.0.0`
- **Hotfix Commands:** 
  1. `CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s -X main.version=2.0.0" -o server-hotfix .`
  2. `docker cp server-hotfix devops-service:/app/server`
  3. `docker restart devops-service`
- **Explanation:** I chose the approach of compiling the static binary on the host/CI and transferring it using `docker cp` followed by a container restart. This method is ideal for production hotfix scenarios because it circumvents the lengthy process of full image rebuilding and registry pulling, keeping the downtime to merely 1-2 seconds during the process restart.

## Part III: CI/CD with Jenkins
- **Pipeline Snapshot:** See `pipeline-success.png` in the repository.
- **Rollback Strategy:** If the deploy stage fails midway (e.g., health check fails after `docker restart`), the pipeline is configured to catch the failure in the `post { failure }` block. An automated rollback can be implemented here by restoring a pre-deploy backup of the binary (`server.bak`) or executing a `docker run` command utilizing the previous stable image tag (`${APP_NAME}:${LAST_STABLE_HASH}`) to instantly restore the service.