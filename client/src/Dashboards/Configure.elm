module Dashboards.Configure exposing (..)

import Dashboards.Lib exposing (getBuildPath)
import Dashboards.Models as DashboardsModel exposing (BuildConfig, BuildConfigForm)
import Html exposing (Html, div, h6, hr, text)
import Html.Attributes exposing (class, id)
import Models exposing (BuildType, Model, Project, initialProject)
import Msgs exposing (DashboardMsg(ChangeBuildAbbreviation, ChangeDashboardName, EditDashboard), Msg(DashboardMsg))
import Pages.Components exposing (textBox, textField)


isFormValid : DashboardsModel.Model -> Bool
isFormValid model =
    let
        isValid =
            not (List.any (\c -> not c.abbreviation.isValid) model.dashboardForm.buildConfigs)
    in
    model.dashboardForm.name.isValid && isValid


view : Model -> List Project -> List BuildType -> Html Msg
view model projects buildTypes =
    let
        content =
            List.map (\i -> configRow i projects buildTypes) model.dashboards.dashboardForm.buildConfigs
    in
    div [] content


configRow : BuildConfigForm -> List Project -> List BuildType -> Html Msg
configRow config projects buildTypes =
    let
        buildPath =
            getBuildPath config.id projects buildTypes
    in
    div [ class "box" ]
        [ text buildPath
        , div [ class "level" ]
            [ div [ class "level-left" ]
                [ textBox config.abbreviation (ChangeBuildAbbreviation config.id >> Msgs.DashboardMsg)
                ]
            ]
        ]
