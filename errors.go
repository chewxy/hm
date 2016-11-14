package hm

const (
	typeMismatch         = "Type Mismatch: %v != %v"
	nyi                  = "%v Not Yet Implemented for %v of %T"
	nu                   = "Types %v and %v are not unifiable"
	undefinedTV          = "Undefined TypeVariable cannot be unified"
	recursiveUnification = "Type %v will cause a recursive unification with %v"
	tvinstance           = "Different instances of TypeVariable %v != %v. Name are the same: %q "
)
