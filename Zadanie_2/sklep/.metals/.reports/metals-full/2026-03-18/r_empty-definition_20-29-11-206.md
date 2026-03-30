error id: file://<WORKSPACE>/app/controllers/ProductController.scala:
file://<WORKSPACE>/app/controllers/ProductController.scala
empty definition using pc, found symbol in pc: 
empty definition using semanticdb
empty definition using fallback
non-local guesses:
	 -javax/inject/id.
	 -javax/inject/id#
	 -javax/inject/id().
	 -play/api/mvc/id.
	 -play/api/mvc/id#
	 -play/api/mvc/id().
	 -play/api/libs/json/id.
	 -play/api/libs/json/id#
	 -play/api/libs/json/id().
	 -id.
	 -id#
	 -id().
	 -scala/Predef.id.
	 -scala/Predef.id#
	 -scala/Predef.id().
offset: 1866
uri: file://<WORKSPACE>/app/controllers/ProductController.scala
text:
```scala
package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._
import models.Product
import scala.collection.mutable.ListBuffer

@Singleton
class ProductController @Inject()(val controllerComponents: ControllerComponents) extends BaseController {

  // Nasza "baza danych" w pamięci w formie listy
  private val productList = new ListBuffer[Product]()
  
  // Dodajemy dwa produkty na start, żeby było co wyświetlać
  productList += Product(1, "Laptop", 3500.00)
  productList += Product(2, "Myszka", 120.50)

  // 1. READ ALL - Pobierz wszystkie produkty (GET)
  def getAll: Action[AnyContent] = Action { implicit request =>
    Ok(Json.toJson(productList))
  }

  // 2. READ ONE - Pobierz produkt po ID (GET)
  def getById(id: Long): Action[AnyContent] = Action { implicit request =>
    productList.find(_.id == id) match {
      case Some(product) => Ok(Json.toJson(product))
      case None => NotFound(Json.obj("message" -> s"Nie znaleziono produktu o id $id"))
    }
  }

  // 3. CREATE - Dodaj nowy produkt (POST)
  def add(): Action[JsValue] = Action(parse.json) { implicit request =>
    request.body.validate[Product].fold(
      errors => BadRequest(Json.obj("message" -> "Błędny format JSON")),
      product => {
        productList += product
        Created(Json.toJson(product))
      }
    )
  }

  // 4. UPDATE - Zaktualizuj produkt (PUT)
  def update(id: Long): Action[JsValue] = Action(parse.json) { implicit request =>
    request.body.validate[Product].fold(
      errors => BadRequest(Json.obj("message" -> "Błędny format JSON")),
      newProduct => {
        val index = productList.indexWhere(_.id == id)
        if (index >= 0) {
          productList.update(index, newProduct)
          Ok(Json.toJson(newProduct))
        } else {
          NotFound(Json.obj("message" -> s"Nie znaleziono produktu o id $id@@"))
        }
      }
    )
  }

  // 5. DELETE - Usuń produkt (DELETE)
  def delete(id: Long): Action[AnyContent] = Action { implicit request =>
    val index = productList.indexWhere(_.id == id)
    if (index >= 0) {
      productList.remove(index)
      Ok(Json.obj("message" -> s"Usunięto produkt o id $id"))
    } else {
      NotFound(Json.obj("message" -> s"Nie znaleziono produktu o id $id"))
    }
  }
}
```


#### Short summary: 

empty definition using pc, found symbol in pc: 