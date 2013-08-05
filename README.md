# Go Ember Token-based Authentication

### Description ###

Simple authentication server to use with [embercasts authentication screencasts](http://www.embercasts.com/).



### Todo: ###

In func ValidTokenProvided:
The server stores the current token in the currentToken variable, this part is working. The current problem is, it is unknown in which area of the request the user token resides. Application works fine on it's own. However, when concept is applied to a real application, missing component will be needed.