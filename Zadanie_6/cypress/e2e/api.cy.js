const API = 'http://localhost:8080'

// request bez rzucania bledu na zly status
const neg = (method, url, body = undefined) =>
  cy.request({
    method, url, body, failOnStatusCode: false,
    ...(body ? { headers: { 'Content-Type': 'application/json' } } : {})
  })

describe('Testy API – Products', () => {

  let createdProductId



  it('POST /products – tworzy produkt', () => {
    cy.request('POST', `${API}/products`, {
      name: 'API Produkt', description: 'Opis API', price: 55.55,
    }).then((res) => {
      expect(res.status).to.eq(201)
      expect(res.body.name).to.eq('API Produkt')
      expect(res.body.description).to.eq('Opis API')
      expect(res.body.price).to.eq(55.55)
      expect(res.body.ID).to.be.greaterThan(0)
      createdProductId = res.body.ID
    })
  })

  it('GET /products – zwraca listę produktów', () => {
    cy.request('GET', `${API}/products`).then((res) => {
      expect(res.status).to.eq(200)
      expect(res.body).to.be.an('array')
      expect(res.body.length).to.be.greaterThan(0)
    })
  })

  it('GET /products/:id – zwraca konkretny produkt', () => {
    cy.request('GET', `${API}/products/${createdProductId}`).then((res) => {
      expect(res.status).to.eq(200)
      expect(res.body.name).to.eq('API Produkt')
      expect(res.body.price).to.eq(55.55)
    })
  })

  it('PUT /products/:id – aktualizuje produkt', () => {
    cy.request('PUT', `${API}/products/${createdProductId}`, {
      name: 'Zmieniony Produkt', price: 99.99,
    }).then((res) => {
      expect(res.status).to.eq(200)
      expect(res.body.name).to.eq('Zmieniony Produkt')
      expect(res.body.price).to.eq(99.99)
    })
  })

  it('DELETE /products/:id – usuwa produkt', () => {
    cy.request('DELETE', `${API}/products/${createdProductId}`).then((res) => {
      expect(res.status).to.eq(204)
    })
  })



  it('POST /products – niepoprawne dane (pusty body)', () => {
    neg('POST', `${API}/products`, 'invalid').then((res) => {
      expect(res.status).to.be.oneOf([400, 500])
    })
  })

  it('GET /products/:id – produkt nie istnieje', () => {
    neg('GET', `${API}/products/999999`).then((res) => {
      expect(res.status).to.eq(404)
      expect(res.body).to.have.property('error')
    })
  })

  it('PUT /products/:id – produkt nie istnieje', () => {
    neg('PUT', `${API}/products/999999`, { name: 'X', price: 1 }).then((res) => {
      expect(res.status).to.eq(404)
    })
  })

  it('DELETE /products/:id – produkt nie istnieje', () => {
    neg('DELETE', `${API}/products/${createdProductId}`).then((res) => {
      expect(res.status).to.eq(404)
      expect(res.body).to.have.property('error')
    })
  })
})

describe('Testy API – Carts', () => {

  let cartId
  let itemId

  it('POST /carts – tworzy koszyk', () => {
    cy.request('POST', `${API}/carts`).then((res) => {
      expect(res.status).to.eq(201)
      expect(res.body.status).to.eq('aktywny')
      expect(res.body.ID).to.be.greaterThan(0)
      cartId = res.body.ID
    })
  })

  it('GET /carts/:id – pobiera koszyk', () => {
    cy.request('GET', `${API}/carts/${cartId}`).then((res) => {
      expect(res.status).to.eq(200)
      expect(res.body.status).to.eq('aktywny')
    })
  })

  it('POST /carts/:id/items – dodaje element do koszyka', () => {
    cy.request('POST', `${API}/products`, {
      name: 'Produkt do koszyka', description: 'Test', price: 25.00,
    }).then((productRes) => {
      cy.request('POST', `${API}/carts/${cartId}/items`, {
        product_id: productRes.body.ID, quantity: 2,
      }).then((res) => {
        expect(res.status).to.eq(201)
        expect(res.body.quantity).to.eq(2)
        expect(res.body.product_id).to.eq(productRes.body.ID)
        itemId = res.body.ID
      })
    })
  })

  it('GET /carts/:id – koszyk z elementami', () => {
    cy.request('GET', `${API}/carts/${cartId}`).then((res) => {
      expect(res.status).to.eq(200)
      expect(res.body.items).to.be.an('array')
      expect(res.body.items.length).to.be.greaterThan(0)
    })
  })

  it('DELETE /carts/:id/items/:itemId – usuwa element z koszyka', () => {
    cy.request('DELETE', `${API}/carts/${cartId}/items/${itemId}`).then((res) => {
      expect(res.status).to.eq(204)
    })
  })

  it('DELETE /carts/:id – usuwa koszyk', () => {
    cy.request('DELETE', `${API}/carts/${cartId}`).then((res) => {
      expect(res.status).to.eq(204)
    })
  })



  it('GET /carts/:id – koszyk nie istnieje', () => {
    neg('GET', `${API}/carts/999999`).then((res) => {
      expect(res.status).to.eq(404)
      expect(res.body).to.have.property('error')
    })
  })

  it('POST /carts/:id/items – koszyk nie istnieje', () => {
    neg('POST', `${API}/carts/999999/items`, { product_id: 1, quantity: 1 }).then((res) => {
      expect(res.status).to.eq(404)
    })
  })

  it('POST /carts/:id/items – produkt nie istnieje', () => {
    cy.seedCart().then((cartRes) => {
      neg('POST', `${API}/carts/${cartRes.body.ID}/items`, { product_id: 999999, quantity: 1 }).then((res) => {
        expect(res.status).to.eq(404)
      })
    })
  })

  it('DELETE /carts/:id/items/:itemId – element nie istnieje', () => {
    cy.seedCart().then((cartRes) => {
      neg('DELETE', `${API}/carts/${cartRes.body.ID}/items/999999`).then((res) => {
        expect(res.status).to.eq(404)
      })
    })
  })

  it('DELETE /carts/:id – koszyk nie istnieje', () => {
    neg('DELETE', `${API}/carts/999999`).then((res) => {
      expect(res.status).to.eq(404)
    })
  })
})

describe('Testy API – Payments', () => {

  it('POST /payments – tworzy płatność', () => {
    cy.request('POST', `${API}/payments`, {
      amount: 100.50, status: 'PAID',
    }).then((res) => {
      expect(res.status).to.eq(201)
      expect(res.body.amount).to.eq(100.50)
      expect(res.body.status).to.eq('PAID')
      expect(res.body.ID).to.be.greaterThan(0)
    })
  })

  it('GET /payments – zwraca listę płatności', () => {
    cy.request('GET', `${API}/payments`).then((res) => {
      expect(res.status).to.eq(200)
      expect(res.body).to.be.an('array')
      expect(res.body.length).to.be.greaterThan(0)
    })
  })



  it('POST /payments – niepoprawne dane', () => {
    neg('POST', `${API}/payments`, 'notjson').then((res) => {
      expect(res.status).to.be.oneOf([400, 500])
    })
  })
})
