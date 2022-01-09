pipeline {
  agent any
    environment {
      imagename = "docker.luispereira.xyz/mark-notes-server"
      registryUrl = "https://docker.luispereira.xyz"
        registryCredential = 'private-docker'
        dockerImage = ''
        APP_NAME = 'mark-notes-server'
        SERVER_SSH_HOST = 'luis@192.168.1.75'
    }
  stages {
    stage('Test application') {
      agent {
          label "ext_agent"
      }
      steps {
        echo "Cloning repo"
        checkout scm
        echo "Testing golang application in ${WORKSPACE}..."
        sh "go test ./..."
      }
    }
    stage('Build image') {
      steps {
        echo 'Building docker image..'
        script {
          dockerImage = docker.build imagename
        }
      }
    }
    stage('Deploy image') {
      steps {
        script {
          docker.withRegistry(registryUrl, registryCredential) {
            dockerImage.push("$BUILD_NUMBER")
            dockerImage.push("latest")
          }
        }
      }
    }
    stage('Remove Unused docker image') {
      steps {
        sh 'docker rmi $imagename:$BUILD_NUMBER'
        sh 'docker rmi $imagename:latest'
      }
    }
    stage('Setup Server') {
      steps {
        echo 'Writing setup script....'
          script {
            def data = "[ ! -d 'docker' ] && mkdir docker\n"
              data = data + "cd docker\n"
              data = data + "[ ! -d ${APP_NAME} ] && mkdir ${APP_NAME}\n"
              writeFile(file: 'setup.sh', text: data)
          }
        sshagent(['jenkins_global']) {
          echo 'Running Setup....'
            sh '''
            ssh -o StrictHostKeyChecking=no ${SERVER_SSH_HOST} 'bash -s' < setup.sh
            '''
            sh "scp docker-compose.yml ${SERVER_SSH_HOST}:docker/${APP_NAME}"
            sh "scp .env.example ${SERVER_SSH_HOST}:docker/${APP_NAME}"
        }
      }
    }
    stage('Prepare to Deploy') {
      steps {
        input "Prepare .env on ${SERVER_SSH_HOST}"
      }
    }
    stage('Deploy to Server') {
      steps {
        echo 'Writing build script....'
          script {
            def data = "cd docker && cd ${APP_NAME}\n"
              data = data + "[ ! -f '.env' ] && echo 'Missing .env' && exit 1\n"
              data = data + "docker-compose pull && docker-compose up -d\n"
              writeFile(file: 'build.sh', text: data)
          }
        sshagent(['jenkins_global']) {
          echo 'Deploying....'
            sh '''
            ssh -o StrictHostKeyChecking=no ${SERVER_SSH_HOST} 'bash -s' < build.sh
            '''
        }
      }
    }
  }
}