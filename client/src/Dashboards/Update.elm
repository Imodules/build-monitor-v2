module Dashboards.Update exposing (..)

import Dashboards.Api as Api
import Dashboards.Lib exposing (configInList, getBuildPath, getDefaultPrefix)
import Dashboards.Models as Dashboards exposing (Branch, BuildConfigForm, ConfigDetail, DashboardDetails, EditTab(Configure, Select), VisibleBranch, buildConfigToForm, initialBuildConfigForm, initialFormModel)
import Lib exposing (createCommand)
import List.Extra exposing (find, findIndex, getAt, swapAt)
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
            let
                newVisibleBranches =
                    case model.details of
                        RemoteData.Success details ->
                            switchBranches details.configs model.visibleBranches

                        _ ->
                            model.visibleBranches
            in
            ( { model | visibleBranches = newVisibleBranches }, Cmd.none )

        ChangeDashboardName value ->
            let
                newDashboardForm old =
                    { old | name = updateDashboardName value, isDirty = True }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        ChangeDashboardColumnCount value ->
            let
                newDashboardForm old =
                    { old | columnCount = updateColumnCount value, isDirty = True }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        ChangeSuccessIcon value ->
            let
                newDashboardForm old =
                    { old | successIcon = updateIcon value, isDirty = True }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        ChangeFailedIcon value ->
            let
                newDashboardForm old =
                    { old | failedIcon = updateIcon value, isDirty = True }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        ChangeRunningIcon value ->
            let
                newDashboardForm old =
                    { old | runningIcon = updateIcon value, isDirty = True }
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
                                , columnCount = initTextFieldValue (toString dashboard.columnCount)
                                , successIcon = initTextFieldValue dashboard.successIcon
                                , failedIcon = initTextFieldValue dashboard.failedIcon
                                , runningIcon = initTextFieldValue dashboard.runningIcon
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

        SortConfigUp i ->
            let
                newDashboardForm old =
                    { old | buildConfigs = sortConfig i -1 old.buildConfigs, isDirty = True }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        SortConfigDown i ->
            let
                newDashboardForm old =
                    { old | buildConfigs = sortConfig i 1 old.buildConfigs, isDirty = True }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )


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


updateColumnCount : String -> TextField
updateColumnCount value =
    let
        intValue =
            Result.withDefault 0 (String.toInt value)

        ( isValid, error ) =
            if intValue < 1 || intValue > 12 then
                ( False, "Column count should be a value >= 1 and <= 12" )
            else
                ( True, "" )
    in
    { value = value
    , isValid = isValid
    , isDirty = True
    , error = error
    }


updateIcon : String -> TextField
updateIcon value =
    let
        ( isValid, error ) =
            if String.length value == 0 then
                ( False, "A value is required" )
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


sortConfig : Int -> Int -> List BuildConfigForm -> List BuildConfigForm
sortConfig idx direction configs =
    let
        newIdx =
            idx + direction

        maybeList =
            swapAt idx newIdx configs
    in
    case maybeList of
        Just list ->
            list

        _ ->
            configs


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


switchBranches : List ConfigDetail -> List VisibleBranch -> List VisibleBranch
switchBranches configs branches =
    List.map
        (\b ->
            if b.size == 1 then
                b
            else
                incrementBranch configs b
        )
        branches


incrementBranch : List ConfigDetail -> VisibleBranch -> VisibleBranch
incrementBranch configs vb =
    let
        cfg =
            find (\x -> x.id == vb.id) configs

        isRunning =
            False

        runningBranchIndex =
            case cfg of
                Just config ->
                    if config.isRunning then
                        getRunningBranchIndex config vb.index
                    else
                        Nothing

                _ ->
                    Nothing
    in
    case runningBranchIndex of
        Just index ->
            { vb | index = index }

        _ ->
            if vb.index < vb.size - 1 then
                { vb | index = vb.index + 1 }
            else
                { vb | index = 0 }


getRunningBranchIndex : ConfigDetail -> Int -> Maybe Int
getRunningBranchIndex config currentIndex =
    let
        maybeCurrentBranch =
            getAt currentIndex config.branches

        runningBranches =
            List.filter (\x -> x.isRunning) config.branches
    in
    case maybeCurrentBranch of
        Just branch ->
            findIndex (\x -> x.isRunning && x.name /= branch.name) config.branches

        _ ->
            findIndex (\x -> x.isRunning) config.branches
