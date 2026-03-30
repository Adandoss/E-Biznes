error id: file://<WORKSPACE>/app/models/Category.scala:Json.
file://<WORKSPACE>/app/models/Category.scala
empty definition using pc, found symbol in pc: 
empty definition using semanticdb
empty definition using fallback
non-local guesses:
	 -play/api/libs/json/Json.
	 -Json.
	 -scala/Predef.Json.
offset: 153
uri: file://<WORKSPACE>/app/models/Category.scala
text:
```scala
package models

import play.api.libs.json._

case class Category(id: Long, name: String)

object Category {
  implicit val format: OFormat[Category] = Js@@on.format[Category]
}
```


#### Short summary: 

empty definition using pc, found symbol in pc: 