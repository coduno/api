# Contributing to Coduno

## Commit messages

  1. First line should be a tagged headline, maximum 50 chars.
    * Tags of what you changed (all lowercase, separated by `" ,"`)
    * `": "`
    * Headline
  2. Second line is empty.
  3. Rest of the message must be a description of what you did. Enumerations are welcome.

There's currently no signing policy in place, but feel free to sign your commits.

### Tags

 * `gae` if you changed AppEngine configuration like `app.yaml`, or modified `Dockerfile` or something else that has to do with running on AppEngine
 * `sec` if your change has security implications in general.
 * `ctrl` if you changed a controller
 * `sess` if you change sessions, authentication or cookie handling.
 * `dep` if you introduce a new dependency

If your tag list does not fit in 20 chars, you are probably doing something wrong.
