This is a basic crud app — a part of the [one-2n-sre-bootcamp](https://one2n.io/sre-bootcamp) curriculum.

## Steps to run

- Clone
- run `make run`

# How to run

```zsh
docker build -t one2n:0.1.0 .
docker run --rm -it -e DB_PATH=students.db -p 8000:8000 one2n:0.1.0
```
