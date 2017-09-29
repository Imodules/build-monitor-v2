module Dashboards.Update exposing (..)

import Dashboards.Api as Api
import Dashboards.Lib exposing (configInList, getBuildPath, getDefaultPrefix)
import Dashboards.Models as Dashboards exposing (BuildConfigForm, ConfigDetail, EditTab(Configure, Select), VisibleBranch, buildConfigToForm, initialBuildConfigForm, initialFormModel)
import Lib exposing (createCommand)
import List.Extra exposing (find)
import Models exposing (Model)
import Msgs exposing (DashboardMsg(..), Msg(ChangeLocation, DashboardMsg))
import RemoteData
import Routes exposing (Route(DashboardsRoute))
import Routing exposing (getToken)
import Types exposing (Id, TextField, Token, initTextFieldValue)


update : DashboardMsg -> Model -> ( Model, Cmd Msg )
update msg model =
    let
        ( newDashboards, cmd ) =
            update_ model.flags.apiUrl (getToken model) msg model
    in
    ( { model | dashboards = newDashboards }, cmd )


update_ : String -> Token -> DashboardMsg -> Model -> ( Dashboards.Model, Cmd Msg )
update_ baseUrl token msg model_ =
    let
        model =
            model_.dashboards
    in
    case msg of
        OnFetchDashboards response ->
            ( { model | dashboards = response }, Cmd.none )

        OnFetchDetails response ->
            let
                newVisibleBranches =
                    if model.visibleBranches == [] then
                        initVisibleBranches model
                    else
                        model.visibleBranches
            in
            ( { model | details = response, visibleBranches = newVisibleBranches }, Cmd.none )

        ChangeBranches _ ->
            ( { model | visibleBranches = switchBranches model.visibleBranches }, Cmd.none )

        ChangeDashboardName value ->
            let
                newDashboardForm old =
                    { old | name = updateDashboardName value, isDirty = True }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        ChangeBuildAbbreviation id value ->
            let
                newDashboardForm old =
                    { old | buildConfigs = updateBuildConfigAbbreviation id value old.buildConfigs, isDirty = True }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        ClickBuildType id ->
            let
                updatedList old =
                    if configInList id old then
                        List.filter (\i -> i.id /= id) old
                    else
                        getNewConfig model_ id :: old

                newDashboardForm old =
                    { old | buildConfigs = updatedList old.buildConfigs, isDirty = True }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        CreateDashboard ->
            ( model, Api.createDashboard baseUrl token model )

        EditDashboard ->
            ( model, Api.editDashboard baseUrl token model )

        OnSelectTabClick ->
            let
                newDashboardForm old =
                    { old | tab = Select }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        OnConfigureTabClick ->
            let
                newDashboardForm old =
                    { old | tab = Configure }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        StartCreateDashboard ->
            ( { model | dashboardForm = initialFormModel }, Cmd.none )

        StartEditDashboard id ->
            let
                dashboards =
                    case model.dashboards of
                        RemoteData.Success dashboards ->
                            dashboards

                        _ ->
                            []

                maybeDashboardToEdit =
                    find (\i -> i.id == id) dashboards

                newDashboardForm old =
                    case maybeDashboardToEdit of
                        Just dashboard ->
                            if String.isEmpty old.id || old.id /= id then
                                { id = dashboard.id
                                , name = initTextFieldValue dashboard.name
                                , buildConfigs = List.map buildConfigToForm dashboard.buildConfigs
                                , isDirty = False
                                , tab = Select
                                }
                            else
                                old

                        _ ->
                            Dashboards.initialFormModel
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        OnCreateDashboard result ->
            case result of
                Ok dashboard ->
                    ( model, createCommand (ChangeLocation DashboardsRoute) )

                Err dashboard ->
                    let
                        x =
                            Debug.log "error saving dashboard" dashboard
                    in
                    ( model, Cmd.none )


updateDashboardName : String -> TextField
updateDashboardName value =
    let
        ( isValid, error ) =
            if String.length value < 5 then
                ( False, "Name must be at least 5 characters" )
            else
                ( True, "" )
    in
    { value = value
    , isValid = isValid
    , isDirty = True
    , error = error
    }


updateBuildConfigAbbreviation : Id -> String -> List BuildConfigForm -> List BuildConfigForm
updateBuildConfigAbbreviation id value configs =
    let
        findAndUpdate config =
            if config.id == id then
                { config | abbreviation = updateAbbreviation value }
            else
                config
    in
    List.map findAndUpdate configs


updateAbbreviation : String -> TextField
updateAbbreviation value =
    let
        ( isValid, error ) =
            if String.length value < 1 then
                ( False, "Abbreviation must be at least 1 character" )
            else
                ( True, "" )
    in
    { value = value
    , isValid = isValid
    , isDirty = True
    , error = error
    }


getNewConfig : Model -> Id -> BuildConfigForm
getNewConfig model id =
    let
        projects =
            case model.projects of
                RemoteData.Success projects ->
                    projects

                _ ->
                    []

        buildTypes =
            case model.buildTypes of
                RemoteData.Success buildTypes ->
                    buildTypes

                _ ->
                    []

        buildPath =
            getBuildPath id projects buildTypes

        prefix =
            getDefaultPrefix buildPath
    in
    initialBuildConfigForm id prefix


initVisibleBranches : Dashboards.Model -> List VisibleBranch
initVisibleBranches model =
    case model.details of
        RemoteData.Success details ->
            buildNewVisibleBranches details.configs

        _ ->
            model.visibleBranches


buildNewVisibleBranches : List ConfigDetail -> List VisibleBranch
buildNewVisibleBranches cds =
    List.map createVisibleBranch cds


createVisibleBranch : ConfigDetail -> VisibleBranch
createVisibleBranch cd =
    { id = cd.id, size = List.length cd.branches, index = 0 }


updateVisibleBranches : List ConfigDetail -> List VisibleBranch -> List VisibleBranch
updateVisibleBranches cds old =
    let
        updateBranch cd =
            createVisibleBranch cd

        getUpdatedBranch cd maybeBranch =
            case maybeBranch of
                Just branch ->
                    let
                        branchCount =
                            List.length cd.branches

                        validIndex =
                            if branch.index >= branchCount then
                                0
                            else
                                branch.index
                    in
                    { branch | size = branchCount, index = validIndex }

                _ ->
                    createVisibleBranch cd
    in
    List.map updateBranch cds


switchBranches : List VisibleBranch -> List VisibleBranch
switchBranches branches =
    List.map
        (\b ->
            if b.size == 1 then
                b
            else
                incrementBranch b
        )
        branches


incrementBranch : VisibleBranch -> VisibleBranch
incrementBranch branch =
    if branch.index < branch.size - 1 then
        { branch | index = branch.index + 1 }
    else
        { branch | index = 0 }
