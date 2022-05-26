#!/bin/bash

set -e

help() {
	echo "build / pack / deploy / clean / run"
}

clean() {
        if [ -d "build/elmasy" ]
        then
                rm -rf "build/elmasy"
        fi
     
	 	if [ -f build/elmasy*.tar ]
		then
        	rm build/elmasy*.tar
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

	echo "Building frontend..."

	cd frontend
	
	OUTPUT=$(npm install 2>&1)
	if [ $? != 0 ]
	then
		echo "Failed to install npm dependencies!"
		echo "$OUTPUT"
		exit 1
	fi

	OUTPUT=$(ng build --prod --output-path=../build/elmasy/static 2>&1)
	if [ $? != 0 ]
	then
		echo "Failed to build elmasy!"
		echo "$OUTPUT"
		exit 1
	fi

	cd ..
}

elmasy-doc() {

	echo "Copying swagger.yaml..."
	cp api/swagger.yaml build/elmasy/static/

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
	tar -cf elmasy_$(git log -n 1 --pretty=format:"%H").tar elmasy/
	rm -rf elmasy/
	cd ..
}

deploy() {
	pack

	echo "Deploying..."
	bash ignore/scripts/deploy.sh
}

run() {
	build
	cd elmasy
	./elmasy
}

$1
