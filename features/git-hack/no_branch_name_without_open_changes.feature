Feature: git hack errors without a branch name without open changes

  Background:
    Given I have a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    |
      | main    | remote   | main commit    | main_file    |
      | feature | local    | feature commit | feature_file |
    And I am on the "feature" branch
    When I run `git hack` while allowing errors


  Scenario: result
    Then it runs no Git commands
    And I get the error "No branch name provided"
    And I am still on the "feature" branch
