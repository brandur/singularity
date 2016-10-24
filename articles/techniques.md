

## Run Less (#run-less)

## Reduce Moving Parts (#moving-parts)

## Embrace Ephemerality (#ephemerality)

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

## Reduce Moving Parts

### Static Sites

Static sites may be the ultimate form of minimalism when it comes to web
services. They're cheap to run computationally and will handle even the largest
volumes of traffic without batting an eye. They're also ideal for horizontal
scaling which makes it especially easy to introduce redundancy into the system.

They're obviously not suitable for many cases in that core services are likely
going to dynamic responses, but even the most database-driven service out there
can still have its company's blog, marketing pages, and status site easily made
into 


### Slinky (#slink)

## Inject Chaos (#inject-chaos)

* The Netflix chaos monkey.

## Rational Microservices

## Don't Write New Software
