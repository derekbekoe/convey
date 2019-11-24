
tmp_file=$(mktemp)
wget -O $tmp_file https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf $tmp_file
rm $tmp_file
echo "PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
