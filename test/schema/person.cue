package foo

import "time"

// ## describe a person
#person: {
	// ### full name
	//
	//   - can include first middle last names
	//   - required
	name: string

	// ### date of birth
	//
	//   - in iso-8601 format
	//   - optional
	dob?: string
	dob?: time.Format("2006-01-02")
}
