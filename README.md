Kpas allows engineers who don't have time for devops to create kubernetes clusters across providers and deploy applications to them in just a couple of commands.

In just two commands kpas provides a fully working cluster with CI, git push deployment along with log aggregation and an familiar CLI for interacting (think logs, consoles etc) with running applications.

## What is it

Kpas (a riff on Kubernetes Platform As a Service) provides a simple CLI for creating standardised clusters across cloud, VM and Bare Metal providers.

It also provides an optional, opinionated layer on top of these cluster using best in class open source software which includes:

- In-cluster Registry
- CI and Git Server
- Nginx Ingress & SSL with Lets Encrypt
- Centralised Logging & Metrics
- A CLI for generating clusters and adding apps

The goal is to provide an opinionated, out of the box "release ready" configuration for a typical modern web application stack. A typical workflow is:

- Provision a cluster with Kpas
- Initialize your app for Kpas deployment using the CLI
- `git push` your web application to have it automatically deployed

## Why does Kpas exist

Kubernetes is rapidly becoming the common runtime for the internet. The benefits of a single abstraction like this are already being realised in the Enterprise. Kpas aims to make the efficiency of a common runtime for everything available to everyone from individual developers wanting to efficiently run lots of side projects to startups wanting to put a solid foundation for future growth in place.