// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package id describes ACT's dot-delimited IDs.
package id

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// SepTag is the identifier tag separator.
	// It is exported for testing and sanitisation purposes.
	SepTag = '.'
)

var (
	// ErrTagHasSep occurs when a tag passed to New contains the separator rune.
	ErrTagHasSep = errors.New("tag contains separator")

	// ErrTagEmpty occurs when a tag passed to New is empty.
	ErrTagEmpty = errors.New("tag empty")
)

// ID represents an ACT ID.
type ID struct {
	tags []string
}

// New tries to construct an ACT ID from tags.
// It fails if any of the tags is empty (unless there is only one such tag), or contains a separator.
func New(tags ...string) (ID, error) {
	// Normalise the empty tag.
	if len(tags) == 1 && tags[0] == "" {
		return ID{}, nil
	}

	vtags, err := validateTags(tags)
	if err != nil {
		return ID{nil}, fmt.Errorf("tag validation failed for %v: %w", tags, err)
	}

	return ID{tags: vtags}, nil
}

func validateTags(tags []string) ([]string, error) {
	vtags := make([]string, len(tags))

	for i, t := range tags {
		vt := strings.TrimSpace(strings.ToLower(t))
		if err := validateTag(vt); err != nil {
			return nil, fmt.Errorf("%w: tag %q", err, vt)
		}
		vtags[i] = vt
	}
	return vtags, nil
}

func validateTag(t string) error {
	// TODO(@MattWindsor91): case folding and trimming
	if t == "" {
		return ErrTagEmpty
	}
	if strings.ContainsRune(t, SepTag) {
		return ErrTagHasSep
	}
	return nil
}

// TryFromString tries to convert a string to an ACT ID.
// It returns any validation error arising.
func TryFromString(s string) (ID, error) {
	return New(strings.Split(s, string(SepTag))...)
}

//FromString converts a string to an ACT ID.
// It returns the empty ID if there is an error.
func FromString(s string) ID {
	id, err := TryFromString(s)
	if err != nil {
		return ID{}
	}
	return id
}

// IsEmpty gets whether this ID is empty.
func (i ID) IsEmpty() bool {
	return len(i.tags) == 0
}

// Tags extracts the tags comprising an ID as a slice.
func (i ID) Tags() []string {
	return i.tags
}

// String converts an ACT ID to a string.
func (i ID) String() string {
	return strings.Join(i.tags, string(SepTag))
}

// Join appends r to this ID, creating a new ID.
func (i ID) Join(r ID) ID {
	if i.IsEmpty() {
		return r
	}
	if r.IsEmpty() {
		return i
	}
	return ID{append(i.tags, r.tags...)}
}