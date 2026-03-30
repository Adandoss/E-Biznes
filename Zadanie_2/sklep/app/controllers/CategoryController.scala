package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._
import scala.collection.mutable.ListBuffer

case class Category(id: Long, name: String)

object Category {
  implicit val format: OFormat[Category] = Json.format[Category]
}

@Singleton
class CategoryController @Inject()(val controllerComponents: ControllerComponents) extends BaseController {

  private val categoryList = ListBuffer(
    Category(1, "Elektronika"),
    Category(2, "Akcesoria")
  )

  def getAll = Action {
    Ok(Json.toJson(categoryList))
  }

  def getById(id: Long) = Action {
    categoryList.find(_.id == id) match {
      case Some(category) => Ok(Json.toJson(category))
      case None => NotFound(Json.obj("message" -> "Nie znaleziono kategorii"))
    }
  }

  def add() = Action(parse.json) { request =>
    request.body.validate[Category].fold(
      _ => BadRequest(Json.obj("message" -> "Błędny format JSON")),
      category => {
        categoryList += category
        Created(Json.toJson(category))
      }
    )
  }

  def update(id: Long) = Action(parse.json) { request =>
    request.body.validate[Category].fold(
      _ => BadRequest(Json.obj("message" -> "Błędny format JSON")),
      newCategory => {
        val index = categoryList.indexWhere(_.id == id)
        if (index >= 0) {
          categoryList.update(index, newCategory)
          Ok(Json.toJson(newCategory))
        } else {
          NotFound(Json.obj("message" -> "Nie znaleziono kategorii"))
        }
      }
    )
  }

  def delete(id: Long) = Action {
    val index = categoryList.indexWhere(_.id == id)
    if (index >= 0) {
      categoryList.remove(index)
      Ok(Json.obj("message" -> "Usunięto kategorię"))
    } else {
      NotFound(Json.obj("message" -> "Nie znaleziono kategorii"))
    }
  }
}