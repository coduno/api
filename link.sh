#!/bin/bash

for it in web codeground dashboard
do
	if ! [ -d ../$it ]
	then
		git clone git@github.com:coduno/${it}.git ../$it
	fi

	ln -s $PWD/../$it $PWD/$it
done
