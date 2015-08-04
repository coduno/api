// Package model groups a lot of types that model the data that
// drives Coduno.
//
// It is designed for Google Datastore, but switching to other
// BigTable implementations is not inconceivable (though it might
// require some more work).
//
// Not all exported types directly translate to entities that are
// persisted, but most of them do.
package model
