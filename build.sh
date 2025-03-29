#!/usr/bin/sudo bash

if [ -f ./middle ]; then
	rm ./middle
fi

echo "Building main..."
sudo -E go build -buildvcs=false .

if [ -f ./middle ]; then
	sudo mv ./middle /usr/local/bin/middle
else
	echo "NOTE: MiddletonScript executable build incomplete."
fi

echo "Building wasm..."
sudo -E GOOS=js GOARCH=wasm go build -buildvcs=false -o middle.wasm .

if [ -f ./middle.wasm ]; then
	sudo mv ./middle.wasm ~/.config/MiddletonScript/middle.wasm
else
	echo "NOTE: MiddletonScript wasm build incomplete."
fi 
