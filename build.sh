#!/usr/bin/sudo bash

if [ -f ./middle ]; then
	rm ./middle
fi

sudo -E go build -buildvcs=false .

if [ -f ./middle ]; then
	sudo mv ./middle /usr/local/bin/middle
else
	echo "NOTE: MiddletonScript executable build incomplete."
fi 
