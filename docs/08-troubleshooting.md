---
title: Help
---

<details>
  <summary>I didn't receive a sign-in or certificate confirmation email</summary>
  <p>AWS email delivery can be slow sometimes. Please give it 30-60s. Otherwise, be sure to check your spam folder.</p>
</details>

<details>
  <summary>My application times out or seems slow</summary>
  <p>Lambda `memory` scales CPU alongside RAM, so if your application is slow to initialize or serve responses, you may want to try `1024` or above. See [Lambda Pricing](https://aws.amazon.com/lambda/pricing/) for options.</p>
  <p>Ensure that all of your dependencies are deployed. You may use `up -v` to view what is added or filtered from the deployment or `up build --size` to output the contents of the zip.</p>
</details>

<details>
  <summary>I'm seeing 404 Not Found responses</summary>
  <p>By default, Up ignores files which are found in `.upignore`. Use the verbose flag such as `up -v` to see if files have been filtered or `up build --size` to see a list of files within the zip sorted by size. See [Ignoring Files](#configuration.ignoring_files) for more information.</p>
</details>

<details>
  <summary>My deployment seems stuck</summary>
  <p>The first deploy also creates resources associated with your project and can take roughly 1-2 minutes. AWS provides limited granularity into the creation progress of these resources, so the progress bar may appear "stuck".</p>
</details>

<details>
  <summary>How do I sign into my team?</summary>
  <p>Run `up team login` if you aren't signed in, then run `up team login --team my-team-id` to sign into any teams you're an owner or member of.</p>
</details>

<details>
  <summary>Unable to associate certificate error</summary>
  <p>If you receive a `Unable to associate certificate` error it is because you have not verified the SSL certificate. Certs for CloudFront when creating a custom domain MUST be in us-east-1, so if you need to manually resend verification emails visit [ACM in US East 1](https://console.aws.amazon.com/acm/home?region=us-east-1).</p>
</details>

<details>
  <summary>I'm seeing 403 Forbidden errors in CI</summary>
  <p>If you run into "403 Forbidden" errors this is due to GitHub's low rate limit for unauthenticated users, consider creating a [Personal Access Token](https://github.com/settings/tokens) and adding `GITHUB_TOKEN` to your CI.</p>
</details>
