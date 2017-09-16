---
title: FAQ
---

<details>
  <summary>Is this a hosted service?</summary>
  <p>There are currently no plans for a hosted version. Up lets you deploy applications to your own AWS account for isolation, security, and longevity, don't worry about a startup going out of business.</p>
</details>

<details>
  <summary>What platforms does Up support?</summary>
  <p>Currently AWS via API Gateway and Lambda are supported, this is the focus until Up is nearing feature completion, after which additional providers such as GCP and Azure will be added.</p>
</details>

<details>
  <summary>How is this different than other serverless frameworks?</summary>
  <p>Most of the AWS Lambda based tools are function-oriented, while Up abstracts this away entirely. Up does not use framework "shims", the servers that you run using Up are regular HTTP servers and require no code changes for Lambda compatibility.</p>

  <p>Up keeps your apps and APIs portable, makes testing them locally easier, and prevents vendor lock-in. The Lambda support for Up is simply an implementation detail, you are not coupled to API Gateway or Lambda. Up uses the API Gateway proxy mode to send all requests (regardless of path or method) to your application.</p>

  <p>If you're looking to manage function-level event processing pipelines, Apex or Serverless are likely better candidates, however if you're creating applications, apis, micro services, or websites, Up is built for you.</p>
</details>

<details>
  <summary>Why run HTTP servers in Lambda?</summary>
  <p>You might be thinking this defeats the purpose of Lambda, however most people just want to use the tools they know and love. Up lets you be productive developing locally as you normally would, Lambda for hosting is only an implementation detail.</p>

  <p>With Up you can use any Python, Node, Go, or Java framework you'd normally use to develop, and deploy with a single command, while maintaining the cost effectiveness, self-healing, and scaling capabilities of Lambda.</p>
</details>

<details>
  <summary>How much does it cost to run an application?</summary>
  <p>AWS API Gateway provides 1 million free requests per month, so there's a good chance you won't have to pay anything at all. Beyond that view the <a href="https://aws.amazon.com/api-gateway/pricing/">AWS Pricing</a> for more information.</p>
</details>

<details>
  <summary>How well does it scale?</summary>
  <p>Up scales to fit your traffic on-demand, you don't have to do anything beyond deploying your code. There's no restriction on the number of concurrent instances, apps, custom domains and so on.</p>
</details>

<details>
  <summary>How much latency does Up's reverse proxy introduce?</summary>
  <p>With a 512mb Lambda function Up introduces an average of around 500µs (microseconds) per request.</p>
</details>

<details>
  <summary>Do the servers stay active while idle?</summary>
  <p>This depends on the platform, and with Lambda being the initial platform provided the current answer is no, the server(s) are frozen when inactive and are otherwise "stateless".</p>

  <p>Typically relying on background work in-process is an anti-pattern as it does not scale. Lambda functions combined with CloudWatch scheduled events for example are a good way to handle this kind of work, if you're looking for a scalable alternative.</p>
</details>

<details>
  <summary>Why is Up licensed as GPLv3?</summary>
  <p>Up is licensed in such a way that myself as an independent developer can continue to improve the product and provide support. Commercial customers receive access to a premium version of Up with additional features, priority support for bugfixes, and of course knowing that the project will stick around! Up saves your team countless hours maintaining infrastructure and custom tooling, so you can get back to what makes your company and products unique.</p>
</details>

<details>
  <summary>Can I donate?</summary>
  <p>Yes you can! Head over to the <a href="https://opencollective.com/apex-up">OpenCollective</a> page. Any donations are greatly appreciated and help me focus more on Up's implementation, documentation, and examples. If you're using the free OSS version for personal or commercial use please consider giving back, even a few bucks buys a coffee :).</p>
</details>
