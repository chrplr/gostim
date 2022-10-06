Experiments with Go and SDL2
============================
<christophe@pallier.org>



Prerequisites
-------------

`github.com/veandco/go-sdl2` will be downloaded during the first compilation.

The library [SDL2](http://libsdl.org) must be installed. 

Under Linux Debian-like:


	sudo apt install libsdl2-*



Compilation
-----------

	cd streaming
	go build streaming.go
        ./streaming


Note: Installation of Go a Raspberry Pi
---------------------------------------


Following the instuctions at <https://www.e-tinkers.com/2019/06/better-way-to-install-golang-go-on-raspberry-pi/>,
save the following script in `~/bin/go_installer.sh`


	export GOLANG="$(curl -s https://go.dev/dl/ | awk -F[\>\<] '/linux-armv6l/ && !/beta/ {print $5;exit}')"
	wget https://golang.org/dl/$GOLANG
	sudo tar -C /usr/local -xzf $GOLANG
	rm $GOLANG
	unset GOLANG

and run it:

     chmod +x ~/bin/go_installer.sh
     ~/bin/go_installer.sh




    echo "PATH=$PATH:/usr/local/go/bin" >>~/.profile
    echo "GOPATH=$HOME/golang" >>~/.profile
    source ~/.profile



