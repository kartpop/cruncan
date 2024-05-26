Feature: Onerequest http handler

  Scenario Outline: Should handle POST onerequest, publish to kafka and save to database
    When a POST request with JSON body from file <onerequest> is sent
    Then the request with message from file <onerequest> is published to kafka
    And the request with message from file <onerequest> is saved to database with correct user id
    And the response status code is <status>

    Examples:
      | onerequest              | status |
      | "data/1r_body.json"     | 201    |
