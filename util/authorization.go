package util

import "google.golang.org/appengine/datastore"

// HasParent returns true if key or any of its parents equals
// parent (recursively), otherwise false.
func HasParent(parent, key *datastore.Key) bool {
	for key != nil {
		if key.Equal(parent) {
			return true
		}
		key = key.Parent()
	}
	return false
}
