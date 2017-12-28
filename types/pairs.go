package hmtypes

// pair types

// Choice is the type of choice of algorithm to use within a class method.
//
// Imagine how one would implement a class in an OOP language.
// Then imagine how one would implement method overloading for the class.
// The typical approach is name mangling followed by having a jump table.
//
// Now consider OOP classes and the ability to override methods, based on subclassing ability.
// The typical approach to this is to use a Vtable.
//
// Both overloading and overriding have a general notion: a jump table of sorts.
// How does one type such a table?
//
// By using  Choice.
//
// The first type is the key of either the vtable or the name mangled table.
// The second type is the value of the table.
//
// TODO: implement hm.Type
type Choice Pair

// Super is the inverse of Choice. It allows for supertyping functions.
//
// Supertyping is typically  implemented as a adding an entry to the vtable/mangled table.
// But there needs to be a separate accounting structure to keep account of the types.
//
// This is where Super comes in.
//
// TODO: implement hm.Type
type Super Pair
