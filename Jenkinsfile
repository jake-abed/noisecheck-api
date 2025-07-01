/* Requires the Docker Pipeline plugin */
pipeline {
    agent { docker { image 'golang:1.24.4-alpine3.22' } }
    stages {
        stage('build') {
            steps {
                sh 'go version'
            }
        }
    }
}
