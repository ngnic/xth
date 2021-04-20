#!/bin/bash

set -e

mkdir ./tmp 
cd ./tmp
cp ../docker-compose.production.yml docker-compose.yml
zip -r $(git rev-parse HEAD).zip docker-compose.yml 
aws s3 cp $(git rev-parse HEAD).zip s3://elasticbeanstalk-ap-southeast-1-461960015384/$(git rev-parse HEAD).zip 
aws elasticbeanstalk create-application-version \            
	--application-name main \            
	--version-label "$(git rev-parse HEAD)" \            
	--description "$(git log -1 --pretty=%s | cut -c -200)" \            
	--source-bundle S3Bucket="elasticbeanstalk-ap-southeast-1-461960015384",S3Key="$(git rev-parse HEAD).zip
cd .. && rm -rf ./tmp
