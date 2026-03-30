package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._
import scala.collection.mutable.ListBuffer

case class Product(id: Long, name: String, price: BigDecimal)

object Product {
  implicit val format: OFormat[Product] = Json.format[Product]
}

@Singleton
class ProductController @Inject()(val controllerComponents: ControllerComponents) extends BaseController {

  private val productList = ListBuffer(
    Product(1, "Laptop", 3500.00),
    Product(2, "Myszka", 120.50)
  )

  def getAll = Action {
    Ok(Json.toJson(productList))
  }

  def getById(id: Long) = Action {
    productList.find(_.id == id) match {
      case Some(product) => Ok(Json.toJson(product))
      case None => NotFound(Json.obj("message" -> s"Nie znaleziono produktu o id $id"))
    }
  }

  def add() = Action(parse.json) { request =>
    request.body.validate[Product].fold(
      _ => BadRequest(Json.obj("message" -> "Błędny format JSON")),
      product => {
        productList += product
        Created(Json.toJson(product))
      }
    )
  }

  def update(id: Long) = Action(parse.json) { request =>
    request.body.validate[Product].fold(
      _ => BadRequest(Json.obj("message" -> "Błędny format JSON")),
      newProduct => {
        val index = productList.indexWhere(_.id == id)
        if (index >= 0) {
          productList.update(index, newProduct)
          Ok(Json.toJson(newProduct))
        } else {
          NotFound(Json.obj("message" -> s"Nie znaleziono produktu o id $id"))
        }
      }
    )
  }

  def delete(id: Long) = Action {
    val index = productList.indexWhere(_.id == id)
    if (index >= 0) {
      productList.remove(index)
      Ok(Json.obj("message" -> s"Usunięto produkt o id $id"))
    } else {
      NotFound(Json.obj("message" -> s"Nie znaleziono produktu o id $id"))
    }
  }
}