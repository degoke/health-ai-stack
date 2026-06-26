// Package registry provides the shared FHIR definition catalog for the monorepo.
//
// It seeds bundled R4 StructureDefinition and SearchParameter resources, persists
// catalog and install overlay state through pkg/store contracts, and compiles an
// immutable in-memory snapshot for fast runtime lookups.
package registry
