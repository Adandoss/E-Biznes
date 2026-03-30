error id: file://<WORKSPACE>/app/controllers/CategoryController.scala:
file://<WORKSPACE>/app/controllers/CategoryController.scala
empty definition using pc, found symbol in pc: 
empty definition using semanticdb
empty definition using fallback
non-local guesses:
	 -javax/inject/category.
	 -javax/inject/category#
	 -javax/inject/category().
	 -play/api/mvc/category.
	 -play/api/mvc/category#
	 -play/api/mvc/category().
	 -play/api/libs/json/category.
	 -play/api/libs/json/category#
	 -play/api/libs/json/category().
	 -category.
	 -category#
	 -category().
	 -scala/Predef.category.
	 -scala/Predef.category#
	 -scala/Predef.category().
offset: 1029
uri: file://<WORKSPACE>/app/controllers/CategoryController.scala
text:
```scala
package controllers

import javax.inject._
import play.api.mvc._
import play.api.libs.json._
import models.Category
import scala.collection.mutable.ListBuffer

@Singleton
class CategoryController @Inject()(val controllerComponents: ControllerComponents) extends BaseController {

  private val categoryList = new ListBuffer[Category]()
  categoryList += Category(1, "Elektronika")
  categoryList += Category(2, "Akcesoria")

  def getAll: Action[AnyContent] = Action { implicit request =>
    Ok(Json.toJson(categoryList))
  }

  def getById(id: Long): Action[AnyContent] = Action { implicit request =>
    categoryList.find(_.id == id) match {
      case Some(category) => Ok(Json.toJson(category))
      case None => NotFound(Json.obj("message" -> "Nie znaleziono kategorii"))
    }
  }

  def add(): Action[JsValue] = Action(parse.json) { implicit request =>
    request.body.validate[Category].fold(
      errors => BadRequest(Json.obj("message" -> "Błędny format JSON")),
      category => {
        categoryList += category@@
        Created(Json.toJson(category))
      }
    )
  }

  def update(id: Long): Action[JsValue] = Action(parse.json) { implicit request =>
    request.body.validate[Category].fold(
      errors => BadRequest(Json.obj("message" -> "Błędny format JSON")),
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

  def delete(id: Long): Action[AnyContent] = Action { implicit request =>
    val index = categoryList.indexWhere(_.id == id)
    if (index >= 0) {
      categoryList.remove(index)
      Ok(Json.obj("message" -> "Usunięto kategorię"))
    } else {
      NotFound(Json.obj("message" -> "Nie znaleziono kategorii"))
    }
  }
}
```


#### Short summary: 

empty definition using pc, found symbol in pc: 