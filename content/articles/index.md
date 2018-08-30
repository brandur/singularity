Software is all around us. It powers the channels we use to
communicate with friends and family, the grids that
energize our cities, the transport networks that underpin
our civilization, the probes that we send into deep space,
and more and more our very lives and habits. With every
passing day it eats a little more of our world.

Software is often an intricately-built machine that's made
up of tens of thousands or even millions of lines of code.
This intricacy leads to complexity, and complexity leads to
fragility. We're used to software handling the most
difficult workloads that we see in our lives, but we're
also used to it being notoriously unreliable. Everyone from
the layman to the most hardened engineers knows to expect
bugs in everything from the simplest todo list app all the
way up to the voting machines that power our democracies.

Sometimes software launches with bugs. Sometimes bugs
appear at the edges that a developer didn't expect.
Sometimes bugs are introduced as updates break existing
features in subtle ways. No matter their origin, they're
everywhere and appearing all the time. Experience helps
produce software that's more reliable on average, but it's
never enough to stem the flow completely.

We can do better. Eliminating all bugs is impossible, but a
sizable number are introduced through reckless practices,
bad tools, and poor design. These broad classes of bugs are
avoidable, and we should aim to avoid them even if it means
shifting firmly-concreted mindsets or reevaluating core
beliefs. The _Self-hosting Singularity_ aims to be a guide
on how to do so.

It will briefly cover some ideas around _why_ software
fails then talk about values that we should aim to achieve.
The work's main body will talk about better design of
software through each aspect of its development: the
[development of software](#develop-soundly) that will
resist failure and entropy, [architecting
systems](#architect-wisely) more likely to succeed, and
[operating those systems](#operate-defensively) in ways to
maximize stability.

<div class="spacer"></div>

## Risks (#risks)

TODO

### Time (#time)

TODO

### Complexity (#complexity)

TODO

### Entropy (#entropy)

TODO

### Bitrot (#bitrot)

TODO

## Principles (#principles)

Abstract concepts.

### The walk away test (#walk-away-test)

Most uninitiated users of Intenet services take stability
for granted. A website is a website, and opening amazon.com
looks no different than janeshomepage.com. It's assumed
that keeping any websites online is a trivial task. Anyone
who's seen the other side can tell you that this is
generally not even remotely true. The larger a service, the
closer it operates to the brink of instability. Keeping
many sites online is a constant unseen battle, with
operations specialists receiving pages and taking care of
problems before they become serious enough to cause a
significant service disruption.

The _Walk Away Test_ is a thought experiment that runs as
follows: if every engineer and operations person at an
organization walked away from their web service today, how
long would it be before it went down? Assume that "walking
away" is total; every technical person has abandoned their
laptops and turned off their pagers.

There are many kinds of failures that can readily take down
a large system:

* A catastrophic event occurs in a critical component like
  a database or server that automated systems cannot
  compensate for.
* An externally-induced incident occurs that brings the
  system to its knees because no one is there to counteract
  for it. A DDOS attack (or its more legitimate equivalent,
  a user running an aggressive load test) is a good example
  of this class of problem.
* Small failures or problems continue to accumulate over
  time until in aggregate their large enough to bring down
  a service.
* The absence of constant development activity allows
  problems that have always been there to manifest. For
  example, a memory leak in a service has been hidden
  because developers are constantly deploying and
  restarting it.

Most operations people would like to think otherwise, but
for many large services, the _Walk Away Test_ would have a
result measured in days. This often includes services who
evangelize throughout the operations community and whose
ideas could generally be considered "best practices" (in
reality, oftentimes what's preached is not actually what's
practiced). If this claim seems incredible, the next time
you meet an engineer who's worked on a core AWS service,
ask them what their pager burden was like.

Somewhat unintuitively, smaller websites and services have
a much better chance of doing well on the _Walk Away Test_.
Between their minimal load, minimal architecture, dead
simple deployment (say a set of HTML files stored in S3),
heavy reliance on external operations services (like
Heroku), it's quite possible for them to stay online even
if their owners are absent for years.

### The skeleton crew (#skeleton-crew)

A basic question

How many people are spending most of their days staving off
disasters versus working on new things?

### The five year test (#five-year)

Clone down a repository five years from now and try to
build it. What's the likelihood that it will succeed?

## Develop soundly (#develop-soundly)

### Use types (#types)

TODO

### Build in layers (#layers)

Packages etc.

TODO

### Use languages that are memory safe (#memory-safe)

### Don't write new software (#new-software)

TODO

### Don't fork software (#forked-software)

TODO

### Use well-maintained software (#maintained-software)

TODO

## Architect wisely (#architect-wisely)

TODO

### Run less (#run-less)

TODO

### Embrace ephemerality (#ephemerality)

Use services. Don't do anything yourself if you can avoid
it.

* Use AWS instead of your own data center.
* Use GitHub instead of Phab.
* Use Rollbar instead of Sentry.
* Use DataDog instead of Graphite.

It's an easy conceptual pitfall to think that running
something yourself will be easy and solve all your
problems. Even if it easy after the initial install,
consider how it gets upgrade and maintained over the next
five years.

Engineers should have close to an allergic reaction when
somebody suggests running a new type of component in-house,
even if it's a technology that's exciting and known to be
mostly stable (e.g. Kafka). There may be a time where it's
appropriate to do so eventually as your organization and
uses reach colossal scale, but stave that off as long as
possible.

### Reduce moving parts (#moving-parts)

TODO

### Use static sites (#static-sites)

Static sites may be the ultimate form of minimalism when it
comes to web services. They're cheap to run computationally
and will handle even the largest volumes of traffic without
batting an eye. They're also ideal for horizontal scaling
which makes it especially easy to introduce redundancy into
the system, and perfect for global distribution.

They're obviously not suitable for many cases in that core
services are likely going to dynamic responses, but even
the most database-driven service out there can still have
its company's blog, marketing pages, and status site easily
made into 

### Design services in moderation (#services)

TODO

### Use ACID databases (#acid)

TODO

### Use relational databases (#relational)

TODO

## Operate defensively (#operate-defensively)

### Manage lifespans (#manage-lifespans)

TODO

### Inject chaos (#inject-chaos)

* The Netflix chaos monkey.

## Final words (#final-words)

<!--
# vim: set tw=59:
-->
