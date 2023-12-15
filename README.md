# pgseeder - PostgreSQL Database seeder 

Efficiently seed your postgres database from files using a hierarchical dependency structure. Supports parallel execution, file streaming, and query batching.

See a basic example [here](./examples/basic/README.md).

## Data file format

Each line in your data files should have the information required to create one database row. The data can be any format but I recommend json. The file then needs to be gzipped.

Example:


**users.json**

```json
[
    {"id": 1, "name": "John"},
    {"id": 2, "name": "Jane"}
]
```

```bash
jq -c '.[]' users.json | gzip > users.gz
```