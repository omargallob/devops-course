/**
 * E2E tests for the exercise API endpoints.
 * These test the Go backend directly via HTTP.
 */

const API_URL = Cypress.env("API_URL") || "http://localhost:8080";

describe("GET /api/exercises/{exerciseId}", () => {
  it("returns exercise definition for a valid ID", () => {
    cy.request(`${API_URL}/api/exercises/m01-hello-world`).then((resp) => {
      expect(resp.status).to.eq(200);
      expect(resp.body).to.have.property("id", "m01-hello-world");
      expect(resp.body).to.have.property("title", "Hello World");
      expect(resp.body).to.have.property("instructions");
      expect(resp.body).to.have.property("starterCode");
      expect(resp.body).to.have.property("validationMode", "exact");
    });
  });

  it("does not expose expected output", () => {
    cy.request(`${API_URL}/api/exercises/m01-hello-world`).then((resp) => {
      expect(resp.body).to.not.have.property("expectedOutput");
    });
  });

  it("returns 404 for unknown exercise", () => {
    cy.request({
      url: `${API_URL}/api/exercises/nonexistent`,
      failOnStatusCode: false,
    }).then((resp) => {
      expect(resp.status).to.eq(404);
      expect(resp.body).to.have.property("error");
    });
  });

  it("returns starter code that is valid Go", () => {
    cy.request(`${API_URL}/api/exercises/m01-hello-world`).then((resp) => {
      expect(resp.body.starterCode).to.include("package main");
    });
  });
});

describe("POST /api/validate", () => {
  it("passes with correct solution", () => {
    const code = `package main

import "fmt"

func main() {
\tfmt.Println("Hello, World!")
}
`;
    cy.request({
      method: "POST",
      url: `${API_URL}/api/validate`,
      body: { exerciseId: "m01-hello-world", code },
      headers: { "Content-Type": "application/json" },
    }).then((resp) => {
      expect(resp.status).to.eq(200);
      expect(resp.body).to.have.property("passed", true);
      expect(resp.body).to.have.property("exerciseId", "m01-hello-world");
    });
  });

  it("fails with incorrect output", () => {
    const code = `package main

import "fmt"

func main() {
\tfmt.Println("Wrong output")
}
`;
    cy.request({
      method: "POST",
      url: `${API_URL}/api/validate`,
      body: { exerciseId: "m01-hello-world", code },
      headers: { "Content-Type": "application/json" },
    }).then((resp) => {
      expect(resp.status).to.eq(200);
      expect(resp.body).to.have.property("passed", false);
      expect(resp.body).to.have.property("expectedOutput");
      expect(resp.body).to.have.property("actualOutput");
      expect(resp.body).to.have.property("diff");
    });
  });

  it("returns compile error for invalid code", () => {
    const code = `package main

func main() {
\tthis is not valid go
}
`;
    cy.request({
      method: "POST",
      url: `${API_URL}/api/validate`,
      body: { exerciseId: "m01-hello-world", code },
      headers: { "Content-Type": "application/json" },
    }).then((resp) => {
      expect(resp.status).to.eq(200);
      expect(resp.body).to.have.property("passed", false);
      expect(resp.body).to.have.property("compileError");
    });
  });

  it("returns 400 for empty code", () => {
    cy.request({
      method: "POST",
      url: `${API_URL}/api/validate`,
      body: { exerciseId: "m01-hello-world", code: "   " },
      headers: { "Content-Type": "application/json" },
      failOnStatusCode: false,
    }).then((resp) => {
      expect(resp.status).to.eq(400);
      expect(resp.body).to.have.property("error");
    });
  });

  it("returns 400 for missing exercise ID", () => {
    cy.request({
      method: "POST",
      url: `${API_URL}/api/validate`,
      body: { code: "package main" },
      headers: { "Content-Type": "application/json" },
      failOnStatusCode: false,
    }).then((resp) => {
      expect(resp.status).to.eq(400);
    });
  });

  it("returns 404 for unknown exercise", () => {
    cy.request({
      method: "POST",
      url: `${API_URL}/api/validate`,
      body: { exerciseId: "nonexistent", code: "package main" },
      headers: { "Content-Type": "application/json" },
      failOnStatusCode: false,
    }).then((resp) => {
      expect(resp.status).to.eq(404);
    });
  });

  it("passes the fizzbuzz exercise with correct solution", () => {
    const code = `package main

import "fmt"

func main() {
\tfor i := 1; i <= 15; i++ {
\t\tswitch {
\t\tcase i%15 == 0:
\t\t\tfmt.Println("FizzBuzz")
\t\tcase i%3 == 0:
\t\t\tfmt.Println("Fizz")
\t\tcase i%5 == 0:
\t\t\tfmt.Println("Buzz")
\t\tdefault:
\t\t\tfmt.Println(i)
\t\t}
\t}
}
`;
    cy.request({
      method: "POST",
      url: `${API_URL}/api/validate`,
      body: { exerciseId: "m01-fizzbuzz", code },
      headers: { "Content-Type": "application/json" },
    }).then((resp) => {
      expect(resp.status).to.eq(200);
      expect(resp.body).to.have.property("passed", true);
    });
  });
});
