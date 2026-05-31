/**
 * E2E tests for the Playground UI component.
 * Requires both the Astro dev server (port 4321) and Go backend (port 8080).
 */

describe("Playground Component", () => {
  beforeEach(() => {
    cy.visit("/playground/");
  });

  it("renders the code editor", () => {
    cy.get(".playground").should("exist");
    cy.get(".playground-editor .cm-editor").should("exist");
    cy.get(".playground-run").should("exist").and("contain.text", "Run");
  });

  it("has pre-filled starter code", () => {
    cy.get(".cm-content").should("contain.text", "package main");
  });

  it("runs code and displays output", () => {
    cy.get(".playground-run").click();
    cy.get(".playground-status").should("contain.text", "Running...");

    // Wait for compile response (Go Playground can be slow)
    cy.get(".output-content", { timeout: 15000 }).should(
      "contain.text",
      "Hello, World!"
    );
    cy.get(".playground-status").should("contain.text", "Success");
    cy.get(".output-content").should("have.class", "success");
  });

  it("supports Ctrl+Enter to run", () => {
    cy.get(".cm-content").type("{ctrl+enter}");
    cy.get(".playground-status").should("contain.text", "Running...");
    cy.get(".output-content", { timeout: 15000 }).should("not.be.empty");
  });

  it("displays compile errors", () => {
    // Clear editor and type invalid code
    cy.get(".cm-content").click().type("{selectAll}");
    cy.get(".cm-content").type(
      "package main{enter}{enter}func main() {{{enter}\tinvalid code{enter}}"
    );
    cy.get(".playground-run").click();

    cy.get(".output-content", { timeout: 15000 }).should("have.class", "error");
    cy.get(".playground-status").should("contain.text", "Compilation error");
  });

  it("disables run button while loading", () => {
    cy.get(".playground-run").click();
    cy.get(".playground-run").should("be.disabled");
    // Eventually re-enables
    cy.get(".playground-run", { timeout: 15000 }).should("not.be.disabled");
  });
});
