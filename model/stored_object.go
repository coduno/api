package model

// StoredObject refers to an object stored in
// Google Cloud Storage.
//
// See https://cloud.google.com/storage/docs/concepts-techniques#concepts
type StoredObject struct {
	// The bucket the object resides in.
	Bucket string

	// Name of the object inside the bucket.
	Name string
}
