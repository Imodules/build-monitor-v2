module Dashboards.Lib exposing (..)

import Dashboards.Models exposing (BuildConfigForm)
import List.Extra exposing (find)
import Models exposing (BuildType, Project, initialProject)
import Types exposing (Id)


configInList : Id -> List BuildConfigForm -> Bool
configInList id configs =
    List.any (\config -> config.id == id) configs


getBuildPath : Id -> List Project -> List BuildType -> String
getBuildPath id projects buildTypes =
    let
        maybeBuildType =
            find (\i -> i.id == id) buildTypes
    in
    case maybeBuildType of
        Just buildType ->
            getProjectPath buildType.projectId projects ++ " / " ++ buildType.name

        _ ->
            ""


getProjectPath : Id -> List Project -> String
getProjectPath id projects =
    let
        maybeParentProject =
            find (\i -> i.id == id) projects

        parentProject =
            case maybeParentProject of
                Just project ->
                    project

                _ ->
                    initialProject
    in
    if parentProject.parentProjectId /= "_Root" then
        parentProject.name ++ " / " ++ getProjectPath parentProject.parentProjectId projects
    else
        parentProject.name


getDefaultPrefix : String -> String
getDefaultPrefix path =
    let
        paths =
            String.split " / " path

        parts =
            List.map (\p -> getPathPart p ++ "-") paths
    in
    String.dropRight 1 (String.concat parts)


getPathPart : String -> String
getPathPart s =
    let
        words =
            String.words s

        letters =
            List.map getFirstLetter words
    in
    String.concat letters


getFirstLetter : String -> String
getFirstLetter s =
    String.slice 0 1 s