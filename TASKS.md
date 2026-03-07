# TASKS

a [meads](https://github.com/jpillora/meads) (`md`) managed task log

* created: 2026-03-05T11:08:31Z
* next-id: 2

## 1. Fix camel2dash mangling short acronyms like IDs

* status: completed
* type: bug
* created: 2026-03-05T11:08:31Z

camel2dash("IDs") produces "i-dss" instead of "ids" - the algorithm mishandles field names where a short uppercase acronym is followed by a lowercase letter
