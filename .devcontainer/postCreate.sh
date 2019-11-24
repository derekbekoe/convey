set -x

# Delete the usual cloned repo contents as we're getting the repo with 'go get' instead.
cd /home/vsonline
rm -rf /home/vsonline/workspace

echo "Installing go..."
tmp_file=$(mktemp)
wget -q -O $tmp_file https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf $tmp_file
rm $tmp_file
echo "go installed."

new_gopath=/home/vsonline/workspace/go
go_convey_src=$new_gopath/src/github.com/derekbekoe/convey
echo "PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
echo "GOPATH=$new_gopath" >> ~/.bashrc
echo "GOPATH_CONVEY=$go_convey_src" >> ~/.bashrc
export GOPATH=$new_gopath
go_exe=/usr/local/go/bin/go
$go_exe get -u github.com/derekbekoe/convey

# Install VS Code Go extension dependencies - https://github.com/Microsoft/vscode-go/wiki/Go-tools-that-the-Go-extension-depends-on
echo "Installing VS Code Go extension dependencies..."
$go_exe get -v github.com/ramya-rao-a/go-outline
$go_exe get -v github.com/rogpeppe/godef
$go_exe get -v github.com/mdempsky/gocode
$go_exe get -v github.com/uudashr/gopkgs/cmd/gopkgs
$go_exe get -v github.com/stamblerre/gocode
$go_exe get -v github.com/sqs/goreturns
$go_exe get -v golang.org/x/lint/golint
echo "Done."

echo "NOTE: Use the command below to be in the correct development source directory."
echo "cd \$GOPATH_CONVEY"
