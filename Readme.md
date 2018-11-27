![](assets/title.png)

Up deploys infinitely scalable serverless apps, APIs, and static websites in seconds, so you can get back to working on what makes your product unique.

With Up there's no need to worry about managing or scaling machines, paying for idle servers, worrying about logging infrastructure or alerting. Just deploy your app with `$ up` and you're done!

Use the free OSS version, or subscribe to [Up Pro](#pro-features) for a small monthly fee for unlimited use within your company, there is no additional cost per team-member or application. Deploy dozens or even hundreds of applications for pennies thanks to AWS Lambda's cost effective nature.

## About

Up focuses on deploying "vanilla" HTTP servers so there's nothing new to learn, just develop with your favorite existing frameworks such as Express, Koa, Django, Golang net/http or others.

Up currently supports Node.js, Golang, Python, Java, Crystal, Clojure and static sites out of the box. Up is platform-agnostic, supporting AWS Lambda and API Gateway as the first targets. You can think of Up as self-hosted Heroku style user experience for a fraction of the price, with the security, isolation, flexibility, and scalability of AWS.

Check out the [documentation](https://up.docs.apex.sh/) for more instructions and links, or try one of the [examples](https://github.com/apex/up-examples), or chat with us in [Slack](https://chat.apex.sh/).

![](assets/screen2.png)

## OSS Features

Features of the free open-source edition.

![Open source edition features](assets/features-community.png)

## Pro Features

Up Pro provides additional features for production-ready applications such as encrypted environment variables, error alerting, unlimited team members, unlimited applications and priority [email support](mailto:support@apex.sh) for **$20/mo USD**. Visit [Subscribing to Up Pro](https://up.docs.apex.sh/#guides.subscribing_to_up_pro) to get started.

![Pro edition features](assets/features-pro.png)

[![](https://gui.apex.sh/component?name=ShadowButton&config=%7B%22text%22%3A%22SUBSCRIBE%22%2C%22color%22%3A%227956EF%22%7D)](https://up.docs.apex.sh/#guides.subscribing_to_up_pro)

## Quick Start

Install Up:

```
$ curl -sf https://up.apex.sh/install | sh
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

Open it in the browser, or copy the url to your clipboard:

```
$ up url -o
$ up url -c
```

## Donations

We also welcome financial contributions for the open-source version on [Open Collective](https://opencollective.com/apex-up). Your contributions help keep this project alive!

### Sponsors

<a href="https://opencollective.com/apex-up#backers" target="_blank"><img src="https://opencollective.com/apex-up/tiers/sponsors.svg?avatarHeight=36&width=600" /></a>

### Backers

<a href="https://opencollective.com/apex-up#backers" target="_blank"><img src="https://opencollective.com/apex-up/tiers/backers.svg?avatarHeight=36&width=600" /></a>


<a href="https://apex.sh"><img src="http://tjholowaychuk.com:6000/svg/sponsor"></a>
