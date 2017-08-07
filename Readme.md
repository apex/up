![](https://dl.dropboxusercontent.com/u/6396913/Apex/Up/Readme/title-fs8.png)

Up deploys infinitely scalable serverless apps, APIs, and static websites in seconds, abstracting away complex infrastructure so you can get back to writing code, all while providing cost effective, scalable, and global services.

Up currently supports Node.js, Golang, Python, Crystal, and static sites out of the box. Up currently targets AWS Lambda and API Gateway as its platform, however more will be available in the future, you can think of Up as a serverless provider-agnostic Heroku style experience.

Check out some of the [examples](https://github.com/apex/up-examples) to get started.

![](https://dl.dropboxusercontent.com/u/6396913/Apex/Up/Readme/screen-koa-fs8.png)

## Features

Open source community edition: Coming soon.

![](https://dl.dropboxusercontent.com/u/6396913/Apex/Up/Readme/up-features-community-fs8.png)

## Pro Features

Close sourced pro edition: Coming less soon.

![](https://dl.dropboxusercontent.com/u/6396913/Apex/Up/Readme/up-features-pro-fs8.png)

## Pricing

Updated as of July 2017 based on public information. Some services offer a restricted free version, or free access for solo developers – this table is based on commercial use.

![pricing table](https://dl.dropboxusercontent.com/u/6396913/Apex/Up/Readme/pricing.png)

## FAQ

<details>
  <summary>Is this a hosted service?</summary>
  <p>There are no plans for a hosted version. Up lets you deploy applications to your own AWS account for isolation, security, and longevity, don't worry about a startup going out of business.</p>
</details>

<details>
  <summary>Why isn't Up licensed MIT?</summary>
  <p>Up is licensed in such a way that myself as an independent developer can continue to improve the product and provide support. Commercial customers receive access to a premium version of Up with additional features, priority support for bugfixes, and of course knowing that the project will stick around! Up saves your team countless hours maintaining infrastructure and custom tooling, so you can get back to what makes your company and products unique.</p>
</details>

<details>
  <summary>How is this different than other serverless frameworks?</summary>
  <p>Most of the AWS Lambda based tools are function-oriented, while Up abstracts this away entirely. Up does not use framework "shims", the servers that you run using Up are regular HTTP servers and require no code changes for Lambda compatibility.</p>
</details>

<details>
  <summary>How much does it cost to run an application?</summary>
  <p>AWS API Gateway provides 1 million free requests per month, so there's a good chance you won't have to pay anything at all. Beyond that view the <a href="https://aws.amazon.com/api-gateway/pricing/">AWS Pricing</a> for more information.</p>
</details>

<details>
  <summary>Do I have to manage instances or container counts?</summary>
  <p>Nope! Up scales to fit your traffic on-demand, you don't have to do anything beyond deploying your code.</p>
</details>

<details>
  <summary>How much latency does Up's reverse proxy introduce?</summary>
  <p>With a 512mb Lambda function Up introduces an average of around 500µs (microseconds) per request.</p>
</details>

<details>
  <summary>Can I donate?</summary>
  <p>I'm glad you asked! Yes you can, head over to the <a href="https://opencollective.com/apex">OpenCollective</a> page. Any donations are greatly appreciated and help me focus more on Up's implementation, documentation, and examples.</p>
</details>

## Community

- [Example applications](https://github.com/apex/up-examples) for Up
- [Slack](https://apex-dev.azurewebsites.net/) to chat with apex(1) and up(1) community members
- [Blog](https://blog.apex.sh/) to follow release posts, tips and tricks

<a href="https://apex.sh"><img src="http://tjholowaychuk.com:6000/svg/sponsor"></a>
