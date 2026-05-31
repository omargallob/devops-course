// e2e support file — runs before each spec
// Add custom commands and global configuration here.

// Prevent Cypress from failing on uncaught exceptions from the app
Cypress.on("uncaught:exception", () => false);
