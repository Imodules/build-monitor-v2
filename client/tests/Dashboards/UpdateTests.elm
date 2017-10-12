module Dashboards.UpdateTests exposing (..)

import Dashboards.Update exposing (..)
import Expect exposing (Expectation)
import Fuzz exposing (Fuzzer, int, list, string)
import Test exposing (..)


updateDashboardNameSuite : Test
updateDashboardNameSuite =
    describe "updateDashboardName"
        [ describe "When we have a name that is >= 5 characters"
            [ test "It should return a valid result" <|
                \_ ->
                    updateDashboardName "abcde"
                        |> Expect.all
                            [ (.value >> Expect.equal "abcde")
                            , (.isValid >> Expect.true "isValid should be true")
                            , (.isDirty >> Expect.true "isDirty should be true")
                            , (.error >> Expect.equal "")
                            ]
            ]
        , describe "When we have < 5 characters in our name"
            [ test "It should return an error with the proper message" <|
                \_ ->
                    updateDashboardName "abcd"
                        |> Expect.all
                            [ (.value >> Expect.equal "abcd")
                            , (.isValid >> Expect.false "isValid should be false")
                            , (.isDirty >> Expect.true "isDirty should be true")
                            , (.error >> Expect.equal "Name must be at least 5 characters")
                            ]
            ]
        ]


getRunningBranchIndexSuite : Test
getRunningBranchIndexSuite =
    describe "getRunningBranchIndex"
        [ describe "When we do not have any running branches"
            [ test "It should return Nothing" <|
                let
                    branches =
                        [ { name = "branch 1", isRunning = False, builds = [] }
                        , { name = "branch 2", isRunning = False, builds = [] }
                        , { name = "branch 3", isRunning = False, builds = [] }
                        ]

                    configDetail =
                        { id = "cfg-id", name = "then cd name", abbreviation = "ab1", isRunning = False, branches = branches }

                    result =
                        getRunningBranchIndex configDetail 0
                in
                    \_ -> Expect.equal Nothing result
            ]
        , describe "When we have a running branch"
            [ test "It should return the index of that branch" <|
                let
                    branches =
                        [ { name = "branch 1", isRunning = False, builds = [] }
                        , { name = "branch 2", isRunning = True, builds = [] }
                        , { name = "branch 3", isRunning = False, builds = [] }
                        ]

                    configDetail =
                        { id = "cfg-id", name = "then cd name", abbreviation = "ab1", isRunning = True, branches = branches }

                    result =
                        getRunningBranchIndex configDetail 0
                in
                    \_ -> Expect.equal (Just 1) result
            ]
        , describe "When we have multiple running branches and the first index is the current"
            [ test "It should return the index of the next running branch" <|
                let
                    branches =
                        [ { name = "branch 1", isRunning = True, builds = [] }
                        , { name = "branch 2", isRunning = False, builds = [] }
                        , { name = "branch 3", isRunning = True, builds = [] }
                        ]

                    configDetail =
                        { id = "cfg-id", name = "then cd name", abbreviation = "ab1", isRunning = True, branches = branches }

                    result =
                        getRunningBranchIndex configDetail 0
                in
                    \_ -> Expect.equal (Just 2) result
            ]
        , describe "When we have multiple running branches and the last index is the current"
            [ test "It should return the index of the first running branch" <|
                let
                    branches =
                        [ { name = "branch 1", isRunning = True, builds = [] }
                        , { name = "branch 2", isRunning = False, builds = [] }
                        , { name = "branch 3", isRunning = True, builds = [] }
                        ]

                    configDetail =
                        { id = "cfg-id", name = "then cd name", abbreviation = "ab1", isRunning = True, branches = branches }

                    result =
                        getRunningBranchIndex configDetail 2
                in
                    \_ -> Expect.equal 0 (Maybe.withDefault -1 result)
            ]
        , describe "When we have multiple running branches and the middle index is the current"
            [ test "It should return the index of the next running branch" <|
                let
                    branches =
                        [ { name = "branch 1", isRunning = True, builds = [] }
                        , { name = "branch 2", isRunning = True, builds = [] }
                        , { name = "branch 3", isRunning = True, builds = [] }
                        , { name = "branch 4", isRunning = True, builds = [] }
                        ]

                    configDetail =
                        { id = "cfg-id", name = "then cd name", abbreviation = "ab1", isRunning = True, branches = branches }

                    result =
                        getRunningBranchIndex configDetail 2
                in
                    \_ -> Expect.equal (Just 3) result
            ]
        ]
