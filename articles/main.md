# Foreword

## Critical Thought (#critical-thought)

!section class="col-style-1"

Firstly, it's well worthwhile reminding readers of the important of critical
thought. This is apparently a platitude of the highest order, but in practice,
technical articles are one place where otherwise great critical thinkers often
forget to apply the technique. There was a time when Node was billed as the
saviour of the engineering world. Before that, Ruby occupied a similar
heroic position. Prior to that, C++.

Historical context reveals that none of these technologies turned out to be absolute solutions.

Like every other technical publication ever written, this one doesn't have all the
answers, and contains a healthy dose of idealism. What's preached isn't
necessarily what's practiced.

!/section

## Application

!section class="col-style-2"

The ideas contained in this publication are not necessarily universally
applicable. There are cases where services will need to be close to hardware,
require very manual operation, or have restrictions in place for compliance
reasons.

However, a cautious architect would be wise not to discount anything wholesale.
Re-examing assumptions that are firmly concreted into the shared mental
landscape at any organization may reveal that 

!/section

# Introduction

!section class="col-style-1"

Abstract concepts.

!/section

## The Walk Away Test (#walk-away-test)

!section class="col-style-1"

Most uninitiated users of Intenet services take stability for granted. A
website is a website, and opening amazon.com looks no different than
janeshomepage.com. It's assumed that keeping any websites online is a trivial
task. Anyone who's seen the other side can tell you that this is generally not
even remotely true. The larger a service, the closer it operates to the brink
of instability. Keeping many sites online is a constant unseen battle, with
operations specialists receiving pages and taking care of problems before they
become serious enough to cause a significant service disruption.

The _Walk Away Test_ is a thought experiment that runs as follows: if every
engineer and operations person at an organization walked away from their web
service today, how long would it be before it went down? Assume that "walking
away" is total; every technical person has abandoned their laptops and turned
off their pagers.

There are many kinds of failures that can readily take down a large system:

* A catastrophic event occurs in a critical component like a database or server
  that automated systems cannot compensate for.
* An externally-induced incident occurs that brings the system to its knees
  because no one is there to counteract for it. A DDOS attack (or its more
  legitimate equivalent, a user running an aggressive load test) is a good
  example of this class of problem.
* Small failures or problems continue to accumulate over time until in
  aggregate their large enough to bring down a service.
* The absence of constant development activity allows problems that have always
  been there to manifest. For example, a memory leak in a service has been
  hidden because developers are constantly deploying and restarting it.

Most operations people would like to think otherwise, but for many large
services, the _Walk Away Test_ would have a result measured in days. This
often includes services who evangelize throughout the operations community and
whose ideas could generally be considered "best practices" (in reality,
oftentimes what's preached is not actually what's practiced). If this claim
seems incredible, the next time you meet an engineer who's worked on a core AWS
service, ask them what their pager burden was like.

Somewhat unintuitively, smaller websites and services have a much better chance
of doing well on the _Walk Away Test_. Between their minimal load, minimal
architecture, dead simple deployment (say a set of HTML files stored in S3),
heavy reliance on external operations services (like Heroku), it's quite
possible for them to stay online even if their owners are absent for years.

!/section

## The Skeleton Crew (#skeleton-crew)

!section class="col-style-1"

A basic question

How many people are spending most of their days staving off disasters versus
working on new things?

!/section

# Risks

!section class="col-style-1"

TODO

!/section

## Time (#time)

!section class="col-style-1"

TODO

!/section

## Entropy (#entropy)

!section class="col-style-1"

TODO

!/section

## Bitrot (#bitrot)

!section class="col-style-1"

TODO

!/section

# Techniques

!section class="col-style-1"

TODO

!/section

## Run Less (#run-less)

!section class="col-style-1"

TODO

!/section

## Embrace Ephemerality (#ephemerality)

!section class="col-style-1"

Use services. Don't do anything yourself if you can avoid it.

* Use AWS instead of your own data center.
* Use GitHub instead of Phab.
* Use Rollbar instead of Sentry.
* Use DataDog instead of Graphite.

It's an easy conceptual pitfall to think that running something yourself will
be easy and solve all your problems. Even if it easy after the initial install,
consider how it gets upgrade and maintained over the next five years.

Engineers should have close to an allergic reaction when somebody suggests
running a new type of component in-house, even if it's a technology that's
exciting and known to be mostly stable (e.g. Kafka).

!/section

## Reduce Moving Parts (#moving-parts)

!section class="col-style-1"

TODO

!/section

### Static Sites (#static-sites)

!section class="col-style-1"

Static sites may be the ultimate form of minimalism when it comes to web
services. They're cheap to run computationally and will handle even the largest
volumes of traffic without batting an eye. They're also ideal for horizontal
scaling which makes it especially easy to introduce redundancy into the system.

They're obviously not suitable for many cases in that core services are likely
going to dynamic responses, but even the most database-driven service out there
can still have its company's blog, marketing pages, and status site easily made
into 


!/section

### Slinky (#slinky)

!section class="col-style-1"

TODO

!/section

## Inject Chaos (#inject-chaos)

!section class="col-style-1"

* The Netflix chaos monkey.

!/section

## Rational Microservices (#rational-microservices)

!section class="col-style-1"

TODO

!/section

## Don't Write New Software (#new-software)

!section class="col-style-1"

TODO

!/section

### Don't Fork Existing Software (#forked-software)

!section class="col-style-1"

TODO

!/section

# Implementation

!section class="col-style-1"

TODO

!/section

## Dynamic Languages (#dynamic-languages)

!section class="col-style-1"

Consider not using them.

!/section
