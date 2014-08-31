Build the scaffold of an API from a JSON schema with best practices baked in.


# Requirements

   - [Go](http://golang.org/dl/)


To build the article:

   - Pandoc


# Usage

Build the application:

```
$ make
```

Generate a new scaffold:

```
$ ./api-from-schema [path_to_schema.json] [path_to_new_dir]
```

This will use the JSON schema found in `path_to_schema.json` and create a scaffold in the directory `path_to_new_dir`.


