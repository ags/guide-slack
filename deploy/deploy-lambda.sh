#!/bin/bash

set -euo pipefail

export AWS_DEFAULT_REGION=ap-southeast-2
export AWS_REGION=ap-southeast-2

lambdaBucketName="ags-${AWS_REGION}-lambda-fns"
stackName=guide-slack-slash-fn
hostedZoneName="sloth.cc."
domainName="guide.sloth.cc"
certificateARN="arn:aws:acm:us-east-1:805705958857:certificate/bcbbcd4f-8539-464f-a31a-f1118eb6799c"

echo "--- Building"
make clean all

echo "--- Uploading package to ${lambdaBucketName}"
aws cloudformation package \
  --template-file deploy/cfn/lambda.sam.yml \
  --s3-bucket $lambdaBucketName \
  --output-template-file tmp/packaged-template.yml

echo "--- Deploying ${stackName}"
aws cloudformation deploy \
  --capabilities CAPABILITY_IAM \
  --stack-name $stackName \
  --template-file tmp/packaged-template.yml \
  --parameter-overrides \
    "CertificateARN=${certificateARN}" \
    "DomainName=${domainName}" \
    "GuideAPIKey=${APP_GUIDE_API_KEY}" \
    "HostedZoneName=${hostedZoneName}" \
    "RollbarToken=${APP_ROLLBAR_TOKEN}" \
    "SlackSigningSecret=${APP_SLACK_SIGNING_SECRET}"

echo "--- Cleaning up"
make clean
