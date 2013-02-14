/*
httpize export methods of a given variable to handle HTTP requests.

It allows you to create a http.Handler tied to a variable whose type 
implements httpize.MethodProvider. Any HTTP request routed to the handler,
will be checked to see if the URL matchs the name of a method of the 
MethodProvider which was exported. If so each parameter to the method 
will be created from the URL parameters, checked and passed to the method.
On returning the exported method returns a Reader used to create the HTTP 
response body.

*/
package httpize
