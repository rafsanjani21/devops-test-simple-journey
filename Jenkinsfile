pipeline {
    agent any
    environment {
        APP_NAME    = 'devops-app'
        PORT        = '8080'
        COMMIT_HASH = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
    }
    stages {
        stage('Checkout') {
            steps { checkout scm }
        }
        stage('Test') {
            steps { sh 'go test -v ./...' }
        }
        stage('Build Image') {
            steps {
                sh "docker build --build-arg VERSION=${COMMIT_HASH} -t ${APP_NAME}:${COMMIT_HASH} ."
            }
        }
        stage('Deploy (Hotfix / Swap Binary)') {
            steps {
                sh """
                    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s -X main.version=${COMMIT_HASH}" -o server-hotfix .
                    if [ \$(docker ps -q -f name=^/${APP_NAME}\$) ]; then
                        docker cp server-hotfix ${APP_NAME}:/app/server
                        docker restart ${APP_NAME}
                    else
                        docker run -d --name ${APP_NAME} --restart unless-stopped -p ${PORT}:8080 ${APP_NAME}:${COMMIT_HASH}
                    fi
                    sleep 2
                    curl -s http://localhost:${PORT}/ | grep "${COMMIT_HASH}"
                """
            }
        }
    }
}