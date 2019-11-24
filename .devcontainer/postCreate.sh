echo "Installing go..."
tmp_file=$(mktemp)
wget -O $tmp_file https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf $tmp_file
rm $tmp_file
echo "PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
echo "GOPATH=/home/vsonline/workspace/go" >> ~/.bashrc
echo "GOPATH_CONVEY=/home/vsonline/workspace/go/src/github.com/derekbekoe/convey" >> ~/.bashrc
source ~/.bashrc

go get -u github.com/derekbekoe/convey

# Install VS Code Go extension dependencies - https://github.com/Microsoft/vscode-go/wiki/Go-tools-that-the-Go-extension-depends-on
echo "Installing VS Code Go extension dependencies..."
go get -v github.com/ramya-rao-a/go-outline
go get -v github.com/rogpeppe/godef
go get -v github.com/mdempsky/gocode
go get -v github.com/uudashr/gopkgs/cmd/gopkgs
go get -v github.com/stamblerre/gocode
go get -v github.com/sqs/goreturns
go get -v golang.org/x/lint/golint



