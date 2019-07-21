#!/bin/sh -x

# Run RequestAddr at bootup
sudo mv /tmp/RequestShort /etc/init.d RequestShort
chmod +x /etc/init.d/RequestShort
