---
title: AWS Credentials
---

Before using Up you need to first provide your AWS account credentials so that Up is allowed to create resources on your behalf.

## AWS Credential Profiles

Most AWS tools support the `~/.aws/credentials` file for storing credentials, allowing you to specify `AWS_PROFILE` environment variable so Up knows which one to reference. To read more on configuring these files view [Configuring the AWS CLI](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html).

Here's an example of `~/.aws/credentials`, where `export AWS_PROFILE=myapp` would activate these settings.

```
[myapp]
aws_access_key_id = xxxxxxxx
aws_secret_access_key = xxxxxxxxxxxxxxxxxxxxxxxx
```

## Best Practices

You may store the profile name in the `up.json` file itself as shown in the following snippet:

```json
{
  "profile": "myapp"
}
```

This is ideal since it ensures that you do not accidentally have a different environment set and deploy to another account, potentially to an app of the same name, as apps are account-specific you may have one named "api" in many.

## IAM Policy for Up CLI

Below is a policy for [AWS Identity and Access Management](https://aws.amazon.com/iam/) which provides Up access to manage your resources. Note that the policy may change as features are added to Up, so you may have to adjust the policy.

If you're using Up for a production application it's highly recommended to configure an IAM role and user(s) for your team, restricting the access to the account and its resources.

<details>
  <summary>Show policy</summary>
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "acm:*",
                "cloudformation:Create*",
                "cloudformation:Delete*",
                "cloudformation:Describe*",
                "cloudformation:ExecuteChangeSet",
                "cloudformation:Update*",
                "cloudfront:*",
                "cloudwatch:*",
                "ec2:*",
                "ecs:*",
                "events:*",
                "iam:AttachRolePolicy",
                "iam:CreatePolicy",
                "iam:CreateRole",
                "iam:DeleteRole",
                "iam:DeleteRolePolicy",
                "iam:GetRole",
                "iam:PassRole",
                "iam:PutRolePolicy",
                "lambda:AddPermission",
                "lambda:Create*",
                "lambda:Delete*",
                "lambda:Get*",
                "lambda:InvokeFunction",
                "lambda:List*",
                "lambda:RemovePermission",
                "lambda:Update*",
                "logs:Create*",
                "logs:Describe*",
                "logs:FilterLogEvents",
                "logs:Put*",
                "logs:Test*",
                "route53:*",
                "route53domains:*",
                "s3:*",
                "ssm:*"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": "apigateway:*",
            "Resource": "arn:aws:apigateway:*::/*"
        }
    ]
}
```
</details>
