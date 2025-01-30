# The big side project

The purpose of this project is to learn a lot. For starters, I want to learn how to manage a Postgres database, either with Docker, or more likely with a Managed Database Provider (MDS) like Supabase or Cloud SQL. I also want to learn the differences between cloud providers, like Google Cloud, DigitalOcean, Railway, Fly.io, AWS, and more. Side goals to this include learning how to setup CI/CD so I can establish a full development cycle. I want to also setup a development environment, potentially with docker compose and a locally run postgres db instance running in a volume.

I'd like to learn about logging, analytics, payments, and authentication, but that might spill over into another project.

## Logs

### 29th-30th Jan 2025

So far, I've created a simple Go server and a minimal Dockerfile. I have built the docker image and ran the container locally. I have spent most of time though learning about Docker, how to run a Postgres database in production, different cloud providers for running my Go backend server, how I can setup proper log aggregation, and other accessory concepts. What I've realised is that as long as I have my app containerised with Docker, migrating to another cloud provider probably won't be too difficult, as long as I don't depend too heavily on cloud services.

My plan currently is to use Railway to spin up the Postgres database for production - if I need to migrate off, there is a standard approach to do so. If I ever do need to migrate off Railway, I would use DigitalOcean, because that is where I will also be hosting my Go backend server. I think this combination works well - while initially the database and server are not close to each other or part of the same provider, this initial setup is free, easy to get started, and the progression from that point onwards is also pretty straightforward. It's also better to not put all your eggs in one basket. Migrating the go backend server is easy, and the database should be too, but only time will tell if I stick with this hosting strategy.
