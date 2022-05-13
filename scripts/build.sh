#!/bin/bash

set -e

help() {
	echo "build / pack / clean / run"
}

clean() {
        if [ -d "build/elmasy" ]
        then
                rm -rf "build/elmasy"
        fi

		

		if [ -d "build/elmasy.com" ]
        then
                rm -rf "build/elmasy.com"
        fi

        if [ -f "build/elmasy.tar" ]
        then
                rm "build/elmasy.tar"
        fi
}


elmasy-dir() {

	echo "Creating required directories..."

	mkdir -p "build/elmasy/static"
}

elmasy-bin() {

	elmasy-dir

	echo "Building elmasy..."


	OUTPUT=$(go build -o build/elmasy/elmasy cmd/elmasy/main.go)
	if [ $? != 0 ]
	then
		echo "Failed to build elmasy!"
		echo "$OUTPUT"
		exit 1
	fi
}

elmasy-frontend() {

	echo "Downloading frontend..."

	OUTPUT=$(git clone --recursive https://github.com/elmasy-com/elmasy.com build/elmasy.com 2>&1)
	if [ $? != 0 ]
	then
		echo "Failed to clone https://github.com/elmasy-com/elmasy.com"
		echo "$OUTPUT"
		exit 1
	fi

	echo "Building frontend..."
	OUTPUT=$(cd build/elmasy.com && hugo -d ../elmasy/static && cd ../.. 2>&1)
	if [ $? != 0 ]
	then
		echo "Failed to build frontend"
		echo "$OUTPUT"
		exit 1
	fi

	rm -rf "build/elmasy.com/"
}

elmasy-doc() {

	echo "Copying swagger-ui..."

	cp -r web/swagger-ui/ build/elmasy/static/
	cp api/swagger.yaml build/elmasy/static/swagger-ui/

}

elmasy-other() {

	echo "Copying other files..."

	cp configs/elmasy.conf build/elmasy/elmasy.conf.example
	cp init/elmasy.service build/elmasy/
}


build() {

	clean

	elmasy-bin
	elmasy-frontend
	elmasy-doc
	elmasy-other
}

pack() {
	build

	echo "Creating tar archive..."
	cd build
	tar -cf elmasy.tar elmasy/
	rm -rf elmasy/
	cd ..
}

run() {
	build
	cd elmasy
	./elmasy
}

$1
