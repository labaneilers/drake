# Drake : Docker make

<img width="400%" align="center" src="drake.jpg">

Drake is a cross-platform command-line build tool that executes build commands in a docker container, so that you can encapsulate all your build dependencies. Regardless of your development platform, your build commands will run in a single linux environment. 

# How does it work?
Drake builds a Docker image from your specified build-time Dockerfile, launches the container with your repository root mounted in the container's file system, and passes the shell command you specify via docker run. You use a simple yaml configuration file to map host commands to container commands.

# Example
Imagine I'm building a dotnet core application, and I'm using node/gulp to build static files. I also have a shell script to deploy.

First, create a configuration file: `drk.yaml`:

```yaml
dockerFile: Dockerfile.build
dockerDir: /code
commands:
  default: 
    command: dotnet
  deploy:
    command: ./deploy.sh
  static:
    command: gulp
```

Then, create a `Dockerfile.build` that contains your build dependencies:

```dockerfile
FROM microsoft/dotnet:sdk

# Install nodejs
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN apt-get install -y nodejs

# Install gulp
RUN npm install gulp -g
```

Then, you can use these same build commands from Windows, Mac, or Linux:

```bash
# Runs "dotnet build" inside the docker container
# The default command is "dotnet", and any command that isn't specifically defined gets appended to "dotnet".
$ drk build

# Runs "gulp"
$ drk static

# Runs ./deploy.sh --credential=12345
$ drk deploy --credential=12345
```

## Installation
Since Drake is a build tool, installation can be done with no dependencies:

### Windows

Run this in powershell:

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/labaneilers/drake/master/Install.ps1'))
```

### Mac OSX

Run this in bash:

```bash
wget -O - https://raw.githubusercontent.com/labaneilers/drake/master/install-darwin-amd64.sh | bash
```