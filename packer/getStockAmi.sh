aws ssm get-parameters --names /aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2 | jq -r '.Parameters | last(.[]).Value'
