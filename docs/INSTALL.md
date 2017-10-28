# Manual installation

Follow these instructions to manually install the dependencies to run, test and build ringo

## Golang installation

Download and install golang for your operating system by following this link: <https://golang.org/dl/>

## Setup your $GOPATH env variable

- **macOS**

	Edit (or create) `$HOME/.profile` and add this line :

		export GOPATH=$HOME/dev/go
		export PATH=$GOPATH/bin:$PATH
		...

- **linux**

	Edit (or create) `$HOME/.bash_profile` and add this line :

		export GOPATH=~/go
		export PATH=$GOPATH/bin:$PATH
		...

## GB installation

	go get github.com/constabulary/gb/...

## Docker installation

Download and install Docker  for your operating system by following this link: <https://docs.docker.com/engine/installation/>

## Configuration file

Copy files from `ansible/dev/conf/*.json.j2` in `conf/` folder to a `.json` file and edit configurations to match your setup.

## Install Xcode command line tools (macOS Only)

	xcode-select --install