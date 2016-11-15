Vagrant.configure(2) do |config|
  config.vm.box = "debian/jessie64"
  config.vm.box_version = "8.2.1"
  config.vm.network "private_network", ip: "192.168.33.10"
  config.vm.synced_folder ".", "/opt/src/github.com/dthtvwls/crossfader"
  config.vm.provision "shell", inline: <<-SHELL
    export GOPATH=/opt
    export GOBIN="$GOPATH/bin"
    export PATH="$GOBIN:$PATH"

    echo 'deb http://ftp.debian.org/debian jessie-backports main' >> /etc/apt/sources.list
    apt-get update
    apt-get install -y -t jessie-backports --no-install-recommends git golang haproxy openjdk-8-jre-headless zookeeperd

    go get github.com/Masterminds/glide
    cd /opt/src/github.com/dthtvwls/crossfader
    glide install
    go install
  SHELL
end
