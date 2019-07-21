#!/bin/sh -x

# Alternate: Running Commands on Your Linux Instance at Launch
#   https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-cloud-init

# SSH key file configured through Terraform? cloud-init?

# Run RequestAddr at bootup
sudo mv /tmp/RequestAddr /etc/init.d RequestAddr
chmod +x /etc/init.d/RequestAddr
sudo systemctl enable /etc/init.d/RequestAddr

# Remove data scripts before creating AMI
# rm /var/lib/cloud/instances/instance-id/*

sudo systemctl enable /etc/init.d/RequestShort

