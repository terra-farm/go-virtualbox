# -*- mode: ruby -*-
# vi: set ft=ruby :

gosrc = File.join(ENV['GOPATH'] || File.join((ENV['HOME'] || ENV['USERPROFILE']), "go"),"src")

Vagrant.configure("2") do |config|

  config.vm.box = "bento/ubuntu-16.04"

  #config.vm.define 'go-virtualbox'
  #config.vm.hostname = 'go-virtualbox'

  config.vm.synced_folder Dir.home, '/home/vagrant/home', create: true
  config.vm.synced_folder gosrc, '/home/vagrant/GO/src', create: true

  #config.vm.provision 'file', source: 'golang-bashrc', destination: '~/.bashrc'
  config.vm.provision :shell, path: "golang-bootstrap.sh"

  config.vm.provider :virtualbox do |vb|
   vb.name = 'go-virtualbox'
  end

end
