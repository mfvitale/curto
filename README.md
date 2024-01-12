# Curto: fast url-shortner service

![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/mfvitale/curto/CI/main)
[![Coverage Status](https://coveralls.io/repos/github/mfvitale/curto/badge.svg?branch=main)](https://coveralls.io/github/mfvitale/curto?branch=main)

# Why 'curto'?

I developed this service as a side project with the main goal to learn Go lang. The name 'curto' in Neapolitan meas short.


# Start the application

Before running the application you need to set these envs: REDIS_ENDPOINT, REDIS_USERNAME, REDIS_PASSWORD, MACHINE_ID

For example if you want to connect to a local Redis
```bash
export REDIS_ENDPOINT=localhost:6379
export REDIS_USERNAME="username"
export REDIS_PASSWORD="password"
export MACHINE_ID=0
```

you can then start application with
```bash
./run.sh
```
## Shorten an URL

```bash
curl -X GET "http://localhost:8080/shorten?url=https://www.mfvitale.me"
```
the response will be the shortened URL

```bash
http://localhost:8080/eBIGnfbLW
```
*Note that the URL returned uses the config 'domain' inside config.yml file. You can also specify it through 'DOMAIN' env*

## Get original url

```bash
curl -X GET "http://localhost:8080/eBIGnfbLW"
```
the response will be a '303 See Other' HTTP status

```bash
<a href="https://www.mfvitale.me">See Other</a>.
```
### Kubernetes deploy
If you want to deploy it on k8s you need to create:
* A secret with name `redis-secrets` and add `REDIS_ENDPOINT`. Optionally you can add `REDIS_USERNAME` and `REDIS_PASSWORD`
* A config map with name `app-config` where you can specify the `domain` with your custom domain.
* Build the image locally naming it `mfvitale/curto:latest` 

You can just run `kubectl apply -f .deploy/`
# Try it
I have deployed the service on [Render](https://render.com/) so you can reach the service [here](https://curto.onrender.com/)
# Design choices

The service has been designed to be fast, scalable, and cloud-native.

So I chose to use:
* [Snowflake algorithm](https://betterprogramming.pub/implementing-snowflake-algorithm-in-golang-c1098fdc73d0) to have a fast and distributed id generation instead using Database auto-increment feature.
* Base62 conversion to avoid reusing know (but long resulting) hash function (SHA-1, CRC32) that requires managing collision
* Use Redis to store <shortURL, longURL> mapping

****
## TL;DR

The service has been implemented as an application and not an API, so it also manages the redirection for a short URL.

The easiest way to implement an URLShortner service is to keep the pair <shortURL, longURL>

A short URL is like http://shortServiceUrl{hashValue} where '**hashValue**' is the result of a hash function.

A first approach is to use SHA-1 hash function that generates something like 'da39a3ee5e6b4b0d3255bfef95601890afd80709' but the problem is that is too long (eventually much longer than the longURL) so to cope with that we need to cut it to a defined size. This method can lead to hash collisons. I don't want to manage it.

Another approach is to use base 62 conversions. With this approach, we convert a base 10 number to its base 62 representation. Where we can take this number? The number can be a generated ID.

Also for the numeric ID generation, there are different approaches:

* use the Database auto-increment feature
* use Snowflake

I discarded the first choice because:
* Hard to scale with multiple data centers
* Don't want to set up a database only to generate IDs
* Single point of failure in case of not replicated Database
* It does not scale well when a server is added or removed ( in case of a multi-master replication)

In the first instance, I deployed the service with a K8s Deployment and as MachineId for the Snowflake algorithm, I used the PID of the process running in the pod. The problem with this approach is that:
* different pods can have the same PID
* PID can change when the pod restarts

So I decided to switch from Deployment to StatefulSet to have a persistent id for the pods. I used this id as MachineId.

# Improvements
* introduce a rate-limiter to prevent resource starvation
* use helm to deploy
