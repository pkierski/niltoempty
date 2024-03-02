# niltoempty
Recursively initialize all nil maps and slices in a given object, so they json.Marshal() as empty object {} or array [] instead of null.

This is more complete solution based on the idea from [nilslice](https://github.com/golang-cz/nilslice). It works not only for nil slices but also for nil maps. niltoempty.Initialize traverses any addressable entity recusively.

