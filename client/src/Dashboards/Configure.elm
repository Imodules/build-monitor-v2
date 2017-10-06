module Dashboards.Configure exposing (..)

import Dashboards.Lib exposing (getBuildPath)
import Dashboards.Models as DashboardsModel exposing (BuildConfig, BuildConfigForm)
import Html exposing (Html, button, div, h6, hr, i, p, span, text)
import Html.Attributes exposing (class, disabled, id)
import Html.Events exposing (onClick)
import Models exposing (BuildType, Model, Project, initialProject)
import Msgs exposing (DashboardMsg(ChangeBuildAbbreviation, ChangeDashboardName, EditDashboard, SortConfigDown, SortConfigUp), Msg(DashboardMsg))
import Pages.Components exposing (textBox, textField)


isFormValid : DashboardsModel.Model -> Bool
isFormValid model =
    let
        isValid =
            not (List.any (\c -> not c.abbreviation.isValid) model.dashboardForm.buildConfigs)
    in
    model.dashboardForm.name.isValid && model.dashboardForm.columnCount.isValid && isValid


view : Model -> List Project -> List BuildType -> Html Msg
view model projects buildTypes =
    let
        listLength =
            List.length model.dashboards.dashboardForm.buildConfigs

        indexedConfigRow i bc =
            configRow i listLength bc projects buildTypes

        content =
            List.indexedMap indexedConfigRow model.dashboards.dashboardForm.buildConfigs
    in
    div [] content


configRow : Int -> Int -> BuildConfigForm -> List Project -> List BuildType -> Html Msg
configRow index total config projects buildTypes =
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
            , div [ class "level-right" ] [ div [ class "level-item" ] [ sortButtons index total ] ]
            ]
        ]


sortButtons : Int -> Int -> Html Msg
sortButtons index total =
    div [ class "field has-addons" ]
        [ p [ class "control" ] [ iconButton "fa fa-angle-up" "Sort Up" (index == 0) (SortConfigUp index) ]
        , p [ class "control" ] [ iconButton "fa fa-angle-down" "Sort Down" (index == total - 1) (SortConfigDown index) ]
        ]


iconButton : String -> String -> Bool -> DashboardMsg -> Html Msg
iconButton icon_ text_ disabled_ msg =
    button [ class "button", disabled disabled_, onClick (DashboardMsg msg) ]
        [ span [ class "icon is-small" ] [ i [ class icon_ ] [] ]
        , span [] [ text text_ ]
        ]
