# Swaggerilla - write tests and generate docs!

Simple test framework for API testing. Should be used against running dev instance of the API.

It allows to write tests in declarative way, then declaration can be used to generate swagger documentation.

It's still pre-pre-pre-alfa, needs to be improved and fixed a lot. Currently, it's just a proof of concept.

## Advantages of such framework:
- API tests and Swagger documentation with examples from the same box.
- Swagger doc can be tested against test instance of the API. If your tests pass, then your doc will work too.
- No need to refactor your API. You can write code of the API as you wish.

## Drawbacks:
- One may forget to define some API parameter in test, so it won't be available in the documentation. And there is no way (other than code review) to verify that.
- Swagger supports one declaration of request for one HTTP return code (1 declaration for code 200, one for 404 and so on). But what if you have different test cases and all of them produce the same 200 response code? Currently, only first test is used in such situation.
- It's difficult to define all properties of the swagger (like validators, formats) and make the code of the tests readable at the same time. Currently many things provided by swagger are ignored for sake of simplicity of the tests