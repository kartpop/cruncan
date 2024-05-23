Feature: Onerequest kafka handler

  Scenario Outline: Should process onerequest kafka message and post to threerequest api
    Given the threerequest API responds with <status_code> on matching body from file <threerequest_body>
    And the token API response with <status_code> on matching body from file <token_body>
    When a onerequest from file <onerequest_body> is ingested
    Then a request is made to the threerequest API
    
    Examples:
      | status_code | threerequest_body         | token_body                | onerequest_body       |
      | 200         | "data/3r_body.json"       | "data/token_body.json"    | "data/1r_body.json"   |

