const API = 'http://localhost:8080'

Cypress.Commands.add('seedProduct', (name, description, price) => {
  return cy.request('POST', `${API}/products`, { name, description, price })
})

Cypress.Commands.add('seedCart', () => {
  return cy.request('POST', `${API}/carts`)
})
