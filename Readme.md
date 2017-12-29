![](assets/title.png)

Up deploys infinitely scalable serverless apps, APIs, and static websites in seconds, so you can get back to working on what makes your product unique.

Up focuses on deploying "vanilla" HTTP servers so there's nothing new to learn, just develop with your favorite existing frameworks such as Express, Koa, Django, Golang net/http or others.

Up currently supports Node.js, Golang, Python, Java, Crystal, and static sites out of the box. Up is platform-agnostic, supporting AWS Lambda and API Gateway as the first targets. You can think of Up as self-hosted Heroku style user experience for a fraction of the price, with the security, isolation, flexibility, and scalability of AWS.

Check out the [documentation](https://up.docs.apex.sh/) for more instructions and links, or try one of the [examples](https://github.com/apex/up-examples), or chat with us in [Slack](https://chat.apex.sh/).

![](assets/screen.png)

## Features

Open source community edition.

![Open source edition features](assets/features-community.png)

## Pro Features

Up Pro is **$20/mo USD** for unlimited use within your company, with no additional cost per team member. Up Pro is currently in early-access alpha, please use the **up-early-adopter-57AAA8693354** coupon for 50% off indefinitely. Head over to [Subscribing to Up Pro](https://up.docs.apex.sh/#guides.subscribing_to_up_pro) to get started.

Note that the following Pro features are currently available:

 - Encrypted env variables
 - Alerting (email, sms, slack)

![Pro edition features](assets/features-pro.png)

## Quick Start

Install Up:

```
$ curl -sf https://up.apex.sh/install | sh
```

Tell up which AWS profile to use:

```
export AWS_PROFILE=example
```

Create an `app.js` file:

```js
require('http').createServer((req, res) => {
  res.end('Hello World\n')
}).listen(process.env.PORT)
```

Deploy the app:

```
$ up
```

Open it in the browser:

```
$ up url --open
```

## Donations

We also welcome financial contributions for the open-source version on [Open Collective](https://opencollective.com/apex-up). Your contributions help keep this project alive!

### Sponsors

<a href="https://opencollective.com/apex-up/sponsor/0/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/0/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/1/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/1/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/2/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/2/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/3/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/3/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/4/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/4/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/5/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/5/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/6/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/6/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/7/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/7/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/8/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/8/avatar.svg"></a>
<a href="https://opencollective.com/apex-up/sponsor/9/website" target="_blank"><img src="https://opencollective.com/apex-up/sponsor/9/avatar.svg"></a>

### Backers

<a href="https://opencollective.com/apex-up#backers" target="_blank"><img src="https://opencollective.com/apex-up/backers.svg?width=890"></a>


<a href="https://apex.sh"><img src="http://tjholowaychuk.com:6000/svg/sponsor"></a>
