---
title: Help
---

<details>
  <summary>I didn't receive a sign-in or certificate confirmation email</summary>
  <p>AWS email delivery can sometimes be slow, please give it 30-60s and check the spam folder.</p>
</details>

<details>
  <summary>My application times out</summary>
  <p>Lambda `memory` scales CPU alongside the RAM – so if your application is booting or serving responses slowly you may want to try `1024` or above, see [Lambda Pricing](https://aws.amazon.com/lambda/pricing/) for options.</p>
</details>

<details>
  <summary>My deployment seems stuck</summary>
  <p>The first deploy also creates resources associated with your project, and can take roughly 1-2 minutes. AWS provides limited granularity into the creation progress of these resources, so the progress bar may appear "stuck".</p>
</details>

<details>
  <summary>How do I get my team back?</summary>
  <p>Run `up team login` if you aren't signed in, then run `up team login --team my-team-id` to sign into any teams you're an owner or member of.</p>
</details>
