/*
Httpize is a web framework, that allows you use methods/types in web requests.

Create a http.Handler tied to a value. Any HTTP request 
routed to the handler, will be checked to see if the URL matchs the name of a 
method of the which was exported. If so each parameter to the method 
will be created from the URL parameters, checked and passed to the method.
On returning the exported method returns a io.WriterTo used to create the HTTP 
response body.
*/
package httpize
