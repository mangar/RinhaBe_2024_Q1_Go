#!/usr/bin/bash

cp ./dockerfile ../../crebito-fiber

cd ../../crebito-fiber

docker build . -t mangar/rinhabe_2024_q1_go:0.0.1

rm dockerfile
