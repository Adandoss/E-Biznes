error id: file://<WORKSPACE>/app/controllers/CartController.scala:Action.
file://<WORKSPACE>/app/controllers/CartController.scala
empty definition using pc, found symbol in pc: 
empty definition using semanticdb
empty definition using fallback
non-local guesses:
	 -javax/inject/Action.
	 -javax/inject/Action#
	 -javax/inject/Action().
	 -play/api/mvc/Action.
	 -play/api/mvc/Action#
	 -play/api/mvc/Action().
	 -play/api/libs/json/Action.
	 -play/api/libs/json/Action#
	 -play/api/libs/json/Action().
	 -Action.
	 -Action#
	 -Action().
	 -scala/Predef.Action.
	 -scala/Predef.Action#
	 -scala/Predef.Action().
offset: 1113
uri: file://<WORKSPACE>/app/controllers/CartController.scala
text:
```scala
package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._
import models.CartItem
import scala.collection.mutable.ListBuffer

@Singleton
class CartController @Inject()(val controllerComponents: ControllerComponents) extends BaseController {

  private val cartList = new ListBuffer[CartItem]()
  // Przykładowy element: w koszyku mamy 1 sztukę produktu o ID 1
  cartList += CartItem(1, 1, 1)

  def getAll: Action[AnyContent] = Action { implicit request =>
    Ok(Json.toJson(cartList))
  }

  def getById(id: Long): Action[AnyContent] = Action { implicit request =>
    cartList.find(_.id == id) match {
      case Some(item) => Ok(Json.toJson(item))
      case None => NotFound(Json.obj("message" -> "Nie znaleziono elementu w koszyku"))
    }
  }

  def add(): Action[JsValue] = Action(parse.json) { implicit request =>
    request.body.validate[CartItem].fold(
      errors => BadRequest(Json.obj("message" -> "Błędny format JSON")),
      item => {
        cartList += item
        Created(Json.toJson(item))
      }
    )
  }

  def update(id: Long): Action[JsValue] = Acti@@on(parse.json) { implicit request =>
    request.body.validate[CartItem].fold(
      errors => BadRequest(Json.obj("message" -> "Błędny format JSON")),
      newItem => {
        val index = cartList.indexWhere(_.id == id)
        if (index >= 0) {
          cartList.update(index, newItem)
          Ok(Json.toJson(newItem))
        } else {
          NotFound(Json.obj("message" -> "Nie znaleziono elementu w koszyku"))
        }
      }
    )
  }

  def delete(id: Long): Action[AnyContent] = Action { implicit request =>
    val index = cartList.indexWhere(_.id == id)
    if (index >= 0) {
      cartList.remove(index)
      Ok(Json.obj("message" -> "Usunięto element z koszyka"))
    } else {
      NotFound(Json.obj("message" -> "Nie znaleziono elementu w koszyku"))
    }
  }
}
```


#### Short summary: 

empty definition using pc, found symbol in pc: 