package model

// StoredObject refers to an object stored in
// Google Cloud Storage.
//
// The project the object belongs to is defined
// by the context the application runs in.
//
// See https://cloud.google.com/storage/docs/overview#blocks
type StoredObject struct {
	// The bucket the object resides in.
	Bucket,

	// Name of the object inside the bucket.
	Name string
}
