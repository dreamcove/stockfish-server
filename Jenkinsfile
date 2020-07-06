pipeline {
  agent {
    docker {
      image 'dreamcove/build-aws:latest'
    }
  }
  options {
    disableConcurrentBuilds()
    buildDiscarder(logRotator(numToKeepStr: '5'))
    timeout(time: 1, unit: 'HOURS')
  }
  stages {
    stage('Setup') {
      steps {
        slackSend channel: "#builds", message: "Build <${env.BUILD_URL}|#${env.BUILD_NUMBER}> of ${env.JOB_NAME} has been started by ${env.GIT_AUTHOR_NAME}"
        sh 'export'
      }
    }
    stage('Build') {
      steps {
        sh 'docker build .'
      }
    }
    stage('Lambdas') {
      when {
        expression {
          env.DEPLOYENV == 'master'
        }
      }
      steps {
        sh 'docker push dreamcove/stockfish-server'
      }
    }
  }
  post {
    success {
      slackSend channel: "#builds", message: "Build <${env.BUILD_URL}|#${env.BUILD_NUMBER}> of ${env.JOB_NAME} completed successfully"
    }
    failure {
      slackSend channel: "#builds", message: "Build <${env.BUILD_URL}|#${env.BUILD_NUMBER}> of ${env.JOB_NAME} failed"
    }
  }
}