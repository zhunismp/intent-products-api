package product

// For now, the domain model is used directly as the database model,
// since the current logic is simple and both representations are identical.
// As the application grows and domain logic becomes more complex,
// the database model may diverge from the domain model to better support
// persistence and business concerns separately.