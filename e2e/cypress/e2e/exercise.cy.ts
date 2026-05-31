/**
 * E2E tests for the Exercise UI component.
 * Requires both the Astro dev server (port 4321) and Go backend (port 8080).
 */

describe("Exercise Component", () => {
  beforeEach(() => {
    cy.visit("/exercises/m01-hello-world/");
  });

  it("renders exercise with title and instructions", () => {
    cy.get(".exercise").should("exist");
    cy.get(".exercise-title").should("contain.text", "Hello World");
    cy.get(".exercise-instructions").should("exist");
  });

  it("has starter code pre-filled in editor", () => {
    cy.get(".cm-content").should("contain.text", "package main");
    cy.get(".cm-content").should("contain.text", "// Your code here");
  });

  it("can run code without submitting", () => {
    // Type a solution
    cy.get(".cm-content").click().type("{selectAll}");
    cy.get(".cm-content").type(
      'package main{enter}{enter}import "fmt"{enter}{enter}func main() {{{enter}\tfmt.Println("Hello, World!"){enter}}'
    );

    cy.get(".playground-run").click();
    cy.get(".output-content", { timeout: 15000 }).should(
      "contain.text",
      "Hello, World!"
    );
  });

  it("submits and shows pass result for correct solution", () => {
    // Type correct solution
    cy.get(".cm-content").click().type("{selectAll}");
    cy.get(".cm-content").type(
      'package main{enter}{enter}import "fmt"{enter}{enter}func main() {{{enter}\tfmt.Println("Hello, World!"){enter}}'
    );

    cy.get(".exercise-submit-btn").click();
    cy.get(".exercise-submit-btn").should("contain.text", "Validating...");

    cy.get(".exercise-result", { timeout: 15000 }).should("be.visible");
    cy.get(".result-badge").should("contain.text", "Passed!");
    cy.get(".result-badge").should("have.class", "pass");
  });

  it("submits and shows fail result for incorrect solution", () => {
    // Type wrong solution
    cy.get(".cm-content").click().type("{selectAll}");
    cy.get(".cm-content").type(
      'package main{enter}{enter}import "fmt"{enter}{enter}func main() {{{enter}\tfmt.Println("Wrong!"){enter}}'
    );

    cy.get(".exercise-submit-btn").click();

    cy.get(".exercise-result", { timeout: 15000 }).should("be.visible");
    cy.get(".result-badge").should("contain.text", "Not quite...");
    cy.get(".result-badge").should("have.class", "fail");
    cy.get(".result-details").should("not.be.empty");
  });

  it("resets editor to starter code", () => {
    // Modify the code
    cy.get(".cm-content").click().type("{selectAll}");
    cy.get(".cm-content").type("modified code");

    // Reset
    cy.get(".exercise-reset-btn").click();
    cy.get(".cm-content").should("contain.text", "// Your code here");
  });

  it("toggles hint visibility", () => {
    cy.get(".exercise-hint").should("not.be.visible");
    cy.get(".exercise-hint-btn").click();
    cy.get(".exercise-hint").should("be.visible");
    cy.get(".exercise-hint-btn").should("contain.text", "Hide Hint");

    cy.get(".exercise-hint-btn").click();
    cy.get(".exercise-hint").should("not.be.visible");
  });

  it("clears result on reset", () => {
    // Submit first
    cy.get(".cm-content").click().type("{selectAll}");
    cy.get(".cm-content").type(
      'package main{enter}{enter}import "fmt"{enter}{enter}func main() {{{enter}\tfmt.Println("Hello, World!"){enter}}'
    );
    cy.get(".exercise-submit-btn").click();
    cy.get(".exercise-result", { timeout: 15000 }).should("be.visible");

    // Reset should hide result
    cy.get(".exercise-reset-btn").click();
    cy.get(".exercise-result").should("not.be.visible");
  });
});
