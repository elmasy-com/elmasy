#!/bin/bash

###
# Install elmasy.tar on a Debian server
###

if [ $EUID != "0" ]
then
    echo "RUN AS ROOT!"
    exit 1
fi

# Remove leftover file, if error occured before
if [ -e elmasy.tar ]
then
    rm elmasy.tar
fi

updateself() {

   
    OUTPUT=$(wget 'https://raw.githubusercontent.com/elmasy-com/elmasy/main/scripts/install.sh' 2>&1)
    if [ $? != 0 ]
    then
        echo "Failed to download install.sh!"
        echo "$OUTPUT"
        exit 1
    fi
}

install() {

    # Create "elmasy" user, if not created before 
    if [ $(id -u elmasy > /dev/null 2>&1; echo $?) != "0" ]
    then
        echo "Creating elmasy user..."
        adduser --no-create-home --gecos "" --disabled-login elmasy
    fi


    OUTPUT=$(wget "https://github.com/elmasy-com/elmasy/raw/main/build/elmasy.tar" 2>&1)
    if [ $? != 0 ]
    then
        echo "Failed to download elmasy.tar!"
        echo "$OUTPUT"
        exit 1
    fi

    # WSL does not use systemd
    if [[ $(systemctl 2>&1) != *"Failed to connect to bus: Host is down"* ]]
    then

        # Check if Elmasy installed before, and stop the running service
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
    fi

    if [ -e "/opt/elmasy-old" ]
    then
        echo "Removing elmasy-old..."
        rm -r "/opt/elmasy-old"
    fi

    if [ -e /opt/elmasy ]
    then
        echo "Moving old files to /opt/elmasy-old ..."
        mv /opt/elmasy /opt/elmasy-old
    fi

    echo "Extracting new files..."
    tar -xf elmasy.tar -C /opt

    if [ -e /opt/elmasy-old ]
    then
        OUTPUT=$(diff /opt/elmasy/elmasy.conf /opt/elmasy-old/elmasy.conf)
        if [ $? != 0 ]
        then
            echo "Changes in the new config:"
            echo "New < | > old"
            echo "$OUTPUT"
        fi

        echo "Copy the old config file to its new place..."
        cp /opt/elmasy-old/elmasy.conf /opt/elmasy
    fi

    echo "Setting executable capabilities to allow raw socket..."
    setcap cap_net_admin,cap_net_raw=eip /opt/elmasy/elmasy

    # WSL does not use systemd
    if [[ $(systemctl 2>&1) != *"Failed to connect to bus: Host is down"* ]]
    then
        echo "installing the new service file..."
        cp /opt/elmasy/elmasy.service /lib/systemd/system/elmasy.service
        systemctl daemon-reload
    fi

    echo "Removing leftover elmasy.tar..."
    rm elmasy.tar
}