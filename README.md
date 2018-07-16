# Addressbook

Endpoints:
POST   /users             Send user object(first_name, last_name, email, phone_number) in body with datatype application/json

GET    /users             Returns all users

GET    /users/:id         Returns user by Id

PUT    /users/:id         Update user by Id. Send updated user object in body

DELETE /users/:id         Delete user by id

GET    /export            Exports all users in db and returns the path

POST   /import            Given a Path in query params users will be imported in to db.

How to Setup:
After mysql is installed in the machine. Create database 'restapp' and table with the following schema

CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `first_name` varchar(50) DEFAULT NULL,
  `last_name` varchar(50) DEFAULT NULL,
  `email` varchar(50) DEFAULT NULL,
  `phone_number` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=latin1;

Username and password is root/root
