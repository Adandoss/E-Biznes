
describe('Testy funkcjonalne sklepu', () => {

  before(() => {
    cy.seedProduct('Laptop Testowy', 'Laptop do testów', 3999.99)
    cy.seedProduct('Myszka Testowa', 'Myszka do testów', 49.99)
  })

  beforeEach(() => {
    cy.visit('/')
  })


  it('1. Strona główna się ładuje', () => {
    cy.url().should('include', '/')
    cy.get('body').should('be.visible')
    cy.get('body').should('not.be.empty')
  })


  it('2. Tytuł "Sklep" jest widoczny', () => {
    cy.get('h1').should('exist')
    cy.get('h1').should('contain.text', 'Sklep')
    cy.get('h1').should('be.visible')
  })


  it('3. Nawigacja zawiera 3 linki', () => {
    cy.get('nav').should('exist')
    cy.get('nav').should('be.visible')
    cy.get('nav a').should('have.length', 3)
  })


  it('4. Link "Produkty" prowadzi do /', () => {
    cy.get('nav a').contains('Produkty').should('exist')
    cy.get('nav a').contains('Produkty').should('have.attr', 'href', '/')
    cy.get('nav a').contains('Produkty').should('be.visible')
  })


  it('5. Link "Koszyk" prowadzi do /cart', () => {
    cy.get('nav a').contains('Koszyk').should('exist')
    cy.get('nav a').contains('Koszyk').should('have.attr', 'href', '/cart')
    cy.get('nav a').contains('Koszyk').should('be.visible')
  })


  it('6. Link "Płatności" prowadzi do /payments', () => {
    cy.get('nav a').contains('Płatności').should('exist')
    cy.get('nav a').contains('Płatności').should('have.attr', 'href', '/payments')
    cy.get('nav a').contains('Płatności').should('be.visible')
  })


  it('7. Nagłówek strony Produkty jest widoczny', () => {
    cy.get('h2').should('be.visible')
    cy.get('h2').should('contain.text', 'Produkty')
  })


  it('8. Lista produktów wyświetla się', () => {
    cy.get('ul').should('exist')
    cy.get('ul li').should('have.length.greaterThan', 0)
    cy.get('ul li').first().should('be.visible')
  })


  it('9. Produkty zawierają cenę w PLN', () => {
    cy.get('ul li').first().should('contain.text', 'PLN')
    cy.get('ul li').first().invoke('text').should('match', /\d/)
    cy.get('ul li').first().should('contain.text', 'Dodaj do koszyka')
  })


  it('10. Każdy produkt ma przycisk "Dodaj do koszyka"', () => {
    cy.get('ul li').first().find('button').should('exist')
    cy.get('ul li').first().find('button').should('contain.text', 'Dodaj do koszyka')
    cy.get('ul li button').should('have.length.greaterThan', 0)
  })


  it('11. Formularz dodawania produktu istnieje', () => {
    cy.get('form').should('exist')
    cy.get('h3').should('contain.text', 'Dodaj nowy produkt')
    cy.get('form').should('be.visible')
  })


  it('12. Formularz ma pola: nazwa, opis, cena', () => {
    cy.get('input[placeholder="Nazwa produktu"]').should('exist')
    cy.get('input[placeholder="Opis"]').should('exist')
    cy.get('input[placeholder="Cena (PLN)"]').should('exist')
    cy.get('form button[type="submit"]').should('contain.text', 'Dodaj produkt')
  })


  it('13. Nawigacja do strony Koszyk', () => {
    cy.get('nav a').contains('Koszyk').click()
    cy.url().should('include', '/cart')
    cy.get('h2').should('contain.text', 'Koszyk')
  })


  it('14. Pusty koszyk wyświetla komunikat', () => {
    cy.get('nav a').contains('Koszyk').click()
    cy.get('main').should('contain.text', 'Koszyk jest pusty')
    cy.get('main').should('be.visible')
  })


  it('15. Koszyk wyświetla sumę', () => {
    cy.get('nav a').contains('Koszyk').click()
    cy.get('main').should('contain.text', 'Suma')
    cy.get('main').should('contain.text', 'PLN')
  })


  it('16. Nawigacja do strony Płatności', () => {
    cy.get('nav a').contains('Płatności').click()
    cy.url().should('include', '/payments')
    cy.get('h2').should('contain.text', 'Płatności')
  })


  it('17. Strona płatności wyświetla kwotę do zapłaty', () => {
    cy.get('nav a').contains('Płatności').click()
    cy.get('main').should('contain.text', 'Do zapłaty')
    cy.get('main').should('contain.text', 'PLN')
  })


  it('18. Przycisk "Kup" jest wyłączony gdy koszyk pusty', () => {
    cy.get('nav a').contains('Płatności').click()
    cy.get('button').contains('Kup').should('exist')
    cy.get('button').contains('Kup').should('be.disabled')
  })


  it('19. Sekcja historii płatności istnieje', () => {
    cy.get('nav a').contains('Płatności').click()
    cy.get('h3').should('exist')
    cy.get('h3').should('contain.text', 'Historia płatności')
  })


  it('20. Dodanie produktu przez formularz', () => {
    cy.get('input[placeholder="Nazwa produktu"]').type('Cypress Produkt')
    cy.get('input[placeholder="Opis"]').type('Opis testowy')
    cy.get('input[placeholder="Cena (PLN)"]').type('123.45')
    cy.get('form button[type="submit"]').click()
    cy.get('ul li').should('contain.text', 'Cypress Produkt')
    cy.get('ul li').should('contain.text', '123.45')
    cy.get('ul li').should('contain.text', 'PLN')
  })
})
