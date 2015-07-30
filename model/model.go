// Package model groups a lot of types that model the data that
// drives Coduno.
//
// It is designed to be agnostic about the underlying persistance
// layer, but with optimizations for Google Datastore in mind.
// So switching to other BigTable implementations is not inconceivable
// (though it might require some more work).
// Implementations (like App Engine and Compute Engine APIs for
// Datastore) are to be switched by using a Go build tag.
//
// Not all exported types directly translate to entities that are
// persisted, but most of them do.
package model
