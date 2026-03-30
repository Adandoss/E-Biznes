package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._
import scala.collection.mutable.ListBuffer

case class CartItem(id: Long, productId: Long, quantity: Int)

object CartItem {
  implicit val format: OFormat[CartItem] = Json.format[CartItem]
}

@Singleton
class CartController @Inject()(val controllerComponents: ControllerComponents) extends BaseController {

  private val cartList = ListBuffer(
    CartItem(1, 1, 1)
  )

  def getAll = Action {
    Ok(Json.toJson(cartList))
  }

  def getById(id: Long) = Action {
    cartList.find(_.id == id) match {
      case Some(item) => Ok(Json.toJson(item))
      case None => NotFound(Json.obj("message" -> "Nie znaleziono elementu w koszyku"))
    }
  }

  def add() = Action(parse.json) { request =>
    request.body.validate[CartItem].fold(
      _ => BadRequest(Json.obj("message" -> "Błędny format JSON")),
      item => {
        cartList += item
        Created(Json.toJson(item))
      }
    )
  }

  def update(id: Long) = Action(parse.json) { request =>
    request.body.validate[CartItem].fold(
      _ => BadRequest(Json.obj("message" -> "Błędny format JSON")),
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

  def delete(id: Long) = Action {
    val index = cartList.indexWhere(_.id == id)
    if (index >= 0) {
      cartList.remove(index)
      Ok(Json.obj("message" -> "Usunięto element z koszyka"))
    } else {
      NotFound(Json.obj("message" -> "Nie znaleziono elementu w koszyku"))
    }
  }
}