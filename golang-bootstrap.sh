#!/bin/bash

# Golang installation variables
VERSION="1.9"
OS="linux"
ARCH="amd64"

# Home of the vagrant user, not the root which calls this script
HOMEPATH="/home/vagrant"

# Updating and installing stuff
sudo apt-get update
sudo apt-get install -y git curl

if [ ! -f "$HOMEPATH/go.tar.gz" ]; then
	# No given go binary
	# Download golang
	FILE="go$VERSION.$OS-$ARCH.tar.gz"
	URL="https://storage.googleapis.com/golang/$FILE"

	echo "Downloading $FILE ..."
	curl $URL -L -o "$HOMEPATH/go.tar.gz" || rm -f "$HOMEPATH/go.tar.gz"
#else
#	# Go binary given
#	echo "Using given binary ..."
#	cp "/vagrant/go.tar.gz" "$HOMEPATH/go.tar.gz"
fi;

echo "Cleaning ..."
rm -rf "$HOMEPATH/.go"

echo "Extracting ..."
tar -C "$HOMEPATH" -xzf "$HOMEPATH/go.tar.gz" || rm -f "$HOMEPATH/go.tar.gz"
mv "$HOMEPATH/go" "$HOMEPATH/.go"
#rm "$HOMEPATH/go.tar.gz"

# Create go folder structure
GP="/vagrant/GO"
mkdir -p "$GP/src"
mkdir -p "$GP/pkg"
mkdir -p "$GP/bin"

cat > $HOMEPATH/.bashrc << EOF
# Golang environments
export GOROOT=$HOMEPATH/.go
export PATH=\$PATH:\$GOROOT/bin
export GOPATH=$HOMEPATH/GO
export PATH=\$PATH:\$GOPATH/bin

# Prompt
export PROMPT_COMMAND=_prompt
_prompt() {
    local ec=\$?
    local code=""
    if [ "\$ec" != "0" ]; then
        code="\[\e[0;31m\][\${ec}]\[\e[0m\] "
    fi
    PS1="\${code}\[\e[0;32m\][\u] \W\[\e[0m\] \$ "
}

# Automatically change to the vagrant dir
cd /vagrant
EOF
