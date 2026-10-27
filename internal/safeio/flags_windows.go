package safeio

// Windows named pipes are not ordinary filesystem files; pre/post file type
// checks apply. Root-scoped calls continue to enforce path containment.
const readFlags = 0
