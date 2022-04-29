#!/bin/bash

###
# Install elmasy.tar on a server
###

if [ $1 == "" ]
then
    echo "Usage: ./$0 /path/to/elmasy.tar"
    exit
elif [[ $(file $1) != *"tar archive"* ]]
then
    echo "$1 is not a TAR file!"
    exit 1
fi

if [ -e "/lib/systemd/system/elmasy.service" ]
then

    if [[ $(systemctl is-active elmasy.service) == "active" ]]
    then
        echo "Elmasy is running!"

        systemctl stop elmasy.service
        if [ $? != 0 ]
        then
            echo "Failed to stop Elmasy!"
            exit 1
        else
            echo "Elmasy stopped!"
        fi
    fi
fi

if [ -e "/opt/elmasy-old" ]
then
    rm -r "/opt/elmasy-old"
fi

echo "Moving old files to /opt/elmasy-old ..."
mv /opt/elmasy /opt/elmasy-old

echo "Extracting new files..."
tar -xf $1 -C /opt

cp /opt/elmasy-old/elmasy.conf /opt/elmasy
cp /opt/elmasy/elmasy.service /lib/systemd/system/elmasy.service
systemctl daemon-reload

systemctl start elmasy