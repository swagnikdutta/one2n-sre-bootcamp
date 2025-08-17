This is a basic crud app — a part of the [one-2n-sre-bootcamp](https://one2n.io/sre-bootcamp) curriculum.

# How to run


- Clone the repo.
- Install `make`.
- Run `make docker-run-backend`
  - This would first run the postgres container. 
  - Once that's up and healthy, it will run the migrations.
  - Once the migrations have completed, the backend container (running the main app) would come up.

