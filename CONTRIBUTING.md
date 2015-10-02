# Contributing to Coduno

## Google App Engine

The Coduno API is designed for Google's App Engine environment. You'll see
[the data model](https://godoc.org/github.com/coduno/api/model) leverage
optimizations for Google Datastore, and lot's of code depending on the
App Engine runtime. Therefore a good understanding of this environment is
critical in order to contribute.

### General

 * [Building Scalable Web Applications with Google App Engine (talk)](https://youtu.be/Oh9_t5W6MTE) ([slides](https://docs.google.com/presentation/d/1-nuc9jOvfHTW-yEP6RrJw-SOFFLrBbtWk8kPU8mwdGo/embed))

### Datastore

 * [Datastore Concepts Overview](https://cloud.google.com/datastore/docs/concepts/overview)
 * [`google.golang.org/appengine/datastore`](https://godoc.org/google.golang.org/appengine/datastore)
 * [Mastering the datastore](https://cloud.google.com/appengine/articles/datastore/overview) (series)
 * [Balancing Strong and Eventual Consistency with Google Cloud Datastore](https://cloud.google.com/datastore/docs/articles/balancing-strong-and-eventual-consistency-with-google-cloud-datastore/)
 * [Under the Covers of the Google App Engine Datastore (talk)](https://youtu.be/tx5gdoNpcZM) ([slides](http://snarfed.org/datastore_talk.html))
 * [How I learned to love the Datastore (talk)](https://youtu.be/WAa1r4BSWAU)

## Git

### Commit messages

  1. First is a tagged headline, maximum 50 chars.
    * Tags of what you changed (all lowercase, separated by `", "`)
    * `": "`
    * Headline
  2. Second line is empty.
  3. Rest of the message must be an exhaustive description of what the changeset contains. Enumerations are welcome.

There's currently no signing policy in place, but feel free to sign your commits.

#### Tags

Tags allow for better querying the history and provide context. Please use them, but to a sensible degree. It's okay to introduce new tags; just be consistent.

 * `ctrl` if you changed a controller.
 * `dep` if you introduce a new dependency.
 * `doc` when touching only documentation.
 * `gae` for App Engine configuration like `app.yaml`, `index.yaml`, `Dockerfile` or something else that has to do with running on App Engine.
 * `gen` if the changes are all automatically generated.
 * `gcs` anything that talks or refers to Google Cloud Storage.
 * `sec` if your change has security implications.
 * `sub` when changing ingestion of submissions.
 * `test` when touching a tester implementation.
 * `ws` regarding WebSockets.

If your tag list does not fit in 20 chars, you are probably doing something wrong.
