# Curto: fast url-shortner service

[![Coverage Status](https://coveralls.io/repos/github/mfvitale/curto/badge.svg?branch=main)](https://coveralls.io/github/mfvitale/curto?branch=main)

# Why 'curto'?

I developed this service as a side project with the main goal to learn Go lang. The name 'curto' in neapolitan meas short.


# Start the application

Before running the application you need to set these envs: REDIS_ENDPOINT, REDIS_USERNAME, REDIS_PASSWORD

For exmaple if you want to connect a local redis
```bash
export REDIS_ENDPOINT=localhost:6379
export REDIS_USERNAME="username"
export REDIS_PASSWORD="password"
```

you can the start application with
```bash
./run.sh
```
## Shorten an URL

```bash
curl -X GET "http://localhost:8080/encode?url=https://www.mfvitale.me"
```
response will be the shorten url

```bash
http://localhost:8080/eBIGnfbLW
```
*Note that the url returned uses the config 'domain' inside config.yml file. You can also specify it through 'DOMAIN' env*

## Get original url

```bash
curl -X GET "http://localhost:8080/eBIGnfbLW"
```
response will a '303 See Other' http status

```bash
<a href="https://www.mfvitale.me">See Other</a>.
```
# Try it
I have deployed the service on [Okteto](https://www.okteto.com/) so you can reach the service [here](https://curto-url-shortner-mfvitale.cloud.okteto.net/)
# Design choices

The service has been designed to be fast, scalable and cloud-native.

So chosed to use:
* [Snowflake algorithm*](https://betterprogramming.pub/implementing-snowflake-algorithm-in-golang-c1098fdc73d0) to have a fast and distributed id generation instead using Database auto-increment feature.
* Base62 conversion to avoid resusing know (but long resulting) hash function (SHA-1, CRC32) that requires to manage collision
* Use Redis to store <shortURL, longURL> mapping

****
## TLDR

The service has been implemented as an application and not an API, so it also manage the redirection for a short URL.

The most easy way to implement an URLShortner servive is to keep the pair <shortURL, longURL>

A short URL is like http://shortServiceUrl{hashValue} where '**hashValue**' is the result of an hash function.

A first approach is to use SHA-1 hash function that generates something like 'da39a3ee5e6b4b0d3255bfef95601890afd80709' but the problem is that is to long (eventually much longer than the longURL) so to cope with that we need to cut it to a defined size. This method can lead to hash collisons. I don't want to manage it.

Another approach is to use base 62 conversion. With this approach we convert a base 10 number to it's base 63 rapresentation. Where we can take this number? The number can be a generated ID.

Also for the numeric ID generation there are different approaches:

* use Database auto-increment feature
* use Snowflake

I discarded the first choice beacuse:
* Hard to scale with multiple data centers
* Don't want to setup a database only for generate IDs
* Single point of failure in case of not replicated Database
* It does not scale well when a server is added or removed ( in case of a multi-master replication)

In first instance I deployed the service with a K8s Deployment and as MachineId for the Snowflake algorithm I used the pid of the process running in the pod. The problem of this approach is that:
* different pod can have the same pid
* pid can change when pod restarted

So I decided to switch from Deployment to StatefulSet to have a persistent id for the pods. I used this id as MachineId.

# Improvements
* introduce a rate-limiter to prevent resource starvation
* use heml to deploy