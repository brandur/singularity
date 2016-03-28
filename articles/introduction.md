# Introduction

Abstract concepts.

## The Walk Away Test

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

## The Skeleton Crew
