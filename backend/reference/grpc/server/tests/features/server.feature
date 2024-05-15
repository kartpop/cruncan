Feature: Test grpc server

  Scenario: Test grpc server
    When a transaction request with <amount> and <interest> is sent to the server for transfer from <source_account> to <target_account>
    Then the server processes and sends a transaction response with <success> and <transferred> amount
    Examples:
      | amount | interest | source_account | target_account | success | transferred |
      | 100    | 15       | "sdf-12"       | "sdf-18"       | "true"  | 115         |
      | 100    | 0        | "sdf-15"       | "sdf-21"       | "true"  | 100         |
