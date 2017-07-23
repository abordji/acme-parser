# ACME Application Test

This test is meant to both give the applicant a good overview of typical
projects he will have to work on, and to give ACME a good overview of
the applicant skills and "way of doing things". The applicant is free
to choose whatever language he wants (Golang being the primary language
used at ACME and being well-suited for the task, using it might be a
good point)

The core goals can be done in less than 180 minutes, stretch goals are
optional but appreciated. This test is guideline-less on purpose to let
the applicant expose all skills he wants. The submission must consist of
a compilable source tree with a README file explaining how to compile,
configure and run it from scratch

## Description of the Task:

* Goal: develop an application to save contact's information loaded from
  an API into a database (MySQL or SQLite). The application has to:
    * handle the stream format:
         * one line per contact
         * each channel of this contact begins with a prefix that defines
           the value (@ = email, cookie = cookie, uid = uid)
    * handle the fact that one contact can have several emails, uids
      (user identification key) and cookies.

* In this test, the internal software can:
    * Call an API
    * Load and read the API response
    * Test if the response respect criteria
    * If it does, save all information in a database
    * Deals with potential errors

* API call:
    * The application will call this url:
      http://acme.co/contacts

## Rating

The rating of the submission will be based on multiple criteria, such as:

* simplicity: how simple is the code
* readability: is the code understandable, readable, and specific
* flexibility: is the application easily reconfigurable
* rapidity: how much time does the application take
* prediction: is it possible to predict this time
* optimization: is the application usable for Big Data
