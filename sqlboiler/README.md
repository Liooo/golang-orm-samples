# sqlboiler

## Getting Started

### Prerequisites
- postgres server is up
- have the following executables in your $PATH
    - `dropdb`
    - `createdb`
    - `sqlboiler` (See [sqlboiler#download](https://github.com/volatiletech/sqlboiler#download))

```sh
cp .envrc.example .envrc # then edit to your needs
direnv allow # or, `source .envrc` if you don't have direnv installed

# drop db, create db, then (re) generates sqlboiler templates under models/ based on schema.sql
sh ./scripts/reset-db.sh

# execute
go run .
```
